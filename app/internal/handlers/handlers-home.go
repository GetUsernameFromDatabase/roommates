package handlers

import (
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/models"
	"roommates/internal/services"
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

// PageSettings is the appearance and notifications of this device,
// inside a household with its nav
func (h *DefaultHandler) PageSettings(ctx *gin.Context) {
	var nav *pages.HouseholdNav
	if m, ok := ctx.Get(g.GinKeyMembership); ok {
		nav = &pages.HouseholdNav{HouseholdID: m.(services.Membership).HouseholdID, Section: pages.SectionSettings}
	}
	RenderTempl(ctx, pages.Settings(h.service.PushPublicKey(), nav, ctx.Request.URL.Path))
}
