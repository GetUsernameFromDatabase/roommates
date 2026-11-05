package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	g "roommates/internal/config"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// guest cases only, signed in it also saves to the account (needs the database)
func TestSetLanguage(t *testing.T) {
	cases := []struct {
		name     string
		language string
		htmx     bool
		status   int
		location string
		refresh  string
		cookie   string
	}{
		// every text changes, so htmx reloads the page
		{"htmx reloads", "en", true, http.StatusNoContent, "", "true", "en"},
		{"plain form goes back to next", "et", false, http.StatusSeeOther, "/login", "", "et"},
		{"unknown language is refused", "xx", true, http.StatusBadRequest, "", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := gin.New()
			r.Use(sessions.Sessions("test", cookie.NewStore([]byte("test-secret"))))
			r.POST(g.RouteLanguage, (&DefaultHandler{}).SetLanguage)

			form := url.Values{"language": {c.language}, g.QueryNext: {"/login"}}
			req := httptest.NewRequest(http.MethodPost, g.RouteLanguage, strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if c.htmx {
				req.Header.Set(g.HHK_HXRequest.String(), "true")
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			res := w.Result()
			if res.StatusCode != c.status {
				t.Errorf("status = %d, want %d", res.StatusCode, c.status)
			}
			if got := res.Header.Get("Location"); got != c.location {
				t.Errorf("Location = %q, want %q", got, c.location)
			}
			if got := res.Header.Get(g.HHK_HXRefresh.String()); got != c.refresh {
				t.Errorf("HX-Refresh = %q, want %q", got, c.refresh)
			}
			got := ""
			for _, ck := range res.Cookies() {
				if ck.Name == g.CookieLanguage.String() {
					got = ck.Value
				}
			}
			if got != c.cookie {
				t.Errorf("language cookie = %q, want %q", got, c.cookie)
			}
		})
	}
}
