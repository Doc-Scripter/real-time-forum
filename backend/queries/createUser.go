package queries

import (

	"fmt"

	"forum/database"
	"time"
	"forum/models"
)

// Insert a new user
func CreateUser(user models.User) (int, error) {
	// Set created_at timestamp
	user.CreatedAt = time.Now()

	// Prepare the INSERT statement
	query := "INSERT INTO users (nickname, age, gender, first_name, last_name, email, password, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

	// Execute the query
	result, err := database.DB.Exec(query,
		user.Nickname,
		user.Age,
		user.Gender,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.CreatedAt,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return int(lastID), nil
}
