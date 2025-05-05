package websocket

import (
    "encoding/json"
    "sync"
)

type Manager struct {
    Clients    map[*Client]bool
    Broadcast  chan []byte
    Register   chan *Client
    Unregister chan *Client
    mu         sync.RWMutex
    // Track online status
    OnlineUsers map[int]string // UserID -> Username
}

func NewManager() *Manager {
    return &Manager{
        Clients:     make(map[*Client]bool),
        Broadcast:   make(chan []byte),
        Register:    make(chan *Client),
        Unregister:  make(chan *Client),
        OnlineUsers: make(map[int]string),
    }
}

func (m *Manager) Run() {
    for {
        select {
        case client := <-m.Register:
            m.mu.Lock()
            m.Clients[client] = true
            m.OnlineUsers[client.UserID] = client.Username
            m.mu.Unlock()
            m.broadcastUserStatus()

        case client := <-m.Unregister:
            m.mu.Lock()
            if _, ok := m.Clients[client]; ok {
                delete(m.Clients, client)
                delete(m.OnlineUsers, client.UserID)
                close(client.Send)
            }
            m.mu.Unlock()
            m.broadcastUserStatus()

        case message := <-m.Broadcast:
            m.mu.RLock()
            for client := range m.Clients {
                select {
                case client.Send <- message:
                default:
                    close(client.Send)
                    delete(m.Clients, client)
                }
            }
            m.mu.RUnlock()
        }
    }
}

func (m *Manager) GetClientByUsername(username string) *Client {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    for client := range m.Clients {
        if client.Username == username {
            return client
        }
    }
    return nil
}

func (m *Manager) UpdateStatus(client *Client, status string) {
    m.mu.Lock()
    defer m.mu.Unlock()

    if status == "online" {
        m.OnlineUsers[client.UserID] = client.Username
    } else {
        delete(m.OnlineUsers, client.UserID)
    }
    m.broadcastUserStatus()
}

func (m *Manager) broadcastUserStatus() {
    m.mu.RLock()
    onlineUsers := make(map[string]interface{})
    for userID, username := range m.OnlineUsers {
        onlineUsers[username] = map[string]interface{}{
            "id":     userID,
            "status": "online",
        }
    }
    m.mu.RUnlock()

    statusMsg := Message{
        Type: "status_update",
        Data: onlineUsers,
    }

    if msgBytes, err := json.Marshal(statusMsg); err == nil {
        m.Broadcast <- msgBytes
    }
}