package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"forum/database"
	"forum/utils"

	"github.com/gorilla/websocket"
)

type Message struct {
	Data any `json:"data"`
	Sender string      `json:"from,omitempty"`
	Receiver   string      `json:"to,omitempty"`
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
		msg := <-broadcast
		messagesMutex.Lock()
		messages = append(messages, msg)
		messagesMutex.Unlock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func getMessages(user_id string) []Message {
	var args interface{user_id  offset}


	query := `
			SELECT 	p.id, p.user_id, u.username, p.title, p.content, p.created_at,
			WHERE pc.category_id = ?
			ORDER BY p.created_at DESC`
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		fmt.Errorf("error querying posts: %v", err)
		return nil
	}
	defer rows.Close()
	return args
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		messagesMutex.Lock()
		defer messagesMutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	case http.MethodPost:
		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			utils.ErrorMessage(w, "Invalid Request", http.StatusBadRequest)
			return
		}
		messagesMutex.Lock()
		messages = append(messages, msg)
		messagesMutex.Unlock()
		broadcast <- msg
		w.WriteHeader(http.StatusCreated)
	default:
		utils.ErrorMessage(w, "Method not allowed", http.StatusBadRequest)
		return
	}
}

func MessageWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	clients[conn] = true
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			delete(clients, conn)
			conn.Close()
			break
		}
		broadcast <- msg
	}
}
