package handlers

import (
	"encoding/json"
	"net/http"

	"forum/logging"
	"forum/middleware"
	"forum/models"
	"forum/queries"
)

func AuthStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)

	if !ok {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
		})
		return
	}

	// Get nickname
	var user models.User
	user.ID = userID
	if err := queries.GetUserByID(&user); err != nil {
		logging.Log("[ERROR] Failed to get user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"nickname":      user.Nickname,
		"user_id":       user.ID,
	})
}
