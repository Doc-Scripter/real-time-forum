package queries

import (
	"database/sql"

	"forum/database"
	"forum/models"
)

// Insert a like/dislike
func CreateCommentLikeDislike(like models.LikeDislike) error {
	// First check if user already liked/disliked
	var existingID int
	err := database.DB.QueryRow(
		"SELECT id FROM comment_likes WHERE user_id = ? AND comment_id = ?",
		like.UserID, like.CommentID,
	).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Create new like
		_, err = database.DB.Exec(
			"INSERT INTO comment_likes (user_id, comment_id, is_like) VALUES (?, ?, ?)",
			like.UserID, like.CommentID, map[bool]int{true: 1, false: 0}[like.IsLike],
		)
		return err
	}

	if err != nil {
		return err
	}

	// Update existing like
	_, err = database.DB.Exec(
		"UPDATE comment_likes SET is_like = ? WHERE id = ?",
		map[bool]int{true: 1, false: 0}[like.IsLike], existingID,
	)
	return err
}
