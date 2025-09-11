package controller

import (
	"errors"
	"net/http"
	"roommates/components"
	"roommates/gintemplrenderer"
	g "roommates/globals"
	"roommates/locales"
	"roommates/middleware"
	"roommates/models"
	"roommates/rdb"
	"roommates/utils"

	"github.com/gin-gonic/gin"
)

func (c *Controller) PageMain(ctx *gin.Context) {
	pc := components.PageMain()
	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageProfile(ctx *gin.Context) {
	pc := components.PageProfile()
	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PagePayments(ctx *gin.Context) {
	pc := components.PagePayments()
	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageNotes(ctx *gin.Context) {
	pc := components.PageNotes()
	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageMessaging(ctx *gin.Context) {
	pc := components.PageMessaging()
	r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
	ctx.Render(r.Status, r)
}

func (c *Controller) PageHouses(ctx *gin.Context) {
	pc := components.PageHouses()
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
	render := func(model models.Login) {
		page := components.PageLogin(model)
		r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, page)
		ctx.Render(r.Status, r)
	}

	switch method {
	case http.MethodGet:
		render(models.Login{Initial: true})
	case http.MethodPost:
		var model models.Login
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
				model.Error = utils.T(
					ctx.Request.Context(),
					locales.LKFormsErrorInvalidCredential,
					// this error is safe to output publically
					err.Error(),
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
	render := func(model models.Register) {
		page := components.PageRegister(model)
		r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, page)
		ctx.Render(r.Status, r)
	}

	switch method {
	case http.MethodGet:
		render(models.Register{Login: models.Login{Initial: true}})
	case http.MethodPost:
		var model models.Register
		ctx.ShouldBind(&model)

		hasError, _ := model.IsValid()
		if hasError {
			render(model)
			return
		}

		userID, err := c.registerUser(ctx, RegisterAccountRequest{
			Email:    model.Email,
			Password: model.Password,
			Username: model.Username,
		})
		if errors.Is(err, g.ErrorAccountAlreadyExists) {
			model.Error = utils.T(
				ctx.Request.Context(),
				locales.LKFormsErrorAlreadyExists,
				// this error is safe to output publically
				err.Error(),
			)
			render(model)
			return
		}
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
