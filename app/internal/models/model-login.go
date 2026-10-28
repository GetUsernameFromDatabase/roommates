package models

import (
	"net/mail"
	l "roommates/internal/locales"
	"strings"
)

type LoginRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

type Login struct {
	ModelBase
	LoginRequest
	// path to return to after signing in, see config.QueryNext
	Next string `form:"next"`
}

func (m *Login) ValidateEmail() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}
	return validateEmail(m.Email)
}

func validateEmail(email string) (msgs []l.LKMessage) {
	// ParseAddress also accepts "Name <a@b.c>", only the bare address is allowed
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != strings.TrimSpace(email) || len(email) > maxEmailLength {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsEmailGeneric})
	}
	return msgs
}

// ValidatePassword only checks that a password was given,
// password rules are for creating one (see [Register.ValidatePassword])
// and should not be shown to someone trying to sign in
func (m *Login) ValidatePassword() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if m.Password == "" {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsPasswordRequired})
	}
	return msgs
}

func (m *Login) GetValidators() []Validator {
	return []Validator{
		m.ValidateEmail,
		m.ValidatePassword,
	}
}

func (m *Login) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// IsValid checks if the form is valid and sets the Initial to false
func (m *Login) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}
