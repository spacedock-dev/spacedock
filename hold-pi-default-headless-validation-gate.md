---
title: Hold the Pi default headless validation gate
status: validation
source: "Staff review M3 for test-behavior-completeness, 2026-08-09"
started: 2026-08-09T20:36:21Z
completed:
verdict:
score: 0.9
sprint: test-behavior-completeness
sprint-readiness: ready
group: common-evidence
worktree: .worktrees/spacedock-ensign-hold-pi-default-headless-validation-gate
issue:
pr:
mod-block:
id: fh6rv0k6wr25zty0jjan4jp7
gates:
    version: 1
    records:
        - id: gate:fh6rv0k6wr25zty0jjan4jp7:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:fh6rv0k6wr25zty0jjan4jp7-backlog-1
              briefing:
                id: briefing:fh6rv0k6wr25zty0jjan4jp7:backlog:attempt-1:revision-1
                digest: sha256:86e1bf356ccc99c08b44e16c2e29dbc399027927b6f42695a64dba0f6ba8f582
                request-digest: sha256:b481c7b7ea716cd7f8fb24b333077ea387e0ac2421928749a261df271f5d6cb7
                room-ref: ./hold-pi-default-headless-validation-gate/review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:fh6rv0k6wr25zty0jjan4jp7:backlog:1
                briefing: briefing:fh6rv0k6wr25zty0jjan4jp7:backlog:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T20:35:54.378022Z"
                decision: approve
                reason: The Captain authorized shaping and requires end-user value; this task owns the real Pi headless gate-stop result.
              application:
                target-stage: ideation
                state: consumed
        - id: gate:fh6rv0k6wr25zty0jjan4jp7:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:fh6rv0k6wr25zty0jjan4jp7-ideation-1
              briefing:
                id: briefing:fh6rv0k6wr25zty0jjan4jp7:ideation:attempt-1:revision-1
                digest: sha256:b7bf620a437c4aba3156522c2eee177d0adb4b1e2ad0236cba43c1f08f04e063
                request-digest: sha256:a538be68b35505144f0506819b7dc69eb0b5d493597c61b71b5636cbbfb7b3c6
                room-ref: ./hold-pi-default-headless-validation-gate/review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:fh6rv0k6wr25zty0jjan4jp7:ideation:1
                briefing: briefing:fh6rv0k6wr25zty0jjan4jp7:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-09T21:33:26.164844Z"
                decision: approve
                reason: Captain approved the Pi headless open-validation-gate result.
              application:
                target-stage: implementation
                state: consumed
---

## Problem

The unchanged Pi `default-headless-gate-stop` journey reaches the validation gate with terminal fields set.

The durable oracle has one first branch for four final-state facts. It checks that the entity changed, that its status is validation, and that `completed` and `verdict` are empty. The branch reports one generic error.

The live Pi artifact at `/private/tmp/spacedock-xp6-artifacts-pi-default-headless/pi-common/default-headless-gate-stop/` ran for `10m37.753288125s` without a process timeout. The root session read this final entity:

```text
status: validation
completed: true
verdict: approved
```

The entity changed. Its status is validation. The gate binding is present and open. The first oracle branch therefore fails on `completed` and `verdict`. The stable semantic code for this branch is `gate-hold-terminal-fields-set`.

The Pi implementation worker wrote the terminal fields in its stage report. The shared ensign contract already forbids worker frontmatter writes. The live run proves that worker prose alone does not protect the First Officer boundary.

## Value

A headless Pi run stops at the first open validation gate with a nonterminal entity and durable evidence.

The final entity has `status: validation`, empty `completed`, empty `verdict`, one open Briefing, no Resolution, no Application, and no successor dispatch.

The current artifact has two non-empty terminal fields. The repaired journey has zero. This final-state count is independent of the repair mechanism.

## Dependencies and order

Task `ts7gq0mr9s3chx2w4wppd1kt` must land first. It supplies strict-XFAIL grading and the typed semantic error path.

Task `98aa776adg66gn823a8gamdq` must land next. It repairs the Sonnet and Codex implementation-worker boundary. This task does not change that repair or assume its mechanism for Pi.

This task owns the Pi strict-XFAIL binding under `fh6rv0k6wr25zty0jjan4jp7`. The runtime label `xp6c9qfe7y4wwp46enc3f85n` is evidence-only.

Run the current Pi journey as target-level XFAIL for owner `fh6rv0k6wr25zty0jjan4jp7` after those dependencies. Start Pi product work only after this semantic XFAIL.

## Proposed approach

### 1. Normalize accidental terminal fields at the Pi First Officer boundary

After the Pi worker returns and its stage report passes, read the entity frontmatter before the First Officer advances the status.

If the current stage is nonterminal and `completed` or `verdict` is non-empty, run the existing command once:

```text
spacedock status --workflow-dir <workflow> --set <slug> completed= verdict=
```

Read the entity again after the command. If the command fails, or either field remains non-empty, stop without gate preparation. Do not use `--force`.

Only after the fields are empty can the First Officer run the existing status transition, state commit, gate preparation, gate commit, and presentation.

This repair is First Officer-owned state cleanup. It does not record a decision, consume a gate, or dispatch a successor.

### 2. Keep the Pi runtime rule explicit

Add the boundary rule to `skills/first-officer/references/pi-first-officer-runtime.md`. The rule names the two fields, the one clear operation, the reread, the fail-closed result, and the `--force` prohibition.

Do not change the shared ensign contract in this task. It already says that a worker updates the entity body and does not modify YAML frontmatter. The observed run makes a First Officer boundary guard necessary.

### 3. Keep the final-state semantic code in the metric

Refine the first `assertGateHeld` branch into separate semantic codes. Use `gate-hold-terminal-fields-set` only when `completed` or `verdict` is non-empty. Keep separate codes for an unchanged entity and a non-validation status.

The existing gate-binding and open-attempt checks remain semantic controls. The Pi strict-XFAIL binding for owner `fh6rv0k6wr25zty0jjan4jp7` names only the terminal-field code.

### 4. Prove the cleanup order in the existing command log

Allow one `completed= verdict=` cleanup line before the validation status change. Require it before the successful gate preparation. Reject a duplicate cleanup, a cleanup after preparation, or a status change after preparation.

This command-log check proves the repair boundary. The final-state assertion proves the value.

### 5. Use the Pi target-level XFAIL binding

Use `liveXFail("pi", "fh6rv0k6wr25zty0jjan4jp7")` for `default-headless-gate-stop`.

Keep the binding until the exact repaired candidate passes. A repaired run must report XPASS while the binding remains. Remove the binding only after the passing run.

## Alternatives and value mapping

| Mechanism | Value AC | Simplest alternative | Why the alternative fails |
|---|---|---|---|
| Pi First Officer cleanup and reread | AC-1 and AC-3 | Add another worker instruction | The existing worker instruction already forbids frontmatter writes. The live child still wrote the fields. |
| Stable final-state code | AC-2 | Match the full error text | Model output and assertion text can change without a semantic change. |
| Command-log order check | AC-3 | Check only the final entity | A final entity cannot prove that cleanup happened before gate preparation. |
| Strict Pi XFAIL binding for owner `fh6rv0k6wr25zty0jjan4jp7` | AC-2 | Keep Pi as TODO | TODO skips the journey and hides a repaired or changed failure. |
| Existing `status --set` clear operation | AC-1 | Add a new normalize command | A new command changes grammar and adds an unnecessary surface. |
| No `--force` path | AC-1 and AC-3 | Bypass all status guards | `--force` can bypass merge-hook protection and hide a workflow error. |

## Acceptance criteria

**AC-1 (VALUE) — The unchanged Pi journey stops with a clean open gate.**

The exact `default-headless-gate-stop` fixture ends at `status: validation` with empty `completed` and `verdict`, one open expected Briefing, no Resolution, no Application, and no successor dispatch. The independent baseline has two non-empty terminal fields.

Verified by: run the exact Pi live selector and inspect the entity, gate room, command log, and journey metric.

**AC-2 — The known Pi failure has one strict semantic code.**

The current candidate runs the real fixture and reports XFAIL with exactly `gate-hold-terminal-fields-set` for owner `fh6rv0k6wr25zty0jjan4jp7`. A repaired candidate reports XPASS and fails the lane while the binding remains.

Verified by: run the strict-XFAIL selector before the Pi repair, then run it again after the repair and remove the binding only after PASS.

**AC-3 — Pi cleanup is bounded and stays before gate authority.**

The command log has at most one successful terminal-field cleanup. The cleanup precedes the validation status transition and gate preparation. No decision, consume, withdrawal, or successor dispatch follows preparation.

Verified by: extend the existing command-log oracle with duplicate, reordered, post-prepare, and missing-cleanup controls. Run the exact live selector.

**AC-4 — The Sonnet and Codex repair remains separate.**

The Sonnet and Codex default-headless journeys retain their `98a` worker-dispatch binding and behavior. The Pi cleanup rule does not change their runtime adapters or the shared worker envelope.

Verified by: run the exact Sonnet and Codex selectors after `98a`, then run the offline and race suites.

## Test plan

First run the focused oracle tests. The `assertGateHeld` controls must distinguish unchanged state, non-validation status, terminal fields, wrong gate binding, resolved attempts, and applied attempts. The terminal-field mutant must return exactly `gate-hold-terminal-fields-set`.

Next run the command-log tests. A clear operation before validation is valid. A duplicate, a post-prepare clear, and a status change after prepare are invalid. The missing-worker branch keeps its existing `implementation-worker-not-dispatched` code.

Then run the Pi strict-XFAIL selector on the current candidate:

```bash
SPACEDOCK_LIVE_RUNTIME=pi \
  go test -tags live -count=1 -timeout 40m \
  -run '^TestLiveCommonDefaultHeadlessGateStop$' ./internal/ensigncycle -v
```

The expected result is one executed XFAIL for owner `fh6rv0k6wr25zty0jjan4jp7`. A skip or infrastructure error fails the proof.

After the Pi runtime rule lands, run the same selector. Expect XPASS failure while the strict binding remains. Remove the binding and run the selector again. Expect PASS with the final-state count at zero.

Run the Sonnet and Codex selectors after `98a`. They must keep their own strict code and worker-dispatch evidence. Finish with `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal`.

No new fixture is needed. The existing `recorded-gate/pre-gate` fixture and `runGateStopScenario` already exercise the unchanged journey.

## Expected surface and estimate

The product behavior must touch the Pi runtime contract. The proof must touch the existing oracle and strict-XFAIL binding. A fixture-only or assertion-only change is out of scope.

| File | Purpose | Gross insertions | Gross deletions | Net |
|---|---|---:|---:|---:|
| `skills/first-officer/references/pi-first-officer-runtime.md` | Bound cleanup, reread, and fail-closed Pi behavior | 8 | 3 | +5 |
| `internal/ensigncycle/gate_assert_impl_test.go` | Stable final-state semantic code | 5 | 2 | +3 |
| `internal/ensigncycle/claude_runtime_helpers_test.go` | Cleanup order and command-log controls | 5 | 3 | +2 |
| `internal/ensigncycle/shared_live_runner_test.go` | Pi strict-XFAIL source binding | 3 | 1 | +2 |
| **Total** |  | **21** | **9** | **+12** |

The tolerance is one file and plus or minus 14 net lines. No new file, fixture, CLI command, or stored format belongs in the surface.

## Semantic scope

- Command grammar: unchanged. The repair uses the existing `status --set` field syntax.
- Stored formats: unchanged. Empty `completed` and `verdict` remain the nonterminal representation.
- Authority: unchanged. The First Officer owns frontmatter and gate state. The worker owns the stage body and report.
- Runtime behavior: changed only for Pi after worker completion and before a nonterminal status advance. The Pi run still presents one open gate and stops without a decision.
- Host behavior: Sonnet and Codex remain owned by `98a`. The Pi rule does not alter their dispatch envelope.

## Riskiest mechanism spike

The existing clear operation was exercised on disposable copies before this plan selected it.

First, a copy of `docs/dev` retained `_mods/pr-merge`. With `status: validation`, `completed: 2026-08-09T00:00:00Z`, and `verdict: approved`, the clear command exited 1 because the merge hook was not run.

Second, a copy without `_mods` used the same fields and status. The exact command exited 0 and printed:

```text
completed: 2026-08-09T00:00:00Z ->
verdict: approved ->
```

The post-state kept `status: validation` and left both fields empty. The implementation therefore fails closed on a refused clear and does not use `--force`. The unchanged live fixture has no merge hook, so this existing command path passes its relevant guard.

## Documentation diff

Apply this wording change to `skills/first-officer/references/pi-first-officer-runtime.md`.

Before:

```text
`«completion-signal»`: For `pi-subagents`, the PRIMARY completion signal is the child return/status (`status: completed`); an optional advisory is only a non-blocking heads-up via raw `intercom send` before return (`contact_supervisor` carries no completion reason). file verification remains the completion gate: the FO reads the entity file and verifies the stage report before advancing state. For `pi-agent-teams`, task/member completion is likewise verified against the entity file.
```

After:

```text
`«completion-signal»`: For `pi-subagents`, the PRIMARY completion signal is the child return/status (`status: completed`); an optional advisory is only a non-blocking heads-up via raw `intercom send` before return (`contact_supervisor` carries no completion reason). Before a nonterminal status advance, the FO reads `completed` and `verdict`. If either field is non-empty, the FO runs `status --set <slug> completed= verdict=` once, reads the fields again, and stops if either field remains non-empty. The FO does not use `--force`. File verification remains the completion gate: the FO reads the entity file and reads the stage report before advancing state. For `pi-agent-teams`, apply the same field rule after task/member completion.
```

This text documents the Pi operator-visible stop rule. It does not change the CLI grammar or the entity format.

## Evidence

The root Pi transcript ends with a gate presentation and a stop message. The durable final-state read shows `status: validation`, `completed: true`, and `verdict: approved`. The nested Pi worker transcript contains the write that introduced those fields. The runtime label `xp6c9qfe7y4wwp46enc3f85n` is evidence-only and is not the task owner. The current shared contract forbids that write, so the evidence supports a First Officer boundary repair rather than a fixture relaxation.

The live process completed without timeout. The expected open Briefing and no-decision result are already exercised by `runGateStopScenario`; only the terminal-field clause fails on this candidate.

## Scope

- Repair the Pi `default-headless-gate-stop` final-state boundary.
- Preserve the existing fixture, gate storage, command grammar, and worker envelope.
- Keep the strict-XFAIL binding for owner `fh6rv0k6wr25zty0jjan4jp7` until the repaired Pi candidate passes.
- Run after target-level XFAIL and `98a` evidence.

## Out of scope

- Repairing the Sonnet or Codex worker-dispatch gap.
- Changing the Pi `gate-guardrail` journey.
- Relaxing `assertGateHeld` or deleting the terminal-field assertion.
- Adding a fixture, simulator, new CLI verb, or new stored field.
- Using `--force` to bypass a workflow guard.

## Stage Report: ideation

- DONE: Identify the exact failed Pi final-state clause before selecting a repair.
  The final entity had `status: validation` but also `completed: true` and `verdict: approved`. The exact clause is the terminal-field branch, coded as `gate-hold-terminal-fields-set`. The strict-XFAIL owner is `fh6rv0k6wr25zty0jjan4jp7`.
- DONE: Define the smallest change that passes the unchanged headless gate-stop journey.
  The Pi First Officer clears accidental terminal fields once with the existing status command, rereads the entity, fails closed on refusal, and then uses the existing gate path. The command-log oracle and the `fh6rv0k6wr25zty0jjan4jp7` strict binding prove the boundary. The runtime label `xp6c9qfe7y4wwp46enc3f85n` remains evidence-only.
- DONE: Give a visible-value statement and gross and net line estimates.
  The visible value is two non-empty terminal fields reduced to zero while one open Briefing remains. The estimate is 21 gross insertions, 9 gross deletions, and +12 net lines across four existing files.

### Summary

Ideation isolates the Pi failure to terminal fields set on an otherwise open validation gate. The plan adds a bounded First Officer cleanup, a stable semantic code, command-order proof, and a strict-XFAIL Pi binding owned by `fh6rv0k6wr25zty0jjan4jp7`. It leaves the fixture unchanged and makes the final-state journey measurable.

## Captain Decision: implementation design reset

The Captain approved a bounded reset after exact Pi run `31382215343` exposed two oracle errors.

- Change only `internal/ensigncycle/claude_runtime_helpers_test.go` and `internal/ensigncycle/pi_shared_live_runner_test.go` in this reset.
- Limit the reset to 37 insertions and 3 deletions. Limit the full candidate to six files, 84 insertions, and 24 deletions.
- Keep cleanup conditional. If cleanup is present, require one successful cleanup before validation.
- Append the exact root Pi session JSONL. Require a correlated implementation spawn, run, and completion before validation.
- Keep missing-completion, dirty-final-state, duplicate, reordered, and post-prepare controls red.
- Preserve the product, fixture, command grammar, authority, n28 bindings, and Claude and Codex branches.

The end value remains a clean open Pi validation gate with correlated implementation lifecycle evidence.

## Captain Decision: exact Pi XPASS accepted

The Captain accepted protected run `31386886972` and artifact `9062678223` on head `ae1c20ae780b43c3e8147ab2526f52e41aa6a6fa`.

- The exact target reported `XPASS pi/default-headless-gate-stop owner=fh6rv0k6wr25zty0jjan4jp7 observed=[]`.
- The measured result was `xpass` in `528147 ms` with model `openai/gpt-5.6-luna:max`.
- Implementation run `4277388b-74fa-48d0-814c-d759ae26480a` reached `State: complete` with a child session before validation.
- The final entity had `status: validation`, empty `completed`, empty `verdict`, and one open gate awaiting the Captain.

The Captain authorized removal of only the fh6 Pi binding and its reconciliation entry.

## Material Finding: Pi rejection-flow timeout

- **Released user:** The Pi rejection-flow maintainer cadence.
- **Harm:** The supported journey times out before second-gate completion and terminal state.
- **Authority:** The zh rejection-flow value target and its required live proof.
- **Trigger:** Run `31388978952`, artifact `9064353436`, timed out after `769.17s`, after second validation PASSED at an open gate.

This finding is not owned by fh6. The Captain accepted the fh6 value proof from the same run.

## Stage Report: implementation

- DONE: Run the exact current Pi target first and retain its result.
  Local attempts retained two exact infrastructure blockers. Protected run `31382215343` then retained the target-level XFAIL evidence.
- DONE: Implement only the approved Pi boundary and proof surface.
  The candidate adds conditional cleanup checks, stable codes, command-order controls, and correlated Pi lifecycle evidence.
- DONE: Preserve the fixture, product boundaries, and other owners.
  The command grammar, stored format, gate authority, fixture, n28 bindings, and Claude and Codex branches remain unchanged.
- DONE: Prove the conditional cleanup and final-state rules.
  Focused controls reject dirty final state, duplicate cleanup, reordered cleanup, post-prepare cleanup, and missing completion.
- DONE: Prove XPASS, remove only fh6 binding, and prove normal PASS.
  Run `31386886972` and artifact `9062678223` reported fh6 XPASS with `observed=[]`.
  Run `31388978952` and artifact `9064353436` reported normal PASS in `475830 ms`.
- DONE: Run the required local verification.
  Focused gate, lifecycle, registry, owner, format, and XPASS-policy checks passed.
  The full and race suites passed before the final binding-only change, as approved.
- DONE: Stop for and follow the approved design reset.
  The reset stayed within its 37-insertion and 3-deletion increment.
  The pre-binding-removal candidate stayed within six files and the approved gross and net limits.
- DONE: Commit and push the exact candidate and durable report.
  The clean candidate is `cae0a19d91f3b8b8f129de83fad65613dd6894aa`, rebased onto `0175d4399a9259acc5c3311ea467d972f1c8351d`.

### Summary

PASSED. Pi now reaches a clean, open validation gate with a correlated implementation lifecycle.
The exact target passes without the fh6 binding. A later zh-owned rejection-flow timeout does not change the fh6 result.
