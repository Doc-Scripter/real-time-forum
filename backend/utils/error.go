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

// func RenderErrorPage(w http.ResponseWriter, status int) {
// 	errorPage, exists := ErrorMessages[status]
// 	if !exists {
// 		errorPage = ErrorPage{
// 			Code:        status,
// 			Title:       "Unexpected Error",
// 			Description: "An unexpected error occurred. Please try again later.",
// 		}
// 	}

// 	tmpl, err := template.ParseFiles(filepath.Join("../frontend/templates", "error.html"))
// 	if err != nil {
// 		http.Error(w, "Error template not found", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(status)
// 	tmpl.Execute(w, errorPage)
// }

func ErrorMessage(w http.ResponseWriter, msg string, errorCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": msg,
	})
}
