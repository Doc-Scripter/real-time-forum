package query

import (
	"forum/database"
	"forum/models"
)

func GetUserByID(user models.User) error {
	return database.DB.QueryRow("SELECT username FROM users WHERE id = ?", user.ID).Scan(&user.Username)
}
