package main

import (
	"net/http"
	"os"
	"root/db"
	"root/server"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {

	// Setup the flags.
	var err error
	var timeout int

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	dbURL, ok := os.LookupEnv("DB_URL")
	if !ok {
		logrus.Fatalf("must have the env-var of DB_URL")
	}

	strTimeout, ok := os.LookupEnv("DB_TIMEOUT")
	if !ok {
		timeout = 9
	} else {
		timeout, err = strconv.Atoi(strTimeout)
		if err != nil {
			logrus.Fatalf("DB_TIMEOUT must be an int")
		}
	}

	// Create the db.
	db, err := db.NewDB(dbURL, timeout)
	if err != nil {
		logrus.Fatalf("cannot create db: %v", err)
	}

	// Create the server.
	s := server.NewServer(db)

	// Start the server.
	logrus.Infof("Listening on :%s...", port)
	http.ListenAndServe(":"+port, s)
}
