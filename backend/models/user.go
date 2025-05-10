package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Nickname  string    `json:"nickname"`
    Age       int       `json:"age"`
    Gender    string    `json:"gender"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    Email     string    `json:"email"`
    Password  string    `json:"password,omitempty"`
    CreatedAt time.Time `json:"createdAt"`
}
var CurrentUser int
