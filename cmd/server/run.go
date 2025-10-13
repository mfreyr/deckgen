package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

func run(logger zerolog.Logger, server *http.Server) {
	serverError := make(chan error, 1)

	go func() {
		logger.Info().Str("address", fmt.Sprintf("http://%s", server.Addr)).Msg("server is starting")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	sigStop := make(chan os.Signal, 1)
	signal.Notify(sigStop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		logger.Error().Err(err).Msg("server error")
	case sig := <-sigStop:
		logger.Info().Str("signal", sig.String()).Msg("received shutdown signal")
	}

	logger.Info().Msg("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("server shutdown error")
		return
	}
	logger.Info().Msg("server exited properly")
}
