---
id: s4ykctf21g60dvfgdd6cy9ny
title: Prepare provider-neutral gate rooms and align canonical Briefing recording
status: ideation
source: "Durable-decisions cross-repo dogfood ruling after xb and Subspace em review, 2026-07-24"
started: 2026-07-24T14:54:10Z
completed:
verdict:
score: "1.0"
worktree:
issue:
sprint: durable-decisions
gates:
    version: 1
    current:
        gate: gate:docs-dev:s4:backlog
    records:
        - id: gate:docs-dev:s4:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:s4-backlog-1
              briefing:
                id: briefing:docs-dev:s4:backlog:attempt-1:revision-1
                digest: sha256:8d6888f2f9d067835f24c8845d703547638ff919f71f709c681e856551cfb80f
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:s4:backlog:1
                briefing: briefing:docs-dev:s4:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-24T14:53:41.564011Z"
                decision: approve
                reason: Captain approved filing the narrow post-em alignment task and dispatching it through the sprint.
              application:
                action: advance
                target-stage: ideation
                state: consumed
                blockers: []
---

Make a prepared gate room mechanical and self-consistent so a First Officer supplies judgment and artifact choices, while Spacedock derives recorder metadata and any presentation provider consumes one frozen package.

## Problem

The xb recorder correctly verifies a prepared room, but no command prepares that room. A First Officer still handcrafts `request.json`, a canonical Briefing at the hardcoded `briefing.json` basename, artifact digests, and output layout. This reproduces the exact dogfood friction the room contract was meant to remove and leaves Spacedock one revision behind Subspace's active `em/Remove invented Briefing review preconditions` contract.

The unreleased-v1 ruling is direct alignment with no compatibility layer. A prepared request binds a readable canonical Briefing by locator, id, and digest. The durable association remains derived from the frozen request, Briefing, exact Result, and presented inventory; no `association.json` is written or required. Spacedock remains provider-neutral and independently rejects duplicate JSON members at every authority-bearing room boundary. Chat is the default presentation channel; an override is usable only when its literal provider-package capability probe succeeds.

## Scope boundaries

- Spacedock owns mechanical room preparation, request/Briefing validation, recorder alignment, provider-neutral capability wording, and duplicate-member rejection.
- Subspace owns its active `em` primitive and the later q0 `/subspace:r gate <gate-room>` convenience mapping. Do not implement Subspace transport in this repository.
- Target the landed `em` contract before implementation completion. Ideation may inspect its active branch, but it must record any delta if `em` changes before landing.
- Do not emit `association.json`, add compatibility wrappers, require an exact provider version, or revive caller-supplied provider output paths.
- Keep the broader lifecycle-next-action, advisory-round scaffolding, and readiness projection in `gate-agent-ergonomics`; they are not prerequisites for this release slice.

## Acceptance criteria

**AC-1 (VALUE) - A First Officer can prepare a complete recorder-ready room without handcrafting recorder metadata.**
Verified by: a CLI fixture starts from a concise gate review plus selected existing artifacts, runs the preparation command, and immediately passes the canonical room/Briefing validation. The assertion fails if the test must prewrite `request.json`, choose ids/digests, or copy a reproducible referenced artifact into the room.

**AC-2 - The frozen request binds an arbitrary readable canonical Briefing filename by locator, id, and digest.**
Verified by: command-level tests record successfully when the canonical Briefing is not named `briefing.json`, then fail byte-cleanly on locator traversal, substitution, or digest mismatch. Reintroducing a basename assumption makes the positive case fail.

**AC-3 - Authority-bearing room JSON rejects duplicate members before mutation.**
Verified by: adversarial fixtures cover the request, located Briefing, Result, and presented inventory with conflicting duplicate members and assert a nonzero exit plus byte-identical entity state. Removing duplicate-member detection makes at least one fixture close or mutate the gate.

**AC-4 - Presentation selection is truthful and capability-based.**
Verified by: skill/integration behavior exercises chat as the default, a declared override whose literal package capability succeeds, and an unavailable override that halts before presentation with a corrective diagnostic. Exact provider version matching or an unconditional Subspace invocation makes the exercise fail.

**AC-5 - The retained evidence has one derived association and no parallel durable artifact.**
Verified by: an end-to-end room fixture prepares, presents, retains Result and inventory, and records the decision while the recorder recomputes their complete association; the room contains no required `association.json`. Deleting or changing any associated frozen input makes recording fail byte-cleanly.

## Test expectations

Ideation must spike the riskiest seam first: construct one room with an arbitrary canonical Briefing filename and prove the current xb recorder fails only at the basename assumption, then state the smallest change that makes the same room recordable. Reuse existing gate/room and skill integration fixtures where possible. Implementation adds focused Go command tests, adversarial duplicate-member fixtures, and the smallest existing live/provider-neutral journey that can be augmented. Validation runs `go test ./...`, `go test ./... -race`, gofmt cleanliness, the high-stakes detached adversarial audit, and every runtime live lane required by the actual diff.
