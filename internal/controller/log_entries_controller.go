package controller

import (
	"coffee-log/db/sqlc"
	"coffee-log/internal/form"
	"coffee-log/internal/middleware"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewLogEntriesController(db *sql.DB) *LogEntriesController {
	return &LogEntriesController{
		db: db,
	}
}

type LogEntriesController struct {
	db *sql.DB
}

type LogEntriesIndexParams struct {
	LogID string `uri:"log_id" binding:"required"`
}

func (controller *LogEntriesController) Index(c *gin.Context) {
	var params LogEntriesIndexParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, errorResponse(err))
		return
	}

	store := middleware.StoreFromCtx(c, controller.db)

	log2, entries, err := store.GetLogAndEntriesBySlugOrderByDateDesc(c, params.LogID)
	if err != nil {
		c.Error(err)
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, "Log '%s' not found", params.LogID)
		} else {
			c.String(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	var lastEntry sqlc.LogEntry
	if len(entries) > 0 {
		lastEntry = entries[0]
	}

	createEntryForm := form.LogEntryForm{
		Coffee:       lastEntry.Coffee,
		Water:        lastEntry.Water.String,
		BrewMethod:   lastEntry.BrewMethod.String,
		GrindNotes:   lastEntry.GrindNotes.String,
		TastingNotes: "",
		AddlNotes:    lastEntry.AddlNotes.String,
		CoffeeGrams:  lastEntry.CoffeeGrams.Int32,
		WaterGrams:   lastEntry.WaterGrams.Int32,
	}

	c.HTML(http.StatusOK, "entries/index.tmpl", gin.H{
		"Log":          log2,
		"Entries":      entries,
		"NewEntryForm": createEntryForm,
	})
}

type ShowLogEntryParams struct {
	LogEntriesIndexParams
	ID int64 `uri:"id" binding:"required"`
}

func (controller *LogEntriesController) Show(c *gin.Context) {
	params := ShowLogEntryParams{}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, errorResponse(err))
		return
	}
}

func (controller *LogEntriesController) Create(c *gin.Context) {
	store := middleware.StoreFromCtx(c, controller.db)

	var params LogEntriesIndexParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, errorResponse(err))
		return
	}

	var entryForm form.LogEntryForm
	if err := c.ShouldBind(&entryForm); err != nil {
		c.Error(err)
		c.String(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	arg := entryForm.CreateParams()

	log_, logEntry, err := store.CreateLogEntry(c, params.LogID, arg)
	if err != nil {
		c.Error(err)
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, "log '%s' not found", params.LogID)
		} else {
			c.String(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	c.HTML(http.StatusOK, "entries/index.tmpl", gin.H{
		"Log":             log_,
		"CreatedLogEntry": logEntry,
	})
}

func (controller *LogEntriesController) Edit(c *gin.Context) {
	// TODO
}

func (controller *LogEntriesController) Update(c *gin.Context) {
	// TODO
}

func (controller *LogEntriesController) Delete(c *gin.Context) {
	// TODO
}
