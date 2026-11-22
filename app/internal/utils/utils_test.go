package utils

import (
	"crypto/rand"
	"testing"

	"roommates/internal/routes"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestParseUUID(t *testing.T) {
	var max [16]byte
	for i := range max {
		max[i] = 0xff
	}
	ids := []pgtype.UUID{{Bytes: [16]byte{}, Valid: true}, {Bytes: [16]byte{15: 1}, Valid: true}, {Bytes: max, Valid: true}}
	for range 1000 {
		var b [16]byte
		_, _ = rand.Read(b[:])
		ids = append(ids, pgtype.UUID{Bytes: b, Valid: true})
	}
	for _, id := range ids {
		s := routes.ShortUUID(id)
		if len(s) != 22 {
			t.Fatalf("ShortUUID(%s) = %q, want 22 characters", id, s)
		}
		for _, in := range []string{s, id.String()} {
			if got, err := ParseUUID(in); err != nil || got != id {
				t.Fatalf("ParseUUID(%q) = %v, %v, want %s", in, got, err, id)
			}
		}
	}
	for _, bad := range []string{
		"", "abc",
		"zzzzzzzzzzzzzzzzzzzzzz", // more than 128 bits
		"-000000000000000000001", "+000000000000000000001",
		"00000000000000000000_1", "000000000000000000000!",
	} {
		if _, err := ParseUUID(bad); err == nil {
			t.Errorf("ParseUUID(%q) accepted", bad)
		}
	}
}
