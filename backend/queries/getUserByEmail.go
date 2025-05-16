package queries

import (
	"forum/database"
	"forum/models"
)

// Get user by email
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := database.DB.QueryRow("SELECT id, , email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}
