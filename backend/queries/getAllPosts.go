package queries

import (
	"fmt"

	"forum/database"
	"forum/models"
)

const (
	PostsPerPage = 8
)

// GetAllPosts with support for pagination
func GetAllPosts(offset int, category, filter string, userID int) ([]models.Post, error) {
	var query string
	var args []interface{}

	// Check if filter is a category filter
	if len(category) > 9 && category[:9] == "category-" {
		categoryID := category[9:]
		query = `
			SELECT DISTINCT
				p.id, p.user_id, u.nickname, p.title, p.content, p.created_at,
				(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
				(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
			FROM posts p
			JOIN users u ON p.user_id = u.id
			JOIN post_categories pc ON p.id = pc.post_id
			WHERE pc.category_id = ?
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?`
		args = []interface{}{categoryID, PostsPerPage, offset}
	} else {
		switch filter {
		case "my-posts":
			query = `
				SELECT DISTINCT
					p.id, p.user_id, u.nickname, p.title, p.content, p.created_at,
					(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
					(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
				FROM posts p
				JOIN users u ON p.user_id = u.id
				WHERE p.user_id = ?
				ORDER BY p.created_at DESC
				LIMIT ? OFFSET ?`
			args = []interface{}{userID, PostsPerPage, offset}

		case "liked-posts":
			query = `
				SELECT DISTINCT
					p.id, p.user_id, u.nickname, p.title, p.content, p.created_at,
					(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
					(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
				FROM posts p
				JOIN users u ON p.user_id = u.id
				JOIN likes l ON p.id = l.post_id
				WHERE l.user_id = ? AND l.is_like = 1
				ORDER BY p.created_at DESC
				LIMIT ? OFFSET ?`
			args = []interface{}{userID, PostsPerPage, offset}

		default:
			query = `
				SELECT DISTINCT
					p.id, p.user_id, u.nickname, p.title, p.content, p.created_at,
					(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
					(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
				FROM posts p
				JOIN users u ON p.user_id = u.id
				ORDER BY p.created_at DESC
				LIMIT ? OFFSET ?`
			args = []interface{}{PostsPerPage, offset}
		}
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %v", err)
	}
	defer rows.Close()

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
			&post.Likes,
			&post.Dislikes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}

		// Get categories for each post
		categories, err := GetCategoriesByPostID(post.ID)
		if err == nil {
			post.Categories = categories
		}

		// Get comments for each post
		comments, err := GetCommentsByPostID(post.ID)
		if err == nil {
			post.Comments = comments
		}

		posts = append(posts, post)
	}

	return posts, nil
}
