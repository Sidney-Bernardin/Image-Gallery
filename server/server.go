package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"root/db"
)

type server struct {
	router *mux.Router
	db     db.DB
	logger *logrus.Logger
}

func NewServer(db db.DB) *server {
	s := &server{mux.NewRouter(), db, logrus.New()}
	s.loadRoutes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
