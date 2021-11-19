package sqlc

import (
	"coffee-log/util"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var ErrRollback = errors.New("rollback")

// Rollback wraps fn in a transaction that is always rolled back. useful for unit testing.
func Rollback(t *testing.T, db *sql.DB, fn func(context.Context, *Store)) {
	ctx := context.Background()
	store := NewStore(db)
	//store.Debug = true
	_ = store.transaction(ctx, func(store *Store) error {
		fn(context.WithValue(ctx, "tx", store.tx), store)
		return ErrRollback
	})
	//t.Log("rolled back transaction")
}

func RandomIP() string {
	ip := rand.Int31()
	return fmt.Sprintf("%d", ip)
}

func RandomUser() CreateUserParams {
	return CreateUserParams{
		DisplayName: "Test Testerson",
		Username:    util.RandomUsername(),
		Password:    util.RandomPassword(),
		TimeZone:    sql.NullString{
			String: "EST",
			Valid: true,
		},
	}
}

func ValidLogEntry(logID int64) CreateLogEntryParams {
	return CreateLogEntryParams{
		LogID: logID,
		EntryDate: time.Now(),
		Coffee: util.RandomString(8),
	}
}
