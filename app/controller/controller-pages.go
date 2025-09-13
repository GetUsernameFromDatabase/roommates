package controller

import (
	"net/http"
	"roommates/components"
	"roommates/gintemplrenderer"
	g "roommates/globals"
	"roommates/middleware"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func (c *Controller) PageMain(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if ctx.GetHeader(string(g.HHXRequest)) == "" {
		pc = components.PageMain(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	} else {
		pc = components.MainPageContent()
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageProfile(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if ctx.GetHeader(string(g.HHXRequest)) == "" {
		pc = components.PageProfile(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	} else {
		pc = components.ProfilePageContent()
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PagePayments(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if ctx.GetHeader(string(g.HHXRequest)) == "" {
		pc = components.PagePayments(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	} else {
		pc = components.PaymentsPageContent()
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageNotes(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if ctx.GetHeader(string(g.HHXRequest)) == "" {
		pc = components.PageNotes(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	} else {
		pc = components.NotesPageContent()
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageMessaging(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if ctx.GetHeader(string(g.HHXRequest)) == "" {
		pc = components.PageMessaging(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	} else {
		pc = components.MessagingPageContent()
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageHouses(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if ctx.GetHeader(string(g.HHXRequest)) == "" {
		pc = components.PageHouses(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	} else {
		pc = components.HousesPageContent()
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}
