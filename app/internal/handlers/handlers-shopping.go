package handlers

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/routes"
	"roommates/internal/services"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateShoppingList opens the new list, an invalid name shows on the household page
func (h *DefaultHandler) CreateShoppingList(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	listID, err := h.service.CreateShoppingList(ctx, actor, ctx.PostForm("name"), ctx.PostForm("category"))
	if errors.Is(err, g.ErrorInvalidInput) {
		h.renderShopping(ctx, http.StatusUnprocessableEntity,
			utils.T(ctx.Request.Context(), locales.LK_ShoppingNameLength, services.MaxShoppingListNameLength))
		return
	}
	if err != nil {
		InternalError(ctx, err)
		return
	}
	Redirect(ctx, routes.URL(routes.ShoppingList, actor.HouseholdID, listID))
}

func (h *DefaultHandler) PageShoppingList(ctx *gin.Context) {
	listID, err := utils.ParseUUID(ctx.Param(routes.ParamListID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	h.renderShoppingList(ctx, http.StatusOK, listID, listForm{})
}

func (h *DefaultHandler) renderShoppingList(ctx *gin.Context, status int, listID pgtype.UUID, form listForm) {
	actor := membershipFromContext(ctx)
	list, err := h.service.GetShoppingList(ctx, actor, listID)
	if errors.Is(err, g.ErrorNotFound) {
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, status, pages.ShoppingList(actor, list, form.errorMessage, form.itemText, form.itemQty, form.categorySaved))
}

// listForm is what [DefaultHandler.renderShoppingList] shows besides the list itself
type listForm struct {
	errorMessage, itemText, itemQty string
	categorySaved                   bool
}

// shoppingAction runs an action on the :listID list and shows the list again with done,
// htmx requests get the page itself so they can swap in the list without a reload
func (h *DefaultHandler) shoppingAction(ctx *gin.Context, done listForm,
	action func(actor services.Membership, listID pgtype.UUID) error) {
	listID, err := utils.ParseUUID(ctx.Param(routes.ParamListID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	err = action(actor, listID)
	switch {
	case err == nil && IsRequestHTMX(ctx):
		h.renderShoppingList(ctx, http.StatusOK, listID, done)
	case err == nil:
		Redirect(ctx, routes.URL(routes.ShoppingList, actor.HouseholdID, listID))
	case errors.Is(err, g.ErrorInvalidInput):
		// only adding an item takes free text, the other actions post no text or qty
		h.renderShoppingList(ctx, http.StatusUnprocessableEntity, listID, listForm{
			errorMessage: utils.T(ctx.Request.Context(), locales.LK_ShoppingItemLength,
				services.MaxShoppingItemTextLength, services.MaxShoppingItemQtyLength),
			itemText:      ctx.PostForm("text"),
			itemQty:       ctx.PostForm("qty"),
			categorySaved: false,
		})
	case errors.Is(err, g.ErrorNotFound):
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
	default:
		InternalError(ctx, err)
	}
}

// shoppingItemAction is [DefaultHandler.shoppingAction] with :itemID parsed
func (h *DefaultHandler) shoppingItemAction(ctx *gin.Context,
	action func(actor services.Membership, listID, itemID pgtype.UUID) error) {
	itemID, err := utils.ParseUUID(ctx.Param(routes.ParamItemID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	h.shoppingAction(ctx, listForm{}, func(actor services.Membership, listID pgtype.UUID) error {
		err := action(actor, listID, itemID)
		if errors.Is(err, g.ErrorNotFound) {
			// deleted or cleared meanwhile by someone else, the fresh list shows that
			return nil
		}
		return err
	})
}

// DeleteShoppingList goes back to the shopping lists
func (h *DefaultHandler) DeleteShoppingList(ctx *gin.Context) {
	listID, err := utils.ParseUUID(ctx.Param(routes.ParamListID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	err = h.service.DeleteShoppingList(ctx, actor, listID)
	switch {
	case err == nil:
		Redirect(ctx, routes.URL(routes.Household, actor.HouseholdID))
	case errors.Is(err, g.ErrorNotFound):
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
	default:
		InternalError(ctx, err)
	}
}

// RenameShoppingList shows the list again, an invalid name shows on it
func (h *DefaultHandler) RenameShoppingList(ctx *gin.Context) {
	listID, err := utils.ParseUUID(ctx.Param(routes.ParamListID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	err = h.service.RenameShoppingList(ctx, actor, listID, ctx.PostForm("name"))
	switch {
	case err == nil:
		stay(ctx, routes.URL(routes.ShoppingList, actor.HouseholdID, listID))
	case errors.Is(err, g.ErrorInvalidInput):
		h.renderShoppingList(ctx, http.StatusUnprocessableEntity, listID, listForm{
			errorMessage:  utils.T(ctx.Request.Context(), locales.LK_ShoppingNameLength, services.MaxShoppingListNameLength),
			itemText:      "",
			itemQty:       "",
			categorySaved: false,
		})
	case errors.Is(err, g.ErrorNotFound):
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
	default:
		InternalError(ctx, err)
	}
}

// SetShoppingListCategory shows the list with a "saved" status next to the category
func (h *DefaultHandler) SetShoppingListCategory(ctx *gin.Context) {
	saved := listForm{errorMessage: "", itemText: "", itemQty: "", categorySaved: true}
	h.shoppingAction(ctx, saved, func(actor services.Membership, listID pgtype.UUID) error {
		return h.service.SetShoppingListCategory(ctx, actor, listID, ctx.PostForm("category"))
	})
}

func (h *DefaultHandler) AddShoppingItem(ctx *gin.Context) {
	h.shoppingAction(ctx, listForm{}, func(actor services.Membership, listID pgtype.UUID) error {
		return h.service.AddShoppingItem(ctx, actor, listID, ctx.PostForm("text"), ctx.PostForm("qty"))
	})
}

func (h *DefaultHandler) ClearCheckedShoppingItems(ctx *gin.Context) {
	h.shoppingAction(ctx, listForm{}, func(actor services.Membership, listID pgtype.UUID) error {
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

// DeleteShoppingItem goes back to the list, the item's own page is gone with it
func (h *DefaultHandler) DeleteShoppingItem(ctx *gin.Context) {
	listID, it, ok := h.loadShoppingItem(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	if err := h.service.DeleteShoppingItem(ctx, actor, listID, it.ID); err != nil {
		pageError(ctx, err)
		return
	}
	Redirect(ctx, routes.URL(routes.ShoppingList, actor.HouseholdID, listID))
}
