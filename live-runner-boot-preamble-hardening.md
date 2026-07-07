---
title: "Live-runner boot preamble hardening — a driven FO's boot fumbles can't reach the scenario under test"
status: backlog
source: "0250 Commander session 2026-07-07 (captain-requested filing). Two same-day sonnet-lane instances of boot-phase nondeterminism failing scenarios UPSTREAM of their subject matter: (1) k7's PR run — TestLiveDefaultHeadlessStopsAtGate, the driven FO's `status --discover` returned empty twice in the fixture and it correctly stopped at no-workflow (rerun green); (2) zm's PR run — TestLiveClaudeSharedScenarios/smallest-sufficient-mechanism, the driven FO prefixed its boot with `cd /home/user` and operated outside the fixture root (caught by the runner's own claude_live_runner_test.go:130 diagnostic; opus variant green on the identical suite both times). Related same-day codex instance: shallow-boot's FO tried four wrong `state sweep` flag shapes (never `state sweep --workflow-dir .`), concluded the verb was broken, and improvised a partial manual advancement that skipped the archive — a call-shape fumble, not a deliverable defect. Each red costs a full lane run + FO triage + rerun; the flake tax erodes lane trust."
started:
completed:
verdict:
score: 0.4
worktree:
issue:
id: sc592grb0w36q8ravzeya70r
---

Live scenario lanes let the driven FO's boot-phase nondeterminism (cwd wandering, discovery fumbles, call-shape guessing) fail scenarios before their subject matter runs — the assertion then reds on a preamble accident, indistinguishable in the checks summary from a real behavioral failure. Direction for ideation: harden the preamble on both sides — (a) runner/prompt side: anchor the launch context explicitly (fixture root named in the drive prompt; runner pins/validates cwd before the scenario clock starts; possibly a boot-preamble retry distinct from scenario retry, since a boot fumble is not scenario evidence); (b) binary side (rider): `state ready` / `state sweep` tolerate the obvious flag shapes an FO guesses (`--workflow-dir` in any position, or a positional dir) instead of exit-2ing a competent caller — cf. tv's default---stage precedent for meeting callers where they are. Constraint: hardening must NOT mask genuine boot-behavior regressions — scenarios whose SUBJECT is boot behavior (k7's greet drives) opt out. Acceptance sketch: value — a seeded boot-fumble drive (cd-away or discovery fumble) either self-corrects or reds with a distinct preamble-failure diagnostic, never as the scenario's own assertion (baseline: the two 2026-07-07 sonnet reds); mechanism — the runner guard + prompt anchor ship, with the flag-shape rider or its explicit rejection recorded.
