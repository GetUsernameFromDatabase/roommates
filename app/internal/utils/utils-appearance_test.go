package utils

import (
	"context"
	"testing"
)

func TestNewAppearance(t *testing.T) {
	cases := []struct {
		theme, mode string
		want        Appearance
	}{
		{"plain", "dark", Appearance{"plain", "dark"}},
		{"", "", Appearance{"home", "system"}},
		{"evil\"><script>", "light", Appearance{"home", "light"}},
		{"plain", "DARK", Appearance{"plain", "system"}},
	}
	for _, c := range cases {
		if got := NewAppearance(c.theme, c.mode); got != c.want {
			t.Errorf("NewAppearance(%q, %q) = %v, want %v", c.theme, c.mode, got, c.want)
		}
	}
	for _, theme := range Themes {
		if light, dark := NewAppearance(theme, "").Background(); light == "" || dark == "" {
			t.Errorf("theme %q has no theme-color backgrounds", theme)
		}
	}
	if got := GetAppearance(context.Background()); got != NewAppearance("", "") {
		t.Errorf("GetAppearance without middleware = %v, want the default", got)
	}
}
