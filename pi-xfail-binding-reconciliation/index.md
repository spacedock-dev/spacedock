---
title: Reconcile the live Pi XFAIL bindings with the owners and evidence that exist today
status: ideation
score: 0.7
source: "pi-ux carve review, 2026-10-03: live Pi bindings name an archived owner, one binding cannot XPASS in code, and one real Pi FAIL carries no binding."
id: d525n1p5zgnz99hmtjq16z57
sprint: pi-ux
group: tooling
sprint-readiness: ready
gates:
    version: 1
    records:
        - id: gate:d525n1p5zgnz99hmtjq16z57:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:d525n1p5zgnz99hmtjq16z57-backlog-1
              briefing:
                id: briefing:d525n1p5zgnz99hmtjq16z57:backlog:attempt-1:revision-1
                digest: sha256:d4571aa8babbfc59dfb8670f73e8bb79bb6c5a7ee2029184c2b3763b2a4a3e40
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:d525n1p5zgnz99hmtjq16z57:backlog:1
                briefing: briefing:d525n1p5zgnz99hmtjq16z57:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-10-03T04:04:34.372071Z"
                decision: approve
                reason: Covers truthful live Pi XFAIL bindings, owners, and evidence, without touching the journey repairs.
                conn:
                    quote: i already said dispatch to ideation, but don't present the ideation gate until staff review finishes
                    source: captain instruction, this session, 2026-10-03
              application:
                target-stage: ideation
                state: consumed
started: 2026-10-03T04:05:02Z
---

The live Pi XFAIL inventory contains four bindings, including one naming an archived
PASSED owner whose validation explicitly leaves a separate topology failure unresolved.
This task restores truthful ownership without claiming that representation repairs fix
Pi conduct or that a green Go test under XFAIL proves the journey passed.

## Problem

Observed at code `cdfa462d1d426febb2391733507c7fde2143fb0a` on 2026-10-03:

- `shared_live_runner_test.go` binds Pi owner-handoff, rejection-flow, keep-moving,
  and smallest-sufficient-mechanism; the registry reconciliation derives `xfail/pi=4`.
- `rejection-flow` names archived PASSED `p17swb3375rt525fn7f8xt7e`. Its timeout
  repair passed, but `_archive/finish-pi-rejection-flow.md`, validation Summary,
  explicitly retains `rejection-worker-topology` as a separate XFAIL. Its reported
  `--- PASS: TestLiveCommonRejectionFlow (1499.96s)` is NOT normal PASS evidence.
- `piSharedLiveDriver.smallestMechanismTrace` calls `claudeMechanismTrace` on Pi
  session data. This is a harness/extractor mismatch that prevents trustworthy
  XPASS assessment, not evidence that the binding is obsolete or the product faulty.
- `recorded-gate-lifecycle` is unbound on Pi. Its active repair owner records a real
  missing-reference FAIL and expressly prohibits adding a masking Pi XFAIL.

## Proposed approach

All restatements below, including replacement AC-2 and the revised surface estimate,
are proposed for the captain's ideation gate; they are not implemented registry changes.
The FO approved this binding-only direction during ideation.

1. Re-anchor only Pi `rejection-flow` from archived `p17swb3375rt525fn7f8xt7e`
   to active `6h3teccccn3qh71yqcmjbjx4` (`own-pi-rejection-worker-topology`). Update
   the adjacent registry comment to describe the remaining topology failure rather
   than imply the completed timeout repair owns it. The FO filed this bounded owner
   after the audit found no active Pi topology owner; the existing
   `persist-codex-rollout-for-rejection-topology` is explicitly Codex-only.
2. Keep the other three Pi bindings and the recorded-gate absence as listed under
   AC-2. Keep-moving's historical XPASS is a removal candidate, not permission to
   delete: exact XPASS plus an unbound normal PASS are required, and no such paired
   evidence is currently retained in its owner. Any later evidence-led clearance
   must update the disposition before validation; no speculative removal is planned.
3. Reuse the existing registry-reading active-owner check, not a new owner manifest,
   parser, allowlist, or test framework. It already produces the exact baseline RED.
   The value served is AC-1; a handwritten list or source grep is simpler but cannot
   demonstrate that actual registered owner IDs resolve against workflow state.

## Scope boundary

Binding attribution, evidence classification, and this disposition record only.
Extractor and product repairs remain with their journey owners in `pi-live-completeness`
(and the newly filed Pi topology owner). Do not edit the Pi driver, shared assertions,
fixtures, XFAIL grading policy, timeout budgets, runtime skills, or Claude/Codex bindings.
Archived PASSED is not proof of XPASS, and a Go PASS under XFAIL is not normal PASS.
No new runtime, quarantine, CI lane, or binding is proposed. No user-visible command,
stored format, runtime authority, or host behavior changes; no docs-site diff is needed.

## Expected surface and tolerance

Proposed replacement baseline: Estimate net LOC change: +0, across 1 file
(`internal/ensigncycle/shared_live_runner_test.go`). Insertions ~+5, deletions ~-5;
tolerance +/-10 net LOC, +/-0 code files. This entity body/report is the separate
state deliverable. Existing tests suffice; no test scaffolding or Pi runner edit.

Declared observable semantic change: the Pi rejection-flow XFAIL owner ID and its
explanatory comment. Other binding states remain unchanged under today's evidence.
If new evidence justifies clearance, present that changed disposition before proceeding.

## Acceptance criteria

The following criteria are proposed restatements for the ideation gate. AC-2 replaces
the seed's extractor-repair criterion and subsumes its duplicate disposition criterion.

**AC-1 (VALUE) — Every live Pi binding names a non-archived, nonterminal owner or is absent with exact passing evidence.**
Verified by the existing `TestRuntimeLiveTODOOwnersAreActive` in
`internal/contractlint/live_registry_reconciliation_test.go`, with the explicit state
checkout supplied through `SPACEDOCK_LIVE_STATE_DIR`. It reads the live declarations
and joins their owner IDs to actual entity state (despite TODO in its name, it checks
XFAIL too). Independent baseline: the 2026-10-03 run reports exactly one inactive Pi
binding, rejection-flow → archived PASSED p17. Finished result: zero inactive Pi
bindings. Falsifying edit: restore p17; the same test must turn RED. Binding absence
alone does not authorize removal; AC-3 supplies the independent clearance bar.

**AC-2 — Every live Pi binding has a disposition and deciding evidence; extractor limitations remain classified, not silently repaired or declared obsolete.**
Verified by a one-off independent inventory reconciliation against the AST-derived
live registry and these rows, including owner bodies. Falsifier: delete a disposition
row or introduce an unrecorded Pi binding; reconciliation is incomplete. These are
proposed implementation dispositions, not claims of edits already shipped:

| Pi journey | Disposition | Owner and deciding evidence |
| --- | --- | --- |
| `owned-conflict-owner-handoff` | **kept** | `fe7bfjz9sb8wyckmnnm3ncjx`, active `repair-pi-owner-handoff.md`: problem records retained XFAIL; AC-1/2 require exact XPASS then normal PASS. No clearing pair supplied. |
| `rejection-flow` | **re-anchored** | p17 → `6h3teccccn3qh71yqcmjbjx4`, active `own-pi-rejection-worker-topology/index.md`. Archived p17 validation Summary explicitly defers `rejection-worker-topology` while passing only the timeout repair. No active Pi topology owner existed before the FO filing; the Codex-only topology task is not a substitute. New owner's AC-2 retains the binding until exact XPASS plus normal PASS. |
| `keep-moving-posture` | **kept** | `x02375wsg6q61xek7p0t36j2`, active `repair-pi-keep-moving-posture.md`: source records CI run `31770740214`, Pi `openai/gpt-5.6-luna:max`, XPASS `observed=[]`, but no exact normal PASS. AC-2 still requires both. Treat as a clearance candidate, not an evidenced clearance. |
| `smallest-sufficient-mechanism` | **kept** | `h30c9jrfcf21fdh2qs5z58sd`, active `repair-pi-smallest-sufficient-mechanism.md`: no passing pair. Pi driver's `smallestMechanismTrace` calls `claudeMechanismTrace`; classify as extractor/harness mismatch requiring owner-led repair and subsequent product assessment. Do not infer obsolescence or a product diagnosis from this mismatch. |

No binding is cleared on the evidence available today. Additional unbound disposition:
`recorded-gate-lifecycle` **remains absent** (not a cleared binding), owned by
`gcmfwfjd9735b58sbzw7xsb8`, `repair-pi-recorded-gate-lifecycle.md`. Its evidence names
CI `31770740214` and `observed=[recorded-gate-lifecycle-violation]`, missing committed
reference. Its AC-2 forbids adding Pi XFAIL; ordinary FAIL remains visible. Its old
Claude-Opus inventory prose is historical, not the authority for today's source map.

**AC-3 — No Pi binding is removed on archived status, timeout-only success, or masked Go PASS.**
Verified by the registry diff and, for any proposed removal, retained exact-target Pi
XPASS with empty observed semantic failures followed by normal PASS with that binding
absent on the candidate. Capture code SHA, target/model, command, exit code, grade and
durable artifacts; a lane-level green alone does not qualify. Falsifier: remove
rejection-flow citing only p17's PASSED verdict or XFAIL-backed Go PASS; reject clearance.

**AC-4 (no-regression) — Journey behavior and other runtimes are unchanged.**
Verified by `go test ./internal/ensigncycle/...`, the existing registry reconciliation,
full/race checks, and a scoped diff showing no assertion, extractor, fixture, grading
policy, Claude/Codex binding, or recorded-gate Pi absence changes. Falsifier: alter a
shared assert or add a masking recorded-gate Pi XFAIL even if the lane becomes green.

## Risk evidence

No spike needed: the registry AST reader, explicit mutable-state join, and grade
classification already exist and were exercised during ideation. With
`SPACEDOCK_LIVE_STATE_DIR` set to the dispatched checkout, registry reconciliation
passed and derived four Pi XFAILs; the owner check failed only for p17. This is the
red-first acceptance baseline, not a new mechanism to invent. `TestGradeLiveTargetMatrix`
and `TestLiveGradeLaneResult` passed: XFAIL, XPASS and normal PASS are distinct grades,
even though all can return a green Go test. Returning lane success for XFAIL therefore
cannot establish clearance. No model-backed run or extractor repair was attempted.

## Test plan

Primary registry-reading proof owner: existing
`internal/contractlint/live_registry_reconciliation_test.go`, specifically
`TestRuntimeLiveRegistryReconciliation` and `TestRuntimeLiveTODOOwnersAreActive`.
Do not add a duplicate test beside the shared runner. Deterministic cost is seconds;
mutation control is a temporary reversal of only the owner ID, then restoration.

- Before editing, run `SPACEDOCK_LIVE_STATE_DIR="$PWD/docs/dev/.spacedock-state" go test
  ./internal/contractlint -run '^(TestRuntimeLiveRegistryReconciliation|TestRuntimeLiveTODOOwnersAreActive)$'
  -count=1 -v`. Already observed RED only for p17. After the one-line re-anchor,
  expect GREEN; restore p17 temporarily and expect the same RED. Record the state
  branch SHA as well as the code SHA. A missing env var SKIPs the owner join and
  is not acceptance; a fresh clone must initialize/fetch the declared dev state checkout.
- One-off manual proof for AC-2/3/4: independently compare every derived Pi binding
  with the disposition table and read each cited owner/evidence record. An omitted
  row falsifies AC-2; missing XPASS/normal-PASS pair falsifies clearance under AC-3;
  unrelated code changes falsify AC-4. No static prose-presence test is proposed.
- Run focused grade controls, `go test ./internal/ensigncycle/...`, `go test ./...`,
  `go test ./... -race`, and formatting validation during implementation. Preserve
  failures honestly; no unrelated repairs to make this registry-only task green.
- One authorized live Pi cadence confirms emitted owner/grade attribution after the
  re-anchor, not product repair. Retain artifacts and distinguish XFAIL/XPASS/PASS,
  infrastructure failure, and timeout. Budget tens of minutes (p17's historical run
  was ~25 minutes); no run-budget changes or retry controller are proposed here.
  Extractor mismatch stays owned by h30 and recorded-gate may still FAIL normally.
  If clearance is proposed later, use the exact target sequence and paired evidence
  above rather than aggregate cadence success. Historical XPASS alone is insufficient.

## Stage Report: ideation

- DONE: Task body requires, for every live Pi XFAIL binding, a recorded disposition of cleared, kept, or re-anchored, with the evidence that decided it, under the exact `## Acceptance criteria` heading.
  Four proposed dispositions appear under AC-2; registry reconciliation independently derives `xfail/pi=4` at code cdfa462d1d426febb2391733507c7fde2143fb0a.
- DONE: A value AC measures that every live Pi binding names a non-archived owner or is absent, with today's binding that names an archived PASSED owner as the independent baseline.
  AC-1 baseline exercised: existing active-owner test fails only on archived p17; FO filed active successor 6h3teccccn3qh71yqcmjbjx4 (state HEAD 5dc6d9a3d7475c95f084807cb63f529ee5de3165 at report preparation).
- DONE: The test plan names the registry-reading proof owner; removal requires exact XPASS or normal PASS evidence, extractor and product repairs stay with their journey owners, and the smallest-mechanism extractor mismatch is classified rather than declared obsolete.
  Existing internal/contractlint/live_registry_reconciliation_test.go owns proof; removal actually requires BOTH exact XPASS and normal PASS per owner contracts. Pi runner stays untouched.
- DONE: AC-1 proof plan and baseline evidence.
  TestRuntimeLiveTODOOwnersAreActive reads registered owners and workflow state; observed RED for p17, planned GREEN for 6h3 then restored-p17 RED mutation control. No registry edit in ideation.
- DONE: AC-2 proposed replacement proof plan and per-binding evidence.
  Independent inventory-to-record reconciliation must cover all four Pi bindings; deleting any row falsifies completeness. Replaces extractor repair and folds duplicate seed disposition AC, pending captain approval.
- DONE: AC-3 proposed clearance proof plan.
  Registry diff plus exact XPASS/normal-PASS pair for any removal; p17's archived PASSED and XFAIL-backed Go PASS cannot clear a binding. No removal proposed.
- DONE: AC-4 proposed no-regression proof plan.
  Focused/full/race checks and scoped diff must preserve all shared asserts, extractors, other-runtime bindings and recorded-gate Pi absence; adding a masking binding fails this criterion.
- DONE: Proposed owned-conflict-owner-handoff disposition — kept.
  fe7bfjz9sb8wyckmnnm3ncjx is active; owner records retained XFAIL, no exact clearing pair.
- DONE: Proposed rejection-flow disposition — re-anchored, not cleared.
  Archived p17 validation explicitly defers rejection-worker-topology; FO-created 6h3teccccn3qh71yqcmjbjx4 owns it. Prior active topology task is Codex-only.
- DONE: Proposed keep-moving-posture disposition — kept.
  x02375wsg6q61xek7p0t36j2 records run 31770740214 XPASS observed=[] but no paired normal PASS; historical XPASS is insufficient for deletion.
- DONE: Proposed smallest-sufficient-mechanism disposition — kept.
  h30c9jrfcf21fdh2qs5z58sd owns repair; claudeMechanismTrace on Pi data is a harness mismatch, not proof of obsolescence or product failure.
- DONE: Proposed recorded-gate-lifecycle disposition — remains unbound.
  gcmfwfjd9735b58sbzw7xsb8 records ordinary missing-reference FAIL and forbids masking Pi XFAIL; absence is intentional, not a cleared binding.
- DONE: Exercise established proof mechanisms; no spike needed.
  Registry derives four Pi XFAILs; grade matrix and lane-result tests distinguish XFAIL/XPASS/PASS, and fail if semantic XFAIL is mislabeled normal PASS. gofmt -w ./cmd ./internal completed; its unrelated two-line formatting delta was reverted, leaving no code changes.
- FAILED: Broad full and race verification attempts.
  go test ./... and go test ./... -race each exceeded this tool's 120-second cap; neither is claimed green. Focused registry and grade checks provide the ideation evidence; full suites remain implementation obligations.
- SKIPPED: Live Pi cadence and production registry edits.
  This dispatch is ideation; no model spend or code edits. The proposed restatements and +0/one-code-file estimate await captain approval.

### Summary

Proposed a one-binding re-anchor using the existing red-first registry proof, with four evidence-backed dispositions and recorded-gate intentionally unbound. The proposed AC-2 replacement keeps extractor/product repair with journey owners; all restatements remain subject to the captain's ideation gate.

AC scan output (all four criteria cited):
```json
{"command":"read","stage":"ideation","acs":[{"id":"AC-1","line":"107","unevidenced":"false","citations":[{"line":"195","text":"  AC-1 baseline exercised: existing active-owner test fails only on archived p17; FO filed active successor 6h3teccccn3qh71yqcmjbjx4 (state HEAD 5dc6d9a3d7475c95f084807cb63f529ee5de3165 at report preparation)."},{"line":"198","text":"- DONE: AC-1 proof plan and baseline evidence."}]},{"id":"AC-2","line":"117","unevidenced":"false","citations":[{"line":"193","text":"  Four proposed dispositions appear under AC-2; registry reconciliation independently derives `xfail/pi=4` at code cdfa462d1d426febb2391733507c7fde2143fb0a."},{"line":"200","text":"- DONE: AC-2 proposed replacement proof plan and per-binding evidence."}]},{"id":"AC-3","line":"137","unevidenced":"false","citations":[{"line":"202","text":"- DONE: AC-3 proposed clearance proof plan."}]},{"id":"AC-4","line":"144","unevidenced":"false","citations":[{"line":"204","text":"- DONE: AC-4 proposed no-regression proof plan."}]}]}
```
