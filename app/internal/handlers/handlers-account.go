package handlers

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/models"
	"roommates/internal/services"
	"roommates/internal/sessions"
	"roommates/internal/utils"
	"roommates/web/pages"
	"strings"

	"github.com/gin-gonic/gin"
)

// renderAccount swaps [pages.AccountSection] in place, errorAt and infoAt are a pages.AccountAt*
func (h *DefaultHandler) renderAccount(ctx *gin.Context, status int, errorAt, errorMessage, infoAt, infoMessage string) {
	path, _ := settingsPath(ctx)
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, status, pages.AccountSection(path, user, errorAt, errorMessage, infoAt, infoMessage))
}

// accountError shows err by the form at, if it is one the user can fix, otherwise it is an internal error
func (h *DefaultHandler) accountError(ctx *gin.Context, at string, err error) {
	keys := map[error]locales.LK{
		g.ErrorInvalidCredential:     locales.LK_AccountWrongPassword,
		services.ErrorEmailUnchanged: locales.LK_AccountEmailSame,
		g.ErrorAccountAlreadyExists:  locales.LK_RegistrationPageAccountExists,
	}
	for e, key := range keys {
		if errors.Is(err, e) {
			h.renderAccount(ctx, http.StatusUnprocessableEntity, at, utils.T(ctx.Request.Context(), key), "", "")
			return
		}
	}
	InternalError(ctx, err)
}

// invalidAccountForm shows the validation messages of the form at, one per line
func (h *DefaultHandler) invalidAccountForm(ctx *gin.Context, at string, msgs []locales.LKMessage) {
	lines := make([]string, len(msgs))
	for i, m := range msgs {
		lines[i] = utils.T(ctx.Request.Context(), m.Key, m.Args...)
	}
	h.renderAccount(ctx, http.StatusUnprocessableEntity, at, strings.Join(lines, "\n"), "", "")
}

// ChangeUsername saves the username of the signed in user
func (h *DefaultHandler) ChangeUsername(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	var model models.ChangeUsername
	_ = ctx.ShouldBind(&model)
	if ok, msgs := model.IsValid(); !ok {
		h.invalidAccountForm(ctx, pages.AccountAtUsername, msgs)
		return
	}
	if err := h.service.SetUsername(ctx, userID, model.Username); err != nil {
		InternalError(ctx, err)
		return
	}
	s := sessions.FromGinContext(ctx)
	s.Session.Set(sessions.KeyUsername, model.Username)
	s.Save()
	h.renderAccount(ctx, http.StatusOK, "", "", pages.AccountAtUsername,
		utils.T(ctx.Request.Context(), locales.LK_CommonSaved))
}

// ChangePassword sets a new password and signs every other device out
func (h *DefaultHandler) ChangePassword(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	var model models.ChangePassword
	_ = ctx.ShouldBind(&model)
	if ok, msgs := model.IsValid(); !ok {
		h.invalidAccountForm(ctx, pages.AccountAtPassword, msgs)
		return
	}
	if err := h.service.ChangePassword(ctx, userID, model.Current, model.Password); err != nil {
		h.accountError(ctx, pages.AccountAtPassword, err)
		return
	}
	// every session was deleted, this one comes back under a new id
	s := sessions.FromGinContext(ctx)
	s.Login(s.SessionInfo)
	h.renderAccount(ctx, http.StatusOK, "", "", pages.AccountAtPassword,
		utils.T(ctx.Request.Context(), locales.LK_AccountPasswordChanged))
}

// ChangeEmail sends a confirmation link to the new email, it changes once the link is opened
func (h *DefaultHandler) ChangeEmail(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	var model models.ChangeEmail
	_ = ctx.ShouldBind(&model)
	if ok, msgs := model.IsValid(); !ok {
		h.invalidAccountForm(ctx, pages.AccountAtEmail, msgs)
		return
	}
	// the request context carries the language of the email
	if err := h.service.RequestEmailChange(ctx.Request.Context(), userID, model.Current, model.Email); err != nil {
		h.accountError(ctx, pages.AccountAtEmail, err)
		return
	}
	h.renderAccount(ctx, http.StatusOK, "", "", pages.AccountAtEmail,
		utils.T(ctx.Request.Context(), locales.LK_AccountEmailSent, services.NormalizeEmail(model.Email)))
}

// SignOutEverywhere ends every session of the user, this one too
func (h *DefaultHandler) SignOutEverywhere(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	if err := h.service.SignOutEverywhere(ctx, userID); err != nil {
		InternalError(ctx, err)
		return
	}
	s := sessions.FromGinContext(ctx)
	s.Logout()
	reload(ctx, g.RouteLogin)
}

// ConfirmEmailChange makes the email of the ?token= link the user's, signed in or not
func (h *DefaultHandler) ConfirmEmailChange(ctx *gin.Context) {
	err := h.service.ConfirmEmailChange(ctx.Request.Context(), ctx.Query("token"))
	switch {
	case errors.Is(err, g.ErrorTokenInvalid):
		notice(ctx, http.StatusNotFound, locales.LK_CommonLinkInvalid,
			locales.LK_ConfirmEmailInvalidDescription, "/", locales.LK_VerifyEmailContinue)
	case errors.Is(err, g.ErrorAccountAlreadyExists):
		notice(ctx, http.StatusConflict, locales.LK_ConfirmEmailTakenTitle,
			locales.LK_ConfirmEmailTakenDescription, "/", locales.LK_VerifyEmailContinue)
	case err != nil:
		InternalError(ctx, err)
	default:
		notice(ctx, http.StatusOK, locales.LK_ConfirmEmailDoneTitle,
			locales.LK_ConfirmEmailDoneDescription, "/", locales.LK_VerifyEmailContinue)
	}
}
