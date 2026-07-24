---
title: Workflow-declared section reads assemble stage context by pointer
status: backlog
source: "Durable-decisions sprint audit of Subspace 198f762, Spacedock 6y, and held js6 stakes-read-through, 2026-07-24."
score: "0.8"
id: 0cj3qf6fefedfj7j9exq62jb
gates:
    version: 1
    current:
        gate: gate:docs-dev:0c:backlog
    records:
        - id: gate:docs-dev:0c:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:0c-backlog-1
              briefing:
                id: briefing:docs-dev:0c:backlog:attempt-1:revision-1
                digest: sha256:116794807857c45a046d56cf9217b822e6c4fbe21077edb36f0951a02995f364
                digest-domain: canonical-bytes
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:docs-dev:0c:backlog:1
                briefing: briefing:docs-dev:0c:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-07-24T03:59:33.82556Z"
                decision: approve
                reason: approve, no dispatch yet
              application:
                action: advance
                target-stage: ideation
                state: pending
                blockers: []
---

A workflow can put a cross-stage policy in its README, but `dispatch show-stage-def` returns only the selected `### stage` subsection. Let workflow metadata point at additional relevant sections and let the dispatched ensign load those exact sections through the existing structured read path.

## Problem

Subspace commit `198f762` added a top-level `## Constraint Authority`, while its stage definitions mention the policy without carrying its three authority rules. The live dispatch packet fetches only the stage subsection, so the worker cannot reliably apply the top-level policy.

Copying the policy into every stage is ad hoc and creates divergent authorities. Injecting the whole README is noisy. Hardcoding `Constraint Authority`, `Stakes`, or another heading into the launcher would create one product mechanism per policy.

The existing substrate is close: workflow stages already live in README frontmatter and appear as structured `stages` in `status --read README.md --json`; the same read helper returns a fence-safe heading index with exact offsets and section lengths. What is missing is a workflow-declared pointer from a stage (or stage defaults) to additional sections, plus assembly into the dispatch fetch instructions.

## Proposed approach

Ideation should spike and choose the smallest schema and command delta that provides this behavior:

- A workflow declares, in stage metadata or its defaults, an ordered list of README section references relevant to that stage.
- `status --read README.md --json` exposes the resolved references alongside the stage data.
- `dispatch build` places fetch instructions or equivalent pointers in the packet; it does not inline the whole README.
- The ensign bootstrap loads each declared section before stage work, using the existing fence-safe heading map/read helper rather than grep-derived ranges.
- Missing, duplicate, or ambiguous headings fail before work with an actionable diagnostic.
- Section names carry no built-in product semantics. The workflow owns what `Constraint Authority`, `Stakes`, or another section means.

The exact frontmatter spelling and whether the read helper gains a narrow section-body projection are ideation decisions, not requirements in this seed. Reuse the held `js6` spike proving byte-exact README section extraction; do not rebuild that parser or revive its hardcoded stakes ontology.

## Obligation delta

- **Bearer:** a dispatched ensign reads only sections explicitly declared for its current stage; the dispatch builder resolves and presents the pointers.
- **Burden:** bounded by the declared ordered list; no implicit whole-README read and no per-heading special case.
- **Authority:** the workflow's committed frontmatter plus captain direction from this audit.
- **No inferred obligation:** a referenced section supplies workflow context; it does not independently widen an entity's accepted ACs.

## Out of scope

- Hardcoding `Constraint Authority`, `Stakes`, development stage names, or any section's meaning.
- Injecting the entire workflow README or copying shared policies into every stage definition.
- A general artifact registry, remote URI fetcher, or dynamic include language.
- Changing gate/recorder semantics; the sibling FO follow-up owns gate-time provenance.

## Acceptance criteria

**AC-1 (VALUE) - A dispatched worker receives the workflow policy that its stage explicitly references, without a copied policy in the stage subsection.**
Verified by: a behavior fixture whose top-level policy contains a unique rule and whose stage declares the reference; `dispatch build` plus the ensign bootstrap returns the exact section, while the current stage-only baseline returns zero policy bytes.

**AC-2 - The mechanism is generic and bounded.**
Verified by: the same fixture using two arbitrary section names and non-development stage names; only the ordered, declared sections are loaded, duplicates are de-duplicated deterministically or rejected as specified, and unrelated README sections never enter worker context.

**AC-3 - Section references are visible through the existing structured read surface.**
Verified by: `status --read <fixture README> --json` exposing the resolved per-stage references together with the existing `stages` array, and changing when the frontmatter declaration changes.

**AC-4 - Broken references fail before stage work.**
Verified by: missing and ambiguous-heading fixtures producing a nonzero dispatch/bootstrap result with the workflow path, stage, and requested heading; no worker implementation action occurs.

**AC-5 - Existing workflows and packets remain unchanged when no section references are declared.**
Verified by: current dispatch golden fixtures remaining byte-identical and the existing status/read and stage-definition suites passing.

## Test plan

First reuse the held `js6` spike and the existing `internal/status/section_read` heading-map fixtures. Then add focused parser/read JSON tests and augment the existing dispatch-build golden/ensign bootstrap fixture with declared, absent, missing, duplicated, fenced-heading, and arbitrary-stage cases. Run the repository's focused tests, `go test ./...`, and `go test ./... -race`; require one live fixture-backed dispatch only if the ensign-loading behavior is not already exercised by the current live lane.
