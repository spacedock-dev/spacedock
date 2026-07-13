---
title: Restore lazy loading for first-officer merge and write cores
status: ideation
source: Captain correction after fresh-session boot trace, 2026-07-13
started: 2026-07-13T15:56:42Z
completed:
verdict:
score:
worktree:
issue:
id: 1kevganrmr2csr539ktfjerh
---

Restore the intended first-officer loading boundary: boot reads the shared core and active runtime adapter, while the write core loads at the first FO-authored mutation and the merge core loads only at terminal or merge-mod recovery handling.

## Problem

The archived `dp` task (`dpwp415wfzj6yrcwbs0krrea`, `fo-deferred-load-point-hunt-vs-skill-addressing`) diagnosed a real Claude failure class in which the model hunted the filesystem for delayed references; its structural addressing direction fed PR #491. PR #495 (`f22360de`) was not `dp`: it explicitly superseded #491's root-addressing direction and bundled `6h`, `p4`, and the `m1y` small-core preload. That preload overcorrected by eagerly importing `fo-merge-core.md` and `fo-write-core.md` from the first-officer entry skill. Commit `6baeed70` moved merge eager and `1e4423e1` moved write eager; later contractlint tests canonized exactly three eager imports. This regressed the earlier split explicitly designed in `shared-merge-dispatch-contract`: an interactive greet that stops should not pay for mutation or terminal ceremony it never uses.

The failure was deterministic discovery, not laziness itself. Restoring vague skill discovery would recreate the original hunt. The delayed cues must name exact canonical reference paths and load them at behaviorally enforced triggers.

## Proposed approach

- Keep `@references/first-officer-shared-core.md` as the entry skill's only eager canonical reference. Keep the selected host runtime adapter's existing boot-time read and the resident smallest-sufficient rule unchanged.
- Remove `@references/fo-write-core.md` and `@references/fo-merge-core.md` from `skills/first-officer/SKILL.md`; do not delete or change either canonical body.
- Add one address rule to the shared core: the skill loader supplies `Base directory for this skill: <dir>` when `spacedock:first-officer` loads. A delayed reference read resolves by joining that already-supplied base with one literal suffix. Read exactly `<base>/references/fo-write-core.md` or `<base>/references/fo-merge-core.md`; never resolve against cwd, invoke another skill, try an alternate path, or search. If that exact read fails, halt and report it instead of hunting.
- Add the write cue at the first FO-authored write intent or state-changing command, including `spacedock new`, `status --set`, `state commit`, dispatch-state mutation, Edit/Write/apply-patch, redirection, archive, and terminal mutation. Read write-core first, then run `«write.classify»(target, intent)`, then mutate. Read-only boot/status calls and a gate presentation do not trigger it.
- Add the merge cue at the first terminal boundary or the first handling of `mod-block=merge:*`. Read merge-core before `«merge.guard»`, `«hooks.run»("merge")`, direct terminal-status mutation, archive/cleanup, or reasoning from a merge mod-block. Dispatch and non-terminal gates do not trigger it.
- Define cross-trigger order. If a terminal operation is also the session's first mutation, load write-core and classify first, then load merge-core, then invoke the merge/terminal operation. If write-core is already resident, do not reload it; merge-core still loads immediately before its first owned boundary.
- Replace #495's eager-topology assumptions with one structural topology/reachability test and host-neutral read-order oracles. The structural test proves exact path resolution; the stream oracles prove the model used the path at the correct time.

### Exact contract wording

The implementation should express the delayed cues in this behaviorally equivalent form (wordsmithing may shorten it without weakening the ordered verbs):

> Resolve delayed FO references only from the `Base directory for this skill` supplied when `spacedock:first-officer` loaded. Append the literal `references/...` suffix; never use cwd, `Skill(...)`, discovery, or search. If the exact read fails, report the path and halt.
>
> - `<base>/references/fo-write-core.md` — before the first FO-authored write intent or state-changing command: READ, then `«write.classify»`, then mutate.
> - `<base>/references/fo-merge-core.md` — at the first terminal boundary or `mod-block=merge:*` recovery: READ before merge guard, merge hook, terminal mutation, archive, or cleanup.

The State Management sentence changes from “under the eagerly loaded `«write.classify»` write-authority scope” to “after the exact write-core read and `«write.classify»` at the first FO-authored mutation.” The Completion and Gates terminal branch and Mod Hook Convention each call the merge load rule before their first merge-owned action; they do not duplicate the core body.

### Observable read-order journeys

| Journey | Required ordered observation |
| --- | --- |
| Cold gate hold (`gate-guardrail`) | shared core → selected runtime adapter → gate presentation; no write-core or merge-core read |
| First mutation (`filing`) | shared/runtime → exact write-core read → `«write.classify»` → `spacedock new`; no merge-core read |
| Non-terminal dispatch (`rejection-flow`) | exact write-core read precedes first dispatch-state mutation; merge-core remains absent through the human gate |
| Terminal guard (`merge-hook-guardrail`) | exact write-core read/classify → exact merge-core read → attempted terminal `status --set`; existing refusal remains the outcome |
| Merged-PR terminalization (`shallow-boot`) | exact write-core read/classify → exact merge-core read → terminal/archive commands; durable merged+archived state remains correct |
| Resume `mod-block=merge:*` | exact write-core read/classify when recovery will mutate → exact merge-core read → merge guard/hook handling; no merge action precedes the core read |

The existing `shallow-boot` scenario is deliberately a positive terminal-load journey, not the absence oracle: it advances a merged PR before greeting and therefore legitimately needs both deferred cores. The mutation-free `gate-guardrail` journey supplies the cold no-load proof.

## Out of scope

- Deferring the active runtime adapter.
- Reintroducing the removed standalone write-core skill wrapper or the rejected callable-skill discovery path.
- Changing merge, mutation, gate, or dispatch behavior beyond reference load timing.
- Deferring `fo-dispatch-core`, status viewer, gate presentation, feedback routing, or recovery differently from their existing triggers.
- Changing the content of `fo-write-core.md` or `fo-merge-core.md`, or introducing a runtime-managed reference loader/manifest.
- Treating lower token use as permission to weaken merge guards, write scope, durable-state assertions, or supported-host behavior.

## Mechanism spike

The riskiest decision is whether a delayed reference can be addressed without the rejected callable-skill wrapper or cwd-relative guessing. Existing real host artifacts prove the required primitive:

- Claude PR #495 sonnet shallow-boot stream: the `Skill(spacedock:first-officer)` result supplied `Base directory for this skill: /tmp/spacedock-live-plugin-2455506351/skills/first-officer`; subsequent `cat` calls to that base plus `references/first-officer-shared-core.md`, `references/fo-merge-core.md`, and the runtime adapter all succeeded. The opus stream used `Read` with the same base-plus-literal-suffix shape successfully.
- Codex PR #496 cycle-2 shallow-boot stream: the first command read shared, merge, write, and the Codex adapter from one installed first-officer base under `.../plugins/cache/.../skills/first-officer/references/`; every exact file resolved. This is also the captured eager baseline the task removes.
- The archived `dp` evidence distinguishes this winning primitive from the failure: absolute reads reconstructed from the announced plugin base succeeded; cwd-relative `references/...` reads and subsequent search were the unstable path. The rejected `Skill(spacedock:fo-*-core)` promotion is unnecessary.

Result: no new runtime mechanism is needed. Implementation uses the already-proven loader-supplied base plus literal suffix and tests timing/order, which is the unproven part. It must begin RED by running the new order oracle against the current eager streams/topology: gate-hold must fail because both cores are read eagerly, while missing/reordered controls must fail for their named reasons.

## Acceptance criteria

**AC-1 (VALUE) - A fresh first-officer run that reaches its first human gate without write or terminal work loads the shared core and selected runtime adapter but reads neither merge-core nor write-core, reducing the cold resident contract by at least 8,000 bytes on every host.**
Verified by: fresh Claude and Codex `gate-guardrail` runs record normalized read events through the gate stop; shared and the selected adapter are present, both deferred cores are absent, the durable gate/final response remains correct, and the fixed byte registry is at least 8,000 bytes below the current eager baseline. Baseline on 2026-07-13: merge 2,830 B + write 5,843 B = 8,673 B eagerly read; the new cues may consume at most 673 B of that saving.

**AC-2 - The first FO-authored mutation deterministically reads the canonical write-core and classifies the target before mutation, without broad search, cwd-relative lookup, wrapper invocation, or retry.**
Verified by: Claude and Codex `filing` traces normalize to exact write-core read → classify evidence → successful `spacedock new`, with no merge-core read and no `find`, recursive grep/glob, `Skill(spacedock:fo-write-core)`, missing-path retry, or earlier write. Offline planted missing, post-mutation, wrong-base, wrapper, and broad-search traces each turn the same oracle red.

**AC-3 - Merge-core first loads at a terminal or merge-mod recovery boundary, after write-core when that boundary is also the first mutation, and nowhere earlier.**
Verified by: the non-terminal `rejection-flow` trace reaches its gate with no merge-core read; `merge-hook-guardrail` and `shallow-boot` traces read the exact merge core before terminal/guard commands and preserve refusal and merged+archived outcomes; a host-dialect fixture for existing `mod-block=merge:*` proves read-before-guard/hook order. Planted eager, missing, and post-merge reads fail.

**AC-4 - The contract has one eager canonical reference plus two deterministic, resolvable delayed cues, with no alternative entry surface.**
Verified by: replace `TestFirstOfficerEagerReferencesKeepDispatchCoreDeferred`, `TestFirstOfficerEagerReferenceTopology`, and stale closure comments/assumptions with a topology test requiring only `@references/first-officer-shared-core.md`, selected-runtime instructions, both base-anchored literal paths, and their ordered trigger verbs. Production `os.Stat` closure verifies each canonical body. Planted eager imports, missing/dangling cues, reversed READ/classify/mutate wording, and restored `skills/fo-{write,merge}-core` wrappers each fail.

**AC-5 - The loading-boundary repair preserves first-officer mutation, gate, dispatch, terminal, archive, refusal, and teardown behavior across supported hosts.**
Verified by: focused contractlint/order-oracle suites and the affected Claude/Codex live scenarios pass without weakening durable assertions; Pi's static topology/runtime-adapter coverage remains green until Pi has shared live runners. Then run `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, `go test ./... -race`, and `go test -tags live -run '^$' ./...`. A detached throwaway audit plants eager, missing, reordered, wrong-base, and wrapper regressions and observes the intended failures.

AC-4 is the mechanism serving AC-1 through AC-3; it does not pass if the live/order evidence misses the value boundary. AC-5 prevents the byte reduction from passing while behavior regresses.

## Test plan

1. **RED topology (small, offline).** Rewrite the eager-reference tests around a one-import registry and delayed-cue records. Run against current main and record failures for eager merge/write plus missing delayed cues. Add non-vacuous controls for eager import, missing file, wrong base/cwd wording, wrapper restoration, and READ/classify/mutate reversal.
2. **RED order oracle (medium, fixture).** Extend the Claude turn parser's existing `ReadTargets` path and add a Codex `item.* command_execution` normalizer. Emit ordered events for exact core reads, broad search, wrapper skill calls, classify markers, mutation commands, merge guard/hook, and gate/final boundaries. Feed committed host-dialect traces covering the six-row journey table. The current eager gate trace must fail; missing/reordered/broad-search fixtures must fail independently.
3. **Contract change (small).** Remove the two eager imports, update shared-core header/Deferred Load Points/State Management/Completion and Gates/Mod Hook Convention with one base-address rule and ordered triggers, and leave both canonical bodies byte-identical. Update closure registries/comments that call them eager; preserve their section-anchor tests.
4. **Focused GREEN (medium).** Run `go test ./internal/contractlint -count=1`, parser/order-oracle packages, `go test ./skills/integration/... -count=1`, and fixed-byte metrics. Require exactly one eager import, both delayed paths resolving, every planted control red, and at least 8,000 B cold-surface reduction.
5. **Live behavior (expensive, existing calls).** Add assertions to existing Claude/Codex `gate-guardrail`, `filing`, `rejection-flow`, `merge-hook-guardrail`, and `shallow-boot` runners rather than inventing a new model-only scenario. Grade normalized read order plus each scenario's existing durable state/final-message assertions. Run protected host jobs at the exact implementation SHA; local SKIP is not green evidence. Pi remains compile/static-covered because its shared live scenario runners are separately backlog-scoped.
6. **Documentation.** Update `docs/runtime-live-ci.md` after the durable-state assertion paragraph: “Reference-loading scenarios also inspect tool-call order because the read boundary is the behavior: write-core before the first FO mutation, merge-core before terminal/mod-block handling, and neither during a mutation-free gate hold. Durable state still grades the resulting mutation, refusal, and archive outcome.” No CLI/user-command docs change: command syntax and output are unchanged.
7. **Full gates and adversarial audit (medium).** Run the repository gates in AC-5. In a detached throwaway checkout, independently plant: eager merge import; eager write import; deleted cue; cue after mutation; cwd-relative/alternate path; wrapper directory; merge read after guard. Run the focused tests after each isolated edit and record the failing test/assertion, then remove the checkout. The validator also compares both core-body hashes to pre-change main to prove timing-only scope.

Estimated implementation complexity is medium: instruction text and structural tests are small; the main work is a reliable cross-host event normalizer and strengthening existing live assertions. No production Go command behavior changes.

## Stage Report: ideation

- DONE: Produce a behavior-first design that restores exact-path deferred write-core and merge-core loads at their mutation and terminal/mod-block triggers without reintroducing search or skill-discovery guessing.
  The design anchors both literal reference suffixes to the loader-supplied first-officer base, defines halt-on-exact-read-failure behavior, and orders write → classify → merge → terminal mutation when triggers coincide.
- DONE: Define measurable acceptance criteria and concrete topology/read-order/live scenarios proving a greet-stop reads neither deferred core while later mutation and merge paths deterministically read the right core first.
  Five ACs bind a mutation-free cold gate hold, filing, non-terminal rejection dispatch, terminal refusal, merged-PR terminalization, and merge-mod recovery to normalized host read events plus preserved durable outcomes; AC-1 records the 8,673 B eager baseline and requires at least 8,000 B reduction.
- DONE: Record the riskiest mechanism decision and a focused-through-full test plan, including planted eager/missing/reordered cues and a detached adversarial audit.
  Real Claude/Codex artifacts prove the loader-supplied base plus literal suffix resolves; the plan starts RED, adds cross-host order oracles and seven planted regressions, reuses protected live journeys, runs full/race/compile gates, and finishes in a throwaway checkout.
- DONE: Distinguish the cold no-load oracle from the existing mutation-bearing shallow-boot scenario.
  `gate-guardrail` proves no deferred reads; `shallow-boot` remains a positive write+merge load case because it terminalizes and archives a merged PR before greeting.
- DONE: Propose the documentation delta for the observable runtime-test contract.
  `docs/runtime-live-ci.md` gains one paragraph explaining why reference-load order is transcript behavior while durable state remains authoritative for mutation, refusal, and archive outcomes.
- DONE: Request independent staff review for this complex skill-integration ideation before gate presentation.
  The First Officer received the review request and the key no-load/terminal-fixture distinction through the Codex mailbox.
- DONE: Run the repository baseline gates required for this stage.
  `go test ./...` and `go test ./... -race` both exited 0 on the unchanged product tree; only this state-checkout entity body is authored by this worker.

### Summary

The design restores the original loading economy without reviving the rejected callable-skill surface: one eager shared core, exact base-anchored delayed reads, and explicit trigger/order rules. It makes the boundary falsifiable through a fixed byte delta, normalized Claude/Codex read traces, preserved durable outcomes, and planted regressions; the existing shallow-boot fixture is correctly retained as a terminal-load positive rather than misused as a no-load test.
