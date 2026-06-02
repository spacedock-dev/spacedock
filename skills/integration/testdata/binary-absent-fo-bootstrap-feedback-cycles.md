<!--
Vendored fragment of the archived #262 entity (binary-absent-fo-bootstrap)
"## Feedback Cycles" section, copied verbatim from
docs/dev/.spacedock-state/_archive/binary-absent-fo-bootstrap/index.md (the
state-checkout orphan branch, which is NOT present in a code-only worktree).

This is the external prior artifact that AC-3 of encode-deliverable-principles
grounds its M1/M2 audit citation against. It is vendored here so the
cross-check ALWAYS runs in CI (no skip-on-absence hole), and so a future drift
between the contract's citation and the real #262 record fails the test rather
than being silently skipped. Re-verify against the archived entity if #262's
record is ever amended.
-->

## Feedback Cycles

**Cycle 1 (FO, 2026-06-02) — validation recommended PASSED; detached adversarial audit found two test-strength holes; routed to implementation for one-line tightenings before merge.**

The shipped prose is correct — the audit refuted nothing material on the contract edit (install lines verbatim vs README:27/37; scope surgical, only step 1's abort clause; AC-3 single-source confirmed by an independent grep). But two `skills/integration/contract_gate_test.go` assertions under-pin their own docstrings, so the suite would green-light a future regression in either direction:

- **M1 (`contract_gate_test.go:144`)** — the Class A no-doctor check is guarded by `if n := strings.Count(classA, doctor); n > 0`. If a future edit deletes the `Do NOT run \`spacedock doctor\`` prohibition entirely (zero doctor mentions), the check is skipped and the test passes — so it does NOT pin that Class A carries the no-doctor guidance. **Fix:** require the prohibition string present in Class A (so its removal fails), then keep the route-count guard for the present case.
- **M2 (`contract_gate_test.go:151`)** — AC-2 is a bare `strings.Contains(classB, "spacedock doctor")`, satisfied by a negated/disclaimer mention (audit verified an edit replacing the Class B route with `(Historically we suggested spacedock doctor but no longer.)` passes). **Fix:** assert the live-route phrasing (e.g. `run \`spacedock doctor\`` or `for the per-class remedy`) so a gutted route phrased as a disclaimer fails.

**Required proof:** each tightened test must FAIL on the auditor's adversarial edit (Class A prohibition deleted; Class B route replaced by a disclaimer) and pass on the real file. Mutation-verify, do not re-read.
