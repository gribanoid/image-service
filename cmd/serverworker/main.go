package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gribanoid/image-service/internal/app/server"
	"github.com/gribanoid/image-service/internal/app/worker"
	"github.com/gribanoid/image-service/internal/pkg/syncer"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	run(ctx)

	slog.InfoContext(ctx, "app stopped")
}

func run(ctx context.Context) {
	var wg sync.WaitGroup

	s := syncer.New()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Run(ctx, s); err != nil {
			slog.ErrorContext(ctx, "server stopped with error: %v", err)

			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := worker.Run(ctx, s); err != nil {
			slog.ErrorContext(ctx, "worker stopped with error: %v", err)

			return
		}
	}()

	wg.Wait()
}
