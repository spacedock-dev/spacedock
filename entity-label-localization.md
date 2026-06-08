---
id: jm0vqtx3j5vfw106kg0nz87b
title: Localize the operating voice to the workflow's entity-label (keep the shared contract generic)
status: implementation
source: "captain (2026-06-08, this session). Captain-facing + workflow-specific prose says generic \"entity\"; it should read as the workflow's declared label (this workflow: \"task\"). Principle: localize the OPERATING VOICE, keep the SHARED ABSTRACTION generic — do NOT rename entity→task in the shared contract (it serves any workflow: ticket/story/experiment)."
started: 2026-06-08T18:17:44Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-entity-label-localization
issue:
sprint: 0199-pre-flip-mechanics
group: dev-quality
sprint-readiness: ready
---

Add a generic **operating-voice convention** to the shared agent contract: when an agent produces captain-facing or operating prose (gate presentations, status reports, the dispatch package, conversation), it refers to entities by the workflow's declared `entity-label` / `entity-label-plural` — resolved from the README it already read at boot — instead of the generic "entity". The convention is generic ("use the workflow's declared label"), never a hardcoded "task". The contract *mechanics* stay generic "entity"; only the human-facing *output* localizes.

## Problem

The fix is **how agents refer to entities in their own generated prose**, not a set of template files to edit. The captain's words: "it is how we refer to the entities."

The README already declares `entity-label` / `entity-label-plural` (here: `task` / `tasks`), and the FO reads them at boot (shared-core Startup step 4). But there is no convention telling the agent to *use* that label when it speaks to the captain, so its generated prose defaults to the generic "entity" everywhere the human reads it.

The original framing of this task — "add `{entity-label}` placeholders to the present-gate and commander-dispatch templates" — targeted surfaces that do not exist as editable templates. `skills/present-gate/SKILL.md`'s only "entity" strings are the placeholder field name `{entity title}` (abstraction-level structure, not voice) and an assembly-rule instruction to the FO ("open the entity file"); its captain-facing `Decision:` line is noun-free. And there is no commander-dispatch template at all — each `docs/roadmap/{sprint}/dispatch-sprint-execution.md` is hand-authored from scratch. "Adding a placeholder to a template" therefore had no real target, and a live drive over it reduced to a tautology (the FO writes a handed label, then observes it). The corrected deliverable is a **generic convention in the shared contract**, proven by observing that the agent's *own generated* prose resolves the README label — an independent source the test never hands it.

## Proposed approach

Add **one generic convention** to the shared agent contract. No template files; no binary changes; no hardcoded label.

### Where the convention lives

- **`skills/first-officer/references/first-officer-shared-core.md`** — the natural home is the `## Working Principles` → **FO posture** list, alongside "Name the end value" and "Lead with a recommendation": these are exactly the rules governing how the FO talks to the captain. Add a posture bullet: when the FO produces captain-facing or operating prose, it refers to the workflow's entities by the declared `entity-label` / `entity-label-plural` (read at Startup step 4), not the generic word "entity". Mechanics stay generic — `entity_path`, the entity-read line, the abstraction — only the human-facing noun localizes.
- **`skills/present-gate/SKILL.md`** (the captain-facing assembly rules) — add a one-line rule that the gate-summary prose the FO writes (the `Decision:` line, the gist roll-up, the chosen-direction sentence) speaks the declared label, so the convention is reinforced at the single highest-traffic captain-facing surface. The template's `{entity title}` placeholder and structural headings stay generic; only the FO-authored prose localizes.

The instruction is generic in both places ("resolve and use the workflow's declared `entity-label`"), never the literal "task". A workflow declaring `entity-label: experiment` gets "experiment"; one declaring `ticket` gets "ticket"; the default `entity` workflow is unchanged.

### Guardrail (load-bearing, unchanged from cycle 1)

The contract *mechanics* stay generic "entity": `first-officer-shared-core.md`'s and `ensign-shared-core.md`'s abstraction prose, the `dispatch build` ensign-package output (`entity_path`, the entity-read line), `spacedock status` field names and machine output. Those serve any workflow and are never localized — they are not operating voice. Abstraction-generic and voice-localized do not conflict: they govern disjoint text (machine/contract mechanics vs. captain-facing generated prose). The no-leak guardrail is the same one cycle 1 carried as AC-3.

### No spike needed

No spike needed — the design composes only already-proven behavior. (1) The FO already reads the README `entity-label` at Startup step 4, so the value is in hand; nothing new is parsed. (2) The FO already presents gates and status live; this only changes the noun in prose it already generates. The one genuine unknown — "does a contract-level convention actually make the agent's *generated* prose track the README label, rather than default to 'entity'?" — is a behavior question that cannot be de-risked by a throwaway spike; it is exactly what the live drive in the test plan proves (per the README `## done` rule: a contract/skill change is PASSED only by a live drive observing the behavior). So the live drive is the proof, not a precondition spike.

## Out of scope

- **The cross-workflow `{wf}#{ref}` qualifier.** Deferred — the captain's cycle-1 downscope stands. This re-ideation is the label convention only; do NOT re-add the qualifier here. Its banked design (separator `#`, source = workflow dir basename / the `StateBranch` suffix, qualify-only-when-≥2-workflows) is recorded in the cycle-1 Feedback entry and the cycle-1 stage report below — a later task recovers it from there. The cycle-2 reframe removed the detailed "Cross-workflow qualification" approach subsection and AC-4 from the live body, leaving this one-line deferral as the pointer.
- **Binary substitution** of the label in `dispatch build` / `status` output. Those stay generic "entity" — they are machine mechanics, not operating voice, and are explicitly protected by the guardrail.
- **A `spacedock:roadmap` graduation bake.** There is no `roadmap` skill in `skills/` yet; nothing to bake.
- Renaming `entity` → `task` (or any label) anywhere in the shared first-officer / ensign contract mechanics. The abstraction stays generic.
- Per-workflow forks of the shared skills.

## Acceptance criteria

- **AC-1 — the FO's generated captain-facing prose refers to entities by the README's declared label.** Driving the FO to present a gate (or a status report) for a fixture workflow whose README declares `entity-label: ticket` / `entity-label-plural: tickets`, the FO's own generated prose calls the entities "ticket" / "tickets" where it would otherwise say "entity" / "entities". **Verified by:** a live FO drive over the `entity-label: ticket` fixture, observing "ticket" in the FO's generated gate/status prose. The expected word "ticket" exists ONLY in the fixture's README `entity-label` field — an independent source the FO must READ and RESOLVE at boot; the test never hands the FO the word. A FO that ignores the convention emits "entity" and the drive FAILS. A grep over the contract/skill files does NOT satisfy this (proof-policy; README `## done`).
- **AC-2 — the resolved label tracks the README, not a hardcoded value (anti-tautology differential).** The same FO drive run over a SECOND fixture whose README declares a DIFFERENT label (`entity-label: experiment` / `experiments`) produces prose saying "experiment", not "ticket" and not "entity". Two different README declarations yielding two different generated nouns proves the FO is resolving the field, not parroting a value the test supplied. **Verified by:** the `ticket`-fixture drive (AC-1) and the `experiment`-fixture drive producing label-matched prose respectively; were the convention hardcoded or the proof tautological, one of the two drives would mismatch its README.
- **AC-3 — the contract mechanics still read generic "entity" (no-leak guardrail).** `first-officer-shared-core.md` / `ensign-shared-core.md` abstraction prose and the `spacedock dispatch build` ensign-package output keep generic "entity" (`entity_path`, the entity-read line). **Verified by:** the existing `dispatch build` golden over the dev fixture (`go test ./internal/dispatch -run TestBuild`) still passing unchanged — its frozen expected text says "entity file"; a leak of localization into the machine package would diverge the golden and fail. This AC's expected value is the frozen golden, independent of the skill edits.

## Test plan

The behavioral half — that the convention actually makes the agent's *generated* prose track the declared label — is proven by a **live FO drive**, never by a substring match over the contract/skill text (proof-policy; this workflow's established discipline for skill/contract-prose changes, cf. `survey-skill-correctness-pass` AC-4: "a grep over SKILL.md does NOT satisfy this — only the rendered output does"). The convention edits themselves (the posture bullet, the assembly-rule line) are authoring work, NOT standalone behavioral ACs.

**The anti-tautology design (the crux the reviewer flagged):** the only place the expected noun lives is the fixture README's `entity-label` field. The FO discovers it by executing Startup step 4 (read the README), then chooses whether to honor the convention when it writes captain-facing prose. The test asserts the *generated* prose against the README field — two independent artifacts that can diverge. The failure mode is real and observable: a FO that ignores the convention writes "entity" and the assertion fails. The two-fixture differential (AC-2: `ticket` vs `experiment`) closes the last loophole — a hardcoded or test-handed label could pass one fixture but cannot match both READMEs.

- **Fixtures (throwaway):**
  - Fixture A: a minimal workflow whose README declares `entity-label: ticket` / `entity-label-plural: tickets`, with one entity carrying a `## Stage Report` and an `## Acceptance criteria` section so the FO can present a gate.
  - Fixture B: the same shape but `entity-label: experiment` / `experiments`, for the differential.
- **AC-1 / AC-2 — live FO drive (the only proof for the generated label).** Build `spacedock`, point the FO at Fixture A, drive it to the gate-present (or a status-report turn), and observe "ticket(s)" in its generated prose; repeat over Fixture B and observe "experiment(s)". Confirm neither drive was handed the noun by the harness (the prompt names the fixture path, not the label). Cost: minutes per drive; needs the FO + present-gate skills and a built binary. Negative control: the default-`entity`/dev workflow drive still says "entity" (the convention is a no-op when label == "entity").
- **AC-3 — guardrail code gate.** `go test ./internal/dispatch -run TestBuild` (and the full `internal/dispatch` golden suite) stays green: the ensign dispatch package's frozen "entity"-bearing text is unchanged, proving localization did not leak into the machine mechanics. Independent expected value = the checked-in golden.
- **Estimated total cost/complexity: LOW.** Two small prose edits (a posture bullet in `first-officer-shared-core.md`; a one-line assembly rule in `present-gate/SKILL.md`); two throwaway fixtures; one reused `dispatch build` golden assertion; two live FO drives (`ticket` + `experiment`) plus the default-`entity` negative control. No binary/Go production changes; no template files; no qualifier.

## Stage Report: ideation

- DONE: Lock Layer-1 scope: which captain-facing surfaces get {entity-label} placeholders (present-gate + commander-dispatch templates) and the cross-workflow {wf}#{ref} qualifier source (dir-basename vs declared short-name) + separator.
  Locked to the two model-filled surfaces — `skills/present-gate/SKILL.md` and `docs/roadmap/{sprint}/dispatch-sprint-execution.md`; qualifier source = workflow dir basename (the `StateBranch` suffix, `internal/status/state.go`), separator `#`, no new README field. See "Layer-1 scope" + "Cross-workflow qualifier source".
- DONE: Produce AC + test plan proven by RENDERING over a fixture (a workflow whose entity-label != "entity" shows the label in a rendered gate; a two-workflow fixture shows qualified refs) — never a substring grep over the template (proof-policy).
  AC-1/AC-2 = live present-gate + commander-dispatch render over an `entity-label: ticket` fixture (expected noun from fixture README); AC-4 = live two-workflow `dev`+`user-testing` render (qualifier from dir basename); AC-3 guardrail = existing `dispatch build` golden + region-disjointness. No grep-as-proof anywhere.
- DONE: Record "no spike needed: {the proven rendering mechanism it composes}" if the design only reuses already-proven rendering.
  Recorded in "No spike needed": composes (1) FO already reads README entity-label at boot, (2) present-gate already renders gates live, (3) `StateBranch` basename resolution — riskiest unknown settled by the live drive, which is the proof not a de-risking spike.

### Summary

Locked Layer-1 to the two captain-facing, FO-model-filled surfaces (present-gate SKILL.md + the commander-dispatch `dispatch-sprint-execution.md` template), explicitly deferring Layer-2 binary substitution and the not-yet-existent `spacedock:roadmap` graduation bake. Pinned the cross-workflow qualifier to the workflow dir basename — reusing the proven `StateBranch` identity rather than adding a new README short-name knob (YAGNI) — with separator `#`. The four ACs prove the behavioral half by live rendering over fixtures whose declared label/basename is the independent expected value, never by a grep over the template; the shared-contract guardrail rides the existing `dispatch build` golden plus region-disjointness. Determination: no spike needed — the design only substitutes into already-proven render paths.

## Feedback Cycles

- **Cycle 1 — captain downscope at the ideation gate (2026-06-08).** Ship **Layer-1 label localization only** in 0.19.9: the present-gate + commander-dispatch templates speak the README `entity-label` (AC-1/AC-2), with the shared contract staying generic (AC-3). The **cross-workflow `{wf}#{ref}` qualifier is DEFERRED to a follow-up** — do NOT build it in 0.19.9. Its design (the "Cross-workflow qualification" approach subsection + **AC-4**) stays in this body as the banked follow-up design, but **AC-4 is out of 0.19.9 scope**: validation requires evidence for **AC-1/AC-2/AC-3 only**, and the two-workflow render + its fixture machinery are not built. Rationale: the qualifier adds two-workflow-fixture machinery for a presently-rare scenario; the label localization (the core ask) ships now, the qualifier when a real multi-workflow-in-one-context need appears.
- **Cycle 2 — captain reframe at the ideation gate (2026-06-08); independent staff review.** Cycle 1's "edit the present-gate / commander-dispatch templates" target was wrong: the review confirmed those surfaces do not exist as editable templates (present-gate's "entity" strings are placeholders + assembly-rule guidance; there is no commander-dispatch template — each `dispatch-sprint-execution.md` is hand-authored), so the placeholder-add had no target and its live drive was a tautology. **Corrected target: a generic operating-VOICE convention in the shared contract** — "when an agent generates captain-facing/operating prose, refer to entities by the workflow's declared `entity-label`, resolved from the README." Home: the FO posture list in `first-officer-shared-core.md` + the present-gate assembly rules. Proof is a live drive observing the FO's GENERATED prose track the README label (the word lives only in the fixture README, never handed to the FO), with a `ticket`-vs-`experiment` two-fixture differential to kill the tautology. The body's Problem/Approach/AC/Test-plan were rewritten to this target; AC-1/AC-2/AC-3 replace the old set; the qualifier stays deferred (cycle-1 downscope unchanged).

## Stage Report: ideation (cycle 2)

- DONE: Re-frame from template-editing to operating-voice convention.
  Confirmed the cycle-1 targets do not exist: `skills/present-gate/SKILL.md` "entity" occurrences are only `{entity title}` (placeholder) + "open the entity file" (FO instruction) — its `Decision:` line is noun-free; `dispatch-sprint-execution.md` has no generator (a proposal-doc artifact, hand-authored). Corrected target = a generic convention in `first-officer-shared-core.md` FO-posture list + present-gate assembly rules: agents refer to entities by the README's declared `entity-label` in captain-facing/operating prose. Guardrail (contract mechanics stay generic "entity") unchanged. Rewrote Problem / Proposed approach / Out of scope.
- DONE: Revise AC + test plan to a non-tautological live-drive proof.
  AC-1 = live FO drive over an `entity-label: ticket` fixture; observe "ticket" in the FO's GENERATED prose; the word lives only in the fixture README (independent source the test never hands the FO); a non-conforming FO says "entity" and fails. AC-2 = `ticket`-vs-`experiment` two-fixture differential — two READMEs → two different generated nouns proves resolution, not parroting. AC-3 = no-leak guardrail via the existing `dispatch build` golden (`go test ./internal/dispatch -run TestBuild`). No grep-as-proof.
- DONE: Spike-or-no-spike determination re-recorded for the corrected target.
  "No spike needed": composes (1) FO already reads `entity-label` at Startup step 4, (2) FO already presents gates/status live. The one genuine unknown (does a contract convention make GENERATED prose track the label?) is a behavior question only the live drive settles — it is the proof, not a precondition spike (README `## done`: contract/skill change PASSED only by a live drive).
- DONE: Append `## Stage Report: ideation (cycle 2)` noting the re-frame.
  This section; cycle-2 entry also added to `## Feedback Cycles`.

### Summary

Re-framed jm from "edit template files" to "add a generic operating-voice convention to the shared contract" after the staff review confirmed the cycle-1 surfaces (a noun-bearing present-gate `Decision:` line; a commander-dispatch template) do not exist. The deliverable is now one generic convention — agents refer to entities by the workflow's declared `entity-label` in captain-facing prose, resolved from the README — homed in the FO-posture list and the present-gate assembly rules, with the contract mechanics staying generic "entity". The proof is made airtight against the flagged tautology: the expected noun lives ONLY in the fixture README that the FO must read and resolve, asserted against the FO's *generated* prose, with a `ticket`-vs-`experiment` two-fixture differential so a hardcoded or test-handed label cannot pass. The cross-workflow qualifier stays deferred per the cycle-1 captain downscope.

## Stage Report: implementation

- DONE: AC-1/AC-2 substrate (the convention): add ONE generic operating-voice convention — a FO-posture bullet in `skills/first-officer/references/first-officer-shared-core.md` (`## Working Principles` → FO posture) ... PLUS a one-line assembly rule in `skills/present-gate/SKILL.md` that the FO-authored gate-summary prose speaks the declared label. Generic in BOTH places, NEVER hardcoded "task"; the `{entity title}` placeholder + structural headings stay generic.
  Worktree commit f3ff718a. FO posture bullet "Speak the workflow's declared label, not the generic 'entity'" added after the three existing posture bullets; present-gate assembly rule 10 added after rule 9 (Chosen direction / Checklist gist / Decision speak the label; `{entity title}` + headings stay generic). Both say "use the workflow's declared `entity-label`", never "task".
- DONE: Two live-drive fixtures: a minimal throwaway workflow whose README declares `entity-label: ticket` / `entity-label-plural: tickets`, and a second declaring `entity-label: experiment` / `experiments` — each with one entity carrying a `## Stage Report` + `## Acceptance criteria` so a FO can present a gate.
  `fixtures/entity-label-drive/{ticket,experiment}-workflow/` (commit f3ff718a). Each entity sits at the gated `review` stage with both sections. Both validate VALID, ORPHANS: none, entity present — exercised the built binary: `spacedock status --workflow-dir … --validate` and `--boot`. The expected noun lives ONLY in each README's `entity-label` field.
- DONE: AC-3 guardrail: `go test ./internal/dispatch -run TestBuild` stays green UNCHANGED — the dispatch-build ensign-package "entity"-bearing golden is untouched, proving localization did NOT leak into the machine mechanics.
  91 passed unchanged (full dispatch suite 195 passed). The golden's `Read the entity file at …` line is untouched; my edits live only in `skills/` and `fixtures/`, neither in `internal/dispatch`.

### Summary

Landed the deliverable as two small prose edits plus two throwaway drive fixtures (worktree branch `spacedock-ensign/entity-label-localization`, commit f3ff718a). The convention is generic in both homes — the FO-posture bullet and present-gate assembly rule 10 both say "use the workflow's declared `entity-label`", never a hardcoded "task"; only the human-facing noun localizes while contract mechanics (`entity_path`, the entity-read line, machine output) stay generic "entity". The two fixtures (`ticket` + `experiment`) validate VALID with no orphans and sit at a gated stage, ready for validation to drive a live FO over both and observe its GENERATED prose track each README label (AC-1/AC-2, anti-tautology differential). AC-3 guardrail confirmed green and unchanged (`TestBuild` 91 passed). PRE-EXISTING UNRELATED FAILURE surfaced separately to team-lead: `go test ./internal/status -run TestMigrationCheckFixturesParseConsistently` fails on `docs/roadmap/0198-pre-flip-hardening/debrief.md` (`session-date` bare-YAML scalar) — fails identically with my changes stashed; the file was committed in the branch base `83a95b86` (= merge-base with main), so it is on main, outside this task's scope, and untouched by my edits.
