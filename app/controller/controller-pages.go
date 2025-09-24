package controller

import (
	"roommates/components"
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

	var tc templ.Component
	if utils.IsRequestHTMX(ctx) {
		tc = components.MainPageContent(houses)
	} else {
		tc = components.PageMain(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		}, houses)
	}
	RenderTempl(ctx, tc)
}

func (c *Controller) PageProfile(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var tc templ.Component
	if utils.IsRequestHTMX(ctx) {
		tc = components.ProfilePageContent()
	} else {
		tc = components.PageProfile(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	}
	RenderTempl(ctx, tc)
}

func (c *Controller) PagePayments(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var tc templ.Component
	if utils.IsRequestHTMX(ctx) {
		tc = components.PaymentsPageContent()
	} else {
		tc = components.PagePayments(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	}
	RenderTempl(ctx, tc)
}

func (c *Controller) PageNotes(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var tc templ.Component
	if utils.IsRequestHTMX(ctx) {
		tc = components.NotesPageContent()
	} else {
		tc = components.PageNotes(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	}
	RenderTempl(ctx, tc)
}

func (c *Controller) PageMessaging(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)

	var tc templ.Component
	if utils.IsRequestHTMX(ctx) {
		tc = components.MessagingPageContent()
	} else {
		tc = components.PageMessaging(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		})
	}
	RenderTempl(ctx, tc)
}

func (c *Controller) PageHouses(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)
	houses, err := c.DB.UserHouses(ctx, authInfo.UserID)
	if err != nil {
		HandleServerError(ctx, err, "error getting houses")
		return
	}

	var tc templ.Component
	if utils.IsRequestHTMX(ctx) {
		tc = components.HousesPageContent(houses)
	} else {
		tc = components.PageHouses(components.SPageWrapper{
			AuthInfo: authInfo,
			PathURL:  ctx.Request.URL.Path,
		}, houses)
	}
	RenderTempl(ctx, tc)
}
