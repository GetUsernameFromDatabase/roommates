package models

import (
	"errors"
	l "roommates/locales"
	"roommates/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type House struct {
	ModelBase
	// non-visual
	HouseID      pgtype.UUID `form:"house_id"`
	RoommateKeys []string    `form:"roommates[]"`
	// make sure to match indices with RoommateKeys
	//  username of the id
	RoommateLabels []string `form:"roommates_labels[]"`

	Name string `form:"name"`
	// used only by htmx to get user suggestions when adding roomates to house
	SearchedUser string `form:"searched_user"`
}

func (m *House) ValidateName() (msgs []l.LKMessage) {
	if m.Initial {
		return
	}

	if m.Name == "" {
		msgs = append(msgs, l.LKMessage{Key: l.LKFormsNameErrorEmpty})
		return msgs
	}

	charProblems := utils.ValidateString(m.Name, utils.RuneValidationRules{
		LettersAllowed:       true,
		DigitsAllowed:        true,
		MaxConsecutiveSpaces: -1,
	})

	digitOrLetterErrorAlreadyAdded := true
	for charProblem := range charProblems {
		switch charProblem {
		case utils.VSDigit, utils.VSLetter:
			if !digitOrLetterErrorAlreadyAdded {
				msgs = append(msgs, l.LKMessage{Key: l.LKFormsErrorsOnlyLettersAndDigits})
			}
			digitOrLetterErrorAlreadyAdded = true
		case utils.VSSpaces:
			msgs = append(msgs, l.LKMessage{Key: l.LKFormsErrorsNoMultipleSpaces})
		}
	}

	return msgs
}

// converts roomatekeys into UUID and filters out invalid keys
//
// if some keys were invalid then error will be set using the i18n with request context
//
//	if bool == true -- invalid UUID
func (m *House) FilterNonValidUUID(ctx *gin.Context) (bool, []pgtype.UUID) {
	hasInvalidUUID := false
	var roomateIDs []pgtype.UUID

	i := 0
	for i < len(m.RoommateKeys) {
		rKey := m.RoommateKeys[i]
		var id pgtype.UUID
		err := id.Scan(rKey)
		if err != nil {
			hasInvalidUUID = true
			m.RoommateKeys = append(m.RoommateKeys[:i], m.RoommateKeys[i+1:]...)
			m.RoommateLabels = append(m.RoommateLabels[:i], m.RoommateLabels[i+1:]...)
			continue
		}
		roomateIDs = append(roomateIDs, id)
		i++
	}

	if hasInvalidUUID {
		m.Error = utils.T(
			ctx.Request.Context(),
			l.LKFormsHouseErrorSomeRoomatesInvalid,
			"",
		)
	}
	return hasInvalidUUID, roomateIDs
}

func (m *House) GetValidators() []Validator {
	return []Validator{
		m.ValidateName,
	}
}

func (m *House) Validate() []l.LKMessage {
	if m.Initial {
		return nil
	}
	return ValidateModel(m)
}

// checks if the form is valid and sets the Initial to false
//
//	NB: this will not check if houseID is valid
func (m *House) IsValid() (bool, []l.LKMessage) {
	m.Initial = false
	return IsModelValid(m)
}

func (m *House) NeedsValidHouseID() error {
	if m.HouseID.Valid {
		return nil
	}
	return errors.New("house_id needs to be a valid UUID")
}
