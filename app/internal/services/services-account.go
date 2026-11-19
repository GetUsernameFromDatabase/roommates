package services

import (
	"context"
	"errors"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/mailer"
	"roommates/internal/utils"

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

// ChangePassword sets a new password if current is right, and signs the user
// out everywhere, the caller signs the current session in again
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
		return q.DeleteUserSessions(ctx, userID)
	})
}

// SignOutEverywhere deletes every session of the user, this one too
func (s *DefaultService) SignOutEverywhere(ctx context.Context, userID pgtype.UUID) error {
	return s.db.DeleteUserSessions(ctx, userID)
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
	link := g.GetAppBaseURL() + g.RouteConfirmEmail + "?token=" + token
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
