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
	"strconv"
)

var (
	logService *LogService
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %+v", err)
	}
	defer db.Close()

	logService = NewLogService(db)

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", LogsIndexController)
	r.GET("/logs", LogsIndexController)
	r.GET("/logs/:log_id", EntriesIndexController)
	r.GET("/logs/:log_id/entries", EntriesIndexController)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func LogsIndexController(c *gin.Context) {
	logs, err := logService.ListLogs(c)
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

func EntriesIndexController(c *gin.Context) {
	log2, entries, err := logService.GetLogAndEntriesBySlugOrderByDateDesc(c, c.Param("log_id"))
	if err != nil {
		c.Error(err)
		if errors.Is(err, ErrRecordNotFound) {
			c.String(http.StatusNotFound, "Log '%s' not found", c.Param("log_id"))
		} else {
			c.String(http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	var lastEntry queries.Entry
	if len(entries) > 0 {
		lastEntry = entries[0]
	}

	createEntryForm := NewLogEntryForm{
		Coffee:      lastEntry.Coffee,
		Water:       lastEntry.Water.String,
		Method:      lastEntry.Method.String,
		Grind:       lastEntry.Method.String,
		Tasting:     lastEntry.Tasting.String,
		AddlNotes:   lastEntry.AddlNotes.String,
		CoffeeGrams: "",
		WaterGrams:  "",
	}

	if lastEntry.CoffeeGrams.Valid {
		createEntryForm.CoffeeGrams = fmt.Sprintf("%d", lastEntry.CoffeeGrams.Int32)
	}
	if lastEntry.WaterGrams.Valid {
		createEntryForm.WaterGrams = fmt.Sprintf("%d", lastEntry.WaterGrams.Int32)
	}

	c.HTML(http.StatusOK, "entries/index.tmpl", gin.H{
		"Log":          log2,
		"Entries":      entries,
		"NewEntryForm": createEntryForm,
	})
}

func CreateEntryController(c *gin.Context) {
	logID, err := strconv.ParseInt(c.Param("log_id"), 10, 64)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, "invalid log ID")
		return
	}

	var form NewLogEntryForm
	c.Bind(&form)


}
