package middleware

import (
	g "roommates/internal/config"
	"roommates/internal/handlers"
	"roommates/internal/locales"

	"github.com/gin-gonic/gin"
)

func NewInternationalization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lang := ""
		// try to get locale from cookie
		if langCookie := g.GetCookie(ctx, g.CookieLanguage); langCookie != "" {
			if locales.Language(langCookie).IsValid() {
				lang = langCookie
			} else {
				// deleting cookie as it stores invalid language code
				g.DeleteCookie(ctx, g.CookieLanguage)
			}
		}
		if lang == "" {
			lang = handlers.GetHeader(ctx, g.HHK_AcceptLanguage)
		}
		ctx.Request = ctx.Request.WithContext(locales.WithLocale(ctx.Request.Context(), lang))
		ctx.Next()
	}
}
