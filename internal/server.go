package internal

import (
	"coffee-log/db/sqlc"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Address string
	Port int32
	db *sql.DB
	router *gin.Engine
}

func NewServer(addr string, port int32, db *sql.DB) *Server {
	server := Server {
		Address: addr,
		Port: port,
		db: db,
	}

	r := gin.Default()
	server.router = r

	r.LoadHTMLGlob("templates/**/*")

	txOptions := sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}

	r.Use(sqlc.WrapInTransaction(db, &txOptions))
	r.Use(AuthMiddleware("Coffee Log", int32(10)))

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

func (server *Server)Run() error {
	address := fmt.Sprintf("%s:%d", server.Address, server.Port)
	return server.router.Run(address)
}

func (server *Server)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	server.router.ServeHTTP(w, req)
}
