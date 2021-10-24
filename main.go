package main

import (
	"coffee-log/queries"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var (
	db *sql.DB
	q *queries.Queries
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")

	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %+v", err)
	}

	q = queries.New(db)

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", LogsIndex)
	r.GET("/logs", LogsIndex)
	r.GET("/logs/:log_id", LogEntriesIndex)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func LogsIndex(c *gin.Context) {
	logs, err := q.ListLogs(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Error(err)
		return
	}

	if len(logs) == 1 {
		c.Redirect(http.StatusFound, fmt.Sprintf("/logs/%s", logs[0].Slug))
		return
	}

	c.HTML(http.StatusOK, "logs/index.tmpl", gin.H{
		"Logs": logs,
	})
}

func LogEntriesIndex(c *gin.Context) {
	log2, err := q.GetLogBySlug(c, c.Param("log_id"))
	if err != nil {
		c.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusNotFound, "Log '%s' not found", c.Param("log_id"))
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	entries, err := q.ListLogEntriesByLogIdOrderByDateDesc(c, log2.ID)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.HTML(http.StatusOK, "entries/index.tmpl", gin.H{
		"Log": log2,
		"Entries": entries,
	})
}
