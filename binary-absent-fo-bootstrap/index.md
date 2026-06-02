---
id: 5vb6mh9kewyh0p68r93mf6m1
title: Harden binary-absent FO startup — install hint in the contract-gate abort, don't route to the missing binary's doctor
status: validation
source: "captain (2026-06-02) — live install-path test (plugin present, binary unlinked): the FO aborts cleanly but its install guidance is model-improvised, and the contract routes the binary-absent case to `spacedock doctor`, which is itself the missing binary"
started: 2026-06-02T14:37:38Z
completed:
verdict:
score: "0.28"
worktree: .worktrees/spacedock-ensign-binary-absent-fo-bootstrap
issue:
---

When a Spacedock plugin is present but the `spacedock` binary is absent from PATH, the FO contract-version gate (`first-officer-shared-core.md` §Startup step 1) runs `spacedock --version`, gets "command not found", and aborts startup — correctly, no silent failure. But the captain's live `claude -p` test surfaced two weaknesses in what happens next:

- **Binary-absent is routed to `spacedock doctor`.** The gate's abort instruction says "run `spacedock doctor` for the per-class remedy." For the *binary-absent* class that is a dead end — `doctor` is the same missing binary. A strong model worked this out and didn't loop, but a weaker model could loop on `doctor` or just print "command not found." `doctor` is the right remedy only for the *binary-present / version-mismatch* class (where the binary can run).
- **Install guidance is model-improvised, not contract-guaranteed.** `first-officer-shared-core.md` carries zero install guidance. The good behavior the captain saw (concrete `brew` + `go build ./cmd/spacedock` steps) was improvised by the model, not guaranteed by the contract. The contract should carry the install hint so even a weak model emits a correct, runnable remedy.

This is the FO-contract sibling of the `init→install` rename drift + unknown-subcommand-silent-exit code fixes folded into `1x` (same live install-path test). It is scaffolding: it edits the shipped FO contract that governs every dispatched session.

## Problem restated

`first-officer-shared-core.md` §Startup step 1 currently funnels three failure modes into ONE abort instruction with ONE remedy:

> If the token is outside the range (binary too old: `<N>` below the lower bound; plugin too old: `<N>` at or above the upper bound), **or `spacedock --version` is unavailable**, ABORT startup with the actionable mismatch message and run `spacedock doctor` for the per-class remedy …

The `or spacedock --version is unavailable` clause is the binary-absent class, and for it `spacedock doctor` is a dead end — `doctor` is the same absent binary. The other two clauses (too old / too new) are the binary-*present* class: a `contract <N>` token was successfully parsed, so the binary ran, so `spacedock doctor` is a live, correct remedy. The fix is to split step 1's single abort sentence into these two classes and give the binary-absent class a self-contained, runnable install hint instead of the doctor bounce.

## Design — two-class abort split

Step 1's gate stays one logical check (`run spacedock --version`, parse `contract <N>`, compare to `>=1,<2`), but its **abort branch** splits by whether the binary produced a parseable version token:

**Class A — binary absent / non-executable (no `contract <N>` token because the command did not run):** `spacedock --version` is "command not found", non-executable, or emits no parseable `contract <N>` token. ABORT with a message that carries a runnable install hint inline and does NOT route to `spacedock doctor` (the binary that doctor would invoke is the one that is missing). The install hint is two concrete, copy-runnable lines drawn verbatim from the repo README's Install section so they cannot drift from the real install path:

- Released lane: `brew install spacedock-dev/homebrew-tap/spacedock`
- Source lane: `go build -o spacedock ./cmd/spacedock` (run from a `git clone` of the repo)

**Class B — binary present but contract out of range (token parsed, `<N>` outside `>=1,<2`):** the binary ran and reported a contract version, but it is below the lower bound (binary too old) or at/above the upper bound (plugin too old). ABORT with the existing actionable mismatch message and KEEP `run spacedock doctor for the per-class remedy` — `doctor` can run and is the correct triage for a present-but-mismatched binary.

Both classes still ABORT before discovery / `--boot`; the split only changes the remedy text per class, never the no-silent-failure guarantee.

### Exact before/after for step 1

The single sentence to change is the abort clause. The current text is:

> If the token is outside the range (binary too old: `<N>` below the lower bound; plugin too old: `<N>` at or above the upper bound), or `spacedock --version` is unavailable, ABORT startup with the actionable mismatch message and run `spacedock doctor` for the per-class remedy — do NOT proceed to discovery or `--boot`.

Replace with two-class wording of this shape (implementation may tune prose, but must preserve: the brew + go-build lines verbatim from README for Class A, the absence of any `doctor` route in Class A, and the retained `spacedock doctor` route in Class B):

> Abort by class. **Binary absent or non-executable** (`spacedock --version` is not found, or emits no parseable `contract <N>` token): ABORT startup and tell the operator the `spacedock` binary is not on PATH, with the runnable install hint — released lane `brew install spacedock-dev/homebrew-tap/spacedock`, or source build `go build -o spacedock ./cmd/spacedock` from a clone of the repo. Do NOT run `spacedock doctor` for this class — `doctor` is the same missing binary. **Binary present but contract out of range** (token parsed, `<N>` below the lower bound = binary too old, or at/above the upper bound = plugin too old): ABORT startup with the actionable mismatch message and run `spacedock doctor` for the per-class remedy. In every class, do NOT proceed to discovery or `--boot`.

## Sync scope (AC-3)

Verified during ideation: `agents/first-officer.md` does NOT mirror the startup-gate text — it delegates ("Then begin the Startup procedure from the shared core"). A repo-wide grep for the gate markers (`spacedock --version`, `Contract version gate`, `per-class remedy`, `spacedock doctor`) across `skills/` and `agents/` returns exactly ONE prose source: `first-officer-shared-core.md`. (`skills/integration/marketplace_manifest_test.go` matches `spacedock doctor` only in a comment about the contract-range bracketing test; it is a Go test, not a mirror of the gate prose, and is out of scope.) Therefore the sync scope is: **edit `first-officer-shared-core.md` step 1 only; assert by grep that no second prose copy of the abort guidance exists.** AC-3 is satisfied by that single-source invariant, not by reconciling two diverging copies.

## Acceptance criteria

**AC-1 — binary-absent emits a runnable install hint, not a doctor bounce.** In `first-officer-shared-core.md` §Startup step 1, the binary-absent / non-executable abort class contains both runnable install lines verbatim (`brew install spacedock-dev/homebrew-tap/spacedock` and `go build -o spacedock ./cmd/spacedock`) and contains no `spacedock doctor` route within that class.
- Verified by: a test asserting both install strings are present in step 1, AND that the binary-absent class text does not contain `spacedock doctor`. See AC-1 framing note below — this is a presence/absence check over the contract text where the text IS the deliverable, not a behavioral claim a code gate could enforce.

**AC-2 — doctor stays the remedy for the binary-present version-mismatch class.** The binary-present-but-out-of-range abort class in step 1 still routes to `spacedock doctor`.
- Verified by: a test asserting the step-1 text retains `spacedock doctor` associated with the version-out-of-range / present-binary class (the string survives and is attached to Class B, distinguishing it from a blanket deletion).

**AC-3 — single-source startup-gate guidance, no drift.** The abort guidance lives in exactly one prose file (`first-officer-shared-core.md`); `agents/first-officer.md` continues to delegate to the shared core rather than restating the gate.
- Verified by: a grep over `skills/` and `agents/` for the gate markers returning `first-officer-shared-core.md` as the sole prose source (the integration `.go` test comment excepted), AND `agents/first-officer.md` still containing only the "begin the Startup procedure from the shared core" delegation.

### AC-1 framing — why a presence/absence check is the right proof here, not a behavioral gate

The workflow forbids ACs whose only proof is "the instruction text says to do X" *when X is a behavioral guarantee a code gate should enforce* — there the ceiling is wording-is-present and the real proof belongs in a binary guard or a failing test. This entity is the legitimate exception, and the reason is mechanical, not a plea for leniency: **the binary is absent.** No `spacedock` command can run to emit the install hint, because the thing that would emit it is exactly what is missing. The FO is an agent reading a contract; for this failure mode the contract prose *is* the executable — it is the only artifact present at the moment of failure. So the claim is not "an agent behaves a certain way" (un-checkable here without a live binary); the claim is "the contract text the FO loads carries clause Y and omits clause Z." That is a property *of the text*, and the workflow's own ideation rules name a property-of-the-text check as legitimate proof at the claim's own level ("a presence check over instruction files proving they carry a required clause or stay free of a banned token is proof at the claim's own level"). The proof is therefore a string-presence test (both install lines present) plus a banned-token test (no `spacedock doctor` inside the binary-absent class), exercising the real file — recorded here so the gate does not bounce AC-1 as a bare wording-is-present criterion.

## Test plan

- **Test type:** a single Go test (or extend an existing `skills/`-area contract-text test) that reads `first-officer-shared-core.md`, isolates §Startup step 1, and asserts the AC-1/AC-2/AC-3 string relationships. Fixture cost: low — file read + substring assertions, no binary build, no live workflow. No CLI or live-workflow test is warranted because the deliverable is loaded contract text, not runtime behavior (the binary is absent by definition of the failure mode).
- **AC-1:** assert step-1 text contains `brew install spacedock-dev/homebrew-tap/spacedock` AND `go build -o spacedock ./cmd/spacedock`; assert the binary-absent class substring does not contain `spacedock doctor`.
- **AC-2:** assert step-1 text still contains `spacedock doctor` (attached to the present-but-out-of-range class).
- **AC-3:** assert a grep of `skills/` + `agents/` for the gate markers yields only `first-officer-shared-core.md` as a prose source; assert `agents/first-officer.md` contains the delegation line and not the gate prose. (This is the single-source invariant, expressible as a test or a documented grep in the stage report.)
- **Spike:** no spike needed. The design composes only already-proven mechanisms — it edits static contract prose the FO loads at startup, with no parser round-trip, no new on-disk format, and no runtime handoff. The one mechanism-level question (does any second file mirror the gate text, which would make a single-file edit insufficient?) was resolved during ideation by the repo-wide grep recorded under "Sync scope" above: exactly one prose source. The riskiest claim — "the install hint can only live in the contract text because the binary is absent" — is settled by the framing note above; it is true by construction.

## Test gates / proof discipline

The deliverable is the FO's own loaded instructions — the binary is *absent*, so no `spacedock` command can run to emit the hint; the guarantee must live in the contract text the FO already has. This is the legitimate doc-as-deliverable case (the contract prose *is* the product for an agent-read contract), distinct from the prose-only-AC antipattern the workflow forbids (an AC claiming a *behavioral* guarantee that a code gate should enforce). Proof is a presence test for the two install lines + an absence test for a doctor-route inside the binary-absent class + the single-source grep invariant, all exercising the real file. The AC-1 framing note above records the reasoning so the gate does not bounce AC-1 as a bare wording-is-present criterion.

## Stage Report: ideation

- DONE: The two-class split is concretely specified: exact abort-message text for binary-absent (runnable brew + go-build install hint, NO doctor route) vs binary-present-out-of-range (keeps the doctor route).
  "Design — two-class abort split" + "Exact before/after for step 1" give Class A (brew `spacedock-dev/homebrew-tap/spacedock` + `go build -o spacedock ./cmd/spacedock`, no doctor) and Class B (keeps `spacedock doctor`); install lines drawn verbatim from README:14-39.
- DONE: The prose-as-deliverable AC tension is addressed head-on: AC-1 framed as grep-checkable presence with the reasoning recorded, so the gate does not bounce it as a bare wording-is-present criterion.
  "AC-1 framing" section: binary is absent so no command can emit the hint → the contract text IS the executable → property-of-the-text check is proof at the claim's own level (matches the workflow's own legitimacy rule); test plan makes it a presence + banned-token test over the real file.
- DONE: agents/first-officer.md (and any other mirror of the startup-gate text) is checked for duplicated abort guidance; the sync scope is named so implementation keeps shared-core and the agent prompt in agreement.
  "Sync scope (AC-3)": repo-wide grep of skills/ + agents/ for gate markers returns one prose source (first-officer-shared-core.md); agents/first-officer.md delegates ("begin the Startup procedure from the shared core"), does not mirror. Scope = edit shared-core step 1 only; AC-3 is a single-source invariant.

### Summary

Fleshed out the binary-absent FO-startup hardening into a concrete two-class abort split for `first-officer-shared-core.md` §Startup step 1, with verbatim before/after wording and install lines lifted from the repo README so they cannot drift. Resolved the doc-as-deliverable AC tension by grounding AC-1's presence/absence proof in the fact that the binary is absent (the contract prose is the only artifact present at failure, so a property-of-the-text check is proof at the claim's own level — the workflow's own named exception). Verified by grep that the gate text has a single prose source, narrowing AC-3 to a single-source invariant rather than a two-file reconciliation; recorded "no spike needed" since the design only edits static contract prose.

## Stage Report: implementation

- DONE: The two-class abort split lands in references/first-officer-shared-core.md Startup step 1 per the ideation before/after: Class A (binary absent/non-executable) carries both install lines verbatim (`brew install spacedock-dev/homebrew-tap/spacedock` + `go build -o spacedock ./cmd/spacedock`) with NO `spacedock doctor` route; Class B (present-but-out-of-range) keeps the doctor route.
  Applied the ideation's exact replacement wording to step 1's abort clause; install lines match README:27,37 verbatim. Class A carries `Do NOT run \`spacedock doctor\` for this class`; Class B retains `run \`spacedock doctor\` for the per-class remedy`. Commit b5e5ced7.
- DONE: A Go test enforces AC-1/AC-2/AC-3 over the real contract file (both install strings present, no `spacedock doctor` in the absent-class span, `spacedock doctor` retained for Class B, single-source grep over skills/+agents/) and is green — written failing-first.
  `TestStartupAbortSplitsByBinaryPresence` (AC-1/AC-2) + `TestStartupGateGuidanceHasSingleProseSource` (AC-3) in skills/integration/contract_gate_test.go, reusing the existing `foSharedCore`/`sectionAfter`/`repoRoot` helpers. Confirmed failing pre-edit (class marker absent), green post-edit. Full `go test ./skills/integration/` passes; `go vet` clean. AC-1's banned-token check distinguishes a doctor ROUTE from the prohibition clause, which the prescribed Class A wording itself contains.
- DONE: agents/first-officer.md still only delegates to the shared core (no gate-prose mirror introduced).
  AC-3 test asserts `agents/first-officer.md` retains `begin the Startup procedure from the shared core` and contains no gate markers (`per-class remedy`, `spacedock doctor`); the skills/+agents/ .md walk returns exactly one prose source (first-officer-shared-core.md). The marketplace_manifest_test.go `spacedock doctor` match is a Go comment, excluded by the .md-only scope.

### Summary

Split Startup step 1's single abort clause into two classes using the ideation's exact wording: binary-absent emits the runnable brew+go-build install hint with no doctor bounce; binary-present-out-of-range keeps the doctor route. Added three contract-text assertions to the existing `contract_gate_test.go` (failing-first), green over the real file; full integration package and vet pass. The single-prose-source invariant and the agents/first-officer.md delegation are guarded by test, not just grep.

## Stage Report: validation

- DONE: Reproduce each AC's 'Verified by' evidence: run `go test ./skills/integration/` and confirm TestStartupAbortSplitsByBinaryPresence + TestStartupGateGuidanceHasSingleProseSource are green; confirm they FAIL when the step-1 edit is reverted.
  `go test ./skills/integration/` = 26 passed; both named tests green (`-run` over the pair = 2 passed). Failing-first validated by `git checkout b5e5ced7~1 -- first-officer-shared-core.md`: TestStartupAbortSplitsByBinaryPresence FAILS ("step 1 has no binary-absent class marker"). TestStartupGateGuidanceHasSingleProseSource correctly stays GREEN on the step-1 revert — it is a single-source/no-drift invariant, not a witness of the two-class wording; proven non-trivial by injecting a gate-marker mirror into agents/first-officer.md → it FAILS as designed. File restored, tree clean.
- DONE: Verify the contract change is surgical and correct: only step 1's abort clause changed; Class A carries both install lines AND no `spacedock doctor` route; Class B retains the doctor route; no other shared-core guidance was altered.
  `git diff b5e5ced7~1 b5e5ced7 --numstat` on first-officer-shared-core.md = 1 insertion / 1 deletion (the abort clause only); test file is purely additive (170 lines, reuses existing helpers). Install lines match README:27,37 verbatim. AC-1 absence check proven real (not a spell-check): injecting a second `spacedock doctor` route into Class A makes the occurrence count exceed the single allowed prohibition form → FAILS. AC-2 proven: removing the Class B doctor route → FAILS. `go vet` clean.
- DONE: PASSED/REJECTED with evidence; confirm the tests prove the real intended behavior, not a spelling check.
  PASSED. Each AC test was exercised by mutation, not re-reading: AC-1 (route injected→fail), AC-2 (route removed→fail), AC-3 (mirror injected→fail). The banned-token check counts `spacedock doctor` in the isolated Class A span and permits only the exact prohibition string, so a routing instruction cannot hide as the prohibition — this is a structural invariant over real parsed spans, not a substring spell-check.

### Summary

PASSED. The deliverable is a surgical 1-line edit to `first-officer-shared-core.md` Startup step 1 (split abort by class) plus an additive contract-text test in `contract_gate_test.go`. Confirmed the failing-first claim by reverting the step-1 edit (AC-1/AC-2 test fails); the AC-3 test correctly remains green on that revert because it pins a single-source invariant, which I separately proved non-trivial by injecting a drift mirror. Each AC relationship was validated by mutating the real file and observing the test flip, not by substring re-reading — AC-1 presence (both install lines in Class A), AC-1 absence (doctor route count > prohibition fails), AC-2 (Class B doctor route deletion fails), AC-3 (mirror injection fails). Install lines match README verbatim; `go vet` clean; full integration package 26 passed; working tree clean.

## Feedback Cycles

**Cycle 1 (FO, 2026-06-02) — validation recommended PASSED; detached adversarial audit found two test-strength holes; routed to implementation for one-line tightenings before merge.**

The shipped prose is correct — the audit refuted nothing material on the contract edit (install lines verbatim vs README:27/37; scope surgical, only step 1's abort clause; AC-3 single-source confirmed by an independent grep). But two `skills/integration/contract_gate_test.go` assertions under-pin their own docstrings, so the suite would green-light a future regression in either direction:

- **M1 (`contract_gate_test.go:144`)** — the Class A no-doctor check is guarded by `if n := strings.Count(classA, doctor); n > 0`. If a future edit deletes the `Do NOT run \`spacedock doctor\`` prohibition entirely (zero doctor mentions), the check is skipped and the test passes — so it does NOT pin that Class A carries the no-doctor guidance. **Fix:** require the prohibition string present in Class A (so its removal fails), then keep the route-count guard for the present case.
- **M2 (`contract_gate_test.go:151`)** — AC-2 is a bare `strings.Contains(classB, "spacedock doctor")`, satisfied by a negated/disclaimer mention (audit verified an edit replacing the Class B route with `(Historically we suggested spacedock doctor but no longer.)` passes). **Fix:** assert the live-route phrasing (e.g. `run \`spacedock doctor\`` or `for the per-class remedy`) so a gutted route phrased as a disclaimer fails.

**Required proof:** each tightened test must FAIL on the auditor's adversarial edit (Class A prohibition deleted; Class B route replaced by a disclaimer) and pass on the real file. Mutation-verify, do not re-read.

## Stage Report: implementation (cycle 1)

- DONE: M1 — Class A no-doctor check is skip-on-zero; require the prohibition string present AND keep the route-count guard.
  `contract_gate_test.go` now asserts `Do NOT run \`spacedock doctor\`` is present in Class A, then asserts the doctor count is exactly 1 (the prohibition only). Mutation-verify: deleting the Class A prohibition sentence → test FAILS on both the presence and count assertions (was previously skipped at n==0). Commit 4268dc1a.
- DONE: M2 — Class B doctor assertion accepts a negated mention; assert the live-route phrasing instead.
  Replaced the bare `strings.Contains(classB, doctor)` with a check for both `run \`spacedock doctor\`` and `for the per-class remedy`. Mutation-verify: replacing the Class B route with `(Historically we suggested spacedock doctor but no longer.)` → test FAILS; real file passes. Commit 4268dc1a.

### Summary

Tightened the two flagged assertions in `skills/integration/contract_gate_test.go`; the shipped contract prose is unchanged. M1 now pins that Class A carries the no-doctor prohibition (presence) and routes nowhere (count==1), closing the skip-on-zero hole. M2 now requires Class B's live-route phrasing, so a disclaimer can no longer satisfy it. Both fixes mutation-verified by editing the real contract file: the Class A-prohibition-deletion adversarial edit fails M1, the Class B-disclaimer adversarial edit fails M2, and the unmodified file passes both — full `go test ./skills/integration/` green, `go vet` clean, contract file restored with no residue.

## Stage Report: validation (cycle 1)

- DONE: M1 — Class A now pins the no-doctor guidance: deleting the prohibition sentence FAILS (presence), adding a real doctor route into Class A still FAILS (count), the lone prohibition passes.
  Independently mutation-verified against commit 4268dc1a on the REAL contract file. Deleted the `Do NOT run \`spacedock doctor\`` sentence from Class A → FAILS at contract_gate_test.go:147 (missing prohibition) AND :150 (count 0≠1). Added `If unsure, run \`spacedock doctor\` to triage.` into Class A → FAILS at :150 (count 2≠1). Unmodified file: the lone prohibition passes (count==1). The old skip-on-zero hole (`if n > 0`) is closed: zero mentions now fails the presence assertion rather than being skipped.
- DONE: M2 — Class B now pins the LIVE route: replacing the route with a disclaimer FAILS; the real route passes.
  Replaced the Class B route with `ABORT startup with the actionable mismatch message. (Historically we suggested spacedock doctor but no longer.)` → FAILS at contract_gate_test.go:163 (missing live route `run \`spacedock doctor\`` + `for the per-class remedy`). The disclaimer still contains the bare substring "spacedock doctor", so the OLD `strings.Contains` assertion would have passed it — the strengthened phrasing-based check correctly rejects it. Unmodified file passes.
- DONE: Full AC suite holds: `go test ./skills/integration/` green, `go vet` clean, prose diff vs cycle 0 is zero (only the test file changed).
  Baseline and post-restore both = 26 passed; `go vet ./skills/integration/` clean. `git diff b5e5ced7 4268dc1a -- first-officer-shared-core.md --numstat` empty (zero prose diff); commit 4268dc1a touches only `contract_gate_test.go` (+22/−11). File restored after every mutation; `git status --short` clean at the end.

### Summary

PASSED (cycle 1). Both test-strength holes the detached audit found are genuinely closed, verified by mutating the real contract file rather than trusting the implementer. M1: deleting the Class A prohibition now fails (presence + count), and an added doctor route still fails (count≠1); the skip-on-zero gap is gone. M2: a disclaimer-only Class B mention now fails because the assertion requires the active routing phrasing (`run \`spacedock doctor\`` + `for the per-class remedy`), where the old bare-substring check would have green-lit it. The shipped contract prose is byte-identical to my cycle-0 PASS (zero diff); only `contract_gate_test.go` changed. Full integration package 26 passed, `go vet` clean, tree clean.
