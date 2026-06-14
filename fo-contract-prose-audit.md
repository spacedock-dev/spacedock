---
id: 95bcs48mr3jemtsb2zq7zbtb
title: FO contract residual-prose audit + comm-officer polish (post-split)
status: validation
source: 0.20.3 / 0203-fo-efficiency sprint (T3); captain 2026-06-13
started: 2026-06-13T18:09:33Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-fo-contract-prose-audit
issue:
sprint: 0203-fo-efficiency
sprint-readiness: ready
---

The T3 of the 0.20.3 (0203-fo-efficiency) sprint: after j9 Phase-1 splits the FO contract references into boot-resident vs deferred modules, sweep the slimmed refs for residual dead/redundant prose and run a comm-officer light-touch polish pass.

## Dependency

BLOCKED on j9 Phase-1 (the contract structural split). The audit cut-list does not exist until the structure is split — do NOT dispatch ideation until Phase-1 lands. Stays a backlog seed until then.

## Notes

Behavior-preserving CLEANUP: tightened/deduped prose in the already-split refs, NOT the structural split (j9 Phase 1) and NOT a behavioral contract change. At ideation, prove a real checkable change — the existing live FO scenarios stay green (behavior-preserving) plus a measurable size reduction of the refs — never a review of its own prose. If the split left no residual prose to cut, this collapses to a recorded roadmap decision, not a shipped task.

## Problem

j9's Phase-1 split reorganizes the FO contract *by when it is needed* — but it MOVES content, it does not rewrite it. After the split lands, three classes of residual prose remain in the slimmed surfaces, all invisible until the structure physically exists:

- **Dead prose** — sentences that only made sense in the pre-split single-file context. A forward/back cross-reference now pointing inside the same (now-deferred) module, a "see the dispatch section below" that no longer has a "below" because dispatch moved to a separate file, a transitional clause whose other half left.
- **Redundant prose** — the same obligation stated twice because the two statements used to sit far apart in one file and now sit adjacent (or one boot-resident, one deferred) after the carve. j9's C1/C5 splits explicitly duplicate a load-point reference across the boot core and the deferred module; the audit checks whether the boot-side and deferred-side phrasings of a shared obligation can be tightened to one canonical statement plus a pointer.
- **Unpolished prose** — the moved content carries the density and rhythm of a contract written to be read whole; once it is read in smaller modules at distinct moments, individual modules can be tightened for clarity/concision without changing meaning. This is the comm-officer light-touch pass.

None of these can be enumerated until the split exists. The cut-LIST is implementation; ideation designs the METHOD that produces it and the proof that the cut changed nothing observable while measurably shrinking the surface.

## Dependency and collapse case (the explicit determination point)

**Implementation is BLOCKED on j9 Phase-1.** Verified at this ideation: no `references/claude-fo-dispatch.md` and no merge reference exist on disk yet (`skills/first-officer/references/` holds only the four pre-split files), and the j9 entity is still at `status: ideation`. The slimmed boot-resident core and the two deferred modules — the modules this audit sweeps — do not exist until Phase-1 lands. Only the audit METHOD + ACs (this ideation) are doable now; the cut-LIST and the does-it-collapse determination need the split to physically exist.

**Collapse case, stated as an explicit fork the implementer hits FIRST.** The first implementation step is the survey pass (below). It produces a residual-prose inventory. That inventory forks the task:

- **Inventory non-empty (meaningful residual prose exists):** T3 proceeds as a code change — cut + comm-officer polish + the green-scenario / size-delta proof. The normal path.
- **Inventory empty or trivial (the split left clean modules — no dead cross-refs, no duplicated obligations, nothing worth a comm-officer pass beyond cosmetic):** T3 does NOT ship a code change. It collapses to a recorded roadmap decision — a one-paragraph note in `docs/roadmap/0203-fo-efficiency/README.md` (T3 section) stating "the Phase-1 split left no meaningful residual prose; audit closed with no cut," with the survey inventory attached as the evidence that the determination was made on the record, not skipped. The entity terminalizes as a decision, not a shipped change.

This fork is the begin-with-the-end discipline: the implementer runs the cheap survey first and lets it decide whether there is a task at all, rather than assuming a cut-list exists.

## Proposed approach (the audit METHOD — designed against j9's PLANNED module structure)

j9's Phase-1 produces a known module set (`docs/dev/.spacedock-state/lazy-teamcreate-shallow-boot/index.md`, Phase-1 section): the **slimmed boot-resident core** (`first-officer-shared-core.md` + the boot-resident runtime sections of `claude-first-officer-runtime.md`), a **deferred dispatch reference** (`references/claude-fo-dispatch.md` — Worker Resolution, Dispatch Adapter, Context Budget, Event Loop, Degraded Mode, standing-teammate mechanics, reuse-conditions/budget-probe), and a **deferred merge reference** (Merge-and-Cleanup, Ship-Local, worktree-removal safety, Mod-Block enforcement). The audit sweeps EACH resulting module. The method is a three-step sequence run per module, plus the comm-officer polish pass; it is deliberately mechanical so the cut-list is reproducible.

### Step 0 — Survey (cheap, decides whether a task exists)

For each of the three modules, produce a residual-prose inventory by reading the module IN ITS POST-SPLIT FORM (not the pre-split file) and recording each candidate against one of the three classes:

- **Dead-cross-reference scan.** For every "see X" / "above" / "below" / "the {section} section" pointer in the module, resolve its target. If the target moved to a different module by the split, the pointer is either dead (the target is unreachable from here at this lifecycle moment) or must become a load-point reference (name the deferred file, on the `present-gate`/`feedback-rejection-flow` precedent). Record each.
- **Duplicated-obligation scan.** j9's carve knowingly states some obligations on both sides of a boundary (C1 MODS-REPORT vs RUN-STARTUP-HOOKS; C5 gate-spine boot-resident vs reuse-conditions deferred). For each obligation that now appears in two modules, record whether the two phrasings can collapse to one canonical statement + a pointer, or whether both are genuinely needed at their two lifecycle moments (a deliberate, kept duplication — recorded as KEEP with rationale).
- **Orphaned-transition scan.** Sentences whose meaning depended on adjacent content that the split relocated (a "first … then …" whose second half left; a qualifier modifying a now-absent clause). Record each.

The union of these three inventories per module IS the cut-list. Empty union across all three modules ⇒ the collapse case above.

### Step 1 — Cut (mechanical, per inventory entry)

Apply each non-KEEP inventory entry as the smallest edit that removes the dead/redundant prose while preserving the obligation: delete a dead pointer, collapse a duplicated obligation to canonical+pointer, repair an orphaned transition. No new obligations, no reworded obligations beyond what removing the dead/redundant text requires. This is where behavior-preservation is on the line — every cut must leave the obligation intact, which AC-1/AC-2 (the green live scenarios) police.

### Step 2 — comm-officer light-touch polish (the polish vehicle)

Route each slimmed module through the workflow's standing `comm-officer` teammate (`docs/dev/_mods/comm-officer.md`) in **file-in-place** mode (`polish this file {absolute_path}`) for a clarity/concision pass. comm-officer's contract is exactly the right tool here: its default is "light-touch — preserve the caller's voice, rhythm, and technical vocabulary; cut empty words, tighten sentences, fix clear grammar; do NOT rewrite for style." It preserves disambiguating parentheticals and semantic qualifiers (the load-bearing nuance in a contract), flags anything bigger than light-touch for human review rather than making it, and applies the `elements-of-style:writing-clearly-and-concisely` skill. The polish is best-effort and non-blocking (its hard rule): if comm-officer is unavailable, the cut from Step 1 still ships and the stage report notes the fallback. Polish NEVER changes an obligation — comm-officer's "do not change semantic qualifiers silently" rule plus AC-1/AC-2's green scenarios are the two guards that keep the polish behavior-preserving.

### How the method honors the simplify-the-contract ethos

j9's ethos ("begin with the end; do the hardest things first, de-risk when it is cheap; choose the simplest approach") governs the method: Step 0's survey is the cheap de-risk that decides whether the task exists at all (begin with the end); the per-module mechanical sweep is the simplest approach that makes the cut-list reproducible rather than a judgment call; comm-officer is reused (the simplest approach — the polish vehicle already ships) rather than building new tooling.

## Out of scope

- **The structural split itself (j9 Phase-1).** T3 sweeps the ALREADY-split modules; it does not move content between files. Moving content is j9's job.
- **The cut-LIST and the does-it-collapse determination.** These are IMPLEMENTATION, BLOCKED on j9 Phase-1 — they need the split to physically exist. Ideation produces the METHOD and ACs only; it does not enumerate specific prose to cut, because the target files do not exist yet.
- **Any behavioral contract change.** This is a behavior-preserving cleanup. If the survey surfaces a clause whose MEANING should change, that is a separate behavioral task, not T3 — record it and route it, do not fold it into the cleanup.
- **The Codex and Pi runtime adapters' residual-prose audit.** j9 Phase-1 splits the Claude adapter (the bulk file) + the shared core; if Phase-1 later extends to Codex/Pi, their audit rides that, not this task. T3 scopes to the modules Phase-1 actually produces.
- **New token-counting or size-measurement tooling.** The size-delta proof uses byte/line counts of the ref files (a `wc`-shaped measurement) and the existing live-boot read sequence — no new binary facility.

## Acceptance criteria

Each AC names an end-state property of the finished task, verified by something OUTSIDE this task body that can fail. **No AC is proven by a string/substring/regex match over the prose of any FO contract reference** — that is banned by this workflow and by `internal/contractlint/boundary_guard_test.go` (a passing prose-grep only asserts the implementer's own text is present; a valid paraphrase fails it and an inverted clause passes it). The behavior-preserving ACs (AC-1, AC-2) are the existing live FO scenarios staying green. The size-reduction AC (AC-3) is a byte/line measurement of the ref files — a `wc`-shaped count, which the boundary guard explicitly admits as the legal "line-floor / size" structural category because the expected value (a smaller count than the pre-audit baseline) comes from the filesystem's measurement, NOT from reading the prose for meaning. AC-4 is the collapse-case decision record (only on the collapse fork).

**AC-1 (behavioral, regression, live) — The audit is behavior-preserving: every existing live FO scenario stays green after the cut + the comm-officer polish.**
Verified by: `gate-guardrail`, `rejection-flow`, `feedback-3-cycle-escalation`, and `merge-hook-guardrail` in `internal/ensigncycle` all PASS against the post-audit refs, run `go test -tags live -count=1 -run TestLiveClaudeSharedScenarios ./internal/ensigncycle`. These four scenarios exercise the boot-resident gate spine (`gate-guardrail`), the deferred dispatch reference loading at first dispatch (`rejection-flow`, `feedback-3-cycle-escalation` — reuse-conditions, feedback routing, the reconcile sweep), and the deferred merge reference loading at terminalization (`merge-hook-guardrail` — mod-block enforcement). A green run is the independent proof that no cut removed a reachable obligation and no polish silently changed one — the scenarios fail if a guardrail obligation went missing or inverted, regardless of the prose wording. (Live, real host, durable on-disk-state assertions — never a contract grep.)

**AC-2 (behavioral, regression, cross-host) — The behavior preservation holds on Codex too: the same shared scenarios pass under the Codex runner.**
Verified by: `TestLiveCodexSharedScenarios ./internal/ensigncycle` passing — the shared-core cuts (the host-neutral half of the audit) are exercised by both hosts via the per-host runners, so a cut that broke an obligation only on one host is caught. The zero-spend parity guards (`TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions`) are run FIRST, before paying for either live suite, to confirm the scenario set still covers both hosts. (Live, real host; rides the existing scenario set — no new scenario authored, because behavior-PRESERVATION is the claim and the existing scenarios already cover the affected guardrails.)

**AC-3 (structural, size measurement) — The audited refs are measurably smaller than the pre-audit baseline.**
Verified by: a byte/line count (`wc -c` / `wc -l`, or the equivalent committed in the stage report) of the audited ref files AFTER the audit being strictly less than the same count taken on the post-j9-Phase-1 / pre-T3-audit baseline, with the per-file deltas recorded in the stage report. The expected value (a smaller count) comes from the filesystem measurement, an independent oracle the prose cannot fake — a paraphrase that preserves meaning but does not shrink the file fails this AC, and an inverted clause that happens to be shorter does NOT pass AC-1, so the two ACs together cannot both be satisfied by a meaning-changing edit. This is the boundary guard's admitted "size floor" category (a measurement of the file, not a match over its prose-for-meaning), NOT a prose-grep. The baseline is captured at the start of implementation (a `wc` snapshot of the slimmed modules immediately after Phase-1 lands, before any T3 cut) and recorded in the entity body so the delta is checkable, not asserted.

**AC-4 (collapse fork only) — If the split left no meaningful residual prose, the determination is recorded as a roadmap decision, not silently skipped.**
Verified by: on the collapse fork (Step 0's inventory empty/trivial), a committed diff to `docs/roadmap/0203-fo-efficiency/README.md` (the T3 section) recording "Phase-1 split left no meaningful residual prose; audit closed with no cut," with the survey inventory attached. The independent oracle is the committed roadmap-doc change (a file the change produces, outside this task body) — the git log shows the decision was recorded. On the non-collapse fork this AC is N/A (AC-1..AC-3 carry the proof). This AC exists so the collapse case cannot be used to silently abandon the task: even "nothing to cut" must leave a checkable on-disk artifact.

**On the comm-officer polish (Step 2):** the polish itself carries no AC of its own. Its quality is not measurable by a "polish applied" or "clarity improved" metric — such a check would be either a banned prose-grep over the refs or a tautology. The polish ships INSIDE the cut, and its behavior-preservation is proven by AC-1/AC-2 (green scenarios) while its size contribution rolls into AC-3 (the byte/line delta). comm-officer's own "flag for review rather than make a big change" rule plus the high-stakes detached audit at the gate are the qualitative guards; there is no standalone polish AC.

## Test plan

- **AC-1 (behavior-preserving, Claude live) — zero new authoring.** Run the four existing live scenarios after the cut + polish: `go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle -v`. Cost: the existing serial-suite wall-time (~27m opus), already budgeted in CI; no new scenario added because behavior-PRESERVATION is the claim and the existing scenarios already cover every affected guardrail (boot spine, deferred dispatch, deferred merge). **Spot-check first:** run the zero-spend parity/definition guards (`TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions`) before paying for the live run.
- **AC-2 (behavior-preserving, Codex live) — zero new authoring.** `go test -tags live -count=1 -timeout 40m -run TestLiveCodexSharedScenarios ./internal/ensigncycle -v`. Catches a shared-core cut that broke an obligation only on Codex. Cost: the existing Codex serial-suite wall-time.
- **AC-3 (size measurement) — no model spend.** Capture a `wc -c`/`wc -l` baseline of the slimmed modules immediately after j9 Phase-1 lands (before any T3 cut), record it in the entity body, then re-measure after the audit and record the per-file delta in the stage report. The check is the after-count being strictly less than the baseline — a filesystem measurement, run with `wc`, no binary facility needed.
- **AC-4 (collapse fork only) — no model spend.** If Step 0's inventory is empty/trivial, the proof is the committed roadmap-doc diff; the git log is the oracle. N/A on the non-collapse fork.
- **High-stakes detached audit:** the FO's own contract refs are high-stakes; a detached adversarial audit of the cut + polish before merge, on top of the live scenarios, confirms no obligation was dropped or inverted by a cut that the scenarios happen not to exercise. This is the read-shape backstop the boundary guard names.
- **Fixture vs live:** AC-1/AC-2 are live (behavior-preservation is the claim — the runtime must still obey every contract obligation after the cut). AC-3 is an offline filesystem measurement. AC-4 is an offline committed-diff record. No AC leans on a prose-grep over the contract.

## Spike determination

**No spike needed.** The method composes only already-proven mechanisms: (1) the live shared-scenario harness is an existing, exercised facility (`internal/ensigncycle`, run today in CI for `gate-guardrail`/`rejection-flow`/`merge-hook-guardrail`/`feedback-3-cycle-escalation`); (2) comm-officer is an existing standing teammate with a proven file-in-place polish mode (`docs/dev/_mods/comm-officer.md`); (3) the size measurement is a plain `wc`. The one thing the method cannot exercise until j9 Phase-1 lands — the actual cut-list — is IMPLEMENTATION, not a mechanism risk: the mechanisms (green-scenario proof, comm-officer polish, byte/line delta) are all proven independently of which specific prose gets cut. The riskiest unknown for THIS task is "does the split leave meaningful residual prose?", and that is answered by Step 0's cheap survey at the START of implementation (the begin-with-the-end / pay-the-small-bill-first discipline), recorded as the collapse fork — not by a throwaway spike now against files that do not exist.

## Stage Report: ideation

- DONE: Design the audit METHOD against j9's planned module structure (boot-resident core + lazy dispatch/merge refs): how to sweep each resulting module for dead/redundant/duplicative prose, plus the comm-officer light-touch polish pass — the cut-LIST itself is implementation (post-Phase-1), not ideation.
  Proposed-approach section defines a per-module four-step method (Step 0 survey → three named scans producing the inventory/cut-list; Step 1 mechanical cut; Step 2 comm-officer file-in-place polish) run against j9's known Phase-1 module set (slimmed boot core, deferred `claude-fo-dispatch.md`, deferred merge ref). The cut-list is explicitly named as the Step-0 OUTPUT (implementation), not enumerated at ideation.
- DONE: Acceptance criteria proven by EXTERNAL evidence: behavior-preserving (the existing live FO scenarios stay green after the polish) + a measurable size reduction of the refs (byte/token delta) — never a review of T3's own prose.
  AC-1/AC-2 are the existing live shared scenarios (`gate-guardrail`/`rejection-flow`/`feedback-3-cycle-escalation`/`merge-hook-guardrail`) staying green on Claude + Codex — the independent behavior-preservation oracle. AC-3 is a `wc -c`/`wc -l` byte/line delta vs a recorded post-Phase-1 baseline (the boundary guard's admitted "size floor" category, filesystem oracle, not a prose-grep). The no-AC-for-the-polish note explains why a "polish quality" metric would be a banned tautology.
- DONE: Record the dependency + collapse case: implementation is BLOCKED on j9 Phase-1 (the split must land first); if the split leaves no meaningful residual prose, T3 ships as a recorded roadmap decision, not a code change — state that determination point explicitly.
  Dependency section records BLOCKED-on-j9-Phase-1, VERIFIED at this ideation (no `claude-fo-dispatch.md`/merge ref on disk; j9 still `status: ideation`). The collapse case is stated as an explicit FORK the implementer hits first (Step 0 survey decides): non-empty inventory ⇒ code change; empty/trivial ⇒ AC-4 records a roadmap-doc decision with the inventory attached as on-disk evidence, terminalizing as a decision not a shipped change.

### Summary

Designed the residual-prose-audit METHOD against j9's PLANNED Phase-1 module set (slimmed boot-resident core + deferred dispatch ref + deferred merge ref): a reproducible per-module four-step sweep — Step 0 survey (dead-cross-reference / duplicated-obligation / orphaned-transition scans whose union IS the cut-list and decides the collapse fork), Step 1 mechanical cut, Step 2 comm-officer light-touch file-in-place polish (the workflow's existing standing polish teammate, voice-preserving and non-blocking). ACs are all external: AC-1/AC-2 = the existing live FO scenarios staying green on Claude + Codex (behavior-preservation oracle), AC-3 = a `wc` byte/line delta vs a recorded post-Phase-1 baseline (the boundary guard's legal size-floor category), AC-4 = a committed roadmap-doc decision on the collapse fork; none is a prose-grep over the contract, and the comm-officer polish carries no standalone AC by design. Implementation is BLOCKED on j9 Phase-1 (verified: the deferred refs do not exist yet); the collapse-to-roadmap-decision case is the explicit Step-0 fork. No spike needed — the method composes only already-proven mechanisms (the live scenario harness, comm-officer's file-in-place mode, plain `wc`).

## Stage Report: implementation

- DONE: Step-0 survey produces a residual-prose inventory across j9's four merged FO-ref modules (the three scans: dead-cross-reference / duplicated-obligation / orphaned-transition) — the inventory's union IS the cut-list and decides the cut-vs-collapse fork.
  Inventory NON-EMPTY (4 cuts, all in `first-officer-shared-core.md`) ⇒ CUT fork. Dead-cross-reference/orphaned-transition: 3x "S7b" pointers (the merged-PR sweep is step 8 in the flat 1-9 Startup list — no "S7b" label exists; residue of a pre-split S*/sub-step numbering scheme) + 1x orphaned "(B6)" label on "Rebase-conflict halt" (referenced nowhere in the repo, verified by repo-wide grep: 1 hit = the label itself). Duplicated-obligation: KEEP — concurrency-safe commit rule (canonical in shared-core §State Management, pointer-only from dispatch/merge) and worktree-ownership (canonical in dispatch, pointer in shared-core:147) are already canonical+pointer. FLAG (separate task, NOT cut here): merge module's two mod-block sections (`## Mod-Block Enforcement` + `## ... at Terminal Transitions`) restate the same mechanism-level invariant + session-resume rule; j9 added the first adjacent to the pre-existing second — collapsing two `##` sections is a judgment-call restructure that risks the live `merge-hook-guardrail` scenario, out of scope for a behavior-preserving mechanical cut.
- DONE: On the CUT fork: each non-KEEP entry applied as the smallest behavior-preserving edit + a comm-officer light-touch file-in-place polish pass, with a `wc -c`/`wc -l` size delta recorded per file vs a baseline captured on the merged refs BEFORE any edit (AC-3).
  4 cuts applied (Step 1, the smallest edit each: 3 S7b pointers retargeted to the live step title "the Merged-PR sweep below" / "sweep-advanced"; "(B6)" dropped from the heading — obligations untouched). All four modules then run through comm-officer file-in-place light-touch polish (Step 2); polish preserved every MUST/NOT and the verbatim `TERMINAL_TEARDOWN_BOUNDED:` marker. comm-officer flagged ONE candidate semantic change — it had lowercased the emphasis "NOT" in the FO Write Scope off-limits clause ("the workflow README ... so it is NOT in this list"); I restored the caps-NOT (load-bearing emphasis distinguishing the FO-owned README from off-limits scaffolding, byte-neutral casing swap). Code committed on worktree branch `spacedock-ensign/fo-contract-prose-audit` @ `41b5858f`. AC-3 delta below.
- DONE: The offline gate stays green — `go test ./...` exits 0 (the cut must not break `internal/contractlint` reference-closure or the marker/hook-consistency tests); the live behavior-preservation proof (AC-1/AC-2) is left to validation — do not burn live runs.
  `go test ./...` = 1385 passed / 0 fail (exit 0). Targeted: `internal/contractlint` 23 passed (reference-closure + boundary guard); marker/hook-consistency 18 passed across 13 packages; zero-spend parity guards `TestSharedScenarioRunnerCoverage|TestSharedRuntimeScenarioDefinitions` passed (scenario set still covers both hosts). AC-1/AC-2 (live Claude+Codex shared scenarios) intentionally deferred to validation — no live runs burned.

### AC-3 size delta (post-Phase-1 baseline → after cut+polish)

Baseline captured on the merged refs (HEAD `fe4261be`) BEFORE any edit; final after cut + comm-officer polish. `wc -c`/`wc -l`:

| module | baseline (lines / bytes) | final (lines / bytes) | Δ bytes |
|---|---|---|---|
| first-officer-shared-core.md | 231 / 27377 | 231 / 27253 | -124 |
| claude-first-officer-runtime.md | 43 / 5304 | 43 / 5225 | -79 |
| claude-fo-dispatch.md | 241 / 31930 | 241 / 31791 | -139 |
| claude-fo-merge.md | 85 / 10877 | 85 / 10839 | -38 |
| **total** | **600 / 75488** | **600 / 75108** | **-380** |

All four files strictly smaller (AC-3 satisfied). Lines unchanged — the polish was within-line concision, the cut a label/pointer repair, neither removed lines.

### Summary

CUT fork (Step-0 inventory non-empty). The survey found 4 dead cross-references in `first-officer-shared-core.md`, all residue of a pre-split S*/B* step-numbering scheme the flat 1-9 Startup list left dangling (3x "S7b" → step 8 Merged-PR sweep; 1x orphaned "(B6)" label). Each repaired as the smallest behavior-preserving edit, then all four modules polished by comm-officer file-in-place (light-touch, every obligation + the verbatim teardown marker preserved). AC-3 delta -380 bytes total, all four files strictly smaller. Offline gate green (1385 pass / 0 fail; contractlint + marker/parity guards pass). AC-1/AC-2 live behavior-preservation left to validation. Flagged for a SEPARATE task (out of scope here): the merge module's two adjacent mod-block-enforcement sections restate one obligation — a section-collapse is a judgment-call restructure, not a behavior-preserving mechanical cut.
