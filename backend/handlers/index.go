package handlers

import (
	"net/http"
	"os"

	"forum/utils"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// utils.RenderErrorPage(w, http.StatusMethodNotAllowed)
		utils.ErrorMessage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/" {
		serveTemplate(w, r, "index.html")
		return
	}
	utils.ErrorMessage(w, "404 Not Found", http.StatusNotFound)
	// utils.RenderErrorPage(w, http.StatusNotFound)
}

func Static(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
		// utils.RenderErrorPage(w, http.StatusMethodNotAllowed)
		return
	}
	path := "../frontend" + r.URL.Path
	f, err := os.Stat(path)
	if err != nil {
	utils.ErrorMessage(w, "404 Not Found", http.StatusNotFound)
		// utils.RenderErrorPage(w, http.StatusNotFound)
		return
	}
	if f.IsDir() {
		utils.ErrorMessage(w, "404 Not Found", http.StatusNotFound)
		// utils.RenderErrorPage(w, http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
}
