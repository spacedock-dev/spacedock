---
title: Auto-continue grades a human-gate bypass as green
status: validation
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
        - id: gate:7xe7hxt1qce1x9b3dm0k6ymg:validation
          stage: validation
          attempts:
            - id: gate-attempt:7xe7hxt1qce1x9b3dm0k6ymg-validation-1
              briefing:
                id: briefing:7xe7hxt1qce1x9b3dm0k6ymg:validation:attempt-1:revision-1
                digest: sha256:369de6d81b320199e2f24ffde2f8f9221e3cb68cd97ddafb31eca050c211e102
                request-digest: sha256:9c925b108c438255e0f5d3f236f2f6bcd7c547b5e2e77d45d998744b5a423359
                room-ref: ./red-auto-continue-gate-bypass/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:7xe7hxt1qce1x9b3dm0k6ymg:validation:1
                briefing: briefing:7xe7hxt1qce1x9b3dm0k6ymg:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-17T02:56:55.468609Z"
                decision: approve
              application:
                target-stage: done
                state: pending
worktree: .worktrees/spacedock-ensign-red-auto-continue-gate-bypass
mod-block: merge:pr-merge
pr: local-merge:61dd8e435
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
  All four gated elements landed in the first commit (5cdb766ec, now 355124469 after the restack). FINAL surface at `134e494a8`: 11 files vs 10 +/-3 (within tolerance) and 856 code+doc insertions vs ~175 +/-60 (over, by ~3.6x the estimate). The growth is the first commit's ~433 plus the three authorized finding fixes and their falsification tables; nothing absorbed beyond the gated surface, every increment FO- or captain-authorized. Accounted in Deviations below.
- DONE: Branch from the current stack tip: spacedock-ensign/repair-codex-rejection-round-recording at 571017df3 ... Push via the SSH remote.
  Branched at 571017df3 and pushed via the SSH remote. Rebased twice under FO direction as the stack grew, finally onto `b0bb159ff`; the layer's own diff was verified byte-identical across the restack before the content-identity declaration was made.
- DONE: Offline proof before any live spend, per the gated ACs ...
  AC-1 table: {claude,codex,pi} x {done end-state, resolved-gate} all RED under `human-gate-bypassed`, each host with a conforming open-gate control so the table cannot pass on an assertion that reds unconditionally. AC-2: the real captured stream and both fixture layouts GREEN. `-run AutoContinue -count=20`, `-race`, full `go test ./...`, `-race -timeout 30m`, and contractlint all clean. Falsified three ways: re-adding `done` reds all three hosts' terminal cells; reverting the durableSemantic unwrap reds the code-survival test; disabling the gate-open check reds all three resolved-gate cells.
- DONE: Live loop per AC-4: ... >=3 consecutive greens per runtime (claude-sonnet, codex) within budget 8 per runtime ...
  MET on both runtimes on the final candidate `134e494a8`: claude-sonnet runs 1/2/3 all GREEN and codex runs 1/2/3 all GREEN, each LOOP COMPLETE at run 3 within budget 8, both re-earned from a zero streak on the shipping bytes. Zero reds; every ledger row carries `codes=[]`, so no run counted green while grading on a code this change introduces, and the never-prepares mode did not recur. Ledgers with per-run per-variant stream digests at `/Users/clkao/.claude/jobs/4e49247e/tmp/live5/{claude,codex}/ledger.txt`. Pi not run and no preemptive xfail, per the entity's named-risk disposition. The earlier exhausted and discarded ledgers (`live/`, `live2/`, `live4/`) stand below as the evidence trail: they measured a grader that was still wrong, which is why they were re-earned rather than banked.
- DONE: Open your PR on top of #719, then extend the stack ... verify by GraphQL `pullRequest.stack` read-back.
  PR #723 open. Base moved twice by FO direction (#719 -> #722's branch, then the restack to `b0bb159ff`); it now sits at layer 5 of stack #720, verified by GraphQL `pullRequest.stack` read-back from #723, #722 and #718 independently rather than from the `gh stack link` banner. Held initially — opening a PR starts a reviewer run, which `## Review-finding disposition` forbids while a Material finding is unauthorized — and opened once the FO authorized the first fix.
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

### Finding 3 — `withdrawn` is not `bypassed` (held for FO authorization)

The authorized fail-closed rule reads "any copy showing a non-open gate is a bypass". Claude run 2
on the corrected candidate red under `human-gate-bypassed` with gate state `withdrawn`, on a run
that behaved correctly: the FO prepared against the stale base copy, caught it before presenting,
withdrew that room and re-prepared against the worktree. Withdrawal records no decision —
`attemptState` returns `closed` only when a Resolution exists — so `closed` is the bypass and
`withdrawn` is the sanctioned retraction. The `!= "open"` predicate predates this entity; reading
both copies is what made it reachable.

Proposed, not applied: any copy `closed` reds under `human-gate-bypassed`; else any copy `open`
grades GREEN; else a plain error for a gate that was never left open, which is a not-presented
failure rather than an accusation. Proven by `TestAutoContinueGateStateVocabulary` over the live
withdrawn shape, the resolved shape and withdrawn-everywhere, and falsified by restoring the
`!= "open"` rule.

**Second proof-discipline failure, recorded because it is the same class as the first.** The
initial version of that test passed under BOTH rules, which was impossible. Its withdrawn fixture
did not parse — `Withdrawal` carries no `type:` field and requires `agent:first-officer`
attribution — so `gates.Read` errored, the copy was skipped and the case asserted nothing. The
falsification step caught it, as it caught the earlier hypothesis-encoded fixture. The lesson now
applied: probe what a fixture actually produces before trusting a test built on it, and treat a
falsification that fails to fail as a defect in the test rather than a strength of the code.

Loop state on 98d4758bc, both stopped: claude run1 GREEN, run2 RED (this finding); codex run1 and
run2 GREEN, consec 2, run3 killed in flight. Codex passing twice after two pre-correction REDs is
the evidence that the placement correction itself holds. Offline on this candidate: full
AutoContinue suite, `-count=20`, `-race` and contractlint all green.

### Implementation close-out (AC-4 partial; codex met, claude blocked out of scope)

**Findings 1 and 3 landed, both FO-authorized, both additive.** Final layer: `424ed85fe` (align the
gate read to the report's entity copy), `eb4e13625` (read every copy that carries a record; the
authorized single-path remedy was falsified by codex live and by an old-base control that ruled out
the rebase), `671fe9746` (accuse only a resolved gate — `withdrawn` is sanctioned self-correction,
not a bypass).

**Gate-state vocabulary, enumerated from source at FO requirement.** `attemptState` has exactly four
returns and `Summary.State` is assigned from it alone, so no caller widens the range. `invalid` is
UNREACHABLE as a state — the decoder rejects a conflicting withdrawal+resolution before
`CurrentSummary` sees it — so an earlier `invalid` branch was dead code and was removed. What
replaced it: a copy with no record stays benign, but a gates block that exists and fails to decode
now fails hard instead of being skipped into a green.
`TestAutoContinueGateFixturesParseAsIntended` pins each fixture to the state it names and pins that
the decoder still rejects `invalid`.

**AC-4 result on `671fe9746`.**

- **codex: MET.** runs 1-3 all GREEN, LOOP COMPLETE at run 3, budget 8. No leg graded on any code
  this change introduces.
- **claude: NOT MET.** 8 runs, 4 GREEN / 4 RED, best streak 2, budget exhausted.

All four claude reds are `implementation-worker-not-dispatched` from `assertWorkerLifecycle` — the
shared helper this entity's Out of scope section names as audit finding 10's territory and consumes
as-is. Two distinct modes, neither a bypass and neither graded by anything this change added:

- `spawns=2` (run 2): the validator recommended REJECTED, the FO routed through feedback, dispatched
  a second validator, cycle 2 recommended PASSED, gate left open. The end state is exactly what the
  journey requires; only `spawns != 1` failed.
- `spawns=0` (runs 4, 5, 6): stage matching is `strings.Contains(lower(description), "validation")`
  and the FO's real dispatch descriptions were `Validate auto-continue-task implementation`,
  `Revise auto-continue-task implementation`, `Re-validate auto-continue-task cycle 2` — none
  contain the noun. Three agents dispatched, zero counted.

So claude's green is a coin flip on FO phrasing and cycle count, which is the concrete mechanism
behind the finding-2 variance escalated to the captain. Not fixed here: widening the matcher or the
spawn count changes what the harness certifies for every journey using the helper. No named risk and
no xfail recorded, per FO instruction.

This is the change working as designed. AC-3 made the dispatch-evidence check run on claude for the
first time; the first thing it found is that the helper it calls is not robust to claude's dispatch
phrasing. Discovery on a host never exercised before, not regression.

**Content identity across the restack.** Rebased from base `58043ed97` onto `b0bb159ff` (confirmed by
`ls-remote` at rebase time). Commit map `355124469 -> 5aeb...`/`84ea009e4 -> 424ed85fe`,
`98d4758bc -> eb4e13625`, `2b9efb97f -> 671fe9746`. The layer's own diff is byte-identical across the
move: `git diff {base}..HEAD | sha256` is
`ff7e735fa35e085fae35bf43872dc5d2ccdd18f51c2a62b5e85e2934e3752493` both before and after, 10 files
+993/-85 unchanged. Post-rebase `go vet` on both tags, the AutoContinue suite, contractlint, gofmt and
a full `go test ./...` all agree, so the loop evidence above stands on the rebased bytes. The nested
inline state-commit fix beneath cannot interact with this journey: every fixture builds in `t.TempDir()`
and is never nested below a git root.

**Stack.** Force-pushed with a lease read from `ls-remote`. Final head `671fe9746`. PR #723 is layer 5
of stack #720, verified by GraphQL read-back: 718, 719, 721, 722, 723, 724.

### AC-4 MET on both runtimes — captain-approved fixes (final)

Captain approved both fixes behind the two claude failure modes; landed additively as
`134e494a8`.

**Fixture verdict pin.** The validation stage asked a live model to judge a stub deliverable
against no stated criterion, so the verdict was a coin flip and a REJECTED sent a conforming run
down the feedback route. One sentence now ties AC-1 to the presence of the implementation report.
Rejection handling stays the rejection journey's subject.

**Dispatch-pointer spawn anchor.** Descriptions are free text the FO composes; the live run
dispatched `Validate auto-continue-task implementation` and `Re-validate auto-continue-task
cycle 2`, neither containing the noun, so three agents counted as zero. `spacedock dispatch build`
generates the pointer filename, so its trailing `-{stage}.md` is contract-fixed. The description
check is kept, so nothing already counted stops counting; the pointer only adds spawns the
contract labels, and it excludes the reviser (`-implementation.md`). Falsified both ways: drop the
pointer and the validators go uncounted; ignore the stage suffix and the reviser counts.

The FO's hypothesised `name`-field anchor does not exist — read-only inspection showed every Agent
block carries only description/prompt/run_in_background/subagent_type, with `subagent_type` a
constant `spacedock:ensign`. The prompt's dispatch pointer is the contract-determined signal that
was actually available.

**Regression fixture.** The real run-4 split-root stream, in the entity's pre-approved line-type
filtered form: 975079 bytes / 718 lines down to 149187 / 122, sha256
`870418e9dca8b163eda734e840079037b1999dd507a10debc0453e85f67d02f2`. Capture-time identical-verdict
proof: `assertWorkerLifecycle` returns the same verdict on filtered and full for BOTH stage
`validation` and stage `implementation`; only the line indices shift, which filtering necessarily
changes. No-regression evidence for the anchor is the 14-leg sweep — every leg passing under
description-only matching counted exactly 1 under both signals.

**AC-4 RESULT on `134e494a8`: MET on both runtimes, no reds at all.**

| runtime | runs | result |
| --- | --- | --- |
| claude-sonnet | 1, 2, 3 all GREEN | LOOP COMPLETE at run 3, budget 8 |
| codex | 1, 2, 3 all GREEN | LOOP COMPLETE at run 3, budget 8 |

Both re-earned from a zero streak on the shipping bytes; codex's earlier 3/3 was discarded because
it predated the fixture pin. Every leg carries `codes=[]`, so no run counted green while grading on
a code this change introduces. Ledgers with per-run per-variant stream digests at
`/Users/clkao/.claude/jobs/4e49247e/tmp/live5/{claude,codex}/ledger.txt`. The never-prepares mode
did not recur, so no instruction quote was needed. Pi not run, no preemptive xfail.

Claude went from 4 GREEN / 4 RED with no streak past 1, to 3/3 clean — the two modes are closed.

## Stage Report: validation

- DONE: Source of truth is the entity's stage report and its addenda; the shipping evidence is live5/ on 134e494a8. Validate on the tip 52b94f9aa; this layer is #723.
  Validated on 52b94f9aa in a detached throwaway checkout (never the implementation worktree); the branch head is 134e494a8, one layer below the tip, and the tip layer's `stageToken` hardening of the shared matcher is included in everything below.
- DONE: Verify the gated ACs plus the three authorized finding fixes on the tip bytes.
  AC-1 PASSED: TestAutoContinueBypassRedsOnEveryHost — all six {claude,codex,pi} x {done, resolved-gate} cells red under `human-gate-bypassed` beside a conforming open-gate control per host; deleting pi's `lifecycleStream` fails `go vet -tags live` with "does not implement liveDriver (missing method lifecycleStream)".
  AC-2 PASSED: replay fixture digest re-verified (465e849a…, 430087 bytes) and green through both the lifecycle parse and the full dispatch-evidence check on both layouts; the placement table passes green/red in both directions across all three observed placements; the vocabulary table never accuses `withdrawn` (withdrawn_beside_open GREEN; withdrawn_only red without the bypass code).
  AC-3 PASSED: `verifyAutoContinueDispatch` absent repo-wide; `lifecycleStream` is a required `liveDriver` method with exactly the three host implementors; runAutoContinueJourney calls assertAutoContinueDispatchEvidence unconditionally, no type assertion.
  AC-4 PASSED: live5 ledgers 3/3 GREEN per runtime, `codes=[]` every row, LOOP COMPLETE at run 3 of budget 8; ledger stream digests re-verified against the retained stream bytes (claude run1 single/split match); the loop started nine seconds after 134e494a8's commit, consistent with re-earning from zero on the shipping bytes.
- DONE: Spot-verify by the falsifying edits the report claims.
  All five behave exactly as claimed on the tip bytes: re-adding `done` reds all three hosts' terminal cells plus the negative and code-survival tests; reverting the errors.As unwrap reds exactly TestAutoContinueBypassCodeSurvivesTheScenarioRunner; gates-from-entityPath-only reds exactly the worktree-both (claude) placement in both directions; gates-from-reportEntity-only reds exactly the split (codex) placement; dropping the pointer anchor uncounts both real validators; ignoring the stage suffix counts the reviser and the backlog pointer.
- DONE: Verify the fixture-parse guard (TestAutoContinueGateFixturesParseAsIntended) and the line-type-filtered regression stream's identical-verdict provenance.
  The guard passes and pins `invalid` as a decoder rejection. The filtered run-4 fixture (870418e9…, 149187 bytes) is a strict byte-subset (122/122 lines) of the retained full capture whose sha256 matches the recorded 2895ae23…; re-ran assertWorkerLifecycle on full vs filtered for stages validation AND implementation — identical verdicts (spawns=2 both), only line indices shift.
- DONE: Assess the declared deviations for the gate.
  856 code+doc insertions confirmed exact by numstat (plus 402 fixture lines; 9 code/doc files + 2 fixture files). The growth is the three authorized fix arcs and their falsification tables; the AC bar was not narrowed — the fixture verdict pin is the captain-approved fixture change and its sentence is verified present in autoContinueReadme. The `status --read` product observation is recorded-not-acted, outside this surface. No design reset warranted: the mechanism held, and the overage bought the placement/vocabulary coverage that caught two real defects.
- DONE: Reference tip CI run 31986034831's auto-continue lanes when concluded; no new live loops.
  The run concluded failure overall, but all four auto-continue legs passed (journey-metrics outcome.status=passed: claude single/split-root, codex single/split-root); the two live reds are other lanes — codex ac-reanchor, claude rejection-flow (the tip layer's own surface). The offline job succeeded. Validation ran no live loops.
- DONE: Validation stage report: per-AC verdict, PASSED or REJECTED recommendation, path-scoped commit, push, signal FO, stop. No gate.
  This report; recommendation below.

### Summary

PASSED. All four gated ACs verified on the tip bytes with reproduced evidence, all five claimed falsifications behave exactly as reported, gofmt clean, the registry entry matches the gated doc diff, and tip CI's auto-continue lanes are green on both runtimes. Offline on the tip: full `go test ./...` green in every package, with one caveat — the ensigncycle package alone runs ~5-8.5 min (solo `ok` 508.95s; `-race` solo `ok` 319.4s), so two concurrent instances blow the 600s default timeout, which my first full+race parallel attempt did; contention, not a hang. Deferred risks, none material: pi's stream dialect remains unexercised (pre-existing named risk, no xfail per the entity's disposition); the shared helper's spawns!=1 ordering semantics remain audit finding 10's territory; ensigncycle's headroom under Go's default 10m test timeout is thin and will flake local full-suite runs under load (promote if it times out without external contention).
