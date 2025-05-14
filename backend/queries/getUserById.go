package queries

import (
	"forum/database"
	"forum/models"
)

func GetUserByID(user *models.User) error {
	return database.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", user.ID).Scan(&user.Nickname)
}
