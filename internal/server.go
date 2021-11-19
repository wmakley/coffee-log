package internal

import (
	"coffee-log/internal/controller"
	"coffee-log/internal/middleware"
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

	r.Use(middleware.RequestTransaction(config.DB, &txOptions, config.Debug))
	r.Use(middleware.AuthMiddleware(middleware.AuthMiddlewareOptions{
		Realm:       "Coffee Log",
		MaxAttempts: 10,
		Debug:       config.Debug,
		DbConn:      config.DB,
	}))

	logsController := controller.NewLogsController(config.DB)

	r.GET("/", logsController.FindOrCreateLogForUserAndRedirectToEntries)

	logs := r.Group("/logs")
	{
		logs.GET("/", logsController.Index)
		logs.GET("/:log_id", logsController.Show)
	}

	logEntries := logs.Group("/:log_id/entries")
	{
		logEntriesController := controller.NewLogEntriesController(config.DB)
		logEntries.GET("/", logEntriesController.Index)
		logEntries.GET("/:id", logEntriesController.Show)
		logEntries.POST("/", logEntriesController.Create)
		logEntries.GET("/:id/edit", logEntriesController.Edit)
		logEntries.PATCH("/:id", logEntriesController.Update)
		logEntries.DELETE("/:id", logEntriesController.Delete)
	}

	return &server
}

func (server *Server) Run(address string, port int32) error {
	return server.router.Run(fmt.Sprintf("%s:%d", address, port))
}

func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	server.router.ServeHTTP(w, req)
}
