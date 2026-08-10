---
title: Cut gate-guardrail turn and tool bloat
status: validation
source: "Journey-metrics audit of PR #643, Runtime Live E2E run 31297186020, Claude job 93204212216, artifact 9033837253. The gate-guardrail journey used 22 assistant turns and 24 tool calls, up from 11 and 11 in the v0.26 Sonnet observation. Captain directed a separate filing on 2026-08-09."
started: 2026-08-09T14:51:32Z
completed:
verdict:
score: 0.9
sprint: durable-decisions
sprint-readiness: ready
group: gate-lifecycle-ux
worktree: .worktrees/spacedock-ensign-cut-gate-guardrail-turn-and-tool-bloat
issue:
pr:
mod-block: merge:pr-merge
id: 5k704rrfk5r75vqv3bwn1yhf
gates:
    version: 1
    records:
        - id: gate:5k704rrfk5r75vqv3bwn1yhf:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:5k704rrfk5r75vqv3bwn1yhf-backlog-1
              briefing:
                id: briefing:5k704rrfk5r75vqv3bwn1yhf:backlog:attempt-1:revision-1
                digest: sha256:bdd851d207fe8f84300b59cbd34359f1cc613df5fee26b869f323db6e63f7956
                request-digest: sha256:ef783c09a07b51873e1f7c8d3fef81b226bfa69901afc1a652389886a6d7b2d2
                room-ref: ./cut-gate-guardrail-turn-and-tool-bloat/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5k704rrfk5r75vqv3bwn1yhf:backlog:1
                briefing: briefing:5k704rrfk5r75vqv3bwn1yhf:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T14:50:42.1944Z"
                decision: approve
                reason: Captain directed ideation dispatch; the measured-call classification constrains the smallest safe change.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:5k704rrfk5r75vqv3bwn1yhf:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:5k704rrfk5r75vqv3bwn1yhf-ideation-1
              briefing:
                id: briefing:5k704rrfk5r75vqv3bwn1yhf:ideation:attempt-1:revision-1
                digest: sha256:966fe77eebf72e5aa68f230fe9733fb73cb5cee62659f85207e60f0791957c2f
                request-digest: sha256:da8824069956b6ee755d3bb86d73a3250f3d727b489467fbddfd66062873c137
                room-ref: ./cut-gate-guardrail-turn-and-tool-bloat/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5k704rrfk5r75vqv3bwn1yhf:ideation:1
                briefing: briefing:5k704rrfk5r75vqv3bwn1yhf:ideation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T15:29:21.646183Z"
                decision: approve
                reason: The direction removes measured contract-caused gate toil through existing structured reads and real gate outputs, with bounded files and live comparison evidence and no new parser, harness, CLI, or lifecycle state.
              application:
                target-stage: implementation
                state: consumed
        - id: gate:5k704rrfk5r75vqv3bwn1yhf:validation
          stage: validation
          attempts:
            - id: gate-attempt:5k704rrfk5r75vqv3bwn1yhf-validation-1
              briefing:
                id: briefing:5k704rrfk5r75vqv3bwn1yhf:validation:attempt-1:revision-1
                digest: sha256:426f9cf05a4aa03f8a58e83b1a8828ab60fe565998f23e992dab11c8a1e5816a
                request-digest: sha256:50ef47bd95988565b70026652460f0f31f64d7c603f38dd2eb525f5f9dbe04f2
                room-ref: ./cut-gate-guardrail-turn-and-tool-bloat/review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:5k704rrfk5r75vqv3bwn1yhf:validation:1
                briefing: briefing:5k704rrfk5r75vqv3bwn1yhf:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-10T04:50:05.224524Z"
                decision: approve
                reason: Validation proves the Captain-accepted 17-turn/17-call path while preserving atomic gate durability and all human decision authority.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:5k704rrfk5r75vqv3bwn1yhf-validation-2
              briefing:
                id: briefing:5k704rrfk5r75vqv3bwn1yhf:validation:attempt-2:revision-1
                digest: sha256:1fb466a96f70b90e4d03b0955c20a3b79a282357877ef083e6940ea040d0c574
                request-digest: sha256:cef968ba5d141869dc45fe648bcbd1b4dafcd5b71536821c10580f9a57a454f1
                room-ref: ./cut-gate-guardrail-turn-and-tool-bloat/review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:5k704rrfk5r75vqv3bwn1yhf:validation:2
                briefing: briefing:5k704rrfk5r75vqv3bwn1yhf:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-10T05:23:04.026874Z"
                decision: approve
                reason: 'Validation proves the Captain-corrected contract-only scope: three files, 25 changed lines, rejected mechanisms absent, and all offline checks green.'
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:5k704rrfk5r75vqv3bwn1yhf-validation-3
              briefing:
                id: briefing:5k704rrfk5r75vqv3bwn1yhf:validation:attempt-3:revision-1
                digest: sha256:4656724fc8b77c469fb270d01933b23e695437f8366ba12930fd1088de900c68
                request-digest: sha256:ade3162a4d5f73d48142dc1d6147de1973cfea4de104b9df45f8d37657b9c089
                room-ref: ./cut-gate-guardrail-turn-and-tool-bloat/review/validation/briefing-3
              resolution:
                type: Resolution
                id: resolution:spacedock:5k704rrfk5r75vqv3bwn1yhf:validation:3
                briefing: briefing:5k704rrfk5r75vqv3bwn1yhf:validation:attempt-3:revision-1
                by: agent:first-officer
                at: "2026-08-10T14:41:29.995795Z"
                decision: approve
                reason: Validation proves the reconciled contract-only candidate stays within three files and 30 lines, excludes rejected mechanisms, and passes all required offline checks.
              application:
                target-stage: done
                state: pending
---

Make the supported gate-guardrail journey reach one committed open gate with materially fewer turns and tool calls.

## Problem

Runtime Live E2E run `31297186020`, Claude job `93204212216`, artifact `9033837253`, used Claude Sonnet 5 at maximum effort. Its successful `gate-guardrail` journey took 22 assistant turns and 24 host tool calls (16 Bash, five Read, three Skill), compared with 11 turns and 11 calls in the v0.26 Sonnet observation. It still produced the right result, but spent calls discovering a path already resolvable by `status --read`, inspecting unrelated state, probing documented command grammar, and projecting boot twice after preparation.

The 24 calls below are numbered in transcript order. `required` means the value transaction or authority boundary needs the call. `contract-caused` means the current FO packaging or instruction text demands it even though the journey could be safe without it. `missing CLI information` means the prompt/CLI did not name enough selected-source information for a direct call. `agent error` means an existing output or contract already made the call unnecessary. The owner is the smallest surface that can remove or retain the call.

| # | Call | Class | Owning surface and reason |
|---:|---|---|---|
| 1 | `Skill(spacedock:first-officer)` | required | Claude host entry: the prompt explicitly selects the FO. |
| 2 | Read `first-officer-shared-core.md` | contract-caused | `skills/first-officer/SKILL.md`: mandatory deferred core load. |
| 3 | Read `claude-first-officer-runtime.md` | contract-caused | `skills/first-officer/SKILL.md`: mandatory runtime-adapter load. |
| 4 | sandbox check plus `spacedock --version` | contract-caused | shared FO startup version gate. |
| 5 | `status --boot --identify --json` | required | CLI: one discovery/taxonomy/readiness projection is the scheduler input. |
| 6 | `state ready`; `state sweep` | contract-caused | shared `engage` convergence contract; already batched into one host call. |
| 7 | `Skill(spacedock:fo-gate-lifecycle)` | contract-caused | shared FO deferred-load boundary. |
| 8 | two `find` searches under the known workflow | agent error | FO execution: boot already supplied `entity_dir` and slug, while `status --read <slug> --json` resolves the canonical file. |
| 9 | Read the entity body | required | gate judgment needs the report and acceptance criteria; the replacement structured read still supplies them. |
| 10 | Read workflow `README.md` | contract-caused | cold-report review needs the stage obligation; boot already supplies taxonomy but not stage prose. |
| 11 | Read `recorder-contract.md` | agent error | FO execution: it read this claimed Reference but omitted it from `gate prepare`; selected inputs must be read once and then used. |
| 12 | Read live-harness `evidence/command.log` | agent error | FO execution: the partial logging shim was neither a selected gate source nor necessary AC evidence. |
| 13 | broad main/state Git history and stats | missing CLI information | prompt/CLI seam: the prompt names committed Artifact/References but not their paths, so one path-scoped committed-tree listing is needed; the broad history bundle is not. |
| 14 | main and state `git status` | agent error | CLI owns selected-source cleanliness validation and `state commit` owns path scope. |
| 15 | repeated Git log/show after a retained-cwd mistake | agent error | FO execution: one path-scoped committed-tree listing should replace both history probes. |
| 16 | Read committed gate review and entity snapshot | required | gate judgment must inspect the exact selected Artifact and References. |
| 17 | `gate prepare --help` | contract-caused | `skills/fo-gate-lifecycle/SKILL.md` says to preflight; the version gate and documented real invocation already establish the surface. |
| 18 | `gate prepare ...` | required | CLI creates exactly one open binding and emits room/Briefing/digest/state. |
| 19 | `state commit recorded-gate-task` | required | CLI makes the open binding durable without sweeping the dirty sibling. |
| 20 | first post-prepare boot projection | contract-caused | gate skill says to re-read `awaiting-captain`; `state=open` plus successful commit already proves that deterministic readiness. |
| 21 | second full entity Read | contract-caused | gate skill asks for a re-read; structured checklist/AC reads can supply only the presentation inputs after commit. |
| 22 | `Skill(spacedock:present-gate)` | required | lazy presentation contract is the authority-preserving captain boundary. |
| 23 | `grep` for report line numbers | agent error | shared core already specifies `status --read --checklist`, whose output includes `start`/`end`. |
| 24 | second post-presentation boot projection | agent error | FO execution: presentation is the stopping condition; it cannot change durable readiness. |

Totals: seven required, nine contract-caused, one missing-CLI-information, and seven agent-error calls. The product does not need a new projection. Existing structured reads contain the missing entity/checklist/AC facts; only selected-source path discovery remains judgmental, and one path-scoped Git tree listing is sufficient for the prompt's unnamed committed inputs.

## Proposed approach

Change the gate-lifecycle contract, and only that behavior owner, to prescribe a bounded cold-gate sequence:

1. Retain `definition_dir`, `entity_dir`, slug, stage, and readiness from the single boot. Resolve the entity with `status --read <slug> --json`; never `find` or broad-search a known workflow. If a prompt names committed gate materials without exact paths, list committed Markdown candidates once with path-scoped `git -C ... ls-tree`, then read and use the selected bytes once. Do not inspect harness logs, broad Git history, or worktree status for a cleanliness decision that `gate prepare` owns. This serves AC-1 and AC-2. The simpler alternative, only saying “be concise,” is insufficient because it does not give the model a direct discovery route.
2. Delete the lifecycle-surface help preflight. The documented real `gate prepare` invocation is the capability check; any nonzero result halts with its exact error. This serves AC-1 and AC-2. Retaining the help call is insufficient because it spends a call without strengthening the subsequent real command.
3. After `gate prepare` emits the required `state=open` and `state commit` exits zero, do not project boot again. In one read-only shell event, run `status --read <slug> --checklist --json` and `status --read <slug> --ac-scan --json`; use their stage, checklist text, line ranges, and citations to cross-check and present. Do not full-read/grep the entity again, and do not project status after presentation. This serves all four ACs. The simpler alternative, merely deleting the second boot, leaves the first boot, reread, and grep waste in place and is unlikely to reach the turn ceiling.

This is instruction routing, not a harness controller: the binary continues to own resolution, preparation, validation, and mutation. No scheduler, new command, state, format, simulator, retry path, or authority is introduced.

## Spike result

The riskiest claim is that the post-prepare boot projection can be removed without losing readiness or presentation evidence. Exercised on 2026-08-09 against the current real CLI:

- `go test ./internal/ensigncycle -run '^TestRecordedGateLifecycle(RealCLIReplay|PhaseDetectionIgnoresHelpProbes)$' -count=1 -v` passed. The replay observed `state=open`, committed the exact open attempt, and its negative-control portion removed all help events without changing lifecycle qualification.
- `go test ./internal/status -run '^(TestDefaultStageMatchesExplicitCurrent|TestBootReadyGatesUsesStatusDerivedSelection)$' -count=1 -v` passed. Omitted-stage checklist and AC scans matched explicit current-stage output in text and JSON, and an open selected current attempt deterministically projected `awaiting-captain`.
- A direct `status --read <slug> --json` against this workflow resolved the canonical absolute entity path. Direct checklist output exposed `start` and `end` line numbers, and AC scan exposed citations/unevidenced state.

Therefore the selected contract-only mechanism rests on exercised CLI behavior. A new lifecycle projection or static prose proof is rejected.

## Expected surface and semantic boundaries

- `skills/fo-gate-lifecycle/SKILL.md`: replace the preflight/cold-report/post-commit wording, at most 12 inserted and 16 deleted physical lines, with net bytes non-positive. Current size is 6,813 bytes; it must remain below the existing 7,000-byte component cap.
- `docs/site/concepts/gates-and-decisions.md`: replace one sentence (two inserted and two deleted wrapped lines) so the public description says successful prepare+commit supplies the open-state proof and the FO performs one structured checklist/AC read before presentation.
- Total tolerance: exactly two files, at most 14 insertions and 18 deletions. No Go, fixture, harness, generated artifact, or new test file.

Observable semantic changes are limited to FO runtime behavior: one path-resolved cold-report inspection, no gate help probe, no post-prepare/post-presentation boot projection, and one grouped structured evidence read. Command grammar, command output, stored formats, gate authority, mutation order, and lifecycle state do not change. The result remains exactly one prepared and committed open attempt presented to the captain; no decision, consume, successor dispatch, or other mutation occurs.

Concrete documentation diff:

```diff
--- a/docs/site/concepts/gates-and-decisions.md
+++ b/docs/site/concepts/gates-and-decisions.md
@@
-question, Artifact, summary, and References, commits the emitted binding,
-re-reads `awaiting-captain`, and presents it.
+question, Artifact, summary, and References, commits the emitted `state=open`
+binding, performs one structured checklist/AC read, and presents it without another boot projection.
```

### Superseded live-mechanism acceptance criteria (historical)

1. The real gate journey requires materially less operator effort.

One candidate run of `TestLiveCommonGateGuardrail` with Claude Sonnet 5 at maximum effort uses at most 16 assistant turns and 18 host tool calls, measured from its archived Claude stream with the same journey-metrics collector used for baseline run `31297186020` (22 turns, 24 calls).

2. Known search and retry waste is absent.

The candidate stream contains no filesystem search for the known workflow/entity/contract path, no `gate ... --help` or failed command-shape probe, no broad/repeated Git or worktree-status inspection, and no boot projection after successful `gate prepare`.

3. Gate authority and final state remain correct.

The candidate's durable split-root state contains exactly one committed current-stage open attempt and its two-file room, preserves the dirty sibling and all preexisting body bytes, and contains no Resolution/decision, withdrawal, application consumption, status advance, successor dispatch, archive, or other unauthorized mutation. The captain-facing review names the bound Briefing/digest and ends at the human decision boundary.

4. The fix changes only the owning contract and its public description.

The candidate diff stays within the two-file/LOC/byte budget above; no CLI grammar/output, Go code, harness controller, fixture, lifecycle state, standing check, simulator, or static prose-presence test is added.

## Acceptance criteria

**AC-1 — The merge-base candidate contains at most the three allowed contract/documentation files and 30 changed lines, all directly reducing the documented First Officer gate path.**
Verified by: validation cycle 3 recorded merge base `a929fcb60` through candidate `339d05a23` as exactly the three allowed paths at +11/-14 (25 changed lines), with a fourth path or sixth additional line as the falsifying boundary.

**AC-2 — The candidate contains no `status --gate-evidence`, `gate prepare-review`, transactional preparation/rollback, related documentation/tests, or prepare-review-only harness/removed-verb compatibility surface.**
Verified by: validation cycle 3 inspected the cumulative diff and recorded no Go, test, fixture, harness file, rejected command, or transaction surface.

**AC-3 — Existing focused gate-lifecycle, status, and contract smoke checks pass, as do formatting, the full suite, the race suite, and the cumulative diff check.**
Verified by: validation cycle 3 recorded green focused lifecycle/status/contractlint checks, `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, and `go test ./... -race`, including the lifecycle and status tests' falsifying conditions.

**AC-4 — Validation requires and triggers no authenticated model journey or GitHub Actions workflow; prior runtime counts remain observations rather than acceptance guarantees.**
Verified by: validation cycle 3 explicitly recorded that no authenticated journey or GitHub workflow was invoked while all offline contract-only evidence passed.

## Test plan

1. Before editing, retain the spike commands/results above as the mechanical baseline. During implementation, apply the exact contract and documentation wording within the declared budget; use `git diff --stat`, `git diff --numstat`, `wc -c`, and `git diff --check` to verify surface and bytes. Estimated cost: small, contract-only.
2. Run existing focused behavioral checks, not a prose-presence test: `go test ./internal/ensigncycle -run '^TestRecordedGateLifecycle(RealCLIReplay|PhaseDetectionIgnoresHelpProbes|MissingEventControls)$' -count=1` and `go test ./internal/status -run '^(TestDefaultStageMatchesExplicitCurrent|TestBootReadyGatesUsesStatusDerivedSelection)$' -count=1`. The lifecycle tests fail on missing/duplicate prepare/commit or unauthorized ordering; the status tests fail if structured reads lose current-stage or open-readiness semantics.
3. Run the repository-required suite: `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`. No files should change from formatting because no Go source is in scope.
4. Run exactly one real comparison at the candidate commit with the same model and journey: `SPACEDOCK_LIVE_RUNTIME=claude SPACEDOCK_LIVE_MODEL=claude-sonnet-5 go test -tags live -count=1 -timeout 20m -run '^TestLiveCommonGateGuardrail$' ./internal/ensigncycle -v`, with the CI maximum-effort setting and artifact capture used by Runtime Live E2E. Count assistant turns and host tool calls from the archived stream; compare to run `31297186020`; inspect the command log and durable state for AC-2/AC-3. Do not rerun merely to obtain a better metric: a miss returns to implementation analysis.
5. No permanent metrics simulator, harness modification, static wording assertion, new runtime lane, or second model run is permitted. The existing live journey is the skill smoke test required for this command-text change.

## Stage Report: ideation

- DONE: Classify every call in the retained 24-call journey as required, contract-caused, missing CLI information, or agent error, with its owning surface.
  Artifact `9033837253` supplied the exact 24-call stream; the table accounts for all calls and totals 7/9/1/7.
- DONE: Exercise the riskiest proposed mechanism before selecting it, and reject any harness controller, static prose proof, or new lifecycle state.
  Real lifecycle and status tests proved help is optional, open readiness is deterministic, and structured reads carry stage evidence; no new product mechanism is proposed.
- DONE: Specify the smallest change, expected files and LOC, acceptance evidence, and a same-model live comparison capped at 16 turns and 18 tool calls.
  The gated baseline is two contract/doc files, at most 14 insertions/18 deletions, plus one Sonnet-5 live comparison against run 31297186020.

### Summary

Selected a bounded gate-lifecycle contract correction: use existing structured reads, trust the real prepare/commit outputs, and stop after one presentation. The design removes the mandated help/reprojection calls and constrains agent discovery waste without adding CLI or harness machinery.

## Stage Report: implementation

- DONE: Deliver only the two-file gate-lifecycle contract and public-doc correction within the approved line and byte budgets, with no CLI, fixture, harness, schema, or new-test changes.
  Commit `7da5d36a5` changes exactly the two approved files by +6/-9 lines; the skill shrank from 6,813 to 6,762 bytes, and `git diff --check` passed.
- DONE: Preserve one prepared and committed open gate presentation while replacing search/help/reprojection waste with the existing path-resolved structured reads.
  The focused lifecycle tests fail on missing/duplicate prepare/commit or unauthorized ordering; status tests fail if structured reads lose current-stage or open-readiness semantics, and both groups passed.
- FAILED: Run the focused lifecycle/status checks, repository-required formatting and test suites, and exactly one Sonnet-5 live comparison capped at 16 turns and 18 tool calls; do not rerun a miss for luck.
  Focused checks, `gofmt -w ./cmd ./internal`, and `go test ./...` passed; `go test ./... -race` hit unrelated load-sensitive `TestSonnetTeamDeleteHangReplay` (the exact race test then passed), and the sole live invocation skipped before model launch because no nonempty benchmark token or CI API key was available.

### Summary

Committed the bounded gate-lifecycle and public-documentation correction without changing runtime code, fixtures, or tests. Mechanical behavior checks passed, but this environment could not supply the required live comparison and the full race suite exposed one unrelated load-sensitive replay failure; both are retained for validation rather than hidden by reruns or out-of-scope edits.

## Stage Report: implementation (cycle 2)

- DONE: Retain the exact approved two-file candidate and focused/full behavior evidence; make no product change unless this run proves the candidate itself is defective.
  Commit `7da5d36a5` remains byte-clean and unchanged; focused checks and `go test ./...` from cycle 1 remain the behavior baseline.
- DONE: Run the repository-required full race suite in a load-safe serial context and require it to pass; read any new failure from this run rather than inheriting the prior label.
  `go test ./... -race -p 1` passed, including `internal/ensigncycle` in 384.486s; the prior parallel timing label was not inherited.
- FAILED: Obtain exactly one actual Sonnet-5 gate-guardrail comparison at the candidate commit through an available authenticated local or GitHub path (the prior pre-launch skip is not a run), verify the 16-turn/18-tool-call and authority-state limits, and append a complete replacement implementation Stage Report with no FAILED item only if all obligations are met.
  The sole actual max-effort Sonnet-5 run passed authority-state assertions but measured 18 turns/27 calls against 16/18; stream `d9b7b6af…`, final message `e06a60be…`, and metrics `0e938606…` remain under `/private/tmp/spacedock-gate-live.3fPRvk/`.
- DONE: Route the material value failure through the implementation finding boundary without unauthorized candidate mutation or rerun.
  FO authorized `route for decision`: calls 8–17 precede lifecycle-skill load at call 18, so the dominant waste lies outside the approved surface and requires Captain design authority.

### Summary

The serial race suite passed and the one authenticated comparison preserved correct open-gate authority: one prepare and commit, no Resolution, consume, advance, or successor dispatch. The value ceiling still failed at 18 turns and 27 calls, so the approved candidate is held unchanged and the pre-load ownership boundary is escalated for a Captain decision.

## Design Reset: Automatic Headless Engage

Captain authorization on 2026-08-09 expands ownership from the lifecycle/doc pair into the boot-resident shared core. This section supersedes the earlier two-file surface, ACs, and test plan; commit `7da5d36a5` remains the unchanged starting candidate until the First Officer separately authorizes implementation.

### Revised ownership and exact boundary

The shared core owns automatic headless engage only through gate selection. The required order is:

1. Run the existing single local boot and enter `«interaction.boundary»()` with its retained workflow, mods, dispatchable, and ready-gate record.
2. Interactive launch presents that record, names each ready gate, hints `engage`, and stops. It does not converge, load `fo-gate-lifecycle`, open a gate/entity/source, or render a review.
3. Headless launch automatically invokes `«engage»()` once per selected workflow. Engage runs `state ready`; on exit 3 it halts. It then runs the separate read-only `state sweep`, followed by every registered startup hook exactly once.
4. After convergence and hooks, the shared core obtains the authoritative `status --next --json` envelope and selects its first `ready_gates` row before considering dispatchable work. This gate-first selection is boot-resident ownership; `fo-dispatch-core.md` remains deferred unless no ready gate exists and a worker dispatch is actually considered.
5. The selected row triggers `Skill(skill="spacedock:fo-gate-lifecycle")` as the first gate action. No gate/entity/source read, filesystem listing/search, Git inspection, capability probe, presenter load, or mutation may precede that load.
6. The lifecycle skill then owns one grouped path-resolved evidence read, the existing write-core boundary, one prepare and commit, one grouped checklist/AC scan, and one presentation. It leaves the attempt open without Resolution, consume, advance, or dispatch when no conn exists.

This adds no scheduler, state, lease, retry loop, command, or lifecycle format. It moves only the automatic headless gate-selection boundary into the already boot-resident core so the lifecycle contract can govern every gate-related action.

### Feasibility budget

The no-startup-hook `gate-guardrail` fixture has this falsifiable lower-bound plan:

| Phase | Host calls | Required events |
|---|---:|---|
| Load and boot | 6 | FO Skill; shared-core read; runtime-marker probe; adapter read; grouped sandbox/version gate; one boot |
| Automatic engage | 2 | `state ready`; then `state sweep` and post-sweep `status --next --json` in one guarded shell event because boot declares zero startup hooks |
| Gate action | 7 | lifecycle Skill; grouped exact-source read; write-core read; prepare; commit; grouped checklist/AC scan; presenter Skill |
| Total | 15 | Three-call margin under 18; 15 dependent tool-bearing turns plus the final presentation equals the 16-turn ceiling |

For a workflow with registered startup hooks, sweep completes before the hook bodies and selection; their declared calls count normally and are not hidden by batching. The 16-turn ceiling is feasible for the specified zero-hook comparison but has no turn margin. If the next sole live run needs an additional dependent turn, the Captain must choose between relaxing the turn threshold or authorizing a deeper mechanism such as a binary-owned convergence/selection fold; this task does not assume either choice.

### Revised expected surface and semantic boundaries

- `skills/first-officer/references/first-officer-shared-core.md`: clarify automatic headless engage, gate-first selection ownership, and the pre-load prohibition in at most 6 inserted and 8 deleted physical lines; net bytes must be negative so the file remains strictly below its 26,900-byte cap.
- `skills/fo-gate-lifecycle/SKILL.md` and `docs/site/concepts/gates-and-decisions.md`: retain commit `7da5d36a5` byte-for-byte.
- Total product diff from the original base: exactly three files, at most 12 insertions and 17 deletions. No CLI, fixture, harness, schema, generated artifact, simulator, new test, command grammar/output, or lifecycle-state change.

Interactive boot-only semantics, gate authority, mutation order, and the decision recorder remain unchanged. The only new semantic owner is the shared core's automatic headless transition from convergence to ready-gate selection; ordinary dispatch still loads the dispatch core before its first dispatch capability.

### Superseded automatic-headless acceptance criteria (historical)

1. Full-path headless gate work meets the operator-effort ceiling.

One new max-effort Sonnet-5 `TestLiveCommonGateGuardrail` run at the authorized implementation commit uses at most 16 assistant turns and 18 host tool calls, including state ready, state sweep, startup-hook handling, gate binding, commit, and presentation.

2. The automatic boundary is ordered and waste-free.

The archived stream orders boot → interaction boundary → state ready → state sweep → zero registered startup hooks → ready-gate selection → lifecycle Skill → gate evidence. It contains no gate read/search/list/Git inspection before lifecycle load, no help/shape probe, no broad history/status inspection, and no boot projection after prepare.

3. Interactive boot remains boot-only.

The interactive branch remains byte-identical: it names ready gates and stops without convergence, lifecycle/presenter load, gate evidence read, preparation, or mutation.

4. Gate authority remains exact.

The comparison passes `assertGateHeld`: exactly one committed current-stage open attempt and bound two-file room, with the dirty sibling and prior bytes preserved and no Resolution, decision, consume, status advance, successor dispatch, archive, or other unauthorized mutation.

5. The design reset changes only the expanded owner.

The diff stays within the three-file/line/byte surface above and adds no runtime mechanism, CLI or test surface.

### Revised test plan

1. Before editing, rerun the focused contract and lifecycle controls: `go test ./internal/contractlint -run '^(TestSharedCoreRemainsBelowPreChangeByteCap|TestBootResidentDeferredLoadPointsResolve|TestFOInstructionComponentCaps)$' -count=1`; the cap test fails on shared-core growth, while closure/cap controls fail on a dangling deferred boundary or oversized component.
2. Run `go test ./internal/ensigncycle -run '^(TestShallowBootFixtureContainsOnlyHeldGate|TestShallowBootPromptIsMutationFreeInteractiveGreet|TestGateGuardrailNegativeBrokenStateTransition|TestRecordedGateLifecycle(RealCLIReplay|PhaseDetectionIgnoresHelpProbes|MissingEventControls))$' -count=1` and the existing status ready-gate/next-envelope checks. These fail if boot-only fixtures mutate/present or if open-gate ordering/authority breaks; no prose-presence test is added.
3. After separate FO implementation authorization, change only the shared-core file, retain the two candidate files byte-for-byte, and verify `git diff --numstat`, `wc -c`, `git diff --check`, and the three-file cap.
4. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and load-safe `go test ./... -race -p 1`.
5. Spend exactly one new authenticated max-effort Sonnet-5 `TestLiveCommonGateGuardrail` run. Grade its archived metrics and ordered tool stream against AC-1/AC-2, and require the real durable `assertGateHeld` state proof for AC-4; do not rerun a miss for luck.
6. Preserve the interactive branch byte-for-byte and use the existing shallow-boot fixture/negative controls as its no-mutation regression evidence. Do not add a second model run, permanent metrics simulator, harness change, or static wording assertion.

## Stage Report: implementation (cycle 3)

- DONE: Define the exact automatic headless-engage boundary so fo-gate-lifecycle is the first gate action after convergence and gate selection, while interactive boot-only still names the gate and stops without loading it.
  The reset assigns gate-first selection to the boot-resident core and forbids every gate action before lifecycle load; the interactive branch remains byte-identical and stop-only.
- DONE: Account for state ready, state sweep, and startup hooks in the real pre-load path; reassess whether the 16-turn/18-call ceiling is feasible and propose a falsifiable measurement plan.
  The zero-hook fixture budgets 15 calls and exactly 16 turns including presentation; the sole next live stream must prove the full ordering and thresholds, with no rerun.
- DONE: Update the task design, acceptance criteria, expected surface, and test plan under the Captain-approved expanded ownership before changing candidate bytes or spending a new live run.
  Revised AC-1–AC-5 and the six-step plan are durable here; candidate `7da5d36a5` is unchanged and no new model run was spent.

### Summary

Reset the design around automatic headless engage: converge fully, select gates before dispatch, and load lifecycle before any gate evidence. The value ceiling is feasible for the zero-hook comparison but has no turn margin; implementation remains blocked on distinct First Officer authorization, and a miss would return to the Captain for a threshold-or-mechanism decision.

## Stage Report: implementation (cycle 4)

- DONE: Implement only the revised shared-core automatic headless-engage and pre-load boundary; retain the lifecycle-skill and documentation bytes from commit 7da5d36a5 unchanged.
  Commit `159ba44a6` changes only the newly authorized shared-core owner beyond retained commit `7da5d36a5`. The cumulative product diff is exactly three files and +11/-14 lines; shared core is 26,739 bytes. The retained lifecycle and documentation SHA-256 values remain `46e16eb5fd431d6d950541a6d2e6e8b725d1b4a0f2bae04ab1a8b8b61a8a748b` and `1be775009fb125eb69791b3754a9feba8a91e535c9b62ca6f74388c48f6a402d`.
- DONE: Run the focused controls, diff and byte limits, formatting, the full suite, and the load-safe race suite; stop on evidence from this candidate.
  Contract, shallow-boot, lifecycle, and status controls passed; `git diff --check`, `gofmt -w ./cmd ./internal`, and `go test ./...` passed. The first serial race attempt exhausted temporary disk after the ensigncycle package passed; after disk capacity recovered, the same `go test ./... -race -p 1` suite passed without a product edit. The implementation branch was clean at commit `159ba44a6` before and after the sole live run.
- FAILED: After all deterministic checks pass, spend exactly one authenticated max-effort Sonnet-5 gate-guardrail run and grade its full ordered path against the 16-turn/18-call and authority limits without rerunning a miss.
  The sole authenticated run passed `TestLiveCommonGateGuardrail` and its real durable authority assertions in 341.58 seconds, but measured 28 assistant turns and 28 host tool calls against 16/18. It used 15 Bash, nine Read, three Skill, and one ListAgents call. Stream SHA-256 is `ce7bb65debfe85c712e10ea8a1e8a7c3a5c2d09aabdbc188d70e485c730efab2`, final-message SHA-256 is `0bf05c079999b22daeac420fc3e5661817d6a2e2d831a03972689a2e44435891`, and metrics SHA-256 is `ab9f540efdb3a2b0737c9ef24c0c85eb38ee12571d930b892c85c82c1fb422c4`; exact files remain under `/private/tmp/spacedock-gate-live.3fPRvk/cycle4-auth-artifacts/` and `/private/tmp/spacedock-gate-live.3fPRvk/cycle4-auth-metrics/`.
- DONE: Preserve the failed candidate and route the material value evidence across the design boundary without an unauthorized edit or rerun.
  The revised boundary itself ordered boot, `state ready`, `state sweep`, `status --next --json`, then `Skill(spacedock:fo-gate-lifecycle)` before any gate evidence. After lifecycle load, however, the model spent 13 discovery/evidence calls before prepare, including `git status`, broad Git history/tree probes, two filesystem `find` operations, harness command-log inspection, and a dirty-sibling read. AC-2 therefore failed with AC-1 even though AC-4 authority passed. This is material to value and requires Captain design authority for any stronger lifecycle enforcement, binary-owned evidence bundle, changed threshold, or expanded surface; candidate `159ba44a6` is unchanged and was not rerun.

### Summary

Implemented and mechanically validated the Captain-approved automatic headless gate-first boundary. The one authenticated comparison proved that boundary and preserved exact gate authority, but missed the value ceiling at 28 turns/28 calls because discovery waste remained after lifecycle load. The failed evidence is preserved and routed for a Captain design decision with no candidate mutation or rerun.

## Review-finding disposition: cycle-4 live value failure

- Exact finding: the sole authenticated max-effort Sonnet-5 `TestLiveCommonGateGuardrail` run at candidate `159ba44a6605213b4efb4f3600b5636cd6a44a31` passed the real durable `assertGateHeld` authority check but measured 28 assistant turns and 28 host tool calls against the approved 16-turn/18-call ceiling. Calls 7–10 correctly ordered `state ready` → `state sweep` → `status --next --json` → `Skill(spacedock:fo-gate-lifecycle)` before gate evidence. Calls 11–23 then spent 13 discovery/evidence calls before prepare, including `git status`, broad Git history/tree inspection, two filesystem `find` operations, harness command-log inspection, and a dirty-sibling read.
- Released user and normal workflow: a headless First Officer engaging a known ready gate follows the supported automatic-engage and gate-lifecycle contracts to inspect committed gate materials, prepare one room, commit it, and stop at the Captain boundary.
- Observable harm: the supported journey consumes 12 more turns and 10 more host calls than the approved ceiling and is also slower in operator-visible interaction than the retained 22-turn/24-call baseline; correct authority is preserved, but the task's promised effort reduction is absent.
- Affected value AC or non-negotiable boundary: `value-ac[AC-1]` requires the full max-effort Sonnet-5 headless gate path, including convergence, preparation, commit, and presentation, to use at most 16 assistant turns and 18 host tool calls.
- Trigger evidence: the measured JSON reports `turns: 28`, `tool_calls: 28`, outcome `passed`, and duration 340,563 ms. Exact stream SHA-256 is `ce7bb65debfe85c712e10ea8a1e8a7c3a5c2d09aabdbc188d70e485c730efab2`, final-message SHA-256 is `0bf05c079999b22daeac420fc3e5661817d6a2e2d831a03972689a2e44435891`, and metrics SHA-256 is `ab9f540efdb3a2b0737c9ef24c0c85eb38ee12571d930b892c85c82c1fb422c4`; exact files remain under `/private/tmp/spacedock-gate-live.3fPRvk/cycle4-auth-artifacts/` and `/private/tmp/spacedock-gate-live.3fPRvk/cycle4-auth-metrics/`.
- Worker proposal: `Material` / current task owns the value failure, but the sufficient remedy is outside the approved shared-core-only reset / `ROUTE FOR DECISION`. Keep candidate `159ba44a6` unchanged; Captain design authority is required before adding a binary surface or changing accepted value.
- First Officer authorization: `ROUTE FOR DECISION`; preserve the candidate and exact evidence, compare the authorized design options, and do not edit, commit, test, or rerun the candidate.

### Design-boundary comparison

| Option | Effect on 16 turns / 18 calls | Ownership and estimated surface | Disposition |
|---|---|---|---|
| Stronger lifecycle prose | The smallest file change, but not a sufficient mechanism. The retained lifecycle already says to use one grouped path-resolved evidence read and forbids the observed searches/status/history; the sole run ignored those imperatives after loading the skill. More wording cannot falsifiably guarantee removal of the 12 excess calls. | `skills/fo-gate-lifecycle/SKILL.md`; about 1 file, +4/-4 lines. | Reject as insufficient for `value-ac[AC-1]`. |
| One binary-owned gate-evidence bundle | Replaces calls 11–23 with one supported read-only call. The measured-path projection is 28 - 13 + 1 = 16 host calls, leaving two calls below the ceiling; because those 13 calls occupied dependent assistant turns, the same fold projects approximately 16 turns. It preserves model judgment over which committed candidates become Artifact/References while moving path resolution, Git-object identity, and bounded content retrieval into the binary. | Add `status --read <slug> --gate-evidence --json`. `internal/status/gate_evidence.go` owns canonical entity/stage resolution and emits the current entity report, current workflow-stage prose, and committed Markdown candidate blobs from the workflow root and resolved entity subtree with repo-relative path, root identity, object ID/digest, and bytes; it excludes untracked/dirty siblings and performs no mutation. `internal/status/native_runner.go` owns flag routing and incompatibility checks. `internal/status/gate_evidence_test.go` owns behavior/adversarial fixtures. `skills/fo-gate-lifecycle/SKILL.md` routes the first post-load action to the bundle and forbids supplementary discovery. `docs/site/reference/command-reference.md` and `docs/site/concepts/gates-and-decisions.md` document the read-only envelope and preparation path. Estimate: 6 files, about +260/-12 lines. | Select: smallest sufficient mechanism. |
| Threshold change | Raising the ceiling to at least 28/28 makes the current run pass only by redefining success. It accepts worse effort than the 22/24 baseline and does not deliver the task's value. | Entity AC/design text only; no product file required, but only Captain authority can change accepted value. | Reject unless the Captain explicitly abandons or replaces `value-ac[AC-1]`. |

### Proposed mechanism and falsifiable proof plan

Propose exactly one mechanism: the read-only `status --read <slug> --gate-evidence --json` bundle above. It is smaller than a scheduler or binary-owned gate lifecycle: the First Officer still judges evidence, writes the summary, calls the existing `gate prepare`, commits, loads the presenter, and stops under the existing authority contract. The binary owns only facts it can resolve deterministically from the two Git roots.

1. Before implementation, add focused status tests that construct the recorded-gate fixture and require one envelope to return exactly the canonical entity/current-stage excerpt, workflow validation prose, committed `gate-review.md`, `entity-snapshot.md`, and `recorder-contract.md`, with stable path/root/object/digest/byte identity. The test must fail if the bundle includes untracked `dirty-sibling.md`, traverses outside either root, accepts a non-gate stage, silently returns a dirty selected blob, or mutates either repository.
2. Add a command-level fixture that feeds the returned exact paths unchanged to existing `gate prepare`, then commits and runs `assertGateHeld`. Mutants removing a blob, changing its object ID/bytes, swapping roots, duplicating candidates, or dirtying a selected source must fail before prepare or at the existing prepare guard. No prose-presence test counts.
3. Update only lifecycle routing and the two public references after the binary behavior is green. Run focused status/lifecycle/contract checks, formatting, full, race, and the detached adversarial audit required for shipped contract/status changes.
4. Spend one new max-effort Sonnet-5 gate-guardrail run only after deterministic proof passes. Require ordered boot → ready → sweep → next → lifecycle → exactly one gate-evidence bundle → prepare → commit → grouped checklist/AC scan → presenter, at most 16 turns/18 calls, no supplementary search/list/Git/status inspection, and a passing real `assertGateHeld`. A metric miss, any extra discovery call, bundle identity mismatch, dirty-sibling inclusion, or authority mutation falsifies the mechanism and returns to this checkpoint without a rerun for luck.

Candidate `159ba44a6605213b4efb4f3600b5636cd6a44a31` and its code worktree HEAD remain unchanged; no test or live workflow was run for this disposition analysis.

## Stage Report: implementation (cycle 5)

- DONE: Add focused status tests before product edits for canonical gate evidence identity, dirty selected blobs, dirty siblings, traversal, non-gate stages, duplicate or swapped roots, and zero mutation.
  `TestGateEvidenceReturnsCanonicalCommittedBundleWithoutMutation` fails if exact root/path/input-path/object/digest/bytes or candidate cardinality diverges, the dirty sibling appears, or either Git status changes; `TestGateEvidenceFailsClosed` fails if any named invalid case exits zero. Both were red before implementation and pass at `9e21db9c3`.
- DONE: Implement the read-only `status --read <task> --gate-evidence --json` bundle across the approved six-file surface; keep gate judgment, preparation, commit, presentation, and recorder authority unchanged.
  Commit `9e21db9c3` changes six mechanism files by +325/-4; cumulative candidate `7da5d36a5` + `159ba44a6` + `9e21db9c3` is seven files and +333/-15, below the 7-file/360-line stop boundary. The envelope pins immutable local Git bytes while lifecycle prose leaves selection, summary, prepare, commit, presentation, and recording with the First Officer.
- FAILED: Run focused command and lifecycle checks, formatting, full, race, detached adversarial audit, then exactly one max-effort Sonnet-5 gate-guardrail run against 16 turns/18 calls and `assertGateHeld`; do not rerun a miss.
  Focused status/lifecycle/contract checks, `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, and serial `go test ./... -race -p 1` passed. A detached exact-head audit made entity inclusion and dirty-byte acceptance mutants; the canonical and dirty-source tests failed on those mutations. The sole authenticated live run passed `assertGateHeld` but measured 21 turns/22 calls, so it was not rerun.

### Summary

The binary bundle shipped at `9e21db9c3` and reduced the prior 28-turn/28-call run to 21/22 while preserving exact gate authority, candidate identity, and dirty-sibling exclusion. It still misses the approved 16/18 value: after the correct bundle call, the model made forbidden plugin grep, help, skill grep, and `pwd`/`ls` probes before prepare, then a post-presenter README read and ListAgents call. Exact stream/final/metrics SHA-256 values are `bedb642e6c43ff30c37eb6d57f30f1d6b29feb49a111fd24c1175ca6903ab1ee`, `d0a62c126a892c9b4bb59e26aa178d105ba5a887163eb91c68ba77c88cb348fe`, and `530f2342fbbe7a36e833e211dd240c996c484e48428e53148972413945b40497` under `/private/tmp/spacedock-gate-live.3fPRvk/cycle5-bundle-{artifacts,metrics}/`; the candidate is frozen with no rerun.

## Design Reset: Atomic Gate-Review Preparation Composite

Captain direction retains the original 16-turn/18-call value and authorizes one binary composite over deterministic post-selection work. Candidate `9e21db9c3` remains the unchanged design baseline.

### Command and ownership

Add one split-root command surface with two modes under the same owner:

```text
spacedock gate prepare-review ENTITY --workflow-dir WORKFLOW_DIR --json
spacedock gate prepare-review ENTITY --workflow-dir WORKFLOW_DIR --publish \
  --question TEXT --artifact INPUT_PATH --summary TEXT --recommendation TEXT \
  [--reference INPUT_PATH ...] --json
```

Inspect mode is read-only. It emits the current `gate-evidence` entity/stage/candidate envelope plus an exact `publish_argv` skeleton and launch cwd. Publish mode re-resolves the envelope, accepts only its unique `input_path` candidates, and treats question, Artifact/References, summary, and recommendation as First Officer inputs. It invokes the existing gate preparation, publishes the entity plus two-file room in one path-scoped state commit/sync, and emits the committed checklist/acceptance-evidence and complete presenter inputs. It does not choose or synthesize evidence, summary, recommendation, or verdict.

The First Officer still loads `spacedock:present-gate`, presents once, and later uses the existing recorder. The composite never presents, records, consumes, advances, dispatches, recommends, decides, or adds a gate state, lease, scheduler, retry loop, room format, or alternate recorder. Existing `status --read --gate-evidence`, `gate prepare`, `state commit`, checklist/AC reads, and recorder commands remain compatible; lifecycle routing adopts the composite.

### Quantitative call projection

| Measured cycle-5 calls | Composite effect |
|---|---|
| 1–12: FO load through lifecycle load and gate evidence | Calls 1–11 stay; call 12 becomes `gate prepare-review` inspect. |
| 13 plugin grep, 14 `gate prepare --help`, 15 skill grep, 16 `pwd; ls` | Removed: inspect output supplies authoritative argv grammar, cwd, and candidate input paths. |
| 17 prepare, 18 binding commit, 19 checklist/AC projection | Replaced by one `gate prepare-review --publish` call. |
| 20 presenter Skill | Stays. |
| 21 post-presenter README read | Removed: publish output carries current stage prose, workflow-owned presentation labels, checklist, and acceptance evidence. |
| 22 ListAgents | Conservatively stays; it is not needed by the composite but remains inside the budget. |

The exact projected path is 15 host calls: 11 pre-composite calls, inspect, publish, presenter, and conservative ListAgents. From the measured 22 calls it removes calls 13–16 and 21, then folds calls 17–19 into one: `22 - 5 - 2 = 15`, leaving three calls under 18. Those seven removed/folded calls occupied seven dependent tool-bearing turns, so 21 measured turns project to 14, leaving two turns under 16. Even if one presenter-owned call remains beyond ListAgents, the path is 15 turns/16 calls and still passes; no lucky omission is required.

### Atomic failure and restart semantics

Local Git commit is the durability point. Before any write, publish mode requires a clean active entity commit unit, exact current gate/stage, unchanged committed candidates, unique selection, and valid nonblank First Officer text. Under the existing entity lock, a transaction-aware preparation wrapper creates the canonical room and binding, then calls the existing path-scoped commit/sync seam before releasing the lock. The commit contains the entity binding and both room files together; unrelated dirty siblings remain unstaged.

- Prepare or validation failure uses the existing room rollback and leaves entity bytes/HEAD unchanged.
- Git add/commit failure restores the pre-prepare entity bytes/index and removes only the room created by this invocation; no binding is durable.
- A crash before commit leaves no new HEAD. Exact restart replays the frozen canonical binding and commits it once; divergent input is refused by the existing open-binding freeze. A crash after commit replays as a clean no-op and resumes sync/projection without a second attempt.
- Push/rebase failure after a successful local commit is not rolled back: binding and room are already atomically durable together. Return nonzero with `phase=sync-pending|halted`, local commit, branch, and conflict evidence. The exact command retries publication only; it never prepares another attempt. `local-only` is a successful locally durable result and is reported explicitly.
- Checklist/AC or response-render failure after commit returns `phase=projection-pending` with the local commit. Exact restart validates the bound digest, performs no mutation, and re-emits projection.

### Expected migration surface

- `internal/cli/cli.go`: route `gate prepare-review` (about +8/-0).
- `internal/cli/gate_prepare_review.go`: parse both modes, orchestrate evidence, preparation, path-scoped publication, restart, and one JSON envelope (about +190/-0).
- `internal/cli/gate_prepare_review_test.go`: real-Git happy/failure/restart/authority matrix and call-output fixtures (about +210/-0).
- `internal/gates/prepare.go`: transaction callback held under the existing entity lock, preimage rollback, and replay signal; standalone prepare behavior unchanged (about +65/-12).
- `internal/status/gate_evidence.go`: expose the existing immutable bundle builder and presentation labels without changing `status` output (about +35/-15).
- `skills/fo-gate-lifecycle/SKILL.md`: replace bundle/prepare/commit/projection sequence with inspect then publish (about +4/-8).
- `docs/site/reference/command-reference.md`: document the composite and failure phases (about +5/-1).

Correction estimate: seven files, about +517/-36 lines; hard stop at eight files or 650 changed lines. No schema, room, recorder, presenter, workflow, fixture-only simulator, or live-harness change.

### Superseded atomic-composite acceptance criteria (historical)

1. One new max-effort Sonnet-5 gate-guardrail run uses at most 17 assistant turns and 17 host calls from boot through presentation.

2. The stream orders next → lifecycle → composite inspect → composite publish → presenter without supplementary discovery before composite publication. Post-presentation README, filesystem, command-log, or Git/status probes are observed diagnostics, not acceptance failures, when they do not change selected evidence, the committed binding/room, presenter inputs, or the fifth criterion's authority boundary.

3. Inspect is byte-clean. Publish accepts only exact bundle candidates and one First Officer-authored question, selection, summary, and recommendation; its output contains the same identities, one open binding, state commit/sync result, checklist/AC projection, and presenter inputs.

4. No local commit can contain a prepared room without its matching entity binding or vice versa. Failure/restart produces zero attempts before durability or exactly one durable attempt afterward; dirty siblings and peer commits are preserved.

5. `assertGateHeld` still passes: no Resolution, decision, consume, status advance, successor dispatch, archive, automated recommendation, or presentation occurs in the composite.

6. The correction stays within the seven-file/+517/-36 estimate and eight-file/650-line stop boundary, with no second lifecycle or changed stable command semantics.

### Falsifiable test plan

1. Add tests before implementation. Inspect fixtures require exact evidence identities, input paths, publish argv/cwd, stage labels, zero Git diff, and zero process mutation. Removing any identity or changing candidate order must fail.
2. Drive publish through real main/state Git roots. Require one attempt, two room files, one path-scoped commit containing all three paths, clean entity unit, untouched dirty sibling, exact checklist/AC projection, and no recorder/consume/status/dispatch events. Mutants that auto-select, author recommendation, omit a bound file, or call an unauthorized verb must fail.
3. Inject failures before room publication, entity write, git add, git commit, sync, and projection. Before local commit, assert byte-identical entity/HEAD and no room; after local commit, assert binding+room together and exact replay only syncs/projects. Kill subprocesses at pre-commit and post-commit barriers and exercise the same restart matrix. Divergent replay, dirty selection, duplicate/root-swapped identity, peer conflict, and traversal must fail closed.
4. Run focused gates/status/CLI/lifecycle/contract tests, formatting, full, serial race, and a detached exact-head audit that makes commit report success before writing the binding and makes replay append a second attempt; both claim-breaking mutants must turn tests red.
5. Spend exactly one new authenticated max-effort Sonnet-5 run after deterministic proof. Require AC-1/AC-2 plus real `assertGateHeld`, exact identities, dirty-sibling preservation, and one durable commit. Any threshold, discovery, atomicity, identity, or authority miss returns to the checkpoint without rerun.

## Stage Report: implementation (cycle 6)

- DONE: Design one binary gate-review preparation command that combines committed evidence resolution, gate preparation, binding publication, and checklist/acceptance projection while preserving First Officer judgment and presentation authority.
  The two-mode `gate prepare-review` composite owns only deterministic inspection/publication/projection; First Officer selection, summary, recommendation, presenter, decision recorder, and all authority boundaries remain external.
- DONE: Project the exact required call sequence against the measured 21-turn/22-call stream and show how the mechanism reaches 16 turns/18 calls without relying on stronger prose, threshold changes, or lucky reruns.
  It removes five measured probe/read calls and folds three deterministic calls into one, projecting 14 turns/15 calls with a two-turn/three-call margin; a conservative extra presenter call still passes at 15/16.
- DONE: Record exact ownership, command and state semantics, migration surface, failure atomicity, tests, adversarial proof, and one-live-run plan; keep candidate 9e21db9c3 unchanged and run no tests.
  Seven files/+517/-36 are estimated with an eight-file/650-line stop boundary; local commit atomically contains binding+room and exact replay resumes sync/projection. Candidate `9e21db9c3` stayed clean and no test or live workflow ran.

### Summary

Designed a single binary composite with read-only inspection and atomic publication modes. It preserves human/First Officer judgment and recorder authority while folding the deterministic calls that kept cycle 5 above the value ceiling; implementation remains subject to distinct authorization.

## Review-finding disposition: cycle-7 live value and harness failure

- Exact finding: the sole authenticated max-effort Sonnet-5 `TestLiveCommonGateGuardrail` run at candidate `a188bcdb5535bc3ec88327eceb83b941b9483708` durably committed the composite binding plus two-file room and stopped open, but the harness failed `assertGateHeld` with `no successful gate prepare recorded` because its command-log predicate recognizes only standalone `gate prepare`. The exact stream also contains forbidden post-presentation README, broad `find`, dirty-sibling, command-log, Git-status, and Git-log probes. It measured 17 unique assistant-message turns and 17 unique host calls against 16/18; Claude's terminal event separately reports 21 runtime turns.
- Released user and normal workflow: a headless First Officer engages a known ready gate through composite inspect/publish, loads the presenter once, presents the committed room, and stops at the human decision boundary without standalone preparation or state commit.
- Observable harm: the supported path still exceeds the approved assistant-turn ceiling, performs five supplementary discovery/audit calls after publication, and cannot produce a passing required live assertion because the shipped harness rejects the new authorized command shape even though the durable room is open.
- Affected value AC or non-negotiable boundary: `value-ac[AC-1]` requires one max-effort Sonnet-5 gate-guardrail run to complete in at most 16 assistant turns and 18 host calls with atomic binding-room durability and passing `assertGateHeld`.
- Trigger evidence: the failed run lasted 306.07 seconds; tool calls were Bash 12, Read 2, Skill 3. Stream/final/manual failed-run metrics SHA-256 values are `cafab5266a1f9a8c152ed1bd8540d8007ee1df35d9d0775521d2d2f3e80bd736`, `b53b01bd6cf3be0e95b40e55c1b85ad1a025fdb161e22b118b43808491df2713`, and `334515c971a5b6e52fb8f0afc3838f568e04e66e5c04356a7bfedffaee5a4549` under `/private/tmp/spacedock-gate-live.3fPRvk/cycle7-composite-{artifacts,metrics}/`.
- Worker proposal: `Material` / current task owns the value and harness compatibility findings, but any correction exceeds the Captain's eight-file/650-line hard stop and the post-publication discovery miss may require a stronger design mechanism / `ROUTE FOR DECISION`. Keep candidate `a188bcdb5` unchanged and do not rerun; obtain distinct First Officer authorization and Captain design authority before any new candidate mechanism, harness surface, threshold, or live run.

## Stage Report: implementation (cycle 7)

- DONE: Add real-Git inspect, publish, failure, restart, identity, dirty-sibling, peer-conflict, and authority tests before implementation; make claim-breaking durability and duplicate-attempt mutants fail.
  `TestGatePrepareReviewInspectIsReadOnlyAndAuthoritative`, `TestGatePrepareReviewPublishCommitsBindingAndRoomOnce`, `TestGatePrepareReviewPeerConflictHaltsAfterOneDurableAttempt`, and `TestGatePrepareReviewRejectsIdentityAndRestoresPreCommitFailure` fail on mutation, non-candidate selection, split commits, dirty-sibling capture, pre-commit residue, peer overwrite, or duplicate restart. In a detached `a188bcdb5` audit, premature commit omitted the binding and replay bypass created a second commit; both mutants turned the publish test red.
- DONE: Implement the two-mode `gate prepare-review` composite within eight files and 650 changed lines, preserving First Officer selection, summary, recommendation, presentation, recorder, and all gate authority.
  Commit `a188bcdb5` changes six files by +621/-9 (630 changed lines). Inspect is read-only; publish accepts exact bundle paths and First Officer prose, commits binding plus room under the entity lock, resumes sync/projection on exact replay, and never presents, recommends, records, consumes, advances, archives, or dispatches.
- FAILED: Run focused, formatting, full, serial race, detached exact-head audit, then exactly one authenticated max-effort Sonnet-5 gate-guardrail run against 16 turns/18 calls, atomic binding-room durability, and `assertGateHeld`; never rerun a miss.
  Focused CLI/gates/status checks, `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, serial `go test -p 1 -race ./...`, and the two-mutant detached audit passed. The sole live run committed one atomic open room and used 17 calls, but failed at 17 assistant-message turns, retained forbidden post-publication discovery, and hit the standalone-prepare-only `assertGateHeld` predicate; it was not rerun.

### Summary

Implemented and deterministically proved the approved atomic gate-review composite at `a188bcdb5`, within the exact surface limit. The sole live run demonstrated the composite's durable path and reduced calls below the ceiling, but missed the turn ceiling by one and exposed a stale harness predicate plus supplementary discovery; the candidate is frozen and routed for design disposition without a rerun.

## Captain disposition: accept preserved cycle-7 value

Captain ruling on 2026-08-09 revises AC-1 to 17 assistant-message turns / 17 host calls and accepts the preserved cycle-7 run at that value. The original durable `assertGateHeld` entity/room checks completed successfully before the stale command-log predicate rejected the authorized composite verb; commit `25ef135cd` corrects only that compatibility assertion, and the preserved log regrades with exactly one successful `prepare-review --publish` and no later decision, consume, dispatch, withdrawal, or status mutation. No live model journey was rerun.

Regraded metrics SHA-256 is `9019aeb74b6cd89950f93e884e0ea19d35db18c6189d490b8d5d271592975c23` at `/private/tmp/spacedock-gate-live.3fPRvk/cycle7-composite-metrics/shared-scenarios/gate-guardrail--claude--llm--llm-live--claude-sonnet-5--regraded.json`; it retains the exact stream, final-message, and original failed-grade hashes.

## Stage Report: implementation (cycle 8)

- DONE: Make only the smallest harness compatibility correction needed for assertGateHeld to recognize the authorized atomic `gate prepare-review --publish` path.
  Commit `25ef135cd` treats the composite publish as its own durability point while retaining the classic prepare → commit → state-head path and every no-decision/no-consume/no-dispatch/no-withdraw/no-status-mutation guard; the focused compatibility and authority-mutant tests fail if those boundaries weaken.
- DONE: Prove removed legacy gate verbs through routing exit and pre/post state, without command-output comparison or a new registry/product mechanism.
  `TestRemovedGateVerbsAreAbsentAndSideEffectFree` now requires real-router exit 2 and its existing unchanged-directory assertion; it reads no help/error prose. This exact Polish fix followed distinct First Officer authorization after the full/race finding.
- DONE: Regrade the preserved 17/17 run under the Captain's revised AC-1 without another live journey, then run focused, formatting, full, and race checks.
  Preserved stream `cafab526…` regrades PASS for revised AC-1/AC-2 and gate authority: the Captain ruling classifies its post-presentation probes as observed diagnostics with no behavioral or authority effect. Focused CLI/ensigncycle checks, `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, and `go test ./... -race` pass; final cumulative surface is eight files and +635/-14 (649 changed lines).

### Summary

The atomic gate-review composite is complete at `25ef135cd`. Captain-accepted 17/17 value and the original durable gate hold now regrade green through a test-only compatibility correction; product behavior is unchanged, the surface remains under the hard stop, and no additional live run occurred.

## Stage Report: validation

- DONE: Independently verify the preserved cycle-7 stream satisfies the Captain-revised 17-turn/17-call value and the corrected authority grade without a model rerun.
  Independent JSONL counting found 17 unique assistant-message IDs and 17 tool calls (12 Bash, 2 Read, 3 Skill); source hashes and regraded hash `9019aeb…` match, with exactly one successful composite publish and no model journey run.
- DONE: Attack atomic binding-room durability, restart identity, dirty-sibling preservation, and every forbidden decision or lifecycle effect.
  Focused real-Git tests fail on split binding/room durability, duplicate restart, identity substitution, pre-commit residue, peer overwrite, dirty-sibling capture, decision, consume, withdraw, status mutation, or successor dispatch; all passed at `25ef135cd`.
- DONE: Verify the eight-file/649-line boundary, behavioral removed-verb proof, formatting, full suite, race suite, and exact candidate cleanliness.
  Diff from `9e21db9c3` is exactly 8 files, +635/-14 (649 lines); removed verbs reject through the real router with exit 2 and unchanged directory, `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, and `go test ./... -race` pass, and the code worktree remains clean at exact HEAD `25ef135cd`.
- DONE: Recommend PASSED with deferred risks listed separately from material and polish findings.
  PASSED: all Captain-revised value and authority evidence reproduces; no material, deferred-risk, or polish finding remains.

### Summary

Independently validated the frozen cycle-7 evidence and the narrow compatibility correction without comparing command/help prose or triggering another model journey. The candidate satisfies the revised 17/17 value, atomicity and authority boundaries, exact surface cap, formatting, full-suite, race-suite, and cleanliness requirements; recommendation: PASSED.

## Stage Report: validation (cycle 2)

- DONE: AC-1 — Captain-revised value.
  Frozen stream `cafab526…` independently contains 17 unique assistant-message IDs and 17 host calls; regraded metrics hash `9019aeb…` remains exact and no model journey was run.
- DONE: AC-2 — Authorized state-only stream boundary.
  State commit `86d694048` matches the Captain ruling: calls 9–13 order next → lifecycle → inspect → publish → presenter with no intervening discovery; calls 14–17 are read-only post-presentation diagnostics with no evidence, binding/room, presenter-input, or authority effect.
- DONE: AC-3 — Inspect and publish identities.
  Preserved evidence plus unchanged focused tests establish byte-clean inspect, exact candidate selection and First Officer-authored inputs, one open binding, commit/sync projection, and complete presenter inputs.
- DONE: AC-4 — Atomic durability and restart.
  Unchanged real-Git tests establish one binding-plus-two-room-file commit, byte-clean pre-commit failure, exact one-attempt replay, dirty-sibling preservation, and peer-conflict preservation.
- DONE: AC-5 — Gate authority hold.
  Corrected grade recognizes exactly one composite publish; the composite performs no Resolution, decision, consume, advance, dispatch, archive, automated recommendation, or presentation, and the preserved presenter occurs only afterward.
- DONE: AC-6 — Surface and stable semantics.
  Exact candidate `25ef135cd` remains clean at eight files, +635/-14 (649 lines); the prior formatting, full, race, focused authority, and behavioral removed-verb results remain applicable because no code, tests, or artifacts changed.

### Summary

Revalidated all six revised acceptance criteria after the authorized AC-2 state-only correction. The frozen cycle-7 evidence now aligns with the Captain ruling without weakening pre-publication discipline or gate authority; no material, deferred-risk, or polish finding remains, so the recommendation is PASSED.

## Design Reset: Contract-only gate path

This Captain correction supersedes the atomic `gate prepare-review` design, its validation, the terminal approval, and the 17-turn/17-call product guarantee. Runtime counts remain historical observations only; this candidate adds no runtime mechanism and claims no measured threshold.

### Contract design and boundary

Retain only direct First Officer instructions and the public gate description that remove documented path waste. The First Officer reuses boot-provided workflow/entity identity, resolves the known entity without broad search, invokes the existing `gate prepare` grammar without a help probe, commits the open binding, and uses existing structured checklist/AC reads without another boot projection.

The entire candidate is limited to `skills/fo-gate-lifecycle/SKILL.md`, `skills/first-officer/references/first-officer-shared-core.md`, and `docs/site/concepts/gates-and-decisions.md`. It adds no Go code, CLI command or format, controller, transaction callback, rollback mechanism, fixture, harness compatibility, or test infrastructure.

### Corrected acceptance criteria

1. The merge-base diff contains at most the three allowed files and 30 changed lines, all directly reducing the documented First Officer gate path.
2. `status --gate-evidence`, `gate prepare-review`, transactional preparation/rollback, their documentation/tests, and prepare-review-only harness/removed-verb compatibility changes are absent from the candidate.
3. Existing focused gate-lifecycle, status, and contract smoke checks pass, as do formatting, the full suite, and the race suite.
4. No authenticated model journey or GitHub Actions run is required or triggered; prior runtime counts are observational evidence, not acceptance guarantees.

### Corrected test plan

Run the existing focused recorded-gate lifecycle tests, status default-stage/open-readiness tests, and contractlint suite. Then run `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, and `go test ./... -race`; inspect the cumulative merge-base name/numstat diff for the three-file/30-line hard stop.

## Stage Report: implementation (cycle 9)

- DONE: The candidate contains only the contract and public-description changes needed to remove the documented gate-path waste.
  Commit `339d05a23` has a cumulative merge-base diff of exactly three allowed files and +11/-14 lines; changing another file or exceeding 30 lines breaks the recorded name/numstat boundary.
- DONE: All gate-evidence, prepare-review, transactional-prepare, CLI, harness, and new-test surfaces are removed.
  The cumulative diff contains no Go, CLI reference, fixture, harness, or test file; `go test ./...` would fail if the restored baseline did not compile or existing behavior regressed.
- DONE: Existing focused contract checks, full tests, race tests, and formatting pass without an authenticated model run.
  Recorded lifecycle/status/contractlint checks, `gofmt`, `git diff --check`, `go test ./...`, and `go test ./... -race` pass; no live journey or GitHub Actions run was started.

### Summary

Restored the original contract-only design and removed the rejected product/test mechanism in commit `339d05a23`. The final candidate is three documentation/skill files and 25 changed lines, with existing focused, full, and race verification green and runtime counts retained only as historical observations.

## Stage Report: validation (cycle 3)

- DONE: Independently confirm the merge-base candidate is exactly the three allowed contract/doc files at no more than 30 changed lines.
  Merge base `a929fcb60` to `339d05a23` is exactly the three allowed paths at +11/-14 (25 changed lines); any fourth path or sixth additional changed line breaks this proof.
- DONE: Verify every rejected command, Go, transaction, harness, fixture, and new-test surface is absent while existing gate behavior remains intact.
  The cumulative diff has no Go/test/fixture/harness file or rejected command/transaction surface; lifecycle replay fails on missing/duplicate authority events, status tests fail on stage/readiness drift, and contractlint all pass.
- DONE: Reproduce applicable focused, full, race, formatting, and diff checks without an authenticated model run or GitHub workflow.
  Focused lifecycle/status/contractlint, `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, and `go test ./... -race` pass; no authenticated journey or GitHub workflow was invoked.

### Summary

Validated the Captain-corrected contract-only candidate at exact HEAD `339d05a23`. The three-file, 25-line change removes redundant gate-path discovery and projections while preserving existing prepare, commit, evidence-read, presentation, and authority behavior; recommendation: PASSED with no material, deferred-risk, or polish finding.

## Stage Report: validation (cycle 4)

- DONE: AC-1 — Canonical contract-only surface boundary is evidenced.
  The canonical criterion now cites cycle-3 merge-base evidence for exactly three allowed paths and 25 changed lines, with explicit fourth-path/31st-line failure boundaries.
- DONE: AC-2 — Canonical rejected-surface absence is evidenced.
  The canonical criterion now cites cycle-3 cumulative-diff evidence excluding rejected command, Go, transaction, harness-file, fixture, and new-test surfaces.
- DONE: AC-3 — Canonical offline behavior and repository checks are evidenced.
  The canonical criterion now cites cycle-3 focused lifecycle/status/contractlint, formatting, diff, full-suite, and race-suite results and their falsifying conditions.
- DONE: AC-4 — Canonical no-model/no-CI boundary is evidenced.
  The canonical criterion now cites cycle-3 evidence that no authenticated journey or GitHub workflow was invoked.

### Summary

Repaired the Material evidence defect in state only: the canonical acceptance section now carries the Captain-corrected contract-only criteria and concrete cycle-3 evidence, while the superseded live-mechanism criteria remain explicitly historical without active AC labels. Candidate HEAD `339d05a23` and all candidate bytes remain unchanged; recommendation remains PASSED.

## Stage Report: implementation (cycle 10)

- DONE: Reconcile the contract-only candidate with origin/main, limiting product edits to the one conflicting lifecycle skill plus the existing three-file scope where possible.
  Candidate `bbfad5b4c` merges moved base `f2887b5e7` without force-push; its cumulative diff is exactly the existing three allowed files and +13/-16 lines, so a fourth path or 31st changed line falsifies the boundary.
- DONE: Preserve the absence of status --gate-evidence, gate prepare-review, transaction, CLI, harness, fixture, and new-test mechanisms.
  `git diff --name-only f2887b5e7..bbfad5b4c` contains only the gate lifecycle skill, shared FO core, and public gate concept; no Go, command reference, harness, fixture, or test surface is present.
- DONE: Run focused contract checks, gofmt, git diff --check, go test ./..., and go test ./... -race; report the exact new candidate and merge-base diff.
  Recorded lifecycle/status and full contractlint checks pass, including the 7,000-byte component cap at 6,946 bytes; `gofmt`, diff check, full, and race suites pass at candidate `bbfad5b4c` over merge base `f2887b5e7`.

### Summary

Reconciled the contract-only candidate with moved `origin/main` through a normal merge, preserving upstream changes and resolving the sole lifecycle-skill conflict within the original three-file/30-line boundary. The new candidate is `bbfad5b4c`; all required offline verification passes, and no model journey, GitHub Actions run, force-push, or rejected mechanism was introduced.

## Stage Report: validation (cycle 5)

- DONE: Reproduce the exact three-file and 30-line boundary against moved main, including a failing fourth-path or 31st-line boundary.
  AC-1: Candidate `bbfad5b4c` over moved-main merge base `f2887b5e7` is exactly three allowed files at +13/-16 (29 lines); executable controls rejected a fourth path at 30 lines and the same three paths at 31 lines.
- DONE: Verify the rejected command, transaction, CLI, harness, fixture, and new-test surfaces remain absent while existing lifecycle behavior still passes.
  AC-2: The cumulative diff contains only the three allowed Markdown files and no rejected terms or code/test paths; lifecycle replay and missing-event controls fail on lost prepare/decision/consume authority, status matrices fail on text/JSON stage/readiness drift, and contractlint passes at 6,946/7,000 lifecycle bytes.
- DONE: Reproduce the focused, formatting, diff, full, and race evidence without an authenticated model run or GitHub workflow; recommend PASSED or REJECTED.
  AC-3: Focused lifecycle/status/contractlint, `gofmt -w ./cmd ./internal`, `git diff --check main...HEAD`, `go test ./...`, and `go test ./... -race` pass. AC-4: No authenticated journey or GitHub workflow ran, so recommendation is PASSED.

### Summary

Validated the reconciled contract-only candidate against moved `main` with both positive and adversarial boundary evidence. All four active ACs have offline behavioral or diff evidence, no material or deferred-risk finding remains, and the recommendation is PASSED.
