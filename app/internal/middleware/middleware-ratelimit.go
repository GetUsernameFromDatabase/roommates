package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// NewRateLimit allows each client IP at most limit requests per window,
// later ones get a plain text 429
//
// shortcut: counters live in this process, move them to a shared store
// (e.g. the database) once more than one app process runs
func NewRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	counts := map[string]int{}
	windowEnd := time.Now().Add(window)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		if now := time.Now(); now.After(windowEnd) {
			// a fresh map every window keeps memory bounded
			counts = map[string]int{}
			windowEnd = now.Add(window)
		}
		counts[ip]++
		over := counts[ip] > limit
		mu.Unlock()

		if over {
			c.String(http.StatusTooManyRequests, "too many requests, try again later")
			c.Abort()
			return
		}
		c.Next()
	}
}
