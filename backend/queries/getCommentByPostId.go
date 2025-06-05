package queries

import (
	"forum/database"
	"forum/models"
)

// Get comments for a post
func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	query := `
		SELECT 
			c.id, c.user_id, u.nickname as username, c.content, c.created_at,
			(SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND is_like = 1) as likes,
			(SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND is_like = 0) as dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC`

	rows, err := database.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.Username,
			&comment.Content,
			&comment.CreatedAt,
			&comment.Likes,
			&comment.Dislikes,
		)
		if err != nil {
			return nil, err
		}
		comment.PostID = postID // Set the post ID
		comments = append(comments, comment)
	}

	return comments, nil
}
