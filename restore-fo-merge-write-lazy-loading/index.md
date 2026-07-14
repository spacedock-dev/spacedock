---
title: Restore lazy loading for first-officer merge and write cores
status: implementation
source: Captain correction after fresh-session boot trace, 2026-07-13
started: 2026-07-13T15:56:42Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-restore-fo-merge-write-lazy-loading
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
- Replace #495's eager-topology assumptions with one structural topology/reachability test and host-neutral read-order oracles. Contractlint is deliberately limited to topology: one eager shared import, two literal delayed paths whose canonical files exist, absence of eager merge/write imports, and absence of wrapper skills. Behavioral stream/live assertions exclusively prove when each path was read, whether a lookup hunted or retried, and whether a mutation or merge action came first.
- Seed a non-empty `mod-block=merge:*` into the existing Claude/Codex positive `shallow-boot` terminal journey. The resumed run must read merge-core before its first merge-guard attempt while retaining the scenario's existing done/PASSED/archive assertions; `merge-hook-guardrail` retains the complementary refusal proof.

### Exact contract wording

The implementation should express the delayed cues in this behaviorally equivalent form (wordsmithing may shorten it without weakening the ordered verbs):

> Resolve delayed FO references only from the `Base directory for this skill` supplied when `spacedock:first-officer` loaded. Append the literal `references/...` suffix; never use cwd, `Skill(...)`, discovery, or search. If the exact read fails, report the path and halt.
>
> - `<base>/references/fo-write-core.md` — before the first FO-authored write intent or state-changing command: READ, then `«write.classify»`, then mutate.
> - `<base>/references/fo-merge-core.md` — at the first terminal boundary or `mod-block=merge:*` recovery: READ before merge guard, merge hook, terminal mutation, archive, or cleanup.

The State Management sentence changes from “under the eagerly loaded `«write.classify»` write-authority scope” to “after the exact write-core read and `«write.classify»` at the first FO-authored mutation.” The Completion and Gates terminal branch and Mod Hook Convention each call the merge load rule before their first merge-owned action; they do not duplicate the core body. This wording is a normative instruction, not a contractlint grammar: tests do not infer execution order from verb order or fail on equivalent prose rearrangement.

`«write.classify»` currently has no runtime event. `internal/ensigncycle/fo_product_edit_guard_impl_test.go` extracts classification-looking text from model narration with `foWriteClassRe`; that narration is not execution proof and is removed from this task's order oracle. The observable write boundary is exact write-core read before the first mutation attempt, paired with the journey's durable allowed or refused outcome. The canonical classifier-block structural tests remain, but emitting a real classification event would require a separate binary/runtime change and is outside this timing-only repair.

### Observable read-order journeys

| Journey | Required ordered observation |
| --- | --- |
| Cold gate hold (`gate-guardrail`) | shared core → selected runtime adapter → gate presentation; no write-core or merge-core read |
| First mutation (`filing`) | shared/runtime → exact write-core read → `spacedock new`; no merge-core read |
| Non-terminal dispatch (`rejection-flow`) | exact write-core read precedes first dispatch-state mutation; merge-core remains absent through the human gate |
| Terminal guard (`merge-hook-guardrail`) | exact write-core read → exact merge-core read → attempted terminal mutation; existing refusal remains the outcome |
| Merged-PR terminalization (`shallow-boot`) | exact write-core read → exact merge-core read → terminal/archive commands; durable merged+archived state remains correct |
| Resume `mod-block=merge:*` (`shallow-boot`, seeded) | boot identifies the non-empty merge mod-block → exact write-core read → exact merge-core read → first merge-guard attempt; no merge action precedes the merge-core read, and durable done/PASSED/archive assertions remain intact |

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

### Cycle 2 live recovery spike

The staff-review risk was exercised end to end on 2026-07-14 with a throwaway commissioned single-root workflow. Both Claude and Codex started from an active `resume-merge` entity with `status: implementation`, `verdict: passed`, `mod-block: merge:local-merge`, and durable landed sentinel `pr: pr-merge:99`. Each host loaded its real first-officer contract, booted the workflow, recognized the in-flight merge, and ran `${SPACEDOCK_BIN} merge guard resume-merge --verdict passed --json` without rerunning the merge hook.

- Claude returned `signal: finalized`, created archive commit `c465283`, and left `_archive/resume-merge.md` at `status: done`, `verdict: passed`, cleared `mod-block`, stamped completion/archive times, no active entity, no dispatchable item, and a clean worktree.
- Codex returned the same finalized signal, created archive commit `b3d7643`, and left the same terminal frontmatter, archive location, empty dispatchable set, and clean worktree.
- The captured RED baseline is explicit: Claude read shared → merge → write → runtime before boot; Codex read shared, merge, write, and runtime in its first command before boot. Thus the recovery mechanics and durable assertions are proven reachable on both hosts, while the current eager timing fails the intended lazy boundary. Implementation must make the same live journeys read merge-core after boot identifies `mod-block=merge:*` and before the first merge-guard attempt, without changing either terminal outcome.

## Acceptance criteria

**AC-1 (VALUE) - A fresh first-officer run that reaches its first human gate without write or terminal work loads the shared core and selected runtime adapter but reads neither merge-core nor write-core, reducing the cold resident contract by at least 8,000 bytes on every host.**
Verified by: fresh Claude and Codex `gate-guardrail` runs record normalized read events through the gate stop; shared and the selected adapter are present, both deferred cores are absent, the durable gate/final response remains correct, and the fixed byte registry is at least 8,000 bytes below the current eager baseline. Baseline on 2026-07-13: merge 2,830 B + write 5,843 B = 8,673 B eagerly read; the new cues may consume at most 673 B of that saving.

**AC-2 - The first FO-authored mutation is preceded by one exact canonical write-core read, without broad search, cwd-relative lookup, wrapper invocation, retry, or an earlier mutation.**
Verified by: Claude and Codex `filing` traces normalize to exact write-core read → successful `spacedock new`, with no merge-core read and no `find`, recursive grep/glob, `Skill(spacedock:fo-write-core)`, missing-path retry, or earlier write. Offline planted missing, post-mutation, wrong-path retry, wrapper, and broad-search traces each turn the same behavioral oracle red. Model narration containing `write.classify` is ignored because no real classification event exists; the canonical classifier block remains structurally tested, and the journey's durable allowed/refused outcome remains authoritative.

**AC-3 - Merge-core first loads at a terminal or live merge-mod recovery boundary, after write-core when that boundary is also the first mutation, and nowhere earlier.**
Verified by: the non-terminal `rejection-flow` trace reaches its gate with no merge-core read; `merge-hook-guardrail` reads the exact merge core before terminal mutation and preserves its refusal; and both Claude and Codex `shallow-boot` fixtures start with a non-empty `mod-block=merge:*`, then normalize to boot identifies mod-block → write-core read → merge-core read → first merge-guard attempt while preserving the existing done/PASSED/cleared-mod-block/archive/clean-state assertions. Planted eager, missing, and post-guard reads fail. The cycle-2 live spike proves this recovery path and durable outcome are reachable on both real hosts before implementation.

**AC-4 - The contract topology has one eager canonical reference plus two deterministic, resolvable delayed paths, with no alternative entry surface.**
Verified by: replace `TestFirstOfficerEagerReferencesKeepDispatchCoreDeferred`, `TestFirstOfficerEagerReferenceTopology`, and stale closure comments/assumptions with a topology test requiring only `@references/first-officer-shared-core.md`, selected-runtime instructions, and the two exact base-anchored literal paths. Production `os.Stat` closure verifies each canonical body. Planted eager imports, missing/dangling literal paths, and restored `skills/fo-{write,merge}-core` wrappers each fail. Contractlint makes no timing or read-order claim and does not inspect ordered/reversed prose; AC-1 through AC-3's behavioral stream/live evidence owns those guarantees.

**AC-5 - The loading-boundary repair preserves first-officer mutation, gate, dispatch, terminal, archive, refusal, and teardown behavior across supported hosts.**
Verified by: focused contractlint/order-oracle suites and the affected Claude/Codex live scenarios pass without weakening durable assertions; Pi's static topology/runtime-adapter coverage remains green until Pi has shared live runners. Then run `gofmt -w ./cmd ./internal`, `git diff --check`, `go test ./...`, `go test ./... -race`, and `go test -tags live -run '^$' ./...`. A detached throwaway audit plants eager, missing, reordered, wrong-path, search, and wrapper regressions and observes the intended failures.

AC-4 is the mechanism serving AC-1 through AC-3; it does not pass if the live/order evidence misses the value boundary. AC-5 prevents the byte reduction from passing while behavior regresses.

## Test plan

1. **RED topology (small, offline).** Rewrite the eager-reference tests around a one-import registry and two literal delayed-path records. Run against current main and record failures for eager merge/write plus missing delayed paths. Add non-vacuous structural controls for eager import, missing/dangling canonical file, altered literal path, and wrapper restoration. Do not parse trigger wording, verb order, or reversed prose.
2. **RED order oracle (medium, fixture).** Extend the Claude turn parser's existing `ReadTargets` path and add a Codex `item.* command_execution` normalizer. Emit ordered events for exact core reads, broad search, wrapper skill calls, wrong-path attempts/retries, mutation commands, merge guard/hook, and gate/final boundaries; do not emit classification events from narration. Feed committed host-dialect traces covering the six-row journey table. The current eager gate trace must fail; independently planted missing read, eager read, read-after-mutation/guard, wrong-path retry, wrapper invocation, and broad-search fixtures must each fail. A specifically reversed event stream is the timing-order control; reversed instruction prose is not.
3. **Contract change (small).** Remove the two eager imports, update shared-core header/Deferred Load Points/State Management/Completion and Gates/Mod Hook Convention with one base-address rule and ordered triggers, and leave both canonical bodies byte-identical. Update closure registries/comments that call them eager; preserve their section-anchor tests.
4. **Focused GREEN (medium).** Run `go test ./internal/contractlint -count=1`, parser/order-oracle packages, `go test ./skills/integration/... -count=1`, and fixed-byte metrics. Require exactly one eager import, both delayed paths resolving, every planted control red, and at least 8,000 B cold-surface reduction.
5. **Live behavior (expensive, existing calls).** Add assertions to existing Claude/Codex `gate-guardrail`, `filing`, `rejection-flow`, `merge-hook-guardrail`, and `shallow-boot` runners rather than inventing a model-only scenario. Seed `shallow-boot` with a registered non-empty merge mod-block and require boot-identify → write-core read → merge-core read → first merge-guard attempt, while retaining its existing done/PASSED/cleared-mod-block/archive/clean-state assertions. Grade normalized read order plus every scenario's existing durable state/final-message assertion; keep `merge-hook-guardrail`'s refusal unchanged. Run protected host jobs at the exact implementation SHA; local SKIP is not green evidence. Pi remains compile/static-covered because its shared live scenario runners are separately backlog-scoped.
6. **Documentation.** Update `docs/runtime-live-ci.md` after the durable-state assertion paragraph: “Reference-loading scenarios also inspect tool-call order because the read boundary is the behavior: write-core before the first FO mutation, merge-core before terminal/mod-block handling, and neither during a mutation-free gate hold. Durable state still grades the resulting mutation, refusal, and archive outcome.” No CLI/user-command docs change: command syntax and output are unchanged.
7. **Full gates and adversarial audit (medium).** Run the repository gates in AC-5. In a detached throwaway checkout, independently plant structural regressions (eager merge import, eager write import, deleted/dangling literal path, wrapper directory) and behavioral trace regressions (eager read, missing read, read after mutation, wrong-path retry/search, merge read after guard). Run the focused test or stream oracle after each isolated edit and record the failing assertion, then remove the checkout. The validator also compares both core-body hashes to pre-change main to prove timing-only scope.

Estimated implementation complexity is medium: instruction text and structural tests are small; the main work is a reliable cross-host event normalizer and strengthening existing live assertions. No production Go command behavior changes.

### Feedback Cycles

**Cycle 1 (independent ideation staff review).** The design direction is sound,
but the evidence plan is not ready for the ideation gate:

- AC-3 claims `mod-block=merge:*` recovery while proving it only with an
  authored host-dialect fixture. Add a live Claude/Codex recovery journey, or
  seed merge mod-block recovery into an existing positive terminal journey and
  preserve its durable assertions. This matters because a resumed FO must load
  the deferred merge core before continuing an already in-flight merge.
- Restrict contractlint to topology: reference closure, literal exact paths,
  structural absence, and wrapper absence. Trigger timing and read order must be
  proved exclusively by behavioral stream/live assertions, not ordered wording
  or reversed-prose tests.
- Clarify the observable `write.classify` evidence. Either require a real
  classification event emitted by the driven behavior, or stop treating model
  narration as proof.

Routed back to ideation: update the approach, AC-2 through AC-4, and test plan
together, exercise the live recovery path, then return for independent staff
re-review.

**Cycle 2 (Roborev-first validation, job 915).** Corrected-guideline Roborev
rejected exact
`557f8df3e6a62d34987edda70533375fc48ba8f6..83af540d2a4a37a789a665c9b3a3c871bdc3ebdf`
before downstream tests, push, PR, or CI:

- MEDIUM: `claude-first-officer-runtime.md` and `fo-status-viewer/SKILL.md`
  retain stale claims that merge/write cores are eagerly loaded.
- MEDIUM: `fo_reference_order_test.go` can treat a path mention or failed read
  as a successful core load and misses shell mutation forms including
  redirection, `sed -i`, `mv`, and Git writes.

Routed back to implementation under the captain's standing instruction. These
are approved-scope contract consistency and behavioral-oracle gaps; no
reframing is needed. Correct every stale eager claim, require successful exact
reads, cover the full mutation boundary with adversarial negatives, and return
to a fresh Roborev-first validation without push or CI.

**Cycle 3 (Roborev-first validation, job 988).** Corrected-guideline Roborev
rejected exact
`557f8df3e6a62d34987edda70533375fc48ba8f6..0e38db2bda324b924b3220aef1c494dc4e33b26e`
before downstream tests, push, PR, or CI:

- MEDIUM: the changed shallow-boot journey replaces automatic merged-PR sweep
  discovery with a pre-seeded sentinel/mod-block and prescribes owner reads and
  the exact merge-guard command. Preserve the original discovery journey and add
  a separate outcome-oriented recovery scenario.
- MEDIUM: canonical-read classification accepts a matching suffix under any
  base path instead of the exact loaded first-officer skill base.
- MEDIUM: a compound command that targets a wrong path and then the canonical
  path can suppress the wrong-path hazard because only the aggregate command is
  classified. Classify every target independently.

Captain standing decision: send back without reframing. These are approved-scope
journey preservation and exact behavioral-oracle defects. Preserve automatic
discovery evidence, bind reads to the exact skill base, and reject every wrong
target even when a later target is canonical. Return to fresh Roborev-first
validation without push, PR, or CI.

**Cycle 4 (Roborev-first validation, job 1167).** Roborev reviewed the full
range through exact implementation head
`664d64c32813a534c790cfb80e38aa20ef82465e` and returned three medium findings:

- A successful shared-core read can be classified as failed because the real
  document contains the phrase `is not found`; structured exit/status and a
  canonical document anchor must distinguish success from failure.
- The command classifier lowercases extracted paths but compares them with the
  original-case loaded skill base, so valid reads under uppercase-containing
  paths can be classified as wrong-path reads.
- The merge-recovery assertion searches the entire Markdown body instead of
  isolating YAML frontmatter, allowing stale frontmatter to false-green when
  matching terminal fields appear later in body text.

Captain clarification (2026-07-14): merged-PR discovery or advancement is not
required during the read-only boot greeting; it may occur on `engage`. Job 1167
did not repeat the obsolete before-greet discovery objection. These findings
are proof-oracle defects, not evidence of an observed runtime lazy-load failure.
After the cycle-3 escalation boundary, the captain explicitly approved one
narrow repair round: keep runtime lazy-loading behavior unchanged, align merge
recovery with `engage`, and repair only the three job-1167 proof oracles before
a fresh full-range Roborev review and exact-head live validation.

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

## Stage Report: ideation (cycle 2)

- DONE: Add live Claude and Codex evidence for a resumed `mod-block=merge:*` recovery while preserving durable terminal assertions.
  A throwaway commissioned workflow seeded `mod-block: merge:local-merge` and `pr: pr-merge:99`. Both real hosts ran the shipped recovery guard to `signal: finalized`; Claude archive commit `c465283` and Codex archive commit `b3d7643` each left status done, verdict passed, cleared mod-block, completion/archive timestamps, no active or dispatchable entity, and a clean worktree. Neither host reran the merge hook.
- DONE: Turn the live evidence into an implementation-level positive journey rather than leaving AC-3 fixture-only.
  The revised plan seeds the existing Claude/Codex `shallow-boot` journey with a registered non-empty merge mod-block and requires boot-identify → write-core read → merge-core read → first merge-guard attempt while retaining its existing done/PASSED/cleared-mod-block/archive/clean-state assertions; the existing merge guardrail still owns refusal evidence.
- DONE: Restrict contractlint to structural topology and move every timing/read-order claim to behavioral evidence.
  AC-4 and the RED topology plan now cover only one eager shared import, exact literal delayed paths, file closure, absence of eager merge/write imports, and absence of wrapper skills. Reversed order, eager/missing reads, post-mutation or post-guard reads, wrong-path retries, wrapper calls, and broad search are planted event-stream failures; ordered or reversed prose is never treated as execution.
- DONE: Resolve `write.classify` observability without inventing evidence.
  The audit found no real runtime event: `foWriteClassRe` recognizes model narration. AC-2 and the order oracle therefore ignore classification narration and prove the observable boundary with exact write-core read before mutation plus the durable allowed/refused result. The canonical classifier-block structural tests remain; a machine event is explicitly outside this timing-only task.
- DONE: Align the proposed approach, observable journeys, AC-2 through AC-4, and the seven-step test plan with all cycle-1 review findings.
  The design now assigns topology to contractlint, ordering to normalized host events/live runs, and merge recovery value to preserved durable state. It retains the 8,000-byte cold-surface target and the timing-only hash guard on both canonical core bodies.
- DONE: Re-run the repository baseline gates after the cycle-2 design revision.
  `go test ./...` and `go test ./... -race` exited 0 on the unchanged product tree; `git diff --check` also exited 0. Only this split-root entity body is authored by this worker.

### Summary

Cycle 2 closes the staff-review gaps with real host recovery proof and a sharper evidence boundary. The current eager streams are now a captured RED baseline, the post-change live journey must defer merge-core until the seeded recovery boundary without changing terminal state, contractlint asserts topology only, and narration no longer masquerades as a `write.classify` event.

## Stage Report: implementation

- DONE: Defer merge-core and write-core from cold first-officer startup while preserving one eager shared core and deterministic exact-path lazy loading at the first mutation or merge boundary.
  Commit `83af540d` leaves only `first-officer-shared-core.md` eager, adds base-anchored exact write/merge cues with halt-on-failure semantics, and leaves both canonical core bodies byte-identical (`a347fc4a…` merge, `a350c1cc…` write).
- DONE: Implement topology-only contract tests plus behavioral/event-stream and protected Claude/Codex live assertions for gate, mutation, terminal, refusal, and seeded mod-block recovery journeys, including adversarial order controls.
  Contractlint now proves only eager/deferred reachability and wrapper absence; host normalizers ignore classification narration and the shared order oracle rejects eager, missing, post-mutation/post-guard, wrong-path, wrapper, broad-search, and reversed recovery traces. Both protected runners call the oracle for gate, filing, rejection, merge refusal, and seeded `mod-block=merge:pr-merge` shallow-boot recovery while retaining their durable assertions.
- DONE: Prove at least 8,000 bytes of cold-contract reduction without behavior regression, apply the concrete runtime documentation diff, and pass gofmt, focused, full, race, compile-live, and detached mutation audit gates.
  Cold entry+shared size is 27,020 B versus the fixed 35,095 B eager baseline, saving 8,075 B. `gofmt -w ./cmd ./internal`, focused contractlint/order/integration suites, `go test ./...`, `go test ./... -race`, `go test -tags live -run '^$' ./...`, and `git diff --check` passed; the detached audit independently made eager import, dangling literal, restored wrapper, and reversed-order plants fail before returning green. `docs/runtime-live-ci.md` records the stream-order evidence boundary.

### Summary

The first-officer cold contract again loads one shared core while write and merge owners resolve once, by literal paths under the loader-supplied skill base, at their first owned boundaries. Commit `83af540d` adds topology-only structural guards, cross-host behavioral order oracles, live-runner assertions, the seeded merge-mod recovery journey, and the runtime evidence documentation without changing either canonical core body or any binary/mod behavior.

## Stage Report: validation

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..83af540d and inspect the stored range and all four guideline sections; on any finding stop REJECTED without downstream tests, push, PR, or CI.
  Roborev job 915 stored exact `git_ref` `557f8df3e6a62d34987edda70533375fc48ba8f6..83af540d2a4a37a789a665c9b3a3c871bdc3ebdf` and all four guideline sections, then found stale eager-loading instructions and a false-green behavioral oracle; recommendation: REJECTED.
- SKIPPED: Only after Roborev PASS independently reproduce AC-1 through AC-5: at least 8,000 cold bytes saved, exact lazy read ordering across the journey matrix, topology-only structural proof, unchanged core bodies, and preserved durable host behavior.
  The mandatory first gate failed, so no acceptance-criterion evidence was reproduced or treated as green.
- SKIPPED: Run gofmt/diff, focused/full/race/compile-live gates plus the detached adversarial audit and applicable exact-head Claude/Codex live journeys; report any external authorization blocker without pushing or engaging CI.
  Per the Roborev-first stop rule, no downstream formatting check, test, audit, live journey, push, PR, or CI action ran.

### Summary

Roborev job 915 rejected the exact implementation range with two medium findings. The Claude runtime and status-viewer still describe merge/write cores as eager, while the order oracle can mistake a command that merely mentions a core path for a successful read and misses shell mutation forms such as redirection, `sed -i`, `mv`, and Git writes; both must be corrected with negative controls before fresh validation.

## Stage Report: implementation (cycle 2)

- DONE: Update every affected Claude runtime and status-viewer reference so no shipped instruction falsely claims merge/write cores remain eager, while preserving exact deferred load points.
  Commit `0e38db2b` replaces stale eager claims in the Claude runtime, status viewer, and feedback-rejection flow with the shared core's exact write/merge boundaries; a contractlint consistency guard covers all three shipped surfaces.
- DONE: Make the order oracle require a successful exact canonical read rather than path mention or failed attempt, with negative controls for wrong/failed/retried reads.
  Claude reads now resolve only from matching non-error tool results; Codex reads require completed exit-zero commands. Both hosts additionally require the canonical body's heading in output, so `echo`/`sed` path mentions, empty output, failed reads, and wrong/failed retries remain red instead of becoming load evidence.
- DONE: Detect shell mutation boundaries including redirection, sed -i, mv, and Git writes; run focused/full/race/compile-live and adversarial gates, committing locally only for fresh Roborev-first validation.
  The event normalizer now detects `>`/`>>`, in-place sed, file commands, and Git mutations alongside Spacedock writes. Focused contractlint/order suites, `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, `go test -tags live -run '^$' ./...`, and `git diff --check` passed. A detached audit made a stale eager instruction, path-only command, and failed exact read independently fail; the branch remains local and unpushed with no PR or CI engagement.

### Summary

Cycle 2 closes both Roborev findings without changing the approved loading design. Commit `0e38db2b` makes every shipped cross-reference agree on deferred exact-path loading and upgrades the behavioral oracle from path-text recognition to successful canonical-body evidence across the wider shell mutation surface; both canonical core hashes and the 8,075-byte cold saving remain unchanged.

## Stage Report: validation (cycle 2)

- FAILED: Run corrected-guideline Roborev first on exact range 557f8df3..0e38db2b and inspect the stored range and all four guideline sections; on any finding stop REJECTED without downstream tests, push, PR, or CI.
  Roborev job 988 stored exact `git_ref` `557f8df3e6a62d34987edda70533375fc48ba8f6..0e38db2bda324b924b3220aef1c494dc4e33b26e` and all four guideline sections, then returned three medium findings; recommendation: REJECTED.
- SKIPPED: Only after Roborev PASS independently reproduce AC-1 through AC-5: at least 8,000 cold bytes saved, successful canonical lazy-read ordering across the journey matrix, topology-only structural proof, unchanged core bodies, and preserved durable host behavior.
  The mandatory first gate failed, so no acceptance-criterion evidence was reproduced or treated as green.
- SKIPPED: Run gofmt/diff, focused/full/race/compile-live gates plus the detached adversarial audit and applicable exact-head Claude/Codex live journeys; report any external authorization blocker without pushing or engaging CI.
  Per the Roborev-first stop rule, no downstream formatting check, test, audit, live journey, push, PR, or CI action ran.

### Summary

Roborev job 988 rejected the exact cycle-2 range. It found that the changed shallow-boot fixture removes automatic merged-PR discovery coverage and prescribes the implementation, that suffix-only checks accept a core under the wrong skill base, and that a compound wrong-path-then-canonical retry can evade the hazard detector; the implementation needs separate discovery/recovery journeys and per-target exact-base classification before fresh validation.

## Stage Report: implementation (cycle 3)

- DONE: Preserve the original automatic merged-PR sweep discovery journey and add a separate outcome-oriented mod-block recovery scenario without pre-seeding away discovery or prescribing exact owner commands.
  Commit `52d4e92f` restores shallow-boot's original open-PR discovery fixture and adds a distinct `merge-mod-block-recovery` journey whose prompt specifies only the durable resumed-merge outcome. Both Claude and Codex runners preserve the terminal/archive/clean-state oracle, and commit `664d64c3` registers the new shared scenario in the locked seed-scenario documentation.
- DONE: Require successful canonical core reads under the exact loaded first-officer skill base, rejecting matching suffixes rooted elsewhere.
  The host normalizers derive the loaded first-officer base from the real skill-load event (with a successful canonical `SKILL.md` fallback) and accept a merge/write core only at the exact base-plus-literal path. Claude and Codex wrong-base/same-suffix negatives fail, while the positive host-order fixtures remain green.
- DONE: Classify every command target independently so a wrong-path access remains hazardous even when the same compound command later reads the canonical target; run focused/full/race/compile-live and detached gates locally only.
  Each shell read target is now classified independently; compound wrong-then-canonical controls retain `wrong-core-path` and fail. Focused order/recovery/doc-lock tests, `gofmt -w ./cmd ./internal`, `go test ./...`, corrected-PATH `go test ./... -race`, `go test -tags live -run '^$' ./...`, `git diff --check`, and the detached adversarial audit passed locally. The detached plants independently proved shallow discovery cannot be replaced, a wrong skill base fails, and a compound wrong-then-canonical read stays hazardous; no branch push, PR, or CI action occurred.

### Summary

Cycle 3 separates automatic merged-PR discovery from resumed mod-block recovery and makes the recovery prompt judge the durable outcome instead of a prescribed command sequence. Exact loaded-base binding and per-target classification close both Roborev false-green paths while preserving the canonical core hashes (`a347fc4a…`, `a350c1cc…`) and the 8,075-byte cold-contract saving.
