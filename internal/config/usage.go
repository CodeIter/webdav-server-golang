package config

import (
	"fmt"
	"io"
)

const usageText = `Usage: webdav-server [--help|-h] [--version|-v] [--listen|-l=host:port] [--serve|-s=DIR] [--username|-u=USERNAME] [--cert|-c=PATH] [--key|-k=PATH]

Options:
  -h, --help        Show this help message and exit
  -v, --version     Show version information and exit
  -l, --listen      Address to listen on. Accepts host:port, host (port 8000),
                    or port (host 0.0.0.0). Default: 0.0.0.0:8000
  -s, --serve       Directory to serve. Required, must exist.
  -u, --username    Enable HTTP basic auth for this user. The password is read
                    from the PASSWORD environment variable. Without this flag
                    the server is public.
  -c, --cert        HTTPS certificate file path. Requires --key/-k too.
  -k, --key         HTTPS private key file path. Requires --cert/-c too.

Environment:
  PASSWORD          Password for --username. Required only when --username is set.

Examples:
  PASSWORD=secret webdav-server -s ./data -u alice
  webdav-server --serve /srv/files --listen 9000
  webdav-server -s . -l 192.168.1.10:8080
  webdav-server --serve /srv/files --cert /etc/ssl/server.crt --key /etc/ssl/server.key
`

// PrintUsage writes the help text to w.
func PrintUsage(w io.Writer) {
	fmt.Fprint(w, usageText)
}
