---
title: 'FO contract: replace bare `spacedock` command examples with `${SPACEDOCK_BIN:-spacedock}`'
status: validation
source: 'boot-forensics (2026-06-16) — FO used homebrew cask (0.20.2) instead of $SPACEDOCK_BIN (dev) throughout a session because contract command blocks model the bare-spacedock shorthand. SPACEDOCK_BIN was set and correct but never used; cap divergence cost ~16k tok in failed --read calls. Fix: find-and-replace executable command positions in first-officer-shared-core.md and claude-first-officer-runtime.md.'
score: 0.5
id: 13f8b12x9f7ba25ywm5wt2x7
started: 2026-07-03T03:03:12Z
worktree: .worktrees/spacedock-ensign-fo-contract-bin-shorthand
mod-block: merge:pr-merge
pr: "#473"
verdict:
completed:
archived:
---

## Problem statement

Contract command examples that model bare `spacedock` in executable positions teach the FO to run whatever binary is on `$PATH`, silently ignoring `SPACEDOCK_BIN`. Surfaced by boot-forensics on 2026-06-16: `$SPACEDOCK_BIN` was the local dev build (has `status --read`); PATH `spacedock` was the homebrew 0.20.2 cask (no `--read`). Both satisfied `contract 1`, so the version gate passed silently and every `status --read` fell through to the full entity dump (~8k tok each, two hits = ~16k tok wasted).

Since the seed was filed, the launcher-invariant work landed the `${SPACEDOCK_BIN:-spacedock}` invariant paragraph in both shared cores, converted 17 call sites, and shipped `internal/contractlint/launcher_invariant_test.go`. But that lint is green today while 8 bare executable-position sites still ship, because of three escape classes:

1. **Scope** — the lint walks only `skills/first-officer/references/`; the deferred-skill SKILL.md files and `skills/ensign/references/` are unscanned.
2. **Code blocks** — the matcher is backtick-anchored; fenced-block command lines carry no backticks, and multi-line continuations (`spacedock dispatch build \`) put the flags on later lines.
3. **Flag allowlist** — the enumerated invocation-flag list misses `--advance`, `--folder`, `--id-seed`, `--id-actor`. This is the same class as the `--name` escape already recorded in the lint's own discriminator comment ("the line-132 escape class"); an allowlist re-opens with every new flag.

## Current-state audit (2026-07-03)

Surface audited: `skills/first-officer/references/*.md` (7 files), `skills/ensign/references/*.md` (4 files), and the deferred-skill files `skills/{fo-status-viewer,fo-write-core,present-gate,feedback-rejection-flow,using-legacy-claude-team,fo-dispatch-recovery}/SKILL.md`.

**Already fixed — 17 resolved-launcher sites:** first-officer-shared-core.md 7/9/10/11/15/17 (version gate, state management), fo-merge-core.md 9/16, fo-dispatch-core.md 40, claude-fo-dispatch.md 52/100/112 (spawn-standing-all, context-budget, reconcile), ensign-shared-core.md 7 (the invariant paragraph), fo-write-core/SKILL.md 13, fo-status-viewer/SKILL.md 11/15/38. The ensign references and present-gate, feedback-rejection-flow SKILL.md files carry zero bare occurrences.

**Bare executable-position sites remaining — N=8, each with its escape class:**

| # | Site | Bare form | Escape class |
|---|------|-----------|--------------|
| 1 | fo-dispatch-core.md:17 | fenced block: `spacedock status --workflow-dir {workflow_dir} --set {slug} …` | code block |
| 2 | fo-dispatch-core.md:121 | fenced block: `spacedock dispatch build \` (+ continuation flags) | code block |
| 3 | claude-first-officer-runtime.md:35 | "Use `spacedock new <slug> [--folder] [--id-seed S --id-actor A]` via Bash" | flag allowlist |
| 4 | claude-fo-dispatch.md:38 | "run `spacedock dispatch build --advance`" | flag allowlist |
| 5 | claude-fo-dispatch.md:44 | "when `spacedock dispatch build --advance` exits non-zero" | flag allowlist |
| 6 | codex-first-officer-runtime.md:20 | "`spacedock dispatch build --advance`'s emitted `output.prompt`" | flag allowlist |
| 7 | fo-write-core/SKILL.md:14 | "via `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < stub`" | scope + flags |
| 8 | fo-write-core/SKILL.md:29 | "`spacedock new <slug> --id-seed "{slug-or-title}"` mints it" | scope + flags |

**Allowed-context bare sites that stay bare (~24):** `→` capability-binding lines (fo-dispatch-core 108/112/135/155, shared-core 113/118/123/127); name-only prose mentions with no invocation flag — naming a command is not invoking it (fo-dispatch-core 116, claude-fo-dispatch 7/23/48/102, shared-core 40/145, pi-first-officer-runtime 7, fo-status-viewer 21, fo-write-core 3/28/30/32, fo-dispatch-recovery 3/49, using-legacy-claude-team 22/23/24/42 — verb spans and delta-flag spans are separate there); diagnostic/install/help contexts (`spacedock new --help` in claude-first-officer-runtime 35, `go build -o spacedock ./cmd/spacedock` in claude-fo-dispatch 116, the captain command `/spacedock bare`, the version-gate PATH fallback probe).

## Proposed approach

Three moves, all in the launcher-invariant lint's existing home:

1. **Convert the 8 sites** (exact wording below): the `spacedock` token at the invocation position becomes `${SPACEDOCK_BIN:-spacedock}`; no other wording changes. Site 5 is converted rather than exempted — the thing whose non-zero exit triggers break-glass IS the resolved-launcher invocation, so the resolved form is the more accurate text and avoids carving a new exemption.
2. **Extend `launcher_invariant_test.go`** so each escape class is closed structurally:
   - **R1 (backtick spans):** a span `` `spacedock <verb> …` `` containing any `--[a-z][a-z-]*` flag is an invocation. The generic flag matcher replaces the enumerated allowlist, retiring that escape class permanently; `--help` moves from "unlisted flag" to an explicit exemption. Verb set gains `contract` (the version-gate verb) alongside status/state/dispatch/merge/new.
   - **R2 (code blocks):** track ``` fence state while scanning; inside a fence, a line matching `^\s*spacedock\s` is command position and must be resolved — no flag requirement, since command position in a block is executable by definition.
   - **Span-scoped exemptions:** already-resolved form and the `--help` context are evaluated per backtick span, not per line. This matters concretely: claude-first-officer-runtime.md:35 carries both a must-fix invocation span and a legitimate `spacedock new --help` span on ONE line — a line-scoped `--help` exemption would swallow the violation. `→`-binding stays line-shaped; `on $PATH`/doctor/brew install/go build stay line contexts.
   - **Scope:** three walked groups — the FO references dir, the ensign references dir, and the six deferred-skill SKILL.md paths — each with a scanned>0 guard so a moved directory cannot make the lint pass vacuously. Commission/refit/survey/debrief/integration are deliberately out of scope: they author operator-facing scaffolding where bare `spacedock` is correct usage documentation.
3. **Extend the discriminator test** with one MUST-flag and one MUST-pass case per new class (fenced-block line, `--advance`-style unlisted flag, deferred-skill path, same-line invocation+`--help` mix).

## Text changes (before → after)

Only the launcher token changes; surrounding wording is untouched.

1. fo-dispatch-core.md:17 — `spacedock status --workflow-dir {workflow_dir} --set {slug} status={next_stage} …` → `${SPACEDOCK_BIN:-spacedock} status --workflow-dir {workflow_dir} --set {slug} status={next_stage} …`
2. fo-dispatch-core.md:121 — `spacedock dispatch build \` → `${SPACEDOCK_BIN:-spacedock} dispatch build \`
3. claude-first-officer-runtime.md:35 — "Use `spacedock new <slug> [--folder] [--id-seed S --id-actor A]`" → "Use `${SPACEDOCK_BIN:-spacedock} new <slug> [--folder] [--id-seed S --id-actor A]`"; the same line's "see `spacedock new --help`" stays bare (help context).
4. claude-fo-dispatch.md:38 — "run `spacedock dispatch build --advance`" → "run `${SPACEDOCK_BIN:-spacedock} dispatch build --advance`"
5. claude-fo-dispatch.md:44 — "when `spacedock dispatch build --advance` exits non-zero" → "when `${SPACEDOCK_BIN:-spacedock} dispatch build --advance` exits non-zero"
6. codex-first-officer-runtime.md:20 — "`spacedock dispatch build --advance`'s emitted `output.prompt`" → "`${SPACEDOCK_BIN:-spacedock} dispatch build --advance`'s emitted `output.prompt`"
7. fo-write-core/SKILL.md:14 — "via `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < stub`" → "via `${SPACEDOCK_BIN:-spacedock} new <slug> [--folder] [--id-seed S --id-actor A] < stub`"
8. fo-write-core/SKILL.md:29 — "`spacedock new <slug> --id-seed "{slug-or-title}"` mints it" → "`${SPACEDOCK_BIN:-spacedock} new <slug> --id-seed "{slug-or-title}"` mints it"

## Acceptance criteria

- **AC-1 (measures the end value, divergeable baseline):** The bare-executable-position count across the shipped contract surface (FO references + ensign references + the six deferred-skill SKILL.md files) is 0, down from the audited baseline of 8, and the extended lint pins it at 0. The baseline can move the wrong way — any contract edit can reintroduce a bare invocation — and the lint is the gate that stops it. Tested by: `go test ./internal/contractlint -run ResolvedLauncher` green on the converted tree; run against the pre-conversion tree the extended lint reports exactly the 8 audited sites (this run recorded in the implementation stage report).
- **AC-2 (falsifiability):** A seeded bare invocation in each newly covered class turns the lint RED: (a) a `spacedock <verb>` command line inside a fenced code block, (b) a backtick invocation whose only flag is outside the old allowlist (e.g. `--advance`), (c) a bare invocation in a deferred-skill SKILL.md. Tested by: discriminator MUST-flag unit cases for each shape, plus one live seeded-violation run (edit a surface file → observe FAIL output → revert) recorded in the implementation stage report.
- **AC-3 (no over-flagging):** Every audited allowed-context bare site is byte-identical after the change and the lint is green: `→` binding lines, name-only prose mentions, the same-line `spacedock new --help` span next to a converted invocation, `go build -o spacedock`, `/spacedock bare`, the version-gate PATH fallback probe. Tested by: discriminator MUST-pass cases per exemption class; full lint green; gate diff review confirms only the 8 sites changed.

## Test plan

- Unit (Go, `internal/contractlint`): extend `TestBareLauncherHelperScannerDiscriminates` with the new leak/pass shapes; extend `TestFOReferencesUseResolvedLauncher` to walk the three surface groups with per-group zero-scan guards. Cost: small — one test file, no fixtures.
- Falsifiability run during implementation: seed one violation per AC-2 class, observe RED, revert; cite the failing output in the stage report.
- No CLI, fixture, or live-workflow tests: no binary behavior changes. The runtime consequence of launcher drift (SPACEDOCK_BIN-vs-PATH divergence) is already proven by the `internal/ensigncycle` live drive the existing lint header cites; this task is doc conversion plus lint coverage.

## Spike determination

No spike needed: the design rests on proven mechanisms — (1) the resolved launcher form is exercised live on every dispatch (the fetch-command launcher shim; the ensigncycle SPACEDOCK_BIN-vs-PATH live drive), (2) line-scan structural linting over the reference surface is the established pattern in `launcher_invariant_test.go`, (3) fence-state tracking is plain string toggling on ``` lines with no parser dependency.

## Documentation changes

None on the docs site: this changes agent-contract skill text only — no CLI output, command surface, or operator-facing behavior. The before/after list above is the complete diff spec; ideation runs without a worktree, so it lives here.

## Stage Report: ideation

- DONE: Current-state audit FIRST: this seed predates the launcher-invariant work — the shared core now carries the `${SPACEDOCK_BIN:-spacedock}` invariant paragraph and many converted call sites. Enumerate every executable command position across the ENTIRE shipped contract surface (first-officer-shared-core.md, both runtime adapters, fo-dispatch-core.md, claude-fo-dispatch.md, fo-merge-core.md, the deferred skills, the ensign contract), classify each bare `spacedock` as executable-position (must fix) vs the invariant's allowed contexts (prose naming, `→` binding lines, fallback probe, install hints), and record what is already fixed vs what remains.
  Audit section records 17 already-resolved sites, 8 remaining bare executable-position sites (each tagged with its lint escape class), ~24 allowed-context sites; ensign references and present-gate/feedback-rejection-flow are clean.
- DONE: The guarantee ships as a code gate, not prose-only: an internal/contractlint structural check (the quarantine's sanctioned home) that fails on a bare `spacedock` in an executable command position, with the allowed-context exemptions expressed structurally — so the class cannot creep back with the next contract edit.
  Approach extends the existing `launcher_invariant_test.go`: generic flag matcher (retires the allowlist escape class), fenced-code-block rule, deferred-skill + ensign scope with zero-scan guards, span-scoped exemptions.
- DONE: At least one AC MEASURES the end value divergeably: bare-executable-position count across the surface goes N→0, the lint pins it at 0, and a seeded violation demonstrably turns the lint RED (the falsifiability run recorded).
  AC-1 measures 8→0 with the lint as the pin; AC-2 requires the seeded-violation RED run recorded; AC-3 guards against over-flagging.

### Summary

Audited the full contract surface with the extended rule simulated via grep: 8 bare executable-position sites remain (2 fenced-block, 4 unlisted-flag, 2 unscanned-scope), all invisible to the current lint, which passes green today. Ideation specifies the exact 8-site conversion, a three-part lint extension (generic flags, fence tracking, scope groups with vacuity guards, span-scoped `--help` exemption for the same-line mix at claude-first-officer-runtime.md:35), and measured ACs with a falsifiability run. No spike needed — all mechanisms are proven in the existing lint and the ensigncycle live drive.

## Stage Report: implementation

- DONE: Exactly the 8 audited sites convert, per the body's before→after spec verbatim — only the launcher token changes; every audited allowed-context bare site stays byte-identical (AC-3's diff-review bar: nothing but the 8 sites in the contract-surface diff).
  Commit `dae763d6`, worktree `spacedock-ensign-fo-contract-bin-shorthand`: `git diff --stat` on the 5 changed doc files shows 8 insertions/8 deletions (one line each), matching the 8-site table exactly; no other line in those files or any other scanned file changed.
- DONE: The lint extension closes all three escape classes structurally (fence-state command-position rule, generic `--flag` matcher retiring the enumerated allowlist with `--help` as an explicit span-scoped exemption, the three walked surface groups each with a scanned>0 vacuity guard) — and the pre-conversion tree run reports exactly the 8 audited sites (AC-1's divergeable baseline, recorded in the stage report).
  `internal/contractlint/launcher_invariant_test.go` rewritten: `scanBareLauncherCalls` tracks fence state (```` ``` ```` toggling) and flags any `^\s*spacedock\s` command-position line inside a fence with no flag requirement; `lineHasBareLauncherHelperCall` now extracts backtick spans per line and applies the generic `--[a-z][a-z-]*` flag matcher plus a per-span `--help`/resolved-form exemption (was line-scoped, swallowing the claude-first-officer-runtime.md:35 mixed-span case); `launcherSurfaceGroups` walks FO references, ensign references, and the six deferred-skill `SKILL.md` paths, each with its own scanned>0 guard. A throwaway test (`git show HEAD~1:{path}` piped through `scanBareLauncherCalls`, deleted after use, not committed) measured the pre-conversion tree at exactly 8 violations across the 5 changed files, matching the audit table 1:1 (verified line-by-line against the site table).
- DONE: AC-2 falsifiability demonstrated live: one seeded violation per new class (fenced-block line, unlisted-flag backtick span, deferred-skill SKILL.md) turns the lint RED with output cited, then reverted clean; discriminator MUST-flag/MUST-pass cases added per class; full `go test ./...` green.
  Live seeded runs (each edited, `go test ./internal/contractlint -run TestLauncherSurfaceUsesResolvedLauncher`, observed FAIL, `git checkout --` reverted): (a) fenced bare `spacedock status --workflow-dir {workflow_dir} --discover` in `ensign-shared-core.md` → FAIL citing line 127; (b) bare `spacedock new <slug> --folder` (unlisted flag) in `first-officer-shared-core.md` → FAIL citing line 173; (c) bare `spacedock status --set {slug} status=validation` in `present-gate/SKILL.md` (deferred-skill scope) → FAIL citing line 49. All three reverted; `git status --short` clean afterward. Discriminator additions: `TestBareLauncherHelperScannerDiscriminates` gained `--advance` unlisted-flag and same-line invocation+`--help`-mix MUST-flag/MUST-pass pairs; new `TestFencedLauncherBlockScannerDiscriminates` (fence MUST-flag/MUST-pass) and `TestDeferredSkillLauncherScopeDiscriminates` (scope MUST-pass, six exact paths). `go test ./...` — 15 packages, all green (`internal/contractlint` 0.2–2s across runs).

### Summary

Converted the 8 audited bare-`spacedock` executable-position sites to `${SPACEDOCK_BIN:-spacedock}` (commit `dae763d6`) and extended `launcher_invariant_test.go` to close the three escape classes structurally: fence-state tracking for code-block command lines, a generic flag matcher replacing the enumerated allowlist with a span-scoped `--help` exemption, and scope widened to ensign references plus the six deferred-skill `SKILL.md` files, each behind a scanned>0 guard. Measured the pre-conversion baseline at exactly 8 (matching the audit) and live-seeded one RED violation per new escape class, reverting cleanly each time; full `go test ./...` is green.

### Feedback Cycles

**Cycle 1 — validation (2026-07-03), detached adversarial audit, group-drop mutation stays green.** The audit's "drop a surface group from the walk" edit — deleting the `{"ensign references", …}` entry from `launcherSurfaceGroups` in `launcher_invariant_test.go` — leaves `go test ./internal/contractlint` fully GREEN. The scanned>0 vacuity guard is per-group, so it cannot fire for a group that is no longer registered, and `TestDeferredSkillLauncherScopeDiscriminates` pins only `deferredSkillPaths()`' contents, not that any group is walked. Silent scope loss is this task's own demonstrated failure mode (the deferred skills went unscanned until this change). All other audit edits red correctly: reverting a converted site, the `--help` line-scope re-widening (caught by the mixed-line discriminator), the extractor-suffix break (scanned>0 guard), a moved references directory (walk Fatalf), and a fenced-block sneak. Ask: extend the scope discriminator to pin `launcherSurfaceGroups(t)` to exactly the three named groups ("first-officer references", "ensign references", "deferred-skill SKILL.md files"), each with a non-empty path set, so deleting a group registration reds.

## Stage Report: validation

- DONE: Every `**AC-N**` verified by REPRODUCING its evidence: AC-1's baseline re-run yourself (lint against the pre-conversion tree reports exactly the 8 audited sites; against head reports 0) — not read off the report; AC-2's seeded-violation RED re-demonstrated for at least one class end-to-end (seed, observe the failing output, revert clean); AC-3's no-over-flagging bar checked by diffing the contract surface — exactly the 8 sites changed, every audited allowed-context bare site byte-identical, full lint + `go test ./...` green.
  AC-1: on a throwaway clone at `dae763d6` with `git checkout 58388c42 -- skills/`, `TestLauncherSurfaceUsesResolvedLauncher` FAILs with exactly 8 violations matching the audit table 1:1 (claude-first-officer-runtime.md:35, claude-fo-dispatch.md:38/44, codex-first-officer-runtime.md:20, fo-dispatch-core.md:17/121, fo-write-core/SKILL.md:14/29); restored to head it passes with 0. AC-2: re-seeded all three classes live (fenced `spacedock state ready --commit` in fo-dispatch-core.md → FAIL at :18; backtick `spacedock dispatch build --advance` in first-officer-shared-core.md → FAIL at :171; backtick `spacedock status --set` in present-gate/SKILL.md → FAIL at :47), each reverted to a clean `git status`. AC-3: `git diff 58388c42 dae763d6 -- 'skills/**'` shows exactly the 8 one-token conversions and nothing else (allowed-context bare sites byte-identical, including the same-line `spacedock new --help` span); `go test -count=1 ./...` green (exit 0, zero non-ok lines).
- DONE: Detached adversarial audit on a THROWAWAY checkout (the shipped contract/scaffolding is a Proof-policy high-stakes surface): construct claim-breaking edits the lint must catch — e.g. revert one converted site, re-widen the `--help` exemption to line scope next to an invocation, drop a surface group from the walk (the vacuity guard must red), sneak a bare invocation into a fenced block — and confirm a test reds on EACH; record findings as `### Feedback Cycles` or "refuted nothing material".
  Ran on a scratchpad clone at `dae763d6` (deleted after). Reds confirmed: converted-site revert (fo-write-core/SKILL.md:29 → lint FAIL), line-scoped `--help` re-widening (mixed-line discriminator FAIL at launcher_invariant_test.go:279), extractor-suffix break (scanned>0 guard: "scanned zero files in group"), moved ensign references dir (walk Fatalf plus sibling lints), fenced-block sneak (FAIL, same run as AC-2a). NOT caught: deleting a group entry from `launcherSurfaceGroups` — suite stays green; recorded as Feedback Cycle 1 above.
- DONE: A PASSED/REJECTED recommendation with honest accounting: the lint's expected values must come from outside the files under test (structural shape rules, not text echoes — this is contractlint's sanctioned quarantine, verify it stays structural); flag anything self-referential.
  Recommendation: REJECTED (cycle 1, narrow) — solely on the group-drop discriminator gap; every AC passes as specified. Structural check clean: the surface lint expects zero scanner hits (rule-derived, not file-echoed), discriminator inputs are synthetic line shapes classified against the invariant and free to diverge from file content, and the guards assert scan-shape, not text. Nothing self-referential found.

### Summary

All three ACs reproduce exactly as claimed: pre-conversion baseline of 8 sites (1:1 with the audit table), 0 at head, three live seeded REDs reverting clean, and a diff containing only the 8 one-token conversions with the full uncached suite green. The detached adversarial audit refuted one thing: dropping a surface-group registration from `launcherSurfaceGroups` is invisible to the suite (per-group vacuity guards cannot defend a deregistered group). Recommending REJECTED with a small, concrete ask — a discriminator pinning the three walked groups — after which this is a PASS.

## Stage Report: implementation (cycle 2)

- DONE: The scope discriminator pins `launcherSurfaceGroups(t)` to exactly the three named groups, each with a non-empty path set — and the cycle-1 group-drop mutation (deleting the ensign-references entry) now demonstrably turns the suite RED (mutation applied → FAIL output cited → reverted clean).
  Added `TestLauncherSurfaceGroupsScopeDiscriminates` (commit `ca13076b`): calls `launcherSurfaceGroups(t)` directly, asserts exactly 3 groups named `"first-officer references"`, `"ensign references"`, `"deferred-skill SKILL.md files"` in that order, each with `len(paths) > 0`. Re-applied the exact cycle-1 mutation (deleted the `{"ensign references", mdFilesIn(t, ensignReferenceDir(t))}` line) and ran `go test ./internal/contractlint -run Launcher -v`: `TestLauncherSurfaceGroupsScopeDiscriminates` FAILed — `launcherSurfaceGroups returned 2 groups, want 3 (a deregistered group would silently stop being scanned)` — while every other test (including the old `TestDeferredSkillLauncherScopeDiscriminates`) stayed green, confirming this new test is what closes the gap. Reverted; `git diff --stat` afterward showed only the net +23 lines of the new test, no residue from the mutation.
- DONE: Everything already green stays green: full `go test ./...` uncached, the 8-site conversion and all discriminator cases untouched.
  `go test ./... -count=1`: 15 packages, all `ok` (`internal/contractlint` 1.7s). `git diff dae763d6 ca13076b` touches only `internal/contractlint/launcher_invariant_test.go` (+23/-0); the 8-site doc conversion and every prior discriminator case are byte-identical to cycle 1.

### Summary

Closed the cycle-1 gap by adding `TestLauncherSurfaceGroupsScopeDiscriminates`, which calls `launcherSurfaceGroups(t)` directly and pins it to exactly the three named groups with non-empty path sets — the prior scope discriminator only checked `deferredSkillPaths()` in isolation and never exercised the group-registration list itself. Re-ran the exact cycle-1 group-drop mutation live: the new test now FAILs with a clear diagnostic, reverted cleanly, and the full uncached `go test ./...` suite stays green. Commit `ca13076b`.

## Stage Report: validation (cycle 2)

- DONE: re-apply your original group-drop mutation on a throwaway copy and confirm the suite now goes RED
  On a fresh throwaway clone at `ca13076b`, deleting the `{"ensign references", …}` entry from `launcherSurfaceGroups` reds `TestLauncherSurfaceGroupsScopeDiscriminates` with the exact diagnostic the cycle-2 report cites ("launcherSurfaceGroups returned 2 groups, want 3"); a variant dropping the deferred-skill group entry instead reds the same way. Both reverted to a clean tree; throwaway deleted after.
- DONE: confirm everything previously green stays green (full uncached suite)
  `go test -count=1 ./...` at `ca13076b` in the worktree: exit 0, zero non-ok lines. `git diff --stat dae763d6 ca13076b` touches only `internal/contractlint/launcher_invariant_test.go` (+23/-0) — the 8-site conversion and all prior discriminators are byte-identical to cycle 1.
- DONE: verify the cycle-2 stage report's claimed falsifiability run holds
  Reproduced verbatim: the claimed FAIL message, the mutation shape, and the clean revert all match the report. The new test stays structural — its expected values (three group names, each non-empty) are spec-derived constants, not echoes of the files under test.

### Summary

Re-review of the cycle-2 fix (`ca13076b`): the group-drop mutation that was invisible in cycle 1 now reds through `TestLauncherSurfaceGroupsScopeDiscriminates`, a variant group-drop reds identically, the full uncached suite is green, and the fix adds nothing beyond the +23-line discriminator. Feedback Cycle 1 is closed. Recommendation: PASSED.
