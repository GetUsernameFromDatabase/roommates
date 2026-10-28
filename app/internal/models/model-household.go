package models

import (
	l "roommates/internal/locales"
	"strings"
	"unicode/utf8"
)

type Household struct {
	ModelBase
	Name string `form:"name"`
}

func (m *Household) ValidateName() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	nameLength := utf8.RuneCountInString(strings.TrimSpace(m.Name))
	if nameLength < minHouseholdNameLength || nameLength > maxHouseholdNameLength {
		msgs = append(msgs, l.LKMessage{
			Key:  l.LK_ValidationErrorsHouseholdNameLength,
			Args: []any{minHouseholdNameLength, maxHouseholdNameLength},
		})
	}
	return msgs
}

func (m *Household) GetValidators() []Validator {
	return []Validator{m.ValidateName}
}

func (m *Household) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// IsValid checks if the form is valid and sets the Initial to false
func (m *Household) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}
