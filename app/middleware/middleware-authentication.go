package middleware

import (
	"net/http"
	"roommates/docs"
	g "roommates/globals"
	"roommates/rdb"
	"roommates/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// will delete cookie, redirect to login (not on api endpoints) and abort future handlers
func unauthorize(ctx *gin.Context) {
	utils.DeleteCookie(ctx, string(g.CSessionToken))
	// API requests should not redirect to login when unauthorized
	if strings.HasPrefix(ctx.FullPath(), docs.SwaggerInfo.BasePath) {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, errors.New("unauthorized access"))
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
	ctx.Abort()
}

// sets auth info into gin context
func setAuthInfo(ctx *gin.Context, value *rdb.UserSessionValue) {
	ctx.Set(g.GAuth, value)
}

// gets auth info into gin context
func GetAuthInfo(ctx *gin.Context) *rdb.UserSessionValue {
	value, ok := ctx.Get(g.GAuth)
	if !ok {
		return nil
	}
	return value.(*rdb.UserSessionValue)
}

// sets authentication information into the context
//
// use block if this should unauthorize access to the endpoint
func NewAuthenticationMiddleware(ah MiddlewareHandlers, block bool) gin.HandlerFunc {
	if !block {
		return func(ctx *gin.Context) {
			token := utils.GetAuthToken(ctx)
			if token == "" {
				return
			}

			rh := ah.GetRH()
			usv, _ := rh.GetUserSession(ctx, token)
			if usv == nil {
				return
			}

			setAuthInfo(ctx, usv)
			ctx.Next()
		}
	}

	return func(ctx *gin.Context) {
		token := utils.GetAuthToken(ctx)
		if token == "" {
			unauthorize(ctx)
			return
		}

		rh := ah.GetRH()
		usv, err := rh.GetUserSession(ctx, token)
		if usv == nil {
			if err != nil {
				utils.ServerErrorResponse(ctx, "unable to get session")
				ctx.Abort()
			} else {
				unauthorize(ctx)
			}
			return
		}

		setAuthInfo(ctx, usv)
		ctx.Next()
	}
}
