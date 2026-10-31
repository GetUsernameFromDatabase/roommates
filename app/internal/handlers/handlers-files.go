package handlers

import (
	"cmp"
	"errors"
	"mime"
	"net/http"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func filesPath(householdID pgtype.UUID) string {
	return householdPath(householdID) + "/files"
}

func shoppingItemPath(householdID, listID, itemID pgtype.UUID) string {
	return shoppingListPath(householdID, listID) + "/items/" + itemID.String()
}

// uploadFile passes the "file" form field to the service, the body is cut off a little
// above [services.MaxFileSize] so a huge upload is not read to the end
func (h *DefaultHandler) uploadFile(ctx *gin.Context, actor services.Membership, t services.FileTarget) error {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, services.MaxFileSize+1<<20)
	header, err := ctx.FormFile("file")
	var tooLarge *http.MaxBytesError
	switch {
	case errors.As(err, &tooLarge):
		return g.ErrorFileTooLarge
	case err != nil:
		// no file chosen or a broken form
		return g.ErrorInvalidInput
	}
	f, err := header.Open()
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck // read only
	_, err = h.service.UploadFile(ctx, actor, t, header.Filename, f, header.Size)
	return err
}

// handleUpload uploads to t and goes to done, or shows the page again with rerender
func (h *DefaultHandler) handleUpload(ctx *gin.Context, t services.FileTarget, done string,
	rerender func(status int, errorMessage string)) {
	actor := membershipFromContext(ctx)
	err := h.uploadFile(ctx, actor, t)
	c := ctx.Request.Context()
	switch {
	case err == nil:
		Redirect(ctx, done)
	case errors.Is(err, g.ErrorFileTooLarge):
		rerender(http.StatusRequestEntityTooLarge, utils.T(c, locales.LK_FilesTooLarge, services.MaxFileSize>>20))
	case errors.Is(err, g.ErrorInvalidInput):
		rerender(http.StatusUnprocessableEntity, utils.T(c, locales.LK_FilesInvalid, services.MaxFileNameLength))
	case errors.Is(err, g.ErrorNotFound):
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
	default:
		InternalError(ctx, err)
	}
}

// -----------------------------------------------------------------------------

func (h *DefaultHandler) PageFiles(ctx *gin.Context) {
	h.renderFiles(ctx, http.StatusOK, "")
}

func (h *DefaultHandler) renderFiles(ctx *gin.Context, status int, errorMessage string) {
	actor := membershipFromContext(ctx)
	filter := services.FileFilter{Query: ctx.Query("q"), Kind: ctx.Query("type"), Place: ctx.Query("place")}
	files, err := h.service.ListFiles(ctx, actor, filter)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, status, pages.HouseholdFiles(actor, files, filter, errorMessage))
}

// UploadFile adds a file that is not attached to anything
func (h *DefaultHandler) UploadFile(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	h.handleUpload(ctx, services.FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}},
		filesPath(actor.HouseholdID), func(status int, errorMessage string) {
			h.renderFiles(ctx, status, errorMessage)
		})
}

func (h *DefaultHandler) UploadNoteFile(ctx *gin.Context) {
	n, ok := h.loadNote(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	h.handleUpload(ctx, services.FileTarget{NoteID: n.ID, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}},
		notePath(actor.HouseholdID, n.ID), func(status int, errorMessage string) {
			h.renderNote(ctx, status, n, errorMessage)
		})
}

// loadShoppingItem gets the :itemID item of the :listID list, or shows the error page and returns false
func (h *DefaultHandler) loadShoppingItem(ctx *gin.Context) (pgtype.UUID, services.ShoppingItem, bool) {
	listID, errList := utils.ParseUUID(ctx.Param(g.ParamListID))
	itemID, errItem := utils.ParseUUID(ctx.Param(g.ParamItemID))
	if errList != nil || errItem != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return listID, services.ShoppingItem{}, false
	}
	it, err := h.service.GetShoppingItem(ctx, membershipFromContext(ctx), listID, itemID)
	if err != nil {
		pageError(ctx, err)
		return listID, it, false
	}
	return listID, it, true
}

func (h *DefaultHandler) PageShoppingItem(ctx *gin.Context) {
	if listID, it, ok := h.loadShoppingItem(ctx); ok {
		h.renderShoppingItem(ctx, http.StatusOK, listID, it, "")
	}
}

func (h *DefaultHandler) renderShoppingItem(ctx *gin.Context, status int, listID pgtype.UUID, it services.ShoppingItem, errorMessage string) {
	actor := membershipFromContext(ctx)
	files, err := h.service.AttachedFiles(ctx, actor,
		services.FileTarget{NoteID: pgtype.UUID{}, ListID: listID, ShoppingItemID: it.ID, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}})
	if err != nil {
		InternalError(ctx, err)
		return
	}
	// the receipt of the shop it was bought in
	var expense *services.Expense
	if it.ExpenseID.Valid {
		e, err := h.service.GetExpense(ctx, actor, it.ExpenseID)
		if err != nil {
			pageError(ctx, err)
			return
		}
		expense = &e
	}
	RenderTemplStatus(ctx, status, pages.ShoppingItemPage(actor, listID, it, expense, files, errorMessage))
}

func (h *DefaultHandler) SetShoppingItemDescription(ctx *gin.Context) {
	listID, it, ok := h.loadShoppingItem(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	err := h.service.SetShoppingItemDescription(ctx, actor, listID, it.ID, ctx.PostForm("description"))
	switch {
	case err == nil:
		Redirect(ctx, shoppingListPath(actor.HouseholdID, listID)+"/items/"+it.ID.String())
	case errors.Is(err, g.ErrorInvalidInput):
		h.renderShoppingItem(ctx, http.StatusUnprocessableEntity, listID, it,
			utils.T(ctx.Request.Context(), locales.LK_ShoppingDescriptionLength, services.MaxShoppingItemDescriptionLength))
	default:
		pageError(ctx, err)
	}
}

func (h *DefaultHandler) UploadShoppingItemFile(ctx *gin.Context) {
	listID, it, ok := h.loadShoppingItem(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	h.handleUpload(ctx, services.FileTarget{NoteID: pgtype.UUID{}, ListID: listID, ShoppingItemID: it.ID, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}},
		shoppingItemPath(actor.HouseholdID, listID, it.ID), func(status int, errorMessage string) {
			h.renderShoppingItem(ctx, status, listID, it, errorMessage)
		})
}

// -----------------------------------------------------------------------------

// PageFile shows a file inside the app, with a way back to where it is attached
func (h *DefaultHandler) PageFile(ctx *gin.Context) {
	fileID, err := utils.ParseUUID(ctx.Param(g.ParamFileID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	h.renderFile(ctx, http.StatusOK, fileID, "")
}

func (h *DefaultHandler) renderFile(ctx *gin.Context, status int, fileID pgtype.UUID, errorMessage string) {
	actor := membershipFromContext(ctx)
	f, err := h.service.GetFile(ctx, actor, fileID)
	if err != nil {
		pageError(ctx, err)
		return
	}
	RenderTemplStatus(ctx, status, pages.FilePage(actor, f, errorMessage))
}

// RenameFile shows the file again, an invalid name shows on it
func (h *DefaultHandler) RenameFile(ctx *gin.Context) {
	fileID, err := utils.ParseUUID(ctx.Param(g.ParamFileID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	err = h.service.RenameFile(ctx, actor, fileID, ctx.PostForm("name"))
	switch {
	case err == nil:
		Redirect(ctx, filesPath(actor.HouseholdID)+"/"+fileID.String())
	case errors.Is(err, g.ErrorInvalidInput):
		h.renderFile(ctx, http.StatusUnprocessableEntity, fileID,
			utils.T(ctx.Request.Context(), locales.LK_FilesNameInvalid, services.MaxFileNameLength))
	default:
		pageError(ctx, err)
	}
}

// DownloadFile streams the file from storage, images and PDFs are shown in the browser,
// anything else is downloaded, so an uploaded HTML or SVG never runs on this origin
func (h *DefaultHandler) DownloadFile(ctx *gin.Context) {
	h.serveFile(ctx, false)
}

// DownloadFileThumb streams the small JPEG of an image, 404 for other files
func (h *DefaultHandler) DownloadFileThumb(ctx *gin.Context) {
	h.serveFile(ctx, true)
}

func (h *DefaultHandler) serveFile(ctx *gin.Context, thumb bool) {
	fileID, err := utils.ParseUUID(ctx.Param(g.ParamFileID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	f, body, err := h.service.OpenFile(ctx, membershipFromContext(ctx), fileID, thumb)
	if err != nil {
		pageError(ctx, err)
		return
	}
	defer body.Close() //nolint:errcheck // read only
	disposition, mimeType, size := "attachment", f.Mime, f.SizeBytes
	if thumb {
		// the size is not stored, -1 sends it chunked
		mimeType, size = "image/jpeg", -1
	}
	if thumb || services.ImageMimes[mimeType] || mimeType == "application/pdf" {
		disposition = "inline"
	}
	ctx.DataFromReader(http.StatusOK, size, mimeType, body, map[string]string{
		// encodes names that are not ASCII as RFC 2231
		"Content-Disposition": cmp.Or(mime.FormatMediaType(disposition, map[string]string{"filename": f.Name}), disposition),
		// the contents never change, its id would, the name can: the Download link has it in the address
		"Cache-Control": "private, max-age=86400",
	})
}

// DeleteFile goes back to the page in the form, the Referer holds only the origin (Referrer-Policy strict-origin)
func (h *DefaultHandler) DeleteFile(ctx *gin.Context) {
	fileID, err := utils.ParseUUID(ctx.Param(g.ParamFileID))
	if err != nil {
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	actor := membershipFromContext(ctx)
	if err := h.service.DeleteFile(ctx, actor, fileID); err != nil {
		pageError(ctx, err)
		return
	}
	Redirect(ctx, safeNextPath(ctx.PostForm(g.QueryNext)))
}
