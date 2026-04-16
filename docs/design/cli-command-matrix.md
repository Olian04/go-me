# me CLI Command Matrix (v1)

## Global CLI Syntax

```text
me [global flags] <command> [command flags] [args...]
```

## Global Flags

| Flag              | Type                         | Default                        | Applies To              | Behavior                                 |
| ----------------- | ---------------------------- | ------------------------------ | ----------------------- | ---------------------------------------- |
| `--text`, `-t`    | bool                         | implicit when no mode flag set | top-level `me`          | Human-readable text output               |
| `--compact`, `-c` | bool                         | `false`                        | top-level `me`          | Deterministic compact fingerprint output |
| `--json`          | bool                         | `false`                        | top-level `me`          | Structured JSON output                   |
| `--yaml`          | bool                         | `false`                        | top-level `me`          | Structured YAML output                   |
| `--no-color`      | bool                         | `false`                        | text output             | Disable ANSI colors                      |
| `--timeout`       | duration                     | `2s`                           | top-level `me`          | Provider deadline budget                 |
| `--source`        | repeatable/comma-list string | see provider table below       | top-level `me`          | Provider list; replaces defaults when set |
| `--strict`        | bool                         | `false`                        | top-level `me`          | Strict validation; unknown sources and provider errors fail the run |
| `--version`       | bool                         | `false`                        | top-level `me`          | Print version/build metadata and exit    |
| `--help`          | bool                         | `false`                        | top-level + subcommands | Print help and exit                      |

When `--strict` is **not** set, aggregate uses **best-effort** semantics (partial results, unknown `--source` names recorded in `errors[]`). There is no separate `--best-effort` flag.

## Commands

### Canonical Help Invocations

```bash
me --help
me whoami --help
me id --help
```

Help output is human-focused text regardless of output mode flags.

### `me` (default action)

- **Intent:** Show identity information for current runtime user/environment.
- **Provider usage:** when `--source` is omitted, all providers run (`osaccount`, `envcontext`, `network`, `sysinfo`, `authproviders`). If any `--source` is present, it **replaces** that full default (not additive).
- **Resolution path:** `pkg/aggregate` + `pkg/identity/model` (no direct platform command execution in handler).
- **Output modes:** text (default), compact, JSON, YAML.
- **Help behavior:** `me --help` prints synopsis and global flag behavior.

Examples:

```bash
me
me --text
me --compact
me --json
me --yaml
me --source osaccount --source network
me --source osaccount,network --source authproviders
me --source osaccount,envcontext,network,sysinfo,authproviders
me --strict
```

### `me whoami`

- **Intent:** Strict GNU Coreutils `whoami` compatibility mode.
- **Behavior:** prints effective username only.
- **Resolution path:** derive effective username from `pkg/identity/model`, then render GNU-compatible `whoami` output.
- **Flag compatibility:** does not accept identity output/reliability/source flags from top-level `me`.
- **Help behavior:** `me whoami --help` only.

Examples:

```bash
me whoami
```

### `me id`

- **Intent:** Strict GNU Coreutils `id` compatibility mode.
- **Behavior:** output and exit semantics mirror GNU `id` behavior.
- **Resolution path:** derive uid/gid/group/name data from `pkg/identity/model`, then render GNU-compatible `id` output and option behavior.
- **Flag compatibility:** accepts only GNU `id` options included in the documented v1 subset below.
- **Help behavior:** `me id --help` only.

v1 supported GNU-compatible subset:

- `-u`, `--user`
- `-g`, `--group`
- `-G`, `--groups`
- `-n`, `--name`
- `-r`, `--real`

Not in v1 subset (usage error exit `2`):

- GNU `id` flags outside the subset above.

Examples:

```bash
me id
me id -u
me id -g -n
me id -G
me id -u -r
```

## Provider Names (CLI Values)

**All providers (when `--source` omitted they all run), in order:**

| Name            | Role |
| --------------- | ---- |
| `osaccount`     | Local OS user account fields (username, ids, home, shell, groups). |
| `envcontext`    | Sudo, SSH, CI environment hints. |
| `network`       | Hostname, FQDN, domain/workgroup, local addresses. |
| `sysinfo`       | `GOOS`/`GOARCH`, OS name/version. |
| `authproviders` | Git and cloud auth hints. |

Unknown source handling:

- best-effort (default): continue execution, include unknown-source diagnostics in output `errors`
- strict: log unknown source and terminate early with usage error (exit `2`)

## Flag Interaction Rules

- Output mode flags are mutually exclusive:
  - `--text/-t`, `--compact/-c`, `--json`, `--yaml`
- If no output-mode flag is provided, behavior is equivalent to `--text`.
- Without `--strict`, runs are best-effort by default.
- `--source` supports mixed forms:
  - repeatable entries
  - comma-separated lists
  - mixed repeatable + comma-separated in same invocation
- `--no-color` applies to text output mode only.
- `me whoami` rejects top-level identity flags to preserve strict compatibility.
- `me id` rejects top-level identity flags to preserve strict compatibility.
- No command handler path should resolve identity data outside `pkg/identity/*`.

## Exit Codes

- `0`: success (includes partial result in best-effort mode)
- `2`: invalid flags/arguments, strict-mode unknown source, invalid combinations
- `3`: strict-mode provider failure or required source unavailable
- `4`: internal/runtime error

Compact mode uses the same exit code contract.

## Man Page Artifacts and Packaging Expectations

Generated man page artifacts (build/release):

- `me.1`
- `me-whoami.1`

Packaging/install expectation:

- Man pages are installed under `share/man/man1`.

Version alignment rule:

- Man page title/version metadata tracks the CLI release tag used for the build.
