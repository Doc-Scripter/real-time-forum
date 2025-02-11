package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	"forum/models"
)

func ValidatePost(p models.Post) error {
	if p.Title == "" || p.Content == "" {
		return errors.New("title and content are required")
	}
	return nil
}

func errorMessage(w http.ResponseWriter, msg string, errorCode int) {
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": msg,
	})
}

func serveTemplate(w http.ResponseWriter, r *http.Request, templatePath string) {
	path := filepath.Join("../frontend/templates", templatePath)
	http.ServeFile(w, r, path)
}
