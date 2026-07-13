---
title: "Independent pre-cut antipattern audit for v0.25.0"
status: implementation
source: "Captain directive 2026-07-13: dispatch the pre-cut antipattern audit now."
score: 1.0
worktree: .worktrees/spacedock-ensign-precut-antipattern-audit-0250
id: 1nv6hngvcx5gj67sqtdxsc2b
started: 2026-07-13T08:23:30Z
---

## Problem

The 0250 FO behavioral-discipline release content is merged on `main` and
`v0.25.0-pre2` is published. Before the stable tag, the release contract
requires an independent pre-cut audit of the assembled surface. The audit must
find real ship blockers, not merely confirm that the intended wording exists.

## Audit scope

Review the assembled `v0.24.0..main` release delta and the release-critical
surfaces it changes: FO contract and host adapters, launcher/bootstrap behavior,
runtime-live evidence, and the stable-cut prerequisites. Trace each concern to
code or observable behavior. Treat the live wrong-root incident and the
unstarted `7v` coordination as release context, but do not implement either.

## Acceptance criteria

**AC-1 (independent behavioral audit):** The report identifies the reviewed
release surfaces and supports every conclusion with a concrete source citation,
runtime artifact, or smallest relevant command result; instruction-text grep or
Contractlint presence alone is never evidence of a behavioral claim.

**AC-2 (actionable verdict):** The audit records one of CLEAN-TO-CUT,
BLOCKED, or CLEAN-WITH-FOLLOWUPS; every blocker names the exact release risk and
the smallest corrective action, while non-blockers are clearly separated as
followups.

**AC-3 (scope discipline):** The auditor does not modify product code, tests,
release configuration, or shipped contract text. It records findings and
evidence in this task's stage report only.

## Test plan

- Start from the actual `v0.24.0..main` diff and follow each high-risk change to
  its behavior and existing evidence.
- Run only focused, decision-relevant checks needed to test a suspected
  antipattern; do not add new Contractlint checks or rerun unrelated suites as
  ceremonial reassurance.
- Verify the release preconditions in `docs/releasing.md`, including the
  stamped-commit/exact-SHA live-E2E ordering.

