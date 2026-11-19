package handlers

import (
	"errors"
	"net/http"

	g "roommates/internal/config"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
)

// SetNotification turns one kind of push notification on or off (form kind, on=1) for the member
// in the household of the route and swaps only the notifications section
func (h *DefaultHandler) SetNotification(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	err := h.service.SetNotification(ctx, actor, ctx.PostForm("kind"), ctx.PostForm("on") == "1")
	if errors.Is(err, g.ErrorNotFound) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		InternalError(ctx, err)
		return
	}
	path, _ := settingsPath(ctx)
	if !IsRequestHTMX(ctx) {
		ctx.Redirect(http.StatusSeeOther, path)
		return
	}
	settings, err := h.service.ListNotificationSettings(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTempl(ctx, pages.NotificationSectionSaved(&actor, settings))
}
