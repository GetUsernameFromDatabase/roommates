package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeaderHSTSOnlyInRelease(t *testing.T) {
	t.Cleanup(func() { gin.SetMode(gin.TestMode) })
	for mode, want := range map[string]bool{gin.ReleaseMode: true, gin.DebugMode: false} {
		gin.SetMode(mode)
		r := gin.New()
		r.GET("/", NewSecurityHeader(), func(c *gin.Context) { c.Status(http.StatusOK) })
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
		if got := w.Header().Get("Strict-Transport-Security") != ""; got != want {
			t.Errorf("%s mode: HSTS sent = %v, want %v", mode, got, want)
		}
	}
}
