package handlers

import (
	"net/http"
	"roommates/internal/services"
	"roommates/web/pages"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// PageActivity shows the newest activity, or with ?before=<RFC 3339 time> the rows older than it,
// the "load more" button takes its rows from that page
func (h *DefaultHandler) PageActivity(ctx *gin.Context) {
	var before pgtype.Timestamptz
	if s := ctx.Query("before"); s != "" {
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
			return
		}
		before = pgtype.Timestamptz{Time: t, Valid: true, InfinityModifier: pgtype.Finite}
	}
	actor := membershipFromContext(ctx)
	rows, err := h.service.ListActivity(ctx, actor, before)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, http.StatusOK, pages.ActivityPage(actor, rows, len(rows) == services.ActivityPageSize))
}
