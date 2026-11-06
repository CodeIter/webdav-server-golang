package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTLSConfig(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "cert.pem")
	keyPath := filepath.Join(dir, "key.pem")
	if err := os.WriteFile(certPath, []byte("CERT"), 0o600); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	if err := os.WriteFile(keyPath, []byte("KEY"), 0o600); err != nil {
		t.Fatalf("write key: %v", err)
	}

	cfg, err := Parse([]string{"--serve", dir, "--cert", certPath, "--key", keyPath}, func(string) string { return "" })
	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if cfg.CertFile != certPath || cfg.KeyFile != keyPath {
		t.Fatalf("Parse() TLS files = (%q, %q), want (%q, %q)", cfg.CertFile, cfg.KeyFile, certPath, keyPath)
	}

	_, err = Parse([]string{"--serve", dir, "--cert", certPath}, func(string) string { return "" })
	if err == nil {
		t.Fatal("Parse() with only cert should fail")
	}
}

func TestParseListenAddr(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "empty uses defaults", input: "", want: "0.0.0.0:8000"},
		{name: "port only", input: "9000", want: "0.0.0.0:9000"},
		{name: "host only", input: "localhost", want: "localhost:8000"},
		{name: "host and port", input: "127.0.0.1:8080", want: "127.0.0.1:8080"},
		{name: "bare IPv6", input: "::1", want: "[::1]:8000"},
		{name: "port out of range", input: "127.0.0.1:99999", wantErr: true},
		{name: "zero port", input: "0", wantErr: true},
		{name: "missing port after colon", input: "host:", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseListenAddr(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseListenAddr(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseListenAddr(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
