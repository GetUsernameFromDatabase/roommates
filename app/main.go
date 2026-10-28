package main

import (
	"context"
	_ "embed"
	"roommates/db"
	g "roommates/globals"
	"roommates/locales"
	"roommates/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var log = logger.Main

func main() {
	ctx := context.Background()
	isDevelopmentMode := g.IsDevelopment()
	serverAddr := ":" + g.EVK_ServerPort.Get("8080")
	logger.InitLogger(isDevelopmentMode)

	dbpool := ensureDatabasePool(ctx)
	defer dbpool.Close()

	// to reset database should it be required
	// db.NewMigrations(dbpool).MigrateTo(0)
	db.MigrateToLatest(dbpool)

	locales.LoadLocales()
	r := InitGinEngine(dbpool)
	// TODO: https://gin-gonic.com/en/docs/middleware/security-guide/#https-and-tls
	err := r.RunTLS(serverAddr, "./certificates/server.pem", "./certificates/server.key")
	log.Fatal().Err(err).Msg("running http server has stopped")
}

// -----------------------------------------------------------------------------

// connects to database
//
// panics if cannot connect to database
func ensureDatabasePool(ctx context.Context) *pgxpool.Pool {
	dbURL := g.GetDatabaseURL()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal().Stack().Err(err).
			Msg("unable to create connection pool")
	}
	// ensure connection
	if err = dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		log.Fatal().Stack().Err(err).
			Msg("unable to ping database")
	}

	return dbpool
}
