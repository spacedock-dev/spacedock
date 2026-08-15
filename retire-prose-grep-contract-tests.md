---
id: v7a6xqh2rm3asjvj8qz1y4p0
title: Retire banned prose-grep contract tests and dedupe surviving pins
status: backlog
source: "Captain review of the 0.27 stack + audit-r2 (2026-08-15); captain directive: file, dispatch off stack tip, PR as stack layer"
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:v7a6xqh2rm3asjvj8qz1y4p0:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:v7a6xqh2rm3asjvj8qz1y4p0-backlog-1
              briefing:
                id: briefing:v7a6xqh2rm3asjvj8qz1y4p0:backlog:attempt-1:revision-1
                digest: sha256:db3d3e2c2eabd9426e3eeec0af2179e2124792c069cecfb280ae50c060165f3e
                request-digest: sha256:799382d8dca6e2b600e6141353217b02e347f0288e2749291aedb9d75fb1c687
                room-ref: ./retire-prose-grep-contract-tests/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:v7a6xqh2rm3asjvj8qz1y4p0:backlog:1
                briefing: briefing:v7a6xqh2rm3asjvj8qz1y4p0:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T18:15:14.081321Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: file, dispatch based off stack tip, PR on top of the stack'
              application:
                target-stage: ideation
                state: pending
---

Delete the committed prose-grep tests the Proof policy bans (paraphrase reds them, inversion passes them), dedupe the two double pins, and resolve the four gray cases. Base all work on the stack tip (branch stack27/09-trim-version-output); the deliverable becomes stack layer 10. This MUST land before make-shipped-contracts-self-contained (layer 11), whose prose rewrites would red several of these pins.

Audit inventory (verified at the ship tree):

Banned, delete (9 functions, 4 files):
- internal/contractlint/feedback_rejection_publication_smoke_test.go:9 and :40 (both functions)
- internal/contractlint/initial_worker_spawn_guard_test.go:30, :58, :67 (nine frozen sentences plus a token count)
- internal/contractlint/pi_spawn_binding_test.go:21 only (the :35 structural-absence sibling is ALLOWED - keep)
- internal/contractlint/version_gate_smoke_test.go:35, :77, :158 (keep the :53 gate --help structural-absence guard and TestVersionGateSandboxRegistry / TestVersionGateDeferredTrigger / TestInstallHintNoDrift)

Duplicated pins, dedupe (2):
- startup_collapse_test.go:20 duplicates fo_function_reference_invariant_test.go:26 byte cap (they disagree at exactly 26900) - delete startup_collapse_test.go
- fo_function_reference_invariant_test.go:37 topology vs first_officer_eager_references_test.go:12 - keep one, delete the other

Gray, resolve with rationale (4):
- fo_write_core_mutation_gate_test.go:17,48 - leaning KEEP (table-as-data, falsifying change exists); record the rationale
- skills/integration/survey_probe_test.go:79,173 - right in spirit (executes what it extracts), violates the read-quarantine letter via documented-undetected indirection; fix by documented exemption or moving extraction behind the boundary
- live_registry_reconciliation_test.go:290,302 - code-shape greps over Go test source; the falsifying tell applies (inverting liveRuntimeRunsParallel keeps them green); cut or rebind to behavior
- live_registry_reconciliation_test.go:270 - triple-copy command literal; reduce to a two-way docs-to-workflow binding, drop the in-test copy

GATE QUESTION this entity must carry to its ideation gate: after deleting version_gate_smoke_test.go:35/:77/:158, the install-gate boot machinery has NO committed check (the stack already deleted the shell-mirror harness). The captain decides at the gate: accept the checkless-but-honest state per the 2026-07-20 ruling, or pair the deletion with a live-journey assertion.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

The prose files themselves (layer 11 owns rewrites). The allowed tests named above.

## Expected surface and tolerance

Estimate net LOC change: strongly negative, ~5-6 test files.

## Acceptance criteria

**AC-1 - No committed test asserts self-sourced prose wording; a paraphrase of any contract sentence leaves the suite green.**
Verified by: one probe paraphrase edit per formerly-pinned file, suite green, probes reverted.

**AC-2 - Each surviving pin names its independent second source; the byte cap and topology are asserted exactly once.**
Verified by: grep inventory in the report plus the deduped files deleted.

**AC-3 - The change removes more lines than it adds and the suite stays green.**
Verified by: negative delta; go test ./internal/contractlint/ ./skills/integration/ plain and -race.

## Test plan

Deletion plus probe-paraphrase evidence both sides; existing allowed lints as the floor.
