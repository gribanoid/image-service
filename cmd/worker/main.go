package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gribanoid/image-service/internal/app/worker"
	"github.com/gribanoid/image-service/internal/pkg/syncer"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := worker.Run(ctx, syncer.New()); err != nil {
		slog.ErrorContext(ctx, "app stopped with error: %v", err)

		return
	}

	slog.InfoContext(ctx, "app stopped")
}
