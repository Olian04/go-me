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

### Identity providers (`me` default command)

The default `me` command loads **identity providers** and merges their output into one document (text, compact, JSON, or YAML).

**When `--source` is omitted**, all of the following run, in order:

| Provider        | Purpose |
| --------------- | ------- |
| `osaccount`     | Local user account: username, uid/gid, home, shell, groups (where supported). |
| `envcontext`    | Environment hints: sudo user, SSH user, CI detection. |
| `network`       | Hostname, FQDN, domain/workgroup, local addresses. |
| `sysinfo`       | OS/runtime facts: `GOOS`/`GOARCH`, friendly OS name/version. |
| `authproviders` | Git user/email and best-effort cloud identity hints (AWS/GCP/Azure). |

**`--source` semantics:** If you pass **any** `--source`, that list **replaces** the full default set entirely—it is not additive. To run a subset only, name only those providers, for example:

```bash
me --source osaccount,sysinfo
```

Use `me --help` for the same provider list and behavior. Unknown provider names are reported in JSON/YAML `errors` in best-effort mode, or fail with exit code `2` when `--strict` is set.

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
