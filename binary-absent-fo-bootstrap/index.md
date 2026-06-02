---
id: 5vb6mh9kewyh0p68r93mf6m1
title: Harden binary-absent FO startup — install hint in the contract-gate abort, don't route to the missing binary's doctor
status: implementation
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
