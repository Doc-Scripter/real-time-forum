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

type Message struct {
	// ID       int    `json:"id,omitempty"`
	Data     string `json:"data"`
	Sender   string `json:"sender"`
	Receiver int    `json:"receiver"`
	Time     string `json:"time"`
}

var (
	messages      []Message
	messagesMutex sync.Mutex
	upgrader      = websocket.Upgrader{}
	clients       = make(map[*websocket.Conn]bool)
	broadcast     = make(chan Message)
)

func init() {
	go handleMessages()
}

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
			// Send an envelope with type and message
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

func SaveMessageToDB(senderID, receiverID int, text string) error {
	_, err := database.DB.Exec(
		"INSERT INTO messages (user_id, receiver_id, message,is_read) VALUES (?, ?, ?, 0)",
		senderID, receiverID, text,
	)
	return err
}

// getMessages retrieves all messages between two users, ordered by creation time.
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

		msg.Time = createdAt.Format("15:04:05")
		allMessages = append(allMessages, msg)
	}

	return allMessages, nil
}

// New function to get all messages for a user
func getLastMessages(userID int) ([]Message, error) {
     // This query gets the last message from each conversation
    // It uses a subquery to find the max created_at for each conversation pair
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
        
        msg.Time = createdAt.Format("15:04:05")
        msg.Receiver = otherUserID // Set the other user's ID as the receiver
        result = append(result, msg)
    }
    
    return result, nil
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.ErrorMessage(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodPut:
        // Mark messages as read
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

		// Check if we're filtering by receiver
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

		// Get messages between these two users
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

func MessageWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Log("Error upgrading to WebSocket: %v", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	clients[conn] = true
	defer func() {
		conn.Close()
		delete(clients, conn)
	}()
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			logging.Log("Error reading message: %v", err)
			break
		}
		if msg.Time == "" {
			msg.Time = time.Now().Format("15:04:05")
		}

		broadcast <- msg
	}
}
