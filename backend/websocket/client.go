package websocket

import (
    "encoding/json"
    "log"
    "sync"
    "time"
    "github.com/gorilla/websocket"
)

type Client struct {
    ID       string
    UserID   int
    Username string
    Conn     *websocket.Conn
    Manager  *Manager
    mu       sync.Mutex
    Send     chan []byte
}

type Message struct {
    Type      string      `json:"type"`
    Content   string      `json:"content,omitempty"`
    From      string      `json:"from,omitempty"`
    To        string      `json:"to,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
    Data      interface{} `json:"data,omitempty"`
}

func (c *Client) ReadPump() {
    defer func() {
        c.Manager.Unregister <- c
        c.Conn.Close()
    }()

    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }

        var msg Message
        if err := json.Unmarshal(message, &msg); err != nil {
            continue
        }

        msg.From = c.Username
        msg.Timestamp = time.Now()

        // Handle different message types
        switch msg.Type {
        case "chat":
            // Find recipient and send message
            if recipient := c.Manager.GetClientByUsername(msg.To); recipient != nil {
                recipientMsg, _ := json.Marshal(msg)
                recipient.Send <- recipientMsg
                
                // Also send back to sender for confirmation
                c.Send <- recipientMsg
            }
        case "status":
            c.Manager.UpdateStatus(c, msg.Content)
        }
    }
}

func (c *Client) WritePump() {
    defer func() {
        c.Conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            c.mu.Lock()
            w, err := c.Conn.NextWriter(websocket.TextMessage)
            if err != nil {
                c.mu.Unlock()
                return
            }

            w.Write(message)
            c.mu.Unlock()

            if err := w.Close(); err != nil {
                return
            }
        }
    }
}