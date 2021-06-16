package db

type MockDB struct{}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func (m *MockDB) PostsCreate(post *Post, b []byte) error           { return nil }
func (m *MockDB) PostsGet(postID string) (*Post, error)            { return nil, nil }
func (m *MockDB) PostsGetAll(offset, limit int64) (*[]Post, error) { return nil, nil }
func (m *MockDB) PostsDelete(postID string) error                  { return nil }
