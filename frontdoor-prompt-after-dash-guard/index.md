---
id: nz2aae5kfk1rb4e4z62vsnv1
title: Front-door silently swallows a positional prompt placed after `--` — bootstrap never prepends
status: ideation
source: captain (2026-06-03, session 12) — launched `spacedock claude --safehouse-enable=… --plugin-dir "$(pwd)" -- --model … --effort ultracode '@/tmp/handoff.md'` and the bootstrap/default prompt never prepended; the `@file` after `--` was treated as host passthrough so hasTask=false
started: 2026-06-04T06:30:55Z
completed:
verdict:
score: 0.28
worktree:
issue:
---

The front-door grammar (`internal/cli/frontdoor.go` `parseFrontDoorArgs`) splits args at `--`: tokens before `--` become the spacedock *task* (`fd.task`, `hasTask=true`), tokens after `--` forward verbatim to the host. The bootstrap prompt is only combined with the operator's prompt when the prompt is the task — `launchPrompt` returns `base + " " + task` only when `hasTask`. When the operator places their positional prompt AFTER `--` (intuitive, since "the prompt goes to claude"), it is classified as host passthrough, `hasTask` stays false, and the bare bootstrap is appended as a separate trailing positional after the operator's prompt. The host then receives two positionals and the spacedock bootstrap is effectively lost — the operator silently loses the launch-and-go preamble with no warning.

The grammar inversion (host flags after `--`, task before) was a deliberate fix for a value-taking host flag swallowing the prompt; this task does NOT propose reverting it. It proposes making the failure mode legible instead of silent.

## Problem

A positional (non-flag) token after `--` is almost always an operator who meant it as the launch prompt. Today it silently degrades to host passthrough, the bootstrap is misplaced, and there is no feedback. The captain hit this directly in session 12.

## Proposed approach

### Direction: warn-to-stderr, launch unchanged

When a bare positional appears after `--` and no task was given before `--`, emit a warning to stderr naming the corrected form (put the prompt before `--`), and launch UNCHANGED. We do NOT refuse.

Why warn, not refuse:
- The captain framed this explicitly as "make the silent failure LEGIBLE, not a behavior change to the grammar." A refuse halts the launch; a warn keeps the launch byte-identical and adds feedback.
- The guard's boundary is a *model* of the host arg grammars (claude/codex). That model is necessarily incomplete (host CLIs evolve). A warn that mis-fires on an unknown-but-legitimate host positional is a harmless extra stderr line; a refuse that mis-fires *blocks a launch the operator intended* — strictly worse. Choosing warn caps the blast radius of a false positive at noise, never at a denied launch.
- This mirrors the established `TestDanglingValueTakingHostFlagStillSwallows` philosophy: the front door preserves the host argv it is handed and surfaces consequences rather than silently re-routing.

The warning text names the corrected form so the operator can immediately fix it, e.g.:

> spacedock claude: warning: positional `@/tmp/handoff.md` after `--` is forwarded to the host as-is; the spacedock launch prompt was NOT prepended to it. To make it the launch prompt, put it BEFORE `--`: `spacedock claude '@/tmp/handoff.md' -- --model …`

### Trigger condition (when the guard fires)

The guard fires only when BOTH hold:
1. `fd.hasTask == false` — the operator gave no task before `--`. (If `hasTask == true`, a positional after `--` is a deliberate host positional; the operator already placed their prompt before `--`. No guard.)
2. The after-`--` passthrough contains at least one *stray positional* (defined below).

### Boundary: stray positional vs. legitimate host positional

The guard scans `fd.passthrough` token by token (the after-`--` slice). A token is a **stray positional** when it is non-flag (does not start with `-`, and is not the bare `--` separator) AND it is NOT accounted for by a recognized host grammar slot:

- **Value of a recognized value-taking host flag.** If the immediately preceding token is a host flag known to take a value (claude `-p`/`--print`, `--model`, `--mcp-config`, etc.; codex `-m`/`--model`, etc.), the current token is that flag's value — NOT stray. The `--flag=value` equals-form carries its value inside the flag token, so the following token is unaffected. We maintain a small per-host set of value-taking flags; a flag not in the set is treated as a boolean (its successor is then scanned independently). This errs toward NOT firing (a value-taking flag we don't know about means we might warn on its value — acceptable, it's only a warning).
- **Argument of a known host subcommand.** If the passthrough LEADS with a known host subcommand (codex `exec`, codex `resume`, claude has none relevant here), the positionals belonging to that subcommand are legitimate (`codex exec <prompt>`). The existing `codexResume` already encodes the "resume is a leading subcommand" fact; the guard extends the same idea to `exec`.

Concretely for the two named cases:
- `claude -p <prompt>` after `--`: `<prompt>` is the value of `-p` → NOT stray → no warning. ✅ negative case.
- `codex exec <prompt>` after `--`: passthrough leads with `exec` → subcommand args legitimate → no warning. ✅ negative case.
- `-- --model X '@/tmp/handoff.md'` (the captain's case): `X` is the value of `--model`; `@/tmp/handoff.md` is a bare positional with no preceding value-taking flag and no leading subcommand → STRAY → warning fires. ✅ positive case.

### Where it hooks into the existing parser

`parseFrontDoorArgs` (internal/cli/frontdoor.go) already computes `fd.passthrough` (the after-`--` slice) and `fd.hasTask`. The classification of after-`--` tokens today is: everything goes to `fd.passthrough` verbatim; nothing distinguishes a flag from a positional. The guard does NOT change that classification or the assembled argv — it is a pure read over the already-parsed `fd`.

Placement: a helper `strayPromptAfterDash(fd, host) (positional string, ok bool)` that returns the first stray positional. `runClaude`/`runCodex` call it after `parseFrontDoorArgs` succeeds and before assembling `inner`, and when `ok` they `fmt.Fprintf(stderr, warning)` — then proceed exactly as today. Host-specific value-taking-flag sets and subcommand names are passed in (claude vs codex) so the one helper serves both launchers, matching how `launchPrompt` is shared with a per-host `base`.

This keeps the guard a thin, side-effect-free-except-stderr layer: the launch argv is identical with or without the guard. That is the property the tests pin.

## Out of scope

- Reverting the host-flags-after-`--` grammar inversion.
- `$0` / binary-path propagation (that is `fc` launcher-binary-path-passthrough).
- Changing the assembled inner argv in any way. The guard only writes to stderr.

## Riskiest-unknown determination

**No spike needed: composes the existing parser.** The design relies only on already-proven mechanisms:
- `parseFrontDoorArgs` already produces `fd.passthrough` (the verbatim after-`--` slice) and `fd.hasTask` — pinned by `frontdoor_parse_test.go` (`host-flags-after-fence-no-task`, `task-then-fenced-host-flag`, etc.).
- The "leading subcommand" recognition is the same shape as the existing, proven `codexResume`.
- The guard reads `fd` and writes stderr; it does not parse host argv beyond a flat token scan over a slice the parser already produced. No round-trip, no on-disk format, no host-tool flag-support assumption is introduced.

The one modeled fact is the host value-taking-flag sets. These are NOT load-bearing for correctness of the launch (the argv is unchanged regardless); they only tune false-positive rate of an advisory warning. The negative-case test (`claude -p <prompt>`) pins the one set membership that matters for the named boundary. No end-to-end host launch is needed because the claim — "the stray prompt is named in a warning and the argv is unchanged" — is fully observable at the parse/launch-seam level the existing tests already exercise.

## Acceptance criteria

**AC-1: A bare positional after `--` (no task before `--`) is named in a stderr warning, and the assembled inner argv is unchanged.**
Verified by: a `runClaude` launch-seam test (style of `frontdoor_test.go`) with args `["--", "--model", "gpt-x", "@/tmp/handoff.md"]` and a `.safehouse`-free temp dir + compatible manifest, asserting (a) `stderr` contains the warning naming the stray positional `@/tmp/handoff.md` and the corrected form (prompt before `--`), and (b) `fake.launchedArg` equals `["claude", "--agent", "spacedock:first-officer", "--model", "gpt-x", "@/tmp/handoff.md", wantBootstrapPrompt]` — byte-identical to the pre-guard argv, proving warn-not-change. A second `runCodex` case proves the same for codex (subcommand-less passthrough leading with a value-taking flag).

**AC-2: A legitimate host positional after `--` does NOT trip the guard (no warning), and the argv is unchanged.**
Verified by: table-driven launch-seam test asserting `stderr` does NOT contain the warning marker for each negative case: (a) `["--", "-p", "do the thing"]` (claude `-p <prompt>` — value of a value-taking flag); (b) `runCodex` with `["--", "exec", "do the thing"]` (codex `exec <prompt>` — known subcommand argument); (c) `["task before", "--", "@/tmp/handoff.md"]` (a task WAS given before `--`, so `hasTask==true` — guard suppressed). Each also asserts `fake.launchedArg` is the expected unchanged argv. The presence-of-warning oracle is a single shared substring marker so positive and negative cases test the same signal.

**AC-3: The stray-positional classifier is unit-tested over the after-`--` token grammar independent of launch.**
Verified by: a table test (style of `frontdoor_parse_test.go`) over the `strayPromptAfterDash(fd, host)` helper asserting `(positional, ok)` for: stray after value-flag-pair, value-flag value (not stray), leading subcommand arg (not stray), `--flag=value` equals-form successor positional (stray — equals-form does not consume the next token), and the `hasTask==true` short-circuit (not stray). This pins the boundary logic directly, so a regression in flag-set membership or subcommand recognition fails a fast unit test, not only the launch-seam tests.

## Test plan

All `internal/cli` Go tests, no host CLI / network / fixtures beyond the existing `fakeHost`, `lookFound`, `compatibleManifest`, `wantBootstrapPrompt`, `wantCodexBootstrapPrompt` helpers. Cost: ~3 small test functions added to the existing `internal/cli` package; runs in the existing `go test ./internal/cli/...` (sub-second). No new fixture, no live workflow test — the claim is parse/launch-seam-level and the existing seam tests already prove the argv-assembly path. The warning text is asserted via a stable substring marker (e.g. the literal `after \`--\`` phrase plus the named positional) so the oracle survives benign rewording of the rest of the sentence while still failing if the warning is dropped.

## Stage Report: ideation

- DONE: Choose the direction for a bare positional appearing AFTER `--` when no task was given before `--`; state the chosen direction and why; define the non-false-positive boundary grounded in real host grammars; map how parseFrontDoorArgs classifies after-`--` tokens today and where the guard hooks in.
  Chose warn-to-stderr (launch unchanged) over refuse — false-positive blast radius is capped at noise, not a denied launch; documented in "Proposed approach". Boundary: stray = non-flag token NOT the value of a recognized value-taking host flag and NOT the arg of a known leading subcommand (claude `-p <prompt>`, codex `exec <prompt>` covered). Mapped to parseFrontDoorArgs producing `fd.passthrough`/`fd.hasTask`; guard is a pure read in runClaude/runCodex via `strayPromptAfterDash(fd, host)`.
- DONE: Write ACs each backed by a front-door parse/launch TEST asserting the warning-or-error AND the assembled inner argv, including a negative case; record the riskiest-unknown determination.
  AC-1 (positive: warning names stray positional + argv byte-identical via runClaude/runCodex launch-seam), AC-2 (negatives: `-p <prompt>`, `exec <prompt>`, hasTask short-circuit — no warning + unchanged argv), AC-3 (unit table over the classifier). Recorded "no spike needed: composes the existing parser" with the proven mechanisms (parseFrontDoorArgs slice/hasTask outputs, codexResume-shaped subcommand recognition).

### Summary

Designed an advisory stderr guard that makes the silent prompt-after-`--` swallow legible without altering the assembled host argv (the captain's session-12 case `-- --model … '@/tmp/handoff.md'` now warns). Key decision: warn rather than refuse, because the guard's host-grammar model is necessarily incomplete and a false-positive warning is harmless noise whereas a false-positive refusal blocks an intended launch. The boundary reuses the existing `codexResume`-style subcommand fact plus a small per-host value-taking-flag set; correctness of the launch never depends on that set (argv is unchanged regardless), so no spike was warranted — the parse/launch-seam tests prove the full claim.
