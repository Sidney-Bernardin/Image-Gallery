package db

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func (m *MockDB) PostsThumbnailGet(postID string) (*bytes.Buffer, error) {

	const operation = "MockDB.PostsThumbnailGet"

	if postID == "dosenotexist" {
		return nil, mongo.ErrNoDocuments
	}

	f, err := os.Open("../for_unit_tests.png")
	if err != nil {
		return nil, errors.Wrap(err, operation)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	buffer := bytes.NewBuffer(make([]byte, 0))
	part := make([]byte, 1024)

	var count int

	for {
		if count, err = reader.Read(part); err != nil {
			break

		}
		buffer.Write(part[:count])

	}
	if err != io.EOF {
		return nil, errors.New(operation + ": " + err.Error())
	}

	return buffer, nil
}
