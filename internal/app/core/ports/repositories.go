package ports

import (
	"context"
	"mime/multipart"
)

type PhotoRepo interface {
	Upload(ctx context.Context, files []*multipart.FileHeader) error
}
