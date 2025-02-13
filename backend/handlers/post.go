package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"

	"forum/database"
	"forum/middleware"
	"forum/models"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.GetUserID(r)
	if !ok {
		errorMessage(w, "Unauthorized: No user ID", http.StatusUnauthorized)
		return
	}

	var postData models.Post
	postData.UserID = userId

	if err := ParseJSONBody(r.Body, &postData); err != nil {
		errorMessage(w, "Cannot decode post body", http.StatusInternalServerError)
		return
	}

	postData.Title = html.EscapeString(postData.Title)
	postData.Content = html.EscapeString(postData.Content)

	if err := ValidatePost(postData); err != nil {
		errorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.CreatePost(postData); err != nil {
		errorMessage(w, "Cannot create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created successfully",
	})
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get page from query parameters
	page := r.URL.Query().Get("page")
	category := r.URL.Query().Get("category")
	filter := r.URL.Query().Get("filter")

	pageNum := 1
	if page != "" {
		if num, err := strconv.Atoi(page); err == nil {
			pageNum = num
		}
	}

	offset := (pageNum - 1) * database.PostsPerPage
	userID, _ := middleware.GetUserID(r)
	fmt.Println("GetPostsHandler: ", category, filter, offset, userID)
	posts, err := database.GetAllPosts(offset, category, filter, userID)
	if err != nil {
		errorMessage(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetSinglePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	postID := r.URL.Path[len("/api/posts/"):]
	id, err := strconv.Atoi(postID)
	if err != nil {
		errorMessage(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := database.GetPostByID(id)
	if err != nil {
		errorMessage(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
