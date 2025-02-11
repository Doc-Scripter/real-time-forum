package main

import (
	"log"
	"net/http"

	"forum/database"
	"forum/handlers"
	"forum/middleware"
)

func main() {
	err := database.InitDB("./forum.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v. Ensure the file exists and is accessible.", err)
	}

	// API routes - Public
	http.HandleFunc("/static/", handlers.Static)
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.HandleFunc("/api/posts", handlers.GetPostsHandler)
	http.HandleFunc("/api/categories", handlers.GetCategoriesHandler)
	http.HandleFunc("/api/posts/", handlers.GetSinglePostHandler)
	http.HandleFunc("/api/stats", handlers.GetForumStatsHandler)

	// API routes - Protected
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/logout", handlers.LogoutHandler)
	protectedMux.HandleFunc("/api/posts", handlers.GetPostsHandler)
	protectedMux.HandleFunc("/api/posts/create", handlers.CreatePostHandler)
	protectedMux.HandleFunc("/api/comments", handlers.CreateCommentHandler)
	protectedMux.HandleFunc("/api/likes", handlers.CreatePostLikeDislikeHandler)
	protectedMux.HandleFunc("/api/auth/status", handlers.AuthStatusHandler)
	protectedMux.HandleFunc("/api/comments/like", handlers.CreateCommentLikeHandler)

	// Registering routes
	http.Handle("/api/protected/", middleware.AuthMiddleware(http.StripPrefix("/api/protected", protectedMux)))

	// Applying CORS middleware to all API routes
	http.Handle("/api/", middleware.CorsMiddleware(http.DefaultServeMux))

	// Start server
	log.Println("Server started on http://localhost:9111")
	log.Fatal(http.ListenAndServe(":9111", nil))
}
