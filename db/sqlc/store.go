package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Store struct {
	*Queries
	db *sql.DB
	tx *sql.Tx
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
		tx:      nil,
	}
}

func StoreWithTx(tx *sql.Tx) *Store {
	return &Store{
		Queries: New(tx),
		db:      nil,
		tx:      tx,
	}
}

// StoreFromCtx get a data store using the context transaction, or create a new one
// using the context database connection. Panics if either of those cannot be
// found in the context.
func StoreFromCtx(ctx *gin.Context) *Store {
	maybeDb, dbSet := ctx.Get("db")
	if !dbSet {
		panic(errors.New("db is not set"))
	}
	db, ok := maybeDb.(*sql.DB)
	if !ok {
		panic(errors.New("db is not *sql.DB"))
	}

	maybeTx, txSet := ctx.Get("tx")
	if !txSet {
		return NewStore(db)
	}

	tx, ok := maybeTx.(*sql.Tx)
	if !ok {
		panic(fmt.Errorf("ctx key 'tx' cannot be converted to *sql.Tx"))
	}

	return StoreWithTx(tx)
}

// ProvideDatabase provides the database connection to any Store instances
// created using StoreFromCtx
//func ProvideDatabase(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.Set("db", db)
//		c.Next()
//	}
//}

func WrapInTransaction(db *sql.DB, options *sql.TxOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)

		tx, err := db.BeginTx(c, options)
		if err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Set("tx", tx)

		c.Next()

		if len(c.Errors) > 0 {
			log.Print("rolling back due to errors")
			if rbErr := tx.Rollback(); rbErr != nil {
				c.Error(rbErr)
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("error rolling back due to prior errors: %v: %v", c.Errors.String(), rbErr))
			}
			return
		}

		if commitErr := tx.Commit(); err != nil {
			c.Error(commitErr)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error committing: %v", commitErr))
		}
	}
}

func (store *Store) transaction(ctx context.Context, fn func(*Queries) error) error {
	if store.tx != nil {
		// we are already in a transaction, some other process
		// will handle the rollback
		return fn(New(store.tx))
	}

	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	store.tx = tx

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction err: %+v, rollback error: %+v", err, rbErr)
		}
		return err
	}

	commitErr := tx.Commit()
	store.tx = nil
	return commitErr
}

var ErrIPBanned = errors.New("ip address is banned")
var ErrBadCredentials = errors.New("invalid credentials")

func (store *Store) CheckAndLogLoginAttempt(
	ctx context.Context,
	ipAddress string,
	username string,
	password string,
	maxAttempts int32,
) (User, error) {
	var user User

	_, err := store.GetBannedIP(ctx, ipAddress)
	if err != nil {
		// ignore not found
		if err != sql.ErrNoRows {
			return user, err
		}
	} else {
		// banned if record exists
		return user, ErrIPBanned
	}

	foundUser := true
	user, err = store.GetUserByUsername(ctx, username)
	if err != nil {
		foundUser = false
		if err != sql.ErrNoRows {
			return user, err
		}
	}

	if foundUser && user.Password == password {
		// success!
		return user, nil
	}

	banned, err := store.CreateOrIncrementLoginAttempt(ctx, ipAddress, maxAttempts)
	if err != nil {
		return user, err
	}
	if banned {
		return user, ErrIPBanned
	}

	return user, ErrBadCredentials
}

func (store *Store) CreateOrIncrementLoginAttempt(ctx context.Context, ip string, maxAttempts int32) (banned bool, err error) {
	var attempts LoginAttempt

	attempts, err = store.GetLoginAttempt(ctx, ip)
	if err != nil {
		if err != sql.ErrNoRows {
			return true, nil
		}

		attempts, err = store.CreateLoginAttempt(ctx, ip)
		if err != nil {
			return
		}
	} else {
		attempts, err = store.IncrementLoginAttempt(ctx, ip)
		if err != nil {
			return
		}
	}

	if maxAttempts > 0 && attempts.Attempts >= maxAttempts {
		banned = true
		_, err = store.CreateBannedIP(ctx, ip)
		if err != nil {
			return
		}
	}

	return
}

func (store *Store) GetLogAndEntriesBySlugOrderByDateDesc(
	ctx context.Context, slug string,
) (Log, []LogEntry, error) {
	var log2 Log
	var entries []LogEntry

	err := store.transaction(ctx, func(q *Queries) error {
		var err error

		log2, err = q.GetLogBySlug(ctx, slug)
		if err != nil {
			return err
		}

		entries, err = q.ListLogEntriesByLogIDOrderByDateDesc(ctx, log2.ID)
		if err != nil {
			if err != sql.ErrNoRows {
				return err
			}
			// return empty slice
		}

		return nil
	})

	return log2, entries, err
}

func (store *Store) CreateLogEntry(
	ctx context.Context, logSlug string, params CreateLogEntryParams,
) (Log, LogEntry, error) {
	var log_ Log
	var logEntry LogEntry

	txErr := store.transaction(ctx, func(queries *Queries) error {
		var err error
		if log_, err = queries.GetLogBySlug(ctx, logSlug); err != nil {
			return err
		}

		params.LogID = log_.ID

		if logEntry, err = queries.CreateLogEntry(ctx, params); err != nil {
			return err
		}

		return nil
	})

	return log_, logEntry, txErr
}
