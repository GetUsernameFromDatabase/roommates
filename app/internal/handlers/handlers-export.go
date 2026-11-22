package handlers

import (
	"errors"
	"net/http"
	"time"

	g "roommates/internal/config"
	"roommates/internal/services"

	"github.com/gin-gonic/gin"
)

// ExportExpenses downloads the expenses as CSV, ?format= is wide or long
func (h *DefaultHandler) ExportExpenses(ctx *gin.Context) {
	data, err := h.service.ExportExpenses(ctx.Request.Context(), membershipFromContext(ctx), ctx.DefaultQuery("format", services.ExportWide))
	sendCSV(ctx, "expenses", data, err)
}

// ExportBalances downloads who is owed what as CSV
func (h *DefaultHandler) ExportBalances(ctx *gin.Context) {
	data, err := h.service.ExportBalances(ctx.Request.Context(), membershipFromContext(ctx))
	sendCSV(ctx, "balances", data, err)
}

func sendCSV(ctx *gin.Context, name string, data []byte, err error) {
	switch {
	case errors.Is(err, g.ErrorInvalidInput):
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	case err != nil:
		InternalError(ctx, err)
		return
	}
	ctx.Header("Content-Disposition", `attachment; filename="`+name+"-"+time.Now().Format(time.DateOnly)+`.csv"`)
	ctx.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}
