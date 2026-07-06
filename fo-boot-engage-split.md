---
id: k74gt0qv3j4b86knvy2rhsta
title: Lighten the interactive boot greet — managed-workflow list + an explicit "engage" verb; no forced workflow pick, no gate render at the greet
status: validation
source: "FO session 2026-07-04: interactive boot ran ~5 minutes; original framing was Startup step 8 rendering a full present-gate review per ready gate before stopping. Captain's subspace-tui review (2026-07-04) redirected the scope wider: the launcher bootstrap prompt (frontdoor.go bootstrapPrompt/codexBootstrapPrompt, both literally ending '...Engage.') should drop that flourish since 'engage' becomes a real captain-invoked verb; Startup step 3's multi-workflow pick goes away in favor of just listing managed workflows; gate assembly moves behind the engage verb. RESOLVED (captain, 2026-07-04): engage is an FO interaction verb, not a new binary command; scope is the current/named workflow only for now, with the «engage»(workflow) signature deliberately left open to a future multi-workflow extension; frontdoor.go referent confirmed. Remaining coordination note (not a blocker): 7v (pi-bootstrap-prompt-parity) hardcodes the CURRENT codexBootstrapPrompt text including 'Engage.' as pi's target — sequence this entity before 7v. Flagship member of the 0.25.0 fo-behavioral-discipline sprint alongside z25/zm/vcm."
started: 2026-07-04T10:38:11Z
completed:
verdict:
score: 0.35
worktree: .worktrees/spacedock-ensign-fo-boot-engage-split
issue:
sprint: 0250-fo-behavioral-discipline
mod-block: merge:pr-merge
pr: "#475"
---

The interactive FO boot is too heavy: the launcher bootstrap prompt (`internal/cli/frontdoor.go`) ends with a flourish "Engage.", Startup step 8 renders a full `present-gate` review for every ready gate before stopping, and step 3 asks the captain to pick when multiple workflows are discovered. Redirect (captain, 2026-07-04 — RESOLVED): make the greet light — list the managed workflow(s), hint "Use engage", no forced workflow pick, no gate render at the greet. Define a new prose-function `«engage»(workflow)`: run the existing `«dispatch.next-action»()` event-loop skeleton (dispatch any ready entity, present any ready gate, advance any completed non-gated stage) to its stopping condition for the NAMED workflow — no new binary command, no new dispatch mechanism, just a captain-invoked entry point into logic the contract already defines. Scope for THIS entity is the current/named workflow only, not all managed workflows at once — no boot-sequencing restructure needed; headless/single-entity mode is unaffected (the no-pick change is interactive-only). Forward-compatible by construction: `«engage»(workflow)` takes a workflow argument now so a later multi-workflow call (`«engage»(workflow1, workflow2, …)` or an all-managed default) extends the same signature rather than requiring a redesign — sweeping multiple workflows at once is explicitly a future extension, not in scope here, and must not be precluded by this entity's design. Coordinate with `7v` (pi-bootstrap-prompt-parity): 7v's Approach hardcodes the CURRENT `codexBootstrapPrompt` text (including "Engage.") as pi's target — sequence this entity's bootstrap-prompt edit BEFORE 7v, or have 7v's implementer re-derive the target from the live constant rather than 7v's written quote. This is the flagship member of the 0.25.0 behavioral-discipline sprint (captain, 2026-07-04) alongside z25/zm/vcm.

## Problem

The interactive FO boot is heavy and ceremonial in three places, none of which the captain asked for before they've given any direction:

1. **Startup step 8 (interactive branch) directs a `present-gate` render at the greet** — "present the summary (and any ready `gate: true` gate as captain-facing text)"; the "## Deferred load points" note reinforces it: "presents any ready gate via `present-gate`." A contract-compliant boot therefore assembles a full gate review (stage-report read + AC cross-check) for EVERY ready gate before stopping, so greet cost scales with the ready-gate count. One observed boot ran ~5 minutes.
2. **Startup step 3 forces a workflow choice** when `status --discover` returns multiple paths ("multiple → present the list"), which reads as "pick which one" before the captain has said what they want.
3. **The launcher bootstrap prompt ends with a throwaway flourish "Engage."** (`internal/cli/frontdoor.go` `bootstrapPrompt` + `codexBootstrapPrompt`) that now collides with `engage` becoming a real captain-invoked verb.

A captain who just wants to see what's managed and then explicitly trigger a sweep can't: the boot front-loads assembly work before there is any per-entity direction.

## Proposed approach

Two waves with different collision profiles (per the 0250 sprint doc's own sequencing).

### Wave 0 — bootstrap-prompt code half (FIRST, parallel-safe, Go-only)

Drop " Engage." from both launcher constants in `internal/cli/frontdoor.go`; `pi.go`'s `piBootstrapPrompt` carries no flourish and is untouched.

- `bootstrapPrompt` (line 25):
  - before: `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage."`
  - after: `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too."`
- `codexBootstrapPrompt` (line 533):
  - before: `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Engage. Assume $spacedock:first-officer for the entire session."`
  - after: `"You totally got this. Take your time. I love you. And tell all subagents and team members you love them too. Assume $spacedock:first-officer for the entire session."`
- mirror both in the launcher oracle constants `wantBootstrapPrompt` / `wantCodexBootstrapPrompt` (`internal/cli/safehouse_frontdoor_test.go` lines 18, 234), which the existing argv assertions (lines 55/78/159/256/317/420) already pin as the appended last token.

Wave 0 goes first because it unblocks the **`7v` (pi-bootstrap-prompt-parity) coordination**: 7v's Approach hard-quotes the CURRENT `codexBootstrapPrompt` (WITH "Engage.") as pi's byte-identity target. Sequence this Wave-0 edit BEFORE 7v, or have 7v re-derive its target from the live constant rather than its written quote — otherwise 7v re-introduces the flourish on pi.

### Wave 1 — contract-prose half (STRICT SERIAL within the sprint — shares `first-officer-shared-core.md` with z25/zm/vcm)

Three edits + one new function block in `skills/first-officer/references/first-officer-shared-core.md`. Ideation may fan out in parallel (each entity writes only its own body); implementation serializes to avoid same-file merge conflicts and because z25/zm/vcm compose into one coherent Working-Principles section — but `k74g` edits Startup + adds a new block (DIFFERENT sections from Working-Principles), so it is collision-light and lands first per the sprint order.

**Edit 1 — Startup step 3 (no forced pick):**
- before: `… multiple → present the list (or fail with an ambiguity error in single-entity mode).`
- after: `… multiple → LIST the managed workflows in the greet and proceed; do NOT ask the captain to pick one before greeting (the captain acts on a chosen workflow later via «engage»(workflow)). In single-entity mode, fail with an ambiguity error.`

**Edit 2 — Startup step 8 (interactive branch; no gate render at greet + engage hint):**
- before: `- **Interactive:** present the summary (and any ready `gate: true` gate as captain-facing text), then STOP for input; do NOT auto-dispatch. The expensive deferrals stay past the greet, reached on the captain's first direction.`
- after: `- **Interactive:** present the summary — the managed workflow(s) with their dispatchable / ready-gate counts — and hint `Use engage <workflow>` to act; then STOP for input. Do NOT auto-dispatch, and do NOT render a `present-gate` review at the greet: NAME any ready `gate: true` gate in the summary, but assemble its review only when «engage» reaches it. The expensive deferrals — gate assembly included — stay past the greet, reached on the captain's first «engage».`

**Edit 3 — Deferred load points note (line 37; present-gate loads at engage, not greet):**
- before: `A greet-and-stop boot loads NONE of these — it composes its summary from `«state.boot»` JSON + README frontmatter (Startup step 8) and presents any ready gate via `present-gate`. Each loads only at its trigger:`
- after: `A greet-and-stop boot loads NONE of these — it composes its summary from `«state.boot»` JSON + README frontmatter (Startup step 8) and NAMES any ready gate without rendering it; `present-gate` loads only when «engage» reaches a gate, not at the greet. Each loads only at its trigger:`

**Edit 4 — new boot-resident function block `## «engage»(workflow)`** (placed after the Startup section; a DIFFERENT section from Working-Principles). It wraps the EXISTING event loop — no new binary, no new dispatch mechanism (captain-resolved semantics, verbatim intent): an FO interaction verb; wraps `«dispatch.next-action»()`; current/named-workflow scope only; forward-compatible signature. Proposed block:

> **## «engage»(workflow): run the event loop to its stopping condition for one named workflow**
> - **trigger:** the captain invokes `engage`, optionally naming a workflow, after the greet. A captain-facing FO INTERACTION VERB — NOT a binary command and NOT a new dispatch mechanism; it names an entry point into logic the contract already defines.
> - **effect:** for the named `workflow` (default: the current / only managed workflow), run `«dispatch.next-action»()` (the deferred dispatch module's event-loop skeleton — loading `references/fo-dispatch-core.md` on first dispatch) to its stopping condition: dispatch each ready entity, advance each completed non-gated stage, present each ready gate via `present-gate`.
> - **scope:** ONE workflow per invocation. The `workflow` argument is present now so a future multi-workflow form (`«engage»(w1, w2, …)` or an all-managed default) EXTENDS this signature rather than replacing it — sweeping several workflows at once is a named future extension (0250 "Out of scope"), not defined here, and this block must not preclude it.
> - **done-when:** `«dispatch.next-action»()` reaches its stopping condition for the named workflow (a gate presented and awaiting the captain, terminal reached, or nothing dispatchable).
> - → **prose** — no binary backs `engage`; it wraps the existing `→ prose` `«dispatch.next-action»()` skeleton (driver binary descoped to roadmap 0222).

Headless / single-entity mode is UNAFFECTED — the no-pick and no-greet-render changes are interactive-only; headless step 8 still authors full gate reviews at its stop (that transcript is the async handoff surface a human picks up).

**Leanness (sprint-wide constraint, not a per-entity gate):** Edits 1–3 are net roughly byte-neutral; Edit 4 adds a boot-resident block. Keep it to the tight form above; the combined Startup + Working-Principles union is measured against the 0.24.0 baseline by the sprint's independent preflight staff review, not by this entity alone.

**Doc-diff determination:** the user-visible surface this entity changes — the greet text, the "Use engage" hint, and the `engage` verb — lives ENTIRELY in the FO contract prose (the skill file edited above), which IS the deliverable. No docs-site page describes the greet, the pick behavior, or FO-driving verbs (`docs/site/advanced/multi-workflow.md` says only "Spacedock finds them all and operates them at the same time," already consistent with a list-only greet — no edit needed). So the "documentation diff" for this behavior change is the contract before/after above; there is no separate `docs/` file to change.

## Spike (riskiest mechanism first) — the live-boot value-proof measurement

**Question (checklist):** can a live interactive boot be driven and its pre-greet tool-call count deterministically asserted as the value proof for AC-1?

**Exercise:** parsed n=91 real FO-boot transcripts (`~/.claude/projects/-Users-clkao-…spacedock-v1/*.jsonl`, boots identified by the bootstrap prompt in the first user message, ensign dispatches excluded); counted and categorized `tool_use` blocks per boot; located `present-gate` invocations by row index; checked `status --discover`.

**Findings:**
1. The measurement ARTIFACT exists and is machine-parseable — every `tool_use` block (name + input) is countable; `present-gate` invocations are detectable by matching `present-gate` in name/input. (Mechanism proven.)
2. A raw **pre-greet TOTAL tool-call count against a threshold is NOT a viable metric.** The greet-stop boundary is not robustly recoverable from a passive transcript: "first text-only assistant turn" fires early because FO boots interleave narration turns with tool turns (yielded false 0-counts, spread 0–17); "first human follow-up message" is confounded by Claude Code's system-injected user-role rows (attachments, system-reminders, queued ops) that collapse the boundary to ~1. And the count is LLM-nondeterministic run-to-run.
3. The heavy work products are **contract-directed but not present at the greet in the natural repo.** Across 91 boots `present-gate` appears in 49, but its FIRST invocation is ALWAYS at row ≥74 (post-direction) and NEVER within the first 30 rows (the boot/greet span) — the current FO in practice defers the render past the greet. And `status --discover` returns exactly ONE workflow here, so step 3's multi-workflow pick branch never fires naturally. So a raw count delta on the natural repo would be ≈0, not "a baseline that moved the wrong way."

**Conclusion for AC-1:** the value proof must be a LIVE interactive drive on a CONSTRUCTED heavy fixture (≥2 discoverable workflows + ≥2 ready `gate: true` gates), asserting CATEGORICAL / structural signals — present-gate-at-greet count (baseline ≥1-per-ready-gate → target 0) and the greet's terminal prose (baseline pick-question → target engage-hint, no pick) — captured at the FO's controlled greet-stop (the drive's stop IS the boundary; no fragile detection). NOT a raw total count. The baseline's heaviness is contract-DIRECTED (step 8 + the deferred-load note both direct a render); the fixture forces a compliant baseline to execute it, and the target's lightness is deterministic-by-contract. This is recorded before AC-1 is finalized, per the checklist.

## Acceptance criteria

**AC-1 (value, measured against a baseline that moves the wrong way) — the interactive greet is deterministically light.** On a heavy fixture (≥2 discoverable workflows; the chosen workflow holding ≥2 ready `gate: true` entities), a live interactive boot on the branch reaches a SINGLE greet-stop that renders ZERO `present-gate` reviews and asks NO workflow-pick question. The current-contract baseline is heavy across TWO sequential stops: step 3's multi-workflow pick short-circuits the FIRST greet-stop with a pick question, and only the resolved-workflow boot that follows it renders one `present-gate` review per ready gate. Measured quantities: the terminal greet prose changes from the baseline's first-stop workflow-pick question to a `Use engage` hint with no pick; and authored gate reviews — rendered `Gate review:` blocks, counted rather than `present-gate` LOADS (one load can author several reviews inline) — drop from ≥ready-gate-count (baseline, one per ready gate on its resolved-workflow boot) to 0 (branch, at its single greet-stop, which also loads `present-gate` zero times).
- *Tested by:* a live interactive FO drive on the same fixture under both the `origin/main` contract and the branch. Baseline (two stops): snapshot the first yield for input and assert a workflow-pick question with zero authored gate reviews there, then pick the workflow and snapshot the resolved-workflow boot and assert authored `Gate review:` blocks ≥ ready-gate count (counting rendered reviews, not `present-gate` loads — one load can author several). Branch (one stop): snapshot its single greet-stop and assert the robust conjunction — zero `present-gate` loads AND zero authored `Gate review:` blocks — corroborated by the greet prose (names-not-renders, no pick, `Use engage` hint). Parse `tool_use` for loads and the transcript for authored review blocks. (Per the spike: categorical signals at controlled greet-stops, not a raw total count.)

**AC-2 (value, behavior) — `«engage»(workflow)` runs the existing loop to its stopping condition with no new mechanism.** A live `«engage»(<workflow>)` after the greet is observed dispatching a ready entity / advancing a completed non-gated stage / presenting a ready gate via `present-gate` — the actions `«dispatch.next-action»()` produces — and stops at the loop's stopping condition (a gate awaiting the captain, terminal, or nothing dispatchable). No binary `engage` command and no new dispatch code exist. Headless / single-entity mode is unchanged.
- *Tested by:* the same live drive: invoke `«engage»(<workflow>)` on the fixture, observe the loop actions and the stop; confirm a headless boot on the same fixture still drives per step 8's headless rule (unchanged).

**AC-3 (mechanism, serves AC-2's "engage is the real verb" value) — the launcher flourish is gone without breaking the launch.** `bootstrapPrompt` and `codexBootstrapPrompt` no longer end with "Engage."; a launched FO still boots into role and the prompt is still appended last (and still suppressed on `--resume`). `piBootstrapPrompt` is untouched.
- *Tested by:* `internal/cli/safehouse_frontdoor_test.go` — the updated `want*BootstrapPrompt` oracles assert the launched inner argv's last token equals the flourish-free string across the safehouse / non-safehouse / resume-suppression / codex cases; neither constant ends with "Engage." (This mechanism serves the value that `engage` is a real captain verb, not a launch flourish that collides with it.)

## Test plan

- **AC-3 (Wave 0 — cheap, deterministic, minutes):** Go unit test only. Update the two `want*BootstrapPrompt` constants; the existing argv oracles already assert the appended last token, so the diff is the constants + a "does not end with Engage." assertion. `go test ./internal/cli/`. No fixture, no live run.
- **AC-1 + AC-2 (Wave 1 — live, moderate; the main cost is fixture setup):** a LIVE interactive FO drive is required by this workflow's proof policy (behavioral contract claims are proven by a live drive, never a prose-grep). Build a fixture repo with ≥2 discoverable workflows, one holding ≥2 ready `gate: true` entities. Baseline drive on `origin/main`, target drive on the branch; capture each greet-stop transcript and assert AC-1's categorical signals; then invoke `«engage»(<workflow>)` and observe AC-2's loop actions + stop. Human-observed or PTY-scripted (the interactive greet stops for input, so it is not a headless `-p` run).
- **Spike status:** DONE this stage (above) — the passive-transcript raw count is ruled out; the categorical controlled-greet-stop signal on a heavy fixture is the metric AC-1 uses.
- **Sequencing / coordination:** Wave 0 lands before `7v` (or 7v re-derives from the live constant); Wave 1's shared-core edit is strict-serial within the sprint and precedes z25/zm/vcm's Working-Principles edits per the 0250 order.

## Stage Report: ideation

- DONE: Write Problem/Proposed approach/Acceptance criteria/Test plan incorporating the captain-resolved engage semantics verbatim (interaction verb wrapping the existing dispatch.next-action loop, current-workflow-only scope, forward-compatible signature for a future multi-workflow call)
  Body sections added; the `«engage»(workflow)` block states interaction-verb / wraps `«dispatch.next-action»()` / one-workflow scope / forward-compatible signature verbatim; settled semantics carried, not re-litigated.
- DONE: Spike the riskiest mechanism first per proof policy: can a live interactive boot be driven and its pre-greet tool-call count deterministically asserted as the value proof; record the result before finalizing AC-1
  Parsed n=91 real boot transcripts: raw pre-greet total count is NOT viable (boundary un-recoverable from a passive log + LLM-nondeterministic); present-gate never in first 30 rows and only 1 workflow discoverable → baseline is contract-directed, needs a heavy fixture; AC-1 uses the categorical present-gate-at-greet + pick-question signals at a controlled greet-stop. Recorded in the "## Spike" section before AC-1.
- DONE: Scope the frontdoor.go bootstrap-prompt removal (Wave 0, code, parallel-safe) separately from the Startup step 3/8 contract-prose change (Wave 1) per the 0250 sprint doc's own sequencing, and note the 7v coordination
  Proposed approach split into Wave 0 (frontdoor.go const + oracle constants, Go-only, first, parallel-safe) and Wave 1 (Startup step 3+8 + deferred-load note + new `«engage»` block, strict-serial in the sprint); 7v coordination noted (Wave 0 before 7v, or 7v re-derives from the live constant).

### Summary

Fleshed out the ideation body for the light-greet + `engage`-verb entity: Problem (3 heavy/ceremonial boot behaviors), a two-wave approach with specific contract before/after wording and the new `«engage»(workflow)` block, three ACs, and a test plan. Riskiest-first spike (recorded before AC-1): a raw pre-greet total tool-call count is not a deterministic value proof — the greet-stop boundary is un-recoverable from a passive transcript and the count is LLM-nondeterministic; the robust metric is the categorical present-gate-at-greet count + workflow-pick-question signal, captured at a controlled greet-stop on a constructed heavy fixture (≥2 workflows + ≥2 ready gates), since the natural repo boot is already near-light (single workflow, present-gate deferred past the greet). No separate docs-site diff — the changed user-visible surface is the FO contract prose itself.

## Stage Report: implementation

- DONE: Wave 0 (Go): bootstrapPrompt and codexBootstrapPrompt no longer end with ' Engage.'; wantBootstrapPrompt/wantCodexBootstrapPrompt oracles mirrored; piBootstrapPrompt untouched; go test ./internal/cli/ green (AC-3)
  frontdoor.go: both constants drop the flourish; test want* oracles mirrored + new TestBootstrapPromptsDropEngageFlourish (asserts neither carries "Engage."); pi.go piBootstrapPrompt has no flourish, untouched; `go test ./internal/cli/` ok. TDD: assertion failed pre-edit, green post-edit. Commit 9f869521.
- DONE: Wave 1: the four shared-core edits land per the body's before/after — Startup step 3 (list managed workflows, no forced pick), step 8 interactive (NAME ready gates, no present-gate render, 'Use engage' hint), the deferred-load-points note, and the new '## «engage»(workflow)' block placed after Startup — AND the block's scope bullet is tightened with its forward-compat meaning fully preserved
  All four edits applied verbatim to skills/first-officer/references/first-officer-shared-core.md; «engage» block sits between Startup and Deferred load points, references the real «dispatch.next-action»() loop (fo-dispatch-core.md); scope bullet trimmed 142 bytes (359→217) keeping one-per-invocation + arg-EXTENDS-not-replaces + 0250-out-of-scope-not-precluded. Commit 5a304b68.
- DONE: Leanness on the record: wc -c of first-officer-shared-core.md before and after with the byte delta (this member's share of the sprint's ~5,600-byte combined resident ceiling)
  before 21663 bytes → after 23435 bytes, delta +1772 (edits 1–3 net a small add from the engage-hint prose; the «engage» block is the bulk). Well under the sprint's ~5,600-byte combined ceiling.

### Summary

Wave 0 (Go, TDD): dropped the " Engage." launch flourish from both launcher bootstrap constants, mirrored the want* oracles, added an explicit AC-3 assertion; `go test ./internal/cli/` green; piBootstrapPrompt untouched. Wave 1 (contract prose): landed the four first-officer-shared-core.md edits (step 3 no-pick + list, step 8 no gate render + engage hint, deferred-load note, new «engage»(workflow) block after Startup) and tightened the block's scope bullet by 142 bytes with the forward-compat meaning intact. Leanness: +1772 bytes this member, within the sprint's ~5,600-byte ceiling. AC-1/AC-2's live heavy-fixture drive is the fresh validation stage's job, not exercised here.

## Stage Report: validation

- DONE: AC-1+AC-2 live heavy-fixture drive — construct the fixture, drive interactive baseline+target greets, assert baseline present-gate-at-greet >= ready-gate count while target == 0 and baseline asks a pick while target shows the engage hint, then engage + a headless boot (unchanged)
  Fixture: a git repo with 2 discoverable workflows (wf-alpha, wf-beta); wf-alpha holds 2 ready gates (001/002-alpha at `review`, `next-suppressed-by=gate`) + 1 dispatchable (003-alpha, backlog→review). Built A/B skill trees differing ONLY in `first-officer-shared-core.md` (baseline 21663 vs target 23435; `diff -rq` reports that one file). Four live FO drives (general-purpose agents handed the exact contract file, interactive-class per the -p/exec rule); each drive's own `subagents/*.jsonl` parsed for `Skill(present-gate)` tool_use, not the agent's self-report. Results: target multi-greet → present-gate 0, engage hint, NO pick; baseline multi-greet → PICK question ("cannot boot until you pick one"); baseline post-pick wf-alpha boot → present-gate 2 (== ready-gate count); `engage wf-alpha` → loop presented both gates via present-gate and STOPPED at "a gate awaiting the captain", recorded the 003 dispatch, no binary / no new mechanism; headless wf-alpha → drove the loop, authored 2 full present-gate reviews at the stop, did not resolve (step-8 HEADLESS branch unchanged). AC-1 signal caveat below.
- DONE: AC-3 reproduced independently — go test ./internal/cli/ green in the worktree; neither bootstrapPrompt nor codexBootstrapPrompt ends with 'Engage.'; piBootstrapPrompt untouched
  `go test ./internal/cli/ -count=1` green (uncached, 51.7s). frontdoor.go:25 `bootstrapPrompt` + :533 `codexBootstrapPrompt` both end flourish-free; pi.go:20 `piBootstrapPrompt` unchanged (`git diff fa14da5f..HEAD -- internal/cli/pi.go` empty). Non-self-referential: a detached adversarial audit on a throwaway worktree re-introduced " Engage." into the PRODUCTION constant → the argv oracles (`TestClaudeNoSafehouseLaunchesPlain`, `TestClaudeSafehousePresentWrapsArgv`, which observe the real launched inner argv) AND `TestBootstrapPromptsDropEngageFlourish` all went RED, proving the checks catch the regression rather than tautologically pass.
- DONE: Every AC-N verdict carries reproduced, non-self-referential evidence; leanness re-measured on the record (21663 -> 23435, +1772); explicit PASSED/REJECTED recommendation
  AC-1 PASS, AC-2 PASS, AC-3 PASS — proven by live observed behavior / running argv oracles / code-structure search; none by prose-grep over the contract. Leanness re-measured from git: origin/main 21663 → branch 23435, delta +1772 (matches the claim; branch parent fa14da5f is also 21663, so the delta is this branch's). AC-2 structural: no `engage` command in `cmd/` or `internal/` (every "engage" hit is English in test comments or the AC-3 flourish test); branch diff is 3 files only (frontdoor.go, safehouse_frontdoor_test.go, first-officer-shared-core.md), no new dispatch code.

### Summary

Recommendation: PASSED. All three ACs are proven with live, transcript-verified, non-self-referential evidence: the target interactive greet is deterministically light (0 present-gate, no workflow-pick, `Use engage` hint), `«engage»(wf-alpha)` runs the existing `«dispatch.next-action»()` loop to its gate-awaiting-captain stop with no binary and no new dispatch code, headless is unchanged (2 full gate reviews authored at the stop), the launch flourish is gone (adversarial-audit-confirmed), and leanness is +1772. One material finding on AC-1's MEASUREMENT WORDING (non-blocking): AC-1 asserts baseline present-gate >= ready-gate-count AND the pick-question both "at the greet-stop" on one >=2-workflow fixture, but the baseline's step-3 pick short-circuits the boot BEFORE step-8's gate render — so at the baseline's FIRST greet-stop present-gate is 0, and the render (2) appears only at the baseline's POST-pick wf-alpha boot. The two baseline heavy behaviors are sequential, not simultaneous. This is an imprecision in the AC's phrasing, not a defect in the deliverable: both categorical signals reproduce (baseline: pick, then 2 renders once a workflow is resolved; target: no pick, 0 renders) and the value — a lighter greet — is proven. The code and contract are correct, so nothing routes to implementation; recommend the captain tighten AC-1's wording (measure the present-gate signal on the baseline's resolved-workflow boot, or on a single-workflow fixture) rather than bounce the deliverable.
