package models

import (
	l "roommates/internal/locales"
)

func validateCurrentPassword(password string) (msgs []l.LKMessage) {
	if password == "" {
		msgs = append(msgs, l.LKMessage{Key: l.LK_ValidationErrorsPasswordRequired})
	}
	return msgs
}

// ChangePassword is the change password form of a signed in user
type ChangePassword struct {
	ModelBase
	Current   string `form:"current_password"`
	Password  string `form:"password"`
	Password2 string `form:"password_2"`
}

func (m *ChangePassword) ValidateCurrent() []l.LKMessage {
	return validateCurrentPassword(m.Current)
}

func (m *ChangePassword) ValidatePassword() []l.LKMessage {
	return validateNewPassword(m.Password)
}

func (m *ChangePassword) ValidatePasswordMatch() []l.LKMessage {
	return validatePasswordMatch(m.Password, m.Password2)
}

func (m *ChangePassword) GetValidators() []Validator {
	return []Validator{m.ValidateCurrent, m.ValidatePassword, m.ValidatePasswordMatch}
}

func (m *ChangePassword) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// IsValid checks if the form is valid and sets the Initial to false
func (m *ChangePassword) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}

// -----------------------------------------------------------------------------

// ChangeEmail is the new email of a signed in user, confirmed with the password
type ChangeEmail struct {
	ModelBase
	Current string `form:"current_password"`
	Email   string `form:"email"`
}

func (m *ChangeEmail) ValidateCurrent() []l.LKMessage {
	return validateCurrentPassword(m.Current)
}

func (m *ChangeEmail) ValidateEmail() []l.LKMessage {
	return validateEmail(m.Email)
}

func (m *ChangeEmail) GetValidators() []Validator {
	return []Validator{m.ValidateCurrent, m.ValidateEmail}
}

func (m *ChangeEmail) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// IsValid checks if the form is valid and sets the Initial to false
func (m *ChangeEmail) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}

// -----------------------------------------------------------------------------

// ChangeUsername is the new username of a signed in user
type ChangeUsername struct {
	ModelBase
	Username string `form:"username"`
}

func (m *ChangeUsername) ValidateUsername() []l.LKMessage {
	return validateUsername(m.Username)
}

func (m *ChangeUsername) GetValidators() []Validator {
	return []Validator{m.ValidateUsername}
}

func (m *ChangeUsername) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// IsValid checks if the form is valid and sets the Initial to false
func (m *ChangeUsername) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}
