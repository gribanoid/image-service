package server

import (
	"context"
	"fmt"
	"github.com/gribanoid/image-service/internal/pkg/syncer"
	"log/slog"
	"sync"

	"github.com/gribanoid/image-service/internal/app/server/http"
	"github.com/gribanoid/image-service/internal/config"
	"github.com/gribanoid/image-service/internal/pkg/service"
)

type server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func Run(ctx context.Context, syncer *syncer.Syncer) error {
	var (
		wg  sync.WaitGroup
		cfg = config.Default
		srv = provideServerHTTP(ctx, cfg.Server, &wg, syncer)
	)

	select {
	case <-ctx.Done():
		stopServerHTTP(ctx, srv)
	}

	wg.Wait()

	return nil
}

func provideServerHTTP(ctx context.Context, cfg *config.Server, wg *sync.WaitGroup, syncer *syncer.Syncer) server {
	var (
		imageSVC = service.NewImage(cfg.ImageDir, syncer)
		handler  = http.NewHandler(imageSVC)
		srvHTTP  = http.NewServer(cfg.Port, handler)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srvHTTP.ListenAndServe(); err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("serve http: %v", err))
		}
	}()

	return srvHTTP
}

func stopServerHTTP(ctx context.Context, srv server) {
	if err := srv.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("stop http server: %v", err))
	}
}
