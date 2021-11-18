package internal

import (
	"coffee-log/db/sqlc"
	"context"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	db     *sql.DB
	server *Server
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	testDatabaseURL := os.Getenv("TEST_DATABASE_URL")

	var err error
	db, err = sql.Open("postgres", testDatabaseURL)
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

	server = NewServer(db)

	os.Exit(m.Run())
}

func TestRootPath(t *testing.T) {
	sqlc.Rollback(t, db, func(ctx context.Context, store *sqlc.Store) {
		_, err := store.CreateUser(ctx, sqlc.CreateUserParams{
			DisplayName: "Test",
			Username:    "test",
			Password:    "test",
		})
		require.NoError(t, err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.SetBasicAuth("test", "test")

		server.ServeHTTP(w, req)

		t.Log("Response", w.Result())

		require.Equal(t, http.StatusFound, w.Code)
	})
}
