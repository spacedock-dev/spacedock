---
title: State-repo verbs — spacedock state ready / sweep / commit <slug> (rebase-HALT enforced by the verb)
status: backlog
source: 0205 carve (2026-06-17, captain "stamp them") — index DoD candidate "state-verbs"; 2y MERGED unblocks it.
score: 0.5
sprint: 0205-layered-fo
group: verb-core
sprint-readiness: ready
id: rgq0m30693ke84zvb8yna078
---

`spacedock state ready` / `state sweep` / `state commit <slug>` behind the «state.ensure-ready» / «state.sweep-merged» / «state.commit» prose-functions. The commit verb is path-scoped, pushes with retry-on-reject-rebase, and REFUSES to return on a rebase conflict (exit non-zero, stderr naming the paths, repo left in rebase-abort) — so the FO cannot proceed even if it wanted to. The rebase-conflict halt is enforced by the verb, not by FO discipline.

Spike must-build: «state.commit» — the `git -C` form drifted 2/3 under Haiku in the w4 spike, and it carries the rebase-HALT trigger w4 could NOT exercise (failure-mode 1, NOT-EXERCISABLE). Depends on 2y (MERGED, the host-neutral cores).

Ideation must resolve: the three verb surfaces + their return contracts; SPIKE the rebase-HALT FIRST (a 2-writer / injected-conflict harness that actually FIRES the halt) + a Go test that «state.commit»'s `git -C` form survives; oracle-based ACs (never a prose-grep).
