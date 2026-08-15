---
id: 9x6xw292fsz1b4648x9hn40y
title: Make shipped contract content self-contained
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
        - id: gate:9x6xw292fsz1b4648x9hn40y:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9x6xw292fsz1b4648x9hn40y-backlog-1
              briefing:
                id: briefing:9x6xw292fsz1b4648x9hn40y:backlog:attempt-1:revision-1
                digest: sha256:de4553c953de929a52bfb362a7d19a5ac077470125db6f1fad8ab8263e978581
                request-digest: sha256:3d3d92fe0a1b05ad75500987f12bc7e7d65f3b1e395d76480f0691267f354cc0
                room-ref: ./make-shipped-contracts-self-contained/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:9x6xw292fsz1b4648x9hn40y:backlog:1
                briefing: briefing:9x6xw292fsz1b4648x9hn40y:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T18:15:19.816976Z"
                decision: approve
                reason: 'Captain directive 2026-08-15: file, dispatch based off stack tip, PR on top of the stack'
              application:
                target-stage: ideation
                state: pending
---

Shipped skills reference artifacts a user's machine does not have. Rewrite the seven audited instances so every shipped sentence resolves within the progressively-disclosed contract set. Base on stack layer 10 (retire-prose-grep-contract-tests); the deliverable becomes stack layer 11 - layer 10 removes the pins these rewrites would red.

The seven, with audited rewrites:
1. skills/first-officer/references/first-officer-shared-core.md:35 - "(driver binary descoped to roadmap 0222)" -> "no driver binary backs the drive yet".
2. skills/first-officer/references/fo-dispatch-core.md:213 - roadmap 0222 + runtime-support.md references -> "no driver binary backs it yet; hand-follow the deterministic skeleton above and do not probe for the unshipped command."
3. skills/first-officer/references/fo-install-gate.md:22 - internal/safehouse/state.go and insideRegistry -> "every sandbox the binary itself recognizes - the --version Sandbox: line names the same registry". CAUTION: TestVersionGateSandboxRegistry requires each env var name AND value in prose; keep the APP_SANDBOX_CONTAINER_ID = agent-safehouse row.
4. skills/first-officer/references/fo-write-core.md:11-12 - drop the docs/dev/README.md literal (covered by {workflow_dir}/README.md) and generalize docs/dev/_mods/** to {workflow_dir}/_mods/**; update internal/contractlint/fo_write_core_mutation_gate_test.go fixtures in step.
5. skills/commission/SKILL.md:678 - "the failure mode #201 addresses" -> "the failure mode this step exists to prevent".
6. skills/survey/SKILL.md:99-108 - delete tracker numbers (#318, #69, #321-#324, #317.2, 9h, za); keep the descriptive halves.
7. skills/commission/references/templates/development.md:129 - state the rule without the dated captain-ruling attribution (line 113's captain-ruling[YYYY-MM-DD] format spec stays).

Keep-verdicts (do not touch): illustrative example ids in commission SKILL.md:385-418 and development.md:64-65; references/... and spacedock:* loader-resolved cross-references; split-root .spacedock-state mentions.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

New content; anything beyond the seven plus their in-step test fixtures.

## Expected surface and tolerance

Estimate net LOC change: near zero, ~7 shipped files + 1 test file.

## Acceptance criteria

**AC-1 - No shipped skill references a roadmap number, tracker number, repo source path, dev-workflow path, or dated ruling a shipped reader cannot resolve.**
Verified by: the audit's pattern grep over skills/ and agents/ returns only the recorded keep-verdicts.

**AC-2 - Contract meaning is preserved: reference closure and the surviving structural lints stay green, and the rewritten sentences carry the same instruction.**
Verified by: go test ./internal/contractlint/ green; before/after table in the report.

**AC-3 - The suite stays green including the updated write-core fixtures.**
Verified by: go test ./internal/contractlint/ ./skills/integration/ plain and -race.

## Test plan

Seven surgical rewrites, fixture update in step, audit-pattern grep as one-off evidence.
