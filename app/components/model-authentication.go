package components

import (
	"net/mail"
	l "roommates/locales"
	"strings"
)

type LoginModel struct {
	// validation errors are shown if false
	Initial bool
	Error   string

	Email    string `form:"email"`
	Password string `form:"password"`
}

func (m *LoginModel) ValidateEmail() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	_, err := mail.ParseAddress(m.Email)
	if err != nil {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsEmailErrorGeneric})
	}
	return msgs
}

func (m *LoginModel) ValidatePassword() (msgs []l.LKMessage) {
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

func (m *LoginModel) GetValidators() []Validator {
	return []Validator{
		m.ValidateEmail,
		m.ValidatePassword,
	}
}

func (m *LoginModel) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// checks if this has errors or not
//
// keep in mind that if Initial == true then this will not check errors
func (m *LoginModel) IsValid() (bool, []l.LKMessage) {
	return IsModelValid(m)
}

// -----------------------------------------------------------------------------

type RegisterModel struct {
	LoginModel
	Username  string `form:"username"`
	Password2 string `form:"password_2"`
}

func (m *RegisterModel) ValidateUsername() (msgs []l.LKMessage) {
	if m.Initial {
		return
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
func (m *RegisterModel) ValidatePasswordMatch() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if m.Password != m.Password2 {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsPasswordErrorMustMatch})
	}

	return msgs
}

func (m *RegisterModel) GetValidators() []Validator {
	return []Validator{
		m.ValidateEmail,
		m.ValidatePassword,
		m.ValidateUsername,
		m.ValidatePasswordMatch,
	}
}

func (m *RegisterModel) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// checks if this has errors or not
//
// keep in mind that if Initial == true then this will not check errors
func (m *RegisterModel) IsValid() (bool, []l.LKMessage) {
	return IsModelValid(m)
}
