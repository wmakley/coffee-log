package controller

import (
	"coffee-log/internal/middleware"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewLogsController(db *sql.DB) *LogsController {
	return &LogsController{
		db: db,
	}
}

type LogsController struct {
	db *sql.DB
}

type ShowLogParams struct {
	LogID string `uri:"log_id" binding:"required"`
}

func (controller *LogsController) FindLogAndRedirectToEntries(c *gin.Context) {
	params := ShowLogParams{}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, errorResponse(err))
		return
	}

	store := middleware.StoreFromCtx(c, controller.db)

	log, err := store.GetLogBySlug(c, params.LogID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, "log not found")
		} else {
			c.Error(err)
			c.String(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	c.Redirect(http.StatusFound, "/logs/"+log.Slug+"/entries/")
}

func (controller *LogsController) FindOrCreateLogForUserAndRedirectToEntries(c *gin.Context) {
	user, exists := middleware.GetCurrentUser(c)
	if !exists {
		panic("user not set")
	}
	store := middleware.StoreFromCtx(c, controller.db)

	log, err := store.FindOrCreateLogForUser(c, user)
	if err != nil {
		c.String(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.Redirect(http.StatusFound, "/logs/"+log.Slug+"/entries/")
}
