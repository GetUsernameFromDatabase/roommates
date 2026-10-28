package mailer

import (
	"net/mail"
	"strings"
	"testing"
)

func TestBuildMessage(t *testing.T) {
	from := &mail.Address{Name: "Roommates", Address: "noreply@example.com"}

	raw, err := buildMessage(from, Message{To: "a@example.com", Subject: "Kinnita e-post", Body: "Tere õde\nhttps://x/y?token=abc"})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := mail.ReadMessage(strings.NewReader(string(raw)))
	if err != nil {
		t.Fatal(err)
	}
	if got := parsed.Header.Get("To"); got != "a@example.com" {
		t.Errorf("To = %q", got)
	}
	if !strings.Contains(string(raw), "@example.com>") {
		t.Error("Message-ID should use the From domain")
	}

	_, err = buildMessage(from, Message{To: "a@example.com\r\nBcc: b@example.com", Subject: "x", Body: "x"})
	if err == nil {
		t.Error("newline in To must be rejected")
	}
}
