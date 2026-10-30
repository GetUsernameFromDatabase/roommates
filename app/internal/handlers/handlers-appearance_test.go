package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	g "roommates/internal/config"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetAppearance(t *testing.T) {
	cases := []struct {
		name     string
		htmx     bool
		status   int
		location string
	}{
		// a reload here would lose a picture being cropped
		{"htmx stays on the page", true, http.StatusNoContent, ""},
		{"plain form goes back to next", false, http.StatusSeeOther, "/settings"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			form := url.Values{"theme": {"plain"}, "mode": {"dark"}, g.QueryNext: {"/settings"}}
			req := httptest.NewRequest(http.MethodPost, g.RouteAppearance, strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if c.htmx {
				req.Header.Set(g.HHK_HXRequest.String(), "true")
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			(&DefaultHandler{}).SetAppearance(ctx)
			ctx.Writer.WriteHeaderNow()

			res := w.Result()
			if res.StatusCode != c.status {
				t.Errorf("status = %d, want %d", res.StatusCode, c.status)
			}
			if got := res.Header.Get("Location"); got != c.location {
				t.Errorf("Location = %q, want %q", got, c.location)
			}
			if got := res.Header.Get(g.HHK_HXRedirect.String()); got != "" {
				t.Errorf("HX-Redirect = %q, want none", got)
			}
			cookies := map[string]string{}
			for _, cookie := range res.Cookies() {
				cookies[cookie.Name] = cookie.Value
			}
			if cookies[g.CookieTheme.String()] != "plain" || cookies[g.CookieMode.String()] != "dark" {
				t.Errorf("cookies = %v, want theme plain and mode dark", cookies)
			}
		})
	}
}
