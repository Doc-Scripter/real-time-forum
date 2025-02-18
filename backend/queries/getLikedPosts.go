package queries

import (
	"forum/database"
	"forum/models"
)

func GetLikedPosts(userID int) ([]models.Post, error) {
	query := `
		SELECT DISTINCT
			p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN likes l ON p.id = l.post_id
		WHERE l.user_id = ? AND l.is_like = 1
		ORDER BY p.created_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}
