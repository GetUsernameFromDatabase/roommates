package models

import (
	l "roommates/locales"
	"strings"
)

type Register struct {
	Login
	Username         string `form:"username"`
	Password2        string `form:"password_2"`
	FullName         string `form:"full_name"`
	IsFullNamePublic bool   `form:"is_full_name_public"`
}

func (m *Register) ValidateUsername() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if m.Username != strings.TrimSpace(m.Username) {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsUsernameErrorSpaces})
	}

	const minLength = 3
	if len(m.Username) < minLength {
		msgs = append(msgs, l.LKMessage{
			Key:  l.LKFormsUsernameErrorLength,
			Args: []any{minLength},
		})
	}

	return msgs
}

// validates if both passwords match
func (m *Register) ValidatePasswordMatch() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if m.Password != m.Password2 {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsPasswordErrorMustMatch})
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

// checks if this has errors or not
//
// keep in mind that if Initial == true then this will not check errors
func (m *Register) IsValid() (bool, []l.LKMessage) {
	return IsModelValid(m)
}
