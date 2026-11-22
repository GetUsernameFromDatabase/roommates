package services

import (
	"context"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/push"
	"roommates/internal/routes"
	"roommates/internal/utils"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// RecurringInput is the expense form with how often it repeats, one of [g.Intervals]
type RecurringInput struct {
	ExpenseInput
	Interval string
}

// RecurringRow is a schedule, Shares are filled in only by [DefaultService.GetRecurring]
type RecurringRow struct {
	dbqueries.ListRecurringRow
	Shares []dbqueries.ListRecurringSharesRow
}

// Input fills the expense form from r, dated on its next due day, with its split
func (r RecurringRow) Input() RecurringInput {
	in := RecurringInput{
		Title:      r.Title,
		Amount:     FormatCentsInput(r.AmountCents),
		Category:   r.Category,
		IncurredOn: r.NextDue.Time.Format(time.DateOnly),
		Custom:     r.Custom,
		Members:    nil,
		Shares:     map[pgtype.UUID]string{},
		ListID:     pgtype.UUID{},
		ItemID:     pgtype.UUID{},
		Interval:   r.Interval,
	}
	for _, sh := range r.Shares {
		in.Members = append(in.Members, sh.UserID)
		if sh.AmountCents != nil {
			in.Shares[sh.UserID] = FormatCentsInput(*sh.AmountCents)
		}
	}
	return in
}

// addIntervals is anchor moved k intervals on, a month without anchor's day gives its last day
func addIntervals(anchor time.Time, interval string, k int) time.Time {
	y, m, d := anchor.Date()
	switch interval {
	case g.IntervalWeek:
		return anchor.AddDate(0, 0, 7*k)
	case g.IntervalMonth:
		m += time.Month(k)
	default:
		y += k
	}
	// day 0 of the next month is the last day of this one, time.Date normalises the month
	last := time.Date(y, m+1, 0, 0, 0, 0, 0, anchor.Location()).Day()
	return time.Date(y, m, min(d, last), 0, 0, 0, 0, anchor.Location())
}

// nextDue is the first due day of a schedule from anchor that comes after after,
// counted from anchor so a clamped day goes back to anchor's day the month after
func nextDue(anchor time.Time, interval string, after time.Time) time.Time {
	for k := 1; ; k++ {
		if d := addIntervals(anchor, interval, k); d.After(after) {
			return d
		}
	}
}

func dateOf(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, InfinityModifier: pgtype.Finite, Valid: true}
}

// insertRecurring starts a schedule of p from its date, call it in the expense's transaction
func insertRecurring(ctx context.Context, q *dbqueries.Queries, p dbqueries.InsertExpenseParams, custom bool, interval string, shares []expenseShare) (pgtype.UUID, error) {
	id, err := q.InsertRecurring(ctx, dbqueries.InsertRecurringParams{
		HouseholdID: p.HouseholdID,
		PayerID:     p.PayerID,
		Category:    p.Category,
		Title:       p.Title,
		AmountCents: p.AmountCents,
		Custom:      custom,
		Interval:    interval,
		Anchor:      p.IncurredOn,
		NextDue:     dateOf(nextDue(p.IncurredOn.Time, interval, p.IncurredOn.Time)),
	})
	if err != nil {
		return id, err
	}
	return id, insertRecurringShares(ctx, q, id, p.HouseholdID, custom, shares)
}

// insertRecurringShares stores who shares a schedule, their amounts only for a custom split
func insertRecurringShares(ctx context.Context, q *dbqueries.Queries, id, householdID pgtype.UUID, custom bool, shares []expenseShare) error {
	for _, sh := range shares {
		var cents *int64
		if custom {
			cents = &sh.cents
		}
		if err := q.InsertRecurringShare(ctx, dbqueries.InsertRecurringShareParams{
			RecurringID: id,
			HouseholdID: householdID,
			UserID:      sh.userID,
			AmountCents: cents,
		}); err != nil {
			return err
		}
	}
	return nil
}

// advanceRecurring moves the actor's schedule id on from due, the draft due then was added or
// skipped: [g.ErrorNotFound] unless the actor pays it, a [FieldError] for [FieldRecurring]
// when due is not its next due day any more, a second tab or a skip got there first
func advanceRecurring(ctx context.Context, q *dbqueries.Queries, actor Membership, id pgtype.UUID, due pgtype.Date) (string, error) {
	r, err := q.LockRecurring(ctx, dbqueries.LockRecurringParams{ID: id, HouseholdID: actor.HouseholdID})
	if err != nil {
		return "", noRowsAs(err, g.ErrorNotFound)
	}
	if r.PayerID != actor.UserID {
		return "", g.ErrorNotFound
	}
	if !due.Valid || !r.NextDue.Time.Equal(due.Time) {
		return "", FieldError{Field: FieldRecurring}
	}
	return r.Title, q.SetRecurringNextDue(ctx, dbqueries.SetRecurringNextDueParams{
		ID:          id,
		HouseholdID: actor.HouseholdID,
		NextDue:     dateOf(nextDue(r.Anchor.Time, r.Interval, r.NextDue.Time)),
	})
}

// lockRecurringToChange locks schedule id for the actor to edit or delete, [g.ErrorNotFound]
// unless they pay it or own the household
func lockRecurringToChange(ctx context.Context, q *dbqueries.Queries, actor Membership, id pgtype.UUID) (dbqueries.RecurringExpense, error) {
	r, err := q.LockRecurring(ctx, dbqueries.LockRecurringParams{ID: id, HouseholdID: actor.HouseholdID})
	if err != nil {
		return r, noRowsAs(err, g.ErrorNotFound)
	}
	if r.PayerID != actor.UserID && actor.Role != g.RoleOwner {
		return r, g.ErrorNotFound
	}
	return r, nil
}

// ListRecurring returns the household's schedules, soonest due first
func (s *DefaultService) ListRecurring(ctx context.Context, actor Membership) ([]RecurringRow, error) {
	return s.listRecurring(ctx, dbqueries.ListRecurringParams{HouseholdID: actor.HouseholdID, ID: pgtype.UUID{}, DueFor: pgtype.UUID{}})
}

// ListDueRecurring returns the schedules the actor pays that are due, oldest first
func (s *DefaultService) ListDueRecurring(ctx context.Context, actor Membership) ([]RecurringRow, error) {
	return s.listRecurring(ctx, dbqueries.ListRecurringParams{HouseholdID: actor.HouseholdID, ID: pgtype.UUID{}, DueFor: actor.UserID})
}

func (s *DefaultService) listRecurring(ctx context.Context, p dbqueries.ListRecurringParams) ([]RecurringRow, error) {
	rows, err := s.db.ListRecurring(ctx, p)
	return utils.Map(rows, func(r dbqueries.ListRecurringRow) RecurringRow { return RecurringRow{ListRecurringRow: r, Shares: nil} }), err
}

// GetRecurring returns a schedule of the actor's household with its shares, or [g.ErrorNotFound]
func (s *DefaultService) GetRecurring(ctx context.Context, actor Membership, id pgtype.UUID) (RecurringRow, error) {
	rows, err := s.listRecurring(ctx, dbqueries.ListRecurringParams{HouseholdID: actor.HouseholdID, ID: id, DueFor: pgtype.UUID{}})
	if err != nil {
		return RecurringRow{}, err
	}
	if len(rows) == 0 {
		return RecurringRow{}, g.ErrorNotFound
	}
	r := rows[0]
	r.Shares, err = s.db.ListRecurringShares(ctx, dbqueries.ListRecurringSharesParams{RecurringID: id, HouseholdID: actor.HouseholdID})
	return r, err
}

// UpdateRecurring saves a schedule, only its payer or an owner can; a new date starts it again
// from that day, see [lockRecurringToChange] and [ExpenseInput.clean] for errors
func (s *DefaultService) UpdateRecurring(ctx context.Context, actor Membership, id pgtype.UUID, in RecurringInput) error {
	if !slices.Contains(g.Intervals, in.Interval) {
		return g.ErrorInvalidInput
	}
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		// the members are read under the lock a leaver takes, see [DefaultService.Leave]
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		r, err := lockRecurringToChange(ctx, q, actor, id)
		if err != nil {
			return err
		}
		members, err := q.ListActiveMembers(ctx, actor.HouseholdID)
		if err != nil {
			return err
		}
		p, shares, err := in.clean(r.PayerID, members)
		if err != nil {
			return err
		}
		// a new interval counts from the next due day too, the old anchor is on the old step
		anchor, next := r.Anchor, r.NextDue
		if !p.IncurredOn.Time.Equal(r.NextDue.Time) || in.Interval != r.Interval {
			anchor, next = p.IncurredOn, p.IncurredOn
		}
		if err := q.UpdateRecurring(ctx, dbqueries.UpdateRecurringParams{
			ID:          id,
			HouseholdID: actor.HouseholdID,
			Category:    p.Category,
			Title:       p.Title,
			AmountCents: p.AmountCents,
			Custom:      in.Custom,
			Interval:    in.Interval,
			Anchor:      anchor,
			NextDue:     next,
		}); err != nil {
			return err
		}
		if err := q.DeleteRecurringShares(ctx, dbqueries.DeleteRecurringSharesParams{RecurringID: id, HouseholdID: actor.HouseholdID}); err != nil {
			return err
		}
		if err := insertRecurringShares(ctx, q, id, actor.HouseholdID, in.Custom, shares); err != nil {
			return err
		}
		// a schedule has no currency of its own, its drafts get the household's; read under the
		// lock, actor.Currency is from before it and a currency change can have committed since
		h, err := q.GetHousehold(ctx, actor.HouseholdID)
		if err != nil {
			return err
		}
		cents := p.AmountCents
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActRecurringEdited, SubjectID: id, TargetUserID: pgtype.UUID{}, Title: p.Title, AmountCents: &cents, Currency: h.Currency})
	})
}

// DeleteRecurring ends a schedule, the expenses it added stay, see [lockRecurringToChange]
func (s *DefaultService) DeleteRecurring(ctx context.Context, actor Membership, id pgtype.UUID) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		// the household is locked before the schedule like a leaver does, see [cleanChore]
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		r, err := lockRecurringToChange(ctx, q, actor, id)
		if err != nil {
			return err
		}
		if err := rowsAffected(q.DeleteRecurring(ctx, dbqueries.DeleteRecurringParams{ID: id, HouseholdID: actor.HouseholdID})); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActRecurringDeleted, SubjectID: id, TargetUserID: pgtype.UUID{}, Title: r.Title, AmountCents: nil, Currency: ""})
	})
}

// SkipRecurring moves the actor's schedule past the draft due on due, see [advanceRecurring]
func (s *DefaultService) SkipRecurring(ctx context.Context, actor Membership, id pgtype.UUID, due pgtype.Date) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		title, err := advanceRecurring(ctx, q, actor, id, due)
		if err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActRecurringSkipped, SubjectID: id, TargetUserID: pgtype.UUID{}, Title: title, AmountCents: nil, Currency: ""})
	})
}

// NotifyRecurringDue pushes each payer once about each due draft of theirs, run it every minute
func (s *DefaultService) NotifyRecurringDue(ctx context.Context) error {
	rows, err := s.db.ListRecurringToNotify(ctx)
	if err != nil {
		return err
	}
	for _, r := range rows {
		// claimed first, so two runs at once cannot both push
		n, err := s.db.SetRecurringNotified(ctx, dbqueries.SetRecurringNotifiedParams{ID: r.ID, NextDue: r.NextDue})
		if err != nil {
			return err
		}
		if n == 0 {
			continue
		}
		payer, err := s.db.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{HouseholdID: r.HouseholdID, UserID: r.PayerID})
		if err != nil {
			// left since the list was read, the trigger ended the schedule
			if err = noRowsAs(err, nil); err != nil {
				return err
			}
			continue
		}
		s.notifyAsync(ctx, payer, g.NotifyRecurring, []pgtype.UUID{r.PayerID}, func(ctx context.Context, _ pgtype.UUID, _, household string) push.Message {
			return push.Message{
				Title: utils.T(ctx, locales.LK_RecurringPushTitle, household),
				Body:  utils.T(ctx, locales.LK_RecurringPushBody, r.Title, utils.FormatDate(ctx, r.NextDue.Time)),
				URL:   routes.URL(routes.Expenses, r.HouseholdID),
				Tag:   "recurring-" + r.ID.String(),
			}
		})
	}
	return nil
}
