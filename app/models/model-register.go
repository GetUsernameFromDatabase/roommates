package models

import (
	l "roommates/locales"
	"strings"
)

type Register struct {
	Login
	Username  string `form:"username"`
	Password2 string `form:"password_2"`
}

func (m *Register) ValidateUsername() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if m.Username != strings.TrimSpace(m.Username) {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsUsernameSpaces})
	}

	if len(m.Username) < minUsernameLength {
		msgs = append(msgs, l.LKMessage{
			Key:  l.LK_ValidationErrorsUsernameLength,
			Args: []any{minUsernameLength},
		})
	}

	return msgs
}

// ValidatePasswordMatch validates if both passwords match
func (m *Register) ValidatePasswordMatch() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if m.Password != m.Password2 {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsPasswordMatch})
	}

	return msgs
}

func (m *Register) GetValidators() []Validator {
	return []Validator{
		m.ValidateEmail,
		m.ValidatePassword,
		m.ValidateUsername,
		m.ValidatePasswordMatch,
	}
}

func (m *Register) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// IsValid checks if the form is valid and sets the Initial to false
func (m *Register) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}
