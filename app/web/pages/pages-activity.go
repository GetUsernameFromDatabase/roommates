package pages

import (
	"context"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"roommates/internal/routes"
	"roommates/internal/services"
	"roommates/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

// activityKey is the sentence of a kind, said of someone else and of the viewer
type activityKey struct{ them, you locales.LK }

var activityKeys = map[string]activityKey{
	g.ActExpenseAdded:         {locales.LK_ActivityExpenseAdded, locales.LK_ActivityExpenseAddedYou},
	g.ActExpenseEdited:        {locales.LK_ActivityExpenseEdited, locales.LK_ActivityExpenseEditedYou},
	g.ActExpenseDeleted:       {locales.LK_ActivityExpenseDeleted, locales.LK_ActivityExpenseDeletedYou},
	g.ActSharePaid:            {locales.LK_ActivitySharePaid, locales.LK_ActivitySharePaidYou},
	g.ActShareUnpaid:          {locales.LK_ActivityShareUnpaid, locales.LK_ActivityShareUnpaidYou},
	g.ActSettlementRequested:  {locales.LK_ActivitySettlementRequested, locales.LK_ActivitySettlementRequestedYou},
	g.ActSettlementConfirmed:  {locales.LK_ActivitySettlementConfirmed, locales.LK_ActivitySettlementConfirmedYou},
	g.ActSettlementDeclined:   {locales.LK_ActivitySettlementDeclined, locales.LK_ActivitySettlementDeclinedYou},
	g.ActRecurringAdded:       {locales.LK_ActivityRecurringAdded, locales.LK_ActivityRecurringAddedYou},
	g.ActRecurringEdited:      {locales.LK_ActivityRecurringEdited, locales.LK_ActivityRecurringEditedYou},
	g.ActRecurringDeleted:     {locales.LK_ActivityRecurringDeleted, locales.LK_ActivityRecurringDeletedYou},
	g.ActRecurringSkipped:     {locales.LK_ActivityRecurringSkipped, locales.LK_ActivityRecurringSkippedYou},
	g.ActListCreated:          {locales.LK_ActivityListCreated, locales.LK_ActivityListCreatedYou},
	g.ActListRenamed:          {locales.LK_ActivityListRenamed, locales.LK_ActivityListRenamedYou},
	g.ActListDeleted:          {locales.LK_ActivityListDeleted, locales.LK_ActivityListDeletedYou},
	g.ActItemBought:           {locales.LK_ActivityItemBought, locales.LK_ActivityItemBoughtYou},
	g.ActItemsCleared:         {locales.LK_ActivityItemsCleared, locales.LK_ActivityItemsClearedYou},
	g.ActNoteCreated:          {locales.LK_ActivityNoteCreated, locales.LK_ActivityNoteCreatedYou},
	g.ActNoteEdited:           {locales.LK_ActivityNoteEdited, locales.LK_ActivityNoteEditedYou},
	g.ActNoteDeleted:          {locales.LK_ActivityNoteDeleted, locales.LK_ActivityNoteDeletedYou},
	g.ActFileUploaded:         {locales.LK_ActivityFileUploaded, locales.LK_ActivityFileUploadedYou},
	g.ActFileRemoved:          {locales.LK_ActivityFileRemoved, locales.LK_ActivityFileRemovedYou},
	g.ActMemberJoined:         {locales.LK_ActivityMemberJoined, locales.LK_ActivityMemberJoinedYou},
	g.ActMemberLeft:           {locales.LK_ActivityMemberLeft, locales.LK_ActivityMemberLeftYou},
	g.ActMemberRemoved:        {locales.LK_ActivityMemberRemoved, locales.LK_ActivityMemberRemovedYou},
	g.ActMemberPromoted:       {locales.LK_ActivityMemberPromoted, locales.LK_ActivityMemberPromotedYou},
	g.ActOwnershipTransferred: {locales.LK_ActivityOwnershipTransferred, locales.LK_ActivityOwnershipTransferredYou},
	g.ActSteppedDown:          {locales.LK_ActivitySteppedDown, locales.LK_ActivitySteppedDownYou},
	g.ActHouseholdRenamed:     {locales.LK_ActivityHouseholdRenamed, locales.LK_ActivityHouseholdRenamedYou},
	g.ActCurrencyChanged:      {locales.LK_ActivityCurrencyChanged, locales.LK_ActivityCurrencyChangedYou},
	g.ActHouseholdDeleted:     {locales.LK_ActivityHouseholdDeleted, locales.LK_ActivityHouseholdDeletedYou},
	g.ActHouseholdRestored:    {locales.LK_ActivityHouseholdRestored, locales.LK_ActivityHouseholdRestoredYou},
	g.ActChoreAdded:           {locales.LK_ActivityChoreAdded, locales.LK_ActivityChoreAddedYou},
	g.ActChoreDeleted:         {locales.LK_ActivityChoreDeleted, locales.LK_ActivityChoreDeletedYou},
	g.ActChoreDone:            {locales.LK_ActivityChoreDone, locales.LK_ActivityChoreDoneYou},
	g.ActChoreSkipped:         {locales.LK_ActivityChoreSkipped, locales.LK_ActivityChoreSkippedYou},
	g.ActChorePassed:          {locales.LK_ActivityChorePassed, locales.LK_ActivityChorePassedYou},
	g.ActEventAdded:           {locales.LK_ActivityEventAdded, locales.LK_ActivityEventAddedYou},
	g.ActEventDeleted:         {locales.LK_ActivityEventDeleted, locales.LK_ActivityEventDeletedYou},
}

// activityText is the sentence of a row, "Mari added the expense Rent · 650,00 €",
// the viewer's own rows use the "you" key ("You added the expense Rent")
func activityText(ctx context.Context, actor services.Membership, r services.ActivityRow) string {
	someone := utils.T(ctx, locales.LK_ActivitySomeone)
	who, target := someone, someone
	key := activityKeys[r.Kind].them
	if r.ActorName != nil {
		who = utils.MemberName(ctx, *r.ActorName)
		if r.ActorID == actor.UserID {
			key = activityKeys[r.Kind].you
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
	last := rows[len(rows)-1]
	return routes.URL(routes.HouseholdActivity, householdID) + "?before=" + services.Cursor{At: last.CreatedAt, ID: last.ID}.String()
}
