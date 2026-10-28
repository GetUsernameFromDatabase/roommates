package config

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestGetDatabaseURLEscapes(t *testing.T) {
	const password = `a b'c\d@e/f?g`
	t.Setenv(EVK_DatabaseHost.String(), "localhost")
	t.Setenv(EVK_DatabasePort.String(), "5432")
	t.Setenv(EVK_DatabaseName.String(), "roommates")
	t.Setenv(EVK_DatabaseUser.String(), "us er")
	t.Setenv(EVK_DatabasePassword.String(), password)

	cfg, err := pgconn.ParseConfig(GetDatabaseURL())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Password != password || cfg.User != "us er" || cfg.Database != "roommates" ||
		cfg.Host != "localhost" || cfg.Port != 5432 {
		t.Fatalf("parsed %q %q %q %q %d", cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)
	}
}
