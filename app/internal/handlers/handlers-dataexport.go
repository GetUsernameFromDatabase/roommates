package handlers

import (
	"net/http"

	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
)

// RequestDataExport queues a ZIP of the user's data, emailed in the request language when
// it is built, and swaps only the data export section
func (h *DefaultHandler) RequestDataExport(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	path, _ := settingsPath(ctx)
	c := ctx.Request.Context()
	if err := h.service.RequestDataExport(ctx, userID, utils.GetLocaleString(c)); err != nil {
		InternalError(ctx, err)
		return
	}
	if !IsRequestHTMX(ctx) {
		ctx.Redirect(http.StatusSeeOther, path)
		return
	}
	status, err := h.service.GetDataExportStatus(ctx, userID)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTempl(ctx, pages.DataExportSection(path, status, utils.T(c, locales.LK_DataExportRequested)))
}

// DownloadDataExport streams the user's own ready export, 404 for anything else
func (h *DefaultHandler) DownloadDataExport(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	exportID, err := utils.ParseUUID(ctx.Param(g.ParamExportID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	body, size, readyAt, err := h.service.OpenDataExport(ctx, userID, exportID)
	if err != nil {
		pageError(ctx, err)
		return
	}
	defer body.Close() //nolint:errcheck // read only
	ctx.DataFromReader(http.StatusOK, size, "application/zip", body, map[string]string{
		"Content-Disposition": `attachment; filename="roommates-export-` + readyAt.Format("2006-01-02") + `.zip"`,
		"Cache-Control":       "private, no-store",
	})
}
