package handlers

import (
	"encoding/json"
	"net/http"

	"forum/database"
	"forum/database/query"
	"forum/middleware"
	"forum/models"
	"forum/utils"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorMessage(w, "Hey! Method not allowed...", http.StatusMethodNotAllowed)
		return
	}
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.ErrorMessage(w, "Please login to continue...", http.StatusUnauthorized)
		return
	}

	var comment models.Comment
	if err := ParseJSONBody(r.Body, &comment); err != nil {
		utils.ErrorMessage(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Override the user_id with the authenticated user's ID
	comment.UserID = userID
	if err := database.CreateComment(comment); err != nil {
		utils.ErrorMessage(w, "Opps! Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment created successfully"})
}

func CreateCommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RenderErrorPage(w, http.StatusMethodNotAllowed)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.ErrorMessage(w, "Please login to continue...", http.StatusUnauthorized)
		return
	}
	var like models.LikeDislike
	like.UserID = userID

	if err := ParseJSONBody(r.Body, &like); err != nil {
		utils.ErrorMessage(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := query.UpdateCommentLikes(like); err != nil {
		utils.ErrorMessage(w, "Sorry! We couldn't update your reaction. Try again.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment like updated successfully"})
}
