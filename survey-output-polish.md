---
id: h5cn1nvzxed69454a6pqk2tx
title: Survey output polish — Codex workstream null-signal + scratch-preamble leak (2nd-round runs)
status: backlog
source: "Captain real-world survey runs, 2026-06-09/10 — a second round over two repos (a public adapter-maintenance repo + a private code repo). Records only the spacedock:survey behavioral feedback; corpus content is deliberately omitted (the private repo is pre-OSS), per the survey-feedback anonymization discipline."
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0202-survey-improvements
group: survey
sprint-readiness: ready
---

Two concrete output defects observed across the second round of real `spacedock:survey` runs, plus cross-references where those runs refine the sibling survey-polish seeds. Corpus specifics stay out of this shared repo — only behavior.

## Feedback (the asks)

**1. The Codex workstream classification is a null signal.** In both runs the Codex `workstreams:` line read `(unlabeled) — N` (one run `(unlabeled) — 2`, the other `(unlabeled) — 17`), i.e. it classified zero Codex workstreams in either repo while reporting high activity (one run: `exec_command 1255 · update_plan 9 · spawn_agent 0`). The anti-over-promise "unlabeled" floor is working (it declines to invent labels), but the Codex `workstreams:` line then carries no information — it presents a breakdown that is always empty. Either give Codex workstream classification real signal (e.g. infer from exec_command/update_plan/branch patterns) or stop presenting it as a breakdown when it would be uniformly unlabeled (collapse to a single honest "N ad-hoc shell-driven sessions, unclassified" line).

**2. Model scratch-reasoning leaks into the user-facing report.** One run's output opened with the agent's pre-amble — `⏺ I have everything I need. Let me cross-check the OPEN forks against the repo before presenting:` followed by its reasoning bullets — before the actual report. The cross-check *content* is valuable (it's the decided/shipped/superseded reclassification — see Notes), but the "I have everything I need / let me…" framing is scratch that should not appear in the rendered survey.

## Out of scope

- The decision-frontier triage logic itself (NEEDS YOU / BACKLOG / RECENT DECISIONS + the open-fork-vs-repo cross-check) — that is a strength to preserve and extend, not polish away (see Notes).
- Reproducing the surveyed corpus content (private; not in this repo).

## Acceptance criteria

(Ideation firms. Verified by the survey's rendered OUTPUT on a constructed fixture — never a prose-grep of the skill.)

**AC-1 (sketch) — the Codex workstream line never presents an empty breakdown as a breakdown.** Verified by: a survey run over a fixture whose Codex sessions are all unclassifiable renders a single honest unclassified line (or a real classification), not a `workstreams: (unlabeled) — N` breakdown row — observed in rendered output.

**AC-2 (sketch) — no model scratch-reasoning precedes the report.** Verified by: a survey run over a fixture renders the report without an `I have everything I need` / `let me…` pre-amble — observed in rendered output (the cross-check's findings still appear, but as report content, not scratch).

## Test plan

(Ideation/implementation firms.) Fixture-driven survey renders over constructed session sets (an all-unclassifiable-Codex fixture; a fixture that triggers the cross-check). Per the survey discipline, a grep over SKILL.md never satisfies the behavioral AC.

## Notes — second-round refinements to sibling survey-polish seeds

The same two runs refine the cluster (carry these into each seed's ideation; not yet folded into their bodies):

- **`zby` #1 (lead with the workdir-attributed Codex count) appears ALREADY SHIPPED.** Both runs led with the workdir-attributed count and demoted name-match to a parenthetical caveat. Verify against the current skill before implementing — `zby` may narrow to just #2 (knowledge-work archetype).
- **`zwg` #1 (iteration vs gate framing) → narrow to MODE-AWARE.** Both repos were overwhelmingly mechanical (issue→worktree→PR loops, no exploration tracks), where the gate-drive offer framing is apt; one run already segmented mechanical→gated vs unlabeled→book-keeping. So the iteration/steering reframe is exploration-mode-specific — do NOT strip gate language from mechanical offers.
- **`zwg` #2 (branch-and-merge work-by-area misattribution) did NOT bite either run** — product code led work-by-area in both. Confirms the misattribution is repo/workflow-specific; the fix is detect-branch-and-merge + caveat (conditional), lower blast-radius than first filed.
- **`za` (report the fact of subagent dispatch) was not exercised** — neither repo was spacedock-orchestrated (one run showed Codex `spawn_agent 0`). Still valid for orchestrated repos.

**Strategic note (feeds 0.21.x, not this seed):** the runs' NEEDS YOU / BACKLOG / RECENT DECISIONS triage + the open-fork-vs-repo cross-check (reclassify decided / shipped / superseded) is a working single-repo prototype of the cross-workflow decision-frontier / ready-room. Preserve and generalize it; it is the 0.21.x decision-abstraction wedge in nascent form.
