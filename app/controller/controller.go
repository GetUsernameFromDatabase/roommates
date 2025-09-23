package controller

import (
	"roommates/db/dbqueries"
	"roommates/logger"
	"roommates/rdb"
	"roommates/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var log = logger.ControllerLoggger

// logs the server error and responds with public error as message
func HandleServerError(ctx *gin.Context, err error, publicError string) {
	// TODO: dump request into a file for easier reproducibility
	log.Error().Err(err).Caller().Msg(publicError)
	utils.ServerErrorResponse(ctx, publicError)
}

type SimpleResponse struct {
	Message string `json:"message" example:"OK"`
}

type Controller struct {
	DB   *dbqueries.Queries
	RH   *rdb.RedisHandler
	Pool *pgxpool.Pool
}

func New(dbpool *pgxpool.Pool, rh *rdb.RedisHandler) *Controller {
	dbHandler := dbqueries.New(dbpool)
	return &Controller{
		DB:   dbHandler,
		RH:   rh,
		Pool: dbpool,
	}
}

func (c *Controller) GetDB() *dbqueries.Queries {
	return c.DB
}

func (c *Controller) GetRH() *rdb.RedisHandler {
	return c.RH
}
