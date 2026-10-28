package models

import (
	"strings"
	"testing"
)

func TestRegisterIsValid(t *testing.T) {
	const pw = "Abcdefg!"
	tests := []struct {
		name                       string
		email, user, pw, pwConfirm string
		want                       bool
	}{
		{"valid", "a@b.co", "abc", pw, pw, true},
		{"display name email", "Bob <a@b.co>", "abc", pw, pw, false},
		{"not an email", "a.b.co", "abc", pw, pw, false},
		{"password 72 bytes", "a@b.co", "abc", "A!" + strings.Repeat("a", 70), "A!" + strings.Repeat("a", 70), true},
		{"password 73 bytes", "a@b.co", "abc", "A!" + strings.Repeat("a", 71), "A!" + strings.Repeat("a", 71), false},
		{"password too short", "a@b.co", "abc", "Abcde!", "Abcde!", false},
		{"password no capital", "a@b.co", "abc", "abcdefg!", "abcdefg!", false},
		{"password no symbol", "a@b.co", "abc", "Abcdefgh", "Abcdefgh", false},
		{"passwords differ", "a@b.co", "abc", pw, pw + "x", false},
		{"username too short", "a@b.co", "ab", pw, pw, false},
		{"username leading space", "a@b.co", " abc", pw, pw, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//exhaustruct:ignore
			m := Register{Username: tt.user, Password2: tt.pwConfirm}
			m.Email, m.Password = tt.email, tt.pw
			if got, msgs := m.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v (messages %v)", got, tt.want, msgs)
			}
		})
	}
}

// login only checks that an email and a password were given,
// not the password rules
func TestLoginIsValid(t *testing.T) {
	tests := []struct {
		name, email, pw string
		want            bool
	}{
		{"valid", "a@b.co", "x", true},
		{"display name email", "Bob <a@b.co>", "x", false},
		{"no password", "a@b.co", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//exhaustruct:ignore
			m := Login{}
			m.Email, m.Password = tt.email, tt.pw
			if got, msgs := m.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v (messages %v)", got, tt.want, msgs)
			}
		})
	}
}
