package utils

import (
	"net/http"
	g "roommates/globals"
	"strings"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// Will respond to request with status code and error message
func ErrorResponse(ctx *gin.Context, status int, err error) {
	res := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, res)
}

// responds with 500 http
func ServerErrorResponse(ctx *gin.Context, publicError string) {
	code := http.StatusInternalServerError
	res := HTTPError{
		Code:    code,
		Message: publicError,
	}
	ctx.JSON(code, res)
}

func DeleteCookie(ctx *gin.Context, name g.CookieKey) {
	ctx.SetCookie(string(name), "", -1, "/", "", false, true)
}

// Gets authentication from header
func GetAuthTokenFromHeader(ctx *gin.Context) string {
	authHeader := ctx.GetHeader(string(g.HAuthorization))
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token
}

func GetAuthTokenFromCookie(ctx *gin.Context) string {
	cookie, err := ctx.Cookie(string(g.CSessionToken))
	if err != nil {
		return ""
	}
	return cookie
}

// will try to get auth token first from header then from cookie
func GetAuthToken(ctx *gin.Context) string {
	token := GetAuthTokenFromHeader(ctx)
	if token != "" {
		return token
	}
	return GetAuthTokenFromCookie(ctx)
}

func IsRequestHTMX(ctx *gin.Context) bool {
	return ctx.GetHeader(string(g.HHXRequest)) != ""
}

// redirects htmx request with header, otherwise uses ctx.Redirect with status see other
func Redirect(ctx *gin.Context, location string) {
	if IsRequestHTMX(ctx) {
		ctx.Header(string(g.HHXRedirect), "/")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/")
	}
}
