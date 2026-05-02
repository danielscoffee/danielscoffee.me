package logging

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("LOG_FORMAT", "")
	t.Setenv("LOG_LEVEL", "")

	cfg := ConfigFromEnv()
	if cfg.Format != "json" {
		t.Fatalf("expected default json format, got %q", cfg.Format)
	}
	if cfg.Level != zerolog.InfoLevel {
		t.Fatalf("expected default info level, got %s", cfg.Level)
	}
}

func TestConfigFromEnv_ParsesValues(t *testing.T) {
	t.Setenv("LOG_FORMAT", "text")
	t.Setenv("LOG_LEVEL", "debug")

	cfg := ConfigFromEnv()
	if cfg.Format != "text" {
		t.Fatalf("expected text format, got %q", cfg.Format)
	}
	if cfg.Level != zerolog.DebugLevel {
		t.Fatalf("expected debug level, got %s", cfg.Level)
	}
}

func TestConfigFromEnv_InvalidFallsBack(t *testing.T) {
	t.Setenv("LOG_FORMAT", "weird")
	t.Setenv("LOG_LEVEL", "invalid")

	cfg := ConfigFromEnv()
	if cfg.Format != "json" {
		t.Fatalf("expected fallback json format, got %q", cfg.Format)
	}
	if cfg.Level != zerolog.InfoLevel {
		t.Fatalf("expected fallback info level, got %s", cfg.Level)
	}
}
