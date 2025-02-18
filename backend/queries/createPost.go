package queries

import (
	"errors"

	"forum/database"
	"forum/models"
)

// Insert a new post
func CreatePost(post models.Post) error {
	res, err := database.DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", post.UserID, post.Title, post.Content)
	if err != nil {
		return errors.New("failed to insert post into post table")
	}
	postID, err := res.LastInsertId()
	if err != nil {
		return errors.New("failed to get last post id")
	}
	if err := InsertPostCategories(int(postID), post.RawCategories); err != nil {
		return err
	}
	return nil
}
