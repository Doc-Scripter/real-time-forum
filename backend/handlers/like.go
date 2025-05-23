package handlers

import (
	"encoding/json"
	"net/http"

	"forum/logging"
	"forum/middleware"
	"forum/models"
	"forum/queries"
	"forum/utils"
)

func CreatePostLikeDislikeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.ErrorMessage(w, "Please login to continue...", http.StatusUnauthorized)
		return
	}

	var like models.LikeDislike

	if err := ParseJSONBody(r.Body, &like); err != nil {
		logging.Log("[ERROR] Error parsing JSON: %v", err)
		utils.ErrorMessage(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set the user ID from the authenticated user
	like.UserID = userID
	if err := queries.CreateLikeDislike(like); err != nil {
		logging.Log("[ERROR] Error creating like/dislike: %v", err)
		utils.ErrorMessage(w, "Opps... Failed to like/dislike", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Like/Dislike created successfully"})
}
