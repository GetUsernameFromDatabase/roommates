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

func chatPath(householdID pgtype.UUID) string {
	return householdPath(householdID) + "/chat"
}

// PageChat shows the newest messages, or with ?before=<message id> the ones older than it
func (h *DefaultHandler) PageChat(ctx *gin.Context) {
	var before pgtype.UUID
	if b := ctx.Query("before"); b != "" {
		var err error
		if before, err = utils.ParseUUID(b); err != nil {
			RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
			return
		}
	}
	h.renderChat(ctx, http.StatusOK, before, "", "")
}

func (h *DefaultHandler) renderChat(ctx *gin.Context, status int, before pgtype.UUID, draft, errorMessage string) {
	actor := membershipFromContext(ctx)
	c, err := h.service.ListChat(ctx, actor, before)
	if err != nil {
		swapMainOnError(ctx)
		InternalError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, status, pages.ChatPage(actor, c, before.Valid, draft, errorMessage))
}

// PostMessage sends the "body" field with the optional "file",
// htmx requests get the chat itself so the composer is cleared without a reload
func (h *DefaultHandler) PostMessage(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, services.MaxFileSize+1<<20)
	err := h.postMessage(ctx, actor)
	body, c := ctx.PostForm("body"), ctx.Request.Context()
	switch {
	case err == nil && IsRequestHTMX(ctx):
		h.renderChat(ctx, http.StatusOK, pgtype.UUID{}, "", "")
	case err == nil:
		Redirect(ctx, chatPath(actor.HouseholdID))
	case errors.Is(err, g.ErrorFileTooLarge):
		h.renderChat(ctx, http.StatusRequestEntityTooLarge, pgtype.UUID{}, body,
			utils.T(c, locales.LK_FilesTooLarge, services.MaxFileSize>>20))
	case errors.Is(err, g.ErrorInvalidInput):
		h.renderChat(ctx, http.StatusUnprocessableEntity, pgtype.UUID{}, body,
			utils.T(c, locales.LK_ChatLength, services.MaxMessageLength))
	default:
		swapMainOnError(ctx)
		InternalError(ctx, err)
	}
}

func (h *DefaultHandler) postMessage(ctx *gin.Context, actor services.Membership) error {
	header, err := ctx.FormFile("file")
	var tooLarge *http.MaxBytesError
	switch {
	case errors.As(err, &tooLarge):
		return g.ErrorFileTooLarge
	case errors.Is(err, http.ErrMissingFile):
		return h.service.PostMessage(ctx, actor, ctx.PostForm("body"), nil)
	case err != nil:
		return g.ErrorInvalidInput
	}
	f, err := header.Open()
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck // read only
	return h.service.PostMessage(ctx, actor, ctx.PostForm("body"),
		&services.Upload{Name: header.Filename, Body: f, Size: header.Size})
}

func (h *DefaultHandler) DeleteMessage(ctx *gin.Context) {
	messageID, err := utils.ParseUUID(ctx.Param(g.ParamMessageID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	err = h.service.DeleteMessage(ctx, actor, messageID)
	switch {
	case err == nil:
		Redirect(ctx, chatPath(actor.HouseholdID))
	case errors.Is(err, g.ErrorNotFound):
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
	default:
		InternalError(ctx, err)
	}
}

// ChatBadge is the nav's unread count, the membership middleware has just counted it
func (h *DefaultHandler) ChatBadge(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	RenderTempl(ctx, pages.ChatBadge(actor.HouseholdID, actor.UnreadMessages))
}
