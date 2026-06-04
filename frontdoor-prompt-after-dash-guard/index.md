---
id: nz2aae5kfk1rb4e4z62vsnv1
title: Front-door silently swallows a positional prompt placed after `--` — bootstrap never prepends
status: validation
source: captain (2026-06-03, session 12) — launched `spacedock claude --safehouse-enable=… --plugin-dir "$(pwd)" -- --model … --effort ultracode '@/tmp/handoff.md'` and the bootstrap/default prompt never prepended; the `@file` after `--` was treated as host passthrough so hasTask=false
started: 2026-06-04T06:30:55Z
completed:
verdict:
score: 0.28
worktree: .worktrees/spacedock-ensign-frontdoor-prompt-after-dash-guard
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

## Stage Report: implementation

- DONE: Implement the advisory stderr guard per the approved ideation: a pure read `strayPromptAfterDash(fd, host)` invoked from runClaude/runCodex that WARNS to stderr (naming the corrected form: task before `--`) when a bare positional appears after `--` and no task was given before it — WITHOUT altering the assembled host argv. Boundary: do NOT warn on a legit host positional. TDD: write the failing test first.
  `strayPromptAfterDash` + `warnStrayPromptAfterDash` added to internal/cli/frontdoor.go; called after the gate in both runClaude/runCodex. Boundary skips a value-taking host flag's value (per-host `valueTakingHostFlags`) and a leading host subcommand's args (`leadingHostSubcommands`: codex exec/resume). Failing test (undefined symbol → build failed) confirmed before impl. Commit a354943f.
- DONE: ACs each proven by a front-door parse/launch test. AC-1 positive (warning names the stray positional AND launch argv byte-identical), AC-2 negatives (`-p <prompt>`, `exec <prompt>`, hasTask short-circuit → NO warning + unchanged argv), AC-3 unit table over the classifier. Full offline `go test ./...` green; `go vet`/`gofmt` clean.
  internal/cli/frontdoor_stray_prompt_test.go: AC-1 = TestClaudeStrayPromptAfterDashWarns + TestCodexStrayPromptAfterDashWarns; AC-2 = TestStrayPromptGuardNegatives (3 cases); AC-3 = TestStrayPromptAfterDashClassifier (5 cases). `go test ./...` = 1014 passed in 12 packages; `gofmt -l internal/cli/` empty; `go vet ./internal/cli/` no issues.

### Summary

Implemented the warn-not-change guard exactly as the ideation specified: a pure-read classifier plus a stderr emitter wired into both launchers after the contract gate, leaving the assembled host argv byte-identical (the AC-1 tests assert the full argv equals the pre-guard shape). The boundary reuses the established `codexResume` leading-subcommand idea (extended to `exec`) and a small per-host value-taking-flag set so `-p <prompt>` / `exec <prompt>` do not warn; membership is advisory-only since the argv never changes. TDD: the classifier symbol was undefined first (build failed), then the 12 new subtests pass and the full offline suite is green (1014/1014).

## Stage Report: validation

- DONE: Reproduce AC-1 (TestClaudeStrayPromptAfterDashWarns / TestCodexStrayPromptAfterDashWarns), AC-2 negatives (TestStrayPromptGuardNegatives — 3 cases), AC-3 (TestStrayPromptAfterDashClassifier — 5 cases). Confirm the guard is advisory-only and NEVER alters argv.
  All 12 subtests pass (`go test -run` on the four functions). Argv-unchanged invariant confirmed structurally (`warnStrayPromptAfterDash` takes `fd` by value, returns nothing, only writes stderr) AND by adversarial mutation: an injected `fd.passthrough = fd.passthrough[:len-1]` after the warn call FAILS AC-1's byte-identical argv assertion. The advisory itself, the value-flag exemption, and the leading-subcommand exemption each fail the suite when individually disabled (4 mutations, all caught) — the tests pin behavior, not spelling.
- FAILED: Confirm the boundary holds: the captain's session-12 case (`-- --model … '@/tmp/handoff.md'`) DOES warn, and legitimate host positionals do NOT. Full offline `go test ./...` green; `go vet`/`gofmt` clean.
  Full offline suite GREEN (12 packages ok), `gofmt -l internal/cli/` empty, `go vet ./internal/cli/` clean. BUT the captain's ACTUAL session-12 invocation — which included `--plugin-dir "$(pwd)"` (see this entity's `source:` line) — produces a WRONG warning. `parseFrontDoorArgs` re-injects `--plugin-dir <dir>` at the FRONT of `fd.passthrough`; the classifier then scans `--plugin-dir` (flag, skipped) → `<dir>` whose preceding token `--plugin-dir` is NOT in `valueTakingHostFlags["claude"]`, so `<dir>` is returned as the FIRST stray positional. The warning names `/tmp/co` (a spacedock-internal re-injection) instead of the operator's `@/tmp/handoff.md`, advises "put `/tmp/co` BEFORE `--`" (wrong — it is already a spacedock flag), and MASKS the real stray prompt the guard exists to surface. Reproduced via real `runClaude` and via `parseFrontDoorArgs`+`strayPromptAfterDash` directly; affects both `--plugin-dir P` and `--plugin-dir=P` (the equals-form is normalized to the space-pair in passthrough).

### Summary

REJECTED. Tests are strong (4 adversarial mutations all caught; argv-unchanged invariant holds in every exercised path) and the simple-`--` boundary cases are correct. The blocker is a material correctness hole the AC tests never exercised: the entity's own motivating case — the captain's session-12 `spacedock claude … --plugin-dir "$(pwd)" -- --model … '@/tmp/handoff.md'` — warns on the WRONG token. `--plugin-dir` is spacedock-owned and re-injected to the front of `fd.passthrough`, but is absent from the classifier's value-taking-flag set, so its value is misclassified as the stray positional and the operator's real after-`--` prompt is shadowed. Fix direction (for implementation): the classifier must account for the re-injected `--plugin-dir <dir>` pair (e.g. add `--plugin-dir` to the per-host value-taking set, or have the classifier skip the spacedock-injected prefix); AC-1 should grow a case carrying `--plugin-dir` before `--` so the regression is pinned. This is advisory-only (argv is still unchanged), so it is not a launch-safety bug — but it defeats the feature's stated purpose for the exact scenario it was built to address.

### Feedback Cycles

**Cycle 1 (2026-06-04) — validation REJECTED + detached audit, one converging material root cause.** Both the validation and the detached audit (`audit-frontdoor-prompt-after-dash-guard`, a354943f) confirmed the argv-unchanged invariant holds and the tests are strong (4 adversarial mutations caught), but `valueTakingHostFlags` is incomplete in two converging ways: (1) **the captain's actual session-12 case** carried `--plugin-dir "$(pwd)"`, which the front-door re-injects at the front of `fd.passthrough`; `--plugin-dir` is absent from the value-taking set, so the classifier names the DIR as the stray positional and shadows the real `@/tmp/handoff.md` — the guard misfires on the exact scenario it exists for; (2) common host value-taking flags (claude `--permission-mode`/`--add-dir`/`--append-system-prompt`/`--settings`/`--session-id`/`--output-format`; codex `--config`/`-c`/`--cd`/`--image`/`--sandbox`/`--profile`) are absent, so the guard warns on their VALUES with actively-wrong "put X before --" advice. Routed to implementation: account for the spacedock-injected `--plugin-dir <dir>` prefix (+ any re-prepended spacedock flags), broaden the value-taking sets to the real host flags, suppress the prescriptive advice when the prior token is an unrecognized `-`-flag, and add AC cases for the `--plugin-dir … -- … @file` session-12 shape + a common unrecognized value-flag negative.

**Cycle 2 (2026-06-04) — re-validation REJECTED, same root-cause family (injected `--plugin-dir` prefix), now the subcommand path.** Cycle 1 closed the session-12 value-flag path and the argv invariant holds, but the injected-prefix fix was only HALF done: it defeats the LEADING-SUBCOMMAND exemption. `spacedock codex --plugin-dir <dir> -- exec <prompt>` (and `resume <id>`) misfires — warns on `exec`/`resume` with wrong advice — because `parseFrontDoorArgs` re-prepends `--plugin-dir <dir>`, pushing the subcommand to index 2, while the exemption only checks `passthrough[0]`. **Structural fix (not another per-path patch):** skip the leading injected `--plugin-dir <dir>` pair ONCE, up front, so BOTH the value-flag scan and the subcommand check see the real after-`--` tokens. Add a codex `--plugin-dir … -- exec <prompt>` negative AC. Polish: `TestClaudeStrayPromptSession12Shape` does NOT red if `--plugin-dir` is dropped from `valueTakingHostFlags` (the fallback (c) independently suppresses the dir), so the structural skip is the load-bearing mechanism — add a pinning test or correct the cycle-1 report's item-(a) credit. (Cycle 3 escalates to the captain.)

## Stage Report: implementation (cycle 1)

- DONE: (a) Account for the spacedock-injected `--plugin-dir <dir>` prefix so it does not shadow the real prompt.
  Added `--plugin-dir` to both hosts' `valueTakingHostFlags`; parseFrontDoorArgs re-injects only `--plugin-dir` into fd.passthrough (safehouse knobs live in the separate fd.safehouseFlags field, booleans are consumed) so this is the complete injected-prefix set. Session-12 case now names `@/tmp/handoff.md`, never `/co`. Commit 571493c4. [CYCLE-2 CORRECTION: validation showed the SET-membership claim here was redundant — the conservative fallback (c), not the set addition, suppressed `/co` for the session-12 shape, AND this value-flag-only approach did NOT cover the leading-subcommand path. Cycle 2 replaces it with a structural `skipInjectedPrefix` and removes `--plugin-dir` from the value set; see Stage Report: implementation (cycle 2).]
- DONE: (b) Broaden valueTakingHostFlags to the real common host value-taking sets.
  claude += --permission-mode/--add-dir/--append-system-prompt/--settings/--session-id/--output-format; codex += --config/-c/--cd/--image/--sandbox/--profile.
- DONE: (c) Conservative fallback: do NOT emit prescriptive advice when the prior token is an unrecognized `-`-prefixed flag (likely its value).
  Classifier suppresses when the predecessor is a space-form `-`-flag not in the value-taking set; an equals-form `--flag=value` carries its own value so its successor stays confidently stray (this distinction fixed the equals-form classifier case that the first fallback cut broke).
- DONE: (d) AC cases for the session-12 shape (warns on @file, not the dir) + a common unrecognized value-flag negative; go test ./... + vet/gofmt green.
  Added TestClaudeStrayPromptSession12Shape, an AC-2 negative `--some-new-flag the-value`, and two AC-3 classifier cases (injected --plugin-dir prefix; unrecognized-flag fallback). `go test ./...` = 1018 passed in 12 packages; `gofmt -l internal/cli/` empty; `go vet ./...` clean.

### Summary

Root cause was a single incomplete set: `--plugin-dir` (re-injected by the parser) and the common host value-taking flags were absent from `valueTakingHostFlags`, so the classifier named the wrong token. Fixed by completing the value-taking sets and adding a conservative fallback that suppresses the prescriptive "put X before --" advice whenever the candidate's predecessor is an unrecognized space-form flag — equals-form predecessors stay confidently stray since they carry their own value. The argv-unchanged invariant is untouched (the guard is still pure-read + stderr); the captain's exact session-12 invocation now warns on `@/tmp/handoff.md`. Verified offline: 1018/1018.

## Stage Report: validation (cycle 1)

- DONE: Confirm the REJECTED blocker is closed — session-12 shape warns on `@/tmp/handoff.md`, never the dir.
  TestClaudeStrayPromptSession12Shape passes; real `runClaude` on `--plugin-dir /co -- --model gpt-x @/tmp/handoff.md` names `@/tmp/handoff.md`, asserts NOT `/co`, and pins the full unchanged argv. Functionally closed.
- DONE: Confirm the argv-unchanged invariant holds across all paths.
  Holds — guard is still pure-read + stderr (`warnStrayPromptAfterDash` takes fd by value, returns nothing). Every probe (session-12, fallback cases, codex misfire below) shows the inner argv intact.
- DONE: Full offline `go test ./...` green; gofmt/vet clean.
  12 packages ok; `gofmt -l internal/cli/` empty; `go vet ./...` clean.
- FAILED: Scrutinize the conservative fallback (c) for over-suppression / false-negatives that re-bury the captain's class.
  Concern 3(i) CONFIRMED but ACCEPTABLE: a genuine stray prompt placed immediately after an unrecognized space-form flag is now silent (e.g. `-- --dangerously-skip-permissions @file` → no warning). This is a deliberate, documented tradeoff and a NARROW residual — the scan continues, so a following positional or any recognized value/boundary flag re-surfaces the real prompt (`--newflag @a @b` warns on `@b`; `--newflag x --model g @file` warns on `@file`). It does NOT re-bury the captain's actual class (whose prompt follows the VALUE of recognized `--model`, so its predecessor is a positional → confidently stray). 3(ii) equals-form predecessors warn confidently — confirmed.
- FAILED: NEW MATERIAL — the injected-`--plugin-dir`-prefix bug is only HALF fixed; it still defeats the leading-subcommand exemption.
  `spacedock codex --plugin-dir <dir> -- exec <prompt>` (and `resume <id>`) MISFIRES: the guard warns on `exec`/`resume` with the actively-wrong advice "put `exec` before `--`". Root cause is the SAME injected-prefix interaction the cycle was meant to close — `parseFrontDoorArgs` re-prepends `--plugin-dir <dir>` so the subcommand lands at index 2, but the exemption only checks `subcommands[fd.passthrough[0]]`. Cycle-1 fixed the prefix interaction for the value-flag path (item a) but missed it for the subcommand path. Reproduced via real `runCodex` (argv intact: `… codex … --plugin-dir /co exec do the thing …`) and via the classifier directly for both `exec` and `resume`. Advisory-only (no launch-safety impact), but a real misfire on a plausible developer invocation, and the implementation report (item a) claims the injected-`--plugin-dir` prefix is fully accounted for.
- FAILED: Test-strength check — does the session-12 test red if `--plugin-dir` is dropped from `valueTakingHostFlags`?
  NO. Mutation: removing `--plugin-dir` from claude's value set leaves TestClaudeStrayPromptSession12Shape GREEN, because the conservative fallback (c) independently suppresses `/co` (its predecessor `--plugin-dir` is an unrecognized space-form flag). So `--plugin-dir`'s SET membership is redundant for the session-12 shape — the fallback is the load-bearing mechanism, not the set addition the implementation report (item a) credits. A regression that drops `--plugin-dir` from the set would not be caught by this test. Polish-tier (behavior stays correct) but worth a pinning test or report correction.

### Summary

REJECTED. The original blocker is functionally closed and the argv-unchanged invariant holds, but a detached-style adversarial probe of the new logic surfaced a fresh MATERIAL hole of the SAME family the cycle was meant to close: the injected `--plugin-dir <dir>` prefix still defeats the leading-subcommand exemption, so `spacedock codex --plugin-dir <dir> -- exec/resume <arg>` misfires with wrong advice. The subcommand check inspects only `passthrough[0]` and must skip a leading injected `--plugin-dir <dir>` pair (mirror the value-flag fix). Fix direction: have the leading-subcommand check (and the index-0 assumption) account for the re-prepended `--plugin-dir <dir>` pair, and add a codex `--plugin-dir … -- exec <prompt>` negative AC. Secondary (Polish): the session-12 test does not actually pin `--plugin-dir` set membership (the fallback masks it) — either add a pinning test or correct the implementation report's claim that the set addition fixes session-12. The fallback false-negative (3i) is accepted as a documented, narrow tradeoff.

## Stage Report: implementation (cycle 2)

- DONE: Structural fix (not another per-path patch): skip the leading injected `--plugin-dir <dir>` pair ONCE up front, then run BOTH the value-flag scan AND the leading-subcommand check against the real after-`--` tokens.
  Added `skipInjectedPrefix(passthrough)` which strips leading `--plugin-dir <dir>` pair(s); `strayPromptAfterDash` runs both the subcommand check and the token scan against the skipped slice. `--plugin-dir` REMOVED from `valueTakingHostFlags` (both hosts) so the structural skip is the single, load-bearing prefix mechanism. Confirmed parseFrontDoorArgs re-prepends ONLY `--plugin-dir` (safehouse knobs → fd.safehouseFlags, booleans consumed), so the skip set is complete. Commit 8ddd957c.
- DONE: Add a codex `--plugin-dir <dir> -- exec <prompt>` negative AC + a `--plugin-dir <dir> -- resume <id>` case (no warning, argv unchanged).
  Both added to TestStrayPromptGuardNegatives (launch-seam) and as classifier cases in TestStrayPromptAfterDashClassifier; a multiple-injected-`--plugin-dir`-pairs case too. Real `runCodex` on both shapes emits no warning and the inner argv is intact.
- DONE: Make the structural skip the pinned, load-bearing mechanism (the session-12 test alone did not red on `--plugin-dir` set-drop).
  Mutation-verified: nooping `skipInjectedPrefix` REDS 6 subtests (the two subcommand-behind-prefix launch negatives + two classifier cases + their parent tables). Also corrected the cycle-1 report's item-(a) credit in place (the set addition was redundant; the fallback was load-bearing in cycle 1, the structural skip is in cycle 2).
- DONE: go test ./... + vet/gofmt green; argv-unchanged invariant intact.
  `go test ./...` = 1022 passed in 12 packages; `gofmt -l internal/cli/` empty; `go vet ./...` clean. Guard remains pure-read + stderr; all launch-seam tests assert the byte-identical inner argv.

### Summary

Replaced cycle 1's piecemeal per-path handling of the injected `--plugin-dir <dir>` prefix with one structural skip: `skipInjectedPrefix` strips the leading pair(s) once, and every classifier check (subcommand exemption, value-flag scan, conservative fallback) runs against the real after-`--` tokens behind it — closing the value-flag, subcommand, and any future per-token path in a single place. Removed `--plugin-dir` from the value-taking sets so the skip is the genuine load-bearing mechanism, mutation-confirmed (nooping it reds the new subcommand-behind-prefix tests). The codex `--plugin-dir … -- exec`/`resume` misfire is closed, the argv-unchanged invariant is untouched, and the full offline suite is green (1022/1022). The cycle-1 item-(a) credit was corrected in place per the polish note.

## Stage Report: validation (cycle 2)

- DONE: Confirm the subcommand misfire is closed — codex `--plugin-dir <dir> -- exec/resume <arg>` no longer warns, argv unchanged.
  Real `runCodex` on `--plugin-dir /co -- exec "do the thing"`, `-- resume abc`, and multi-`--plugin-dir … -- exec p` all emit NO warning and forward the inner argv verbatim (`… codex … --plugin-dir /co exec do the thing …`). TestStrayPromptGuardNegatives' two new subcommand-behind-prefix cases + the classifier cases pass.
- DONE: Confirm the structural skip is genuinely load-bearing (not a tautology) and session-12 still warns.
  Mutation on a throwaway copy: nooping `skipInjectedPrefix` (return passthrough unchanged) REDS exactly 6 FAIL lines — 4 leaf subtests (the two subcommand-behind-prefix launch negatives + the two classifier cases `subcommand behind injected --plugin-dir prefix` and `multiple injected pairs`) plus their 2 parent table rollups. Confirms the team-lead's "6 subtests" claim. Under the noop, the codex `exec`-behind-prefix path misfires again (returns `("exec", true)`), proving the skip is what closes it. Session-12 still warns on `@/tmp/handoff.md` (real runCodex `--plugin-dir /co -- -m gpt @file` → warned=true, argv intact).
- DONE: argv-unchanged invariant intact across all paths; multiple-injected-pairs handled; cycle-1 item-(a) credit corrected.
  Guard remains pure-read + stderr; every launch-seam case asserts the byte-identical inner argv and all pass. Multiple `--plugin-dir` pairs are stripped (the `for len>=2` loop) and forwarded intact. The cycle-1 item-(a) `[CYCLE-2 CORRECTION: …]` annotation is present in the entity (line 135).
- DONE: Adversarial probe of the structural skip — does it ever wrongly strip a NON-injected `--plugin-dir` typed after `--`?
  It CAN, but the outcome is correct by design. Operator-typed `-- --plugin-dir /host/dir @file` is textually identical to the injected shape; `skipInjectedPrefix` strips `--plugin-dir /host/dir` and warns on the real trailing `@file` — correct, because a leading `--plugin-dir <dir>` is the flag's value regardless of origin (the code comment states this; the dir is never a stray prompt). Only contrived residual: `-- --plugin-dir @file` (operator's prompt typed where the dir-value belongs) strips both tokens and stays silent — a false-negative, but advisory-only and the operator's own argv genuinely reads `@file` as the plugin-dir value. Not blocking.
- DONE: Full offline `go test ./...` green; gofmt/vet clean.
  12 packages ok; `gofmt -l internal/cli/` empty; `go vet ./...` clean.

### Summary

PASSED. The cycle-1 subcommand misfire is closed by a clean structural fix: `skipInjectedPrefix` strips the spacedock-injected leading `--plugin-dir <dir>` pair(s) once, so the subcommand exemption, value-flag scan, and conservative fallback all run against the operator's real after-`--` grammar. Mutation-verified the skip is genuinely load-bearing for the subcommand path (nooping it reds 6 subtests and re-opens the `exec`/`resume` misfire), the argv-unchanged invariant holds in every path (guard stays pure-read + stderr), multiple injected pairs are handled, and the cycle-1 item-(a) credit was corrected in place. The adversarial probe confirmed the skip is sound regardless of `--plugin-dir` origin (the dir is always the flag's value); the only residual is a contrived, advisory-only false-negative. No material findings.
