---
id: zwgx5e0beh3kv1n02g8vwamm
title: Survey skill — reframe the commission offer toward iteration/exploration + branch-merge-aware work attribution
status: backlog
source: "Captain real-world survey run on a private code repo (a worktree → PR → merge branch-and-merge workflow, superpowers scaffold), 2026-06-09. Records only the spacedock:survey behavioral feedback — the surveyed corpus (project, collaborators, workstreams, the user's open forks, PR/issue numbers) is deliberately omitted, since the state repo is shared via origin."
started:
completed:
verdict:
score:
worktree:
issue:
sprint: 0202-survey-improvements
group: survey
sprint-readiness: defer
---

Second real-world feedback round on `spacedock:survey`, this time from a **code** repo (branch-and-merge: file/triage issue → worktree branch → implement → test → PR → merge, on a superpowers scaffold). The awareness improvements are landing; two sharp critiques remain. Corpus content is private and intentionally not captured here — only the survey-skill behavior to improve. Companion to `zbysa` (survey-knowledge-work-archetype).

## Feedback

**0. (Keep — the awareness direction works.)** The user's reaction: the survey "feels more aware of me, so I feel more heard." The per-user framing + the attribution precision are landing — preserve that direction; the rest of this is refinement, not reversal.

**1. The commission offer's "explicit approval gate" framing doesn't link to an iteration/exploration mental model.** The survey pitches spacedock's value as "turn the decision points into explicit approval gates — the agent advances and stops only where you marked a gate." The user's reaction: *"explicit approval still doesn't fully link with my mental model of iteration and exploration."* Even on a repo with a clear mechanical drive signature (worktree → PR → merge), the gate metaphor reads as transactional approve/veto and misses how the work actually feels — iterative and exploratory, with the human steering a loop, not stamping gates. The offer should bridge to iteration/exploration: the agent iterates/explores and you steer; an approval gate is ONE shape that loop can take, not the whole model. Reframe so the pitch speaks to iterative/exploratory work rather than leading with gates.

**2. Work-by-area edit-counts misattribute under a branch-and-merge workflow.** The survey's work-by-area note concluded that the scaffolding/config directory had far more edits than the product-code directory and inferred "you're spending session energy on scaffolding, not product code — the machinery is a significant workstream." The user: *"This is mistaken — we know that I'm working in branches and then merging into main."* The edit count over-weights directly-edited scaffolding on the working branch and under-counts product code that lands via **PR merges** from feature branches. So the "scaffolding > product" signal is FALSE for a branch-and-merge workflow. The survey should compute work-by-area in a branch-aware way (merged-PR diffstats / across-branch edits), not raw single-branch or working-dir edit counts — or, at minimum, detect a branch-and-merge workflow and caveat/correct the edit-count signal so it never falsely reads "scaffolding over product."

## Proposed direction (ideation fills in)

- (#1) Reframe the commission/discovery-bridge offer so the lead is iteration + steering, not "explicit approval gates." Gates are a tool inside the loop; the pitch should connect to how an iterative/exploratory operator actually experiences the work. (Touches the survey's commission offer and likely the broader spacedock positioning — gates are one expression of steering, not the headline.)
- (#2) Make work-by-area branch-aware: when the inferred workflow is branch-and-merge (worktree → PR → merge), attribute product work from merged-PR diffs rather than raw edits on the surveyed branch, so directly-edited scaffolding does not dominate the picture. Failing a full fix, caveat the signal explicitly under a detected branch-and-merge workflow.

## Out of scope

- A redesign of spacedock's gate mechanism itself (this is about how the survey FRAMES and OFFERS it).
- Reproducing the surveyed corpus content (private; not in this repo).

## Acceptance criteria

(Ideation fills in. Each verified by the survey's rendered OUTPUT on a constructed code-repo FIXTURE with a branch-and-merge workflow — never the user's private repo, and never a prose-grep of the skill.)

**AC-1 (sketch) — the commission offer reflects iteration/steering, not gates-as-the-whole-model.** Verified by: the rendered offer on a fixture does not lead with "explicit approval gates" as the value proposition and frames the agent-iterates/you-steer loop; checked against the produced output.

**AC-2 (sketch) — work-by-area is branch-aware and does not false-signal "scaffolding > product".** Verified by: a survey run over a fixture where product code lands via merged PRs (and scaffolding is edited directly on the branch) produces a work-by-area that attributes the product work and does not conclude scaffolding dominates — observed in the rendered report.

## Notes

Companion to `zbysa` (survey-knowledge-work-archetype), related to `za` (survey-report-subagent-dispatch-fact) and the done `47rx`/`69`. Both critiques came from a real user run; the corpus specifics stay out of this shared repo. #1 is a framing/positioning improvement that may reach beyond the survey into how spacedock describes its own value; #2 is a concrete work-by-area correctness fix for branch-and-merge repos.
