package logging

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Format string
	Level  zerolog.Level
}

func ConfigFromEnv() Config {
	format := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))
	if format != "text" && format != "json" {
		format = "json"
	}

	level := zerolog.InfoLevel
	if raw := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))); raw != "" {
		if parsed, err := zerolog.ParseLevel(raw); err == nil {
			level = parsed
		}
	}

	return Config{Format: format, Level: level}
}

func New(cfg Config, w io.Writer) zerolog.Logger {
	if w == nil {
		w = os.Stdout
	}

	var writer io.Writer = w
	if cfg.Format == "text" {
		writer = zerolog.ConsoleWriter{Out: w, TimeFormat: time.RFC3339}
	}

	return zerolog.New(writer).Level(cfg.Level).With().Timestamp().Logger()
}
