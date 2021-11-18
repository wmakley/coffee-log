package internal

import (
	"coffee-log/db/sqlc"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	db *sql.DB
	router *gin.Engine
}

func NewServer(db *sql.DB) *Server {
	server := Server {
		db: db,
	}

	r := gin.Default()
	server.router = r

	r.LoadHTMLGlob("../templates/**/*")

	txOptions := sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}

	r.Use(sqlc.WrapInTransaction(db, &txOptions))
	r.Use(AuthMiddleware(AuthMiddleOptions{
		Realm:       "Coffee Log",
		MaxAttempts: 10,
		Debug:       true,
	}))

	logsController := NewLogsController()

	r.GET("/", logsController.FindOrCreateLogForUserAndRedirectToEntries)

	logs := r.Group("/logs")
	{
		logs.GET("/", logsController.FindOrCreateLogForUserAndRedirectToEntries)
		logs.GET("/:log_id", logsController.FindLogAndRedirectToEntries)
	}

	logEntries := logs.Group("/:log_id/log_entries")
	{
		logEntriesController := NewLogEntriesController()
		logEntries.GET("/", logEntriesController.Index)
		logEntries.POST("/", logEntriesController.Create)
	}

	return &server
}

func (server *Server)Run(address string, port int32) error {
	return server.router.Run(fmt.Sprintf("%s:%d", address, port))
}

func (server *Server)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	server.router.ServeHTTP(w, req)
}
