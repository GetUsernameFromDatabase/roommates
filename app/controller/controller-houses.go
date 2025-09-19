package controller

// TODO: API-s for
// - making a new house
// - inviting a new user to the house
// - changing the name of the house
// - changing the picture of the house

import (
	"fmt"
	"net/http"
	"roommates/components"
	"roommates/db/dbqueries"
	"roommates/gintemplrenderer"
	g "roommates/globals"
	"roommates/models"
	"roommates/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c *Controller) HtmxRoomateSearch(ctx *gin.Context) {
	if !utils.IsRequestHTMX(ctx) {
		utils.ErrorResponse(ctx, http.StatusForbidden, g.ErrorHtmxRequired)
		return
	}

	method := ctx.Request.Method
	var model models.House
	ctx.ShouldBind(&model)

	switch method {
	case http.MethodGet:
		render := func(foundUsers []dbqueries.UsersLikeExcludingExistingRow) {
			pc := components.HouseRoomatesInputSearchResults(model.SearchedUser, foundUsers)
			r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
			ctx.Render(r.Status, r)
		}

		if model.SearchedUser == "" {
			render(nil)
			return
		}

		users, err := c.DB.UsersLikeExcludingExisting(ctx, dbqueries.UsersLikeExcludingExistingParams{
			Username:      model.SearchedUser,
			ExistingUsers: model.RoommateLabels,
		})
		if err != nil {
			utils.ServerErrorResponse(ctx, "could not find users")
			return
		}
		render(users)
	case http.MethodPost:
		userId := ctx.PostForm("user_id")
		userLabel := ctx.PostForm("user_label")
		if userId == "" || userLabel == "" {
			err := fmt.Errorf("%s or %s missing from form", strconv.Quote("user_id"), strconv.Quote("user_label"))
			utils.ErrorResponse(ctx, http.StatusForbidden, err)
			return
		}

		model.RoommateKeys = append(model.RoommateKeys, userId)
		model.RoommateLabels = append(model.RoommateLabels, userLabel)

		pc := components.HouseForm(model)
		r := gintemplrenderer.New(ctx.Request.Context(), http.StatusOK, pc)
		ctx.Render(r.Status, r)
	default:
		ctx.String(http.StatusMethodNotAllowed, "method %s not allowed", method)
	}
}
