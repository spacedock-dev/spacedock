# 0203 — FO efficiency (0.20.3)

**Sprint:** the entities matching `sprint: 0203-fo-efficiency` — list current members with `spacedock status --workflow-dir docs/dev --where sprint=0203-fo-efficiency`. Membership and per-task state are the query, never enumerated or tracked in this doc.
**Theme:** make the first officer cheap to boot and run.

## Goal (success criterion)

An FO reaches interactive readiness — greet + state summary + *able to present a gate* — in seconds at **< ~60k** context, versus today's minutes at 126k+. Proven by a live FO-boot drive that **measures** the saving (j9 AC-6), never a grep over the restructured contract.

## Why

Boot forensics (`boot-analysis.md`) measured ~160k peak context and ~13.6 min to greet — with no team created and no worker dispatched. Structural, not a bug.

## Cost levers (ranked)

| Lever | ~boot cost removed | Needs the split? |
|-------|-------------------:|------------------|
| Lazy-TeamCreate (defer the team-mode prefix re-cache) | **~89k** | no |
| Defer contract reads at greet | ~16k | yes (minimal) |
| Defer the human status-table render | ~8.7k | no |
| Defer mod-file reads | ~6.5k | no |

## Definition of Done

0.20.3 ships when, merged to `main` (the post-flip integration trunk):
- **j9** — the FO contract is split into a boot-resident core + deferred dispatch/merge references; `TeamCreate` deferred off the boot/greet path; shallow-boot-then-greet off `status --boot --json`. AC-1..AC-6 green, including the live shallow-boot scenario, the offline gate staying green post-split, the `contractlint` closure test, and the **measured-saving drive** (greet context < ~60k, no pre-greet ~89k spike).
- **#344** — the context-budget probe emits no spurious `config_drift`/`mixed_models` warnings on healthy members and reads the correct window.
- **T3** — the slimmed FO refs are audited + comm-officer-polished, behavior-preserving (live scenarios green) and measurably smaller — or a recorded roadmap decision if the split left nothing to cut.
- **sr + 87** — the stale-trunk cluster closed: no helper/mod/ref/doc resolves the integration trunk to `next` post-flip, settled as ONE trunk-config source. `sr` de-conflates `dispatch reconcile`'s git-hygiene from team-management; `87` refits the `pr-merge` mod base `next`→`main`.
- **xf** — the standing-teammate lifecycle replaced by on-demand one-shot polish; feature usage prose moves into the `comm-officer` mod, the FO contract keeps only a generic hook.
- `v0.20.3` cut after the pre-cut antipattern audit is clean.

## Out of scope

p2/vc (0.20.4 binary-simplification line); xp (cross-session FO↔Commander comms — the coordination gap this sprint surfaced); ey (proof-policy port to shipped scaffolding).

## Artifacts

- `staff-review.md` — preflight readiness gap analysis (verdict: READY)
- `dispatch-sprint-execution.md` — cold-boot Commander dispatch package
- `boot-analysis.md` — the boot forensics (evidence base)
