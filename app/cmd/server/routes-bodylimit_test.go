package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLimitBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(limitBody(4))
	r.POST("/", func(c *gin.Context) {
		if _, err := io.ReadAll(c.Request.Body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})
	for _, tc := range []struct {
		body   string
		length int64
		status int
	}{
		{"abcd", 4, http.StatusOK},
		{"abcde", 5, http.StatusRequestEntityTooLarge},
		// chunked: no declared size, the reader cuts it
		{"abcde", -1, http.StatusBadRequest},
	} {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
		req.ContentLength = tc.length
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.status {
			t.Errorf("%d bytes, length %d: status %d, want %d", len(tc.body), tc.length, w.Code, tc.status)
		}
	}
}
