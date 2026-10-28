package handlers

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/models"
	"roommates/internal/services"
	"roommates/internal/sessions"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
)

func (h *DefaultHandler) PageLogin(ctx *gin.Context) {
	method := ctx.Request.Method
	render := func(model models.Login) {
		page := pages.Login(model)
		RenderTempl(ctx, page)
	}

	switch method {
	case http.MethodGet:
		//exhaustruct:ignore
		render(models.Login{
			ModelBase: models.ModelBase{Initial: true},
			Next:      ctx.Query(g.QueryNext),
		})
	case http.MethodPost:
		var model models.Login
		_ = ctx.ShouldBind(&model)

		isValid, _ := model.IsValid()
		if !isValid {
			render(model)
			return
		}

		credsInDB, err := h.service.ShouldUserBeSignedIn(ctx, model.LoginRequest)
		if err != nil {
			if errors.Is(err, g.ErrorInvalidCredential) {
				model.Error = utils.T(
					ctx.Request.Context(),
					locales.LK_LoginPageInvalidCredentials,
				)
				render(model)
				return
			}
			InternalError(ctx, err)
			return
		}

		s := sessions.FromGinContext(ctx)
		s.Login(sessions.SessionInfo{
			UserID:   credsInDB.ID.String(),
			Email:    credsInDB.Email,
			Username: credsInDB.Username,
		})
		Redirect(ctx, safeNextPath(model.Next))
	default:
		ctx.String(http.StatusMethodNotAllowed, "method %s not allowed", method)
	}
}

func (h *DefaultHandler) PageRegister(ctx *gin.Context) {
	method := ctx.Request.Method
	render := func(model models.Register) {
		page := pages.Register(model)
		RenderTempl(ctx, page)
	}

	switch method {
	case http.MethodGet:
		//exhaustruct:ignore
		render(models.Register{Login: models.Login{
			ModelBase: models.ModelBase{Initial: true},
			Next:      ctx.Query(g.QueryNext),
		}})
	case http.MethodPost:
		var model models.Register
		_ = ctx.ShouldBind(&model)

		isValid, _ := model.IsValid()
		if !isValid {
			render(model)
			return
		}

		userID, err := h.service.RegisterUser(ctx, dbqueries.InsertUserParams{
			Email:    model.Email,
			Username: model.Username,
			Password: model.Password,
		})
		if err != nil {
			if errors.Is(err, g.ErrorAccountAlreadyExists) {
				model.Error = utils.T(
					ctx.Request.Context(),
					locales.LK_RegistrationPageAccountExists,
				)
				render(model)
				return
			}
			InternalError(ctx, err)
			return
		}

		// the account exists already, a failed email can be resent from the banner
		if err := h.service.SendVerificationEmail(ctx.Request.Context(), *userID); err != nil {
			log.Error().Err(err).Msg("could not send verification email")
		}

		s := sessions.FromGinContext(ctx)
		s.Login(sessions.SessionInfo{
			UserID:   userID.String(),
			Email:    services.NormalizeEmail(model.Email),
			Username: model.Username,
		})
		Redirect(ctx, safeNextPath(model.Next))
	default:
		ctx.String(http.StatusMethodNotAllowed, "method %s not allowed", method)
	}
}

// Logout ends the session and sends the user to the login page
func (h *DefaultHandler) Logout(ctx *gin.Context) {
	s := sessions.FromGinContext(ctx)
	s.Logout()
	Redirect(ctx, g.RouteLogin)
}
