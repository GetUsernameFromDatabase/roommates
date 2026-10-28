package services

import (
	"context"
	"regexp"
	g "roommates/internal/config"
	"roommates/internal/locales"
	"strings"
	"testing"

	"github.com/invopop/ctxi18n"
	"golang.org/x/crypto/bcrypt"
)

var tokenInLink = regexp.MustCompile(`token=([A-Z2-7]+)`)

// tokenOf returns the token of the link in an email body
func tokenOf(t *testing.T, body string) string {
	t.Helper()
	m := tokenInLink.FindStringSubmatch(body)
	if m == nil {
		t.Fatalf("no token link in %q", body)
	}
	return m[1]
}

// englishContext gives emails a language, like the i18n middleware does
func englishContext(t *testing.T) context.Context {
	t.Helper()
	locales.LoadLocales()
	ctx, err := ctxi18n.WithLocale(t.Context(), "en")
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func TestVerifyEmail(t *testing.T) {
	s := newTestService(t)
	ctx := englishContext(t)
	userID := newTestUser(t, s)

	if err := s.SendVerificationEmail(ctx, userID); err != nil {
		t.Fatal(err)
	}
	msg := nextEmail(t, s)
	if !strings.Contains(msg.Body, g.RouteVerifyEmail+"?token=") {
		t.Fatalf("no verify link in %q", msg.Body)
	}
	token := tokenOf(t, msg.Body)

	wantErr(t, "wrong token", s.VerifyEmail(ctx, "WRONG"), g.ErrorTokenInvalid)
	if err := s.VerifyEmail(ctx, token); err != nil {
		t.Fatal(err)
	}
	user, err := s.GetUser(ctx, userID)
	if err != nil || !user.EmailValidated {
		t.Fatalf("email not validated: %v", err)
	}
	wantErr(t, "token reused", s.VerifyEmail(ctx, token), g.ErrorTokenInvalid)

	// nothing is sent once confirmed
	if err := s.SendVerificationEmail(ctx, userID); err != nil {
		t.Fatal(err)
	}
	noEmail(t, s)
}

func TestResetPassword(t *testing.T) {
	s := newTestService(t)
	ctx := englishContext(t)
	userID := newTestUser(t, s)
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}

	// unknown accounts get no email and no error, so they cannot be told apart
	if err := s.RequestPasswordReset(ctx, "nobody-"+user.Email); err != nil {
		t.Fatal(err)
	}
	noEmail(t, s)

	// the typed email is normalized like on sign in
	if err := s.RequestPasswordReset(ctx, " "+strings.ToUpper(user.Email)); err != nil {
		t.Fatal(err)
	}
	msg := nextEmail(t, s)
	if msg.To != user.Email {
		t.Fatalf("reset sent to %q", msg.To)
	}
	token := tokenOf(t, msg.Body)

	if err := s.CheckPasswordResetToken(ctx, token); err != nil {
		t.Fatal(err)
	}
	// a verify token is not a reset token
	if err := s.SendVerificationEmail(ctx, userID); err != nil {
		t.Fatal(err)
	}
	verifyToken := tokenOf(t, nextEmail(t, s).Body)
	wantErr(t, "verify token used for reset", s.ResetPassword(ctx, verifyToken, "New-pass1!"), g.ErrorTokenInvalid)

	if err := s.ResetPassword(ctx, token, "New-pass1!"); err != nil {
		t.Fatal(err)
	}
	creds, err := s.db.GetUserCredentials(ctx, user.Email)
	if err != nil {
		t.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte("New-pass1!")) != nil {
		t.Fatal("password was not changed")
	}
	if user, _ := s.GetUser(ctx, userID); !user.EmailValidated {
		t.Fatal("a reset proves the email, it should be confirmed")
	}
	wantErr(t, "token reused", s.ResetPassword(ctx, token, "Other-pass1!"), g.ErrorTokenInvalid)
	wantErr(t, "check used token", s.CheckPasswordResetToken(ctx, token), g.ErrorTokenInvalid)
}

func TestExpiredTokenIsInvalid(t *testing.T) {
	s := newTestService(t)
	ctx := englishContext(t)
	userID := newTestUser(t, s)

	if err := s.SendVerificationEmail(ctx, userID); err != nil {
		t.Fatal(err)
	}
	token := tokenOf(t, nextEmail(t, s).Body)
	if _, err := s.pool.Exec(ctx,
		`UPDATE email_tokens SET expires_at = now() - interval '1 second' WHERE token_hash = $1`,
		hashToken(token)); err != nil {
		t.Fatal(err)
	}
	wantErr(t, "expired token", s.VerifyEmail(ctx, token), g.ErrorTokenInvalid)
}

func TestInviteNeedsConfirmedEmail(t *testing.T) {
	s := newTestService(t)
	ctx := englishContext(t)
	owner := newTestUser(t, s)

	householdID, err := s.CreateHousehold(ctx, owner, "Flat")
	if err != nil {
		t.Fatal(err)
	}
	m := mustMembership(t, s, householdID, owner)
	wantErr(t, "unconfirmed creates invite", s.CreateInvite(ctx, m), g.ErrorEmailNotVerified)
	wantErr(t, "unconfirmed sends invite", s.SendInviteEmail(ctx, m, "friend@example.com"), g.ErrorEmailNotVerified)
	noEmail(t, s)

	if err := s.db.SetEmailValidated(ctx, owner); err != nil {
		t.Fatal(err)
	}
	if err := s.SendInviteEmail(ctx, m, " Friend@Example.com "); err != nil {
		t.Fatal(err)
	}
	msg := nextEmail(t, s)
	if msg.To != "friend@example.com" || !strings.Contains(msg.Body, g.RouteInvite+"?code=") {
		t.Fatalf("unexpected invite email %+v", msg)
	}
	if !strings.Contains(msg.Subject, "Flat") {
		t.Fatalf("subject should name the household: %q", msg.Subject)
	}
}
