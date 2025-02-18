package handlers

import (
	"encoding/json"
	"net/http"

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

	// Get username
	var user models.User
	user.ID = userID
	if err := queries.GetUserByID(&user); err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"username":      user.Username,
		"user_id":       user.ID,
	})
}
