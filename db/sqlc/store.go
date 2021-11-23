package sqlc

import (
	"coffee-log/util"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Store struct {
	*Queries
	db    *sql.DB
	tx    *sql.Tx
	Debug bool
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

/// transaction is an internal method using for manually creating transactions
func (store *Store) transaction(ctx context.Context, fn func(*Store) error) error {
	if store.tx != nil {
		if store.Debug {
			log.Print("store transaction: store already has a transaction, re-using that instead of beginning new")
		}
		return fn(store)
	}

	if store.Debug {
		log.Print("store transaction: beginning")
	}
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	store.tx = tx

	innerStore := Store{
		Queries: New(tx),
		db:      nil,
		tx:      tx,
	}
	err = fn(&innerStore)
	if err != nil {
		if store.Debug {
			log.Print("store transaction: rolling back")
		}
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("transaction err: %+v, rollback error: %+v", err, rbErr)
		}
		return err
	}

	if store.Debug {
		log.Print("store transaction: committing")
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

	banned, err := store.createOrIncrementLoginAttempt(ctx, ipAddress, maxAttempts)
	if err != nil {
		return user, err
	}
	if banned {
		return user, ErrIPBanned
	}

	return user, ErrBadCredentials
}

func (store *Store) createOrIncrementLoginAttempt(ctx context.Context, ip string, maxAttempts int32) (banned bool, err error) {
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

	err := store.transaction(ctx, func(store *Store) error {
		var err error

		log2, err = store.GetLogBySlug(ctx, slug)
		if err != nil {
			return err
		}

		entries, err = store.ListLogEntriesByLogIDOrderByDateDesc(ctx, log2.ID)
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
	ctx context.Context,
	logSlug string,
	params CreateLogEntryParams,
) (
	Log, LogEntry, error,
) {
	var log_ Log
	var logEntry LogEntry

	txErr := store.transaction(ctx, func(store *Store) error {
		var err error
		if log_, err = store.Queries.GetLogBySlug(ctx, logSlug); err != nil {
			return err
		}

		params.LogID = log_.ID

		if logEntry, err = store.Queries.CreateLogEntry(ctx, params); err != nil {
			return err
		}

		return nil
	})

	return log_, logEntry, txErr
}

func (store *Store) DeleteAllLoginAttemptsAndBans(ctx context.Context) error {
	return store.transaction(ctx, func(store *Store) error {
		err := store.DeleteAllBannedIPs(ctx)
		if err != nil {
			return err
		}

		return store.DeleteAllLoginAttempts(ctx)
	})
}

func (store *Store) FindOrCreateLogForUser(ctx context.Context, user *User) (Log, error) {
	var userLog Log

	err := store.transaction(ctx, func(store *Store) error {
		var err error
		userLog, err = store.GetLogByUserId(ctx, user.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				var slug string
				if slug, err = util.Sluggify(user.Username); err != nil {
					return err
				}
				userLog, err = store.CreateLog(ctx, CreateLogParams{
					UserID: user.ID,
					Slug:   slug,
					Title:  fmt.Sprintf("%s's Log", user.DisplayName),
				})
			}
			return err
		}

		return nil
	})

	return userLog, err
}
