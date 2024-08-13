package ports

import (
	"context"
	"mime/multipart"
)

type PhotoService interface {
	Upload(ctx context.Context, files []*multipart.FileHeader) error
}
