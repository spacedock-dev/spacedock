---
id: wg2r4fwrdrv82hegd8x4rpqc
title: Sweep simplify the FO + ensign operating contracts (Phase 0.A of binary-simplification-roadmap)
status: ideation
source: "captain (2026-06-02) — binary-simplification-roadmap Phase 0.A; the qs cycle-3 reconcile prose simplification (PR #273) was the proof-of-concept (55% reduction with load-bearing meaning preserved). Apply the same pattern across the operating-contract surface that every ensign dispatch re-reads."
score: "0.45"
worktree:
started: 2026-06-03T00:51:44Z
completed:
verdict:
issue:
---

Generalize the qs cycle-3 prose-simplification pattern across the FO + ensign operating contracts that every dispatch re-loads. The qs cycle-3 sweep cut 880 → 400 words (55%) in the 4 reconcile sections while preserving all load-bearing mechanism; the same three rhetorical inflators (audit-trail exposition, cross-file restatement, over-qualification of slack) are visible elsewhere in the contract surface and should be pruned uniformly.

## Problem

Every ensign load re-pays the cost of the contract prose. With ~M ensign dispatches per session (this session: 15+), small contract-prose inefficiencies multiply faster than FO-side cognitive load (the FO loads ~1× per session). The qs cycle-3 simplification proved a 55% reduction is achievable WITHOUT semantic loss when the AI-engineer-pattern + oracle-tests are the verification gate. Other prose surfaces have accreted the same inflators and have not been swept.

## Proposed approach

Apply the qs-cycle-3 sweep pattern to four contract files:

| File | Current LOC | Owner |
|------|-------------|-------|
| `skills/first-officer/references/first-officer-shared-core.md` | 401 | FO universal |
| `skills/ensign/references/ensign-shared-core.md` | 157 | ensign universal |
| `skills/ensign/references/claude-ensign-runtime.md` | 46 | claude ensign |
| `skills/ensign/references/codex-ensign-runtime.md` | 69 | codex ensign |

**Total surface: 673 LOC of contract prose.** Conservative reduction target: 25–30% (~170 LOC removed from every ensign dispatch's read).

Sweep mechanism (per the qs cycle-3 dispatch):
1. AI engineer agent (general-purpose, senior-engineer persona) reads each file's diff context + adjacent prose, proposes verbatim-replacement rewrites per the three-inflator framework.
2. Captain reviews proposed rewrites; approves verbatim or names exceptions.
3. Implementation ensign applies APPROVED rewrites verbatim — no creative editing, no "improvement."
4. Validation ensign runs the oracle test set LOCALLY (`go test ./skills/integration/... ./internal/hostneutrality/...`) AND full repo `go test ./...`, confirms each oracle still PASSES, verifies word-count delta matches the proposal.

## Riskiest mechanism — exercised at ideation

The qs cycle-3 sweep IS the spike. It exercised every step of the mechanism (AI engineer dispatched, captain reviewed, ensign applied verbatim, oracles validated) at smaller scale (4 sections vs this entity's 4 files) and landed cleanly. The mechanism holds; this entity scales it.

`no spike needed` for the broader sweep: the proven mechanisms relied on are (a) `host-neutrality oracle` (`internal/hostneutrality/prose_neutrality_test.go`) — pins Claude-only tokens out of shared-core, (b) `portability oracle` (`skills/integration/portability_test.go`) — pins `~/.claude` paths into Claude adapter only, (c) `split-root-sync contract oracles` (`internal/hostneutrality/split_root_sync_contract_test.go`) — pins B3/B5/B6 load-bearing clauses, (d) `working-principles oracle` — pins insider-jargon-free + Working Principles section.

## Acceptance criteria

**AC-1 — Net word-count reduction of ≥25% across the four files, summed.**
Verified by: `wc -w` on each file pre- and post-sweep at the worktree HEAD; the validation stage report records before/after for each file + the percentage delta. AC fails if sum-delta < 25%.

**AC-2 — Every existing oracle test stays green at the swept HEAD.**
Verified by: `go test ./skills/integration/... ./internal/hostneutrality/... -v` runs with all of the ~25 contract-prose oracles PASSING (named in the test plan). AC fails if any oracle that passed pre-sweep fails post-sweep.

**AC-3 — Full repo `go test ./...` PASSES at the swept HEAD.**
Verified by: `go test ./...` output shows N/N PASS across all packages. AC fails on ANY failure. This is the broadest guarantee — no test outside the contract-surface oracles regresses.

**AC-4 — Two new lock-in oracles are added to permanently catch re-introduction of the swept inflators.**
Verified by: a new `TestNoAuditTrailExposition` (in `internal/hostneutrality/prose_neutrality_test.go` or a sibling file) that bans phrases like "audit-trail", "cycle-{N}", "the {entity} audit said", "now do X because of Y", or other history-as-meta-comment patterns from shared-core + both runtime adapters; AND a new `TestNoCrossFileRestatement` that flags n-gram overlap above a threshold (e.g., any 12-word phrase shared between shared-core and the runtime adapters, excluding intentional contrast spans). Both tests must PASS at the swept HEAD AND must fail when fed a deliberately-inserted regression (positive proof of the lock).

## Test plan

The validation stage runs these test sets LOCALLY (not just CI; the captain's directive — local-test-in-validation):

| Suite | Time | Asserts |
|-------|------|---------|
| `go test ./internal/hostneutrality/... -v` | <2s | `TestSharedCoreHasNoUnqualifiedClaudeHelpers`, `TestClaudeAdapterOwnsRelocatedCommands`, `TestSpanHostQualifiedRequiresContrast`, `TestFOHaltGateProse`, `TestFOSyncProse`, `TestEnsignSyncProse`, `TestCommissionJourneyProse`, `TestNoClaudeHomeReadsInGenericPackages`, `TestCodexRuntimeAdaptersAreLoadable`, `TestCodexAwaitingCompletionPinsMailboxSemantics` + the new AC-4 oracles |
| `go test ./skills/integration/... -v` | <5s | `TestShippedSurfaceHasNoHiddenMachineDependency`, `TestPortabilityCheckDiscriminatesHostSpecific`, `TestShippedInstructionsCarryNoInsiderJargon`, `TestFOContractCarriesWorkingPrinciplesSection`, `TestShipLocalCeremonyBlockExists`, `TestNoPluginStatusPathInVendoredSkills`, `TestNoPRMergeOrModBehaviorIntroduced`, `TestStartupGateGuidanceHasSingleProseSource`, `TestCommissionStateBackendDecisionRule`, `TestNoPluginPrivateStatusPathInContracts`, etc. |
| `go test ./...` | <60s | full repo |
| AC-4 negative test | <1s | the AC-4 lock-in oracles must FAIL on a deliberately-reverted regression (one of the swept files temporarily re-wraps a phrase) — proves the lock catches backsliding |
| Word-count delta | <1s | `wc -w` pre/post; sum-delta ≥25% |

Cost: low. All in-process Go tests. No live host, no network. Expected single implementation cycle.

## Out of scope

- The four files qs cycle-3 already swept (`first-officer-shared-core.md` step 10/supersede + claude/codex `first-officer-runtime.md` reconcile sections). Those are already simplified.
- FO + ensign SKILL.md files (under 30 lines each; low signal).
- Any commands moved into the binary (`state sync`, `dispatch advance`, `pr complete`, etc.) — those are roadmap Phase 1+ work, not Phase 0.A.
- Adding a portability test new-rules surface — the existing portability + host-neutrality oracles are sufficient for this sweep.
- ANY semantic change beyond pure rephrasing. If the AI engineer or implementer believes a clause needs SEMANTIC adjustment (not just word reduction), they MUST flag it and ask the captain — do NOT silently semantic-edit.

## Notes

- The qs cycle-3 pattern is the proof; this entity scales it.
- Per the captain's prioritization framing in the roadmap doc: ensign + shared contract has HIGHER per-load impact than FO event-loop changes because every dispatch re-pays the cost.
- AC-4 (the lock-in oracles) is the durable artifact — the prose simplification gains are protected against future backsliding for free.
- Implementation worktree branch: `spacedock-ensign/shared-contract-prose-sweep`.
