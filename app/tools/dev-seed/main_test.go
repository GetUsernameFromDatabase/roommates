package main

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	g "roommates/internal/config"
	"roommates/internal/db"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/push"
	"roommates/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type nullStore struct{}

func (nullStore) Put(context.Context, string, io.ReadSeeker, int64, string) error { return nil }
func (nullStore) Delete(context.Context, string) error                            { return nil }
func (nullStore) Get(context.Context, string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("")), nil
}

func TestCheckDevelopment(t *testing.T) {
	mode := gin.Mode()
	t.Cleanup(func() { gin.SetMode(mode) })

	gin.SetMode(gin.DebugMode)
	if err := checkDevelopment(); err != nil {
		t.Fatalf("debug mode: %v", err)
	}
	gin.SetMode(gin.ReleaseMode)
	if err := checkDevelopment(); err == nil {
		t.Fatal("release mode: want an error")
	}
}

// the first seed may already be in the test database, rows are left behind
func TestSeedRunsOnce(t *testing.T) {
	testDB := os.Getenv("TEST_POSTGRES_DB")
	if testDB == "" {
		t.Skip("TEST_POSTGRES_DB not set")
	}
	t.Setenv(g.EVK_DatabaseName.String(), testDB)
	ctx := t.Context()
	pool, err := pgxpool.New(ctx, g.GetDatabaseURL())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	db.LockTestDB(t, pool)
	db.MigrateToLatest(pool)
	locales.LoadLocales()
	s := services.NewDefaultService(pool, nil, nullStore{}, push.New())
	q := dbqueries.New(pool)

	seed(ctx, s, q)
	if seed(ctx, s, q) {
		t.Fatal("second seed created data")
	}
	mari, err := q.GetUserCredentials(ctx, emailMari)
	if err != nil {
		t.Fatal(err)
	}
	households, err := s.ListHouseholds(ctx, mari.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(households) != 1 {
		t.Fatalf("got %d households, want 1", len(households))
	}
}
