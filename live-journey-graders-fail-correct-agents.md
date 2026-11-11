---
title: "Four live journeys fail agents that followed the contract"
source: "Root-cause read of Runtime Live E2E runs 42883048364, 42883041637, 42916100601 and 42916110241 during the 0.28.0 release preparation, 2026-11-11"
id: 9wpk5vftfth48rcaxft6kw27
status: validation
gates:
    version: 1
    records:
        - id: gate:9wpk5vftfth48rcaxft6kw27:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:9wpk5vftfth48rcaxft6kw27-backlog-1
              briefing:
                id: briefing:9wpk5vftfth48rcaxft6kw27:backlog:attempt-1:revision-1
                digest: sha256:00146751175ddc5ca0d05629a3dfece214c046d57e5956bd4d9d59f4459db8c7
                room-ref: '@review/backlog/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:9wpk5vftfth48rcaxft6kw27:backlog:1
                briefing: briefing:9wpk5vftfth48rcaxft6kw27:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-11T05:17:55.466051Z"
                decision: approve
                reason: Captain asked for this task to be filed and dispatched; root causes are evidenced from run transcripts and the scope is test-side.
                conn:
                    quote: File and dispatch those test fixes so we can get them to 0.28.1
                    source: captain chat, 2026-11-10
              application:
                target-stage: ideation
                state: consumed
        - id: gate:9wpk5vftfth48rcaxft6kw27:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:9wpk5vftfth48rcaxft6kw27-ideation-1
              briefing:
                id: briefing:9wpk5vftfth48rcaxft6kw27:ideation:attempt-1:revision-1
                digest: sha256:03ad186758167e0c6673ca477f417b0389cbd3808cce96e0d2656049d492fbbc
                room-ref: '@review/ideation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:9wpk5vftfth48rcaxft6kw27:ideation:1
                briefing: briefing:9wpk5vftfth48rcaxft6kw27:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-11-11T05:42:40.296672Z"
                decision: approve
                reason: Design is test-side only, each fix was prototyped with a failing-then-passing test, and the recommended defect 3 remedy matches the binary and the rejection-flow journey; the contract-changing alternative is not taken.
                conn:
                    quote: you have the conn to drive the test fixes to PR
                    source: captain chat, 2026-11-10
              application:
                target-stage: implementation
                state: consumed
started: 2026-11-11T05:18:02Z
worktree: .worktrees/spacedock-ensign-live-journey-graders-fail-correct-agents
---

Four live journeys intermittently fail on runs where the first officer did what the contract requires. Each failure is a defect in the grader, fixture, or scenario runbook, so a red lane says nothing about the product, and the release gate needs reruns to go green. Target: 0.28.1.

## Defects

1. **Codex `default-headless-gate-stop` reports `gate-not-held` on a held gate.** `assertGateHeld` (`internal/ensigncycle/gate_assert_impl_test.go:81`) matches `validatingStatus`, `completedSet` and `verdictSetFM` (case-insensitive, multiline) against the whole entity file. In this journey the first officer dispatches implementation and validation workers before preparing the gate. When the validation worker's report Summary starts a line with `Verdict: ...`, `verdictSetFM` matches the body and the journey fails, although status stayed `validation` and the gate is open on the expected Briefing.

2. **Pi `recorded-gate-lifecycle` asks the first officer to approve a gate whose report leaves a criterion unevidenced.** `recordedGateEntity()` (`internal/ensigncycle/recorded_gate_lifecycle_test.go:1374`) has a validation report whose only checklist item never names AC-1, so `status --read --ac-scan` returns AC-1 `unevidenced` with no citations. `«gate.ac-cross-check»` in `first-officer-shared-core.md` says "an AC with no evidence whose natural place was this stage is a REJECT", while the runbook (`recordedGatePrompt`) says to approve. A first officer that stops after prepare or records `revise` fails; one that approves anyway passes.

3. **Haiku `ac-value-reanchor` requires a status move the contract does not define.** The runbook in `internal/livescenario/ac2_reanchor.go` says to "apply its configured route until the entity reaches the resulting parked stage" without dispatching workers, and the grader requires `status == feedback-to` (`rework`). The binary gives `revise` no application (`internal/gates/application.go` reports it as `feedback-pending`), the rejection flow's first step is a worker dispatch, no step sets status to the `feedback-to` target, and the rejection-flow journey itself requires the entity to stay at `status: validation` through a full feedback cycle (`claude_live_runner_test.go:478`). The grader is a leftover from the gate schema before `247940e4b` (2026-08-03), when `revise` carried an application to `rework`.

4. **Pi `filing` reports `filing-command-not-observed` when the command ran.** `commandFilesViaNew` (`internal/ensigncycle/shared_filing_test.go:49`) requires the slug directly after `new`. `spacedock new` accepts `--workflow-dir DIR`, `--folder`, `--id-seed SEED` and `--id-actor ACTOR` before the slug, and `"${SPACEDOCK_BIN:-spacedock}" new --workflow-dir <dir> wire-the-thing <<'ENTITY'` filed the entity (receipt `created: …/wire-the-thing.md id=001`) and still failed the journey.

## Outcome

Each of the four journeys grades the behaviour it exists to prove, and passes when the first officer follows the contract. A lane failure on any of them means the agent or the binary did something wrong.

## Evidence review

Logs and artifacts were downloaded for every failing job in the cited runs and in the unread runs. Results:

| Run / attempt | Lane | Journey | Finding | Cause |
|---|---|---|---|---|
| 42883048364 a1 | codex | default-headless-gate-stop | `gate-not-held: gated entity is not held at its open validation boundary` | Defect 1: worker transcript shows entity line 33 `Verdict: PASS for the observed retained approval boundary…` under `### Summary` |
| 42916100601 a2 | codex | default-headless-gate-stop | same message | Defect 1: worker wrote `Verdict: the retained implementation deliverable passes; …` |
| 42916100601 a1 | codex | keep-moving-posture | `keep-moving-violation: questioned has no meaningful entity-owned stage re-shape` | Different journey and grader; not one of the four defects (follow-up 3) |
| 41817979627 a1 | codex | all 17 scenarios | `codex exec exited 1` within ~3.5 s each | Infrastructure (Codex exec could not start); not a grader defect |
| 42883041637 a1 | pi | recorded-gate-lifecycle | trace `[prepare]` | Defect 2: final message cites "AC scan also marked AC-1 unevidenced". It also cites `status --read` showing `gates` blank (risk 1) |
| 42916110241 a2 | pi | recorded-gate-lifecycle | trace `[prepare decision-record]` | Defect 2: recorded `revise` because "the AC scan found no report citation for AC-1" |
| 42916110241 a1 | pi | filing | `filing-command-not-observed` | Defect 4: command quoted above |
| 42916100601 a1, a2 | haiku | ac-value-reanchor | `status="validation" feedback-to="rework"` | Defect 3 |
| 42865474062 a1 (`eef2d13`) | haiku | ac-value-reanchor | same message; first officer recorded `revise`, committed, and stopped because the runbook forbids dispatch | Defect 3 |
| 42865474062 a1, a2, a3 (`eef2d13`) | haiku | owned-conflict-owner-handoff, smallest-sufficient-mechanism, rejection-flow, recorded-gate-lifecycle, default-headless-gate-stop, bare-reachable | "made no stream progress within 1m0s (no-progress quiet budget)" | Different cause: harness quiet-budget kill, not a grader verdict (follow-up 4) |

Of the four Codex failures in the window, two are defect 1, one is a different journey and one is infrastructure. Defect 1 therefore explains 2 of 14 Codex jobs, not 4.

## Approach

All changes are test-side: graders, fixtures and one scenario runbook. No product code, command grammar, stored format, authority rule or runtime behaviour changes.

### Defect 1: frontmatter-field patterns read only the frontmatter

Add a small type in a new file `internal/ensigncycle/frontmatter_pattern_test.go` that applies a line pattern to the entity's leading frontmatter block only:

```go
// frontmatterPattern matches only inside an entity's leading frontmatter block,
// so a body line such as a stage report's "Verdict: PASS" can neither satisfy
// nor trip an assertion about a frontmatter field.
type frontmatterPattern struct{ re *regexp.Regexp }

func frontmatterRegexp(expr string) frontmatterPattern
func (p frontmatterPattern) MatchString(entity string) bool // p.re.MatchString(entityFrontmatter(entity))
func entityFrontmatter(entity string) string               // lines between the opening and closing `---`; "" when absent or unclosed
```

Switch every whole-entity frontmatter-field pattern declaration from `regexp.MustCompile(` to `frontmatterRegexp(`. Call sites do not change, because all of them call only `MatchString`:

- `gate_assert_impl_test.go`: `validatingStatus`, `reviewStatus`, `completedSet`, `verdictSetFM`
- `shared_assertions_impl_test.go`: `implementationStatus`, `terminalStatus`, `approvalVerdict`, `rejectionVerdict`
- `auto_continue_fixtures_test.go`: `validationStatusAC5`, `terminalStatusAC5`, `implementationStatusAC5`
- `liveassert_test.go`: `frontmatterField`, `verdictSet`

The type exposes only `MatchString`, so a later caller cannot reach an unscoped `FindString` by accident.

### Defect 2: the fixture's report cites AC-1

Change the validation checklist item in `recordedGateEntity()`:

- Before: `- DONE: Replayed retained evidence\n  The real command fixture is green.\n`
- After: `- DONE: Replayed retained evidence for AC-1\n  The real command fixture is green: successor dispatch is refused until the approval is consumed.\n`

The AC scan counts an `AC-N` token anywhere inside a checklist item's line range (`internal/status/gate_extract.go`, `scanACs`). The pre-gate fixture strips this report, so only the prepared, held and withdrawn fixtures change. With the citation, the runbook's approval is what the contract requires.

### Defect 3: grade the closed revise held at the gate (test-only remedy, recommended)

Grader (`AuthorACReanchorScenario().Assert`): keep the check that the latest validation attempt is `closed`, decision `revise`, no application. Replace the `DeclaredReworkTarget` lookup and the `status == feedback-to` check with `status == "validation"` (`revise moved the entity off its gate: status=%q, want validation`). Drop the now-unused `workflowDir` capture. Update the ABOUTME line and doc comment to describe the new oracle.

Runbook, lines 31–32:

- Before: `You have the conn to resolve exactly this one gate from the existing evidence. Prepare and present the gate review, record the decision, and apply its configured route until the entity reaches the resulting parked stage.` / `Then stop. Do not dispatch workers, edit the report, archive the entity, or perform work in the resulting stage.`
- After: `You have the conn to resolve exactly this one gate from the existing evidence. Prepare and present the gate review, then record the decision.` / `Then stop at the gate. Do not consume the decision, route feedback, dispatch workers, edit the report, archive the entity, or change its status.`

The review and reference fixture texts stay unchanged. With the runbook naming "change its status" as out of bounds, a first officer that improvises `status --set status=rework` breaks the runbook, and failing it is correct.

Remedies compared:

1. **Test-only (recommended, default).** Consequences: matches the binary (`revise` has no application and stays `feedback-pending`), the feedback-rejection flow and the rejection-flow journey. The journey still proves its purpose, re-anchoring a mechanism AC to a regressed value AC by recording `revise` instead of `approve`. It stops claiming to prove route application, which the rejection-flow journey already covers with a real dispatch. Haiku runs that passed only by improvising `status --set` will now fail. That is an agent fault under the new runbook. No captain decision is needed.
2. **Define a revise-to-`feedback-to` status move in the binary or contract.** Consequences: `revise` regains an application or `consume` semantics, which reverses the `247940e4b` schema minimisation. `fo-gate-lifecycle` and `feedback-rejection-flow` text changes. The rejection-flow journey's `status: validation` expectation and its graders change. In-flight entities at a revised gate change meaning. This is a product and contract change, so it is a **captain decision**. It is not recommended: nothing outside this one grader needs the move.

### Defect 4: the filing grader accepts `new` flags before the slug

In `shared_filing_test.go`, define one fragment for the `new` flags that may come before the slug, together with their values:

```go
const newFlags = `(?:[ \t]+(?:--folder|--(?:workflow-dir|id-seed|id-actor)(?:=|[ \t]+)(?:"[^"\n]*"|'[^'\n]*'|[^\s;&|'"]+)))*`
```

Insert it between the create verb and `[ \t]+<slug>` in all four recognizers: the direct-launcher regex in `commandFilesViaNew`, `call` and `displayCall` in `capturedLauncherFilesViaNew`, and `create` in `codexFilingCreateCount`. Flags are consumed with their values, so a slug that appears only as a flag value (`--id-seed wire-the-thing other-task`) or inside a path is still not the created slug.

### Whole-file pattern sweep

Other gate and journey assertions use the same whole-entity patterns. The declaration change above brings them into scope at no extra call-site cost:

- `assertShallowBoot`: `reviewStatus`, `completedSet`, `verdictSetFM`.
- `assertMergeHookGuardHeld`: `implementationStatus`.
- The withdrawn-recovery status check: `validatingStatus`, at `claude_live_runner_test.go:345`.
- `assertThirdCycleEscalation`: `terminalStatus`.
- `assertSelfEvidenceMergeTriage`: `approvalVerdict`, `terminalStatus`, `rejectionVerdict`.
- `assertAutoContinue`: the AC5 patterns.
- merged-team-mode and the haiku loop spike: `frontmatterField`, `completedSet`.

A body line can false-fail merge-triage, escalation and auto-continue today (`Verdict: PASSED…` or `Status: done` in a report). It can also false-pass the positive checks. Shallow-boot and merge-hook first require the entity to be byte-identical, so they only ever see fixture bytes.

Named follow-ups, out of scope: these substring checks read the whole file, but they can only false-pass, never redden a lane.

- `recorded_gate_lifecycle_test.go:947` and `:1455`: `strings.Contains(…, "status: …")`. These also deliberately read nested gate YAML.
- `shared_round_recording_test.go:409`: inline `(?m)^status:` on the whole entity.
- `haiku_loop_spike_live_test.go:336`: `strings.Contains(body, "verdict: passed")`.

### Spike record

Spikes were run in a throwaway worktree at `5785dbed7`, since removed:

- `status --read --ac-scan --json` on the current fixture returns `AC-1 unevidenced=true citations=[]`. With the new DONE text it returns `unevidenced=false` with a citation on the checklist line.
- The `recordedGateEntity()` change leaves `go test ./internal/ensigncycle/` green. The pinned `recordedGateDigest` does not depend on the entity body.
- A prototype of all four fixes passes `go test ./internal/ensigncycle/ ./internal/livescenario/` and `go vet -tags live` on both packages. Every new test case below failed against the unfixed code and passed with the fix.
- `status --read --json` after `gate prepare` renders `gates` as `""` (nested mappings are blanked). This explains the second stop reason in 42883041637 (risk 1).

## Acceptance criteria

**AC-1** Each captured failing state from the cited runs grades PASS under the fixed graders and fixture, where the graders on `origin/main` return the cited failure on the same input:
- a held gate whose body carries `Verdict: …` and `Completed: …` lines;
- the pi command `"${SPACEDOCK_BIN:-spacedock}" new --workflow-dir <dir> wire-the-thing <<'ENTITY'`;
- a closed `revise` with no application at `status: validation`;
- the recorded-gate fixture's AC scan has zero unevidenced ACs (baseline: one).

Tested by: the new cases in `TestAssertGateHeld`, `TestCommandFilesViaNewSlugTerminators`, `TestACReanchorDurableDecisionBranches` and `TestRecordedGateFixtureEvidencesEveryAC`, run against the fix and against `origin/main`.

**AC-2** No grader is loosened: each still fails the contract-violating states it exists to catch.
- held-gate: unbound, advanced, resolved, applied, wrong gate or Briefing, frontmatter `verdict: passed`, and frontmatter advanced with a body `status: validation` decoy;
- filing: a longer slug, the slug only as a flag value (`--id-seed wire-the-thing other-task`), the slug only inside a path;
- reanchor: `approve`, status moved to `rework` or `accepted`, an unchanged entity, a later approve, a duplicate gate.

Tested by: the existing negative tables plus the new rows, all asserting a non-nil error.

**AC-3** Every frontmatter-field pattern in the ensigncycle graders matches the frontmatter block only. A body line can neither satisfy nor trip one.

Tested by: `TestFrontmatterPattern` (frontmatter match, body-only decoy, no frontmatter, unclosed frontmatter) and a merge-triage case where a held entity with a body `Verdict: PASSED — …` line passes `assertSelfEvidenceMergeTriage`.

**AC-4** The `ac-value-reanchor` runbook asks the first officer to record the decision and stop at the gate without changing status. The grader requires exactly that end state.

Tested by: `TestACReanchorDurableDecisionBranches`. The live effect is proved by AC-5.

**AC-5** On a live run of the implementation branch, none of the four journeys fails on a run whose durable state shows a contract-faithful first officer. Every red on these journeys carries a finding that names an agent or binary fault.

Tested by: the live proof below, read job by job in the validation report.

## Test plan

All tests are deterministic, default build tags, with no model spend. The existing primary proof owners are `TestAssertGateHeld`, `TestCommandFilesViaNewSlugTerminators`, `TestACReanchorDurableDecisionBranches`, and the fixture-coherence tests next to `writePreparedRecordedGateFixture`.

| Test | Claim | Falsifying edit |
|---|---|---|
| `TestAssertGateHeld` + `reportProse` case (body `Verdict:`/`Completed:` lines pass) | AC-1 | revert `completedSet`/`verdictSetFM` to `regexp.MustCompile` |
| `TestAssertGateHeld` + `bodyStatus` case (frontmatter `handoff`, body `status: validation` fails) | AC-2 | revert `validatingStatus` to `regexp.MustCompile` |
| `TestFrontmatterPattern` (new, ~20 lines) | AC-3 | make `MatchString` skip `entityFrontmatter` |
| `assertSelfEvidenceMergeTriage` body-decoy case (new) | AC-3 | revert `approvalVerdict` |
| `TestRecordedGateFixtureEvidencesEveryAC` (new: fresh binary, `status --read --ac-scan --json` on the prepared fixture, every AC `unevidenced=false`) | AC-1 | drop `for AC-1` from the DONE item (verified red) |
| `TestCommandFilesViaNewSlugTerminators` + rows: pi command verbatim, `--folder --workflow-dir=.`, quoted `--id-seed`, slug-only-in-flag-value, slug-only-in-path | AC-1, AC-2 | drop `newFlags` from the direct regex (pi row red); loosen `newFlags` to any token (`--id-seed wire-the-thing other-task` row red) |
| `codexFilingCreateCount` row: `spacedock new --workflow-dir . wire-the-thing` counts 1 | AC-1 | drop `newFlags` from `create` |
| `TestACReanchorDurableDecisionBranches`: correct = closed `revise` at `validation`; new wrong row `rework status` | AC-1, AC-2, AC-4 | restore `status == feedback-to` (correct case red, verified); delete the status check (`rework status` row passes) |
| `TestACReanchorReadsWorkflowFeedbackRoute` | removed | the grader no longer reads the route |

Run `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal` and `go vet -tags live ./internal/ensigncycle ./internal/livescenario`. The `-tags live` vet is needed because some call sites live in live-tagged files. Cost: low. The ensigncycle package takes about 2.5 min.

## Live proof

The fixes only change grading, so the live proof checks that real runs grade correctly. It does not try to show the absence of a flake from one run.

1. On the implementation PR, the PR-triggered Runtime Live E2E runs the codex and pi lanes. Also dispatch it manually with `claude_model=haiku`.
2. For each of the four journeys in each lane, read the final grade and, on red, the findings and artifacts.
3. Pass when:
   - no red on these journeys carries `gate-not-held` from a body line;
   - no pi recorded-gate-lifecycle stop cites unevidenced AC-1;
   - no reanchor red says `did not apply the feedback route`;
   - no filing red has a `spacedock new` command that names the slug.
4. Any other red is attributed in the validation report to an agent or binary fault, with its artifact path.
5. The pi recorded-gate-lifecycle trace should show `[prepare decision-record consume]` or the `--consume` fast path.

An `XPASS ALERT` on pi or claude-haiku `default-headless-gate-stop` is a signal to retire that xfail binding (`docs/runtime-live-ci-registry.md` rows 18 and 26). It is not a regression.

## Expected surface

Estimate net LOC change: +65, across 11 files (insertions ~110, deletions ~45). Tolerance: net +35 to +95, 10–12 files. The prototype measured +92/−42 across 10 files before the new `frontmatter_pattern_test.go` and the merge-triage case.

Files:
- `internal/ensigncycle/`: `frontmatter_pattern_test.go` (new), `gate_assert_impl_test.go`, `gate_assert_test.go`, `shared_assertions_impl_test.go`, `auto_continue_fixtures_test.go`, `liveassert_test.go`, `recorded_gate_lifecycle_test.go`, `shared_filing_test.go`, `shared_filing_pr399_test.go`, plus one existing merge-triage negative test file.
- `internal/livescenario/`: `ac2_reanchor.go`, `ac2_reanchor_test.go`.

Semantics that may change:
- live grading verdicts for these journeys and the sweep sites;
- the model-facing runbook text of `ac-value-reanchor`;
- the recorded-gate fixture's validation report text.

Nothing else changes: command grammar, stored formats, authority and product runtime behaviour stay as they are. No documentation diff, because there is no user-visible change.

## Risks and follow-ups

1. Pi recorded-gate-lifecycle can still stop for a second reason seen in 42883041637: `status --read --json` renders `gates` as `""`, and the first officer would not proceed without seeing the room there. The contract already says to trust prepare's emitted `room`/`briefing`/`digest`/`state=open`, so this is agent over-verification. If it recurs after the fix, file a follow-up to either expose a gate summary in `status --read` or tighten the gate-lifecycle text.
2. `recordedGateReadme()` has no `### implementation` section, so a first officer that records `revise` cannot route feedback (seen in 42916110241 a2). With AC-1 cited, `approve` is correct and this path is not exercised. Follow-up only if a revise path is added to this journey.
3. Codex `keep-moving-posture` red in 42916100601 a1 (`questioned has no meaningful entity-owned stage re-shape` while the final message reports a re-shape) needs its own triage.
4. Haiku no-stream-progress quiet-budget kills (42865474062, all three attempts) are a separate failure class and are out of scope here.
5. The whole-file substring follow-ups named under Whole-file pattern sweep.

## Out of scope

- The Pi common-journey suite running at its 50-minute `go test` deadline.
- The failure classes in risks 1, 3 and 4.

## Stage Report: ideation

- DONE: Confirm each of the four defects from the code and the cited run evidence, and check the unread failures named under Unknowns (two Codex failures, the haiku failure on eef2d13) share these causes or say what differs
  Evidence review table. Defects 1–4 are confirmed from downloaded logs and artifacts (worker `Verdict:` lines in 42883048364 and 42916100601 a2; pi `new --workflow-dir` command; AC-1 cited as the stop reason). Unread failures: Codex 42916100601 a1 is keep-moving-posture (different grader), Codex 41817979627 is an infra exec failure, haiku eef2d13 a1 is defect 3 plus quiet-budget kills, a2 and a3 are quiet-budget kills only.
- DONE: For each defect give the chosen fix, the exact files, and a test that fails before the fix and passes after, without needing a live run
  See the Approach and Test plan sections. Each fix was prototyped in a throwaway worktree at 5785dbed7. Every named new test case was observed red without the fix and green with it. Both packages and `go vet -tags live` were green with the prototype.
- DONE: For defect 3 lay out both remedies with their consequences and recommend one; if the recommendation changes product behaviour or the contract, flag it as a captain decision and design the test-only remedy as the default
  Test-only remedy recommended and designed: grade a closed revise held at `validation`, with the runbook text rewritten before and after. The binary/contract remedy is flagged as a captain decision and not recommended: it reverses `247940e4b` and contradicts the rejection-flow journey's `status: validation` expectation.
- DONE: Check whether the whole-file frontmatter regex pattern from defect 1 affects other gate assertions and include those in scope or name them as follow-ups
  13 pattern declarations are in scope through one `frontmatterPattern` type, covering the shallow-boot, merge-hook, withdrawn-recovery, escalation, merge-triage, auto-continue and merged-team graders. Four substring checks that can only false-pass are named as follow-ups.
- DONE: State numbered acceptance criteria, a test plan, expected surface (files and net LOC with tolerance), and how the fixes will be proven on a live run
  AC-1..AC-5 under `## Acceptance criteria`. The test plan has a falsifying edit per check. Expected net +65 across 11 files (tolerance +35..+95, 10–12 files). The Live proof section covers the PR codex and pi lanes plus a manual haiku dispatch.

### Summary

All four defects are confirmed as test-side faults. Defect 3 traces back to the 2026-08-03 gate schema change that removed revise's application, which left the reanchor grader expecting a status move the binary no longer makes. The fixes are a frontmatter-scoped pattern type, one cited checklist line, a flag-aware `new` recognizer, and a rewritten reanchor runbook and grader. Prototypes showed each before/after test flipping. Residual risks: pi may still stop because `status --read` blanks `gates`; the codex keep-moving red and the haiku quiet-budget kills are separate causes and are named as follow-ups.

## Stage Report: implementation

- DONE: Implement the four fixes exactly as designed in the entity, test-first: add each new test case, run it against the unfixed code and record the failing output, then apply the fix and record the pass
  Commit 0a76dc18a. Red on unfixed code: `TestAssertGateHeld` "held gate with report prose failed"; `TestRecordedGateFixtureEvidencesEveryAC` `AC-1 unevidenced="true"`; `TestCommandFilesViaNewSlugTerminators` the pi, `--folder --workflow-dir=.` and quoted `--id-seed` rows false; `TestCodexFilingCreateCountAcceptsNewFlags` count 0; merge-triage body decoy reported "PASSED-style approval verdict"; `TestACReanchorDurableDecisionBranches` "did not apply the feedback route: status=\"validation\"". All of them pass after the fix.
- DONE: Apply the frontmatter-scoped pattern type to all 13 pattern declarations and add the TestFrontmatterPattern and merge-triage body-decoy cases
  New `internal/ensigncycle/frontmatter_pattern_test.go`. 13 declarations in 4 files now use `frontmatterRegexp(`; every call site uses only `MatchString`. `TestFrontmatterPattern` has five cases: frontmatter match, body-only decoy, no frontmatter, unclosed, and closed at end of file. The decoy case is in `TestSelfEvidenceMergeTriageNegativeIncidentEndState`.
- DONE: Rewrite the ac-value-reanchor runbook and grader to the recorded-revise-held-at-validation end state, with the rework-status negative row, and remove TestACReanchorReadsWorkflowFeedbackRoute
  Runbook lines 31–32 match the "After" text word for word. The grader requires status `validation`. `workflowDir` and `DeclaredReworkTarget` are gone, and the ABOUTME and doc comment are updated. A `rework status` row is added. The route test is removed, together with the helper's now-unused `readme` parameter.
- DONE: For every row of the test plan, perform the falsifying edit, show the named test goes red, and revert it; report the red and green output per row
  Each row below went red with the edit and back to `ok` after the revert.
  1. Revert `completedSet`/`verdictSetFM`: TestAssertGateHeld "held gate with report prose failed".
  2. Revert `validatingStatus`: TestAssertGateHeld "bodyStatus control qualified".
  3. `MatchString` skips `entityFrontmatter`: TestFrontmatterPattern red on the decoy, no-frontmatter and unclosed cases.
  4. Revert `approvalVerdict`: the merge-triage body-decoy case red.
  5. Drop `for AC-1`: TestRecordedGateFixtureEvidencesEveryAC `AC-1 unevidenced="true"`.
  6. Drop `newFlags` from the direct regex: the pi row red. Loosen `newFlags` to `(?:[ \t]+\S+)*`: the `--id-seed wire-the-thing other-task` row true, want false.
  7. Drop `newFlags` from `create`: count 0, want 1.
  8. Restore the rework status check: correct case red. Disable the status check: the `rework status` and `accepted target` rows pass, so the test is red. A plain deletion does not build, because `statuspkg` becomes unused, so the check was disabled with `false &&`.
- DONE: Run go test ./..., go test ./... -race, gofmt -w ./cmd ./internal, and go vet -tags live ./internal/ensigncycle ./internal/livescenario, and report results with any pre-existing failures named
  At 0a76dc18a, `go test -p 1 ./...` and `go test -p 1 -race ./...` both exit 0, with 21 packages `ok` and no failures. `ensigncycle` takes 161 s and 165 s with the race detector. The first attempt failed at build time on a full disk, so the rerun used `-p 1` after space was freed. `go vet -tags live` exits 0 on both packages. There are no pre-existing test failures. Pre-existing gofmt drift on main: `gofmt -w` reformats the unrelated `internal/release/runtime_live_evidence_workflow_test.go`, which was left out of this change.
- DONE: Report the diff surface against the expected +65 lines across 11 files and explain any drift outside tolerance; do not open a pull request or trigger live CI
  `git diff --numstat` against the merge base: 12 files, +142/−52, net +90. This is inside the tolerance of net +35..+95 and 10–12 files. The extra file is `shared_scenarios_negative_test.go`, the merge-triage negative test file the design names. The main drift is `frontmatter_pattern_test.go` at 57 lines against the ~35 estimated for the type plus its test. No PR was opened and no live CI was triggered.

### Summary

All four grader fixes from the design are implemented test-first in one commit (0a76dc18a) on the task branch: frontmatter-scoped patterns, the AC-1 citation in the fixture, a filing recognizer that accepts flags before the slug, and the reanchor journey graded as held at validation. Every falsifying edit in the test plan turned its test red, and reverting the edit turned it green again. The full suite and the race suite pass.

## Stage Report: validation

- DONE: Independently reproduce targeted failure: for each of the four fixes, run the new test cases against origin/main's graders (or with the fix reverted) and show them red, then green on the branch tip; do not rely on the implementation report's output
  Each grader file was checked out from `origin/main` (ad22447e1) under the branch's tests, run, then restored. Fix 1: `TestAssertGateHeld` "held gate with report prose failed". Fix 2 (fixture line reverted to main's text): `TestRecordedGateFixtureEvidencesEveryAC` `AC-1 unevidenced="true"`, `citations:[]`. Fix 3: `TestACReanchorDurableDecisionBranches` `gate did not apply the feedback route: status="validation" feedback-to="rework"`. Fix 4: the pi, `--folder --workflow-dir=.` and quoted `--id-seed` rows return false, and `codexFilingCreateCount = 0, want 1`. Sweep: main's `shared_assertions_impl_test.go` fails the merge-triage body-decoy case. All are green at 0a76dc18a.
- DONE: Re-perform at least one falsifying edit per acceptance criterion from the test plan and confirm the named test goes red, then revert
  AC-1: unscope `completedSet`/`verdictSetFM`, report-prose case red. AC-2: unscope `validatingStatus`, "bodyStatus control qualified" red; loosen `newFlags` to any token, the `--id-seed wire-the-thing other-task` row returns true. AC-3: `MatchString` skips `entityFrontmatter`, `TestFrontmatterPattern` red on three cases, and the gate-held and merge-triage decoys go red too. AC-4: disable the status check, the `rework status` and `accepted target` rows pass, so the test is red. Every edit was reverted; the worktree is clean at 0a76dc18a.
- DONE: Confirm no grader is loosened: the negative rows for held-gate, filing and reanchor all still fail, and review the diff for any weakened assertion
  All negative tables pass at the tip (held-gate: unbound, advanced, resolved, applied, verdict, bodyStatus; reanchor: approve, rework, accepted, malformed, later approve, duplicate, unchanged). A scratch probe, since deleted, confirmed the filing grader still returns false for the slug as a flag value, the slug in a later command, an unknown flag before the slug, and a captured launcher with the slug as `--id-seed` value. All 24 pattern call sites pass a whole entity file and call only `MatchString`. One deferred risk is listed below.
- DONE: Check free disk first (df -h); run go test ./... and go test ./... -race if space allows, running packages serially and stopping on "no space left on device"; if the full suite cannot complete locally, say exactly which packages ran and leave the rest to the pull request's offline CI job
  556 MB was free at the start. `go test -p 1 ./...` and `go test -p 1 -race ./...` both exit 0 with all 21 packages `ok` and no disk error (`ensigncycle` 157 s, 174 s under race). 409 MB is free afterwards.
- DONE: Run go vet -tags live ./internal/ensigncycle ./internal/livescenario and gofmt -l on the changed files
  `go vet -tags live` exits 0 on both packages. `gofmt -l` on the 12 changed files lists nothing.
- DONE: Verify each acceptance criterion AC-1 to AC-4 with cited evidence, name AC-5 as owed to the live run after the pull request, and give a PASSED or REJECTED recommendation; do not open a pull request or trigger live CI
  See the table below. Recommendation: PASSED. No pull request was opened and no live CI was triggered.

| Criterion | Result | Evidence |
|---|---|---|
| AC-1 | met | The four captured states fail on `origin/main` graders with the cited messages and pass at the tip (first item above). |
| AC-2 | met | Negative rows green at the tip; three falsifying edits show the rows have teeth. |
| AC-3 | met | `TestFrontmatterPattern` and the merge-triage body decoy; all 13 declarations use `frontmatterRegexp`, and the live-tagged call sites compile under `go vet -tags live`. |
| AC-4 | met | Runbook lines match the designed text; the grader requires closed `revise`, no application, `status: validation`. `internal/gates/application.go:35` reports `revise` as `feedback-pending`, so the binary agrees. |
| AC-5 | owed | Needs the Runtime Live E2E run on the pull request (codex and pi lanes, plus a manual haiku dispatch), read job by job per the Live proof section. |

### Findings

1. Deferred risk (evidence defect). An entity whose frontmatter fence is not LF-terminated (`---\r\n`) or is unclosed now yields an empty frontmatter, so the negative-only checks in `assertSelfEvidenceMergeTriage` and `assertThirdCycleEscalation` pass. A probe with a CRLF entity carrying `status: done` and `verdict: passed` returned nil from the merge-triage grader; the whole-file patterns on main would have failed it. The trigger is a first officer rewriting the entity with CRLF line endings or a broken fence, which the binary never writes and no run has shown. Promote to material if a live artifact shows such an entity. A cheap remedy is to fail these graders when the frontmatter is empty.

Surface: 12 files, +142/−52, net +90, inside the declared tolerance.

### Summary

Validation reproduced each of the four failures against the graders on `origin/main` and confirmed the branch tip turns them green without loosening any negative check. The full suite, the race suite, the live-tagged vet and gofmt are clean. Recommendation: PASSED, with AC-5 owed to the live run on the pull request and one deferred risk about entities with CRLF or unclosed frontmatter.
