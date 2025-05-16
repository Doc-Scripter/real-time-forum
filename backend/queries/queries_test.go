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
		{name: "test", user: models.User{Nickname: "john", Email: "john@gmail.com", Password: "mpass"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := CreateUser(tt.user); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name    string
		post    models.Post
		wantErr bool
	}{
		{name: "create post", post: models.Post{Title: "test title", Content: "test content", RawCategories: []int{1, 2}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreatePost(tt.post); (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCategories(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "category", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCategories()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAllPosts(t *testing.T) {
	type args struct {
		offset   int
		category string
		filter   string
		userID   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "get all posts",
			args:    args{offset: 8, category: "1", filter: "", userID: 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAllPosts(tt.args.offset, tt.args.category, tt.args.filter, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCategoriesByPostID(t *testing.T) {
	tests := []struct {
		name    string
		postID  int
		wantErr bool
	}{
		{name: "post", postID: 1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCategoriesByPostID(tt.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategoriesByPostID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetForumStats(t *testing.T) {
	tests := []struct {
		name    string
		notNil  bool
		wantErr bool
	}{
		{name: "test", notNil: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetForumStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetForumStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.notNil {
				t.Errorf("GetForumStats() = %v, want %v", got, tt.notNil)
			}
		})
	}
}
