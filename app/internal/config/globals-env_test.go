package config

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestGetDatabaseURLEscapes(t *testing.T) {
	const password = `a b'c\d@e/f?g`
	t.Setenv(EVK_DatabaseHost.String(), "localhost")
	t.Setenv(EVK_DatabasePort.String(), "5432")
	t.Setenv(EVK_DatabaseName.String(), "roommates")
	t.Setenv(EVK_DatabaseUser.String(), "us er")
	t.Setenv(EVK_DatabasePassword.String(), password)

	pool, err := pgxpool.ParseConfig(GetDatabaseURL())
	if err != nil {
		t.Fatal(err)
	}
	cfg := pool.ConnConfig
	if cfg.Password != password || cfg.User != "us er" || cfg.Database != "roommates" ||
		cfg.Host != "localhost" || cfg.Port != 5432 {
		t.Fatalf("parsed %q %q %q %q %d", cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)
	}
	if pool.MaxConns != 10 {
		t.Fatalf("default MaxConns = %d", pool.MaxConns)
	}

	t.Setenv(EVK_DatabaseMaxConns.String(), "3")
	pool, err = pgxpool.ParseConfig(GetDatabaseURL())
	if err != nil {
		t.Fatal(err)
	}
	if pool.MaxConns != 3 {
		t.Fatalf("MaxConns = %d", pool.MaxConns)
	}
}
