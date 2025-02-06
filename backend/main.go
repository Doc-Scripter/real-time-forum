package main

import (
	"log"
	"net/http"
	"path/filepath"

	"forum/database"
	"forum/handlers"
	"forum/middleware"
	"forum/utils"
)

func serveTemplate(w http.ResponseWriter, r *http.Request, templatePath string) {
	// Get the absolute path to the template file
	path := filepath.Join("../frontend/templates", templatePath)
	http.ServeFile(w, r, path)
}

func main() {
	// Initialize database
	err := database.InitDB("./forum.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../frontend/static"))))

	// Public HTML routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			serveTemplate(w, r, "index.html")
			return
		}
		utils.RenderErrorPage(w, http.StatusNotFound)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.RenderErrorPage(w, http.StatusMethodNotAllowed)
			return
		}
		serveTemplate(w, r, "login.html")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		serveTemplate(w, r, "register.html")
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		serveTemplate(w, r, "post.html")
	})

	http.HandleFunc("/comment", func(w http.ResponseWriter, r *http.Request) {
		serveTemplate(w, r, "comment.html")
	})

	// API routes - Public
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.HandleFunc("/api/posts", handlers.GetPostsHandler) // Public GET endpoint
	http.HandleFunc("/api/categories", handlers.GetCategoriesHandler)
	// Protected API routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/logout", handlers.LogoutHandler)
	protectedMux.HandleFunc("/api/posts", handlers.GetPostsHandler)
	protectedMux.HandleFunc("/api/posts/create", handlers.CreatePostHandler)
	protectedMux.HandleFunc("/api/comments", handlers.CreateCommentHandler)
	protectedMux.HandleFunc("/api/likes", handlers.CreateLikeDislikeHandler)
	protectedMux.HandleFunc("/api/auth/status", handlers.AuthStatusHandler)
	protectedMux.HandleFunc("/api/comments/like", handlers.CreateCommentLikeHandler)

	http.Handle("/api/protected/", middleware.AuthMiddleware(http.StripPrefix("/api/protected", protectedMux)))

	// Start server
	log.Println("Server started on http://localhost:9111")
	log.Fatal(http.ListenAndServe(":9111", nil))
}
