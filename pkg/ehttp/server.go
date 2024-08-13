package ehttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/lissteron/photoloader/pkg/helpers"
)

const (
	_defaultIdleTimeout       = time.Minute
	_defaultReadTimeout       = time.Minute
	_defaultWriteTimeout      = time.Minute
	_defaultGlobalTimeout     = time.Minute
	_defaultCloseTimeout      = 10 * time.Second
	_defaultReadHeaderTimeout = 10 * time.Second
)

type ServerConfig struct {
	ListenAddr        string
	IdleTimeout       time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	GlobalTimeout     time.Duration
	CloseTimeout      time.Duration
	ReadHeaderTimeout time.Duration
}

type ServerPKG struct {
	config   *ServerConfig
	server   *http.Server
	listener net.Listener
	router   chi.Router

	logger Logger
	meter  Meter

	slowRequestDuration time.Duration
}

func NewServer(cfg *ServerConfig, options ...ServerOption) *ServerPKG {
	server := &ServerPKG{
		config: cfg,
		router: chi.NewRouter(),
	}

	server.initDefaultConfig()

	server.server = &http.Server{
		ReadHeaderTimeout: _defaultReadHeaderTimeout,
		ReadTimeout:       server.config.ReadTimeout,
		WriteTimeout:      server.config.WriteTimeout,
		IdleTimeout:       server.config.IdleTimeout,
		Handler:           server.router,
	}

	for _, opt := range options {
		opt.apply(server)
	}

	server.router.Use(server.recoverMiddleware)

	return server
}

func (s *ServerPKG) Router() chi.Router {
	return s.router
}

func (s *ServerPKG) Start(_ context.Context) error {
	var err error

	if s.listener == nil {
		s.listener, err = net.Listen("tcp", s.config.ListenAddr)
		if err != nil {
			return fmt.Errorf("init listner: %w", err)
		}
	}

	if err := s.server.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http failed: %w", err)
	}

	return nil
}

func (s *ServerPKG) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.config.CloseTimeout)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	return nil
}

func (s *ServerPKG) initDefaultConfig() {
	if s.config.GlobalTimeout <= 0 {
		s.config.GlobalTimeout = _defaultGlobalTimeout
	}

	if s.config.IdleTimeout <= 0 {
		s.config.IdleTimeout = _defaultIdleTimeout
	}

	if s.config.ReadTimeout <= 0 {
		s.config.ReadTimeout = _defaultReadTimeout
	}

	if s.config.WriteTimeout <= 0 {
		s.config.WriteTimeout = _defaultWriteTimeout
	}

	if s.config.CloseTimeout <= 0 {
		s.config.CloseTimeout = _defaultCloseTimeout
	}

	if s.config.ReadHeaderTimeout <= 0 {
		s.config.ReadHeaderTimeout = _defaultReadHeaderTimeout
	}
}

func (s *ServerPKG) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(ctx context.Context) {
			if rec := recover(); rec != nil {
				if helpers.IsNotNil(s.logger) {
					s.logger.Error(
						ctx,
						fmt.Sprintf("http server recovered from panic: %v", rec),
						"stack",
						string(debug.Stack()),
					)
				}

				w.WriteHeader(http.StatusInternalServerError)
			}
		}(r.Context())

		next.ServeHTTP(w, r)
	})
}
