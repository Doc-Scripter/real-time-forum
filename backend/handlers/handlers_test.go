package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/middleware"
)

func TestRegisterHandler(t *testing.T) {
	// mock response
	res := httptest.NewRecorder()

	// mock registration request request
	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "testpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/register", &buffer)
	res.Code = http.StatusOK

	RegisterHandler(res, req)
	if res.Code != http.StatusCreated {
		t.Errorf("RegisterHandler() got = %v, want %v", res.Code, http.StatusCreated)
	}
}

func TestLoginHandler(t *testing.T) {
	// mock response
	res := httptest.NewRecorder()

	// mock registration request request
	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "testpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/login", &buffer)
	res.Code = http.StatusOK

	LoginHandler(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("LoginHandler() got = %v, want %v", res.Code, http.StatusOK)
	}
	if res.Result().Cookies()[0].Value == "" {
		t.Error("No session token in response")
	}
}

func TestCreatePostHandler(t *testing.T) {
	// Mock response
	res := httptest.NewRecorder()

	// Create a test post
	var post bytes.Buffer
	json.NewEncoder(&post).Encode(map[string]string{
		"title":   "Test post",
		"content": "This is a test post",
	})

	// Create a new request
	req := httptest.NewRequest(http.MethodPost, "/api/protected/api/posts/create", &post)

	// Set up authentication context
	// This assumes you're using a middleware to check for a session token
	// You may need to adjust this based on your actual authentication mechanism
	cookie := &http.Cookie{
		Name:  "session_token",
		Value: "test_session_token",
	}
	req.AddCookie(cookie)

	// You might also need to set up the user context
	// This is an example and may need to be adjusted based on your implementation
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)

	// Call the handler
	CreatePostHandler(res, req)

	// Check the response
	if res.Code != http.StatusCreated {
		t.Errorf("CreatePostHandler() got = %v, want %v", res.Code, http.StatusCreated)
	}

	// Optionally, you can check the response body to ensure the post was created correctly
	var responseBody map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBody)
}

func TestGetPostsHandler(t *testing.T) {
	// test response recorder
	res := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)

	GetPostsHandler(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("CreatePostHandler() got = %v, want %v", res.Code, http.StatusOK)
	}
}

func TestLogoutHandler(t *testing.T) {
	// mock response
	res := httptest.NewRecorder()

	// mock registration request request
	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "testpassword",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/register", &buffer)
	res.Code = http.StatusOK

	RegisterHandler(res, req)
	req = httptest.NewRequest(http.MethodGet, "/api/login", &buffer)
	res.Code = http.StatusOK

	LoginHandler(res, req)
	req = httptest.NewRequest(http.MethodGet, "/api/protected/api/logout", &buffer)
	res.Code = http.StatusOK

	LogoutHandler(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("LoginHandler() got = %v, want %v", res.Code, http.StatusOK)
	}
	if len(res.Result().Cookies()) != 0 {
		t.Error("Session token still present in response")
	}
}
