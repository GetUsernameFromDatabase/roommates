// Package models https://templ.guide/syntax-and-usage/forms#create-a-view-model
package models

import (
	l "roommates/internal/locales"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 32
	minPasswordLength = 8
	// bcrypt (golang.org/x/crypto/bcrypt) refuses passwords longer than 72 bytes
	maxPasswordLength = 72
	// RFC 5321 limit for a forward path
	maxEmailLength = 254

	passwordSpecialCharaters = "!@#$%^&*()-_=+[]{}|;:',.<>?/`~"
)

// ModelBase fields are form:"-" so a posted form cannot set them
type ModelBase struct {
	// validation errors are shown if false
	Initial bool   `form:"-"`
	Error   string `form:"-"`
}

type Validator func() []l.LKMessage

type Validatable interface {
	Validate() []l.LKMessage
	GetValidators() []Validator
	IsValid() (bool, []l.LKMessage)
}

func IsModelValid(m Validatable) (bool, []l.LKMessage) {
	msgs := m.Validate()
	value := len(msgs) == 0
	return value, msgs
}

func ValidateModel(m Validatable) (msgs []l.LKMessage) {
	validators := m.GetValidators()
	for _, validator := range validators {
		msgs = append(msgs, validator()...)
	}
	return msgs
}
