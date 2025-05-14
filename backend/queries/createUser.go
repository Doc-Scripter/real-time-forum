package queries

import (
	"forum/database"
	"forum/models"
)

// Insert a new user
func CreateUser(user models.User) error {
	_, err := database.DB.Exec("INSERT INTO users (, email, password) VALUES (?, ?, ?)",  user.Email, user.Password)
	return err
}
