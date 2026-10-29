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

// PageChat shows the newest messages, or with ?before=<message id> the ones older than it,
// ?reply=<message id> answers that message
func (h *DefaultHandler) PageChat(ctx *gin.Context) {
	before, errBefore := optionalUUID(ctx.Query("before"))
	replyTo, errReply := optionalUUID(ctx.Query("reply"))
	if errBefore != nil || errReply != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	h.renderChat(ctx, http.StatusOK, before, replyTo, "", "")
}

// optionalUUID is a zero UUID for an empty s
func optionalUUID(s string) (pgtype.UUID, error) {
	if s == "" {
		return pgtype.UUID{}, nil
	}
	return utils.ParseUUID(s)
}

// renderChat shows the composer answering replyTo, a message that is gone is no longer answered
func (h *DefaultHandler) renderChat(ctx *gin.Context, status int, before, replyTo pgtype.UUID, draft, errorMessage string) {
	actor := membershipFromContext(ctx)
	c, err := h.service.ListChat(ctx, actor, before)
	if err != nil {
		swapMainOnError(ctx)
		InternalError(ctx, err)
		return
	}
	var reply *services.ReplyPreview
	if replyTo.Valid {
		p, err := h.service.GetReplyPreview(ctx, actor, replyTo)
		switch {
		case err == nil:
			reply = &p
		case !errors.Is(err, g.ErrorNotFound):
			swapMainOnError(ctx)
			InternalError(ctx, err)
			return
		}
	}
	RenderTemplStatus(ctx, status, pages.ChatPage(actor, c, before.Valid, reply, draft, errorMessage))
}

// PostMessage sends the "body" field with the optional "file",
// htmx requests get the chat itself so the composer is cleared without a reload
func (h *DefaultHandler) PostMessage(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, services.MaxFileSize+1<<20)
	replyTo, err := optionalUUID(ctx.PostForm("reply_to"))
	if err == nil {
		err = h.postMessage(ctx, actor, replyTo)
	} else {
		err = g.ErrorInvalidInput
	}
	body, c := ctx.PostForm("body"), ctx.Request.Context()
	switch {
	case err == nil && IsRequestHTMX(ctx):
		// a ?reply= left in the address bar would answer the same message again after a reload
		ctx.Header(g.HHK_HXReplaceURL.String(), chatPath(actor.HouseholdID))
		h.renderChat(ctx, http.StatusOK, pgtype.UUID{}, pgtype.UUID{}, "", "")
	case err == nil:
		Redirect(ctx, chatPath(actor.HouseholdID))
	case errors.Is(err, g.ErrorFileTooLarge):
		h.renderChat(ctx, http.StatusRequestEntityTooLarge, pgtype.UUID{}, replyTo, body,
			utils.T(c, locales.LK_ChatFilesTooLarge, services.MaxFileSize>>20))
	case errors.Is(err, g.ErrorNotFound):
		// the answered message is gone, the draft comes back without its banner and is sent as a plain message
		h.renderChat(ctx, http.StatusUnprocessableEntity, pgtype.UUID{}, pgtype.UUID{}, body, "")
	case errors.Is(err, g.ErrorInvalidInput):
		h.renderChat(ctx, http.StatusUnprocessableEntity, pgtype.UUID{}, replyTo, body,
			utils.T(c, locales.LK_ChatLength, services.MaxMessageLength))
	default:
		swapMainOnError(ctx)
		InternalError(ctx, err)
	}
}

func (h *DefaultHandler) postMessage(ctx *gin.Context, actor services.Membership, replyTo pgtype.UUID) error {
	form, err := ctx.MultipartForm()
	var tooLarge *http.MaxBytesError
	switch {
	case errors.As(err, &tooLarge):
		return g.ErrorFileTooLarge
	case errors.Is(err, http.ErrNotMultipart):
		// a form without files, the composer always sends multipart
		return h.service.PostMessage(ctx, actor, ctx.PostForm("body"), replyTo, nil)
	case err != nil:
		return g.ErrorInvalidInput
	}
	var files []services.Upload
	for _, header := range form.File["file"] {
		f, err := header.Open()
		if err != nil {
			return err
		}
		defer f.Close() //nolint:errcheck // read only
		files = append(files, services.Upload{Name: header.Filename, Body: f, Size: header.Size})
	}
	return h.service.PostMessage(ctx, actor, ctx.PostForm("body"), replyTo, files)
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
	case err == nil && IsRequestHTMX(ctx):
		// the form swaps only the messages out of it
		h.renderChat(ctx, http.StatusOK, pgtype.UUID{}, pgtype.UUID{}, "", "")
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
