package internal

import (
	"coffee-log/db/sqlc"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	ShowLogParams
}

func (o *LogEntriesController) Index(ctx *gin.Context) {
	params := LogEntriesIndexParams{}
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.Error(err)
		ctx.String(http.StatusNotFound, errorResponse(err))
		return
	}

	store := StoreFromCtx(ctx, o.db)

	log2, entries, err := store.GetLogAndEntriesBySlugOrderByDateDesc(ctx, params.logID)
	if err != nil {
		ctx.Error(err)
		if err == sql.ErrNoRows {
			ctx.String(http.StatusNotFound, "Log '%s' not found", params.logID)
		} else {
			ctx.String(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	var lastEntry sqlc.LogEntry
	if len(entries) > 0 {
		lastEntry = entries[0]
	}

	createEntryForm := NewLogEntryForm{
		Coffee:       lastEntry.Coffee,
		Water:        lastEntry.Water.String,
		BrewMethod:   lastEntry.BrewMethod.String,
		GrindNotes:   lastEntry.GrindNotes.String,
		TastingNotes: "",
		AddlNotes:    lastEntry.AddlNotes.String,
		CoffeeGrams:  lastEntry.CoffeeGrams.Int32,
		WaterGrams:   lastEntry.WaterGrams.Int32,
	}

	ctx.HTML(http.StatusOK, "entries/index.tmpl", gin.H{
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

	var form NewLogEntryForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.Error(err)
		ctx.String(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}

	store := StoreFromCtx(ctx, o.db)

	arg := sqlc.CreateLogEntryParams{
		LogID:        0,
		EntryDate:    time.Now(),
		Coffee:       form.Coffee,
		Water:        blankToNullString(form.Water),
		BrewMethod:   blankToNullString(form.BrewMethod),
		GrindNotes:   blankToNullString(form.GrindNotes),
		TastingNotes: blankToNullString(form.TastingNotes),
		AddlNotes:    blankToNullString(form.AddlNotes),
		CoffeeGrams: sql.NullInt32{
			Int32: form.CoffeeGrams,
			Valid: form.CoffeeGrams > 0,
		},
		WaterGrams: sql.NullInt32{
			Int32: form.WaterGrams,
			Valid: form.WaterGrams > 0,
		},
	}

	log_, logEntry, err := store.CreateLogEntry(ctx, params.logID, arg)
	if err != nil {
		ctx.Error(err)
		if err == sql.ErrNoRows {
			ctx.String(http.StatusNotFound, "log '%s' not found", params.logID)
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
