package handlers

import (
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/utils"

	"github.com/gin-gonic/gin"
)

const appearanceCookieMaxAge = 365 * 24 * 60 * 60

// SetAppearance stores the posted theme and mode in cookies and goes back to ?next=,
// unknown values become the defaults
func (h *DefaultHandler) SetAppearance(ctx *gin.Context) {
	a := utils.NewAppearance(ctx.PostForm("theme"), ctx.PostForm("mode"))
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(g.CookieTheme.String(), a.Theme, appearanceCookieMaxAge, "/", "", true, true)
	ctx.SetCookie(g.CookieMode.String(), a.Mode, appearanceCookieMaxAge, "/", "", true, true)
	Redirect(ctx, safeNextPath(ctx.PostForm(g.QueryNext)))
}
