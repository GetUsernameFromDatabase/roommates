package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/services"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func choresPath(householdID pgtype.UUID) string {
	return householdPath(householdID) + "/chores"
}

// choreFromContext is the :choreID of the request, or shows the not found page and returns false
func choreFromContext(ctx *gin.Context) (pgtype.UUID, bool) {
	id, err := utils.ParseUUID(ctx.Param(g.ParamChoreID))
	if err != nil {
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return id, false
	}
	return id, true
}

// choreFormError is the message of a failed chore form and the field it belongs to
func choreFormError(ctx *gin.Context, err error) (message, field string) {
	var fe services.FieldError
	errors.As(err, &fe)
	c := ctx.Request.Context()
	message = map[string]string{
		services.FieldChoreTitle:    utils.T(c, locales.LK_ChoresTitleInvalid, services.MaxChoreTitleLength),
		services.FieldChoreRotation: utils.T(c, locales.LK_ChoresRotationInvalid),
		services.FieldChoreSchedule: utils.T(c, locales.LK_ChoresScheduleInvalid),
		services.FieldChoreFirstDue: utils.T(c, locales.LK_ChoresFirstDueInvalid),
		services.FieldChoreMembers:  utils.T(c, locales.LK_ChoresMembersInvalid),
		services.FieldChoreAssignee: utils.T(c, locales.LK_ChoresAssigneeInvalid),
	}[fe.Field]
	if message == "" {
		return utils.T(c, locales.LK_ChoresInvalid), ""
	}
	return message, fe.Field
}

// choreInputFromForm reads the chore form, the members in the order they are posted; "schedule" is
// weekdays or every, the "unit" select only counts for every
func choreInputFromForm(ctx *gin.Context) services.ChoreInput {
	in := services.ChoreInput{
		Title:    ctx.PostForm("title"),
		Rotation: ctx.PostForm("rotation"),
		Unit:     utils.IfElse(ctx.PostForm("schedule") == g.ChoreUnitWeekdays, g.ChoreUnitWeekdays, ctx.PostForm("unit")),
		FirstDue: strings.TrimSpace(ctx.PostForm("first_due")),
		Every:    0,
		Weekdays: nil,
		Members:  nil,
		Assignee: pgtype.UUID{},
	}
	in.Every, _ = strconv.Atoi(strings.TrimSpace(ctx.PostForm("every")))
	for _, v := range ctx.PostFormArray("weekday") {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 && n <= int(time.Saturday) {
			in.Weekdays = append(in.Weekdays, time.Weekday(n))
		}
	}
	for _, v := range ctx.PostFormArray("member") {
		if id, err := utils.ParseUUID(v); err == nil {
			in.Members = append(in.Members, id)
		}
	}
	// "none" and anything else that is no member is nobody's turn
	if id, err := utils.ParseUUID(ctx.PostForm("assignee")); err == nil {
		in.Assignee = id
	}
	return in
}

// choreInput is the form of an existing chore, its next due date is the first one
func choreInput(c services.Chore) services.ChoreInput {
	in := services.ChoreInput{
		Title:    c.Title,
		Rotation: c.Rotation,
		Unit:     c.Unit,
		Every:    c.Every,
		Weekdays: c.Weekdays,
		FirstDue: c.NextDue.Format(services.DateLayout),
		Members:  nil,
		Assignee: c.AssigneeID,
	}
	for _, m := range c.Members {
		in.Members = append(in.Members, m.UserID)
	}
	return in
}

// renderChores shows the list, notice is the toast of what just happened, a nil form is the empty
// new chore form, folded
func (h *DefaultHandler) renderChores(ctx *gin.Context, status int, notice string, form *pages.ChoreForm) {
	actor := membershipFromContext(ctx)
	chores, err := h.service.ListChores(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	members, err := h.service.ListMembers(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	if form == nil {
		// without JS the picker sends no day, so the date is a good one to start from
		today, err := h.service.HouseholdToday(ctx, actor)
		if err != nil {
			InternalError(ctx, err)
			return
		}
		in := services.ChoreInput{
			Title:    "",
			Rotation: g.RotationFixed,
			Unit:     g.ChoreUnitWeekdays,
			Every:    1,
			Weekdays: nil,
			FirstDue: today.Format(services.DateLayout),
			Members:  nil,
			Assignee: pgtype.UUID{},
		}
		for _, m := range members {
			in.Members = append(in.Members, m.UserID)
		}
		form = &pages.ChoreForm{Input: in, Message: "", Field: "", Saved: false, Open: false}
	}
	RenderTemplStatus(ctx, status, pages.ChoresPage(actor, chores, members, *form, notice))
}

// PageChores lists the household's chores
func (h *DefaultHandler) PageChores(ctx *gin.Context) {
	h.renderChores(ctx, http.StatusOK, "", nil)
}

// CreateChore adds the chore and shows the list with it, an invalid form shows again with its error
func (h *DefaultHandler) CreateChore(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	in := choreInputFromForm(ctx)
	_, err := h.service.CreateChore(ctx, actor, in)
	switch {
	case err == nil:
		h.renderChores(ctx, http.StatusOK, utils.T(ctx.Request.Context(), locales.LK_ChoresNoticeAdded), nil)
	case errors.Is(err, g.ErrorInvalidInput):
		message, field := choreFormError(ctx, err)
		h.renderChores(ctx, http.StatusUnprocessableEntity, "", &pages.ChoreForm{Input: in, Message: message, Field: field, Saved: false, Open: true})
	default:
		InternalError(ctx, err)
	}
}

// renderChore shows the chore's page, a nil form is the chore as it is stored, with "Saved" when saved
func (h *DefaultHandler) renderChore(ctx *gin.Context, status int, id pgtype.UUID, notice string, form *pages.ChoreForm, saved bool) {
	actor := membershipFromContext(ctx)
	d, err := h.service.GetChore(ctx, actor, id)
	if err != nil {
		pageError(ctx, err)
		return
	}
	members, err := h.service.ListMembers(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	if form == nil {
		form = &pages.ChoreForm{Input: choreInput(d.Chore), Message: "", Field: "", Saved: saved, Open: false}
	}
	RenderTemplStatus(ctx, status, pages.ChorePage(actor, d, members, *form, notice))
}

// PageChore shows the :choreID chore
func (h *DefaultHandler) PageChore(ctx *gin.Context) {
	id, ok := choreFromContext(ctx)
	if !ok {
		return
	}
	h.renderChore(ctx, http.StatusOK, id, "", nil, false)
}

// UpdateChore saves the form and shows the page again with "Saved" by the button
func (h *DefaultHandler) UpdateChore(ctx *gin.Context) {
	id, ok := choreFromContext(ctx)
	if !ok {
		return
	}
	in := choreInputFromForm(ctx)
	err := h.service.UpdateChore(ctx, membershipFromContext(ctx), id, in)
	switch {
	case err == nil:
		h.renderChore(ctx, http.StatusOK, id, "", nil, true)
	case errors.Is(err, g.ErrorInvalidInput):
		message, field := choreFormError(ctx, err)
		h.renderChore(ctx, http.StatusUnprocessableEntity, id, "", &pages.ChoreForm{Input: in, Message: message, Field: field, Saved: false, Open: false}, false)
	default:
		pageError(ctx, err)
	}
}

// DeleteChore deletes the chore and goes to the list, the page it was on is gone
func (h *DefaultHandler) DeleteChore(ctx *gin.Context) {
	id, ok := choreFromContext(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	if err := h.service.DeleteChore(ctx, actor, id); err != nil {
		pageError(ctx, err)
		return
	}
	Redirect(ctx, choresPath(actor.HouseholdID))
}

// choreTurnAction runs an action of the card on the posted turn ("due" and "assignee") and shows the
// list with a toast of what it did, an action on a turn that moved on or that cannot be done only
// shows the list
func (h *DefaultHandler) choreTurnAction(ctx *gin.Context,
	run func(actor services.Membership, id pgtype.UUID, turn services.ChoreTurn) error,
	notice func(c services.Chore) string) {
	id, ok := choreFromContext(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	var turn services.ChoreTurn
	if day, err := time.Parse(services.DateLayout, ctx.PostForm("due")); err == nil {
		turn.Due = pgtype.Date{Time: day, InfinityModifier: pgtype.Finite, Valid: true}
	}
	// "" is nobody's turn, which stays an invalid UUID
	turn.Assignee, _ = utils.ParseUUID(ctx.PostForm("assignee"))
	err := run(actor, id, turn)
	if errors.Is(err, g.ErrorInvalidInput) || errors.Is(err, services.ErrChoreTurnMoved) {
		h.renderChores(ctx, http.StatusOK, "", nil)
		return
	}
	if err != nil {
		pageError(ctx, err)
		return
	}
	d, err := h.service.GetChore(ctx, actor, id)
	if err != nil {
		pageError(ctx, err)
		return
	}
	h.renderChores(ctx, http.StatusOK, notice(d.Chore), nil)
}

// choreNextNotice is "Done. Next: Jaan, 21.11.", without anyone's turn only the date
func choreNextNotice(ctx *gin.Context, withWho, nobody locales.LK, c services.Chore) string {
	c2 := ctx.Request.Context()
	date := utils.FormatDate(c2, c.NextDue)
	if !c.AssigneeID.Valid {
		return utils.T(c2, nobody, date)
	}
	return utils.T(c2, withWho, c.AssigneeName, date)
}

// DoneChore marks the turn done
func (h *DefaultHandler) DoneChore(ctx *gin.Context) {
	h.choreTurnAction(ctx, func(actor services.Membership, id pgtype.UUID, turn services.ChoreTurn) error {
		return h.service.DoneChore(ctx, actor, id, turn)
	}, func(c services.Chore) string {
		return choreNextNotice(ctx, locales.LK_ChoresNoticeDone, locales.LK_ChoresNoticeDoneNobody, c)
	})
}

// SkipChore drops the turn's due date
func (h *DefaultHandler) SkipChore(ctx *gin.Context) {
	h.choreTurnAction(ctx, func(actor services.Membership, id pgtype.UUID, turn services.ChoreTurn) error {
		return h.service.SkipChore(ctx, actor, id, turn)
	}, func(c services.Chore) string {
		return choreNextNotice(ctx, locales.LK_ChoresNoticeSkipped, locales.LK_ChoresNoticeSkippedNobody, c)
	})
}

// PassChore hands the turn on to the next in the rota
func (h *DefaultHandler) PassChore(ctx *gin.Context) {
	h.choreTurnAction(ctx, func(actor services.Membership, id pgtype.UUID, turn services.ChoreTurn) error {
		return h.service.PassChore(ctx, actor, id, turn)
	}, func(c services.Chore) string {
		return utils.T(ctx.Request.Context(), locales.LK_ChoresNoticePassed, c.AssigneeName)
	})
}

// AssignChore gives the turn to the posted "user" and shows the chore's rota again
func (h *DefaultHandler) AssignChore(ctx *gin.Context) {
	id, ok := choreFromContext(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	userID, err := utils.ParseUUID(ctx.PostForm("user"))
	if err != nil {
		h.renderChore(ctx, http.StatusOK, id, "", nil, false)
		return
	}
	err = h.service.AssignChore(ctx, actor, id, userID)
	if errors.Is(err, g.ErrorInvalidInput) {
		h.renderChore(ctx, http.StatusOK, id, "", nil, false)
		return
	}
	if err != nil {
		pageError(ctx, err)
		return
	}
	d, err := h.service.GetChore(ctx, actor, id)
	if err != nil {
		pageError(ctx, err)
		return
	}
	h.renderChore(ctx, http.StatusOK, id, utils.T(ctx.Request.Context(), locales.LK_ChoresNoticeAssigned, d.Chore.AssigneeName), nil, false)
}

// MoveChoreMember moves :userID one place up or down, the posted "dir", and shows the rota again
func (h *DefaultHandler) MoveChoreMember(ctx *gin.Context) {
	id, ok := choreFromContext(ctx)
	if !ok {
		return
	}
	userID, err := utils.ParseUUID(ctx.Param(g.ParamUserID))
	if err != nil {
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return
	}
	err = h.service.MoveChoreMember(ctx, membershipFromContext(ctx), id, userID, ctx.PostForm("dir") == "up")
	if err != nil && !errors.Is(err, g.ErrorInvalidInput) {
		pageError(ctx, err)
		return
	}
	h.renderChore(ctx, http.StatusOK, id, "", nil, false)
}
