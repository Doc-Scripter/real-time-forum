package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"forum/database"
	"forum/logging"
	"forum/middleware"
)

// UserStatus represents the current status of a user in the forum
type UserStatus struct {
	Receiver         int        `json:"receiver"`
	Nickname         string     `json:"nickname"`
	Online           bool       `json:"online"`
	LastConversation *time.Time `json:"last_conversation,omitempty"`
}

// GetForumStatusHandler retrieves the status of all users in the forum
// Returns a JSON array of UserStatus containing online status and last conversation time
func GetForumStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	query := `
		SELECT DISTINCT u.id, u.nickname,
			   CASE WHEN s.expires_at > datetime('now') THEN 1 ELSE 0 END as online,
			   m.last_conversation
		FROM users u
		LEFT JOIN (
			SELECT user_id, MAX(expires_at) as expires_at
			FROM sessions 
			GROUP BY user_id
		) s ON u.id = s.user_id
		LEFT JOIN (
			SELECT 
				CASE 
					WHEN user_id = ? THEN receiver_id 
					ELSE user_id 
				END as conversation_partner,
				MAX(created_at) as last_conversation
			FROM messages
			WHERE user_id = ? OR receiver_id = ?
			GROUP BY conversation_partner
		) m ON u.id = m.conversation_partner
		WHERE u.id != ?
		ORDER BY 
			CASE WHEN m.last_conversation IS NOT NULL THEN 0 ELSE 1 END,
			m.last_conversation DESC,
			u.nickname ASC
	`

	rows, err := database.DB.Query(query, userID, userID, userID, userID)
	if err != nil {
		logging.Log("[ERROR] : fetching users: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []UserStatus
	for rows.Next() {
		var id int
		var nickname string
		var online bool
		var lastConversationStr sql.NullString

		if err := rows.Scan(&id, &nickname, &online, &lastConversationStr); err != nil {
			logging.Log("[ERROR] : fetching user status: %v", err)
			continue
		}

		user := UserStatus{
			Receiver: id,
			Nickname: nickname,
			Online:   online,
		}

		if lastConversationStr.Valid && lastConversationStr.String != "" {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05", lastConversationStr.String); err == nil {
				user.LastConversation = &parsedTime
			} else {
				logging.Log("[ERROR] : fetching user status: %v", err)
			}
		}

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
