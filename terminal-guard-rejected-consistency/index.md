---
id: 6b8k79kbmzn8n48g2amf6q4m
title: Terminal-guard verdict=rejected consistency — align contract prose + the --set/--archive asymmetry
status: validation
source: sprint-end antipattern reviews (2026-06-01) — 0.19.3 minor-findings bucket
started: 2026-06-03T07:09:59Z
completed:
verdict:
score: "0.26"
worktree: .worktrees/spacedock-ensign-terminal-guard-rejected-consistency
issue:
---

Two related inconsistencies in how the terminal-transition guard treats `verdict=rejected`.

## Problem statement

The merge-hook terminal guard refuses a terminal transition / archival when a `merge` hook is registered AND `pr` is empty AND `mod-block` is empty (force / `merge: local` are the documented escapes). The `--set` surface adds one more escape — `verdict=rejected` — that `--archive` does not, and the FO contract prose names neither surface's verdict escape.

### Inconsistency 1 — `--set` vs `--archive` asymmetry (the bug)

`--set` exempts `verdict=rejected`:

> `handlers.go:163` — `if !force && policy != mergeLocal && isTerminalUpdate() && modBlock == "" && postUpdatePR == "" && postUpdateVerdict != "rejected" {`

`--archive` does NOT — `mutate.go:324` has the same guard MINUS the verdict clause:

> `mutate.go:324` — `if !force && policy != mergeLocal && modBlock == "" && pr == "" {`

So on a `verdict=rejected` entity with empty `pr`/`mod-block` under the default `merge: pr` policy, reject-then-terminalize-then-archive HALF-PASSES: the terminal `--set` goes through (exit 0) and the immediately following `--archive` refuses (exit 1) — the entity stalls terminalized-but-unarchived, and the FO's only escape is `--force`, which the contract says "is never part of the happy path."

**Empirically confirmed** (scratch run on a copy of `testdata/merge-pr-workflow` with a `verdict: rejected` entity, `pr`/`mod-block` empty):
- `status --set 040-rejected status=done` → `status: implementation -> done`, exit 0
- `status --archive 040-rejected` → `Error: entity 040-rejected cannot be archived — workflow has merge hook(s) [local-merge] that have not run (pr field is empty and mod-block is empty)...`, exit 1

### Inconsistency 2 — contract prose vs code drift

The FO operating-contract enumerations of the terminal guard never state the `verdict=rejected` exemption that `--set` actually applies:

- `first-officer-shared-core.md:214` — "`status --set` and `status --archive` refuse terminal updates while merge hooks exist with both `pr` and `mod-block` empty." (names no verdict escape; implies both surfaces are identical)
- `first-officer-shared-core.md:324` (Mod-Block Enforcement) — "`status --set` and `status --archive` also refuse terminal transitions and archival when merge hooks ... are registered AND `pr` is empty AND `mod-block` is empty. `--force` bypasses. `merge: local` exempts only the pr-requirement..." (names `--force` and `merge: local` as the only escapes; `verdict=rejected` unmentioned)

A reader trusting the prose believes a rejected entity is blocked at `--set` exactly as at `--archive`. The code disagrees on `--set`.

## Target consistent behavior

**Exempt both surfaces on `verdict=rejected`.** A rejected entity never went through the merge ceremony — there is no PR to require, no merge to gate on — so the merge-hook requirement is vacuous for it. The guard's purpose is to stop an entity from terminalizing as *accepted/merged* without the ceremony; a *rejected* entity is the opposite case. `--set` already encodes this; `--archive` should match. Add the `verdict != "rejected"` clause to `mutate.go:324` so both surfaces let a rejected entity through. The mod-block-pending guard (`mutate.go:306`, policy-independent) is untouched — a rejected entity with a live mod-block still refuses, preserving the ceremony-separation invariant.

Direction rationale (vs the alternative "guard both" — drop the `--set` exemption): guarding both would force every rejected entity through a no-op merge ceremony (set mod-block → invoke a hook with nothing to merge → clear → archive) or `--force`, which contradicts "`--force` is never part of the happy path." Exempting both keeps reject-and-archive on the happy path. **The gate / captain confirms direction; both ACs below are written for "exempt both."**

## Acceptance criteria

**AC-1 — `--set` and `--archive` agree on `verdict=rejected`.** Under the default `merge: pr` policy with a registered merge hook, a `verdict=rejected` entity with empty `pr`/`mod-block` terminalizes via `--set` (exit 0) AND archives via `--archive` (exit 0) — both succeed without `--force`. (Per chosen direction "exempt both"; if the gate flips direction to "guard both," both surfaces refuse with exit 1 instead — the test asserts the surfaces AGREE either way.)
Verified by: `internal/status/merge_policy_guard_test.go` — a new `TestRejectedVerdictArchiveMatchesSet` (Go unit test driving the native binary on a new `merge-pr-workflow` fixture entity carrying `verdict: rejected`, empty `pr`/`mod-block`). The test runs `--set ... status={terminal}` then `--archive` on that entity and asserts both exit codes match the `--set` outcome. On today's code this test FAILS: `--set` → exit 0, `--archive` → exit 1 (confirmed by the scratch run above; the new fixture + assertion is the regression test).

**AC-2 — mod-block-pending still refuses for a rejected entity.** A `verdict=rejected` entity with a non-empty `mod-block` is still refused by `--archive` (the policy-independent block survives; the verdict escape relaxes only the merge-hook pr-requirement, mirroring how `merge: local` relaxes only the pr-requirement).
Verified by: `merge_policy_guard_test.go` — assertion (extend the new fixture set or add a sibling case) that `--archive` of a rejected-but-mod-block-pending entity exits 1 naming the pending mod-block.

**AC-3 — Contract prose states the `verdict=rejected` exemption.** The FO terminal-guard enumerations (`first-officer-shared-core.md:214` and the Mod-Block Enforcement bullet `:324`) name the `verdict=rejected` escape alongside `--force` and `merge: local`, so prose matches code on BOTH surfaces.
Verified by: the contract reference text states it (prose-presence is legitimate here — the contract text IS the deliverable, per the workflow's legitimacy boundary; the behavioral claim it documents is independently gated by AC-1).

## Test plan

- **AC-1 / AC-2:** Go unit tests in `internal/status/merge_policy_guard_test.go`, following the existing `assertMergeGolden` harness (golden-envelope freeze + exit-code assertion). New fixture entities under `internal/status/testdata/merge-pr-workflow/` (e.g. `040-rejected.md` with `verdict: rejected`, and a `050-rejected-pending.md` with `verdict: rejected` + `mod-block: merge:local-merge` for AC-2). Cost: low — mirrors `TestMergePrDefaultNoSentinelArchiveStillRefuses` line-for-line, just with a rejected-verdict entity and the inverted expectation. New goldens needed for the new envelope names.
- **AC-3:** prose edit to the two enumerations; verification is the reference text itself (no behavioral test — the behavior is AC-1's job).
- Fixture vs CLI vs live: Go unit/fixture tests only. No live workflow test needed — the claim is command-level guard behavior, fully exercised by driving the native binary against a fixture.

### Riskiest-unknown check

"What would invalidate the rest of this work if it broke?" — the claim that the asymmetry actually manifests as `--set` passes / `--archive` refuses on a default-policy rejected entity. **Exercised first** (scratch run above on a copy of `merge-pr-workflow`): confirmed `--set` exit 0, `--archive` exit 1. The fix itself (adding one clause to the `mutate.go:324` condition) composes only already-proven guard machinery — same `resolveMergePolicy`/`scanMods` path `--set` already uses. No further spike needed beyond the confirmed asymmetry.

## Notes
- `internal/status` serialized lane (handlers.go/mutate.go) — coordinate with #251 and any other status-lane implementation (one impl worktree at a time). The contract-prose edit is protected scaffolding → dispatched worker only, authored in the vendored copy.
- The entity body's earlier reference to "Python oracle `:2648`" is stale — there is no Python oracle in the tree; the Go implementation in `internal/status` is authoritative. Dropped from the problem statement above.

## Stage Report: ideation

- DONE: Pin the exact inconsistency between contract prose and the --set/--archive asymmetry on verdict=rejected, name the target consistent behavior, and the test that fails on the asymmetry.
  Asymmetry pinned to `handlers.go:163` (has `postUpdateVerdict != "rejected"`) vs `mutate.go:324` (lacks the verdict clause); prose drift pinned to `first-officer-shared-core.md:214` and `:324`; target = exempt both surfaces; failing test = `TestRejectedVerdictArchiveMatchesSet` in `merge_policy_guard_test.go` (--set exit 0 / --archive exit 1 on today's code).

### Summary

Confirmed the asymmetry empirically (scratch run on a copy of `testdata/merge-pr-workflow` with a `verdict: rejected` entity, empty `pr`/`mod-block`, default `merge: pr`): `--set status=done` → exit 0, `--archive` → exit 1 — reject-then-archive half-passes. Pinned the bug to the one missing `verdict != "rejected"` clause at `mutate.go:324` and the prose drift to two `first-officer-shared-core.md` enumerations (lines 214, 324) that name only `--force`/`merge: local` as escapes. Chose direction "exempt both" (a rejected entity never ran the merge ceremony, so the requirement is vacuous; "guard both" would push rejected entities onto a `--force` path the contract forbids), wrote three ACs (set/archive agree; mod-block-pending still refuses; prose names the escape) with a Go-unit-test plan mirroring the existing `assertMergeGolden` cases, and recorded the riskiest-unknown as already exercised. Flagged the stale "Python oracle :2648" reference as nonexistent and dropped it.

## Stage Report: implementation

- DONE: Make --set and --archive AGREE on verdict=rejected (chosen direction: exempt both). Add the missing `verdict != "rejected"` clause to the terminal-guard at internal/status/mutate.go.
  Added the clause at the runArchive merge-hook guard (was mutate.go:343 on next; sibling already present at handlers.go:203). Commit f07babee.
- DONE: Author the failing Go test FIRST (TDD): TestRejectedVerdictArchiveMatchesSet in internal/status/merge_policy_guard_test.go.
  New fixture 040-rejected.md (verdict: rejected, empty pr/mod-block); test stages one root, runs --set status=done then --archive on the same entity, asserts exit codes agree. Confirmed RED on pre-fix code (set=0/archive=1 disagreement) and on a fix-revert (golden+exit mismatch); GREEN after fix.
- DONE: AC-2 — a rejected entity with non-empty mod-block is STILL refused by --archive (policy-independent block survives).
  New fixture 050-rejected-pending.md (verdict: rejected, mod-block: merge:local-merge); TestRejectedVerdictModBlockPendingArchiveRefuses asserts --archive exits 1 naming the pending mod-block.
- DONE: Align contract prose (AC-3): the two FO terminal-guard enumerations must name the verdict=rejected escape alongside --force and merge: local.
  Updated first-officer-shared-core.md at the terminalize-step bullet (:214) and the Mod-Block Enforcement bullet (:324) to name verdict=rejected on both surfaces. Commit f07babee.

### Summary

Added the one missing `verdict != "rejected"` clause to runArchive's merge-hook guard in internal/status/mutate.go so --archive matches the --set finalize path; both surfaces now let a rejected entity (empty pr/mod-block, default merge: pr) through without --force, while the policy-independent mod-block-pending guard still refuses a rejected entity with a live block. TDD strictly: wrote TestRejectedVerdictArchiveMatchesSet first (drives the native runner --set-then-archive on one staged root, asserts surfaces agree) and confirmed it reds on today's code AND on a fix-revert (golden + exit-code mismatch), then green after the fix. Aligned both FO contract enumerations to name the verdict=rejected escape. Full suite green: go build/vet clean, 980 tests passed.

## Stage Report: validation

- DONE: Independently reproduce AC-1 and AC-2 from a fresh run; confirm AC-1 would RED on a fix-revert (guard clause genuinely load-bearing).
  `go test ./internal/status/ -run 'TestRejectedVerdictArchiveMatchesSet|TestRejectedVerdictModBlockPendingArchiveRefuses'` → 2 passed. AC-1 RED on fix-revert (removed verdict clause at mutate.go:347): exit mismatch native=1 golden=0. Both succeed without --force on today's code.
- DONE: AC-2 — rejected + non-empty mod-block still refused by --archive, exit 1 naming the block.
  `TestRejectedVerdictModBlockPendingArchiveRefuses` GREEN (golden merge-pr-rejected-pending-archive: exit 1, stderr names `pending mod-block (merge:local-merge)`). Also confirmed it stays green with the verdict clause reverted — proving the escape is narrow (relaxes only the merge-hook pr-requirement, not the policy-independent mod-block guard).
- DONE: AC-3 — both FO contract enumerations name the verdict=rejected escape; guard clause landed at runArchive mirroring the handlers.go sibling; full offline suite green.
  Contract: `first-officer-shared-core.md:214` (terminalize-step) and `:324` (Mod-Block Enforcement) both name `verdict=rejected` alongside `--force`/`merge: local`. Code: handlers.go:203 (`--set`) and mutate.go:347 (`--archive`) mirror each other. Only one contract copy exists (no vendored duplicate missed). `go test ./...` → 980 passed/12 pkgs, build+vet clean, raw exit 0.
- DONE: Reject if the new guard test is weak — prove it reds on a regression, not merely passes on today's code.
  Decisive proof: reverted the clause AND regenerated goldens (`-update`) so the golden-freeze passes against broken behavior — the test STILL failed at merge_policy_guard_test.go:173 (`--set and --archive must agree on verdict=rejected, got set=0 ... archive=1`). Test strength rests on the behavioral agreement assertion, not the golden spelling-freeze. All files restored to clean HEAD (git status: ok).

### Summary

PASSED. All 3 ACs reproduced from fresh runs with evidence outside the task body. AC-1/AC-2 are real Go unit tests driving the native binary; AC-3 is a contract-prose deliverable whose behavioral claim is independently gated by AC-1. The guard test is genuinely strong: it reds on a fix-revert even with goldens regenerated to match the broken behavior, because the load-bearing assertion is the explicit `setCode != archiveCode` agreement check, not the golden freeze. handlers.go (--set) and mutate.go (--archive) guards now mirror each other; both contract enumerations match code. Note for the FO: this is a HIGH-STAKES status mutation/guard surface — recommend the detached adversarial audit before merge per the validation stage's policy.

### Feedback Cycles

**Cycle 1 (2026-06-04) — detached adversarial audit, material coverage hole.** The detached audit (`audit-terminal-guard-rejected-consistency`, checkout f07babee) refuted the validation's test-strength: the `--archive` merge-hook exemption is not regression-proofed. Adversarial edit — widen `verdict != "rejected"` to `verdict == ""` at `internal/status/mutate.go:347` (so any non-empty verdict, e.g. `verdict=passed`, escapes the guard) — left `go test ./internal/status/` fully GREEN. The `--set` side is pinned (caught) but `--archive` is not; no test drives a non-rejected non-empty-verdict entity through `--archive` expecting refusal. CODE is correct; coverage hole only. Routed back to implementation: add a fixture mirroring `040-rejected` with `verdict: passed` + an assertion that `--archive` exits 1 (must red under the widening probe, green on correct code).

## Stage Report: implementation (cycle 2)

- DONE: Close the audit's --archive coverage hole — regression-proof the verdict exemption against over-widening.
  Added fixture 060-passed-nosentinel.md (verdict: passed, empty pr/mod-block, default merge: pr) and TestNonRejectedVerdictArchiveStillRefuses asserting --archive exits 1 ("cannot be archived"). Commit 18f6fb9b.
- DONE: Prove the new test bites — RED under the `verdict==""` over-widening probe, GREEN on correct `verdict != "rejected"` code; existing tests stay green.
  Verified: under the probe the new test fails (native exit 0 vs golden 1) while the pre-existing 020-no-sentinel default-refuse test stays green (confirming it gave no signal — the audit's root cause). Restored code; go build/vet clean, go test ./... 981 passed (+1).

### Summary

Cycle 1 feedback (detached adversarial audit) found the --archive verdict exemption was not regression-proofed: widening `verdict != "rejected"` to `verdict == ""` (any non-empty verdict escapes the merge-hook guard) left the suite green, because the only default-refuse archive fixture had an empty verdict. Added a verdict=passed no-sentinel fixture (060) and TestNonRejectedVerdictArchiveStillRefuses asserting --archive still refuses it. Proved it bites: RED under the over-widening probe (with the pre-existing 020 test staying green, demonstrating the gap), GREEN on the correct code. Code unchanged — this is test-only coverage. Full suite 981 passed.
