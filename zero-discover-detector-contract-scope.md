---
title: "zero-discover broad-search detector flags shapes the contract does not ban (flat ls of cwd)"
status: backlog
group: tooling
source: "PR #466 claude-live sonnet (2026-07-02, run 28576842676): TestLiveZeroDiscoverReportsAndStops failed because the FO ran `ls -la {project root}` after a zero discover — the detector (detectBroadSearchAtBoot, asserted at zero_discover_live_test.go:113) flags ANY ls over the project root, while the contract's block clause (first-officer-shared-core.md Startup step 3) enumerates only `find` / `grep -r` / `ls -R` / recursive Glob/Grep. A flat, non-recursive listing of the launch cwd is not on the banned list and does not hunt a workflow. Third over-broad live-heuristic instance of the day (after the wrong-root cd detector, #462, and the codex narration negation scope, task mq) — the class is now a recurring merge-blocker."
id: 4t8ej1rmpmk2hzzpshtrty0s
---

## Problem
The live zero-discover scenario's detector is stricter than the contract it enforces: the contract bans recursive/broad filesystem hunts after a zero `status --discover`; the detector reds on any `ls` touching the project root, including a flat `ls -la` of the FO's own cwd — plausible harmless orientation before report-and-stop. Correct FO behavior (report no workflow, stop) plus one innocuous listing = lane failure.

## Desired direction (for ideation to refine)
Reconcile detector and contract deliberately — either narrow the detector to the contract's enumerated shapes (find / grep -r / ls -R / recursive glob over the root, plus obvious equivalents like `ls **/`), with the real PR #466 stream as a must-NOT-match fixture; or, if the intent is genuinely zero post-discover filesystem ops, tighten the CONTRACT prose to say so and keep the detector — but then the contract change is the deliverable and the live lane proves it. The determination of intent (is a flat cwd listing acceptable FO behavior at zero-discover?) is the ideation's first question; align both artifacts to whichever answer, never leave them disagreeing.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- The captured PR #466 sonnet stream replays through the detector with the chosen semantics (must-NOT-match if narrowed; must-match if the contract tightens and the scenario prompt/oracle changes accordingly).
- Genuine broad-search shapes (find, grep -r, ls -R, recursive glob) still red — RED/GREEN fixture pair.
- Contract prose and detector agree verbatim on the banned set; contractlint or a fixture binds them if feasible per the read-quarantine rules.
