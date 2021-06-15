package server

import (
	"net/http"
)

func (s *server) PostsCreate() http.HandlerFunc {

	type Request struct{}
	type Response struct{}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *server) PostsGet() http.HandlerFunc {

	type Response struct{}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *server) PostsGetAll() http.HandlerFunc {

	type Response struct{}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *server) PostsDelete() http.HandlerFunc {

	type Response struct{}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}
