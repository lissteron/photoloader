package services

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/lissteron/photoloader/internal/app/core/ports"
	"github.com/lissteron/photoloader/pkg/elog"
)

type Photo struct {
	logger    elog.Logger
	photoRepo ports.PhotoRepo
}

func NewPhoto(logger elog.Logger, photoRepo ports.PhotoRepo) *Photo {
	return &Photo{
		logger:    logger,
		photoRepo: photoRepo,
	}
}

func (s *Photo) Upload(ctx context.Context, files []*multipart.FileHeader) error {
	if err := s.photoRepo.Upload(ctx, files); err != nil {
		return fmt.Errorf("photo repo upload: %w", err)
	}

	return nil
}
