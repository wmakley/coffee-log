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

	server = NewServer(db, false)

	os.Exit(m.Run())
}

func TestRootPathRedirectsToLogEntries(t *testing.T) {
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
	req := util.NewTestRequest("GET", "/", nil, user.BasicCredentials())

	server.ServeHTTP(w, req)
	res := w.Result()

	util.AssertRedirectedTo(t,"/logs/" + user.Username + "/entries/", 302, res)

	// new request to follow redirect
	req = util.FollowRedirect(t, res, user.BasicCredentials())
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	res = w.Result()
	_ = util.ReadAndLogBody(t, res)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, "text/html", res.Header.Get("content-type"))
}
