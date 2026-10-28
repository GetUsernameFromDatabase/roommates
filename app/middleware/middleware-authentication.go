package middleware

import (
	"net/http"
	g "roommates/globals"
	"roommates/sessions"

	"github.com/gin-gonic/gin"
)

// NewOnlyAuthenticated redirects unauthenticated users to [g.RouteLogin]
func NewOnlyAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.FromGinContext(ctx)
		if s.IsLoggedIn() {
			ctx.Next()
			return
		}

		ctx.Redirect(http.StatusSeeOther, g.RouteLogin)
		ctx.Abort()
	}
}

// NewGuestOnly redirects logged in users to /
func NewGuestOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.FromGinContext(ctx)
		if !s.IsLoggedIn() {
			ctx.Next()
			return
		}

		// deliberately using ctx.Redirect
		ctx.Redirect(http.StatusSeeOther, "/")
		ctx.Abort()
	}
}
