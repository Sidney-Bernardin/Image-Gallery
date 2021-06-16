package server

import (
	"io/ioutil"
	"net/http"
	"root/db"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) PostsCreate() http.HandlerFunc {

	type Response struct {
		Success bool `json:"success"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Verify that the title and description are in the request.
		required := []string{"title", "description"}
		for _, v := range required {
			if r.FormValue(v) == "" || r.FormValue(v) == " " {
				err := errors.New(v + " is required")
				s.err(w, err, http.StatusUnprocessableEntity)
				return
			}
		}

		// Setup the post.
		post := &db.Post{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
		}

		// Parse the multipart/form-data input.
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			s.err(w, err, http.StatusUnprocessableEntity)
			return
		}

		// Retrieve the file from the form-data.
		file, _, err := r.FormFile("thumbnail")
		if err != nil {
			s.err(w, err, http.StatusUnprocessableEntity)
			return
		}
		defer file.Close()

		// Read the file into bytes.
		b, err := ioutil.ReadAll(file)
		if err != nil {
			s.err(w, err, http.StatusInternalServerError)
			return
		}

		// Create the post.
		if err := s.db.PostsCreate(post, b); err != nil {
			s.err(w, err, http.StatusInternalServerError)
			return
		}

		// Respond with success.
		s.respond(w, http.StatusCreated, Response{true})
	}
}

func (s *server) PostsGet() http.HandlerFunc {

	type Response struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnail   string `json:"thumbnail"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Get the postID.
		postID := mux.Vars(r)["postID"]

		// Get the post.
		post, err := s.db.PostsGet(postID)
		if err != nil {
			s.err(w, err, http.StatusInternalServerError)
			return
		}

		// Create the response.
		response := Response{
			post.ID.Hex(),
			post.Title,
			post.Description,
			"/posts/thumbnailfor" + post.Title,
		}

		// Respond.
		s.respond(w, http.StatusOK, response)
	}
}

func (s *server) PostsGetAll() http.HandlerFunc {

	type Post struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnail   string `json:"thumbnail"`
	}

	type Response struct {
		Posts []Post `json:"posts"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Get the offset and the limit.
		offset := mux.Vars(r)["offset"]
		limit := mux.Vars(r)["limit"]

		// Convert offset into an int64.
		offset64, err := strconv.ParseInt(offset, 0, 64)
		if err != nil {
			msg := "offset must be an int64"
			err = errors.Wrap(err, msg)
			s.err(w, err, http.StatusUnprocessableEntity)
			return
		}

		// Convert limit into an int64.
		limit64, err := strconv.ParseInt(limit, 0, 64)
		if err != nil {
			msg := "limit must be an int64"
			err = errors.Wrap(err, msg)
			s.err(w, err, http.StatusUnprocessableEntity)
			return
		}

		// Get all of the posts.
		posts, err := s.db.PostsGetAll(offset64, limit64)
		if err != nil {
			s.err(w, err, http.StatusInternalServerError)
			return
		}

		// Copy each post into the response array and give each of them a thumbnail field.
		var response Response
		for _, v := range *posts {
			response.Posts = append(response.Posts, Post{
				v.ID.Hex(),
				v.Title,
				v.Description,
				"/thumbnailfor" + v.Title,
			})
		}

		// Respond.
		s.respond(w, http.StatusOK, response)
	}
}

func (s *server) PostsDelete() http.HandlerFunc {

	type Response struct{}

	return func(w http.ResponseWriter, r *http.Request) {

		// Get the postID.
		postID := mux.Vars(r)["postID"]

		// Delete the post.
		if err := s.db.PostsDelete(postID); err != nil {
			s.err(w, err, http.StatusInternalServerError)
			return
		}

		// Respond.
		s.respond(w, http.StatusNoContent, nil)
	}
}
