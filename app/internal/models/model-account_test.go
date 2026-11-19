package models

import (
	"strings"
	"testing"
)

func TestChangeUsernameIsValid(t *testing.T) {
	tests := []struct {
		name, user string
		want       bool
	}{
		{"valid", "Mari", true},
		{"too short", "ab", false},
		{"too long", strings.Repeat("a", maxUsernameLength+1), false},
		{"trailing space", "Mari ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//exhaustruct:ignore
			m := ChangeUsername{Username: tt.user}
			if got, msgs := m.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v (messages %v)", got, tt.want, msgs)
			}
		})
	}
}

// the current password is only required, the new one follows the rules
func TestChangePasswordIsValid(t *testing.T) {
	const pw = "Abcdefg!"
	tests := []struct {
		name, current, pw, pw2 string
		want                   bool
	}{
		{"valid", "old", pw, pw, true},
		{"no current", "", pw, pw, false},
		{"weak new", "old", "abcdefgh", "abcdefgh", false},
		{"differ", "old", pw, pw + "x", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//exhaustruct:ignore
			m := ChangePassword{Current: tt.current, Password: tt.pw, Password2: tt.pw2}
			if got, msgs := m.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v (messages %v)", got, tt.want, msgs)
			}
		})
	}
}
