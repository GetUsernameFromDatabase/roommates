package pages

import (
	"strings"
	"testing"
)

func TestNoteHTMLIsSafe(t *testing.T) {
	var b strings.Builder
	src := "**hi**\n<script>alert(1)</script>\n[x](javascript:alert(1)) <img src=x onerror=alert(1)>"
	if err := noteHTML(src).Render(t.Context(), &b); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if !strings.Contains(out, "<strong>hi</strong>") {
		t.Errorf("markdown not rendered: %s", out)
	}
	for _, bad := range []string{"<script", "javascript:", "onerror"} {
		if strings.Contains(out, bad) {
			t.Errorf("output contains %q: %s", bad, out)
		}
	}
}
