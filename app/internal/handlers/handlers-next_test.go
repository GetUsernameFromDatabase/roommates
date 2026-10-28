package handlers

import "testing"

func TestSafeNextPath(t *testing.T) {
	cases := map[string]string{
		"":                    "/",
		"/":                   "/",
		"/invite?code=ABC":    "/invite?code=ABC",
		"//evil.com":          "/",
		`/\evil.com`:          "/",
		"https://evil.com":    "/",
		"evil.com":            "/",
		"/ok\r\nSet-Cookie:x": "/",
		"/\t/evil.com":        "/",
		"/\x00/evil.com":      "/",
	}
	for next, want := range cases {
		if got := safeNextPath(next); got != want {
			t.Errorf("safeNextPath(%q) = %q, want %q", next, got, want)
		}
	}
}
