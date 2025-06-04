package main

import (
	//"log"
	"fmt"
	"net/http"
	"os"

	"forum/database"
	"forum/handlers"
	"forum/logging"
	"forum/middleware"
)

func main() {
	// Check for flags
	if len(os.Args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: No flags are allowed when running the program\n")
		os.Exit(1)
	}

	// Initialize logger
	logging.InitializeLogger()

	// Initialize database
	if err := database.InitDB(); err != nil {
		logging.TerminalLog("Database initialization failed: %v", err)
		return
	}

	// Create main mux
	mux := http.NewServeMux()

	// Register public routes
	mux.HandleFunc("/static/", handlers.Static)
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/register", handlers.RegisterHandler)
	mux.HandleFunc("/", handlers.Index)

	// Create protected mux
	protectedMux := http.NewServeMux()
	// protectedMux.HandleFunc("/api/posts", handlers.GetPostsHandler)
	protectedMux.HandleFunc("/api/categories", handlers.GetCategoriesHandler)
	protectedMux.HandleFunc("/api/posts/", handlers.GetSinglePostHandler)
	protectedMux.HandleFunc("/api/messaging", handlers.MessageWebSocketHandler)
	protectedMux.HandleFunc("/api/messages", handlers.MessageHandler)
	protectedMux.HandleFunc("/api/unread", handlers.UnreadHandler)
	protectedMux.HandleFunc("/api/status", handlers.GetForumStatusHandler)
	protectedMux.HandleFunc("/api/auth/status", handlers.AuthStatusHandler)
	protectedMux.HandleFunc("/api/logout", handlers.LogoutHandler)
	protectedMux.HandleFunc("/api/posts", handlers.GetPostsHandler)
	protectedMux.HandleFunc("/api/posts/create", handlers.CreatePostHandler)
	protectedMux.HandleFunc("/api/comments", handlers.CreateCommentHandler)
	protectedMux.HandleFunc("/api/likes", handlers.CreatePostLikeDislikeHandler)
	protectedMux.HandleFunc("/api/comments/like", handlers.CreateCommentLikeHandler)

	// Register protected routes
	mux.Handle("/api/protected/", middleware.AuthMiddleware(http.StripPrefix("/api/protected", protectedMux)))

	// Apply CORS middleware
	handler := middleware.CorsMiddleware(mux)

	// Start server
	logging.TerminalLog("Server started on http://localhost:9111")
	if err := http.ListenAndServe(":9111", handler); err != nil {
		logging.TerminalLog("Server failed to start: %v", err)
		return
	}
}
