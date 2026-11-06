// Package config parses and validates command-line arguments and environment.
package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	// PasswordEnvVar is the environment variable holding the password.
	PasswordEnvVar = "PASSWORD"

	programName = "webdav-server"
)

var (
	// ErrHelp is returned when -h/--help is requested.
	ErrHelp = errors.New("help requested")
	// ErrVersion is returned when -v/--version is requested.
	ErrVersion = errors.New("version requested")
)

// Config holds validated runtime settings.
type Config struct {
	ListenAddr string
	ServeDir   string
	Username   string
	Password   string
}

// AuthEnabled reports whether HTTP basic authentication is required.
func (c *Config) AuthEnabled() bool {
	return c.Username != ""
}

// Parse reads arguments and the environment, then validates them.
// getenv is injected so the function can be tested without touching the process environment.
func Parse(args []string, getenv func(string) string) (*Config, error) {
	var (
		showHelp    bool
		showVersion bool
		listen      string
		serveDir    string
		username    string
	)

	fs := flag.NewFlagSet(programName, flag.ContinueOnError)
	fs.SetOutput(io.Discard) // usage is printed by the application, not the flag package

	fs.BoolVar(&showHelp, "help", false, "")
	fs.BoolVar(&showHelp, "h", false, "")
	fs.BoolVar(&showVersion, "version", false, "")
	fs.BoolVar(&showVersion, "v", false, "")
	fs.StringVar(&listen, "listen", "", "")
	fs.StringVar(&listen, "l", "", "")
	fs.StringVar(&serveDir, "serve", "", "")
	fs.StringVar(&serveDir, "s", "", "")
	fs.StringVar(&username, "username", "", "")
	fs.StringVar(&username, "u", "", "")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	if fs.NArg() > 0 {
		return nil, fmt.Errorf("unexpected argument: %q", fs.Arg(0))
	}
	if showHelp {
		return nil, ErrHelp
	}
	if showVersion {
		return nil, ErrVersion
	}

	addr, err := ParseListenAddr(listen)
	if err != nil {
		return nil, err
	}

	dir, err := validateServeDir(serveDir)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ListenAddr: addr,
		ServeDir:   dir,
	}

	if username != "" {
		password := getenv(PasswordEnvVar)
		if password == "" {
			return nil, fmt.Errorf("%s environment variable is required when --username is set", PasswordEnvVar)
		}
		cfg.Username = username
		cfg.Password = password
	}

	return cfg, nil
}

// validateServeDir checks that the provided directory exists and is a directory, returning its absolute path.
func validateServeDir(dir string) (string, error) {
	if dir == "" {
		return "", errors.New("--serve is required")
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("resolve serve directory: %w", err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("serve directory: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("serve path %q is not a directory", abs)
	}

	return abs, nil
}
