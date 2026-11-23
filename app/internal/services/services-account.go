package services

import (
	"context"
	"errors"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/mailer"
	"roommates/internal/routes"
	"roommates/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// ErrorEmailUnchanged is returned when the new email is the current one
var ErrorEmailUnchanged = errors.New("email unchanged")

// checkPassword returns [g.ErrorInvalidCredential] if password is not the user's
func (s *DefaultService) checkPassword(ctx context.Context, userID pgtype.UUID, password string) error {
	hashed, err := s.db.GetUserPassword(ctx, userID)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) != nil {
		return g.ErrorInvalidCredential
	}
	return nil
}

// ChangePassword sets a new password if current is right, cancels a pending email change
// and signs the user out everywhere, the caller signs the current session in again
func (s *DefaultService) ChangePassword(ctx context.Context, userID pgtype.UUID, current, password string) error {
	if err := s.checkPassword(ctx, userID, current); err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		if err := q.UpdateUserPassword(ctx, dbqueries.UpdateUserPasswordParams{
			ID:       userID,
			Password: string(hashed),
		}); err != nil {
			return err
		}
		// asked for with the old password, which may be the stolen one
		if err := q.DeleteEmailTokensOfUser(ctx, dbqueries.DeleteEmailTokensOfUserParams{
			UserID: userID,
			Kind:   tokenKindEmailChange,
		}); err != nil {
			return err
		}
		return signOutEverywhere(ctx, q, userID)
	})
}

// SignOutEverywhere ends every session of the user, this one too
func (s *DefaultService) SignOutEverywhere(ctx context.Context, userID pgtype.UUID) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		return signOutEverywhere(ctx, q, userID)
	})
}

// signOutEverywhere deletes the user's sessions and push subscriptions, so a signed out
// device stops getting notifications too; one still signed in sends its own again (push.js)
func signOutEverywhere(ctx context.Context, q *dbqueries.Queries, userID pgtype.UUID) error {
	if err := q.DeleteUserSessions(ctx, userID); err != nil {
		return err
	}
	return q.DeletePushSubscriptionsOfUser(ctx, userID)
}

// AccountDeleteBlocked names the households that stop [DefaultService.DeleteAccount]
type AccountDeleteBlocked struct {
	// shared households the user is the only owner of
	LastOwner []string
	// shared households where the user owes or is owed, or a settlement waits
	Unsettled []string
}

func (e AccountDeleteBlocked) Error() string {
	return "account delete blocked"
}

// DeleteAccount leaves the user's shared households, deletes the ones only they live in
// (restorable, the membership stays), signs them out everywhere and marks the account to be
// erased after [g.AccountEraseAfter], signing in before that restores it;
// returns [g.ErrorInvalidCredential] or [AccountDeleteBlocked], then nothing changed
//
// ctx gives the language of the email
func (s *DefaultService) DeleteAccount(ctx context.Context, userID pgtype.UUID, password string) error {
	if err := s.checkPassword(ctx, userID, password); err != nil {
		return err
	}
	user, err := s.db.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	var left []Membership
	err = s.inTx(ctx, func(q *dbqueries.Queries) error {
		households, err := q.ListHouseholdsOfUser(ctx, userID)
		if err != nil {
			return err
		}
		var blocked AccountDeleteBlocked
		for _, h := range households {
			if err := q.LockHousehold(ctx, h.ID); err != nil {
				return err
			}
			actor, err := q.GetActiveMembership(ctx, dbqueries.GetActiveMembershipParams{HouseholdID: h.ID, UserID: userID})
			if errors.Is(err, pgx.ErrNoRows) {
				continue // removed meanwhile
			} else if err != nil {
				return err
			}
			if actor.DeletedAt.Valid {
				continue
			}
			counts, err := q.CountActiveMembers(ctx, h.ID)
			if err != nil {
				return err
			}
			if counts.Members == 1 {
				// as DeleteHousehold, the membership stays so the household can still be restored
				if err := q.DeleteInvitesOfHousehold(ctx, h.ID); err != nil {
					return err
				}
				if err := q.SoftDeleteHousehold(ctx, h.ID); err != nil {
					return err
				}
				if err := logActivity(ctx, q, actor, ActivityEntry{Kind: g.ActHouseholdDeleted, SubjectID: pgtype.UUID{}, TargetUserID: pgtype.UUID{}, Title: h.Name, AmountCents: nil, Currency: ""}); err != nil {
					return err
				}
				continue
			}
			// keeps going after a blocked one so the page can name them all, the rollback undoes the rest
			switch err := leaveTx(ctx, q, actor); {
			case errors.Is(err, g.ErrorLastOwner):
				blocked.LastOwner = append(blocked.LastOwner, h.Name)
				// leaveTx stops before the balances, both are named so one try shows everything
				switch err := checkSettled(ctx, q, h.ID, userID); {
				case errors.Is(err, g.ErrorBalanceNotZero):
					blocked.Unsettled = append(blocked.Unsettled, h.Name)
				case err != nil:
					return err
				}
			case errors.Is(err, g.ErrorBalanceNotZero):
				blocked.Unsettled = append(blocked.Unsettled, h.Name)
			case err != nil:
				return err
			default:
				left = append(left, actor)
			}
		}
		if len(blocked.LastOwner) > 0 || len(blocked.Unsettled) > 0 {
			return blocked
		}
		if err := q.MarkAccountDeleted(ctx, userID); err != nil {
			return err
		}
		return signOutEverywhere(ctx, q, userID)
	})
	if err != nil {
		return err
	}
	for _, actor := range left {
		s.notifyLeft(ctx, actor)
	}
	s.sendAsync(mailer.Message{
		To:      user.Email,
		Subject: utils.T(ctx, locales.LK_MailAccountDeletedSubject),
		Body:    utils.T(ctx, locales.LK_MailAccountDeletedBody, user.Username, utils.FormatDate(ctx, time.Now().Add(g.AccountEraseAfter))),
	})
	return nil
}

// RestoreAccount undoes [DefaultService.DeleteAccount], signing in calls it,
// left households need a new invite, deleted ones are restored from home;
// returns [g.ErrorInvalidCredential] when the account was erased meanwhile
func (s *DefaultService) RestoreAccount(ctx context.Context, userID pgtype.UUID) error {
	n, err := s.db.RestoreAccount(ctx, userID)
	if err == nil && n == 0 {
		return g.ErrorInvalidCredential
	}
	return err
}

// EraseDeletedAccounts erases the accounts deleted longer than [g.AccountEraseAfter] ago,
// one transaction each, a failed one is logged and tried again next run
func (s *DefaultService) EraseDeletedAccounts(ctx context.Context) error {
	before := timestamptz(time.Now().Add(-g.AccountEraseAfter))
	ids, err := s.db.ListAccountsToErase(ctx, before)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.eraseAccount(ctx, id, before); err != nil {
			log.Error().Err(err).Str("user", id.String()).Msg("could not erase deleted account")
		}
	}
	return nil
}

func (s *DefaultService) eraseAccount(ctx context.Context, userID pgtype.UUID, before pgtype.Timestamptz) error {
	var user dbqueries.LockAccountToEraseRow
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		var err error
		if user, err = q.LockAccountToErase(ctx, dbqueries.LockAccountToEraseParams{ID: userID, DeletedBefore: before}); err != nil {
			return err
		}
		steps := []func(context.Context, pgtype.UUID) error{
			q.EndMembershipsOfUser,
			q.DeletePrivateNotesOfUser,
			q.DeleteAllEmailTokensOfUser,
			q.DeleteDataExportsOfUser,
			q.DeleteUserSessions,
			q.DeletePushSubscriptionsOfUser,
			q.EraseUser,
		}
		for _, step := range steps {
			if err := step(ctx, userID); err != nil {
				return err
			}
		}
		return nil
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil // signed in meanwhile
	} else if err != nil {
		return err
	}
	lang := ""
	if user.Language != nil {
		lang = *user.Language
	}
	lctx := locales.WithLocale(ctx, lang)
	s.sendAsync(mailer.Message{
		To:      user.Email,
		Subject: utils.T(lctx, locales.LK_MailAccountErasedSubject),
		Body:    utils.T(lctx, locales.LK_MailAccountErasedBody, user.Username),
	})
	return nil
}

// SetUsername saves a validated username, usernames are not unique
func (s *DefaultService) SetUsername(ctx context.Context, userID pgtype.UUID, username string) error {
	return s.db.UpdateUsername(ctx, dbqueries.UpdateUsernameParams{ID: userID, Username: username})
}

// RequestEmailChange sends a link to newEmail that makes it the user's email,
// returns [g.ErrorInvalidCredential], [ErrorEmailUnchanged] or [g.ErrorAccountAlreadyExists]
//
// ctx gives the language of the email
func (s *DefaultService) RequestEmailChange(ctx context.Context, userID pgtype.UUID, password, newEmail string) error {
	if err := s.checkPassword(ctx, userID, password); err != nil {
		return err
	}
	email := NormalizeEmail(newEmail)
	user, err := s.db.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	if email == user.Email {
		return ErrorEmailUnchanged
	}
	// tells no more than register does
	if _, err := s.db.GetUserCredentials(ctx, email); err == nil {
		return g.ErrorAccountAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	token, err := newEmailToken(ctx, s.db, userID, tokenKindEmailChange, emailChangeTokenValidFor, &email)
	if err != nil {
		return err
	}
	link := g.GetAppBaseURL() + routes.ConfirmEmail + "?token=" + token
	s.sendAsync(mailer.Message{
		To:      email,
		Subject: utils.T(ctx, locales.LK_MailEmailChangeSubject),
		Body:    utils.T(ctx, locales.LK_MailEmailChangeBody, user.Username, link),
	})
	return nil
}

// ConfirmEmailChange makes the new email of the token the user's and tells the old address,
// returns [g.ErrorTokenInvalid] or [g.ErrorAccountAlreadyExists] (taken since the link was sent)
func (s *DefaultService) ConfirmEmailChange(ctx context.Context, token string) error {
	var old dbqueries.GetUserRow
	var newEmail string
	err := s.inTx(ctx, func(q *dbqueries.Queries) error {
		row, err := q.ConsumeEmailChangeToken(ctx, hashToken(token))
		if err != nil {
			return noRowsAs(err, g.ErrorTokenInvalid)
		}
		if old, err = q.GetUser(ctx, row.UserID); err != nil {
			return err
		}
		newEmail = row.NewEmail
		// the unique email is the check, a failed one rolls back and the link stays usable
		if err := q.UpdateUserEmail(ctx, dbqueries.UpdateUserEmailParams{ID: row.UserID, Email: newEmail}); err != nil {
			if isUniqueViolation(err) {
				return g.ErrorAccountAlreadyExists
			}
			return err
		}
		// links sent to the old address must not work any more
		for _, kind := range []string{tokenKindEmailChange, tokenKindReset, tokenKindVerify} {
			if err := q.DeleteEmailTokensOfUser(ctx, dbqueries.DeleteEmailTokensOfUserParams{
				UserID: row.UserID,
				Kind:   kind,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if old.Email != newEmail {
		s.sendAsync(mailer.Message{
			To:      old.Email,
			Subject: utils.T(ctx, locales.LK_MailEmailChangedSubject),
			Body:    utils.T(ctx, locales.LK_MailEmailChangedBody, old.Username, newEmail),
		})
	}
	return nil
}
