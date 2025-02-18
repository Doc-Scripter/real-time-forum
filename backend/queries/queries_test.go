package queries

import (
	"testing"

	"forum/models"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{name: "test", user: models.User{Username: "john", Email: "john@gmail.com", Password: "mpass"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateUser(tt.user); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
