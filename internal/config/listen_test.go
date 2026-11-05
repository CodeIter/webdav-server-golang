package config

import "testing"

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
