package internal

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	db     *sql.DB
	router *gin.Engine
}

func NewServer(db *sql.DB, debug bool) *Server {
	server := Server{
		db: db,
	}

	r := gin.Default()
	server.router = r

	r.LoadHTMLGlob("../templates/**/*")

	txOptions := sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}

	r.Use(RequestTransaction(db, &txOptions, debug))
	r.Use(AuthMiddleware(AuthMiddlewareOptions{
		Realm:       "Coffee Log",
		MaxAttempts: 10,
		Debug:       debug,
		DbConn:      db,
	}))

	logsController := NewLogsController(db)

	r.GET("/", logsController.FindOrCreateLogForUserAndRedirectToEntries)

	logs := r.Group("/logs")
	{
		logs.GET("/", logsController.FindOrCreateLogForUserAndRedirectToEntries)
		logs.GET("/:log_id", logsController.FindLogAndRedirectToEntries)
	}

	logEntries := logs.Group("/:log_id/entries")
	{
		logEntriesController := NewLogEntriesController(db)
		logEntries.GET("/", logEntriesController.Index)
		logEntries.POST("/", logEntriesController.Create)
	}

	return &server
}

func (server *Server) Run(address string, port int32) error {
	return server.router.Run(fmt.Sprintf("%s:%d", address, port))
}

func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	server.router.ServeHTTP(w, req)
}
