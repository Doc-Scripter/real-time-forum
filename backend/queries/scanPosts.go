package queries

import (
	"database/sql"
	"fmt"

	"forum/models"
)

// Helper function to scan posts from rows
func scanPosts(rows *sql.Rows) ([]models.Post, error) {
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Nickname,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}

		// Get categories for each post
		categories, err := GetCategoriesByPostID(post.ID)
		if err == nil {
			post.Categories = categories
		}

		posts = append(posts, post)
	}
	return posts, nil
}
