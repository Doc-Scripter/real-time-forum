package handlers

import (
	"encoding/json"
	"html"
	"net/http"

	"forum/database"
	"forum/database/query"
	"forum/middleware"
	"forum/models"
	"forum/utils"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RenderErrorPage(w, http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		errorMessage(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var comment models.Comment
	if err := ParseJSONBody(r.Body, &comment); err != nil {
		errorMessage(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// **Escape the comment content to neutralize scripts**
	comment.Content = html.EscapeString(comment.Content)

	// **Reject comments that are empty after escaping**
	if comment.Content == "" {
		errorMessage(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	comment.UserID = userID
	if err := database.CreateComment(comment); err != nil {
		errorMessage(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment created successfully"})
}

func CreateCommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var like models.LikeDislike
	like.UserID = userID

	if err := ParseJSONBody(r.Body, &like); err != nil {
		errorMessage(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := query.UpdateCommentLikes(like); err != nil {
		http.Error(w, "Failed to update comment like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment like updated successfully"})
}
