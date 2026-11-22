// Package services coordinates API and database activity to carry out application logic
package services

import (
	"roommates/internal/db/dbqueries"
	"roommates/internal/logger"
	"roommates/internal/mailer"
	"roommates/internal/push"
	"roommates/internal/routes"
	"roommates/internal/storage"
	"roommates/internal/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var log = logger.Main.With().Str("package", "services").Logger()

type DefaultService struct {
	db *dbqueries.Queries
	// intended to be used for transactions
	// https://docs.sqlc.dev/en/latest/howto/transactions.html
	pool   *pgxpool.Pool
	mailer mailer.Mailer
	store  storage.Store
	pusher push.Pusher
	// notifications being sent, tests wait for them
	pushing sync.WaitGroup
}

func NewDefaultService(dbpool *pgxpool.Pool, m mailer.Mailer, store storage.Store, pusher push.Pusher) *DefaultService {
	dbHandler := dbqueries.New(dbpool)
	return &DefaultService{
		db:      dbHandler,
		pool:    dbpool,
		mailer:  m,
		store:   store,
		pusher:  pusher,
		pushing: sync.WaitGroup{},
	}
}

// Cursor is where a newest-first page ended, the time and id of its last row; the next page
// compares against these values and not against the row, which may be gone by then.
// The zero Cursor is the first page
type Cursor struct {
	At pgtype.Timestamptz
	ID pgtype.UUID
}

// String is the ?before= form of c, "<unix microseconds>-<short id>", read by [ParseCursor]
func (c Cursor) String() string {
	return strconv.FormatInt(c.At.Time.UnixMicro(), 10) + "-" + routes.ShortUUID(c.ID)
}

// ParseCursor reverses [Cursor.String], "" is the zero Cursor
func ParseCursor(s string) (Cursor, error) {
	if s == "" {
		return Cursor{}, nil
	}
	micros, id, _ := strings.Cut(s, "-")
	us, err := strconv.ParseInt(micros, 10, 64)
	if err != nil {
		return Cursor{}, err
	}
	uid, err := utils.ParseUUID(id)
	if err != nil {
		return Cursor{}, err
	}
	at := pgtype.Timestamptz{Time: time.UnixMicro(us), Valid: true, InfinityModifier: pgtype.Finite}
	return Cursor{At: at, ID: uid}, nil
}
