# T1 — FO contract structural split + lazy-load

**Milestone:** 0.20.3 (0203 FO efficiency)
**Stage:** ideation
**Blocks:** T2 (shallow-boot-then-greet); T3 files along after.

Restructure the first-officer contract so a boot reads only what it needs to greet, report state, and present a gate. Extract the dispatch/team material and the merge material into already-lazy load points, leaving a slim boot-resident core. This is the structural split that T2's boot-flow reshape and T3's prose audit both depend on — neither cut-list exists until the split lands.

## Problem

Boot forensics on a live FO session (`/tmp/boot-analysis-spacedock-v1.md`) measured the cost of reaching an interactive greet **with no team created and no worker dispatched** — a 100% pre-dispatch session:

- Peak context **160,594 tokens** (event 148); the session crossed 100k at event 97 (~363s) and never fell back.
- Wall-clock to the captain-facing gate question **~511s (~8.5 min)**; to the tail ~13.6 min. The wall-clock is dominated by generation latency that grows with loaded context — the two slowest think-turns (128.6s, 100.1s) both fired *above 97k context*.

The single biggest avoidable structural waste in that picture is the two FO reference files, **read whole, back-to-back, at events 23 and 25 (t≈8–9s), immediately after the skill loaded, before any work began:**

| Forensics rank | Tokens | File |
|---|---|---|
| #1 | ~9,310 (37,249 ch) | `first-officer-shared-core.md` — single biggest ingest |
| #3 | ~6,900 (27,602 ch) | `claude-first-officer-runtime.md` — read back-to-back with #1 |

That is **~16,200 tokens** of contract read on every boot. The forensics call out the runtime adapter specifically: "gate `claude-first-officer-runtime.md` behind 'am I creating a team this turn' — which never happened this session, so ~6.9k was pure waste." The runtime adapter is ~70% team/dispatch material that a boot which never dispatches never uses, yet it loads first.

The cost is structural, not a bug: the loader reads the entire contract up front. Generation latency scales with loaded context, so trimming the boot read trims the wall-clock directly.

## Spike (riskiest-first) — is the two-tier split structurally sound?

**The riskiest unknown:** does the boot-resident path have a hard dependency on deferred content that would break if that content is not loaded at boot? This is static mechanism analysis — dependency tracing over the contract sections plus confirming the lazy hooks already exist. It is the right depth for ideation. The *live* FO-boot drive that proves the FO still greets and presents a gate correctly with the modules unloaded is the implementation/validation proof (the `shallow-boot` scenario in AC-1), not this spike.

### Spike step 1 — every top-level section mapped to a tier

Token estimates are char/4 over the section body. The two files total **~9,014 + ~6,700 ≈ 15,700 tokens**, matching the forensics' ~16.2k.

**`first-officer-shared-core.md` (17 sections, ~9,014 tok)**

| Section | Lines | ~tok | Tier |
|---|--:|--:|---|
| Startup | 24 | 963 | **boot-resident** |
| Status Viewer | 39 | 655 | **boot-resident** |
| ID Styles | 10 | 256 | **boot-resident** |
| Single-Entity Mode | 12 | 188 | **boot-resident** |
| Working Directory | 4 | 44 | **boot-resident** |
| Dispatch | 24 | 746 | **dispatch-deferred** |
| Completion and Gates | 44 | 1340 | **split** — gate-decision spine boot-resident; reuse-conditions dispatch-deferred (see step 5) |
| Merge and Cleanup (+ Ship-Local, Worktree-removal) | 49 | 1785 | **merge-deferred** |
| State Management | 6 | 66 | **boot-resident** (rebase-conflict halt referenced by Startup) |
| Worktree Ownership (+ Split-Root Worktree Contract) | 30 | 750 | **dispatch-deferred** |
| FO Write Scope | 15 | 430 | **boot-resident** |
| Mod Hook Convention (+ Mod-Block Enforcement) | 21 | 390 | **merge-deferred** (mod-block is a merge-ceremony concept) |
| Standing Teammates | 9 | 402 | **dispatch-deferred** (folds behind `using-claude-team`) |
| Clarification and Communication | 6 | 115 | **boot-resident** |
| Working Principles | 13 | 524 | **boot-resident** |
| Probe and Ideation Discipline | 7 | 308 | **boot-resident** |
| Issue Filing | 4 | 19 | **boot-resident** |

**`claude-first-officer-runtime.md` (11 sections, ~6,700 tok)**

| Section | Lines | ~tok | Tier |
|---|--:|--:|---|
| Team Creation (+ standing-teammate discovery/lazy-spawn/declaration) | 44 | 1622 | **dispatch-deferred** → `using-claude-team` |
| Worker Resolution | 10 | 153 | **dispatch-deferred** |
| Dispatch Adapter (+ break-glass) | 58 | 1941 | **dispatch-deferred** |
| Degraded Mode (spacedock seams) | 7 | 159 | **dispatch-deferred** → `using-claude-team` |
| Context Budget and Dead Ensign Handling | 23 | 586 | **dispatch-deferred** |
| Captain Interaction | 16 | 468 | **split** — gate-guardrail + greet boot-resident; team-mode chat hint + single-entity gate-resolution dispatch-deferred |
| Feedback Rejection Flow (bare mode) | 6 | 86 | **dispatch-deferred** (the routing skill `feedback-rejection-flow` is already lazy) |
| Event Loop (incl. reconcile sweep step 0) | 24 | 988 | **dispatch-deferred** |
| Mod-Block Enforcement at Terminal Transitions | 23 | 503 | **merge-deferred** |
| Agent Back-off | 6 | 84 | **boot-resident** (cheap; captain-interaction adjacent) |
| Entity-Body Inspection | 4 | 82 | **boot-resident** (points at shared-core Probe discipline) |

**Tier totals (approximate):**

- **Boot-resident:** shared-core ~3,800 tok (Startup, Status Viewer, ID Styles, Single-Entity Mode, Working Directory, State Management, FO Write Scope, Clarification, Working Principles, Probe discipline, Issue Filing, plus the gate-decision spine of Completion-and-Gates) + runtime ~700 tok (Captain Interaction greet/guardrail, Agent Back-off, Entity-Body Inspection) ≈ **~4,500 tok.**
- **Dispatch-deferred:** ~6,000 tok (the bulk of the runtime adapter + shared-core Dispatch/Worktree-Ownership/Standing-Teammates/reuse-conditions).
- **Merge-deferred:** ~2,700 tok (Merge-and-Cleanup, Ship-Local, Mod-Hook/Mod-Block, runtime Mod-Block-at-Terminal).

**Boot read drops from ~15,700 tok to ~4,500 tok — roughly an 11k-token cut on every boot, ~70% of the contract-read cost.** (Net win is slightly less than the gross because the deferred modules re-load when a real dispatch/merge happens — but on a no-dispatch boot, the session the forensics measured, the full ~11k is saved.)

### Spike step 2 — loader claim confirmed

`skills/first-officer/SKILL.md` instructs reading **both** reference files at startup:

- Line 18: `@references/first-officer-shared-core.md` — the `@`-directive inlines the shared core into the skill body at load.
- Lines 23–25: "Load the runtime adapter for your platform: … read `references/claude-first-officer-runtime.md`" (Claude branch).
- Line 27: "Then begin the Startup procedure from the shared core."

The live forensics timeline corroborates this exactly: row 1 = `Skill spacedock:first-officer` (event 13, t=4s), row 2 = Read `first-officer-shared-core.md` (event 23, t=8s), row 3 = Read `claude-first-officer-runtime.md` (event 25, t=9s). Both files are read within one second of each other, at the very top of the session. `agents/first-officer.md` adds nothing — it only delegates to the skill ("invoke the `spacedock:first-officer` skill now to load it. Then begin the Startup procedure"). The skill is the single loader.

### Spike step 3 — coupling trace (boot-resident → deferred)

The crux question: does any boot-resident step *depend on* team/dispatch/merge knowledge before the greet? Findings, with the resolution for each:

**C1 — shared-core Startup (the boot procedure itself) has ZERO reference into deferred content.** Startup steps 1–7 are: contract-version gate, `git rev-parse` root, `status --discover`, README read, `status --boot`, split-root halt-gate, split-root pull-on-boot. None mentions team, reconcile, dispatch, or merge. Its only forward reference is "follow the rebase-conflict halt in **State Management**" — and State Management is itself boot-resident. **This is the clean island that makes the split viable.** The boot procedure stands alone.

**C2 — the runtime adapter "Team Creation" says "At startup (after reading the README, before dispatch)".** This is the load-bearing phrase. Read literally, "at startup" couples team creation into boot. But its own next sentence reframes it: "Invoke it before the first team-mode tool call in the session." And it already delegates entirely to the lazy skill: `Skill(skill="spacedock:using-claude-team")`. **Resolution:** retire "at startup" in favor of the truthful trigger — team creation fires at *first dispatch*, which is when `using-claude-team` is meant to load. This is a one-clause wording change, not a structural break. It aligns the contract with what the forensics already shows happening (the measured session never created a team because it never dispatched).

**C3 — the reconcile sweep (Event Loop step 0) runs "(a) at boot, AFTER the split-root pull --rebase and BEFORE the first dispatch".** This is the genuine boot-adjacent step inside otherwise-deferred content. It needs a `team_name` (the A/B/C drift classes are roster-derived). **It does NOT need to run before the greet** — it runs "before the first dispatch," and the shallow-boot flow (T2) greets before any dispatch. So reconcile rides into the dispatch-deferred Event Loop module that loads at first dispatch, alongside the `team_name` it requires. There is no boot-resident step that calls reconcile; the greet happens off `status --boot --json` alone (forensics row 8 confirms `--boot --json` is already run at t=27s, before all the heavy reads). **No stub needed in the boot-resident core.**

**C4 — standing-teammate discovery pass ("after team creation … BEFORE entering the normal dispatch event loop").** This is wholly inside the team-creation flow; it lazy-spawns at first dispatch already ("No spawn calls at boot. Spawn is deferred to the first team-mode dispatch"). It travels with the dispatch module. No boot coupling.

**C5 — Completion-and-Gates "decide reuse-or-fresh" references the reuse-conditions and the runtime context-budget probe.** The gate-decision spine (never self-approve, present-gate, AC cross-check, the gated-stage branch) is boot-resident — a shallow boot must be able to present a gate. But the reuse-conditions block and the budget probe are only reached *after a worker completes*, which cannot happen before a dispatch. **Resolution:** split the section — the gate-presentation/AC-cross-check spine stays boot-resident; the reuse-or-fresh machinery (reuse conditions 0–4, the model-mismatch diagnostic, SendMessage advancement, supersede-shutdown) moves to the dispatch module. `present-gate` and `feedback-rejection-flow` are already lazy skills the boot-resident spine invokes by name — that precedent is exactly the shape.

**C6 — Mod-Block Enforcement (shared-core) and Mod-Block at Terminal (runtime) are referenced from Merge-and-Cleanup only.** Both are merge-ceremony concepts. They travel with the merge module. The boot-resident core needs to know merge hooks *exist* (the MODS section of `status --boot` reports them), but not the enforcement mechanics — those are read at terminalization. No boot coupling.

**Summary:** the only boot-resident step that touches deferred concepts is C2's "at startup" wording, and that is a wording-truth fix, not a structural dependency. C3 (reconcile) is boot-adjacent but fires before-first-dispatch, not before-greet, so it lives cleanly in the dispatch module. No boot-resident step genuinely needs team/dispatch/merge knowledge before the greet.

### Spike step 4 — is the fold target non-duplicative?

`spacedock:using-claude-team` already carries the generic team lifecycle: Deferred Team Tools (the ToolSearch hop), Team Creation (TeamCreate-first sequencing, naming, bare-mode fallback), the TeamCreate recovery procedure, the failure-recovery ladder, Degraded Mode (triggers/effects/captain-report/shutdown-sweep), Awaiting Completion, and Terminal Team Teardown. The runtime adapter's "Team Creation" section **already invokes this skill** (`Skill(skill="spacedock:using-claude-team")`) and explicitly states "the generic blocks they reference (`## Degraded Mode`, `## Awaiting Completion`, `## Terminal Team Teardown`) live in that skill, not in this file."

So the split is already partly done: the *generic* team lifecycle is in the lazy skill; the *spacedock-specific* adapter sections (Worker Resolution, Dispatch Adapter, standing-teammate discovery/lazy-spawn, the spacedock Degraded-Mode seams, Context Budget, Event Loop/reconcile) are what still load eagerly in the runtime adapter. **Folding the runtime adapter's team/dispatch sections behind the same lazy load point as `using-claude-team` is non-duplicative** — they are the spacedock specializations that the generic skill leaves to the consumer. The clean mechanism: a new lazily-loaded reference (e.g. `references/claude-fo-dispatch.md`) holding the spacedock dispatch/team sections, read at the same first-dispatch moment `using-claude-team` is invoked. No content overlaps with the generic skill; the generic skill keeps the lifecycle, the deferred reference keeps the spacedock adapter.

### Spike step 5 — VERDICT

**Viable with one wording tweak.** The clean two-tier split is structurally sound:

- The shared-core Startup procedure is a self-contained boot island with no forward dependency into deferred content (C1) — this is what makes the cut clean.
- The greet runs off `status --boot --json` alone, already executed early in the live session (forensics row 8), so nothing team/dispatch/merge is needed before the greet.
- Every coupling resolves without a boot-resident stub: reconcile (C3) and standing-teammates (C4) fire before-first-dispatch (not before-greet) and travel with the dispatch module; mod-block (C6) travels with the merge module; the gate spine splits cleanly from reuse-or-fresh (C5) along the present-gate/feedback-rejection precedent.
- **The one required tweak:** retire the runtime adapter's "At startup … before dispatch" framing for team creation (C2) in favor of "at first dispatch," matching the lazy `using-claude-team` invocation that is already the team's load point. One clause, not a redesign.

No boot-path step genuinely needs team/dispatch/merge knowledge before the greet. The split is not blocked.

## Proposed approach

The section→module assignment is the table in spike step 1. The mechanism:

1. **Slim the boot-resident core.** `first-officer-shared-core.md` keeps its boot-resident sections (Startup, Status Viewer, ID Styles, Single-Entity Mode, Working Directory, State Management, FO Write Scope, Clarification, Working Principles, Probe discipline, Issue Filing) plus the gate-presentation spine extracted from Completion-and-Gates. The runtime adapter keeps Captain Interaction's greet/guardrail, Agent Back-off, and Entity-Body Inspection.

2. **Extract the dispatch/team module.** Move the runtime adapter's Worker Resolution, Dispatch Adapter, Context Budget, Event Loop (incl. reconcile), Degraded-Mode seams, standing-teammate discovery/lazy-spawn/declaration, and the shared-core Dispatch / Worktree-Ownership / Standing-Teammates / reuse-conditions into a lazily-loaded reference read at first dispatch. The generic team lifecycle stays in `using-claude-team` (already lazy); this new reference holds only the spacedock adapter. The first-dispatch load point is the existing `Skill(skill="spacedock:using-claude-team")` invocation — extend it to also pull the spacedock dispatch reference. **C2 fix:** the team-creation trigger reads "at first dispatch," not "at startup."

3. **Extract the merge module.** Move Merge-and-Cleanup, Ship-Local Ceremony, Worktree-removal safety, Mod-Hook/Mod-Block Enforcement (shared-core), and Mod-Block-at-Terminal (runtime) into a lazily-loaded reference read at terminalization. The boot-resident core reaches it the same way it reaches `present-gate` / `feedback-rejection-flow`: by naming the load point at the terminal boundary.

4. **Loader change.** `SKILL.md` line 18's `@references/first-officer-shared-core.md` inlines only the slimmed core. The runtime-adapter read (lines 23–25) loads only the boot-resident runtime sections. The dispatch and merge references are NOT read at startup — they load at their phase boundary, extending the existing lazy pattern (`present-gate`, `feedback-rejection-flow`, `using-claude-team`). No new mechanism.

**Coupling resolutions (from step 3):** C1 — no change needed (clean island). C2 — wording fix to "at first dispatch." C3/C4 — reconcile + standing-teammates travel with the dispatch module (before-first-dispatch, not before-greet). C5 — split Completion-and-Gates: gate spine boot-resident, reuse-or-fresh deferred. C6 — mod-block travels with the merge module.

This task is the structural move only. It is a content reorganization that must be **behavior-preserving** — the same instructions, reachable at the same moments, just loaded lazily. The regression scenarios (AC-2) are the proof of no behavior change; they are the load-bearing safety net for "this is a move, not a rewrite."

## Out of scope

- **The j9 boot-flow change (T2).** Greet-first sequencing, deferred mod-reads, deferred status-table render — the *new* shallow-boot behavior — is T2, blocked-by this task. T1 makes the modules loadable lazily; T2 changes when the FO greets. The `shallow-boot` live scenario is authored against T1's structure but its greet-first behavior is T2's claim.
- **p2 / vc binary-simplification** (`spacedock pr complete`, `reconcile --act`) — parked to 0.20.4.
- **The residual-prose audit + comm-officer polish (T3)** — the cut-list of leftover prose to trim does not exist until this split lands; T3 files along after.
- **The Codex and Pi runtime adapters.** This task splits the Claude adapter (the bulk file the forensics measured). The shared-core split applies cross-host, but the per-host dispatch/merge reference extraction for `codex-first-officer-runtime.md` (~4.9K) and `pi-first-officer-runtime.md` (~5.9K) is a follow-on if their boot cost warrants it — they are an order smaller than the Claude adapter.

## Acceptance criteria

Each AC names an end-state property of the finished split, verified by something outside this task body that can fail. No AC is proven by a string/substring/regex match over an instruction file the model reads — that is banned by this workflow (a passing match only asserts the implementer's own text is present).

**AC-1 — With the dispatch and merge modules NOT loaded, a freshly-booted FO greets the captain, reports workflow state, and presents a human gate without self-approving, mutating, or archiving the entity.**
Verified by: a new live shared-runtime scenario `shallow-boot` in `internal/ensigncycle` (added to `sharedRuntimeScenarios()` with Claude + Codex runners per the README's add-a-scenario procedure). The runner launches the real host front door against a fixture sitting at a human gate, and the host-neutral assertion over `(before, after, observed)` confirms the FO presented the gate (durable: entity still at the gate stage, not archived, no `verdict`/`completed` set; final message carries a gate review + decision prompt) — the gate-guardrail behavior surviving with the deferred modules unloaded. Behavioral and live, not a contract grep. An offline negative case in `shared_scenarios_negative_test.go` builds the broken end-state (entity self-approved/archived) and proves the assertion goes red.

**AC-2 — The split is behavior-preserving: the deferred modules load and function correctly when a real dispatch or merge happens.**
Verified by: the existing live scenarios `gate-guardrail`, `rejection-flow`, and `merge-hook-guardrail` in `internal/ensigncycle` still pass after the split. `gate-guardrail`/`rejection-flow` exercise the dispatch module loading at first dispatch (reuse-conditions, feedback routing); `merge-hook-guardrail` exercises the merge module loading at terminalization (mod-block enforcement). A green run of all three is the proof that lazy-loading did not drop a reachable instruction — run via `go test -tags live -run TestLiveClaudeSharedScenarios ./internal/ensigncycle`.

**AC-3 — The boot-resident core has no reference dependency on deferred-only content: every reference the boot-resident core makes resolves either within the boot-resident set or to a known lazy load point, never into the body of the dispatch/merge modules.**
Verified by: a new reference-closure structural guard in `internal/contractlint` (the allowed quarantine), extending `TestUserSkillReferenceClosureResolves`. The check parses the boot-resident files and the deferred-module manifest as real artifacts, builds the set of section anchors the deferred modules own, and fails if a boot-resident `@`/read reference resolves into a deferred-module section that is not one of the declared lazy load points (`using-claude-team`, `present-gate`, `feedback-rejection-flow`, the new dispatch/merge references). This tests a relationship between two independent values — the boot-resident reference set vs. the deferred-module section set — which can diverge (a future edit that points the boot core at a moved section makes them disagree), so it can fail; it is not a spelling check over a single file. A control test plants a boot-resident reference into a deferred section and proves the guard goes red.

**AC-4 — The boot read shrinks: the files loaded at startup (the slimmed core + boot-resident runtime sections) no longer carry the dispatch/team or merge sections.**
Verified by: a structural-absence check in `internal/contractlint` confirming the dispatch-module section anchors (Dispatch Adapter, Worker Resolution, Event Loop, Context Budget, standing-teammate spawn) and merge-module anchors (Merge and Cleanup, Ship-Local Ceremony, Mod-Block Enforcement) are ABSENT from the boot-resident files and PRESENT in their deferred references — the same structural-absence shape as `TestRetiredPluginPrivatePathsAbsent`. The expected value (which anchors belong where) comes from the deferred-module manifest, an independent source the boot files can diverge from. This is a structural location check, not a behavioral claim; the behavioral win (lower loaded context) is AC-1's live scenario. The check fails if a deferred section is left behind in the boot core or a boot section is wrongly moved out.

## Test plan

- **AC-1 (`shallow-boot` live scenario):** the costly item. Following the README's 4-step add-a-scenario procedure — host-neutral entry in `sharedRuntimeScenarios()`, fixture + prompt in `shared_fixtures_test.go`, host-neutral assertion + offline negative in `shared_scenarios_negative_test.go`, runner entries in BOTH `claudeScenarioRunners()` and `codexScenarioRunners()`. Live, real host, durable-state assertion. Cost: one model-spend scenario added to the serial suite (~5–7 min Claude opus); the offline negative and the parity/definition guards run at zero spend. **Spot-check first:** run the parity/definition guards (`TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions`) before paying for the live run.
- **AC-2 (regression):** zero new authoring — run the three existing live scenarios after the split. Cost: the existing serial-suite wall-time, already budgeted in CI.
- **AC-3 + AC-4 (structural guards):** Go tests in `internal/contractlint`, no model spend, run in the offline gate job (`go test ./...`). Cost: low — they extend existing reference-closure and structural-absence patterns. Each ships with a control test (planted violation goes red) so the guard is proven able to fail, not vacuous.
- **Fixture vs live:** AC-1 is live (the runtime integration — does the FO behave with modules unloaded — is the claim). AC-2 is live (same reason). AC-3/AC-4 are structural fixtures over the shipped surface. No AC leans on a prose-grep over the contract.

## Spike result

**Verdict: viable with one wording tweak.** The two-tier split is structurally sound. Evidence:

- **Loader claim confirmed** (step 2): `SKILL.md` reads both reference files at startup (`@references/first-officer-shared-core.md` + the runtime-adapter read), corroborated by the live forensics timeline (Skill → shared-core Read → runtime Read within one second, t=8–9s).
- **Boot island is clean** (C1): shared-core Startup steps 1–7 have zero forward reference into team/dispatch/merge content; their only forward reference (rebase-conflict halt) targets the boot-resident State Management section.
- **Greet needs no deferred knowledge:** the greet runs off `status --boot --json`, already executed at t=27s in the live session before any heavy read.
- **Every coupling resolves** (step 3): C2 (one "at startup"→"at first dispatch" clause), C3/C4 (reconcile + standing-teammates fire before-first-dispatch, travel with the dispatch module), C5 (gate spine splits from reuse-or-fresh on the present-gate precedent), C6 (mod-block travels with the merge module). No boot-resident step needs team/dispatch/merge knowledge before the greet, so no boot-resident stub is forced.
- **Fold target is non-duplicative** (step 4): `using-claude-team` already carries the generic team lifecycle and the runtime adapter already invokes it; the spacedock-specific adapter sections fold behind the same first-dispatch load point without overlapping the generic skill.

**Boot read drops ~15,700 → ~4,500 tokens (~11k cut, ~70% of contract-read cost) on a no-dispatch boot** — the exact ~16.2k the forensics flagged as read-whole-up-front, of which ~6.9k (the runtime adapter) was pure waste in the measured session.

**Honesty on spike depth:** this is static mechanism analysis (dependency tracing + confirming the lazy hooks exist), the right depth for ideation. It does NOT prove the FO still behaves correctly with the modules unloaded — that is AC-1's live `shallow-boot` drive at implementation/validation. The spike establishes the split *can* be clean; the live scenario establishes that it *is*.
