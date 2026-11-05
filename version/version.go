// Package version exposes build-time version information.
// The values are injected with -ldflags by the Makefile.
package version

import "fmt"

var (
	// Version is the release version or git describe output.
	Version = "dev"
	// Commit is the short git commit hash.
	Commit = "none"
	// Date is the UTC build timestamp.
	Date = "unknown"
)

// String returns a human-readable version line.
func String() string {
	return fmt.Sprintf("webdav-server %s (commit %s, built %s)", Version, Commit, Date)
}
