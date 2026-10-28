package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	g "roommates/internal/config"
	"roommates/internal/db/dbqueries"
	"roommates/internal/locales"
	"roommates/internal/mailer"
	"roommates/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// email_tokens.kind values, the table CHECKs these
const (
	tokenKindVerify = "verify"
	tokenKindReset  = "reset"
)

const (
	verifyTokenValidFor = 48 * time.Hour
	resetTokenValidFor  = time.Hour
)

func hashToken(token string) []byte {
	h := sha256.Sum256([]byte(token))
	return h[:]
}

func timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, InfinityModifier: pgtype.Finite, Valid: true}
}

// newEmailToken stores the hash of a new token and returns the token
func newEmailToken(ctx context.Context, q *dbqueries.Queries, userID pgtype.UUID, kind string, validFor time.Duration) (string, error) {
	token := rand.Text()
	return token, q.InsertEmailToken(ctx, dbqueries.InsertEmailTokenParams{
		TokenHash: hashToken(token),
		UserID:    userID,
		Kind:      kind,
		ExpiresAt: timestamptz(time.Now().Add(validFor)),
	})
}

// sendAsync so the response does not wait for SMTP, and its timing does
// not reveal whether an email was sent, failures are only logged
func (s *DefaultService) sendAsync(msg mailer.Message) {
	go func() {
		if err := s.mailer.Send(msg); err != nil {
			log.Error().Err(err).Str("subject", msg.Subject).Msg("could not send email")
		}
	}()
}

// -----------------------------------------------------------------------------

func (s *DefaultService) GetUser(ctx context.Context, userID pgtype.UUID) (dbqueries.GetUserRow, error) {
	return s.db.GetUser(ctx, userID)
}

// SendVerificationEmail sends a link that confirms the email, nothing if already confirmed
//
// ctx gives the language of the email
func (s *DefaultService) SendVerificationEmail(ctx context.Context, userID pgtype.UUID) error {
	user, err := s.db.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	if user.EmailValidated {
		return nil
	}
	token, err := newEmailToken(ctx, s.db, userID, tokenKindVerify, verifyTokenValidFor)
	if err != nil {
		return err
	}
	link := g.GetAppBaseURL() + g.RouteVerifyEmail + "?token=" + token
	s.sendAsync(mailer.Message{
		To:      user.Email,
		Subject: utils.T(ctx, locales.LK_MailVerifySubject),
		Body:    utils.T(ctx, locales.LK_MailVerifyBody, user.Username, link),
	})
	return nil
}

// VerifyEmail returns [g.ErrorTokenInvalid] for unknown, expired or used tokens
func (s *DefaultService) VerifyEmail(ctx context.Context, token string) error {
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		userID, err := q.ConsumeEmailToken(ctx, dbqueries.ConsumeEmailTokenParams{
			TokenHash: hashToken(token),
			Kind:      tokenKindVerify,
		})
		if err != nil {
			return noRowsAs(err, g.ErrorTokenInvalid)
		}
		if err := q.SetEmailValidated(ctx, userID); err != nil {
			return err
		}
		return q.DeleteEmailTokensOfUser(ctx, dbqueries.DeleteEmailTokensOfUserParams{
			UserID: userID,
			Kind:   tokenKindVerify,
		})
	})
}

// RequestPasswordReset sends a reset link if an account has this email,
// the caller must respond the same either way so accounts cannot be probed
func (s *DefaultService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.db.GetUserCredentials(ctx, NormalizeEmail(email))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	token, err := newEmailToken(ctx, s.db, user.ID, tokenKindReset, resetTokenValidFor)
	if err != nil {
		return err
	}
	link := g.GetAppBaseURL() + g.RouteResetPassword + "?token=" + token
	s.sendAsync(mailer.Message{
		To:      user.Email,
		Subject: utils.T(ctx, locales.LK_MailResetSubject),
		Body:    utils.T(ctx, locales.LK_MailResetBody, user.Username, link),
	})
	return nil
}

// CheckPasswordResetToken returns [g.ErrorTokenInvalid] if the token cannot be used
func (s *DefaultService) CheckPasswordResetToken(ctx context.Context, token string) error {
	_, err := s.db.GetEmailToken(ctx, dbqueries.GetEmailTokenParams{
		TokenHash: hashToken(token),
		Kind:      tokenKindReset,
	})
	return noRowsAs(err, g.ErrorTokenInvalid)
}

// ResetPassword sets the password and confirms the email (the link proved access to it)
//
// shortcut: other signed in sessions stay valid, the session store cannot be
// queried by user, add a per user session version once sessions are checked against the db
func (s *DefaultService) ResetPassword(ctx context.Context, token, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	return s.inTx(ctx, func(q *dbqueries.Queries) error {
		userID, err := q.ConsumeEmailToken(ctx, dbqueries.ConsumeEmailTokenParams{
			TokenHash: hashToken(token),
			Kind:      tokenKindReset,
		})
		if err != nil {
			return noRowsAs(err, g.ErrorTokenInvalid)
		}
		if err := q.UpdateUserPassword(ctx, dbqueries.UpdateUserPasswordParams{
			ID:       userID,
			Password: string(hashed),
		}); err != nil {
			return err
		}
		if err := q.SetEmailValidated(ctx, userID); err != nil {
			return err
		}
		return q.DeleteEmailTokensOfUser(ctx, dbqueries.DeleteEmailTokensOfUserParams{
			UserID: userID,
			Kind:   tokenKindReset,
		})
	})
}

// -----------------------------------------------------------------------------

// SendInviteEmail creates an invite and sends it to email, see [DefaultService.CreateInvite]
func (s *DefaultService) SendInviteEmail(ctx context.Context, actor Membership, email string) error {
	code, err := s.createInvite(ctx, actor)
	if err != nil {
		return err
	}
	inviter, err := s.db.GetUser(ctx, actor.UserID)
	if err != nil {
		return err
	}
	household, err := s.db.GetHousehold(ctx, actor.HouseholdID)
	if err != nil {
		return err
	}
	link := g.GetAppBaseURL() + g.RouteInvite + "?code=" + code
	s.sendAsync(mailer.Message{
		To:      NormalizeEmail(email),
		Subject: utils.T(ctx, locales.LK_MailInviteSubject, inviter.Username, household.Name),
		Body:    utils.T(ctx, locales.LK_MailInviteBody, inviter.Username, household.Name, link, code),
	})
	return nil
}
