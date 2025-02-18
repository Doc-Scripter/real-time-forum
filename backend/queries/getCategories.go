package queries

import (
	"forum/database"
	"forum/models"
)

// Get all categories
func GetCategories() ([]models.Category, error) {
	rows, err := database.DB.Query("SELECT id, name FROM categories ORDER BY name")
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
