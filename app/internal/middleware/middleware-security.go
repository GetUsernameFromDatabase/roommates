package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// all scripts, styles, fonts and images are served from /assets
//
// img-src blob: is the chat composer's previews of picked images
//
// style-src keeps 'unsafe-inline' since htmx injects its indicator styles
// and components may use style attributes
var contentSecurityPolicy = strings.Join([]string{
	"default-src 'self'",
	"script-src 'self'",
	"style-src 'self' 'unsafe-inline'",
	"img-src 'self' data: blob:",
	"font-src 'self'",
	"connect-src 'self'",
	"form-action 'self'",
	"frame-ancestors 'none'",
	"base-uri 'self'",
	"object-src 'none'",
}, "; ")

// NewSecurityHeader https://gin-gonic.com/en/docs/middleware/security-headers/#example
func NewSecurityHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if c.Request.Host != expectedHost {
		// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
		// 	return
		// }
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", contentSecurityPolicy)
		// the XSS auditor is removed from browsers and could be abused, CSP replaces it
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/X-XSS-Protection
		c.Header("X-XSS-Protection", "0")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	}
}

// NewCrossOriginProtection rejects cross-origin non-safe requests (CSRF)
// using Sec-Fetch-Site / Origin headers, so no CSRF token is needed
//
// https://pkg.go.dev/net/http#CrossOriginProtection
func NewCrossOriginProtection() gin.HandlerFunc {
	cop := http.NewCrossOriginProtection()
	return func(c *gin.Context) {
		if err := cop.Check(c.Request); err != nil {
			c.String(http.StatusForbidden, "cross-origin request rejected")
			c.Abort()
			return
		}
		c.Next()
	}
}
