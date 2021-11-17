package internal

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var server *Server

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	dbUrl := os.Getenv("TEST_DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	portInt, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		log.Fatalf("error parsing port number: %+v", err)
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("error connecting to database: %+v", err)
	}
	defer func(db *sql.DB) {
		closeErr := db.Close()
		if closeErr != nil {
			log.Fatalf("error closing database conn: %+v", closeErr)
		}
	}(db)

	rand.Seed(time.Now().UnixMilli())

	server = NewServer("0.0.0.0", int32(portInt), db)
	if err = server.Run(); err != nil {
		log.Fatalf("error running server: %+v", err)
	}

	os.Exit(m.Run())
}

func TestLogsIndex(t *testing.T) {

}
