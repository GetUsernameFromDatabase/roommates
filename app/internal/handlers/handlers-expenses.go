package handlers

import (
	"cmp"
	"errors"
	"maps"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"
	"roommates/web/pages"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func expensesPath(householdID pgtype.UUID) string {
	return householdPath(householdID) + "/expenses"
}

func expensePath(householdID, expenseID pgtype.UUID) string {
	return expensesPath(householdID) + "/" + expenseID.String()
}

func (h *DefaultHandler) PageExpenses(ctx *gin.Context) {
	h.renderExpenses(ctx, http.StatusOK, "")
}

// renderExpenses shows balances, pending settlements and expenses, errorMessage is the outcome of a settlement
func (h *DefaultHandler) renderExpenses(ctx *gin.Context, status int, errorMessage string) {
	actor := membershipFromContext(ctx)
	filter := services.ExpenseFilter{Query: ctx.Query("q"), Show: ctx.Query("show")}
	o, err := h.service.GetExpensesOverview(ctx, actor, filter)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, status, pages.HouseholdExpenses(actor, o, filter, errorMessage))
}

// PageNewExpense starts empty, as a copy of ?from=<expenseID>, or for the bought items of ?list=<listID>,
// only ?item=<itemID> when given, their title and category come from the list
func (h *DefaultHandler) PageNewExpense(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	in := services.ExpenseInput{
		Title:      "",
		Amount:     "",
		Category:   g.ExpenseCategories[0],
		IncurredOn: time.Now().Format(services.DateLayout),
		Custom:     false,
		Members:    nil,
		Shares:     nil,
		ListID:     pgtype.UUID{},
		ItemID:     pgtype.UUID{},
	}
	if list := ctx.Query("list"); list != "" {
		var errList, errItem error
		in.ListID, errList = utils.ParseUUID(list)
		if item := ctx.Query("item"); item != "" {
			in.ItemID, errItem = utils.ParseUUID(item)
		}
		if errList != nil || errItem != nil {
			RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
			return
		}
		var err error
		if in.Title, in.Category, err = h.service.ShoppingExpenseDraft(ctx, actor, in.ListID, in.ItemID); err != nil {
			pageError(ctx, err)
			return
		}
	}
	if from := ctx.Query("from"); from != "" {
		expenseID, err := utils.ParseUUID(from)
		if err != nil {
			RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
			return
		}
		e, err := h.service.GetExpense(ctx, actor, expenseID)
		if err != nil {
			pageError(ctx, err)
			return
		}
		in = e.CopyInput()
	}
	h.renderExpenseEditor(ctx, http.StatusOK, in, "")
}

// renderExpenseEditor shows the expense form, an equal split with no one ticked ticks everyone
func (h *DefaultHandler) renderExpenseEditor(ctx *gin.Context, status int, in services.ExpenseInput, errorMessage string) {
	actor := membershipFromContext(ctx)
	members, err := h.service.ListMembers(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	// a copy can name people who have left since, the form has no box to untick them
	active := map[pgtype.UUID]bool{}
	for _, m := range members {
		active[m.UserID] = true
	}
	in.Members = slices.DeleteFunc(in.Members, func(id pgtype.UUID) bool { return !active[id] })
	maps.DeleteFunc(in.Shares, func(id pgtype.UUID, _ string) bool { return !active[id] })
	if len(in.Members) == 0 && !in.Custom && errorMessage == "" {
		for _, m := range members {
			in.Members = append(in.Members, m.UserID)
		}
	}
	RenderTemplStatus(ctx, status, pages.ExpenseEditor(actor, members, in, errorMessage))
}

// expenseFromForm reads the ticked "member" values and the typed shares named "share-<userID>"
func expenseFromForm(ctx *gin.Context) services.ExpenseInput {
	in := services.ExpenseInput{
		Title:      ctx.PostForm("title"),
		Amount:     ctx.PostForm("amount"),
		Category:   ctx.PostForm("category"),
		IncurredOn: ctx.PostForm("incurred_on"),
		Custom:     ctx.PostForm("split") == "custom",
		Members:    nil,
		Shares:     map[pgtype.UUID]string{},
		ListID:     pgtype.UUID{},
		ItemID:     pgtype.UUID{},
	}
	// a bad id only links nothing, the expense itself is what was asked for
	if id, err := utils.ParseUUID(ctx.PostForm("list")); err == nil {
		in.ListID = id
	}
	if id, err := utils.ParseUUID(ctx.PostForm("item")); err == nil {
		in.ItemID = id
	}
	for _, v := range ctx.PostFormArray("member") {
		if id, err := utils.ParseUUID(v); err == nil {
			in.Members = append(in.Members, id)
		}
	}
	for key, values := range ctx.Request.PostForm {
		if rest, ok := strings.CutPrefix(key, "share-"); ok {
			if id, err := utils.ParseUUID(rest); err == nil {
				in.Shares[id] = values[0]
			}
		}
	}
	return in
}

// CreateExpense opens the new expense
func (h *DefaultHandler) CreateExpense(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	in := expenseFromForm(ctx)
	expenseID, err := h.service.CreateExpense(ctx, actor, in)
	switch {
	case err == nil:
		Redirect(ctx, expensePath(actor.HouseholdID, expenseID))
	case errors.Is(err, g.ErrorInvalidInput):
		h.renderExpenseEditor(ctx, http.StatusUnprocessableEntity, in,
			utils.T(ctx.Request.Context(), locales.LK_ExpensesInvalid, services.MaxExpenseTitleLength))
	default:
		InternalError(ctx, err)
	}
}

// loadExpense gets the :expenseID expense, or shows the error page and returns false
func (h *DefaultHandler) loadExpense(ctx *gin.Context) (services.Expense, bool) {
	expenseID, err := utils.ParseUUID(ctx.Param(g.ParamExpenseID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return services.Expense{}, false
	}
	e, err := h.service.GetExpense(ctx, membershipFromContext(ctx), expenseID)
	if err != nil {
		pageError(ctx, err)
		return e, false
	}
	return e, true
}

func (h *DefaultHandler) PageExpense(ctx *gin.Context) {
	if e, ok := h.loadExpense(ctx); ok {
		RenderTempl(ctx, pages.ExpensePage(membershipFromContext(ctx), e, ""))
	}
}

func (h *DefaultHandler) UploadExpenseFile(ctx *gin.Context) {
	e, ok := h.loadExpense(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	h.handleUpload(ctx, services.FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: e.Expense.ID, MessageID: pgtype.UUID{}},
		expensePath(actor.HouseholdID, e.Expense.ID), func(status int, errorMessage string) {
			RenderTemplStatus(ctx, status, pages.ExpensePage(actor, e, errorMessage))
		})
}

func (h *DefaultHandler) DeleteExpense(ctx *gin.Context) {
	expenseID, err := utils.ParseUUID(ctx.Param(g.ParamExpenseID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	deleteExpense := h.service.DeleteExpense
	if ctx.PostForm("keep") != "" {
		deleteExpense = h.service.DeleteExpenseKeepingFiles
	}
	if err := deleteExpense(ctx, actor, expenseID); err != nil {
		pageError(ctx, err)
		return
	}
	Redirect(ctx, expensesPath(actor.HouseholdID))
}

// SetExpenseSharePaid takes the wanted state, not a toggle, so a stale page cannot undo it,
// it goes back to the form's "next" page, by default the expense
func (h *DefaultHandler) SetExpenseSharePaid(ctx *gin.Context) {
	expenseID, errExpense := utils.ParseUUID(ctx.Param(g.ParamExpenseID))
	userID, errUser := utils.ParseUUID(ctx.Param(g.ParamUserID))
	if errExpense != nil || errUser != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	if err := h.service.SetExpenseSharePaid(ctx, actor, expenseID, userID, ctx.PostForm("paid") == "1"); err != nil {
		pageError(ctx, err)
		return
	}
	Redirect(ctx, safeNextPath(cmp.Or(ctx.PostForm(g.QueryNext), expensePath(actor.HouseholdID, expenseID))))
}

// -----------------------------------------------------------------------------

// CreateSettlement records that the actor paid the "to" member
func (h *DefaultHandler) CreateSettlement(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	toUserID, err := utils.ParseUUID(ctx.PostForm("to"))
	if err != nil {
		err = g.ErrorInvalidInput
	} else {
		err = h.service.CreateSettlement(ctx, actor, toUserID, ctx.PostForm("amount"))
	}
	switch {
	case err == nil:
		Redirect(ctx, expensesPath(actor.HouseholdID))
	case errors.Is(err, g.ErrorInvalidInput):
		h.renderExpenses(ctx, http.StatusUnprocessableEntity, utils.T(ctx.Request.Context(), locales.LK_ExpensesSettleInvalid))
	default:
		InternalError(ctx, err)
	}
}

// settlementAction runs an action on the :settlementID settlement and goes back to the expenses
func (h *DefaultHandler) settlementAction(ctx *gin.Context, action func(actor services.Membership, settlementID pgtype.UUID) error) {
	settlementID, err := utils.ParseUUID(ctx.Param(g.ParamSettleID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	if err := action(actor, settlementID); err != nil {
		pageError(ctx, err)
		return
	}
	Redirect(ctx, expensesPath(actor.HouseholdID))
}

func (h *DefaultHandler) ConfirmSettlement(ctx *gin.Context) {
	h.settlementAction(ctx, func(actor services.Membership, settlementID pgtype.UUID) error {
		return h.service.ConfirmSettlement(ctx, actor, settlementID)
	})
}

func (h *DefaultHandler) DeleteSettlement(ctx *gin.Context) {
	h.settlementAction(ctx, func(actor services.Membership, settlementID pgtype.UUID) error {
		return h.service.DeleteSettlement(ctx, actor, settlementID)
	})
}
