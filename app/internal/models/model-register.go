package models

import (
	l "roommates/internal/locales"
	"strings"
	"unicode/utf8"
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

	usernameLength := utf8.RuneCountInString(m.Username)
	if usernameLength < minUsernameLength {
		msgs = append(msgs, l.LKMessage{
			Key:  l.LK_ValidationErrorsUsernameLength,
			Args: []any{minUsernameLength},
		})
	}
	if usernameLength > maxUsernameLength {
		msgs = append(msgs, l.LKMessage{
			Key:  l.LK_ValidationErrorsUsernameMaxLength,
			Args: []any{maxUsernameLength},
		})
	}

	return msgs
}

// ValidatePassword checks the password rules for a new password
func (m *Register) ValidatePassword() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}
	return validateNewPassword(m.Password)
}

// validateNewPassword checks the rules for choosing a password
func validateNewPassword(password string) (msgs []l.LKMessage) {
	if len(password) < minPasswordLength {
		msgs = append(msgs, l.LKMessage{
			Key:  l.LK_ValidationErrorsPasswordLength,
			Args: []any{minPasswordLength},
		})
	}

	// bcrypt refuses longer passwords, counted in bytes
	if len(password) > maxPasswordLength {
		msgs = append(msgs, l.LKMessage{
			Key:  l.LK_ValidationErrorsPasswordMaxLength,
			Args: []any{maxPasswordLength},
		})
	}

	if strings.ToUpper(password) == password || strings.ToLower(password) == password {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsPasswordCase})
	}

	if !strings.ContainsAny(password, passwordSpecialCharaters) {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsPasswordSymbol})
	}

	return msgs
}

// ValidatePasswordMatch validates if both passwords match
func (m *Register) ValidatePasswordMatch() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	return validatePasswordMatch(m.Password, m.Password2)
}

func validatePasswordMatch(password, password2 string) (msgs []l.LKMessage) {
	if password != password2 {
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
