# Guide for AI agents working on me

This document summarizes repository structure and naming conventions so changes stay consistent as the `me` CLI evolves.

## What me is

`me` is a command-line application written in Go. The domain behavior is intentionally minimal at bootstrap time; contributors should keep architecture clean and incremental as features are introduced.

## Layout

| Path | Role |
| --- | --- |
| `cmd/me` | CLI entrypoint and argument handling. |
| `cmd/me/version` | Version/build metadata populated via ldflags. |
| `pkg/` | Library and domain packages used by `cmd`. |
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
