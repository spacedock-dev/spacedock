---
title: Recognize native Pi worker completion evidence
status: implementation
source: Captain binding resolution:binding-1789576723761874000; tip CI 35058669297
started: 2026-09-16T17:00:12Z
completed:
verdict:
score: 0.95
worktree: .worktrees/spacedock-ensign-pi-native-completion-evidence
issue:
pr:
mod-block:
id: cpx5q07b7wqmd6xv6wc7khvd
gates:
    version: 1
    records:
        - id: gate:cpx5q07b7wqmd6xv6wc7khvd:validation
          stage: validation
          attempts:
            - id: gate-attempt:cpx5q07b7wqmd6xv6wc7khvd-validation-1
              briefing:
                id: briefing:cpx5q07b7wqmd6xv6wc7khvd:validation:attempt-1:revision-1
                digest: sha256:5ff3627d7887abcb5cd217b4fa43e91af634041417b147048da3308f562e96b1
                room-ref: '@review/validation/briefing-1'
              resolution:
                type: Resolution
                id: resolution:spacedock:cpx5q07b7wqmd6xv6wc7khvd:validation:1
                briefing: briefing:cpx5q07b7wqmd6xv6wc7khvd:validation:attempt-1:revision-1
                by: person:captain
                at: "2026-09-16T18:44:54.981178Z"
                decision: approve
                reason: Captain binding resolution:binding-1789584263853542000 approves validated candidate ba8e6eeb for PR delivery and final tip CI; merge remains unauthorized.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:cpx5q07b7wqmd6xv6wc7khvd-validation-2
              briefing:
                id: briefing:cpx5q07b7wqmd6xv6wc7khvd:validation:attempt-2:revision-1
                digest: sha256:7d3dcc2c5e07cd29556669324a638df3224bb4ec963cecf0b6b539ec1351f156
                room-ref: '@review/validation/briefing-2'
              resolution:
                type: Resolution
                id: resolution:spacedock:cpx5q07b7wqmd6xv6wc7khvd:validation:2
                briefing: briefing:cpx5q07b7wqmd6xv6wc7khvd:validation:attempt-2:revision-1
                by: person:captain
                at: "2026-09-17T15:04:54.800792Z"
                decision: approve
                reason: Captain binding resolution:binding-1789657416015652000 approves validation attempt 2 and combined tip 0997b6b7 for corrected-stack publication and full tip CI. Reviewed artifact sha256:090f392bf3872fe4169655453f82ff67e908d3c287051a65df70a7be0c2247ac. Merge and failure waiver remain unauthorized.
              application:
                target-stage: done
                state: superseded
---

Recognize Pi native completion notifications so valid worker lifecycles are graded correctly without accepting another worker or hiding downstream failures.

## Problem

Tip CI 35058669297 reports incomplete implementation/validation workers although retained parent and child evidence proves their identity, committed reports, completion, and notification before advancement. The existing grader recognizes status/wait responses only. Auto-continue separately fails gate preparation on a mistyped reference path; that recorded whole journey must remain red.

## Proposed approach

Extend the existing Pi lifecycle evidence path. Join the parent dispatch result run ID to the notification-referenced child session identity and original assignment. Credit completion at its position in the parent timeline. Keep existing report/gate assertions and historical completion forms. Reuse existing parsers and retained sessions. Establish the smallest concrete surface and a failing captured replay before candidate edits.

Captain approved this bounded follow-up, independent validation, and a PR above the current stack through binding resolution:binding-1789576723761874000 at 2026-09-16T16:38:43.761876Z. Approved artifact /private/tmp/pi-native-notification-fix-review.md, sha256:fb4c087fbde4cad12db9080553c5a85a84a7ada2f504da760151d8c03fbcc3d5. This seed records the approved brief; no earlier stage report is manufactured.

## Risk evidence

Existing read-only triage: ../fo-continuation/artifacts/validation/pi-tip-ci-triage/report.md and retained-events.json/source-manifest.json. Downloaded artifacts: /tmp/spacedock-tip-ci/pi. Exact parent toolCallId/result run ID joins child session_info run identity, assignment, successful commit and terminal event. The native notification precedes advancement. Child directory UUID differs from run ID, so generic labels or filename guessing are insufficient.

## Out of scope

Product runtime changes, process controllers, new CI lanes, new public commands, naming changes, weakening gate checks, or fixing a one-off path typo by adding a mechanism. No live model runs locally. Expensive verification stays at the final tip.

## Expected surface and tolerance

Use existing lifecycle helper, Pi evidence adapter, shared native-stream bridge, and existing test owners with a compact captured fixture. Before editing, report concrete paths and size plus the failing replay. The approved boundary is one narrow correlation path; return for design review if it requires a new framework or materially broader scope. Do not invent an arbitrary numerical cap absent from the approved brief.

## Acceptance criteria

**AC-1 — Native Pi completion credits only the dispatched worker at the correct boundary.**
Verified by: captured parent/child replay from CI credits default implementation and auto validation, while wrong run/assignment, missing or unsuccessful child, and out-of-order notification fail. Removing identity correlation must make a negative case pass incorrectly and be caught.

**AC-2 — Correct completion observation does not hide outcome failures or change other hosts.**
Verified by: the original auto-continue whole journey still fails because its gate was never prepared; historical Pi status/wait and existing Claude/Codex lifecycle controls continue to pass. Bypassing gate assertions fails the retained negative case.

**AC-3 — The repaired stack supplies honest final CI evidence.**
Verified by: deterministic checks locally and required tip PR host lanes after independent validation. Retain concrete failures and expected-failure policy; never call skipped/cancelled checks green. No local model or duplicate middle-layer live run.

## Test plan

Start with the exact captured failure replay through existing lifecycle test owners. Add narrowly targeted correlation/ordering negatives, preserving durable state assertions. Run go test ./..., go test ./... -race, gofmt -w ./cmd ./internal. Known local installed-manifest resolver failure requires exact-log confirmation; it never waives new failures. Independent validation must test attribution and incomplete-outcome boundaries. Full expensive host checks run at the stack tip only.

### Feedback Cycles


## Stage Report: implementation

- DONE: Credit the captured native Pi completions only after matching the dispatched worker, child evidence and parent event ordering; wrong or missing evidence fails.
  Commit `ba8e6eeb4e28e310b17897a0a8b12e151366b8f9`: captured default implementation and auto validation replay pass; 27 negatives reject run/agent/epoch/task/cwd, missing/error evidence, escaped/symlink locators and ordering/duplicate conflicts. Removing correlation admits wrong workers; ignoring conflicting epochs reproduced a red test before its authorized fix.
- DONE: Preserve downstream gate failure and historical Pi/Claude/Codex behavior through existing replay and negative controls.
  Full retained auto lifecycle now passes; captured absent-gate reconstruction still returns `entity has no gates record`. Existing Git check rejects an uncommitted report; bypassing either durable check breaks the replay. Historical Pi status/wait, Claude replay, Codex lifecycle and cross-host bypass controls pass (`artifacts/implementation/focused.log`).
- DONE: Deliver a minimal committed correction with required local checks and clear tip-CI deferral.
  Code committed locally on the registered branch; focused tests, live-tag compile-only check and formatting completed. Final normal/race ran sequentially; ensigncycle passed in 354.805s / 335.058s. Native host lanes remain explicitly deferred below.

### Summary

`go test ./...` and `go test ./... -race` did not pass. Both exit 1 only at the anticipated `TestCodexResolveManifestAgainstInstalledHost`, `codex_resolve_test.go:44`: `spacedock@spacedock not installed in codex, but resolver returned "/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json"`. All other packages pass; no race diagnostic. Exact final/initial logs are in `artifacts/implementation/`. The FO explicitly DECLINED the exact pre-existing installed-manifest resolver defect for this task, as anticipated in the assignment; it is not an unresolved correction finding. This disposition does not change either exit-1 result or waive new failures.

`go test -tags live ./internal/ensigncycle -run '^$'` passes; it selects no tests. Removing artifactDir propagation from either existing call site would leave that live journey unable to credit native notices.

Final tip PR/native host CI and independent validation remain deferred, not passed. FO owns these after implementation handback; no local native/model, network-auth, PR, push, CI or rebase was performed. The original auto gate-path typo is retained as an outcome failure, not fixed or declared green. AC-3 remains pending final tip CI.

One 150-line Pi helper correlates the parent dispatch result with the notification-referenced child’s exact assignment and run/epoch, and credits completion at the native parent position. Both existing callers pass optional retained-artifact roots; existing parsers and independent Git/report/gate checks remain authoritative, with no synthetic event, child-call concatenation or second commit parser.

Actual surface is 12 files, +424/-11, within the proposed helper/test estimate: seven Go harness/test files plus four projected JSONL fixtures and provenance (14,818 bytes). Exact paths/numstat are in `artifacts/implementation/surface.txt`; they are `internal/ensigncycle/{claude_runtime_helpers_test.go,pi_native_completion_test.go,pi_rejection_extractors_test.go,auto_continue_fixtures_test.go,claude_live_runner_test.go,shared_promoted_live_test.go,pi_auto_continue_double_dispatch_replay_test.go}` and `internal/ensigncycle/testdata/pi_native_completion/{default-parent.jsonl,default-child.jsonl,auto-parent.jsonl,auto-child.jsonl,provenance.json}`.

Both pre-edit reproductions, separate FO FIX authorizations (G1 and conflicting epoch), command details and final/initial logs are retained under `artifacts/implementation/`. The captured fixture records are field projections with source SHA-256/line provenance; the Git end state is explicitly reconstructed because no original Git repository was retained. `gofmt -w ./cmd ./internal` ran, final touched files were formatted again, and the unrelated baseline release-test spacing change was restored; code worktree is clean. State report synchronization is left to the FO under the dispatch’s explicit no-push direction.


## Stage Report: validation

- DONE: Verify captured native completion and wrong-worker/order negatives at the actual grader boundary without accepting completion prose alone.
  Candidate `ba8e6eeb4e28e310b17897a0a8b12e151366b8f9`: both captured/full retained parents pass, 27 existing negatives plus seven independent overlay probes pass; removing identity correlation makes the wrong-run test fail (artifacts/validation/).
- DONE: Confirm downstream gate and commit checks still fail correctly, and historical Pi/Claude/Codex controls remain valid.
  Existing replay preserves `entity has no gates record` and uncommitted-report rejection; bypassing either check makes its test fail. Historical status/wait, Claude replay, Codex lifecycle and cross-host gate controls pass in focused.log.
- DONE: Audit exact diff and local checks against every AC; record honest tip-CI deferral and any findings before recommending delivery.
  AC-1/AC-2 pass locally; source hashes/projections independently match all four originals. AC-3 remains partial: final normal/race logs each exit 1 only at the FO-declined installed-manifest defect; exact 12-file +424/-11 diff, formatting and candidate cleanliness checked.
- SKIPPED: Final tip PR/native host CI for AC-3.
  Explicitly deferred by dispatch to FO after independent validation; no model/native run, CI, PR, push or rebase. No skipped or unrun check is represented as green.

### Summary

Recommend PASSED for local validation of `ba8e6eeb4e28e310b17897a0a8b12e151366b8f9`, with final live/tip CI explicitly pending. Checklist: 3 DONE, 1 SKIPPED, 0 FAILED; no new material, deferred-risk or polish finding. Approval/scoping, per-AC evidence, commands, source verification, seven independent probes and three falsifying mutation logs are recorded in `artifacts/validation/review.md` and adjacent files.

Candidate and frontmatter remained unchanged. Existing final normal/race logs were reused, preserving their actual exit-1 installed-manifest failure and FO DECLINED disposition; state report/artifacts are committed locally, with synchronization left to FO under the explicit no-push dispatch.


## Stage Report: implementation (cycle 2)

- DONE: Credit verified synchronous Pi completion at its real tool-result boundary and verified async completion in the existing rejection-route observer.
  Code `8fca2b12b3ed272f59ea6ea465ed3ac4d574c7ae`: shared child verification credits sync index66 before advancement73 and exact async spawn112/done116 before gate135. Both captured tests exercise lifecycle and route boundaries; removing the sync branch or native route integration restores their captured failures.
- DONE: Reproduce both captured false failures before fixing; preserve identity/order/error and durable gate/commit negative controls afterward.
  `artifacts/implementation/p1p2/red.log` retains both pre-edit reds. Final tests preserve 27 native negatives, add 21 sync negatives and nine same-stage cases; accepting missing exit, stale result, another worker’s notice or completion after advancement/gate breaks the negative controls. Original absent-gate and uncommitted-report checks still reject through independent existing owners.
- DONE: Report exact scope and focused normal/race evidence, leaving combined full suites and all six native variants for final tip CI.
  Focused normal/race pass in 5.146s / 7.647s; live-tag compile-only check passes, with no model invocation. Exact logs, commands, provenance and surface are in `artifacts/implementation/p1p2/`; full suites/native coverage remain deferred below, not passed.

### Summary

P1/P2 were FO-authorized Material task-owned AC-1 fixes from `artifacts/validation/tip35149242496/report.md`. One shared native-child verifier now supplies per-dispatch completion identities and indexes to both observers. Sync requires a unique correlated result, explicit success/exit0 and matching run/owner/agent/task/cwd/epoch/terminal stop; fresh same-stage route identities include the exact dispatch call so an earlier correction cannot complete a later reviewer. Existing Git/report/gate checks and legacy Pi status/wait plus Claude/Codex controls remain in place.

Scope is five existing Go test/harness files (`pi_native_completion_test.go`, `pi_rejection_extractors_test.go`, `claude_live_runner_test.go`, `pi_auto_continue_double_dispatch_replay_test.go`, `pi_rejection_extractors_test_test.go`, all under `internal/ensigncycle`) and five new files in `testdata/pi_native_completion` (`sync-parent.jsonl`, `sync-child.jsonl`, `route-parent.jsonl`, `route-child.jsonl`, `tip35149242496-provenance.json`). Total +700/-83 across ten files includes 162 blank parent-record placeholders preserving exact indexes; fixtures/provenance total 13,513 bytes. The existing helper grows from 150 to 235 lines, below the proposed +100–140-line estimate. Exact numstat is retained in `surface.txt`; no skill, CLI, workflow or recorder policy changed.

Source head `d0a6f413af32fa69ac9da4af2c74c0086fcca812` and CI merge `3f23437d91c98bcf4cba93e52d0f1da2761d2cb1` share full tree `717d412d4b2b0f2ede3246c6f25f8a49d8d8d252`; raw source hashes and selected record positions accompany the fixtures. Original default durability remains transcript/shim evidence without a Git bundle; the same-stage bundle and its independent checks remain retained in the validation artifacts. This correction does not reclassify CI 35149242496 as green.

Exact `gofmt -w ./cmd ./internal` ran; the unrelated pre-existing release-test spacing was restored and final touched files formatted. `git diff --check` passes. Full `go test ./...` and race suites are deliberately deferred to the FO’s combined final tip after #801; this worker did not run them, touch #801, rebase or push. All six same-stage native variants remain pending final tip CI: `plain` previously failed; `review-required`, `separate-review-required`, `round-required`, `round-missing`, and `cycle-limit` were unstarted (not skipped or passed). Independent validation and final native CI are still owed; no frontend/outcome threshold was relaxed.


## Stage Report: implementation (cycle 3)

- DONE: Reconcile the owned runner conflict while preserving both independently validated #801 semantics and Pi native completion routing.
  Rebased head `0997b6b7590a5a89ba979194238e93f7f47d94e4` preserves native spawn counting and actual recorder exit interception alongside Pi artifact-root routing; the two-hunk range diff is retained in `artifacts/implementation/reconciliation/range-diff.txt`.
- DONE: Verify the affected observers and live-tag compilation without broad or native model runs.
  Focused normal checks pass (exit 0, 322.229s), including native identity/reuse, recorder exit controls, Pi sync/async attribution/order and downstream gate checks; live-tag compile-only passes (exit 0, 1.168s, no tests selected). Logs and exact selector are retained in `artifacts/implementation/reconciliation/`.
- DONE: Report exact rebased head, scope, focused evidence and remaining combined validation obligations.
  Old head `8fca2b12b3ed272f59ea6ea465ed3ac4d574c7ae` is reconciled onto `c572d62f05ae57d643dbaab5ecfb7de4a2ccff31`; rewritten commits are `d91dcc9dada97779546082f2ee4cb962881c94df` and `0997b6b7590a5a89ba979194238e93f7f47d94e4`. Combined validation and CI remain pending below.

### Summary

The first Pi commit applied unchanged. The second conflicted only in `internal/ensigncycle/claude_live_runner_test.go`: the same-stage check retains #801's `countRouteEvents(routes, routeSpawn) > 1` predicate with Pi `routeErr` propagation; rejection routing retains #801's non-Pi actual command-log/recorder-exit interception with Pi artifact-root completion routing. No new mechanism, relaxed control or unrelated code change was introduced; `conflict.diff`, `range-diff.txt`, `surface.txt` and `commands.txt` retain the reconciliation evidence.

The focused identity/reuse tests reject missing or wrong native completion identity; accepting a shared display name as worker identity would break those controls. Recorder execution and the launcher test verify actual successful publication and reject echo-only, failed-inner-recorder, duplicate and second-round calls; the launcher test uses a temporary stub host, with no model invocation. Pi sync/async negatives still reject wrong dispatch, run, epoch, missing child or late completion, and the original absent-gate/uncommitted-report controls still reject independently; crediting completion prose or bypassing durable checks would break these controls.

`gofmt -w ./cmd ./internal` completed with exit 0. Only the unrelated pre-existing `internal/release/runtime_live_evidence_workflow_test.go` spacing changed; its exact HEAD bytes were restored. `git diff --check` passes and the code worktree is clean. No broad normal/race suite, native/model run, CI or push was performed during this reconciliation. Independent validation owns the final combined normal/race commands; this report does not claim those suites passed. The prior full-suite exit-1 installed-manifest resolver failure and FO DECLINED disposition remain unchanged in the earlier report.

Final tip CI and all six same-stage native variants remain owed: `plain` previously failed; `review-required`, `separate-review-required`, `round-required`, `round-missing`, and `cycle-limit` were unstarted, not passed. AC-3 remains pending. State synchronization is left to the FO under the explicit no-push assignment.


## Stage Report: validation (cycle 2)

- DONE: Independently verify Pi sync and async-route completion at the captured failure boundaries, retaining strict identity/order/error and incomplete-outcome controls.
  Frozen `0997b6b7590a5a89ba979194238e93f7f47d94e4`: full retained streams prove sync65/66, validation80/81 and async112/116; ten detached subtests reject cross-worker/stale/order/assignment/artifact failures. Existing missing-report/gate/Git controls remain authoritative; see `artifacts/validation/final-tip-0997/`.
- DONE: Verify the reconciled shared runner preserves #801 controls and run final combined normal and race suites on the frozen stack tip.
  Reused #801 independent evidence; exact runner diff preserves native spawn count and actual recorder exit. Original normal exit1 includes two10m timeouts; FO-authorized isolated probes pass, then serial normal/race `-p 1 -timeout 30m` each finish exit1 only for the exact declined resolver defect; no race diagnostic.
- DONE: Report exact candidate, per-AC evidence, scope and findings, with native CI pending and no skipped or failed check described as green.
  Canonical `final-tip-0997/report.md` records AC-1/AC-2 local proof, actual red exits, timeout finding/dispositions, source hashes, unchanged candidate and the cross-stack roadmap-authorization HOLD. AC-3 final native acceptance remains pending.

### Summary

Recommend PASSED for deterministic #806 correction and #801 shared-runner reconciliation at `0997b6b7590a5a89ba979194238e93f7f47d94e4` atop `c572d62f05ae57d643dbaab5ecfb7de4a2ccff31`; checklist 3 DONE, 0 SKIPPED, 0 FAILED counts completed validation work, not green commands. Both completed broad serial commands remain exit1 for the exact FO-declined installed-manifest failure; original default-timeout failure is prominently retained, with cumulative slowdown supported and isolated persistent deadlock refuted.

Final tip native CI must execute all six same-stage variants on each host; previously unstarted variants are not passed. Claude's retained roadmap-authorization outcome failure remains HOLD/route for decision with its strict case unchanged and no waiver. No candidate/frontmatter, push/rebase/CI/PR/gate action or local native/model run occurred; current-head compile-only and owned green checks were reused, and state evidence is committed locally for FO synchronization.


## Stage Report: implementation (cycle 4)

- DONE: Credit the current gate only against its own completed validation and committed report, retaining explicit withdrawal of superseded attempts.
  Code `02ef599d016f7c3f48def7f7437f35a3f99163e2` joins actual Pi gate calls/results to the preceding native validator, retains matching successful withdrawal, and compares the current validation report section with committed HEAD through the existing durability owner.
- DONE: Demonstrate the captured repaired-gate chronology passes while premature gates, omitted withdrawal, missing/error completion, wrong identity, and missing/uncommitted reports still fail.
  Captured red was `spawns=2 completed=141 validation=99`; the corrected replay credits repair completion141 before replacement144 while checking first-attempt ordering. Two positives and 16 targeted negatives pass, alongside existing sync/async, durable gate/commit and cross-host controls; deleting withdrawal checks or accepting an uncommitted report breaks the negative cases.
- DONE: Record exact candidate changes, focused and required deterministic checks, surface estimate, and remaining held live outcomes without claiming full acceptance.
  Focused final passes in 7.918s. Required serial full normal/race both exit1 only at the known resolver defect; ensigncycle passes 293.059s/289.350s. Exact commands, red/green and full logs, dispositions and surface are retained in `artifacts/implementation/gate-repair35238049892/`.

### Summary

FO-authorized finding3 is fixed on parent `0997b6b7590a5a89ba979194238e93f7f47d94e4`, within AC-1/AC-2. The original recovered report omission remains part of the captured evidence: attempt1 was withdrawn, a repair validator completed, then attempt2 was prepared. Each observed prepare uses its own preceding completed validator; the observer no longer combines the latest completion with the first gate. Successful withdrawal also permits re-preparing the same validated report without another worker. The unchanged-stage-report/gate-metadata-only positive would fail under a mandatory-new-worker rule. FO explicitly DECLINED that invented requirement; no Briefing-read dependency, whole-entity comparison or new source parser was added.

The required commands were `go test -p 1 -timeout 30m ./...` then `go test -p 1 -timeout 30m -race ./...`, using the authorized slow-host flags. Both final commands exited1 at `TestCodexResolveManifestAgainstInstalledHost`, `codex_resolve_test.go:44`: `spacedock@spacedock not installed in codex, but resolver returned "/Users/clkao/.codex/plugins/cache/spacedock-local/spacedock/0.28.0-pre0/.codex-plugin/plugin.json"`. All other packages passed, with no race diagnostic. The exact pre-existing resolver defect remains FO DECLINED, not an unresolved correction finding; neither full suite is claimed green. The initial full normal result is also retained (same sole failure; ensigncycle224.128s) before the authorized multi-prepare correction and final reruns.

Actual scope is eight files, +438/-7, including 144 blank parent records preserving original event indices. Existing Go owners are `internal/ensigncycle/{pi_native_completion_test.go,claude_runtime_helpers_test.go,auto_continue_fixtures_test.go,pi_auto_continue_double_dispatch_replay_test.go}`; compact additions are `testdata/pi_native_completion/{repair-parent.jsonl,repair-child-1.jsonl,repair-child-2.jsonl,repair-provenance.json}` under the same package. The native helper grows by101 net lines, below the proposed90–130 helper-line estimate; four fixture/provenance files total14,698 bytes, below the proposed15–25KB. Raw parent/child source hashes accompany projections. No continuation Git bundle existed; tests reconstruct durable state through the existing fixture/Git owners rather than claiming original objects were replayed.

`gofmt -w ./cmd ./internal` completed; unrelated pre-existing release-test spacing was restored exactly. `git diff --check` passes and the committed code worktree is clean. No local native/model run, CI, YAML edit, rebase, push, merge or entity-frontmatter mutation was performed. Independent validation and final integration/native CI remain owed. CI35238049892 is not regraded green: Pi plain remains incomplete, five later same-stage variants were unstarted, and workflowScript/runs.run/revival support, fixture host wording and Claude recording/recovery findings remain HOLD outside this patch. FO owns the separate fail-fast coverage correction and later integration; AC-3/full live acceptance remains pending.
