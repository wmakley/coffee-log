package internal

import (
	"coffee-log/db/sqlc"
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
	logID string `uri:"log_id" binding:"required"`
}

func (con *LogsController) FindLogAndRedirectToEntries(c *gin.Context) {
	params := ShowLogParams{}
	if err := c.ShouldBindUri(&params); err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, errorResponse(err))
		return
	}

	store := StoreFromCtx(c, con.db)

	log, err := store.GetLogBySlug(c, params.logID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, "log not found")
		} else {
			c.Error(err)
			c.String(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	c.Redirect(http.StatusFound, "/logs/"+log.Slug+"/entries")
}

func (con *LogsController) FindOrCreateLogForUserAndRedirectToEntries(c *gin.Context) {
	userRaw, ok := c.Get("user")
	if !ok {
		panic("user not set")
	}
	user, ok := userRaw.(*sqlc.User)
	if !ok {
		panic("user is not *sqlc.User")
	}

	store := StoreFromCtx(c, con.db)

	log, err := store.GetLogByUserId(c, user.ID)

	if err == sql.ErrNoRows {
		log, err = store.CreateLog(c, sqlc.CreateLogParams{
			UserID: user.ID,
			Slug:   user.Username,
		})
	}

	if err != nil {
		c.String(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.Redirect(http.StatusFound, "/logs/"+log.Slug+"/entries")
}
