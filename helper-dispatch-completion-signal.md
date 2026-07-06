---
title: "Ad-hoc helper dispatches carry the ensign completion-signal contract (no silent-idle helpers)"
status: backlog
source: "FO Commander session 2026-07-06/07 (0250 drive): the k7 detached-audit helper (named background general-purpose agent, hand-rolled prompt) emitted its 9,215-char report as final assistant TEXT (transcript agent-ak7-detached-audit-e5f62bb29dfae423.jsonl event 132) and idled silently; the report arrived only after an explicit FO prod (first-ever SendMessage at event 140). Ensigns never hit this: dispatch build pins a '### Completion Signal' SendMessage(to=team-lead) block into every envelope. Ad-hoc helpers — including the detached adversarial audits the Proof policy REQUIRES per high-stakes member — get hand-rolled prompts with no completion contract. Blocks reliable autonomous/headless operation."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
id: ehvqaf7makrvc38zkb4fh0qt
---

A named background teammate's plain-text output is invisible to the team; only SendMessage delivers. `spacedock dispatch build` guarantees this for entity-stage dispatches via the emitted Completion Signal block, but ad-hoc helper dispatches (detached audits, staff reviewers, pre-cut auditors) are hand-rolled and silently idle on completion, requiring an FO prod — incompatible with autonomous mode. Proposed direction: a helper mode on the shipped builder (e.g. `dispatch build --helper --role-file <prompt>`) that emits the standard envelope — name, run_in_background, and the same pinned Completion Signal block stage dispatches get — so every background worker, stage or helper, carries an enforced delivery contract. Alternative considered and disfavored: a prose-only template in claude-fo-dispatch.md (unenforced, the exact prose-vs-code-gate gap this workflow's Proof policy warns about). Acceptance sketch: value — a named background helper spawned via the shipped mechanism delivers its final report to team-lead with zero FO prods, proven by a live helper dispatch; mechanism — the helper envelope carries the Completion Signal block, unit-tested at the envelope level.
