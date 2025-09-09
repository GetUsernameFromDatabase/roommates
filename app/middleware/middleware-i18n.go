package middleware

import (
	"net/http"
	"roomates/locales"
	"roomates/logger"
	"roomates/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/invopop/ctxi18n"
	"github.com/pkg/errors"
)

func NewLanguageMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var lang locales.Language
		pathSegments := strings.Split(ctx.Request.URL.Path, "/")
		if len(pathSegments) > 1 {
			lang = locales.Language(pathSegments[1])
		}
		if !lang.Valid() {
			lang = locales.Default
		}

		cc, err := ctxi18n.WithLocale(ctx.Request.Context(), string(lang))
		if err != nil {
			publicErr := "error setting locale"
			logger.Main.Error().Err(err).Msg(publicErr)
			utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New(publicErr))
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(cc)
		ctx.Next()
	}
}
