---
title: The Pi install offer is suppressed in every session after one failed attempt
status: backlog
score: 0.8
source: "FO review of task ekw, 2026-10-03: the install-attempt sentinel has no session key, and its only measured effect is to suppress the offer."
sprint: pi-ux
sprint-readiness: ready
id: 271f46crset81jwf2rgast0c
---

## Problem

A failed First Officer install writes one file:

    ${TMPDIR:-/tmp}/spacedock-install-attempted

The file has no session key and no project key. The install-and-resume flow in
`skills/first-officer/references/fo-install.md` offers the install only when the
file is absent. The First Officer writes the file before it runs the install, so
both a failed install and a good one consume the single attempt.

The file has one effect. After one failed install, no later Pi session on this
machine receives the install offer, in any project. The First Officer prints the
manual command instead and says the file exists. The operator must delete a temp
file to get the automatic path back. The operator has no reason to know that the
file is there.

The file does not bound a loop. Step 5 of the same flow says "Fall back, never
loop". The First Officer runs the install one time and then stops. A stored
attempt count is not necessary for that rule. When the install is good, the
binary is present and no session gets an offer. So the file is a failure latch
with machine-wide scope, and nothing else.

The runtime cannot supply a session key for the file. The first officer's root
shell carries no `PI_SUBAGENT_PARENT_SESSION`, and that variable names a detached
child's parent. Task ekw measured this and is superseded by this task.

## Visible value

After one failed install, a new Pi session gets the install offer again. The same
session gets no second offer, because the flow still stops after a failure.
Measured against baseline: before, a session that starts after one failed install
prints the manual command and no offer; after, that session gets the offer.

## Out of scope

- The install commands, the channel selection, and the sandbox rule.
- Session identity. No key is necessary when the file is absent. Do not build a
  session-scoped marker in its place.
- The convergence step and the version re-check.
- The three failing repository tests on the Codex path.

## Expected surface and tolerance

1. `skills/first-officer/references/fo-install.md`: remove the sentinel step and
   the sentinel words from the step-4 fallback message.
2. `internal/contractlint/version_gate_smoke_test.go` line 79: remove the guard.
3. `skills/integration/testdata/version_gate_flow.sh` line 17: remove the sentinel.
4. `skills/integration/version_gate_fixture_test.go` line 126: remove the sentinel.
5. `docs/site/get-started/install.md`: keep the manual command. Remove the sentinel
   text and the "re-enables the offer" sentence.

Estimate net LOC change: about -25, across 5 files. Tolerance: net 0 to -60, and
at most 6 files. A larger change needs a revised gate.

Semantic changes: the install offer no longer has an attempt file. No other
behaviour changes.

## Acceptance criteria

**AC-1 (VALUE) - A new session gets the offer after an earlier failed install.**
Verified by: a test that runs the install flow, fails the install, then starts a
new session and requires the offer. Falsifier: restore the sentinel check and the
new session must lose the offer.

**AC-2 - The same session cannot start a second install.**
Verified by: the existing "never loop" behaviour. The flow stops after one attempt
with the manual command. Falsifier: add a retry to the step-5 fallback and require
the flow to stop instead.

**AC-3 - No file under TMPDIR records an install attempt.**
Verified by: the flow fixture test must pass with no sentinel file present.
Falsifier: write the file and require the test to fail.

**AC-4 (no-regression) -** The install commands, the channel selection, the sandbox
rule and the documented manual command do not change, and the sentinel-related
tests in `internal/contractlint` and `skills/integration` pass.

## Test plan

Check first whether a test requires this rule: "a bad install stays suppressed".
If such a test is there, it protects the fault. Remove it in the same change.
Primary proof owner: `skills/integration/version_gate_fixture_test.go` together
with the `internal/contractlint` version-gate smoke test. Deterministic tests
only. No live run.
