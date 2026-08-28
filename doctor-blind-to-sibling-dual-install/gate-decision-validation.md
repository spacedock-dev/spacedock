# APPROVE

Approve candidate `20608bba42fd9348e968d8f0837cbb52f33a8277` for terminal delivery.

## Decision inputs

- Main Artifact: corrected candidate `20608bba42fd9348e968d8f0837cbb52f33a8277`.
- Reference: `doctor-blind-to-sibling-dual-install.md`, entity `x0petxt7xvr459b6zh4vf4wj`.

## What changed

The launcher and doctor now scan every matching sibling record. An installed and enabled record triggers conflict handling, regardless of inventory order.

Order-sensitive launcher and doctor tests cover repeated sibling IDs. The tests include disabled-first and enabled-first inventory order.

## Why V-1 is closed

V-1 found that the first matching sibling record could hide a later enabled record. The corrected scan no longer stops at a disabled record.

The detached adversarial matrix passed five inventory arrangements. Focused launcher and doctor tests also passed against the corrected candidate.

No validation finding remains.

## Acceptance evidence

| Criterion | Evidence |
| --- | --- |
| AC-1 | Order-sensitive tests prove that an enabled sibling triggers repair before launch. Detached inventory permutations and installed-package journeys also passed. |
| AC-2 | `--no-install` exits 1, performs no install, performs no launch, and prints the host-correct repair command. The real Claude host reproduced this result. |
| AC-3 | Doctor reports both plugin identities, exits 0, and performs no install. Claude, Codex, repeated-scope, and real Claude inventory cases passed. |
| AC-4 | Healthy single-channel behavior and command references remain correct. Full, race, documentation, disabled-sibling, and ordinary front-door evidence passed. |

## Required checks

- `go test ./...` passed.
- `go test ./... -race` passed with no race report.
- Changed-file `gofmt -d` produced no output.
- `git diff --check` passed.
- `mkdocs build --strict` passed.

The installed stable-stamped Codex shallow boot passed through the ordinary front door in 36.91 seconds.

The bounded Claude journey loaded `spacedock@spacedock` version `0.28.0-pre0` through the ordinary front door. External Claude OAuth then returned HTTP 401 because its access token had expired.

This external authentication limitation occurred after bootstrap and provider selection. It does not indicate a candidate defect.

## Surface

The exact surface is +410/-83 = +327 net across 9 files. This is +72 LOC and +2 files versus the approved +255/7 estimate.

The variance is within the approved ±80 LOC and ±2 file limits.

## Decision effect

Approval authorizes the terminal delivery path. It does not merge the candidate or complete the task. The merge guard must prove delivery first.
