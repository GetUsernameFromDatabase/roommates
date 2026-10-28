// Package middleware -- collection of [gin.HandlerFunc]s
package middleware

import (
	"roommates/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
)

var log = logger.Main

// NewZerolog https://gin-gonic.com/en/docs/logging/structured-logging/#using-zerolog
func NewZerolog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Str("client_ip", c.ClientIP()).
			Msg("request")
	}
}
