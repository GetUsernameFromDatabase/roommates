package controller

import (
	"net/http"
	"roommates/components"
	"roommates/db/dbqueries"
	g "roommates/globals"
	"roommates/middleware"
	"roommates/models"
	"roommates/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// did the authenticated user make the house
//
// will only log the error if one occurs
func isNoteMaker(ctx *gin.Context, q *dbqueries.Queries, noteID int32) bool {
	authInfo := middleware.GetAuthInfo(ctx)
	isMaker, err := q.IsUserNoteMaker(ctx, dbqueries.IsUserNoteMakerParams{
		NoteID: noteID,
		UserID: authInfo.UserID,
	})

	if err != nil {
		log.Error().Err(err).Caller().
			Int32("note_id", noteID).
			Str("user_id", authInfo.UserID.String()).
			Msg("")
	}
	return isMaker
}

func renderNoteForm(ctx *gin.Context, model *models.Note) {
	tc := components.NoteForm(model)
	RenderTempl(ctx, tc)
}

// populates the note model with info from database
func (c *Controller) populateNoteModel(ctx *gin.Context, noteID int32, houseID pgtype.UUID) (*models.Note, error) {
	if noteID == 0 {
		house, err := c.DB.SelectHouse(ctx, houseID)
		if err != nil {
			return nil, err
		}
		model := models.NewNoteOnlyHouse(house)
		return &model, nil
	}

	note, err := c.DB.SelectNote(ctx, noteID)
	if err != nil {
		return nil, err
	}
	model := models.NewNote(note)
	return &model, nil
}

type ReqHxNoteInHouseAccordion struct {
	ID int32 `uri:"id" binding:"required"`
}

func (c *Controller) HxNoteInHouseAccordion(ctx *gin.Context) {
	var req ReqHxNoteInHouseAccordion
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err)
		return
	}

	note, err := c.DB.SelectNote(ctx, req.ID)
	if err != nil {
		HandleServerError(ctx, err, "could not get note")
		return
	}

	tc := components.NoteInHouseAccordion(note)
	RenderTempl(ctx, tc)
}

func (c *Controller) GetHxNoteModal(ctx *gin.Context) {
	houseID := requirePgUUID(ctx, "id")
	if houseID == nil {
		return
	}

	qNoteID := ctx.Query("note_id")
	if qNoteID == "" {
		model, err := c.populateNoteModel(ctx, 0, *houseID)
		if err != nil {
			HandleServerError(ctx, err, "could not get note data")
			return
		}
		tc := components.NoteModal(model)
		RenderTempl(ctx, tc)
		return
	}

	id, err := strconv.ParseInt(qNoteID, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err)
		return
	}
	model, err := c.populateNoteModel(ctx, int32(id), *houseID)
	if err != nil {
		HandleServerError(ctx, err, "could not get note data")
		return
	}
	tc := components.NoteModal(model)
	RenderTempl(ctx, tc)
}

func (c *Controller) PostHxNote(ctx *gin.Context) {
	houseID := requirePgUUID(ctx, "id")
	if houseID == nil {
		return
	}

	var model models.Note
	ctx.ShouldBind(&model)
	isValid, _ := model.IsValid()
	if !isValid {
		renderNoteForm(ctx, &model)
		return
	}

	authInfo := middleware.GetAuthInfo(ctx)
	_, err := c.DB.InsertNote(ctx, dbqueries.InsertNoteParams{
		Title:   model.Title,
		Content: model.Content,
		MakerID: authInfo.UserID,
		HouseID: *houseID,
	})
	if err != nil {
		HandleServerError(ctx, err, "could not save note")
		return
	}

	utils.Redirect(ctx, "")
}

// intended to be used with RNoteID
type ReqDeleteNote struct {
	ID int32 `uri:"id" binding:"required"`
}

func (c *Controller) DeleteNote(ctx *gin.Context) {
	var req ReqDeleteNote
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err)
		return
	}

	if isMaker := isNoteMaker(ctx, c.DB, req.ID); !isMaker {
		utils.ErrorResponse(ctx, http.StatusForbidden, g.ErrorNotAllowedToModify)
		return
	}

	if err := c.DB.DeleteNote(ctx, req.ID); err != nil {
		HandleServerError(ctx, err, "could not delete note")
		return
	}
	utils.Redirect(ctx, "")
}

// intended to be used with RNoteID
type ReqPutHxNote struct {
	ID int32 `uri:"id" binding:"required"`
}

func (c *Controller) PutHxNote(ctx *gin.Context) {
	var req ReqPutHxNote
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, err)
		return
	}

	if isMaker := isNoteMaker(ctx, c.DB, req.ID); !isMaker {
		utils.ErrorResponse(ctx, http.StatusForbidden, g.ErrorNotAllowedToModify)
		return
	}

	var model models.Note
	ctx.ShouldBind(&model)
	isValid, _ := model.IsValid()
	if !isValid {
		renderNoteForm(ctx, &model)
		return
	}

	err = c.DB.UpdateNote(ctx, dbqueries.UpdateNoteParams{
		ID:      req.ID,
		Title:   model.Title,
		Content: model.Content,
	})
	if err != nil {
		HandleServerError(ctx, err, "could not update note")
		return
	}
	utils.Redirect(ctx, "")
}
