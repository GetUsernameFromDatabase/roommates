package middleware

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/services"
	"roommates/internal/sessions"
	"roommates/internal/utils"
	"roommates/web/pages"
	"roommates/web/renderer"

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
		householdID, errHousehold := utils.ParseUUID(ctx.Param(g.ParamHouseholdID))
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

		// the nav names the household and switches to the others
		households, err := service.ListHouseholds(ctx, userID)
		if err != nil {
			log.Error().Err(err).Msg("could not list households")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.Set(g.GinKeyMembership, m)
		rctx := utils.WithUnreadChat(ctx.Request.Context(), m.UnreadMessages)
		ctx.Request = ctx.Request.WithContext(utils.WithHouseholds(rctx, households))
		ctx.Next()
	}
}

func notFound(ctx *gin.Context) {
	r := renderer.New(ctx.Request.Context(), http.StatusNotFound, pages.NotFound())
	ctx.Render(r.Status, r)
	ctx.Abort()
}
