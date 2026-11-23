package services

import (
	"context"
	"crypto/rand"
	"errors"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/models"
	"roommates/internal/push"
	"roommates/internal/routes"
	"roommates/internal/utils"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Membership is the active membership of the signed in user in a household,
// it is the actor for every household action
type Membership = dbqueries.GetActiveMembershipRow

const inviteValidFor = 7 * 24 * time.Hour

// Crockford base32, no I L O U so typed codes are hard to misread
const inviteCodeAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
const inviteCodeLength = 10

func newInviteCode() string {
	b := make([]byte, inviteCodeLength)
	_, _ = rand.Read(b) // never returns an error
	for i := range b {
		// 256 is a multiple of 32, so there is no modulo bias
		b[i] = inviteCodeAlphabet[b[i]%32]
	}
	return string(b)
}

// NormalizeInviteCode makes a typed code match the stored one,
// a whole pasted invite link is accepted too
func NormalizeInviteCode(code string) string {
	code = strings.TrimSpace(code)
	if u, err := url.Parse(code); err == nil && u.Query().Has("code") {
		code = u.Query().Get("code")
	}
	return strings.ToUpper(strings.TrimSpace(code))
}

// inTx runs fn in a transaction, committing if fn returns nil
func (s *DefaultService) inTx(ctx context.Context, fn func(q *dbqueries.Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op after commit

	if err := fn(s.db.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func noRowsAs(err, as error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return as
	}
	return err
}

// -----------------------------------------------------------------------------

// CreateHousehold creates a household with userID as its owner, in timeZone (the creator's
// browser's) when that is one of [HouseholdSettings.TimeZones], else in the column's default
func (s *DefaultService) CreateHousehold(ctx context.Context, userID pgtype.UUID, name, timeZone string) (pgtype.UUID, error) {
	var householdID pgtype.UUID
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		var err error
		householdID, err = q.InsertHousehold(ctx, dbqueries.InsertHouseholdParams{
			Name:      strings.TrimSpace(name),
			CreatedBy: userID,
		})
		if err != nil {
			return err
		}
		if ok, err := q.IsTimeZone(ctx, timeZone); err != nil {
			return err
		} else if ok {
			if err := q.SetHouseholdTimeZone(ctx, dbqueries.SetHouseholdTimeZoneParams{ID: householdID, TimeZone: timeZone}); err != nil {
				return err
			}
		}
		return q.UpsertMembership(ctx, dbqueries.UpsertMembershipParams{
			HouseholdID: householdID,
			UserID:      userID,
			Role:        g.RoleOwner,
		})
	})
	return householdID, err
}

func (s *DefaultService) ListHouseholds(ctx context.Context, userID pgtype.UUID) ([]dbqueries.ListHouseholdsOfUserRow, error) {
	return s.db.ListHouseholdsOfUser(ctx, userID)
}

// GetMembership returns [g.ErrorNotFound] if the user is not an active member
func (s *DefaultService) GetMembership(ctx context.Context, householdID, userID pgtype.UUID) (Membership, error) {
	m, err := s.db.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{
		HouseholdID: householdID,
		UserID:      userID,
	})
	return m, noRowsAs(err, g.ErrorNotFound)
}

type HouseholdOverview struct {
	Household dbqueries.GetHouseholdRow
	Members   []dbqueries.ListActiveMembersRow
	// invites need a confirmed email, see [DefaultService.CreateInvite]
	ActorEmailValidated bool
	// only filled for owners, members cannot manage invites
	Invites       []dbqueries.ListActiveInvitesRow
	ShoppingLists []dbqueries.ListShoppingListsRow
}

func (s *DefaultService) GetHouseholdOverview(ctx context.Context, actor Membership) (HouseholdOverview, error) {
	var o HouseholdOverview
	var err error
	if o.Household, err = s.db.GetHousehold(ctx, actor.HouseholdID); err != nil {
		return o, err
	}
	if o.Members, err = s.db.ListActiveMembers(ctx, actor.HouseholdID); err != nil {
		return o, err
	}
	actorUser, err := s.db.GetUser(ctx, actor.UserID)
	if err != nil {
		return o, err
	}
	o.ActorEmailValidated = actorUser.EmailValidated
	if o.ShoppingLists, err = s.db.ListShoppingLists(ctx, actor.HouseholdID); err != nil {
		return o, err
	}
	if actor.Role == g.RoleOwner {
		o.Invites, err = s.db.ListActiveInvites(ctx, actor.HouseholdID)
	}
	return o, err
}

// -----------------------------------------------------------------------------

// CreateInvite needs an owner with a confirmed email ([g.ErrorEmailNotVerified]),
// so throwaway accounts cannot send invites
func (s *DefaultService) CreateInvite(ctx context.Context, actor Membership) error {
	_, err := s.createInvite(ctx, actor)
	return err
}

func (s *DefaultService) createInvite(ctx context.Context, actor Membership) (string, error) {
	if actor.Role != g.RoleOwner {
		return "", g.ErrorForbidden
	}
	user, err := s.db.GetUser(ctx, actor.UserID)
	if err != nil {
		return "", err
	}
	if !user.EmailValidated {
		return "", g.ErrorEmailNotVerified
	}
	code := newInviteCode()
	_, err = s.db.InsertInvite(ctx, dbqueries.InsertInviteParams{
		HouseholdID: actor.HouseholdID,
		Code:        code,
		CreatedBy:   actor.UserID,
		ExpiresAt:   timestamptz(time.Now().Add(inviteValidFor)),
	})
	return code, err
}

func (s *DefaultService) DeleteInvite(ctx context.Context, actor Membership, inviteID pgtype.UUID) error {
	if actor.Role != g.RoleOwner {
		return g.ErrorForbidden
	}
	return s.db.DeleteInvite(ctx, dbqueries.DeleteInviteParams{
		ID:          inviteID,
		HouseholdID: actor.HouseholdID,
	})
}

// GetInvite returns [g.ErrorInviteInvalid] for unknown or expired codes
func (s *DefaultService) GetInvite(ctx context.Context, code string) (dbqueries.GetInviteByCodeRow, error) {
	invite, err := s.db.GetInviteByCode(ctx, NormalizeInviteCode(code))
	return invite, noRowsAs(err, g.ErrorInviteInvalid)
}

// JoinWithInvite makes userID a member, an existing member keeps their role
func (s *DefaultService) JoinWithInvite(ctx context.Context, code string, userID pgtype.UUID) (pgtype.UUID, error) {
	code = NormalizeInviteCode(code)
	var householdID pgtype.UUID
	var joined Membership
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		invite, err := q.GetInviteByCode(ctx, code)
		if err != nil {
			return noRowsAs(err, g.ErrorInviteInvalid)
		}
		if err := q.LockHousehold(ctx, invite.HouseholdID); err != nil {
			return err
		}
		// read again under the lock, a removal or the last member leaving revokes the invites
		if invite, err = q.GetInviteByCode(ctx, code); err != nil {
			return noRowsAs(err, g.ErrorInviteInvalid)
		}
		householdID = invite.HouseholdID

		_, err = q.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{
			HouseholdID: householdID,
			UserID:      userID,
		})
		if err == nil {
			return nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		err = q.UpsertMembership(ctx, dbqueries.UpsertMembershipParams{
			HouseholdID: householdID,
			UserID:      userID,
			Role:        g.RoleMember,
		})
		if err != nil {
			return err
		}
		joined, err = q.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{
			HouseholdID: householdID,
			UserID:      userID,
		})
		if err != nil {
			return err
		}
		return logActivity(ctx, q, joined, ActivityEntry{Kind: g.ActMemberJoined, SubjectID: pgtype.UUID{}, TargetUserID: pgtype.UUID{}, Title: "", AmountCents: nil, Currency: ""})
	})
	if err != nil {
		return pgtype.UUID{}, err
	}
	// an existing member joined nobody new
	if joined.UserID.Valid {
		s.notifyHousehold(ctx, joined, nil, routes.HouseholdPeople, func(ctx context.Context, _ pgtype.UUID, actorName string) string {
			return utils.T(ctx, locales.LK_OwnershipPushJoined, actorName)
		})
	}
	return householdID, nil
}

// notifyHousehold sends a [g.NotifyHousehold] push about a change made by actor to userIDs,
// or to the other members when nil, the notification opens the household page at the route
func (s *DefaultService) notifyHousehold(ctx context.Context, actor Membership, userIDs []pgtype.UUID, route string,
	body func(ctx context.Context, userID pgtype.UUID, actorName string) string) {
	if userIDs == nil {
		var err error
		if userIDs, err = s.otherMemberIDs(ctx, actor); err != nil {
			// the change is made, only the notifications are missing
			log.Error().Err(err).Msg("could not list members for household push")
			return
		}
	}
	s.notifyAsync(ctx, actor, g.NotifyHousehold, userIDs, func(ctx context.Context, userID pgtype.UUID, actorName, household string) push.Message {
		return push.Message{
			Title: household,
			Body:  body(ctx, userID, actorName),
			URL:   routes.URL(route, actor.HouseholdID),
			Tag:   "household-" + actor.HouseholdID.String(),
		}
	})
}

// -----------------------------------------------------------------------------

// Leave ends the actor's membership, see [g.ErrorLastOwner] and [g.ErrorBalanceNotZero]
func (s *DefaultService) Leave(ctx context.Context, actor Membership) error {
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		return leaveTx(ctx, q, actor)
	})
	if err == nil {
		s.notifyLeft(ctx, actor)
	}
	return err
}

// leaveTx is [DefaultService.Leave] inside q's transaction, without the push
func leaveTx(ctx context.Context, q *dbqueries.Queries, actor Membership) error {
	if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
		return err
	}
	counts, err := q.CountActiveMembers(ctx, actor.HouseholdID)
	if err != nil {
		return err
	}
	if actor.Role == g.RoleOwner && counts.Owners == 1 && counts.Members > 1 {
		return g.ErrorLastOwner
	}
	if err := checkSettled(ctx, q, actor.HouseholdID, actor.UserID); err != nil {
		return err
	}
	// an invite into an empty household would make a member nobody can manage
	if counts.Members == 1 {
		if err := q.DeleteInvitesOfHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
	}
	if err := q.EndMembership(ctx, dbqueries.EndMembershipParams{
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	}); err != nil {
		return err
	}
	return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActMemberLeft, SubjectID: pgtype.UUID{}, TargetUserID: pgtype.UUID{}, Title: "", AmountCents: nil, Currency: ""})
}

// notifyLeft tells the other members that actor left
func (s *DefaultService) notifyLeft(ctx context.Context, actor Membership) {
	s.notifyHousehold(ctx, actor, nil, routes.HouseholdPeople, func(ctx context.Context, _ pgtype.UUID, actorName string) string {
		return utils.T(ctx, locales.LK_OwnershipPushLeft, actorName)
	})
}

// ownerActsOnMember checks that actor is an owner and target an active
// member (not an owner), returns [g.ErrorNotFound] or [g.ErrorForbidden]
func ownerActsOnMember(ctx context.Context, q *dbqueries.Queries, actor Membership, targetID pgtype.UUID) error {
	if actor.Role != g.RoleOwner {
		return g.ErrorForbidden
	}
	if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
		return err
	}
	// again under the lock, a transfer or step down of the actor's that ran at the same time
	// (a double submit) can have ended it since the request loaded it
	self, err := q.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{
		HouseholdID: actor.HouseholdID,
		UserID:      actor.UserID,
	})
	if err != nil {
		return noRowsAs(err, g.ErrorForbidden)
	}
	if self.Role != g.RoleOwner {
		return g.ErrorForbidden
	}
	target, err := q.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{
		HouseholdID: actor.HouseholdID,
		UserID:      targetID,
	})
	if err != nil {
		return noRowsAs(err, g.ErrorNotFound)
	}
	if target.Role == g.RoleOwner {
		return g.ErrorForbidden
	}
	return nil
}

// RemoveMember lets an owner remove a member, owners cannot be removed,
// all invites are revoked so the removed member cannot rejoin with one,
// see [g.ErrorBalanceNotZero]
func (s *DefaultService) RemoveMember(ctx context.Context, actor Membership, targetID pgtype.UUID) error {
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := ownerActsOnMember(ctx, q, actor, targetID); err != nil {
			return err
		}
		if err := checkSettled(ctx, q, actor.HouseholdID, targetID); err != nil {
			return err
		}
		if err := q.DeleteInvitesOfHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		if err := q.EndMembership(ctx, dbqueries.EndMembershipParams{
			HouseholdID: actor.HouseholdID,
			UserID:      targetID,
		}); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActMemberRemoved, SubjectID: pgtype.UUID{}, TargetUserID: targetID, Title: "", AmountCents: nil, Currency: ""})
	})
	if err != nil {
		return err
	}
	// the member is removed, a failure below only loses the notifications
	target, err := s.db.GetUser(ctx, targetID)
	if err != nil {
		log.Error().Err(err).Msg("could not get removed member for push")
		return nil
	}
	ids, err := s.otherMemberIDs(ctx, actor)
	if err != nil {
		log.Error().Err(err).Msg("could not list members for household push")
		return nil
	}
	s.notifyHousehold(ctx, actor, ids, routes.HouseholdPeople, func(ctx context.Context, _ pgtype.UUID, actorName string) string {
		return utils.T(ctx, locales.LK_OwnershipPushRemoved, actorName, target.Username)
	})
	// the removed member cannot open the household any more, so theirs links home
	s.notifyAsync(ctx, actor, g.NotifyHousehold, []pgtype.UUID{targetID}, func(ctx context.Context, _ pgtype.UUID, actorName, household string) push.Message {
		return push.Message{
			Title: household,
			Body:  utils.T(ctx, locales.LK_OwnershipPushRemovedYou, actorName),
			URL:   "/",
			Tag:   "household-" + actor.HouseholdID.String(),
		}
	})
	return nil
}

// otherMemberIDs are the active members of actor's household but actor
func (s *DefaultService) otherMemberIDs(ctx context.Context, actor Membership) ([]pgtype.UUID, error) {
	members, err := s.db.ListActiveMembers(ctx, actor.HouseholdID)
	ids := []pgtype.UUID{}
	for _, m := range members {
		if m.UserID != actor.UserID {
			ids = append(ids, m.UserID)
		}
	}
	return ids, err
}

// PromoteMember lets an owner make a member an owner
func (s *DefaultService) PromoteMember(ctx context.Context, actor Membership, targetID pgtype.UUID) error {
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := ownerActsOnMember(ctx, q, actor, targetID); err != nil {
			return err
		}
		if err := q.SetMembershipRole(ctx, dbqueries.SetMembershipRoleParams{
			HouseholdID: actor.HouseholdID,
			UserID:      targetID,
			Role:        g.RoleOwner,
		}); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActMemberPromoted, SubjectID: pgtype.UUID{}, TargetUserID: targetID, Title: "", AmountCents: nil, Currency: ""})
	})
	if err == nil {
		s.notifyMadeOwner(ctx, actor, targetID)
	}
	return err
}

func (s *DefaultService) notifyMadeOwner(ctx context.Context, actor Membership, targetID pgtype.UUID) {
	s.notifyHousehold(ctx, actor, []pgtype.UUID{targetID}, routes.HouseholdPeople, func(ctx context.Context, _ pgtype.UUID, actorName string) string {
		return utils.T(ctx, locales.LK_OwnershipPushOwner, actorName)
	})
}

// TransferOwnership makes a member an owner and the actor a member in one go,
// returns [g.ErrorForbidden] unless an owner hands over to an active member
func (s *DefaultService) TransferOwnership(ctx context.Context, actor Membership, targetID pgtype.UUID) error {
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := ownerActsOnMember(ctx, q, actor, targetID); err != nil {
			return err
		}
		err := q.SetMembershipRole(ctx, dbqueries.SetMembershipRoleParams{
			HouseholdID: actor.HouseholdID,
			UserID:      targetID,
			Role:        g.RoleOwner,
		})
		if err != nil {
			return err
		}
		if err := q.SetMembershipRole(ctx, dbqueries.SetMembershipRoleParams{
			HouseholdID: actor.HouseholdID,
			UserID:      actor.UserID,
			Role:        g.RoleMember,
		}); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActOwnershipTransferred, SubjectID: pgtype.UUID{}, TargetUserID: targetID, Title: "", AmountCents: nil, Currency: ""})
	})
	if err == nil {
		s.notifyMadeOwner(ctx, actor, targetID)
	}
	return err
}

// StepDown makes an owner a member while another owner is left, see [g.ErrorLastOwner]
func (s *DefaultService) StepDown(ctx context.Context, actor Membership) error {
	if actor.Role != g.RoleOwner {
		return g.ErrorForbidden
	}
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		counts, err := q.CountActiveMembers(ctx, actor.HouseholdID)
		if err != nil {
			return err
		}
		if counts.Owners < 2 {
			return g.ErrorLastOwner
		}
		if err := q.SetMembershipRole(ctx, dbqueries.SetMembershipRoleParams{
			HouseholdID: actor.HouseholdID,
			UserID:      actor.UserID,
			Role:        g.RoleMember,
		}); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActSteppedDown, SubjectID: pgtype.UUID{}, TargetUserID: pgtype.UUID{}, Title: "", AmountCents: nil, Currency: ""})
	})
}

// PurgeDate is when a household deleted at deletedAt is purged, until then it can be restored
func PurgeDate(deletedAt pgtype.Timestamptz) time.Time {
	return deletedAt.Time.Add(g.HouseholdPurgeAfter)
}

// DeleteHousehold lets an owner delete the household after typing its name, a mismatch is
// a [FieldError] of [FieldName]; open debts do not stop it, the page warns about them,
// the household stays read-only and restorable until [PurgeDate], its invites are gone at once
func (s *DefaultService) DeleteHousehold(ctx context.Context, actor Membership, confirmName string) error {
	if actor.Role != g.RoleOwner {
		return g.ErrorForbidden
	}
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		h, err := q.GetHousehold(ctx, actor.HouseholdID)
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		if strings.TrimSpace(confirmName) != h.Name {
			return FieldError{Field: FieldName}
		}
		if err := q.DeleteInvitesOfHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		if err := q.SoftDeleteHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActHouseholdDeleted, SubjectID: pgtype.UUID{}, TargetUserID: pgtype.UUID{}, Title: h.Name, AmountCents: nil, Currency: ""})
	})
	if err == nil {
		purge := time.Now().Add(g.HouseholdPurgeAfter)
		s.notifyHousehold(ctx, actor, nil, routes.Household, func(ctx context.Context, _ pgtype.UUID, actorName string) string {
			return utils.T(ctx, locales.LK_OwnershipPushDeleted, actorName, utils.FormatDate(ctx, purge))
		})
	}
	return err
}

// RestoreHousehold lets an active owner undo [DefaultService.DeleteHousehold] before it is purged,
// returns [g.ErrorNotFound] for non-members and [g.ErrorForbidden] for members
func (s *DefaultService) RestoreHousehold(ctx context.Context, userID, householdID pgtype.UUID) error {
	var actor Membership
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := q.LockHousehold(ctx, householdID); err != nil {
			return err
		}
		var err error
		actor, err = q.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{
			HouseholdID: householdID,
			UserID:      userID,
		})
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		if actor.Role != g.RoleOwner {
			return g.ErrorForbidden
		}
		if !actor.DeletedAt.Valid {
			return nil
		}
		if err := q.RestoreHousehold(ctx, householdID); err != nil {
			return err
		}
		return logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActHouseholdRestored, SubjectID: pgtype.UUID{}, TargetUserID: pgtype.UUID{}, Title: "", AmountCents: nil, Currency: ""})
	})
	if err == nil && actor.DeletedAt.Valid {
		s.notifyHousehold(ctx, actor, nil, routes.Household, func(ctx context.Context, _ pgtype.UUID, actorName string) string {
			return utils.T(ctx, locales.LK_OwnershipPushRestored, actorName)
		})
	}
	return err
}

// PurgeDeletedHouseholds removes the households deleted longer than [g.HouseholdPurgeAfter] ago
// with everything in them, their files are queued for [DefaultService.PurgeFiles]
func (s *DefaultService) PurgeDeletedHouseholds(ctx context.Context) error {
	n, err := s.db.PurgeDeletedHouseholds(ctx, timestamptz(time.Now().Add(-g.HouseholdPurgeAfter)))
	if n > 0 {
		log.Info().Int64("households", n).Msg("purged deleted households")
	}
	return err
}

// -----------------------------------------------------------------------------

// what a [FieldError] of [DefaultService.UpdateHousehold] names
const (
	FieldName         = "name"
	FieldCurrency     = "currency"
	FieldTimeZone     = "time_zone"
	FieldReminderHour = "reminder_hour"
)

// HouseholdCurrencies are the currencies a household can use
var HouseholdCurrencies = []string{"EUR", "USD", "GBP", "SEK", "NOK", "DKK", "PLN", "CHF"}

type HouseholdSettings struct {
	Name     string
	Currency string
	// an IANA name, one of TimeZones
	TimeZone  string
	TimeZones []string
	// local hour 0..23 of the chore and recurring expense reminders
	ReminderHour int
	// Settled is true while the currency can change, nobody owes anybody
	Settled bool
	// who owes whom, deleting the household warns about them
	Debts []Debt
}

func (s *DefaultService) GetHouseholdSettings(ctx context.Context, actor Membership) (HouseholdSettings, error) {
	h, err := s.db.GetHousehold(ctx, actor.HouseholdID)
	if err != nil {
		return HouseholdSettings{}, noRowsAs(err, g.ErrorNotFound)
	}
	settled, err := householdSettled(ctx, s.db, actor.HouseholdID)
	if err != nil {
		return HouseholdSettings{}, err
	}
	debts, err := s.Debts(ctx, actor)
	if err != nil {
		return HouseholdSettings{}, err
	}
	zones, err := s.db.ListTimeZones(ctx)
	return HouseholdSettings{Name: h.Name, Currency: h.Currency, TimeZone: h.TimeZone, TimeZones: zones, ReminderHour: int(h.ReminderHour), Settled: settled, Debts: debts}, err
}

// ValidReminderHour is whether h is an hour of the day
func ValidReminderHour(h int) bool { return h >= 0 && h <= 23 }

// UpdateHousehold lets an owner rename the household and change its currency, time zone and
// reminder hour, returns [g.ErrorForbidden] for others, a [FieldError] for a bad field
// and for a currency change while a debt is open, amounts would change their meaning
func (s *DefaultService) UpdateHousehold(ctx context.Context, actor Membership, name, currency, timeZone string, reminderHour int) error {
	if actor.Role != g.RoleOwner {
		return g.ErrorForbidden
	}
	name, ok := trimToLength(name, 1, models.MaxHouseholdNameLength)
	if !ok {
		return FieldError{Field: FieldName}
	}
	if !slices.Contains(HouseholdCurrencies, currency) {
		return FieldError{Field: FieldCurrency}
	}
	if !ValidReminderHour(reminderHour) {
		return FieldError{Field: FieldReminderHour}
	}
	if ok, err := s.db.IsTimeZone(ctx, timeZone); err != nil {
		return err
	} else if !ok {
		return FieldError{Field: FieldTimeZone}
	}
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := q.LockHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		h, err := q.GetHousehold(ctx, actor.HouseholdID)
		if err != nil {
			return noRowsAs(err, g.ErrorNotFound)
		}
		if h.Name != name {
			if err := logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActHouseholdRenamed, SubjectID: pgtype.UUID{}, TargetUserID: pgtype.UUID{}, Title: name, AmountCents: nil, Currency: ""}); err != nil {
				return err
			}
		}
		if h.Currency != currency {
			if err := logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActCurrencyChanged, SubjectID: pgtype.UUID{}, TargetUserID: pgtype.UUID{}, Title: currency, AmountCents: nil, Currency: ""}); err != nil {
				return err
			}
			settled, err := householdSettled(ctx, q, actor.HouseholdID)
			if err != nil {
				return err
			}
			if !settled {
				return FieldError{Field: FieldCurrency}
			}
			if err := closeEvenDebts(ctx, q, actor.HouseholdID); err != nil {
				return err
			}
		}
		return q.UpdateHousehold(ctx, dbqueries.UpdateHouseholdParams{
			Name:         name,
			Currency:     currency,
			TimeZone:     timeZone,
			ReminderHour: int16(reminderHour),
			ID:           actor.HouseholdID,
		})
	})
}
