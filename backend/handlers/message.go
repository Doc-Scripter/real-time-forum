package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"forum/database"
	"forum/logging"
	"forum/middleware"
	"forum/utils"
	"time"

	"github.com/gorilla/websocket"
)

// Message represents a chat message with sender, receiver and content information
type Message struct {
	Data     string `json:"data"`
	Sender   string `json:"sender"`
	Receiver int    `json:"receiver"`
	Time     string `json:"time"`
}

var (
	messages      []Message
	messagesMutex sync.Mutex
	upgrader      = websocket.Upgrader{
		HandshakeTimeout: 0,
		ReadBufferSize:   0,
		WriteBufferSize:  0,
		WriteBufferPool:  nil,
		Subprotocols:     []string{},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			panic("TODO")
		},
		CheckOrigin: func(r *http.Request) bool {
			panic("TODO")
		},
		EnableCompression: false,
	}
	clients       = make(map[*websocket.Conn]bool)
	broadcast     = make(chan Message)
)

// init starts the message handling goroutine
func init() {
	go handleMessages()
}

// handleMessages processes incoming messages and broadcasts them to connected clients
func handleMessages() {
	for {
		msg, ok := <-broadcast
		if !ok {
			fmt.Println("broadcast channel closed")
			return
		}
		messagesMutex.Lock()
		messages = append(messages, msg)
		messagesMutex.Unlock()
		for client := range clients {
			envelope := map[string]interface{}{
				"type":    "message",
				"message": msg,
			}
			err := client.WriteJSON(envelope)
			if err != nil {
				logging.Log("[ERROR] Error sending message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// MarkMessagesAsRead updates messages as read for a given receiver and sender
func MarkMessagesAsRead(receiverID, senderID int) error {
    query := `
        UPDATE messages 
        SET is_read = 1 
        WHERE receiver_id = ? AND user_id = ?
    `
    _, err := database.DB.Exec(query, receiverID, senderID)
    if err != nil {
		logging.Log("Error marking messages as read: %v", err)
        return err
    }
    return nil
}

// SaveMessageToDB stores a new message in the database
func SaveMessageToDB(senderID, receiverID int, text string) error {
	_, err := database.DB.Exec(
		"INSERT INTO messages (user_id, receiver_id, message,is_read) VALUES (?, ?, ?, 0)",
		senderID, receiverID, text,
	)
	return err
}

// getMessages retrieves all messages between two users, ordered by creation time
func getMessages(loggedInUserID, otherUserID int) ([]Message, error) {
	rows, err := database.DB.Query(
		`SELECT m.message, u.nickname AS  sender, m.receiver_id, m.created_at
		 FROM messages m
		 JOIN users u ON m.user_id = u.id
		 WHERE (m.user_id = ? AND m.receiver_id = ?) OR (m.user_id = ? AND m.receiver_id = ?)
		 ORDER BY m.created_at ASC`,
		loggedInUserID, otherUserID,
		otherUserID, loggedInUserID,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying messages: %v", err)
	}
	defer rows.Close()

	var allMessages []Message
	for rows.Next() {
		var msg Message
		var createdAt time.Time

		if err := rows.Scan(&msg.Data, &msg.Sender, &msg.Receiver, &createdAt); err != nil {
			return nil, fmt.Errorf("error scanning message row: %v", err)
		}

		msg.Time = createdAt.Format("3:04:05 PM")
		allMessages = append(allMessages, msg)
	}

	return allMessages, nil
}

// getLastMessages retrieves the most recent message from each conversation for a user
func getLastMessages(userID int) ([]Message, error) {
    query := `
        SELECT m.message, u.nickname as sender_name, 
               CASE 
                   WHEN m.user_id = ? THEN m.receiver_id 
                   ELSE m.user_id 
               END as other_user_id,
               m.created_at
        FROM messages m
        JOIN users u ON m.user_id = u.id
        JOIN (
            SELECT 
                CASE 
                    WHEN user_id = ? THEN receiver_id 
                    ELSE user_id 
                END as conversation_partner,
                MAX(created_at) as latest_time
            FROM messages
            WHERE user_id = ? OR receiver_id = ?
            GROUP BY conversation_partner
        ) latest ON (
            (m.user_id = ? AND m.receiver_id = latest.conversation_partner) OR
            (m.user_id = latest.conversation_partner AND m.receiver_id = ?)
        ) AND m.created_at = latest.latest_time
        ORDER BY m.created_at DESC
    `
    
    rows, err := database.DB.Query(query, userID, userID, userID, userID, userID, userID)
    if err != nil {
		fmt.Println("error querying last messages: ",err)
        return nil, fmt.Errorf("error querying last messages: %v", err)
    }
    defer rows.Close()
    
    var result []Message
    for rows.Next() {
        var msg Message
        var createdAt time.Time
        var otherUserID int
        
        if err := rows.Scan(&msg.Data, &msg.Sender, &otherUserID, &createdAt); err != nil {
            return nil, fmt.Errorf("error scanning message row: %v", err)
        }
        
        msg.Time = createdAt.Format("3:04:05 PM")
        msg.Receiver = otherUserID
        result = append(result, msg)
    }
    
    return result, nil
}

// MessageHandler handles HTTP requests for message operations (GET, POST, PUT)
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.ErrorMessage(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodPut:
        var req struct {
            SenderID int `json:"senderId"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logging.Log("Error decoding request body: %v", err)
            utils.ErrorMessage(w, "Invalid request", http.StatusBadRequest)
            return
        }
        
        userID, ok = middleware.GetUserID(r)
        if !ok {
            utils.ErrorMessage(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if err := MarkMessagesAsRead(userID, req.SenderID); err != nil {
			logging.Log("Error marking messages as read: %v", err)
            utils.ErrorMessage(w, "Failed to mark messages as read", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        return
	case http.MethodGet:
		if r.URL.Query().Has("receiver"){
		receiverIDStr := r.URL.Query().Get("receiver")

		if receiverIDStr == "" {
			utils.ErrorMessage(w, "Invalid receiver ID", http.StatusBadRequest)
			return
		}
		receiverID, err := strconv.Atoi(receiverIDStr)
		if err != nil {
			logging.Log("Error parsing receiver ID: %v", err)
			utils.ErrorMessage(w, "error parsing the receiverid", http.StatusBadRequest)
			return
		}

		msgs, err := getMessages(userID, receiverID)
		if err != nil {
			logging.Log("Error getting messages: %v", err)
			utils.ErrorMessage(w, "Failed to fetch messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msgs)
		return
	}else{
		msgs, err := getLastMessages(userID)
		if err != nil {
			logging.Log("Error getting messages: %v", err)
			utils.ErrorMessage(w, "Failed to fetch messages", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msgs)
		return
	}

	case http.MethodPost:
		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			logging.Log("Error decoding message: %v", err)
			utils.ErrorMessage(w, "Invalid Request", http.StatusBadRequest)
			return
		}
		senderID, ok := middleware.GetUserID(r)
		if !ok {
			utils.ErrorMessage(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err := SaveMessageToDB(senderID, msg.Receiver, msg.Data); err != nil {
			logging.Log("Error saving message to DB: %v", err)
			utils.ErrorMessage(w, "Failed to save message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	default:
		utils.ErrorMessage(w, "Method not allowed", http.StatusBadRequest)
		return
	}
}
