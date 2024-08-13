package repositories

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/lissteron/photoloader/config"
	"github.com/lissteron/photoloader/pkg/elog"
)

type Photo struct {
	logger elog.Logger
	config *config.PhotoConfig
}

func NewPhoto(logger elog.Logger, cfg *config.PhotoConfig) *Photo {
	return &Photo{
		logger: logger,
		config: cfg,
	}
}

func (r *Photo) Upload(_ context.Context, files []*multipart.FileHeader) error {
	for _, fileHeader := range files {
		if err := r.uploadFile(fileHeader); err != nil {
			return fmt.Errorf("upload file: %w", err)
		}
	}

	return nil
}

func (r *Photo) uploadFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("file header open: %w", err)
	}

	defer file.Close()

	// Открываем файл для записи на диск
	dst, err := os.Create(filepath.Join(r.config.Dir, fileHeader.Filename))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer dst.Close()

	// Копируем содержимое загружаемого файла в файл на диске
	if _, err := io.Copy(dst, file); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	return nil
}
