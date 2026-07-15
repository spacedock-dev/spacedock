---
title: Restore lazy loading for first-officer merge and write cores (clean reset)
status: done
source: Clean reset from rejected 1k implementation, captain direction 2026-07-14
started: 2026-07-14T14:12:36Z
completed: 2026-07-15T03:28:37Z
verdict: passed
score:
worktree: .worktrees/spacedock-ensign-restore-fo-merge-write-lazy-loading-reset
issue:
milestone: 0.25.0
id: gk7ceyrs4496jgp535w3awfp
mod-block:
pr: pr-merge:512
archived: 2026-07-15T03:28:37Z
---

Restore the intended first-officer loading boundary: boot reads the shared core and active runtime adapter, while write authority loads at the first FO-authored mutation and merge handling loads only at a terminal or merge-mod recovery boundary.

## Problem

The earlier fix for delayed-reference filesystem hunting overcorrected by eagerly importing `fo-write-core.md` and `fo-merge-core.md` from the first-officer entry skill. A mutation-free interactive run that stops at a gate now pays for mutation and terminal ceremony it never uses. The defect was nondeterministic discovery, not lazy loading: deferred reads must use the loader-supplied first-officer base plus one literal canonical suffix, never cwd, wrapper-skill discovery, alternate paths, or search.

The prior implementation branch and its proof harness were rejected. This task intentionally carries no prior stage reports, feedback cycles, parser design, or implementation plan.

## Required outcome

- `skills/first-officer/SKILL.md` eagerly imports only the shared core; the active runtime adapter remains a boot read, while write and merge remain canonical deferred references.
- The contract keeps the intended triggers: write authority at the first FO mutation, merge handling at the first terminal or `mod-block=merge:*` recovery boundary, with write before merge when both apply.
- Merged-PR discovery is not required on boot; it may happen on `engage`.
- Do not change mutation resolution, merge resolution, routing, or terminal semantics. If a resolution is broken, that is a separate defect shared by eager and lazy loading.
- Do not claim semantic runtime order from a bespoke transcript or command classifier. Release value is established by the post-PR same-resolved-model shallow-boot ledger comparison.

## Mechanism/value trace

- Served value: reduce the real shallow-boot token window while preserving supported-host outcomes.
- Simplest route: keep the small import/load-cue change, enforce only eager-import/reference closure structurally, retain existing durable scenario assertions, and let post-PR Claude CI compare the measured `shallow-boot-window` ledger against the published v0.24.0 baseline for the same resolved model.
- Rejected mechanism: the contract-specific semantic trace observer is deleted. Do not replace it with shell/path/payload classifiers, runtime instrumentation, a command parser, operation-language interpreter, controller, lifecycle layer, daemon, lease, or recovery protocol.

## Acceptance criteria

- **AC-1 (VALUE / RELEASE PROOF):** Post-PR Claude CI emits `shallow-boot-window` and compares the same resolved model against the published v0.24.0 ledger (Sonnet total 47,936, cache creation 638, pre-greet peak 9,864). The comparison must demonstrate the intended cold-window improvement before merge; record client/model drift, and do not treat a cross-model delta as evidence.
- **AC-2:** Structural coverage proves the entry skill has one eager canonical shared-core import and that the write/merge deferred references resolve to their canonical non-empty files. It does not infer runtime order from instruction prose or transcripts.
- **AC-3:** Existing gate, filing, and terminal/recovery scenario assertions retain their durable success/refusal/archive outcomes without a contract-specific semantic trace oracle.
- **AC-4:** Focused structural checks, `go test ./...`, `go test ./... -race`, and relevant exact-head local Codex journeys are green before Roborev. Local Claude live is not required; the same-model Claude ledger comparison remains the post-PR merge gate.

## Feedback Cycles

### Cycle 5 — PR live evidence correction

- **Evidence defect (1k scope):** Exact-head PR run `29382760645`, Opus rerun job `87252929149`, invoked `fo-status-viewer` before the shallow-boot greet. The pre-existing `assertGreetInvokesNoDeferredFOSkill` stream oracle then failed. That oracle claims per-skill runtime order from a transcript and conflicts with the captain's direction that the same-model shallow-boot token ledger—not a contract-specific harness—proves this task's deferred-loading value.
- **Correction:** Delete the pre-greet deferred-FO-skill assertion, its dedicated fixtures/controls, and stale runner/spec commentary. Do not replace it with another Skill/Read/path/command/event classifier. Retain the durable mutation-free gate/team assertions and the quantitative `shallow-boot-window` ledger.
- **Separate evidence defect:** The same Opus job correctly advanced and dispatched `approved-gate`, then the unchanged keep-moving regex treated `no "want me to advance?" pause` as an actual permission request. That free-narration false positive is outside 1k and must be fixed in a separate, sequenced test-infrastructure task rather than folded into this branch.
- **Return conditions:** Focused/full/race and relevant exact-head local Codex checks green, followed by an exact-head Roborev request. Local Claude remains unnecessary. Return without touching the separate keep-moving oracle.

## Stage Report: implementation

- DONE: Deliver the must-have deferred write/merge loading boundaries while preserving real gate, mutation, and terminal outcomes.
  Commit `09288711` leaves only the shared core eager, defers exact write/merge references, and keeps gate hold, atomic filing, and terminal refusal outcomes green on both hosts.
- DONE: Use existing supported-host event streams and scenario runners; do not implement a shell parser, runtime controller, or parallel lifecycle, and obtain exact-head local Codex live green before Roborev.
  Existing gate, filing, and merge-hook runners now grade emitted Claude/Codex path and action events; final exact-head Codex was 3/3 green before any Roborev request.
- DONE: Keep the change narrow; if another controller/lifecycle layer appears or no visible product proof emerges within 90 minutes, stop and request a design reset.
  The change is limited to load cues, structural guards, and one bounded event-field observer; it adds no runtime product code, controller, lifecycle, or shell grammar.
- DONE: AC-1 (VALUE)
  Fresh exact-head Claude and Codex gate journeys read shared core plus their adapter and neither deferred core; deferring the unchanged write+merge files saves `5,843 + 2,830 = 8,673` cold bytes.
- DONE: AC-2
  Exact-head Claude and Codex were each 3/3 green: filing reads write before `new`, terminal reads write then merge before `status --set`, and durable success/refusal/no-archive outcomes remain enforced.
- DONE: AC-3
  Contractlint proves the sole eager import, the two canonical loader-base suffixes, reference closure, ceremony anchors, and absence of wrapper core skills; runtime order comes only from host events.
- DONE: AC-4
  Focused trace controls, `go test ./...`, `go test ./... -race`, exact-head Codex 3/3, and exact-head Claude 3/3 all passed on commit `09288711`.

### Summary

Restored the intended cold boundary without changing mutation, merge, routing, or terminal resolution semantics. The existing supported-host journeys now make the load order visible and falsifiable while retaining their durable outcome assertions.

### Roborev

- Exact-head review requested for `09288711f8c9069eb7afefe5057c7c737236f4d2` after the recorded Codex and Claude live-green runs. No prior review existed for that SHA.
- Request: Roborev job `1423`, agent `codex`, reasoning `thorough`; status `done`, verdict `F` (finished 2026-07-14T23:27:09+08:00).
- Medium finding at `internal/ensigncycle/fo_deferred_load_trace_test.go:39`: `observeLoad` can count any successful event containing a matching path suffix as a completed read. Roborev recommends requiring the exact loader-supplied base and host-native read operation, correlating a genuine full-file Codex read, and adding negative cases for path-only mentions, partial reads, and same-suffix alternate roots.
- Medium finding at `internal/ensigncycle/fo_deferred_load_trace_test.go:172`: the checks cover atomic filing and `status=done`, but not the broader set of FO-authored mutations or every existing mutating live scenario. Roborev recommends detecting the first general FO mutation across those scenarios while retaining the specialized terminal merge-order assertion.
- No code was changed in response; both concrete findings await first-officer feedback routing.

### Feedback Cycles

- Cycle 1 — detached adversarial audit of exact head `09288711` on throwaway detached worktrees.
  - Material (AC-1/AC-2/AC-4 evidence defect): planted path-only `printf`, partial `head`, and same-suffix alternate-root events; all were accepted as completed canonical reads, so the new negative controls went RED.
  - Material (AC-2/AC-4 evidence defect): planted an ordinary `status --set` before the write read and a later instrumented filing; the early mutation was invisible, so the new negative control went RED.
  - Current raw Claude/Codex events show the intended behavior, so no outcome defect was observed. This is a narrow evidence fix: bind the exact installed base and full host-native read result, add the negative controls, and apply first-mutation checks to existing mutating scenarios. No parser, controller, or lifecycle is required; this is not a design reset.

## Stage Report: validation

- DONE: Independently verify AC-1: fresh Claude and Codex gate holds omit both deferred cores, retain the adapter/shared core, and save at least 8,000 cold bytes.
  Exact-head raw streams read the full shared core plus selected adapter and neither deferred core; `5,843 + 2,830 = 8,673` cold bytes are removed from the eager baseline.
- FAILED: Independently verify AC-2 and AC-3: existing host events show exact write/merge load order before real owned actions, durable filing/terminal outcomes remain unchanged, and structural checks prove only canonical imports without prose-derived runtime order.
  AC-3 and today's filing/terminal outcomes pass, but AC-2 evidence is invalid: the observer accepts fake/partial/alternate-root reads and ignores an earlier ordinary mutation.
- DONE: Reproduce focused, full, race, and exact-head local live checks at 09288711; run the required detached adversarial audit for shipped contract/scaffolding, account for the exact-head Roborev result, and issue PASSED or a precisely classified REJECTED recommendation.
  Focused structural/trace controls passed; isolated `go test ./...` and `go test ./... -race` passed; Codex and Claude were each 3/3 green; detached audits reproduced both Roborev Medium findings.
- DONE: AC-3 structural topology.
  The sole eager shared-core import, resolvable write/merge cores, 8,673-byte value, and reference closure passed without using instruction prose to infer runtime order.
- FAILED: AC-4 behavioral proof boundary.
  Exact-head live journeys are green, but the harness can false-pass claim-breaking traces; Roborev job `1423` independently returned `F` with the same two Medium evidence findings.
- DONE: Recommendation: REJECTED — narrow evidence fix.
  Harden the existing host-event observer and extend it across existing mutating scenarios; no product-semantics change or architecture reset is warranted.

### Summary

The current Claude and Codex journeys visibly perform the intended lazy reads and preserve durable outcomes, and all deterministic gates pass at `09288711`. Validation rejects the shipped evidence boundary because two detached adversarial edits proved it can certify behavior it did not observe; route a narrow test-harness correction back to implementation.

## Stage Report: implementation (cycle 1)

- DONE: Bind deferred-load evidence to the exact loader-supplied installation and a successful host-native full-file read.
  Commit `c7def2b7` requires correlated, non-partial Claude `Read` results and full canonical Codex command output at exact installed paths; commit `f745bba1` binds Codex to its unique loader-visible isolated cache entry. Path-only, partial-read, and same-suffix alternate-root controls were RED on `09288711` and are green at exact head.
- DONE: Detect the first general FO-authored mutation across existing mutating scenarios while retaining specialized filing and terminal order checks.
  Fixed supported-event markers now catch ordinary early `status --set` and other established mutation surfaces; Claude and Codex rejection, escalation, merge-triage, smallest-mechanism, keep-moving, and shallow-boot runners apply the write-before-first-mutation boundary. Filing still requires write before `new`, and terminal flow still requires write then merge before its action.
- DONE: Preserve product behavior and keep the correction inside the existing observer and scenario runners.
  No product semantics, shell parser, runtime controller, or lifecycle layer changed. Focused controls, live-tag compilation, `go test ./...`, and `go test ./... -race` passed at `f745bba1592837aa8697d321177cc3a07a224115`; exact-head Codex and Claude were each 3/3 green.
- DONE: Prove the Claude gate through supported stored-login setup without changing repository or operator credential state.
  The exact-head Claude pass used a disposable profile populated from the already-working operator login and normal external Go caches. Later unchanged attempts that used `Bash cat` or skipped the write-core read were correctly rejected by the hardened observer rather than false-passing.
- DONE: Obtain an exact-head Roborev result after all required gates.
  Roborev job `1433` reviewed `f745bba1592837aa8697d321177cc3a07a224115` with agent `codex`, reasoning `thorough`; status `done`, verdict `P`, with no issues found.

### Summary

Cycle 1 closes both false-green evidence gaps without changing lazy-loading or workflow semantics. Exact installed-base/full-read proof and first-mutation coverage are now adversarially falsifiable, with deterministic gates, supported-host live evidence, and exact-head Roborev green.

### Feedback Cycles

- Cycle 2 — detached adversarial re-review of exact head `f745bba1592837aa8697d321177cc3a07a224115`.
  - Material (AC-2/AC-4 evidence defect): a single Codex command that reads the exact write-core path and also retries a same-suffix alternate-root path is accepted because one exact-path match masks the noncanonical occurrence.
  - Material (AC-2/AC-4 evidence defect): early supported `state commit`, `status --archive`, and `merge guard` mutations are absent from the fixed first-mutation classifier and therefore false-pass before the write-core read.
  - Material (AC-4 evidence defect): a lexical `/tmp` installed base does not match host events reported through macOS's canonical `/private/tmp` alias; the same three Codex scenarios were rejected 0/3 through the alias and passed 3/3 through the canonical path.
  - No outcome defect was observed. Keep the correction in the observer: reject every non-exact core occurrence, cover supported mutation families, and compare canonical installed paths; do not change product semantics, parsing, controllers, or lifecycle.
- Cycle 3 — Roborev job 1445, exact head `5e78f71c`, ESCALATED.
  - Material evidence defect (AC-2): common terminal workflows can execute `merge guard` or `status --archive` after write-core but before merge-core and still pass a later terminal action.
  - Material evidence defect (AC-2/AC-4): supported stdin JSON dispatch with `"advance":true` remains invisible to the mutation classifier.
  - Material evidence defect (AC-4): valid common checkout paths containing spaces are truncated and falsely rejected.
  - Release-scope triage: all three affect named ACs and common or explicitly exercised workflows; the Linux-only non-distinct `/tmp` alias control is a deferred risk with synthetic-symlink revisit condition.
  - Design-reset trigger: another observer patch would require shell quoting, JSON payload semantics, and further command-family interpretation. The branch is now 749 additions across 12 files, including a 575-line observer, while the shipped contract delta is small. No automatic cycle 4 dispatch.
  - Captain resolution: skip the contract-specific runtime trace harness. Preserve the deferred-load contract and minimal structural coverage; use the same-model shallow-boot token comparison against the published v0.24.0 ledger as the release proof. Remove the semantic observer instead of repairing its command, path, or payload classifiers.
- Post-PR release gate — Runtime Live E2E run `29378927729`, exact rebased head `cee47a55`, NOT MERGEABLE.
  - AC-1 value result: every CI job passed, but same-model Sonnet `claude-sonnet-5` shallow-boot-window usage was 61,781 total tokens across 22 turns versus v0.24.0's 47,936 across 18 turns. Cache creation was 639 versus 638 and pre-greet peak improved to 9,334 from 9,864, but the requested return to the v0.24.0 total-token level did not occur.
  - Evidence-workload defect: the existing `shallowBootPrompt`, fixture, assertions, and scenario description still require S7b merged-PR finalization before greet. The captured run obeyed that stale requirement, read both deferred cores, mutated and archived state, then greeted. This contradicts the captain-approved outcome that merged-PR discovery may wait until `engage`, so the metric is not a mutation-free shallow boot.
  - Routed correction: reshape the existing shallow-boot scenario into a mutation-free interactive greet/gate hold and remove its before-greet merged-PR requirement. Do not add a deferred-load trace observer or replacement classifier. Preserve merged-PR handling at `engage` through existing supported outcome coverage, or identify the smallest separate coverage gap without widening 1k into new infrastructure.

## Stage Report: validation (cycle 2)

- FAILED: Reproduce the path-only, partial-read, and alternate-root negative controls and verify only an exact loader-base host-native full-file read satisfies AC-2/AC-4.
  The shipped individual controls and exact-read positive passed, but a detached combined exact-plus-alternate-root control false-passed at `f745bba1`; one exact occurrence currently masks the noncanonical retry.
- FAILED: Reproduce the early ordinary-mutation control across existing mutating scenarios, retain filing/terminal ordering, and confirm no product semantic, shell parser, controller, or lifecycle change.
  Existing ordinary-status, filing, and terminal controls pass, and the two correction commits change tests only. Detached early `state commit`, `status --archive`, and `merge guard` controls all false-passed because those supported mutators are not classified.
- DONE: Independently verify exact head `f745bba1` with applicable focused/full/race and Claude/Codex live evidence, account for Roborev job `1433`, and issue PASSED or precisely classified REJECTED.
  Focused controls and live-tag compilation passed; isolated `go test ./...` and `go test ./... -race` passed. Exact-head Codex passed 3/3 from a canonical artifact path; the `/tmp` alias run exposed the observer false-negative above. Implementer-owned Claude 3/3 is recorded; local Claude rerun is not a 1k gate.
- DONE: Account for independent review evidence.
  Roborev job `1433` reviewed exact head `f745bba1`, completed with verdict `P`, and reported no issue; the detached controls above exercise combinations and mutation categories outside that review.
- DONE: Preserve the quantitative shallow-boot condition as a post-PR CI gate.
  CI must compare the same resolved model against the published v0.24.0 `shallow-boot-window` baseline (Sonnet total 47,936, cache creation 638, pre-greet peak 9,864), recording client/model drift; the 8,673-byte core reduction alone is not token proof. Its absence locally is not a validation failure.
- DONE: Recommendation: REJECTED — narrow evidence fix.
  Three observer gaps can still reject valid canonical behavior or certify claim-breaking traces. Correct the evidence harness without a product-semantics change or design reset, then rerun the affected controls and gates.

### Summary

Exact-head deterministic gates, canonical-path Codex outcomes, implementer-owned Claude evidence, and Roborev are green, but independent adversarial controls still falsify AC-2/AC-4's proof boundary. Validation rejects `f745bba1` for three narrow observer defects and routes an evidence-only correction back to implementation; the quantitative same-model shallow-boot comparison remains a post-PR CI requirement.

## Stage Report: implementation (cycle 2)

- FAILED: Canonicalize the loader base and observed paths, require successful full-file exact reads, and reject every noncanonical core occurrence even when an exact occurrence also exists; cover combined exact+alternate and `/tmp`↔`/private/tmp` cases.
  Commit `5e78f71c` makes the assigned combined-path and macOS alias controls pass after they failed at `f745bba1`, but Roborev found a valid loader path containing spaces is truncated and falsely rejected.
- FAILED: Extend the existing bounded supported-mutation classifier to catch early `state commit`, `status --archive`, and `merge guard` alongside current mutations; add adversarial controls and stop if this requires a general shell or operation-language parser.
  All six Claude/Codex controls for the three requested families went RED then green without a parser, but Roborev found JSON `"advance":true` remains invisible and merge/archive are not graded as merge-core terminal actions.
- FAILED: Preserve product semantics, run focused/full/race and exact-head local Codex gates, request exact-head Roborev only after green, and return any rejection without further mutation for cycle-3 escalation.
  Product semantics were untouched; focused/live-tag, `go test ./...`, `go test ./... -race`, and exact-head Codex gate/merge-hook/filing 3/3 passed at `5e78f71cd63503f9e49c63cc400221db0b6cb304`. Roborev job `1445` then returned verdict `F`; no code changed afterward.

### Roborev and release-scope triage

- MATERIAL: `merge guard` and `status --archive` can occur after write-core but before merge-core and still pass a later terminal-status action. These are common terminal workflows under AC-2's explicit terminal ordering promise.
- MATERIAL: supported stdin JSON dispatch with `"advance":true` is present in existing reviewer-reuse fixtures and can false-pass the first-mutation boundary; this is an explicitly exercised supported workflow.
- MATERIAL: quoted loader paths containing spaces are valid common checkout paths and currently false-reject exact reads, undermining AC-4 proof portability.
- DEFERRED RISK: the `/tmp` alias regression control is non-distinct on Linux. Revisit with an explicit synthetic symlink if non-macOS CI must independently exercise alias resolution; the assigned macOS `/tmp`↔`/private/tmp` case was distinct and green locally.

### Summary

Cycle 2 closes the three routed examples with a bounded test-only observer change and all required local gates green. Exact-head Roborev nevertheless rejected three material proof gaps, so the last automatic correction stops without another code mutation and returns for cycle-3 captain escalation.

## Stage Report: implementation (cycle 3)

- DONE: Delete the semantic trace observer and its scenario integrations without replacing it.
  Commit `ea60aa7c1a44686b0d5bfe3e4dccb63b8ab62c8f` deletes the 575-line `fo_deferred_load_trace_test.go` oracle and removes its runner seam, cache-path resolver, scenario assertions, and observer-only stream fields. The Claude, Codex, PTY, and streamwatch runner files now match pre-task baseline `6e99954f`; no shell, path, payload, mutation, or terminal classifier replaces the oracle.
- DONE: Preserve the smallest deferred write/merge product contract and only structural closure/import tests.
  The retained branch delta is seven product/contractlint files with 59 insertions and 79 deletions against `6e99954f`: one eager shared-core import, canonical deferred write/merge references, adapter boot routing, and structural reference closure. The local 8,000-byte proxy test was removed so repository byte count cannot stand in for host token-window evidence; focused contract and live-fixture checks passed.
- DONE: Use the same-resolved-model `shallow-boot-window` ledger comparison as release proof, with focused/full/race and exact-head Codex green before Roborev; local Claude is not required.
  AC-1 records the published v0.24.0 Sonnet ledger (total 47,936, cache creation 638, pre-greet peak 9,864) as the post-PR comparison baseline and rejects cross-model deltas as proof. `gofmt -w ./cmd ./internal`, standalone `go test ./...`, and standalone `go test ./... -race` completed successfully; exact-head Codex gate-guardrail, merge-hook-guardrail, and filing journeys passed 3/3 in 171.99s before review. Roborev job `1446` reviewed the same head and returned `F` only because it requested host-specific runtime-event checks recreating the explicitly prohibited semantic observer; the finding was intentionally declined and no code changed after review.

### Summary

Cycle 3 removes the contract-specific semantic trace architecture and leaves the minimal deferred-core contract guarded only by structural closure/import checks and existing durable scenario outcomes. All required local gates and exact-head Codex journeys are green; the quantitative same-model Claude ledger comparison remains the post-PR release proof, and the reviewer request to restore semantic tracing was not implemented.

## Stage Report: validation (cycle 3)

- DONE: Verify the exact branch contains the minimal deferred-core contract and structural closure/import coverage, with the semantic trace oracle and all replacement classifiers absent.
  Exact head `ea60aa7c` differs from pre-task `6e99954f` only in seven contract/contractlint files: one eager shared-core import, one canonical write cue, one canonical merge cue, and non-empty target checks. `internal/ensigncycle` is byte-identical to baseline; the 575-line oracle, runner seams, path/payload/mutation classifiers, and observer-only fields are absent.
- DONE: Reproduce the applicable focused/full/race and exact-head local Codex evidence, then assess Roborev job 1446 against the captain-revised ACs and release-scope policy.
  Focused structural and durable-outcome controls, live-tag compilation, `go test ./...`, and `go test ./... -race` passed. Exact-head Codex gate hold, merge-hook refusal/no-archive, and atomic filing passed 3/3 in 208.76s; a detached eager-import plus dangling-target audit correctly went RED.
- DONE: Issue PASSED or a precisely classified REJECTED recommendation, preserving the same-resolved-model post-PR shallow-boot ledger comparison as the unresolved merge gate rather than demanding a contract-specific trace harness.
  Recommendation: PASSED — PR-ready, not yet merge-ready. The post-PR Claude ledger must still demonstrate the cold-window improvement against the same resolved v0.24.0 model before merge; no cross-model delta qualifies.
- SKIPPED: AC-1 final release proof.
  Intentionally unavailable before the PR live lane. The published `journey-costs-v0.24.0.json` independently confirms Sonnet `claude-sonnet-5`, total 47,936, cache creation 638, and pre-greet peak 9,864; CI already emits `shallow-boot-window`, preserves resolved-model/client metadata, and posts the delta for the required pre-merge decision.
- DONE: AC-2 structural proof.
  The focused guards independently read the entry import topology and filesystem targets; the detached audit proved both an extra eager import and a broken canonical target fail rather than self-certify prose.
- DONE: AC-3 durable outcome proof.
  The live runners retain gate entity equality/status/refusal, archive absence, and filed-entity existence/atomic-create assertions while carrying no contract-specific deferred-load observation.
- DONE: AC-4 local gate proof.
  All required deterministic gates and applicable exact-head Codex journeys are green; local Claude is not required by the revised criterion.
- DONE: Release-scope triage of Roborev job `1446` (verdict `F`).
  Its sole Medium asks for host-event proof the revised ticket explicitly rejects. For supported-host FO operators, eager cold cost is caught by AC-1 and outcome harm by AC-3; no current AC or safety/compatibility boundary promises per-event semantic order. Revisit instrumentation only if the same-model ledger or a durable supported-host outcome regresses.

### Summary

The bounded reset is structurally minimal, independently falsifiable, and preserves the existing durable supported-host outcomes at exact head `ea60aa7c`. Validation recommends PASSED for PR creation while holding merge on the quantitative same-resolved-model Claude comparison; Roborev 1446 requests the superseded mechanism rather than identifying a current release-boundary failure.

## Stage Report: implementation (cycle 4)

- DONE: Reshape the existing shallow-boot scenario, fixture, assertions, and specification into a mutation-free interactive greet/gate hold; remove the S7b merged-PR-before-greet requirement without adding a new observer or harness.
  Commits `4d6f9379` and `ea63bcd4` remove the merged-PR entity, pr-merge mod, merged-`gh` shim, S7b prompt/outcome checks, and stale specification text from the existing workload. The fixture now contains only the held gate; its structured greet rejects resolved/already-engaged messages while durable checks require unchanged gate state, no archive/worktree, and no team artifact.
- DONE: Preserve merged-PR discovery at `engage` through the smallest existing supported-outcome coverage, documenting any genuine separate gap instead of widening 1k into new infrastructure.
  `TestFOEngageRetainsStartupPRAdvancement` pins `engage` → `«hooks.run»("startup")` exactly once and startup-mod ownership of live PR advancement; existing `TestSweepNonEmptyNamesRegisteredStartupModNextStep` proves a merged PR routes to the registered startup mod. The one-shot Codex journey proves persisted no-mutation/no-artifact outcomes but cannot exclude transient engage/dispatch commands; that command-sequence gap is documented and intentionally receives no semantic event classifier or interactive controller.
- DONE: Run focused/full/race and relevant exact-head local Codex checks, request exact-head Roborev after green, and leave the branch ready for a new same-model Claude shallow-boot CI measurement.
  Prompt/fixture/oracle, engage-boundary, sweep-outcome, and boot-identify focused checks passed; `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` passed at exact head `ea63bcd44c93de44483924f404e21f5141867cee`. Exact-head Codex shallow-boot passed in 29.46s. Roborev job `1469` returned `F` only by requesting a genuine-interactive transport or both-host command-stream classifiers, the two mechanisms this correction explicitly excludes; no code changed afterward. Local Claude was not required, and the branch is ready for the same-resolved-model Sonnet `shallow-boot-window` CI comparison.

### Summary

Cycle 4 corrects the measured workload rather than rebuilding the deleted trace architecture: shallow boot now ends after local identify, a structured held-gate greet, and no persisted mutation or dispatch artifact. Engage retains startup-mod PR advancement through existing structural/outcome coverage; all local gates are green, and the unresolved release proof is the new PR's same-model Claude ledger comparison.

## Stage Report: validation (cycle 4)

- DONE: Independently verify the revised shallow-boot workload is a mutation-free interactive greet/gate hold and that its durable negative controls cannot accept an engaged, resolved, dispatched, archived, or team-created state.
  Focused prompt/fixture/identify/oracle controls passed. A detached audit removed the team and prospective-engage checks and reintroduced `merged-pr.md`; all three planted regressions went RED. Exact-head Codex passed 1/1 in 31.02s with unchanged gate/no archive/worktree; ancillary raw inspection showed only the version gate and `status --boot --identify --json`, followed by the held-gate/engage-next greet and no collaboration event.
- DONE: Verify merged-PR advancement remains owned by `engage` startup coverage without a new semantic trace or controller, and triage Roborev 1469 against the revised ACs and release scope.
  `TestFOEngageRetainsStartupPRAdvancement` retains the single startup-hook ownership, and `TestSweepNonEmptyNamesRegisteredStartupModNextStep` exercised a merged PR and returned the registered `_mods/pr-merge.md` next step. The deleted trace oracle remains absent and no transport/controller/classifier was added.
- DONE: Reproduce applicable focused/full/race and exact-head Codex evidence; issue PASSED/REJECTED while retaining the new same-model Sonnet ledger as the hard pre-merge proof.
  Focused checks, live-tag fixture compilation, `go test ./...`, and `go test ./... -race` passed at `ea63bcd4`; exact-head Codex shallow boot passed. Recommendation: PASSED — PR-update-ready, with merge forbidden pending the corrected same-resolved-model Sonnet ledger.
- SKIPPED: AC-1 corrected-workload release proof.
  Run `29378927729` measured the superseded S7b workload at 61,781 versus v0.24.0's 47,936 and cannot satisfy AC-1. PR #512 still points at `cee47a55`; after updating it to `ea63bcd4`, the new Sonnet `claude-sonnet-5` comparison must reach the v0.24.0 level before merge.
- DONE: AC-2 structural proof.
  The unchanged deferred-core topology/closure guards remain green in the full suite; cycle 4 adds no semantic-order inference.
- DONE: AC-3 durable outcome proof.
  The corrected shallow fixture contains only the held gate, and the live/offline assertions independently reject mutated/resolved gate state, archive, worktree, and team artifacts while retaining existing gate/filing/terminal outcome suites.
- DONE: AC-4 local gate proof.
  Applicable focused, full, race, and exact-head local Codex gates are green; local Claude is not required.
- DONE: Release-scope triage of Roborev job `1469` (verdict `F`).
  Its two Mediums request a genuine interactive transport or host command-stream classifier, both excluded by the revised mechanism boundary. No current AC or safety/compatibility boundary promises that instrumentation; the exact Codex spot-check performed only identify, while the same-model Claude ledger catches costly transient engage. Treat this as a deferred risk and revisit if the corrected ledger is green while archived artifacts show engage/dispatch before greet.

### Summary

Cycle 4 restores a controlled mutation-free shallow-boot workload and its durable guards, keeps merged-PR advancement at engage, and passes every required local gate at `ea63bcd4`. Validation recommends PASSED for updating PR #512, not for merge: the corrected same-resolved-model Sonnet ledger remains the hard release proof.
