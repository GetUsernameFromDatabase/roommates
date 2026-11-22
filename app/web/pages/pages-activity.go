package pages

import (
	"context"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/routes"
	"roommates/internal/services"
	"roommates/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var activityKeys = map[string]locales.LK{
	g.ActExpenseAdded:         locales.LK_ActivityExpenseAdded,
	g.ActExpenseEdited:        locales.LK_ActivityExpenseEdited,
	g.ActExpenseDeleted:       locales.LK_ActivityExpenseDeleted,
	g.ActSharePaid:            locales.LK_ActivitySharePaid,
	g.ActShareUnpaid:          locales.LK_ActivityShareUnpaid,
	g.ActSettlementRequested:  locales.LK_ActivitySettlementRequested,
	g.ActSettlementConfirmed:  locales.LK_ActivitySettlementConfirmed,
	g.ActSettlementDeclined:   locales.LK_ActivitySettlementDeclined,
	g.ActRecurringAdded:       locales.LK_ActivityRecurringAdded,
	g.ActRecurringEdited:      locales.LK_ActivityRecurringEdited,
	g.ActRecurringDeleted:     locales.LK_ActivityRecurringDeleted,
	g.ActRecurringSkipped:     locales.LK_ActivityRecurringSkipped,
	g.ActListCreated:          locales.LK_ActivityListCreated,
	g.ActListRenamed:          locales.LK_ActivityListRenamed,
	g.ActListDeleted:          locales.LK_ActivityListDeleted,
	g.ActItemBought:           locales.LK_ActivityItemBought,
	g.ActItemsCleared:         locales.LK_ActivityItemsCleared,
	g.ActNoteCreated:          locales.LK_ActivityNoteCreated,
	g.ActNoteEdited:           locales.LK_ActivityNoteEdited,
	g.ActNoteDeleted:          locales.LK_ActivityNoteDeleted,
	g.ActFileUploaded:         locales.LK_ActivityFileUploaded,
	g.ActFileRemoved:          locales.LK_ActivityFileRemoved,
	g.ActMemberJoined:         locales.LK_ActivityMemberJoined,
	g.ActMemberLeft:           locales.LK_ActivityMemberLeft,
	g.ActMemberRemoved:        locales.LK_ActivityMemberRemoved,
	g.ActMemberPromoted:       locales.LK_ActivityMemberPromoted,
	g.ActOwnershipTransferred: locales.LK_ActivityOwnershipTransferred,
	g.ActSteppedDown:          locales.LK_ActivitySteppedDown,
	g.ActHouseholdRenamed:     locales.LK_ActivityHouseholdRenamed,
	g.ActCurrencyChanged:      locales.LK_ActivityCurrencyChanged,
	g.ActHouseholdDeleted:     locales.LK_ActivityHouseholdDeleted,
	g.ActHouseholdRestored:    locales.LK_ActivityHouseholdRestored,
	g.ActChoreAdded:           locales.LK_ActivityChoreAdded,
	g.ActChoreDeleted:         locales.LK_ActivityChoreDeleted,
	g.ActChoreDone:            locales.LK_ActivityChoreDone,
	g.ActChoreSkipped:         locales.LK_ActivityChoreSkipped,
	g.ActChorePassed:          locales.LK_ActivityChorePassed,
	g.ActEventAdded:           locales.LK_ActivityEventAdded,
	g.ActEventDeleted:         locales.LK_ActivityEventDeleted,
}

// activityText is the sentence of a row, "Mari added the expense Rent · 650,00 €",
// the viewer's own rows use the "-you" variant of the key ("You added the expense Rent")
func activityText(ctx context.Context, actor services.Membership, r services.ActivityRow) string {
	someone := utils.T(ctx, locales.LK_ActivitySomeone)
	who, target := someone, someone
	key := activityKeys[r.Kind]
	if r.ActorName != nil {
		who = utils.MemberName(ctx, *r.ActorName)
		if r.ActorID == actor.UserID {
			key += "-you"
		}
	}
	if r.TargetName != nil {
		target = utils.MemberName(ctx, *r.TargetName)
	}
	s := utils.T(ctx, key, who, r.Title, target)
	if r.AmountCents != nil && r.Currency != nil {
		s += " · " + utils.FormatMoney(ctx, *r.AmountCents, *r.Currency)
	}
	return s
}

// activityLink is where the row's subject is, "" when it is gone or has no page
func activityLink(householdID pgtype.UUID, r services.ActivityRow) string {
	id := r.SubjectID
	switch {
	case r.ExpenseExists:
		return routes.URL(routes.Expense, householdID, id)
	case r.ListExists:
		return routes.URL(routes.ShoppingList, householdID, id)
	case r.NoteExists:
		return routes.URL(routes.Note, householdID, id)
	case r.FileExists:
		return routes.URL(routes.File, householdID, id)
	case r.RecurringExists:
		return routes.URL(routes.Recurring, householdID)
	case r.ChoreExists:
		return routes.URL(routes.Chore, householdID, id)
	case r.EventExists:
		return routes.URL(routes.CalendarEvent, householdID, id)
	}
	return ""
}

// activityMoreURL loads the rows older than the last of rows
func activityMoreURL(householdID pgtype.UUID, rows []services.ActivityRow) string {
	before := rows[len(rows)-1].CreatedAt.Time.UTC().Format(time.RFC3339Nano)
	return routes.URL(routes.HouseholdActivity, householdID) + "?before=" + before
}
