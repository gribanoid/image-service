package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gribanoid/image-service/internal/pkg/service"
)

var (
	errInvalidDateFormat    = fmt.Errorf("invalid date format")
	errInvalidMultipartForm = fmt.Errorf("invalid form-data")
	errInvalidImageBody     = fmt.Errorf("invalid image body")
	errInternal             = fmt.Errorf("internal")
)

type Handler struct {
	imageSVC *service.Image
}

func NewHandler(imageSVC *service.Image) *Handler {
	return &Handler{imageSVC: imageSVC}
}

func (h *Handler) UploadImage(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Error(errInvalidMultipartForm))

		return
	}

	fileHeaders := form.File["image"]
	if len(fileHeaders) != 1 {
		c.JSON(http.StatusBadRequest, c.Error(errInvalidImageBody))

		return
	}

	file, err := fileHeaders[0].Open()
	if err != nil {
		slog.Log(c, slog.LevelError, err.Error())
		c.JSON(http.StatusInternalServerError, c.Error(errInternal))

		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			slog.ErrorContext(c, "download image close file: %v", err)
		}
	}()

	fileContent := make([]byte, fileHeaders[0].Size)

	if _, err = file.Read(fileContent); err != nil {
		slog.Log(c, slog.LevelError, err.Error())
		c.JSON(http.StatusInternalServerError, c.Error(errInternal))

		return
	}

	path, err := h.imageSVC.Upload(fileContent)
	if err != nil {
		slog.Log(c, slog.LevelError, err.Error())
		c.JSON(http.StatusInternalServerError, c.Error(errInternal))

		return
	}

	c.JSON(http.StatusOK, gin.H{"filepath": path})
}

func (h *Handler) DownloadImagesByDate(c *gin.Context) {
	date := c.Query("date")

	if _, err := time.Parse(time.DateOnly, date); err != nil {
		c.JSON(http.StatusBadRequest, c.Error(errInvalidDateFormat))

		return
	}

	zipContent, err := h.imageSVC.DownloadZIP(c, date)
	if err != nil {
		if errors.Is(err, service.ErrDirNotFound) {
			c.JSON(http.StatusNoContent, nil)
		}

		slog.Log(c, slog.LevelError, err.Error())
		c.JSON(http.StatusInternalServerError, c.Error(errInternal))

		return
	}

	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", date))

	if _, err = c.Writer.Write(zipContent); err != nil {
		slog.Log(c, slog.LevelError, err.Error())
		c.JSON(http.StatusInternalServerError, c.Error(errInternal))
	}
}
