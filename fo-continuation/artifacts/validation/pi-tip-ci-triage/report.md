# Pi tip CI triage — proposal, not fix authorization

Inspected stack tip `ac3c06a445fcfbc773639254ff7242965248ba78` and retained run `35058669297`. No code changes, tests, native runs, CI, pushes, or rebases. The approved fo-continuation report, gate, and frontmatter remain unchanged.

## Conclusion

The two `completed=-1` failures miss a supported Pi completion notification. Attribution is available without trusting completion prose: parent dispatch call/result identity joins a child session's exact run-name metadata, assignment pointer, entity/cwd, committed report, and successful terminal event. Parent notifications occur before advancement/preparation, and parents read the resulting reports before those boundaries.

The default-headless-gate-stop journey has a grader evidence defect. Auto-continue additionally has a real downstream outcome failure: its only gate preparation used a misspelled reference path and failed, leaving no prepared approval gate. Correcting the completion grader must expose that remaining failure, not turn the entire auto-continue run green.

## Exact evidence and ordering

Line numbers below are 1-based in each original session, not concatenated grader offsets. Original files, SHA-256 hashes and sizes are in `source-manifest.json`; exact selected native records with provenance are in `retained-events.json`. Full original parent/child JSONL remains at the manifest paths.

| Journey / stage | Parent dispatch and result | Retained child identity / durable result | Completion and next boundary |
| --- | --- | --- | --- |
| default-headless / implementation | Parent lines 39/40, matching toolCallId, `details.runId=c7f1d5e7-ef00-4235-a860-790fda66500a` | Child directory `22943409-f624-4946-bdb2-7aff5c3f9f6a/run-0`; line 4 `session_info.name=subagent-worker-c7f1d5e7-ef00-4235-a860-790fda66500a-1`; line 5 exact implementation pointer; lines 33/34 path-scoped commit `43c5ba0`; lines 36–38 report count/state/staging checks | Child terminal `stop` line 39 at 05:26:52.921Z; parent notification line 43 at 05:26:53.116Z names that child; parent report reads lines 45/47; validation mutation line 48 at 05:27:26.242Z |
| default-headless / validation | Parent lines 52/53, run `fbe89c19-4218-4448-abb0-a9ebf10ebb6c` | Child directory `cd9a7854-9e56-497f-b9f8-6a61d46dbc63/run-0`; line 4 exact run-name metadata; validation assignment; line 76 committed report `305bd81` and exact entity content | Terminal `stop` line 78 at 05:30:15.291Z; notification line 56 at 05:30:15.487Z; parent report verification follows; prepare line 73 succeeds, command log retains state commit `9bf0901f8e4dd2c509c83806392f06f21f113f9d` |
| auto-continue single-root / validation | Parent lines 37/38, run `f4560401-c862-462b-9594-522cbc325a40` | Child directory `4c1a88ae-3f7b-4ce8-8539-c922c068d407/run-0`; line 4 exact run-name metadata; line 5 exact validation pointer; lines 28/29 commit `120a280` on worktree branch, lines 31/32 clean status and report | Terminal `stop` line 33 at 05:29:37.124Z; notification line 41 at 05:30:01.460Z; parent exact report read line 48; gate prepare line 101 at 05:34:38.625Z |

Child directory UUIDs are **not** run IDs. The exact `session_info` run-name join is necessary; guessing from the directory name or the generic `worker` label is unsound. Parent `Session file:` locators point inside their own retained session tree. Each child cwd and pointer read agrees with its parent assignment, and the child report/commit targets the dispatched entity (the auto-continue worktree copy, not its stale base copy).

No direct Git repositories/bundles for these Linux temp roots are present in the downloaded artifact set. Durable claims here are supported by retained native commit/status/show results plus existing CI assertions, not a newly executed local `git show`. The default journey ran its independent gate-state/log graders and reported only the implementation lifecycle failure. Auto-continue's gate-state checks are short-circuited by `assertWorkerLifecycle`, so those checks have not yet delivered their downstream failure classification.

## Finding G1 — missed supported notification

- Released user and normal workflow: Pi first officers using the shipped async `pi-subagents` completion wake in normal shared CI journeys.
- Observable harm: correctly dispatched, committed and completed workers are reported as not dispatched; live compatibility evidence fails for the wrong reason and masks later outcome checks.
- Authority: `contract[skills/first-officer/references/pi-first-officer-runtime.md#runtime-implementation]` declares the native background completion notice as the primary completion signal and requires report verification before advancement.
- Trigger evidence: the parent/child chains above, together with `claude_runtime_helpers_test.go:185–263`, whose Pi branches recognize only status/`subagent_wait` result text and never `custom_message` / `subagent-notify`.

Proposed classification: **Material evidence defect**. Proposed ownership: existing Pi runtime evidence/lifecycle grader owner in `internal/ensigncycle`, separate from fo-continuation's four-file instruction scope and separate from dispatch naming. Proposed disposition: **FIX in that owner, only after distinct FO authorization**. This proposal does not authorize edits.

### Smallest safe correction surface

Keep the existing lifecycle grader and its durable report/gate checks. Add one narrow Pi notification correlation path using retained parent and child evidence; do not create a controller, fake a status result, concatenate child tool calls into the parent timeline, or accept `Background task completed` alone.

1. At the existing Pi evidence boundary, supply the scenario's known artifact/session root alongside the parent stream, or a small typed map of validated child records keyed by exact notification identity. `liveResult.artifactDir` already exists. `piSharedLiveDriver.run` already retains child sessions but currently folds only the root transcript into `result.stream`; `lifecycleStream` simply returns it. Ensure both existing callers reach the same resolver: auto-continue uses `driver.lifecycleStream`, whereas default-headless uses `nativeLifecycleStream` through the Codex-oriented helper. Do not fix only one caller.
2. In the existing `assertWorkerLifecycle` owner, resolve parent spawn toolCallId → successful result `details.runId`; then recognize the actual native notification at its **parent-stream position**. Map its exact session locator beneath the captured parent-session tree, remapping only the known artifact-root prefix for portable replay. Reject traversal, external/symlink targets, missing children and ambiguous identities.
3. Require the referenced child's `session_info` exact agent/run/epoch identity, exact original task pointer, matching cwd and dispatched entity, and a successful terminal assistant `stop` after the report's successful commit evidence. Preserve `spawn < child terminal <= notification < next boundary`, report verification and existing durable-state assertions. Reject duplicate/conflicting completion evidence rather than letting a later notice overwrite the first completion.
4. Do not require all historical Pi forms to gain a new dependency. Preserve the existing supported status/wait cases and Claude/Codex behavior. No new inference from the notice's commit abbreviation or success phrase.

Likely edit surface is the existing lifecycle helper/test owner, Pi evidence adapter, the existing default-headless native-stream bridge, and one captured replay fixture set. Prefer a small optional typed evidence parameter/loader over a new shared runner abstraction. Exact code shape belongs to implementation; no claim that text-only branching in one function is sufficient.

Reusable pieces: `piSessionRecord`, `piSessionMessage`, `piToolCallBlock`, `piTextContent` and `piToolResults` already parse Pi records and toolCallId joins in `pi_rejection_extractors_test.go`; extend or reuse these types narrowly where appropriate. None currently validates native notification-to-child-run identity. `readPiSessionEvidence` grades model/time/cost, not lifecycle attribution. `onePiSession` only enforces root cardinality. Do not reuse `piAsyncRunID` as the authoritative join: its legacy `Async workflow/run [...]` text does not match this run's `Async: worker [...]`, while structured `details.runId` does. The old dispatch-prefix regex remains the naming owner's separate concern.

### Bounded proof owed by an authorized correction

Use the existing `TestImplementationLifecycleAndObserverNegativeControls` owner and captured replay pattern (`auto_continue_replay_test.go`, `TestPiAutoContinueReplayDoubleDispatch`). Add a compact retained parent-plus-child replay from these exact run bytes, with provenance/hashes. It must credit default implementation completion and auto validation completion without generating synthetic completion events. The auto **whole journey must still fail** for its absent prepared gate.

Negative matrix: unchanged generic completion text with wrong/missing run metadata; unrelated child or different entity/cwd/task pointer; sibling/escaped locator; missing or nonterminal/error child; stale epoch; notification before spawn or after next boundary; duplicate/conflicting notice; missing report/commit. Preserve status/wait positives and existing Claude/Codex controls. These are extensions of the current owned tests, not a new standing lane. No test or native run was executed during this investigation.

## Finding O1 — invalid auto-continue gate reference

- Released user and normal workflow: a Pi first officer continuing a successful validation to its promised unresolved captain gate.
- Observable harm: validation committed, but no approval package was prepared/presented; the run stops before the requested gate outcome. No unauthorized approval or successor dispatch is observed.
- Authority: `contract[skills/fo-gate-lifecycle/SKILL.md#first-officer-gate-lifecycle]` requires retaining resolved paths for one successful preparation before presentation, and requires stopping on a nonzero result.
- Trigger evidence: parent line 101 passes `/tmp/TestLiveCommonAuto-ContinueAfterImplementation131399588/004/README.md`, although successful reads at lines 97/99 use `/tmp/TestLiveCommonAutoContinueAfterImplementation131399588/004/README.md`. The inserted hyphen makes the reference nonexistent; line 102 exits 1, and final line 103 reports the blocker. There is no later successful prepare/commit in the complete parent trace. Process status is `error=<nil> timeout=false`.

Proposed classification: **Material observed outcome failure for this journey**, not a completion-grader defect or an authorization bypass. Proposed ownership: existing auto-continue/gate-path runtime journey owner; not fo-continuation and not naming. Proposed disposition per FO's current direction: retain and re-observe at the repaired tip; do not introduce a new mechanism or product change from this single path-typing occurrence. The fail-closed stop after the typo follows the contract; changing that stop rule would weaken it. A corrected grader alone must not claim the original run green.

## Checklist and summary

- DONE: Determine from retained native and child evidence whether completion was missed by the grader or the runtime violated lifecycle ordering.
  Both flagged completions have exact run/child identity and correct report/commit/notification/next-boundary order; auto-continue independently fails later preparation.
- DONE: Propose smallest safe correction and exact scope with ownership and separate downstream outcome findings; leave candidate unchanged.
  Proposed one retained-evidence correlation path in the existing Pi lifecycle observer, with bounded captured replay/negative controls; separate outcome typo retained for repaired-tip observation.

Return this proposal for a distinct FO disposition. No change to the approved fo-continuation validation verdict or gate is proposed.
