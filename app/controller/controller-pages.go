package controller

import (
	"net/http"
	"roommates/components"
	"roommates/gintemplrenderer"
	"roommates/middleware"
	"roommates/utils"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func (c *Controller) PageMain(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)
	houses, err := c.DB.UserHouses(ctx, authInfo.UserID)
	if err != nil {
		HandleServerError(ctx, err, "error getting houses")
		return
	}

	var pc templ.Component
	if utils.IsRequestHTMX(ctx) {
		pc = components.MainPageContent(houses)
	} else {
		pc = components.PageMain(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		}, houses)
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageProfile(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if utils.IsRequestHTMX(ctx) {
		pc = components.ProfilePageContent()
	} else {
		pc = components.PageProfile(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PagePayments(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if utils.IsRequestHTMX(ctx) {
		pc = components.PaymentsPageContent()
	} else {
		pc = components.PagePayments(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageNotes(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if utils.IsRequestHTMX(ctx) {
		pc = components.NotesPageContent()
	} else {
		pc = components.PageNotes(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageMessaging(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var pc templ.Component
	if utils.IsRequestHTMX(ctx) {
		pc = components.MessagingPageContent()
	} else {
		pc = components.PageMessaging(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageHouses(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)
	houses, err := c.DB.UserHouses(ctx, authInfo.UserID)
	if err != nil {
		HandleServerError(ctx, err, "error getting houses")
		return
	}

	var pc templ.Component
	if utils.IsRequestHTMX(ctx) {
		pc = components.HousesPageContent(houses)
	} else {
		pc = components.PageHouses(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		}, houses)
	}

	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}
