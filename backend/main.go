package main

import (
	"log"
	"net/http"

	"forum/database"
	"forum/handlers"
	"forum/middleware"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// API routes - Public
	mux:=http.NewServeMux()
	mux.HandleFunc("/static/", handlers.Static)
	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/posts", handlers.GetPostsHandler)
	mux.HandleFunc("/api/categories", handlers.GetCategoriesHandler)
	mux.HandleFunc("/api/posts/", handlers.GetSinglePostHandler)
	mux.HandleFunc("/api/status", handlers.GetForumStatusHandler)
	mux.HandleFunc("/api/messages", handlers.MessageHandler)
	mux.HandleFunc("/api/messaging", handlers.MessageWebSocketHandler)
	mux.HandleFunc("/api/CurrentUser", handlers.CurrentUserHandler)


	// API routes - Protected
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/auth/status", handlers.AuthStatusHandler)
	protectedMux.HandleFunc("/api/logout", handlers.LogoutHandler)
	protectedMux.HandleFunc("/api/posts", handlers.GetPostsHandler)
	protectedMux.HandleFunc("/api/posts/create", handlers.CreatePostHandler)
	protectedMux.HandleFunc("/api/comments", handlers.CreateCommentHandler)
	protectedMux.HandleFunc("/api/likes", handlers.CreatePostLikeDislikeHandler)
	protectedMux.HandleFunc("/api/comments/like", handlers.CreateCommentLikeHandler)

	// Registering routes
	mux.Handle("/api/protected/", middleware.AuthMiddleware(http.StripPrefix("/api/protected", protectedMux)))

	// Applying CORS middleware to all API routes
	handler:= middleware.CorsMiddleware(mux)

	// Start server
	log.Println("Server started on http://localhost:9111")
	log.Fatal(http.ListenAndServe(":9111", handler))
}
