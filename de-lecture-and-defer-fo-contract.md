---
id: wn3dg7txnrte0jrcxf56b859
title: De-lecture the FO contract and defer non-boot shared-core sections (folds s6q)
status: ideation
source: "Captain CL, 2026-08-24: contract audit ruling after the #757 fo-install overbuilt review - 'file and fold s6q. dispatch. keep it light in ideation as this is clear about what to cut'; folds defer-shared-core-non-boot-sections (s6qamkh7efky9zh5jh6ba6xq, superseded into this task)"
started:
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:wn3dg7txnrte0jrcxf56b859:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:wn3dg7txnrte0jrcxf56b859-backlog-1
              briefing:
                id: briefing:wn3dg7txnrte0jrcxf56b859:backlog:attempt-1:revision-1
                digest: sha256:b879d0be71511f84291a906609ac6f12c9d70175e7665971caf86203204700c5
                request-digest: sha256:e0160113b31146057ccfa33cf4ef0a3ece370f9b5cc84ad1c782dd9516f8467a
                room-ref: ./de-lecture-and-defer-fo-contract/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:wn3dg7txnrte0jrcxf56b859:backlog:1
                briefing: briefing:wn3dg7txnrte0jrcxf56b859:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-24T23:37:42.284397Z"
                decision: approve
                reason: 'Captain CL in chat 2026-08-24: ''file and fold s6q. dispatch. keep it light in ideation as this is clear about what to cut and we should have stacked PR of the two to run live CI'' - accepts the seed with light-ideation and stacked-delivery directives'
              application:
                target-stage: ideation
                state: consumed
---

Two movements over the FO contract with one lens: the contract carries rules and commands; rationale, narration, design history, and repeated justification live in entity bodies and tests. Movement 1 (delete): the 2026-08-24 FO audit enumerated the lecture, file by file. Movement 2 (defer, folded from s6q): move boot-unreachable shared-core sections behind their existing triggers, measured by boot occupancy in tokens (s6q's archived body carries the full measured profile: shared core ~10.1k tok of a 40.4k greet boot; candidate span 9,105 of 26,092 bytes = 35%).

The audit cut list (chat, 2026-08-24, all files read in full):
- first-officer-shared-core.md: the "future multi-workflow form" design-history sentence in engage's scope; the stopping-point paragraph's re-argument (~40% compressible, exception list is the content); "Dense durable boundaries keep unhinted auto-compaction survivable"; Working Principles justification tails ("...is the asymmetry that lets a means-accurate, end-missed stage pass", "a scheme minted in a dispatch prompt becomes every downstream artifact's vocabulary" second clause, third restatements in the smallest-mechanism bullet). KEEP the banned-excuse enumerations - prohibitions are behavior-shaping, their justifications are the lecture.
- fo-dispatch-core.md: dispatch.checklist's "This is not a work breakdown..." justification; consent-stop tail ("The license hangs off... bites hardest in non-dev workflows"); second-verifier's third restatement ("N agents... agreement raises cost, not confidence"); fan-out checkpoint why-clauses (~25%).
- claude-fo-dispatch.md: Awaiting Completion states one rule four ways with narrated failure psychology - compress the narration, KEEP the anti-pattern list (bought with a real incident; captain-reviewed compression, not deletion); reconcile binding's interleaved rationale ("Cost of a miss..." etc.) - spec stands alone.
- fo-write-core.md: Workflow Fit Gate's two arguing paragraphs (rule = check fit, name the output's existing home, ask when ambiguous); the --next-id-preview/new-closes-window lecture stated three times - once plus pointers.
- fo-merge-core.md: "the toil merge.guard eliminates" tail.

The deferral groups and open Working Principles residency question are s6q's, verbatim in its archived body - ideation rules on Working Principles explicitly rather than leaving it unexamined.

## Problem

{Ideation fills this in - lightly, per captain order: the cut list above and s6q's measured profile are the spec.}

## Proposed approach

{Ideation fills this in - lightly. Sequencing question worth one paragraph: deletes first then moves, or one pass per file. s6q's starting hypothesis for the moves: first-dispatch group rides the existing dispatch-core read; engage group gets one trigger; every move leaves a one-line pointer, and pointer overhead counts against the measured delta.}

## Out of scope

fo-install.md (#757 already rebuilt it under this lens). Any change to what the FO does - load timing and prose weight only; every rule keeps its semantics and trigger. The cold baseline (system prompt, tool schemas). Orphan cleanup (tracked separately).

## Expected surface and tolerance

Estimate net LOC change: -90 net across ~7 files (deletions ~-100 to -140; deferral moves ~0 net plus pointer lines; contractlint test updates). Tolerance +/-50 net, +/-2 files. Ideation refines with the per-file table.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified. Seeded; ideation refines lightly.

**AC-1 (value) - A greet-and-stop FO boot occupies measurably fewer context tokens than the recorded 40,450-tok baseline.**
Verified by: boot-forensics per-turn occupancy on a fresh greet-and-stop session after the change, in tokens against the baseline (s6q AC-1 verbatim - the number lives outside every edited file and moves the wrong way on pointer bloat or an added boot-time read).

**AC-2 (value) - The contract's total prose weight goes DOWN, not sideways: net negative line delta across skills/first-officer/references/** with behavior parity.**
Verified by: git diff --numstat against the merge base (negative net), plus the parity evidence: contractlint suites green, and the touched hosts' live lanes green at the stack tip. Fails if the sweep relocates lecture instead of deleting it, or if any lane shows a behavior regression.

**AC-3 - No section is moved behind a trigger it can be reached before.**
Verified by: {ideation designs, inheriting s6q AC-2's constraint: not a prose-grep; contractlint structural absence plus reference closure is the static form, a live-lane boot-through-gate run is the behavioral form.}

## Test plan

{Ideation fills this in - lightly. Known required: the diff is the shipped-contract high-stakes surface, so the detached adversarial audit applies, and ALL host live lanes are required green (host-neutral core). Delivery below puts those lanes at the stack tip.}

## Delivery

Stacked layer on the existing stack, per the captain's order: branch off spacedock-ensign/install-gate-channel-aware-hint (#757, which already carries the shared-core edits this task touches next). Stack becomes #756 -> #757 -> this; the live lanes run once at this tip and prove all three layers. gh pr create + gh stack link per the pr-merge mod's Stacked mode.
