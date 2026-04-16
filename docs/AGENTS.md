# Guide for AI agents working on me

This document summarizes repository structure and naming conventions so changes stay consistent as the `me` CLI evolves.

## What me is

`me` is a command-line application written in Go. The domain behavior is intentionally minimal at bootstrap time; contributors should keep architecture clean and incremental as features are introduced.

## Layout

| Path | Role |
| --- | --- |
| `cmd/me` | CLI entrypoint and argument handling (`urfave/cli/v3`). |
| `cmd/me/version` | Version/build metadata (ldflags + `runtime/debug` fallback). |
| `cmd/me/render` | Human-readable text formatting for the default command. |
| `pkg/identity/model` | Canonical identity payload for JSON/YAML. |
| `pkg/identity/provider` | Provider interface for identity sources. |
| `pkg/aggregate` | Provider orchestration, timeouts, merge rules. Default run order (`DefaultSources`): `osaccount`, `envcontext`, `network`, `sysinfo` (`authproviders` is opt-in via `--source`). |
| `pkg/sysinfo` | Best-effort host OS name/version (`/etc/os-release`, `sw_vers`, etc.). |
| `pkg/identity/osaccount` | OS account / user lookup provider. |
| `pkg/identity/envcontext` | Environment context (sudo/ssh/CI hints). |
| `pkg/identity/network` | Hostname and local network hints. |
| `pkg/identity/sysinfo` | Identity provider: `GOOS`/`GOARCH`, OS name/version via `pkg/sysinfo` (default-on). |
| `pkg/identity/authproviders` | Best-effort git/cloud identity hints. |
| `pkg/compact` | v1 compact fingerprint slots. |
| `pkg/gnu` | GNU `whoami`/`id` text projections. |
| `pkg/` | Other library packages used by `cmd`. |
| `test/unit` | Unit tests by package focus. |
| `test/integration` | Integration/system tests when needed. |
| `.cursor/` | Agent rules, skills, and command docs for AI workflows. |
| `goreleaser/` | GoReleaser and container image packaging assets. |

## Naming and abstractions

- Prefer clear, action-oriented function names.
- Keep package names singular and lowercase.
- Avoid catch-all `util`/`helpers` packages.
- Wrap errors with context (`fmt.Errorf("...: %w", err)`).
- Define interfaces where they are consumed.

## What to read before large changes

1. `README.md`
2. This file (`docs/AGENTS.md`)
3. `makefile`
4. `.github/workflows/release.yml`
5. `goreleaser/.goreleaser.yaml`

---

_Maintainers: update this file when architecture or public behavior changes significantly._
