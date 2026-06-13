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

## Acceptance criteria

(Ideation firms. Verified by the survey's rendered OUTPUT on a constructed fixture — never a prose-grep of the skill.)

**AC-1 (sketch) — the report does not present its recent window as the project's whole history.** Verified by: a survey run over a fixture whose early-window and late-window sessions differ in mode/workstream renders either a phase/evolution signal or an explicit recent-window framing — observed in rendered output.

**AC-2 (sketch) — the report carries a single honest partial-lens statement.** Verified by: a survey run over a fixture with off-corpus gaps (e.g. blank-cwd sessions and/or dispatched subagents) renders one clear "this reflects the agent-log corpus; not captured: …" statement, not only scattered technical asides — observed in rendered output, with the named gaps derived from the fixture's session rows (independent source), not skill prose.

## Test plan

(Ideation/implementation firms.) Fixture-driven renders: an early-vs-late-window-divergent fixture for AC-1; an off-corpus-gap fixture for AC-2. Per the survey discipline, a grep over SKILL.md never satisfies the behavioral AC.

## Notes

- **Validation to preserve:** a second author confirmed the recent-workflow inference is accurate (nothing surprising/untrue). The keep-signal (`zby` #3 / `zwg` #0) now has two-author support.
- **Relation to `za`:** `za` covers the subagent-dispatch subset of the corpus-incompleteness; this seed is the broader honest-lens framing (#2) plus the temporal axis (#1).
- **Cross-cutting (feeds 0.21.x):** "the agent log is a partial lens" applies equally to `orient` and the cross-workflow ready-room, which are also reconstructed from logs. The honest-partial-lens principle is a property of the whole log-derived decision layer, not just the survey.
