# me

`me` is a Go-based CLI application scaffolded for fast feature development, release automation, and agent-assisted workflows.

---

## Using the CLI

### Install or download

**Download prebuilt binaries:**

[GitHub Releases](https://github.com/olian04/go-me/releases) — the release workflow publishes cross-compiled binaries for Linux, macOS, and Windows (amd64 and arm64) plus checksums.

**Install with the Go toolchain:**

```bash
go install github.com/Olian04/go-me/cmd/me@latest
```

**Build from source:**

```bash
make build
# or
go build -o me ./cmd/me
```

### Running the CLI

```bash
./me version
```

---

## Development

| Command              | Purpose                                                 |
| -------------------- | ------------------------------------------------------- |
| `make build`         | Build `./me` from source                                |
| `make test`          | Run unit tests                                          |
| `make lint`          | Run vet, module verify, vuln scan, gosec, golangci-lint |
| `make format`        | Run go formatters                                       |
| `make build-release` | Local GoReleaser snapshot build into `dist/`            |

Contributor and agent-oriented notes on layout and conventions: `docs/AGENTS.md`.
