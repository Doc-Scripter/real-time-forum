package queries

import (
	"errors"

	"forum/database"
)

// Insert categories for a post in categories table
func InsertPostCategories(postID int, categories []int) error {
	for _, category := range categories {
		_, err := database.DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?,?)", postID, category)
		if err != nil {
			return errors.New("failed to add categories")
		}
	}
	return nil
}
