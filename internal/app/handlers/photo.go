package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lissteron/photoloader/internal/app/core/ports"
	"github.com/lissteron/photoloader/pkg/elog"
)

const _fileKey = "photos"

type Photo struct {
	logger       elog.Logger
	photoService ports.PhotoService
	path         string
}

func NewPhoto(logger elog.Logger, photoService ports.PhotoService, path string) *Photo {
	return &Photo{
		logger:       logger,
		photoService: photoService,
		path:         path,
	}
}

func (h *Photo) SetRounte(router chi.Router) {
	router.Post(h.path, h.Upload)
}

func (h *Photo) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Парсим multipart форму, не устанавливая лимит на размер
	err := r.ParseMultipartForm(0)
	if err != nil {
		h.logger.Errorf(ctx, "parse multipart form: %v", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Получаем файлы из формы по ключу
	files := r.MultipartForm.File[_fileKey]
	if len(files) == 0 {
		http.Error(w, "no files", http.StatusBadRequest)

		return
	}

	if err := h.photoService.Upload(ctx, files); err != nil {
		h.logger.Errorf(ctx, "photo service upload: %v", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("upload complete"))
}
