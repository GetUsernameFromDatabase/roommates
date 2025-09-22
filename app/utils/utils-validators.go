package utils

import "unicode"

type RuneValidationRules struct {
	LettersAllowed bool
	DigitsAllowed  bool
	// if -1 then will not check how many spaces are "chained" on after another
	MaxConsecutiveSpaces int
}

type StringValidatorProblems int

const (
	VSLetter StringValidatorProblems = iota
	VSDigit
	VSSpaces
)

// validates string checking only allowing things given in parameters
func ValidateString(s string, aris RuneValidationRules) map[StringValidatorProblems]bool {
	invalids := make(map[StringValidatorProblems]bool)
	// this is so that I will not have unecessary -1 check
	checkChainedSpaces := func(r rune) bool {
		return false
	}
	if aris.MaxConsecutiveSpaces != -1 {
		// I could make it a closure but I won't
		spaceCount := 0
		// returns if rune was ' '
		checkChainedSpaces = func(r rune) bool {
			if r == ' ' {
				spaceCount++
				if spaceCount > aris.MaxConsecutiveSpaces {
					invalids[VSSpaces] = true
				}
				return true
			}
			spaceCount = 0
			return false
		}
	}

	for _, r := range s {
		if checkChainedSpaces(r) {
			continue
		}
		isLetter := unicode.IsLetter(r)
		isDigit := unicode.IsDigit(r)

		// skip allowed cases
		if aris.LettersAllowed && isLetter {
			continue
		} else if aris.DigitsAllowed && isDigit {
			continue
		}

		if aris.LettersAllowed && !isLetter {
			invalids[VSLetter] = true
			continue
		}
		if aris.DigitsAllowed && !isDigit {
			invalids[VSDigit] = true
			continue
		}
	}
	return invalids
}
