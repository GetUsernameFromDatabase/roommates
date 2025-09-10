package middleware

import (
	"roommates/db/dbqueries"
	"roommates/rdb"
)

type MiddlewareHandlers interface {
	GetDB() *dbqueries.Queries
	GetRH() *rdb.RedisHandler
}
