package handlers

import (
	"forum/logging"
	"net/http"
	"time"
)

// MessageWebSocketHandler handles WebSocket connections for real-time messaging.
// It upgrades HTTP connection to WebSocket, manages client connections,
// and broadcasts received messages to all connected clients.
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
			logging.Log("[INFO] : reading message: %v", err)
			break
		}
		if msg.Time == "" {
			msg.Time = time.Now().Format("15:04:05")
		}

		broadcast <- msg
	}
}
