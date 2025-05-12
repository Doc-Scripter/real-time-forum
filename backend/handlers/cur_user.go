package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum/middleware"
)

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	fmt.Println("userID: ", userID)
	json.NewEncoder(w).Encode(userID)
}
