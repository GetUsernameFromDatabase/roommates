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

func TestNewDefaultPort(t *testing.T) {
	t.Setenv("SMTP_PORT", "")
	for host, want := range map[string]string{"mailpit": "1025", "localhost": "1025", "smtp.example.com": "587"} {
		t.Setenv("SMTP_HOST", host)
		if got := New().(*SMTPMailer).Port; got != want {
			t.Errorf("SMTP_HOST=%s: port %s, want %s", host, got, want)
		}
	}
}

func TestIsLocalHost(t *testing.T) {
	for host, want := range map[string]bool{
		"localhost": true, "127.0.0.1": true, "::1": true, "mailpit": true, "MailPit": true,
		"smtp.example.com": false, "10.0.0.5": false, "localhost.example.com": false,
	} {
		if got := isLocalHost(host); got != want {
			t.Errorf("isLocalHost(%q) = %v", host, got)
		}
	}
}
