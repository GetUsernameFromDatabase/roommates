package middleware

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/routes"
	"roommates/internal/services"
	"roommates/internal/sessions"
	"roommates/internal/utils"
	"roommates/web/pages"
	"roommates/web/renderer"
	"slices"
	"strings"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

// NewHouseholdMember aborts with 404 unless the signed in user is an active
// member of the :householdID household, the membership is stored under
// [g.GinKeyMembership]
//
// 404 and not 403 so household ids of others cannot be probed
func NewHouseholdMember(service *services.DefaultService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.FromGinContext(ctx)
		userID, errUser := utils.ParseUUID(s.UserID)
		householdID, errHousehold := utils.ParseUUID(ctx.Param(routes.ParamHouseholdID))
		if errUser != nil || errHousehold != nil {
			notFound(ctx)
			return
		}

		m, err := service.GetMembership(ctx, householdID, userID)
		if err != nil {
			if !errors.Is(err, g.ErrorNotFound) {
				log.Error().Err(err).Msg("could not get membership")
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			notFound(ctx)
			return
		}

		// a deleted household is read-only until it is purged, restoring it is outside this group
		if m.DeletedAt.Valid && ctx.Request.Method != http.MethodGet && ctx.Request.Method != http.MethodHead &&
			!slices.Contains(userSettings, strings.TrimPrefix(ctx.FullPath(), routes.HouseholdSettings)) {
			readOnly(ctx, m)
			return
		}

		// the nav names the household and switches to the others
		households, err := service.ListHouseholds(ctx, userID)
		if err != nil {
			log.Error().Err(err).Msg("could not list households")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Set(g.GinKeyMembership, m)
		rctx := utils.WithUnreadChat(ctx.Request.Context(), m.UnreadMessages)
		rctx = utils.WithPendingSettlements(rctx, m.PendingSettlements)
		rctx = pages.WithMembership(rctx, m)
		ctx.Request = ctx.Request.WithContext(utils.WithHouseholds(rctx, households))
		ctx.Next()
	}
}

// userSettings are the forms of a household's settings page that change the user, not the
// household, so a deleted household's read-only check lets them through (see accountRoutes)
var userSettings = []string{
	routes.AccountAvatar,
	routes.AccountAvatarDelete,
	routes.AccountUsername,
	routes.AccountPassword,
	routes.AccountEmail,
	routes.AccountSignOutEverywhere,
	routes.AccountPayment,
}

// readOnly refuses a change to a deleted household, htmx shows the message in <main>
// like handlers.swapMainOnError, as a form may target only a part of the page
func readOnly(ctx *gin.Context, m services.Membership) {
	rctx := ctx.Request.Context()
	if ctx.GetHeader(g.HHK_HXRequest.String()) != "" {
		ctx.Header(g.HHK_HXRetarget.String(), "main")
		ctx.Header(g.HHK_HXReselect.String(), "main > *")
		ctx.Header(g.HHK_HXReswap.String(), "innerHTML")
		ctx.Header(g.HHK_HXPushURL.String(), "false")
	}
	page := pages.Notice(utils.T(rctx, locales.LK_OwnershipReadOnlyTitle),
		utils.T(rctx, locales.LK_OwnershipReadOnly, utils.FormatDate(rctx, services.PurgeDate(m.DeletedAt))),
		templ.SafeURL(routes.URL(routes.Household, m.HouseholdID)), utils.T(rctx, locales.LK_OwnershipBack))
	r := renderer.New(rctx, http.StatusForbidden, page)
	ctx.Render(r.Status, r)
	ctx.Abort()
}

func notFound(ctx *gin.Context) {
	r := renderer.New(ctx.Request.Context(), http.StatusNotFound, pages.NotFound())
	ctx.Render(r.Status, r)
	ctx.Abort()
}
