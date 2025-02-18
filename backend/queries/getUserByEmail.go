package queries

import (
	"forum/database"
	"forum/models"
)

// Get user by email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	return user, err
}
