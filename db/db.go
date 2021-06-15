package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	PostsCreate(*Post) error
	PostsGet(string) (*Post, error)
	PostsGetAll(string, string) (*[]*Post, error)
	PostsDelete(string) error
}

type db struct {
	client *mongo.Client
}

func NewDB(url string, timeout int) (DB, error) {

	// Create client.
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	// Create context.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeout)*time.Second,
	)
	defer cancel()

	// Connent the client and return.
	err = client.Connect(ctx)
	return &db{client}, err
}

func (db *db) PostsCreate(post *Post) error {
	return nil
}

func (db *db) PostsGet(postID string) (*Post, error) {
	return nil, nil
}

func (db *db) PostsGetAll(offset, limit string) (*[]*Post, error) {
	return nil, nil
}

func (db *db) PostsDelete(postID string) error {
	return nil
}
