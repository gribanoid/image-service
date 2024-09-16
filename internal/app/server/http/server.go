package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	srv     *http.Server
	handler *Handler
}

func NewServer(port int, handler *Handler) *server {
	return &server{
		handler: handler,
		srv:     &http.Server{Addr: fmt.Sprintf(":%d", port)},
	}
}

func (s *server) ListenAndServe() error {
	r := gin.Default()
	r.GET("/", ping)
	r.GET("/download", s.handler.DownloadImagesByDate)
	r.POST("/upload", s.handler.UploadImage)

	s.srv.Handler = r

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("ListenAndServe: %w", err)
	}

	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// TODO: instead of healtz check.
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
