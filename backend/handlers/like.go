package handlers

import (
	"encoding/json"
	"net/http"

	"forum/database"
	"forum/middleware"
	"forum/models"
)

func CreatePostLikeDislikeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		errorMessage(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var like models.LikeDislike

	if err := ParseJSONBody(r.Body, &like); err != nil {
		errorMessage(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set the user ID from the authenticated user
	like.UserID = userID
	if err := database.CreateLikeDislike(like); err != nil {
		http.Error(w, "Failed to like/dislike", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Like/Dislike created successfully"})
}
