// Package handlers handle HTTP requests
package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/events"
	"roommates/internal/logger"
	"roommates/internal/services"
	"roommates/web/pages"
	"roommates/web/renderer"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var log = logger.Main.With().Str("package", "handlers").Logger()

// -----------------------------------------------------------------------------

type DefaultHandler struct {
	service *services.DefaultService
	hub     *events.Hub
}

func NewDefaultHandler(service *services.DefaultService, hub *events.Hub) *DefaultHandler {
	return &DefaultHandler{
		service: service,
		hub:     hub,
	}
}

// -----------------------------------------------------------------------------

// InternalError logs the error with stack and renders [pages.Error] with status 500
func InternalError(ctx *gin.Context, err error) {
	// WithStack records the stack here, errors from pgx, bcrypt etc. carry none
	log.Error().Stack().Ctx(ctx).Err(errors.WithStack(err)).Msg(g.EC_InternalError.String())
	swapMainOnError(ctx)
	RenderTemplStatus(ctx, http.StatusInternalServerError, pages.Error())
	ctx.Abort()
}

// RenderTempl renders templ component with status 200
func RenderTempl(ctx *gin.Context, component templ.Component) {
	RenderTemplStatus(ctx, http.StatusOK, component)
}

// RenderTemplStatus renders templ component with given status
func RenderTemplStatus(ctx *gin.Context, status int, component templ.Component) {
	// a boosted form would push its action URL: after an error a reload must show the form again
	if status >= http.StatusBadRequest && IsRequestHTMX(ctx) {
		ctx.Header(g.HHK_HXPushURL.String(), "false")
	}
	r := renderer.New(ctx.Request.Context(), status, component)
	ctx.Render(r.Status, r)
}

// IsRequestHTMX checks if request is made by htmx by looking at [HX-Request] header
//
// [HX-Request]: https://four.htmx.org/reference/headers/HX-Request
func IsRequestHTMX(ctx *gin.Context) bool {
	return ctx.GetHeader(g.HHK_HXRequest.String()) != ""
}

// Redirect navigates to location: htmx swaps the new page into body, pushes its URL and scrolls
// to the top like a boosted link, other requests get see other
//
//	if location == "" then ctx.Request.Referer() will be used
func Redirect(ctx *gin.Context, location string) {
	hxLocation(ctx, location, "innerHTML show:top", "true")
}

// reload is [Redirect] with a full page load for htmx, for changes outside <body>
// such as <html lang>, the theme or the session
func reload(ctx *gin.Context, location string) {
	if location == "" {
		location = ctx.Request.Referer()
	}
	if IsRequestHTMX(ctx) {
		ctx.Header(g.HHK_HXRedirect.String(), location)
	} else {
		ctx.Redirect(http.StatusSeeOther, location)
	}
}

// stay sends an action back to its own page, htmx swaps the page in like a boosted link
// where it was scrolled, as HX-Redirect would reload it
func stay(ctx *gin.Context, location string) {
	hxLocation(ctx, location, "innerHTML show:none", "false")
}

// backTo is [stay] when next is the page the htmx form is on, keeping its ?query such as a filter,
// otherwise [Redirect] to next
func backTo(ctx *gin.Context, next string) {
	if u, err := url.Parse(ctx.GetHeader("HX-Current-URL")); err == nil && IsRequestHTMX(ctx) && u.Path == next {
		stay(ctx, u.RequestURI())
		return
	}
	Redirect(ctx, next)
}

func hxLocation(ctx *gin.Context, location, swap, push string) {
	if location == "" {
		location = ctx.Request.Referer()
	}
	if !IsRequestHTMX(ctx) {
		ctx.Redirect(http.StatusSeeOther, location)
		return
	}
	header, _ := json.Marshal(map[string]string{
		"path": location, "target": "body", "swap": swap, "push": push,
	})
	ctx.Header(g.HHK_HXLocation.String(), string(header))
}

// swapMainOnError makes htmx fill the page's <main> with the error page's message,
// as forms that select or target a part of the page would miss it on an error page,
// live-resync.js does the same for a live region that gets a 404
func swapMainOnError(ctx *gin.Context) {
	if IsRequestHTMX(ctx) {
		ctx.Header(g.HHK_HXRetarget.String(), "main")
		ctx.Header(g.HHK_HXReselect.String(), "main > *")
		ctx.Header(g.HHK_HXReswap.String(), "innerHTML")
	}
}

// stayInPlace keeps the scroll position of a boosted page swap, which would go to the top,
// so an outcome shown by its form stays in view
func stayInPlace(ctx *gin.Context) {
	ctx.Header(g.HHK_HXReswap.String(), "innerHTML show:none")
}

// GetHeader is a wrapper of [gin.Context.GetHeader]
func GetHeader(ctx *gin.Context, key g.HTTPHeaderKey) string {
	return ctx.GetHeader(key.String())
}
