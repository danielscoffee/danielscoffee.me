package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/danielscoffee/danielscoffee.me/internal/app"
)

func gracefulShutdown(apiServer *http.Server, logger zerolog.Logger, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	logger.Info().Str("signal", ctx.Err().Error()).Msg("shutdown signal received")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
	} else {
		logger.Info().Msg("server shutdown complete")
	}

	done <- true
}

func main() {
	runtime, err := app.NewRuntime()
	if err != nil {
		panic(fmt.Sprintf("bootstrap error: %s", err))
	}

	runtime.Logger.Info().
		Int("port", runtime.Port).
		Str("site_url", runtime.SiteURL).
		Int("post_count", runtime.PostCount).
		Int("project_count", runtime.ProjectCount).
		Str("log_format", runtime.LogCfg.Format).
		Str("log_level", runtime.LogCfg.Level.String()).
		Msg("starting server")

	done := make(chan bool, 1)
	go gracefulShutdown(runtime.Server, runtime.Logger, done)

	err = runtime.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		runtime.Logger.Error().Err(err).Msg("http server error")
		os.Exit(1)
	}

	<-done
}
