package handlers

import (
	"encoding/json"
	"net/http"

	"forum/queries"
	"forum/utils"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := queries.GetCategories()
	if err != nil {
		utils.RenderErrorPage(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
