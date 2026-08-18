---
id: 0mdhjk9jdf10h1qnhvx5tvn8
title: Remove the tautological workflow-file tests
status: implementation
source: "Captain directive, CL 2026-08-17, after measuring the edge-advance case: \"the goal of the task is removal of those tautological tests\". Evidence from this session: 1307 lines guarding a 173-line file were green across three consecutive releases while the mechanism they covered failed silently, and validation later found two material defects in that same area."
started: 2026-08-18T00:32:37Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-remove-tautological-workflow-tests
issue:
---

Delete the tests that assert a workflow file contains the text someone wrote into it. The goal is removal, not replacement.

## Problem

A family of tests in this repository reads a `.github/workflows/*.yml` file and string-matches its contents. The shape is:

```go
strings.Contains(strings.ReplaceAll(ifCond, " ", ""),
    "steps.decision.outputs.advance=='true'")
```

That asserts a string is present in a file the author just edited. It cannot observe a release, and it cannot fail for a real reason.

The cost is measured, not asserted. `edge_reconcile_test.go` (648 lines), `edge_advance_wiring_test.go` (341), and `edge_advance_noregress_test.go` (318) guarded `edge_advance_decision.go` (173 lines) at a ratio of 7.5 to 1. All three were green on 2026-08-15, 08-16 and 08-17. Across those same three releases the `edge-advance` job skipped silently inside green runs, `next` froze at `v0.27.0-pre4` and fell 99 commits behind, and the captain ran a four-week-old first-officer contract against a current binary. The tests could not have caught it: the workflow step was present and merely evaluated false, and reading the YAML never evaluates anything.

Task `2d` already deletes those three as collateral, because the mechanism they covered is gone. At least eight more files follow the same shape and roughly 2000 lines remain, the largest being `internal/release/journey_workflow_test.go` at 827 lines. Nothing removes them.

The class is self-perpetuating: this kind of test is cheap to write and never goes red, so it accumulates until the ratio looks like coverage.

## Proposed approach

Ideation audited every file in the candidate set plus the grep-boundary extras (see Audit record below), applying the discriminating question per file: **can this fail for a reason other than someone editing the file it reads?** The audit sorts every assertion into four legitimate classes and one deletable class:

- **S1 Executes**: extracts the real step script from the workflow and runs it under bash, observing exit code, output, and on-disk state. These evaluate — the exact thing the edge-advance tautologies never did.
- **S2 Independent oracle**: compares parsed values against an oracle mirroring a source outside the file that moves on its own (GitHub's node24 schedule, the FO write-scope policy fixtures, the GHA permission model, the shipped-platform set).
- **S3 Cross-artifact agreement**: derives both sides from independently-edited locations and hardcodes neither (release.yml stamp target vs .goreleaser.yaml devBranch; Locate/Post steps' shared metrics dir; workflow vs docs command shape).
- **S4 Prohibition tripwire**: asserts the ABSENCE of a specific retired or banned surface (re-pinned SPACEDOCK_PINNED_CLAUDE_VERSION, trunk-to-`next`, `-v | tee` firehose, `@latest` floats, goreleaser→journey-ledger `needs:` edge, retired PTY/tmux surfaces). The ban is invisible in the file itself — a future editor cannot learn it from the diff — and each encodes a recorded decision or measured defect. Falsifiability comes from each ban's existing injection discriminator (mutate a COPY to contain the banned surface; the guard must red), not from an oracle.
- **DELETE — positive mirrors**: presence, ordering, or exact-equality assertions whose expected value is text the implementer wrote into the file under test (wiring `needs:` edges, command strings, `if:` gating expressions, cadence consts copied from the file, selector-set equality). These red only when the author's own text is edited — an event self-evident in the diff — and stayed green through the measured present-but-evaluates-false failure. This is the edge-advance class, including the same `if:`-expression-presence shape (`assertReleaseLedgerStepsSkipWhenNoProducerRun`, the manifest-tag pre-release skip check).

Deletion steps, ordered so each lands independently green:

1. Delete `internal/release/e2egate_workflow_test.go` and `internal/release/manifest_tag_gate_workflow_test.go` whole — pure positive wiring mirrors; the `e2e-gate` and `manifest-tag-gate` predicates themselves stay fully unit-tested in `e2egate_test.go` and `manifest_tag_gate_test.go`.
2. `internal/release/journey_workflow_test.go`: delete the seven positive-mirror guard tests (lines 12–110), the full job-graph pin `TestReleaseWorkflowJobGraphMatchesGitHubActions` + `keysOf`, and the skip-gating presence pair (`TestReleaseWorkflowSkipsLedgerWhenNoProducerRun`, `TestReleaseWorkflowGuardRejectsUngatedLedgerBuild`). Keep the goreleaser↛journey-ledger ban system (S4), all script-execution tests (S1), and the Locate/Post agreement test (S3).
3. `internal/release/workflow_exec_guard_test.go`: re-point the kept ban twins and their before-mutation sanity calls from `assertReleaseWorkflowPublishesJourneyCosts` to `assertGoreleaserDoesNotNeedJourneyLedger` directly, then delete the dead positive predicates (`assertRuntimeLiveWorkflowUploadsRawJourneyMetrics`, the positive body of `assertReleaseWorkflowPublishesJourneyCosts`, `assertReleaseLedgerStepsSkipWhenNoProducerRun`) and their now-unused helpers (`findExecutableStep`, `hasJourneyMetricsUploadAfter`, `pathBlockContainsLine`, `hasExecutableYAMLLine`, `workflowHasExecutableCommandContaining`). `parseWorkflowSteps`, `parseWorkflowJobs`, `needsList`, `executableShellCommands`, `isJourneyCostBuilder` stay — the keepers use them.
4. `internal/release/runtime_live_evidence_workflow_test.go`: reduce the two composite predicates to their ban clauses — the retained-dead-surface list (`TestLivePty`, `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS`, `pty-team-mode`, `Install tmux`, `inputs.effort`, …), the offline-controls must-stay-untagged ban, the job-level `CODEX_HOME` ban (runner.temp is unavailable in job-level env — a GHA-semantics fact), the retired `PI_OPENAI_CODEX_AUTH_JSON` and `SPACEDOCK_PI_LIVE_CHILD_MODEL` bans, and the no-unowned-selector rule over the `liveClaims` ledger. Delete the cadence/matrix/if exact-match mirrors, the copied consts, the exact selector-set equality, and the mutation controls of deleted mirrors; keep the two mutation controls that exercise kept bans ("legacy PTY flag", "offline live tag").
5. `internal/release/cilog_clean_output_workflow_test.go`: delete `TestLiveWorkflowStepsUseGotestsumOneRunShape` + `transformedLiveSteps` and the four-installs count; keep the firehose ban, the `@latest` float bans, and the install-script pin/sha256-verify policy checks.
6. Clause-level trims: `journey_delta_workflow_test.go` keeps the `pull-requests: write` platform-conformance assertion, drops the `if:`-equality and CLI-invocation-presence clauses; `claude_candidate_binary_workflow_test.go` drops the artifact-upload presence mirror (keeps the obsolete `./spacedock` ban); `live_registry_reconciliation_test.go` drops the `-run '^TestLiveCommon'` count-mirror lines (keeps the runtime-selector uniqueness that the S3 workflow-vs-docs extraction depends on, and all cross-artifact checks).
7. `gofmt -w ./cmd ./internal`; `go test ./...`; `go test ./... -race` after each step.

No new mechanisms are introduced: every change is deletion or re-pointing existing predicates at existing functions, serving AC-1 directly. No spike needed beyond the one already run: the survivor-reddening experiment below exercised the riskiest judgment (that a keeper really can red without its input file being touched) and the whole approach otherwise relies on mechanisms the existing suite already proves (yaml.v3 parsing, script extraction, `go test`).

The largest kept judgment block is the goreleaser↛journey-ledger ban system (~250 lines of shape/identity/multi-carrier armor in `journey_workflow_test.go` plus the ~50-line predicate). Recommendation: keep — the ban encodes the measured cut-blocked-on-never-fired-run incident, its armor tests exercise predicate code against adversarial inputs (they fail if `parseWorkflowJobs`'s alias/quoting resolution regresses, not only if release.yml is edited), and trimming legitimate-but-thorough tests is not this task. If the captain rules the armor disproportionate, deleting it removes a further ~250 lines; that is a gate decision, not a default.

## Audit record (ideation, 2026-08-18)

Verdict, evidence, and estimated surface per file. "Mirror" = expected value is text the implementer wrote into the file under test.

| File | Lines | Verdict | Est. − | Evidence (discriminating question) |
|---|---|---|---|---|
| release/e2egate_workflow_test.go | 151 | DELETE | 151 | Wiring mirror: `needs: e2e-gate` edge + SHA-bound command presence. Reds only on release.yml edits. Same shape as deleted `edge_advance_wiring_test.go`. |
| release/manifest_tag_gate_workflow_test.go | 163 | DELETE | 163 | Wiring mirror + `if: !contains(github.ref, '-')` presence — the exact `if:`-expression shape of the measured incident. |
| release/journey_workflow_test.go | 827 | MIXED | ~166 | Delete: 7 positive-mirror guards, full job-graph pin, skip-gating presence pair. Keep: S1 script executions (download/locate steps run under bash with stubbed `gh`, asserting exit codes, `$GITHUB_OUTPUT`, on-disk layout), S3 Locate/Post dir agreement, S4 ledger-edge ban system. |
| release/workflow_exec_guard_test.go | 440 | MIXED | ~200 | Helpers-only file; the positive predicates and their private helpers go dead once step 2 lands. Parser + ban predicate stay (used by keepers). ~10 insertions to re-point twins. |
| release/runtime_live_evidence_workflow_test.go | 298 | MIXED | ~178 | Cadence consts are copied file text compared by exact equality — mirror. Keep only the bans (dead surfaces, untagged offline controls, job-level CODEX_HOME, retired Pi secrets, unowned selectors). ~28 insertions to restructure. |
| release/cilog_clean_output_workflow_test.go | 138 | MIXED | ~55 | gotestsum-shape presence + install count are mirrors. Firehose ban, @latest bans, pin/verify policy are S4/S2 keepers. Found by the same grep; absent from the candidate list. |
| release/journey_delta_workflow_test.go | 41 | MIXED | ~8 | `pull-requests: write` is required by the GHA permission model (workflow default is contents:read) — external oracle, keep. `if:` equality and invocation presence are mirrors. Grep extra. |
| release/claude_candidate_binary_workflow_test.go | 210 | KEEP | 3 | S1: runs the extracted install script against real git+go fixtures; asserts fail-closed behavior by exit code and output. Only the 3-line artifact-upload presence mirror goes; the obsolete `./spacedock` ban stays. |
| contractlint/live_registry_reconciliation_test.go | 626 | KEEP | ~6 | S3: workflow-vs-docs command-shape agreement hardcodes neither side. Only the `-run '^TestLiveCommon'` count-mirror lines go. Found via `filepath.Join` sweep — the FO's literal grep missed it. |
| release/node24_actions_guard_test.go | 152 | KEEP | 0 | S2: `node24MinMajor` mirrors GitHub's deprecation schedule. Proven by reddening below. |
| release/claude_version_float_guard_test.go | 144 | KEEP | 0 | S4: bans the retired `SPACEDOCK_PINNED_CLAUDE_VERSION` re-pin (measured defect #395: CI frozen at 2.1.177 while users floated). The value under guard — which Claude version CI runs — tracks the installer's `latest`, a source that moves on its own; the pin ban is the offline-checkable coupling to it. Injection discriminator present and green. |
| contractlint/workflow_trunk_test.go | 189 | KEEP | 0 | S4: bans pre-flip trunk resolution (`branches: [next]`, `gh run list --branch next`) across ALL workflow files including future ones; discriminator control runs against its own temp fixtures. |
| contractlint/fo_write_core_mutation_gate_test.go | 171 | KEEP | 0 | S2: parses the classifier table and EVALUATES it against write-scope policy fixtures (two workflow dirs, 16 paths); expected values come from the FO write-scope policy, not the file's text. Explicit-pattern requirement defeats default-deny vacuity. |
| release/notes_extract_test.go | 142 | KEEP | 0 | S1: extracts the real empty-body guard from release.yml and EXECUTES it under `sh` against four tag shapes; found the dead `[ ! -s ]` form. Grep extra. |
| release/channel_agreement_guard_test.go | 284 | KEEP | 0 | S3/S2: release.yml stamp target vs .goreleaser.yaml devBranch (two artifacts, either drifts); stable≠edge distinctness; bridge manifest guarded for the RELEASED v0.20.0 binary — an independent released artifact. Found via `filepath.Join` sweep. |
| ensigncycle/codex_liveenv_test.go | 485 | KEEP | 0 | S1: extracts the claude/codex shim heredocs from the live workflow and EXECUTES them, asserting argv behavior. Grep extra. |

Also audited and cleared (adjacent same-shape suspects over `.goreleaser.yaml`): `goreleaser_guard_test.go` (platform-set oracle written independently of the YAML + header-vs-config S3 drift check). The three edge-advance files are excluded — task `2d` owns them.

**Survivor proof (checklist item 2, executed 2026-08-18).** Bumped the oracle `node24MinMajor["actions/checkout"]` from 5 to 6 — the exact edit GitHub's next deprecation wave forces — with `git status .github/` empty. `TestNode24ActionsPinnedAtMinimum` went RED at 12 pin sites across 5 untouched workflow files (`runtime-live-e2e.yml:58,121,331,598,819; docs.yml:31; install-e2e.yml:28; next-publish.yml:19; release.yml:29,129,173,300`). Reverted; suite green again. The survivor survives the reddening test; the candidate set does contain genuine survivors.

## Out of scope

The three edge-advance files, which task `2d` already removes. Do not touch that worktree or those files.

Adding new tests of any kind. If a real gap is found behind a deleted tautology, record it and let the captain decide; do not fill it in this task.

## Not every file in the sweep is the same — this is the trap

A blanket deletion would be wrong, and the FO's initial survey found likely survivors. `node24_actions_guard_test.go` and `claude_version_float_guard_test.go` appear to check a config value against an independent rule — that an action version is pinned rather than floating. That is a real value which can diverge on its own, and the workflow's own proof policy admits it: "a static check counts only when it tests a real value against an independent source that can diverge from it, not as a spelling check over a file the model reads."

The discriminating question for each file is therefore: **can this fail for a reason other than someone editing the file it reads?** If the expected value is just the text the implementer wrote into the file under test, it is a tautology and goes. If it compares the file against an independent source of truth — a pinned version policy, a published manifest, a released artifact — it stays.

Candidate set (line counts as of 2026-08-17, all under `internal/`):
`release/journey_workflow_test.go` 827, `release/runtime_live_evidence_workflow_test.go` 298, `release/claude_candidate_binary_workflow_test.go` 210, `contractlint/workflow_trunk_test.go` 189, `contractlint/fo_write_core_mutation_gate_test.go` 171, `release/manifest_tag_gate_workflow_test.go` 163, `release/node24_actions_guard_test.go` 152, `release/e2egate_workflow_test.go` 151, `release/claude_version_float_guard_test.go` 144.

## Expected surface and tolerance

Estimate net LOC change: −890, across 9 files (tolerance: net −890 ± 150, files 9 ± 2). Separately: ≈930 deletions, ≈40 insertions (the insertions re-point the kept ban twins and restructure the two mixed predicates; nothing new is built). This revises the seed's provisional −1500/9: the audit found five whole-file survivors and three more grep-boundary keepers — the exact mixed verdict the trap warned about — so the deletable mass is smaller than the sweep's raw line count. No gross tolerance is declared. Semantics changed: none — test-only changes; no production code, no command grammar, no stored format, no release behavior, no user-visible surface, and therefore no doc diff.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - No test in the repository asserts only that a workflow file contains a string the author put there.**
This is the measuring AC: the count of surviving positive-mirror assertions — presence, ordering, or exact-equality checks whose expected value is text the implementer wrote into the file under test — must be ZERO, counted over the audited set (the Audit record's 16 files, which extends the seed's 9-file candidate list with the grep-boundary extras). The count moves the wrong way the moment such a test is re-added. Verified by: the audit's per-file record, `go test -list` showing the deleted test functions gone, plus a reviewer reproducing the discriminating question on each surviving file.

**AC-2 - Every surviving check is falsifiable by something other than deleting its own expected text.**
Refined at ideation: the seed's phrasing ("can fail for a reason other than an edit to the file it reads") would strike the very survivors the seed names as the trap, because a prohibition tripwire by nature reds only when the banned surface is re-added to the file. Two arms, each with its own falsification proof:
(a) Oracle/behavior/agreement checks (S1–S3) can fail with the file untouched — verified by mutating the independent source, as executed for node24 (oracle bump reddened 12 pin sites across 5 untouched workflow files), or inherently by executing extracted scripts / comparing independently-edited artifacts.
(b) Prohibition tripwires (S4) assert the ABSENCE of a specific retired or banned surface, cite the decision or measured defect that banned it, and carry an injection discriminator proving the guard reds when the banned surface is added to a copy of the input — e.g. `TestClaudeVersionFloatGuardRejectsReintroducedPin`, `TestWorkflowTrunkLintDiscriminates`.
Fails if a surviving check satisfies neither arm — a positive assertion of authored text falsifiable only by deleting that text.

**AC-3 - No behavior lost with the deletions.**
Verified by: `go test ./...` green, and a named record of anything deleted that was guarding a live behavior, with the captain's decision on each. Fails if a deleted test was the only thing standing between a real regression and main — the same failure mode `2d`'s validator found when the reconcile deletion silently dropped a never-force-push guard.

## Test plan

The deliverable is deletion, so the proof is the surviving set and the green suite — no new tests, no fixtures, no live runs, and no tautology-detecting lint (a new standing enforcement mechanism needing its own captain-approved task).

- AC-1: after each deletion step, `go test ./internal/release/ ./internal/contractlint/ -list '.*'` shows the deleted test functions gone and the kept ones present; the reviewer re-applies the discriminating question to each survivor against the Audit record. Cost: minutes, offline.
- AC-2 arm (a): re-run the recorded node24 oracle experiment (bump `node24MinMajor["actions/checkout"]` to 6, observe RED with `.github/` untouched, revert). Cost: one edit + one `go test -run TestNode24`, recorded in the Audit record with the observed failure sites.
- AC-2 arm (b): every kept S4 ban retains its injection discriminator, and those discriminators run green inside the ordinary suite — their green IS the proof the bans still red on injection. Verified by the same `go test` run; no extra harness.
- AC-3: `go test ./...` and `go test ./... -race` green after every step and at the end; `gofmt -w ./cmd ./internal` clean. The named record of deleted-but-live-guarding checks is the "Gaps recorded for the captain" section below; each entry is a captain decision, not silent loss.
- Surface: `git diff --stat` insertions/deletions and file count measured against the declared net −890 / 9 files at the correction round.

## Gaps recorded for the captain (AC-3 ledger)

Deleting the mirrors removes the only OFFLINE checks of these wirings. In each case the predicate the wiring invokes remains fully unit-tested; what is lost is the textual check that the workflow still invokes it — a check the measured incident showed cannot distinguish present from effective anyway. Decide per row whether that residual tripwire value warrants a follow-up task (an evaluating check would be new work, not this task):

1. goreleaser `needs: e2e-gate` edge + SHA-bound `spacedock-release e2e-gate` invocation (e2egate_workflow_test.go): after deletion, a dropped edge or invocation is caught only by PR review until the next release run.
2. manifest-tag-gate wiring + its pre-release `if:` skip (manifest_tag_gate_workflow_test.go): same shape; the tag-vs-manifest predicate stays tested in `manifest_tag_gate_test.go`.
3. Journey-cost builder/publish presence, builder-before-goreleaser ordering, producer-found `if:` gating, and the journey-metrics upload paths (journey_workflow_test.go positive guards): the behavioral script tests that EXECUTE the download/locate steps remain; the step-wiring text checks go.
4. The live-cadence exact-match mirrors (runtime_live_evidence_workflow_test.go): cadence exclusivity as a spend policy loses its textual enforcement. Correction (validation cycle 1): this row previously also claimed "the kept bans cover the retired surfaces only" — false for the duration of one correction round, because the job-level `CODEX_HOME` / `PI_OPENAI_CODEX_AUTH_JSON` / `SPACEDOCK_PI_LIVE_CHILD_MODEL` bans extracted in step 4 were never applied to the real workflow file (Finding 1, validation cycle 1), proven 3-for-3 against main and 3-for-3 RED on real-file injection with the bans disconnected. Fixed by `TestRuntimeLiveWorkflowSecretBansHoldOnRealFile` (runtime_live_evidence_workflow_test.go), which now runs `assertLiveSecretsBansHold` against `readWorkflow(t, "runtime-live-e2e.yml")` directly; re-proven 3-for-3 RED by injecting each banned surface into the real on-disk `.github/workflows/runtime-live-e2e.yml` and reverting. As of this correction round the kept bans do cover the retired surfaces. Cadence exclusivity itself is still unenforced offline; a new task if the captain wants it.
5. `journey-delta-comment`'s `if: github.event_name == 'pull_request'` PR-only gating (journey_delta_workflow_test.go, if:-equality clause deleted step 6): nothing offline catches this job running on a `workflow_dispatch` trigger where there is no PR to comment on. Found by validation injecting `if: always()` into the real file; suite stayed green.
6. The journey-delta CLI invocation `go run ./cmd/spacedock-release journey-delta` (journey_delta_workflow_test.go's CLI-invocation-presence clause, and `workflowHasExecutableCommandContaining`, both deleted step 6): nothing offline catches the step's `run:` being swapped for an inert `echo`. Found by validation.
7. The Claude artifact upload path `${{ runner.temp }}/spacedock-live-bin/spacedock` (claude_candidate_binary_workflow_test.go, artifact-upload presence mirror deleted step 6): nothing offline catches the upload path being repointed to a nonexistent location, which would silently drop the candidate binary from the run's artifacts. Found by validation.
8. The gotestsum `--jsonfile <name>-detail.jsonl` archive flag on each live step (cilog_clean_output_workflow_test.go's `TestLiveWorkflowStepsUseGotestsumOneRunShape`, deleted step 5): nothing offline catches `--jsonfile` being dropped from a live step, which would silently lose the archive the FO's root-cause procedure (`skills/first-officer/references/first-officer-shared-core.md:73`) and the four `-detail.jsonl` artifact-upload paths depend on. Abandoning gotestsum entirely is still caught by `TestRuntimeLiveCommonSuiteTimeouts`'s workflow-vs-docs run-shape agreement; dropping only the archive flag is not. Found by validation.

## Stage Report: ideation

- DONE: Apply one discriminating question per file and record the verdict with evidence: can this fail for a reason other than someone editing the file it reads? Do not sort by filename or by intuition about what a guard "seems to" protect.
  Audit record table: 16 files, each with the question's answer as evidence; every file in the candidate set was read in full, plus 7 boundary files the seed grep missed or under-listed (found via a `filepath.Join`-style `workflows` sweep and the FO grep's own unlisted matches).
- DONE: Prove at least one survivor is genuinely a survivor by reddening it WITHOUT touching the file it reads — mutate the independent source instead. If nothing in the candidate set survives that test, say so plainly.
  node24 oracle bumped checkout 5→6 with `git status .github/` empty: RED at 12 pin sites in 5 untouched workflow files; reverted, suite green. The candidate set contains genuine survivors.
- DONE: Produce a per-file delete/keep list with line counts and the net figure, so the gate can approve a concrete deletion set rather than an intention.
  Per-file table with line counts and per-file deletion estimates; net −890 (≈930 deletions, ≈40 insertions) across 9 touched files, replacing the seed's provisional −1500/9.

### Summary

Audited the full workflow-file test surface with the discriminating question, assertion-level where files are mixed. Verdicts: 2 whole-file deletions (e2egate/manifest-tag wiring mirrors — the edge-advance shape exactly), 5 mixed files where positive mirrors go and executing/oracle/agreement/ban assertions stay, and 9 keepers. Key decisions for the gate: AC-2 refined with a prohibition-tripwire arm (the seed's literal phrasing would delete the very survivors it names as the trap — each ban instead proves itself by injection discriminator); the goreleaser↛journey-ledger ban system (~300 lines) is the largest kept judgment block, recommendation keep, alternative stated. Falsifiability of the proof: the node24 experiment reds if the oracle no longer diverges independently; the AC-3 gap ledger names all four places where deletion removes the only offline wiring check, each a captain decision. Baseline `go test ./internal/release/ ./internal/contractlint/` green before and after the reverted experiment; no repo code changed by this stage.

### Feedback Cycles

- Cycle 1: REJECTED — validation; surface 9 files / net −871 vs estimate 9 files / net −890 (98%, inside the declared ±150 / ±2); AC unchanged.

Surface is INSIDE tolerance, so no design-reset decision is required and none is recorded. The rejection is on two narrow findings, neither a mechanism failure.

**FO note on finding 1, for the record.** The three secret bans that are defined but never applied to the real workflow were introduced by the step-4 discriminator work, which the FO explicitly ruled in (option (a), "add the injection discriminator, and the ban survives PROVEN"). That ruling asked for proof the ban fires on injection and did not ask that the ban also run against the live file — so the gap follows from the ruling's wording as much as from the implementation. Recorded so the correction round is not read as worker error.

## Stage Report: implementation

- DONE: Land the seven deletion steps in order, running the suite after each so a bad step is isolated to its own commit rather than buried in a 900-line diff.
  6 commits on spacedock-ensign/remove-tautological-workflow-tests (steps 3's re-point+cleanup landed as one commit): d45566dbc (step1: whole-file deletes), 8f275384d (step2: journey_workflow_test.go mirrors+job-graph-pin), 412dff02b (step3: re-point ban twins, drop dead predicates), ff0d850ad (step4: runtime_live_evidence to ban clauses), 1a543252e (step5: cilog gotestsum mirrors), 81ce1fb08 (step6: last three clause-level trims). `go test ./...` ran green after every commit; final `go test ./... -race` also green (exit 0, all packages ok).
- DONE: Keep the goreleaser ban armor — the captain ruled keep. Do not trim it opportunistically while working in that file.
  Untouched except the mechanical re-point at 412dff02b (callers of assertReleaseWorkflowPublishesJourneyCosts switched to call assertGoreleaserDoesNotNeedJourneyLedger directly, same assertion, since the wrapper it called through was deleted as dead). All 6 ban tests (TestReleaseWorkflowGuardRejects/ToleratesGoreleaserNeedsJourneyLedger* and the multi-carrier variants) and assertGoreleaserDoesNotNeedJourneyLedger's body are byte-identical to before.
- DONE: Hold the declared net -890 across 9 files, and if the real figure drifts past the tolerance, stop and report rather than absorbing it.
  Final: net -871 (54 insertions, 925 deletions) across exactly 9 files — `git diff --numstat main..HEAD`. Inside the declared -890 ± 150 / 9 ± 2 tolerance; no drift to report.

### Deviations from the plan, both disclosed for the reviewer

1. Step 4 (runtime_live_evidence_workflow_test.go) surfaced an AC-2(b) gap not caught at ideation: the three kept bans (job-level CODEX_HOME, retired PI_OPENAI_CODEX_AUTH_JSON, retired SPACEDOCK_PI_LIVE_CHILD_MODEL) had zero injection discriminator anywhere in the package before this change (grepped `internal/release` and the whole repo) — they were only ever run against the real clean file, never proven to red on the banned surface. Escalated to team-lead before acting; proceeded after a ~9 min wait under the "do not idle" rule with the lower-risk, well-reasoned option. Added 3 mutation-control cases to the existing `TestRuntimeLiveWorkflowNamedEvidenceMutationControls` table (job-level CODEX_HOME / PI_OPENAI_CODEX_AUTH_JSON / SPACEDOCK_PI_LIVE_CHILD_MODEL reintroduced), extracting the three bans into a new `assertLiveSecretsBansHold`. Falsified by temporarily neutering `assertLiveSecretsBansHold` to `return nil` and re-running: exactly the 3 new subtests went RED (others stayed green), proving the discriminator is real, not vacuous; reverted before committing (ff0d850ad).
2. `workflowHasExecutableCommandContaining`'s deletion was moved from step 3 (where the plan listed it) to step 6, because its last caller (`journey_delta_workflow_test.go`'s CLI-invocation-presence clause) wasn't removed until step 6 — deleting it in step 3 would have broken the build mid-sequence, violating "each step lands independently green."

### Summary

Landed all 7 deletion steps from the ideation-approved plan, each its own commit, `go test ./...` green after every one. Final `go test ./...` and `go test ./... -race` both green. Surface: net -871 across 9 files vs. declared -890 ± 150 / 9 ± 2 — on target. `go test ./internal/release/ ./internal/contractlint/ -list` confirms every deleted mirror test function is gone and every kept S1-S4 test survives. The goreleaser↛journey-ledger ban armor is untouched per the captain's ruling. One judgment call is flagged above for the reviewer: closing a pre-existing AC-2(b) discriminator gap on 3 kept bans in runtime_live_evidence_workflow_test.go with 3 new mutation-control cases (not new coverage — proving bans this task chose to keep actually red).

### Addendum: per-ban discriminator rationale (team-lead ruling, option (a) confirmed)

Team-lead independently grepped and confirmed the same fact, ruled option (a), and asked for the per-ban rationale and an explicit audit-gap note. Recording both:

- job-level CODEX_HOME ban: earns its discriminator. Encodes a real GHA semantics fact — `runner.temp` is unavailable in job-level `env:`, only at step level — a platform constraint invisible in a diff, not a spelling convention. Covered by the "job-level CODEX_HOME reintroduced" case.
- retired PI_OPENAI_CODEX_AUTH_JSON ban: earns its discriminator. Bans a specific retired secret name; the retirement decision isn't visible in the file to a future editor. Covered by "retired PI_OPENAI_CODEX_AUTH_JSON reintroduced".
- retired SPACEDOCK_PI_LIVE_CHILD_MODEL ban: earns its discriminator, same reasoning (retired override invisible in the file). Covered by "retired SPACEDOCK_PI_LIVE_CHILD_MODEL reintroduced".

All three: each discriminator cost one map entry plus one anchor string (~3 lines), well under the cost of an unproven ban per team-lead's stated tradeoff, so delete-the-ban was not the better answer for any of the three.

Audit-gap note for the validator: the ideation Audit record classified these three as S4 keepers and cited AC-2(b) coverage without checking whether a discriminator already existed for them — it didn't (grepped independently by both me and team-lead). This is a gap in the ideation record, not scope creep introduced at implementation; the fix stays inside "prove the ban this task already chose to keep," not new coverage.

Surface: the discriminator addition is already counted in the net -871 / 9-file figure reported above (commit ff0d850ad, part of runtime_live_evidence_workflow_test.go's +42/-149).

## Review-finding disposition

Entered by validation, 2026-08-18. Reviewer observation authority only — classification is proposed, not authorized.

### Finding 1 (proposed Material, outcome defect) — the three kept secret bans no longer guard the real workflow

`assertLiveSecretsBansHold` (runtime_live_evidence_workflow_test.go:27) is invoked from exactly one place: the escape condition of the mutation-control loop at line 102, against a *mutated copy*. No test passes the real `readWorkflow(t, "runtime-live-e2e.yml")` to it. On main the same three bans lived inside `assertOneClaudeCadence`, which `TestRuntimeLiveWorkflowHasOneExplicitClaudeCadence` ran against the real file (`t.Fatal(err)`); that test was deleted as a cadence mirror and the bans went with its call site.

Proof (injection into the real `.github/workflows/runtime-live-e2e.yml`, run on throwaway checkouts of `main` and of HEAD):

| Injected surface | main | candidate |
|---|---|---|
| `CODEX_HOME: ${{ runner.temp }}/codex-home` in `codex-live` job env | RED | **GREEN** |
| `PI_OPENAI_CODEX_AUTH_JSON` in `pi-live` job env | RED | **GREEN** |
| `SPACEDOCK_PI_LIVE_CHILD_MODEL` in `pi-live` job env | RED | **GREEN** |

The disclosed mid-stage deviation is itself sound — neutering `assertLiveSecretsBansHold` to `return nil` REDs exactly the 3 new subtests and leaves the other 2 green, so the discriminators are real. But a discriminator proves the *predicate*; it does not connect the predicate to the real file. The predicate works and is never applied.

Four evidence fields: (1) any maintainer editing `runtime-live-e2e.yml` on the normal PR path; (2) a re-added job-level `CODEX_HOME` resolves empty because `runner.temp` is unavailable in job-level `env:` — the GHA-semantics defect the ban encodes — and the two retired Pi secrets can return, with nothing red offline; (3) `value-ac[AC-3]` "No behavior lost with the deletions", and `value-ac[AC-2]` arm (b), which requires the surviving check to assert the ABSENCE of the banned surface — it currently asserts nothing about the real file; (4) the 3-for-3 main-vs-candidate differential above.

Narrow fix, inside the approved surface: call `assertLiveSecretsBansHold` from a real-file test alongside the two at lines 65-75. No new mechanism, no scope change.

### Finding 2 (proposed Material, evidence defect) — the AC-3 ledger is incomplete and row 4 misstates

AC-3 promises "a named record of anything deleted that was guarding a live behavior, with the captain's decision on each." Four such losses are absent from `## Gaps recorded for the captain`. Each was confirmed by injecting the regression into the real file and observing `internal/release` + `internal/contractlint` (and `internal/ensigncycle` for the first three) stay green:

1. `journey-delta-comment`'s `if: github.event_name == 'pull_request'` gating — replaced with `if: always()`, nothing reds. The job would run on `workflow_dispatch` where there is no PR to comment on.
2. The journey-delta CLI invocation `go run ./cmd/spacedock-release journey-delta` — replaced with `echo`, nothing reds.
3. The Claude artifact upload path `${{ runner.temp }}/spacedock-live-bin/spacedock` — repointed to a nonexistent path, nothing reds; the candidate binary silently stops being uploaded.
4. The gotestsum `--jsonfile <name>-detail.jsonl` archive — dropping `--jsonfile` from a live step leaves the suite green. This is the archive `skills/first-officer/references/first-officer-shared-core.md:73` directs the FO to fetch for root cause, and the four `upload-artifact` paths name the `-detail.jsonl` files. (Abandoning gotestsum *entirely* is still caught, by `TestRuntimeLiveCommonSuiteTimeouts`' workflow-vs-docs run-shape agreement; dropping only the archive is not.)

Ledger row 4 additionally states "the kept bans cover the retired surfaces only." Per Finding 1 they cover nothing on the real file; the row needs correcting, not just extending.

Narrow fix: four ledger rows plus a row-4 correction. No code change.

### Finding 3 (proposed Polish) — AC-1's literal count is 2, not 0

Applying AC-1's own wording ("presence, ordering, or exact-equality checks whose expected value is text the implementer wrote into the file under test") to every survivor myself, two clauses count:

- `journey_delta_workflow_test.go:32` — `job.Permissions["pull-requests"] != "write"`. Exact equality; `write` is text in the file; it cannot red with the file untouched and carries no injection discriminator, so it satisfies neither AC-2 arm strictly.
- `live_registry_reconciliation_test.go:338` — pi "must retain `-failfast`". Presence of authored text. Pre-existing and untouched by this task; its claude/codex half is a genuine prohibition.

Under the ideation Audit record's approved S1-S4 taxonomy both are classified keepers (S2 GHA-permission-model oracle; spend policy), so the count is 0 under the interpretation the gate approved and 2 under AC-1's letter. Naming it so the captain owns the reading; no harm either way, ~4 lines.

Cleared boundary cases I checked and rejected as findings: `cilog_clean_output_workflow_test.go`'s positive install-script clauses evaluate `executableShellCommands`, so commenting out the `shasum -a 256 -c` line REDs the test with the text still present in the file — not a mirror; `channel_agreement_guard_test.go`'s `stableChannelBranch` const is a channel-policy oracle paired with a real cross-artifact comparison; `claude_version_float_guard_test.go`'s `${CLAUDE_VERSION:-latest}` clause is the positive face of a ban with a proven-RED discriminator.

## Stage Report: validation

- DONE: Re-apply the discriminating question yourself to every SURVIVING file, not to the audit's record of them — AC-1's zero count is the measuring criterion and it must be counted, not accepted.
  Swept `./internal` + `./cmd` for every test reading `.github/workflows/*.yml` and confirmed the 16-file audit set is complete (all other `workflows` matches are Spacedock workflow entities, a homonym). Read each survivor's assertions and counted: 2 positive-presence clauses survive under AC-1's letter, 0 under the approved S1-S4 taxonomy — Finding 3.
- DONE: Prove AC-2 both ways on real code: mutate an independent source for an oracle keeper and watch it red with its workflow untouched, and neuter a kept ban to confirm its discriminator fires. The three secret-ban discriminators were added mid-stage as a disclosed deviation — verify them rather than trusting the disclosure.
  Arm (a) reproduced independently: bumping `node24MinMajor["actions/checkout"]` 5→6 REDs `TestNode24ActionsPinnedAtMinimum` at 12 pin sites across 5 workflow files with `runtime-live-e2e.yml` md5-identical before and after. Arm (b): neutering `assertLiveSecretsBansHold` to `return nil` REDs exactly the 3 new subtests and leaves `legacy PTY flag` and `offline live tag` green — the disclosed discriminators are real, not vacuous. Extended arm (b) to all 11 kept bans by injecting each banned surface into the real file: 8 RED, 3 GREEN — Finding 1.
- FAILED: Check the AC-3 ledger against the 925 deleted lines: name anything removed that was guarding a live behavior and is not already recorded, since cycle 1 of a sibling task found exactly that failure when a deletion silently dropped a never-force-push guard.
  Exactly that failure recurred. Enumerated all 19 deleted test functions and 17 deleted helpers against the ledger's 4 rows; injection-tested each candidate loss on a throwaway checkout. Found the three orphaned secret bans (Finding 1, unrecorded and explicitly promised as kept) plus 4 unrecorded live-behavior losses and a misstatement in row 4 (Finding 2).

### Summary

Verified surface and hygiene clean: 54 insertions / 925 deletions / 9 files (net -871, inside the declared -890 ± 150 / 9 ± 2), `go test ./...` and `go test ./... -race` both exit 0, `gofmt -l ./cmd ./internal` empty, all 19 deleted test functions absent from `go test -list` and every named keeper present. The deletion judgment is sound: I re-derived it rather than accepting the audit, and 8 of 11 kept bans red on real-file injection. Recommendation is **REJECTED** on two material findings, both narrow fixes inside the approved surface and neither a mechanism failure: the three secret bans kept by step 4 are no longer applied to the real workflow (proven 3-for-3 against main, one-line fix), and the AC-3 ledger omits 4 live-behavior losses while row 4 asserts the bans still cover the retired surfaces. Deferred risks: none. Polish: AC-1's literal count is 2 rather than 0, a reading the captain owns.

## Stage Report: implementation (cycle 2)

Correction round for validation's REJECTED verdict (Finding 1 and Finding 2, both Material). Finding 3 (AC-1 count polish) explicitly left alone per the FO's routing.

- DONE: Wire the three secret bans to the real workflow so they can fail on the file they guard, and prove it by injecting the banned surface into the real path rather than only into the adversarial copy.
  Added `TestRuntimeLiveWorkflowSecretBansHoldOnRealFile` (runtime_live_evidence_workflow_test.go), running `assertLiveSecretsBansHold` against `readWorkflow(t, "runtime-live-e2e.yml")` directly, mirroring the two existing real-file tests for the other kept guards. Commit 356c4ccac. Proof: injected `CODEX_HOME`, then `PI_OPENAI_CODEX_AUTH_JSON`, then `SPACEDOCK_PI_LIVE_CHILD_MODEL` one at a time into the real on-disk `.github/workflows/runtime-live-e2e.yml` and re-ran the new test — RED 3-for-3 with the exact ban error each time; reverted the file after each injection (`git status` clean) and confirmed GREEN again.
- DONE: Complete the AC-3 ledger with the four omitted live-behavior losses and correct row 4, which claims coverage that finding 1 disproves.
  Added rows 5-8 (journey-delta-comment's PR-only `if:` gating, the journey-delta CLI invocation presence, the Claude artifact-upload path presence, the gotestsum `--jsonfile` archive flag) to `## Gaps recorded for the captain (AC-3 ledger)`, matching the validator's four named losses. Rewrote row 4 to state the bans were disconnected from the real file for the duration of this correction round (Finding 1) and now cover it again after the fix above, rather than asserting uninterrupted coverage.
- DONE: Re-measure and re-declare the surface after the fix; this round should be about one line of code plus ledger prose.
  Code delta this round: +6/-0, one file (runtime_live_evidence_workflow_test.go), commit 356c4ccac. Cumulative surface vs. main: net -865 (60 insertions, 925 deletions) across 9 files — `git diff --numstat main..HEAD`. Still inside the declared -890 ± 150 / 9 ± 2 tolerance (shifted from -871 to -865, the +6 lines of this round's fix).

### Summary

Both material findings closed. Finding 1: the three secret bans extracted in step 4 are now applied to the real workflow file via a dedicated real-file test, re-proven 3-for-3 RED by injecting each banned surface into the actual on-disk `.github/workflows/runtime-live-e2e.yml` (not a string copy) and reverting. Finding 2: the AC-3 ledger now names all 8 live-behavior losses this task's deletions expose, with row 4 corrected to state the true post-fix coverage rather than the false claim validation caught. `go test ./...` and `go test ./... -race` both exit 0 after the fix; `gofmt -l ./cmd ./internal` empty. Surface re-declared at net -865/9 files, inside tolerance. Finding 3 (AC-1's literal count) left untouched per the FO's routing — the captain's reading to make.
