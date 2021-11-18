package internal

import (
	"coffee-log/db/sqlc"
	"coffee-log/util"
	"context"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

	server = NewServer(db, true)

	os.Exit(m.Run())
}

func TestRootPath(t *testing.T) {
	ctx := context.Background()
	store := sqlc.NewStore(db)

	user, err := store.CreateUser(ctx, sqlc.CreateUserParams{
		DisplayName: "Test",
		Username:    util.RandomUsername(),
		Password:    util.RandomPassword(),
	})
	require.NoError(t, err)
	t.Logf("created test user: %+v", user)

	err = store.DeleteAllLoginAttemptsAndBans(ctx)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	req.SetBasicAuth(user.Username, user.Password)

	server.ServeHTTP(w, req)

	t.Logf("Response: %+v", w.Result())

	require.Equal(t, http.StatusFound, w.Code)
}
