package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	// DefaultHost is used when no host is given.
	DefaultHost = "0.0.0.0"
	// DefaultPort is used when no port is given.
	DefaultPort = 8000
)

// ParseListenAddr normalises a listen value into a host:port string.
// Accepted forms: "host:port", "host" (port 8000), "port" (host 0.0.0.0).
// An empty value yields the default 0.0.0.0:8000.
func ParseListenAddr(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return net.JoinHostPort(DefaultHost, strconv.Itoa(DefaultPort)), nil
	}

	var host, port string
	switch {
	case net.ParseIP(value) != nil:
		// Bare IP address such as "192.168.1.10" or "::1".
		host, port = value, strconv.Itoa(DefaultPort)
	case strings.Contains(value, ":"):
		h, p, err := net.SplitHostPort(value)
		if err != nil {
			return "", fmt.Errorf("invalid listen address %q: %w", value, err)
		}
		host, port = h, p
	case isDigits(value):
		host, port = DefaultHost, value
	default:
		host, port = value, strconv.Itoa(DefaultPort)
	}

	if host == "" {
		host = DefaultHost
	}
	if err := validatePort(port); err != nil {
		return "", fmt.Errorf("invalid listen address %q: %w", value, err)
	}

	return net.JoinHostPort(host, port), nil
}

func validatePort(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 65535 {
		return fmt.Errorf("port %q must be a number between 1 and 65535", s)
	}
	return nil
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
