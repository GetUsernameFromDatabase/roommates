package models

import (
	"net/mail"
	l "roommates/locales"
	"strings"
)

type Login struct {
	ModelBase
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (m *Login) ValidateEmail() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	_, err := mail.ParseAddress(m.Email)
	if err != nil {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsEmailErrorGeneric})
	}
	return msgs
}

func (m *Login) ValidatePassword() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	const minPasswordLength = 8
	if len(m.Password) < minPasswordLength {
		msgs = append(msgs, l.LKMessage{
			Key:  l.LKFormsPasswordErrorLength,
			Args: []any{minPasswordLength},
		})
	}

	if strings.ToUpper(m.Password) == m.Password || strings.ToLower(m.Password) == m.Password {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsPasswordErrorCase})
	}

	const specialChars = "!@#$%^&*()-_=+[]{}|;:',.<>?/`~"
	if !strings.ContainsAny(m.Password, specialChars) {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsPasswordErrorSymbol})
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

// checks if the form is valid and sets the Initial to false
func (m *Login) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}
