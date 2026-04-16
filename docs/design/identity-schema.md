# me Identity Schema (v1)

## Canonical Payload

This schema defines the contract for `me --json` and `me --yaml` outputs on the top-level command.

Schema ownership:

- Canonical identity structures are produced by `pkg/identity/model`.
- Command outputs (`--text`, `--compact`, `--json`, `--yaml`, GNU compatibility text for `whoami`/`id`) are projections/serializations of the same canonical model.
- Compatibility commands do not define separate identity acquisition paths.

```json
{
  "subject": {
    "username": "string",
    "display_name": "string",
    "uid": "string",
    "gid": "string",
    "home_dir": "string",
    "shell": "string"
  },
  "sources": [
    {
      "name": "osaccount|envcontext|network|authproviders",
      "status": "ok|partial|error|unavailable",
      "duration_ms": 0,
      "data": {},
      "warnings": ["string"]
    }
  ],
  "meta": {
    "platform": "string",
    "arch": "string",
    "hostname": "string",
    "timestamp": "RFC3339",
    "duration_ms": 0,
    "best_effort": true
  },
  "errors": [
    {
      "source": "string",
      "code": "string",
      "message": "string"
    }
  ]
}
```

## Field Definitions

### `subject`

Normalized primary identity. Fields may be empty when unavailable on a platform; absence is represented by empty string values in v1 for consistency.

### `sources[]`

Per-provider output envelope:

- `name`: provider identifier
- `status`:
  - `ok`: full provider success
  - `partial`: provider ran but some fields unavailable
  - `error`: provider failed
  - `unavailable`: provider unsupported on current platform/runtime
- `duration_ms`: provider runtime
- `data`: provider-specific payload
- `warnings`: non-fatal provider-level notes

### `meta`

Run-level metadata for diagnostics and observability.

### `errors[]`

Aggregated non-fatal errors (best-effort mode) or fatal context before strict-mode exits.

## Provider-Specific `data` Contracts

## `osaccount`

```json
{
  "username": "string",
  "uid": "string",
  "gid": "string",
  "home_dir": "string",
  "shell": "string",
  "groups": ["string"]
}
```

## `envcontext`

```json
{
  "sudo_user": "string",
  "sudo_uid": "string",
  "ssh_user": "string",
  "ci": {
    "is_ci": true,
    "provider": "string",
    "actor": "string"
  }
}
```

## `network`

```json
{
  "hostname": "string",
  "fqdn": "string",
  "domain": "string",
  "workgroup": "string",
  "local_addresses": ["string"]
}
```

## `authproviders`

```json
{
  "git": {
    "user_name": "string",
    "user_email": "string"
  },
  "cloud": {
    "aws": {
      "configured": true,
      "account_id": "string",
      "arn": "string"
    },
    "gcp": {
      "configured": true,
      "account": "string",
      "project": "string"
    },
    "azure": {
      "configured": true,
      "tenant_id": "string",
      "subscription_id": "string",
      "user": "string"
    }
  }
}
```

## Command-Specific Structured Responses

- `me --json` / `me --yaml`: canonical payload, with provider set controlled by `--source`.
- `me whoami`: strict compatibility command, plain username output only (outside this schema).
- `me --version`: separate small schema/value output path:
  - `version`
  - `revision`
  - `build_time`

## GNU `id` Compatibility Contract

`me id` is a GNU Coreutils compatibility command and is primarily text-output oriented, following GNU `id` behavior and flags. This section defines implementation-level required/optional data mapping for cross-platform support.

Normative reference:

- [GNU Coreutils](https://www.gnu.org/software/coreutils/)

Required compatibility fields (must be modeled; may be unresolved on some platforms):

- effective user id
- effective group id
- effective user name
- effective group name
- supplementary groups (id + name pairs where available)

Optional compatibility fields (omit when unavailable or not requested by active GNU-compatible flag subset):

- real user id / name
- real group id / name
- additional formatting variants not in v1 subset

Cross-platform fallback rule:

- If a required compatibility field cannot be resolved on a supported platform, preserve GNU-facing output semantics as closely as possible and record compatibility diagnostics for best-effort paths.
- Optional fields that cannot be resolved are omitted unless GNU semantics for the selected option require their presence.

## Compact Fingerprint Format

`--compact`/`-c` returns a deterministic single-line slash-separated sequence. It is not JSON/YAML, but this schema defines how compact slots derive from canonical fields.

### v1 Slot Model

v1 uses **7 fixed slots** (6 slash characters), broadest to narrowest:

1. `platform_scope` (OS + distro family or platform class)
2. `host_scope` (hostname/FQDN fallback)
3. `account_scope` (primary username)
4. `principal_scope` (UID/SID/platform principal id)
5. `group_scope` (primary gid/group id)
6. `context_scope` (sudo/ssh/ci actor hint)
7. `auth_scope` (best available external identity: git/cloud principal)

Compact string shape:

```text
slot1/slot2/slot3/slot4/slot5/slot6/slot7
```

If slot 2 and slot 6 are unresolved:

```text
slot1//slot3/slot4/slot5//slot7
```

### Slot Resolution and Fallback Policy

Each slot has deterministic source precedence:

- `platform_scope`: runtime platform metadata -> `network` domain/workgroup class hint
- `host_scope`: `network.fqdn` -> `network.hostname`
- `account_scope`: `osaccount.username` -> `subject.username`
- `principal_scope`: `osaccount.uid` (or platform SID equivalent)
- `group_scope`: `osaccount.gid`
- `context_scope`: `envcontext.sudo_user` -> `envcontext.ssh_user` -> `envcontext.ci.actor`
- `auth_scope`: cloud principal (`aws arn`/`gcp account`/`azure user`) -> `git.user_email` -> `git.user_name`

If all fallback candidates fail for a slot, that slot is empty.

### Population Quality Bar

Target for healthy systems: approximately **90% slot population**.

- Empty slots are valid and expected in some environments.
- Repeated empty slots on managed systems should be treated as potential hygiene/observability issues.
- Diagnostics should identify which source/slot failed to resolve and why.

### Diagnostics Guidance

When compact slots are empty, human/structured diagnostics should expose:

- affected slot name
- attempted providers in precedence order
- failure category (`unavailable`, `permission_denied`, `not_configured`, `lookup_failed`)

This guidance supports operational follow-up while preserving compact deterministic output.

## Help and Man Page Boundary

- `--help` output is documentation/help text and is not part of the JSON/YAML identity schema contract.
- Man pages are generated documentation artifacts and are also outside schema guarantees.
- When `--help` is invoked, schema payload output is bypassed.
- `--compact` also bypasses JSON/YAML schema output and emits compact text only.

## Unknown `--source` Error Semantics

Top-level `me` source filtering uses these rules:

- Best-effort: unknown source names are recorded in `errors[]` and returned with normal output.
- Strict: unknown source names are logged and terminate execution early with usage error.

## Compatibility and Stability

- v1 guarantees top-level keys: `subject`, `sources`, `meta`, `errors`.
- New optional nested fields may be added in minor releases.
- Existing keys and semantic meanings are stable through v1.
- Breaking schema changes require major version bump.
- Compact slot count/order/meaning are contract-stable through v1.
