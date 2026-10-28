package config

import "github.com/gin-gonic/gin"

type CookieKey string

const (
	CookieLanguage CookieKey = "language"
	// see utils.Appearance
	CookieTheme CookieKey = "theme"
	CookieMode  CookieKey = "mode"
)

func (c CookieKey) String() string {
	return string(c)
}

// -----------------------------------------------------------------------------

// GetCookie is a wrapper of [gin.Context.Cookie]
//
// error ErrNoCookie is ignored and empty string is returned instead
func GetCookie(ctx *gin.Context, name CookieKey) string {
	cookie, err := ctx.Cookie(name.String())
	if err != nil {
		return ""
	}
	return cookie
}

// DeleteCookie is a wrapper of [gin.Context.SetCookie]
//
// sets cookie value to "" and maxAge to -1 to delete it
func DeleteCookie(ctx *gin.Context, name CookieKey) {
	ctx.SetCookie(name.String(), "", -1, "/", "", false, true)
}
