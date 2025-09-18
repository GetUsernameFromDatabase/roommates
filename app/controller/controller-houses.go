package controller

// TODO: API-s for
// - making a new house
// - inviting a new user to the house
// - changing the name of the house
// - changing the picture of the house

import (
	"net/http"
	"roommates/components"
	"roommates/db/dbqueries"
	"roommates/gintemplrenderer"
	g "roommates/globals"
	"roommates/models"
	"roommates/utils"

	"github.com/gin-gonic/gin"
)

func (c *Controller) HtmxRoomateSearch(ctx *gin.Context) {
	if !utils.IsRequestHTMX(ctx) {
		utils.ErrorResponse(ctx, http.StatusForbidden, g.ErrorHtmxRequired)
		return
	}

	var model models.House
	ctx.ShouldBind(&model)
	render := func(foundUsers []dbqueries.UserLikeRow) {
		pc := components.HouseRoomatesInputSearchResults(model.SearchedUser, foundUsers)
		r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
		ctx.Render(r.Status, r)
	}

	if model.SearchedUser == "" {
		render(nil)
		return
	}

	users, err := c.DB.UserLike(ctx, model.SearchedUser)
	if err != nil {
		utils.ServerErrorResponse(ctx, "could not find users")
		return
	}
	render(users)
}
