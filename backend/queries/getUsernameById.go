package queries

import (
    "forum/database"
)

func GetUsernameByID(userID int, username *string) error {
    return database.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(username)
}
