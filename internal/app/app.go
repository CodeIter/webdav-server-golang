// Package app wires configuration, logging, and the server together.
package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/CodeIter/webdav-server-golang/internal/config"
	"github.com/CodeIter/webdav-server-golang/internal/server"
	"github.com/CodeIter/webdav-server-golang/version"
)

const (
	exitOK      = 0
	exitFailure = 1
	exitUsage   = 2
)

// Run executes the CLI and returns a process exit code.
func Run(args []string, stdout, stderr io.Writer, getenv func(string) string) int {
	cfg, err := config.Parse(args, getenv)
	switch {
	case errors.Is(err, config.ErrHelp):
		config.PrintUsage(stdout)
		return exitOK
	case errors.Is(err, config.ErrVersion):
		fmt.Fprintln(stdout, version.String())
		return exitOK
	case err != nil:
		fmt.Fprintf(stderr, "error: %v\n\n", err)
		config.PrintUsage(stderr)
		return exitUsage
	}

	logger := slog.New(slog.NewTextHandler(stderr, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("configuration loaded",
		"listen", cfg.ListenAddr,
		"dir", cfg.ServeDir,
		"auth", cfg.AuthEnabled(),
	)

	handler := server.NewHandler(cfg, logger)
	if err := server.Run(ctx, cfg.ListenAddr, handler, logger); err != nil {
		logger.Error("server stopped with error", "error", err)
		return exitFailure
	}
	return exitOK
}
