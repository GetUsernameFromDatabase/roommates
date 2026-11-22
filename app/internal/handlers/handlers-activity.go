package handlers

import (
	"net/http"
	"roommates/internal/services"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
)

// PageActivity shows the newest activity, or with ?before=<[services.Cursor]> the rows after it,
// the "load more" button takes its rows from that page
func (h *DefaultHandler) PageActivity(ctx *gin.Context) {
	before, err := services.ParseCursor(ctx.Query("before"))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	rows, err := h.service.ListActivity(ctx, actor, before)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, http.StatusOK, pages.ActivityPage(actor, rows, len(rows) == services.ActivityPageSize))
}
