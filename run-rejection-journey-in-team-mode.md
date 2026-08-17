---
title: Run the rejection journey in team mode
status: validation
source: "Captain ruling, 2026-08-16: do not invoke a journey in single-entity bare mode unless it is specifically testing bare-mode behavior"
id: zqb683j8jth0tyr2eme231e2
gates:
    version: 1
    records:
        - id: gate:zqb683j8jth0tyr2eme231e2:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:zqb683j8jth0tyr2eme231e2-ideation-1
              briefing:
                id: briefing:zqb683j8jth0tyr2eme231e2:ideation:attempt-1:revision-1
                digest: sha256:76ec26b36edcd5ba5fd2fd3cbbe637869ee5eb8705540f8913f822ec390b7d1e
                request-digest: sha256:c56f2bba1b5866a1faea66a0f32e016f75bb41578fb3a6e6f5b50cf0f2bd6fe7
                room-ref: ./run-rejection-journey-in-team-mode/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zqb683j8jth0tyr2eme231e2:ideation:1
                briefing: briefing:zqb683j8jth0tyr2eme231e2:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-16T17:46:40.418163Z"
                decision: revise
                reason: 'Captain annotation on the determined-shape table (claude headless probe row): ''is this true? this used to work'' - verify the probe-failure claim and its history before the claude branch is accepted'
            - id: gate-attempt:zqb683j8jth0tyr2eme231e2-ideation-2
              briefing:
                id: briefing:zqb683j8jth0tyr2eme231e2:ideation:attempt-2:revision-1
                digest: sha256:2359dd68ea2549cd8c8ded29e0e9acbf91701ddc636571ff694843627d3ee5a1
                request-digest: sha256:065aad3085105141c550bcac8c8fdd0613cff7bf97b8e197bb9137ad507e4f65
                room-ref: ./run-rejection-journey-in-team-mode/review/ideation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:zqb683j8jth0tyr2eme231e2:ideation:2
                briefing: briefing:zqb683j8jth0tyr2eme231e2:ideation:attempt-2:revision-1
                by: person:captain
                at: "2026-08-16T18:01:42.457451Z"
                decision: approve
              application:
                target-stage: implementation
                state: consumed
        - id: gate:zqb683j8jth0tyr2eme231e2:validation
          stage: validation
          attempts:
            - id: gate-attempt:zqb683j8jth0tyr2eme231e2-validation-1
              briefing:
                id: briefing:zqb683j8jth0tyr2eme231e2:validation:attempt-1:revision-1
                digest: sha256:7bb7c613a706aa23cf182ffe91a2eea5836a76282781a64cc5cd4e373d264b1a
                request-digest: sha256:6162d7e29a40e917e55f2a187b0e11a41f0d2ae31c5de00097abb970c2ad6edd
                room-ref: ./run-rejection-journey-in-team-mode/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:zqb683j8jth0tyr2eme231e2:validation:1
                briefing: briefing:zqb683j8jth0tyr2eme231e2:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-17T02:57:06.229627Z"
                decision: revise
                reason: 'Captain approve-all: ratifies the staged revise round (attempt-N generalization of the main-side pins; ac-reanchor fixture identity; heading ruling = STRICT + pin the level in fixture instructions). Re-validate and re-gate after the fixes land and the tip is green.'
worktree: .worktrees/spacedock-ensign-run-rejection-journey-in-team-mode
pr: "#724"
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
| `«context-budget»` | ABSENT → condition satisfied (`codex-first-officer-runtime.md:28`) | PRESENT, and it cannot succeed on headless `-p` — three independent blockers, verified below |
| Reuse route | `followup_task` to the live task path | none survives condition 0 |
| implementation spawns | 1 | 2 (round 1 + fresh rework) |
| validation spawns | 1 | 2 (cycle 1 + fresh re-review) |
| rework routing | `followup_task` → the implementation handle | `shutdown_request` to the superseded worker, then a fresh `Agent` |
| re-review routing | `followup_task` → the validation handle | `shutdown_request`, then a fresh validation `Agent` |

Host-neutral, in both branches:

- Ordered chain: impl dispatch → impl `Done:` → `status=validation` → validation dispatch →
  REJECTED report → rework routing → fix `Done:` (the worker's own entries close the round's
  review log) → `- Cycle 1:` line → `gate record --round validation/1` (reporting all four
  accumulated entries) → `state commit` → `status=validation` → re-review routing → PASSED
  report → `gate prepare` → stop nonterminal. (`feedback-rejection-flow/SKILL.md` steps 1-5.)

  **Corrected at implementation (FO-authorized, implementation-discovered truth).** This chain
  previously placed `gate record --round validation/1` and the Cycle line BEFORE the rework
  routing, which inverts the skill it cites. `SKILL.md` step 1 (deliver the correction) is
  "**Done when** the correction is complete in durable workflow state: the target worker's own
  entries closing this round's review log"; only then does step 2 (record this round, once)
  append the Cycle line and invoke the recorder, "**Done when** the recorder exits successfully
  reporting the complete round summary, counting every entry this round accumulated". So the
  correction lands first and the single recording counts all four entries — the reviewer's two
  plus the rework worker's two. Recording before the rework reaches only `entries=2`, an
  incomplete round the skill tells the FO to hold on rather than claim. Claude AC-1 loop run 2
  followed the old, wrong order here, recorded at `entries=2`, and was correctly graded red by
  the pre-existing `entries=4` recognizer (`rejectionRoundSuccess`), which was left untouched:
  relaxing it to accept the incomplete round would have loosened the grade to match a defective
  chain. Nothing in the implementation depended on the wrong ordering — the branch-keyed
  topology grader orders only worker routing events (`spawn`/`reuse`/`done`) and never places
  `gate record` among them.
- Four durable report sections: `## Stage Report: implementation`, `## Stage Report: validation`,
  `## Stage Report: implementation (cycle 2)`, `## Stage Report: validation (cycle 2)` — the
  `(cycle N)` form is the ensign contract's own rule (`ensign-shared-core.md:93`).
- Exactly one `### Feedback Cycles` line, byte-equal to `inputs/feedback-cycle.txt`. The composed
  tree's fixture says a passing re-validation "closes its cycle without adding a line".
- Exactly one successful `gate record --round validation/1`, one `round-1` room, fix marker present,
  entity nonterminal at exactly one prepared open validation gate.
- The re-review reaches a worker whose stage is validation and NEVER the fix producer
  (`SKILL.md` step 4; registry required outcome "independently checked").

## What this journey certifies, and what it does not

Stated plainly, because the revise round found it rather than validation finding it later:

- **Certified on Codex:** the full reuse path — the correction and the re-review both routed to the
  live worker handles through `followup_task`.
- **Certified on Claude:** the HEADLESS shape only — team-mode dispatch with the contract's
  fail-safe fresh branch, because the reuse probe cannot succeed under `-p` (correction 1). This is
  a real, contract-conforming shape and worth grading, but it is NOT the reuse path.
- **Not certified anywhere:** Claude interactive reuse. It works — the probe returns `reuse_ok` on
  the interactive layout — but no live lane runs interactively, so nothing here proves it.

The grading is branch-keyed on the probe result rather than hard-coded to fresh dispatch, so if the
blockers below are ever fixed, the Claude leg starts requiring the reuse chain without another
grading change. That is the whole reason to key the branch on an observable instead of pinning the
shape this environment happens to produce.

**Spun off, not absorbed here (product code, out of this entity's declared scope).** The three
blockers are a user-facing defect, not just harness plumbing: any operator who sets
`CLAUDE_CONFIG_DIR` — a supported Claude Code setting — silently loses worker reuse, standing-teammate
injection, and roster reconcile, because every `~/.claude` reader in `internal/claudeteam`
hardcodes `$HOME/.claude`. The failure is silent by design (reuse condition 0 is fail-safe), so it
degrades to fresh dispatch and never surfaces. Recommend a separate entity covering all three
blockers together; this one neither fixes them nor depends on them being fixed.

## Spike: the riskiest unverified mechanism, run first

One live run per runtime on the composed tree (571017df3) with the team-mode prompt and a persisted
native rollout. Distilled, committed evidence: `_evidence/zq-team-mode-spike/` (topology traces,
durable heading shapes, graded verdicts, and how to reproduce). The raw streams stayed in the job
scratch dir and are not durable.

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

1. **The reuse path is not reachable on Claude headless** — but not for the reason the first
   ideation round gave, and not because reuse is broken generally. Verified at the revise gate
   against source and three controlled runs (`_evidence/zq-team-mode-spike/context-budget-probe-verification.txt`).
   `internal/dispatch/dispatch.go:76` passes `os.Getenv("HOME")` into
   `claudeteam.ContextBudget`, which globs `$HOME/.claude/{projects,teams}` and nothing else; no
   product code reads `CLAUDE_CONFIG_DIR` at all. Three blockers stand between a headless `-p` run
   and a successful probe, each sufficient alone and hit in this order:

   1. **Config-dir blindness.** The live runner gives each scenario its own
      `CLAUDE_CONFIG_DIR` (`filepath.Join(base, scenario.name)`, landed 2026-06-13 in #358's
      parallel fan-out) and CI sets an archivable one
      (`runtime-live-e2e.yml:263-264`), so Claude Code writes its state where the probe never
      looks. Proven by an A/B on one binary and one name: identical fixture state under
      `$HOME/.claude` returns `reuse_ok: true`, and under `$HOME/.claude/rejection-flow` returns
      the exact `no subagent jsonl found for '<name>'` the spike hit.
   2. **`agentType` mismatch.** `scanSubagentMeta` matches `meta.agentType == name`. A headless
      run records `agentType` as the *subagent type* and puts the name in a separate field —
      measured directly: a named background dispatch wrote
      `{"agentType":"general-purpose","name":"probe-worker-x"}`. Interactive sessions record
      `agentType` as the member NAME (14467 meta files on this machine, e.g.
      `spacedock-ensign-p14yaypeav-validation`), which is exactly why it works there.
   3. **No team config.** Headless writes no `~/.claude/teams/` at all, so `lookupModel` fails and
      the probe exits 1 even once 1 and 2 are fixed. `claude-fo-dispatch.md:119` already records
      this for reconcile.

   So Claude worker reuse is real and works interactively; it has never worked in this harness or
   in CI. The reuse assert `assertClaudeReviewerReuse` was on this journey until 2026-06-13
   (#365, `fe4261be8`) and was removed for a DIFFERENT reason, recorded in
   `_archive/lazy-teamcreate-shallow-boot/index.md:370`: "the Claude `-p` run can never satisfy
   [it] — bare mode hard-fails reuse-condition-1". That reason is now obsolete — `-p` is
   team-capable and condition 1 passes — but condition 0 blocks reuse in its place, so restoring
   that assert today would still false-red a conforming run.
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

## Stage Report: ideation (cycle 2)

Scoped to the gate's one annotation: "is this true? this used to work."

- DONE: Verify the probe-failure claim beyond the single spike observation, from source.
  `internal/dispatch/dispatch.go:76` passes `os.Getenv("HOME")`; `internal/claudeteam/contextbudget.go` globs `$HOME/.claude/{projects,teams}`; no non-test code in `internal/` or `cmd/` reads `CLAUDE_CONFIG_DIR`. A/B on one binary and one name: identical state under `$HOME/.claude` → `reuse_ok: true` exit 0; under `$HOME/.claude/rejection-flow` → the exact spike error, exit 1.
- DONE: Answer the history question — did Claude worker reuse ever work, and in what mode.
  Yes, interactively, and it still does. `agentType` carries the FO-assigned member name in real sessions (14467 meta files, e.g. `spacedock-ensign-p14yaypeav-validation`) and `~/.claude/teams/*/config.json` exists. `assertClaudeReviewerReuse` was on this journey until 2026-06-13 (#365, `fe4261be8`), removed because `-p` was BARE then — condition 1, per `_archive/lazy-teamcreate-shallow-boot/index.md:370`. That reason is obsolete; condition 0 blocks it now.
- DONE: Determine interactive vs headless, and state the limitation in the entity.
  Two controlled headless runs settled what `-p` writes: `{"agentType":"general-purpose","name":"probe-worker-x"}` and no `teams/` directory. So beyond the config-dir mismatch there are two further blockers — `scanSubagentMeta` matches `agentType`, not `name`, and `lookupModel` needs a team config headless has never written. New section `## What this journey certifies, and what it does not` states the limitation.
- DONE: Update the entity body and record what changed.
  Determined-shape table row corrected; correction 1 rewritten with the three blockers, the A/B, and the history; certification-scope section added; the three blockers spun off as a recommended separate entity rather than absorbed.

### Summary

The claim's CONCLUSION survives — the probe cannot succeed on headless `-p`, so the Claude branch of
the determined shape stands — but its stated CAUSE was wrong, and the captain was right that this
used to work. It works interactively today; it has never worked in this harness or in CI. What I
called a headless property is three separate defects: the probe ignores `CLAUDE_CONFIG_DIR` while
the harness and CI both redirect it, headless records `agentType` as the subagent type rather than
the worker name, and headless writes no team config for the model lookup. None is in this entity's
scope to fix, and the design does not depend on fixing them, because the grading branch is keyed on
the observed probe result rather than on the shape this environment happens to produce. The one
design change is honesty about coverage: this journey certifies the Codex reuse path and the Claude
headless fail-safe path, and certifies Claude interactive reuse not at all.

## Stage Report: implementation

- DONE: Execute the gated design exactly — all nine mechanisms, the registry doc diff, the certification-scope section stands as shipped truth.
  `9cb13c7c7`. Team-mode cue (host-neutral `rejectionTeamModeCue` + per-host `rejectionHostRealization`, the `merged_team_mode_live_test.go:100-107` split); reuse condition 5 verbatim as drafted; branch-keyed `assertRejectionWorkerTopology`; codex native-rollout extractor on task paths only; per-run digest; fixture `(cycle 2)` heading; host-neutral cycle-line and gate-prepared asserts; tautology removal anchored to `finalMessage`; the finding-10 repairs.
- DONE: Grading stays branch-keyed; the three `CLAUDE_CONFIG_DIR` blockers untouched.
  No file under `internal/claudeteam` or `internal/dispatch` is in the diff. The Claude branch key is read from the probe's own tool_result, so a blocker fix flips the leg to the reuse chain with no grading change: `TestRejectionTopologyProbeFailSafe` proves both directions today, and would fail if the branch were pinned to `fresh`.
- DONE: Offline first — AC-2, AC-3, AC-4 by table test.
  `TestRejectionTopologyExtractsCapturedCodexChain` asserts the 8-event determined chain comes out of REAL captured rollout bytes; it fails if the extractor stops normalizing `/root/`-rooted against bare handles, counts a refused spawn's error-string output as a dispatch, or counts a worker's non-final `agent_message` as a completion. `TestRejectionTopologyRedsNonConformingShapes` (6 cases) fails if any of the self-review chain, fresh-rework, re-review-before-rework, reviewer-collapsed-onto-producer, missing-re-review, or two-validators-up-front shapes stops producing a `rejection-worker-topology` red. `TestRejectionTopologyBranchesAreMutuallyExclusive` fails if either branch's chain also passes under the other's. `TestRejectionCycleLineAndGateShapes` fails if a second Cycle line, a missing one, a reworded one, or an unprepared gate stops redding. `TestAssertRejectionFlow` fails if the exact-duplicate-heading entity passes, or if a final message naming only one outcome passes.
- DONE: `go test ./...` plain and `-race`; contractlint reconciliation.
  20 packages ok plain, 20 ok under `-race`, no data races, no failures. `go test -count=1 ./internal/contractlint/` ok — the registry entry still binds the journey/fixture ids after the required-outcome edit.
- DONE: Live loop per AC-1 — 3 consecutive conforming greens per runtime, ledger row per run.
  Codex 3/3 runs; claude-sonnet 3 consecutive in 5 runs; both inside the 8-run budget. A conforming green required exit 0 AND a persisted digest carrying that run's branch chain in order. Every codex digest is the reuse chain, every claude digest the fail-safe fresh chain. Evidence: `_evidence/zq-team-mode-ac1/` (ledgers, 8 per-run digests, the one red's test log, the driver).
- DONE: FO adherence residual rule — one hardening round at a quoted instruction.
  Mode 1 (ends without `gate prepare`) recurred on the first codex run. Round 1 was spent at the instruction the stream shows the FO following: the prompt's own goal line ("stop after … the second validation passes; leave the entity nonterminal"), which contradicted `feedback-rejection-flow/SKILL.md` step 5 ("Done when exactly one fresh open gate has been prepared and presented") that the same run had loaded. The goal line now names the prepared gate as the bounded stop. Mode 1 did not recur in the 8 following runs. The grade was not loosened.
- SKIPPED: Surface within 11 files / ~406 changed lines ±30%.
  Declared over tolerance and FO-ruled to ride to the gate as declared, with the no-loosening rule outranking the estimate; not trimmed. The measure: 12 files (within the ±1 file tolerance) but **1115 changed lines against a 528 ceiling — +174%**; captured data 13 lines against ~120, under. Largest overruns: the topology table test 384 vs ~110, `shared_reviewer_reuse_test.go` +245 vs ~90, `shared_assertions_impl_test.go` +58 vs ~20. Cause: the estimate was made before the real rollout format was known, and three real-byte findings (a follow-up names its worker under `target` not `task_name`; a refused spawn returns a plain error string, not JSON; a running worker emits non-final `agent_message`s) each needed extractor code, a guard, and a regression case. AC coverage was not narrowed to fit. The captain rules at the gate.
- SKIPPED: Retire the `AUDIT(2026-08-16)` block above `TestLiveCommonRejectionFlow`.
  FO ruled LEAVE UNTOUCHED: the block is absent from this stack's lineage, so there is nothing in the candidate to delete. It lives on `main` (`df0bd50d9`), which is not an ancestor of the composed tree `571017df3` this layer was cut from (they diverge at `0c6a2c32a`). Superseding it belongs to the stack-merge conflict resolution on `shared_live_runner_test.go`, which the FO owns at the merge window — stack's version wins, the rejection-flow block drops as superseded by this layer, the other journeys' blocks are preserved.
- DONE: Signal the FO for the base, rebase, push, stack-link, GraphQL read-back.
  Layer 6 of stack #720. Rebased onto `spacedock-ensign/red-auto-continue-gate-bypass` (#723) after the FO's base reassignment was overtaken by #723 opening on the same base; clean rebase, no conflicts, both layers' asserts green together. Pushed via SSH, PR #724 open against that branch, linked with `gh stack link 718 719 721 722 723 724`. Verified by GraphQL `pullRequest.stack` read-back from four PRs — stack #720, size 6, `main → 718 → 719 → 721 → 722 → 723 → 724` — never the banner or `gh stack view`. One final mechanical rebase onto the FO's restacked head remains, pre-verified clean on a scratch branch.
- DONE: File the implementation stage report, commit path-scoped, push, signal the FO, stop for validation.
  This section; no gate preparation or resolution.

### Summary

The journey now runs in team mode on both runtimes and grades the ordered worker
chain its branch owes, and AC-1 is met: 3 consecutive conforming greens on codex (3
runs) and on claude-sonnet (5 runs), against a measured baseline of 0. The load-bearing
correction to the design came from checking the extractor against real multi-agent
rollout bytes before trusting it: a `followup_task` names its worker under `target`,
not the `task_name` key the spawn uses, so the parser the design described would have
silently dropped every follow-up and reported a conforming reuse chain as two events
short. Two further real shapes — a refused spawn returning an error string instead of
JSON, and non-final worker messages before FINAL_ANSWER — would each have invented or
doubled an event.

Two things need the captain rather than me. The surface is +174% over the ideation
estimate, declared above with its cause; I did not trim AC coverage to reach the
number. And the entity's own determined-shape chain orders `gate record --round` BEFORE
the rework routing while citing `feedback-rejection-flow/SKILL.md` steps 1-5, but the
skill orders the correction first — step 1 is done only when the target worker's
entries close the review log, and step 2 only when the recorder counts every entry the
round accumulated. Claude run 2 followed the entity's order, recorded at `entries=2`,
and was correctly redded by the pre-existing `entries=4` recognizer. The entity text is
wrong on that one ordering; nothing in this layer depends on it, because the topology
grader orders only worker routing events and never places `gate record` among them. I
left the recognizer alone rather than relax it.

One deviation inside the design: item 6 names only the rework bullet, but the
determined shape lists four distinct report sections including
`## Stage Report: validation (cycle 2)`, so the validation bullet got the same
cycle-aware instruction. Without it the two cycle-2 sections are not distinct and the
duplicate-heading shape AC-2 requires to red would not. I also deleted
`assertClaudeSingleEntityRejectionFlow` and its table test: its only production caller
was the line this layer replaces, and it accepted both a reuse and a fresh topology —
the multiple-path acceptance the captain ruled out.

### Amendment: FO dispositions on the three carried items

Recorded after the implementation report was filed; no candidate code changed.

- DONE: Correct the determined-shape chain to the skill's actual order (FO-authorized,
  implementation-discovered-truth class; captain ratifies at the gate).
  `## Determined shape` now orders the correction and its review-log closure BEFORE the single
  round recording, quoting both `SKILL.md` done-conditions, and records that Claude run 2
  followed the old order, recorded at `entries=2`, and was correctly redded by the untouched
  `entries=4` recognizer. This was the FAILED-adjacent finding raised in the summary above; it
  is now closed in the entity text rather than left as a note against it.
- SKIPPED: Retire the `AUDIT(2026-08-16)` block — FO ruled LEAVE UNTOUCHED.
  Superseding the rejection-flow block belongs to the stack-merge conflict resolution on
  `shared_live_runner_test.go`, which the FO owns at the merge window (stack's version wins, the
  rejection-flow block drops as superseded by this layer, the other journeys' blocks are
  preserved). Reclassifies the FAILED item above: not my lineage, not mine to delete.
- SKIPPED: Bring the surface back inside the declared ceiling — FO ruled it rides to the gate as
  declared, with the no-loosening rule outranking the estimate. Unchanged at 12 files / 1115
  changed lines; the captain rules at the gate.

## Stage Report: validation

- DONE: Source of truth is the entity's stage report and amendment sections; validate on the tip 52b94f9aa (PR #724; this layer's diff is one commit above #723's 134e494a8).
  Verified: worktree tip 52b94f9aa = PR #724 head (base = #723's branch), exactly one commit above 134e494a8; layer diff 12 files, 986 insertions / 151 deletions.
- DONE: Verify the nine gated mechanisms on the tip bytes, spot-verifying by the falsifying edits the report claims.
  All nine shipped as gated: team-mode cue split, reuse condition 5 byte-equal to the gated diff, branch-keyed ordered chain, native-rollout extractor on task paths, per-run digest (written pass or fail, before grading), `(cycle 2)` headings for both stages, host-neutral cycle-line and gate-prepared asserts, finalMessage-only outcome grading with exact per-section counts, finding-10 repairs. AC-2: all six non-conforming topologies red as gradedErr `rejection-worker-topology` with the named detail, plus cycle-line/gate/duplicate-heading shapes. AC-3: branches mutually exclusive. AC-4: all four probe/routing cases. Spot mutations on a throwaway checkout: refused-spawn counting, followup-`target` drop, and independence-clause removal are each caught; two claimed falsifying edits are NOT caught, and the stageToken narrowing is unpinned (findings 1-2).
- DONE: AC-1: verify the ledgers and per-run topology digests in _evidence/zq-team-mode-ac1/ (codex 3/3 in 3, sonnet 3-consec in 5, baseline 0), and that a persisted digest accompanies every counted green.
  Committed at 08242fe75; ledger rows match all 8 digests (codex all reuse chains, claude all fresh chains); the claude run-2 red also carries its digest, and its kept test log shows a real premature-recording red (`entries=2` against the untouched `entries=4` recognizer).
- DONE: Verify the semantic conflict resolution in claude_runtime_helpers_test.go: layer 5's dispatch-pointer widening AND this layer's token-anchor narrowing both live in the composed function.
  Both edits are present in the composed `claudeSpawnIsForStage`; layer 5's TestAutoContinueRevalidateStreamCountsBothValidators and every new topology test pass in one focused run. Caveat: the checklist's "stage-token negative cases" do not exist anywhere in the package — see finding 2.
- DONE: Assess the declared deviations for the gate: the +174% surface overrun, the FO-authorized entity ordering correction, reuse condition 5's all-workflow semantic, and the certification-scope statement.
  Overrun accounting accurate within ~1% (measured 1137 changed lines incl. 13 data vs declared 1115+13; both ≈ +175% over the 528 ceiling); the three real-byte causes are visible in the committed fixture bytes. The ordering correction quotes `feedback-rejection-flow/SKILL.md` steps 1-2 done-conditions verbatim and matches them — correction first, one recording counting all four entries. Condition 5 shipped verbatim; for this repo's own dev workflow it is a no-op (validation already declares `fresh: true` + `feedback-to`). Certification scope is consistent with the shipped branch-keyed grading; TestRejectionTopologyProbeFailSafe proves both directions of the key.
- DONE: Reference tip CI run 31986034831's rejection-flow lanes when concluded; do not run new live loops.
  Concluded. Offline: green on the tip SHA. Codex lane: rejection-flow GREEN (the lane's one failure is TestLiveCommonACValueReanchor, a different journey). Claude lane: rejection-flow RED — the digest shows a fully conforming 8-event fresh chain and a prepared gate, but the FO never recorded the round; the new graders named it (`rejection-cycle-line`: 0 Cycle entries want 1; `rejection-round-missing`). A real adherence deviation correctly caught, not a false red (finding 3).
- DONE: Validation stage report: per-AC verdict, PASSED or REJECTED recommendation, path-scoped commit, push, signal FO, stop. No gate.
  This section; committed path-scoped to the state checkout and pushed.

### Per-AC verdicts

- AC-1 PASSED. AC-2 PASSED. AC-4 PASSED. AC-3 PASSED with an evidence nuance: only the codex chain is captured rollout bytes; the claude fresh branch is constructed routes plus synthetic dialect streams plus the five live digests, not a committed captured-stream replay ("two captured topology fixtures" as written is half-met).
- Checks run: focused offline suite green in one run; CI offline job green on the tip; `go test -race ./internal/ensigncycle` green at tip (clean solo run, exit 0); gofmt clean; contractlint green inside the suite. Two earlier local full-package timeouts were validator-side parallel-load artifacts — a clean solo rerun and CI offline both pass.

### Findings (recommendations for `## Review-finding disposition`; classification authority is the FO's)

1. Evidence defect, deferred risk + polish. Two of the implementation report's three claimed falsifying edits for TestRejectionTopologyExtractsCapturedCodexChain do not falsify: `codexBareHandle` reduced to identity passes every topology test (all handles in the committed fixture are `/root/`-rooted — the test comment's "one side rooted, the other bare … would report zero reuse" is wrong against its own bytes), and neutering the FINAL_ANSWER/`Done:` completion filter passes (the awaiting-debt map alone dedupes; the fixture's intermediate message payload is truncated before any `Done:`). The refused-spawn and followup-`target` claims do falsify as reported. Promote-to-material: a mixed-spelling or `Done:`-carrying-intermediate rollout entering the fixture set without a strengthened guard.
2. Evidence defect, deferred risk. The finding-10 stageToken narrowing is unpinned: reverting all three call sites to bare substring leaves the entire ensigncycle package green (clean solo run, exit 0). The false-positive direction the repair exists for has no regression case. Promote: any fixture stream where substring over-match flips a verdict.
3. FO adherence recurrence, for the entity's residual decision rule — not a candidate defect. Tip CI claude-live red as above: a round-recording omission mode (gate prepared, round never recorded), distinct from mode 1 (no prepare) and from the loop's record-before-correction red. Per the entity's rule it earns one quoted-instruction hardening round or escalate-and-own; the run's stream is in CI run 31986034831's claude-live artifact.
4. Polish. The AC-1 README's "Mode 3 (records the round before the correction lands)" collides with the entity's "Mode 3 (interrupt_agent teardown)" numbering.

### Summary

Recommend PASSED. All four ACs have valid, reproduced evidence; the nine gated mechanisms are on the tip bytes; and the strict grading proved itself in both directions — offline against every non-conforming shape in hand, and live in tip CI, where codex ran the reuse chain green and claude was red-flagged with named codes for a real FO round-recording omission. No material finding against the candidate: findings 1-2 are test-strength deferred risks with recorded promote conditions, finding 3 is residual-rate data the entity's own decision rule routes, and the captain rules at the gate on the declared +174% overrun and ratifies the FO-authorized ordering correction.

### Amendment: tip CI claude-live red re-attributed (validation, post-report)

Routed by the FO after my report; verified independently from the run's own artifacts
(stream + digest + final entity reads), not from the graded messages or the message thread.

**Retraction.** My finding 3 ("FO round-recording omission") is WRONG and withdrawn. The run's
bytes show: correction landed first, Cycle line appended byte-identical to the fixture's
`feedback-cycle.txt` (stream line 464), `gate record --round validation/1` succeeded with
`entries=4` (line 473), the re-review was a fresh validation worker (conforming 8-event fresh
digest), all four distinct report sections present, and the FO stopped at one presented open
gate. I graded the FO from the assert messages alone — the exact trap this checklist's
source-of-truth clause names.

**Finding 5 — material (proposed): the composed grading reds the sanctioned withdraw→re-prepare
recovery.** Sequence verified: premature `gate prepare` attempt-1 right after the round-1
REJECTED report (line 301, before the rejection skill loaded — the step-5-before-step-4 trap),
then the full correct flow, then `gate withdraw` of the stale attempt-1 with a reasoned message
(line 589; launcher reports `readiness: withdrawn-awaiting-prepare` — a modeled state), then
`gate prepare` attempt-2, state=open (line 596). End state satisfies this entity's own determined
shape ("nonterminal at exactly one prepared open validation gate") and the registry's required
outcome (rejected authority withdrawn, never presented). The red comes from main-side
`assertRejectionRoundGateBoundary` (c0e3a14e3, on origin/main, untouched by this layer), which
pins THREE things the recovery cannot satisfy: `len(Attempts) == 1`, attempt ID
`gate-attempt:rejection-task-validation-1`, and briefing `…attempt-1:revision-1`. Expected shape
per the sanctioned-recovery vocabulary (`withdrawn-gate-recovery` journey in this same suite;
layer 5's accuse-only-a-resolved-gate rule): exactly one OPEN attempt, withdrawn siblings
tolerated, briefing bound to the open attempt. Classification proposal: outcome defect in the
composed journey grading, Material via contract[docs/runtime-live-ci-registry.md § rejection-flow
required outcome] with the observed tip-CI trigger; narrow fix as stated; ownership (main-side
bytes vs this entity's grading scope) and the tolerance question (may a transient, self-recovered
premature prepare grade red at all?) are FO/captain dispositions. New FO residual mode to record:
premature round-1 prepare, self-recovered via withdraw, never presented.

**Finding 6 — the cycle-line red's actual cause (this layer's assert, working as declared).** The
FO wrote the byte-exact Cycle line under `## Feedback Cycles` (H2); the workflow README declares
the projection under `### Feedback Cycles` (H3), which is what `feedbackCyclesSection` keys on.
Recommend keeping the strict grade (the declared grammar is the contract; no loosening) and
treating heading-level drift as a residual mode hardened at the quoted declaration; polish: when
the count is 0 but the line exists under a level-mismatched heading, the diagnostic should say so
instead of implying the round was never recorded. Level tolerance would be a captain-owned
loosening.

**Amended recommendation: REJECTED (was PASSED).** Every written AC keeps its verdict and the
nine gated mechanisms stand verified; the rejection is the stage definition's over-specification
clause applied to the composed grading — the tip run satisfies the registry outcome and the
entity's determined shape yet grades red on three attempt-1 pins. If the captain rules that a
transient premature prepare itself grades red (tolerance, theirs alone), finding 5 becomes a true
red caught for the wrong reason with the wrong message, and this recommendation reverts to PASSED
with findings 5/6 reclassified as caught FO deviations plus a diagnostics polish item.

### Amendment 2: FO-directed correction — product-bug hypothesis tested and refuted (validation)

The FO routed two re-attributed items to replace retracted finding 3, and asked for a root cause
on item 2 (a product verb possibly eating committed body bytes) plus a per-AC re-verification.

**Full byte citation for the retraction.** Stream indices: 464 Cycle-line Edit (byte-exact) →
472 `gate record --round validation/1` SUCCESS `entries=4` → 475 `state commit` → 532 second
Edit KEEPS the line and appends the cycle-2 validation report after it → 536-540 `git add` +
commit → 560 heading listing shows the line durable in the entity → 589 sanctioned withdraw →
593 `state commit` → 595 re-prepare attempt-2 (open) → 598 `state commit`. Every write committed.

**Item 2 root-caused: NO product verb eats body bytes; the hypothesis is refuted.** Offline
reproduction with the tip-built binary on a faithful fixture replica, running the FO's exact
sequence — `gate prepare` (attempt-1), append `## Feedback Cycles` + the byte-exact line,
`gate record --round validation/1` (reproduced `entries=4`), append the cycle-2 report,
`gate withdraw`, `state commit`, `gate prepare` (attempt-2, open) — leaves the entity body
IDENTICAL through every verb (diff-verified after each), with the H2 section and line intact at
the end and the frontmatter carrying exactly the run's recovered two-attempt state. The graded 0
has one cause: `feedbackCyclesSection` anchors on the literal `### Feedback Cycles`
(shared_assertions_impl_test.go:56); the FO wrote `## Feedback Cycles` (H2), so the finder
returns empty. Not a wrong-rev read (the assert reads the working file at test end), not
serialization loss. Finding 6 stands as the complete explanation; no Material product finding.

**Item 1 attribution, held to the bytes.** The attempt-cardinality red is main-side:
`assertRejectionRoundGateBoundary` pins `len(Attempts)==1`, the attempt-1 ID, and the attempt-1
briefing (shared_round_recording_test.go:246-255), introduced by c0e3a14e3 — on origin/main, not
in this stack — and it ran host-neutrally on this lane before this layer. The Material-candidate
classification and the fix shape stand as filed in finding 5 (select the OPEN attempt; withdrawn
attempts are not violations; bind the open attempt's briefing); only the "this layer's assert"
attribution is corrected.

**Per-AC re-verification: no verdict flips.** AC-1 PASSED unchanged. AC-2 PASSED as written —
all six shapes still red with named codes; the tip run is the CONVERSE direction (a conforming
run redding), which no AC-2 case covers and which finding 5 carries. AC-3 PASSED on its fixtures;
its generalized promise ("a conforming run grades green") has the one observed live
counterexample, produced by the main-side boundary assert — this layer's branch-keyed topology
grading graded the same run GREEN on a conforming digest. AC-4 PASSED unchanged. The REJECTED
recommendation stands on finding 5 (Material, narrow fix named), with finding 6 the strict-grade
residual mode and the captain owning both tolerance rulings.

### Amendment 3: tip CI attempt-2 codex red — Cycle-line target pinned in the fixture (FO-recorded revise item)

**Attribution (run 31991864922 attempt-2, codex lane, from stream bytes not doc text).** The
codex FO drove the repaired journey correctly: team-mode reuse topology exact (spawn
implementation → spawn validation → reuse implementation → reuse validation), `gate record
--round validation/1` exit 0 with entries=4, single publication, gate prepared and left open.
The one miss: it appended the authorized Cycle line to the workflow README.md (a file_change on
README.md plus `git commit -m 'feedback: rejection-task record cycle 1'`) instead of the
entity's `### Feedback Cycles` section. That single write produced BOTH observed codes:
`rejection-cycle-line` (entity holds 0 entries) and `rejection-round-missing` (README.md moved
off its pinned bytes — the round oracle's boundary check reports under that code; the
label-honesty nit is on the captain's filing list).

**Root cause: fixture underdetermination.** The rejection README said "appends one Cycle line
to this section" and never named the file; read from the README, "this section" resolves to the
README's own section. The sibling escalation fixture names its target file explicitly. On the
revised text codex resolved it entity-side once (attempt 1, green) and README-side once
(attempt 2, red) — a nondeterministic reading of an unpinned input, the audit's disease.

**Captain ruling 2026-08-17: "fold and rerun".** The pin folds into this entity's open revise
window; ratification rides validation gate attempt-2.

**Correction landed: 9b1c86a60** (parent 43f7537b0; one file, 2 lines in `rejectionReadme()`):
the first sentence names the `### Feedback Cycles` section in `rejection-task/index.md`; the
next paragraph reads "that section's exact heading"; the H3 strict/pin grammar is byte-identical.
The original implementer was no longer addressable, so the correction ran the contract's
fresh-dispatch path. Worker proof: `go test ./...` and `go test ./... -race` green from the
worktree, `gofmt` clean, and a worktree-wide grep shows no surviving copy of the
self-referential text. FO read back the pushed diff against the authorized shape — exact.
Declared limit (worker-flagged, FO-concurred): offline runs cannot reproduce the codex FO's
Cycle-line placement; only a live codex rejection-flow run proves the pin changes behavior —
that proof is the new tip CI run's codex lane.

### Amendment 4: delta re-verification 52b94f9aa..9b1c86a60 — REJECTED lifted, recommendation PASSED (validation)

Re-verified under the captain's revise ratification condition, on the FO's delta window. Three
commits: 3f6f3e935 (withdraw tolerance + ac-reanchor identity), 43f7537b0 (heading-drift
diagnostic), 9b1c86a60 (fixture Cycle-line target pin).

- Finding 5 SATISFIED. `assertRejectionRoundGateBoundary` now selects the ONE open attempt,
  tolerates withdrawn siblings, still reds a resolved/applied attempt, reds zero open attempts,
  and binds attempt-N to briefing attempt-N with N-agreement — exactly the fix shape filed. The
  regression drives REAL gate verbs (withdraw → re-prepare to an open attempt-2) and the oracle
  grades the recovered state green; falsifying probe verified on a throwaway checkout — restoring
  the literal attempt-1 pin reds the recovery regression. Live proof both directions: claude lane
  17/17 on 43f7537b0 (run 31991864922, the lane my finding 5 red came from), codex 17/17 on
  9b1c86a60 (run 31996696789 — first fully green codex lane, rejection-flow included).
- Finding 6 SATISFIED as ruled. Captain ruled STRICT: the H2 byte-exact line still reds
  (negative case added, byte-exact under `## Feedback Cycles`), the diagnostic now names the
  drift ("byte-exact … the round WAS recorded") instead of the misleading bare count — probe
  verified: removing the diagnostic branch reds the phrase assertions. The fixture pins the
  target file and the exact H3 level (9b1c86a60), the quoted-instruction hardening the residual
  rule requires; Amendment 3 attributes the codex attempt-2 red to that underdetermination and
  the pinned tip's green codex lane proves the pin live.
- ac-reanchor identity fix VERIFIED: repo-local `git config` identity replaces the per-invocation
  `-c` flags, with a read-back guard the old form cannot pass
  (TestACReanchorFixtureRepoCommitsFromAnyProcess); ac-value-reanchor is green inside the 17/17
  codex lane on the tip.
- Claude rejection red on 9b1c86a60 VERIFIED as adherence variance, not a branch defect: the
  run's own stream shows `gate record --round validation/1` at index 288 returning `entries=2`
  (289) BEFORE the rework spawn at 310 — the early-record mode already documented in this
  entity's corrected determined shape and redded by the deliberately-strict `entries=4`
  recognizer. The auto-continue red is a true positive outside this surface. Both have filed
  owners: own-claude-early-rejection-round-record.md (zf7rymtke3b6xp7r0337hjj4, which also
  carries the "never invoked" label-honesty defect this red's message exhibits) and
  own-sonnet-gate-conn-bypass.md (rcpa3nnkmgy9tm9hand0jkf6). Note the FO's window message had
  those two ids swapped; the files are authoritative.
- Checks run at 9b1c86a60: focused rejection + recovery + identity tests green; focused `-race`
  on both changed packages green; gofmt clean; offline CI green on both delta runs.

**Recommendation: PASSED** (replaces Amendment 1's REJECTED). Finding 5's material defect is
fixed and regression-pinned, finding 6 is closed by captain ruling with the diagnostic and pin
landed, findings 1-2 remain deferred risks as filed, finding 4 (mode-numbering collision)
remains polish. Composite green context: each lane is 17/17 on a delta commit containing all
assert fixes; the two tip reds are owned adherence-variance entities, not this branch. The
captain ratifies at gate attempt-2 with the +174% overrun and ordering correction still riding.
