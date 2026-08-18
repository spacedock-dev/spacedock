---
id: 0mdhjk9jdf10h1qnhvx5tvn8
title: Remove the tautological workflow-file tests
status: ideation
source: "Captain directive, CL 2026-08-17, after measuring the edge-advance case: \"the goal of the task is removal of those tautological tests\". Evidence from this session: 1307 lines guarding a 173-line file were green across three consecutive releases while the mechanism they covered failed silently, and validation later found two material defects in that same area."
started: 2026-08-18T00:32:37Z
completed:
verdict:
score:
worktree:
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
4. The live-cadence exact-match mirrors (runtime_live_evidence_workflow_test.go): cadence exclusivity as a spend policy loses its textual enforcement; the kept bans cover the retired surfaces only. If the captain wants cadence enforcement that actually evaluates the `if:` expressions, that is a new task.

## Stage Report: ideation

- DONE: Apply one discriminating question per file and record the verdict with evidence: can this fail for a reason other than someone editing the file it reads? Do not sort by filename or by intuition about what a guard "seems to" protect.
  Audit record table: 16 files, each with the question's answer as evidence; every file in the candidate set was read in full, plus 7 boundary files the seed grep missed or under-listed (found via a `filepath.Join`-style `workflows` sweep and the FO grep's own unlisted matches).
- DONE: Prove at least one survivor is genuinely a survivor by reddening it WITHOUT touching the file it reads — mutate the independent source instead. If nothing in the candidate set survives that test, say so plainly.
  node24 oracle bumped checkout 5→6 with `git status .github/` empty: RED at 12 pin sites in 5 untouched workflow files; reverted, suite green. The candidate set contains genuine survivors.
- DONE: Produce a per-file delete/keep list with line counts and the net figure, so the gate can approve a concrete deletion set rather than an intention.
  Per-file table with line counts and per-file deletion estimates; net −890 (≈930 deletions, ≈40 insertions) across 9 touched files, replacing the seed's provisional −1500/9.

### Summary

Audited the full workflow-file test surface with the discriminating question, assertion-level where files are mixed. Verdicts: 2 whole-file deletions (e2egate/manifest-tag wiring mirrors — the edge-advance shape exactly), 5 mixed files where positive mirrors go and executing/oracle/agreement/ban assertions stay, and 9 keepers. Key decisions for the gate: AC-2 refined with a prohibition-tripwire arm (the seed's literal phrasing would delete the very survivors it names as the trap — each ban instead proves itself by injection discriminator); the goreleaser↛journey-ledger ban system (~300 lines) is the largest kept judgment block, recommendation keep, alternative stated. Falsifiability of the proof: the node24 experiment reds if the oracle no longer diverges independently; the AC-3 gap ledger names all four places where deletion removes the only offline wiring check, each a captain decision. Baseline `go test ./internal/release/ ./internal/contractlint/` green before and after the reverted experiment; no repo code changed by this stage.
