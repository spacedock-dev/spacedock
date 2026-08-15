---
id: x8g3dnqndfa1m85d8ga2cgem
title: Trim version output decoration without a reader
status: ideation
source: "0.27 cut audit (2026-08-14), adversarially verified; captain directed filing"
started: 2026-08-15T02:55:38Z
completed:
verdict:
score:
worktree:
issue:
gates:
    version: 1
    records:
        - id: gate:x8g3dnqndfa1m85d8ga2cgem:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:x8g3dnqndfa1m85d8ga2cgem-backlog-1
              briefing:
                id: briefing:x8g3dnqndfa1m85d8ga2cgem:backlog:attempt-1:revision-1
                digest: sha256:29b040ac52ea6ac3b5ec088df13bd65c14cf47c40a23c4fbd9c953003d046676
                request-digest: sha256:f9aed9bec822660eeb6e0eaac5f049b2e1c6707f2fdcb404ba072ef674111704
                room-ref: ./trim-version-output-decoration/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:x8g3dnqndfa1m85d8ga2cgem:backlog:1
                briefing: briefing:x8g3dnqndfa1m85d8ga2cgem:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-15T02:53:43.914595Z"
                decision: approve
                reason: 'Captain ruling 2026-08-14 (dispatch them): approved into ideation'
              application:
                target-stage: ideation
                state: consumed
---

Remove three decorations from `spacedock --version` output. Each has no reader.

1. The `OS:` line (internal/cli/cli.go:883). The shared contract orders the FO to use `uname -s`, not this line. Doctor computes GOOS itself.
2. The session short-id segment on the Runtime line (cli.go:893). Session matching everywhere reads the env var, never the printed prefix.
3. The `pass --host` remedy suffix on the ambiguous arm (cli.go:910-916). The one ambiguity-affected command, dispatch build, prints its own complete remedy.

Keep the Runtime line itself (host and marker have a recorded live read), the ambiguous reporting arm (the nested-marker leak occurs live), and the Sandbox line (fo-install-gate corroboration).

Captain override recorded here: the OS line was a captain-ratified annotation this cycle. Its ratification record names only speculative consumers. Approval of this entity reverses that annotation.

## Problem

{Ideation fills this in.}

## Proposed approach

{Ideation fills this in.}

## Out of scope

The Runtime detection redesign. The Sandbox line. Contract prose beyond the docs examples that quote the removed lines.

## Expected surface and tolerance

Estimate net LOC change: -NN, across ~3 files. Observable semantics change: --version output shrinks by one line and two segments.

## Acceptance criteria

**AC-1 - The change removes more lines than it adds.**
Verified by: cumulative line delta against origin/main is negative.

**AC-2 - --version output carries no OS line, no session segment, and no --host suffix.**
Verified by: the CLI version test asserts the exact output shape.

**AC-3 - The suite stays green.**
Verified by: go test ./... and go test ./... -race pass.

## Test plan

Deletion, updated version-output assertions, docs example prose reconciled.
