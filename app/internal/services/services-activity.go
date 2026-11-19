package services

import (
	"context"
	"fmt"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"

	"github.com/jackc/pgx/v5/pgtype"
)

// ActivityPageSize is how many activity rows a page has, older ones are a "load more" away
const ActivityPageSize = 50

type ActivityRow = dbqueries.ListActivityRow

// ActivityEntry is one thing that happened, Title is a name to show when the subject is gone
type ActivityEntry struct {
	Kind         string
	SubjectID    pgtype.UUID
	TargetUserID pgtype.UUID
	Title        string
	AmountCents  *int64
	Currency     string
}

// logActivity records e as done by actor, callers pass the q of their own transaction so a
// change and its entry stand or fall together; it is an error for e.Kind not to be in [g.ActivityGroup]
func logActivity(ctx context.Context, q *dbqueries.Queries, actor Membership, e ActivityEntry) error {
	grp, ok := g.ActivityGroup[e.Kind]
	if !ok {
		return fmt.Errorf("unknown activity kind %q", e.Kind)
	}
	var currency *string
	if e.Currency != "" {
		currency = &e.Currency
	}
	return q.InsertActivity(ctx, dbqueries.InsertActivityParams{
		HouseholdID:  actor.HouseholdID,
		ActorID:      actor.UserID,
		Grp:          grp,
		Kind:         e.Kind,
		SubjectID:    e.SubjectID,
		TargetUserID: e.TargetUserID,
		Title:        e.Title,
		AmountCents:  e.AmountCents,
		Currency:     currency,
	})
}

// expenseOf reads the expense a log entry is about, [g.ErrorNotFound] if it is not in the household
func expenseOf(ctx context.Context, q *dbqueries.Queries, householdID, id pgtype.UUID) (dbqueries.GetExpenseRow, error) {
	e, err := q.GetExpense(ctx, dbqueries.GetExpenseParams{ID: id, HouseholdID: householdID})
	return e, noRowsAs(err, g.ErrorNotFound)
}

// listName reads the name of a shopping list, [g.ErrorNotFound] if it is not in the household
func listName(ctx context.Context, q *dbqueries.Queries, householdID, id pgtype.UUID) (string, error) {
	l, err := q.GetShoppingList(ctx, dbqueries.GetShoppingListParams{ID: id, HouseholdID: householdID})
	return l.Name, noRowsAs(err, g.ErrorNotFound)
}

// logShare logs that the share of userID in an expense was marked paid or unpaid
func logShare(ctx context.Context, q *dbqueries.Queries, actor Membership, kind string, expenseID, userID pgtype.UUID, cents int64) error {
	e, err := expenseOf(ctx, q, actor.HouseholdID, expenseID)
	if err != nil {
		return err
	}
	return logActivity(ctx, q, actor, ActivityEntry{Kind: kind, SubjectID: expenseID, TargetUserID: userID, Title: e.Title, AmountCents: &cents, Currency: e.Currency})
}

// ListActivity returns [ActivityPageSize] rows of the actor's household, newest first, older than before
// unless it is not valid; rows about what the actor may not see are left out
//
// shortcut: pages are cut by time, rows with the very same time (one transaction) can fall between
// two pages, when that shows switch to a (created_at, id) cursor
func (s *DefaultService) ListActivity(ctx context.Context, actor Membership, before pgtype.Timestamptz) ([]ActivityRow, error) {
	return s.db.ListActivity(ctx, dbqueries.ListActivityParams{
		HouseholdID: actor.HouseholdID,
		BeforeAt:    before,
		UserID:      actor.UserID,
		PageSize:    ActivityPageSize,
	})
}

// PurgeActivity deletes the activity older than its group's retention in [g.ActivityRetention]
func (s *DefaultService) PurgeActivity(ctx context.Context) error {
	return s.db.PurgeActivity(ctx, dbqueries.PurgeActivityParams{
		MoneySecs:    g.ActivityRetention[g.ActivityMoney].Seconds(),
		ShoppingSecs: g.ActivityRetention[g.ActivityShopping].Seconds(),
		NotesSecs:    g.ActivityRetention[g.ActivityNotes].Seconds(),
		ChoresSecs:   g.ActivityRetention[g.ActivityChores].Seconds(),
		MembersSecs:  g.ActivityRetention[g.ActivityMembers].Seconds(),
	})
}
