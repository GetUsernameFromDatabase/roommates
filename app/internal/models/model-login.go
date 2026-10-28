package models

import (
	"net/mail"
	l "roommates/internal/locales"
	"strings"
)

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Login struct {
	ModelBase
	LoginRequest
}

func (m *Login) ValidateEmail() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	_, err := mail.ParseAddress(m.Email)
	if err != nil {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsEmailGeneric})
	}
	return msgs
}

func (m *Login) ValidatePassword() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if len(m.Password) < minPasswordLength {
		msgs = append(msgs, l.LKMessage{
			Key:  l.LK_ValidationErrorsPasswordLength,
			Args: []any{minPasswordLength},
		})
	}

	if strings.ToUpper(m.Password) == m.Password || strings.ToLower(m.Password) == m.Password {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsPasswordCase})
	}

	if !strings.ContainsAny(m.Password, passwordSpecialCharaters) {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsPasswordSymbol})
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
