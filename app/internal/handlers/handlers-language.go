package handlers

import (
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/sessions"

	"github.com/gin-gonic/gin"
)

// setLanguageCookie is what [middleware.NewInternationalization] reads first
func setLanguageCookie(ctx *gin.Context, lang locales.Language) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(g.CookieLanguage.String(), lang.String(), appearanceCookieMaxAge, "/", "", true, true)
}

// SetLanguage stores the posted language in a cookie and, when signed in, on the account,
// every text on the page changes so htmx reloads it, a plain form goes back to ?next=
func (h *DefaultHandler) SetLanguage(ctx *gin.Context) {
	lang := locales.Language(ctx.PostForm("language"))
	if !lang.IsValid() {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if s := sessions.FromGinContext(ctx); s.IsLoggedIn() {
		userID, err := sessionUserID(ctx)
		if err != nil {
			InternalError(ctx, err)
			return
		}
		if err := h.service.SetUserLanguage(ctx, userID, lang.String()); err != nil {
			InternalError(ctx, err)
			return
		}
	}
	setLanguageCookie(ctx, lang)
	if IsRequestHTMX(ctx) {
		ctx.Header(g.HHK_HXRefresh.String(), "true")
		ctx.Status(http.StatusNoContent)
		return
	}
	Redirect(ctx, safeNextPath(ctx.PostForm(g.QueryNext)))
}
