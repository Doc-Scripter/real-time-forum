package handlers

import (
	"encoding/json"
	"forum/models"
	"net/http"
)

func CurrentUserHandler(w http.ResponseWriter,r *http.Request){
	user:=models.CurrentUser
	json.NewEncoder(w).Encode(user)
}