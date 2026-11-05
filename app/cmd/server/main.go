package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"os/signal"
	g "roommates/internal/config"
	"roommates/internal/db"
	"roommates/internal/locales"
	"roommates/internal/logger"
	"syscall"
	"time"

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
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	r, service := InitGinEngine(ctx, dbpool)
	// no WriteTimeout, it would cut the SSE streams
	//exhaustruct:ignore
	server := &http.Server{
		Addr:              serverAddr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	serveErr := make(chan error, 1)
	go func() {
		if r.TrustedPlatform != "" {
			serveErr <- server.ListenAndServe()
		} else {
			serveErr <- server.ListenAndServeTLS("./certificates/server.pem", "./certificates/server.key")
		}
	}()
	select {
	case err := <-serveErr:
		log.Fatal().Err(err).Msg("running http server has stopped")
	case <-ctx.Done():
	}

	log.Info().Msg("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("http server did not shut down cleanly")
	}
	service.WaitPushes()
	log.Info().Msg("server stopped")
}

// -----------------------------------------------------------------------------

// connects to database
//
// panics if cannot connect to database
func ensureDatabasePool(ctx context.Context) *pgxpool.Pool {
	dbURL := g.GetDatabaseURL()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal().Err(err).
			Msg("unable to create connection pool")
	}
	// ensure connection
	if err = dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		log.Fatal().Err(err).
			Msg("unable to ping database")
	}

	return dbpool
}
