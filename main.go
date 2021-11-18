package main

import (
	"coffee-log/internal"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8080"
	}
	portInt64, err := strconv.ParseInt(portString, 10, 32)
	if err != nil {
		log.Fatalf("error parsing port number: %+v", err)
	}
	port := int32(portInt64)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %+v", err)
	}
	defer func(db *sql.DB) {
		closeErr := db.Close()
		if closeErr != nil {
			log.Fatalf("error closing db conn: %+v", closeErr)
		}
	}(db)

	server := internal.NewServer(db, false)

	err = server.Run("",  port)
	if err != nil {
		log.Fatalf("run error: %+v", err)
	}
}
