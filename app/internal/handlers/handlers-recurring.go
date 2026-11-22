package handlers

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/routes"
	"roommates/internal/services"
	"roommates/internal/utils"
	"roommates/web/pages"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// PageRecurring lists the household's schedules
func (h *DefaultHandler) PageRecurring(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	rows, err := h.service.ListRecurring(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTempl(ctx, pages.HouseholdRecurring(actor, rows))
}

// loadRecurring gets schedule id, or shows the error page and returns false
func (h *DefaultHandler) loadRecurring(ctx *gin.Context, id string) (services.RecurringRow, bool) {
	recurringID, err := utils.ParseUUID(id)
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return services.RecurringRow{}, false
	}
	r, err := h.service.GetRecurring(ctx, membershipFromContext(ctx), recurringID)
	if err != nil {
		pageError(ctx, err)
		return r, false
	}
	return r, true
}

// PageEditRecurring shows the :recurringID schedule's form, only to its payer and owners
func (h *DefaultHandler) PageEditRecurring(ctx *gin.Context) {
	r, ok := h.loadRecurring(ctx, ctx.Param(routes.ParamRecurringID))
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	if r.PayerID != actor.UserID && actor.Role != g.RoleOwner {
		pageError(ctx, g.ErrorNotFound)
		return
	}
	h.renderRecurringEditor(ctx, http.StatusOK, r, r.Input(), "", "")
}

func (h *DefaultHandler) renderRecurringEditor(ctx *gin.Context, status int, r services.RecurringRow, in services.RecurringInput, errorMessage, errorField string) {
	actor := membershipFromContext(ctx)
	members, err := h.service.ListMembers(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, status, pages.RecurringEditor(actor, members, r, in, errorMessage, errorField))
}

// UpdateRecurring saves the schedule form and goes back to the list
func (h *DefaultHandler) UpdateRecurring(ctx *gin.Context) {
	r, ok := h.loadRecurring(ctx, ctx.Param(routes.ParamRecurringID))
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	in := services.RecurringInput{ExpenseInput: expenseFromForm(ctx), Interval: ctx.PostForm("repeat")}
	err := h.service.UpdateRecurring(ctx, actor, r.ID, in)
	switch {
	case err == nil:
		Redirect(ctx, routes.URL(routes.Recurring, actor.HouseholdID))
	case errors.Is(err, g.ErrorInvalidInput):
		message, field := expenseInputError(ctx, err)
		h.renderRecurringEditor(ctx, http.StatusUnprocessableEntity, r, in, message, field)
	default:
		pageError(ctx, err)
	}
}

func (h *DefaultHandler) DeleteRecurring(ctx *gin.Context) {
	recurringID, err := utils.ParseUUID(ctx.Param(routes.ParamRecurringID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	if err := h.service.DeleteRecurring(ctx, actor, recurringID); err != nil {
		pageError(ctx, err)
		return
	}
	Redirect(ctx, routes.URL(routes.Recurring, actor.HouseholdID))
}

// SkipRecurring moves the schedule past the posted "due" day and shows the expenses again,
// a day already added or skipped elsewhere only shows them
func (h *DefaultHandler) SkipRecurring(ctx *gin.Context) {
	recurringID, err := utils.ParseUUID(ctx.Param(routes.ParamRecurringID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	var due pgtype.Date
	if day, err := time.Parse(services.DateLayout, ctx.PostForm("due")); err == nil {
		due = pgtype.Date{Time: day, InfinityModifier: pgtype.Finite, Valid: true}
	}
	if err := h.service.SkipRecurring(ctx, actor, recurringID, due); err != nil && !errors.Is(err, g.ErrorInvalidInput) {
		pageError(ctx, err)
		return
	}
	stay(ctx, routes.URL(routes.Expenses, actor.HouseholdID))
}
