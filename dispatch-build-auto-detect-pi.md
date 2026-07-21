---
title: spacedock dispatch build auto-detects host pi when running under Pi agent harness
status: ideation
score: 0.85
id: 769mybp649pj160n17x13r8g
---

## Problem

When running `spacedock dispatch build` inside a Pi coding agent session without passing `--host pi` explicitly, the helper currently exits 1 with an error:

```
error: missing host source: pass --host, set JSON host, or run under CODEX_THREAD_ID or CLAUDECODE
```

Because Pi sets environment variables such as `PI_CODING_AGENT=true` (and `PI_CODING_AGENT_DIR`), `spacedock dispatch build` should automatically detect `host: "pi"` in the same way it detects Claude Code (`CLAUDECODE`) and Codex (`CODEX_THREAD_ID`). Requiring an explicit `--host pi` flag under Pi creates unnecessary First Officer runtime friction and breaks host-autodetect parity.

## Proposed Approach

1. **Auto-detect Pi Host (`internal/dispatch/` or `cmd/spacedock`):**
   Update the host resolution logic in `spacedock dispatch build` to inspect environment variables (`PI_CODING_AGENT` or `PI_CODING_AGENT_DIR`). When set (and `--host` is not explicitly passed), resolve `host` to `"pi"`.

2. **Unit Tests:**
   Add unit test cases covering Pi environment variable host auto-detection in `internal/dispatch` or `internal/status`.

## Acceptance Criteria

- **AC-1 (Auto-detect Pi host from environment):** Running `spacedock dispatch build` with `PI_CODING_AGENT=true` set in `env` (and no explicit `--host` flag) automatically resolves `host` to `"pi"`. *Verified by:* unit test.
- **AC-2 (Explicit --host overrides env):** When `--host claude` or `--host codex` is explicitly passed, it overrides the `PI_CODING_AGENT` environment detection. *Verified by:* unit test.

## Directives

- Perform ideation dispatch for this entity.

