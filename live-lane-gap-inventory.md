---
id: 1496q9vd177hxtdeafcp7hvj
title: Inventory the live-lane failures and say which are tracked, which are flakes, and what to do
status: ideation
source: "Captain CL, 2026-08-18, after a night of live-lane reds during the 0.27.0-pre8 cut: \"have another ensign inventory of recent live gap and if they are tracked as known flakes, and propose recommendations\". Ten distinct live failures were observed across four runs on three hosts; a prerelease gate was waived without being able to read them."
started: 2026-08-18T15:20:42Z
completed:
verdict:
score:
worktree:
issue:
---

The live lanes failed repeatedly across one night. Nobody can currently say which failures are known, which are tracked, which repeat, and which are environmental. Produce that inventory and a recommended disposition for each.

## Problem

During the 0.27.0-pre8 cut the live lanes failed on every host, and the FO could not distinguish a real regression from an environmental stall. The `e2e-gate` was waived on captain decision without a diagnosis.

The observed failures, with the exact `observed=` markers and run ids, are listed below. They fall into at least three shapes, and the shapes matter more than the count:

- **Timeouts with no assertion.** `made no stream progress within 1m0s (no-progress quiet budget) — a hung stage; killed the subprocess`. Nothing failed; the harness killed a quiet subprocess.
- **Agent contract violations.** `observed=[smallest-mechanism-violation]`, `[recorded-gate-lifecycle-violation]`, `[human-gate-bypassed validation-worker-lifecycle]`, `[implementation-worker-not-dispatched]`, `[filing-command-not-observed]`, `[rejection-worker-topology]`. These are real: a live agent disobeyed its own contract.
- **Infrastructure faults.** `fatal: unable to read <object>` from `git log --follow` inside a self-contained `t.TempDir` fixture, on the deterministic `offline` lane, non-reproducible locally and green on immediate re-run.

The backlog already carries entries whose titles match several of these areas. Whether those entries actually cover the observed failures, or merely sound like them, is exactly what nobody has checked.

## Observed failures to inventory

Treat this list as the starting evidence, not the boundary. Pull the real logs; do not trust this summary.

| Run | Lane | Failure |
|---|---|---|
| 32092321763 attempt 1 | claude-live | `rejection-flow`, `break-glass-shim-selected-team` — both 1m0s no-progress timeouts |
| 32092321763 attempt 2 | claude-live | `smallest-sufficient-mechanism` = `smallest-mechanism-violation`; `auto-continue-after-implementation/split-root` = `human-gate-bypassed validation-worker-lifecycle`; diagnostic: `FO broad-searched the filesystem at boot` |
| 32105482382 | claude-live | `smallest-sufficient-mechanism` (repeat); `recorded-gate-lifecycle` |
| 32105482382 | codex-live | `filing` = `filing-command-not-observed`; `rejection-flow` = `rejection-worker-topology` |
| 32047943955 | pi-live | `default-headless-gate-stop` and `auto-continue-after-implementation` = `implementation-worker-not-dispatched`; `recorded-gate-lifecycle` = `recorded-gate-lifecycle-violation`; `rejection-flow` |
| 32105482382 | offline | `TestDurableKeepMovingRequiresOverlappingJourneys`: `git log --follow ... questioned.md: exit status 128`, stderr `fatal: unable to read 8354e03b...`; passed on immediate re-run |

One failure repeated: `smallest-sufficient-mechanism`, in two consecutive claude-live runs. Every other failure appeared once.

## What to produce

For each observed failure: the marker, the host, the run ids, whether it repeats, its shape, whether an existing backlog entity genuinely covers it, and a recommended disposition.

Dispositions should be concrete: an existing entity already covers this and needs no action; an existing entity is close but its scope must change; this needs a new entity; this is environmental and needs no entity; this cannot be classified until `r5` lands.

## Out of scope

Fixing any of them. Filing the follow-up entities — recommend, and let the captain decide what gets filed. Changing the live lanes, the harness, or the quiet budget.

## Expected surface and tolerance

Estimate net LOC change: 0 code. The deliverable is the inventory and its recommendations in this entity body. Declare the body's size instead, and do not add code, tests, or CI.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - Every observed failure has a disposition backed by evidence, and none is called a flake without proof.**
This is the measuring AC: the count of failures classified as "known flake" or "environmental" WITHOUT a cited reproduction, re-run result, or root cause must be ZERO. Verified by reading each such classification for its evidence citation. Fails the moment a failure is dismissed by assertion — the exact move the workflow's own proof policy forbids and the FO made twice during this incident.

**AC-2 - Each claimed coverage is checked against the entity's real content, not its title.**
Verified by: for every observed failure mapped to an existing backlog entity, the mapping cites the specific text in that entity that covers this failure mode. Fails if an entity is credited because its title matches the area — for example, crediting a rejection-flow entity for a rejection-flow timeout when it addresses agent behavior and not the quiet budget.
