---
id: 3kn22a1d5v5x7rdnrcva0h5a
title: Commissioned seed tasks carry an Acceptance criteria section
status: ideation
source: launch probe 2026-11-17, release defect 4 (CL)
started: 2026-11-18T01:26:38Z
completed:
verdict:
score:
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:3kn22a1d5v5x7rdnrcva0h5a:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3kn22a1d5v5x7rdnrcva0h5a-backlog-1
              briefing:
                id: briefing:3kn22a1d5v5x7rdnrcva0h5a:backlog:attempt-1:revision-1
                digest: sha256:82123b6b0dd21245098c471638edfb4809b48d3821078c332e427a040137d1a2
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:3kn22a1d5v5x7rdnrcva0h5a:backlog:1
                briefing: briefing:3kn22a1d5v5x7rdnrcva0h5a:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-11-18T01:26:07.250128Z"
                decision: approve
                reason: 'Captain in chat 2026-11-17: ''file defect 4. and check if any existing live tests cover the commission skill first. then dispatch.'' then ''override. you can push and continue dispatching.'' Accepts the seed and directs ideation.'
              application:
                target-stage: ideation
                state: consumed
---

A workflow made with `spacedock:commission` must give its seed tasks, and every task its templates produce, a `## Acceptance criteria` section, so the first officer's `status --read --ac-scan` works at the first gate.

## Problem

Commission writes the 2–3 seed tasks as frontmatter and a prose body only (`skills/commission/SKILL.md` ~505-521 on origin/main). The README task template it writes does carry `## Acceptance criteria` (~479), and so does the development template. The experiment and refinement templates (`skills/commission/references/templates/{experiment,refinement}.md`) carry no AC heading in their task snippet, so every later task of those workflows also lacks it.

The first officer runs `status --read <ref> --ac-scan` at every gate (`skills/fo-gate-lifecycle/SKILL.md`). It exits non-zero by design when the section is absent (`internal/status/gate_extract.go:366`, "no ## Acceptance criteria section in this file"). So the first gate on a seed task of any freshly commissioned workflow fails, and on experiment and refinement workflows every gate fails. A new user meets this on their first workflow.

## Proposed approach

{Ideation fills this in.}

Candidate already prepared, not yet reviewed through this workflow: local branch `fix/commission-seed-ac` at `c3895378d` (from origin/main `647d55163`). It adds the AC stub to the seed body and to the experiment and refinement task snippets, plus `internal/contractlint/commission_entity_ac_scan_test.go`, which feeds each generated body into the real `--ac-scan` (fails before on seed, experiment, refinement; passes after). Ideation decides whether to keep it.

## Risk evidence

Backlog check, 2026-11-17: no live test drives the commission skill. The live fixtures in `internal/ensigncycle` hand-write a workflow README with `commissioned-by: spacedock@1`; the "commissioned" live scenario (`claude_live_runner_test.go:568`) uses such a prepared fixture. The only test that reads the commission skill is the static `internal/contractlint/boundary_guard_control_test.go`. The layer-assembly fallback path (no template) is not yet checked.

## Out of scope

- Making `--ac-scan` tolerate a missing section; its loud failure is the intended contract.
- Repairing workflows that users already commissioned.
- A new standing live lane for commission, unless ideation shows no cheaper check can fail.

## Expected surface and tolerance

Estimate net LOC change: {ideation}, across {M} files.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - The first gate on a seed task of a freshly commissioned workflow passes `--ac-scan`, for each template and for the layer-assembly fallback.**
Verified by: a test that builds each entity body commission instructs (seed, README task template, each template's task snippet, fallback) and runs the real `status --read --ac-scan` on it with a stage report citing AC-1; it exits 0. Removing the AC block from the seed body in `skills/commission/SKILL.md` makes it fail.

## Test plan

{Ideation fills this in: name the proof owner, and decide whether a deterministic test is enough or a one-off live commission run is needed, given no live test covers commission today.}

### Feedback Cycles
