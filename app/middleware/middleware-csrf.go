package middleware

import (
	"net/http"
	"roommates/globals"
	"roommates/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
)

// takes auth key from env
func NewCSRFMiddleware() gin.HandlerFunc {
	bAuthKey := []byte(utils.MustGetEnv("CSRF"))
	if len(bAuthKey) != 32 {
		panic("auth key should be 32 bytes")
	}

	// TODO: get this working, helpful links below
	// - https://github.com/gorilla/csrf
	// - https://stackoverflow.com/questions/69860669/how-set-middleware-gorillas-csrf-in-framework-gin-from-go
	// - https://templ.guide/integrations/web-frameworks/#githubcomgorillacsrf
	// - https://htmx.org/docs/#csrf-prevention

	return func(ctx *gin.Context) {
		gorillaCSRF := csrf.Protect(
			bAuthKey,
			csrf.ErrorHandler(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				utils.ErrorResponse(ctx, http.StatusForbidden, errors.New("CSRF token mismatch"))
				ctx.Abort()
			})),
			csrf.FieldName(globals.Csrf),
		)

		gorillaCSRF(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			ctx.Next()
		})).ServeHTTP(ctx.Writer, ctx.Request)
	}
}
