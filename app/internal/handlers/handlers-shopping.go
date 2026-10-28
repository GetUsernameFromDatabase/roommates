package handlers

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func shoppingListPath(householdID, listID pgtype.UUID) string {
	return householdPath(householdID) + "/lists/" + listID.String()
}

// CreateShoppingList opens the new list, an invalid name shows on the household page
func (h *DefaultHandler) CreateShoppingList(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	listID, err := h.service.CreateShoppingList(ctx, actor, ctx.PostForm("name"))
	if errors.Is(err, g.ErrorInvalidInput) {
		h.renderHousehold(ctx, http.StatusUnprocessableEntity, "",
			utils.T(ctx.Request.Context(), locales.LK_ShoppingNameLength, services.MaxShoppingListNameLength))
		return
	}
	if err != nil {
		InternalError(ctx, err)
		return
	}
	Redirect(ctx, shoppingListPath(actor.HouseholdID, listID))
}

func (h *DefaultHandler) PageShoppingList(ctx *gin.Context) {
	listID, err := utils.ParseUUID(ctx.Param(g.ParamListID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	h.renderShoppingList(ctx, http.StatusOK, listID, "")
}

func (h *DefaultHandler) renderShoppingList(ctx *gin.Context, status int, listID pgtype.UUID, errorMessage string) {
	actor := membershipFromContext(ctx)
	list, err := h.service.GetShoppingList(ctx, actor, listID)
	if errors.Is(err, g.ErrorNotFound) {
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	if err != nil {
		swapMainOnError(ctx)
		InternalError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, status, pages.ShoppingList(actor, list, errorMessage))
}

// swapMainOnError makes htmx swap in the <main> of an error page,
// the list forms select #shopping-list, which error pages do not have
func swapMainOnError(ctx *gin.Context) {
	if IsRequestHTMX(ctx) {
		ctx.Header(g.HHK_HXRetarget.String(), "main")
		ctx.Header(g.HHK_HXReselect.String(), "main")
		ctx.Header(g.HHK_HXReswap.String(), "outerHTML")
	}
}

// shoppingAction runs an action on the :listID list and shows the list again,
// htmx requests get the page itself so they can swap in the list without a reload
func (h *DefaultHandler) shoppingAction(ctx *gin.Context, action func(actor services.Membership, listID pgtype.UUID) error) {
	listID, err := utils.ParseUUID(ctx.Param(g.ParamListID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	err = action(actor, listID)
	switch {
	case err == nil && IsRequestHTMX(ctx):
		h.renderShoppingList(ctx, http.StatusOK, listID, "")
	case err == nil:
		Redirect(ctx, shoppingListPath(actor.HouseholdID, listID))
	case errors.Is(err, g.ErrorInvalidInput):
		h.renderShoppingList(ctx, http.StatusUnprocessableEntity, listID,
			utils.T(ctx.Request.Context(), locales.LK_ShoppingItemLength,
				services.MaxShoppingItemTextLength, services.MaxShoppingItemQtyLength))
	case errors.Is(err, g.ErrorNotFound):
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
	default:
		swapMainOnError(ctx)
		InternalError(ctx, err)
	}
}

// shoppingItemAction is [DefaultHandler.shoppingAction] with :itemID parsed
func (h *DefaultHandler) shoppingItemAction(ctx *gin.Context,
	action func(actor services.Membership, listID, itemID pgtype.UUID) error) {
	itemID, err := utils.ParseUUID(ctx.Param(g.ParamItemID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	h.shoppingAction(ctx, func(actor services.Membership, listID pgtype.UUID) error {
		return action(actor, listID, itemID)
	})
}

// DeleteShoppingList goes back to the household page
func (h *DefaultHandler) DeleteShoppingList(ctx *gin.Context) {
	h.householdTargetAction(ctx, g.ParamListID, func(actor services.Membership, listID pgtype.UUID) error {
		return h.service.DeleteShoppingList(ctx, actor, listID)
	})
}

func (h *DefaultHandler) AddShoppingItem(ctx *gin.Context) {
	h.shoppingAction(ctx, func(actor services.Membership, listID pgtype.UUID) error {
		return h.service.AddShoppingItem(ctx, actor, listID, ctx.PostForm("text"), ctx.PostForm("qty"))
	})
}

func (h *DefaultHandler) ClearCheckedShoppingItems(ctx *gin.Context) {
	h.shoppingAction(ctx, func(actor services.Membership, listID pgtype.UUID) error {
		return h.service.ClearCheckedShoppingItems(ctx, actor, listID)
	})
}

// SetShoppingItemClaim takes the wanted state, not a toggle, so a stale page cannot undo others
func (h *DefaultHandler) SetShoppingItemClaim(ctx *gin.Context) {
	h.shoppingItemAction(ctx, func(actor services.Membership, listID, itemID pgtype.UUID) error {
		return h.service.SetShoppingItemClaim(ctx, actor, listID, itemID, ctx.PostForm("claim") == "1")
	})
}

func (h *DefaultHandler) SetShoppingItemChecked(ctx *gin.Context) {
	h.shoppingItemAction(ctx, func(actor services.Membership, listID, itemID pgtype.UUID) error {
		return h.service.SetShoppingItemChecked(ctx, actor, listID, itemID, ctx.PostForm("checked") == "1")
	})
}

func (h *DefaultHandler) DeleteShoppingItem(ctx *gin.Context) {
	h.shoppingItemAction(ctx, func(actor services.Membership, listID, itemID pgtype.UUID) error {
		return h.service.DeleteShoppingItem(ctx, actor, listID, itemID)
	})
}
