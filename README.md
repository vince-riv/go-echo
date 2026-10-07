# go-echo

Minimal HTTP server. Responds with `text/plain` containing the full request: request line, headers, and body.

- Standard library only, runs in the foreground, no command line args
- Access log to stdout, one line per request
- Graceful shutdown on `SIGINT`/`SIGTERM`
- Multi-arch (`linux/amd64`, `linux/arm64`) `scratch` image

## Configuration

| Env var                | Default    | Description                                        |
|------------------------|------------|----------------------------------------------------|
| `GO_ECHO_PORT`         | `8080`     | Listen port                                        |
| `GO_ECHO_HEALTHZ_PATH` | `/healthz` | Path excluded from the access log (e.g. k8s probes) |

The health path still returns the echo response; it is just not logged.

## Run

```sh
go run .
# or
docker run -p 8080:8080 ghcr.io/vince-riv/go-echo:latest
```

```sh
curl -d hello -H 'X-Test: 1' localhost:8080/foo?a=b
```

## Access log format

```
<time> <remote addr> <method> <uri> <status> <bytes> <duration>
```

## Build

```sh
docker build -t go-echo .
```

Pushes to `main` build and publish `ghcr.io/vince-riv/go-echo:latest` via GitHub Actions.
