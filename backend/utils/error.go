package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorPage struct {
	Code        int
	Title       string
	Description string
}

var ErrorMessages = map[int]ErrorPage{
	400: {Code: 400, Title: "Bad Request", Description: "The server cannot process your request due to invalid syntax."},
	401: {Code: 401, Title: "Unauthorized", Description: "You need to be logged in to access this resource."},
	403: {Code: 403, Title: "Forbidden", Description: "You don't have permission to access this resource."},
	404: {Code: 404, Title: "Page Not Found", Description: "The page you're looking for doesn't exist or has been moved."},
	500: {Code: 500, Title: "Internal Server Error", Description: "Something went wrong on our end. Please try again later."},
}


func ErrorMessage(w http.ResponseWriter, msg string, errorCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": msg,
	})
}
