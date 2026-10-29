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

// PageSettings is the user's picture and the appearance and notifications of this device,
// inside a household with its nav
func (h *DefaultHandler) PageSettings(ctx *gin.Context) {
	h.renderSettings(ctx, http.StatusOK, "")
}

// settingsPath is /settings, or the household's settings on a household route, with its nav
func settingsPath(ctx *gin.Context) (string, *pages.HouseholdNav) {
	if m, ok := ctx.Get(g.GinKeyMembership); ok {
		id := m.(services.Membership).HouseholdID
		return householdPath(id) + g.RouteSettings, &pages.HouseholdNav{HouseholdID: id, Section: pages.SectionSettings}
	}
	return g.RouteSettings, nil
}

// renderSettings shows the settings at /settings, or inside the household of the route
func (h *DefaultHandler) renderSettings(ctx *gin.Context, status int, errorMessage string) {
	path, nav := settingsPath(ctx)
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
	RenderTemplStatus(ctx, status, pages.Settings(h.service.PushPublicKey(), nav, path, user, errorMessage))
}
