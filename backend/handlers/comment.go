package handlers

import (
	"encoding/json"
	"html"
	"net/http"

	"forum/logging"
	"forum/middleware"
	"forum/models"
	"forum/queries"
	"forum/utils"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorMessage(w, "Hey! Method not allowed...", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.ErrorMessage(w, "Please login to continue...", http.StatusUnauthorized)
		return
	}

	var comment models.Comment
	if err := ParseJSONBody(r.Body, &comment); err != nil {
		logging.Log("[ERROR] Error parsing JSON: %v", err)
		utils.ErrorMessage(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// **Escape the comment content to neutralize scripts**
	comment.Content = html.EscapeString(comment.Content)

	// **Reject comments that are empty after escaping**
	if comment.Content == "" {
		utils.ErrorMessage(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}

	comment.UserID = userID
	if err := queries.CreateComment(comment); err != nil {
		utils.ErrorMessage(w, "Opps! Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment created successfully"})
}

func CreateCommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorMessage(w, "Method Not Found", http.StatusMethodNotAllowed)

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
		logging.Log("[ERROR] Error parsing JSON: %v", err)
		utils.ErrorMessage(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := queries.UpdateCommentLikes(like); err != nil {
		logging.Log("[ERROR] Error updating comment likes: %v", err)
		utils.ErrorMessage(w, "Sorry! We couldn't update your reaction. Try again.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment like updated successfully"})
}
