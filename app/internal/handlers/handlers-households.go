package handlers

import (
	"errors"
	"net/http"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/models"
	"roommates/internal/services"
	"roommates/internal/sessions"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// membershipFromContext panics if [middleware.NewHouseholdMember] did not run on the route
func membershipFromContext(ctx *gin.Context) services.Membership {
	return ctx.MustGet(g.GinKeyMembership).(services.Membership)
}

// householdPath is the shopping section, the first page of a household
func householdPath(id pgtype.UUID) string {
	return g.RouteHousehold + "/" + id.String()
}

func peoplePath(id pgtype.UUID) string {
	return householdPath(id) + "/people"
}

func sessionUserID(ctx *gin.Context) (pgtype.UUID, error) {
	return utils.ParseUUID(sessions.FromGinContext(ctx).UserID)
}

// PageInvite shows the household of ?code=, POST joins it
func (h *DefaultHandler) PageInvite(ctx *gin.Context) {
	code := services.NormalizeInviteCode(ctx.Query("code"))
	invite, err := h.service.GetInvite(ctx, code)
	if err != nil {
		if errors.Is(err, g.ErrorInviteInvalid) {
			RenderTemplStatus(ctx, http.StatusNotFound, pages.Invite(code, "",
				utils.T(ctx.Request.Context(), locales.LK_InvitePageInvalid)))
			return
		}
		InternalError(ctx, err)
		return
	}

	if ctx.Request.Method != http.MethodPost {
		RenderTempl(ctx, pages.Invite(code, invite.HouseholdName, ""))
		return
	}

	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	householdID, err := h.service.JoinWithInvite(ctx, code, userID)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	Redirect(ctx, householdPath(householdID))
}

// set after [DefaultHandler.SendInviteEmail], so a refresh does not resend
const queryInvited = "invited"

// PageHousehold shows the shopping lists
func (h *DefaultHandler) PageHousehold(ctx *gin.Context) {
	h.renderShopping(ctx, http.StatusOK, "")
}

func (h *DefaultHandler) renderShopping(ctx *gin.Context, status int, errorMessage string) {
	actor := membershipFromContext(ctx)
	overview, err := h.service.GetHouseholdOverview(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	if errorMessage != "" {
		stayInPlace(ctx)
	}
	RenderTemplStatus(ctx, status, pages.HouseholdShopping(actor, overview, errorMessage))
}

// PagePeople shows members and, for owners, invites
func (h *DefaultHandler) PagePeople(ctx *gin.Context) {
	info := ""
	if invited := ctx.Query(queryInvited); invited != "" {
		info = utils.T(ctx.Request.Context(), locales.LK_HouseholdInviteEmailSent, invited)
	}
	h.renderPeople(ctx, http.StatusOK, info, "", "")
}

// renderPeople shows infoMessage as a toast and errorMessage by the action at (a pages.PeopleAt*)
func (h *DefaultHandler) renderPeople(ctx *gin.Context, status int, infoMessage, at, errorMessage string) {
	actor := membershipFromContext(ctx)
	overview, err := h.service.GetHouseholdOverview(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	// an action answers htmx with the page itself, at the people address and where it was scrolled
	if ctx.Request.Method == http.MethodPost && IsRequestHTMX(ctx) {
		ctx.Header(g.HHK_HXReplaceURL.String(), peoplePath(actor.HouseholdID))
		stayInPlace(ctx)
	}
	inviteURL := g.GetAppBaseURL() + g.RouteInvite + "?code="
	RenderTemplStatus(ctx, status, pages.HouseholdPeople(actor, overview, inviteURL, infoMessage, at, errorMessage))
}

// householdAction runs an action of the people page and goes back to it,
// an error shows by the action at (a pages.PeopleAt*)
func (h *DefaultHandler) householdAction(ctx *gin.Context, at string, action func(actor services.Membership) error) {
	actor := membershipFromContext(ctx)
	err := action(actor)
	switch {
	// HX-Redirect would reload the whole page, htmx gets the page in place instead
	case err == nil && IsRequestHTMX(ctx):
		h.renderPeople(ctx, http.StatusOK, "", "", "")
	case err == nil:
		Redirect(ctx, peoplePath(actor.HouseholdID))
	case errors.Is(err, g.ErrorLastOwner):
		h.renderPeople(ctx, http.StatusConflict, "", at,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdLastOwner))
	case errors.Is(err, g.ErrorBalanceNotZero):
		h.renderPeople(ctx, http.StatusConflict, "", at,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdBalanceNotZero))
	case errors.Is(err, g.ErrorEmailNotVerified):
		h.renderPeople(ctx, http.StatusForbidden, "", at,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdVerifyToInvite))
	case errors.Is(err, g.ErrorNotFound):
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
	case errors.Is(err, g.ErrorForbidden):
		RenderTemplStatus(ctx, http.StatusForbidden, pages.Error())
	default:
		InternalError(ctx, err)
	}
}

// householdTargetAction is [householdAction] with the :param UUID parsed
func (h *DefaultHandler) householdTargetAction(ctx *gin.Context, param, at string,
	action func(actor services.Membership, target pgtype.UUID) error) {
	target, err := utils.ParseUUID(ctx.Param(param))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	h.householdAction(ctx, at, func(actor services.Membership) error {
		return action(actor, target)
	})
}

func (h *DefaultHandler) CreateInvite(ctx *gin.Context) {
	h.householdAction(ctx, pages.PeopleAtInvites, func(actor services.Membership) error {
		return h.service.CreateInvite(ctx, actor)
	})
}

// SendInviteEmail creates an invite and emails it, the page says to whom
func (h *DefaultHandler) SendInviteEmail(ctx *gin.Context) {
	var model models.EmailForm
	_ = ctx.ShouldBind(&model)
	if isValid, _ := model.IsValid(); !isValid {
		h.renderPeople(ctx, http.StatusUnprocessableEntity, "", pages.PeopleAtInviteEmail,
			utils.T(ctx.Request.Context(), locales.LK_ValidationErrorsEmailGeneric))
		return
	}
	actor := membershipFromContext(ctx)
	// the request context carries the language of the email
	err := h.service.SendInviteEmail(ctx.Request.Context(), actor, model.Email)
	if err != nil {
		h.householdAction(ctx, pages.PeopleAtInviteEmail, func(services.Membership) error { return err })
		return
	}
	// htmx gets the page with the toast in place, without JS the redirect carries the address
	if IsRequestHTMX(ctx) {
		h.renderPeople(ctx, http.StatusOK,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdInviteEmailSent, model.Email), "", "")
		return
	}
	Redirect(ctx, peoplePath(actor.HouseholdID)+"?"+url.Values{queryInvited: {model.Email}}.Encode())
}

func (h *DefaultHandler) DeleteInvite(ctx *gin.Context) {
	h.householdTargetAction(ctx, g.ParamInviteID, pages.PeopleAtInvites, func(actor services.Membership, inviteID pgtype.UUID) error {
		return h.service.DeleteInvite(ctx, actor, inviteID)
	})
}

func (h *DefaultHandler) RemoveMember(ctx *gin.Context) {
	h.householdTargetAction(ctx, g.ParamUserID, pages.PeopleAtMembers, func(actor services.Membership, userID pgtype.UUID) error {
		return h.service.RemoveMember(ctx, actor, userID)
	})
}

func (h *DefaultHandler) PromoteMember(ctx *gin.Context) {
	h.householdTargetAction(ctx, g.ParamUserID, pages.PeopleAtMembers, func(actor services.Membership, userID pgtype.UUID) error {
		return h.service.PromoteMember(ctx, actor, userID)
	})
}

// LeaveHousehold sends the user home, the household page is no longer theirs
func (h *DefaultHandler) LeaveHousehold(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	err := h.service.Leave(ctx, actor)
	if errors.Is(err, g.ErrorLastOwner) {
		h.renderPeople(ctx, http.StatusConflict, "", pages.PeopleAtLeave,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdLastOwner))
		return
	}
	if errors.Is(err, g.ErrorBalanceNotZero) {
		h.renderPeople(ctx, http.StatusConflict, "", pages.PeopleAtLeave,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdBalanceNotZero))
		return
	}
	if err != nil {
		InternalError(ctx, err)
		return
	}
	Redirect(ctx, "/")
}
