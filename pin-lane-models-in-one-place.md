---
title: The lane model ids are written down five times, and the Pi lane runs an older Luna
status: implementation
score: 0.8
source: "Captain, 2026-10-04: pin gpt-6-luna for CI, and move the Claude lane to sonnet-5.5. Local smoke on openai-codex/gpt-6-luna passed (root and child, 166s, grade artifact written)."
id: q67rfhn33tncbvf9xqbq960r
started: 2026-10-04T05:48:52Z
worktree: .worktrees/spacedock-ensign-pin-lane-models-in-one-place
---

## Problem

Two problems with one cause.

Each lane model id is written down in several places. The Pi lane's
`openai-codex/gpt-5.6-luna:max` appears in `internal/ensigncycle/pi_liveenv.go`, in
three branches of the codex shim inside `.github/workflows/runtime-live-e2e.yml`, in
its step summary, in `docs/runtime-live-ci.md`, and in fixtures. The Claude lane's
`claude-sonnet-5` appears in the workflow matrix, in the harness's accepted-value
list, in the context-budget table, and in an artifact name.

And the Pi lane runs an older Luna. A local front-door smoke on
`openai-codex/gpt-6-luna:max` passed today: 5 reads, ensign skill at read 2, root and
child both `openai-codex/gpt-6-luna`, 166s, grade artifact written, no model error.

## Design

**One place for the lane models.** Put the ids in `internal/release` or
`internal/ensigncycle` as named constants, and print them for the workflow the way
`stamp-paths` prints the stamp targets: a small command that emits the model id per
lane and per auth path. The workflow matrix, the codex shim, the step summary, the
docs and the fixtures then consume that output instead of repeating literals. No
model id may remain in the workflow file.

**Pin the Pi lane to gpt-6-luna**, both auth paths: the OAuth spelling used with the
subscription, `openai-codex/gpt-6-luna:max`, and the API-key spelling beside it. Verify
the API-key spelling against the provider before pinning it; do not assume it from the
OAuth one.

**Claude moves only when the id is confirmed.** `claude-sonnet-5.5` is refused by the
installed CLI today: "There's an issue with the selected model (claude-sonnet-5.5). It
may not exist or you may not have access to it." Keep `claude-sonnet-5` until an
exact id is confirmed, and add the new value in the same single source when it is.

## Out of scope

- Recorded fixtures that name a past model on purpose. `auto_continue_replay_test.go`
  references real `claude-sonnet-5` streams and the CI artifact they came from; those
  are evidence, not configuration, and renaming them falsifies the record.
- Auth, quota, and thinking level.

## Expected surface and tolerance

`internal/ensigncycle/pi_liveenv.go`, the new print command, the workflow, the docs,
and the fixtures that assert the shim's argv. Estimate net +70, across 6 files.
Tolerance net +30 to +140, at most 8 files.

## Acceptance criteria

**AC-1 (VALUE) - One place holds each lane model.**
Verified by: the printed model list is what the workflow and the harness use, and a
test fails when they diverge. Falsifier: change the constant and require the workflow
site to follow.
**AC-2 - The Pi lane runs gpt-6-luna on both auth paths.**
Verified by: the harness constants and the shim's argv name it; a local front-door
smoke passes with root and child on it. Falsifier: leave the old literal and the shim
`--model` argument must differ from the constant.
**AC-3 - No model literal remains in the workflow.**
Verified by: searching the workflow for each id and requiring no match outside the
command that prints them. Falsifier: add one and the check must fail.
**AC-4 (no-regression) -** Real Pi discovery is unaffected, the Claude lane keeps its
current model until the id is confirmed, and the offline suite passes.

## Test plan

Behavioural tests only, with the model supplied independently. Every test must state
the reason it could fail. No test may assert a file's own text against a copy of it,
and no expected value may come from the thing it tests. The smoke is a local run with
the captain's credentials, never a CI lane run.

## Stage Report: implementation

- DONE: Put the lane model ids in one Go source and print them for the workflow, the way stamp-paths prints the stamp targets. Then make the workflow matrix, the codex shim, the step summary, the docs and the fixtures consume that output, so no model literal remains in .github/workflows/runtime-live-e2e.yml.
  `internal/release/livemodels.go` holds the five ids and `LiveModels()`; `spacedock live-models` prints `lane.auth=id` (`--get <key>` prints one) and gated `internal/release/livemodels_test.go` proves no literal remains and every resolved key exists. `grep -nE 'claude-sonnet-5|claude-opus-4-8|gpt-6-luna|gpt-5.6-luna' .github/workflows/runtime-live-e2e.yml` = no match; the offline job resolves the five ids into `needs.offline.outputs.*`, consumed by the claude matrix, the Codex shim (`SPACEDOCK_LIVE_CODEX_MODEL`), the Pi summary, and the journey-delta artifact name. Commit fbb6cc8df.
- DONE: Pin the Pi lane to gpt-6-luna on both auth paths. The OAuth spelling openai-codex/gpt-6-luna:max passed a local front-door smoke today with root and child on it. Verify the API-key spelling against the provider before pinning it; do not assume it from the OAuth one.
  `release.PiOAuthModel`/`PiAPIKeyModel` = `openai-codex/gpt-6-luna:max` / `openai/gpt-6-luna:max`; `pi_liveenv.go` aliases them. API-key spelling verified against the installed pi-ai catalog: `providers/data/openai.json` declares `gpt-6-luna` with `thinkingLevelMap.max="max"` (and `openai-codex.json` the same), so `openai/gpt-6-luna:max` is real, not assumed.
- DONE: Keep the Claude lane on claude-sonnet-5. The id claude-sonnet-5.5 is refused by the installed CLI today, so record it as an input the captain must confirm and do not pin it.
  `release.ClaudeSonnetModel = "claude-sonnet-5"` (opus stays `claude-opus-4-8`); `docs/runtime-live-ci.md` records `claude-sonnet-5.5` as refused by the installed Claude CLI and pending captain confirmation. Not pinned.
- DONE: Do not rename the recorded fixtures that reference real claude-sonnet-5 streams; they are evidence of past runs.
  `auto_continue_replay_test.go`, `shared_filing_pr399_test.go`, `journey_workflow_test.go`, `dispatch/contextbudget_parity_test.go`, `journeymetrics/claude_test.go`, `shallow_boot_measure_unit_test.go` untouched (`git show --stat fbb6cc8df` lists none of them).
- DONE: Prove it: the focused tests, the print command's output, and one local front-door smoke on the pinned Pi model with a distinct and EMPTY SPACEDOCK_LIVE_ARTIFACT_DIR. State the reason each test you add could fail.
  Print output: `claude.sonnet=claude-sonnet-5`, `claude.opus=claude-opus-4-8`, `codex.exec=gpt-6-luna`, `pi.oauth=openai-codex/gpt-6-luna:max`, `pi.api-key=openai/gpt-6-luna:max`. Live smoke is PENDING (see Summary). Falsifiers, each run:
  - `TestLiveModelsCommandPrintsPinnedLaneModels` (cli): fails if a lane constant changes without the print — verified by reverting `PiOAuthModel` to `gpt-5.6-luna`, which failed the exact-stdout assert.
  - `TestLiveModelsCommandGetPrintsOneID` (cli): fails if a key is renamed or an unknown key exits 0/prints empty.
  - `TestRuntimeLiveWorkflowCarriesNoLaneModelLiteral` (release): fails if a literal is re-added — verified by injecting `claude-sonnet-5` into the workflow, which failed.
  - `TestRuntimeLiveWorkflowResolvedKeysExist` (release): fails if a workflow `--get` key has no `LiveModels()` entry — verified by renaming `--get claude.sonnet` to `claude.sonnetX`, which failed.
  - `TestRuntimeLiveWorkflowLaneModelWiring` (release): fails if an offline output is renamed on one side only — verified by suffixing the `pi_api_key` consumer, which failed.
  - `TestCodexLiveWorkflowPinsOnlyExecToLuna` (updated): supplies `SPACEDOCK_LIVE_CODEX_MODEL=luna-sentinel-model`; fails if the shim hardcodes a model — verified by replacing `$live_model` with `hardcoded-luna`, which failed.
  - `TestRuntimeLiveClaudeShimSetsMaximumEffort` (updated): uses `release.ClaudeSonnetModel` only as the arbitrary passthrough arg; fails if the Claude shim stops prepending `--effort max`.

### Summary

Pinned the lane model ids in one Go source (`internal/release/livemodels.go`), added `spacedock live-models` to print them, and rewired the workflow matrix, the Codex exec shim, the Pi step summary, the docs, and the harness/fixtures to consume that source — the live workflow now carries no model literal. The Pi lane runs `gpt-6-luna` on both auth paths; the API-key spelling is verified against the installed provider catalog, and the Claude lane stays on `claude-sonnet-5` with `claude-sonnet-5.5` recorded as pending captain confirmation.

**Live proof is PENDING, not passed.** A local front-door smoke (`TestLivePiFrontDoorSmoke`, a distinct empty `SPACEDOCK_LIVE_ARTIFACT_DIR="/tmp/spacedock-pin-lane-smoke-1791093534"`) finished `ok` in 193.660s and its grade artifact read `verdict: pass` with `root_session.model = openai-codex/gpt-6-luna` (179506ms) and `child_session.model = openai-codex/gpt-6-luna` (77719ms), 5 reads and ensign skill at read 2. The FO's mid-run steering says the machine was saturated (1-min load ~118 earlier, ~85 at steering) and that a live result taken now must be reported PENDING; the smoke had already completed when that arrived, so it was not killed and its numbers are not treated as a quiet-machine conclusion. Load at stop: `8.48 21.01 21.19` on 10 cores. The captain should re-run the smoke on a quiet machine.

Deviation: 13 files, net +289 (design estimate +70, tolerance +30..+140, at most 8 files). The extra surface is forced by the mandated behavioral proof (new print command + constants + two new test files) and by the fixtures the Problem names; no scope was added beyond the design's named sites.


## Review-finding disposition

Validator observations against `fbb6cc8df`; proposals only, not FO authorization. Candidate code and HEAD remain unchanged.

### V1 — documentation still duplicates current lane pins

- Evidence: `docs/runtime-live-ci.md:153–154,177–179` embeds the OAuth/API-key, Sonnet, Opus and Codex ids; adding a printer paragraph did not make the docs consume it. Line 154 also incorrectly calls the changed model ID “unchanged.”
- Released user and normal workflow: maintainers change a lane pin in `internal/release`, and users consult the live-CI documentation.
- Observable harm: the promised one-edit update still requires hand-syncing documentation; the implementation's all-consumers claim is false.
- Authority: value-ac[AC-1] One place holds each lane model; the design explicitly includes docs as a consumer.
- Trigger evidence: inspected the committed docs; unlike the workflow, they still repeat current configuration rather than directing readers to the relevant `live-models --get` outputs.
- Proposal: **outcome defect / material**, owned by this task; narrow fix to remove current-config literals in favor of printer keys/commands (retain the explicitly pending Claude 5.5 input and historical evidence). FO disposition pending.

### V2 — wiring guards do not establish lane identity

- Evidence: `internal/release/livemodels_test.go:51–86` checks independent names and global presence, but not which resolved key populates each output. On an isolated archive, exchanging only the resolver's `--get claude.sonnet` and `--get claude.opus` values leaves all three new release tests passing.
- Released user and normal workflow: routine PR Sonnet lane versus explicit Opus cadence, through the supported workflow resolver.
- Observable harm: the asserted proof accepts a PR lane receiving Opus and the Opus lane receiving Sonnet; it cannot establish the claimed producer-to-consumer correspondence.
- Authority: value-ac[AC-1] The printed list is what each workflow lane uses, and a test fails when they diverge.
- Trigger evidence: actual `go test ./internal/release -run '^TestRuntimeLiveWorkflow(CarriesNoLaneModelLiteral|ResolvedKeysExist|LaneModelWiring)$' -count=1` in a disposable archive returned 0 after that swap; the unchanged candidate wiring itself is correct.
- Proposal: **evidence defect / material**, owned by this task; narrow independent expected key-to-output and output-to-site proof, executing the resolver where practical. No new controller/lifecycle layer or architectural reset is needed. FO disposition pending.

### V3 — finite no-literal guard coverage

- Evidence: `livemodels_test.go:20` authors its forbidden set independently of the workflow and `LiveModels()`, so it is **not tautological**. Known `claude-sonnet-5` insertion fails. An unlisted `gpt-7-luna` shim literal passes the three release guards (the separate sentinel shim test covers a hardcoded shim).
- Proposal: **evidence defect / deferred risk**, not a current pin failure. Exact trigger is a future, presently unsupported model id inlined at a site without its own behavioral oracle; all currently promised lane ids are absent and known-id injection fails. Promote to material when such a new id becomes supported or this guard is offered as universal arbitrary-model coverage. Narrow the overbroad test comment or extend independent coverage when adding pins. FO disposition pending.

### V4 — declared surface tolerance exceeded

- Evidence: `git diff --numstat main...HEAD` = 13 files, +327/-38, net +289; approved estimate 6 files/+70, tolerance at most 8 files and net +30..+140. The sites are related to the design, but that does not waive the explicit tolerance.
- Proposal: **evidence/authorization defect / material gate hold, Needs decision** under `contract[docs/dev/README.md#validation]`: captain-visible scope/design-reset disposition is required before another pass beyond declared tolerance. Ownership is captain/FO, not a validator-authorized implementation expansion. Reduce surface or obtain the captain's explicit revised tolerance; do not silently declare the extra surface “forced.”

## Stage Report: validation

- DONE: Verify the single-source claim: run the print command yourself and report its exact output; confirm each consumer uses it, and confirm that no model literal remains in .github/workflows/runtime-live-e2e.yml. Name each workflow site you checked.
  `go run ./cmd/spacedock live-models` exited 0 with the exact stdout below; additionally executed the actual workflow resolver shell and independently asserted all five GITHUB_OUTPUT records.
```text
claude.sonnet=claude-sonnet-5
claude.opus=claude-opus-4-8
codex.exec=gpt-6-luna
pi.oauth=openai-codex/gpt-6-luna:max
pi.api-key=openai/gpt-6-luna:max
```
  Checked offline step/job outputs (67–93); Claude cadence matrix (122), model env (307), summary (281), artifact/config paths (139–141,293,317–322); Codex env export (376) and all three exec prefixes (388,392,397); both Pi summary branches (801,803); journey-delta `needs` and artifact lookup (868,899). Workflow search for Claude/GPT numeric ids returned no matches. Docs fail the consumer claim (V1).
- DONE: Verify the API-key spelling claim independently. The worker says it verified openai/gpt-6-luna:max against the installed pi-ai catalog rather than assuming it from the OAuth spelling. Check that against the catalog yourself and report what you found.
  Installed `@earendil-works/pi-coding-agent` and nested `pi-ai` are both 1.0.2. Under `/Users/clkao/.local/share/fnm/node-versions/v24.13.1/installation/lib/node_modules/@earendil-works/pi-coding-agent/node_modules/@earendil-works/pi-ai/dist/providers/data/`, `openai.json` independently declares id `gpt-6-luna`, provider `openai`, API `openai-responses`, and `thinkingLevelMap.max="max"`; `openai-codex.json` separately declares the OAuth provider/id/max. Catalog validation proves spelling/support metadata, not account entitlement or a successful API call.
- DONE: Confirm the Pi harness constants now alias the single source instead of repeating literals, and that the fixtures consume it.
  `pi_liveenv.go:19–20` aliases `release.PiOAuthModel/PiAPIKeyModel`; auth selection, local OAuth fallback (`pi_live_runner_test.go:250`), and metrics fixture (`journey_metrics_live_test.go:21`) consume them. Claude accepted-model/context-budget inputs use release constants; Codex argv fixture injects an independent sentinel. Recorded historical streams are untouched.
- DONE: Judge the Claude part: claude-sonnet-5.5 is refused by the installed CLI and is recorded as pending. Decide whether keeping claude-sonnet-5 and recording the 5.5 id as a captain input is the right disposition.
  Correct under the explicit design: retain the confirmed Sonnet 5 pin, record 5.5 as captain input, and do not assert global unavailability from one account/CLI rejection. No new Claude/API/live probe was launched.
- DONE: Audit for tautological tests in the captain's sense: no test may assert a file's own text against a copy of it, and no expected value may come from the thing it tests. The workflow no-literal guard deserves particular scrutiny under this rule.
  CLI expected stdout and no-literal forbidden set are independently authored; resolver-key compatibility compares separate sources; shim tests execute extracted code against independent argv expectations, not copied workflow text. Constants used as passthrough inputs are not a model-pin oracle. Existing Pi auth-selection comparisons reuse production aliases, so they establish branch selection only, not independent pin values; CLI exact output provides that independent pin proof. Wiring presence checks are not tautologies but have the identity hole V2; V3 records finite-list limits.
- DONE: Judge the recorded live evidence honestly: one smoke finished ok with root and child on openai-codex/gpt-6-luna, but the box was saturated, so the worker recorded the live proof as pending. Decide whether that framing is right.
  Read `/tmp/spacedock-pin-lane-smoke-1791093534/pi-frontdoor-smoke/run/pi-ensign-boot-grade.json`: verdict pass, both models correct, root 179506ms/child 77719ms, 5 reads, skill rank 2. This is valid model/boot evidence, not quiet-machine timing evidence; pending acceptance is right under FO steering. FO owns the quiet-machine rerun.
- DONE: Run focused checks and perform the semantic adversarial pass.
  CLI `TestLiveModelsCommand*`, Pi front-door/config/package-discovery fixtures, three release guards, Codex/Claude shim tests, Pi auth/env tests, context-limit table, and the live-tag metrics seed all passed. Pin changes break exact CLI stdout; wrong shim forwarding breaks independent sentinel argv; broken env/auth selection fails fixture state. Focused `-race` across the four affected packages also passed; no host was launched by the metrics-only seed.
  Disposable-archive mutation matrix: known literal -> guard fails; unknown key -> guard fails; Sonnet/Opus key swap -> guards incorrectly pass (V2); future unlisted literal -> release guards pass (V3). Actual resolver output is correct. Unknown CLI key -> exit 2/empty stdout; positional arg -> exit 2; explicit empty --get -> full list. Fixed five-entry lookup has no new scaling/hot-path concern.
  `gofmt -l` on changed Go files and `git diff --check main...HEAD` were clean. Candidate code/HEAD untouched; temporary archives removed.
- SKIPPED: Repository-wide offline suite and new live run.
  Explicit FO prohibition: only focused tests were run; full-suite AC-4 proof is not newly established here. Live AC-2 acceptance remains pending the separately owned quiet-machine rerun, not failed or claimed passed.
- FAILED: Report each finding with the defect kind and release scope, and whether the change can proceed.
  **AC-1 (VALUE): REJECTED** (V1 outcome/material, V2 evidence/material). **AC-2: pins/catalog/argv PASS; quiet-machine live proof PENDING. AC-3: current workflow PASS**, independently guarded (V3 deferred limitation). **AC-4: focused no-regression PASS, repository-wide suite NOT RUN**. Surface V4 is a material authorization hold; no change may proceed to acceptance yet.

### Summary

**REJECTED / do not proceed yet.** The runtime wiring, printed values, provider spellings and Claude disposition are correct, and the no-literal guard is independent rather than tautological. Documentation still violates single-source intent, the identity-wiring test accepts swapped lanes, and the 13-file/+289 surface exceeds approved tolerance; obtain FO dispositions and captain scope resolution before a focused follow-up, while FO supplies the quiet-machine evidence separately.

## Stage Report: implementation (validation-finding fix pass)

FO-authorized narrow fix of the four findings in `Stage Report: validation` against `fbb6cc8df`. Code commit `207b2d980`; no new product/architecture decision.

- DONE: Remove every lane model id from `docs/runtime-live-ci.md` and point at the print command.
  The prose paragraphs that named the OAuth/API-key Pi ids, the Sonnet/Opus ids, the Codex id, and the pending Claude candidate now say only how to obtain them: `spacedock live-models`, `spacedock live-models --get <lane>.<auth>`. The pending Claude candidate and its CLI rejection are recorded here in the entity, not in the doc. Evidence: `grep -nE 'claude-sonnet-5|claude-opus-4-8|gpt-6-luna|gpt-5.6-luna|claude-sonnet-5\.5' docs/runtime-live-ci.md` = no match. Illustrative `SPACEDOCK_PI_LIVE_CHILD_MODEL` override examples for third-party providers (not lane pins, not named in the finding) are left untouched.
- DONE: Strengthen the guards so a wrong value fails, using independent sources where one exists.
  Added `TestRuntimeLiveWorkflowResolverBindsEachOutputToItsKey`: it parses the `Resolve live lane models` step and requires each published output to be filled from its OWN `--get` key, against an independently authored output->key map. This closes the V2 hole: exchanging the resolver's `--get claude.sonnet`/`--get claude.opus` now fails (verified by mutation). Added `TestPiLaneModelsExistInInstalledCatalog`, the independent value oracle for `pi.oauth`/`pi.api-key`, which splits each pinned id into provider/model/thinking and requires the installed pi-ai catalog to declare that model with that thinking level (verified by mutating `PiAPIKeyModel` to a non-existent id, which failed). The test skips when no pi catalog is installed (host state). Lanes with no independent oracle are named plainly in the test comment rather than claimed: `codex.exec` (no installed registry declares Codex `exec --model` ids) and `claude.sonnet`/`claude.opus` (the installed Claude CLI validates a model only after auth; an isolated home short-circuits with "Not logged in" before validation, and real credentials would spend an API call, so no CLI probe is deterministic offline). Their independent oracle remains the authored exact-output CLI test plus recorded rejection evidence.
- DONE: Replace the finite forbidden-id list with a rule that does not need enumeration.
  `TestRuntimeLiveWorkflowCarriesNoLaneModelLiteral` (the finite id list) is replaced by `TestRuntimeLiveWorkflowModelSitesResolveThroughPrinter`, which keys off site shape: every `--model` argument must be a shell expansion, and every matrix `model:` / summary `Model:` value must be a `${{ ... }}` expression. A literal model id at any of those sites fails, including a future id never listed. Verified by injecting `--model gpt-6-luna` and `model: claude-opus-4-8`, each of which failed. Residual deferred risk recorded with its exact trigger: a model literal inlined at a *different* workflow site shape (e.g. a new `SPACEDOCK_*_MODEL:` env assignment) still passes; the site list must be extended when such a site is added.
- RECORDED (Finding 4, authorization) - the surface is not reduced.
  As validated: `git diff --numstat main...HEAD` at `fbb6cc8df` = 13 files, +327/-38, net +289, against an approved 8 files / +140. After this fix pass the branch is 13 files, +511/-39, net +472 (the increase is the doc rewrite and the strengthened guards). No file or scope was removed. This remains a captain-visible authorization record; the extra surface is not self-authorized.

### Verification (this pass)

- `go test ./internal/release -run 'TestRuntimeLiveWorkflow(ModelSitesResolveThroughPrinter|ResolverBindsEachOutputToItsKey|ResolvedKeysExist|LaneModelWiring)|TestPiLaneModelsExistInInstalledCatalog' -count=1` -> PASS.
- `go test ./internal/cli -run 'TestLiveModelsCommand' -count=1` -> PASS.
- Mutation matrix, each run and reverted: resolver `--get` swap -> `TestRuntimeLiveWorkflowResolverBindsEachOutputToItsKey` FAIL; `--model gpt-6-luna` literal -> site rule FAIL; `model: claude-opus-4-8` literal -> site rule FAIL; `PiAPIKeyModel` = non-existent id -> catalog oracle FAIL; Sonnet/Opus constant swap -> CLI exact-output oracle FAIL.
- Boundary: every test reads a real repo file or the installed catalog and states the reason it can fail; no expected value comes from the thing it tests; no file is asserted against a copy of itself. The repository-wide suite was not run (FO prohibition) and no CI lane was started.

### Summary

All four material/deferred/authorization findings are addressed narrowly: the doc restates no id, the resolver identity and Pi values now have failing oracles, the forbidden list is replaced by a site-shape rule with the residual risk recorded, and the surface deviation is recorded for the captain. Focused tests and the mutation matrix pass; the branch surface grew but was not reduced.

## Review-finding disposition (validation cycle 2)

Observations against `207b2d980051f18e71ef0e2c124ed359465b9f8a`; advisory classifications only, not authorization to modify the candidate. All mutations ran in disposable archives beneath the assigned worktree, restored individually and removed. Candidate bytes and HEAD are unchanged.

### R2-V2 — the resolver map is tautological; lane identity still has no behavioral guard

Released user / normal workflow: maintainers editing the resolver and the routine-Sonnet versus pre-release-Opus matrix. Observable harm: the asserted proof accepts swapped cadence consumers; it cannot establish which model actually reaches a lane.
Authority: value-ac[AC-1] The printed model list must be what each lane uses, with a valid test that fails on divergence.
Trigger evidence: `internal/release/livemodels_test.go:68–98` extracts the workflow's output/key spellings and compares them to the same output/key mapping authored in `want`; these private output names have no independent oracle. It catches the requested `--get` swap, but is a restatement of the workflow in the captain's sense and must be rejected. Swapping only the two `needs.offline.outputs.claude_*` operands on workflow line 122 passes all five release guards, both CLI tests, and both shim tests. Appending `echo "claude_sonnet=gpt-7-luna"` to the resolver also passes those checks: the regex ignores that write, so its map does not establish the complete output record.
Proposal: **evidence defect / material**, task-owned; **do not proceed** on this evidence. The unchanged candidate resolves correctly when its actual shell is executed. Replace the restatement with an independently supplied behavioral lane-identity proof, not another copied implementation map; FO authorization required. The CLI exact-output tests pin contract values, not workflow wiring: their comments claiming workflow changes are covered also overstate what they execute.

### R2-V3 — the site-shape residual is incomplete even at existing sites

Released user / normal workflow: maintainers editing the existing Claude matrix or Codex shim, not introducing a new workflow site. Observable harm: an inlined model can escape the claimed single-source guard at a currently supported model site.
Authority: value-ac[AC-1] A current lane model must derive from the single source rather than a workflow literal; the claimed proof must detect divergence.
Trigger evidence: `livemodels_test.go:18–29,35–52,62–63` accepts any `--model` token beginning `"$` and any matching key line containing `${{` anywhere. Replacing the Sonnet operand of the existing matrix expression with `'gpt-7-luna'` passes release, CLI and shim checks. `model: gpt-7-luna # <original expression>` passes all release guards. `--model "$(printf gpt-7-luna)"` and the existing `live_model="gpt-7-luna"` assignment pass release guards but fail the independent Codex sentinel shim test. A matrix expression containing only the literal passes the site guard but is caught by the separate wiring-presence guard (missing Opus); this is not an all-guard escape.
Proposal: **evidence defect / material** for the present broad proof claim; **do not proceed** on it. The recorded different/new-site trigger is real but incomplete: expressions, comments, shell substitutions/defaults, assignments feeding existing sites, and unparsed/duplicate resolver writes also need boundaries. The unchanged workflow has zero current model-id literals. Genuinely future unhandled syntax remains a **deferred risk**, promoted when adopted at a supported model site or offered as universally covered; do not confuse that residual with the reproduced current-shape proof gaps.

### R2-V5 — no implemented provider oracle is not the same as none being reachable

Released user / normal workflow: maintainers assessing whether Codex/Claude pins have independently verified provider metadata. Observable harm: the test/report incorrectly rule out a reachable independent Codex source and blur contract regression checks with provider validity.
Authority: contract[docs/dev/README.md#validation] State the actual observation boundary rather than claiming coverage or impossibility without evidence.
Trigger evidence: `livemodels_test.go:173–185` correctly has no implemented provider-backed oracle for Codex or either Claude pin, but its Codex rationale is too strong. The named installed file `/Users/clkao/.codex/models_cache.json` (fetched `2026-10-04T06:27:41.390507Z`, client `0.158.0`) declares `slug=gpt-6-luna`, `visibility=list`, and max reasoning; the same-form typo `gpt-6-luan` is absent. `codex exec --help` declares `--model <MODEL>`. This is independent installed model metadata, not proof of account entitlement or successful exec. No credentials were read and no API call was made. The recorded smoke is Pi, not independent Codex-exec evidence. No deterministic provider-backed Claude check was supplied or run; the pending candidate rejection does not validate both retained ids.
Proposal: **evidence defect / deferred risk**, task-owned reporting correction; not an additional current-pin blocker because the reachable metadata confirms the current Codex id and contract-value assertions catch drift. Promote if an unverified changed pin is accepted as provider-valid on this rationale. State that these lanes lack an implemented provider-backed oracle, not that no independent Codex metadata exists. The authored CLI output test is a contract oracle, not a provider or workflow oracle.

### R2-V4 — surface authorization remains unresolved

Released user / normal workflow: captain approving the declared implementation tolerance. Observable harm: proceeding would bypass the explicit surface/design gate.
Authority: contract[docs/dev/README.md#validation] Surface beyond declared tolerance needs a captain-visible design-reset disposition, not self-authorization by implementation.
Trigger evidence: `git diff --numstat main...HEAD` independently totals **13 files, +511/-39, net +472**, versus estimate 6 files/+70 and tolerance at most 8 files/net +30..+140; the previous branch was 13 files/net +289. Nothing was reduced.
Proposal: **evidence/authorization defect / material gate hold, Needs decision**, captain/FO-owned; **do not proceed** without the captain's explicit disposition. This validator neither reduces nor authorizes the surface.

## Stage Report: validation (cycle 2)

- DONE: Reproduce the mutation matrix yourself, one at a time, reverting each: swap the resolver's --get values; put a literal after --model; put a literal after model:; introduce a typo in the Pi API-key constant; swap the Sonnet and Opus constants. Report which guard fails for each, and name any mutation that no guard catches.
  All five independently reproduced on disposable `git archive HEAD` copies under the assigned worktree; each source restored before the next mutation. None of the five escapes all targeted guards; this does not make the resolver-map proof non-tautological.
- DONE: Swap the resolver's --get values.
  `--get claude.sonnet` exchanged with `--get claude.opus`: `TestRuntimeLiveWorkflowResolverBindsEachOutputToItsKey` exits 1 with both mismatched output/key pairs; other release guards pass. R2-V2 rejects this copied-map oracle despite that failure.
- DONE: Put a literal after --model.
  First `--model "$live_model"` changed to `--model gpt-7-luna`, an unlisted id: `TestRuntimeLiveWorkflowModelSitesResolveThroughPrinter` fails; `TestCodexLiveWorkflowPinsOnlyExecToLuna` independently fails its executed argv sentinel comparison.
- DONE: Put a literal after model:.
  Matrix value replaced with `model: gpt-7-luna`: site-shape guard fails; `TestRuntimeLiveWorkflowLaneModelWiring` also fails because the sole Opus consumer disappears.
- DONE: Introduce a typo in the Pi API-key constant.
  `PiAPIKeyModel="openai/gpt-6-luan:max"`: catalog test fails specifically “not declared by openai.json with thinking max,” not parsing/string form; CLI exact-output test also fails. Provider, slash and thinking suffix remain well-formed.
- DONE: Swap the Sonnet and Opus constants.
  Constants exchanged: `TestLiveModelsCommandPrintsPinnedLaneModels` fails with reversed Claude ids; all release guards pass, accurately showing that the catalog test covers only Pi.
- DONE: Verify the documentation claim by searching docs/runtime-live-ci.md for every lane id and the pending Claude candidate, and report the match count.
  File-bounded literal occurrence counts: `claude-sonnet-5=0`, `claude-opus-4-8=0`, `gpt-6-luna=0`, `openai-codex/gpt-6-luna:max=0`, `openai/gpt-6-luna:max=0`, pending `claude-sonnet-5.5=0`; old `gpt-5.6-luna=0` too. V1 is resolved: docs direct readers to printer keys; candidate/rejection remain in this entity.
- DONE: Verify the catalog oracle independently: does it interrogate the installed pi-ai catalog, and does it fail for a wrong model id rather than for a wrong string form?
  PASS without skip: the test follows installed `pi` to nested `@earendil-works/pi-ai/dist/providers/data`; independently read the installed 1.0.2 `openai.json` and `openai-codex.json`, both declare `gpt-6-luna` and `thinkingLevelMap.max="max"`. The same-form typo fails membership. Metadata proves support spelling, not entitlement, execution, or that any other catalog-valid model is the intended pin; absent catalog means SKIP, including an offline image without Pi.
- DONE: Judge the site-shape rule on its own terms: does it catch a model id it never listed, and at which site shapes does it have false negatives? Decide whether the recorded residual trigger is accurate and complete.
  Unlisted `gpt-7-luna` is caught at direct flag/key sites. Extra probes establish expression/comment/substitution/assignment false negatives at existing sites (R2-V3); the “different site shape” residual is accurate as one example but not complete. Matrix literal operand and matrix consumer swap pass release+CLI+shim guards; duplicate literal resolver write does too. A first all-literal matrix-expression probe was caught by wiring presence, so that broader escape expectation was corrected rather than reported passed.
- DONE: Check the honest boundaries: the test comment names codex.exec and both Claude ids as lanes with no independent oracle. Confirm that is accurate, and that no guard claims coverage it does not have.
  Accurate only as “no implemented provider-backed oracle”; Codex's “none in reach/no installed registry” rationale is contradicted by the independently probed installed model cache (R2-V5). CLI contract-value assertions are not provider validity or workflow identity tests; Pi smoke does not prove Codex exec; Claude rejection evidence does not verify both retained pins. Site/wiring/CLI comments overclaim coverage (R2-V2/R2-V3).
- DONE: Audit for tautological tests, and scrutinise the resolver output-to-key map especially: is it independent of the workflow, or a restatement of it?
  **REJECT the resolver map**: duplicating internal output names and their workflow key mapping in `want` is a restatement, not an independent oracle; failure under a text swap does not cure that. Pi catalog membership and executed sentinel argv are independent; CLI fixed ids can regress the specified pins but cannot establish consumer wiring. No additional copied-map guard was added.
- DONE: Reproduce AC evidence and run focused checks.
  Actual resolver shell executed successfully; output was `claude_sonnet=claude-sonnet-5`, `claude_opus=claude-opus-4-8`, `codex_exec=gpt-6-luna`, `pi_oauth=openai-codex/gpt-6-luna:max`, `pi_api_key=openai/gpt-6-luna:max` (one record per line). `go run ./cmd/spacedock live-models` printed the same ids under dotted keys in that order; each `--get` returned its id, unknown `pi.nope` exited 2 with empty stdout. Pi harness aliases and auth/fixture consumers use release constants; retained Claude pin and historical evidence disposition remain correct.
- DONE: Run targeted guard, regression and hygiene checks without a repository-wide suite.
  `go test ./internal/release -run 'TestRuntimeLiveWorkflow(ModelSitesResolveThroughPrinter|ResolverBindsEachOutputToItsKey|ResolvedKeysExist|LaneModelWiring)|TestPiLaneModelsExistInInstalledCatalog' -count=1 -v` and `go test ./internal/cli -run '^TestLiveModelsCommand' -count=1 -v` pass baseline. Focused `go test -race ./internal/release ./internal/cli ./internal/ensigncycle ./internal/claudeteam -run '^(TestRuntimeLiveWorkflow(ModelSitesResolveThroughPrinter|ResolverBindsEachOutputToItsKey|ResolvedKeysExist|LaneModelWiring)|TestPiLaneModelsExistInInstalledCatalog|TestLiveModelsCommand.*|TestCodexLiveWorkflowPinsOnlyExecToLuna|TestRuntimeLiveClaudeShimSetsMaximumEffort|TestPiLiveAuth.*|TestPiLiveEnv.*|TestPiIntercomPackageRootDefaultsBesideSubagents|TestContextLimitForModelBoundary)$' -count=1` passes. Auth tests exercise selection/seeding; env/discovery fixtures exercise scrubbing/package paths; shim sentinel catches forwarding drift; context table catches changed limits. `git diff --check main...HEAD` and changed-Go `gofmt -l` are clean; workflow numeric Claude/GPT-id search yields zero matches.
- DONE: Record the surface deviation for the captain and do not reduce it.
  Recomputed **13 files, +511/-39, net +472** versus **8 files/+140 maximum** (estimate 6/+70); R2-V4 remains a material authorization hold. No candidate files, tests, branch or HEAD were changed; disposable mutation archives were removed.
- SKIPPED: Repository-wide suite, new live smoke and CI lane run.
  Explicit FO prohibition/ownership: no `go test ./...`, repository-wide race suite or new host/CI lane was launched. Existing live acceptance remains separately owned; this round adds no quiet-machine proof and does not claim AC-4's full offline-suite clause.
- DONE: Report each finding with the defect kind and release scope, and whether the change can proceed.
  **REJECTED / do not proceed. AC-1 (VALUE): REJECTED proof** (V1 fixed, R2-V2 evidence/material, R2-V3 evidence/material). **AC-2: constants/catalog/argv PASS, new live proof NOT RUN. AC-3: candidate has no model literal; broad guard proof insufficient. AC-4: focused PASS, full suite NOT RUN.** R2-V5 is evidence/deferred risk, not an additional pin blocker; R2-V4 is material authorization hold/Needs decision. Recommendations do not authorize a repair cycle.

### Summary

All five required planted mutations fail a guard, documentation is now free of current/pending ids, and the Pi oracle genuinely consults installed catalog values. Nevertheless the resolver map is a forbidden restatement and additional current-site/cadence mutations pass the guards; **do not proceed**, and keep the unreduced 13-file/net +472 surface before the captain. Only this state report is committed; no candidate change, full suite, or new live run was made.

## Stage Report: implementation (validation cycle 3)

FO/captain fix round against `207b2d980`: delete-not-chase. Code commit `829e6c049`; no new product/architecture decision.

- DONE: Delete the tautological resolver output-to-key map. Do not replace it with a stronger version of the same comparison; prefer deleting a guard over chasing one.
  Removed `TestRuntimeLiveWorkflowResolverBindsEachOutputToItsKey` and its `resolverLinePattern`. No output-to-key comparison remains anywhere; the resolver is covered only by a structural rule (each `echo` line's exported value must come from `live-models --get`), which asserts no key mapping. `internal/release/livemodels_test.go`.
- DONE: Keep only guards with an independent source. The Pi catalog oracle qualifies; Codex metadata is independently available even though the test comment claims it is not — use it and correct that comment.
  Kept `TestPiLaneModelsExistInInstalledCatalog` (installed pi-ai catalog); added `TestCodexLaneModelExistsInInstalledCache`, reading `${CODEX_HOME:-$HOME/.codex}/models_cache.json`, which independently declares `gpt-6-luna` with the `max` level. Corrected the Pi test's coverage comment, which had wrongly listed `codex.exec` as "none in reach." Falsifiers run: `CodexExecModel="gpt-6-luan"` fails the cache oracle ("not declared"); `PiAPIKeyModel="openai/gpt-6-luan:max"` fails the catalog oracle ("not declared by openai.json with thinking max").
- DONE: Keep the site-shape rule only if it is genuinely structural (a literal cannot satisfy the required shape); extend it to every model-bearing site and name every site it cannot cover with the reason.
  Now genuinely structural at each covered site: `--model` must be exactly `"$live_model"`; `live_model=` must read `${SPACEDOCK_LIVE_CODEX_MODEL...}`; `model:`/`Model:` lines must contain a `${{ ... }}` reference with no model-id token before it (a YAML comment cannot supply it); `SPACEDOCK_LIVE_CODEX_MODEL=`/`SPACEDOCK_LIVE_MODEL:` env lines must reference the resolver output/`matrix.model`; resolver `echo` lines must resolve through `live-models --get`. Reproduced escapes now all fail: `--model "$(printf gpt-7-luna)"`, `--model gpt-7-luna`, `model: gpt-7-luna`, `model: gpt-7-luna # ${{ ... }}`, `live_model="gpt-7-luna"`, literal env lines, and a duplicate literal resolver write. Named uncovered in the test comment, with reasons: a literal nested inside a `${{ ... }}` operand; a literal appended to a derived string already interpolating `${{ needs.offline.outputs.* }}` (e.g. the journey-delta artifact name); any site shape the workflow does not use today.
- DONE: Record the honest coverage outcome per model-bearing site; no guard and no doc sentence may claim more than the tests do.
  Pi values = independent (installed pi-ai catalog). Codex value = independent (installed Codex cache). Claude values = authored contract only, no provider oracle (`TestLiveModelsCommandPrintsPinnedLaneModels`). Workflow literal placement at flag/key/env/resolver/offline sites = structural (`TestRuntimeLiveWorkflowModelSitesResolveThroughPrinter`, `TestRuntimeLiveWorkflowLaneModelWiring`). Workflow `--get` key existence = independent against the Go source (`TestRuntimeLiveWorkflowResolvedKeysExist`). NOTHING covers lane identity — which offline output feeds which matrix cadence lane: the deleted map was tautological, and exchanging the two `needs.offline.outputs.claude_*` operands on matrix line 122 still passes every guard (reproduced). Corrected the CLI test comment that implied the print test exercises workflow wiring; it pins contract values only. No doc sentence claims test coverage.
- DONE: Report the newly measured surface and the remaining deviation for the captain; do not reduce code arbitrarily to fit a number.
  `git diff --numstat main...HEAD` = 13 files, +618/-39, net +579 (was 13 files/net +472 at `207b2d980`). It did not shrink: the deleted resolver map removed ~40 test lines, but the required Codex oracle and the extended site coverage add more (`internal/release/livemodels_test.go` 265→370). Ideation estimate 6 files/+70, tolerance ≤8 files/net +30..+140: still far over, reported not self-authorized.
- SKIPPED: Repository-wide suite, any CI lane run, and a quiet-machine live rerun (FO prohibition / separately owned).
  Ran only focused packages: `go test ./internal/release ./internal/cli ./internal/ensigncycle ./internal/claudeteam -count=1` and `go test -race ./internal/release -run 'TestRuntimeLiveWorkflow|TestPiLane|TestCodexLane'`. Two CLI failures (`TestCodexResolveManifestAgainstInstalledHost`, `TestVersionAmbiguousMarkersExitZero`) are pre-existing host/ambient-env failures unrelated to this change (installed codex plugin state; ambient `PI_CODING_AGENT` marker).

### Summary

The tautological resolver output-to-key map is deleted, not replaced. Guards now rest on independent sources where available — the installed pi-ai catalog for Pi and the installed Codex model cache for Codex — and the workflow's model placement is guarded structurally at every site it uses, with the un-coverable sites named. Honest coverage is recorded per site, including the one site with nothing (matrix lane identity, the hole the deleted map used to hide). The measured surface grew rather than shrank because the Codex oracle and the extended structural coverage outweigh the deleted map; that deviation is reported for the captain rather than reduced to fit a number.

## Review-finding disposition (validation cycle 3)

Observations against `829e6c049`; recommendations only, not authorization to repair. All planted mutations ran one at a time in a disposable `git archive HEAD` beneath the assigned code worktree and were restored; candidate bytes and HEAD remain unchanged.

### R3-V2 — explicit map deleted, but LaneModelWiring remains a naming-convention oracle

Released user / normal workflow: maintainers editing the resolver, offline outputs and Claude cadence consumers. Observable harm: a green guard does not establish producer existence or model delivery, and rejects harmless producer aliases.
Authority: value-ac[AC-1] The printed model list must reach the intended lane, with a valid test that fails on divergence.
Trigger evidence: the explicit output-to-`--get` map and comparison are gone; `livemodels_test.go:173–179` instead builds `want` from the workflow's own output name and requires the same-named step output. Deleting only the `claude_opus` resolver echo passes all focused checks. Renaming that echo's output to `opus_alias` and its job-output reference to `steps.live_models.outputs.opus_alias` preserves delivery but fails LaneModelWiring. Swapping the matrix's two Claude operands, and separately swapping the resolver's two Claude keys, passes all focused checks.
Proposal: **evidence defect / material**, task-owned proof gap; **Needs captain decision**, not an automatic new guard. Derivation removes the enumerated map, not the authored same-name assumption: it moves that restatement into string construction. This is a falsifiable naming check, not a mathematically always-true assertion, but not independent behavioral evidence under the captain's rule. The matrix-gap honesty claim is TRUE; acknowledging the gap does not satisfy unchanged AC-1. Recommend a scope/design disposition rather than another copied-map repair.

### R3-V3 — literal rejection works, but universal structural coverage is still false

Released user / normal workflow: maintainers editing the existing Claude model env or existing resolver shell. Observable harm: literal configuration at promised covered sites passes every focused guard while supplying a wrong model.
Authority: value-ac[AC-1] Current lane models must derive from the single source rather than workflow literals.
Trigger evidence: all 20 direct literal substitutions at listed flag/assignment/key/env/resolver/job-output sites fail; nevertheless `SPACEDOCK_LIVE_MODEL: gpt-7-luna # ${{ matrix.model }}` passes every focused check and YAML decoding yields exactly `gpt-7-luna`. Appending `echo 'claude_sonnet=gpt-7-luna'` inside the existing resolver output block also passes; executing its shell produces a duplicate literal output record. Replacing the Sonnet echo with `echo "claude_sonnet=gpt-7-luna" # "live-models --get claude.sonnet"` passes and actual shell execution emits the wrong value. `livemodels_test.go:59–65,95–105,117–125` uses substring/quote matching, not a complete value-source rule. An entire journey artifact suffix changed to `gpt-7-luna` also passes; the comment only names appended literals, not removal of that site's reference.
Proposal: **evidence defect / material**, task-owned; **do not proceed on the all-sites claim**. The rule is partly structural (direct literal sites), not a copied model-value map, but comments/quoting still satisfy it without supplying the model. Existing acknowledged nested-expression escape reproduces too. A literal `live_model=` plus a reference only in its comment passes the release guards but FAILS the independent Codex shim sentinel test: this is explicitly not an all-guard escape. Future unused syntax is a deferred risk only until adopted; these reproduced env/resolver triggers are already within the promised shapes. Captain decides scope/design; this report authorizes no parser expansion or repair.

### R3-V4 — surface authorization remains unresolved

Released user / normal workflow: captain approving the implementation tolerance. Observable harm: accepting the candidate bypasses the explicit surface gate.
Authority: contract[docs/dev/README.md#validation] Exceeding declared tolerance requires captain-visible scope/design disposition.
Trigger evidence: `git diff --numstat main...HEAD` independently totals **13 files, +618/-39, net +579**, versus **at most 8 files/net +140** (estimate 6/+70), up from net +472. No reduction was made.
Proposal: **evidence/authorization defect / material hold, Needs decision**, captain/FO-owned; do not proceed without explicit disposition. New Codex metadata evidence resolves R2-V5's incorrect no-oracle rationale; it does not waive this gate.

## Stage Report: validation (cycle 3)

- DONE: Verify the deletion: the tautological resolver output-to-key map must be gone, and nothing may have replaced the output-to-key comparison. Report what covers that surface now.
  TRUE: inspected `207b2d980..829e6c049` and repo-bounded Go searches; deleted test/map/pattern absent, no replacement output-to-`--get` comparison. Remaining coverage is echo-string placement, key membership against Go, same-name job-output references and global consumer presence—not resolver identity or executed complete output records. Resolver key swap passes.
- DONE: Verify the Codex oracle independently: it reads the installed models_cache.json and requires the exec model to declare max. Confirm it fails for a wrong value rather than a wrong string form, and that the file it reads is independent of this repository.
  TRUE on this host: `/Users/clkao/.codex/models_cache.json`, outside the entire repository, fetched `2026-10-04T06:41:12.292460Z`, client `0.158.0`; baseline passes without skip. Same-form `CodexExecModel="gpt-6-luan"` fails “not declared”; installed valid slug `gpt-5.5` fails “does not declare the max reasoning level.” Independent membership/level evidence, not pin intent, entitlement or execution; CODEX_HOME can redirect the path and missing cache skips.
- DONE: Verify the structural site rule genuinely rejects a literal at each covered site, and that it is structural rather than a restatement. Plant a literal at one covered site and report the failure.
  PARTLY TRUE: independently planted `gpt-7-luna` at workflow lines 67–71,89–93,122,281,307,376,382,388,392,397,801,803; all 20 fail site-shape or job-output checks. Example line 388 fails “--model gpt-7-luna is not the resolved shell variable.” Not a copied model-id oracle, but universal structural claim FALSE: R3-V3's existing env/comment and resolver/quote escapes pass all focused guards; the literal really reaches parsed env or shell output. Artifact whole-reference replacement escapes too.
- DONE: Verify the honesty claim about the residual gap: swap the two needs.offline.outputs.claude_* operands and confirm that every guard still passes, as the report states. If any guard does catch it, the claim is wrong and you must say so.
  TRUE within all relevant focused guards: exchanged only matrix line 122's Sonnet/Opus operands; release workflow/provider checks, CLI output checks, both executed shim checks and Pi/context regressions all exit 0. No guard caught it. This confirms the honest residual, not AC-1; no repository-wide claim is made.
- DONE: Audit for tautological tests once more, including the rewritten LaneModelWiring helper: does deriving output names from the workflow remove the authored map, or move it?
  Explicit map removed; same-name mapping assumption moved into `want := "${{ steps.live_models.outputs." + output + " }}"`. It is structural convention evidence only, not independent delivery evidence: missing Opus producer passes, legitimate producer alias fails (R3-V2). Installed Pi/Codex catalogs and shim sentinels are independent; CLI fixed values are contract oracles, not workflow/provider oracles. `LiveModels()`'s “cannot disagree” comment also overstates what wiring proves.
- DONE: Record the surface deviation for the captain: 13 files, net +579 against an approved 8 files and +140. Do not reduce it.
  TRUE: recomputed 13 files, +618/-39, net +579; R3-V4 material authorization hold remains. Candidate code/HEAD unchanged; only this state report is committed.
- DONE: Reproduce acceptance evidence and run focused checks; report each finding with defect kind, release scope and proceed/decision disposition.
  Printer output exactly: `claude.sonnet=claude-sonnet-5; claude.opus=claude-opus-4-8; codex.exec=gpt-6-luna; pi.oauth=openai-codex/gpt-6-luna:max; pi.api-key=openai/gpt-6-luna:max` (one record per line). Executed actual resolver shell yields the same five values under underscored output names. Installed pi-ai `openai.json` and `openai-codex.json` independently declare `gpt-6-luna` and max; docs contain zero current/old/pending lane-id matches; Pi aliases and retained Claude disposition remain correct.
- DONE: Run targeted baseline, mutation and hygiene checks without widening validation scope.
  Ran `go test ./internal/release ./internal/cli ./internal/ensigncycle ./internal/claudeteam -run '^(TestRuntimeLiveWorkflow.*|TestPiLaneModelsExistInInstalledCatalog|TestCodexLaneModelExistsInInstalledCache|TestLiveModelsCommand.*|TestCodexLiveWorkflowPinsOnlyExecToLuna|TestRuntimeLiveClaudeShimSetsMaximumEffort|TestPiLiveAuth.*|TestPiLiveEnv.*|TestPiIntercomPackageRootDefaultsBesideSubagents|TestContextLimitForModelBoundary)$' -count=1 -v` and the same selector with `-race`: baseline PASS. Mutations ran the same selector or isolated site/oracle checks, uncached. Direct literals/typos fail as above; env/reference/producer escapes pass. `git diff --check main...HEAD` and changed-Go `gofmt -l` clean; disposable archive removed.
- SKIPPED: Repository-wide suite, CI lane run, and new live smoke.
  FO prohibition/separate ownership: no repository-wide command or host/CI run started. No new quiet-machine or full AC-4 offline-suite evidence is claimed; catalog metadata is not a live run.
- DONE: State the acceptance recommendation without authorizing another repair cycle.
  **REJECTED / Needs captain decision. AC-1 (VALUE): proof insufficient** (R3-V2/R3-V3 evidence/material). **AC-2: pins/catalog/argv pass; existing live acceptance remains pending. AC-3: actual candidate has no model-id literal; universal structural proof false. AC-4: focused regressions/race pass; full suite not run.** R3-V4 is evidence/authorization material hold. Deletion, Codex oracle, matrix-gap honesty and +579 surface claims are verified; the all-sites structural and independent LaneModelWiring claims are not.

### Summary

Third-round adversarial validation confirms the genuine deletion and external Codex value oracle, and confirms that swapping Claude cadence operands still passes. Direct literals fail, but current-site comment/quote mutations still evade the purported structural rule, and LaneModelWiring merely enforces a private naming convention; recommend captain scope/design disposition, not another automatic guard chase, with the unreduced 13-file/net +579 surface explicitly held.

## Stage Report: implementation (reduction round)

Directed reduction round against `829e6c049`. Code commit `87c6b2ac1`. FO direction: keep the single source, the print command, and the two independent provider oracles; delete every workflow-placement guard; record the coverage outcome; report the measured surface; change no pinned value or workflow behaviour.

- DONE: Keep the single source in `internal/release` and the print command that exposes it.
  `internal/release/livemodels.go` (the five constants, `LiveModels()`, `LiveModel()`) and `internal/cli/live_models.go` (`spacedock live-models`, `--get`) are untouched except for two comments that stopped claiming placement enforcement. Pinned values (`claude-sonnet-5`, `claude-opus-4-8`, `gpt-6-luna`, `openai-codex/gpt-6-luna:max`, `openai/gpt-6-luna:max`) unchanged; `git diff 829e6c049..87c6b2ac1 -- internal/release/livemodels.go` shows only comment edits.
- DONE: Keep the two independent provider oracles — the Pi catalog check and the Codex models_cache check.
  `TestPiLaneModelsExistInInstalledCatalog` (reads the installed pi-ai `openai.json`/`openai-codex.json`) and `TestCodexLaneModelExistsInInstalledCache` (reads `${CODEX_HOME:-$HOME/.codex}/models_cache.json`) are retained verbatim with their helpers. Both ran and PASSED on this host without skip: `go test ./internal/release -run 'TestPiLaneModelsExistInInstalledCatalog|TestCodexLaneModelExistsInInstalledCache' -count=1 -v`.
- DONE: Delete every workflow-placement guard; delete rather than strengthen; let nothing replace them.
  Removed from `internal/release/livemodels_test.go`: `TestRuntimeLiveWorkflowModelSitesResolveThroughPrinter` (the structural site rule) and its five regexes plus `stripLineComment`; `TestRuntimeLiveWorkflowLaneModelWiring`; `TestRuntimeLiveWorkflowResolvedKeysExist` (a workflow-vs-source `--get` key comparison). No replacement was added; `grep -nE 'TestRuntimeLiveWorkflow(ModelSitesResolveThroughPrinter|LaneModelWiring|ResolvedKeysExist)|modelFlagArgPattern|modelKeyLinePattern|modelIDTokenPattern|codexLiveModelAssignPattern|resolverEchoPattern|liveModelGetPattern|stripLineComment' internal/release/livemodels_test.go` = NONE.
- DONE: Record the coverage outcome plainly, with the exact residual and the reason placement cannot be enforced.
  Recorded in the `livemodels_test.go` header/`TestPiLaneModelsExistInInstalledCatalog` comment and in `livemodels.go`: model values have independent provider oracles for Pi (installed catalog) and Codex (installed cache); Claude ids rest on the authored exact-output print test only; the placement of any id in the workflow has no enforcement. Exact residual: a swap of the two `needs.offline.outputs.claude_*` operands passes every remaining check. Reason placement cannot be enforced by an independent oracle: both sides of that comparison are authored by this same change, so any placement test restates the change instead of checking it.
- DONE: Correct the test comments that claimed removed coverage. No test comment claims more than the recorded outcome.
  `internal/cli/live_models_test.go` now says the print test fixes contract values only and that workflow consumption has no test enforcement (was: "checked separately and structurally by the release workflow guards"); the `--get` comment no longer references a workflow rename. `internal/ensigncycle/codex_liveenv_test.go` shim-sentinel comment no longer cites "the workflow no-literal test"; it cites the print test and the installed-catalog oracles. `livemodels.go` `LiveModels()` no longer claims the workflow "cannot disagree".
- DONE: Report the newly measured surface without trimming to fit.
  `git diff --numstat main...HEAD` at `87c6b2ac1` = **13 files, +440/-39, net +401**, versus the approved estimate 6 files/+70 and tolerance at most 8 files/net +30..+140. It does NOT fit: file count is 13 (>8) and net is +401 (>+140). This pass reduced net by 178 lines (from +579 at `829e6c049`); the file count did not fall because the deleted guards lived inside a file that still holds the two kept oracles. Reported as measured; nothing further was trimmed to approach the number.
- DONE: Do not change the pinned values or the workflow's behaviour.
  `git diff 829e6c049..87c6b2ac1 -- .github/workflows/runtime-live-e2e.yml docs/runtime-live-ci.md` = empty (no workflow or docs change this pass). Constants unchanged (comment-only edit in `livemodels.go`).
- DONE: Verify the reduction with focused checks (no repository-wide suite, no CI lane).
  `go build ./...` = OK. `gofmt -l` on the four edited files = clean. `go test ./internal/release -count=1` = ok (13.4s). `go test ./internal/cli -run 'TestLiveModelsCommand' -count=1 -v` = PASS. `go test ./internal/ensigncycle -run 'TestCodexLiveWorkflowPinsOnlyExecToLuna|TestRuntimeLiveClaudeShimSetsMaximumEffort|TestContextLimitForModelBoundary' -count=1` = ok. `go test ./internal/claudeteam -run 'TestContextLimitForModelBoundary' -count=1` = ok. No `go test ./...` and no CI lane were run per FO direction.
- SKIPPED: Repository-wide offline suite, `-race` across all packages, and any live/CI lane run.
  FO prohibited the repository-wide suite and a CI lane; existing live acceptance remains separately owned. No quiet-machine live proof is claimed here.

### Summary

Reduced to the directed deliverable: the single `internal/release` source and the `spacedock live-models` print command, plus the two independent provider oracles (installed pi-ai catalog for Pi, installed Codex cache for Codex). Every workflow-placement guard is deleted — structural site rule, `LaneModelWiring`, and the workflow-vs-source `--get` key comparison — with no replacement. The coverage outcome is recorded plainly, including the exact residual (a swapped Claude output operand passes every remaining check) and why it cannot be closed by an independent oracle. Measured surface is 13 files, net +401, which does not fit the approved 8 files/net +140; it is reported rather than trimmed. No pinned value or workflow behaviour changed.

## Review-finding disposition (light reduction validation)

### R4-C1 — coverage wording exceeds the observed boundary

Released user / normal workflow: maintainers using the reduced checks and coverage record to judge lane-model changes. Observable harm: retained Claude ids appear supported by rejection evidence about a different candidate, and an absent placement oracle is presented as a proved impossibility.
Authority: contract[docs/dev/README.md#validation] Report the actual observation boundary; do not claim proof beyond what the checks establish.
Trigger evidence: `internal/release/livemodels_test.go:33–34` says the Claude ids rest on the authored oracle **and the recorded rejection evidence**; rejection of Sonnet 5.5 does not establish retained Sonnet 5 or Opus. Lines 39–41 say placement **cannot** have an independent oracle; `livemodels.go:38` generalizes to **any** placement comparison; reduction-report lines 356 and 370 repeat that impossibility. Deleting these guards proves absence, not impossibility. Reject those sentences, not the directed deletion.
Proposal: **evidence/reporting defect, material to the explicit no-overclaim checklist**, task-owned wording correction only; no guard, parser, architecture or scope change proposed. State simply: Pi/Codex have independent metadata oracles; Claude ids rest on the authored print test; workflow placement has no enforcement; swapping the two `needs.offline.outputs.claude_*` operands passes the remaining checks. The previously recorded surface decision remains captain-owned.

## Stage Report: validation (light reduction round)

- DONE: Run the grep yourself and confirm all three workflow-placement guards and their helpers are gone, with nothing replacing them.
  Repo-bounded `git grep -nE 'TestRuntimeLiveWorkflow(ModelSitesResolveThroughPrinter|LaneModelWiring|ResolvedKeysExist|ResolverBindsEachOutputToItsKey)|modelFlagArgPattern|modelKeyLinePattern|modelIDTokenPattern|codexLiveModelAssignPattern|resolverEchoPattern|liveModelGetPattern|resolverLinePattern|stripLineComment' -- .` returned exit 1 / zero matches; reviewed the entire `829e6c049..87c6b2ac1` diff and related test references: no replacement, only deletions/comment edits.
- DONE: Confirm the Pi catalog check and the Codex models_cache check pass and fail for a wrong value.
  Both pass without skips; independently read installed pi-ai `providers/data/{openai,openai-codex}.json` and `/Users/clkao/.codex/models_cache.json`: `gpt-6-luna`/max present, `gpt-6-luan` absent. One-at-a-time disposable-archive mutations of each Pi pin and the Codex pin to that same-form typo exit 1 specifically for missing metadata membership; source restored after each.
- DONE: Confirm the remaining oracles are genuinely independent.
  Tests read installed provider metadata rather than expected ids copied from the pin source; retained provider tests/helpers are byte-identical to `829e6c049`. This proves catalog/cache membership and reasoning metadata, not intended-pin identity, entitlement or successful execution; missing host metadata causes skips.
- DONE: Read the coverage comments and the report for overclaiming.
  REJECT the sentences identified in R4-C1: Claude retained ids have only the authored print oracle, not support from a different candidate's rejection; no placement enforcement exists, but impossibility of an independent oracle is unproved. Historical reports remain historical, not renewed claims. The requested reduction itself is accepted; no guard is proposed.
- DONE: Confirm the exact residual: a swapped needs.offline.outputs.claude_* operand passes every remaining check.
  Swapped only the two Claude matrix operands in a disposable archive; both provider tests, both print tests, Codex sentinel shim, Claude effort shim and context-boundary test all passed without skips. These are the remaining focused checks exercised, not a claim that the prohibited repository-wide suite ran; placement has no enforcement.
- DONE: Confirm the pinned values and the workflow's behaviour are unchanged by this pass.
  Workflow, print-command implementation and docs are byte-identical to `829e6c049`; non-comment production source is identical. `go run ./cmd/spacedock live-models` printed, in order: `claude.sonnet=claude-sonnet-5`, `claude.opus=claude-opus-4-8`, `codex.exec=gpt-6-luna`, `pi.oauth=openai-codex/gpt-6-luna:max`, `pi.api-key=openai/gpt-6-luna:max`. No numeric Claude/GPT model literal matches in the current workflow.
- DONE: Measure the surface yourself: files and net lines, against the approved 8 files and net +140. Report the numbers.
  Independently summed `git diff --numstat main...HEAD` at `87c6b2ac1`: **13 files, +440/-39, net +401**; exceeds limits by **5 files / 261 net lines**. Reduction alone: 4 files, +21/-199, net -178. No trimming performed or proposed.
- DONE: Run only focused checks and preserve the candidate.
  `go test ./internal/release ./internal/cli ./internal/ensigncycle ./internal/claudeteam -run '^(TestPiLaneModelsExistInInstalledCatalog|TestCodexLaneModelExistsInInstalledCache|TestLiveModelsCommand.*|TestCodexLiveWorkflowPinsOnlyExecToLuna|TestRuntimeLiveClaudeShimSetsMaximumEffort|TestContextLimitForModelBoundary)$' -count=1 -v` passed baseline and operand swap; wrong-value probes failed as intended. `git diff --check main...HEAD` and changed-file `gofmt -l` clean; disposable archive removed; code HEAD/worktree/index unchanged.
- DONE: State acceptance evidence at the reduced boundary.
  AC-1: source/print PASS, placement enforcement deliberately absent (original divergence clause not proved). AC-2: both Pi spellings/metadata PASS, no new live proof. AC-3: current workflow literal search clear, no future-placement enforcement. AC-4: reduction leaves runtime bytes unchanged and focused checks PASS; full offline suite not newly established.
- SKIPPED: Repository-wide suite and any live/CI lane run.
  Explicit FO prohibition on full suite/CI; no new host or quiet-machine smoke launched, and no new live acceptance claimed.
- DONE: Say whether the change can proceed, and name the single decision the captain owes.
  **Not ready for unconditional acceptance as written:** reduction and independent-oracle evidence PASS, but reject R4-C1's wording overclaims. After reporting-only correction, the sole captain decision is whether to authorize **13 files/net +401 instead of 8 files/net +140** for this already-directed reduced deliverable. No renewed guard-design decision or automatic guard repair is requested.

### Summary

The reduction is real, the two metadata oracles pass and reject wrong values, and the reproduced Claude operand swap confirms the deliberately unenforced placement boundary. Correct only the identified coverage overclaims; the captain still owes the single over-tolerance surface disposition (13 files/net +401), with no pinned-value or workflow-behaviour change required.

## Stage Report: implementation (wording repair)

FO-directed wording-only repair of the two R4-C1 coverage-record overclaims against `87c6b2ac1`. Code commit `5ff137093`. No guard, pinned value, or workflow change.

- DONE: Correct the overclaim that the Claude model rejection validates the retained Claude ids.
  `internal/release/livemodels_test.go:32–38` now reads "claude.sonnet, claude.opus: no independent oracle" and states that the recorded rejection "establishes only that one candidate id (claude-sonnet-5.5) is refused by the installed CLI; it does not validate the retained ids. The retained claude.sonnet and claude.opus ids rest only on the authored print test, with no independent oracle." Removed the phrase that had folded "the recorded rejection evidence" into the retained ids' support.
- DONE: Correct the overclaim that placement enforcement is impossible by an independent oracle.
  `internal/release/livemodels_test.go:40–44` now reads "Where the values land has no enforcement today" and "A comparison between the workflow and the source is not independent here, because this change authors both sides; that is an absence of enforcement today, not proof that none is possible." Removed "Placement cannot be checked by an independent oracle", which asserted impossibility from absence. The reproduced residual is unchanged: a swap of the two `needs.offline.outputs.claude_*` operands passes every remaining check.
- DONE: Change nothing else — no code, no guard, no pinned value, no workflow change.
  `git diff 87c6b2ac1..5ff137093 --stat` = 1 file, +12/-11, comment lines only (`internal/release/livemodels_test.go`); `.github/workflows/runtime-live-e2e.yml`, `internal/release/livemodels.go`, `internal/cli/live_models.go` and all pinned constants byte-identical. Rationale: the two sentences were the identified defect; no adjacent wording was edited.
- DONE: Verify the wording edit compiles and the retained oracles still pass and still reject a wrong value.
  `gofmt -l internal/release/livemodels_test.go` = clean. `go test ./internal/release -run 'TestPiLaneModelsExistInInstalledCatalog|TestCodexLaneModelExistsInInstalledCache' -count=1` = ok 0.344s; `go test ./internal/cli -run 'TestLiveModelsCommand' -count=1` = ok 0.258s. Both comments are non-executable; the tests exercising them passed, and neither wrong-value falsifier behavior changed (comment-only diff).
- DONE: Report the final measured surface, files and net lines, against the approved 8 files and net +140, and say plainly that the only remaining item is the captain's authorization of that number.
  `git diff --numstat main...HEAD` at `5ff137093` = **13 files, +441/-39, net +402** (was 13 files/net +401 at `87c6b2ac1`; the +1 is the net comment line). Versus the approved **at most 8 files / net +30..+140**: over by 5 files and 262 net lines. The only remaining item is the captain's authorization of this surface (13 files, net +402); no other decision is requested.
- SKIPPED: Repository-wide suite, `-race` across all packages, and any live/CI lane run.
  FO prohibition / separate ownership carried over from the reduction round. Comment-only change; focused package checks were run, and no quiet-machine or full offline-suite evidence is claimed.

### Summary

Corrected the two R4-C1 coverage-record sentences in `internal/release/livemodels_test.go`: the Claude ids now rest explicitly on the authored print test alone with no independent oracle (the rejection proves only that claude-sonnet-5.5 is refused), and placement is stated as having no enforcement today with the reason a workflow-vs-source comparison is not independent — rather than as an impossibility. Nothing else changed: a one-file comment edit. Final measured surface is 13 files, net +402, still over the approved 8 files/net +140; the sole remaining item is the captain's authorization of that number.
