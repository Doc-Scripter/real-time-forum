package query

import "forum/database"

// Retrieve the stats data from the tables
func GetForumStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get total likes
	var totalLikes int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM likes WHERE is_like = 1
	`).Scan(&totalLikes)
	if err != nil {
		return nil, err
	}

	// Get total comments
	var totalComments int
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM comments
	`).Scan(&totalComments)
	if err != nil {
		return nil, err
	}

	// Get total posts
	var totalPosts int
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM posts
	`).Scan(&totalPosts)
	if err != nil {
		return nil, err
	}

	// Get most active category
	var mostActiveCategory string
	var categoryCount int
	err = database.DB.QueryRow(`
		SELECT c.name, COUNT(pc.post_id) as post_count
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		GROUP BY c.id
		ORDER BY post_count DESC
		LIMIT 1
	`).Scan(&mostActiveCategory, &categoryCount)
	if err != nil {
		mostActiveCategory = "None"
		categoryCount = 0
	}

	stats["total_likes"] = totalLikes
	stats["total_comments"] = totalComments
	stats["total_posts"] = totalPosts
	stats["most_active_category"] = mostActiveCategory
	stats["category_post_count"] = categoryCount

	return stats, nil
}
