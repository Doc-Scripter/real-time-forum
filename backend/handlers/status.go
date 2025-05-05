package handlers

import (
	"encoding/json"
	"forum/database"
	"net/http"
	"strconv"
	"time"
)

type UserStatus struct {
	receiver string `json:"receiver"`
	username string `json:"username"`
	Online   bool   `json:"online"`
}

func GetForumStatusHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, username FROM users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []UserStatus
	for rows.Next() {
		var id int
		var username string
		if err := rows.Scan(&id, &username); err != nil {
			continue
		}
		var expiresAt time.Time
		err := database.DB.QueryRow(
			"SELECT expires_at FROM sessions WHERE user_id = ? ORDER BY expires_at DESC LIMIT 1", id,
		).Scan(&expiresAt)
		online := err == nil && expiresAt.After(time.Now())
		receiver := strconv.Itoa(id)
		
		users = append(users, UserStatus{receiver: receiver, username: username, Online: online})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
