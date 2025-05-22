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
		utils.ErrorMessage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
