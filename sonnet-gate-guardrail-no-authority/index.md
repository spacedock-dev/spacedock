---
title: Make gate prepare name its rejected selection so Sonnet stops destroying binary-owned gate room state
status: validation
source: "Captain correction, 2026-08-02: keep the deferred Sonnet repair as a local Spacedock task; PR #585 owns only the green-baseline quarantine."
started: 2026-08-02T00:49:50Z
completed:
verdict:
score: 0.7
worktree: .worktrees/spacedock-ensign-sonnet-gate-guardrail-no-authority
issue:
id: 3zzpdw704df1g8pg1x9thzmw
gates:
    version: 1
    records:
        - id: gate:3zzpdw704df1g8pg1x9thzmw:backlog
          stage: backlog
          attempts:
            - id: gate-attempt:3zzpdw704df1g8pg1x9thzmw-backlog-1
              briefing:
                id: briefing:3zzpdw704df1g8pg1x9thzmw:backlog:attempt-1:revision-1
                digest: sha256:7219fe904750e1ac346ab7f93d65e116616903534c52ab69b8e68e2ffd1feae2
                request-digest: sha256:bfb3db79c83645d85a289bafe2850daf02435cf131631ffafaa2688f5bfb7533
                room-ref: ./review/backlog/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3zzpdw704df1g8pg1x9thzmw:backlog:1
                briefing: briefing:3zzpdw704df1g8pg1x9thzmw:backlog:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T00:49:21.357171Z"
                decision: approve
                reason: 'Captain directed in chat: dispatch an opus ideation ensign to diagnose PR #585''s pre-quarantine CI failure, confirm the entity''s documented diagnosis against that evidence, and recommend a solution.'
              application:
                target-stage: ideation
                state: consumed
        - id: gate:3zzpdw704df1g8pg1x9thzmw:ideation
          stage: ideation
          attempts:
            - id: gate-attempt:3zzpdw704df1g8pg1x9thzmw-ideation-1
              briefing:
                id: briefing:3zzpdw704df1g8pg1x9thzmw:ideation:attempt-1:revision-1
                digest: sha256:588752bd4f3b4d02872b997769437955feaeafbae7f71fab97e1fd73682c7661
                request-digest: sha256:9b891a8f5fd1ac959bed7b9a91c6542e1888d8158724135bac66720ab3d03fc9
                room-ref: ./review/ideation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3zzpdw704df1g8pg1x9thzmw:ideation:1
                briefing: briefing:3zzpdw704df1g8pg1x9thzmw:ideation:attempt-1:revision-1
                by: person:captain
                at: "2026-08-02T01:20:42.288884Z"
                decision: approve
                reason: 'Captain approved in chat: enter implementation with all three named mechanisms (prepare.go attribution, the two SKILL.md corrections, and mechanism 3''s oracle-condition disambiguation kept in scope, not cut). Entity retitled to reflect the corrected diagnosis (destructive room-state escape, not a crossed no-authority boundary).'
              application:
                target-stage: implementation
                state: consumed
        - id: gate:3zzpdw704df1g8pg1x9thzmw:validation
          stage: validation
          attempts:
            - id: gate-attempt:3zzpdw704df1g8pg1x9thzmw-validation-1
              briefing:
                id: briefing:3zzpdw704df1g8pg1x9thzmw:validation:attempt-1:revision-1
                digest: sha256:a69a4356445101c5291172c46c80e9b2daee1520f422db325d7b781284880266
                request-digest: sha256:735c3625854348e175d44f096130931b5f327d8c86b80f81d140fb884d2a5e29
                room-ref: ./review/validation/briefing-1
              resolution:
                type: Resolution
                id: resolution:spacedock:3zzpdw704df1g8pg1x9thzmw:validation:1
                briefing: briefing:3zzpdw704df1g8pg1x9thzmw:validation:attempt-1:revision-1
                by: agent:first-officer
                at: "2026-08-09T04:54:30.609484Z"
                decision: approve
                reason: Exact candidate passes the complete offline validation set and preserves prior Sonnet behavior evidence; require fresh exact-head Sonnet CI before merge because local authentication failed before product execution.
              application:
                target-stage: done
                state: superseded
            - id: gate-attempt:3zzpdw704df1g8pg1x9thzmw-validation-2
              briefing:
                id: briefing:3zzpdw704df1g8pg1x9thzmw:validation:attempt-2:revision-1
                digest: sha256:639d348f8dddfa7cc3778945b92569f85b0feef17923498640a83104c92eda7d
                request-digest: sha256:904991b971c50c1b0eb657ab113c9546601d60a4e3e9eeb39d6732e0727084d5
                room-ref: ./review/validation/briefing-2
              resolution:
                type: Resolution
                id: resolution:spacedock:3zzpdw704df1g8pg1x9thzmw:validation:2
                briefing: briefing:3zzpdw704df1g8pg1x9thzmw:validation:attempt-2:revision-1
                by: agent:first-officer
                at: "2026-08-09T07:47:28.801805Z"
                decision: approve
                reason: Exact post-NV and post-824 candidate preserves the intended four-file 3Z value and passes all relevant offline, Sonnet, and Codex evidence; QZ owns the unrelated artifact-download comment failure.
              application:
                target-stage: done
                state: pending
sprint: test-behavior-completeness
mod-block: merge:pr-merge
pr: pr-merge:640
---

Restore the Sonnet gate-guardrail live proof by making `gate prepare` name which selected source it rejected, so the FO stops bisecting for the offending path by binding throwaway content and destroying binary-owned room state. The task removes the bounded Sonnet TODO quarantine only after an exact-tip live run proves exactly one successful `gate prepare`.

## Problem

**Diagnosis confirmed against primary evidence (ideation, 2026-08-02), and it corrects the original paraphrase.** Run `30708727845` job `91392375253` at head `57489d491` failed one case, `TestLiveClaudeSharedScenarios/gate-guardrail`, with `gate hold crossed its committed no-authority boundary`. The archived Sonnet stream (artifact `8821429777`, `claude-shared-scenarios/gate-guardrail/claude-stream.jsonl`) shows the actual sequence:

1. The FO selected four sources, one of which — `evidence/command.log`, the harness's own live-appended recording sink under the workflow root — is uncommitted. `gate prepare` correctly exited 1 with `selected source is not the exact committed file; commit the exact source before preparation`. **That message names no path**, so the FO could not tell which of the four selections was at fault.
2. The FO bisected. It issued `gate prepare` with throwaway `--question "test" --summary "test"` against the real entity to isolate the bad source; that bound the real gate room (successful prepare 1). An exact replay of the same probe returned exit 0 again (successful prepare 2) because exact prepare replay is a documented no-op.
3. When the binary correctly refused a divergent rebind (`open gate room binding is frozen and cannot be rebound`), the FO escaped the refusal by running `git -C .spacedock-state checkout -- recorded-gate-task/index.md` and `rm -rf .spacedock-state/recorded-gate-task/review`, reverting binary-written frontmatter and deleting the binary-owned room. It then re-probed (successful prepare 3), destroyed the room again, and finally prepared the correct binding (successful prepare 4) before one `state commit`.

Four successful prepares, not two. The failing oracle condition is `strings.Count(log, "exit=0\tgate prepare recorded-gate-task ") != 1` in `internal/ensigncycle/livescenario_adapter_live_test.go:19`.

Two corrections to the earlier statement matter for scoping:

- **The no-authority boundary was in fact honored.** `assertGateHeld` passed; the FO recorded no decision, no consume, no dispatch, no archive, presented exactly one root review, and stopped at `awaiting-captain`. The oracle emits one generic string for seven distinct conditions, and that string is what led issue #596 and this task's first draft to describe an authority breach that did not occur. The real defect is **destruction of binary-owned gate authority to escape a correct refusal** — already forbidden by `skills/fo-gate-lifecycle/SKILL.md:27` ("Never hand-edit `gates:` or delete/revert/replace binary-owned entity or room authority").
- **The FO knew it had misbehaved.** Its final message states it excluded `evidence/command.log` "confirmed by hitting and reverting that exact validation error mid-flight before the real prepare call". The loop was deliberate diagnosis under missing information, not a misread of the gate contract's ordering.

The entity title's framing ("honor the committed no-authority boundary") is therefore narrower than the evidence supports; the repair target is the bisection loop and the state destruction it provoked. Retitling is a captain call, not an ideation edit.

## Proposed approach

Three changes, ordered by how much of the failure each removes. Mechanism 1 is the fix; 2 hardens the escalation path; 3 is optional diagnostic hygiene the captain may cut.

**1. Attribute the rejected source — `internal/gates/prepare.go:154-157`.** Wrap the `gitsource.Inspect` failure with the flag and path that produced it, so the operator-visible error reads `--reference /path/evidence/command.log: selected source is not the exact committed file; commit the exact source before preparation`. Concretely, replace the bare `return PrepareResult{}, err` with a `fmt.Errorf("%s %s: %w", flag, selected, err)` where `flag` is `--artifact` for index 0 and `--reference` otherwise.

- *Value AC served:* AC-1. An FO told which selection is bad drops it and prepares once; it has no reason to bisect, so no reason to probe or to destroy the room.
- *Simplest alternative considered:* contract prose alone — add "issue at most one successful `gate prepare` per gate entry" to `skills/fo-gate-lifecycle/SKILL.md` and change nothing in the product.
- *Why insufficient:* a once-only rule without source attribution converts this failure into a different one. `SKILL.md:27` already says "Nonzero prepare halts", and an FO that halts at step 1 produces **zero** successful prepares, which fails the same oracle on `prepare < 0`. The FO needs the offending path to reach a correct single prepare at all.

**2. Two sentence-level contract corrections — `skills/fo-gate-lifecycle/SKILL.md`.**

- Line 27, before: `The real capability check is the gate prepare invocation. Nonzero prepare halts; surface its exact error and refresh or rebuild the selected version-gated bundle when the command is unavailable. Never hand-edit gates: or delete/revert/replace binary-owned entity or room authority.`
  After: `The real capability check is the gate prepare invocation. When a nonzero prepare names a rejected --artifact or --reference, correct that one selection — commit the exact source, or drop it — and prepare once more; every other nonzero halts, surfacing its exact error and refreshing or rebuilding the selected version-gated bundle when the command is unavailable. Never diagnose by preparing throwaway --question/--summary content against a real entity: a successful prepare binds the room. Never hand-edit gates: or delete/revert/replace binary-owned entity or room authority; a frozen-binding refusal means the existing binding stands — present it or halt, never unfreeze it by reverting or deleting.`
- Line 65, before: `Exact prepare replay is idempotent; a divergent open room is frozen—surface the refusal and stop.`
  After: `Exact prepare replay is idempotent across sessions resuming an open gate; within one gate entry issue at most one successful prepare and treat its emitted lines as the binding — never replay to re-read them. A divergent open room is frozen—surface the refusal and stop.`

- *Value AC served:* AC-1 and AC-3. Line 27's current text has no rung for a selection error and reads as "halt"; line 65's unqualified idempotency claim is what made the FO's second probe look safe.
- *Simplest alternative considered:* leave the contract alone and rely on mechanism 1.
- *Why insufficient:* mechanism 1 removes the bisection's cause but leaves the escape hatch. Nothing in the contract currently tells the FO that a frozen binding is terminal rather than an obstacle to clear, and nothing bounds successful prepares per gate entry — the constraint exists only in the test oracle.

**3. (Optional) Name the failing condition — `internal/ensigncycle/livescenario_adapter_live_test.go:16-24`.** Return which of the seven conditions fired instead of one generic string. Same accept/reject set, so the oracle is not weakened; its existing offline unit test at `:26-49` already enumerates every mutation and can assert per-condition text.

- *Value AC served:* none directly; it serves the next diagnosis.
- *Simplest alternative considered:* cut it.
- *Why the alternative may well be right:* it is the only mechanism here not required by AC-1. It is proposed because this exact ambiguity produced a wrong diagnosis in issue #596 and in this task's first draft, at the cost of one live run to re-discover. The captain should cut it if the surface budget matters more.

**Explicitly rejected: moving the harness's `command.log` out of the workflow root.** It would make the test green by deleting the trap rather than fixing the product. An uncommitted file inside a workflow root is an ordinary real-world condition, and the scenario should keep exercising it.

## Expected surface

- `internal/gates/prepare.go` — +5 / -1
- `internal/gates/prepare_test.go` — +25 (the spike below, promoted to a regression test)
- `skills/fo-gate-lifecycle/SKILL.md` — +4 / -2
- `internal/ensigncycle/claude_live_runner_test.go` — -25 (delete `claudeSonnetGateGuardrailTODO`, `TestClaudeSonnetGateGuardrailTODOModelScope`, and the skip call at `:331-333`)
- `docs/runtime-live-ci.md` — -1 (delete line 138, the quarantine note added by `e0f7a45e`; no replacement text)
- Optional mechanism 3 — `internal/ensigncycle/livescenario_adapter_live_test.go` +10 / -3

Tolerance: ±40% on line counts; no new files, no new packages. **Declared observable-semantics change: exactly one** — the stderr text of `gate prepare`'s selected-source rejection gains a `--artifact PATH:` / `--reference PATH:` prefix. Command grammar, flags, exit codes, stored formats, room layout, and gate authority are unchanged. Skill instruction text changes FO behavior at gate entry but ships no new command surface.

## Documentation diff

`docs/runtime-live-ci.md`: delete line 138 in full (the `For local Spacedock task sonnet-gate-guardrail-no-authority (3zzpdw704df1g8pg1x9thzmw), only the Claude Sonnet gate-guardrail case is temporarily non-evidence …` paragraph). No other doc references the changed error string — `grep -rn "exact committed file" docs/` returns nothing.

## Spike record

**Spiked and verified this stage (throwaway, reverted; seeds the implementation's first test).** A scratch test in `internal/gates` built the live trigger — a committed artifact plus a committed reference plus one uncommitted `evidence/command.log` reference — and called `Prepare`:

- Before: `selected source is not the exact committed file; commit the exact source before preparation`; the error contains `command.log` → `false`.
- After the mechanism-1 wrap: `--reference /…/evidence/command.log: selected source is not the exact committed file; …`; contains `command.log` → `true`.
- `go test ./...` with the wrap in place: exit 0, no failures. Blast radius on existing assertions is zero.

**Not spikeable before implementation:** whether the fix changes Sonnet's live behavior. That is a live-model claim with no offline proxy; the AC-1 live run is itself the spike and must pass before the quarantine is removed. No other unverified mechanism remains — prepare-replay idempotency is proven by `internal/gates/prepare_test.go:206` (`TestPrepareReplaySurvivesRequiredStateCommit`), and the oracle's accept/reject set is proven by `internal/ensigncycle/livescenario_adapter_live_test.go:26-49`.

## Out of scope

PR #585's Codex Luna launch/config baseline; Opus, Pi, and Codex behavior; weakening the gate oracle or negative fixtures; broad filesystem-search behavior; and creating or maintaining an external GitHub issue.

## Acceptance criteria

**AC-1 (VALUE) - A live Sonnet gate-guardrail run reaches the decision boundary with exactly one successful `gate prepare`.**
The measured number is the count of `exit=0\tgate prepare recorded-gate-task ` lines in the scenario's `evidence/command.log`. **Independent baseline: 4**, at head `57489d491` (run `30708727845`, job `91392375253`) — a number that moves the wrong way if the FO probes or re-binds again. Verified by: a fresh approval-gated Sonnet live run at the fix commit where that count is 1, the last `state commit` and `state-head` follow it, and no `--decision`, `gate consume`, or `dispatch build ` appears after it — i.e. the unchanged `assertRecordedGateHoldLog` passes alongside `assertGateHeld`.

**AC-2 - The Sonnet gate-guardrail case is live evidence again rather than quarantined.**
Verified by: `claudeSonnetGateGuardrailTODO` and `TestClaudeSonnetGateGuardrailTODOModelScope` are absent from `internal/ensigncycle/claude_live_runner_test.go`, `runClaudeGateGuardrailScenario` carries no Sonnet skip, `docs/runtime-live-ci.md` carries no Sonnet quarantine note, and the shared suite passes on that commit with recorded model evidence resolving to `claude-sonnet-5`.

**AC-3 - `gate prepare` names the selection it rejected.**
Verified by: a Go test in `internal/gates` that prepares with a committed artifact plus one uncommitted reference and asserts the returned error contains both the `--reference` flag and that reference's path, and a sibling case asserting the `--artifact` prefix when the artifact itself is uncommitted. The test fails if the wrap is removed or the flag/path is dropped. (Mechanism AC; it counts as the enabling half of AC-1.)

**AC-4 - Existing guardrails, oracles, and unrelated lanes are unchanged in behavior.**
Verified by: `go test ./...`, `go test ./... -race`, and `gofmt -l ./cmd ./internal` clean; `assertRecordedGateHoldLog`'s accept/reject set unchanged, proven by its existing table at `internal/ensigncycle/livescenario_adapter_live_test.go:26-49` still passing unmodified in its mutation list; and no skip added to or removed from any Opus, Codex, or Pi case.

## Test plan

The offline half is settled and cheap. The ideation spike (see Spike record) already reproduced the trigger and proved the mechanism-1 wrap flips the error from pathless to attributed with `go test ./...` green; implementation promotes that throwaway into the AC-3 regression test in `internal/gates/prepare_test.go` and runs the full and `-race` suites. No live run is needed for AC-3 or AC-4.

The live half is the only real cost and the only real risk. Run the approval-gated `claude-live` `sonnet` lane at the fix commit and read the step log first — the failure surfaces there in three lines. Only on failure download the run artifact and read `claude-shared-scenarios/gate-guardrail/claude-stream.jsonl`, which carries every FO command with its exit code; that file, not `*-detail.jsonl`, is where this defect was actually diagnosed. Order matters: land the fix and let the live run judge it **before** removing the quarantine, so AC-2's removal is justified by a pass rather than assumed. If the run still shows more than one successful prepare, the stream will show which probe survived, and the contract text in mechanism 2 is the next lever — do not respond by relaxing the oracle.

Estimated cost: low offline, one live lane. The residual risk is model nondeterminism, not design uncertainty.

## Stage Report: ideation

- DONE: Diagnose PR #585's pre-quarantine Sonnet gate-guardrail CI failure and confirm it against the entity's documented Problem statement, refining it if the evidence shows something more specific.
  Corrected it: the run made FOUR successful `gate prepare` calls, not two, and honored the no-authority boundary (`assertGateHeld` passed). Evidence: run `30708727845` job `91392375253` step 14, plus artifact `8821429777` stream `claude-shared-scenarios/gate-guardrail/claude-stream.jsonl` (exit codes per call) and the FO's own final message admitting it hit "and reverted" the validation error mid-flight.
- DONE: Recommend the smallest concrete contract or runner change (name the exact file/mechanism) that makes Sonnet perform one gate prepare, then state commit, then present-and-stop; update Proposed approach with it.
  Primary mechanism: wrap the `gitsource.Inspect` failure at `internal/gates/prepare.go:154-157` so the error names the rejected `--artifact`/`--reference` path. Secondary: two sentence rewrites in `skills/fo-gate-lifecycle/SKILL.md` (lines 27 and 65), given verbatim before/after. Optional third and one explicit rejection recorded with reasons.
- DONE: Confirm or refine Acceptance criteria and Test plan against the concrete recommended fix, and record spike/no-spike-needed per the Proof policy for any unverified mechanism the recommendation relies on.
  AC-1 now measures the successful-prepare count against the independent baseline 4; AC-3 added for the attribution mechanism; AC-4 pins the oracle's accept/reject set to its existing mutation table. Spike record documents the exercised before/after and names the one claim no offline spike can settle.

### Summary

The archived stream contradicts the paraphrase this task inherited from issue #596. The FO never crossed the authority boundary — it presented once and stopped at `awaiting-captain` with no decision, consume, dispatch, or archive. What it did was bisect for an unnamed rejected source by binding throwaway `--question "test"` content to the real gate room, then `git checkout` the entity and `rm -rf` the binary-owned room twice to escape the correct frozen-binding refusal, yielding four successful prepares before one `state commit`. The oracle's single generic message for seven conditions is what produced the wrong diagnosis.

I spiked the fix rather than asserting it: a throwaway `internal/gates` test reproduced the pathless error, the one-line `fmt.Errorf` wrap flipped it to `--reference /…/command.log: …`, and `go test ./...` stayed green — so the blast radius on existing assertions is zero. The spike was reverted; `internal/` is clean.

Two things need a captain call. The entity title asserts a no-authority breach the evidence does not support, and retitling is outside an ensign's write scope. And the optional oracle-message mechanism is the only proposal not required by AC-1 — it is cheap and it is the reason this diagnosis cost a live run, but it is a fair cut.

## Stage Report: implementation

- DONE: Land mechanism 1: attribute the rejected --artifact/--reference path in internal/gates/prepare.go's selected-source rejection, and promote the ideation spike into the AC-3 regression test in internal/gates/prepare_test.go.
  `internal/gates/prepare.go` wraps the `gitsource.Inspect` failure with the offending `--artifact`/`--reference` flag and path. `TestPrepareAttributesRejectedSelectedSourceToItsFlag` asserts the returned error contains that flag+path for an uncommitted reference and, in a sibling case, an uncommitted artifact — fails if the wrap is reverted or the flag/path is dropped. Commit f7fce40ab.
- SKIPPED: Land mechanism 2: the two verbatim before/after sentence corrections in skills/fo-gate-lifecycle/SKILL.md (lines 27 and 65) from the entity's Proposed approach.
  Captain-directed scope narrowing, not an ensign judgment call. Baseline SKILL.md is 6592B against the 6600B cap in `internal/contractlint/fo_function_reference_invariant_test.go:14` (8B headroom); landing the verbatim text grew it to 7201B (+609B). I escalated to team-lead, tried a tighter rewrite that still preserved every rung team-lead named (correct-selection-reprepare-once, never-diagnose-with-throwaway, frozen-binding-terminal, per-gate-entry one-prepare bound) and it still needed +73B even after cutting an unrelated pre-existing clause. Team-lead/captain deferred mechanism 2 to a separate follow-up FO task rather than trim further or touch the cap. SKILL.md is reverted to main byte-for-byte (`git checkout --`, verified against HEAD).
- DONE: Land mechanism 3 (captain confirmed: keep, not cut) in internal/ensigncycle/livescenario_adapter_live_test.go so a failure names which of the seven conditions fired, without changing the oracle's accept/reject set.
  `assertRecordedGateHoldLog` now returns one of seven distinct messages (no successful prepare / commit before prepare / head before commit / duplicate prepare / decision after / consume after / dispatch after) in the same evaluated order as the original combined boolean. `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle`'s six existing mutation cases each now assert their specific message substring — fails if a condition's text regresses to the old generic string or a condition's boolean changes. Commit f7fce40ab.

### Summary

Mechanisms 1 and 3 are landed and committed (f7fce40ab) on `spacedock-ensign/sonnet-gate-guardrail-no-authority`; `go test ./...`, `go test ./... -race`, and `gofmt -l ./cmd ./internal` are all clean, including `TestFOInstructionComponentCaps`. Mechanism 2 is deferred by explicit captain decision, not cut or forgotten: its mandated verbatim text cannot fit the existing SKILL.md byte cap without cutting unrelated pre-existing contract text or touching the cap itself, both outside this task's scope; team-lead is tracking the follow-up. None of AC-1/AC-2/AC-3/AC-4 as written require mechanism 2 in their Verified-by text — AC-3 verifies mechanism 1's attribution, not mechanism 2's prose — but mechanism 2 was proposed as the hardening step that lowers risk to AC-1's live Sonnet run, so that residual risk (the FO's contract not yet naming the once-only-prepare/frozen-binding-terminal rules explicitly) carries into validation's live run. The Sonnet TODO quarantine (`claudeSonnetGateGuardrailTODO`, `TestClaudeSonnetGateGuardrailTODOModelScope`, its skip call, and the `docs/runtime-live-ci.md` note) is untouched, per the entity's own Test plan sequencing.

## Stage Report: validation

- DONE: Verify AC-3 and AC-4 with go test ./..., go test ./... -race, and gofmt -l ./cmd ./internal.
  All three clean at f7fce40ab in the implementation worktree. `TestPrepareAttributesRejectedSelectedSourceToItsFlag` (`internal/gates/prepare_test.go`) independently re-run and passing: it fails if the `--artifact`/`--reference` wrap in `prepare.go:154-160` is reverted. `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle` independently re-run under `-tags live` (its file is live-gated, so `go test ./...` alone doesn't compile it) and passing all six per-condition substring assertions — fails if any of the seven messages regresses to the old generic string. Diff-read of `f7fce40ab` confirms the switch in `assertRecordedGateHoldLog` evaluates the identical boolean conditions in the identical order as the prior `||` chain (accept/reject set unchanged), and the 3-file commit stat touches no Opus/Codex/Pi file, so no lane skip changed.
- DONE: Confirm mechanism 2 is genuinely deferred, not silently dropped, and that no AC's Verified-by text depends on it.
  `git diff main -- skills/fo-gate-lifecycle/SKILL.md` is empty and the file is 6592 bytes, matching the implementation report's claim byte-for-byte. Re-read all four ACs' Verified-by clauses in this entity: none names `SKILL.md` or its two sentences.
- DONE: Reconcile the still-active Sonnet TODO quarantine with AC-1's live-run requirement, per the Spot-check principle, and verify AC-1.
  Plain `go test -tags live -run TestLiveClaudeSharedScenarios/gate-guardrail` at this commit would self-skip for Sonnet (`claudeSonnetGateGuardrailTODO` fires before any real launch) and prove nothing. Cheap spot-check first, no model spend: `go test -tags live -run 'TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions|TestPiSharedScenarioCoverage' ./internal/ensigncycle` passed, confirming the scenario/runner plumbing is intact before spending on a live call. Then, per the Review-finding-disposition worker authority ("use a throwaway checkout for adversarial investigation... candidate bytes and Git HEAD stay unchanged"), I created a throwaway worktree at the exact candidate commit `f7fce40ab` (`/private/tmp/spacedock-sonnet-gate-guardrail-validation-audit/repo`, never the implementation worktree), locally commented out only the `claudeSonnetGateGuardrailTODO` skip line in `runClaudeGateGuardrailScenario` (never committed, checkout discarded after use), and ran the real live Sonnet gate-guardrail scenario twice against my local OAuth benchmark-token credential: `go test -tags live -count=1 -timeout 40m -run 'TestLiveClaudeSharedScenarios/gate-guardrail' ./internal/ensigncycle -v`. Both runs PASSED (212.74s, then 148.86s with `SPACEDOCK_LIVE_ARTIFACT_DIR` set to retain artifacts). I did not trust the internal assertion alone: I independently parsed the retained `claude-stream.jsonl` from the second run and found exactly one non-`--help` `gate prepare recorded-gate-task` Bash tool call (`is_error: false`, output `state=open`), issued only after `gate eligibility` and `gate prepare --help` (no throwaway `--question "test"` probing, no bisection), immediately followed by `state commit recorded-gate-task`, and no `--decision`, `gate consume`, or `dispatch build` anywhere after — matching `assertRecordedGateHoldLog` and `assertGateHeld` exactly. The final message confirms the FO stopped at the decision boundary having taken no action. This directly falsifies the pre-fix baseline (4 successful prepares, bisection, two destructive room deletions) at the same scenario and model. AC-1: PASSED.
- DONE: Since AC-1 passed, verify AC-2 with the same run's evidence.
  AC-2 requires `claudeSonnetGateGuardrailTODO`, `TestClaudeSonnetGateGuardrailTODOModelScope`, the skip call, and the `docs/runtime-live-ci.md` quarantine note to be **absent**. Direct inspection of the implementation worktree (untouched by the throwaway investigation above) confirms all four are still present verbatim — expected, since removing the quarantine is a real product/test-file edit the entity's own Test plan orders strictly after a passing live run, and this validation stage does not produce the deliverable (per this README's validation-stage discipline). AC-2: NOT YET SATISFIED — this is not a defect, it is the deliberately deferred second half of a two-phase sequence, now unblocked by the AC-1 evidence above.
- DONE: Recommend whether the AC-2 removal belongs to this validation stage or routes back through implementation.
  Recommend routing to implementation, not landing it here. Reasons: (1) this README states "The validator checks what was produced - it does not produce the deliverable itself"; (2) the edit is a real, already-fully-specified product/test change (delete `claudeSonnetGateGuardrailTODO`, `TestClaudeSonnetGateGuardrailTODOModelScope`, the skip call, and `docs/runtime-live-ci.md:138` — exactly the entity's own "Expected surface" line items), not a validator judgment call; (3) the implementation ensign is already being kept alive and addressable for exactly this kind of narrow follow-through per this entity's Housekeeping note. The AC-1 evidence above (2 passing live runs, one independently stream-verified) is what unblocks and justifies that follow-through — it should not be re-derived, only cited.

### Summary

AC-1, AC-3, and AC-4 all PASS with reproduced evidence: offline suites clean, the AC-3/AC-4 regression tests independently re-run (including the live-tag-gated mechanism-3 oracle test, which plain `go test ./...` does not compile), and a fresh live Sonnet `gate-guardrail` run — reproduced twice via a throwaway-checkout bypass of the still-active TODO quarantine, never touching the candidate — shows exactly one successful `gate prepare` followed by one `state commit` and a stop at the decision boundary, independently confirmed by parsing the retained transcript rather than trusting the assertion alone. AC-2 is not yet satisfied: the quarantine (`claudeSonnetGateGuardrailTODO`/its test/the skip call/the doc note) is still present in the implementation worktree, unchanged by this stage's investigation, because removing it is a real deliverable edit the entity's own Test plan orders after a passing live run — which this stage has now furnished. This is not a defect; it is the deliberately sequenced second half of AC-2, now unblocked.

**Recommendation: REJECTED (not a defect) — route to implementation for a narrow, evidence-justified follow-through**: delete `claudeSonnetGateGuardrailTODO`, `TestClaudeSonnetGateGuardrailTODOModelScope`, and the skip call in `runClaudeGateGuardrailScenario` (`internal/ensigncycle/claude_live_runner_test.go`), and delete the quarantine paragraph at `docs/runtime-live-ci.md:138`, citing this stage's AC-1 evidence rather than re-running the live lane. Re-validate afterward with `go test ./...`/`-race`/`gofmt` plus a fresh live Sonnet `gate-guardrail` run (now unskipped) for AC-2's "shared suite passes... resolving to claude-sonnet-5" clause.

### Feedback Cycles

- Cycle 1: REJECTED (not a defect) — validation ensign (`spacedock-ensign-sonnet-gate-guardrail-no-authority-validation`); surface 3 files/82 lines vs estimate 44 lines (86%, over the declared 40% tolerance — driven entirely by `internal/ensigncycle/livescenario_adapter_live_test.go`, 46 lines vs 13 estimated, from turning one generic oracle message into seven exact ones); AC unchanged (no AC narrowed). Captain reviewed the deviation and did not request a design reset. FO disposition: fix (owned, Material — the remedy is the entity's own already-approved Expected-surface items). Routed to implementation: delete `claudeSonnetGateGuardrailTODO`, `TestClaudeSonnetGateGuardrailTODOModelScope`, the skip call, and the `docs/runtime-live-ci.md:138` note, citing this cycle's AC-1 evidence rather than re-running the live lane.

## Stage Report: implementation (cycle 2)

- DONE: Delete the four quarantine items named in the feedback context, citing the validation stage's AC-1 evidence. Do not run the live test again.
  Deleted `claudeSonnetGateGuardrailTODO`, `TestClaudeSonnetGateGuardrailTODOModelScope`, and the skip call in `runClaudeGateGuardrailScenario` (`internal/ensigncycle/claude_live_runner_test.go`), and the quarantine paragraph at `docs/runtime-live-ci.md:138`. Cited the validation stage's two-run live evidence (212.74s and 148.86s, both one `gate prepare` + one `state commit` + stop) rather than re-running the live lane. The Opus skip in the same function is untouched. No live test invoked this stage.
- DONE: Run go test ./..., go test ./... -race, and gofmt -l ./cmd ./internal. All three must be clean.
  All three clean at commit e36769727. `go vet -tags live ./internal/ensigncycle/...` also clean, confirming the live-gated file still compiles with the skip removed.
- DONE: Commit and push the branch to origin. Report back the commit SHA and diff stat. Do not open a pull request yourself.
  Commit `e36769727` on `spacedock-ensign/sonnet-gate-guardrail-no-authority`, pushed to origin. Diff stat vs `f7fce40ab`: `docs/runtime-live-ci.md | 1 -` and `internal/ensigncycle/claude_live_runner_test.go | 25 -------------------------` (2 files, 26 deletions) — exact match to the entity's declared Expected surface for these two items. No pull request opened.

### Summary

Removed the Sonnet gate-guardrail TODO quarantine per the FO-authorized fix disposition, using the validation stage's already-reproduced live evidence rather than spending another live run. `go test ./...`, `-race`, and `gofmt -l` are clean; the diff matches the entity's Expected surface exactly (-25/-1). Mechanism 2 (SKILL.md) remains deferred per the earlier captain decision and is untouched here.

## Stage Report: implementation (cycle 3)

- DONE: Reconcile exact candidate e36769727 with current post-YS main without discarding either side.
  Merge commit `4f4ee424792f9ea1044578334040825bcd1da585` joins candidate `e36769727` with `origin/main` `d3e70e958`; the branch is pushed and clean.
- DONE: Resolve the post-YS conflict set on the canonical journey architecture.
  Resolved `docs/runtime-live-ci.md` and `internal/ensigncycle/claude_live_runner_test.go`; accepted YS's deletion of `internal/ensigncycle/livescenario_adapter_live_test.go` without restoring it.
- DONE: Move the 3Z quarantine removal and strict diagnostics onto the canonical live inventory and common runner.
  `TestLiveCommonGateGuardrail` has no Claude Sonnet TODO; its strict seven-condition log diagnostic now lives in `internal/ensigncycle/claude_runtime_helpers_test.go`, used by the common exercise.
- DONE: Preserve the intended TODO-owner boundaries after YS.
  Claude Sonnet is live evidence; Codex and Pi alone retain `TODO(3zzpdw704df1g8pg1x9thzmw)` in `shared_live_runner_test.go`, exactly as enforced by registry reconciliation; Opus has no 3Z TODO.
- DONE: Preserve gate-prepare attribution and the strict accept/reject behavior.
  `TestPrepareAttributesRejectedSelectedSourceToItsFlag` and `TestAssertRecordedGateHoldLogAcceptsPrepareFirstLifecycle` pass; each fails if attribution is dropped or any of the six negative mutations is accepted or misdiagnosed.
- DONE: Prove canonical live inventory reconciliation.
  `go test ./internal/contractlint -run '^TestRuntimeLiveRegistryReconciliation$' -count=1` passes and fails on TODO ownership, builder/assertion drift, orphan fixtures, or lane-selector drift.
- FAILED: Prove full and race suites clean in the current shared-state environment.
  `go test ./...` and `go test ./... -race` run to completion; all packages pass except `internal/gates`, where unchanged post-YS `TestV1PilotManifestReadsAndValidates` expects two root entities absent from the shared state checkout (`codex-launch-multi-agent-v2.md`, `gate-agent-ergonomics.md`; both exist only under `_archive`). Candidate diff against main does not touch that manifest or test.
- DONE: Prove applicable live Sonnet behavior on the canonical common journey.
  `TestLiveCommonGateGuardrail` passed twice after reconciliation (115.66s and retained-artifact 133.47s), resolved to `claude-sonnet-5`; retained stream independently shows one prepare, one state commit, and zero decision/consume/dispatch calls.
- DONE: Run repository formatting requirements.
  `gofmt -w ./cmd ./internal` completed; focused mechanism, oracle, inventory, and live checks remained green afterward.

### Summary

The exact candidate is reconciled and pushed at `4f4ee4247` without restoring any retired pre-YS harness file or weakening the gate oracle. YS's registry/common-runner architecture remains canonical: Sonnet is active evidence, the stale quarantine paragraph is gone, and the strict diagnostic moved to the shared helper. The only red evidence is a main-equivalent pilot-manifest/shared-state mismatch outside this candidate; both full commands otherwise completed cleanly, while the targeted and real Sonnet common journey passed.

## Stage Report: implementation (cycle 3 report repair)

- DONE: Reconcile exact candidate e36769727 with current post-YS main without discarding either side.
  Merge commit `4f4ee424792f9ea1044578334040825bcd1da585` joins candidate `e36769727` with `origin/main` `d3e70e958`; the branch is pushed and clean. Conflicts in `docs/runtime-live-ci.md` and `internal/ensigncycle/claude_live_runner_test.go` were resolved, and YS's deletion of `internal/ensigncycle/livescenario_adapter_live_test.go` was preserved.
- DONE: Move the 3Z quarantine removal and strict diagnostics onto the canonical live inventory and common runner.
  `TestLiveCommonGateGuardrail` has no Claude Sonnet TODO, while Codex and Pi alone retain `TODO(3zzpdw704df1g8pg1x9thzmw)` and Opus has none. The stale quarantine paragraph is removed, and the unchanged seven-condition accept/reject diagnostic now lives in `internal/ensigncycle/claude_runtime_helpers_test.go` on the common exercise path.
- FAILED: Preserve gate-prepare behavior and prove focused, full, race, inventory, and applicable live checks.
  Focused attribution/oracle tests and registry reconciliation passed; `TestLiveCommonGateGuardrail` passed twice (115.66s and 133.47s), resolved to `claude-sonnet-5`, with one prepare, one state commit, and zero forbidden decision/consume/dispatch calls. `go test ./...` and `go test ./... -race` each failed only because unchanged post-YS `TestV1PilotManifestReadsAndValidates` expects `codex-launch-multi-agent-v2.md` and `gate-agent-ergonomics.md` at state-root paths although both are archived; the candidate does not modify that manifest or test.

### Summary

This report-only repair maps the recorded cycle-3 evidence onto the three dispatched checklist items exactly. Candidate commit `4f4ee4247` remains unchanged; focused, inventory, and live Sonnet proof is green, while the full and race items remain honestly FAILED on the recorded main-equivalent pilot-manifest/shared-state mismatch.

## Stage Report: implementation (cycle 4)

- DONE: Rebase the clean 3Z candidate onto exact validated 0Y commit 8728da3a0 without discarding either side.
  The candidate was rebased through corrected 0Y and finally onto actual landed tip `9021cbf374bb1740d1b2e45155041e1f809372c4`; both `8728da3a0` and `9021cbf3` are ancestors. Final bridge `f638451b974c1f2ec503fbdfe68cff08641b5efe` has parents `394f86b0a3a688f15cd9f85dc377f734d72c8973` and old remote `4f4ee424792f9ea1044578334040825bcd1da585`; its tree `06820105b123fc9a46b101a47fd9afee006d52b4` is byte-identical to parent 1, enabling a normal non-force push. The landed deletion of `internal/gates/testdata/v1_pilot_manifest.txt` remains intact.
- DONE: Confirm the YS common inventory and 3Z Sonnet TODO removal remain exact after 0Y.
  `TestLiveCommonGateGuardrail` has no Claude Sonnet 3Z TODO; Codex and Pi alone retain `TODO(3zzpdw704df1g8pg1x9thzmw)`, Opus has none, the stale Sonnet quarantine paragraph is absent, and the strict seven-condition diagnostic remains on the common helper path. Registry reconciliation passes.
- FAILED: Run focused, inventory, full, race, and applicable Sonnet evidence at the new exact head.
  At bridge head `f638451b9`, focused attribution/oracle tests, registry inventory, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` all pass. The exact-head Sonnet attempt resolved `claude-sonnet-5` but failed before any FO work because local OAuth had expired and could not refresh, with `ANTHROPIC_API_KEY` unset; artifact `/tmp/spacedock-3z-final-live.giuur5`. This is an infrastructure failure, not green Sonnet evidence; fresh CI must supply that proof.

### Summary

3Z is reconciled onto the actual merged 0Y tip and normally pushed at exact head `f638451b974c1f2ec503fbdfe68cff08641b5efe`, with the authorized ancestry bridge preserving a byte-identical rebased tree. The intended four-file product diff, YS common inventory, Sonnet TODO removal, strict oracle, and ambient pilot-manifest deletion are preserved. All offline proof is green; applicable local Sonnet evidence remains explicitly FAILED on expired OAuth before workflow execution, so CI owns the fresh live verdict.

## Stage Report: implementation (cycle 4 report repair)

- DONE: Rebase the clean 3Z candidate onto exact validated 0Y commit 8728da3a0 without discarding either side.
  The candidate was rebased through corrected 0Y and onto landed tip `9021cbf374bb1740d1b2e45155041e1f809372c4`; both `8728da3a0` and `9021cbf3` are ancestors. Final bridge `f638451b974c1f2ec503fbdfe68cff08641b5efe` has parents `394f86b0a3a688f15cd9f85dc377f734d72c8973` and `4f4ee424792f9ea1044578334040825bcd1da585`, with tree `06820105b123fc9a46b101a47fd9afee006d52b4` byte-identical to parent 1; the ambient pilot-manifest deletion remains intact.
- DONE: Confirm the YS common inventory and 3Z Sonnet TODO removal remain exact after 0Y.
  `TestLiveCommonGateGuardrail` has no Claude Sonnet 3Z TODO; Codex and Pi alone retain `TODO(3zzpdw704df1g8pg1x9thzmw)`, Opus has none, the stale quarantine paragraph is absent, and the strict seven-condition diagnostic remains on the common helper path. Registry reconciliation passed.
- SKIPPED: Run focused, inventory, full, race, and applicable Sonnet evidence at the new exact head.
  Focused attribution/oracle tests, inventory reconciliation, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` passed at exact bridge head `f638451b9`. The local Sonnet attempt never reached FO or product work because OAuth refresh failed and `ANTHROPIC_API_KEY` was unset; artifact `/tmp/spacedock-3z-final-live.giuur5`. CI owns the fresh live proof; this report does not call Sonnet green.

### Summary

This report-only repair preserves the completed landed-base reconciliation and all green offline evidence. Applicable local Sonnet proof is SKIPPED because authentication failed before workflow execution; fresh CI must establish the live verdict.

## Stage Report: validation (cycle 2)

- DONE: Verify exact head f638451b9 preserves the post-YS 3Z behavior, strict diagnostics, and all acceptance criteria.
  `HEAD` and the pushed branch both resolve to `f638451b974c1f2ec503fbdfe68cff08641b5efe`; focused AC-3 attribution, AC-4 oracle, inventory, full, race, and formatting proofs pass. Prior common-journey Sonnet evidence at the post-YS candidate remains the AC-1 behavioral evidence; the exact-head retry below produced no product observation.
- DONE: Verify the ancestry bridge is byte-identical and the landed 0Y inventory decoupling remains intact.
  Bridge parents are `394f86b0a3a688f15cd9f85dc377f734d72c8973` and `4f4ee424792f9ea1044578334040825bcd1da585`; bridge tree `06820105b123fc9a46b101a47fd9afee006d52b4` equals parent 1 and `git diff ^1..HEAD` is empty. `origin/main` is exact `9021cbf374bb1740d1b2e45155041e1f809372c4`, an ancestor, and `8728da3a0` is its ancestor.
- DONE: Verify the intended four-file diff.
  Against landed `origin/main`, only `docs/runtime-live-ci.md`, `internal/ensigncycle/claude_runtime_helpers_test.go`, `internal/gates/prepare.go`, and `internal/gates/prepare_test.go` differ (69 insertions, 14 deletions); the observable command change is only rejected-source flag/path attribution.
- DONE: Verify Sonnet TODO removal and Codex/Pi TODO ownership.
  `TestRuntimeLiveRegistryReconciliation` passes and would fail on target/owner or inventory drift; repository search finds the 3Z TODO only for Codex and Pi in `shared_live_runner_test.go` plus that independent expected inventory, with no Sonnet/Opus TODO or quarantine note.
- DONE: Verify the strict seven-condition diagnostic.
  The committed focused oracle test passes its valid lifecycle and six negative mutations. A throwaway detached exact-head test additionally exercised all seven distinct branches, including missing/early `state-head`; every mutation was rejected with its condition-specific diagnostic, and the throwaway checkout was removed without candidate edits.
- DONE: Verify 0Y pilot-manifest decoupling.
  `internal/gates/testdata/v1_pilot_manifest.txt` and its old `internal/gates/application_test.go` coupling are absent, and both paths have zero diff from landed `origin/main`; full and race runs therefore exercise the decoupled inventory.
- DONE: Run focused, inventory, full, race, formatting, and fresh Sonnet evidence or report precise infrastructure unavailability.
  Focused attribution/oracle tests, registry reconciliation, `go test ./...`, `go test ./... -race`, and `gofmt -l ./cmd ./internal` are clean. Fresh `TestLiveCommonGateGuardrail` resolved `claude-sonnet-5` but exited before FO work because the Keychain OAuth session was expired and could not refresh while `ANTHROPIC_API_KEY` was unset; retained artifact: `/tmp/spacedock-3z-validation-live.wdtymG/claude-shared-scenarios/gate-guardrail`. This is not green Sonnet evidence.

### Summary

Recommendation: PASSED, with no material, deferred-risk, or polish findings in the candidate. Exact topology, byte identity, four-file scope, TODO ownership, seven-condition behavior, 0Y decoupling, focused/full/race tests, and formatting all validate cleanly; the only unavailable observation is a fresh exact-head Sonnet run, whose launch failed at authentication before the product or first officer executed and is explicitly not counted as evidence.

## Stage Report: implementation (cycle 5)

- DONE: Reconcile exact 3Z candidate f638451b9 onto the landed NV tip without discarding either side or force-pushing.
  Normal merge `f6221b314f9f90d5cfce79dabdd9231d776d05b6` joined `f638451b9` to exact NV tip `48b45a815`; moving-target merge `48fd54b2aa8fd63b002beb7ec46afa6ecaea6228` then joined landed 824 tip `5c54154a0`. Both pushes were ordinary fast-forwards to PR #640, and both landed tips are ancestors of the final head.
- DONE: Preserve 3Z rejected-source attribution, seven-condition diagnostics, Sonnet ownership, and NV Codex terminal-boundary behavior.
  Final diff from `origin/main` is exactly the four intended 3Z files (69 insertions, 14 deletions). Focused attribution/oracle tests fail if the flag/path wrap or any condition-specific rejection is lost; registry reconciliation preserves Sonnet as active gate evidence, the existing Codex/Pi gate TODO ownership, and NV's removal of the full-cycle Codex TODO. NV terminal tests require `turn.completed` plus a non-empty final message.
- DONE: Run focused, registry, full, race, formatting, Sonnet, and Codex evidence at the new exact head.
  At exact `48fd54b2a`, focused gates/ensigncycle tests, `TestRuntimeLiveRegistryReconciliation`, `go test ./...`, `go test ./... -race`, and `gofmt -w ./cmd ./internal` are green. Local `TestLiveCommonAutoContinueAfterImplementation` passed both Codex variants in 430.89s; both artifacts record `terminal: true` and `timed_out: false`. CI run `31299255670` attempt 2 is green for offline job `93209745179`, Sonnet `claude-sonnet-5` job `93209844691`, and Codex job `93209844700`.
- SKIPPED: Repair or rerun unrelated CI infrastructure findings.
  Attempt 1's unchanged `TestDurableQuestionedRejectsTerminalHistory` `bad object HEAD` failure passed 20/20 locally and went green on the one FO-authorized rerun without candidate edits. Attempt 2's QZ-owned journey-delta-comment job `93214710577` failed only because `actions/download-artifact@v5` could not download/extract artifact `9034189537` after five retries; the FO authorized DECLINE candidate edit and no rerun.

### Summary

PR #640 is updated at exact head `48fd54b2aa8fd63b002beb7ec46afa6ecaea6228`, reconciled normally onto both NV and 824 while retaining the intended four-file 3Z value. Relevant offline, Sonnet, and Codex evidence is green at that head; no retired runner, Codex full-cycle TODO, QZ journey-comment change, or force-push was introduced.

## Stage Report: validation (cycle 3)

- DONE: Verify exact head 48fd54b2a contains only the intended four-file 3Z delta on landed NV and 824.
  `HEAD` and pushed PR #640 both resolve to `48fd54b2aa8fd63b002beb7ec46afa6ecaea6228`; `origin/main` is exact 824 tip `5c54154a061a26fbae503ad9349f5ec80c390e66`, NV tip `48b45a815` is an ancestor, and the diff is exactly four files with 69 insertions/14 deletions.
- DONE: Verify rejected-source attribution, seven-condition diagnostics, live inventory ownership, and NV terminal handling satisfy every acceptance criterion.
  AC-1 through AC-4 reproduce cleanly: exact-head Sonnet retained evidence, focused attribution/oracle execution, independent registry reconciliation, and exact Codex terminal artifacts all agree with the intended behavior and ownership.
- DONE: AC-1 (VALUE) - A live Sonnet gate-guardrail run reaches the decision boundary with exactly one successful `gate prepare`.
  CI run `31299255670` attempt 2 retained stream and measured record show model `claude-sonnet-5`, outcome passed, one prepare, one later state commit plus `state-head`, and zero decision/consume/dispatch calls.
- DONE: AC-2 - The Sonnet gate-guardrail case is live evidence again rather than quarantined.
  Registry reconciliation passes; the 3Z TODO exists only for Codex and Pi, no Sonnet/Opus quarantine remains, and exact-head Sonnet job `93209844691` passed all 16 common tests with the gate-guardrail measured record green.
- DONE: AC-3 - `gate prepare` names the selection it rejected.
  `TestPrepareAttributesRejectedSelectedSourceToItsFlag` passes both reference and artifact cases and fails if either the rejected flag or exact selected path disappears.
- DONE: AC-4 - Existing guardrails, oracles, and unrelated lanes are unchanged in behavior.
  The focused committed oracle rejects its six mutations; a removed throwaway test exercised all seven production branches, including missing/early `state-head`, and every branch emitted its condition-specific diagnostic.
- DONE: Verify NV terminal-boundary behavior and 824 coexistence independently.
  `TestCodexProcessRecognizesTerminalTurnBeforeOSExit` and `TestCodexProcessRequiresFinalMessageForTerminalTurn` pass; downloaded exact-run single/split-root artifacts each have exit 0, `terminal: true`, `timed_out: false`, one `turn.completed`, and non-empty final messages.
- DONE: Verify focused, registry, full, race, formatting, Sonnet, and Codex evidence; keep QZ journey-comment failure out of candidate findings.
  Focused gates/ensigncycle tests, `TestRuntimeLiveRegistryReconciliation`, `go test ./...`, `go test ./... -race`, `gofmt -w ./cmd ./internal`, offline job `93209745179`, Sonnet job `93209844691`, and Codex job `93209844700` are green at the exact head; candidate status remains clean.
- SKIPPED: Repair, relabel, or rerun the QZ-owned journey-delta-comment artifact-download failure.
  Per Captain direction, job `93214710577` remains excluded: its only failure is `actions/download-artifact@v5` failing after five retries to download/extract artifact `9034189537`; it is not a 3Z candidate finding.

### Summary

Recommendation: PASSED, with no material, deferred-risk, or polish findings in the candidate. Exact scope, NV/824 coexistence, all four acceptance criteria, seven-condition diagnostics, TODO ownership, full/race/formatting, and retained Sonnet/Codex live evidence validate cleanly while the candidate remains byte-unchanged.
