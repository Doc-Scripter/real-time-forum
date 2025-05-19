package models

type Post struct {
	ID            int        `json:"id"`
	UserID        int        `json:"user_id"`
	Nickname      string     `json:"nickname"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	CreatedAt     string     `json:"created_at"`
	Likes         int        `json:"likes"`
	Dislikes      int        `json:"dislikes"`
	Comments      []Comment  `json:"comments"`
	RawCategories []int      `json:"raw_categories"`
	Categories    []Category `json:"categories,omitempty"` // Add this field
}
