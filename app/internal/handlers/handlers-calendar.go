package handlers

import (
	"cmp"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/routes"
	"roommates/internal/services"
	"roommates/internal/utils"
	"roommates/web/pages"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// calendarAgendaDays is how many days the phone's agenda shows, and what "load more" adds
const calendarAgendaDays = 30

// calendarDayPath is the calendar on the month of day with day picked
func calendarDayPath(householdID pgtype.UUID, day time.Time) string {
	return routes.URL(routes.Calendar, householdID) + "?month=" + day.Format(pages.MonthLayout) + "&day=" + day.Format(time.DateOnly)
}

// eventFromContext is the :eventID of the request, or shows the not found page and returns false
func eventFromContext(ctx *gin.Context) (pgtype.UUID, bool) {
	id, err := utils.ParseUUID(ctx.Param(routes.ParamEventID))
	if err != nil {
		swapMainOnError(ctx)
		RenderTemplStatus(ctx, http.StatusNotFound, pages.NotFound())
		return id, false
	}
	return id, true
}

// eventOn is the ?on= day of the request, zero for none or a bad one: the event's first
func eventOn(ctx *gin.Context) time.Time {
	on, err := time.Parse(time.DateOnly, ctx.Query("on"))
	if err != nil {
		return time.Time{}
	}
	return on
}

// eventFormError is the message of a failed event form and the field it belongs to
func eventFormError(ctx *gin.Context, err error) (message, field string) {
	var fe services.FieldError
	errors.As(err, &fe)
	c := ctx.Request.Context()
	message = map[string]string{
		services.FieldEventTitle:       utils.T(c, locales.LK_CalendarTitleInvalid, services.MaxEventTitleLength),
		services.FieldEventLocation:    utils.T(c, locales.LK_CalendarLocationInvalid, services.MaxEventLocationLength),
		services.FieldEventNotes:       utils.T(c, locales.LK_CalendarNotesInvalid, services.MaxEventNotesLength),
		services.FieldEventStart:       utils.T(c, locales.LK_CalendarStartInvalid),
		services.FieldEventEnd:         utils.T(c, locales.LK_CalendarEndInvalid),
		services.FieldEventSchedule:    utils.T(c, locales.LK_ChoresScheduleInvalid),
		services.FieldEventRepeatUntil: utils.T(c, locales.LK_CalendarUntilInvalid),
		services.FieldEventRemind:      utils.T(c, locales.LK_CalendarRemindInvalid),
		services.FieldEventMembers:     utils.T(c, locales.LK_CalendarMembersInvalid),
	}[fe.Field]
	if message == "" {
		return utils.T(c, locales.LK_CalendarInvalid), ""
	}
	return message, fe.Field
}

// eventInputFromForm reads the event form: "schedule" is none, weekdays or every, the "unit" select
// only counts for every; the end day and the last day of a repeat only count with their tick,
// and the reminder is the one of the all-day or the timed set
func eventInputFromForm(ctx *gin.Context) services.EventInput {
	in := services.EventInput{
		Title:       ctx.PostForm("title"),
		Location:    ctx.PostForm("location"),
		Notes:       ctx.PostForm("notes"),
		StartDate:   strings.TrimSpace(ctx.PostForm("start_date")),
		AllDay:      ctx.PostForm("all_day") == "1",
		StartTime:   strings.TrimSpace(ctx.PostForm("start_time")),
		EndDate:     "",
		EndTime:     strings.TrimSpace(ctx.PostForm("end_time")),
		Unit:        "",
		Every:       0,
		Weekdays:    nil,
		RepeatUntil: "",
		Remind:      nil,
		Members:     nil,
	}
	if ctx.PostForm("end_on") == "1" {
		in.EndDate = strings.TrimSpace(ctx.PostForm("end_date"))
	}
	switch ctx.PostForm("schedule") {
	case g.ChoreUnitWeekdays:
		in.Unit = g.ChoreUnitWeekdays
	case "every":
		in.Unit = ctx.PostForm("unit")
	}
	in.Every, _ = strconv.Atoi(strings.TrimSpace(ctx.PostForm("every")))
	for _, v := range ctx.PostFormArray("weekday") {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 && n <= int(time.Saturday) {
			in.Weekdays = append(in.Weekdays, time.Weekday(n))
		}
	}
	if in.Unit != "" && ctx.PostForm("until_on") == "1" {
		in.RepeatUntil = strings.TrimSpace(ctx.PostForm("repeat_until"))
	}
	// "none" and anything else that is no minute count is no reminder
	if m, err := strconv.Atoi(ctx.PostForm(utils.IfElse(in.AllDay, "remind_allday", "remind_timed"))); err == nil {
		in.Remind = &m
	}
	for _, v := range ctx.PostFormArray("member") {
		if id, err := utils.ParseUUID(v); err == nil {
			in.Members = append(in.Members, id)
		}
	}
	return in
}

// newEventInput is the new event form: today, timed, not repeating
func newEventInput(today time.Time) services.EventInput {
	day := today.Format(time.DateOnly)
	return services.EventInput{
		Title:       "",
		Location:    "",
		Notes:       "",
		StartDate:   day,
		AllDay:      false,
		StartTime:   "",
		EndDate:     day,
		EndTime:     "",
		Unit:        "",
		Every:       1,
		Weekdays:    nil,
		RepeatUntil: "",
		Remind:      nil,
		Members:     nil,
	}
}

// renderCalendar shows the agenda of the next days and the month of ?month= with ?day= picked,
// notice is the toast of what just happened, a nil form is the empty new event form, folded
func (h *DefaultHandler) renderCalendar(ctx *gin.Context, status int, notice string, form *pages.EventForm) {
	actor := membershipFromContext(ctx)
	today, err := h.service.HouseholdToday(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	today = services.DateOnly(today)
	agenda, err := h.service.ListCalendar(ctx, actor, today, today.AddDate(0, 0, calendarAgendaDays))
	if err != nil {
		InternalError(ctx, err)
		return
	}
	month := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.UTC)
	if m, err := time.Parse(pages.MonthLayout, ctx.Query("month")); err == nil {
		month = m
	}
	next := month.AddDate(0, 1, 0)
	day := month
	if !today.Before(month) && today.Before(next) {
		day = today
	}
	if d, err := time.Parse(time.DateOnly, ctx.Query("day")); err == nil && !d.Before(month) && d.Before(next) {
		day = d
	}
	from, to := pages.CalendarGrid(month)
	items, err := h.service.ListCalendar(ctx, actor, from, to)
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
		form = &pages.EventForm{Input: newEventInput(today), Message: "", Field: "", Saved: false, Open: false, Scope: ""}
	}
	view := pages.CalendarView{
		Today:      today,
		Agenda:     agenda,
		AgendaNext: today.AddDate(0, 0, calendarAgendaDays),
		Month:      month,
		Day:        day,
		MonthItems: items,
	}
	RenderTemplStatus(ctx, status, pages.CalendarPage(actor, view, members, *form, notice))
}

// PageCalendar shows the household's calendar
func (h *DefaultHandler) PageCalendar(ctx *gin.Context) {
	h.renderCalendar(ctx, http.StatusOK, "", nil)
}

// CalendarAgenda is the next days of the phone's agenda from ?from= after the day ?last=, "load more"
// appends it
func (h *DefaultHandler) CalendarAgenda(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	today, err := h.service.HouseholdToday(ctx, actor)
	if err != nil {
		InternalError(ctx, err)
		return
	}
	today = services.DateOnly(today)
	from, err := time.Parse(time.DateOnly, ctx.Query("from"))
	if err != nil || from.Before(today) {
		from = today
	}
	items, err := h.service.ListCalendar(ctx, actor, from, from.AddDate(0, 0, calendarAgendaDays))
	if err != nil {
		InternalError(ctx, err)
		return
	}
	// the last day already shown, its month is not named again
	last, _ := time.Parse(time.DateOnly, ctx.Query("last"))
	RenderTempl(ctx, pages.CalendarMore(actor, today, from.AddDate(0, 0, calendarAgendaDays), last, items))
}

// CreateEvent adds the event and shows the calendar with it, an invalid form shows again with its error
func (h *DefaultHandler) CreateEvent(ctx *gin.Context) {
	actor := membershipFromContext(ctx)
	in := eventInputFromForm(ctx)
	_, err := h.service.CreateEvent(ctx, actor, in)
	c := ctx.Request.Context()
	switch {
	case err == nil:
		h.renderCalendar(ctx, http.StatusOK, utils.T(c, locales.LK_CalendarNoticeAdded), nil)
	case errors.Is(err, services.ErrorTooManyEvents):
		h.renderCalendar(ctx, http.StatusUnprocessableEntity, "",
			&pages.EventForm{Input: in, Message: utils.T(c, locales.LK_CalendarTooMany, g.MaxUpcomingEvents), Field: "", Saved: false, Open: true, Scope: ""})
	case errors.Is(err, g.ErrorInvalidInput):
		message, field := eventFormError(ctx, err)
		h.renderCalendar(ctx, http.StatusUnprocessableEntity, "", &pages.EventForm{Input: in, Message: message, Field: field, Saved: false, Open: true, Scope: ""})
	default:
		InternalError(ctx, err)
	}
}

// renderEvent shows the event's page for its occurrence on on, a nil form is the event as it is
// stored, with "Saved" when saved; fileError is the outcome of an upload
func (h *DefaultHandler) renderEvent(ctx *gin.Context, status int, id pgtype.UUID, on time.Time, form *pages.EventForm, saved bool, fileError string) {
	actor := membershipFromContext(ctx)
	d, err := h.service.GetOccurrence(ctx, actor, id, on)
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
		form = &pages.EventForm{Input: d.Input, Message: "", Field: "", Saved: saved, Open: false, Scope: ""}
	}
	RenderTemplStatus(ctx, status, pages.EventPage(actor, d, members, *form, fileError))
}

// PageEvent shows the :eventID event on ?on=
func (h *DefaultHandler) PageEvent(ctx *gin.Context) {
	id, ok := eventFromContext(ctx)
	if !ok {
		return
	}
	h.renderEvent(ctx, http.StatusOK, id, eventOn(ctx), nil, false, "")
}

// UpdateEvent saves the form for the occurrence on ?on= with the posted "scope" and shows the page
// again with "Saved" by the button; an occurrence that moved to another series or ended with the
// edit is not on its page any more, so that goes to the calendar on its day
func (h *DefaultHandler) UpdateEvent(ctx *gin.Context) {
	id, ok := eventFromContext(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	d, err := h.service.GetOccurrence(ctx, actor, id, eventOn(ctx))
	if err != nil {
		pageError(ctx, err)
		return
	}
	in := eventInputFromForm(ctx)
	scope := cmp.Or(ctx.PostForm("scope"), services.EventScopeAll)
	err = h.service.UpdateEvent(ctx, actor, id, d.On, scope, in)
	switch {
	case err == nil:
		moved, perr := time.Parse(time.DateOnly, in.StartDate)
		if perr != nil {
			moved = d.On
		}
		if d.Repeats && scope != services.EventScopeAll {
			Redirect(ctx, calendarDayPath(actor.HouseholdID, moved))
			return
		}
		// a weekdays series may start on another day than the one posted
		if _, err := h.service.GetOccurrence(ctx, actor, id, moved); err != nil {
			Redirect(ctx, calendarDayPath(actor.HouseholdID, moved))
			return
		}
		h.renderEvent(ctx, http.StatusOK, id, moved, nil, true, "")
	case errors.Is(err, g.ErrorInvalidInput):
		message, field := eventFormError(ctx, err)
		h.renderEvent(ctx, http.StatusUnprocessableEntity, id, d.On, &pages.EventForm{Input: in, Message: message, Field: field, Saved: false, Open: false, Scope: scope}, false, "")
	default:
		pageError(ctx, err)
	}
}

// DeleteEvent deletes the occurrence on ?on= with the posted "scope" and goes to the calendar, the
// page it was on is gone
func (h *DefaultHandler) DeleteEvent(ctx *gin.Context) {
	id, ok := eventFromContext(ctx)
	if !ok {
		return
	}
	actor := membershipFromContext(ctx)
	d, err := h.service.GetOccurrence(ctx, actor, id, eventOn(ctx))
	if err != nil {
		pageError(ctx, err)
		return
	}
	if err := h.service.DeleteEvent(ctx, actor, id, d.On, cmp.Or(ctx.PostForm("scope"), services.EventScopeAll)); err != nil {
		pageError(ctx, err)
		return
	}
	Redirect(ctx, calendarDayPath(actor.HouseholdID, d.On))
}

// UploadEventFile attaches an upload to the :eventID event and goes back to the occurrence the form
// was on, the upload form itself has no ?on=
func (h *DefaultHandler) UploadEventFile(ctx *gin.Context) {
	t, back, ok := h.EventFileTarget(ctx)
	if !ok {
		return
	}
	var on time.Time
	if u, err := url.Parse(ctx.Request.Referer()); err == nil && u.Path == back {
		if d, err := time.Parse(time.DateOnly, u.Query().Get("on")); err == nil {
			on = d
			back += "?on=" + d.Format(time.DateOnly)
		}
	}
	h.handleUpload(ctx, t, back, func(status int, errorMessage string) {
		h.renderEvent(ctx, status, t.EventID, on, nil, false, errorMessage)
	})
}

// EventFileTarget is the :eventID event for the routes that attach files, the service says whether
// the actor may change it
func (h *DefaultHandler) EventFileTarget(ctx *gin.Context) (services.FileTarget, string, bool) {
	id, ok := eventFromContext(ctx)
	return services.FileTarget{NoteID: pgtype.UUID{}, ListID: pgtype.UUID{}, ShoppingItemID: pgtype.UUID{}, ExpenseID: pgtype.UUID{}, MessageID: pgtype.UUID{}, EventID: id},
		routes.URL(routes.CalendarEvent, membershipFromContext(ctx).HouseholdID, id), ok
}
