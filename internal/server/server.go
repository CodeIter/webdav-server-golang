// Package server builds the WebDAV HTTP handler and manages its lifecycle.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/net/webdav"

	"github.com/CodeIter/webdav-server-golang/internal/config"
	"github.com/CodeIter/webdav-server-golang/pkg/auth"
)

const (
	readHeaderTimeout = 10 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 10 * time.Second
	authRealm         = "webdav"
)

// NewHandler returns the WebDAV handler, wrapped in basic auth when configured.
func NewHandler(cfg *config.Config, logger *slog.Logger) http.Handler {
	dav := &webdav.Handler{
		Prefix:     "/",
		FileSystem: webdav.Dir(cfg.ServeDir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				logger.Warn("webdav request failed",
					"method", r.Method,
					"path", r.URL.Path,
					"error", err,
				)
			}
		},
	}

	if !cfg.AuthEnabled() {
		return dav
	}
	return auth.BasicAuth(authRealm, cfg.Username, cfg.Password)(dav)
}

// Run serves handler on addr until ctx is cancelled, then shuts down gracefully.
// Read and write timeouts are deliberately not set, since large uploads and downloads are expected.
func Run(ctx context.Context, addr string, handler http.Handler, logger *slog.Logger) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
	}

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
