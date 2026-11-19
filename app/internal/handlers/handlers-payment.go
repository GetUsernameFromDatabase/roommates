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
)

// SetPaymentDetails saves the user's account holder and IBAN from the settings page,
// at /settings or in a household, and swaps only the payment section
func (h *DefaultHandler) SetPaymentDetails(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	path, _ := settingsPath(ctx)
	name, iban := ctx.PostForm("name"), ctx.PostForm("iban")
	err = h.service.SetPaymentDetails(ctx, userID, name, iban)

	c := ctx.Request.Context()
	var fieldErr services.FieldError
	switch {
	case err == nil:
		if !IsRequestHTMX(ctx) {
			ctx.Redirect(http.StatusSeeOther, path)
			return
		}
		d, getErr := h.service.GetPaymentDetails(ctx, userID)
		if getErr != nil {
			InternalError(ctx, getErr)
			return
		}
		show, showErr := h.paymentShow(ctx)
		if showErr != nil {
			InternalError(ctx, showErr)
			return
		}
		RenderTempl(ctx, pages.PaymentSection(path, d, show, "", utils.T(c, locales.LK_CommonSaved)))
	case errors.As(err, &fieldErr):
		if !IsRequestHTMX(ctx) {
			ctx.Redirect(http.StatusSeeOther, path)
			return
		}
		msg := utils.T(c, locales.LK_PaymentIbanInvalid)
		if fieldErr.Field == services.FieldPayeeName {
			msg = utils.T(c, locales.LK_PaymentNameInvalid)
		}
		show, showErr := h.paymentShow(ctx)
		if showErr != nil {
			InternalError(ctx, showErr)
			return
		}
		// what was typed stays in the fields
		RenderTemplStatus(ctx, http.StatusUnprocessableEntity,
			pages.PaymentSection(path, services.PaymentDetails{Name: name, IBAN: iban}, show, msg, ""))
	default:
		InternalError(ctx, err)
	}
}

// paymentShow is the household switch of the route's membership, nil outside a household
func (h *DefaultHandler) paymentShow(ctx *gin.Context) (*bool, error) {
	if _, ok := ctx.Get(g.GinKeyMembership); !ok {
		return nil, nil
	}
	show, err := h.service.GetShowPaymentDetails(ctx, membershipFromContext(ctx))
	return &show, err
}

// SetShowPaymentDetails switches whether the user's details show in the household of the route
func (h *DefaultHandler) SetShowPaymentDetails(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	path, _ := settingsPath(ctx)
	show := ctx.PostForm("show") == "1"
	if err := h.service.SetShowPaymentDetails(ctx, actor, show); err != nil {
		InternalError(ctx, err)
		return
	}
	if !IsRequestHTMX(ctx) {
		ctx.Redirect(http.StatusSeeOther, path)
		return
	}
	d, err := h.service.GetPaymentDetails(ctx, userID)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTempl(ctx, pages.PaymentSwitch(path, show, d.IBAN != ""))
}
