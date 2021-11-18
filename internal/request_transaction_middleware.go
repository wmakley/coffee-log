package internal

import (
	"coffee-log/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	txKey = "tx"
)

var (
	ErrNoTx = errors.New("no transaction")
)

func getTransactionFromContext(ctx context.Context) (*sql.Tx, error) {
	maybeTx := ctx.Value(txKey)
	if maybeTx == nil {
		return nil, ErrNoTx
	}

	tx, ok := maybeTx.(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("ctx key '%s' cannot be converted to *sql.Tx", txKey)
	}

	return tx, nil
}

// StoreFromCtx creates a Store using the context transaction, falling back
// to dbConn if the current context has no transaction
func StoreFromCtx(ctx context.Context, dbConn *sql.DB) *sqlc.Store {
	tx, err := getTransactionFromContext(ctx)
	if err != nil {
		if err == ErrNoTx {
			log.Print("context has no transaction, using original db conn")
			return sqlc.NewStore(dbConn)
		}
		panic(err)
	}

	return sqlc.StoreWithTx(tx)
}

func RequestTransaction(db *sql.DB, options *sql.TxOptions, debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := getTransactionFromContext(c)
		if err == nil {
			if debug {
				log.Print("request already has a transaction, using that instead of beginning new")
			}
			c.Next()
			return
		} else if err != ErrNoTx {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if debug {
			log.Print("beginning new request transaction")
		}
		tx, err := db.BeginTx(c, options)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error starting transaction: %w", err))
			return
		}

		c.Set(txKey, tx)
		c.Next()

		if len(c.Errors) > 0 {
			if debug {
				log.Print("rolling back due to errors")
			}
			if rbErr := tx.Rollback(); rbErr != nil {
				c.AbortWithError(http.StatusInternalServerError,
					fmt.Errorf("error rolling back due to prior errors: %v: %w", c.Errors, rbErr))
			}
			return
		}

		if commitErr := tx.Commit(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error committing: %w", commitErr))
		}
	}
}
