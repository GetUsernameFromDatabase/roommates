package middleware

import (
	"net/http"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/routes"
	"roommates/internal/sessions"

	"github.com/gin-gonic/gin"
)

// NewOnlyAuthenticated redirects unauthenticated users to [routes.Login]
func NewOnlyAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.FromGinContext(ctx)
		if s.IsLoggedIn() {
			ctx.Next()
			return
		}

		location := routes.Login
		// only a page can be returned to, not a form submission
		if ctx.Request.Method == http.MethodGet && ctx.Request.URL.Path != routes.Home {
			location += "?" + url.Values{g.QueryNext: {ctx.Request.URL.RequestURI()}}.Encode()
		}
		ctx.Redirect(http.StatusSeeOther, location)
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
		ctx.Redirect(http.StatusSeeOther, routes.Home)
		ctx.Abort()
	}
}
