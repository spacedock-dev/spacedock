# Flag the CLAUDE_CODE_SUBPROCESS_ENV_SCRUB Permission-Mode Incompatibility

Status: approved design, ready for implementation.

Fixes: https://github.com/spacedock-dev/spacedock/issues/504

## Problem

Claude Code's subprocess credential-scrubbing hardening (`CLAUDE_CODE_SUBPROCESS_ENV_SCRUB=1`) forces a launched session's permission mode back to `default` unless the launcher declares `--allowedTools` explicitly. Spacedock launches Claude Code as a subprocess in both the unsandboxed (`--permission-mode auto`) and `--safehouse` (`--dangerously-skip-permissions`) forms, so a spacedock user who has this hardening set globally (a reasonable, recommended posture) gets a dispatched first officer silently downgraded to prompt-on-everything. The only warning today is Claude Code's own generic message, printed after spacedock's launch banner, with no acknowledgment from spacedock that this is a known incompatibility.

Spacedock should not try to silently route around the hardening by guessing an `--allowedTools` allowlist — that is a security-relevant decision that belongs to the operator. Spacedock's job is to flag the incompatibility clearly enough that the operator understands what happened and what their options are.

## Scope

Claude-only. `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB` is a Claude Code env var with no Codex or Pi analog.

## Design

### Detection helper

A new helper in `internal/cli/frontdoor.go`, alongside the existing `hasEnv`/`withoutEnv` env helpers:

```go
func subprocessEnvScrubActive(env []string) bool
```

Reports true when `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB` is present in `env` and set to a truthy value (non-empty, not `"0"`).

### 1. Launch-time warning (`spacedock claude`)

In `runClaude`, before the launch banner: if `subprocessEnvScrubActive(os.Environ())` is true and the operator has not already declared `--allowedTools`/`--allowed-tools` in `fd.passthrough` (checked the same way `passthroughHasFlag` checks other flags), print a spacedock-attributed stderr warning.

Fires on **both** the unsandboxed and `--safehouse` launch paths — we have no evidence `--dangerously-skip-permissions` is exempt from the hardening, so the warning errs toward showing rather than silently assuming safety.

The warning names the mechanism and the two remedies Claude Code's own message points at: declare `--allowedTools` explicitly, or unset `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB` for the launch. It is suppressed when the operator already declared `--allowedTools` themselves (the documented workaround), mirroring the existing "operator wins" suppression pattern used for `--permission-mode`.

### 2. `spacedock doctor` check

In `runDoctorWithPi`'s non-pi branch, when `host == "claude"`: after the existing manifest-compatibility report, print the same advisory note if `subprocessEnvScrubActive(env)` is true. This is informational only — it does not change doctor's exit code, which stays governed solely by manifest compatibility. It exists so the incompatibility surfaces during setup/CI, not only mid-launch.

## Testing

- `frontdoor_test.go`: table-driven cases on the warning — appears when the var is truthy and no `--allowedTools` given; suppressed when the var is unset, `"0"`, or `--allowedTools` is already present; appears on both the wrap and unwrap launch paths.
- Doctor test: the advisory note appears/is absent analogously, without disturbing the existing exit-code assertions for manifest verdicts.

## Out of scope

- Guessing or injecting an `--allowedTools` allowlist on the operator's behalf.
- Any change to Codex's launch path.
- Changing doctor's exit code based on this env var.
