package handlers

import (
	"net/http"
	"roommates/internal/models"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
)

// PageHome lists the user's households, POST creates one
func (h *DefaultHandler) PageHome(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}

	//exhaustruct:ignore
	model := models.Household{ModelBase: models.ModelBase{Initial: true}}
	if ctx.Request.Method == http.MethodPost {
		_ = ctx.ShouldBind(&model)
		if isValid, _ := model.IsValid(); isValid {
			householdID, err := h.service.CreateHousehold(ctx, userID, model.Name)
			if err != nil {
				InternalError(ctx, err)
				return
			}
			Redirect(ctx, householdPath(householdID))
			return
		}
	}

	households, err := h.service.ListHouseholds(ctx, userID)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTempl(ctx, pages.Home(user.Username, user.EmailValidated, households, model))
}
