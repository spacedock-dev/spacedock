---
title: Cut gate-guardrail turn and tool bloat
status: implementation
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
mod-block:
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

## Acceptance criteria

**AC-1 (VALUE) — The real gate journey requires materially less operator effort.**

One candidate run of `TestLiveCommonGateGuardrail` with Claude Sonnet 5 at maximum effort uses at most 16 assistant turns and 18 host tool calls, measured from its archived Claude stream with the same journey-metrics collector used for baseline run `31297186020` (22 turns, 24 calls).

**AC-2 — Known search and retry waste is absent.**

The candidate stream contains no filesystem search for the known workflow/entity/contract path, no `gate ... --help` or failed command-shape probe, no broad/repeated Git or worktree-status inspection, and no boot projection after successful `gate prepare`.

**AC-3 — Gate authority and final state remain correct.**

The candidate's durable split-root state contains exactly one committed current-stage open attempt and its two-file room, preserves the dirty sibling and all preexisting body bytes, and contains no Resolution/decision, withdrawal, application consumption, status advance, successor dispatch, archive, or other unauthorized mutation. The captain-facing review names the bound Briefing/digest and ends at the human decision boundary.

**AC-4 — The fix changes only the owning contract and its public description.**

The candidate diff stays within the two-file/LOC/byte budget above; no CLI grammar/output, Go code, harness controller, fixture, lifecycle state, standing check, simulator, or static prose-presence test is added.

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
