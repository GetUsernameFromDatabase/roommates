package utils

import (
	"roommates/internal/locales"
	"strings"
	"testing"

	"github.com/invopop/ctxi18n"
)

func TestTUnescapesTextButNotArgs(t *testing.T) {
	locales.LoadLocales()
	ctx, err := ctxi18n.WithLocale(t.Context(), "en")
	if err != nil {
		t.Fatal(err)
	}

	body := T(ctx, locales.LK_MailVerifyBody, `a\nb "c"`, "https://x/y")
	if !strings.HasPrefix(body, "Hi a\\nb \"c\",\n\n") {
		t.Fatalf("text newlines must be real, the username kept as is: %q", body)
	}

	if got := unescapeJSON(`a&b\"c\\d`); got != `a&b"c\d` {
		t.Errorf("unescapeJSON = %q", got)
	}
	if got := unescapeJSON(`broken\`); got != `broken\` {
		t.Errorf("invalid escapes should be returned unchanged, got %q", got)
	}
}
