package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"forum/logging"
	"forum/models"
	"forum/utils"
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
	_, err := os.Stat(path)
	if err != nil {
		logging.Log("[ERROR] :Error serving template: %v", err)
		w.WriteHeader(http.StatusNotFound)
		utils.ErrorMessage(w, "Template not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
}
