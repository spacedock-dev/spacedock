---
id: zbysa49erepjs7442m8a27mv
title: Survey skill — name a knowledge-work archetype + lead with workdir-attributed Codex count
status: backlog
source: "Captain real-world survey run on a private knowledge-work repo, 2026-06-09. This records only the spacedock:survey behavioral feedback — the surveyed corpus content (areas, skill names, the user's open forks, counts) is deliberately NOT reproduced here, since the state repo is shared via origin."
started:
completed:
verdict:
score:
worktree:
issue:
sprint:
group:
sprint-readiness:
---

Real-world feedback on `spacedock:survey` from a run on a **knowledge-work** repo (a notes/memo/writing/ops shop, not an issue→PR code repo). The workflow-inference was genuinely good; two gaps surfaced because the repo is knowledge work, not code. The surveyed corpus is private and intentionally not captured here — only the survey-skill behavior to improve. Builds on `69` (survey-codex-cwd-workaround, done — workdir attribution exists) and `47rx`/`survey-codex-and-sandbox-followups` (done — the G anti-tautology "unlabeled" floor).

## Feedback (the asks)

**1. Lead with the workdir-attributed Codex count, not the name-matched one (presentation).** On this run, matching Codex sessions by project NAME yielded an order-of-magnitude more "hits" than matching by actual working directory — the surplus was a same-basename-repo collision, not the user's sessions. The survey *does* explain the distinction, but the inflated name-match number is the first thing a user reads, and it's exactly the kind of figure misread at a glance. Surface the workdir-attributed (sibling-free) count FIRST/loudest; demote the name-match to an explicit caveat. The data is already computed (`69`); the framing buries the precise number under the loose one.

**2. Name a knowledge-work archetype (classifier).** Every workstream classified **unlabeled** — the most interesting finding. A knowledge-work repo has no mechanical issue→PR signature (so no gate-and-drive automation to offer) and no veto-heavy creative signature, so the classifier correctly DECLINES to over-promise (`47rx`'s floor working as intended). But the report then ends flat at "generic book-keeping." There is a real third archetype: a **knowledge-work loop** — intake → process → file → log → close, where the gates are "confirm this batch / approve this write / scope this draft." The survey should RECOGNIZE and NAME this mode (a third option beside mechanical-invocation-and-drive and exploration-and-steering) and frame an honest offer for it, instead of falling through to a generic fallback. Not a new automation pitch — a named mode that gives the anti-over-promise floor a positive landing.

**3. (Keep — do not regress.) The workflow inference is genuinely good.** The survey reconstructed the intake→process→file→log→close loop from the decision points alone, and the per-area edit profile was a sharp identity signal (it correctly read the repo as a knowledge-work shop, not a software project). This part sells itself — preserve it.

## Proposed direction (ideation fills in)

- (#1) In the report's Codex section, lead with the workdir-attributed count; render the name-match only as a labeled caveat. The first number the user sees must be the sibling-free one.
- (#2) Add a knowledge-work archetype to the mode classifier. Detection signal: a process/file/log loop + a content/ops edit profile + no issue→PR signature + no veto-heavy creative signature. On a match, name the mode and frame book-keeping as a recognized knowledge-work offer, not a generic fallback. The existing anti-over-promise floor stays.

## Out of scope

- Re-opening `69`'s attribution mechanism (the workdir data is correct; this is framing).
- Any automation pitch for knowledge-work repos beyond honest, opt-in book-keeping.
- Reproducing the surveyed corpus content (private; not in this repo).

## Acceptance criteria

(Ideation fills in. Each verified by the survey's rendered OUTPUT on a constructed knowledge-work FIXTURE — never the user's private repo, and never a prose-grep of the skill.)

**AC-1 (sketch) — the Codex section leads with the workdir-attributed count.** Verified by: a survey run over a fixture with a deliberate same-basename collision renders the workdir count first and the name-match as a caveat — checked against the rendered report.

**AC-2 (sketch) — a knowledge-work repo is named, not left generic.** Verified by: a survey run over a knowledge-work fixture (a process/file/log loop, no issue→PR signature) classifies + names the knowledge-work archetype and frames an honest offer, rather than "generic book-keeping" — observed in the rendered output.

## Notes

Related: `69` survey-codex-cwd-workaround (done — attribution data), `47rx` survey-codex-and-sandbox-followups (done — the G unlabeled floor this gives a positive landing), `za` survey-report-subagent-dispatch-fact (backlog — sibling follow-up). Provenance is a real user run on a private knowledge-work repo, so the archetype gap is a real finding, not hypothetical — but the corpus specifics stay out of this shared repo.
