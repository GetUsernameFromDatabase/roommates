package controller

import (
	"fmt"
	"net/http"
	"roommates/db/dbqueries"
	"roommates/gintemplrenderer"
	"roommates/logger"
	"roommates/rdb"
	"roommates/utils"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var log = logger.ControllerLoggger

type SimpleResponse struct {
	Message string `json:"message" example:"OK"`
}

// ---| CONTROLLER ---

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

// --- CONTROLLER |---

func RenderTempl(ctx *gin.Context, component templ.Component) {
	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, component)
	ctx.Render(r.Status, r)
}

// logs the server error and responds with public error as message
func HandleServerError(ctx *gin.Context, err error, publicError string) {
	// TODO: dump request into a file for easier reproducibility
	log.Error().Err(err).Caller(1).Msg(publicError)
	utils.ServerErrorResponse(ctx, publicError)
}

// will get UUID from uri (:id)
//
// will also write a response to when UUID is not valid,
// return will be nil when that occurs
//
// this function only exists since binding pgtype.UUID is problematic, text marshal no bueno :(
func requirePgUUID(ctx *gin.Context, param string) *pgtype.UUID {
	reqParam := ctx.Param(param)
	var id pgtype.UUID
	id.Scan(reqParam)

	if !id.Valid {
		err := fmt.Errorf("invalid id (%s) in the url param", strconv.Quote(param))
		utils.ErrorResponse(ctx, http.StatusForbidden, err)
		return nil
	}
	return &id
}
