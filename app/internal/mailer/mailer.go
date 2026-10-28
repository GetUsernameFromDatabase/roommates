// Package mailer sends plain text emails over SMTP
package mailer

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	g "roommates/internal/config"
	"roommates/internal/logger"
	"strings"
	"time"
)

var log = logger.Main.With().Str("package", "mailer").Logger()

type Message struct {
	To      string
	Subject string
	Body    string
}

type Mailer interface {
	Send(msg Message) error
}

// New returns an SMTP mailer from env, see [g.EVK_SMTPHost]
//
// in development without SMTP_HOST, emails are logged instead
func New() Mailer {
	host := g.EVK_SMTPHost.Get("")
	if host == "" && g.IsDevelopment() {
		log.Warn().Msg("SMTP_HOST not set, emails are written to the log")
		return LogMailer{}
	}
	return &SMTPMailer{
		Host:     g.EVK_SMTPHost.MustGet(),
		Port:     g.EVK_SMTPPort.Get("587"),
		Username: g.EVK_SMTPUsername.Get(""),
		Password: g.EVK_SMTPPassword.Get(""),
		From:     g.EVK_SMTPFrom.Get("Roommates <noreply@localhost>"),
	}
}

// LogMailer writes emails to the log, for development only as emails contain secrets
type LogMailer struct{}

func (LogMailer) Send(msg Message) error {
	log.Info().Str("to", msg.To).Str("subject", msg.Subject).Msg("email not sent:\n" + msg.Body)
	return nil
}

type SMTPMailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

const smtpTimeout = 30 * time.Second

// Send uses STARTTLS when the server offers it, net/smtp refuses to
// send the password unencrypted unless the host is localhost
func (m *SMTPMailer) Send(msg Message) error {
	from, err := mail.ParseAddress(m.From)
	if err != nil {
		return fmt.Errorf("SMTP_FROM: %w", err)
	}
	raw, err := buildMessage(from, msg)
	if err != nil {
		return err
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(m.Host, m.Port), smtpTimeout)
	if err != nil {
		return err
	}
	_ = conn.SetDeadline(time.Now().Add(smtpTimeout))
	c, err := smtp.NewClient(conn, m.Host)
	if err != nil {
		_ = conn.Close()
		return err
	}
	defer c.Close() //nolint:errcheck // Quit below reports the real outcome

	if ok, _ := c.Extension("STARTTLS"); ok {
		//exhaustruct:ignore
		if err := c.StartTLS(&tls.Config{ServerName: m.Host}); err != nil {
			return err
		}
	}
	if m.Username != "" {
		if err := c.Auth(smtp.PlainAuth("", m.Username, m.Password, m.Host)); err != nil {
			return err
		}
	}
	if err := c.Mail(from.Address); err != nil {
		return err
	}
	if err := c.Rcpt(msg.To); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(raw); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

func buildMessage(from *mail.Address, msg Message) ([]byte, error) {
	// header injection: a newline in To or Subject would add headers
	if strings.ContainsAny(msg.To+msg.Subject, "\r\n") {
		return nil, fmt.Errorf("newline in email header")
	}
	domain := from.Address[strings.LastIndex(from.Address, "@")+1:]

	var b bytes.Buffer
	fmt.Fprintf(&b, "From: %s\r\n", from.String())
	fmt.Fprintf(&b, "To: %s\r\n", msg.To)
	fmt.Fprintf(&b, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", msg.Subject))
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	fmt.Fprintf(&b, "Message-ID: <%s@%s>\r\n", rand.Text(), domain)
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	b.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	qp := quotedprintable.NewWriter(&b)
	if _, err := qp.Write([]byte(strings.ReplaceAll(msg.Body, "\n", "\r\n"))); err != nil {
		return nil, err
	}
	if err := qp.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
