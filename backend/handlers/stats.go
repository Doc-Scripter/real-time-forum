package handlers

import (
	"encoding/json"
	"net/http"

	"forum/database/query"
)

func GetForumStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := query.GetForumStats()
	if err != nil {
		errorMessage(w, "Failed to fetch forum stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
