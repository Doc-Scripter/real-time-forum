package queries

import (
	"forum/database"
	"forum/models"
)

// Insert a new user
func CreateUser(user models.User) error {
	_, err := database.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	return err
}
