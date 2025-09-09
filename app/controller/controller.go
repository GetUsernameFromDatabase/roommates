package controller

import (
	"roomates/db/dbqueries"
	"roomates/logger"
	"roomates/rdb"
	"roomates/utils"

	"github.com/gin-gonic/gin"
)

var log = logger.ControllerLoggger

// logs the server error and responds with public error as message
func HandleServerError(ctx *gin.Context, err error, publicError string) {
	log.Error().Err(err).Caller().Msg(publicError)
	utils.ServerErrorResponse(ctx, publicError)
}

type SimpleResponse struct {
	Message string `json:"message" example:"OK"`
}

type Controller struct {
	DB *dbqueries.Queries
	RH *rdb.RedisHandler
}

func New(db *dbqueries.Queries, rh *rdb.RedisHandler) *Controller {
	return &Controller{
		DB: db,
		RH: rh,
	}
}

func (c *Controller) GetDB() *dbqueries.Queries {
	return c.DB
}

func (c *Controller) GetRH() *rdb.RedisHandler {
	return c.RH
}
