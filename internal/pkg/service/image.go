package service

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/gribanoid/image-service/internal/pkg/syncer"
)

const extPNG = ".png"

var ErrDirNotFound = fmt.Errorf("directory not found")

type Image struct {
	imageDir string
	syncer   *syncer.Syncer

	mu *sync.Mutex
}

func NewImage(imageDir string, syncer *syncer.Syncer) *Image {
	return &Image{imageDir: imageDir, syncer: syncer, mu: &sync.Mutex{}}
}

// Upload TODO: add image mime type check.
func (i *Image) Upload(content []byte) (string, error) {
	var (
		date     = time.Now().Format(time.DateOnly)
		dirPath  = filepath.Join(i.imageDir, date)
		filePath = filepath.Join(dirPath, fmt.Sprintf("%s.jpg", uuid.New().String()))
	)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("create directories: %w", err)
	}

	if exists := i.syncer.Lock(date); !exists {
		i.mu.Lock()
		i.syncer.AddDate(date)
		i.mu.Unlock()
		i.syncer.Lock(date)
	}
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	i.syncer.Unlock(date)

	return filePath, nil
}

func (i *Image) DownloadZIP(ctx context.Context, date string) ([]byte, error) {
	dirPath := filepath.Join(i.imageDir, date)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, ErrDirNotFound
	}

	var (
		buf = new(bytes.Buffer)
		zw  = zip.NewWriter(buf)
	)

	i.syncer.Lock(date)
	defer i.syncer.Unlock(date)

	walkFunc := func(path string, info os.FileInfo, errIn error) error {
		if errIn != nil {
			return errIn
		}

		if filepath.Ext(path) == extPNG {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("os read file: %w", err)
			}

			f, err := zw.Create(path)
			if err != nil {
				return fmt.Errorf("zip writer create: %w", err)
			}

			if _, err = f.Write(data); err != nil {
				return fmt.Errorf("zip file write: %w", err)
			}
		}

		return nil
	}

	if err := filepath.Walk(dirPath, walkFunc); err != nil {
		slog.ErrorContext(ctx, "filepath walk: %v", err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("close zip writer: %w", err)
	}

	return buf.Bytes(), nil
}
