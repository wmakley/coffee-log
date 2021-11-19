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

func (o *LogEntriesController) Index(c *gin.Context) {
	var params LogEntriesIndexParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, errorResponse(err))
		return
	}

	store := middleware.StoreFromCtx(c, o.db)

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
	ShowLogParams
	id int64 `uri:"id"`
}

func (o *LogEntriesController) Show(ctx *gin.Context) {
	params := ShowLogEntryParams{}
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.Error(err)
		ctx.String(http.StatusNotFound, errorResponse(err))
		return
	}
}

func (o *LogEntriesController) Create(ctx *gin.Context) {
	var params LogEntriesIndexParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.Error(err)
		ctx.String(http.StatusNotFound, errorResponse(err))
		return
	}

	var entryForm form.LogEntryForm
	if err := ctx.ShouldBind(&entryForm); err != nil {
		ctx.Error(err)
		ctx.String(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	store := middleware.StoreFromCtx(ctx, o.db)

	arg := entryForm.CreateParams()

	log_, logEntry, err := store.CreateLogEntry(ctx, params.LogID, arg)
	if err != nil {
		ctx.Error(err)
		if err == sql.ErrNoRows {
			ctx.String(http.StatusNotFound, "log '%s' not found", params.LogID)
		} else {
			ctx.String(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.HTML(http.StatusOK, "entries/index.tmpl", gin.H{
		"Log":             log_,
		"CreatedLogEntry": logEntry,
	})
}
