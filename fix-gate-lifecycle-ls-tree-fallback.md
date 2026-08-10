---
id: qbppyssy13pyh1gtkh2n8vp5
title: Fix the gate-lifecycle ls-tree fallback command
status: ideation
source: Captain intake; recovered from deleted public issue spacedock-dev/spacedock#669
started: 2026-08-10T22:53:32Z
completed:
verdict:
score: 0.9
worktree:
issue:
pr:
mod-block:
gates:
    version: 1
    records:
        - id: gate:qbppyssy13pyh1gtkh2n8vp5:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:qbppyssy13pyh1gtkh2n8vp5-backlog-1
              briefing:
                id: briefing:qbppyssy13pyh1gtkh2n8vp5:backlog:attempt-1:revision-1
                digest: sha256:7ca8a11685cfba977d897cea574e18061a7fe81a1cb8a565f4c272f5a09f7593
                request-digest: sha256:90c2970ab0e6992546e9f0ae718f393fc3bc37336d7a3a48669896d5b76cfb08
                room-ref: ./fix-gate-lifecycle-ls-tree-fallback/review/backlog/briefing-1
              withdrawal:
                by: agent:first-officer
                at: "2026-08-10T22:37:04.823517Z"
                reason: Backlog record has no Stage Report, so checklist and acceptance evidence cannot be assembled.
            - id: gate-attempt:qbppyssy13pyh1gtkh2n8vp5-backlog-2
              briefing:
                id: briefing:qbppyssy13pyh1gtkh2n8vp5:backlog:attempt-2:revision-1
                digest: sha256:819dd15e568e3a55efc09967edd4acfaf7485728be52ca08b89b892103d09a25
                request-digest: sha256:fc1a79f889bcfcc19d8666e6dae5f108ce568c2e56121a2999103e1f436e45e5
                room-ref: ./fix-gate-lifecycle-ls-tree-fallback/review/backlog/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:qbppyssy13pyh1gtkh2n8vp5:backlog:2
                briefing: briefing:qbppyssy13pyh1gtkh2n8vp5:backlog:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-10T22:53:09.147634Z"
                decision: approve
                reason: 'Captain directed dispatch of QBP; the recovered seed is bounded to the PR #659 fallback regression and defines falsifiable fixture and exact-Codex proof.'
              application:
                target-stage: ideation
                state: consumed
---

When gate Artifact or Reference paths are absent, the First Officer must discover committed Markdown with a complete, path-scoped Git command. The current instruction abbreviates this as `git -C ... ls-tree`, which deterministically exits 129 because it omits the required tree-ish.

## Problem

Codex can follow the fallback instruction literally during gate preparation. The command fails before artifact selection, even though offline checks and all supplied-path routes can pass. This makes a deterministic product instruction defect look intermittent.

The defect entered through PR #659, commit `bbfad5b4c7886dbdee797e66e34e67a348d05cfd`. Local reproduction:

```console
$ git -C . ls-tree
usage: git ls-tree [<options>] <tree-ish> [<path>...]
$ echo $?
129
```

## Proposed approach

Specify one complete read-only command shape with an explicit tree-ish and intended path. Keep discovery limited to committed Markdown. Prove the documented command by executing it in a fixture and through the exact Codex missing-path fallback.

## Out of scope

- Changing gate authority, preparation, digest, or consume semantics.
- Changing the supplied Artifact/Reference path.
- Broad repository discovery or uncommitted-file discovery.

## Acceptance criteria

**AC-1 (VALUE) - A Codex First Officer with missing Artifact and Reference paths discovers the intended committed Markdown and prepares the gate without a Git usage failure.**
Verified by: an exact local Codex fallback journey that passes with paths omitted; omitting the tree-ish or broadening discovery makes the journey fail.

**AC-2 - The fallback command is complete, read-only, path-scoped, and excludes uncommitted Markdown.**
Verified by: a fixture-backed command test with committed intended Markdown, committed unrelated Markdown, and uncommitted Markdown; changing the tree-ish, path, or committed-only rule makes the test fail.

## Test plan

First run the exact documented command against a small Git fixture. Then run the targeted local Codex missing-path gate journey. Run the applicable focused Go tests, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` before completion.

## Stage Report: backlog

Seed outcome: When gate paths are absent, Codex discovers the intended committed Markdown and prepares the gate without a Git usage failure.

Included scope: Complete the read-only, path-scoped `git ls-tree` command with an explicit tree-ish. Keep discovery limited to committed Markdown. Prove the command in a Git fixture and in the Codex missing-path journey.

Excluded scope: Do not change gate authority, gate record semantics, supplied-path behavior, or discovery of uncommitted files.

Proof needed for ideation: Reproduce the PR #659 regression, run the proposed command against a controlled Git fixture, and confirm that the Codex journey reaches gate preparation. The evidence must show runtime behavior, not compare instruction text.
