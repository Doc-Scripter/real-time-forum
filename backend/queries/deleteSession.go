package queries

import "forum/database"

// DeleteUserSessions removes all existing sessions for a user
func DeleteUserSessions(userID int) error {
	_, err := database.DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}
