# webdav-server

A small, dependency-light WebDAV server written in Go. It serves a local
directory over HTTP using the WebDAV protocol, so it can be mounted by Finder,
Windows Explorer, `davfs2`, `rclone`, and other WebDAV clients. It
supports optional HTTP basic authentication and graceful shutdown.

## Features

- Serves any existing local directory over WebDAV (PROPFIND, MKCOL, COPY, MOVE, LOCK, and related methods).
- Public mode by default, or password protection with a single user.
- Flexible listen address: `host:port`, `host` alone (port 8000), or `port` alone (host 0.0.0.0).
- Passwords are read from the `PASSWORD` environment variable, never from command-line arguments, so they do not appear in process listings.
- Constant-time credential comparison.
- Graceful shutdown on `SIGINT` or `SIGTERM`.
- Build-time version information injected with `-ldflags`.

## Requirements

- Go 1.27 or later
- `make` (optional, but recommended)
- `git` (optional, used for version metadata)

## Installation

### From source with Make

```bash
git clone https://github.com/CodeIter/webdav-server-golang.git
cd webdav-server
make build
```

The binary is written to `bin/webdav-server`. Copy it somewhere on your `PATH` if you want to use it globally:

```bash
sudo install -m 0755 bin/webdav-server /usr/local/bin/
```

### With the Go toolchain

```bash
go install github.com/example/webdav-server/cmd@latest
```

Note that `go install` does not inject version metadata, so `--version` reports `dev`.

## Usage

```
webdav-server [--help|-h] [--version|-v] [--listen|-l=host:port] [--serve|-s=DIR] [--username|-u=USERNAME]
```

| Flag | Description |
| --- | --- |
| `-h`, `--help` | Show help and exit. |
| `-v`, `--version` | Show version, commit, and build date, then exit. |
| `-l`, `--listen` | Address to listen on. Default: `0.0.0.0:8000`. |
| `-s`, `--serve` | Directory to serve. Required, and it must exist. |
| `-u`, `--username` | Enable basic auth for this user. The password is read from `$PASSWORD`. |

### Listen address formats

- `8080` listens on `0.0.0.0:8080`.
- `localhost` listens on `localhost:8000`.
- `192.168.1.10:9000` listens on that specific host and port.
- `::1` listens on `[::1]:8000`.

Ports must be between 1 and 65535.

### Examples

Serve the `./data` directory on port 8000 without authentication:

```bash
webdav-server -s ./data
```

Serve `/srv/files` on port 9000:

```bash
webdav-server --serve /srv/files --listen 9000
```

Require a login for user `alice`. The password comes from the environment:

```bash
PASSWORD='secret' webdav-server -s ./data -u alice -l 0.0.0.0:8443
```

Show the version:

```bash
webdav-server --version
```

### Environment

- `PASSWORD`: the password for the user given with `--username`. It is required only when `--username` is set. If `--username` is omitted, this variable is ignored.

### Exit codes

- `0`: success, including `--help` and `--version`.
- `1`: a runtime error, such as a failure to bind the port.
- `2`: invalid arguments or configuration.

## Connecting clients

Use the URL `http://<host>:<port>/` (or `https://` if you terminate TLS in front of the server).

- **macOS Finder:** Go > Connect to Server, then enter the URL.
- **Windows Explorer:** Map a network drive and enter the URL.
- **Linux:** Mount with `davfs2`, or open the URL in a file manager such as Nautilus or Dolphin.
- **rclone:** Create a remote with type `webdav` and vendor `other`.

When authentication is enabled, enter the username and the password from `$PASSWORD`.

## Security notes

- Without `--username`, anyone who can reach the port can read and write files in the served directory. Bind to `127.0.0.1` or a private interface if the server should not be public.
- Basic auth sends credentials on every request. Run the server behind HTTPS (for example, a reverse proxy such as Caddy or nginx with TLS) when it is reachable from untrusted networks.
- The server does not follow the served directory's boundaries beyond what `golang.org/x/net/webdav` enforces. Do not serve directories that contain sensitive files you do not want exposed.
- Locks are kept in memory, so they are lost when the server restarts.

## Development

The project uses the following layout:

```
.
├── cmd/            Entry point (main.go)
├── internal/
│   ├── app/        Wiring of configuration, logging, and server
│   ├── config/     Argument parsing, listen address parsing, usage text
│   └── server/     WebDAV handler construction and lifecycle
├── pkg/
│   └── auth/       Reusable HTTP basic auth middleware
├── version/        Build-time version variables
└── Makefile
```

Common Make targets:

| Target | Purpose |
| --- | --- |
| `make fmt` | Format code with `gofmt -s` and `go fmt`. |
| `make vet` | Run `go vet`. |
| `make test` | Run tests with the race detector. |
| `make build` | Build `bin/webdav-server` with version metadata. |
| `make run ARGS="..."` | Run from source with the given arguments. |
| `make clean` | Remove the `bin/` directory. |
| `make all` | Format, vet, test, and build. |

Run the server locally for development:

```bash
mkdir -p data
make run ARGS="-s ./data -l 8000"
```

## License

MIT License.

