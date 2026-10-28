package middleware

import (
	"net/http"
	g "roommates/globals"
	"roommates/handlers"
	"roommates/locales"
	"roommates/logger"

	"github.com/gin-gonic/gin"
	"github.com/invopop/ctxi18n"
)

func NewInternationalization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// try to get locale from cookie
		langCookie := g.GetCookie(ctx, g.CookieLanguage)
		if langCookie != "" {
			lang := locales.Language(langCookie)
			if lang.IsValid() {
				err := setLocaleContext(ctx, lang.String())
				if err != nil {
					localeMissing(ctx, err, langCookie)
					return
				}
				ctx.Next()
				return
			}

			// deleting cookie as it stores invalid language code
			g.DeleteCookie(ctx, g.CookieLanguage)
		}

		// try to set locale with ones from header
		langHeader := handlers.GetHeader(ctx, g.HHK_AcceptLanguage)
		err := setLocaleContext(ctx, langHeader)
		if err == nil {
			ctx.Next()
			return
		}
		localeMissing(ctx, err, langHeader)
	}
}

// logs the error and sends response with [handlers.Fail]
// and aborts pending handlers with [gin.Context.Abort]
func localeMissing(ctx *gin.Context, err error, locale string) {
	logger.Main.Error().Stack().Err(err).
		Str("locale", locale).Str("DefaultLocale", ctxi18n.DefaultLocale.String()).
		Msg(g.EC_MissingLocale.String())
	handlers.Fail(ctx, http.StatusInternalServerError, g.EC_MissingLocale)
	ctx.Abort()
}

// tries to set locale for [gin.Context.Request] context
//
// [ctxi18n.DefaultLocale] is set in initGinEngine
func setLocaleContext(ctx *gin.Context, language string) error {
	cc, err := ctxi18n.WithLocale(ctx.Request.Context(), language)
	if err != nil {
		return err
	}

	ctx.Request = ctx.Request.WithContext(cc)
	return nil
}
