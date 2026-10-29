// Package handlers handle HTTP requests
package handlers

import (
	"net/http"
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
	RenderTemplStatus(ctx, http.StatusInternalServerError, pages.Error())
	ctx.Abort()
}

// RenderTempl renders templ component with status 200
func RenderTempl(ctx *gin.Context, component templ.Component) {
	RenderTemplStatus(ctx, http.StatusOK, component)
}

// RenderTemplStatus renders templ component with given status
func RenderTemplStatus(ctx *gin.Context, status int, component templ.Component) {
	r := renderer.New(ctx.Request.Context(), status, component)
	ctx.Render(r.Status, r)
}

// IsRequestHTMX checks if request is made by htmx by looking at [HX-Request] header
//
// [HX-Request]: https://four.htmx.org/reference/headers/HX-Request
func IsRequestHTMX(ctx *gin.Context) bool {
	return ctx.GetHeader(g.HHK_HXRequest.String()) != ""
}

// Redirect htmx request with header, otherwise uses [gin.Context.Redirect] with status see other
//
//	if location == "" then ctx.Request.Referer() will be used
func Redirect(ctx *gin.Context, location string) {
	if location == "" {
		location = ctx.Request.Referer()
	}

	if IsRequestHTMX(ctx) {
		ctx.Header(g.HHK_HXRedirect.String(), location)
	} else {
		ctx.Redirect(http.StatusSeeOther, location)
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
