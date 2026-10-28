package db

import (
	"os"
	g "roommates/internal/config"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TestMigrateUpAndDown walks the migrations up and down against a real database
//
// it drops every table, so it only runs against the database named in
// TEST_POSTGRES_DB (the other DATABASE_* and POSTGRES_* variables are reused)
func TestMigrateUpAndDown(t *testing.T) {
	testDB := os.Getenv("TEST_POSTGRES_DB")
	if testDB == "" {
		t.Skip("TEST_POSTGRES_DB not set")
	}
	t.Setenv(g.EVK_DatabaseName.String(), testDB)

	pool, err := pgxpool.New(t.Context(), g.GetDatabaseURL())
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	maxVersion := NewMigrations(pool).maxVersion
	if maxVersion < 1 {
		t.Fatal("no migrations found")
	}

	// new [migrations] each step, like a restart of the app would
	steps := []int{0, maxVersion, 0, 1, maxVersion, 1, maxVersion}
	for _, version := range steps {
		NewMigrations(pool).MigrateTo(version)
		if got := SelectMigrationVersion(pool); got != version {
			t.Fatalf("after MigrateTo(%d) version is %d", version, got)
		}

		var usersExists bool
		err := pool.QueryRow(t.Context(),
			`SELECT to_regclass('users') IS NOT NULL`).Scan(&usersExists)
		if err != nil {
			t.Fatal(err)
		}
		if usersExists != (version > 0) {
			t.Fatalf("at version %d users table exists = %v", version, usersExists)
		}
	}
}
