package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/", NewRateLimit(2, time.Hour), func(c *gin.Context) { c.Status(http.StatusOK) })

	post := func(remoteAddr string) int {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.RemoteAddr = remoteAddr
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w.Code
	}

	for i, want := range []int{http.StatusOK, http.StatusOK, http.StatusTooManyRequests} {
		if got := post("1.2.3.4:1000"); got != want {
			t.Fatalf("request %d from same IP: status %d, want %d", i+1, got, want)
		}
	}
	if got := post("5.6.7.8:1000"); got != http.StatusOK {
		t.Fatalf("other IP: status %d, want %d", got, http.StatusOK)
	}
}
