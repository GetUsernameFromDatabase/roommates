package utils

import (
	"roommates/internal/locales"
	"testing"
	"time"

	"github.com/invopop/ctxi18n"
)

func TestFormatPerLanguage(t *testing.T) {
	locales.LoadLocales()
	day := time.Date(2026, 11, 19, 23, 0, 0, 0, time.UTC)
	for lang, want := range map[string][]string{
		"et": {"19.11.2026", "12,50 €", "1 234,50 €", "512 B", "2 KB", "1,5 MB"},
		"en": {"19 Nov 2026", "€12.50", "€1,234.50", "512 B", "2 KB", "1.5 MB"},
	} {
		ctx, err := ctxi18n.WithLocale(t.Context(), lang)
		if err != nil {
			t.Fatal(err)
		}
		got := []string{FormatDate(ctx, day), FormatMoney(ctx, 1250, "EUR"), FormatMoney(ctx, 123450, "EUR"),
			FormatSize(ctx, 512), FormatSize(ctx, 2<<10), FormatSize(ctx, 3<<19)}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("%s: got %q, want %q", lang, got[i], want[i])
			}
		}
	}
}
