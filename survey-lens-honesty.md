---
id: 5xwch4ddhrzpyfbmsjgdsjy9
title: Survey owns its lens — recent-window snapshot + partial agent-log corpus
status: ideation
source: "Captain relayed author feedback, 2026-06-10, after asking an author whether their survey told them anything surprising or untrue. Author: nothing surprising/untrue, it reflected recent workflows well — BUT their relationship with the projects had changed a lot over time, and their agent conversations don't capture all their working context. Anonymized; corpus omitted."
started: 2026-06-13T04:49:28Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0202-survey-improvements
group: survey
sprint-readiness: ready
---

An author reviewing their own survey output confirmed it reflected recent workflows **well** and surfaced nothing surprising or untrue — strong validation of the core inference (keep, do not regress). They named two honest limits of the lens, neither fully owned today.

## Feedback (the asks)

**1. Surface evolution, or own the snapshot. ("my relationship with each of them has changed a ton over time")** The survey reports a `{date range}` but does no trajectory analysis — it is a flat snapshot of a recent window, so a project that matured or pivoted reads only as its recent state. The decided→shipped→superseded reclassification (a pivot abandoning an older branch) is a weak temporal seed but not surfaced as an arc. Either surface workflow evolution / phase-shifts where the signal exists (e.g. earlier-window vs later-window mode or workstream shift), or — at minimum — frame the report explicitly as "recent-window snapshot, not the project's history" so the reader isn't misled into reading it as the whole story.

**2. An honest partial-lens caveat. ("my claude/codex conversations don't capture all of that context for whatever reason")** The survey reflects only what is in the agent-session corpus, and that corpus is incomplete. Today's caveats are scattered and technical — `blank_cwd` sessions, the Codex name-vs-workdir sibling gap, "Claude-only for now" — and don't add up to a clear, user-facing "this is a partial lens; here is what it can't see." The omissions include: work done outside agent sessions entirely (manual edits, other tools, off-log thinking and discussion), the pre-window past, no-cwd sessions, and (per `za`) subagent-dispatched work. The survey should state its lens honestly so its accuracy-on-what-it-sees is not mistaken for completeness.

## Out of scope

- Folding subagent CONTENT into the body — `za` already scopes the subagent case to the FACT of dispatch, not content.
- Reproducing surveyed corpus content (private/omitted).

## Chosen direction (ideation, 2026-06-13)

Both honest-lens outputs render into the captain-locked **value-&-numbers-first** spine
(`0202-survey-improvements/index.md`), specifically the mock's two-part **parenthetical subhead**
directly under the title:

```
SpaceDock survey — your last 30 days
(recent-window snapshot · agent logs only · 12 sessions had no working dir, not placed)
```

- **The subhead's first clause `(recent-window snapshot …)` is AC-1's framing.** The survey already
  reads a `{date range}` span (the `scoping` query returns `span`). The fix is to STOP letting that
  span read as the project's whole history: the title carries the window (`your last {N} days`) and the
  subhead's first clause states plainly that this is a *recent-window snapshot, not the project's
  history*. No trajectory/evolution analysis is built (see decision below) — the snapshot FRAMING is
  the AC-1 deliverable, and it is the arm the captain already chose in the mock.
- **The subhead's remaining clauses `· agent logs only · {blank_cwd} … no working dir` are AC-2's
  partial-lens statement,** plus the closing `↓ full analysis: … what we can't see` pointer to a
  consolidated "what this can't see" detail block. Today's caveats are scattered (`blank_cwd` asides,
  the Codex name-vs-workdir gap line, "Claude-only for now") and never add up to one user-facing lens
  statement. The fix is ONE consolidated statement: a top-line `agent logs only` + the countable gap
  (`{blank_cwd} … no working dir`) in the subhead, expanded once in the detail block with the full
  not-captured list (work outside agent sessions, the pre-window past, no-cwd sessions, and — via `za`
  — dispatched-subagent work). Scattered technical asides collapse INTO this one statement.

The two outputs share one rendered surface (the subhead) so the lens is stated once, honestly, at the
top — not reconstructed by the reader from technical asides scattered through the body.

## Riskiest-path determination

**No spike needed.** Both outputs compose already-proven reads:

- AC-1's snapshot framing reuses the `span` (`MIN(started_at) .. MAX(ended_at)`) already returned by
  the `scoping` query — the framing is a render-side label on data the survey already has.
- AC-2's `blank_cwd` count is already returned by `scoping`; the Codex name-vs-workdir superset gap is
  already computed (`codex-presence` vs `codex-scoped`); the dispatched-subagent gap is `za`'s
  dispatch-fact, not new here.

The one path that WOULD be unverified — surfacing an early-vs-late-window **mode/workstream shift**
(AC-1's "evolution" arm) — needs a NEW windowed query (`mode-classification` aggregates the whole
window with no timestamp split). The captain-locked mock chose the snapshot-framing arm and demoted
evolution to nothing, so this seed does NOT build the windowed query (YAGNI). The evolution arm is
recorded as a deferred option (see Notes) should a future fixture show the signal is worth the new
query. With evolution out, nothing here is unverified — hence no spike.

## Acceptance criteria

(Verified by the survey's rendered OUTPUT on a constructed fixture — never a prose-grep of the skill.)

**AC-1 — the report frames its recent window as a snapshot, not the project's whole history.**
Verified by: a survey run over a fixture whose sessions span a date range renders, directly under the
title, an explicit recent-window-snapshot framing (the window in the title + the subhead's
`recent-window snapshot` clause) — observed in rendered output, with the window derived from the
fixture's session timestamps (`span`, independent source), not skill prose.

**AC-2 — the report carries a single consolidated partial-lens statement.** Verified by: a survey run
over a fixture with off-corpus gaps (blank-cwd sessions and a name-only Codex superset) renders ONE
"agent logs only … not captured: …" statement (in the subhead + the detail block), not only scattered
technical asides — observed in rendered output, with the named/counted gaps derived from the fixture's
session rows (`blank_cwd`, the `codex-presence` − `codex-scoped` gap — independent sources), not skill
prose.

## Doc diff (SKILL.md — applied by implementation)

The change is to the **step-4 report template** in `skills/survey/SKILL.md`. Two edits:

**Edit 1 — title + subhead carry the snapshot framing and the lens statement (new, above the
`PROJECT:` fence header).** Today the body opens with the synthesis-fence headline
(`Found {N} sessions … here's the lay of the land:`) and the `PROJECT:` line. Per the captain-locked
spine the rendered report opens with the title + parenthetical subhead:

```
SpaceDock survey — your last {N} days
(recent-window snapshot · agent logs only{if blank_cwd>0: · {blank_cwd} sessions had no working dir, not placed})
```

— `{N} days` derives from the `scoping` `span`; `{blank_cwd}` from `scoping`. The subhead is dropped to
just `(recent-window snapshot · agent logs only)` when `blank_cwd=0` (drop-empty-slot rule).

**Edit 2 — a consolidated "what this can't see" detail line replaces the scattered asides.** The
existing `{if blank_cwd>0: {blank_cwd} uncaptured-cwd sessions}` line under `PROJECT:` and the loose
Codex name-only-superset aside fold into ONE closing pointer + detail block:

```
  ↓ full analysis: modes, work-by-area, what this can't see

WHAT THIS CAN'T SEE   (the agent-log corpus is a partial lens)
  · work done outside agent sessions (manual edits, other tools, off-log discussion)
  · the project's history before this {N}-day window
  {if blank_cwd>0: · {blank_cwd} sessions with no recorded working dir}
  {if codex-presence>codex-scoped: · {codex-presence − codex-scoped} Codex sessions matched by name only (possible same-named sibling)}
  {if any dispatched subagents (za): · work inside dispatched subagents (orchestration counted, content not — see "{N} subagents dispatched")}
```

Each `·` line is dropped when its slot is empty (drop-empty-slot rule). The `WHAT THIS CAN'T SEE`
block is the single consolidated home for the lens caveats; the prior scattered asides
(`uncaptured-cwd sessions` under `PROJECT:`, the inline Codex sibling caveat) are REMOVED so the lens
is stated once. This composes with `za` (the dispatched-subagent line) and `9h` (which owns the
title/lede and the `↓ full analysis` pointer) — those seeds render into the SAME subhead/block, so
ideation coordinates the wording rather than three seeds each editing the header independently.

## Test plan

Fixture-driven survey renders against a constructed agentsview DB (Go fixture / golden render — the
established survey-feedback proof shape; a grep over SKILL.md never satisfies a behavioral AC). Two
fixtures:

- **AC-1 fixture — date-spanning sessions.** A fixture whose Claude session rows span a clear date
  range (earliest `started_at` to latest `ended_at` ⇒ a non-trivial `{N}` days). AC: the rendered
  output carries the recent-window-snapshot framing (window in title + `recent-window snapshot`
  subhead clause), with `{N}` matching the fixture's computed span — the expected `{N}` comes from the
  fixture rows, an independent source, NOT from skill prose. (Reusing the divergent early/late-window
  rows from the dropped evolution arm is fine as fixture shape, but the assertion is on the snapshot
  framing + span, not on a mode shift.)
- **AC-2 fixture — off-corpus gaps.** A fixture with (a) ≥1 blank-cwd Claude session and (b) a Codex
  presence count exceeding the workdir-scoped count (name-only superset). AC: the rendered output
  carries ONE consolidated `WHAT THIS CAN'T SEE` statement naming `{blank_cwd}` and the
  `{codex-presence − codex-scoped}` name-only gap, with those counts matching the fixture rows
  (independent sources) — and NOT a second scattered copy of those caveats elsewhere in the body.

Cost/complexity: low — both reuse the existing survey golden-fixture harness and existing query
outputs; no new query, no live-drive required (snapshot framing + lens statement are render-side
labels on already-read data). A live smoke (run the survey on this repo) is a cheap confidence add but
not the AC proof.

## Notes

- **Validation to preserve:** a second author confirmed the recent-workflow inference is accurate (nothing surprising/untrue). The keep-signal (`zby` #3 / `zwg` #0) now has two-author support.
- **Relation to `za`:** `za` covers the subagent-dispatch subset of the corpus-incompleteness; this seed is the broader honest-lens framing (#2) plus the temporal axis (#1).
- **Cross-cutting (feeds 0.21.x):** "the agent log is a partial lens" applies equally to `orient` and the cross-workflow ready-room, which are also reconstructed from logs. The honest-partial-lens principle is a property of the whole log-derived decision layer, not just the survey.
- **Deferred — the evolution arm (AC-1's other OR branch).** Surfacing an early-vs-late-window mode/workstream SHIFT (not just the snapshot framing) needs a new windowed query (`mode-classification` has no timestamp split). The captain-locked mock chose the snapshot-framing arm, so this seed does not build it (YAGNI). If a future fixture/run shows a project's trajectory is worth surfacing, file a follow-up that adds the windowed split — the data (per-session `started_at`/`ended_at`) is present; only the query+render is missing.

## Stage Report: ideation

- DONE: Firm the two honest-lens outputs against the banked structure
  Both outputs render into the captain-locked subhead (`Chosen direction`): clause 1 = AC-1 snapshot framing (from `scoping.span`); clauses 2+ = AC-2 consolidated partial-lens statement (`agent logs only` + `blank_cwd` + name-only Codex gap), with gaps derived from fixture session rows per AC-2.
- DONE: Record "no spike needed" naming the proven mechanisms
  `Riskiest-path determination`: snapshot framing reuses `scoping.span`; partial-lens reuses `scoping.blank_cwd` + `codex-presence`−`codex-scoped` + `za`'s dispatch fact. The only unverified path (windowed mode-shift) is the dropped evolution arm — deferred, not built — so nothing here is unverified.
- DONE: Propose the doc diff and a test plan (AC-1 date-span fixture + AC-2 off-corpus-gap fixture; never a prose-grep)
  `Doc diff` section gives two concrete step-4 SKILL.md edits (title+subhead, consolidated `WHAT THIS CAN'T SEE` block replacing scattered asides). `Test plan` gives both fixtures with the expected values bound to independent fixture-row sources, on the survey golden-fixture harness.

### Summary

Firmed AC-1 (recent-window snapshot framing) and AC-2 (one consolidated partial-lens statement) to render into the captain-locked value-&-numbers-first subhead, not a new surface. Key decision: take the snapshot-framing arm of AC-1's OR (per the mock) and DROP the early-vs-late evolution arm, which alone would need a new windowed query — recorded as deferred, so "no spike needed" holds. Doc diff coordinates with `za` (dispatched-subagent line) and `9h` (title/`↓ full analysis` pointer), which render into the same subhead/block.
