---
title: A failed dispatch is retried once before anything halts — replace Degraded Mode with a per-entity retry
source: "captain (CL), 2026-07-21 — ruling: an API transport failure must not trigger degradation; nudge the worker for retry. Split out of q4 on the evidence of a two-seat adversarial scope review, which found the defect is the missing rung BELOW the trigger, not the trigger itself."
status: implementation
id: 9q4x5hvyxthc41mt73txr23k
started: 2026-07-21T13:26:28Z
worktree: .worktrees/spacedock-ensign-retire-legacy-and-retry-rung
---

Today the FO contract has no first-failure response at all. It goes straight from one observed dispatch failure to "any SECOND dispatch failure within the session" tripping Degraded Mode — a session-wide, irreversible fallback to sequential bare dispatch. There is no retry, no nudge, and no way back within the session.

## Problem

**Observed twice, both times against a healthy substrate.** On 2026-07-20 two 2ae ensigns died with `API Error: Connection closed mid-response` / `"error":"server_error"` (read from the subagent jsonl, not inferred). Degraded Mode tripped on the second and the session lost concurrent dispatch for its remaining life. The same error pair, ~2 seconds apart, is recorded at `_debriefs/2026-07-04-01.md` issue 2.

**The zombie premise the irreversibility rests on does not describe this runtime.** Measured on 2026-07-21, not argued: the killed 2ae ensign was nudged and RESUMED FROM ITS TRANSCRIPT with 138 turns / 620 KB of context intact and continued its work. Three agents the FO had recorded as "confirmed terminated" were still alive a day later and approved shutdown in ~1 second each; the roster then went 5 → 2 with `reconcile` drift empty. Dead workers are recoverable, the roster is accurate, and shutdown reconciles cleanly.

**BARE MODE IS NOT THE DEFECT AND MUST SURVIVE.** Bare mode — sequential blocking `Agent()`, no `team_name` — is the legitimate teams-unavailable path, live-proven to work by `e3z`. Only the mid-session irreversible TRIGGER is wrong. Any change here must keep bare mode reachable.

**The counter has no mechanism advantage.** `claude-fo-dispatch.md:85` states the trigger has "no time window, no counter — the FO tracks it by its own observation". The contract's stated rationale for bluntness is that an FO under context pressure cannot reliably classify failures — but the counter is itself an unverified FO judgement from the same store. It is a smaller judgement, not a mechanism.

## The open question that decides the design

**Was the failure caused by our own fan-out?** The 2026-07-04 debrief correlates the identical error pair with 6 concurrent live sessions and says root cause was never determined (resource/concurrency ceiling vs. provider-side blip). If the cause is a fan-out ceiling, then retrying at unchanged concurrency reproduces the failure and Degraded Mode's concurrency reduction was the accidentally-correct remedy for the wrong reason. **Resolve this before designing the rung.** Sibling `r5y6` (live-ci-api-error-log-capture) states in its own Problem that a stall "cannot be classified: transient weather vs a CLI retry/hang defect" even with full logs on disk — so any classification-dependent design is currently unevidenced. Land r5y6 first, or accept a design that never classifies.

## Proposed approach

**A failure-kind-AGNOSTIC bounded re-attempt, backed by a durable per-entity ledger. The retry IS the classifier — the FO never parses error strings.**

**What counts as a dispatch failure** (one definition, no error-string vocabulary): `Agent()` returns an error, OR a dispatched worker's session ends with NO completion signal (per `## Awaiting Completion`'s three signals) AND no stage report on the entity. Nothing else. It is read off the dispatch surface the FO already watches; it introduces no transient-vs-structural classification (kept out of scope).

**The rung.** On a dispatch failure of `(entity, stage)`, the FO reads that entity's durable retry ledger:
- **No retry recorded for this `(entity, stage)`** → record one, then re-attempt ONCE:
  - if the failed worker is still addressable (a transport STALL — the live-proven case), the re-attempt is a NUDGE: `SendMessage` resume-from-transcript, preserving the worker's accumulated context (138 turns / 620 KB, measured 2026-07-21);
  - if no live worker remains (an `Agent()` error, or a terminated session), the re-attempt is a FRESH `Agent()` dispatch of the same `(entity, stage)` carrying a distinct `-retry` suffix (AC-1's "incremented suffix"), under dead-ensign handling: mark the prior worker dead in session memory, do NOT cooperatively shut it down.
- **A retry already recorded for this `(entity, stage)`** → this is the second consecutive failure. Hold that entity un-dispatched, surface it to the captain, and stop re-attempting it. Other entities keep running; the session dispatch mode NEVER changes; nothing session-wide degrades.

**Only ONE entity is ever affected, and only reversibly.** A worker death no longer trips any session-wide, irreversible state. Per captain ruling #1 the three Degraded-Mode triggers dissolve to their honest shapes: (1) the second-consecutive-`(entity,stage)`-failure hold above; (2) `/spacedock bare` becomes a plain captain instruction to dispatch bare from that point; (3) `Agent`/`SendMessage` unavailable is the ordinary teams-unavailable condition selecting bare dispatch at the dispatch site. Bare mode itself is untouched.

**The durable retry ledger (the load-bearing mechanism — full design in `## Ideation resolution`).** The bound lives in an FO-written `### Dispatch Retries` subsection on the entity body — durable across compaction and session restart, git-committed, captain-visible, mechanically re-readable before every re-attempt — NOT in FO session memory alone (session memory is a fast-path cache; the ledger is the authority). It reuses the FO body-write pattern that already owns `### Feedback Cycles`, on a DISTINCT axis from the `-cycleN` feedback counter, so a retry never burns a feedback cycle.

## Out of scope

- Retiring bare mode, or any change that makes it unreachable.
- Introducing transport-vs-structural error vocabulary into the contract (the repo has zero occurrences of "transient" under `skills/`).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - A dispatch failure is retried once before anything degrades, and the retry is bounded.**
Verified by: an offline stream fixture asserting that after one `Agent` error the next `Agent` call carries the incremented suffix and NO third attempt appears. Trigger (1) has never had any test in either direction — the replacement must not ship equally unverified.

**AC-2 (VALUE) - Bare mode remains reachable after the change.**
Verified by: a live no-team drive reaching bare dispatch, the same way `e3z` established its baseline — an independent check that can move the wrong way.

**AC-3 - No session-wide irreversible degradation fires on a single worker's death.**
Verified by: the fixture in AC-1 plus a review-time read of the surviving triggers; no committed prose-grep.

## Known co-edits and collisions (found by scope review, verified)

- `internal/contractlint/boot_resident_closure_test.go:53` pins the literal `## Degraded Mode` anchor on disk — mandatory same-commit co-edit.
- `internal/contractlint/dispatch_recovery_value_binding_test.go:35-58` binds the captain-report blockquote VERBATIM to Go. Deleting or editing that prose breaks it.
- **`-cycleN` suffix is already double-booked** — zombie replacement (`fo-dispatch-recovery/SKILL.md:72`) and feedback rework (`fo-dispatch-core.md:55`). A retry burning cycle 2 can trip `feedback-rejection-flow/SKILL.md:17`'s "on cycle 3, escalate to the human" a cycle early. A distinct `-retry` axis avoids it but contends with the 64-char name cap.
- **`claude-fo-dispatch.md:68` bans re-dispatching a replacement ensign** ("you have no evidence the first ensign failed"). A retry rung needs an explicit carve-out for the observed-failure case or the two rules contradict and the FO improvises — the exact behavior the `## Awaiting Completion` anti-pattern list exists to suppress.
- **Shutdown collision:** `fo-dispatch-core.md:55` demands `«worker.shutdown»` of the prior cohort before a `-cycleN` dispatch; `fo-dispatch-recovery/SKILL.md:70` forbids cooperative shutdown to a dead ensign. A transport-killed worker fires both. Getting it wrong yields two live ensigns on one worktree — a correctness failure, worse than the throughput failure degradation costs.
- **`internal/claudeteam/claudeteam.go:68-79` `BareModeAdvisory`** tells the FO to run `ToolSearch select:TeamCreate` on every bare dispatch lacking recent team evidence. After q4 retires TeamCreate that names a tool which no longer exists, on the legitimate bare path.
- **The retry bookkeeping needs a durable answer, not prose.** "Retry once" is bounded only if the FO reliably remembers it retried; otherwise it is an unbounded respawn loop against a dead API with no captain-visible event. This is the load-bearing design question.

## Sequencing against q4

**q4 deletes the repo's only existing retry rung.** `using-legacy-claude-team/SKILL.md:50` («legacy-team.recover» rung 1, "Attempt one new TeamCreate with a fresh name") is the sole try-once-before-degrading step in the contract, and q4 retires that file. If q4 lands first without this entity, the contract is left with strictly ZERO retry surface. Land together, or land this first.

**Byte contention.** q4's AC-4 re-tightens `foFunctionReferenceBaselineBytes` to measured+1 at zero slack. A scope seat measured this change at roughly +164 to +204 B net (fits today's 507 B headroom with ~300 spare), but ~+900 B if the retry bookkeeping gets its own skill section — which does NOT fit and would require q4's removal of `using-legacy-claude-team/SKILL.md` (14065 B) from the measured surface first. Budget the two jointly.

## Test plan

**The offline oracle is the load-bearing test — it is the ONLY thing that runs in CI.** Coverage vacuum CONFIRMED this ideation: `TestLiveDegradedBareRecovery` (and `TestLiveBreakGlassShimRecovery`) are `//go:build live` and named in ZERO `-run` filters in `.github/workflows/runtime-live-e2e.yml` (the live lanes run `TestLiveEnsignCycle` / `TestLiveDefaultHeadlessStopsAtGate` / `TestLiveZeroDiscoverReportsAndStops` / `TestLiveClaudeSharedScenarios` / pty / merged-team — never degraded-bare or break-glass), and line 111 pins `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`, so the teams-unavailable branch never runs. Only the default-tagged unit oracle executes in CI.

1. **AC-1 offline oracle (net-new, default-tagged — RUNS IN CI):** `assertBoundedRetryObservables` + a positive fixture (an `Agent()` error → exactly one `-retry` re-dispatch → NO third attempt) + RED controls (a third-attempt stream must fail; a no-retry stream must fail). Mirrors the existing `internal/ensigncycle/dispatch_recovery_assert_unit_test.go` shape. Falsifiability: dropping the bound reds the "no third attempt" control; dropping the retry reds the "retry once" control.
2. **AC-2 live no-team drive (rewritten):** `TestLiveDegradedBareRecovery` retargeted to the post-retirement expectation — a plain `/spacedock bare` instruction yields bare-shaped `Agent()` calls (no `name`, no `run_in_background`) WITHOUT the retired Degraded-Mode captain report and WITHOUT the `spacedock:fo-dispatch-recovery` load. `assertDegradedBareObservables` is rewritten to assert bare-shape-minus-ceremony; a drive still emitting the report or loading the recovery skill is now a FAILURE. Still `//go:build live`; adding it to a CI `-run` filter needs the teams-unavailable lane un-pinned (out of scope; flagged).
3. **Contractlint co-edits:** `TestDeferredSkillCoresResolveAndCarryCeremony` stays green after the `deferredSkillCores` anchor drop; `TestDegradedModeCaptainReportPrefixBindsSkillBlockquote` deleted with its `degradedModeCaptainReportPrefix` const; `TestFOFunctionPromptSurfaceShrinks` green at the jointly-re-tightened baseline (q4 AC-4). Both Degraded-Mode bindings were spike-verified to red on an anchor/blockquote change and reverted clean.
4. No committed prose-grep (per proof policy); the surviving-trigger check (AC-3) is a review-time read.

## Captain rulings (2026-07-21, recorded at shaping)

1. **Degraded Mode is retired entirely — no session-wide backstop survives.** This answers the open ideation question (the two review seats' disagreement). The three triggers dissolve, each to its honest shape: (1) a second consecutive failure of the SAME entity/stage halts THAT entity and reports to the captain — never the session; (2) `/spacedock bare` becomes a plain captain instruction to dispatch bare from that point, not an irreversible mode transition; (3) `Agent`/`SendMessage` unavailable is not a mode at all — it is the ordinary teams-unavailable condition selecting bare dispatch, evaluated where dispatch happens. Bare mode itself survives untouched (already a hard constraint above). Co-edits stand as mapped: `fo-dispatch-recovery`'s Degraded Mode section, the two contractlint anchors, and the recovery-report verbatim binding.
2. **Title de-minted** (the metaphor family is banned vocabulary): the slug stays stable per the no-slug-churn precedent; the title now reads plain.
3. **Recommended to the captain, pending a nod:** accept the failure-kind-agnostic design both seats converged on (the retry IS the classifier; no error-string vocabulary enters the contract), which removes the r5y6 sequencing dependency — the "was it our own fan-out" question stops mattering because a retry that reproduces the failure halts that entity either way.

## Ideation resolution (2026-07-21)

### Open question — fan-out ceiling vs transient blip: NO cap required

**Verdict: the design does NOT need a fan-out cap.** THIS session is direct behavioral counter-evidence to the fan-out-ceiling hypothesis. It hit ~4 transport stalls (`Connection closed mid-response` / `stream watchdog stall, no progress 600s`), and EVERY ONE was cleared by a single FO nudge (`SendMessage` resume-from-transcript). A nudge resumes the stalled worker at UNCHANGED total concurrency — it drops no fan-out; every other worker keeps running. A hard concurrency ceiling would re-stall a worker resumed while the rest of the fan-out is still live; instead every resume succeeded and the worker continued. That is the signature of an isolated transient stall a bounded retry clears, not a standing ceiling that retry-at-unchanged-fan-out reproduces.

Honest caveat: this rests on the live session summary, not a raw per-stall concurrency/timestamp table, so I cannot certify the four stalls were non-simultaneous. But the argument holds either way — resuming ANY stalled worker without lowering concurrency succeeded, which a hard ceiling forbids. Contrast the lone 2026-07-04 data point (identical error ~2 s apart, right after a 4-way opus spawn): a correlated signature, but root cause was never determined, and it stands against this session's four independent recoveries.

**The design is robust to the answer regardless** — confirming captain ruling #3's lean and retiring the r5y6 dependency. If a ceiling ever DID bind, the failure-kind-agnostic retry reproduces it exactly once, then holds THAT entity — never the session, never an unbounded loop. So the fan-out question stops gating the design; a cap would be premature complexity against evidence pointing the other way. The bounded retry is the safety net either way.

### The load-bearing question — durable retry bookkeeping

**How "retry once" survives context pressure.** The Degraded-Mode counter was deliberately counter-free because "the FO cannot reliably track failures across context pressure" (`claude-fo-dispatch.md:85`). A retry bound stored ONLY in FO session memory inherits that exact unreliability: after a compaction the FO forgets it retried and can respawn a dead API unboundedly, with no captain-visible event. So the authority MUST be durable.

**Design: an FO-written `### Dispatch Retries` ledger on the entity body**, one line per retry:

    - Retry 1: {stage} — {agent-error | no-completion-signal}; {nudged | re-dispatched -retry}

It is durable (survives compaction and session restart), git-committed and captain-visible, and re-read before every re-attempt. Session memory stays a fast-path cache to avoid a re-read each turn, but the ledger is the tiebreak. It reuses the existing `### Feedback Cycles` write pattern (FO-owned body subsection, worktree-side when `worktree:` is set, main-side otherwise), so it needs NO new Go/status field and NO new frontmatter key — the write mechanism already ships. The HOLD of a twice-failed `(entity, stage)` needs no new state either: the entity is held un-advanced (its ledger shows the retry exhausted) and surfaced to the captain; the existing `mod-block` hold is available if a machine-readable "do not re-dispatch" flag is later wanted, but ledger + hold already bound it. The distinct `### Dispatch Retries` axis keeps a retry from advancing `### Feedback Cycles` and tripping `feedback-rejection-flow/SKILL.md:17`'s cycle-3 escalation early.

### Co-edit resolutions

1. **`boot_resident_closure_test.go:53` anchor** — same-commit co-edit: remove `## Degraded Mode` from the `deferredSkillCores` entry for `fo-dispatch-recovery/SKILL.md` when the section is removed. The entry keeps `## Break-Glass Manual Dispatch` and `## Context Budget Failure and Dead Ensign Handling` (both survive). GROUNDED: renaming the heading in the spike red `TestDeferredSkillCoresResolveAndCarryCeremony` with `missing section anchor "## Degraded Mode"`; reverted.
2. **`dispatch_recovery_value_binding_test.go:35-58` blockquote** — delete `TestDegradedModeCaptainReportPrefixBindsSkillBlockquote` AND the `degradedModeCaptainReportPrefix` const it binds (`internal/ensigncycle/dispatch_recovery_assert_test.go:17`): the verbatim captain report is retired, so there is nothing left to bind. GROUNDED: paraphrasing the blockquote in the spike red the value-binding test; reverted.
3. **`-cycleN` double-booking** — resolved by a DISTINCT `-retry` axis. The retry re-dispatch carries a `-retry` suffix, never a `-cycleN` increment, and the bound lives in `### Dispatch Retries`, not `### Feedback Cycles` — so a retry cannot advance the feedback counter or trip cycle-3 early. Name-cap: `-retry` is one axis token (not `-retryN`) over the same `{worker_key}-{slug}-{stage}` stem, capped the way `spacedock dispatch build` already caps (`fo-dispatch-recovery/SKILL.md:50`).
4. **`claude-fo-dispatch.md:68` re-dispatch ban** — the `## Awaiting Completion` ban ("re-dispatch a replacement ensign — you have no evidence the first ensign failed") is scoped to the NO-evidence case (a still-pending worker). The retry rung fires ONLY on an OBSERVED dispatch failure (the two-part definition), which is exactly the evidence the ban says is absent. Add a one-clause carve-out at :68 naming the observed-failure retry as the exception, so the two rules read as complementary.
5. **Shutdown-vs-dead-ensign contradiction** — a transport-killed worker is DEAD, so dead-ensign handling governs (`fo-dispatch-recovery/SKILL.md:70`: do NOT cooperatively shut down a dead ensign), NOT supersede-shutdown (`fo-dispatch-core.md:55`, which demands cooperative shutdown of a LIVE prior cohort before a `-cycleN` dispatch). The retry re-dispatch is explicitly NOT a `-cycleN` increment and its predecessor is dead, so the supersede precondition never holds; the rung states: mark the failed worker dead, skip cooperative shutdown, fresh-dispatch — avoiding two-live-ensigns-on-one-worktree. (For a STALL where the worker is still addressable, the re-attempt is a nudge to the SAME worker — no second worker is spawned at all.)
6. **`claudeteam.go:68-79` `BareModeAdvisory` TeamCreate probe** — OUT of 9q4's scope. It is a q4 co-edit (retiring TeamCreate naming); 9q4 retires Degraded Mode, which does not touch `claudeteam.go`. Since the two land jointly, q4 owns the `BareModeAdvisory` rewrite and 9q4 leaves it untouched.

### Sequencing against q4 — land WITH or BEFORE q4; budget the byte ratchet jointly

Two couplings force the ordering:
- **Retry-surface coupling:** q4 deletes `using-legacy-claude-team/SKILL.md`, whose `«legacy-team.recover»` rung 1 (:50) is the contract's ONLY try-once-before-degrading step. q4-alone-first leaves zero retry surface — a regression. 9q4 must be in place first or in the same change.
- **Binding coupling (found this ideation):** 9q4's removal of `## Degraded Mode` from `fo-dispatch-recovery` reds `TestUsingLegacyClaudeTeamDegradedModePointersNameRealAnchor` — using-legacy's `**Fall back to Degraded Mode**` pointer (:51) would name a section that no longer exists, and the `deferredSkillCores` anchor it reads is being removed. That test's subject is deleted by q4. So removing the anchor strictly BEFORE q4 deletes using-legacy would require throwaway edits to a file q4 discards.

**Firm verdict: land 9q4 and q4 as ONE coordinated change (joint commit / same session)** — q4 deleting using-legacy and 9q4 removing the Degraded-Mode section together, so neither contractlint test is transiently red. Never q4-alone-first.

**Byte budget (jointly, grounded on THIS disk).** Measured FO surface = **123236 B** against baseline **123323 B** — only **87 B of headroom**, far tighter than q4's cited 507 B (this tree is already at bw's 123323 re-baseline; re-measure at implementation). 9q4 alone is net-NEGATIVE: it removes the ~3131 B `## Degraded Mode` section from `fo-dispatch-recovery` and replaces the ~1.15 KB `## Degraded Mode (trigger)` in `claude-fo-dispatch.md` with the compact retry rung (est. +0.8–1.5 KB), net ≈ −2.4 to −2.9 KB. q4 additionally removes `using-legacy-claude-team/SKILL.md` (14065 B). q4's AC-4 re-tightens the baseline to the post-both measured+1 at zero slack. The retry-bookkeeping MUST stay resident/compact — giving it "its own skill section" (~+900 B) is unnecessary because it reuses the `### Feedback Cycles` write pattern and adds no new skill.

### Estimate

**Expected surface** (prose + Go co-edits + net-new fixture):
- Prose (measured FO-surface files): `fo-dispatch-recovery/SKILL.md` — remove the `## Degraded Mode` section (~−3.1 KB) and drop Degraded-Mode vocabulary from its frontmatter `description`. `claude-fo-dispatch.md` — replace `## Degraded Mode (trigger)` (:83-85) with the resident retry rung + the two dissolved-trigger clauses; add the :68 carve-out; note the `### Dispatch Retries` ledger. `first-officer-shared-core.md:49` — reword the fo-dispatch-recovery one-liner (drop "Degraded Mode"). `using-legacy-claude-team/SKILL.md` Degraded-Mode pointers (:25, :51) — subsumed by q4's file deletion (joint land).
- Go co-edits: `boot_resident_closure_test.go` (drop the `## Degraded Mode` anchor); `dispatch_recovery_value_binding_test.go` (delete the blockquote-binding test); `dispatch_recovery_assert_test.go` (delete `degradedModeCaptainReportPrefix` + `assertDegradedBareObservables`; add the bounded-retry + bare-reachability oracles); `dispatch_recovery_assert_unit_test.go` (rewrite offline positive + RED controls); `dispatch_recovery_live_test.go` + `dispatch_recovery_fixtures_test.go` (retarget `TestLiveDegradedBareRecovery` / `degradedBarePrompt` to the post-retirement expectation); `fo_function_reference_invariant_test.go` (baseline byte number — via q4's AC-4, joint).
- Net-new offline fixture (AC-1): `assertBoundedRetryObservables` + synthetic streams (positive + two RED controls), default-tagged so it runs in CI.

**Tolerance:** ~8–12 files (6–8 Go, 3–4 prose). Prose net −2.4 to −2.9 KB for 9q4 alone; a much larger drop jointly with q4. The breach to watch: a mis-scoped verbose rung landing net-POSITIVE on the surface — the rung must stay under ~1.5 KB.

**Effort:** ~1.5–2 implementation sessions. Prose rung + co-edit deletions are mechanical (~½ session); the oracle rewrite (retire degraded-bare, author the bounded-retry oracle + RED controls + retarget the live test/prompt) is the bulk (~1 session). Joint coordination with q4 adds coordination overhead, not independent effort.

**Test-coverage note:** the change lands with NO live oracle in CI (`TestLiveDegradedBareRecovery` is in zero `-run` filters and the teams-unavailable branch is never exercised — line 111 pins `TEAMS=1`). The default-tagged AC-1 offline oracle is therefore the SOLE CI guard for the new behavior in either direction, and MUST ship in the same commit as the rung.

### Acceptance criteria — firmed

- **AC-1 (firm):** keep. Concretely a default-tagged `assertBoundedRetryObservables(stream)`: an `Agent()` error is followed by exactly ONE re-dispatch carrying the `-retry` suffix and NO third `Agent()` call for that `(entity, stage)`; plus RED controls (third-attempt stream fails; no-retry stream fails). This fills the coverage vacuum and runs in CI.
- **AC-2 (firm, VALUE):** keep the live no-team drive, but the oracle CHANGES — post-retirement `/spacedock bare` produces bare-shaped `Agent()` WITHOUT the Degraded-Mode captain report and WITHOUT the recovery-skill load; a drive still emitting either is now a FAILURE (the wrong-way check preserved).
- **AC-3 (firm):** keep. The AC-1 fixture proves the retry-then-hold-that-entity path affects one entity reversibly; a review-time read of the three dissolved triggers confirms none flips session-wide mode. No committed prose-grep.

## Stage Report: ideation

- DONE: Resolve the OPEN QUESTION (fan-out ceiling vs transient blip)
  Verdict: NO cap. This session's ~4 transport stalls each cleared by a nudge at UNCHANGED concurrency — a hard ceiling would re-stall on resume; it never did. Design is robust either way (a ceiling just holds one entity after one failed retry). See `## Ideation resolution`.
- DONE: Design the failure-kind-agnostic bounded retry rung
  Filled `## Proposed approach`: one dispatch-failure definition (Agent error OR session-end-without-signal), one bounded re-attempt (nudge if addressable, else fresh `-retry` dispatch), second consecutive same-`(entity,stage)` failure holds THAT entity + reports captain; session mode never changes.
- DONE: Resolve THE LOAD-BEARING QUESTION (durable retry bookkeeping)
  Durable FO-written `### Dispatch Retries` ledger on the entity body is the authority (session memory a cache) — survives compaction/restart, git-committed, captain-visible, reuses the `### Feedback Cycles` write pattern so no new Go/frontmatter surface. Distinct axis from `-cycleN`.
- DONE: Resolve the co-edits (2 contractlint bindings spike-verified)
  All 6 resolved in `## Ideation resolution`. GROUNDED: perturbing `## Degraded Mode` red `TestDeferredSkillCoresResolveAndCarryCeremony`; paraphrasing the captain-report blockquote red `TestDegradedModeCaptainReportPrefixBindsSkillBlockquote`; both reverted clean in the spike worktree.
- DONE: Firm the sequencing with q4
  Land as ONE coordinated change (never q4-alone-first): retry-surface coupling + a newly-found binding coupling (removing the Degraded-Mode anchor reds `TestUsingLegacyClaudeTeamDegradedModePointersNameRealAnchor`, whose subject q4 deletes). Byte budget jointly — measured 123236 B / baseline 123323 B = 87 B headroom on disk.
- DONE: Produce the estimate + firm AC-1..AC-3
  Surface, tolerance (~8–12 files, prose net −2.4 to −2.9 KB alone), effort (~1.5–2 sessions), and the coverage-vacuum note (AC-1's default-tagged offline oracle is the SOLE CI guard). AC-1/AC-2/AC-3 firmed in `## Ideation resolution`.

### Summary

Firmed the design to captain ruling #1 (Degraded Mode retired entirely; triggers dissolve to honest shapes). The open question is resolved AGAINST a fan-out cap on this session's direct evidence — every transport stall cleared by a nudge at unchanged concurrency — and the design is robust to the answer regardless, so the r5y6 dependency is retired. The load-bearing bookkeeping is a durable FO-written `### Dispatch Retries` ledger (not session memory alone), reusing the `### Feedback Cycles` write pattern for zero new Go/frontmatter surface. Both contractlint bindings were spike-verified to red on an anchor/blockquote change and revert clean; a NEWLY-found binding coupling forces 9q4 and q4 to land as one coordinated change. Byte surface measured on disk at 87 B headroom (tighter than q4's cited 507 B) — 9q4 alone is net-negative, so it fits. AC-1's offline oracle is the sole CI guard and must ship with the rung.

## Gate: ideation — APPROVED (FO, captain-ratified, JOINT with q4)

- **Verdict:** approved for JOINT implementation with `q4` — ONE worktree, ONE commit, never q4-alone-first.
- **Captain ruling (CL, 2026-07-22):** Degraded Mode is retired; the no-TeamCreate auto-team model is the go-forward floor (q4's governance decision). This unblocks the joint landing.
- **Validation:** the bare no-team live drive (9q4 AC-2, runnable on this session) is the shared live proof; the AC-1 default-tagged offline oracle is the sole CI guard and MUST ship in the same commit as the rung; detached adversarial audit applies (shipped FO contract).
- **Base:** worktree off `origin/main` (`ca136f83`).
