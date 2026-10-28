// Package models https://templ.guide/syntax-and-usage/forms#create-a-view-model
package models

import (
	l "roommates/locales"
)

const (
	minUsernameLength = 3
	minPasswordLength = 8

	passwordSpecialCharaters = "!@#$%^&*()-_=+[]{}|;:',.<>?/`~"
)

type ModelBase struct {
	// validation errors are shown if false
	Initial bool
	Error   string
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
