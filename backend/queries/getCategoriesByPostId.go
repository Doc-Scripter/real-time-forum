package queries

import (
	"forum/database"
	"forum/models"
)

// Get categories for a post
func GetCategoriesByPostID(postID int) ([]models.Category, error) {
	rows, err := database.DB.Query(`
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
