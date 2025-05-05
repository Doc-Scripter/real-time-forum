package handlers

import (
    "log"
    "net/http"
    "forum/middleware"
    "forum/websocket"
    "forum/queries"

    ws "github.com/gorilla/websocket"
    "github.com/google/uuid"
)

var upgrader = ws.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // In production, implement proper origin checking
    },
}

var wsManager = websocket.NewManager()

func init() {
    go wsManager.Run()
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.GetUserID(r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Get username from the session
    var username string
    if err := queries.GetUsernameByID(userID, &username); err != nil {
        http.Error(w, "Failed to get user info", http.StatusInternalServerError)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }

    client := &websocket.Client{
        ID:       uuid.New().String(),
        UserID:   userID,
        Username: username,
        Conn:     conn,
        Manager:  wsManager,
        Send:     make(chan []byte, 256),
    }

    wsManager.Register <- client

    // Start goroutines for reading and writing
    go client.ReadPump()
    go client.WritePump()
}