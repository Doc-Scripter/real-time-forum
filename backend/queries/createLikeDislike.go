package queries

import (
	"database/sql"

	"forum/database"
	"forum/models"
)

// Insert a like/dislike
func CreateLikeDislike(like models.LikeDislike) error {
	var existingID int
	err := database.DB.QueryRow(
		"SELECT id FROM likes WHERE user_id = ? AND post_id = ?",
		like.UserID, like.PostID,
	).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Create new like
		_, err = database.DB.Exec(
			"INSERT INTO likes (user_id, post_id, is_like) VALUES (?, ?, ?)",
			like.UserID, like.PostID, like.IsLike,
		)
	}

	if err != nil {
		return err
	}

	// Update existing like
	_, err = database.DB.Exec(
		"UPDATE likes SET is_like = ? WHERE id = ?",
		like.IsLike, existingID,
	)
	return err
}
