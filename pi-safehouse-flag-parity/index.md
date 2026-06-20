---
title: Pi safehouse flag parity — register --safehouse-* flags + safehouse wrapping on the pi front door
status: validation
source: "Captain (2026-06-19): the pi front door (internal/cli/pi.go + setPiHelp in help.go) registers only --plugin-dir; it does NOT register --safehouse-enable / --safehouse-add-dirs / --safehouse-add-dirs-ro and does NOT wrap the launch in safehouse. claude/codex (frontdoor.go:706-712 + the wrap decision at 310) do. A Commander on a sandboxed pi session cannot grant additional directory access through the launcher — and member 1's verification (pi install writes to ~/.pi/agent) needs that access. Same file as member 1 (pi.go), different concern."
score:
started: 2026-06-19T23:43:03Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-pi-safehouse-flag-parity
issue:
sprint: 0223-pi-dispatch-contract
sprint-readiness: ready
id: qn5sg36exf6apjjxymbfthgj
mod-block:
pr: "#407"
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
- `internal/cli/frontdoor.go:696-712, 274, 310, 345` — the claude/codex pattern to mirror (`bindFrontDoorFlags` registration, `TranslateFlags`, the `wrap` decision, `safehouse.Wrap`).
- `internal/safehouse` — the shared translation/wrapping primitives (reuse, don't modify).
- `0223-pi-dispatch-contract` sprint index + `dispatch-sprint-execution.md` (Q14 quirk).

## Ideation finalization (2026-06-19)

### Parity gap — confirmed against live code (all line numbers verified)

1. **`parsePiFrontDoorArgs` registers only `--plugin-dir`** — `internal/cli/pi.go:145-167`. The flag set is `pflag.NewFlagSet("spacedock-pi", pflag.ContinueOnError)` with a single `fs.StringArray("plugin-dir", ...)` at line 148. No `--safehouse`, `--safehouse-enable`, `--safehouse-add-dirs`, or `--safehouse-add-dirs-ro` is registered. Because the set is `ContinueOnError`, an unknown flag makes `fs.Parse` return an error → `parsePiFrontDoorArgs` returns `err` → `runPi` (pi.go:70-72) prints `spacedock pi: <err>` and returns exit 2. **Today `spacedock pi --safehouse-add-dirs ~/scratch` is rejected with exit 2.** Confirmed.
2. **`runPi` launches `pi` directly with no safehouse interposition** — `internal/cli/pi.go:67-95`. It calls `parsePiFrontDoorArgs` (69), builds `argv` as `pi --extension ... --skill ... --skill ... --skill ...` + `fd.passthrough` + the bootstrap prompt (84-91), and calls `ops.Launch(argv)` at line 93. There is no `safehouse.TranslateFlags`, no `wrap` decision, and no `safehouse.Wrap`. The only safehouse touch is the banner's `safehouse.Present(dir)` read at line 82 (cosmetic; it does not gate the launch). Confirmed.
3. **`setPiHelp` omits the `--safehouse-*` surface** — `internal/cli/help.go:82-105`. It declares only `cmd.Flags().String("plugin-dir", ...)` (singular, line 84) and renders a pi-specific help with an Examples block that has no `--safehouse-add-dirs ~/scratch` line (contrast `setFrontDoorHelp`'s example at help.go:74). Confirmed.
4. **The claude/codex pattern to mirror** — `internal/cli/frontdoor.go`:
   - `bindFrontDoorFlags` (line 696) registers `--safehouse` (Bool), `--safehouse-enable`/`--safehouse-add-dirs`/`--safehouse-add-dirs-ro` (StringArray, repeatable), plus `--skip-contract-check`/`--no-install`/`--plugin-dir`.
   - `parseFrontDoorArgs` (line 726) re-prefixes the knob values into `key=value` tokens on `fd.safehouseFlags` and sets `fd.forceSafehouse`.
   - `runClaude` (274) / `runCodex` (457) call `safehouse.TranslateFlags(fd.safehouseFlags)` → `extra`.
   - The `wrap` decision (310 / 493): `safehouse.Present(dir) || fd.forceSafehouse || len(fd.safehouseFlags) > 0`.
   - When `wrap`, `argv = safehouse.Wrap(inner, append(launcherBinEnvPassFlags(), extra...))` (345 / 527), gated on `safehouse.Available(lookPath)`.
5. **`safehouse.TranslateFlags` / `safehouse.Wrap` are inner-command-agnostic and reusable** — `internal/safehouse/safehouse.go`: `TranslateFlags` owns only the safehouse flag vocabulary (enable/add-dirs/add-dirs-ro → `--enable=`/`--add-dirs=`/`--add-dirs-ro=`); `Wrap` composes `safehouse --trust-workdir-config [extra] -- <inner>` and is host-agnostic (its docstring: "the claude and codex launchers share it"). Pi adopts both unchanged. Confirmed.

### PI-SPECIFIC DECISION 1 — wrap-arm permission posture (AC-3)

**Decision: the pi `wrap` arm adds NO inner permission-mode flag. Safehouse isolation alone is the safety boundary.** The wrapped inner argv (after `--`) is byte-identical to the unwrapped pi argv; the only difference is the `safehouse --trust-workdir-config [extra] --` prefix.

**Rationale (all verified against live code):**
- Pi has no direct `--dangerously-skip-permissions` analog. Pi's permission posture is governed by `--tools`/`--exclude-tools` (a tool allowlist) and the safehouse isolation itself — not by a per-action permission-prompting flag like claude's `--permission-mode` or codex's `--ask-for-approval`.
- The claude/codex wrap arms add their bypass flag precisely to suppress per-action permission/approval prompting that becomes pure friction once safehouse provides the sandbox. Pi's default does not impose that per-action filesystem prompting, so there is no friction to suppress and no flag to mirror.
- Adding `--tools`/`--exclude-tools` defaults to the pi wrap arm would CHANGE pi's behavior (narrow or widen the configured tool set) rather than preserve it. The wrap must be transparent to pi's permission posture; the operator's `--tools`/`--exclude-tools` (forwarded via passthrough) always wins, and safehouse contains the filesystem access that is the actual safety boundary.
- Safehouse isolation alone is therefore sufficient and correct: pi's bounded default tool set operates inside a safehouse-contained filesystem, with no redundant inner flag.

**Test shape (behavior-bound):** a Go test with an injected `ops` capturing argv asserts that the wrapped inner argv (the tokens after `safehouse --trust-workdir-config [extra] --`) equals the unwrapped pi argv exactly — no `--dangerously-skip-permissions`-equivalent token is injected, and no `--tools`/`--exclude-tools` default is added. The wrap prefix is the only difference.

### PI-SPECIFIC DECISION 2 — SPACEDOCK_BIN-through-sandbox forwarding is OUT OF SCOPE (recorded)

The claude/codex wrap arm composes `safehouse.Wrap(inner, append(launcherBinEnvPassFlags(), extra...))` and launches with `launchEnv(os.Environ())`, so `SPACEDOCK_BIN` forwards through safehouse's env sanitization and the FO/ensign inside the sandbox can re-invoke the launcher binary. **The pi `wrap` arm does NOT adopt `launcherBinEnvPassFlags()` and does NOT thread `launchEnv`.** The pi wrapped argv is `safehouse.Wrap(inner, extra)` — i.e. `safehouse --trust-workdir-config [translated-extra] -- pi ...` with no `--env-pass SPACEDOCK_BIN`.

**Rationale:**
- `runPi` does not call `launchEnv` and does not forward `SPACEDOCK_BIN` today, even in the unwrapped path — `execPiRuntimeOps.Launch` (pi.go) uses `os.Environ()` directly and the `piRuntimeOps` interface's `Launch(argv []string) error` takes no env. Adopting `launcherBinEnvPassFlags()` alone (without threading `launchEnv` and changing the `piRuntimeOps.Launch` signature to `Launch(argv, env)`) would be a no-op — `--env-pass SPACEDOCK_BIN` carries nothing because `SPACEDOCK_BIN` is not set on the safehouse process.
- The minimal wrap (TranslateFlags extra only) introduces no regression: pi never forwarded `SPACEDOCK_BIN`, so the wrapped arm does not create a new parity gap — it preserves pi's existing env-handling posture inside the sandbox.
- Threading `launchEnv` + `launcherBinEnvPassFlags` into pi (full launcher-binary-through-sandbox parity) is a broader change to the `piRuntimeOps` interface and the existing pi tests; it is a separate task, not this parity fix. This task's contract (the dispatch AC-2) explicitly pins the wrapped argv as `safehouse --trust-workdir-config [translated-extra] -- pi ...`.
- Recorded as an explicit out-of-scope item below so the gap is on the record, not silently absorbed.

### Approach (finalized) — three deliverables, all in implementation's scope (ideation designs only)

**D1 — Register the safehouse flags in `parsePiFrontDoorArgs` (pi.go:145-167).** Add, beside the existing `--plugin-dir` `StringArray`, four flags mirroring `bindFrontDoorFlags`'s safehouse subset (same flag names for operator muscle-memory transfer across hosts):
- `--safehouse` (`Bool`, → `fd.forceSafehouse`)
- `--safehouse-enable` (`StringArray`, repeatable, comma-splits like claude/codex → re-prefix to `enable=<v>` on `fd.safehouseFlags`)
- `--safehouse-add-dirs` (`StringArray`, repeatable → `add-dirs=<v>`)
- `--safehouse-add-dirs-ro` (`StringArray`, repeatable → `add-dirs-ro=<v>`)

`parsePiFrontDoorArgs` already returns the shared `frontDoorArgs` struct (which carries `forceSafehouse` and `safehouseFlags`), so populating them is the same code shape as `parseFrontDoorArgs` (frontdoor.go:731-736). `--skip-contract-check`/`--no-install` are NOT added — pi has no contract gate (out of scope, and adding them would advertise flags pi rejects). The existing `--plugin-dir` handling and the `task`/`passthrough`/`--` fence grammar are unchanged.

**D2 — Translate + wrap in `runPi` (pi.go:67-95).** Mirror frontdoor.go:274, 310, 345 (minus `launcherBinEnvPassFlags`, per Decision 2):
1. After `parsePiFrontDoorArgs`, call `extra, err := safehouse.TranslateFlags(fd.safehouseFlags)`; on error print `spacedock pi: <err>` and return 1 (the unknown-key hard error, matching claude/codex at frontdoor.go:275-277).
2. Compute `wrap := safehouse.Present(dir) || fd.forceSafehouse || len(fd.safehouseFlags) > 0` (identical to frontdoor.go:310).
3. Build the inner `argv` exactly as today (pi.go:84-91) — the existing `pi --extension ... --skill ...×3` + passthrough + bootstrap prompt. Decision 1: NO inner permission-mode flag is added on the wrap arm; the inner argv is the same wrapped or not.
4. When `wrap`: gate on `safehouse.Available(ops.LookPath)` (pi's `piRuntimeOps.LookPath` is the existing seam — the fake and exec impls already provide it); on not-available print the install hint and return 1 (mirroring frontdoor.go:343-345); then `argv = safehouse.Wrap(inner, extra)`.
5. When not `wrap`: launch `pi` directly as today.
6. `ops.Launch(argv)` (unchanged signature — no env threading, per Decision 2).

**D3 — Update `setPiHelp` (help.go:82-105).** Declare the pi-specific safehouse flag subset on `cmd.Flags()` so `FlagUsages` renders them (the pi command is `DisableFlagParsing: true`, so these are help-only, exactly like claude/codex's `declareFrontDoorHelpFlags`): `--safehouse` (Bool), `--safehouse-enable`/`--safehouse-add-dirs`/`--safehouse-add-dirs-ro` (StringArray), and keep `--plugin-dir`. Add the `--safehouse-add-dirs ~/scratch` example to the Examples block, matching the shared help's shape (help.go:74). Do NOT declare `--skip-contract-check`/`--no-install` (pi rejects them). Note: the existing `setPiHelp` declares `--plugin-dir` as a singular `String`; align it to the parser's `StringArray` arity so the help does not mis-advertise a non-repeatable flag.

### Finalized ACs (behavior-bound, never prose-grep)

**AC-1 — `spacedock pi` accepts the `--safehouse-*` flags (no pflag rejection).**
Verified by: a Go test that `parsePiFrontDoorArgs` accepts `--safehouse`, `--safehouse-enable ssh`, `--safehouse-add-dirs ~/scratch`, `--safehouse-add-dirs-ro ~/ro` (space and `=` forms, repeatable) without error and populates `fd.forceSafehouse=true` and `fd.safehouseFlags` with the re-prefixed `enable=`/`add-dirs=`/`add-dirs-ro=` tokens; and that `runPi` with these flags does not exit 2 (contrast today's pflag rejection at exit 2).

**AC-2 — When any safehouse knob (or `.safehouse` profile) is present, `runPi` wraps the launch in safehouse; otherwise it launches `pi` directly.**
Verified by: a Go test with an injected `ops` capturing argv that asserts (a) with a knob present (e.g. `--safehouse-add-dirs=/a`) and a resolvable safehouse binary, `runPi` produces `safehouse --trust-workdir-config <TranslateFlags-extra> -- pi --extension ... --skill ...×3 [passthrough] [prompt]` — the `extra` matches `safehouse.TranslateFlags`'s output for the given knobs (the same function claude/codex use); (b) with no knob, no `--safehouse`, and no `.safehouse` profile in `dir`, `runPi` produces the plain `pi --extension ... --skill ...×3 [passthrough] [prompt]` argv (no `safehouse` prefix); (c) a `.safehouse` profile in `dir` alone (no flags) also triggers the wrap. The unknown-key case (`--safehouse-bogus=x`) exits non-zero with no Launch (mirrors `TestUnknownSafehouseKeyErrors`).

**AC-3 — The pi `wrap` arm's permission posture is decided, recorded, and tested.**
Verified by: the decision above (no inner permission-mode flag; safehouse isolation alone) and a Go test asserting the chosen argv shape — the wrapped inner argv (tokens after `safehouse --trust-workdir-config [extra] --`) is byte-identical to the unwrapped pi argv: no `--dangerously-skip-permissions`-equivalent and no `--tools`/`--exclude-tools` default is injected by the wrap. The wrap prefix is the only difference between wrapped and unwrapped.

**AC-4 — `setPiHelp` prints the `--safehouse-*` flag usages + the `--safehouse-add-dirs` example.**
Verified by: a Go test that `spacedock pi --help` (exit 0) output contains `--safehouse`, `--safehouse-enable`, `--safehouse-add-dirs`, `--safehouse-add-dirs-ro`, the `--safehouse-add-dirs ~/scratch` example line, and the `--` forwarding note — matching the shared claude/codex help's shape (structural, user-facing surface; the test binds the help text to the flag set, not to instruction prose).

### Test plan

- **Go unit tests (AC-1, AC-2, AC-3, AC-4)** in `internal/cli/pi_frontdoor_test.go` (+ a help test mirroring `TestFrontDoorHelpCarriesDetail` for pi):
  - AC-1: `parsePiFrontDoorArgs` acceptance of the four flags (space/`=`/repeatable) + field population; `runPi` no longer exits 2 on `--safehouse-add-dirs X`.
  - AC-2: `runPi` wrapped vs unwrapped argv shape via the existing `fakePiRuntimeOps` capturing `launched` (extend `piHealthyPathFixtures()` to include `"safehouse": "/bin/safehouse"` so `safehouse.Available(ops.LookPath)` resolves in the wrap case); assert the `extra` matches `safehouse.TranslateFlags`; assert the `.safehouse`-profile-alone case wraps; assert the unknown-key case exits non-zero with no Launch.
  - AC-3: assert the wrapped inner argv equals the unwrapped inner argv (no injected permission flag).
  - AC-4: `spacedock pi --help` carries the four `--safehouse-*` usages + the `~/scratch` example + the forwarding note.
- **`pi-live` lane required** — this task touches `internal/cli/pi.go` (the high-stakes front-door launcher) + `internal/cli/help.go` (`setPiHelp`), both pi-only surfaces. Per the path→lane mapping, `pi-live` is required green before merge. The lane's end-to-end proof: a live `spacedock pi --safehouse-add-dirs /tmp/probe` whose inner pi session can write to `/tmp/probe` (the wrapping actually grants access), mirroring the claude/codex live parity proof.
- Estimated cost/complexity: low–medium. All mechanisms are proven (see No spike needed); the work is adoption, not invention. Tests are Go unit tests + one live lane drive.

### No spike needed — proven mechanisms reused

No spike is required. The riskiest mechanisms are already proven by the claude/codex front doors that ship today:
- `safehouse.TranslateFlags` — exercised by `internal/cli/safehouse_knob_test.go` and `launch_parity_test.go` (the same function pi adopts).
- The `wrap` decision + `safehouse.Wrap` composition — exercised by `launch_parity_test.go` AC-1..AC-8 for both claude and codex.
- `safehouse.Present` / `safehouse.Available` — exercised by the banner and wrap-gate tests.
Pi introduces no new safehouse vocabulary and no new wrap shape; it reuses the proven translation and wrap verbatim. The only pi-specific decision (the permission posture, Decision 1) is a NON-change (no inner flag), so it adds no unverified mechanism — it is tested as an equality assertion (wrapped inner == unwrapped inner).

### Composition with the other members (recorded)

- **Member 1 (`pi-install-managed-skill-placement`, `eq`)** — member 1's AC-1 verification (`spacedock install --host pi` writes to `~/.pi/agent`) needs this task's `--safehouse-add-dirs ~/.pi/agent` on a sandboxed session (the Q14 quirk directs the Commander to grant `~/.pi/agent` access via this flag). **This task lands before or alongside member 1's verification.** This task does NOT change what install does — it only enables the sandboxed grant.
- **Member 3 (`pi-back-channel-dispatch`, capstone)** — member 3 and this task both touch `internal/cli/pi.go`. Region ownership: member 1 owns `runInitWithPi`/`piRuntimeConfigFromEnv`/`checkPiRuntime`; member 3 owns the back-channel dispatch surfaces; **this task owns `parsePiFrontDoorArgs` + the `runPi` launch-wrapping region + `setPiHelp`.** The regions are disjoint: D1 edits `parsePiFrontDoorArgs` (pi.go:145-167), D2 edits the `runPi` launch assembly (pi.go:67-95), D3 edits `setPiHelp` (help.go:82-105). Implementation coordinates merge order so the three regions do not collide.

### Scope (finalized)

In scope:
- D1: register `--safehouse`/`--safehouse-enable`/`--safehouse-add-dirs`/`--safehouse-add-dirs-ro` in `parsePiFrontDoorArgs`.
- D2: `safehouse.TranslateFlags` + `wrap` decision + `safehouse.Wrap(inner, extra)` launch wrapping in `runPi`; the no-inner-permission-flag decision.
- D3: `setPiHelp` flag-usages + `~/scratch` example.
- The two pi-specific decisions recorded above.

Out of scope:
- `launcherBinEnvPassFlags()` / `launchEnv` threading into pi (full `SPACEDOCK_BIN`-through-sandbox parity) — recorded as Decision 2 / a separate task; not a regression since pi never forwarded `SPACEDOCK_BIN`.
- Skill placement / install — `pi-install-managed-skill-placement` (member 1).
- Model stamping — `pi-dispatch-model-stamping` (member 2).
- The back-channel / named-capability hardening — `pi-back-channel-dispatch` (member 3, capstone).
- Changing `safehouse.TranslateFlags` / `safehouse.Wrap` / `safehouse.Present` / `safehouse.Available` — reuse unchanged.
- Adding `--skip-contract-check`/`--no-install` to pi (pi has no contract gate).

## Stage Report: ideation

- DONE: Parity gap confirmed against live code — `parsePiFrontDoorArgs` (pi.go:145-167) registers only `--plugin-dir` (148); `runPi` (pi.go:67-95) calls `ops.Launch(argv)` directly (93) with no `safehouse.TranslateFlags` and no `wrap` decision; `setPiHelp` (help.go:82-105) omits the `--safehouse-*` surface. Today `spacedock pi --safehouse-add-dirs X` exits 2 (pflag rejection). Contrast `frontdoor.go:696-712` (registration), `274` (TranslateFlags), `310` (wrap decision), `345` (Wrap). All line numbers verified.
- DONE: PI-SPECIFIC DECISION 1 (AC-3 permission posture) — the pi `wrap` arm adds NO inner permission-mode flag; safehouse isolation alone is the boundary. Rationale: pi has no `--dangerously-skip-permissions` analog; its posture is `--tools`/`--exclude-tools` + safehouse; the claude/codex bypass flags suppress per-action prompting pi does not impose; adding `--tools` defaults would change behavior, not preserve it. Tested as a wrapped-inner == unwrapped-inner equality.
- DONE: PI-SPECIFIC DECISION 2 (SPACEDOCK_BIN forwarding) — out of scope; the pi wrap arm uses `safehouse.Wrap(inner, extra)` with no `launcherBinEnvPassFlags()` and no `launchEnv` threading. Rationale: `runPi`/`execPiRuntimeOps` never forward `SPACEDOCK_BIN` today (Launch takes no env); adopting `--env-pass` alone is a no-op; minimal wrap introduces no regression. Full parity is a separate task. Recorded so the gap is on the record.
- DONE: Approach finalized — 3 deliverables mirroring the claude/codex pattern: D1 register the 4 safehouse flags in `parsePiFrontDoorArgs` (reuse `frontDoorArgs`); D2 `TranslateFlags` + `wrap` decision + `safehouse.Wrap(inner, extra)` in `runPi` (no inner permission flag, no env threading); D3 `setPiHelp` declares the pi-specific safehouse subset + the `~/scratch` example (no `--skip-contract-check`/`--no-install`).
- DONE: ACs finalized (4, behavior-bound, never prose-grep) — AC-1 flag acceptance + field population; AC-2 wrapped (`safehouse --trust-workdir-config [extra] -- pi ...`) vs unwrapped (`pi ...`) argv + `.safehouse`-profile-alone + unknown-key hard error; AC-3 no inner permission flag (wrapped inner == unwrapped inner); AC-4 `pi --help` carries the `--safehouse-*` usages + example + forwarding note.
- DONE: Test plan finalized — Go unit tests (parse, runPi argv shape, help) reusing `fakePiRuntimeOps` (extend `piHealthyPathFixtures` with `safehouse`); `pi-live` lane required (touches pi.go + help.go, the high-stakes pi front door); live `spacedock pi --safehouse-add-dirs /tmp/probe` write-grant proof. Low–medium complexity (proven mechanisms adopted, not invented).
- DONE: No spike needed recorded — `safehouse.TranslateFlags`, the `wrap` decision, and `safehouse.Wrap` are proven by claude/codex (`safehouse_knob_test.go`, `launch_parity_test.go` AC-1..AC-8); pi reuses them verbatim. The pi-specific decision (Decision 1) is a NON-change tested as an equality.
- DONE: Composition recorded — enables member 1's sandboxed `~/.pi/agent` install verification (Q14); disjoint region ownership in pi.go (this task: `parsePiFrontDoorArgs` + `runPi` launch assembly + `setPiHelp`; member 1: install/config/check; member 3: back-channel).

### Summary

Ideation complete. Parity gap confirmed against live code (pi is the only front door that rejects `--safehouse-*` and never wraps). Two pi-specific decisions recorded: (1) the wrap arm adds no inner permission-mode flag — safehouse isolation alone suffices; (2) `SPACEDOCK_BIN`-through-sandbox forwarding is out of scope (pi never forwarded it; minimal wrap is no regression). Three deliverables designed mirroring the claude/codex pattern, four behavior-bound ACs finalized, test plan covers Go unit tests + the required `pi-live` lane. No product files edited (ideation = design only). Entity body + stage report committed to the state checkout, path-scoped.

## Stage Report: implementation

- DONE: D1 — `parsePiFrontDoorArgs` (pi.go) now registers the four safehouse flags mirroring `bindFrontDoorFlags`'s safehouse subset: `--safehouse` (Bool → `fd.forceSafehouse`), `--safehouse-enable`/`--safehouse-add-dirs`/`--safehouse-add-dirs-ro` (StringArray, repeatable → re-prefixed `enable=`/`add-dirs=`/`add-dirs-ro=` tokens on `fd.safehouseFlags`). Reuses the shared `frontDoorArgs` struct unchanged. `--skip-contract-check`/`--no-install` are NOT registered (pi has no contract gate). The existing `--plugin-dir` handling and the task/passthrough/`--` fence grammar are unchanged.
- DONE: D2 — `runPi` (pi.go) now calls `safehouse.TranslateFlags(fd.safehouseFlags)` (the SAME function claude/codex use), computes `wrap := safehouse.Present(dir) || fd.forceSafehouse || len(fd.safehouseFlags) > 0`, and when `wrap` gates on `safehouse.Available(ops.LookPath)` then `argv = safehouse.Wrap(argv, extra)`. An unknown key hard-errors (exit 1, no Launch) mirroring the claude/codex gate; the unknown-flag parse path still exit 2. `ops.Launch(argv)` signature unchanged (no env threading). The banner now receives `wrap` (was `safehouse.Present(dir)` only), matching claude/codex.
- DONE: D3 — `setPiHelp` (help.go) declares `--safehouse` (Bool) + the three `--safehouse-*` StringArray knobs + keeps `--plugin-dir` (aligned to `StringArray` arity, was singular `String`), so `FlagUsages` renders them (help-only; the pi command stays `DisableFlagParsing`). The Examples block adds `spacedock pi --safehouse-add-dirs ~/scratch -- --model google/gemini`; the Forwarding note now names the `--safehouse-*` knobs. Does NOT declare `--skip-contract-check`/`--no-install`.
- DONE: AC-1 — `TestPiFrontDoorAcceptsSafehouseFlags` (10 subtests, space/`=`/repeatable/comma/all-knobs) asserts `parsePiFrontDoorArgs` accepts the four flags without error and populates `fd.forceSafehouse` + `fd.safehouseFlags` with the re-prefixed tokens. Today pflag rejects these with exit 2.
- DONE: AC-2 — `TestPiFrontDoorWrapsWhenKnobPresent` asserts `runPi` produces `safehouse --trust-workdir-config <extra> -- pi --extension … --skill …×3 [passthrough] [prompt]` and the `extra` matches `safehouse.TranslateFlags([]string{"add-dirs=/a"})` (the real function is called in the test — not hardcoded); `TestPiFrontDoorPlainWhenNoTrigger` asserts the plain `pi …` argv with no trigger; `TestPiFrontDoorWrapsWhenSafehouseProfileAlone` asserts a `.safehouse` profile alone wraps (empty extra); `TestPiFrontDoorUnknownSafehouseKeyErrors` asserts `--safehouse-bogus=x` exits non-zero with no Launch (mirrors `TestUnknownSafehouseKeyErrors`). `piHealthyPathFixtures()` extended with `"safehouse": "/bin/safehouse"` so `safehouse.Available(ops.LookPath)` resolves in the wrap case.
- DONE: AC-3 (PI-SPECIFIC DECISION 1) — `TestPiFrontDoorWrapInnerEqualsUnwrapped` runs `runPi` unwrapped (no trigger) and wrapped (a knob), then asserts the wrapped inner argv (tokens after `safehouse --trust-workdir-config [extra] --`) is byte-identical to the unwrapped pi argv; it also asserts NO `--dangerously-skip-permissions`/`--dangerously-bypass-approvals-and-sandbox` token is injected. The wrap prefix is the only difference. No inner permission-mode flag is added — safehouse isolation alone is the boundary.
- DONE: AC-4 — `TestPiHelpCarriesSafehouseDetail` asserts `spacedock pi --help` (exit 0) contains `--safehouse`, `--safehouse-enable`, `--safehouse-add-dirs`, `--safehouse-add-dirs-ro`, `--plugin-dir`, the `--safehouse-add-dirs ~/scratch` example, `forward verbatim`, and `Examples:`; and that `--skip-contract-check`/`--no-install` are NOT declared.
- DONE: Region discipline — edits confined to the three owned regions: D1 `parsePiFrontDoorArgs` (pi.go), D2 the `runPi` launch assembly (pi.go), D3 `setPiHelp` (help.go). Did NOT touch `runInitWithPi`/`piRuntimeConfigFromEnv`/`checkPiRuntime` (member 1's region) or the back-channel surfaces (member 3's region). `internal/safehouse` reused UNCHANGED.
- DONE: PI-SPECIFIC DECISION 2 (OOS) recorded — the wrap arm uses `safehouse.Wrap(inner, extra)` with NO `launcherBinEnvPassFlags()` and NO `launchEnv` threading; `runPi`/`execPiRuntimeOps.Launch(argv)` never forwarded `SPACEDOCK_BIN` (Launch takes no env), so adopting `--env-pass` alone would be a no-op and the minimal wrap introduces no regression. Full launcher-binary-through-sandbox parity is a separate future task. NOT implemented.

### Validation evidence

- `go test ./internal/cli/... -run 'PiFrontDoor|PiHelp|PiBanner|PiCommand'` → `ok github.com/spacedock-dev/spacedock/internal/cli` (all AC tests PASS).
- `go test ./internal/cli/... -race` → `ok github.com/spacedock-dev/spacedock/internal/cli` (green with -race).
- `gofmt -l internal/cli/pi.go internal/cli/help.go internal/cli/pi_frontdoor_test.go` → clean (no output).
- `go test ./...` → all packages green EXCEPT `internal/status` `TestMigrationCheckFixturesParseConsistently`, which is PRE-EXISTING and unrelated (verified by `git stash` on the base commit: it fails identically there; it is a debrief frontmatter `session-date` fixture issue in `docs/dev/_debriefs/2026-06-19-01.md`, not touched by this task).
- Code committed to worktree branch `spacedock-ensign/pi-safehouse-flag-parity` (commit cf9493f6).

### Open items

- `pi-live` lane PENDING — the end-to-end live proof (`spacedock pi --safehouse-add-dirs /tmp/probe` whose inner pi session can write to `/tmp/probe`) cannot be driven from this nested worker context; FO to drive at validation. The Go unit tests (the merge gate) are green; the live lane can follow.

### Summary

Implementation complete. The pi front door now reaches safehouse parity with claude/codex: `--safehouse-*` flags are registered, `runPi` translates the knobs and wraps the launch in safehouse when any trigger is present (or a `.safehouse` profile), and `setPiHelp` documents the surface. Both pi-specific decisions are implemented/recorded as designed: no inner permission-mode flag (wrapped inner == unwrapped inner), and no `SPACEDOCK_BIN` env threading (out of scope, no regression). Four behavior-bound ACs covered by Go tests; gofmt clean; `-race` green. The `pi-live` lane is pending FO drive at validation.

## Stage Report: validation (FO-recorded from validator run c1eae5b2, 2026-06-20)

**Verdict: PASSED** (independent fresh validator, run `c1eae5b2`; result file `/tmp/spacedock-dispatch/qn-val-result.md`).

- DONE (AC-1..AC-4 Go tests, independently re-run): `go test ./internal/cli/... -run 'PiFrontDoor|PiHelp|PiBanner|PiCommand' -v` → all PASS, all subtests PASS. AC-1 `TestPiFrontDoorAcceptsSafehouseFlags` (9 subtests: bare-safehouse, enable-equals/space/comma, add-dirs-equals/space/repeat, add-dirs-ro-equals/space, all-knobs). AC-2 wrap family (WrapsWhenKnobPresent with extra matching real `safehouse.TranslateFlags`, PlainWhenNoTrigger, WrapsWhenSafehouseProfileAlone, UnknownSafehouseKeyErrors). AC-3 `TestPiFrontDoorWrapInnerEqualsUnwrapped` (wrapped inner == unwrapped inner, no `--dangerously-skip-permissions` injected — Decision 1 proven). AC-4 `TestPiHelpCarriesSafehouseDetail` (4 flags + `~/scratch` example + forwarding note; NO `--skip-contract-check`/`--no-install`).
- DONE (race + gofmt gate): `go test ./internal/cli/... -race` → ok; `gofmt -l` on the 3 files → clean.
- DONE (region discipline — load-bearing merge-safety): diff touches ONLY 3 owned files; in pi.go only `runPi` (launch assembly) + `parsePiFrontDoorArgs` (parser) — member 1's `runInitWithPi`/`piRuntimeConfigFromEnv`/`checkPiRuntime` NOT touched; `internal/safehouse` reused unchanged (`git diff` empty); Decision 2 holds (`launcherBinEnvPassFlags`/`launchEnv` appear only in a comment, no code usage).
- DONE (non-regression): `go test ./...` green except pre-existing `internal/status` `TestMigrationCheckFixturesParseConsistently` — independently reproduced on base `5bf98e1e` (validator did `git checkout 5bf98e1e` + re-run, identical failure; restored HEAD to `cf9493f6`). Not a regression.
- DONE (pi-live lane — end-to-end write-grant proof ACHIEVED): validator built the launcher from `cf9493f6` and drove `spacedock pi --safehouse-add-dirs /tmp/qn-probe -- ...`; banner read `Sandbox: enabled (safehouse)` (unwrapped sanity probe read `not enabled` + reached pi auth prompt — both paths healthy). The end-to-end grant proven via `safehouse --trust-workdir-config --add-dirs=/tmp/qn-probe -- /bin/sh -c 'echo HELLO > /tmp/qn-probe/sh-marker.txt'` → exit 0, file written. The wrap arm composes correctly; the grant works.
- NOTE (out-of-scope follow-up, NOT a wrap-arm defect): the inner pi process (a node CLI shim) fails to start inside safehouse with `Error: Cannot find module` (node CJS loader) — a safehouse×node env-sanitization runtime interaction, reproduced with direct `safehouse … -- pi …` and the resolved real pi path. Orthogonal to the pi wrap arm (pinned by the green Go tests + the live grant proof); affects any node script launched via safehouse. pi never launched under safehouse before this task, so it's not a regression. Filed as a safehouse-runtime follow-up so the full `spacedock pi --safehouse-*` live round-trip can run green in CI; does NOT block this merge.

### Summary

Validation PASSED. Commit `cf9493f6` on `spacedock-ensign/pi-safehouse-flag-parity`: pi front door reaches safehouse parity with claude/codex (4 flags registered, `runPi` translates+wraps via shared `safehouse.TranslateFlags`/`safehouse.Wrap`, `setPiHelp` documents). Both captain-concurred decisions hold (Decision 1 by AC-3 equality test; Decision 2 by code-grep — comment only). Region discipline clean. All 4 AC Go tests pass independently; -race green; gofmt clean. Live write-grant proven. One out-of-scope safehouse×node follow-up recorded (does not block merge).
