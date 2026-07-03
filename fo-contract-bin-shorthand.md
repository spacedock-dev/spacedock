---
title: 'FO contract: replace bare `spacedock` command examples with `${SPACEDOCK_BIN:-spacedock}`'
status: ideation
source: 'boot-forensics (2026-06-16) — FO used homebrew cask (0.20.2) instead of $SPACEDOCK_BIN (dev) throughout a session because contract command blocks model the bare-spacedock shorthand. SPACEDOCK_BIN was set and correct but never used; cap divergence cost ~16k tok in failed --read calls. Fix: find-and-replace executable command positions in first-officer-shared-core.md and claude-first-officer-runtime.md.'
score: 0.5
id: 13f8b12x9f7ba25ywm5wt2x7
started: 2026-07-03T03:03:12Z
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
