---
title: Run known live behavior gaps as strict XFAIL
status: backlog
source: "Captain decision after live-test-truth sprint close, 2026-08-09"
started:
completed:
verdict:
score: 0.95
worktree:
issue:
pr:
mod-block:
sprint: test-behavior-completeness
group: common-evidence
sprint-readiness: ready
id: ts7gq0mr9s3chx2w4wppd1kt
---

## Problem

A `liveTODO(...)` stops before the fixture or journey runs. It records ownership,
but it gives no current behavior evidence. It also cannot report that a product
repair made the journey pass.

Go has no native XFAIL result. A broad inversion of the test process can hide
authentication, launch, timeout, fixture, and unrelated assertion failures.

## Value

Known product gaps run in their required live lanes. CI reports the expected
failure, rejects a different failure, and detects a repaired journey immediately.

## Required semantics

A strict XFAIL applies only at the durable semantic assertion boundary.

- PASS means that a normal journey passed.
- XFAIL means that the named semantic failure occurred.
- XPASS means that an expected failure disappeared. XPASS fails the lane until
  the source binding is removed.
- FAIL means that infrastructure failed or a different semantic failure
  occurred.

The source binding owns the target, active task ID, and stable failure code. The
desired-state registry does not copy this actual-state information.

Keep TODO only when the journey cannot run. Examples include a missing adapter,
fixture, selector, or authentication path.

## Acceptance criteria

**AC-1 (VALUE) — Sonnet and Codex run the known headless behavior gap.**

Replace the two `default-headless-gate-stop` TODO cells with strict XFAIL
bindings. Each lane must run the real fixture, host, exercise, and durable
assertion.

Verified by: each lane reports
`implementation-worker-not-dispatched` as XFAIL for task `98a`. Neither lane
reports the journey as skipped.

**AC-2 (VALUE) — A repaired journey reports XPASS.**

If the durable assertion passes for an XFAIL target, the live lane must fail
with an XPASS result. The result must name the journey, target, and active task.

Verified by: one focused exercise uses the existing journey boundary and proves
that the pass cannot remain hidden as XFAIL.

**AC-3 — Unexpected failures remain failures.**

Authentication, launch, timeout, fixture, and artifact failures must fail
normally. A semantic failure with a different stable code must also fail.

Verified by: code inspection and the existing live failure paths. Do not add a
new meta-test suite for the test runner.

**AC-4 — Reconciliation reports execution state and checks ownership.**

Reconciliation must distinguish TODO from XFAIL. The mutable owner join must
check active owners for both states.

Verified by: the existing source reconciliation and mutable owner join. Do not
add a copied scenario table or recorded gap ledger.

**AC-5 — Journey metrics include the strict outcome.**

The existing journey result must include `pass`, `xfail`, `xpass`, or
`fail`. An XFAIL result must include its owner and stable failure code.

Verified by: the normal live artifact for the converted headless journey. Do not
create a second metrics format.

## Scope

- Add the smallest grade result needed at the semantic assertion boundary.
- Preserve immediate failures for infrastructure and setup errors.
- Convert the two evidenced `98a` cells in the same change.
- Keep the registry as desired state.
- Use the existing journey metrics artifact.

## Out of scope

- Inverting the exit status of the complete test process.
- Matching unstable CLI text or complete error strings.
- Converting an unexecuted TODO without stable failure evidence.
- Adding tests that only test test infrastructure.
- Repairing the product behavior owned by `98a`.
