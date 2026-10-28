package handlers

import (
	"errors"
	"net/http"
	"roommates/web/pages"
	"roommates/internal/db/dbqueries"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/models"
	"roommates/internal/sessions"
	"roommates/internal/utils"

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
		render(models.Login{Initial: true})
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
			UserID: credsInDB.ID.String(),
			Email:    credsInDB.Email,
			Username: credsInDB.Username,
		})
		Redirect(ctx, "/")
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
		render(models.Register{Initial: true})
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
		if errors.Is(err, g.ErrorAccountAlreadyExists) {
			model.Error = utils.T(
				ctx.Request.Context(),
				locales.LK_RegistrationPageAccountExists,
			)
			render(model)
			return
		}

		s := sessions.FromGinContext(ctx)
		s.Login(sessions.SessionInfo{
			UserID:   userID.String(),
			Email:    model.Email,
			Username: model.Username,
		})
		ctx.Redirect(http.StatusSeeOther, "/")
	default:
		ctx.String(http.StatusMethodNotAllowed, "method %s not allowed", method)
	}
}
