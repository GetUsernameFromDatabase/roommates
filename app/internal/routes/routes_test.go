package routes

import (
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestURL(t *testing.T) {
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	if got, want := URL(ChoreDone, id, id), "/h/01SZwviYzes2mjOamuMJWw/chores/01SZwviYzes2mjOamuMJWw/done"; got != want {
		t.Errorf("URL(ChoreDone) = %q, want %q", got, want)
	}
	if got := URL(Settings); got != "/settings" {
		t.Errorf("URL(Settings) = %q", got)
	}
	if got := URL(Member, "a/b c", "x"); got != "/h/a%2Fb%20c/members/x" {
		t.Errorf("URL escapes = %q", got)
	}
	for _, args := range [][]any{{}, {id}, {id, id, id}, {id, 3}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("URL(ChoreDone, %v) did not panic", args)
				}
			}()
			URL(ChoreDone, args...)
		}()
	}
}

func TestJSSource(t *testing.T) {
	for k, v := range JS {
		if !strings.HasPrefix(v, "/") {
			t.Errorf("JS[%q] = %q, want an absolute path", k, v)
		}
	}
	if s := string(JSSource); !strings.HasPrefix(s, "self.ROUTES = {") {
		t.Errorf("JSSource = %q", s)
	}
	if JSHash == "" || JSURL != RoutesJS+"?v="+JSHash {
		t.Errorf("JSURL = %q", JSURL)
	}
}
