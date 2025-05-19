package queries

import (
	"forum/database"
	"forum/models"
)

func GetPostByID(postID int) (models.Post, error) {
	query := `
		SELECT 
			p.id, p.user_id, u.nickname, p.title, p.content, p.created_at,
			(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
			(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`

	var post models.Post
	err := database.DB.QueryRow(query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Nickname,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.Likes,
		&post.Dislikes,
	)
	if err != nil {
		return post, err
	}

	// Get categories
	categories, err := GetCategoriesByPostID(post.ID)
	if err == nil {
		post.Categories = categories
	}

	// Get comments
	comments, err := GetCommentsByPostID(post.ID)
	if err == nil {
		post.Comments = comments
	}

	return post, nil
}
