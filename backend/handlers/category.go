package handlers

import (
	"encoding/json"
	"net/http"

	"forum/database"
	"forum/utils"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := database.GetCategories()
	if err != nil {
		utils.RenderErrorPage(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
