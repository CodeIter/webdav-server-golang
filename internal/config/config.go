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
	CertFile   string
	KeyFile    string
}

// AuthEnabled reports whether HTTP basic authentication is required.
func (c *Config) AuthEnabled() bool {
	return c.Username != ""
}

// TLSEnabled reports whether HTTPS is configured with both a certificate and key.
func (c *Config) TLSEnabled() bool {
	return c.CertFile != "" && c.KeyFile != ""
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
		certFile    string
		keyFile     string
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
	fs.StringVar(&certFile, "cert", "", "")
	fs.StringVar(&certFile, "c", "", "")
	fs.StringVar(&keyFile, "key", "", "")
	fs.StringVar(&keyFile, "k", "", "")

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

	certPath, keyPath, err := validateTLSFiles(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ListenAddr: addr,
		ServeDir:   dir,
		CertFile:   certPath,
		KeyFile:    keyPath,
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

// validateTLSFiles checks that both certFile and keyFile are provided together, and that they exist.
func validateTLSFiles(certFile, keyFile string) (string, string, error) {
	if certFile == "" && keyFile == "" {
		return "", "", nil
	}
	if certFile == "" || keyFile == "" {
		return "", "", errors.New("both --cert and --key must be provided together")
	}

	certAbs, err := filepath.Abs(certFile)
	if err != nil {
		return "", "", fmt.Errorf("resolve certificate path: %w", err)
	}
	if _, err := os.Stat(certAbs); err != nil {
		return "", "", fmt.Errorf("certificate file: %w", err)
	}

	keyAbs, err := filepath.Abs(keyFile)
	if err != nil {
		return "", "", fmt.Errorf("resolve private key path: %w", err)
	}
	if _, err := os.Stat(keyAbs); err != nil {
		return "", "", fmt.Errorf("private key file: %w", err)
	}

	return certAbs, keyAbs, nil
}
