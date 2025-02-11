package database

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/models"
)

const (
	PostsPerPage = 8
)

// Insert a new user
func CreateUser(user models.User) error {
	_, err := DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	return err
}

// Get user by email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := DB.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	return user, err
}

// Insert a new post
func CreatePost(post models.Post) error {
	res, err := DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", post.UserID, post.Title, post.Content)
	if err != nil {
		return errors.New("failed to insert post into post table")
	}
	postID, err := res.LastInsertId()
	if err != nil {
		return errors.New("failed to get last post id")
	}
	if err := InsertPostCategories(int(postID), post.RawCategories); err != nil {
		return err
	}
	return nil
}

// Insert categories for a post in categories table
func InsertPostCategories(postID int, categories []int) error {
	for _, category := range categories {
		_, err := DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?,?)", postID, category)
		if err != nil {
			return errors.New("failed to add categories")
		}
	}
	return nil
}

// GetAllPosts with support for pagination
func GetAllPosts(offset int, filter string, userID int) ([]models.Post, error) {
	var query string
	var args []interface{}

	// Check if filter is a category filter
	if len(filter) > 9 && filter[:9] == "category-" {
		categoryID := filter[9:]
		query = `
			SELECT DISTINCT
				p.id, p.user_id, u.username, p.title, p.content, p.created_at,
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
					p.id, p.user_id, u.username, p.title, p.content, p.created_at,
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
					p.id, p.user_id, u.username, p.title, p.content, p.created_at,
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
					p.id, p.user_id, u.username, p.title, p.content, p.created_at,
					(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
					(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
				FROM posts p
				JOIN users u ON p.user_id = u.id
				ORDER BY p.created_at DESC
				LIMIT ? OFFSET ?`
			args = []interface{}{PostsPerPage, offset}
		}
	}

	rows, err := DB.Query(query, args...)
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
			&post.Username,
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

func GetUserPosts(userID int) ([]models.Post, error) {
	query := `
		SELECT DISTINCT
			p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

// Insert a new comment
func CreateComment(comment models.Comment) error {
	_, err := DB.Exec("INSERT INTO comments (user_id, post_id, content) VALUES (?, ?, ?)", comment.UserID, comment.PostID, comment.Content)
	return err
}

// Get comments for a post
func GetCommentsByPostID(postID int) ([]models.Comment, error) {
	query := `
		SELECT 
			c.id, c.user_id, u.username, c.content, c.created_at,
			(SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND is_like = 1) as likes,
			(SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND is_like = 0) as dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC`

	rows, err := DB.Query(query, postID)
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
		comments = append(comments, comment)
	}

	return comments, nil
}

// Insert a like/dislike
func CreateLikeDislike(like models.LikeDislike) error {
	var existingID int
	err := DB.QueryRow(
		"SELECT id FROM likes WHERE user_id = ? AND post_id = ?",
		like.UserID, like.PostID,
	).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Create new like
		_, err = DB.Exec(
			"INSERT INTO likes (user_id, post_id, is_like) VALUES (?, ?, ?)",
			like.UserID, like.PostID, like.IsLike,
		)
		return err
	}

	if err != nil {
		return err
	}

	// Update existing like
	_, err = DB.Exec(
		"UPDATE likes SET is_like = ? WHERE id = ?",
		like.IsLike, existingID,
	)
	return err
}

// Insert a like/dislike
func CreateCommentLikeDislike(like models.LikeDislike) error {
	// First check if user already liked/disliked
	var existingID int
	err := DB.QueryRow(
		"SELECT id FROM comment_likes WHERE user_id = ? AND comment_id = ?",
		like.UserID, like.CommentID,
	).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Create new like
		_, err = DB.Exec(
			"INSERT INTO comment_likes (user_id, comment_id, is_like) VALUES (?, ?, ?)",
			like.UserID, like.CommentID, map[bool]int{true: 1, false: 0}[like.IsLike],
		)
		return err
	}

	if err != nil {
		return err
	}

	// Update existing like
	_, err = DB.Exec(
		"UPDATE comment_likes SET is_like = ? WHERE id = ?",
		map[bool]int{true: 1, false: 0}[like.IsLike], existingID,
	)
	return err
}

// Get likes/dislikes for a post
func GetLikesDislikes(postID int) ([]models.LikeDislike, error) {
	query := "SELECT id, user_id, post_id, is_like, created_at FROM likes WHERE post_id = ?"

	rows, err := DB.Query(query, postID)
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

// Get likes/dislikes for a post
func GetCommentLikesDislikes(postID int) ([]models.LikeDislike, error) {
	query := "SELECT id, user_id, comment_id, is_like, created_at FROM likes WHERE comment_id = ?"

	rows, err := DB.Query(query, postID)
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

func GetLikedPosts(userID int) ([]models.Post, error) {
	query := `
		SELECT DISTINCT
			p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN likes l ON p.id = l.post_id
		WHERE l.user_id = ? AND l.is_like = 1
		ORDER BY p.created_at DESC`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

// Get all categories
func GetCategories() ([]models.Category, error) {
	rows, err := DB.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// Get categories for a post
func GetCategoriesByPostID(postID int) ([]models.Category, error) {
	rows, err := DB.Query(`
		SELECT c.id, c.name 
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetPostsByCategory(categoryID string) ([]models.Post, error) {
	query := `
		SELECT DISTINCT
			p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		ORDER BY p.created_at DESC`

	rows, err := DB.Query(query, categoryID)
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

// Helper function to scan posts from rows
func scanPosts(rows *sql.Rows) ([]models.Post, error) {
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Username,
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

// DeleteUserSessions removes all existing sessions for a user
func DeleteUserSessions(userID int) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}

func GetPostByID(postID int) (models.Post, error) {
	query := `
		SELECT 
			p.id, p.user_id, u.username, p.title, p.content, p.created_at,
			(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 1) as likes,
			(SELECT COUNT(*) FROM likes WHERE post_id = p.id AND is_like = 0) as dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?`

	var post models.Post
	err := DB.QueryRow(query, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Username,
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

// Retrieve the stats data from the tables
func GetForumStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get total likes
	var totalLikes int
	err := DB.QueryRow(`
		SELECT COUNT(*) FROM likes WHERE is_like = 1
	`).Scan(&totalLikes)
	if err != nil {
		return nil, err
	}

	// Get total comments
	var totalComments int
	err = DB.QueryRow(`
		SELECT COUNT(*) FROM comments
	`).Scan(&totalComments)
	if err != nil {
		return nil, err
	}

	// Get total posts
	var totalPosts int
	err = DB.QueryRow(`
		SELECT COUNT(*) FROM posts
	`).Scan(&totalPosts)
	if err != nil {
		return nil, err
	}

	// Get most active category
	var mostActiveCategory string
	var categoryCount int
	err = DB.QueryRow(`
		SELECT c.name, COUNT(pc.post_id) as post_count
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		GROUP BY c.id
		ORDER BY post_count DESC
		LIMIT 1
	`).Scan(&mostActiveCategory, &categoryCount)
	if err != nil {
		mostActiveCategory = "None"
		categoryCount = 0
	}

	stats["total_likes"] = totalLikes
	stats["total_comments"] = totalComments
	stats["total_posts"] = totalPosts
	stats["most_active_category"] = mostActiveCategory
	stats["category_post_count"] = categoryCount

	return stats, nil
}
