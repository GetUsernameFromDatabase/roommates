package services

import (
	"context"
	"crypto/rand"
	"errors"
	"net/url"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
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

// CreateHousehold creates a household with userID as its owner
func (s *DefaultService) CreateHousehold(ctx context.Context, userID pgtype.UUID, name string) (pgtype.UUID, error) {
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
	invite, err := s.GetInvite(ctx, code)
	if err != nil {
		return pgtype.UUID{}, err
	}

	_, err = s.GetMembership(ctx, invite.HouseholdID, userID)
	if err == nil {
		return invite.HouseholdID, nil
	}
	if !errors.Is(err, g.ErrorNotFound) {
		return pgtype.UUID{}, err
	}

	return invite.HouseholdID, s.db.UpsertMembership(ctx, dbqueries.UpsertMembershipParams{
		HouseholdID: invite.HouseholdID,
		UserID:      userID,
		Role:        g.RoleMember,
	})
}

// -----------------------------------------------------------------------------

// Leave ends the actor's membership, see [g.ErrorLastOwner] and [g.ErrorBalanceNotZero]
func (s *DefaultService) Leave(ctx context.Context, actor Membership) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
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
		return q.EndMembership(ctx, dbqueries.EndMembershipParams{
			HouseholdID: actor.HouseholdID,
			UserID:      actor.UserID,
		})
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
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := ownerActsOnMember(ctx, q, actor, targetID); err != nil {
			return err
		}
		if err := checkSettled(ctx, q, actor.HouseholdID, targetID); err != nil {
			return err
		}
		if err := q.DeleteInvitesOfHousehold(ctx, actor.HouseholdID); err != nil {
			return err
		}
		return q.EndMembership(ctx, dbqueries.EndMembershipParams{
			HouseholdID: actor.HouseholdID,
			UserID:      targetID,
		})
	})
}

// PromoteMember lets an owner make a member an owner
func (s *DefaultService) PromoteMember(ctx context.Context, actor Membership, targetID pgtype.UUID) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := ownerActsOnMember(ctx, q, actor, targetID); err != nil {
			return err
		}
		return q.SetMembershipRole(ctx, dbqueries.SetMembershipRoleParams{
			HouseholdID: actor.HouseholdID,
			UserID:      targetID,
			Role:        g.RoleOwner,
		})
	})
}
