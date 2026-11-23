package handlers

import (
	"errors"
	"net/http"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/models"
	"roommates/internal/routes"
	"roommates/internal/services"
	"roommates/internal/sessions"
	"roommates/internal/utils"
	"roommates/web/pages"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// membershipFromContext panics if [middleware.NewHouseholdMember] did not run on the route
func membershipFromContext(ctx *gin.Context) services.Membership {
	return ctx.MustGet(g.GinKeyMembership).(services.Membership)
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
	Redirect(ctx, routes.URL(routes.Household, householdID))
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
	// a failed create keeps the typed list name
	RenderTemplStatus(ctx, status, pages.HouseholdShopping(actor, overview, errorMessage, ctx.PostForm("name")))
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
		ctx.Header(g.HHK_HXReplaceURL.String(), routes.URL(routes.HouseholdPeople, actor.HouseholdID))
		stayInPlace(ctx)
	}
	inviteURL := g.GetAppBaseURL() + routes.Invite + "?code="
	// a failed invite email keeps the typed address, a sent one clears the field
	inviteEmail := ""
	if at == pages.PeopleAtInviteEmail {
		inviteEmail = ctx.PostForm("email")
	}
	RenderTemplStatus(ctx, status, pages.HouseholdPeople(actor, overview, inviteURL, infoMessage, at, errorMessage, inviteEmail))
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
		Redirect(ctx, routes.URL(routes.HouseholdPeople, actor.HouseholdID))
	case errors.Is(err, g.ErrorLastOwner):
		h.renderPeople(ctx, http.StatusConflict, "", at,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdLastOwner))
	case errors.Is(err, g.ErrorBalanceNotZero):
		h.renderPeople(ctx, http.StatusConflict, "", at,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdBalanceNotZero))
	case errors.Is(err, g.ErrorEmailNotVerified):
		h.renderPeople(ctx, http.StatusForbidden, "", at,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdVerifyToInvite))
	case errors.Is(err, services.ErrorTooManyInvites):
		h.renderPeople(ctx, http.StatusTooManyRequests, "", at,
			utils.T(ctx.Request.Context(), locales.LK_HouseholdInviteTooMany))
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
	Redirect(ctx, routes.URL(routes.HouseholdPeople, actor.HouseholdID)+"?"+url.Values{queryInvited: {model.Email}}.Encode())
}

func (h *DefaultHandler) DeleteInvite(ctx *gin.Context) {
	h.householdTargetAction(ctx, routes.ParamInviteID, pages.PeopleAtInvites, func(actor services.Membership, inviteID pgtype.UUID) error {
		return h.service.DeleteInvite(ctx, actor, inviteID)
	})
}

func (h *DefaultHandler) RemoveMember(ctx *gin.Context) {
	h.householdTargetAction(ctx, routes.ParamUserID, pages.PeopleAtMembers, func(actor services.Membership, userID pgtype.UUID) error {
		return h.service.RemoveMember(ctx, actor, userID)
	})
}

func (h *DefaultHandler) PromoteMember(ctx *gin.Context) {
	h.householdTargetAction(ctx, routes.ParamUserID, pages.PeopleAtMembers, func(actor services.Membership, userID pgtype.UUID) error {
		return h.service.PromoteMember(ctx, actor, userID)
	})
}

// TransferOwnership makes a member an owner and the actor a member, the page then shows the new role
func (h *DefaultHandler) TransferOwnership(ctx *gin.Context) {
	h.householdTargetAction(ctx, routes.ParamUserID, pages.PeopleAtMembers, func(actor services.Membership, userID pgtype.UUID) error {
		return h.reloadActor(ctx, actor, h.service.TransferOwnership(ctx, actor, userID))
	})
}

// StepDown makes the actor a member while another owner is left
func (h *DefaultHandler) StepDown(ctx *gin.Context) {
	h.householdAction(ctx, pages.PeopleAtStepDown, func(actor services.Membership) error {
		return h.reloadActor(ctx, actor, h.service.StepDown(ctx, actor))
	})
}

// reloadActor stores the actor's membership again after a role change without an error,
// so the page rendered next is the one of their new role
func (h *DefaultHandler) reloadActor(ctx *gin.Context, actor services.Membership, err error) error {
	if err != nil {
		return err
	}
	m, err := h.service.GetMembership(ctx, actor.HouseholdID, actor.UserID)
	if err != nil {
		return err
	}
	ctx.Set(g.GinKeyMembership, m)
	return nil
}

// DeleteHousehold deletes the household after its name was typed, the page is shown again
// read-only, an error is swapped into the open dialog by the name field
func (h *DefaultHandler) DeleteHousehold(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	err := h.service.DeleteHousehold(ctx, actor, ctx.PostForm("name"))
	var field services.FieldError
	switch {
	case err == nil:
		stay(ctx, routes.URL(routes.HouseholdSettings, actor.HouseholdID))
	case errors.As(err, &field):
		if IsRequestHTMX(ctx) {
			ctx.Header(g.HHK_HXRetarget.String(), "#household-delete-error")
			ctx.Header(g.HHK_HXReswap.String(), "innerHTML")
		}
		RenderTemplStatus(ctx, http.StatusUnprocessableEntity,
			pages.DeleteHouseholdError(utils.T(ctx.Request.Context(), locales.LK_OwnershipDeleteNameMismatch)))
	case errors.Is(err, g.ErrorForbidden):
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusForbidden, pages.Error())
	default:
		InternalError(ctx, err)
	}
}

// RestoreHousehold undoes a deletion, the page the restore button is on is shown again
func (h *DefaultHandler) RestoreHousehold(ctx *gin.Context) {
	userID, err := sessionUserID(ctx)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	householdID, err := utils.ParseUUID(ctx.Param(routes.ParamHouseholdID))
	if err != nil {
		err = g.ErrorNotFound
	} else {
		err = h.service.RestoreHousehold(ctx, userID, householdID)
	}
	switch {
	case err == nil:
		// the household page the banner is on, or home
		next := routes.Home
		if u, err := url.Parse(ctx.GetHeader("HX-Current-URL")); err == nil && strings.HasPrefix(u.Path, routes.URL(routes.Household, householdID)) {
			next = u.RequestURI()
		}
		stay(ctx, next)
	case errors.Is(err, g.ErrorForbidden):
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusForbidden, pages.Error())
	case errors.Is(err, g.ErrorNotFound):
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
	default:
		InternalError(ctx, err)
	}
}

// LeaveHousehold sends the user home, the household page is no longer theirs
func (h *DefaultHandler) LeaveHousehold(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	err := h.service.Leave(ctx, actor)
	if err == nil {
		Redirect(ctx, routes.Home)
		return
	}
	// the errors are mapped as in any other people action
	h.householdAction(ctx, pages.PeopleAtLeave, func(services.Membership) error { return err })
}

// UpdateHousehold shows the settings section again, a changed name also refreshes the page
// as the name is in the household switcher, there the new name in the field is the outcome
func (h *DefaultHandler) UpdateHousehold(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	path := routes.URL(routes.HouseholdSettings, actor.HouseholdID)
	before, err := h.service.GetHouseholdSettings(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	name := ctx.PostForm("name")
	timeZone := ctx.PostForm("time_zone")
	// a hour that is no number is out of range, the service rejects it
	hour, convErr := strconv.Atoi(ctx.PostForm("reminder_hour"))
	if convErr != nil {
		hour = -1
	}
	err = h.service.UpdateHousehold(ctx, actor, name, ctx.PostForm("currency"), timeZone, hour)
	show := func(status int, msg, info string) {
		// keep what was typed, a zone or hour the selects cannot show stays as it was saved
		before.Name = name
		if services.ValidReminderHour(hour) {
			before.ReminderHour = hour
		}
		if slices.Contains(before.TimeZones, timeZone) {
			before.TimeZone = timeZone
		}
		RenderTemplStatus(ctx, status, pages.HouseholdSettingsSection(actor, before, msg, info))
	}
	var field services.FieldError
	switch {
	case err == nil && strings.TrimSpace(name) != before.Name:
		stay(ctx, path)
	case err == nil:
		before.Name = strings.TrimSpace(name)
		before.Currency = ctx.PostForm("currency")
		show(http.StatusOK, "", utils.T(ctx.Request.Context(), locales.LK_CommonSaved))
	case errors.Is(err, g.ErrorForbidden):
		RenderTemplStatus(ctx, http.StatusForbidden, pages.Error())
	case errors.As(err, &field) && field.Field == services.FieldName:
		show(http.StatusUnprocessableEntity, utils.T(ctx.Request.Context(), locales.LK_ValidationErrorsHouseholdNameLength, 1, models.MaxHouseholdNameLength), "")
	case errors.As(err, &field) && field.Field == services.FieldTimeZone:
		show(http.StatusUnprocessableEntity, utils.T(ctx.Request.Context(), locales.LK_HouseholdSettingsTimeZoneInvalid), "")
	case errors.As(err, &field) && field.Field == services.FieldReminderHour:
		show(http.StatusUnprocessableEntity, utils.T(ctx.Request.Context(), locales.LK_HouseholdSettingsReminderHourInvalid), "")
	case errors.As(err, &field):
		show(http.StatusUnprocessableEntity, utils.T(ctx.Request.Context(), locales.LK_HouseholdSettingsCurrencyLocked), "")
	default:
		InternalError(ctx, err)
	}
}
