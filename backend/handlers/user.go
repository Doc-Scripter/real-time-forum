package handlers

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/logging"
	"forum/models"
	"forum/queries"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

// LoginHandler authenticates a user and creates a new session
// Accepts POST requests with email/username and password
// Returns a session token on successful login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}

	var credentials struct {
		LoginIdentifier string `json:"loginIdentifier"`
		Password        string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		logging.Log("[ERROR] Failed to decode login credentials: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	var user models.User
	err = database.DB.QueryRow(`
		SELECT id,  nickname, age, gender, first_name, last_name, email, password 
		FROM users 
		WHERE email = ? OR nickname = ?`,
		credentials.LoginIdentifier, credentials.LoginIdentifier).
		Scan(&user.ID, &user.Nickname, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.Password)

	if err != nil {
		logging.Log("[WARNING] User not found: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid email/username or password",
		})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		logging.Log("[WARNING] Invalid password for user: %s", user.Nickname)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid email/username or password",
		})
		return
	}

	err = queries.DeleteUserSessions(user.ID)
	if err != nil {
		logging.Log("[ERROR] Failed to delete user sessions: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to manage sessions",
		})
		return
	}

	sessionToken := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = database.DB.Exec("INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		sessionToken, user.ID, expiresAt)
	if err != nil {
		logging.Log("[ERROR] Failed to create session: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to create session",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"user": map[string]interface{}{
			"nickname": user.Nickname,
			"email":    user.Email,
		},
	})
}

// LogoutHandler ends a user's session and clears their session cookie
// Accepts POST requests only
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		logging.Log("[WARNING] Failed to get session token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var userID int
	err = database.DB.QueryRow(`
        SELECT user_id FROM sessions 
        WHERE token = ?`, cookie.Value).Scan(&userID)
	if err != nil {
		logging.Log("[ERROR] Failed to get user ID from session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = queries.DeleteUserSessions(userID)
	if err != nil {
		logging.Log("[ERROR] Failed to delete user sessions: %v", err)
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

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

// RegisterHandler creates a new user account
// Accepts POST requests with user registration data
// Returns the created user on success
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logging.Log("[ERROR] Failed to decode registration request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	if user.Nickname == "" || user.FirstName == "" || user.LastName == "" ||
		user.Email == "" || user.Password == "" || user.Age < 13 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "All fields are required and age must be at least 13",
		})
		return
	}

	if len(user.Nickname) > 32 || len(user.FirstName) > 32 || len(user.LastName) > 32 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Nickname, first name, and last name must be 32 characters or less",
		})
		return
	}

	if !isValidEmail(user.Email) {
		log.Printf("[ERROR] Invalid email format: %s", user.Email)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid email format",
		})
		return
	}

	var exists bool
	err = database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE nickname = ? OR email = ?)
	`, user.Nickname, user.Email).Scan(&exists)
	if err != nil {
		logging.Log("[ERROR] Database error checking username/email: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Database error checking username/email: %v", err),
		})
		return
	}
	if exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Username or email already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.Log("[ERROR] Error hashing password: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Error hashing password: %v", err),
		})
		return
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	userID, err := queries.CreateUser(user)
	if err != nil {
		logging.Log("[ERROR] Failed to create user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Failed to create user: %v", err),
		})
		return
	}

	user.ID = userID
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
		"message": "Registration successful",
	})
}

// isValidEmail validates if the given string is a valid email address
// Returns true if the email contains @ and . characters
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	return true
}
