package queries

import (
	"database/sql"

	"forum/database"
	"forum/models"
)

func GetPostsByCategory(categoryID string) ([]models.Post, error) {
	query := `
		SELECT DISTINCT
			p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		ORDER BY p.created_at DESC`

	rows, err := database.DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts, err := scanPosts(rows)
	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}
	return posts, err
}
