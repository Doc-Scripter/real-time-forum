package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"forum/database"
	"forum/middleware"
	"forum/queries"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch the user by email
	user, err := queries.GetUserByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Delete any existing sessions for this user
	err = queries.DeleteUserSessions(user.ID)
	if err != nil {
		http.Error(w, "Failed to manage sessions", http.StatusInternalServerError)
		return
	}

	// Generate a session token
	sessionToken := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	// Store the session token in the database
	_, err = database.DB.Exec("INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		sessionToken, user.ID, expiresAt)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set the session token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Login successful",
		"username": user.Username,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete all sessions for this user
	err := queries.DeleteUserSessions(userID)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Check if username exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", user.Username).Scan(&exists)
	if err != nil {
		http.Error(w, "Failed to check username", http.StatusInternalServerError)
		return
	}
	if exists {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Username already taken"})
		return
	}

	// Check if email exists
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		http.Error(w, "Failed to check email", http.StatusInternalServerError)
		return
	}
	if exists {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	_, err = database.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		user.Username, user.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
