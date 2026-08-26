---
id: 9n02rsw1s4tztqzgmwb07n1k
title: gate prepare resolves operator-supplied artifact paths without doubling
status: ideation
source: "email-triage codex audit 2026-08-26: three of six gate-prepare attempts across two days failed with a doubled state path (.../.spacedock-state/docs/triage/.spacedock-state/...); under the no-retry rule the third failure left a batch's gate unprepared for the rest of the window"
started: 2026-08-26T21:19:07Z
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
        - id: gate:9n02rsw1s4tztqzgmwb07n1k:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9n02rsw1s4tztqzgmwb07n1k-backlog-1
              briefing:
                id: briefing:9n02rsw1s4tztqzgmwb07n1k:backlog:attempt-1:revision-1
                digest: sha256:a1db71cb8071b29305e8b027b71b0c14747ac5c4c62eaee6631f73b1de70ede5
                room-ref: ./gate-prepare-accepts-operator-paths/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:9n02rsw1s4tztqzgmwb07n1k:backlog:1
                briefing: briefing:9n02rsw1s4tztqzgmwb07n1k:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-26T21:18:44.298895Z"
                decision: approve
              application:
                target-stage: ideation
                state: consumed
---

`gate prepare --artifact` and `--reference` resolve supplied paths relative to the state entity directory. An operator who supplies a project-relative path (`docs/triage/.spacedock-state/x.md`) gets it re-prefixed under the state root and the command fails: "selected source must be a readable non-symlink regular file" with the doubled path. Three of six live preparations hit this, and the gate skill's cwd-path wording is ambiguous against the binary's resolution rule.

## Problem

{Ideation fills this in. Seeded verbatim failures from the audit: --artifact docs/triage/.spacedock-state/2026-08-17-personal-2.md -> "Error: --artifact .../.spacedock-state/docs/triage/.spacedock-state/2026-08-17-personal-2.md: selected source must be a readable non-symlink regular file"; same shape on --reference docs/triage/README.md; recurrence on a third prepare 2026-08-26T17:18Z. The no-retry rule converts each failure into a withheld gate, so the cost is a stalled batch, not one turn. The gate-lifecycle skill's Prepare wording ("resolve intended root-relative P; state R uses the engaged task directory, never `.`") and the binary's actual resolution disagree in the operator's hands.}

## Proposed approach

{Ideation fills this in. Candidate directions, smallest first: accept absolute and cwd-relative paths by trying the supplied path as-is before prefixing; or refuse with an error that names the exact path form the command accepts and the path it resolved; and align the skill's Prepare wording with whichever the binary does. Ideation picks the smallest sufficient set.}

## Risk evidence

{Backlog: three verbatim transcript failures across two days on the same misresolution decide design should start.}

## Out of scope

Room layout, briefing schema, the no-retry rule itself.

## Expected surface and tolerance

Estimate: production +20 across 2 files; proof +30 across 1 file. {Backlog seed; ideation refines.}

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 (VALUE) - A gate prepare given an absolute path, a cwd-relative path, or a state-relative path to the same readable artifact succeeds identically — or, if a form stays unsupported, the refusal names the resolved path and the accepted form.**
Verified by: {ideation refines — seed: a table test over the three path forms against a fixture room; failing-today baseline: the project-relative form errors with the doubled path.}

## Test plan

{Ideation fills this in. internal/gates is the status-guard high-stakes surface; the detached audit applies.}

### Feedback Cycles

{First officer appends one `- Cycle {N}: ...` line per correction round; the validation gate reads reviewer findings from here.}
