package handlers

import (
	"encoding/json"
	"net/http"

	"forum/logging"
	"forum/queries"
	"forum/utils"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	categories, err := queries.GetCategories()
	if err != nil {
		logging.Log("[ERROR] Error getting categories: %v", err)
		utils.ErrorMessage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
