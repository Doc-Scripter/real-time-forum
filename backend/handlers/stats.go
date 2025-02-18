package handlers

import (
	"encoding/json"
	"net/http"

	"forum/queries"
	"forum/utils"
)

func GetForumStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := queries.GetForumStats()
	if err != nil {
		utils.ErrorMessage(w, "Ooops! Our post categories don't work! Try again later...", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
