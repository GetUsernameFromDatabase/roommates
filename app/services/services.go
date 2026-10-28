// Package services coordinates API and database activity to carry out application logic
package services

import (
	"roommates/db/dbqueries"
	"roommates/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var log = logger.Main.With().Str("package", "services").Logger()

type DefaultService struct {
	db *dbqueries.Queries
	// intended to be used for transactions
	// https://docs.sqlc.dev/en/latest/howto/transactions.html
	pool *pgxpool.Pool
}

func NewDefaultService(dbpool *pgxpool.Pool) *DefaultService {
	dbHandler := dbqueries.New(dbpool)
	return &DefaultService{
		db:   dbHandler,
		pool: dbpool,
	}
}
