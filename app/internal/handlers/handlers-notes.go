package handlers

import (
	"errors"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func notesPath(householdID pgtype.UUID) string {
	return householdPath(householdID) + "/notes"
}

func notePath(householdID, noteID pgtype.UUID) string {
	return notesPath(householdID) + "/" + noteID.String()
}

func noteFromForm(ctx *gin.Context) services.NoteInput {
	return services.NoteInput{
		Title:   ctx.PostForm("title"),
		BodyMD:  ctx.PostForm("body"),
		Private: ctx.PostForm("private") == "1",
	}
}

// noteError shows the page for an error of the notes service
func noteError(ctx *gin.Context, err error) {
	if errors.Is(err, g.ErrorNotFound) {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	InternalError(ctx, err)
}

func noteLengthMessage(ctx *gin.Context) string {
	return utils.T(ctx.Request.Context(), locales.LK_NotesLength, services.MaxNoteTitleLength, services.MaxNoteBodyLength)
}

func (h *DefaultHandler) PageNotes(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	notes, err := h.service.ListNotes(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTempl(ctx, pages.HouseholdNotes(actor, notes))
}

func (h *DefaultHandler) PageNewNote(ctx *gin.Context) {
	RenderTempl(ctx, pages.NoteEditor(membershipFromContext(ctx), pages.NoteForm{NoteID: pgtype.UUID{}, Input: services.NoteInput{}, IsAuthor: true, ErrorMessage: ""}))
}

// CreateNote opens the new note
func (h *DefaultHandler) CreateNote(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	in := noteFromForm(ctx)
	noteID, err := h.service.CreateNote(ctx, actor, in)
	switch {
	case err == nil:
		Redirect(ctx, notePath(actor.HouseholdID, noteID))
	case errors.Is(err, g.ErrorInvalidInput):
		RenderTemplStatus(ctx, http.StatusUnprocessableEntity, pages.NoteEditor(actor,
			pages.NoteForm{NoteID: pgtype.UUID{}, Input: in, IsAuthor: true, ErrorMessage: noteLengthMessage(ctx)}))
	default:
		InternalError(ctx, err)
	}
}

// loadNote gets the :noteID note, or shows the error page and returns false
func (h *DefaultHandler) loadNote(ctx *gin.Context) (services.Note, bool) {
	noteID, err := utils.ParseUUID(ctx.Param(g.ParamNoteID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return services.Note{}, false
	}
	n, err := h.service.GetNote(ctx, membershipFromContext(ctx), noteID)
	if err != nil {
		noteError(ctx, err)
		return n, false
	}
	return n, true
}

func (h *DefaultHandler) PageNote(ctx *gin.Context) {
	if n, ok := h.loadNote(ctx); ok {
		RenderTempl(ctx, pages.NotePage(membershipFromContext(ctx), n))
	}
}

func (h *DefaultHandler) PageEditNote(ctx *gin.Context) {
	n, ok := h.loadNote(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	RenderTempl(ctx, pages.NoteEditor(actor, pages.NoteForm{
		NoteID: n.ID,
		Input: services.NoteInput{
			Title:   n.Title,
			BodyMD:  n.BodyMd,
			Private: n.Visibility == g.NoteVisibilityPrivate,
		},
		IsAuthor:     n.AuthorID == actor.UserID,
		ErrorMessage: "",
	}))
}

// UpdateNote overwrites the note with the form, last write wins
func (h *DefaultHandler) UpdateNote(ctx *gin.Context) {
	noteID, err := utils.ParseUUID(ctx.Param(g.ParamNoteID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	in := noteFromForm(ctx)
	err = h.service.UpdateNote(ctx, actor, noteID, in)
	switch {
	case err == nil:
		Redirect(ctx, notePath(actor.HouseholdID, noteID))
	case errors.Is(err, g.ErrorInvalidInput):
		// loaded only here, to know whether to show the private switch
		n, ok := h.loadNote(ctx)
		if !ok {
			return
		}
		RenderTemplStatus(ctx, http.StatusUnprocessableEntity, pages.NoteEditor(actor, pages.NoteForm{
			NoteID:       n.ID,
			Input:        in,
			IsAuthor:     n.AuthorID == actor.UserID,
			ErrorMessage: noteLengthMessage(ctx),
		}))
	default:
		noteError(ctx, err)
	}
}

// noteAction runs an action on the :noteID note and goes to next
func (h *DefaultHandler) noteAction(ctx *gin.Context, next func(householdID, noteID pgtype.UUID) string,
	action func(actor services.Membership, noteID pgtype.UUID) error) {
	noteID, err := utils.ParseUUID(ctx.Param(g.ParamNoteID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	if err := action(actor, noteID); err != nil {
		noteError(ctx, err)
		return
	}
	Redirect(ctx, next(actor.HouseholdID, noteID))
}

// SetNotePinned takes the wanted state, not a toggle, so a stale page cannot undo others
func (h *DefaultHandler) SetNotePinned(ctx *gin.Context) {
	h.noteAction(ctx, notePath, func(actor services.Membership, noteID pgtype.UUID) error {
		return h.service.SetNotePinned(ctx, actor, noteID, ctx.PostForm("pinned") == "1")
	})
}

func (h *DefaultHandler) DeleteNote(ctx *gin.Context) {
	h.noteAction(ctx, func(householdID, _ pgtype.UUID) string { return notesPath(householdID) },
		func(actor services.Membership, noteID pgtype.UUID) error {
			return h.service.DeleteNote(ctx, actor, noteID)
		})
}
