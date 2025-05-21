package handlers

import (
	"net/http"
	"os"

	"forum/utils"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorMessage(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/" {
		serveTemplate(w, r, "index.html")
		return
	}
	w.Header().Set("X-Status-Code", "404")
	serveTemplate(w, r, "index.html")


}

func Static(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorMessage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := "../frontend" + r.URL.Path
	f, err := os.Stat(path)
	if err != nil {
		utils.ErrorMessage(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if f.IsDir() {
		utils.ErrorMessage(w, "404 Not Found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
}
