package main

import (
	"flag"
	"net/http"
	"root/db"
	"root/server"

	"github.com/sirupsen/logrus"
)

func main() {

	// Setup the flags.
	dbURL := flag.String("dburl", "mongodb://0.0.0.0:27017/", "")
	timeout := flag.Int("dbtimeout", 9, "")

	// Create the db.
	db, err := db.NewDB(*dbURL, *timeout)
	if err != nil {
		logrus.Fatalf("cannot create db: %v", err)
	}

	// Create the server.
	s := server.NewServer(db)

	// Start the server.
	logrus.Info("Listening on :8080...")
	http.ListenAndServe(":8080", s)
}
