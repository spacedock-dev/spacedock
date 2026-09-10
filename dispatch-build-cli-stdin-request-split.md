---
title: Separate dispatch invocation controls from stdin assignment data
status: ideation
source: Captain discussion of checklist transport and dispatch assembly, 2026-09-10
started:
completed:
verdict:
score: 0.5
worktree:
issue:
pr:
id: 09ptxgq8zma6qnx0h1w8a6wp
gates:
    version: 1
    records:
        - id: gate:09ptxgq8zma6qnx0h1w8a6wp:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:09ptxgq8zma6qnx0h1w8a6wp-backlog-1
              briefing:
                id: briefing:09ptxgq8zma6qnx0h1w8a6wp:backlog:attempt-1:revision-1
                digest: sha256:cb3ae42ccb7d64430aef08321f5cb896c7ad8566738a23db978598647dde5adc
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:09ptxgq8zma6qnx0h1w8a6wp:backlog:1
                briefing: briefing:09ptxgq8zma6qnx0h1w8a6wp:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-09-10T20:26:00.677523Z"
                decision: approve
                reason: 'Captain requested: assume this behavior already works in this session. now dispatch it. Authorization is to begin ideation, not implementation.'
              application:
                target-stage: ideation
                state: consumed
---

Design and implement a clear dispatch build interface with invocation controls in CLI flags and assignment data in stdin JSON. Captain requested filing only; do not dispatch or implement from this filing.

## Problem

FO dispatch currently writes checklist text to a temporary file and passes --checklist-file. The helper reads it only during assembly and embeds its contents in the generated assignment. No downstream reader needs the original checklist file.

Stdin JSON already exists, but request flags such as --stamp and --advance select flag/file mode. That prevents a straightforward combination of CLI invocation controls and structured assignment data. Avoid creating a second prompt assembler: the helper must continue producing the assignment and pointer envelope.

The original stdin interface came from the Python helper and was preserved by native commit 39454b644. Commit cfa3b671c added flag/file mode to avoid fragile shell quoting of prose. Any new default must preserve safe transport of arbitrary checklist, Markdown, scope, and feedback content. Removing a scratch file is a small simplification, not a claimed large latency improvement: writing the file and invoking build can already share one tool call.

## Proposed direction

CLI controls: --workflow-dir; --stamp or --advance; explicit --host when needed; transport switches such as --bare-mode; schema/validation controls where applicable.

Stdin assignment JSON: entity_path, stage, checklist array, optional scope_notes and feedback_context. Decide schema-version handling and placement of remaining existing request fields during ideation. Each value has one authoritative location in the new form; reject conflicting or ambiguous mixed inputs rather than silently choosing precedence.

Preserve existing stdin-only and flag/file callers. Define explicit input-mode selection that does not silently reinterpret existing invocations. Programmatically serialize arbitrary inserted prose; a quoted heredoc prevents shell expansion but does not itself escape JSON.

The helper retains sole assignment assembly. Its output remains a JSON envelope with a short file-pointer prompt; the worker reads the generated assignment. The FO does not concatenate a second full prompt.

--stamp keeps its current meaning: require status already matching the dispatched stage; stamp started/worktree fields; commit/synchronize dispatch state; create the stage worktree when required; then emit the assignment envelope. It neither advances stage nor spawns the worker. Preserve incompatibility with --advance and existing failure ordering.

## Acceptance criteria

**AC-1 — FO dispatch can supply assignment text without scratch input files.**
Verified by: a real build fixture combines invocation controls with stdin assignment JSON and emits the expected assignment/envelope without checklist, scope, or feedback input files. Requiring a checklist file again must fail the fixture.

**AC-2 — Assignment content survives transport unchanged.**
Verified by: existing byte-hazard fixtures cover multiline Markdown, quotes, backticks, dollar expressions, Unicode, and trailing newlines through the new input path. Assert generated assignment content against independently specified input bytes.

**AC-3 — Existing callers and dispatch effects remain compatible.**
Verified by: existing stdin-only, flag/file, advance, and stamp fixtures retain expected behavior. New tests reject conflicting duplicate values and preserve status-match checks, stamp/advance incompatibility, durability ordering, and no spawn on build alone.

## Test plan

Use existing internal/dispatch input-mode, JSON ergonomics, byte-hazard, advance, and stamp test owners. Add focused failing behavior tests before implementation. Update FO command text only with appropriate skill smoke and live-journey proof; do not substitute prose-presence tests. Run required Go and race checks for implementation.

## Out of scope

A new prompt assembler, embedding full assignments in spawn messages, changing checklist judgment or report semantics, worker pools, provider batching, stage-completion metadata, and changing --stamp lifecycle authority.
