package handlers

import (
	"encoding/json"
	"errors"
	"io"
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

// Get JSON body content
func ParseJSONBody(r io.Reader, model any) error {
	return json.NewDecoder(r).Decode(model)
}

// Render error page with a custom error message
func serveTemplate(w http.ResponseWriter, r *http.Request, templatePath string) {
	path := filepath.Join("../frontend/templates", templatePath)
	http.ServeFile(w, r, path)
}
