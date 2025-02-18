package queries

import (
	"forum/database"
	"forum/models"
)

func GetUserPosts(userID int) ([]models.Post, error) {
	query := `
		SELECT DISTINCT
			p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}
