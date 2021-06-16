package db

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	PostsCreate(*Post, []byte) error
	PostsGet(string) (*Post, error)
	PostsGetAll(int64, int64) (*[]Post, error)
	PostsDelete(string) error
}

type db struct {
	client  *mongo.Client
	timeout int
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
	return &db{client, timeout}, err
}

func (db *db) PostsCreate(post *Post, b []byte) error {

	const operation = "db.PostsCreate"

	// Create context.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(db.timeout)*time.Second)
	defer cancel()

	// Insert the post.
	_, err := db.client.Database("gallery").Collection("posts").InsertOne(ctx, post)
	if err != nil {
		return errors.Wrap(err, operation)
	}

	// Create a new bucket in the posts database.
	bucket, err := gridfs.NewBucket(db.client.Database("gallery"))
	if err != nil {
		return errors.Wrap(err, operation)
	}

	// Open an upload stream.
	uploadStream, err := bucket.OpenUploadStream("thumbnailfor" + post.Title)
	if err != nil {
		return errors.Wrap(err, operation)
	}
	defer uploadStream.Close()

	// Write to the stream.
	_, err = uploadStream.Write(b)
	return errors.Wrap(err, operation)
}

func (db *db) PostsGet(postID string) (*Post, error) {

	const operation = "db.PostsGet"

	// Get the database and the posts collection.
	posts := db.client.Database("gallery").Collection("posts")

	// Create context.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(db.timeout)*time.Second)
	defer cancel()

	// Generate an objectID.
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, errors.Wrap(err, operation)
	}

	// Find a post with maching ids and return it.
	var post Post
	err = posts.FindOne(ctx, bson.M{"_id": objID}).Decode(&post)
	return &post, errors.Wrap(err, operation)
}

func (db *db) PostsGetAll(offset, limit int64) (*[]Post, error) {

	const operation = "db.PostsGetAll"

	// Create context.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(db.timeout)*time.Second)
	defer cancel()

	// Find all of the posts in the posts collection within the offset and limit.
	cursor, err := db.client.Database("gallery").Collection("posts").Find(
		ctx,
		bson.M{},
		options.Find().
			SetSkip(offset).
			SetLimit(limit))

	if err != nil {
		return nil, errors.Wrap(err, operation)
	}

	// Put the posts into a slice and return it.
	var posts []Post
	err = cursor.All(ctx, &posts)
	return &posts, errors.Wrap(err, operation)
}

func (db *db) PostsDelete(postID string) error {

	const operation = "db.PostsDelete"

	// Create context.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(db.timeout)*time.Second)
	defer cancel()

	// Get the posts and the fs.files collection.
	posts := db.client.Database("gallery").Collection("posts")
	fsFiles := db.client.Database("gallery").Collection("fs.files")

	// Generate an objectID.
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return errors.Wrap(err, operation)
	}

	// Find a post with the same id.
	var post Post
	if err := posts.FindOne(ctx, bson.M{"_id": objID}).Decode(post); err != nil {
		return errors.Wrap(err, operation)
	}

	// Delete the post.
	if _, err := posts.DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
		return errors.Wrap(err, operation)
	}

	// Delete the thumbnail for the post.
	_, err = fsFiles.DeleteOne(ctx, bson.M{"name": "thumbnailfor" + post.Title})
	return errors.Wrap(err, operation)
}
