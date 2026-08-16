---
title: Auto-continue grades a human-gate bypass as green
status: implementation
source: "Live-harness audit finding 3 (2026-08-16); captain order: file and fast-track on the stack, local focused test first"
id: 7xe7hxt1qce1x9b3dm0k6ymg
gates:
    version: 1
    records:
        - id: gate:7xe7hxt1qce1x9b3dm0k6ymg:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:7xe7hxt1qce1x9b3dm0k6ymg-ideation-1
              briefing:
                id: briefing:7xe7hxt1qce1x9b3dm0k6ymg:ideation:attempt-1:revision-1
                digest: sha256:67b5b390d6e86a6ef0e86959b8a72b38207d0cd9a904018acadafe492a7caa36
                request-digest: sha256:abaa582a10c85b4690c7ff1001a39e2d9bcd2f20263d6410d6186e8703a45e0e
                room-ref: ./red-auto-continue-gate-bypass/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:7xe7hxt1qce1x9b3dm0k6ymg:ideation:1
                briefing: briefing:7xe7hxt1qce1x9b3dm0k6ymg:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T17:49:21.399649Z"
                decision: approve
              application:
                target-stage: implementation
                state: consumed
worktree: .worktrees/spacedock-ensign-red-auto-continue-gate-bypass
---

## Problem

The auto-continue journey's fixture pins validation as a human gate, and its own comment
says the correct FO advances to validation, dispatches a fresh validator, and presents the
gate. The graded regex accepts `done` anyway, so an FO that resolves a human gate nobody
approved grades GREEN. The check that would catch it is wired through an optional interface
only the codex driver implements, so on claude and pi a gate bypass is invisible.

Every citation re-verified from source during this ideation:

| Claim | Source | Verified |
| --- | --- | --- |
| Fixture pins validation as a human gate | `auto_continue_fixtures_test.go:81-85` (`gate: true`, also `fresh: true`, `feedback-to: implementation`) | yes |
| Fixture's comment names the correct end state | `auto_continue_fixtures_test.go:101-103` — "advances to validation + dispatches a fresh validator (then presents the validation gate)" | yes |
| Graded regex accepts `done` | `auto_continue_fixtures_test.go:24` — `(?im)^status:\s*(validation\|done)\s*$` | yes |
| Gate-open check exists | `assertAutoContinueDispatchEvidence`, `claude_live_runner_test.go:280-311`; gate-open at 303-309 | yes |
| ...but is wired through an optional interface | `shared_promoted_live_test.go:153-157` — a `driver.(interface{...})` type assertion | yes |
| ...implemented only by codex | `codex_live_runner_test.go:56-58`; repo-wide grep finds no other implementor | yes |

Two defects the audit did not name, found while verifying:

1. `assertAutoContinueDispatchEvidence` takes a `workflowRoot` parameter its body never
   references (`claude_live_runner_test.go:280`) — dead.
2. The runbook (`autoContinuePrompt`, `auto_continue_fixtures_test.go:133-140`) grants no
   conn and no auto-approve coaching. So `done` is not merely "beyond validation"; on this
   fixture it is only reachable by the FO self-approving. There is no legitimate path to it.

## Replay of retained artifacts: did any past green contain a bypass?

**No.** Ten real legs across three green CI runs, zero bypasses.

Where the corpus actually is — the two sources the assignment named are both empty:

- `.live-artifacts-*` at the repo root hold wait-matrix scenarios only; no auto-continue.
- `SPACEDOCK_LIVE_ARTIFACT_DIR` retention under `/private/tmp` has 12 auto-continue
  directories, but the OS tmp cleaner purged every file inside them (0 files across all 12;
  sibling scenario dirs are equally empty, dir mtimes bumped 2026-08-14 09:56). Nothing to
  replay locally.

The recoverable corpus is unexpired GitHub Actions artifacts. Replayed runs 31915540750
(main, green), 31899041131 (green), 31903299621 (green) — claude-sonnet-5, claude-opus-4-8,
and codex, both fixture variants each:

| Signal | Result across all 10 legs |
| --- | --- |
| `status=done` set | 0 |
| `status=validation` set | 10 |
| `gate prepare` invoked | 10 |

So tightening the regex flips no observed green. The hole is latent, not currently
exercised — which is why it survived: nothing green had to change for it to be wrong.

## Spike

Riskiest unverified mechanism: whether the gate-open/dispatch-evidence check can be
host-neutral across three drivers. Exercised, not assumed:

- **The gate-open and committed-report halves are host-neutral by construction.** They read
  durable state only — `gates.Read(entityPath)` plus `CurrentSummary(doc,"validation").State`
  (`claude_live_runner_test.go:303-309`) and `git log -S'## Stage Report: validation'`
  (300-302). No stream, no dialect. `attemptState` (`internal/gates/model.go:368`) returns
  `open` only for an unwithdrawn, unresolved attempt, so a resolved gate cannot read `open`.
- **The worker-lifecycle half parses stream dialects**, so it was exercised per runtime shape.
  `assertWorkerLifecycle(stream, report, "validation", "gate prepare")` run against real
  captured streams:

  | Runtime shape | Streams | Verdict |
  | --- | --- | --- |
  | claude-sonnet-5 | single-root + split-root (run 31915540750, 31899041131) | PASS |
  | claude-opus-4-8 | single-root + split-root (run 31903299621) | PASS |
  | codex | already runs in production; green today | PASS in CI |
  | pi | none retained anywhere | **unexercised** |

  This is the first time the claude branch of `assertWorkerLifecycle` has been run on real
  auto-continue data. It works — making the check host-neutral for claude is proven, not
  speculative. The codex *public* stream reds on replay because the codex leg feeds
  `codexNativeLifecycleStream` (public stream + the CODEX_HOME rollout) and CI does not
  upload CODEX_HOME: an artifact-retention gap, not a defect.

- **Both grading directions falsified offline, zero model spend.** Against a bypass-shaped
  end state (`status: done` plus a validation report): current grading GREEN, tightened
  grading RED. Against a conforming end state (`status: validation` plus the report): both
  GREEN.

Provenance (sha256, for the replay fixture): claude-sonnet-5 single-root
`465e849afda41a39c8443bae98119e8ab806a0af0a116f54f5321761d7eab287` (430087 bytes);
split-root `12ccfdd71964eede7a358ae4f5a4409058a15f379c067f78d70c6c8ed003510e`.

## Proposed approach

1. **Narrow the accepted end state to the one the fixture pins.** Split
   `validationStatusOrBeyond` into `validationStatusAC5` (`^status:\s*validation\s*$`) and
   `terminalStatusAC5` (`^status:\s*done\s*$`). `assertAutoContinue` reds a `done` end state
   with its own code, `human-gate-bypassed`, distinct from the stall it reds today. The
   existing `*gradedErr` plumbing carries it: `durableSemantic` (`claude_live_runner_test.go:210-218`)
   preserves a `*gradedErr`'s own code rather than overwriting it, so no new machinery.
2. **Make the evidence check unconditional and compiler-enforced.** Delete the
   `verifyAutoContinueDispatch` optional interface. Add `lifecycleStream(t, result) string`
   to the `liveDriver` interface (`claude_live_runner_test.go:85-94`): claude and pi return
   `result.stream`, codex returns `nativeLifecycleStream(t, d, result)`. `runAutoContinueJourney`
   then calls `assertAutoContinueDispatchEvidence` for every driver with no type assertion.
   Drop its dead `workflowRoot` parameter while there.
3. **Commit one real captured claude stream as an offline replay control**, following
   `zero_discover_replay_test.go` (the audit's own cited example of the good pattern), so the
   claude dialect's parse is regression-guarded without a live run.
4. **Update the registry entry** so its required outcome names the open gate.

### Why each new mechanism, and what simpler thing was rejected

- `lifecycleStream` on the `liveDriver` interface — serves AC-1/AC-3. Simplest alternative:
  keep the optional type-assertion interface and just add claude and pi implementations.
  Insufficient: an optional interface is satisfiable by omission, which is the precise defect
  finding 3 names. A driver that forgets it grades green with no check and nothing fails. An
  interface method makes omission a build error.
- Committing a real captured stream as testdata — serves AC-2. Simplest alternative: a
  hand-written synthetic stream. Insufficient: a synthetic stream proves the parser matches
  what I imagined the runtime emits, not what it emits; finding 11 is exactly that class of
  blindness. Second alternative: no replay control, rely on the live run. Insufficient: it
  prices every future regression at a live run.
- A distinct `human-gate-bypassed` grade code — serves AC-1. Simplest alternative: reuse the
  generic `auto-continue-state`. Insufficient: the audit's complaint is that the failure was
  invisible; a generic code leaves a bypass indistinguishable from a stall in the journey
  metrics, so the next reader has to redo the audit to tell them apart.

## Acceptance criteria

**AC-1 (value)** On every runtime the journey grades, a human-gate bypass is a RED. The
count of hosts on which a bypass grades green is 0 of 3, down from 2 of 3 today (claude and
pi, via both the `done` regex and the absent gate-open check).
Verified by: an offline table enumerating {claude, codex, pi} x {`status: done` end state,
validation gate resolved rather than open} and asserting every cell RED. The change that
would make it fail: re-adding `done` to the accepted regex, or re-conditioning the evidence
check on a host. Plus compile-time proof no driver can opt out — deleting a driver's
`lifecycleStream` fails the build.

**AC-2 (value, non-regression)** The tightening produces no false red on real observed
conforming behavior: every real captured auto-continue leg from a green CI run still grades
GREEN.
Verified by: the committed replay test over the real captured claude stream from run
31915540750, plus the ten-leg replay recorded above. The change that would make it fail:
tightening the status regex past `validation`, or an ordering rule the real stream violates.

**AC-3 (means, serves AC-1)** The dispatch-evidence check runs on claude, codex, and pi with
no optional interface and no host conditional; `verifyAutoContinueDispatch` does not exist.
Verified by: the build (the `liveDriver` method), and AC-1's table covering all three.

**AC-4 (value)** The live journey is green on claude and codex under the tightened grading.
Verified by: the live loop under the CI shim replica, >=3 consecutive green runs per runtime
with a ledger and stream digest per run; budget 8 runs per runtime. Loop criterion: 3/3
green, no leg graded on a code introduced by this change.

## Test plan

| Test | Proves | Cost |
| --- | --- | --- |
| Offline bypass table, 3 drivers x 2 bypass shapes (default tag) | AC-1, AC-3 | ~1 h to write; runs in ms, no model spend |
| Offline replay of the real captured claude stream through the lifecycle check | AC-2 | ~40 min; runs in ms |
| Existing `auto_continue_negative_test.go` positive baseline, extended with the conforming end state | AC-2 | ~15 min |
| `go test ./internal/ensigncycle/ -run AutoContinue -count=20` and `-race` | the focused local loop is reliably green, not flaky | ~2 min |
| `go test ./...` and `-race` | no collateral damage; all changes are `_test.go` plus one doc | ~5 min |
| `contractlint` registry reconciliation | the registry entry still binds the journey and fixture ids | seconds |
| Live loop, claude-sonnet, CI shim, ledger + digest per run | AC-4 | ~6 min/run, >=3 runs, budget 8 |
| Live loop, codex, CI shim, ledger + digest per run | AC-4 | ~5 min/run, >=3 runs, budget 8 |
| Pi | not run on PR (pi-live is cadence-only). See Risks — no preemptive xfail. | none |

Live budget: ~40 min wallclock with both runtimes in parallel, ~1.5 h if both hit the cap.
All offline directions are falsified before any CI spend, per the captain's order.

## Expected surface

| File | Change | Insertions |
| --- | --- | --- |
| `internal/ensigncycle/auto_continue_fixtures_test.go` | regex split, coded bypass red | ~25 |
| `internal/ensigncycle/auto_continue_negative_test.go` | bypass negative + per-host table | ~45 |
| `internal/ensigncycle/auto_continue_replay_test.go` (new) | real-stream replay control | ~55 |
| `internal/ensigncycle/testdata/…auto_continue…stream.jsonl` (new) | captured claude stream | 279 lines / 430 KB |
| `internal/ensigncycle/shared_promoted_live_test.go` | unconditional wiring | ~10 |
| `internal/ensigncycle/claude_live_runner_test.go` | interface method, split assert, drop dead param | ~20 |
| `internal/ensigncycle/codex_live_runner_test.go` | `verifyAutoContinueDispatch` → `lifecycleStream` | ~5 |
| `internal/ensigncycle/pi_shared_live_runner_test.go` | `lifecycleStream` | ~2 |
| `internal/ensigncycle/shared_live_runner_test.go` | retire the finding-3 AUDIT note | ~5 |
| `docs/runtime-live-ci-registry.md` | required outcome + evidence | ~8 |

10 files, ~175 insertions excluding the captured stream. Tolerance: +/-3 files, +/-60
insertions on code and docs, +/-1 fixture file. The captured stream is one file and is
counted separately; if a reviewer objects to 430 KB, the fallback is a line-type filter
(keep only `tool_use` and `task_notification` rows) with a capture-time proof that the
filtered and full streams produce an identical verdict.

## Observable semantics changed

- **Grading vocabulary.** The auto-continue journey's accepted end state narrows from
  {`validation`, `done`} to {`validation`}. A new grade code `human-gate-bypassed` can appear
  in live grade output and journey metrics.
- **Accepted end states.** `status: done` on this journey becomes a RED. `resolveAutoContinueEndState`'s
  `_archive` lookup stays — a terminalized run gets its body read so it reds honestly on the
  bypass code rather than on an empty after-state.
- **Harness driver contract.** `liveDriver` gains a required `lifecycleStream` method. Test
  surface only.
- **No product change.** No binary, skill, or FO-instruction behavior is touched. Only what
  the harness certifies changes.

## Documentation diff

`docs/runtime-live-ci-registry.md`, entry `auto-continue-after-implementation` (lines 197-207):

```diff
-- **Required outcome:** After observing a completed implementation report, the
-  first officer advances to validation and dispatches a fresh validator instead
-  of stopping.
+- **Required outcome:** After observing a completed implementation report, the
+  first officer advances to validation, dispatches a fresh validator, and
+  presents the validation gate, leaving it open. The fixture's validation stage
+  is `gate: true` and the runbook grants no conn, so a run that reaches `done`
+  resolved a human gate nobody approved: that is a failure, not a success.
 - **Fixtures:**
   - `auto-continue/single-root` — a single-root workflow parked at completed
     implementation with worktree-backed validation.
   - `auto-continue/split-root` — the same invariant in a separate state checkout
     with non-worktree validation.
-- **Evidence:** Every fixture variant must pass the shared `assertAutoContinue`
-  durable-state assertion.
+- **Evidence:** Every fixture variant, on every runtime, must pass the shared
+  `assertAutoContinue` durable-state assertion — end state `status: validation`
+  only, with `done` red under code `human-gate-bypassed` — and the unconditional
+  dispatch-evidence check: an open validation gate, a committed
+  `## Stage Report: validation`, and one fresh validator dispatched and completed
+  before `gate prepare`.
```

## Risks

- **Pi's stream dialect is unexercised for this journey.** No pi auto-continue stream is
  retained anywhere; pi-live is cadence-only and the committed pi testdata is cost-only with
  no tool calls. If pi's real stream lacks the shapes `assertWorkerLifecycle`'s pi branch
  reads (`claude_runtime_helpers_test.go:178-192`), the next pi cadence run reds on
  `validation-worker-lifecycle`. Disposition: **do not preemptively xfail.** A gap declared
  without observing it would mask a genuine pi pass. If pi reds, the red is information —
  finding 11's blindness applies here too — and the repair is a pi-dialect extractor, not a
  loosening of this grading. Recorded so the next reader does not mistake it for a product
  regression.
- **The captured-stream fixture ages.** It pins the claude tool-use dialect. If the dialect
  changes, the replay control reds without any FO misbehavior. Accepted: that red is a true
  signal that the live grader has gone blind, which is the failure this entity exists to
  prevent.

## Stack position

This entity rides the current PR stack (#720). Ideation runs without a worktree. Implementation
branches from the stack tip current at dispatch time — the base is deliberately not pinned
here, since the team-mode rejection layer in ideation may land above
`spacedock-ensign/repair-codex-rejection-round-recording` (571017df3) first.

## Out of scope

- Other audit findings; rejection-flow surfaces (owned by run-rejection-journey-in-team-mode).
- Product binary changes.
- Hardening `assertWorkerLifecycle`'s ordering semantics (audit finding 10) — a separate
  defect on a shared helper; this entity consumes the helper as-is.

## Stage Report: ideation

- DONE: Fast-track: this entity rides the current PR stack (#720 ...) — note that in the plan, do not pick the base now.
  "Stack position" section records #720 and defers the base to dispatch time; ideation ran without a worktree.
- DONE: Read the entity, audit finding 3 ..., the AUDIT note ..., and verify every citation from source yourself (regex, fixture gate pin, codex-only interface wiring).
  All six citations re-verified and tabulated in "Problem"; the fixture gate pin is at `auto_continue_fixtures_test.go:81-85` (audit said 84-86 — off by three, substantively correct).
- DONE: Derive the fixture-determined end state per runtime ... Design grading to exactly that shape ... No multiple-path acceptance anywhere the fixture pins one path.
  End state derived as `status: validation` + open validation gate + dispatched-validator evidence, cited to fixture and driver sources; `done` reds under code `human-gate-bypassed`, and the optional interface is replaced by a `liveDriver` method so no host can opt out.
- DONE: Replay the retained run artifacts ... did any past green contain a bypass?
  No. Ten legs across three green CI runs: zero `status=done`, ten `status=validation`, ten `gate prepare`. The two named local sources were both empty (repo-root dirs are wait-matrix only; the 12 `/private/tmp` dirs were purged by the OS tmp cleaner), so the corpus came from unexpired CI artifacts.
- DONE: Value AC per the captain's order ... Both directions falsified before any CI spend.
  AC-1 measures hosts-where-a-bypass-greens (2 of 3 today, 0 of 3 after); AC-4 sets the loop criterion at 3/3 consecutive green per runtime under the CI shim, budget 8. Falsified offline: the bypass end state grades GREEN today and RED under the tightened regex; the conforming end state grades GREEN under both.
- DONE: Spike or no-spike: the host-neutral gate-open check across three drivers is the riskiest mechanism ...
  Split by risk and recorded in "Spike": the gate-open and committed-report halves read durable state only (`gates.Read` + `git log -S`), so host-neutrality is proven by construction; the stream-parsing half was exercised per runtime shape — four real claude streams across two models PASS, codex passes in production, pi is unexercised and filed as a named risk.
- DONE: Declare expected surface (files, insertions, tolerance) and observable semantics changed ... Test plan per AC with costs.
  10 files / ~175 insertions plus one captured stream, tolerance +/-3 files and +/-60 insertions; semantics limited to grading vocabulary, accepted end states, and the test-only driver contract. Nine-row test plan with per-AC mapping and costs.
- DONE: Write the ideation stage report in the entity and stop for the gate. Do not start implementation.
  This report; no implementation started, no worktree, no product file touched.

### Summary

Every audit citation held up, and verifying them surfaced two things the audit missed: a dead
`workflowRoot` parameter, and the fact that the runbook grants no conn — so `done` is not
"beyond validation" on this fixture, it is only reachable by self-approval. The replay
answered the decision-relevant question cheaply: no past green contained a bypass, so
dropping `done` flips nothing and the hole is latent rather than load-bearing. The spike
changed the design: the gate-open half needs no dialect work at all (it reads durable state),
and the stream half already parses real claude auto-continue streams correctly despite never
having been wired for claude — so host-neutrality costs almost nothing and is proven rather
than hoped. The one honest gap is pi, which has no retained stream anywhere; I recommend
shipping it unconditional without a preemptive xfail, since a gap declared without observing
it would mask a genuine pi pass.

## Stage Report: implementation

- DONE: Execute the gated design exactly ... Nothing absorbed beyond the gated surface (10 files, ~175 insertions + the captured stream, tolerance +/-3 files / +/-60 insertions).
  Commit 5cdb766ec. All four gated elements landed; 10 files (within tolerance) but ~433 code+doc insertions vs ~175 +/-60 — over tolerance, accounted in Deviations below.
- DONE: Branch from the current stack tip: spacedock-ensign/repair-codex-rejection-round-recording at 571017df3 ... Push via the SSH remote.
  Branch spacedock-ensign/red-auto-continue-gate-bypass branched at 571017df3; push held with the PR (see the SKIPPED item).
- DONE: Offline proof before any live spend, per the gated ACs ...
  AC-1 table: {claude,codex,pi} x {done end-state, resolved-gate} all RED under `human-gate-bypassed`, each host with a conforming open-gate control so the table cannot pass on an assertion that reds unconditionally. AC-2: the real captured stream and both fixture layouts GREEN. `-run AutoContinue -count=20`, `-race`, full `go test ./...`, `-race -timeout 30m`, and contractlint all clean. Falsified three ways: re-adding `done` reds all three hosts' terminal cells; reverting the durableSemantic unwrap reds the code-survival test; disabling the gate-open check reds all three resolved-gate cells.
- FAILED: Live loop per AC-4: ... >=3 consecutive greens per runtime (claude-sonnet, codex) within budget 8 per runtime ...
  Codex reached 3/3 consecutive GREEN, but on the pre-fix candidate, so it must be re-run once the candidate changes. Claude did not close: run 2 hit finding 1 below. Ledgers with per-run stream digests at `/Users/clkao/.claude/jobs/4e49247e/tmp/live/{claude,codex}/ledger.txt`. Pi not run and no preemptive xfail, per the entity's named-risk disposition.
- SKIPPED: Open your PR on top of #719, then extend the stack ... verify by GraphQL `pullRequest.stack` read-back.
  Held: opening the PR starts a reviewer run, which `## Review-finding disposition` forbids before FO authorization on an open Material finding.
- DONE: File the implementation stage report in the entity with every declared deviation, commit path-scoped, push, signal the FO, and stop for validation. Do not prepare or resolve any gate.
  This report. No gate prepared or resolved.

### Findings held for FO authorization

1. **Material, remedy proven, NOT applied.** `assertAutoContinueDispatchEvidence` resolves the report through `autoContinueWorktreeDir` but reads gates from the unresolved `entityPath`. When validation is worktree-backed and the FO keeps gate state in the worktree copy, a conforming run reds `entity has no gates record`. Observed live (claude run 2) with the FO's own final message confirming the gate prepared and open; the gate room was under `.worktrees/spacedock-ensign-auto-continue-task/`. This falsifies ideation's "host-neutral by construction ... `gates.Read(entityPath)`" claim, which survived only because the check had never run on claude. Harms `value-ac[AC-2]` (no false red on conforming behavior). Remedy proven in a throwaway: read gates from `reportEntity`, plus a two-direction regression guard (open-in-worktree GREEN, resolved-in-worktree still RED under `human-gate-bypassed`), falsified by reverting the one line.
2. **Needs decision, not owned by this task.** AC-4's claude bar inherits the fixture's validator nondeterminism: on the fix-validation run the validator recommended REJECTED on the stub deliverable, so the FO took the `feedback-to: implementation` route and never ran `gate prepare`. Note that run would have graded GREEN on claude before this change, so the red is the feature working. Remedies (fixture deliverable, or narrowing AC-4) are captain-owned.

### Deviations

- **Surface over tolerance:** ~433 code+doc insertions vs ~175 +/-60; 10 files, within the +/-3 file tolerance. Roughly 56 lines are code MOVED from live-tagged files to the default tag (offset by matching deletions). The rest is work ideation did not cost: the offline gate-resolved cell needs real durable state (gates frontmatter, git repos, both layouts), and the durableSemantic defect was unknown at ideation. I cut the one test with no AC mapping and trimmed comments; further cuts would delete AC-mapped coverage.
- **11th surface, required by AC-1:** `durableSemantic` now unwraps with `errors.As`. `livescenario.Run` wraps its Assert's error, so the bare type assertion relabelled `human-gate-bypassed` to the generic `auto-continue-state` — the exact stall/bypass ambiguity AC-1 exists to remove. Every other call site passes a bare error, so behavior there is unchanged.
- **Assertion moved to the default build tag:** `assertAutoContinueDispatchEvidence` and the two fixture writers it needs. AC-1 requires the resolved-gate cell to be verified offline, which is impossible while the assertion is live-tagged.
- **Declared file with nothing to change:** the expected-surface row for `shared_live_runner_test.go` ("retire the finding-3 AUDIT note") found no such note — `git log -S AUDIT` over `internal/ensigncycle/` returns nothing. No edit made.
- **Gate-resolved case graded under `human-gate-bypassed`:** ideation named the code only for the status half; both bypass shapes are the same defect, so both carry it.
- **One unrelated flake:** `TestCodexProcessRecognizesTerminalTurnBeforeOSExit` (250ms no-progress budget) failed once under full-repo parallel load, then passed 3/3 isolated, once per-package, and on a clean re-run of `go test ./...`. The base commit 571017df3 is green on the same package.

### Summary

The gated design landed and its offline ACs are proven and falsifiable, but verifying them surfaced two things ideation could not have known, both because this check had never run on claude before. The first is a Material false-red: the gate read does not follow the worktree resolution the report read already does, so a correct run reds. Its remedy is one line plus a two-direction regression guard, proven in a throwaway and held for authorization rather than applied. The second is a scope question — AC-4's claude bar depends on a stub-deliverable fixture whose validator sometimes rejects, and the honest remedies belong to the captain. The candidate is unchanged at 5cdb766ec, no PR is open, no gate touched, and both live loops are stopped.

### Implementation addendum (post-authorization cycle)

**Correction to an ideation claim.** The Spike recorded the gate-open half as "host-neutral by
construction ... `gates.Read(entityPath)`". That claim was only ever checked against codex,
because the optional interface meant no other host reached the check. The first claude run
through the now-unconditional check falsified it. Recorded here so the next reader does not
inherit the claim.

**The authorized remedy was itself falsified, by running it.** FO authorized reading gates from
the worktree-resolved `reportEntity`. Applied as 84ea009e4; codex then went from 3/3 GREEN to
RED twice on `entity has no gates record`. A controlled run of the same two commits cherry-picked
onto the OLD base (571017df3) also red, which rules out the rebase and the layers below.

Cause: the gate record and the validation report do not reliably live in the same copy of a
worktree-backed entity, and the placement differs by host — codex files the report in the
worktree copy and the gate record in the base copy; claude has been observed putting both in the
worktree copy. The original code and the authorized fix are each half right, and each reds a
conforming run on the host it does not match.

My first proof was insufficient: the throwaway fixture encoded the placement I hypothesised
rather than the placements observed, so it confirmed the expectation instead of testing the
variable. The corrected remedy — read the gate state from every entity copy carrying a record and
fail closed — is proven across all three placements in both directions, and falsified
symmetrically: gates-from-`entityPath`-only reds the claude shape, gates-from-`reportEntity`-only
reds the codex shape, both-copies passes all three while every resolved gate still reds under
`human-gate-bypassed`. Held for FO authorization; NOT applied to the candidate.

**Product-side observation, not acted on.** Codex's FO reported: "Gate preparation succeeded, but
presentation stopped because both required structured reads failed ... The report exists in the
registered worktree, but the status reader inspected the main entity copy." `spacedock status
--read` does not follow the worktree either. Outside this entity's surface; recorded, not fixed.

**Stack.** Base moved from #719 to `spacedock-ensign/prepare-initial-gated-stage-from-seed`
(58043ed97) per FO. Rebased cleanly, PR #723 opened as layer 5, linked `718 719 721 722 723`.
The `gh stack link` banner is not evidence, so verified by GraphQL `pullRequest.stack` read-back
from #723, #722 and #718 independently — all three report stack #720, size 5, correct order.

**Live loop accounting (AC-4 NOT met).** Pre-fix: codex 3/3 GREEN, claude 1 GREEN + 1 attributed
FALSE red + 1 killed. Post-authorized-fix: codex 2 RED (the falsification above), claude 1 killed
in flight. Old-base control: codex 1 RED. No streak stands on the current candidate; all pre-fix
greens need re-earning once the remedy is settled. Ledgers with per-run stream digests under
`/Users/clkao/.claude/jobs/4e49247e/tmp/live{,2}/`. Pi not run, no preemptive xfail.
