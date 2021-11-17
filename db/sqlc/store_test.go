package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	db *sql.DB
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	dbUrl := os.Getenv("TEST_DATABASE_URL")

	var openErr error
	if db, openErr = sql.Open("postgres", dbUrl); openErr != nil {
		log.Fatalf("error connecting to database: %+v", openErr)
	}
	defer func(db *sql.DB) {
		if closeErr := db.Close(); closeErr != nil {
			log.Fatalf("error closing database conn: %+v", closeErr)
		}
	}(db)

	rand.Seed(time.Now().UnixMilli())

	os.Exit(m.Run())
}

func TestCheckAndLogLoginAttempt_ValidCredentials(t *testing.T) {
	Rollback(t, db, func(ctx context.Context, store *Store) {
		user, err := store.CreateUser(ctx, CreateUserParams{
			DisplayName: "Tester",
			Username:    fmt.Sprintf("%d", rand.Int31()),
			Password:    "password",
			TimeZone:    sql.NullString{},
		})
		if err != nil {
			t.Fatalf("error creating test user: %+v", err)
		}

		ip := fmt.Sprintf("%d", rand.Int31())
		loggedInUser, err := store.CheckAndLogLoginAttempt(ctx, ip, user.Username, user.Password, 5)
		if err != nil {
			t.Fatalf("unexpected error testing correct username and password: %+v", err)
		}

		if loggedInUser != user {
			t.Errorf("expected %v to be equal to %v", loggedInUser, user)
		}
	})
}


func TestCheckAndLogLoginAttempt_BadCredentials(t *testing.T) {
	Rollback(t, db, func(ctx context.Context, store *Store) {
		ip := fmt.Sprintf("%d", rand.Int31())
		username := "fubar"
		password := ""
		maxAttempts := int32(10)

		var err error

		for i := int32(0); i < maxAttempts-1; i++ {
			_, err = store.CheckAndLogLoginAttempt(ctx, ip, username, password, maxAttempts)
			if err != ErrBadCredentials {
				t.Fatalf("expected ErrBadCredentials, but got: %+v", err)
			}
		}

		_, err = store.CheckAndLogLoginAttempt(ctx, ip, username, password, maxAttempts)
		if err != ErrIPBanned {
			t.Errorf("expected ErrIPBanned, but got %+v", err)
		}
	})
}
