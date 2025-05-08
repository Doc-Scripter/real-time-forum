package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"forum/database"
	"forum/middleware"
	"forum/utils"

	"github.com/gorilla/websocket"
)

type Message struct {
	// ID       int    `json:"id,omitempty"`
	Data     string `json:"data"`
	Sender   string `json:"sender"`
	Receiver int `json:"receiver"`
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
				client.Close()
				delete(clients, client)
			}
		}
	}
}
/*
func getMessages(user_id string, receiver string) {
	args := []any{user_id, receiver}

	query := `
			SELECT message FROM messages
			WHERE user_id = ? AND receiver_id = ?
			ORDER BY created_at DESC`
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		fmt.Errorf("error querying posts: %v", err)
		return
	}
	defer rows.Close()
	// var Messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Text); err != nil {
			fmt.Errorf("error scanning rows: %v", err)
			return
		}
		msg.Sender = user_id
		msg.Receiver = receiver
		messages = append(messages, msg)
	}
}
*/
func SaveMessageToDB(senderID, receiverID int, text string) error {
    _, err := database.DB.Exec(
        "INSERT INTO messages (sender_id, receiver_id, text) VALUES (?, ?, ?)",
        senderID, receiverID, text,
    )
    return err
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(messages)
fmt.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		messagesMutex.Lock()
		defer messagesMutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		// Ensure each message contains Data, Sender, and Receiver fields
		json.NewEncoder(w).Encode(messages)
	case http.MethodPost:
		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			utils.ErrorMessage(w, "Invalid Request", http.StatusBadRequest)
			return
		}
		senderID,ok:=middleware.GetUserID(r)
		if !ok {
            utils.ErrorMessage(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
		if err := SaveMessageToDB(senderID, msg.Receiver, msg.Data); err != nil {
            utils.ErrorMessage(w, "Failed to save message", http.StatusInternalServerError)
            return
        }
		// messagesMutex.Lock()
		// messages = append(messages, msg)
		// messagesMutex.Unlock()
		// broadcast <- msg
		w.WriteHeader(http.StatusCreated)
	default:
		utils.ErrorMessage(w, "Method not allowed", http.StatusBadRequest)
		return
	}
}

func MessageWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("we are online")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
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
			fmt.Println("WebSocket read error:", err)
			break
		}
		broadcast <- msg
	}
}
