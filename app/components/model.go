package components

import l "roomates/locales"

type Validator func() []l.LKMessage

type ValidatableModel interface {
	Validate() []l.LKMessage
	GetValidators() []Validator 
}

func IsModelValid(m ValidatableModel) (bool, []l.LKMessage) {
	msgs := m.Validate()
	value := false

	if len(msgs) > 0 {
		value = true
	}
	return value, msgs
}

func ValidateModel(m ValidatableModel) (msgs []l.LKMessage) {
	validators := m.GetValidators()
	for _, validator := range validators {
		msgs = append(msgs, validator()...)
	}
	return msgs
}

// --- keeping this here since it might be useful in the future
// check if this would error when not initial
// func (m *LoginModel) WouldError() (bool, []l.LKMessage) {
// 	// tested on https://go.dev/play/
// 	prevInitital := m.Initial
// 	defer func() {
// 		m.Initial = prevInitital
// 	}()

// 	m.Initial = false
// 	return m.HasErrors()
// }