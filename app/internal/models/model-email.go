package models

import (
	l "roommates/internal/locales"
)

// EmailForm is a form with only an email, e.g. forgot password or invite by email
type EmailForm struct {
	ModelBase
	Email string `form:"email"`
}

func (m *EmailForm) ValidateEmail() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}
	return validateEmail(m.Email)
}

func (m *EmailForm) GetValidators() []Validator {
	return []Validator{m.ValidateEmail}
}

func (m *EmailForm) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// IsValid checks if the form is valid and sets the Initial to false
func (m *EmailForm) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}

// -----------------------------------------------------------------------------

// ResetPassword is the new password form of a password reset link
type ResetPassword struct {
	ModelBase
	Token     string `form:"token"`
	Password  string `form:"password"`
	Password2 string `form:"password_2"`
}

func (m *ResetPassword) ValidatePassword() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}
	return validateNewPassword(m.Password)
}

func (m *ResetPassword) ValidatePasswordMatch() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}
	return validatePasswordMatch(m.Password, m.Password2)
}

func (m *ResetPassword) GetValidators() []Validator {
	return []Validator{m.ValidatePassword, m.ValidatePasswordMatch}
}

func (m *ResetPassword) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// IsValid checks if the form is valid and sets the Initial to false
func (m *ResetPassword) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}
