package handlers

import (
	"encoding/json"
	"html"
	"net/http"
	"strconv"

	"forum/logging"
	"forum/middleware"
	"forum/models"
	"forum/queries"
	"forum/utils"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorMessage(w, "Hey! Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.GetUserID(r)
	if !ok {
		utils.ErrorMessage(w, "Unauthorized: No user ID", http.StatusUnauthorized)
		return
	}

	var postData models.Post
	postData.UserID = userId

	if err := ParseJSONBody(r.Body, &postData); err != nil {
		logging.Log("Error parsing JSON body: %v", err)
		utils.ErrorMessage(w, "Sorry! We didnt get those post details. Try again.", http.StatusInternalServerError)
		return
	}

	postData.Title = html.EscapeString(postData.Title)
	postData.Content = html.EscapeString(postData.Content)

	if err := ValidatePost(postData); err != nil {
		utils.ErrorMessage(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	if err := queries.CreatePost(postData); err != nil {
		logging.Log("Error creating post: %v", err)
		utils.ErrorMessage(w, "Ooops! We couldn't get your post. Try again...", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created successfully",
	})
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorMessage(w, "Hey! Method not allowed", http.StatusMethodNotAllowed)
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
		} else {
			logging.Log("Error converting page number to integer: %v", err)
		}
	}

	offset := (pageNum - 1) * queries.PostsPerPage
	userID, _ := middleware.GetUserID(r)
	posts, err := queries.GetAllPosts(offset, category, filter, userID)
	if err != nil {
		logging.Log("Error fetching posts: %v", err)
		utils.ErrorMessage(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetSinglePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	postID := r.URL.Path[len("/api/posts/"):]
	id, err := strconv.Atoi(postID)
	if err != nil {
		logging.Log("Error converting post ID to integer: %v", err)
		utils.ErrorMessage(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := queries.GetPostByID(id)
	if err != nil {
		logging.Log("Error fetching post: %v", err)
		utils.ErrorMessage(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
