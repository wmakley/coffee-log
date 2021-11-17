package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

var ErrRollback = errors.New("rollback")

// Rollback wraps fn in a transaction that is always rolled back. useful for unit testing.
func Rollback(t *testing.T, db *sql.DB, fn func (context.Context, *Store)) {
	ctx := context.Background()
	store := NewStore(db)
	outerErr := store.transaction(ctx, func(store *Store) error {
		fn(ctx, store)
		return ErrRollback
	})
	if outerErr != ErrRollback {
		t.Fatalf("unexpected error rolling back transaction: %+v", outerErr)
	}
}
