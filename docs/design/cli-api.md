# me CLI API (Top-Down Design)

## Purpose

`me` provides a unified cross-platform interface to identify who the current user is, with consistent behavior across operating systems and distros. It is `whoami`-compatible where needed, but adds richer identity context for users, scripts, and debugging.

## Design Goals

- Predictable defaults for human usage.
- Stable machine-readable output for automation.
- Provider-driven architecture where each identity source is isolated in `pkg/`.
- Cross-platform behavior with explicit partial-failure semantics.
- Thin CLI layer (`urfave/cli/v3`) over provider and model packages.

## Command Model

### Default command: `me`

Default behavior is a compact, human-friendly identity summary.

Default summary (with the **default provider set**: all of `osaccount`, `envcontext`, `network`, `sysinfo`, `authproviders`) includes:

- Primary username (resolved from OS account source)
- UID/GID (or platform equivalent when available)
- Home directory
- Current host name (`meta.hostname` and/or network source)
- Context hints (sudo/ssh/ci) when present
- OS/architecture and friendly OS name/version (`sysinfo`)
- Git/cloud identity hints when present (`authproviders`)

This default is intentionally short and easy to scan.

### Subcommands

- `me whoami`
  - Strict GNU Coreutils `whoami` compatibility mode.
  - Prints effective username only.
  - Does not accept `me` identity/output flags.
- `me id`
  - Strict GNU Coreutils `id` compatibility mode.
  - Supports a documented v1 subset of GNU `id` flags.
  - Does not accept `me` identity/output flags.

## Global Flags

Global flags apply to `me` and all subcommands unless overridden.

- `--text`, `-t`
  - Explicit human-readable output mode.
- `--compact`, `-c`
  - Deterministic slash-separated fingerprint output mode.
- `--json`
  - Structured JSON output mode.
- `--yaml`
  - Structured YAML output mode.
- `--no-color`
  - Disable ANSI color in human output.
- `--timeout <duration>`
  - Provider deadline budget (e.g. `2s`, `500ms`).
  - Applies to aggregate provider execution.
- `--source <name>`
  - Selects which providers run, **in order**. When **omitted**, all providers run: `osaccount`, `envcontext`, `network`, `sysinfo`, `authproviders`.
  - When **any** `--source` is present, it **replaces** that full default list entirely (not additive).
  - Supports mixed input forms:
    - repeatable: `--source osaccount --source network`
    - comma-separated: `--source osaccount,network`
    - mixed: `--source osaccount,network --source authproviders`
- `--strict`
  - Strict validation and failure handling.
  - Unknown source names fail with exit `2`.
- Best-effort vs strict
  - Unless `--strict` is set, aggregate runs in best-effort mode (partial results, unknown sources recorded in `errors[]`).
- `--version`
  - Print version/build metadata and exit.
- `--help`
  - Print help and exit.

Output-mode selection:

- `--text`, `--compact`, `--json`, `--yaml` are mutually exclusive.
- If none is provided, default output mode is human text (equivalent to `--text`).

## Human Output Contract

Human output must use stable section ordering and labels so users can visually parse quickly and docs remain accurate.

### Default `me` human layout (concise)

Rendering uses label/value columns (same tabwriter style as `--version` text): aligned labels on the left, values on the right.

Order:

1. Core subject (always, missing values use `<unknown>`): `Username`, optional `Display name`, `UID`, `GID`, `Home`, `Shell`, `Hostname`
2. When the `sysinfo` source ran: `OS` (friendly name from `os_name` when present, otherwise `platform` / `GOOS`), optional `OS version`, `Architecture` (`arch` / `GOARCH`), optional `Platform` when both `os_name` and `platform` are set
3. Optional env context when present: `Sudo user`, `Sudo UID`, `SSH user`, `CI actor`, `CI provider`, or `CI` / `active` when CI is detected but actor/provider are empty
4. Optional network identity hints when present: `FQDN`, `Domain`, `Workgroup`
5. Optional auth identity hints when present: `Git user`, `Git email`, cloud rows (`AWS ARN`, `AWS account`, `GCP account`, `GCP project`, `Azure user`, `Azure tenant`, `Azure subscription`)
6. Optional `Warnings:` block (only when partial data/failures exist)

Timestamps, aggregate duration, and strict/best-effort flags are omitted from human text; use `--json` or `--yaml` for the full document.

Missing-value rules:

- Unknown field in human mode shows `<unknown>` (core subject rows only).
- Optional rows are omitted entirely when the underlying value is empty.
- Compact format uses empty segments instead of `<unknown>`.
- Empty values should not reorder or remove labels.

### Extended human detail behavior

When `--source` selects a subset of providers, human mode includes only those sources' detail sections while preserving the default base ordering for rows that apply.

Color rules:

- Color may emphasize status and section headers when enabled.
- `--no-color` disables all ANSI styling.
- Color is additive only; it must not carry meaning unavailable in plain text.

## Compact Format Contract

`compact` is a deterministic fingerprint-oriented output format returned as a single line string.

### Core rules

- Slash-separated fixed-slot sequence from broadest identifier to narrowest.
- Fixed slash count invariant across all runs and platforms.
- Slot order is stable and contract-bound in v1.
- Unresolved slot value is empty, but separators remain.
- No surrounding whitespace, labels, or commentary.

Example shape (illustrative only):

```text
<slot1>/<slot2>/<slot3>/<slot4>/<slot5>/<slot6>/<slot7>
```

If slots 3 and 6 are unresolved:

```text
<slot1>/<slot2>//<slot4>/<slot5>//<slot7>
```

### Determinism and normalization

- Normalize case to lowercase where semantics are case-insensitive.
- Trim leading/trailing whitespace.
- Replace internal `/` with escaped or normalized safe representation before join.
- Provider fallback precedence is fixed per slot and must not be dynamic.
- Same source data must produce same compact output.

### Error policy

- Best-effort mode: emit compact string with empty unresolved slots and report source issues through normal warnings/errors channels where applicable.
- Strict mode: any required slot resolution failure results in exit `3` (provider/slot failure), and compact output is not emitted as success output.

## Help UX Contract

Help behavior uses standard `urfave/cli/v3` output with minimal customization in v1.

- `me --help` shows:
  - command synopsis
  - global flags
  - `whoami` subcommand
- `me <command> --help` shows:
  - command synopsis
  - command-specific flags
  - command examples
- Invalid command/flag behavior:
  - print concise error message
  - include a hint to run `--help`
- No separate long-help mode in v1; default urfave help flow is the contract.

## Provider Boundaries

Each identity source is implemented as its own package:

- `pkg/identity/osaccount`
- `pkg/identity/envcontext`
- `pkg/identity/network`
- `pkg/identity/sysinfo` (GOOS/GOARCH, OS name/version)
- `pkg/identity/authproviders`

Shared packages:

- `pkg/identity/model`: canonical data structures
- `pkg/identity/provider`: provider interface and capability metadata
- `pkg/aggregate`: orchestration, timeout handling, result merge, and error policy

## Implementation Contract

All command paths must resolve identity data through `pkg/identity` packages.

- `cmd/me` command handlers are orchestration-only (parse flags, choose output mode, route execution).
- Identity acquisition must be delegated to `pkg/identity/*` and shared model/projection code.
- `me whoami` and `me id` are compatibility-format adapters over the same resolved identity model.
- Command handlers must not shell out to platform tools (`id`, `whoami`, etc.) as primary resolution logic.

### Anti-patterns (do not implement)

- Command-specific OS probing outside `pkg/identity`.
- Duplicated lookup logic across `me`, `me whoami`, and `me id`.
- Compatibility command data paths that bypass provider normalization in `pkg/identity`.

## Output and Exit Semantics

Canonical payload contains:

- `subject`: normalized primary identity
- `sources[]`: per-provider records with status
- `meta`: runtime metadata (platform, host, timestamp, duration)
- `errors[]`: non-fatal issues in best-effort mode

Exit codes:

- `0`: success (including partial success in best-effort)
- `2`: usage/flag/schema error
- `3`: strict-mode provider failure or required source missing
- `4`: internal/runtime failure

Compact output follows the same exit code model.

Source validation semantics:

- Best-effort (default): unknown `--source` values are ignored.
- Strict: unknown `--source` values are logged, then treated as usage errors (exit `2`) with early termination.
- Best-effort unknown-source events must be recorded in output `errors` and returned with normal output.

Compatibility command note:

- `me whoami` and `me id` follow GNU Coreutils semantics and are not affected by top-level `me` identity flags (`--source`, `--strict`, output modes, etc.).

## GNU Compatibility Reference

`me whoami` and `me id` use GNU Coreutils as the canonical compatibility baseline:

- [GNU Coreutils](https://www.gnu.org/software/coreutils/)

When platform-native behavior (macOS BSD tools, Windows equivalents) differs, `me` maps platform data to GNU-visible behavior as closely as possible. Deviations must be documented as compatibility gaps, not silent semantic changes.

## urfave/cli/v3 Mapping

- App bootstrap in `cmd/me` defines global flags and command tree.
- Each command action delegates to a use-case layer in `pkg/aggregate`.
- Error-to-exit-code mapping is centralized (single translator).
- Output rendering is centralized by mode (`text`, `compact`, `json`, `yaml`) and reused by commands.
- Help rendering stays on default urfave help printer path for v1.
- Compatibility command rendering (`whoami`, `id`) uses GNU-compatible formatters fed by `pkg/identity/model` data.

## Man Page Policy (v1)

Man pages are generated at build/release time from the command tree.

- One page per command:
  - `me.1`
  - `me-whoami.1`
- Generated man pages are build artifacts, not hand-authored source files.
- Source of truth remains:
  - CLI command definitions
  - design docs under `docs/design/`

## Version Module Parity Requirement

`cmd/me/version/version.go` must mimic bifrost version resolution behavior:

- Expose `Version`, `Revision`, `BuildTime` variables for ldflags injection.
- Expose `Info()` returning resolved values.
- Resolve values in priority order:
  1. Non-`unknown` ldflags values
  2. `runtime/debug.ReadBuildInfo()` keys:
     - `vcs.tag` for version
     - `vcs.revision` for revision
     - `vcs.time` for build time
  3. Module/dependency fallback when available
  4. `unknown`

## V1 Scope

In scope:

- Default summary output + command set (`whoami`, `id`)
- Deterministic compact fingerprint format via `--compact`
- Multi-source provider abstraction with partial failures
- Stable JSON/YAML contracts
- Structured exit code behavior
- Cross-platform best effort for Linux/macOS/Windows

Not in scope (v1):

- Remote identity providers requiring network auth
- Interactive login/session mutation flows
- TTY dashboards or watch mode
- Pluggable third-party provider SDK
- Policy engine / RBAC evaluation

## Compact Compatibility Policy

- Slot count and order are stable through v1.
- Slot semantics are stable through v1.
- New slots may only be introduced in a major version, unless appended under an explicit backwards-compatible policy that preserves prior slot meanings and separator invariants.
