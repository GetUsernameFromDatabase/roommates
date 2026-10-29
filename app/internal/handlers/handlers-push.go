package handlers

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/push"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
)

// SubscribePush stores the browser's push subscription for the signed in user
func (h *DefaultHandler) SubscribePush(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	var sub push.Subscription
	if err := ctx.ShouldBindJSON(&sub); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	err = h.service.SubscribePush(ctx, userID, sub, utils.GetLocaleString(ctx.Request.Context()))
	if errors.Is(err, g.ErrorInvalidInput) {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if err != nil {
		InternalError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// UnsubscribePush forgets the browser's subscription, sent before logging out too
func (h *DefaultHandler) UnsubscribePush(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	var body struct {
		Endpoint string `json:"endpoint"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if err := h.service.UnsubscribePush(ctx, userID, body.Endpoint); err != nil {
		InternalError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Manifest is the web app manifest in the request's language
//
// https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/Manifest
func (h *DefaultHandler) Manifest(ctx *gin.Context) {
	name := utils.T(ctx.Request.Context(), locales.LK_AppTitle)
	ctx.Header("Content-Type", "application/manifest+json")
	ctx.JSON(http.StatusOK, gin.H{
		"name":       name,
		"short_name": name,
		"lang":       utils.GetLocaleString(ctx.Request.Context()),
		"start_url":  "/",
		"scope":      "/",
		"display":    "standalone",
		"icons": []gin.H{
			{"src": "/icons/web-app-manifest-192x192.png", "sizes": "192x192", "type": "image/png"},
			{"src": "/icons/web-app-manifest-512x512.png", "sizes": "512x512", "type": "image/png"},
		},
	})
}

// PageOffline is cached by the service worker when it installs
func (h *DefaultHandler) PageOffline(ctx *gin.Context) {
	RenderTempl(ctx, pages.Offline())
}
