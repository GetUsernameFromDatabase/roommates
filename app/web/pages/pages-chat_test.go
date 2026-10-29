package pages

import (
	"roommates/internal/db/dbqueries"
	"roommates/internal/services"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestChatRuns(t *testing.T) {
	mari, jaan := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}, pgtype.UUID{Bytes: [16]byte{2}, Valid: true}
	at := time.Date(2026, 10, 29, 14, 0, 0, 0, time.UTC)
	msg := func(author pgtype.UUID, minutes int) services.Message {
		return services.Message{ListMessagesRow: dbqueries.ListMessagesRow{
			ID:              pgtype.UUID{},
			AuthorName:      "",
			Body:            "",
			ReplyToID:       pgtype.UUID{},
			ReplyAuthorName: nil,
			ReplyBody:       "",
			AuthorID:        author,
			CreatedAt:       pgtype.Timestamptz{Time: at.Add(time.Duration(minutes) * time.Minute), InfinityModifier: pgtype.Finite, Valid: true},
		}, Files: nil}
	}
	// newest first: jaan 14:20, mari 14:10, mari 14:06, mari 14:00
	ms := []services.Message{msg(jaan, 20), msg(mari, 10), msg(mari, 6), msg(mari, 0)}
	wantStarts := []bool{true, false, true, true}
	wantEnds := []bool{true, true, false, true}
	for i := range ms {
		if got := startsRun(ms, i); got != wantStarts[i] {
			t.Errorf("startsRun(%d) = %v", i, got)
		}
		if got := endsRun(ms, i); got != wantEnds[i] {
			t.Errorf("endsRun(%d) = %v", i, got)
		}
	}
	for name, want := range map[string]string{"mari": "M", " õie": "Õ", "": "?"} {
		if got := initial(name); got != want {
			t.Errorf("initial(%q) = %q, want %q", name, got, want)
		}
	}
}
