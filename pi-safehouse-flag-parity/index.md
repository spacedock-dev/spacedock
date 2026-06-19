---
title: Pi safehouse flag parity — register --safehouse-* flags + safehouse wrapping on the pi front door
status: ideation
source: "Captain (2026-06-19): the pi front door (internal/cli/pi.go + setPiHelp in help.go) registers only --plugin-dir; it does NOT register --safehouse-enable / --safehouse-add-dirs / --safehouse-add-dirs-ro and does NOT wrap the launch in safehouse. claude/codex (frontdoor.go:706-712 + the wrap decision at 310) do. A Commander on a sandboxed pi session cannot grant additional directory access through the launcher — and member 1's verification (pi install writes to ~/.pi/agent) needs that access. Same file as member 1 (pi.go), different concern."
score:
started: 2026-06-19T23:43:03Z
completed:
verdict:
worktree:
issue:
sprint: 0223-pi-dispatch-contract
sprint-readiness: ready
id: qn5sg36exf6apjjxymbfthgj
---

# Pi safehouse flag parity

## End value

`spacedock pi` accepts the same `--safehouse-*` flags as `spacedock claude`/`spacedock codex` (`--safehouse-enable`, `--safehouse-add-dirs`, `--safehouse-add-dirs-ro`, plus the bare `--safehouse` force) and wraps the pi launch in safehouse when any is given (or a `.safehouse` profile is present), mirroring the claude/codex front-door pattern. A Commander on a sandboxed pi session can grant additional directory access (e.g. `~/.pi/agent` for pi-install writes) through the launcher in one command, with no external `safehouse ... -- spacedock pi ...` wrapping and no claude/codex/pi parity gap.

## Problem — root cause already determined

The pi front door is deliberately a separate front door from the shared claude/codex one:

- `internal/cli/pi.go` `parsePiFrontDoorArgs` (line 145) registers **only** `--plugin-dir` (pi.go:148). It does **not** register `--safehouse-enable` / `--safehouse-add-dirs` / `--safehouse-add-dirs-ro` / bare `--safehouse`.
- `runPi` (pi.go:67) calls `ops.Launch(argv)` directly (pi.go:93) — no `safehouse.TranslateFlags`, no `wrap` decision. The `argv` is just `pi` + passthrough + bootstrap prompt.
- `setPiHelp` (`internal/cli/help.go`) prints a pi-specific help that omits the `--safehouse-add-dirs ~/scratch` example the shared `setHostHelp` prints for claude/codex (help.go:75).

Contrast the shared `frontdoor.go`:
- `parseFrontDoorArgs` (frontdoor.go ~700-712) registers `--safehouse-enable`, `--safehouse-add-dirs`, `--safehouse-add-dirs-ro`, `--plugin-dir`, `--skip-contract-check`, `--no-install`, plus bare `--safehouse`.
- `runClaude`/`runCodex` call `safehouse.TranslateFlags(fd.safehouseFlags)` (frontdoor.go:274) to translate the knobs into safehouse `extra` slot.
- The `wrap` decision (frontdoor.go:310): `safehouse.Present(dir) || fd.forceSafehouse || len(fd.safehouseFlags) > 0` — when true, the launch becomes `safehouse --trust-workdir-config [extra] -- claude --dangerously-skip-permissions ...`.

So pi is the only front door that cannot grant safehouse directory access through the launcher. `spacedock pi --safehouse-add-dirs ~/scratch` is rejected by pflag today.

## Why it matters for this sprint

Member 1 (`pi-install-managed-skill-placement`) D3 makes `spacedock install --host pi` run `pi install`, which writes to `~/.pi/agent/` (the pi package store + `settings.json packages`) — **outside** the repo. A sandboxed Commander session that cannot grant `--safehouse-add-dirs ~/.pi/agent` cannot verify member 1's AC-1 ("spacedock install --host pi actually installs"). The worktree stages themselves (`.worktrees/` under the repo) may work without this, but the install verification does not. This is a blocking dependency for member 1's verification on a sandboxed session.

It is also a general pi/claude/codex front-door parity gap independent of this sprint — pi is the only host where the launcher can't grant sandbox access.

## Approach (mirror the claude/codex pattern)

The fix is to bring `pi.go`'s front-door parsing + launch wrapping to parity with `frontdoor.go`:

1. **Register the safehouse flags** in `parsePiFrontDoorArgs` (pi.go:145): `--safehouse` (force bool), `--safehouse-enable` (StringArray), `--safehouse-add-dirs` (StringArray), `--safehouse-add-dirs-ro` (StringArray). Reuse the same flag names so operator muscle memory transfers across hosts.
2. **Translate + wrap** in `runPi` (pi.go:67): call `safehouse.TranslateFlags(fd.safehouseFlags)` → `extra`; compute `wrap := safehouse.Present(dir) || fd.forceSafehouse || len(fd.safehouseFlags) > 0`; when `wrap`, prepend `safehouse --trust-workdir-config [extra] --` to the `argv` (mirroring frontdoor.go:310+). When not `wrap`, launch `pi` directly as today.
3. **Update `setPiHelp`** to print the `--safehouse-*` flag usages + the `--safehouse-add-dirs ~/scratch` example, matching the shared help.

Pi-specific note: the claude/codex `wrap` arm adds `--dangerously-skip-permissions` to the inner command (the safehouse isolation is the safety boundary). Pi does not have a direct `--dangerously-skip-permissions` analog — pi's permission posture is governed by `--tools`/`--exclude-tools` and the safehouse isolation itself. Ideation confirms whether the pi `wrap` arm needs an inner permission-mode flag or whether safehouse isolation alone is sufficient (likely the latter — pi's default tool set is already bounded; safehouse contains the filesystem access). Record the decision.

## Scope

In scope:
- Register `--safehouse`/`--safehouse-enable`/`--safehouse-add-dirs`/`--safehouse-add-dirs-ro` on the pi front door.
- `safehouse.TranslateFlags` + `wrap` decision + `safehouse --trust-workdir-config [extra] -- pi ...` launch wrapping in `runPi`.
- `setPiHelp` flag-usages + example update.
- The pi `wrap`-arm permission-mode decision (recorded).

Out of scope:
- Skill placement / install — `pi-install-managed-skill-placement` (member 1). This task is the safehouse wrapping that *enables* member 1's verification on a sandboxed session; it does not change what install does.
- Model stamping — `pi-dispatch-model-stamping` (member 2).
- The back-channel / named-capability hardening — `pi-back-channel-dispatch` (member 3, capstone).
- Changing safehouse itself or `TranslateFlags` — reuse the existing translation; pi just adopts it.

## Composition with member 1 (load-bearing)

Member 1's AC-1 verification (`spacedock install --host pi` writes to `~/.pi/agent`) requires this task's `--safehouse-add-dirs ~/.pi/agent` on a sandboxed session. **This task lands before or alongside member 1's verification.** The Commander package's Q14 quirk directs the Commander to grant `~/.pi/agent` access for member 1's install verification via this task's flag.

## Acceptance criteria (provisional — ideation finalizes; proof = behavior, never prose-grep)

**AC-1 — `spacedock pi` accepts the `--safehouse-*` flags (no pflag rejection).**
Verified by: a Go test that `parsePiFrontDoorArgs` accepts `--safehouse`, `--safehouse-enable ssh`, `--safehouse-add-dirs ~/scratch`, `--safehouse-add-dirs-ro ~/ro` without error and populates the corresponding fields; and that `spacedock pi --safehouse-add-dirs X` does not reject (contrast today's pflag rejection).

**AC-2 — When any safehouse knob (or `.safehouse` profile) is present, `runPi` wraps the launch in safehouse.**
Verified by: a Go test (with an injected `ops` capturing the argv) that `runPi` produces `safehouse --trust-workdir-config [translated-extra] -- pi ...` when a knob is present, and plain `pi ...` when none is present and no `.safehouse` profile exists. The translated `extra` matches `safehouse.TranslateFlags`'s output for the given knobs (the same function claude/codex use).

**AC-3 — The pi `wrap` arm's permission posture is decided and recorded.**
Verified by: the ideation records whether the pi `wrap` arm adds an inner permission-mode flag or relies on safehouse isolation alone (with the rationale), and a Go test confirms the chosen behavior. (Behavior-bound: the test asserts the argv shape, not a prose claim.)

**AC-4 — `setPiHelp` prints the `--safehouse-*` flag usages + the `--safehouse-add-dirs` example.**
Verified by: a Go test that the pi help output contains the flag usages and the example line (matching the shared help's shape). Structural — the help is a user-facing surface; the test binds the help text to the flag set.

## Test plan

- Go unit tests (AC-1, AC-2, AC-4): `parsePiFrontDoorArgs` acceptance, `runPi` argv shape (wrapped vs unwrapped), help output.
- AC-3: the ideation decision + a Go test asserting the chosen argv shape.
- `pi-live` lane (this task touches `internal/cli/pi.go` + `help.go` — pi-only surfaces; per the path→lane mapping, `pi-live` required). A live `spacedock pi --safehouse-add-dirs /tmp/probe` that confirms the inner pi session can write to `/tmp/probe` (the end-to-end proof the wrapping grants access).

## Related

- `pi-install-managed-skill-placement` (`eqrcrxcyye56nfwm997bj33d`, member 1) — this task enables member 1's AC-1 verification on a sandboxed session (install writes to `~/.pi/agent`).
- `internal/cli/pi.go` (`parsePiFrontDoorArgs`, `runPi`) + `internal/cli/help.go` (`setPiHelp`) — the source of truth.
- `internal/cli/frontdoor.go:700-712, 274, 310` — the claude/codex pattern to mirror (`TranslateFlags`, `wrap` decision, `safehouse --trust-workdir-config [extra] -- ...`).
- `internal/safehouse` — the shared translation/wrapping primitives (reuse, don't modify).
- `0223-pi-dispatch-contract` sprint index + `dispatch-sprint-execution.md` (Q14 quirk).
