---
id: j903f6f1vgckk3kj6j6zbmt3
title: Lazy-TeamCreate + shallow-boot-then-greet — reach interactive readiness fast
status: ideation
source: "captain (2026-06-04) — boot analysis (this session) measured boot at ~7:36, ~80% model-compose; TeamCreate is the single largest write (89k cache_creation — the whole prompt prefix re-cached to the 1h cache when team-mode activates) and sits on the critical path BEFORE the gate that decides whether dispatch even happens."
score: "0.34"
started: 2026-06-13T17:02:55Z
completed:
verdict:
worktree:
issue:
---

Make the first officer cheap to boot and reach interactive readiness fast. Boot forensics on a live FO session (`docs/roadmap/0203-fo-efficiency/boot-analysis.md`) measured **~160k peak context and ~13.6 min wall-clock to reach an interactive greet — with no team created and no worker dispatched.** A 100% pre-dispatch session paid full deep-boot cost. The FO front-loads its entire contract (both reference files ~16.2k), the full workflow README (~8.1k), both mod files (~6.5k), and renders the human status table (twice, ~8.7k) — then thinks at 100k+ context, where the two slowest turns (128.6s, 100.1s) fired. None of that was needed to greet the captain.

This is the **backbone of the 0.20.3 "FO efficiency" milestone** (`docs/roadmap/0203-fo-efficiency/README.md`). It is ONE task, THREE sequential phases: P1 contract structural split → P2 lazy-TeamCreate → P3 shallow-boot-then-greet. The ideation is consolidated from two completed static-analysis spikes (`T1-ideation.md` = Phase 1, verdict VIABLE; `j9-phases-2-3-ideation.md` = Phases 2-3, verdict VIABLE).

## Problem

- **Time-to-interactive is gated on the full deep-boot, not a fast orient pass.** Generation latency scales with loaded context, so the wall-clock is dominated by thinking at 100k+ — most of it over material the greet never uses.
- **TeamCreate's ~89k cache-creation is on the critical path before the dispatch gate.** Creating a Claude team re-caches the whole conversation prefix under the new team context. The milestone's Cost-levers table ranks this as the single biggest lever (~89k removed) — far larger than the contract reads (~16k), the status-table render (~8.7k), or the mod reads (~6.5k). The measured session never created a team because it never dispatched, so the entire ~89k is avoidable on a no-dispatch boot — *if* the contract stops telling the FO to create a team at startup.
- **The contract is read whole, up front, before any work begins.** ~70% of the Claude runtime adapter is team/dispatch material a no-dispatch boot never uses, yet it loads first (forensics rows 2-3, t=8-9s).

## Cost levers (why this task is the backbone)

| Lever | ~boot cost removed | Needs the contract split? |
|-------|-------------------:|---------------------------|
| Lazy-TeamCreate (defer the team-mode prefix re-cache) | **~89k cache-creation** | no (Phase 2) |
| Defer contract reads at greet | ~16k | yes — minimal (Phase 1) |
| Defer the human status-table render | ~8.7k | no (Phase 3) |
| Defer mod-file reads | ~6.5k | no (Phase 3) |
| README frontmatter-only read (defer the body) | ~7.7k | no (Phase 3) |

The ~89k lazy-TeamCreate dwarfs the rest and sits on the critical path *before the dispatch gate*. The contract split (Phase 1) is the enabling refactor for the one lever that needs it, and doubles as the "audit and clean up the FO contract" ask.

## Proposed approach

Three sequential phases. Each is behavior-preserving where it moves content; the *new* behavior (greet-first sequencing) is Phase 3 and rests on Phase 1's split.

### Phase 1 — contract structural split (enabler + the contract-cleanup ask)

Reorganize the FO contract by *when* it is needed so a boot reads only what it needs to greet, report state, and present a gate. The split is the structural enabler for the deferred-module unload; it is the contract-audit ask itself.

- **Slim the boot-resident core.** `first-officer-shared-core.md` keeps its boot-resident sections (Startup, Status Viewer, ID Styles, Single-Entity Mode, Working Directory, State Management, FO Write Scope, Clarification, Working Principles, Probe discipline, Issue Filing) plus the gate-presentation/AC-cross-check spine extracted from Completion-and-Gates. The runtime adapter keeps Captain Interaction's greet/guardrail, Agent Back-off, and Entity-Body Inspection.
- **Add a top "Operating principles (ethos)" section** the shipped skill currently lacks. Its absence lets a host (Codex in particular) drift from the `agents/first-officer.md` ethos. This is boot-resident GUIDANCE content (behavior-shaping prose, not a testable contract clause) and the existing `## Working Principles` folds under it. These principles also govern *how* the contract is simplified across all three phases (begin-with-the-end, de-risk-cheap, simplest-approach). Verbatim text to add at the top of the boot-resident core:
  > You are dispatcher and responsible for making sure the work is done by the crew. What awesome looks like for the crew:
  > - Begin with the end, be clear about the value.
  > - Do the hardest things first, de-risk when it is cheap.
  > - Communicate and act concisely, choose the simplest approach, JFDI.
- **Drop the unnecessary `agents/first-officer.md` cross-reference** (shared-core line 3, "Keep aligned with …") — not load-bearing; the skill is the single loader and the ethos now lives in-skill (above).
- **Extract the dispatch/team module** into a lazily-loaded reference (e.g. `references/claude-fo-dispatch.md`) read at first dispatch: Worker Resolution, Dispatch Adapter, Context Budget, Event Loop (incl. the reconcile sweep), Degraded-Mode seams, standing-teammate discovery/lazy-spawn/declaration, and the shared-core Dispatch / Worktree-Ownership / Standing-Teammates / reuse-conditions. The generic team lifecycle stays in the already-lazy `using-claude-team` skill (non-duplicative — the new reference holds only the spacedock adapter, which the generic skill leaves to the consumer). The first-dispatch load point is the existing `Skill(skill="spacedock:using-claude-team")` invocation, extended to also pull the dispatch reference.
- **Extract the merge module** into a lazily-loaded reference read at terminalization: Merge-and-Cleanup, Ship-Local Ceremony, Worktree-removal safety, Mod-Hook/Mod-Block Enforcement (shared-core), and Mod-Block-at-Terminal (runtime). The boot-resident core reaches it the same way it reaches `present-gate` / `feedback-rejection-flow` — by naming the load point at the terminal boundary.
- **Loader change.** `SKILL.md`'s `@references/first-officer-shared-core.md` inlines only the slimmed core; the runtime-adapter read loads only the boot-resident runtime sections. The dispatch/merge references are NOT read at startup. No new pattern — this extends the existing lazy precedent (`present-gate`, `feedback-rejection-flow`, `using-claude-team`).

**Coupling resolutions (Phase-1 spike step 3):** C1 — the shared-core Startup procedure is a self-contained boot island with ZERO forward reference into deferred content (its only forward reference, the rebase-conflict halt, targets boot-resident State Management) — this is what makes the cut clean. C2 — the runtime adapter's "At startup (after reading the README, before dispatch)" team-creation framing is retired in favor of the truthful "at first dispatch" trigger (one clause; the next clause already says "before the first team-mode tool call"). C3 — the reconcile sweep fires "before the first dispatch," not before-greet, and needs a roster-derived `team_name`, so it travels with the dispatch module. C4 — standing-teammate discovery/spawn is wholly inside the team-creation flow ("no spawn at boot"); it travels with the dispatch module. C5 — Completion-and-Gates splits: the gate-presentation spine stays boot-resident; reuse-conditions + the budget probe (only reached after a worker completes) move to the dispatch module, on the present-gate/feedback-rejection precedent. C6 — Mod-Block enforcement (both files) is referenced only from Merge-and-Cleanup; it travels with the merge module.

**Effect:** boot read drops ~15,700 → ~4,500 tokens (~11k cut, ~70% of contract-read cost) on a no-dispatch boot.

### Phase 2 — lazy-TeamCreate (the headline ~89k lever)

Change the runtime adapter's `## Team Creation` trigger so the FO creates no team at boot.

- **The single clause** (runtime adapter `## Team Creation`, line 7). Current: "At startup (after reading the README, before dispatch), invoke the generic Claude-team-harness discipline: `Skill(skill="spacedock:using-claude-team")`". New: **"Before the first team-mode dispatch (the first `Agent()` call that uses a `team_name`), invoke …"**. Drop "at startup"; the existing line-11 clause ("before the first team-mode tool call in the session") becomes the consistent trigger.
- **Companion timings (already aligned by Phase 1's placement).** Standing-teammate discovery/spawn is already "no spawn at boot" / "before the first `Agent()` with a `team_name`"; the reconcile sweep is already "before the first dispatch." Phase 1's placement of these in the deferred dispatch module is what realizes the deferral — the FO never reads them, or creates a team, at boot. No binary change.

### Phase 3 — shallow-boot-then-greet (greet off `status --boot --json`)

Reshape the Startup procedure so the FO reaches the greet through exactly the before-greet step set, then stops for input. The new shallow-boot sequence:

```
contract-gate (S1: version mismatch aborts before any state read)
  → git root (S2)
  → workflow discovery (S3)
  → README FRONTMATTER read (S4-slim: entity-label, stage names/ordering, gate/terminal flags)
  → status --boot --json (S5: state-summary source — orphans, pr_state w/ LIVE merge state, dispatchable, team_state, state_backend)
  → split-root halt-gate (S6: prevents a false-empty greet off an uninitialized checkout)
  → split-root pull-on-boot (S7: state freshness; rebase-conflict HALT)
  → GREET the captain (state summary off the boot JSON; able to present a gate) and STOP for input
```

- **Greet off the boot JSON.** The FO greets with a state summary built from `status --boot --json` (orphans, pr_state, dispatchable, team_state, state_backend) and the README frontmatter (entity-label, stage taxonomy, gate flags), then stops. It can present a gate from this state without a team (gates are captain-facing text, not team messages).
- **README frontmatter-only read (S4-slim).** Startup step 4 reads the README frontmatter (~175 tokens: entity-label, stage names/ordering, gate/terminal flags) and defers the ~7.7k body (per-stage prose, proof policy, templates, CI docs) to the phase that consumes it. The boot JSON cannot substitute for the frontmatter — confirmed in `internal/status/json_commands.go`: the boot envelope carries NO `"stages"`/`"labels"`/`"gates"` key, so the frontmatter read stays before-greet; only the body defers.
- **Defer mod-file reads.** The comm-officer and pr-merge mod FILES are read when their hooks fire — comm-officer at first dispatch (spawn needs a team), pr-merge advancement on the first event-loop pass. The greet reports PR/mod/orphan state off the boot JSON, never the mod files.
- **Defer the human status-table render.** The FO never renders the human-formatted table for its own reasoning (it has the boot JSON); it renders it to the captain at most once, on explicit request, per the Status Viewer's existing rules.

**Correctness hunt (Phases 2-3 spike step 3) — no guarantee is dropped by deferring past the greet.** Each guarantee carried by a deferred step is either (a) re-asserted as a state READ in `status --boot --json` so the greet stays accurate, or (b) an ACTION the first event-loop pass / first dispatch catches before it matters: merged-PR advancement (caught by Event-Loop step 1 + the pr-merge idle hook; greet reports merge state off the boot JSON); orphan surfacing (the boot JSON's `orphans` array, before-greet); mod-block resume (Event-Loop step 2, a first-action obligation); supersede/teardown (vacuous at a shallow boot — no team exists to clean). **The one genuinely-cannot-defer step is the split-root halt-gate + pull-on-boot (S6/S7)** — and it is already a before-greet shared-core step, so it is part of the sequence, not an obstacle.

### Documentation changes

No user-facing CLI output, startup banner, or docs-site surface changes. The boot-flow change is internal to the FO contract (skill references read by the model). The 0203 milestone README in `docs/roadmap/` is a roadmap artifact, not docs-site content; it already records this design. No doc diff is owed.

## Out of scope

- **The residual-prose audit + comm-officer polish (T3 in the milestone).** The cut-list of leftover prose to trim does not exist until the Phase-1 split lands; T3 files along after. SEPARATE task.
- **p2 / vc — `spacedock pr complete` + `reconcile --act`** (the binary-simplification line). Higher ROI, heavier lift → 0.20.4. SEPARATE task.
- **Extending `status --boot --json` to carry labels + stage taxonomy** (the README-read shape that would remove the frontmatter read entirely). A heavier binary change with its own golden-fixture/schema test surface; a follow-on, not Phase 3.
- **The Codex and Pi runtime adapters' lazy-TeamCreate + dispatch/merge extraction.** This task splits the Claude adapter (the bulk file the forensics measured, ~6.9k); the Codex (~4.9k) and Pi (~5.9k) adapters are an order smaller — a follow-on if their boot cost warrants it. The shared-core split applies cross-host.
- **xp** (cross-session FO↔Commander comms) and **ey** (proof-policy port to shipped scaffolding) — own tracks.

## Acceptance criteria

Each AC names an end-state property of the finished task, verified by something OUTSIDE this task body that can fail. **No AC is proven by a string/substring/regex match over any instruction file the model reads** — that is banned by this workflow (a passing match only asserts the implementer's own text is present; a valid paraphrase fails it and an inverted clause passes it). The behavioral ACs are live drives; the structural ACs test a relationship between independent values.

**On the Phase-1 GUIDANCE content (the Operating-principles/ethos section, the `agents/first-officer.md` cross-reference drop):** this is behavior-shaping prose, NOT a testable acceptance criterion. Its proof is the existing live scenarios (AC-3) still passing + the high-stakes detached audit + gate review — there is no "ethos present" or "drift reduced" metric, because such a check would be either a banned prose-grep over the contract or a tautology. The ethos is real authoring work that ships with the split; it carries no AC of its own.

**AC-1 (behavioral, live) — A freshly-booted FO greets the captain and reports accurate workflow state with NO team created and the dispatch/merge modules NOT loaded, then stops for input rather than auto-dispatching.**
Verified by: a new live shared-runtime scenario `shallow-boot` in `internal/ensigncycle` (added to `sharedRuntimeScenarios()` with Claude + Codex runners per the README's 4-step add-a-scenario procedure). The runner launches the real host front door against a fixture with at least one dispatchable entity sitting at a human gate; the host-neutral assertion over `(before, after, observed)` confirms: (a) the FO produced a greet with a state summary in its final message and a gate review + decision prompt; (b) durable state shows NO team artifact created and NO worker dispatched — entity still at its gate stage, not advanced/archived, no `verdict`/`completed` set, no worktree created; (c) the FO stopped for input. The team-not-created observation is the lazy-TeamCreate proof; the greet-without-dispatch is the shallow-boot proof; the gate-still-at-stage is the gate-guardrail surviving with the deferred modules unloaded. An offline negative case in `shared_scenarios_negative_test.go` builds the broken end-state (a team artifact present / an entity dispatched or self-approved at boot) and proves the assertion goes red. (Live, real host, durable-state assertion — never a contract grep.)

**AC-2 (behavioral, live) — No team-mode prefix re-cache occurs before the greet.**
Verified by: the `shallow-boot` live scenario's captured host artifacts (stream jsonl / session transcript) show NO `TeamCreate` (or team-tool) call in the pre-greet window — a behavioral observation over the real run's tool-call sequence and ordering relative to the greet message, NOT a grep over the contract. The negative control: a deep-boot run (or a fixture forcing eager team creation) shows the `TeamCreate` call before the greet, proving the assertion distinguishes the two. (Rides AC-1's live run; zero additional model spend.)

**AC-3 (behavioral, regression, live) — The split is behavior-preserving: the deferred modules load and function correctly when a real dispatch or merge happens, and the deferred startup steps still fire on first action.**
Verified by: the existing live scenarios `gate-guardrail`, `rejection-flow`, `merge-hook-guardrail`, and `feedback-3-cycle-escalation` in `internal/ensigncycle` still pass after all three phases. `gate-guardrail`/`rejection-flow` exercise the dispatch module loading at first dispatch (reuse-conditions, feedback routing, the reconcile sweep); `merge-hook-guardrail` exercises the merge module loading at terminalization (mod-block enforcement, merged-PR advancement). A green run is the proof that lazy-loading + deferring the startup steps past the greet did not drop a reachable instruction or a first-action obligation. Run via `go test -tags live -run TestLiveClaudeSharedScenarios ./internal/ensigncycle`. (Zero new authoring beyond AC-1's scenario.)

**AC-4 (structural) — The boot-resident core has no reference dependency on deferred-only content, and the dispatch/merge sections are absent from the boot-resident files.**
Verified by: a reference-closure + structural-absence guard in `internal/contractlint` (the allowed quarantine), extending `TestUserSkillReferenceClosureResolves` and `TestRetiredPluginPrivatePathsAbsent`. (a) The closure check parses the boot-resident files and the deferred-module manifest as real artifacts and fails if a boot-resident `@`/read reference resolves into a deferred-module section that is not one of the declared lazy load points (`using-claude-team`, `present-gate`, `feedback-rejection-flow`, the new dispatch/merge references). (b) The absence check confirms the dispatch-module anchors (Dispatch Adapter, Worker Resolution, Event Loop, Context Budget, standing-teammate spawn) and merge-module anchors (Merge and Cleanup, Ship-Local Ceremony, Mod-Block Enforcement) are ABSENT from the boot-resident files and PRESENT in their deferred references. The expected value (which anchors belong where, which references are lazy load points) comes from the deferred-module manifest — an independent source the boot files can diverge from, so the check can fail on a future edit that points the boot core at a moved section. This tests a relationship between two independent values, not a spelling check over a single file. A control test plants a boot-resident reference into a deferred section and proves the guard goes red.

**AC-5 (structural) — The greet reads only the README frontmatter, not the full body.**
Verified by: a structural check in `internal/contractlint` confirming the boot-resident Startup step-4 instruction targets the README frontmatter slice (entity-label, stage taxonomy, gate flags) and that the README-body-only anchors (proof-policy prose, the add-a-scenario procedure, the PR-body template, the task template) are reachable only from the deferred phase modules, not from the boot-resident core. The expected value (which README sections are greet-relevant vs phase-deferred) comes from the deferred-module manifest, an independent source the boot files can diverge from. A control test plants a boot-resident reference into a body-only section and proves the guard goes red. (Structural; pairs with AC-1, the behavioral proof the greet is accurate off the slim read.)

## Implementation pin (the one accuracy dependency to confirm)

**The greet must report a freshly-merged PR off the boot JSON without reading the pr-merge mod.** AC-1's accuracy and the deferral of M1 (the pr-merge startup hook) both rest on `status --boot --json`'s `pr_state.entries[].state` reflecting **LIVE** merge state, not just the stored `pr:` field. If it carried only the stored field, the greet could not report a freshly-merged PR accurately and the pr-merge `gh pr view` would have to stay before-greet.

**Status at this ideation:** CONFIRMED PRESENT in the current binary. `internal/status/boot.go` `checkPRStates` (line 114) runs `gh pr view {pr} --json state --jq .state` live for each PR-bearing non-terminal entity at boot, and `internal/status/json_commands.go` (lines 170-177) serializes that live `state` into `pr_state.entries[].state`. So the dependency holds today. The implementation must NOT regress it; pin it with a fixture adjacent to `internal/status/boot_probe_parity_test.go` (the existing parity test confines team-state probing, not pr_state — a new fixture is what this pin needs).

## Test plan

- **AC-1 (`shallow-boot` live scenario) — the costly item.** Follow the README's 4-step add-a-scenario procedure: host-neutral entry in `sharedRuntimeScenarios()`, fixture + prompt in `shared_fixtures_test.go`, host-neutral assertion (greet present + no-team + no-dispatch + gate-still-at-stage + stopped-for-input) and offline negative in `shared_scenarios_negative_test.go`, runner entries in BOTH `claudeScenarioRunners()` and `codexScenarioRunners()`. Cost: one model-spend scenario added to the serial suite (~5-7 min Claude opus). **Spot-check first:** run the parity/definition guards (`TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions`) at zero spend before paying for the live run.
- **AC-2 (no team re-cache before greet):** rides AC-1's live run — a transcript assertion over the tool-call sequence (no `TeamCreate` before the greet). Zero additional model spend.
- **AC-3 (regression):** zero new authoring — run the existing live scenarios after all three phases. Cost: the existing serial-suite wall-time, already budgeted in CI.
- **AC-4 + AC-5 (structural guards):** Go tests in `internal/contractlint`, no model spend, run in the offline gate job (`go test ./...`). Extend the existing reference-closure and structural-absence patterns. Each ships with a control test (planted violation goes red) so the guard is proven able to fail, not vacuous.
- **Implementation-time pin (see above):** before authoring AC-1, confirm `pr_state.entries[].state` still carries live merge state with a fixture adjacent to `boot_probe_parity_test.go`. Confirmed present at ideation; the test guards against regression.
- **High-stakes detached audit:** the FO's own Startup procedure + the runtime team-creation trigger are high-stakes — a detached adversarial audit before merge, on top of the live scenarios.
- **Fixture vs live:** AC-1, AC-2, AC-3 are live (the runtime integration — does the FO greet correctly with the team and modules deferred — is the claim). AC-4, AC-5 are structural fixtures over the shipped surface. No AC leans on a prose-grep over the contract.

## Spike status (both VIABLE; static analysis only)

Both source spikes are **static mechanism analysis** — dependency tracing, step enumeration, and correctness tracing over the contract/loader/mod files, with one binary-fact check of the boot JSON schema. This is the right depth for ideation. Neither proves the FO actually greets correctly with the team and modules unloaded; **the live `shallow-boot` drive (AC-1) is the implementation/validation proof, not ideation's.**

- **Phase 1 (contract split) — VIABLE with one wording tweak** (`T1-ideation.md` step 5). The shared-core Startup procedure is a clean boot island (C1); every coupling resolves without a boot-resident stub (C2-C6); the fold target (`using-claude-team`) is non-duplicative. Re-verified under fresh read: the cited test anchors `TestUserSkillReferenceClosureResolves` and `TestRetiredPluginPrivatePathsAbsent` exist in `internal/contractlint/structural_checks_test.go`; the regression scenarios exist in `internal/ensigncycle/shared_scenarios_test.go`.
- **Phases 2-3 (lazy-TeamCreate + shallow-boot) — VIABLE with adjustments A1-A5** (`j9-phases-2-3-ideation.md` step 6). No correctness guarantee is dropped by deferring past the greet; lazy-TeamCreate is one clause; the README read slims to its frontmatter. Re-verified under fresh read: the boot JSON omits labels/stage taxonomy (so the frontmatter read stays before-greet), and the live-pr_state implementation pin is CONFIRMED present in `internal/status/boot.go` + `json_commands.go` (above) — closing the single accuracy dependency the spike flagged.

No gap found under fresh read. Both verdicts hold.

## Stage Report: ideation

- DONE: Consolidate the 0203 roadmap-doc ideation (Phase-1 split in T1-ideation.md + Phases 2-3 in j9-phases-2-3-ideation.md) into this entity's body as the formal ideation deliverable: Problem, the three-phase Proposed approach, Out-of-scope, Acceptance criteria, Test plan — preserving both spike VIABLE verdicts.
  Body now carries Problem, Cost-levers table, the three sequential phases (P1 contract split + ethos/cross-ref cleanups, P2 lazy-TeamCreate one-clause, P3 shallow-boot sequence), Out-of-scope (T3/p2/vc/Codex-Pi/xp/ey), AC-1..AC-5, implementation pin, test plan; both spikes' VIABLE verdicts preserved in the Spike-status section.
- DONE: Verify every AC names EXTERNAL proof (a live shallow-boot scenario / the existing gate-guardrail+rejection-flow+merge-hook-guardrail regression / an internal/contractlint reference-closure guard) and NONE is a string-match over an instruction file; carry the live-PR-state implementation pin.
  AC-1/AC-2/AC-3 are live drives in internal/ensigncycle; AC-4/AC-5 are internal/contractlint reference-closure + structural-absence guards bound to the deferred-module manifest (independent source). No AC is a contract prose-grep. Implementation pin recorded with CONFIRMED-present status (boot.go checkPRStates runs live `gh pr view`; json_commands.go serializes the live state into pr_state.entries[].state).
- DONE: Confirm both spikes' VIABLE verdicts hold under a fresh read, or flag any gap — note these spikes are static analysis and the live shallow-boot drive is the implementation/validation proof, not ideation's.
  Both VIABLE verdicts hold; no gap. Fresh read confirmed: cited contractlint anchors and ensigncycle regression scenarios exist; boot JSON omits labels/stage taxonomy (frontmatter read stays before-greet); the live-pr_state pin is satisfied in the current binary. Spike-status section states explicitly that the static analysis establishes the sequence *can* be correct; the live shallow-boot drive (AC-1) is the implementation/validation proof.
- DONE: (coordinator-relayed Phase-1 scope) Fold the captain-audit additions — drop the agents/first-officer.md cross-reference; add a top Operating-principles (ethos) section folding Working Principles under it; mark the ethos as boot-resident GUIDANCE, not a testable AC.
  Both folded into Phase 1's Proposed approach with verbatim ethos text; AC section notes the ethos/cross-ref-drop is guidance proven by AC-3 + audit + review, with no tautological "drift reduced" metric.

### Summary

Consolidated two completed static-analysis spikes into the formal three-phase ideation deliverable for the 0.20.3 FO-efficiency backbone: P1 contract structural split (now including the captain-audit ethos section + cross-ref cleanup as boot-resident guidance), P2 one-clause lazy-TeamCreate (the ~89k lever), P3 shallow-boot-then-greet off `status --boot --json`. The single accuracy dependency the spikes flagged — that the boot JSON's `pr_state` reflects LIVE merge state so the greet can report a freshly-merged PR without reading the pr-merge mod — was verified present in the current binary (`internal/status/boot.go` runs `gh pr view` live), so the implementation pin is a regression guard, not an open risk. All five ACs are external proofs (live `shallow-boot` + existing regression scenarios + contractlint guards bound to an independent manifest); none is a contract prose-grep; both spike VIABLE verdicts hold under fresh read.
