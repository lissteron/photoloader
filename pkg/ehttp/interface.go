package ehttp

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/lissteron/photoloader/pkg/ehttp/internal"
)

type (
	Meter  = internal.Meter
	Logger = internal.Logger
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Get(ctx context.Context, uri string) (*http.Response, error)
	Post(ctx context.Context, uri, contentType string, body io.Reader) (*http.Response, error)
	PostForm(ctx context.Context, uri string, data url.Values) (*http.Response, error)
}

type Server interface {
	Router() chi.Router
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
