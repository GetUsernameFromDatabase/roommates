// Package handlers handle HTTP requests
package handlers

import (
	"net/http"
	"roommates/gintemplrenderer"
	g "roommates/globals"
	"roommates/logger"
	"roommates/services"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

var log = logger.Main.With().Str("package", "handlers").Logger()

// -----------------------------------------------------------------------------

type DefaultHandler struct {
	service *services.DefaultService
}

func NewDefaultHandler(service *services.DefaultService) *DefaultHandler {
	return &DefaultHandler{
		service: service,
	}
}

// -----------------------------------------------------------------------------

type BaseResponse struct {
	Success bool `json:"success"`
	//exhaustruct:optional
	Data any `json:"data,omitempty"`
	//exhaustruct:optional
	Error *ErrorInfo `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    g.ErrorCode `json:"code"`
	Message string      `json:"message"`
}

// -----------------------------------------------------------------------------

// OK json [BaseResponse] with
//
//	Success: true,
//	Data:    data,
func OK(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, BaseResponse{
		Success: true,
		Data:    data,
	})
}

// Fail json [BaseResponse] with
//
//	Success: false
//	Error: [ErrorInfo]
func Fail(ctx *gin.Context, httpStatus int, code g.ErrorCode) {
	// TODO: translate Message based on code
	errInfo := &ErrorInfo{Code: code, Message: code.String()}
	ctx.JSON(httpStatus, BaseResponse{Success: false, Error: errInfo})
}

// InternalError logs the error with stack and sends [Fail]
func InternalError(ctx *gin.Context, err error) {
	log.Error().Stack().Ctx(ctx).Err(err).Msg(g.EC_InternalError.String())
	Fail(ctx, http.StatusInternalServerError, g.EC_InternalError)
}

// RenderTempl renders templ component
func RenderTempl(ctx *gin.Context, component templ.Component) {
	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, component)
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

// GetHeader is a wrapper of [gin.Context.GetHeader]
func GetHeader(ctx *gin.Context, key g.HTTPHeaderKey) string {
	return ctx.GetHeader(key.String())
}
