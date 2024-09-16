package worker

import (
	"context"
	"fmt"
	"image/jpeg"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gribanoid/image-service/internal/pkg/syncer"
)

const (
	extJPG = ".jpg"
	extPNG = ".png"
)

type imageConvertor struct {
	imageDir       string
	frequencyS     int
	execAfterStart bool
	syncer         *syncer.Syncer

	chDone chan struct{}
}

func NewImageConvertor(imageDir string, frequencyS int, execAfterStart bool, syncer *syncer.Syncer) Worker {
	return &imageConvertor{
		imageDir:       imageDir,
		frequencyS:     frequencyS,
		execAfterStart: execAfterStart,
		syncer:         syncer,
		chDone:         make(chan struct{}),
	}
}

func (ic *imageConvertor) Start(ctx context.Context) error {
	if ic.execAfterStart {
		ic.convertImages(ctx)
	}

	ticker := time.NewTicker(time.Duration(ic.frequencyS) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ic.chDone:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ic.convertImages(ctx)
		}
	}
}

func (ic *imageConvertor) Stop() error {
	close(ic.chDone)

	return nil
}

func (ic *imageConvertor) convertImages(ctx context.Context) {
	walkFunc := func(path string, info os.FileInfo, errIn error) error {
		if errIn != nil {
			return errIn
		}

		if !info.IsDir() {
			return nil
		}

		// TODO: Можно лочить только когда встречается jpg?
		ic.syncer.Lock(info.Name())
		if err := filepath.Walk(ic.imageDir, walkDateDir); err != nil {
			return fmt.Errorf("filepath walk folder: %s err: %w", info.Name(), err)
		}
		ic.syncer.Unlock(info.Name())

		return nil
	}

	if err := filepath.Walk(ic.imageDir, walkFunc); err != nil {
		slog.ErrorContext(ctx, "filepath walk: %v", err)
	}
}

func walkDateDir(path string, info os.FileInfo, errIn error) error {
	if errIn != nil {
		return errIn
	}

	if filepath.Ext(path) == extJPG {
		if err := createPNGFromJPG(path); err != nil {
			return fmt.Errorf("create png file from jpg: %w", err)
		}

		if err := os.Remove(path); err != nil {
			return fmt.Errorf("remove file: %s err: %w", path, err)
		}
	}

	return nil
}

func createPNGFromJPG(path string) error {
	fileJPG, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open jpg file: %w", err)
	}
	defer func() {
		if err = fileJPG.Close(); err != nil {
			slog.Error("close jpg file: %v", err)
		}
	}()

	imgJPG, err := jpeg.Decode(fileJPG)
	if err != nil {
		return fmt.Errorf("decode jpeg: %w", err)
	}

	filenamePNG := strings.TrimSuffix(fileJPG.Name(), filepath.Ext(fileJPG.Name())) + extPNG

	filePNG, err := os.Create(filenamePNG)
	if err != nil {
		return fmt.Errorf("create png file: %w", err)
	}
	defer func() {
		if err = filePNG.Close(); err != nil {
			slog.Error("close png file: %v", err)
		}
	}()

	if err = png.Encode(filePNG, imgJPG); err != nil {
		return fmt.Errorf("encode png: %w", err)
	}

	return nil
}
