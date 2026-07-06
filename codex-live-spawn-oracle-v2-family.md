---
title: "codex-live rejection-flow oracle accepts the current multi_agent_v2 spawn-evidence family (lane red on unchanged main)"
status: ideation
sprint: 0250-fo-behavioral-discipline
source: "0250 Commander session 2026-07-06/07. Timeline: Jul-4 main run green; 2026-07-06 k7 branch run, branch rerun, AND a baseline rerun of the Jul-4 green main run (same code) all red at codex_live_runner_test.go:42 'no validation spawn_agent found — the FO never created a cycle-1 reviewer to reuse'. codex-cli 0.142.5 identical across green and red runs, so not a CLI version change: service-side behavior drift under --enable multi_agent_v2 (the harness opts in deliberately and asserts the flag, codex_live_runner_test.go:283,399-400). Transcript evidence (artifacts, run 28806717191 + run 28693890413 attempt 2): the FO ran a REAL two-worker flow — separate implementation worker and validation reviewer, followup_task cycle-2 reuse, populated agents_states on collab wait calls, correct end-state markers — but zero spawn_agent-shaped events; its spawn mechanism now emits evidence the oracle does not recognize. Blocks every 0250 member's merge (all-lanes-green is the DoD's proof policy)."
started: 2026-07-06T23:10:29Z
completed:
verdict:
score: 0.6
worktree:
issue:
id: 3v5cd2fcx8y5bk1dcxt7swzf
---

The codex-live shared-scenario runner proves worker separation by grepping the exec stream for spawn_agent-shaped events; under current multi_agent_v2 service behavior a genuinely-separated flow emits a different spawn/collab evidence family, so the lane is red on unchanged main and cannot go green for any diff. Direction: harden the spawn-evidence detection (internal/ensigncycle codex runner, with internal/dispatch/codex_v2_adapter.go's v2 modeling) to accept the current v2 family as reviewer-existence proof while REMAINING a real separation oracle — it must still go red when no separate reviewer existed (its whole point is catching an FO that validates inline; narration must never satisfy it). Evidence anchors for ideation: green-run artifacts (run 28693890413 attempt 1, Jul-4) vs red-run artifacts (attempt 2 + run 28806717191 attempts 1-2) — diff what spawn evidence looked like green vs red before writing the new oracle. High-stakes surface (a live lane's own tests → CI machinery): detached adversarial audit required, including a constructed no-separate-reviewer transcript the hardened oracle must reject. Acceptance sketch: value — codex-live returns green on unchanged main scenarios with the hardened oracle (baseline: red today), and 0250 members' merges unblock; mechanism — the oracle change ships with the anti-tautology negative case.
