package config

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestEmptyEnvIsUnset(t *testing.T) {
	t.Setenv(EVK_S3Region.String(), "")
	if got := EVK_S3Region.Get("us-east-1"); got != "us-east-1" {
		t.Fatalf("empty value: Get = %q, want the fallback", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("empty value: MustGet did not panic")
		}
	}()
	EVK_S3Region.MustGet()
}

func TestCheckPublicEnv(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	t.Cleanup(func() { gin.SetMode(gin.DebugMode) })
	tests := []struct {
		platform, baseURL, smtpHost string
		ok                          bool
	}{
		{"", "https://127.0.0.1:8000", "mailpit", true}, // docker compose on a dev machine
		{"cloudflare", "https://flat.example.com", "smtp.resend.com", true},
		{"cloudflare", "https://127.0.0.1:8000", "smtp.resend.com", false},
		{"cloudflare", "https://localhost", "smtp.resend.com", false},
		{"cloudflare", "https://[::1]:8000", "smtp.resend.com", false},
		{"cloudflare", "http://flat.example.com", "smtp.resend.com", false},
		{"cloudflare", "", "smtp.resend.com", false},
		{"cloudflare", "https://flat.example.com", "mailpit", false},
	}
	for _, tt := range tests {
		t.Setenv(EVK_TrustedPlatform.String(), tt.platform)
		t.Setenv(EVK_AppBaseURL.String(), tt.baseURL)
		t.Setenv(EVK_SMTPHost.String(), tt.smtpHost)
		if err := CheckPublicEnv(); (err == nil) != tt.ok {
			t.Errorf("%+v: err = %v", tt, err)
		}
	}
}

func TestGetDatabaseURLEscapes(t *testing.T) {
	const password = `a b'c\d@e/f?g`
	t.Setenv(EVK_DatabaseHost.String(), "localhost")
	t.Setenv(EVK_DatabasePort.String(), "5432")
	t.Setenv(EVK_DatabaseName.String(), "roommates")
	t.Setenv(EVK_DatabaseUser.String(), "us er")
	t.Setenv(EVK_DatabasePassword.String(), password)

	pool, err := pgxpool.ParseConfig(GetDatabaseURL())
	if err != nil {
		t.Fatal(err)
	}
	cfg := pool.ConnConfig
	if cfg.Password != password || cfg.User != "us er" || cfg.Database != "roommates" ||
		cfg.Host != "localhost" || cfg.Port != 5432 {
		t.Fatalf("parsed %q %q %q %q %d", cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)
	}
	if pool.MaxConns != 10 {
		t.Fatalf("default MaxConns = %d", pool.MaxConns)
	}

	t.Setenv(EVK_DatabaseMaxConns.String(), "3")
	pool, err = pgxpool.ParseConfig(GetDatabaseURL())
	if err != nil {
		t.Fatal(err)
	}
	if pool.MaxConns != 3 {
		t.Fatalf("MaxConns = %d", pool.MaxConns)
	}
}
