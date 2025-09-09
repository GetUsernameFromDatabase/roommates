package middleware

import (
	"roomates/db/dbqueries"
	"roomates/rdb"
)

type MiddlewareHandlers interface {
	GetDB() *dbqueries.Queries
	GetRH() *rdb.RedisHandler
}
