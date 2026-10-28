package handlers

import (
	"roommates/internal/sessions"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
)

// PageHome is the landing page for signed in users
func (h *DefaultHandler) PageHome(ctx *gin.Context) {
	s := sessions.FromGinContext(ctx)
	RenderTempl(ctx, pages.Home(s.Username))
}
