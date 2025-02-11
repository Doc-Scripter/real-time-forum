package query

import (
	"forum/database"
	"forum/models"
)

func UpdateCommentLikes(like models.LikeDislike) (err error) {
	// Insert or update the comment like in the database
	_, err = database.DB.Exec(`
		INSERT INTO comment_likes (user_id, comment_id, is_like)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, comment_id)
		DO UPDATE SET is_like = ?`,
		like.UserID, like.CommentID, like.IsLike, like.IsLike)
	return
}
