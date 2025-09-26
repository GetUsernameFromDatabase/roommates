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
	g "roommates/globals"
	"roommates/middleware"
	"roommates/models"
	"roommates/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func insertUsersToHouse(ctx *gin.Context, q *dbqueries.Queries, roomateIDs []pgtype.UUID, houseID pgtype.UUID) error {
	for _, roomateID := range roomateIDs {
		// not expecting hundreds of assignements here so should be fine
		err := q.InsertUserIntoHouse(ctx, dbqueries.InsertUserIntoHouseParams{
			UserID:  roomateID,
			HouseID: houseID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func renderHouseForm(ctx *gin.Context, model models.House) {
	tc := components.HouseForm(model)
	RenderTempl(ctx, tc)
}

// populates the house model with info from database
//
// used when rendering a form for editing a house
func (c *Controller) populateHouseModel(ctx *gin.Context, model *models.House) error {
	var houseID pgtype.UUID
	model.ScanHouseID(&houseID)
	if !houseID.Valid {
		// will signal to have new house
		model.HouseID = ""
		return nil
	}

	name, err := c.DB.SelectHouse(ctx, houseID)
	if err != nil {
		return err
	}
	model.Name = name

	roommates, err := c.DB.SelectHouseRoommates(ctx, houseID)
	if err != nil {
		return err
	}

	for _, roommate := range roommates {
		model.RoommateKeys = append(model.RoommateKeys, roommate.ID.String())
		model.RoommateLabels = append(model.RoommateLabels, roommate.Username)
	}
	return nil
}

func (c *Controller) HtmxRoomateSearch(ctx *gin.Context) {
	var model models.House
	ctx.ShouldBind(&model)
	model.Initial = true

	method := ctx.Request.Method
	switch method {
	case http.MethodGet:
		render := func(foundUsers []dbqueries.UsersLikeExcludingExistingRow) {
			tc := components.HouseRoommatesInputSearchResults(model.SearchedUser, foundUsers)
			RenderTempl(ctx, tc)
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
			HandleServerError(ctx, err, "could not find users")
			return
		}
		render(users)
	case http.MethodPost:
		userId := ctx.PostForm("user_id")
		userLabel := ctx.PostForm("user_label")
		if userId == "" || userLabel == "" {
			err := fmt.Errorf("%s or %s missing from request form", strconv.Quote("user_id"), strconv.Quote("user_label"))
			utils.ErrorResponse(ctx, http.StatusForbidden, err)
			return
		}

		model.RoommateKeys = append(model.RoommateKeys, userId)
		model.RoommateLabels = append(model.RoommateLabels, userLabel)

		renderHouseForm(ctx, model)
	default:
		ctx.String(http.StatusMethodNotAllowed, "method %s not allowed", method)
	}
}

func (c *Controller) GetHtmxHouseModal(ctx *gin.Context) {
	var model models.House
	model.HouseID = ctx.Query("house_id")

	model.Initial = true
	if err := c.populateHouseModel(ctx, &model); err != nil {
		HandleServerError(ctx, err, "could not get house data")
		return
	}

	tc := components.HouseModal(model)
	RenderTempl(ctx, tc)
}

func (c *Controller) PostHtmxHouseForm(ctx *gin.Context) {
	var model models.House
	ctx.ShouldBind(&model)

	isValid, _ := model.IsValid()
	if !isValid {
		renderHouseForm(ctx, model)
		return
	}
	conversionIssueOccured, roomateIDs := model.FilterNonValidUUID(ctx)
	if conversionIssueOccured {
		renderHouseForm(ctx, model)
		return
	}

	tx, err := c.Pool.Begin(ctx.Request.Context())
	if err != nil {
		HandleServerError(ctx, err, "no business pool party :(")
		return
	}
	defer tx.Rollback(ctx)
	qtx := c.DB.WithTx(tx)

	houseID, err := qtx.InsertHouse(ctx, model.Name)
	if err != nil {
		// currently there should not be unique violation issues
		HandleServerError(ctx, err, "unable to create this house")
		return
	}

	authInfo := middleware.GetAuthInfo(ctx)
	// it is assumed the user making the house want's to be in the house
	roomateIDs = append(roomateIDs, authInfo.UserID)

	err = insertUsersToHouse(ctx, qtx, roomateIDs, houseID)
	if err != nil {
		HandleServerError(ctx, err, "error assigning users to house")
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		HandleServerError(ctx, err, "error commiting transaction")
		return
	}

	// TODO: self repairing url with house name + id, where only id is of importance
	utils.Redirect(ctx, strings.ReplaceAll(g.RHouseID, ":id", houseID.String()))
}

func (c *Controller) DeleteHouse(ctx *gin.Context) {
	var houseID pgtype.UUID
	id := ctx.Query("id")
	err := houseID.Scan(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err)
		return
	}

	if err := c.DB.DeleteHouse(ctx, houseID); err != nil {
		HandleServerError(ctx, err, "error commiting transaction")
		return
	}
	utils.Redirect(ctx, g.RHouses)
}

// Delete and put for house form
func (c *Controller) PutHtmxHouseForm(ctx *gin.Context) {
	var model models.House
	ctx.ShouldBind(&model)

	isValid, _ := model.IsValid()
	if !isValid {
		renderHouseForm(ctx, model)
		return
	}
	var houseID pgtype.UUID
	if err := model.ScanHouseID(&houseID); err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err)
		return
	}
	conversionIssueOccured, roomateIDs := model.FilterNonValidUUID(ctx)
	if conversionIssueOccured {
		renderHouseForm(ctx, model)
		return
	}

	tx, err := c.Pool.Begin(ctx.Request.Context())
	if err != nil {
		HandleServerError(ctx, err, "no business pool party :(")
		return
	}
	defer tx.Rollback(ctx)
	qtx := c.DB.WithTx(tx)

	// pointless to keep the house around when it empty
	if len(roomateIDs) == 0 {
		if err := c.DB.DeleteHouse(ctx, houseID); err != nil {
			HandleServerError(ctx, err, "error commiting transaction")
			return
		}
	} else {
		qtx.UpdateHouse(ctx, dbqueries.UpdateHouseParams{
			Name: model.Name,
			ID:   houseID,
		})
		qtx.DeleteHouseUsers(ctx, houseID)

		err = insertUsersToHouse(ctx, qtx, roomateIDs, houseID)
		if err != nil {
			HandleServerError(ctx, err, "error assigning users to house")
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		HandleServerError(ctx, err, "error commiting transaction")
		return
	}
	utils.Redirect(ctx, ctx.Request.Referer())
}

func (c *Controller) HtmxHouseCardUser(ctx *gin.Context) {
	pId := ctx.Param("id") // pgtype.UUID does not do well with uri marshalling

	var userID pgtype.UUID
	err := userID.Scan(pId)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err)
		return
	}

	username, err := c.DB.SelectUsername(ctx, userID)
	if err != nil {
		HandleServerError(ctx, err, "could not get username")
		return
	}

	tc := components.HouseCardUserLink(userID, username)
	RenderTempl(ctx, tc)
}
