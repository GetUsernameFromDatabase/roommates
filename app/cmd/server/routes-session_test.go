package main

import (
	"net/http"
	"roommates/internal/routes"
	"roommates/internal/sessions"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestSignOutEverywhereEndsRealSession signs in through the real route, sign-out everywhere and
// a password reset find the session by the user_id that sessions.Login sets on its row,
// and an open live stream sees it gone through Session.Stored
func TestSignOutEverywhereEndsRealSession(t *testing.T) {
	r, service, signIn := newRouteTest(t)
	const storedPath = "/test-session-stored"
	r.GET(storedPath, func(c *gin.Context) {
		s := sessions.FromGinContext(c)
		stored, err := s.Stored(c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.String(http.StatusOK, strconv.FormatBool(stored))
	})
	stored := func(cookies []*http.Cookie) string {
		t.Helper()
		return send(r, http.MethodGet, storedPath, nil, cookies).Body.String()
	}

	userID, cookies := signIn("user")
	if w := send(r, http.MethodGet, routes.Home, nil, cookies); w.Code != http.StatusOK {
		t.Fatalf("signed in home got %d, want 200", w.Code)
	}
	if got := stored(cookies); got != "true" {
		t.Fatalf("Stored after sign-in got %q, want true", got)
	}
	if err := service.SignOutEverywhere(t.Context(), userID); err != nil {
		t.Fatal(err)
	}
	if got := stored(cookies); got != "false" {
		t.Fatalf("Stored after sign-out everywhere got %q, want false", got)
	}
	if w := send(r, http.MethodGet, routes.Home, nil, cookies); w.Code != http.StatusSeeOther || w.Header().Get("Location") != routes.Login {
		t.Fatalf("after sign-out everywhere got %d to %q, want a redirect to %s", w.Code, w.Header().Get("Location"), routes.Login)
	}
}
