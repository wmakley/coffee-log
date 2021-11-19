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

type ServerConfig struct {
	DB           *sql.DB
	TemplateRoot string
	Debug        bool
}

func NewServer(config *ServerConfig) *Server {
	server := Server{
		db: config.DB,
	}

	r := gin.Default()
	server.router = r

	r.LoadHTMLGlob(config.TemplateRoot + "templates/**/*")

	txOptions := sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	}

	r.Use(RequestTransaction(config.DB, &txOptions, config.Debug))
	r.Use(AuthMiddleware(AuthMiddlewareOptions{
		Realm:       "Coffee Log",
		MaxAttempts: 10,
		Debug:       config.Debug,
		DbConn:      config.DB,
	}))

	logsController := NewLogsController(config.DB)

	r.GET("/", logsController.FindOrCreateLogForUserAndRedirectToEntries)

	logs := r.Group("/logs")
	{
		logs.GET("/", logsController.FindOrCreateLogForUserAndRedirectToEntries)
		logs.GET("/:log_id", logsController.FindLogAndRedirectToEntries)
	}

	logEntries := logs.Group("/:log_id/entries")
	{
		logEntriesController := NewLogEntriesController(config.DB)
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
