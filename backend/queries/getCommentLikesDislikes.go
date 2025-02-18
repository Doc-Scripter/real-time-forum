package queries

import (
	"forum/database"
	"forum/models"
)

// Get likes/dislikes for a post
func GetCommentLikesDislikes(postID int) ([]models.LikeDislike, error) {
	query := "SELECT id, user_id, comment_id, is_like, created_at FROM likes WHERE comment_id = ?"

	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []models.LikeDislike
	for rows.Next() {
		var like models.LikeDislike
		err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.IsLike, &like.CreatedAt)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}
	return likes, nil
}
