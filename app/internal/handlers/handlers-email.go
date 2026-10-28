package handlers

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/models"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// notice renders [pages.Notice] with translated texts
func notice(ctx *gin.Context, status int, title, description locales.LK, href string, linkText locales.LK, args ...any) {
	c := ctx.Request.Context()
	RenderTemplStatus(ctx, status, pages.Notice(
		utils.T(c, title), utils.T(c, description, args...),
		templ.SafeURL(href), utils.T(c, linkText),
	))
}

// VerifyEmail confirms the email of the ?token= link, signed in or not
func (h *DefaultHandler) VerifyEmail(ctx *gin.Context) {
	err := h.service.VerifyEmail(ctx, ctx.Query("token"))
	if errors.Is(err, g.ErrorTokenInvalid) {
		notice(ctx, http.StatusNotFound, locales.LK_VerifyEmailInvalidTitle,
			locales.LK_VerifyEmailInvalidDescription, "/", locales.LK_VerifyEmailContinue)
		return
	}
	if err != nil {
		InternalError(ctx, err)
		return
	}
	notice(ctx, http.StatusOK, locales.LK_VerifyEmailDoneTitle,
		locales.LK_VerifyEmailDoneDescription, "/", locales.LK_VerifyEmailContinue)
}

// ResendVerificationEmail sends a new link to the signed in user
func (h *DefaultHandler) ResendVerificationEmail(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	// the request context carries the language of the email
	if err := h.service.SendVerificationEmail(ctx.Request.Context(), userID); err != nil {
		InternalError(ctx, err)
		return
	}
	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	notice(ctx, http.StatusOK, locales.LK_VerifyEmailResentTitle,
		locales.LK_VerifyEmailResentDescription, "/", locales.LK_VerifyEmailContinue, user.Email)
}

// PageForgotPassword sends a reset link, the answer is the same whether the account exists or not
func (h *DefaultHandler) PageForgotPassword(ctx *gin.Context) {
	//exhaustruct:ignore
	model := models.EmailForm{ModelBase: models.ModelBase{Initial: true}}
	if ctx.Request.Method != http.MethodPost {
		RenderTempl(ctx, pages.ForgotPassword(model))
		return
	}

	_ = ctx.ShouldBind(&model)
	if isValid, _ := model.IsValid(); !isValid {
		RenderTempl(ctx, pages.ForgotPassword(model))
		return
	}
	if err := h.service.RequestPasswordReset(ctx.Request.Context(), model.Email); err != nil {
		InternalError(ctx, err)
		return
	}
	notice(ctx, http.StatusOK, locales.LK_ForgotPasswordPageSentTitle,
		locales.LK_ForgotPasswordPageSentDescription, g.RouteLogin, locales.LK_LoginPageSignIn, model.Email)
}

// PageResetPassword asks for a new password for the ?token= link
func (h *DefaultHandler) PageResetPassword(ctx *gin.Context) {
	invalid := func() {
		notice(ctx, http.StatusNotFound, locales.LK_ResetPasswordPageInvalidTitle,
			locales.LK_ResetPasswordPageInvalidDescription, g.RouteForgotPassword, locales.LK_ResetPasswordPageNewLink)
	}

	if ctx.Request.Method != http.MethodPost {
		token := ctx.Query("token")
		err := h.service.CheckPasswordResetToken(ctx, token)
		if errors.Is(err, g.ErrorTokenInvalid) {
			invalid()
			return
		}
		if err != nil {
			InternalError(ctx, err)
			return
		}
		//exhaustruct:ignore
		RenderTempl(ctx, pages.ResetPassword(models.ResetPassword{
			ModelBase: models.ModelBase{Initial: true},
			Token:     token,
		}))
		return
	}

	var model models.ResetPassword
	_ = ctx.ShouldBind(&model)
	if isValid, _ := model.IsValid(); !isValid {
		RenderTempl(ctx, pages.ResetPassword(model))
		return
	}
	err := h.service.ResetPassword(ctx, model.Token, model.Password)
	if errors.Is(err, g.ErrorTokenInvalid) {
		invalid()
		return
	}
	if err != nil {
		InternalError(ctx, err)
		return
	}
	notice(ctx, http.StatusOK, locales.LK_ResetPasswordPageDoneTitle,
		locales.LK_ResetPasswordPageDoneDescription, g.RouteLogin, locales.LK_LoginPageSignIn)
}
