package services

import (
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestUserLimit(t *testing.T) {
	l := userLimit{mu: sync.Mutex{}, limit: 2, window: time.Hour, windowEnd: time.Time{}, counts: map[pgtype.UUID]int{}}
	a := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	b := pgtype.UUID{Bytes: [16]byte{2}, Valid: true}
	for range 3 {
		if !l.allow(a) {
			t.Fatal("a check without a send counts")
		}
	}
	for i, want := range []bool{true, true, false} {
		if got := l.allow(a); got != want {
			t.Errorf("send %d: allow = %v, want %v", i, got, want)
		}
		if want {
			l.add(a)
		}
	}
	if !l.allow(b) {
		t.Error("another user shares the limit")
	}
}
