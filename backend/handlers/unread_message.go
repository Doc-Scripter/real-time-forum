package handlers

import (
	"encoding/json"
	"net/http"

	"fmt"
	"forum/database"
	"forum/middleware"
	"forum/utils"
	"forum/logging"
)

func UnreadHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.ErrorMessage(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		unreadCount, err := getUnreadCount(userID)
		if err != nil {
			fmt.Println(err)
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

func getUnreadCount(userID int) (int, error) {
	// Query the database to count unread messages for the given user
	query := `
        SELECT COUNT(DISTINCT user_id) 
        FROM messages 
        WHERE receiver_id = ? 
        AND is_read = 0
    `
	var count int
    err := database.DB.QueryRow(query, userID).Scan(&count)
    logging.Log("[DEBUG] Unread count for user %d: %d", userID, count)
    if err != nil {
        logging.Log("[ERROR] Error getting unread count: %v", err)
        return 0, err
    }
    return count, nil
}
