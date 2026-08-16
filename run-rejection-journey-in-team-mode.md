---
title: Run the rejection journey in team mode
status: ideation
source: "Captain ruling, 2026-08-16: do not invoke a journey in single-entity bare mode unless it is specifically testing bare-mode behavior"
id: zqb683j8jth0tyr2eme231e2
---

## Problem

The live rejection-flow journey starts the FO in single-entity `-p` bare mode. In bare mode the contract's feedback flow is sequential fresh dispatch: the addressable-worker reuse path the rejection skill routes through cannot happen at all. So the journey runs in a mode that disables the mechanism it exists to prove. The lifecycle assertion (internal/ensigncycle/claude_runtime_helpers_test.go:209) compounds it: it demands one implementation dispatch and one implementation report section, while the fixture itself scripts a second report section on the rework round — correct two-cycle runs graded red (`implementation-worker-not-dispatched`) in 3 of 8 observed live runs.

## Proposed approach

Invoke the rejection-flow scenario in team mode on both runtimes, so the routing flow under test can actually run. Derive the team-mode determined shape from the fixture and skill (the fixture declares no context-budget probe, so reuse is the required path; the fixture mandates a second implementation report section) and align the lifecycle assertion and grading to that shape: strict, no multiple-path acceptance. A conforming two-cycle run must grade green; a fresh-dispatch-when-reuse-required run must grade red. If bare-mode rejection semantics deserve their own coverage, that is a separately named scenario decided at ideation, not this journey's default.

## FO adherence residuals (folded in — captain ruling 2026-08-16)

The journey's proof includes the FO actually following the determined shape, so this entity also owns the two real adherence failures the blind analysis proved (\_debriefs/2026-08-16-03), rather than a separate behavior entity:

1. The FO records both validation rounds, then ends without ever invoking `gate prepare` — violating the rewritten skill's step-5 done-condition. Composed-tree loop run 1.
2. The codex FO fresh-dispatches the fix worker while the followup surface is live and every reuse condition passes — the route the codex adapter calls the load-bearing exception to casual fresh dispatch. Composed-tree loop run 3 (spawns=2, with `--advance` used before and after in the same stream).

Both were measured under the bare `-p` invocation this entity removes, so they must be re-measured under team mode: the team-mode proof loop IS the residual-rate measurement. If a mode survives there, ideation decides between targeted hardening of the skill/adapter text at the exact point the deviation enters (with falsifying evidence per mode) and holding this entity as the active owner of the residual rate until N clean cadence runs. Evidence pointers: loop run-1/run-3 streams under repair-codex-rejection-round-recording/evidence/ and the job scratch ledgers.

## Out of scope

- Product code: scenario invocation, grading, and FO instruction text only.

## Evidence pointers (FO, post-filing)

Two captain-ordered investigations sharpen this entity's scope; ideation must read both:

- `_debriefs/2026-08-16-02-live-harness-audit.md` — full live-suite audit. Findings 1, 2, 7, 8, 10 are this journey's surface: the codex leg's two graded assertions are mutually contradictory on the fixture's own scripted end-state (exact duplicate heading mandated by the fixture, hard-errored by the section selector); the Claude leg's either/or reviewer acceptance has no ordering anchor and no implementation-producer check; two `observed` checks are tautologies; the fixture-determined Cycle 2 line is never asserted.
- `_debriefs/2026-08-16-03-rejection-shape-attribution.md` — blind two-sided shape analysis of the five preserved runs. Key corrections to this entity's filed text: on codex the followup surface was LIVE, so the determined fix route is REUSE (spawns=1), not bare-mode fresh dispatch — the fresh-dispatch derivation governs the Claude `-p` leg only. Both codex "green" runs used a single-worker chain that re-reviewed its own fix, violating the skill's reviewer-independence ban — undetected because the codex branch swapped out the reviewer-flow assert; the one fixture-literal ensign was graded red by the duplicate-heading error. Minimal corrections identified there (prefix-count sections, fixture-side distinct cycle-2 title, reviewer-topology assert for codex) are candidate design inputs, not pre-approved decisions.

## Determined shape (derived from source, proven by spike)

The shape is determined per host by ONE contract-observable fact: whether the reuse route is
available for a completed worker. `fo-dispatch-core.md:49` (reuse condition 0) is fail-safe —
"if it reports the worker over budget, **or the probe is unavailable**, dispatch fresh".

| | Codex | Claude (headless `-p`) |
|---|---|---|
| `«context-budget»` | ABSENT → condition satisfied (`codex-first-officer-runtime.md:28`) | PRESENT, and it **exits 1** headless: `no subagent jsonl found for '<name>'` |
| Reuse route | `followup_task` to the live task path | none survives condition 0 |
| implementation spawns | 1 | 2 (round 1 + fresh rework) |
| validation spawns | 1 | 2 (cycle 1 + fresh re-review) |
| rework routing | `followup_task` → the implementation handle | `shutdown_request` to the superseded worker, then a fresh `Agent` |
| re-review routing | `followup_task` → the validation handle | `shutdown_request`, then a fresh validation `Agent` |

Host-neutral, in both branches:

- Ordered chain: impl dispatch → impl `Done:` → `status=validation` → validation dispatch →
  REJECTED report → `- Cycle 1:` line → `gate record --round validation/1` → `state commit` →
  rework routing → fix `Done:` → `status=validation` → re-review routing → PASSED report →
  `gate prepare` → stop nonterminal. (`feedback-rejection-flow/SKILL.md` steps 1-5.)
- Four durable report sections: `## Stage Report: implementation`, `## Stage Report: validation`,
  `## Stage Report: implementation (cycle 2)`, `## Stage Report: validation (cycle 2)` — the
  `(cycle N)` form is the ensign contract's own rule (`ensign-shared-core.md:93`).
- Exactly one `### Feedback Cycles` line, byte-equal to `inputs/feedback-cycle.txt`. The composed
  tree's fixture says a passing re-validation "closes its cycle without adding a line".
- Exactly one successful `gate record --round validation/1`, one `round-1` room, fix marker present,
  entity nonterminal at exactly one prepared open validation gate.
- The re-review reaches a worker whose stage is validation and NEVER the fix producer
  (`SKILL.md` step 4; registry required outcome "independently checked").

## Spike: the riskiest unverified mechanism, run first

One live run per runtime on the composed tree (571017df3) with the team-mode prompt and a persisted
native rollout. Artifacts: `/Users/clkao/.claude/jobs/4e49247e/tmp/spike-artifacts/{codex,claude}/`.

- **Codex (459s, RED on `rejection-gate-not-prepared` only).** The native rollout carries the exact
  determined reuse shape: `spawn_agent` ×2 (`…_rejection_task_implementation`,
  `…_rejection_task_validation`), then `followup_task` → implementation handle for the rework, then
  `followup_task` → validation handle for the re-review, each followed by that worker's `Done:`.
  The followup surface is live and the FO uses it. The rework ensign wrote
  `## Stage Report: implementation (cycle 2)`.
- **Claude sonnet (528s, GREEN).** Team mode engaged (4 named background `Agent`s), and every reuse
  attempt hit `dispatch context-budget` exit 1 (`no subagent jsonl found`) → fail-safe fresh
  dispatch with a `shutdown_request` to the superseded worker first. Fully conforming. The FO
  prepared the gate. The rework ensign wrote an EXACT duplicate `## Stage Report: implementation`,
  following the fixture's literal sentence.

## Corrections this design carries against its inputs

1. **The reuse path is not reachable on Claude headless at all** — not because of bare mode, but
   because `«context-budget»` is PRESENT and structurally unavailable in `-p` (no team
   `config.json`; `claude-fo-dispatch.md:119` says the same for reconcile). Grading Claude to a
   reuse shape would false-red a conforming run. The entity's Problem statement assumed team mode
   restores reuse on both hosts; it restores it on Codex only.
2. **The Codex public stream is topology-blind on the current build.** Its only `collab_tool_call`
   items are `wait` (verified across 16 preserved streams and the spike). The in-tree
   `assertCodexReviewerReuse` reads `collab_tool_call` `spawn_agent`/`followup_task` from the public
   stream, so it returns `reviewerUnsupported` — a non-`gradedErr`, i.e. an infrastructure FAIL —
   on every current run. It cannot be wired as-is. Topology lives in the native rollout, whose
   `spawn_agent` arguments are an encrypted blob EXCEPT the plaintext task path — so identity must
   be graded on task paths, never prompt content.
3. **Audit finding 8 is half-landed and half-stale.** The round-2 pointer alternation was already
   removed by stack layer 2 (`assertRejectionRecordedRound` pins `validation/1` and exactly one
   `round-1` room). The "fixture-determined `Cycle 2: PASSED` line" no longer exists: layer 1
   rewrote that fixture prose to "a re-validation that passes closes its cycle without adding a
   line". The assert to add is therefore EXACTLY ONE Cycle line, not a second one.
4. **`assertRejectionGatePrepared` is wired Codex-only** (`claude_live_runner_test.go:402`), and
   `assertRejectionRoundGateBoundary` returns nil when the entity has no gates record. FO residual
   mode 1 (ends without `gate prepare`) is therefore invisible on Claude and Pi today. It is a
   durable on-disk check and becomes host-neutral here.
5. **The fixture, not the model, is the heading variance.** The fixture's rework bullet instructs an
   exact duplicate heading; `ensign-shared-core.md:93` instructs `(cycle N)`. The Codex ensign
   followed the contract, the Claude ensign followed the fixture — same run, opposite bytes.

## Proposed approach

Each mechanism names the AC it serves, the simplest alternative, and why that alternative fails.

1. **Team-mode invocation, host-neutral cue + host realization** (AC-1). `rejectionPrompt` drops
   "Process only `rejection-task`" (which triggers single-entity mode, `first-officer/SKILL.md:7`,
   and thence bare dispatch, `claude-fo-dispatch.md:9`) and carries the mode requirement; the runner
   appends the host's realization sentence, the way it already appends `antiShutdownOverride`. The
   proven reference is `merged_team_mode_live_test.go:100-107`. *Alternative:* keep the single-entity
   prompt and amend `claude-fo-dispatch.md:9`. *Insufficient:* it rewrites a global dispatch default
   to fix one journey's invocation.
2. **Reuse condition 5 in `fo-dispatch-core.md`** (AC-2): a completed worker is never reuse-advanced
   into a stage whose `feedback-to` names its own stage. `fo-dispatch-core.md:57` already presumes
   it ("**If fresh dispatch:** if the next stage's `feedback-to` points at the completed stage, keep
   that agent alive"); the registry requires "independently checked". Without it BOTH the
   4-worker shape and the observed 1-worker self-reviewing chain are contract-conforming.
   *Alternative:* grade only the cycle-2 self-review ban. *Insufficient:* it leaves two conforming
   topologies — the multiple-path acceptance the captain ruled out — and still blesses a reviewer
   grading its own cycle-1 implementation.
3. **Branch-keyed per-cycle lifecycle grading** (AC-2, AC-3), replacing `assertWorkerLifecycle`'s
   `spawns != 1` for this journey: read the run's branch key (Codex: a `followup_task` route exists;
   Claude: the probe's exit status per worker), then require that branch's exact ordered chain —
   spawn counts per stage, routing target per round, and each round's report section, in order.
   *Alternative:* keep counting spawns without ordering. *Insufficient:* audit finding 2 — a run that
   spawns two validators up front, or never re-reviews the corrected candidate, counts the same.
4. **Codex topology extractor over the native rollout** (AC-2), correlating `spawn_agent`
   `function_call` → `function_call_output.task_name` → later `followup_task` target paths and
   `agent_message` authors. *Alternative:* `assertCodexReviewerReuse` on the public stream.
   *Insufficient:* proven blind (correction 2).
5. **Persist a per-run topology digest** into the artifact dir (AC-1's evidence requirement): the
   ordered `(index, event, stage, target)` rows the grader already extracts, ~2 KB. *Alternative:*
   keep reading the isolated `CODEX_HOME` rollout in-process only. *Insufficient:* `t.Cleanup`
   deletes it, which is why the shape-attribution debrief had to infer topology from `dispatch
   build` commands and why `_codex-home` was empty in every preserved run.
6. **Fixture heading repair** (AC-2, AC-3): the rework bullet instructs
   `## Stage Report: implementation (cycle 2)`, matching `ensign-shared-core.md:93`.
   *Alternative:* assert-side prefix counting that tolerates both forms. *Insufficient:* it keeps
   the fixture contradicting the ensign contract and lets the two hosts write different bytes for
   the same journey.
7. **Cycle-line and gate-prepared asserts, host-neutral** (AC-2): exactly one `- Cycle N:` entry in
   `### Feedback Cycles`, byte-equal to `rejectionFeedbackCycle`; exactly one prepared open gate.
8. **Tautology removal** (AC-2): `assertRejectionFlow`'s `observed` checks are graded against
   `finalMessage` alone and require the rejection AND the pass to be reported (the fixture prompt
   determines "report both outcomes"); the `"implementation"` token is dropped. The report-count
   checks become exactly 2 per stage.
9. **Shared-helper repairs** (audit finding 10, AC-2): one ordering semantic (first successful
   transition) for all three dialects, match the command field rather than the raw JSON line, anchor
   stage tokens instead of substring-matching descriptions, delete the callerless
   `successfulStatusSet(allowed ...string)` variant.

## Bare-mode coverage: dropped here, already owned

Dropped from this journey. A dedicated forced-bare scenario is already the filed backlog entity
`bare-mode-coverage` (2ehx77vn94s20k8cqs2p8vhw), whose own seed records that unpinned `claude -p` is
team-enabled by default and that bare therefore needs an EXPLICIT trigger rather than an incidental
one. Adding a second bare scenario here would duplicate that entity's deliverable. This entity's
change is recorded in that entity's `## Problem` as one more `-p` scenario that flipped away from
incidental bare coverage.

## FO adherence residuals: measurement and decision rule

The team-mode proof loop IS the residual measurement; both filed modes are re-measured by it, and
mode 1 already recurred once under team mode (the Codex spike red).

- Mode 1 — ends without `gate prepare`. Detected host-neutrally after change 7 (today it is Codex-only).
- Mode 2 — fresh-dispatches the fix target while the followup route is live. Detected by change 3's
  branch check: on Codex, a `spawn_agent` for the rework while a `followup_task` route exists reds.
- Mode 3 (new, from the spike) — the Codex FO tore workers down with `interrupt_agent`, which
  `codex-first-officer-runtime.md:13` explicitly refuses to bless. Recorded, not graded: teardown is
  outside this journey's contract surface.

**Rule.** A mode that appears in the AC-1 loop gets ONE targeted hardening round at the exact
instruction the run's own stream shows the FO following into the deviation — quoted in the entity as
that mode's falsifying evidence — and the loop is re-run. If the mode survives two hardening rounds,
this entity escalates to the captain with the measured rate and stands as active owner of it until N
clean cadence runs; it does NOT loosen the grade to reach green. Hardening without a quoted
instruction is guesswork and is not permitted by this rule.

## Expected surface and tolerance

| File | Change | Lines |
|---|---|---|
| `internal/ensigncycle/shared_fixtures_test.go` | team-mode prompt, rework heading instruction | ~12 |
| `internal/ensigncycle/claude_live_runner_test.go` | host cue, assert wiring, digest persistence, host-neutral gate check | ~45 |
| `internal/ensigncycle/claude_runtime_helpers_test.go` | branch-keyed per-cycle lifecycle, finding-10 repairs, dead-variant delete | ~85 |
| `internal/ensigncycle/shared_reviewer_reuse_test.go` | Codex native-rollout topology extractor, branch-keyed Claude assert | ~90 |
| `internal/ensigncycle/shared_round_recording_test.go` | Cycle-line assert | ~25 |
| `internal/ensigncycle/shared_assertions_impl_test.go` | exact report counts, final-message anchoring | ~20 |
| `internal/ensigncycle/shared_reviewer_reuse_table_test.go` | negative controls, both hosts | ~110 |
| `internal/ensigncycle/testdata/rejection_topology/*.jsonl` | captured spike topology + a self-review mutant | ~120 (data) |
| `internal/ensigncycle/shared_live_runner_test.go` | retire the AUDIT block | ~12 |
| `skills/first-officer/references/fo-dispatch-core.md` | reuse condition 5 | ~3 |
| `docs/runtime-live-ci-registry.md` | required-outcome clause | ~4 |

11 files, ~406 changed lines plus ~120 lines of captured test data. Tolerance ±30% on lines, ±1 file.
No product code: every Go file above is `_test.go`.

**Observable semantics this task may change**

1. Scenario invocation — the journey proves team-mode routing; bare-mode rejection routing loses its
   incidental coverage (owner: 2ehx77vn94s20k8cqs2p8vhw).
2. Fixture instruction text — the durable end-state gains a distinct cycle-2 implementation heading.
3. Grading vocabulary — new codes `rejection-worker-topology` and `rejection-cycle-line`;
   `rejection-gate-not-prepared` starts firing on Claude and Pi.
4. **FO dispatch behavior, all workflows** — reuse condition 5 makes every review stage whose
   `feedback-to` names the completed worker's stage a fresh dispatch, including this repo's own dev
   workflow (`docs/dev/README.md:25-28`). This is the one change outside the live harness.
5. Artifact surface — one new per-run digest file under the scenario artifact dir.

## Documentation diff

`docs/runtime-live-ci-registry.md`, `### rejection-flow`:

```diff
 - **Required outcome:** A rejected candidate is corrected and independently
   checked before a fresh final gate is presented. Rejected authority cannot
   satisfy the final approval.
+  The journey runs in team mode: the correction is routed to the producer and
+  the re-review to a worker that did not produce the fix, through whichever
+  route the host's reuse conditions leave available.
```

`skills/first-officer/references/fo-dispatch-core.md`, `**Reuse conditions**`:

```diff
 4. `«reuse.model-match»` — the reused worker's stamped model matches `next_stage.effective_model`.
+5. The completed worker's stage is NOT the next stage's `feedback-to` target. A review of the
+   worker's own output is not an independent check, so a review stage is dispatched fresh; the
+   producer stays alive for the correction routing (see **If fresh dispatch** below).
```

`internal/ensigncycle/shared_live_runner_test.go`: the `AUDIT(2026-08-16)` block above
`TestLiveCommonRejectionFlow` is deleted; the Pi `liveXFail` keeps a one-line note that it is
structural harness blindness (finding 11), owned elsewhere.

## Acceptance criteria

- **AC-1 (value).** On the composed tree plus this layer, under the CI codex shim, the focused
  journey reaches **3 consecutive conforming greens per runtime** (codex, claude-sonnet) within 8
  runs per runtime. A conforming green is `go test -tags live -run TestLiveCommonRejectionFlow`
  exit 0 AND a persisted topology digest showing that run's branch chain in order. Any red resets
  the count; every run is one ledger row. *Baseline that can move the wrong way:* the pre-change
  measure is **0 consecutive conforming greens on Codex** — the retained loop ledger shows 2 passes
  in 4 runs and both passes certified a single-worker chain that re-reviewed its own fix — and
  ungraded on Claude, where no topology check exists.
  *Proof:* the ledger plus per-run digests, committed under the entity's evidence directory.
- **AC-2 (a shape-violating run grades red).** Every non-conforming topology in hand grades red with
  a named code: the preserved self-review greens (`loop/run-2`, `loop/run-4`), a fresh-rework mutant
  of the Codex spike rollout, a re-review-before-rework mutant, a duplicate-heading entity, a
  two-Cycle-line entity, an entity with no prepared gate. *Proof:* table tests over the captured
  streams and their mutants; each case names the code it must produce.
- **AC-3 (a conforming run grades green).** Both spike runs' captured artifacts grade green under
  the new asserts on their own branch — the Codex reuse chain and the Claude fail-safe-fresh chain —
  and neither passes under the other branch's expectations. *Proof:* replay table test over the two
  captured topology fixtures.
- **AC-4 (no false red from the contract's own fail-safe).** A run whose `dispatch context-budget`
  probe exits non-zero and which then fresh-dispatches is green; the same run is red if it reuses
  without a successful probe. *Proof:* two synthetic Claude streams differing only in the probe's
  result and the routing that follows.

## Test plan

| What | Proves | Cost |
|---|---|---|
| Table tests over captured + mutated topology fixtures (both hosts) | AC-2, AC-3, AC-4 | ~2 h to write; runs in seconds, no live cost |
| `go test ./...` and `-race` | no collateral damage; all changes are `_test.go` plus two prose files | ~5 min |
| `contractlint` registry reconciliation | the registry entry still binds the journey/fixture ids | seconds |
| Live loop, codex, CI shim, ledger + digest per run | AC-1 | ~8 min/run, ≥3 runs, budget 8 |
| Live loop, claude-sonnet, ledger + digest per run | AC-1 | ~9 min/run, ≥3 runs, budget 8 |
| Pi | not run — `liveXFail` stays; the Pi lane's graders read the Claude dialect (finding 11), so no behavior change can XPASS it. Owner stays as filed. | none |

Live budget: ~2.5 h wallclock with the two runtimes in parallel, ~5 h if both hit the 8-run cap.

## Stage Report: ideation

- DONE: Read, in order: the entity body, audit findings 1/2/7/8/10/11, the shape-attribution debrief in full, the AUDIT note, and the two landed stack layers.
  Composed tree read at `.worktrees/spacedock-ensign-repair-codex-rejection-round-recording` (571017df3 = #718's 3 commits + #719's 3); `gh pr view 720` does not resolve in this repo — the two layers were read from the branch, not from a stack PR.
- DONE: Derive the team-mode determined shape per runtime from source, with citations.
  `## Determined shape` above; the debrief's corrected derivation was checked against source and holds for Codex, and its Claude half is superseded by the probe finding (correction 1).
- DONE: Design the invocation change for both runtimes; decide the fate of bare-mode coverage.
  Host-neutral cue in `rejectionPrompt` plus a per-host realization sentence, modelled on `merged_team_mode_live_test.go:100-107`; bare-mode coverage dropped here and left to the already-filed `bare-mode-coverage` (2ehx77vn94s20k8cqs2p8vhw).
- DONE: Design grading strictly to the determined shape, no multiple-path acceptance.
  Changes 3, 4, 6, 7, 8, 9 in `## Proposed approach`; the only branch is keyed on an observed host fact (the probe's exit), not on accepting two behaviors.
- DONE: The value AC as a measurable loop criterion, both grading directions falsifiable.
  AC-1 (3 consecutive conforming greens per runtime vs a measured 0/4 conforming baseline), AC-2 (red on every non-conforming shape in hand), AC-3/AC-4 (green on both conforming branches).
- DONE: Design the decision rule for surviving FO-adherence residuals.
  `## FO adherence residuals`: one evidence-bound hardening round per mode, quoted instruction required, escalate-and-own after two rounds, never loosen the grade.
- DONE: Spike the riskiest unverified mechanism first and record the result.
  Two live runs on the composed tree. Codex: `spawn_agent` ×2 then `followup_task` to the implementation handle and to the validation handle, in order — the reuse route is live and used. Claude: team mode engaged, every `dispatch context-budget` probe exited 1, so all four dispatches were fail-safe fresh. Falsifying observation for each: had the Codex rollout shown a second `spawn_agent` for the rework, the reuse derivation would be wrong; had the Claude probe exited 0, the fresh-dispatch branch would be a deviation rather than the determined shape.
- DONE: Declare the expected surface and every observable semantic the task may change; test plan per AC with cost estimates.
  `## Expected surface and tolerance` (11 files, ~406 lines, ±30%, no product code) and `## Test plan`. The load-bearing declaration is semantic 4: reuse condition 5 changes FO dispatch behavior for every workflow with a `feedback-to` review stage, including this repo's own.
- DONE: Write the ideation stage report and stop for the gate; no implementation started.
  The spike ran in a throwaway detached worktree under the job scratch dir; nothing was written to the stack branch.

### Summary

The design pins one determined shape per host with a single branch keyed on an observable — whether
the reuse route survives reuse condition 0 — and grades the whole ordered chain rather than a spawn
count. Two live spike runs settled the two facts the design rests on: on Codex the `followup_task`
reuse route is live and used for both the rework and the re-review, and on Claude headless the
context-budget probe fails, so fail-safe fresh dispatch is the conforming shape and a "must reuse"
assert would have false-redded a green run. Three corrections came out of the spikes that change
what implementation must build: the Codex public stream carries no topology at all (only the native
rollout does, and the harness deletes it), the in-tree `assertCodexReviewerReuse` is therefore
unusable as written, and the fixture's duplicate-heading sentence — not model variance — is what
made two hosts write different bytes for the same round. One change reaches outside the harness: a
fifth reuse condition making a review stage always fresh, without which both the 4-worker shape and
the observed self-reviewing chain remain contract-conforming.
