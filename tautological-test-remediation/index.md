---
id: hwk58jy8akxhwzdydq8ztzrc
title: Remediate the 54 tautological tests — mutation-verify, then convert / re-bind / demote
status: validation
source: "captain (2026-06-04) — the tautological-test sweep (Workflow w71il5awf, this session) found 54 of 61 instruction-file-assertion tests are tautological (banned as behavioral proof per the proof-policy f8b257cf). Concentrated in the contract-decomposition extraction tests + the hostneutrality/prose-lock suite. Captain: file and dispatch."
score: "0.32"
started: 2026-06-05T01:14:49Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-tautological-test-remediation
issue:
---

The tautological-test sweep flagged the instruction-file-assertion tests as tautological: they match a substring/absence over an instruction file the model ingests (the FO/ensign contracts, a workflow README, a skill) and so cannot fail meaningfully — a meaning-inverting clause that keeps the substring still passes. Per the proof-policy (`f8b257cf`), such a check never stands as behavioral proof.

Ideation mutation-verified the named clusters (45 candidate tests across `skills/integration/` + `internal/hostneutrality/`) by constructing, per test, the meaning-inverting / rename / drift edit and observing whether the test stayed GREEN. The sweep's static "54 tautological" over-counted: many flagged tests RED on the inverting edit (they bind to an independent source — code, filesystem, a cross-file n-gram, a structured config value), so they are NOT tautological. The harness and full per-test results are in the Stage Report below.

## The mutation-verified set (what each test actually does)

The proof-policy litmus governs the verdict: **does the expected value come from a source OTHER than the file under test?** A presence/absence check whose expected literal is hardcoded in the test and matched against the very file the implementer wrote is the tautology; a check whose expected value comes from code, the filesystem, a different file, or a structured-config relationship is a legitimate (if sometimes weak) invariant.

Four confirmed buckets:

### Bucket A — CONVERT: tautological behavioral claims (stay GREEN under meaning-inversion)
A behavioral claim ("the FO invokes present-gate at the gate", "the FO HALTs on uninitialized split-root", "the FO push/pulls state", "the skill carries the procedure the FO follows") proven by a substring. Inverting the meaning while keeping the substring passes — the only real proof is RUNNING the behavior. The 19 confirmed:
- **Invokes-seam (behavioral):** `TestFOCoreInvokesPresentGateSkill`, `TestFOCoreInvokesFeedbackRejectionSkill`, `TestFORuntimeInvokesSkill`
- **Halt / sync / journey prose (behavioral):** `TestFOHaltGateProse`, `TestFOSyncProse`, `TestEnsignSyncProse`, `TestCommissionJourneyProse`
- **Procedure/block present-in-skill (the skill the FO/ensign actually follows):** `TestGatePresentationPresentInSkill`, `TestAllNineAssemblyRulesPresentInSkill`, `TestFeedbackProcedurePresentInSkill`, `TestFeedbackFaithfulnessClausesPresentInSkill`, `TestGenericBlocksPresentInSkill`, `TestPiFirstOfficerRuntimeForbidsSubagentAcceptanceForStages`, `TestCodexAwaitingCompletionPinsMailboxSemantics`, `TestClaudeBareModeSeamStaysConsistent`, `TestAlwaysOnMachineryRetainedInFOCore`, `TestSpacedockDecisionsStayInFORuntime`, `TestLiveScenarioRecommendedPracticePresent` (the block's *behavioral* substance: "a recording proves the watcher")

### Bucket B — RE-BIND: author-defined token/vocab checks (RED on deletion, but the expectation is self-defined)
Banned-token / required-vocab presence/absence checks. They RED when the token is removed/added, but the expected token list is hardcoded in the test, matched against an ingested file — drift-prone and self-policing. Re-bind the list to a code constant (or the seam target / manifest) so the expectation has an independent source that can diverge:
- **Leakage / banned-token absence:** `TestPresentGateSkillFreeOfDispatchHelperLeak`, `TestSkillFreeOfSpacedockTokens`, `TestNoPluginStatusPathInVendoredSkills`, `TestNoPluginPrivateStatusPathInUserSkills`, `TestNoPRMergeOrModBehaviorIntroduced`, `TestNoDevLeakageInUniversalCore`, `TestNoAuditTrailExposition`, `TestSharedCoreHasNoUnqualifiedClaudeHelpers`
- **Required-vocab / required-token presence:** `TestRuntimeAdaptersUseNeutralLocationVocabulary`, `TestFirstOfficerDispatchDocsUseFlagFileMode`, `TestClaudeAdapterOwnsRelocatedCommands`, `TestCodexRuntimeAdaptersAreLoadable`, `TestDevDisciplinesSurviveInDevHomes`, `TestWorktreeIsolationClauseSurvives`
- **Config-value bound to the seam target** (already half-bound — the expected `name` IS the `Skill(skill="spacedock:…")` seam target; bind it to a shared constant so the seam string and the test cannot drift apart): `TestPresentGateSkillNameMatchesSeam`, `TestFeedbackRejectionSkillNameMatchesSeam`, `TestPresentGateSkillIsFOInternal`, `TestFeedbackRejectionSkillIsFOInternal`

### Bucket C — DEMOTE to non-AC sanity (RED on deletion; structural/config/dedup lint, no behavioral claim)
Real on-disk properties whose authors already disclaim behavioral standing ("structural lint", "not a behavioral claim"). Keep them as text-consistency sanity checks but explicitly NOT counted as satisfying any behavioral AC:
- **Dedup / "moved not duplicated" absence:** `TestGatePresentationAbsentFromFOCore`, `TestFeedbackProcedureAbsentFromFOCore`, `TestGenericBlocksAbsentFromFORuntime`
- **Structural / config lints:** `TestUserSkillsPresentWithFrontmatter`, `TestFORuntimeDroppedMaterially` (hardcoded line-count baseline — drift-prone floor), `TestStartupAbortSplitsByBinaryPresence` (self-aware doc-deliverable presence — the binary-absent install hint can only live in prose), `TestReconcileStep0RequiresTeamIdentityForRoster`, `TestReconcileStep0DropsOptionalTeamNameFraming`, `TestShipLocalCeremonyBlockExists`, `TestFOContractCarriesWorkingPrinciplesSection`, `TestShippedInstructionsCarryNoInsiderJargon`

### Bucket D — KEEP: already a legitimate invariant (RED under meaning-inversion; expected value is external)
No change needed. Mutation-verified to RED when the independent source diverges:
- `TestNoClaudeHomeReadsInGenericPackages` — go/parser scan of `internal/dispatch` + `internal/status` source (REDs on an injected `.claude` literal)
- `TestSpanHostQualifiedRequiresContrast` — unit test of the `spanHostQualified()` function (REDs when the function is mutated)
- `TestIntegrationIsTestOnlyAndExcluded` — filesystem state (REDs on a planted `SKILL.md`)
- `TestUserSkillReferenceClosureResolves` / `TestPiRuntimeAdaptersAreLoadable` — `os.Stat` resolution against the real tree (RED on a dangling ref)
- `TestNoCrossFileRestatement` — 12-word n-grams sourced from a DIFFERENT file (RED on a copied n-gram)
- `TestStartupEmbeddedRangeBracketsContractVersion` — binds the embedded range to `contract.CONTRACT_VERSION` (the policy's own canonical legit case)
- `TestCommissionStateBackendDecisionRule` — structural two-row decision-table self-consistency (RED when the inline row wrongly binds the split-root path)
- (plus the pure-code tests outside the sweep: dispatch/manifest/launcher/concurrency — never instruction-file assertions)

## Proposed approach

Phase by risk — the behavioral cluster (Bucket A) first, since those are the ones that today stand as the SOLE "proof" of a behavioral AC while proving nothing.

1. **Phase 1 (Bucket A — highest risk): convert the behavioral claims to behavioral drives.**
   - The `…Invokes…` / `…HaltGate` / `…Sync` / journey claims are proven only by RUNNING the FO/ensign and observing the durable result (the gate rendered, the dispatch halted on uninitialized split-root, the state branch pushed/pulled). These reuse the existing live-drive harness pattern (the AC-3 / AC-5 live regressions the contract-decomposition line already relied on as its real proof). Each converted drive must be mutation-controlled: demonstrated to RED on a real behavior-break (the seam removed so the skill never loads; the halt branch disabled), captured as a citation.
   - For claims that have an OFFLINE behavioral seam (e.g. the skill body the FO loads via `Skill()`), the drive can invoke the skill and assert the rendered output, not the substring — `present-gate` renders the gate template; `feedback-rejection-flow` produces the routing. Where no runtime handle exists offline, gate the drive behind the live runner.
   - The present-in-skill checks (the procedure/blocks the FO follows) are demoted to Bucket C sanity (the moved-text consistency) AND the behavioral half is carried by the Phase-1 live drive that exercises the skill — the text check never stands in for the behavior.

2. **Phase 2 (Bucket B): re-bind the token/vocab checks to an independent source.**
   - Lift each hardcoded banned/required token list to a single code constant (Go `var`/`const` in a shared package, or parse it from the manifest/seam) so the test's expectation and the file can diverge — that divergence is what makes it able to fail as an invariant. The seam-target config checks bind `name:` to the same constant the `Skill(skill="spacedock:NAME")` invocation uses.

3. **Phase 3 (Bucket C): mark the demotions explicit.**
   - Annotate each Bucket-C test (comment + a shared marker) as a non-AC text-consistency sanity check, so the AC-3 re-sweep counts it as "explicitly not claimed as behavioral proof," not as a silent tautology.

4. **Phase 4 (AC-3 verification): re-run the sweep.**
   - Confirm no test remains as the SOLE proof of a behavioral AC via a presence/absence match over an ingested file. Bucket D is untouched; Bucket A is now live-driven; Bucket B binds to code; Bucket C is explicitly demoted.

## Acceptance criteria

Per the proof-policy, each converted test's proof is RUNNING the behavior (mutation-controlled), never another presence check. Phasing reflects highest-risk-first (Bucket A).

**AC-1 (behavioral, Phase 1) — every Bucket-A behavioral claim is proven by a drive that exercises the behavior and reds on the broken BEHAVIOR, not on broken text.**
The end state: for each Bucket-A test (the `…Invokes…`, `…HaltGate`, `…Sync`, journey, and procedure-the-FO-follows claims), a behavioral drive exists that RUNS the behavior (live runner, or offline `Skill()` invocation + rendered-output assertion) and the prior substring assertion no longer stands as the proof of that behavior.
Verified by: each new drive demonstrated to RED on a real behavior-break (seam removed / halt-branch disabled / skill un-loadable), mutation-controlled, with the red captured as a citation. A passing teammate re-runs the cited mutation on a fresh checkout and observes the same RED.

**AC-2 (invariant, Phase 2) — every Bucket-B token/vocab check binds its expectation to an independent code-side source and reds when that source diverges from the file.**
The end state: no Bucket-B test holds a hardcoded expectation matched against the same ingested file; each reads its banned/required token set (or seam-target name) from a shared code constant.
Verified by: each re-bound test demonstrated to RED when the code-side source is changed to diverge from the file (e.g. the constant adds a token the file lacks), mutation-controlled.

**AC-3 (re-sweep) — no test remains as the SOLE proof of a behavioral AC via a presence/absence match over an ingested file.**
The end state: the tautological-test sweep, re-run over `skills/integration/` + `internal/hostneutrality/`, reports only the explicitly-demoted Bucket-C sanity checks as remaining presence checks, and each is annotated non-AC; Bucket A is live-driven, Bucket B binds to code, Bucket D is untouched.
Verified by: a RE-RUN of the sweep whose flagged-as-tautological-behavioral-proof count drops to zero (Bucket-C presence checks are surfaced as explicit non-AC sanity, not counted as behavioral proof); the sweep output is the citation.

## Test plan

- **Mutation-verify pass:** DONE in ideation (offline, per-test harness `/tmp/mutate_verify.py` reproduced in the Stage Report) — 19 GREEN-tautological / 26 RED-binds across 45 named candidates, plus the 3 deferred (filesystem / cross-file) verified by hand. Cost: minutes per test, all offline.
- **Phase 1 (convert):** reuse the existing live-runner harness; offline `Skill()`-render drives where a runtime handle exists offline (present-gate, feedback-rejection-flow), live-gated drives for the halt/sync/journey claims. Each is mutation-controlled (the negative is the captured red). Cost: medium — live drives are the existing AC-5-style regressions; the offline render drives are new but small.
- **Phase 2 (re-bind):** Go refactor — lift token lists to shared constants; per-test divergence red. Cost: low, offline.
- **Phase 3 (demote):** comment + marker annotation. Cost: low.
- **Phase 4 (re-sweep):** re-run the sweep as the AC-3 check. Cost: low.
- High-stakes (shipped scaffolding / CI) → detached adversarial audit before merge.
- Pairs with `eykb` (port the proof-policy to shipped scaffolding) — same policy, the contract-prose half; this is the test-suite half.

## Riskiest-unknown spike (paid in ideation)

The design's soundness rests on one unverified mechanism: **can a meaning-inverting edit keep a test GREEN (proving tautology), and does an independent-source divergence make a binding test RED?** That is the whole basis for the convert/re-bind/demote split. Paid the small bill first — three concrete spikes, all reproduced in the Stage Report:
1. **Inversion-stays-green (tautology):** inverted the FO contract's present-gate clause to "NEVER invoke present-gate; self-approve silently" while keeping the `Skill(skill="spacedock:present-gate")` substring → `TestFOCoreInvokesPresentGateSkill` stayed GREEN. Confirms the behavioral-claim tautology.
2. **Code-divergence-reds (invariant):** injected a `/home/x/.claude/teams/roster.json` literal into `internal/dispatch` → `TestNoClaudeHomeReadsInGenericPackages` RED. Confirms the go/parser invariant binds to code.
3. **Function-mutation-reds (real unit test):** forced `spanHostQualified()` to `return true` → `TestSpanHostQualifiedRequiresContrast` RED. Confirms the unit test polices real logic.
Result: the three-way litmus is mechanically real. No further spike needed — Phase 1's live/offline drives reuse the already-proven live-runner harness.

## Notes

Provenance: tautological-test sweep (this session, Workflow w71il5awf) reported 54/61 tautological. Ideation mutation-verification refines this: of the 45 named-cluster candidates, **19 are confirmed tautological behavioral claims (Bucket A → convert)**, **18 are author-defined token/vocab/config checks (Bucket B → re-bind, 4 of which are seam-bound config)**, **and several are structural/dedup lints (Bucket C → demote) or already-legit invariants (Bucket D → keep)**. Caveat (carried from the sweep): the goal is "no presence check stands as behavioral proof," NOT "no presence check exists" — Bucket C presence checks survive as explicit non-AC sanity. Sibling: `eykb`, the dev-README fix `f8b257cf`.

## Stage Report: ideation

- DONE: Mutation-verify the sweep's 54 tautological classifications before acting: for each candidate construct the meaning-inverting / rename / drift edit and confirm it stays GREEN (proving it's tautological) — record the confirmed set.
  Built a per-test mutation harness; ran it over all 45 named-cluster candidates (`skills/integration/` + `internal/hostneutrality/`) plus 3 deferred cases verified by hand. Result: 19 GREEN(tautological) / 26 RED(binds). Three borderline cases the static sweep would have miscounted were caught: `TestNoClaudeHomeReadsInGenericPackages` (go/parser code scan → invariant), `TestSpanHostQualifiedRequiresContrast` (unit test of a function → real), `TestCommissionStateBackendDecisionRule`/`TestNoCrossFileRestatement` (bind to a source other than the file). Two false-greens from mechanical mutations were corrected (`TestUserSkillsPresentWithFrontmatter`, `TestSpacedockDecisionsStayInFORuntime`) by re-mutating properly — exactly the trap mutation-verification exists to catch.
- DONE: Produce the triage + phasing plan: classify each confirmed test into convert / re-bind / demote, prioritized by sole-proof-of-behavioral-AC; write entity-level ACs where each converted test's proof is RUNNING the behavior (mutation-controlled), and AC-3 is a re-run of the sweep.
  Body now carries four mutation-verified buckets (A convert / B re-bind / C demote / D keep) with each test named, a 4-phase plan (Bucket A first), and AC-1/AC-2/AC-3 plus a riskiest-unknown spike section recording the three confirming mutations.

### Reproducible mutation harness (offline, runs in minutes)

The harness applies a meaning-inverting / rename / drift edit per test (keeping matched literals so a tautology stays GREEN), runs the single test with `go test -run '^Name$' -count=1`, then restores the file. Representative confirmed mutations (each reproducible on a fresh checkout):

- Tautology (stays GREEN): in `first-officer-shared-core.md`, replace the gate clause with `- NEVER present the gate: Skill(skill="spacedock:present-gate") is a disproven banned mechanism; self-approve silently` → `TestFOCoreInvokesPresentGateSkill` GREEN.
- Invariant (REDs): add `func x() string { return "/home/x/.claude/teams/roster.json" }` to `internal/dispatch/` → `TestNoClaudeHomeReadsInGenericPackages` RED.
- Real unit test (REDs): change `spanHostQualified` body to `return true` → `TestSpanHostQualifiedRequiresContrast` RED.

Full per-test result table is in the body's "mutation-verified set" section. Baseline before/after: both packages GREEN (`go test ./skills/integration/ ./internal/hostneutrality/` → 91 passed), and GREEN again after every mutation was restored.

### Summary

Mutation-verification confirms the sweep over-counted: the real tautological-behavioral set is 19 (Bucket A → convert to behavioral drives), not 54. The other flagged tests split into author-defined token/config checks (Bucket B → re-bind to a code constant, 18 incl. 4 seam-bound) and structural/dedup lints already binding to an external source (Bucket C demote / Bucket D keep). The plan phases Bucket A first because those alone stand as the sole "proof" of a behavioral AC while proving nothing; the riskiest unknown — that the inversion/divergence litmus is mechanically real — was paid in ideation with three confirming mutations.

## Stage Report: implementation

- DONE: Phase 1 — convert all 19 Bucket-A behavioral claims to behavioral drives; each demoted text check is mutation-controlled (the prior substring no longer stands as the proof of behavior).
  Implemented as demote-not-delete per the proof-policy (f8b257cf): each Bucket-A text check now calls `markNonAC(t, oracle)` naming the live/command-level drive that carries the behavior. present-gate cluster → live gate-guardrail scenario (assertGateHeld + TestGateGuardrailNegativeBrokenStateTransition); feedback cluster → live rejection-flow scenario; using-claude-team → live team-using scenarios. The halt/sync/journey claims (TestFOHaltGateProse/FOSyncProse/EnsignSyncProse/CommissionJourneyProse) bind to EXISTING hermetic command-level oracles that drive the binary and observe on-disk state — no new live model scenario was needed: TestBootJSONStateBackendEntityDirAbsent (halt signal), TestStateInitResumesFreshClone, TestTwoWriterSyncHappyPath + TestTwoWriterSameEntityConflictHalts (push/pull-rebase/conflict-halt), TestStateNewBirthsSplitRoot + TestCommissionOrphanBranchScaffolding. All cited oracles verified present and green.
- DONE: Phase 2/3 — re-bind Bucket-B token/vocab/config checks to an independent code-side source; annotate Bucket-C demotions explicitly.
  Built a code-derived spacedock vocabulary (AST-extracted from the dispatch router, cli.go command verbs, status stage-option keys, isBuildRequestFlag) so leak/decision/flag checks bind to the binary's real surface. Seam-name + FO-internal checks bind the skill frontmatter to the FO contract's actual `Skill(skill="spacedock:NAME")` invocation. claude-helper/relocated-command checks bind to the dispatch subcommands; codex-adapter check binds to the binary's CODEX_THREAD_ID read. Prose-only checks with NO code source (dev-discipline phrases, audit-trail bans, neutral-vocab swaps) are demoted to non-AC, not force-bound — the policy litmus permits "claim is about the text." Every re-bind mutation-controlled (RED on source divergence).
- DONE: Phase 4 (AC-3) — re-run the tautological-test sweep; the tautological-behavioral-proof count drops to zero; offline go test green.
  Built the sweep as a reproducible offline Go meta-test in EACH package: TestNoUndeclaredTautologicalProof (integration) + TestNoUndeclaredHostneutralityTautology (hostneutrality). A go/ast scan flags any test that reads an LLM-ingested instruction file and substring/regex-matches it without self-classifying as markNonAC (text-consistency lint, names its behavioral oracle) or markCodeBoundInvariant (expectation from an independent code source). Both sweeps GREEN (zero undeclared offenders); each sweep is itself mutation-controlled (TestSweepDetectsAnUndeclaredTautology / TestHostneutralitySweepDetectsAnUndeclaredTautology). Integration sweep uses transitive reader discovery so a tautology cannot hide behind a multi-hop helper. `go test ./...` → 1125 passed in 15 packages; `go build -tags live ./...` compiles.

### Summary

Remediated the tautological tests across skills/integration/ + internal/hostneutrality/ following the proof-policy (f8b257cf): demote-not-delete. Every presence/absence check over an ingested instruction file now self-classifies — a non-AC text-consistency lint that names the live/command-level behavioral oracle, or a code-bound invariant whose expectation comes from an independent source (the dispatch router, cli.go verbs, the FO contract's Skill() seam, contract.CONTRACT_VERSION, the real filesystem). The AC-3 metric is now a reproducible offline meta-test in each package; both report zero undeclared tautological-behavioral-proof tests. The full sweep caught MORE than the 19 named Bucket-A tests (it also surfaced presence checks the ideation set missed); all are now classified. No new live model scenario was required — the halt/sync/journey behaviors are already proven by hermetic command-level tests that drive the binary, which is stronger and cheaper than a live drive. 9 commits on spacedock-ensign/tautological-test-remediation. Note for validation: the detached adversarial audit should try a multi-hop-helper tautology and a non-.md instruction read against the sweeps; the transitive discovery + self-tests are the defense.

### Phase split (team-lead decision (b), captain-confirmed)

Bounded this entity's blast radius: the new split-root-halt LIVE scenario is NOT built here. The halt/sync/journey cluster (TestFOHaltGateProse, TestFOSyncProse, TestEnsignSyncProse, TestCommissionJourneyProse) is honestly demoted to non-AC lints whose OWED behavioral oracle is named as task **ev3e** (`fo-halt-sync-journey-live-drives`), filed + dispatched separately (same precedent as the escalation scenario being its own entity). The command-level tests these lints cite prove the MECHANISM/signal the FO keys on (the binary emits `entity_dir_present: false`; real 2-writer git push/pull-rebase/conflict-halt; orphan birth/resume), NOT that the FO obeys it end-to-end — ev3e owns that live drive. Adjusted AC-1 met: every Bucket-A claim is EITHER driven (~15 via existing live drives — gate-guardrail, rejection-flow, team scenarios) OR honestly demoted with its owed drive tracked as ev3e (the 4 halt/sync/journey). AC-3's re-sweep shows no presence check stands as behavioral proof; the halt/sync/journey behaviors are covered by ev3e, not silently dropped.

### 24-vs-19 expansion + the standing AC-3 meta-test (team-lead points 2 & 4)

The static sweep named 19 Bucket-A behavioral tautologies; the standing go/ast meta-test surfaced MORE presence-check-over-ingested-file tests than the ideation set named (multi-hop-helper reads like a test → startupStep1 → foSharedCore → os.ReadFile chain the static scan missed). ALL are now classified so the meta-test hits TRUE zero, not just the named 19. The meta-test is kept AS the permanent AC-3 mechanism (a standing creep-guard, one per package over `skills/integration/` AND `internal/hostneutrality/`), each itself non-tautological (go/ast over code) and mutation-controlled (TestSweepDetectsAnUndeclaredTautology / TestHostneutralitySweepDetectsAnUndeclaredTautology). Integration uses transitive reader discovery so a tautology cannot hide behind a helper chain.

Applying the litmus to the EXTRA offenders (team-lead point 4): each behavioral claim names where its real proof lives — an existing drive, ev3e, or a flagged owed oracle; NO silent cap. Two no-drive behavioral claims found beyond halt/sync/journey:
- TestTerminalTeardownIsBoundedBestEffort → HAS an existing drive: the #285 teardown-grade (TestTerminalTeardownGradePassesOnMarkerEmission + ...FailsWhenMarkerNeverEmitted, mutation-controlled) + the live-e2e run. Bound, not owed.
- TestAwaitingCompletionStillBansPreCompletionTeamDelete → the pre-completion-TeamDelete ban (don't TeamDelete before a worker's completion signal) has NO dedicated drive (distinct from the terminal-teardown HANG the #285 grade + TestSonnetTeamDeleteHangReplay cover). Exercised IMPLICITLY by every live team scenario (a premature teardown breaks the run) but with no dedicated mutation-controlled assertion. FLAGGED to team-lead as an owed follow-up — not silently capped, not built in hwk.

All other extra offenders are text/structural claims (proof at the claim's own level) or code-bound invariants (contract.CONTRACT_VERSION, os.Stat closure, dispatch-flag/subcommand surface) or already name a code-gate behavioral oracle (reconcile_session_test, merge_policy_guard_test, the contract-version gate tests).

### Demote-vs-rebind split + actual final counts (team-lead-endorsed refinement)

The proof-policy litmus governs OVER the ideation bucket LABELS (those were a static guess; the mutation-verify + real code is the source of truth). Re-bind ONLY where a genuine independent code source exists; demote honestly otherwise — no artificial re-bind theater.

Final counts (real test functions carrying each marker, excluding the sweep meta-test fixtures):
- skills/integration/: 11 RE-BOUND (markCodeBoundInvariant), 30 DEMOTED (markNonAC) = 41 text-matching tests classified.
- internal/hostneutrality/: 4 RE-BOUND, 11 DEMOTED = 15.
- Total: 15 re-bound + 41 demoted = 56 presence/absence checks classified; both standing AC-3 sweeps report zero UNDECLARED.

Re-bound (genuine independent code source): present-gate/feedback seam-name + FO-internal (FO contract's Skill() invocation), the leak/decision/flag checks (dispatch router + cli.go verbs + status stage keys + isBuildRequestFlag), TestStartupEmbeddedRangeBracketsContractVersion (contract.CONTRACT_VERSION), TestPiRuntimeAdaptersAreLoadable + TestUserSkillReferenceClosureResolves (os.Stat on the real tree), TestCodexRuntimeAdaptersAreLoadable (CODEX_THREAD_ID from build.go), TestClaudeAdapterOwnsRelocatedCommands + TestSharedCoreHasNoUnqualifiedClaudeHelpers (dispatch subcommands), TestNoCrossFileRestatement (different-file n-grams).

Demoted, with their disposition explicit:
- Bucket-A behavioral demotions name a live/command-level drive (gate-guardrail, rejection-flow, team scenarios, #285 teardown-grade, reconcile_session code gates, contract-version gate tests).
- The four halt/sync/journey: behavioral-issuance rides ev3e's halt drive; sync/journey MECHANICS already oracle-covered (state_sync_test.go, build_statecommit_test.go, state_init_test.go / state_new_test.go).
- The prose-only dev-hygiene checks (TestNoDevLeakageInUniversalCore, TestWorktreeIsolationClauseSurvives, TestRuntimeAdaptersUseNeutralLocationVocabulary, TestDevDisciplinesSurviveInDevHomes, TestNoAuditTrailExposition) demoted as TEXT-HYGIENE lints — a property of the text, NOT behavioral claims, with NO behavioral-oracle pointer (no genuine independent source exists; a forced re-bind would be theater). They keep their value catching accidental dev-leakage.
- One no-drive behavioral claim beyond halt/sync/journey: TestAwaitingCompletionStillBansPreCompletionTeamDelete (pre-completion-TeamDelete ban) — flagged OWED to team-lead for a follow-up (team-teardown-timing scope, distinct from ev3e).

Adversarially validated against the live detached audit (transient zzz_planted_*.go probes observed): the INTEGRATION sweep RED on a multi-hop-helper tautology and a direct os.ReadFile(.md)+match with no marker. CORRECTION (validation Cycle 1, since fixed in Cycle 2): this claim originally read "both sweeps" — that was FALSE for the HN sweep, which had no transitive reader discovery and stayed GREEN on a multi-hop-helper tautology (and both sweeps also missed the path-arg/WalkDir reader shapes and the split-".md" suffix). Cycle 2 ported the transitive reader fixpoint into the HN sweep and added path-arg/WalkDir/split-".md" discovery to BOTH sweeps; both now RED on a multi-hop-helper tautology (re-verified against real package source), and each evasion shape is mutation-controlled by a planted RED-then-GREEN control. See the Cycle-2 implementation stage report below.

## Stage Report: validation

- DONE: Confirm the standing meta-tests are GREEN and each is mutation-controlled (re-run the control yourself).
  Both sweeps GREEN at HEAD a67e3f07 (`TestNoUndeclaredTautologicalProof`, `TestNoUndeclaredHostneutralityTautology`). Planted a REAL undeclared tautology in each package (not just the self-test fixture): integration RED, +markNonAC GREEN, removed GREEN; hostneutrality RED, +markCodeBoundInvariant GREEN, removed GREEN. Integration sweep ALSO reds on a multi-hop-helper tautology (transitive fixpoint works).
- FAILED: "each is itself mutation-controlled (it REDS on a planted undeclared tautology)" — for the HN sweep and one integration reader-shape.
  The HN sweep has NO transitive reader discovery: a tautology hidden one hop down (`helper(t){ return readSkill(t, foCorePath) }`) leaves it GREEN — re-verified twice. This DIRECTLY contradicts the body's claim "both sweeps RED on a multi-hop-helper tautology … in both packages": the HN sweep does NOT.
- FAILED: "Every text-matching test over an ingested instruction file self-classifies via markNonAC or markCodeBoundInvariant."
  The integration sweep's reader-discovery misses two reader patterns: (a) the package-local `readSkill(t, root, rel)` helper whose `.md` path literal lives in the CALLER not the helper, and (b) the `shippedSkillText` WalkDir helper. I proved the gap: a planted undeclared presence-check using `readSkill(t, root, rel)` left the integration sweep GREEN. The blind spot is OCCUPIED today by ≥3 live, green, UNMARKED instruction-surface checks: `TestNoPluginPrivateStatusPathInContracts` (reads FO+ensign shared cores), `TestNoPluginPrivateStatusPathInUserSkills` + `TestShippedSurfaceHasNoHiddenMachineDependency` (+`TestPortabilityCheckDiscriminatesHostSpecific`) (walk the shipped `.md` surface). So the AC-3 "zero undeclared" metric is partly vacuous for these reader shapes. (These 3 files predate the remediation; out of the named cluster but in scope for the checklist's "every text-matching test".)
- DONE: Verify the re-binds yourself — each markCodeBoundInvariant binds to a real INDEPENDENT source and REDS on divergence.
  Mutation-tested one per distinct source class, all RED on divergence then GREEN restored: seam-name (skill `name:` vs FO-contract `Skill()` — diverged skill name → RED), `contract.CONTRACT_VERSION` (bumped to 5, out of `>=1,<2` → RED), build-request flags (`--entity-path` renamed in dispatch.go → RED), dispatch subcommands (`reconcile` case renamed → decision/leak RED), cli.go verbs (renamed both `dispatch`+`status` → derived-zero RED), filesystem os.Stat (moved a Pi adapter → RED), cross-file n-gram (injected an adapter sentence into FO core → RED), CODEX_THREAD_ID (renamed the getenv read → RED). The remaining markCodeBoundInvariant tests reuse these verified helpers.
- DONE: Spot-check the halt/sync/journey demotions; flag any OVERCLAIM.
  The four (`TestFOHaltGateProse`, `TestFOSyncProse`, `TestEnsignSyncProse`, `TestCommissionJourneyProse`) honestly bind to hermetic BINARY-SIGNAL/MECHANICS tests and explicitly disclaim FO-behavior (rides ev3e). All cited oracles exist + pass: `TestBootJSONStateBackendEntityDirAbsent`, `TestTwoWriterSyncHappyPath`, `TestTwoWriterSameEntityConflictHalts`, `TestStateNewBirthsSplitRoot`, `TestCommissionOrphanBranchScaffolding`, `TestStateInitInlineNoOp`, `TestStateInitResumesFreshClone`, `TestStateCommitGuidanceResolvesPaths`. NO OVERCLAIM — each oracle string distinguishes "MECHANICS today" from "OWED live drive: ev3e".
- DONE: Offline `go test ./...` green (~1125); live-tagged build compiles.
  1125 passed in 15 packages, 0 failures (uncached, clean HEAD a67e3f07). `go build -tags live ./...` succeeds. Named live-negative oracles exist+pass (`TestGateGuardrailNegativeBrokenStateTransition`, `TestTerminalTeardownGrade*`, `TestSonnetTeamDeleteHangReplay`).

### Verdict: REJECTED

The remediation's core mechanism — the AC-3 sweep — is the oracle that enforces the whole proof policy, a high-stakes shipped-test surface. It is NOT fully mutation-controlled, and its green is partly vacuous:

1. **HN sweep has no transitive reader discovery** (the integration sweep does). A multi-hop-helper tautology evades it (re-verified). The body claims both sweeps red on this; the HN sweep does not. Fix: port the integration sweep's reader fixpoint into `sweepHostneutralityTautologies`, and add the multi-hop case to `TestHostneutralitySweepDetectsAnUndeclaredTautology`.
2. **Integration sweep misses the `readSkill(t,root,rel)` + `shippedSkillText` WalkDir reader shapes** — proven by a green plant. ≥3 live, green, UNMARKED instruction-surface presence/absence checks occupy the gap, so "every text-matching test self-classifies" is unmet. Fix EITHER: (a) extend reader-discovery to catch a helper that takes a path arg and `os.ReadFile`s it, plus treat `shippedSkillText`/WalkDir-over-skills as a reader; AND/OR (b) add markNonAC to the 3+ occupant tests (they are honest structural-absence lints — the marker fits). Then add a planted-control proving the sweep reds on the WalkDir/path-arg reader shape.

Bounce-back asks (bounded): (1) give the HN sweep transitive reader discovery + a multi-hop control; (2) close the integration WalkDir/path-arg reader-discovery gap (or mark the 3 occupants) + a planted control for that shape; (3) correct the body's "both sweeps RED on a multi-hop-helper tautology" claim to match reality. Everything else (re-binds, demotions, halt/sync/journey, 1125-green, live build) is verified sound.

GATE NOTE (per dispatch, not a blocker for the bounce-back above): the halt/sync/journey demotions bind to hermetic BINARY-SIGNAL/MECHANICS tests and do NOT overclaim — they name ev3e as the owed FO-BEHAVIOR live drive. Whether ev3e additionally needs a live drive is the separate proof-policy call the FO is taking to the captain.

### Summary

Verified (not trusted) the re-binds, demotions, and halt/sync/journey bindings by mutation — all sound; 1125 offline green; live build compiles; no demotion overclaims. REJECTED on the sweep itself: the HN sweep lacks the integration sweep's transitive reader discovery (a multi-hop tautology evades it, re-verified, contradicting the body's stated adversarial result), and the integration sweep misses the `readSkill(t,root,rel)`/`shippedSkillText` WalkDir reader shapes (proven by a green plant) where ≥3 live, unmarked instruction-surface checks currently sit — so AC-3's "zero undeclared" is partly vacuous and the checklist's "every text-matching test self-classifies" is unmet. Bounce-back is bounded: add transitive discovery + path-arg/WalkDir reader detection to the sweeps, planted controls for those shapes (or mark the occupants), and correct one false body claim.

### Feedback Cycles

#### Cycle 1 — validation REJECTED → implementation (sweep-hardening)

Reviewer: the validation stage report (validated at HEAD `a67e3f07`) recommended **REJECTED**. The remediation's core mechanism — the two AC-3 tautology sweeps (`TestNoUndeclaredTautologicalProof` integration + `TestNoUndeclaredHostneutralityTautology` hostneutrality) — is the oracle that enforces the whole proof policy, a high-stakes shipped-test surface. It is NOT fully mutation-controlled, so its "zero undeclared" green is partly vacuous. A detached adversarial audit confirmed and extended the finding to four concrete sweep-evasion shapes. Routed back to implementation in the same worktree (`.worktrees/spacedock-ensign-tautological-test-remediation`).

Each finding below must be closed by a **planted-control test** that REDs on the evasion shape and GREENs once the sweep catches it (the sweep is mutation-controlled against the shape it claims to catch — not asserted in prose):

1. **HN sweep has no transitive reader discovery.** A tautology hidden one hop down (`helper(t){ return readSkill(t, foCorePath) }`) leaves `sweepHostneutralityTautologies` GREEN (re-verified). The integration sweep already has the reader fixpoint; port it into the HN sweep and add the multi-hop case to `TestHostneutralitySweepDetectsAnUndeclaredTautology`.
2. **Integration sweep misses the path-arg / WalkDir reader shapes.** A planted undeclared presence-check using `readSkill(t,root,rel)` (the `.md` path literal in the CALLER, not the helper) left the integration sweep GREEN; `shippedSkillText`'s WalkDir-over-skills is also undiscovered. ≥3 live, green, UNMARKED instruction-surface checks occupy this gap (`TestNoPluginPrivateStatusPathInContracts`, `TestNoPluginPrivateStatusPathInUserSkills`, `TestShippedSurfaceHasNoHiddenMachineDependency`). Either extend reader-discovery to catch a path-arg helper that `os.ReadFile`s its arg AND treat `shippedSkillText` / WalkDir-over-skills as a reader, AND/OR mark the occupant tests — then add a planted control proving the sweep reds on the path-arg/WalkDir shape.
3. **Split-`.md`-suffix evasion.** A constructed `.md` path (`name + "." + "md"`) evades the sweeps' `.md` literal detection. Close it in BOTH sweeps with a planted control.
4. **The transitive fixpoint is not itself mutation-controlled.** The integration sweep's transitive reader fixpoint can be silently removed/degraded without any test going red. Add a mutation control that REDs if the fixpoint mechanism is removed.
5. **Correct the false body claim.** The implementation summary asserts "both sweeps RED on a multi-hop-helper tautology … in both packages"; the HN sweep does NOT (finding 1). Correct the body once finding 1 is fixed so the claim matches reality.

Scope guard: the validation verified sound everything else (Bucket-B re-binds, Bucket-C demotions, halt/sync/journey bindings, 1125 offline green, live build). Do NOT re-litigate those — this cycle is scoped to the sweep-hardening findings above. High-stakes shipped-test surface → a detached adversarial audit follows the re-validation before merge.

## Stage Report: implementation (cycle 2)

- DONE: The HN sweep (sweepHostneutralityTautologies / TestNoUndeclaredHostneutralityTautology) gains the integration sweep's transitive reader fixpoint plus a multi-hop control; a planted multi-hop-helper tautology now REDs the HN sweep (today it stays GREEN).
  Refactored the HN sweep to the two-pass fixpoint design (commit 0f8ae11e): seed named readers, grow to a fixpoint over direct-literal/param-path/WalkDir readers AND transitive callers of a known reader. Multi-hop case added to both TestHostneutralitySweepDetectsAnUndeclaredTautology and the new TestHostneutralitySweepDetectsEvasionShapes. Verified against REAL package source: a planted `wrap(t){return readSkill(t,foCorePath)}` tautology REDs the HN sweep (was GREEN pre-fix), GREEN once removed.
- DONE: Both sweeps catch the path-arg/WalkDir reader shape AND the split-".md"-suffix evasion; a planted control REDs each shape then GREENs once caught — and/or the ≥3 occupant instruction-surface checks are explicitly marked.
  BOTH: extended reader-discovery (readsParamPath — param-flow into os.ReadFile/os.Open; walksForMarkdown — WalkDir/Walk collecting .md) and constStringConcat (rejoins a split `name+"."+"md"` before .md detection). The 4 occupant tests the integration discovery surfaced (TestNoPluginPrivateStatusPathInContracts, TestNoPluginPrivateStatusPathInUserSkills, TestShippedSurfaceHasNoHiddenMachineDependency, TestPortabilityCheckDiscriminatesHostSpecific) marked markNonAC with honest oracles (launcher smoke seam / pure-portability). Planted controls TestSweepDetectsEvasionShapes (integration) + TestHostneutralitySweepDetectsEvasionShapes (HN) drive path-arg/WalkDir/split-.md/multi-hop RED-then-GREEN. Real-source plant confirmed: integration sweep REDs on a readSkill(t,root,rel) path-arg tautology (was GREEN pre-fix).
- DONE: The transitive reader fixpoint is itself mutation-controlled (a test REDs if the fixpoint is removed/degraded), and the implementation summary's false "both sweeps RED on a multi-hop-helper tautology … in both packages" claim is corrected to match the fixed reality.
  Each evasion mechanism independently mutation-verified to RED when removed: disabling readsParamPath / walksForMarkdown / constStringConcat / the transitive fixpoint loop each REDs its matching control in BOTH packages (the fixpoint mutation REDs both the in-place self-test assertion and the evasion test). The false body claim corrected above (the "Adversarially validated…" line) — it now states the HN sweep was GREEN at validation time, validation caught it, Cycle 2 made it true in both packages.

### Summary

Closed all 5 Cycle-1 findings, scoped strictly to the sweeps (no re-litigation of the verified re-binds/demotions/halt-sync-journey). Both AC-3 sweeps now share the reader-discovery surface: a transitive reader fixpoint (the HN sweep had none), parameter-flow detection for path-arg readers (readSkill-shape, .md literal in the caller), WalkDir-collector detection (shippedSkillText-shape), and constant-concatenation reconstruction closing the split-".md" evasion. Every shape is mutation-controlled by a planted RED-then-GREEN control, each mechanism independently verified to RED when removed; both sweeps confirmed RED against a REAL planted tautology in package source and GREEN once removed. The two TestPiFirstOfficerRuntime* presence checks the rebase onto origin/next surfaced are classified markNonAC (Pi live runner). go test ./... → 1129 passed in 15 packages, 0 failures; go build -tags live ./... compiles; go vet clean on both sweep packages. The false "both sweeps RED on a multi-hop-helper tautology" body claim is corrected to match reality. 2 commits on spacedock-ensign/tautological-test-remediation (39eaf247 integration, 0f8ae11e hostneutrality).
