package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// arbitrary, only has to be the same for every test using TEST_POSTGRES_DB
const testDBLockKey = 7_201_026

// LockTestDB makes tests of different packages that share TEST_POSTGRES_DB
// run one at a time (go test runs packages in parallel), the lock is
// released when t ends
//
// a connection outside the pool, pool.Close waits for acquired connections
// and a deferred pool.Close runs before t.Cleanup
func LockTestDB(t testing.TB, pool *pgxpool.Pool) {
	t.Helper()
	conn, err := pgx.ConnectConfig(t.Context(), pool.Config().ConnConfig)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(t.Context(), "SELECT pg_advisory_lock($1)", testDBLockKey); err != nil {
		t.Fatal(err)
	}
	// the lock belongs to the session, closing it releases the lock
	t.Cleanup(func() { _ = conn.Close(context.Background()) })
}
