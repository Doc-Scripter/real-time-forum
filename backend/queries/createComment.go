package queries

import (
	"forum/database"
	"forum/models"
)

// Insert a new comment
func CreateComment(comment models.Comment) error {
	_, err := database.DB.Exec("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)", comment.UserID, comment.PostID, comment.Content)
	return err
}
