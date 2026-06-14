---
id: 5ew2jxagk11mr0fzd0rtpdp0
title: In-module restructure of the FO contract refs (consolidate duplicated obligations + collapse redundant sections)
status: ideation
source: "T3 (fo-contract-prose-audit) deferred this (2026-06-14): T3 shipped the mechanical-safe subset (4 dead-ref cuts + comm-officer concision); the substantive restructure was scoped out (duplicated obligations marked KEEP; the merge mod-block section-collapse FLAGGED out-of-scope). Captain: the audit \"would imagine a bigger cleanup and in-module restructure.\""
started: 2026-06-14T18:42:01Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0203-fo-efficiency
---

The substantive in-module cleanup of the post-j9 FO contract refs that T3 deferred. T3 did the safe mechanical subset (dead-ref repairs) and a low-value within-line concision polish (~75 of 79 changed lines, which over-reached on meaning 3× and cost two amend cycles + a rejection). The real value — consolidating genuinely-duplicated obligations and collapsing redundant sections within the modules — was explicitly punted as "judgment-call restructure, not a behavior-preserving mechanical cut."

## Problem statement

The post-j9 FO contract refs accreted duplicated obligation prose across cross-cycle edits. T3 (#367, merged) shipped the mechanically-safe subset and surveyed the duplications but marked them all KEEP. This task does the judgment-call restructure T3 deferred, collapsing the genuine redundancy without dropping or inverting a single obligation.

## Re-scope against the actual tree (post-#367)

The dispatch checklist asked: enumerate exactly which collapses/consolidations genuinely remain vs already-canonical, given #367 landed part of the work. Verified against `origin/main` (#367 = `f87107b1`, merged `81e423ee`):

**#367 did NOT collapse anything structural.** Its 28-line merge-ref diff was within-line concision only — both mod-block sections survive fully (`claude-fo-merge.md:54` and `:64`). The section-collapse and the obligation-consolidation are genuinely outstanding.

**Genuinely remains (in scope):**

1. **Collapse the merge ref's two mod-block sections.** `skills/first-officer/references/claude-fo-merge.md` carries `## Mod-Block Enforcement` (lines 54–62) AND `## Mod-Block Enforcement at Terminal Transitions` (lines 64–85). Both restate the same mechanism-level invariant (`status --set`/`--archive` refuse terminal transitions when merge hooks registered AND `pr` empty AND `mod-block` empty, `--force` bypasses) plus the session-resume rule (j9 added the first adjacent to the pre-existing second). Collapse to ONE canonical section carrying: the set/clear/guard bullets, the mechanism-level enforcement bullet with its `merge: local` / `verdict=rejected` exemption carve-outs, the recovery options (set+invoke, let-hook-set-`pr`, `--force`-with-the-do-NOT-force-to-clear-the-guard warning), the session-resume scan, and the missing-mod-file recovery. Behavior-preserving. (T3's explicit FLAGGED-out-of-scope item.)

2. **Collapse the C1 MODS-REPORT restatement to a pointer.** `first-officer-shared-core.md:27` is the canonical boot-greet operational instruction (the MODS bullet under the `--boot --json` read: what the greet reads and that reading the map opens no mod file). `first-officer-shared-core.md:199` restates the same fact ("The MODS-REPORT at boot reads the boot JSON `mods` map … without opening a mod file") inside the deferred-mods conceptual section. Collapse line 199's restatement to a pointer back to the canonical greet bullet; keep only the line-199 content that is NOT at the greet site (the deferred-mod-block-travels-with-merge-module fact). Behavior-preserving.

3. **Relocate the FO operating-principles ethos up to the skill entry point.** (Captain-added.) Move the high-level ethos out of `first-officer-shared-core.md` (the `## Operating principles (ethos)` block, lines 5–12 on origin/main — confirmed anchor) UP into the always-first-loaded entry point `skills/first-officer/SKILL.md`, which is currently pure mechanism (single-entity rules → `## Operating contract` @-include → runtime adapter → Startup). Two captain requirements drive the shape: (a) the ethos frames the FO from the entry point, not buried in a reference; (b) it reads self-explanatory — DROP the `(ethos)` label and the self-referential meta-framing ("These principles govern…", "the Working Principles below fold under them"), stating the principles plainly as how the FO operates.

   **Seam decision (recommended): move only the high-level ethos; the detailed `## Working Principles` posture stays in shared-core.** The two are distinct: the lines 5–12 block is the high-level "what awesome looks like" framing (three bullets); the lines 209–220 `## Working Principles` block is the detailed FO posture (prefer-code-gate, name-end-value, lead-with-recommendation, reversible-work, speak-declared-label). The detailed posture stays in the boot-resident core because it references mechanism the core owns — the `entity-label` read at Startup step 4, the gate's AC cross-check, `status` guards — and belongs adjacent to that mechanism, not in entry-point framing. Hauling it up would bloat the always-loaded entry point and split the declared-label rule from the Startup-step-4 read it depends on. The captain's requirement is that the *ethos frames* from the entry point, satisfied by moving the high-level block; the habits that implement it stay where their mechanism lives.

   **No duplication results — one canonical home each.** The three ethos bullets land in SKILL.md and are REMOVED from shared-core (not copied). The forward-pointer at shared-core:12 ("The Working Principles below fold under them") is deleted with the block. `## Working Principles`'s own self-reference at shared-core:211 ("These habits govern how the FO frames work and adjudicates gates") is rewritten to point UP to the entry-point ethos, so the seam reads coherently without re-stating the three bullets.

   **Before/after wording.**

   SKILL.md — INSERT after the single-entity-mode block (current line 14, before `## Operating contract`):
   ```
   ## How the first officer operates

   You are dispatcher, responsible for making sure the work is done by the crew. What awesome looks like:
   - Begin with the end; be clear about the value.
   - Do the hardest things first; de-risk while it is cheap.
   - Communicate and act concisely, choose the simplest approach, JFDI.
   ```
   (Label dropped, meta-framing dropped, stated plainly as operating posture.)

   shared-core — DELETE the entire `## Operating principles (ethos)` block (lines 5–12): the `## Operating principles (ethos)` header, the "You are dispatcher…" lede, the three bullets, and the "These principles govern… fold under them." closer.

   shared-core:211 — REWRITE the `## Working Principles` lede from "These habits govern how the FO frames work and adjudicates gates." to point up to the entry-point ethos, e.g. "These habits implement the operating posture stated in the skill entry point — they are how the FO frames work and adjudicates gates in practice." (exact wording is the implementer's, but it MUST reference the entry-point ethos, not silently orphan the seam).

**Already canonical — do NOT touch (the rest of T3's KEEP survey, confirmed correct):**

- **Concurrency-safe-commit:** canonical in `first-officer-shared-core.md:153–156` (State Management). `claude-first-officer-runtime.md:43` and `claude-fo-merge.md:30/34/35` are genuine *pointers* ("per the shared core's State Management rule" / "commit path-scoped") — already canonical + pointer. No collapse.
- **Worktree-ownership:** single canonical site `claude-fo-dispatch.md:105–113`. No duplication. No collapse.
- **C5 reuse-conditions / gate-spine:** canonical in the deferred dispatch module; `first-officer-shared-core.md:125/127/135` reference it ("per the dispatch module's reuse conditions"). Already canonical + pointer. No collapse.
- **RUN-STARTUP-HOOKS:** stated once at the canonical greet site (`:27`); no second site. No collapse.

So the in-module collapse scope is exactly TWO (item 1 + item 2), not the open-ended "re-examine each duplication" the seed implied; item 3 is the captain-added entry-point ethos relocation. The seed's original "module-level coherence re-org" is dropped: the re-scope found no incoherent accretion beyond these, and an open-ended re-org without a concrete target invites the meaning-drift that sank T3's polish.

## Why this is risky (and how to prove it)

A section-collapse or obligation-consolidation can DROP or INVERT an obligation that the live scenarios don't exercise — exactly the class the detached adversarial audit exists for (it caught T3's dropped NEVER-qualifier).

The riskiest unknown is the mod-block collapse: it touches the exact invariant `merge-hook-guardrail` grades, and it carries the `TERMINAL_TEARDOWN_BOUNDED` verbatim marker. **No spike needed** — the mechanism this rests on is already proven: the live `merge-hook-guardrail` scenario (`internal/ensigncycle/shared_scenarios_test.go:39`, intent "FO cannot bypass a registered merge hook by terminalizing without pr, mod-block, or force") plus the detached word-level baseline diff together exercise the only path a bad collapse could break. The validation order pays the small bill first: run the detached diff and `merge-hook-guardrail` before the full scenario sweep.

The ethos relocation (item 3) carries the same drop/invert risk one level milder: it is a move, not a collapse, so the words should survive verbatim — the detached audit (AC2) confirms each ethos bullet lands intact in SKILL.md, and the live scenarios (AC3) confirm the FO still behaves under the relocated framing. The one thing the relocation could silently break is the seam: orphaning the `## Working Principles` self-reference (line 211) or leaving a stale forward-pointer at the old ethos site. AC5's non-duplication check catches a half-done move; the seam-rewrite is called out in the before/after wording so the implementer doesn't leave a dangling reference.

## Acceptance criteria

- **AC1 (structural — the collapse landed).** `claude-fo-merge.md` carries exactly ONE mod-block section (`## Mod-Block Enforcement at Terminal Transitions` removed; its non-redundant content folded into `## Mod-Block Enforcement`), and `first-officer-shared-core.md:199`'s MODS-REPORT restatement is a pointer.
  *Verified by:* a header-count assertion in the validation diff (the collapsed file has one `## Mod-Block` header, down from two) — an on-disk fact outside the task body, checkable by `grep -c '^## Mod-Block' skills/first-officer/references/claude-fo-merge.md` returning `1`.

- **AC2 (high-stakes detached audit — no obligation lost).** A word-level diff of every collapsed/consolidated/relocated obligation against the pre-restructure baseline (the files at #367's merged tree) confirms no MUST / MUST-NOT / qualifier (NEVER, only, unless, except) dropped or inverted across all three changes, the `TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher.` marker is byte-intact, AND every operating-principle bullet (the three ethos lines) survives verbatim in its new SKILL.md home with no principle dropped or weakened in the move.
  *Verified by:* the detached audit's diff output (an independent reviewer / `git diff {#367-tree} -- claude-fo-merge.md SKILL.md first-officer-shared-core.md` whose every removed MUST-bearing clause and every removed ethos bullet is shown to survive verbatim in its destination). The baseline is an independent source (the prior committed tree), not the task's own prose.

- **AC5 (structural — ethos has exactly one canonical home).** The three operating-principle bullets ("Begin with the end…", "Do the hardest things first…", "Communicate and act concisely… JFDI") are present in `SKILL.md` and ABSENT from `first-officer-shared-core.md` — relocated, not duplicated. The `(ethos)` label and the self-referential meta-framing ("These principles govern…", "fold under them") are gone from shared-core.
  *Verified by:* a non-duplication on-disk fact — the bullet text appears in exactly one of the two files. Checkable e.g. by `grep -l 'Begin with the end' skills/first-officer/SKILL.md skills/first-officer/references/first-officer-shared-core.md` returning only the SKILL.md path. (This is a structural count, not a prose-quality check.)

- **AC3 (behavioral, live).** The existing live shared scenarios (`gate-guardrail` / `rejection-flow` / `feedback-3-cycle-escalation` / `merge-hook-guardrail`, Claude + Codex) stay green after the restructure. The mod-block collapse specifically keeps `merge-hook-guardrail` green.
  *Verified by:* the live runner exit code — `go test ./internal/ensigncycle/...` (the Claude + Codex live runners) passing.

- **AC4 (structural gate).** `internal/contractlint` reference-closure + the offline gate stay green.
  *Verified by:* `go test ./...` exit code 0 (includes `internal/contractlint`).

## Test plan

| AC | What verifies it | Kind | Cost |
|----|------------------|------|------|
| AC1 | `grep -c '^## Mod-Block' claude-fo-merge.md == 1`; pointer present at shared-core MODS-report restatement site | on-disk grep | seconds |
| AC2 | Detached word-level diff vs #367 merged tree (`git show f87107b1:…` baseline) across all three files; MUST/MUST-NOT/qualifier survival + marker byte-check + every ethos bullet survives in SKILL.md. Run FIRST (riskiest path). | detached adversarial review | minutes |
| AC3 | `go test ./internal/ensigncycle/...` — live Claude + Codex shared scenarios (the FO still behaves under the relocated framing) | live workflow test | many minutes |
| AC4 | `go test ./...` (offline gate incl. contractlint reference-closure) | Go unit/structural | minutes |
| AC5 | `grep -l 'Begin with the end' SKILL.md first-officer-shared-core.md` → SKILL.md only (ethos has one canonical home) | on-disk grep | seconds |

No fixture authoring needed — every AC binds to an existing test or an on-disk fact against the prior committed tree. The text changes themselves are NOT acceptance criteria; the proofs are the live runner, the offline gate, the header-count fact, the non-duplication fact, and the detached baseline diff.

**Self-explanatory wording is a gate-review property, not an AC.** Requirement (b) — the relocated ethos reads self-explanatory without the "(ethos)" label or meta-framing — is a wording-quality judgment, not a behavioral property a check can falsify. A prose-grep over SKILL.md ("the bullets are present") is banned as proof (the same trap that sank a literal-string AC). It is flagged for the human gate review instead: the reviewer confirms the relocated block frames the FO plainly and needs no label. AC3 (live scenarios green) is what proves the relocation preserved *behavior*; AC5 proves it created no duplication.

Implementation runs WITHOUT a worktree at the ideation level (this is a contract-ref edit, no deliverable branch); the actual edit lands at the implementation stage per the workflow's stage flags.

## Doc-diff note

No user-visible surface changes (CLI output, banners, docs-site copy). These are FO contract-ref and skill-entry instruction files the model reads (`SKILL.md` included), not docs the site describes — no doc diff required.

## Out of scope

- The team-vs-bare dispatch-mode determination (separate task `7e` / `headless-dispatch-mode-intent`).
- A comm-officer concision polish (T3 did it; it is low-value + meaning-change-risky on the contract — do NOT repeat it here; if comm-officer is used at all, harden its guard per the xf note first).
- Module-level coherence re-org (seed item 3): dropped — the re-scope found no incoherent accretion beyond the two collapses, and an open-ended re-org without a concrete target re-invites meaning drift.
- The already-canonical duplications (concurrency-commit, worktree-ownership, C5 reuse-conditions, RUN-STARTUP-HOOKS): confirmed canonical+pointer, NOT collapse candidates.

## Notes

Fast-follow, not a v0.20.3 blocker (T3's behavior-preservation shipped). A `comm-officer` polish-over-reach guard (it changed contract meaning 3× under "light-touch") should be folded into `xf` (which moves comm-officer usage prose into its mod) — a hard "never touch MUST/MUST-NOT/qualifiers in contract prose" rule.

## Stage Report: ideation

- DONE: Re-scope against the ACTUAL tree — enumerate which section-collapses and obligation-consolidations genuinely remain vs already-canonical.
  Verified #367 (`f87107b1`, merged `81e423ee`) did within-line concision only; both mod-block sections survive (`claude-fo-merge.md:54`+`:64`). Scope narrowed to exactly two collapses; concurrency-commit / worktree-ownership / C5 / RUN-STARTUP-HOOKS confirmed already canonical+pointer, dropped from scope.
- DONE: High-stakes detached-audit AC — word-level diff of every collapsed/consolidated obligation vs pre-restructure baseline (no MUST/MUST-NOT/qualifier dropped or inverted; TERMINAL_TEARDOWN_BOUNDED marker byte-intact).
  Recorded as AC2, bound to an independent source (the #367-merged tree via `git show f87107b1:…`), with the marker byte-check and qualifier-survival check named.
- DONE: Behavioral AC — live shared scenarios stay green; mod-block collapse keeps merge-hook-guardrail green; contractlint + offline gate green.
  Recorded as AC3/AC4; scenario names confirmed present in `internal/ensigncycle/shared_scenarios_test.go` (merge-hook-guardrail at :39, intent matches the mod-block invariant). "No spike needed" recorded — rests on the proven merge-hook-guardrail scenario + detached baseline diff.

### Summary

Re-scoped the seed's open-ended "re-examine each duplication" down to exactly two behavior-preserving collapses against the live tree: the merge ref's two mod-block sections (genuinely outstanding — #367 left both intact) and the C1 MODS-REPORT restatement→pointer. The other four KEEP-survey duplications were verified already-canonical+pointer and dropped, as was the open-ended module-coherence re-org (no concrete incoherence found; an untargeted re-org re-invites the meaning drift that sank T3's polish). ACs are entity-level with proofs bound to independent sources — the live runner, the offline gate, an on-disk header-count fact, and the #367-tree baseline diff — never the task's own prose.

## Stage Report: ideation (cycle 2)

- DONE: Incorporate captain feedback — add the entry-point ethos relocation as scope item 3 (additive; kept the two existing collapses).
  Scope item 3 added with the seam decision, no-duplication design, and concrete before/after wording for SKILL.md insert, shared-core delete, and the shared-core:211 self-reference rewrite.
- DONE: Make the design decision the captain left open (what moves: just the ethos vs also the detailed posture) and present it concretely.
  Recommended moving ONLY the high-level ethos (lines 5–12); the detailed `## Working Principles` posture stays in shared-core, justified by the seam (it references Startup-step-4 `entity-label`, the gate AC cross-check, `status` guards — mechanism the boot-resident core owns).
- DONE: Ensure no duplication; one canonical home + pointer if needed.
  Ethos REMOVED from shared-core (not copied); forward-pointer at :12 deleted with the block; `## Working Principles` lede at :211 rewritten to point UP to the entry-point ethos. AC5 enforces the single-home fact.
- DONE: AC guidance — bind relocation proof to live scenarios + a structural non-duplication fact; do NOT invent a tautological prose-grep for "self-explanatory".
  AC3 (live scenarios green) proves behavior-preservation under the relocated framing; AC5 (grep returns SKILL.md only) proves single canonical home. "Self-explanatory" flagged as a gate-review wording-quality property, explicitly NOT an AC (prose-grep banned).
- DONE: Keep the detached-audit AC, extended to the ethos move; re-confirm origin/main anchors.
  AC2 now diffs all three files vs the #367 tree and asserts every ethos bullet survives verbatim in SKILL.md. Anchors re-confirmed against origin/main: ethos block at shared-core lines 5–12; SKILL.md has no ethos (so the move is clean, no pre-existing collision).

### Summary

Folded the captain's additive scope: relocate the FO operating-principles ethos from `first-officer-shared-core.md` up into the always-first-loaded `SKILL.md`, dropping the "(ethos)" label and self-referential meta-framing so it reads as plain operating posture. Made the open design call — move only the high-level three-bullet ethos; the detailed `## Working Principles` posture stays in the core adjacent to the mechanism it references — and gave concrete before/after wording plus the shared-core:211 seam-rewrite so no dangling reference is left. Proofs stay behavioral: AC3 (live scenarios green under the relocated framing) + AC5 (non-duplication grep, SKILL.md is the single home) + the extended detached audit (AC2); the "self-explanatory" requirement is a gate-review property, deliberately not a tautological prose-grep AC.
