---
id: 5ew2jxagk11mr0fzd0rtpdp0
title: In-module restructure of the FO contract refs (consolidate duplicated obligations + collapse redundant sections)
status: validation
source: "T3 (fo-contract-prose-audit) deferred this (2026-06-14): T3 shipped the mechanical-safe subset (4 dead-ref cuts + comm-officer concision); the substantive restructure was scoped out (duplicated obligations marked KEEP; the merge mod-block section-collapse FLAGGED out-of-scope). Captain: the audit \"would imagine a bigger cleanup and in-module restructure.\""
started: 2026-06-14T18:42:01Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-fo-refs-inmodule-restructure
issue:
sprint: 0203-fo-efficiency
mod-block: merge:pr-merge
pr: "#373"
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

## Principle-cascade map (item 4 — the deep goal, planned not executed)

The captain's deeper goal: once the operating principles are stated authoritatively and self-explanatory at the entry point, deep contract prose that merely RESTATES what a principle already implies can be COMPRESSED ("inferred from the principles"), fighting contract bloat. This is the highest-risk edit in the contract, so ideation produces a **map**, not blind cuts. **Decision rule: default to KEEP. When derivability is uncertain, keep and mark "kept — ambiguous derivability." The audit's bias is against silent loss.**

**Honest headline finding:** the compressible surface is SMALL — and reporting that plainly is the valid finding the lead asked for, not a failure. The FO contract is already overwhelmingly load-bearing mechanism (guards, markers, sequencing, MUST/MUST-NOT the live scenarios grade), not posture restatement. The posture that *is* prose lives in two spots: `## Clarification and Communication` (C-a) and the `## Working Principles` `**FO posture:**` bullets (C-b). Once the ethos is hoisted, those restatements compress — but the teeth they carry (reversible-vs-hard-to-reverse, the gate lede-first spine) either move to a single canonical home or already live in `present-gate`. **Estimated delta: ~40–45 words trimmed total (C-a ~12, C-b ~30), roughly 8–12 lines of the ~230-line shared-core, i.e. a low-single-digit-percent reduction.** The cascade's real payoff is not the line count — it is that the entry-point ethos becomes the cited source so future contract prose can lean on it instead of re-explaining posture. The bloat-fighting win is real but bounded; the deep prose is mostly load-bearing mechanism that the principles cannot replace.

### principle-derivable → COMPRESS

- **C-a: `## Clarification and Communication` (shared-core:205–207), the "don't ask to take an allowed step" prose.**
  Derives from: **"Do obvious reversible work without ceremony"** (the hoisted posture's reversible-work principle) + **"Communicate and act concisely, JFDI"** (entry-point ethos bullet 3).
  Before (:207, first clause): *"Do not ask whether to take a step this contract already allows — proceed. If one entity is blocked on clarification, keep dispatching other ready entities. Report workflow state once on idle or at a gate; do not repeat status updates while waiting."*
  After: *"Don't ask permission for a step the contract already allows (the reversible-work principle); keep dispatching other ready entities when one blocks. Report state once on idle or at a gate, not repeatedly while waiting."*
  Compression: folds the "don't ask to take an allowed step" sentence into a one-clause pointer to the principle. The keep-dispatching-other-entities rule and the report-once rule are operational specifics (NOT derivable) and stay. NET small — one sentence of restatement compressed, the two operational rules retained.

- **C-b: the `## Working Principles` `**FO posture:**` bullets (shared-core:215–219, the bullet block = 104 words / 672 chars).** Per the lead's steer, examined harder against BOTH the hoisted entry-point ethos AND the `present-gate` skill (where the gate-time teeth actually live). The bullets split three ways:

  - **name-end-value (:217, 31 words) → COMPRESS.** Derives from ethos bullet 1 ("Begin with the end; be clear about the value"). The only non-ethos content is the rationale "end-value framing is judgeable; step-framing has to be reverse-engineered" — a justification, not a mechanism. Before: *"**Name the end value before starting.** State the outcome — the change in the world the captain gets — before mechanism. End-value framing is judgeable; step-framing has to be reverse-engineered."* After: *"**Name the end value before starting** (entry-point principle 1) — state the outcome before mechanism; end-value framing is judgeable, step-framing is not."* Saves ~12 words.

  - **lead-with-recommendation (:218, 34 words) → COMPRESS.** This is NOT in the terse ethos — but its operational teeth ALREADY live in the `present-gate` skill, which this task does not touch: present-gate rule 1 ("Lede first, decision last… the first three lines are title, chosen direction, recommend… if the captain stops reading after line three they can still vote") IS "approvable in a single yes," and the skill's nine-rule set names "single recommendation" explicitly. So the shared-core bullet is a middle-layer restatement of a rule present-gate enforces at render time. Before: *"**Lead with a recommendation the captain can say yes to.** Open with one clear recommended direction approvable in a single 'yes,' then supply detail. Do not bury under a menu of equally-weighted options."* After: *"**Lead with a recommendation the captain can say yes to** — one recommended direction, not a menu; the gate rendering enforces the lede-first spine (see `present-gate`)."* Saves ~10 words and removes the duplicated "single yes" teeth (now sourced from present-gate, its real home).

  - **reversible-work (:219, 37 words) → COMPRESS, but only the restatement half.** Derives from ethos bullet 3 ("Communicate and act concisely… JFDI"). The reversible/hard-to-reverse distinction is operational teeth the ethos lacks AND is the same teeth C-a leans on — so it should have ONE home. Recommend: keep the reversible-vs-hard-to-reverse distinction as the canonical statement HERE (C-a points to it), compress only the example padding. Before: *"**Do obvious reversible work without ceremony.** Obvious reversible steps (a dispatch the contract already allows, a status read, a routine state transition) just happen. Reserve asking for choices that are hard to reverse or genuinely matter."* After: *"**Do obvious reversible work without ceremony** (entry-point principle 3) — reversible steps the contract allows just happen; reserve asking for choices that are hard to reverse or genuinely matter."* Saves ~8 words.

  - **speak-label (:220, 84 words) → KEEP VERBATIM.** Pure mechanism — references the `entity-label` read at Startup step 4 and enumerates exactly what localizes (human-facing noun) vs stays generic (`entity_path`, machine output). Not principle-derivable. Stays.

  - **prefer-code-gate (:213, 120 words) → KEEP VERBATIM.** The load-bearing prose-vs-code-gate rule the AC cross-check enforces ("an AC of the form 'the contract says X' is satisfied only by 'the binary or a test enforces X'"). This is the anti-tautology spine of the whole proof discipline — not posture, and not derivable from the terse ethos. Stays.

  Net for C-b: ~30 words trimmed across three bullets, two bullets kept verbatim. The earlier "kept — ambiguous" verdict is now RESOLVED: present-gate carries lead-with-recommendation's teeth, so that bullet safely compresses; the other two compress only restatement/padding, keeping their distinct teeth.

### load-bearing mechanism → KEEP VERBATIM (cannot be inferred from any principle)

These stay regardless of how "obvious" they look — each is graded by a live scenario or a test, or is a byte-exact marker / guard / sequence:

- **Never self-approve a gate** (shared-core:130; runtime:19, with the single-entity-mode auto-resolve exception). Graded by `gate-guardrail`. A principle ("the captain decides") does NOT encode the absolute "infer-approval-from-silence is forbidden" guard.
- **`TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher.`** (merge:24) — byte-exact marker a watcher greps. Not derivable; KEEP byte-intact (already AC2).
- **Mod-block set→invoke→clear sequence + the combine-with-terminal-fields refusal** (merge:9–19, 54–85). Graded by `merge-hook-guardrail`. The set/clear ordering and the standalone-clear rule are mechanism, not posture.
- **Split-root rebase-conflict HALT** (shared-core:166–169) — HALT + `rebase --abort` + escalate, never `--force`. A specific recovery sequence; not derivable.
- **ID-style / `--next-id` / `spacedock new` atomic-create rules** (shared-core:175–188; runtime:43). Drift-avoidance mechanism.
- **Dispatch numbered procedure** (dispatch:62–76) — the `status --set ... status={next} worktree=... started` shape, commit-then-create-worktree ordering, ≤3-item checklist cap. Mechanism + exact CLI.
- **AC coverage cross-check at every gate** (shared-core:121). A specific gate obligation; the prefer-code-gate principle motivates it but does not encode "scan `## Acceptance criteria`, name any AC without evidence, REJECT if this stage was the natural place."
- **"Completed non-gated non-terminal stage is not a stopping point"** (shared-core:125) — the MUST-advance rule + the enumerated legitimate halts. Graded behaviorally; the enumerated halt conditions are not principle-derivable.
- **Context-budget probe / dead-ensign handling** (dispatch:190–211), **Reuse conditions** (dispatch:82–104), **Standing-teammate lifecycle** (dispatch:17–57). Mechanism + sequencing.
- **Contract-version gate** (shared-core:18–22), **Status Viewer / Captain-Facing State Display** (shared-core:41–79), **Mod Hook Convention** (shared-core:190–201). Mechanism.

### Staging recommendation

The COMPRESS set is C-a (one `## Clarification and Communication` compression) plus C-b (three `## Working Principles` bullets compressed, two kept verbatim) — ~40–45 words / ~8–12 lines total. **One implementation pass suffices** — the set is small and all in shared-core, well under the size that would warrant staged passes. The pass pays the full audit price: live-green (AC3) + detached word-level diff (AC2) proving no MUST/MUST-NOT/qualifier dropped or inverted, covering C-a and all three C-b bullet rewrites. The audit's specific job for C-b: confirm name-end-value / lead-with-recommendation / reversible-work lost only restatement, that lead-with-recommendation's "single-yes" teeth are genuinely still enforced by `present-gate` (not silently dropped), and that speak-label and prefer-code-gate are byte-untouched. Were a future effort to attempt a wider sweep beyond these two sections, THAT would stage per-section with its own audit; this task does not open that door — the deep mechanism prose is KEEP.

## Why this is risky (and how to prove it)

A section-collapse or obligation-consolidation can DROP or INVERT an obligation that the live scenarios don't exercise — exactly the class the detached adversarial audit exists for (it caught T3's dropped NEVER-qualifier).

The riskiest unknown is the mod-block collapse: it touches the exact invariant `merge-hook-guardrail` grades, and it carries the `TERMINAL_TEARDOWN_BOUNDED` verbatim marker. **No spike needed** — the mechanism this rests on is already proven: the live `merge-hook-guardrail` scenario (`internal/ensigncycle/shared_scenarios_test.go:39`, intent "FO cannot bypass a registered merge hook by terminalizing without pr, mod-block, or force") plus the detached word-level baseline diff together exercise the only path a bad collapse could break. The validation order pays the small bill first: run the detached diff and `merge-hook-guardrail` before the full scenario sweep.

The ethos relocation (item 3) carries the same drop/invert risk one level milder: it is a move, not a collapse, so the words should survive verbatim — the detached audit (AC2) confirms each ethos bullet lands intact in SKILL.md, and the live scenarios (AC3) confirm the FO still behaves under the relocated framing. The one thing the relocation could silently break is the seam: orphaning the `## Working Principles` self-reference (line 211) or leaving a stale forward-pointer at the old ethos site. AC5's non-duplication check catches a half-done move; the seam-rewrite is called out in the before/after wording so the implementer doesn't leave a dangling reference.

## Acceptance criteria

- **AC1 (structural — the collapse landed).** `claude-fo-merge.md` carries exactly ONE mod-block section (`## Mod-Block Enforcement at Terminal Transitions` removed; its non-redundant content folded into `## Mod-Block Enforcement`), and `first-officer-shared-core.md:199`'s MODS-REPORT restatement is a pointer.
  *Verified by:* a header-count assertion in the validation diff (the collapsed file has one `## Mod-Block` header, down from two) — an on-disk fact outside the task body, checkable by `grep -c '^## Mod-Block' skills/first-officer/references/claude-fo-merge.md` returning `1`.

- **AC2 (high-stakes detached audit — no obligation lost).** A word-level diff of every collapsed/consolidated/relocated/compressed obligation against the pre-restructure baseline (the files at #367's merged tree) confirms no MUST / MUST-NOT / qualifier (NEVER, only, unless, except) dropped or inverted across all change classes (mod-block collapse, MODS-REPORT pointer, ethos hoist, the C-a `## Clarification and Communication` compression, and the C-b `## Working Principles` bullet compressions), the `TERMINAL_TEARDOWN_BOUNDED: best-effort teardown exhausted; member(s) stuck in registry; holding for launcher.` marker is byte-intact, AND every operating-principle bullet (the three ethos lines) survives verbatim in its new SKILL.md home with no principle dropped or weakened in the move. The audit specifically confirms: C-a dropped only restatement (not the keep-dispatching-other-entities rule or the report-once rule); the C-b bullets name-end-value / lead-with-recommendation / reversible-work lost only restatement/padding; lead-with-recommendation's "single-yes" teeth remain enforced by `present-gate` rule 1 (not silently dropped); and the speak-label and prefer-code-gate bullets are byte-untouched.
  *Verified by:* the detached audit's diff output (an independent reviewer / `git diff {#367-tree} -- claude-fo-merge.md SKILL.md first-officer-shared-core.md` whose every removed MUST-bearing clause and every removed ethos bullet is shown to survive verbatim in its destination or to be a pure restatement of a hoisted principle). The baseline is an independent source (the prior committed tree), not the task's own prose.

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
| AC2 | Detached word-level diff vs #367 merged tree (`git show f87107b1:…` baseline) across all changed files; MUST/MUST-NOT/qualifier survival across all four change classes + marker byte-check + every ethos bullet survives in SKILL.md + C-a dropped only restatement. Run FIRST (riskiest path). | detached adversarial review | minutes |
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

## Stage Report: ideation (cycle 3)

- DONE: Produce a concrete principle-cascade map (not an open-ended "simplify everything") classifying deep sections COMPRESS vs KEEP VERBATIM.
  Added `## Principle-cascade map (item 4)` surveying all FO refs (shared-core + dispatch/merge/runtime); each candidate classified with the principle it derives from or the scenario/marker that makes it load-bearing.
- DONE: For COMPRESS items, concrete before/after + which principle it derives from.
  C-a (the `## Clarification and Communication` "don't ask to take an allowed step" prose) given before/after, deriving from the reversible-work principle + the concise/JFDI ethos bullet. C-b (the FO-posture bullets) marked "kept — ambiguous derivability" per the decision rule.
- DONE: For KEEP items, name the guard/marker/test that grades them.
  Enumerated never-self-approve (gate-guardrail), TERMINAL_TEARDOWN_BOUNDED byte-marker, mod-block set→invoke→clear (merge-hook-guardrail), rebase-conflict HALT, ID-style/next-id rules, dispatch numbered procedure, AC cross-check, not-a-stopping-point, context-budget, contract-version gate.
- DONE: Apply the default-to-keep decision rule; mark ambiguous derivability as KEEP.
  C-b explicitly kept-ambiguous (the entry-point ethos lacks the operational teeth — "single yes", "not a menu" — that present-gate leans on).
- DONE: Recommend one pass vs staged audited passes.
  ONE pass — the genuine COMPRESS set is a single surgical compression (C-a); not large enough to stage. The pass still pays full audit price (AC2 detached diff + AC3 live-green), now covering C-a.
- DONE: Extend ACs to the compression; keep the detached-audit AC covering no MUST/MUST-NOT/qualifier dropped/inverted.
  AC2 widened to four change classes incl. C-a, asserting C-a dropped only restatement (not the keep-dispatching / report-once rules).

### Summary

Built the principle-cascade map the captain asked for and reported the honest headline: the compressible surface is SMALL because the FO contract is already overwhelmingly load-bearing mechanism, not posture restatement. Exactly one surgical COMPRESS lands (C-a, the "don't ask to take an allowed step" clause → a one-clause pointer to the reversible-work principle); the FO-posture bullets are kept-ambiguous under the default-to-keep rule because the terse entry-point ethos lacks their operational teeth. Ten-plus load-bearing mechanisms are enumerated KEEP-VERBATIM, each named with the scenario/marker/guard that grades it. Recommended ONE implementation pass (the COMPRESS set is too small to stage), still paying the full detached-audit + live-green price. The deliverable is the PLAN at the gate, not the cuts executed blind.

## Stage Report: ideation (cycle 4)

(Continuation of the cycle-3 cascade-map work, deepened per the lead's follow-up steer: examine `## Working Principles` harder, add a line-count delta, resolve the C-b ambiguity.)

- DONE: Examine the `## Working Principles` block as a prime compression candidate (lead's explicit ask).
  C-b split bullet-by-bullet: name-end-value / lead-with-recommendation / reversible-work → COMPRESS (restatement only); speak-label / prefer-code-gate → KEEP VERBATIM (mechanism / anti-tautology spine). Each with concrete before/after.
- DONE: Resolve the cycle-3 "kept — ambiguous" verdict on the posture bullets.
  Resolved via a `present-gate` survey: lead-with-recommendation's operational teeth ("single yes", lede-first spine, "single recommendation") ALREADY live in present-gate rule 1 + the nine-rule set, so the shared-core bullet is a middle-layer restatement that safely compresses. The teeth aren't lost — they have a real home this task doesn't touch.
- DONE: Rough line-count delta estimate for the compressible set.
  ~40–45 words total (C-a ~12, C-b ~30), ~8–12 lines of ~230-line shared-core — low-single-digit-percent. Stated plainly that the payoff is the cited-source framing, not the line count.
- DONE: One pass vs staged passes recommendation, with the honest small-set finding.
  ONE pass — small and all in shared-core. Full audit price retained (AC2 detached diff now covers C-a + all three C-b rewrites, incl. verifying lead-with-recommendation's teeth still enforced by present-gate). The honest "cascade yields little beyond the ethos hoist because the deep prose is mostly load-bearing mechanism" finding is stated in the headline.

### Summary

Deepened the cascade map per the lead's steer. The headline honest finding stands and is now quantified: the compressible surface is ~40–45 words / ~8–12 lines, low-single-digit-percent of shared-core, because the deep prose is overwhelmingly load-bearing mechanism the principles cannot replace. The cycle-3 "kept — ambiguous" verdict on the FO-posture bullets is RESOLVED by surveying `present-gate`: lead-with-recommendation compresses because its gate-time teeth already live there; name-end-value and reversible-work compress only restatement/padding while keeping their distinct teeth; speak-label and prefer-code-gate stay verbatim as mechanism. Recommended ONE audited implementation pass. The real cascade payoff is making the entry-point ethos the cited source future prose leans on, not the modest line saving.

## Stage Report: implementation

- DONE: (a) ethos HOIST to SKILL.md — self-explanatory, label + meta-framing dropped; deleted from shared-core; :211 lede rewritten to point UP.
  SKILL.md `## How the first officer operates` inserted after single-entity block; `grep -l 'Begin with the end'` → SKILL.md only. Commit cf562ad5.
- DONE: (b) the two collapses — merge-ref two mod-block sections → ONE; C1 MODS-REPORT restatement → pointer.
  `grep -c '^## Mod-Block' claude-fo-merge.md` == 1; shared-core:199 now "(Startup step 5)" pointer. Recovery/resume/missing-mod prose preserved.
- DONE: (c) cascade compressions — C-a (Clarification → reversible-work pointer) + three Working-Principles posture bullets; speak-label + prefer-code-gate kept verbatim.
  prefer-code-gate and speak-label byte-identical to f87107b1 baseline (diff empty).
- DONE: Re-confirmed all line anchors against origin/main first.
  Worktree HEAD == origin/main (3779370f); ethos 5–12, mod-block 54/64, MODS 199, lede 211, C-a 207, C-b 217–220 all matched.
- DONE: MANDATORY high-stakes detached audit (self-run before signaling).
  Word-level diff vs #367 baseline (`git show f87107b1:…`) on a detached /tmp checkout, never mutating the worktree. No MUST/MUST-NOT/qualifier dropped or inverted across all five change classes; every removed mod-block obligation (4 mechanism conditions, recovery options, session-resume, missing-mod recovery) survives in the canonical section; `TERMINAL_TEARDOWN_BOUNDED: …` marker byte-intact (grep==1); three ethos bullets verbatim in SKILL.md; lead-with-recommendation's lede-first teeth confirmed live in present-gate SKILL.md rule 1.
- DONE: Gates green.
  `grep -c '^## Mod-Block'` == 1; ethos single-home grep → SKILL.md only; `go test ./...` exit 0 incl. internal/contractlint reference-closure + internal/ensigncycle offline. (Live Claude+Codex shared-scenarios are the validation-stage AC — noted, not run here.)

### Summary

Executed the cycle-4 plan exactly: hoisted the FO operating ethos into the always-loaded SKILL.md entry point (label + meta-framing dropped, plain posture), deleted it from shared-core, and rewrote the Working-Principles lede to point up. Collapsed claude-fo-merge.md's two mod-block sections into one canonical `## Mod-Block Enforcement` and the shared-core MODS-REPORT restatement into a pointer. Compressed C-a and three posture bullets while keeping speak-label and prefer-code-gate byte-verbatim. The detached word-level audit against the #367 baseline confirms no obligation lost or inverted and the teardown marker byte-intact; the offline gate (contractlint reference-closure + ensigncycle) is green. Net: 14 insertions / 33 deletions across three files. Live shared-scenarios remain the validation-stage bar.

## Stage Report: validation

- DONE: Independently reproduce the MANDATORY high-stakes detached audit on a SEPARATE detached checkout (never the impl worktree).
  Created `/tmp/fo-refs-audit-$$` via `git worktree add --detach cf562ad5`; ran all word-level diffs there; removed it after. Worktree never mutated.
- DONE: Word-level diff of every collapsed/compressed/relocated obligation vs the #367 baseline (`git show f87107b1:…`); confirm no MUST/MUST-NOT/qualifier dropped or inverted across all five change classes.
  Isolated THIS commit's changes (cf562ad5 vs parent 3779370f) — the `BASE=$(spacedock dispatch trunk)` merge-ref change in the f87107b1 diff is PRE-EXISTING on the parent (baseline drift between #367 and trunk), NOT this task's. Mod-block collapse: all 4 enumerated conditions, set-before-invoke, clear-after-resolve, standalone-clear refusal, recovery options, session-resume scan, and missing-mod recovery survive in the canonical `## Mod-Block Enforcement`. The deleted "the FO MUST" numbered procedure was a prose restatement of the mechanism invariant whose teeth (the guard) are intact. Ethos hoist, MODS pointer, C-a, C-b all dropped only restatement. C-b lead-with-recommendation's lede-first/single-recommendation teeth confirmed live in `present-gate` SKILL.md rule 1.
- DONE: `TERMINAL_TEARDOWN_BOUNDED: …` marker byte-intact (grep==1); all three ethos bullets verbatim in SKILL.md; speak-label + prefer-code-gate byte-identical to baseline.
  Marker grep==1 (lives under `## Merge and Cleanup`). Three ethos bullets in SKILL.md, absent from shared-core; `(ethos)` label + meta-framing gone. prefer-code-gate + speak-label: empty diff vs f87107b1.
- DONE: Structural gates from the worktree — `grep -c '^## Mod-Block' claude-fo-merge.md` == 1; ethos single-home → SKILL.md only; `go test ./...` incl. internal/contractlint reference-closure green.
  mod-block header count == 1; `grep -l 'Begin with the end'` → SKILL.md only; `go test ./...` exit 0, contractlint forced `-count=1` ok 0.221s, ensigncycle offline green.
- DONE: Confirm structural consistency with the live shared scenarios; NOTE live green gates at PR CI, do not fabricate.
  `merge-hook-guardrail` scenario present (shared_scenarios_test.go:39, intent matches the surviving invariant). The collapse keeps the graded invariant byte-equivalent. Live Claude+Codex shared scenarios (AC3) run in the PR live-e2e gate — NOT run here, NOT fabricated.
- FAILED: detached audit refutation — one Material finding.
  `skills/first-officer/references/claude-first-officer-runtime.md:13` (merge-ref load-point inventory) still names "**Mod-Block Enforcement at Terminal Transitions** (the `TERMINAL_TEARDOWN_BOUNDED` bounded-teardown marker)" as a section in claude-fo-merge.md. The collapse REMOVED that section (`grep -c '^## Mod-Block Enforcement at Terminal Transitions'` == 0). On the baseline both sections existed so the inventory was correct; post-collapse it is a stale cross-file section reference. contractlint stays GREEN (it does file-path reference-closure, not prose section-name enumeration) — exactly the test-strength hole the detached audit exists to catch. AC3 also stays green (marker + guard behavior intact), so no behavior test sees it. The AC2/AC5 seam checks were scoped to shared-core's lede + old ethos site; they missed this runtime-adapter inventory.

### Summary

REJECTED. The five in-scope change classes are clean — every collapsed/compressed/relocated obligation survives with no MUST/MUST-NOT/qualifier dropped or inverted, the `TERMINAL_TEARDOWN_BOUNDED` marker is byte-intact, the three ethos bullets land verbatim in SKILL.md with a single canonical home, prefer-code-gate and speak-label are byte-identical, and the offline gate (incl. contractlint reference-closure) is green. But the detached adversarial audit refuted clean-pass: the mod-block section collapse orphaned a cross-file inventory reference at `claude-first-officer-runtime.md:13`, which still names the removed `Mod-Block Enforcement at Terminal Transitions` section as living in the merge ref. This is a Material seam orphan — contractlint and the live scenarios both stay green over it, so only the audit catches it. Fix: update the runtime adapter's merge-ref load-point inventory to name only the surviving `## Mod-Block Enforcement` section while keeping the `TERMINAL_TEARDOWN_BOUNDED` marker attribution accurate to its real `## Merge and Cleanup` home, then re-validate. AC3 (live Claude+Codex shared scenarios) gates at PR CI and is NOT fabricated here.

### Feedback Cycles

- **Cycle 1 (validation REJECTED, 2026-06-14):** the detached adversarial audit caught a Material dangling reference the impl self-audit missed — the mod-block section collapse removed `## Mod-Block Enforcement at Terminal Transitions` from claude-fo-merge.md, but `claude-first-officer-runtime.md:13`'s merge-ref load-point inventory still names that now-removed section. contractlint + live scenarios stay green over it (only the word-level audit catches it). Fix: update the runtime adapter's merge-ref inventory to name only the surviving `## Mod-Block Enforcement`, keeping TERMINAL_TEARDOWN_BOUNDED attributed to its real `## Merge and Cleanup` home. Routed back to implementation.

## Stage Report: implementation (cycle 1)

- DONE: Fix the dangling runtime-adapter merge-ref inventory (validation REJECTED cycle 1).
  `claude-first-officer-runtime.md:13` named the removed `## Mod-Block Enforcement at Terminal Transitions` section; rewrote the inventory to name only the surviving `## Mod-Block Enforcement` and re-homed the marker to the bounded terminal teardown, matching `first-officer-shared-core.md:130`'s phrasing.
- DONE: Re-run detached audit — cross-file grep for the removed section name, clean.
  `grep -rn 'Mod-Block Enforcement at Terminal Transitions' skills/first-officer/` → zero hits tree-wide. `TERMINAL_TEARDOWN_BOUNDED` marker byte-intact (grep==1) at its real home `claude-fo-merge.md:24` (Merge-and-Cleanup step 10).
- DONE: `go test ./...` incl. contractlint green.
  Re-run after the fix; exit 0.

### Summary

Cycle-1 fix: the detached validation audit caught a Material dangling reference the implementation self-audit missed — the `claude-first-officer-runtime.md` merge-ref inventory still named the `## Mod-Block Enforcement at Terminal Transitions` section my collapse removed. De-dangled the inventory to name only the surviving `## Mod-Block Enforcement`, re-homed the `TERMINAL_TEARDOWN_BOUNDED` marker to the bounded-teardown step (its real home), and confirmed by tree-wide grep that no other cross-file reference names the removed section. Offline gate green.

## Stage Report: validation (cycle 2)

- DONE: Re-verify the cycle-1 fix from the worktree (not the report) — dangling reference gone tree-wide.
  `grep -rn 'Mod-Block Enforcement at Terminal Transitions' skills/first-officer/` → ZERO hits. `claude-first-officer-runtime.md:13` now names "Mod-Block Enforcement, and the bounded terminal teardown (the `TERMINAL_TEARDOWN_BOUNDED` marker)" — removed section name gone, marker re-homed. Fix diff (449b8502 vs cf562ad5) touches ONLY the runtime adapter, 1 line, no scope creep.
- DONE: Re-run the full detached adversarial audit on a SEPARATE checkout — five original change classes still clean vs #367 baseline.
  Fresh `git worktree add --detach 449b8502` in /tmp (removed after). This task's complete diff vs the #367 tree parent (3779370f) = 4 files, 15 ins / 34 del. No MUST/MUST-NOT/qualifier dropped or inverted across the five classes (re-confirmed: mod-block collapse, MODS pointer, ethos hoist, C-a, C-b all unchanged from cycle-1). prefer-code-gate + speak-label STILL byte-identical to f87107b1 (empty diff).
- DONE: `TERMINAL_TEARDOWN_BOUNDED: …` marker byte-intact at its real home (grep==1); ethos single-home (SKILL.md only); NO new dangling reference introduced.
  Marker grep==1 under `## Merge and Cleanup` (its real home, not the removed section). `grep -l 'Begin with the end'` → SKILL.md only. `grep -rn 'at Terminal Transitions' skills/` → ZERO hits; fix introduced no new dangling ref.
- DONE: `grep -c '^## Mod-Block' claude-fo-merge.md` == 1; `go test ./...` incl. contractlint reference-closure green.
  Header count == 1. `go test ./...` exit 0; contractlint forced `-count=1` ok 0.221s on the actual fixed tree (the cached pass was the prose edit not touching a Go input — forced run confirms reference-closure green). merge-hook-guardrail scenario present (structural consistency intact).
- DONE: AC3 live Claude+Codex shared scenarios — gate at PR CI, not fabricated.
  Noted: AC3 runs in the PR live-e2e gate. NOT run here, NOT fabricated. The collapse keeps the merge-hook-guardrail-graded invariant byte-equivalent.

### Summary

PASSED. The cycle-1 fix (449b8502) closed exactly the Material finding from cycle 1: the stale `claude-first-officer-runtime.md:13` merge-ref inventory now names only the surviving `## Mod-Block Enforcement` section and re-homes `TERMINAL_TEARDOWN_BOUNDED` to "the bounded terminal teardown". The fix is a surgical one-line prose edit touching only the runtime adapter — no scope creep, no regression in the five original change classes. Re-ran the full detached audit on a fresh checkout: no MUST/MUST-NOT/qualifier dropped or inverted vs the #367 baseline, marker byte-intact at its real `## Merge and Cleanup` home, ethos single-home in SKILL.md, prefer-code-gate + speak-label byte-identical, no new dangling reference. Structural gates green (mod-block header count == 1, `go test ./...` exit 0 incl. contractlint reference-closure forced). AC3 (live Claude+Codex shared scenarios) gates at PR CI and was deliberately not fabricated.
