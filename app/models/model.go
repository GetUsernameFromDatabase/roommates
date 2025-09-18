package models

import l "roommates/locales"

type ModelBase struct {
	// validation errors are shown if false
	Initial bool
	Error   string
}

type Validator func() []l.LKMessage

type Validatable interface {
	Validate() []l.LKMessage
	GetValidators() []Validator
}

func IsModelValid(m Validatable) (bool, []l.LKMessage) {
	msgs := m.Validate()
	value := false

	if len(msgs) > 0 {
		value = true
	}
	return value, msgs
}

func ValidateModel(m Validatable) (msgs []l.LKMessage) {
	validators := m.GetValidators()
	for _, validator := range validators {
		msgs = append(msgs, validator()...)
	}
	return msgs
}
