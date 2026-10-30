package handlers

import (
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/utils"

	"github.com/gin-gonic/gin"
)

const appearanceCookieMaxAge = 365 * 24 * 60 * 60

// SetAppearance stores the posted theme and mode in cookies, unknown values become the defaults,
// htmx stays on the page that already shows them, a plain form goes back to ?next=
func (h *DefaultHandler) SetAppearance(ctx *gin.Context) {
	a := utils.NewAppearance(ctx.PostForm("theme"), ctx.PostForm("mode"))
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(g.CookieTheme.String(), a.Theme, appearanceCookieMaxAge, "/", "", true, true)
	ctx.SetCookie(g.CookieMode.String(), a.Mode, appearanceCookieMaxAge, "/", "", true, true)
	if IsRequestHTMX(ctx) {
		ctx.Status(http.StatusNoContent)
		return
	}
	Redirect(ctx, safeNextPath(ctx.PostForm(g.QueryNext)))
}
