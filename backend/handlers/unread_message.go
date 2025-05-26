package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"forum/database"
	"forum/logging"
	"forum/middleware"
	"forum/utils"
)

func UnreadHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.ErrorMessage(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get the optional exclude_sender parameter
		excludeSenderStr := r.URL.Query().Get("exclude_sender")
		var excludeSender *int
		if excludeSenderStr != "" {
			if id, err := strconv.Atoi(excludeSenderStr); err == nil {
				excludeSender = &id
			}else{
				logging.Log("[ERROR] : Invalid exclude_sender parameter: %s", excludeSenderStr)
			}
		}

		unreadCount, err := getUnreadCount(userID, excludeSender)
		if err != nil {
			logging.Log("[ERROR] Error getting unread count: %v", err)
			utils.ErrorMessage(w, "Failed to fetch unread message count", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(unreadCount)
		return

	default:
		utils.ErrorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func getUnreadCount(userID int, excludeSender *int) (int, error) {
	// Query the database to count unread messages for the given user
	// Exclude messages from the currently opened chat if specified
	var query string
	var args []interface{}

	if excludeSender != nil {
		query = `
			SELECT COUNT(DISTINCT user_id) 
			FROM messages 
			WHERE receiver_id = ? 
			AND is_read = 0
			AND user_id != ?
		`
		args = []interface{}{userID, *excludeSender}
	} else {
		query = `
			SELECT COUNT(DISTINCT user_id) 
			FROM messages 
			WHERE receiver_id = ? 
			AND is_read = 0
		`
		args = []interface{}{userID}
	}

	var count int
	err := database.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		logging.Log("[ERROR] Error getting unread count: %v", err)
		return 0, err
	}
	return count, nil
}
