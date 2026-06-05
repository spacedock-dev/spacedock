---
id: bqqr8vzz152n8bk2dqf1rw4q
title: Release-gate fix — job-separation + tag-body OPTION B so the 0.19.6/0.20.0 cut doesn't fail like 0.19.5
status: ideation
source: "captain + handoff (2026-06-05) — THE cut blocker, hard prerequisite for the main-flip/0.20.0 milestone. `release.yml` hard-requires a Runtime-Live-E2E on `next` that never auto-runs there, so the cut fails (like 0.19.5). Grounded analysis: Workflow task w569jug9c (session #08)."
score: "0.42"
started: 2026-06-05T19:05:55Z
completed:
verdict:
worktree:
issue:
---

`release.yml` hard-requires a Runtime-Live-E2E run on `next` that never auto-runs there, so a release cut fails the gate (this is what failed the 0.19.5 cut and will fail 0.20.0). This is the hard blocker before the main-flip milestone.

## Direction (for ideation — grounded in the w569jug9c analysis)

- **JOB-SEPARATION fix:** add a SIBLING `journey-ledger` job that goreleaser does NOT `needs:`. NOT "move steps after goreleaser" — that breaks `internal/release/workflow_exec_guard_test.go`. The separation lets the live-e2e/journey-ledger gate exist without blocking the goreleaser cut on a run that doesn't auto-fire on `next`.
- **Tag-body OPTION B:** codify single-`-m`/awk for the tag body, and fix the now-false Go doc/tests that assume the old shape.
- **Coordinate with the Codex peer:** 4n owns the journey-ledger CONSUMER — this fix must align with that consumer; do not author across the boundary, coordinate.
- Spike-first: the riskiest unknown is the goreleaser/needs DAG + the workflow_exec_guard_test interaction — exercise the smallest end-to-end of the separated-job DAG (a dry cut, or the guard test against the new YAML) before the full plan.

## Out of scope

The actual 0.20.0 cut + marketplace flip — that's the `main-flip-0200-marketplace` milestone (this fix is its prerequisite).

## Acceptance criteria

**AC-1 — a release cut succeeds with the live-e2e/journey-ledger gate present but not blocking goreleaser on a never-fired run.**
Verified by: the separated-job DAG (goreleaser does not `needs:` the journey-ledger sibling) proven by a dry/dry-run cut or the workflow exec-guard test against the new `release.yml`; `internal/release/workflow_exec_guard_test.go` stays green (not broken by the change).

**AC-2 — the tag body uses OPTION B and the Go doc/tests match.**
Verified by: the tag-body construction (single-`-m`/awk) tested against the produced tag (the notes land in the body, not folded into the subject — cf. the `7h` release-notes catch), and the now-corrected Go doc/tests pass.

## Test plan

The workflow exec-guard test + a dry cut to exercise the DAG; the tag-body test over the produced tag. Coordinate the journey-ledger consumer alignment with the Codex peer (4n). High-stakes CI/release machinery → detached adversarial audit before merge.
