package controller

import (
	"errors"
	"net/http"
	"roommates/components"
	"roommates/gintemplrenderer"
	g "roommates/globals"
	"roommates/locales"
	"roommates/middleware"
	"roommates/rdb"

	"github.com/gin-gonic/gin"
	"github.com/invopop/ctxi18n/i18n"
)

func (c *Controller) PageMain(ctx *gin.Context) {
	pc := components.PageMain()
	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

// --- AUTH -- login and register ---

func (c *Controller) PageLogin(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)
	if authInfo != nil {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	method := ctx.Request.Method
	render := func(model components.LoginModel) {
		page := components.PageLogin(model)
		r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, page)
		ctx.Render(r.Status, r)
	}

	switch method {
	case http.MethodGet:
		render(components.LoginModel{Initial: true})
	case http.MethodPost:
		var model components.LoginModel
		ctx.ShouldBind(&model)

		hasError, _ := model.IsValid()
		if hasError {
			render(model)
			return
		}

		credsInDb, err := c.shouldUserBeSignedIn(ctx, SignInRequest{
			Email:    model.Email,
			Password: model.Password,
		})
		if err != nil {
			if errors.Is(err, g.ErrorInvalidCredential) {
				model.Error = i18n.T(
					ctx.Request.Context(),
					string(locales.LKFormsErrorInvalidCredential),
					i18n.Default(err.Error()),
				)
				render(model)
				return
			}
			HandleServerError(ctx, err, "error fetching user credentials")
			return
		}

		a := c.signUserIn(ctx, rdb.UserSessionValue{
			UserID:   credsInDb.ID.String(),
			Username: credsInDb.Username,
		})
		log.Info().Msg(a.String())
		ctx.Redirect(http.StatusSeeOther, "/")
	default:
		ctx.String(http.StatusMethodNotAllowed, "method %s not allowed", method)
	}
}

func (c *Controller) PageRegister(ctx *gin.Context) {
	authInfo := middleware.GetAuthInfo(ctx)
	if authInfo != nil {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	method := ctx.Request.Method
	render := func(model components.RegisterModel) {
		page := components.PageRegister(model)
		r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, page)
		ctx.Render(r.Status, r)
	}

	switch method {
	case http.MethodGet:
		render(components.RegisterModel{LoginModel: components.LoginModel{Initial: true}})
	case http.MethodPost:
		var model components.RegisterModel
		ctx.ShouldBind(&model)

		hasError, _ := model.IsValid()
		if hasError {
			render(model)
			return
		}

		userID := c.registerUser(ctx, RegisterAccountRequest{
			Email:    model.Email,
			Password: model.Password,
			Username: model.Username,
		})
		if userID == "" {
			return
		}

		c.signUserIn(ctx, rdb.UserSessionValue{
			UserID:   userID,
			Username: model.Username,
		})
		ctx.Redirect(http.StatusSeeOther, "/")
	default:
		ctx.String(http.StatusMethodNotAllowed, "method %s not allowed", method)
	}
}
