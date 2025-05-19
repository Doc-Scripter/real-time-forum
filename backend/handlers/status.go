package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"forum/database"
	"forum/middleware"
)

type UserStatus struct {
	Receiver int    `json:"receiver"`
	Nickname string `json:"nickname"`
	Online   bool   `json:"online"`
}

func GetForumStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	rows, err := database.DB.Query("SELECT id, nickname FROM users WHERE id != ?", userID)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []UserStatus
	for rows.Next() {
		var id int
		var nickname string
		if err := rows.Scan(&id, &nickname); err != nil {
			continue
		}

		var expiresAt time.Time
		err := database.DB.QueryRow(
			"SELECT expires_at FROM sessions WHERE user_id = ? ORDER BY expires_at DESC LIMIT 1", id,
		).Scan(&expiresAt)
		online := err == nil && expiresAt.After(time.Now())

		users = append(users, UserStatus{Receiver: id, Nickname: nickname, Online: online})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
