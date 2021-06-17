package db

import "go.mongodb.org/mongo-driver/mongo"

type MockDB struct{}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func (m *MockDB) PostsCreate(post *Post, b []byte) error {
	return nil
}

func (m *MockDB) PostsGet(postID string) (*Post, error) {
	if postID == "exists" {
		return &Post{
			Title:       "exists",
			Description: "exists",
		}, nil
	}

	return nil, mongo.ErrNoDocuments
}

func (m *MockDB) PostsGetAll(offset, limit int64) (*[]Post, error) {
	return &[]Post{}, nil
}

func (m *MockDB) PostsDelete(postID string) error {
	return nil
}
