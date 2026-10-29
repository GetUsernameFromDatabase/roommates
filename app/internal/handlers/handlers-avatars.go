package handlers

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"

	"github.com/gin-gonic/gin"
)

// SetAvatar makes the "avatar" form field the user's picture and goes back to the settings,
// the body is cut off a little above [services.MaxFileSize] like [DefaultHandler.uploadFile]
func (h *DefaultHandler) SetAvatar(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, services.MaxFileSize+1<<20)
	err = g.ErrorInvalidInput
	var tooLarge *http.MaxBytesError
	header, formErr := ctx.FormFile("avatar")
	switch {
	case errors.As(formErr, &tooLarge):
		err = g.ErrorFileTooLarge
	case formErr == nil:
		f, openErr := header.Open()
		if openErr != nil {
			InternalError(ctx, openErr)
			return
		}
		defer f.Close() //nolint:errcheck // read only
		err = h.service.SetAvatar(ctx, userID, f, header.Size)
	}
	c := ctx.Request.Context()
	switch {
	case err == nil:
		path, _ := settingsPath(ctx)
		Redirect(ctx, path)
	case errors.Is(err, g.ErrorFileTooLarge):
		h.renderSettings(ctx, http.StatusRequestEntityTooLarge, utils.T(c, locales.LK_FilesTooLarge, services.MaxFileSize>>20))
	case errors.Is(err, g.ErrorInvalidInput):
		h.renderSettings(ctx, http.StatusUnprocessableEntity, utils.T(c, locales.LK_AvatarInvalid))
	default:
		InternalError(ctx, err)
	}
}

func (h *DefaultHandler) RemoveAvatar(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	if err := h.service.RemoveAvatar(ctx, userID); err != nil {
		InternalError(ctx, err)
		return
	}
	path, _ := settingsPath(ctx)
	Redirect(ctx, path)
}

// DownloadAvatar sends the :userID user's picture to them and the people they live with,
// a new picture has a new ?v= so it is cached for good
func (h *DefaultHandler) DownloadAvatar(ctx *gin.Context) {
	viewerID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	userID, err := utils.ParseUUID(ctx.Param(g.ParamUserID))
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	body, err := h.service.OpenAvatar(ctx, viewerID, userID)
	switch {
	case errors.Is(err, g.ErrorNotFound):
		ctx.Status(http.StatusNotFound)
		return
	case err != nil:
		InternalError(ctx, err)
		return
	}
	defer body.Close() //nolint:errcheck // read only
	ctx.DataFromReader(http.StatusOK, -1, "image/jpeg", body, map[string]string{
		"Cache-Control": "private, max-age=31536000, immutable",
	})
}
