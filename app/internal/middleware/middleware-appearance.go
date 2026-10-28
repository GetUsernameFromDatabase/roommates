package middleware

import (
	g "roommates/internal/config"
	"roommates/internal/utils"

	"github.com/gin-gonic/gin"
)

// NewAppearance puts the theme and mode from the cookies in the request context,
// see [utils.GetAppearance]
func NewAppearance() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		a := utils.NewAppearance(g.GetCookie(ctx, g.CookieTheme), g.GetCookie(ctx, g.CookieMode))
		ctx.Request = ctx.Request.WithContext(utils.WithAppearance(ctx.Request.Context(), a))
		ctx.Next()
	}
}
