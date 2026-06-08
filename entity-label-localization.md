---
id: jm0vqtx3j5vfw106kg0nz87b
title: Localize the operating voice to the workflow's entity-label (keep the shared contract generic)
status: ideation
source: "captain (2026-06-08, this session). Captain-facing + workflow-specific prose says generic \"entity\"; it should read as the workflow's declared label (this workflow: \"task\"). Principle: localize the OPERATING VOICE, keep the SHARED ABSTRACTION generic — do NOT rename entity→task in the shared contract (it serves any workflow: ticket/story/experiment)."
started: 2026-06-08T18:17:44Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0199-pre-flip-mechanics
group: dev-quality
sprint-readiness: ready
---

Make captain-facing and workflow-specific prose speak the workflow's own noun. A workflow that declares `entity-label: task` should present gates, dispatch packages, sprint docs, and status reports in terms of "task", while the shared contract keeps the generic "entity" that serves any workflow (ticket / story / experiment).

## Problem

The README already declares `entity-label` / `entity-label-plural` (here: `task` / `tasks`), and the FO already reads them at boot — but the operating voice still says generic "entity" everywhere the captain reads it. The fix is not to rename the abstraction; it is to *use* the label that already exists, only at the layer the human sees.

There is a second, related ambiguity: when a context spans **more than one workflow** (e.g. a session touching both `dev` and a `user-testing` workflow), a bare reference like "task foo" doesn't say *which* workflow's task. Localization must therefore also carry a workflow qualifier when workflows are mixed, and drop it when only one is in play.

## Proposed approach

This task ships **Layer 1 only**. Layer 2 (binary substitution in `dispatch build` / `status`) and the graduation bake are explicitly carried to follow-ups (see Out of scope) — Layer 1 is the smallest self-contained increment and unblocks the rest.

### Layer-1 scope — the two captain-facing surfaces (locked)

The localization lands in exactly the two captain-facing templates the FO/Commander fills, where prose names the workflow noun:

1. **present-gate** — `skills/present-gate/SKILL.md`. The gate template and assembly rules name the noun in the `Decision:` line prose (e.g. "approve to enter implementation" is noun-free, but "reject to bounce this **task** back to …" names it) and in the assembly-rule guidance that tells the FO how to talk about the entity to the captain. Add an `{entity-label}` / `{entity-label-plural}` placeholder where the prose currently says generic "entity", with a one-line instruction that the FO resolves it from the README `entity-label` it already read at boot. The abstraction-level scaffolding lines (the field name `{entity title}`, the `Checklist:` heading) stay generic — they are structure, not operating voice.
2. **commander-dispatch** — the cold-boot Commander package, `docs/roadmap/{sprint}/dispatch-sprint-execution.md` (the `spacedock:roadmap` graduation's `## Commander package` artifact). Today it names the noun as "members"; localize the noun-bearing prose ("drive its members" → "drive its {entity-label-plural}") via the same `{entity-label}` placeholder + resolve-from-README instruction.

Both surfaces are filled by the FO **model** (no binary renders them), so the discipline is "use the label you already read at boot", carried as a placeholder + a one-line resolve instruction — never a second README read bolted on.

### Cross-workflow qualifier source + separator (locked)

- **Separator:** `#` (given by the seed).
- **Source:** the **workflow directory basename** — the same identity that already names the state branch. `StateBranch` (`internal/status/state.go`) resolves `spacedock-state/<basename>` by default with a README `state-branch:` override winning verbatim; the qualifier is exactly that branch suffix (`dev` for `docs/dev`). **No new README short-name field** is introduced — reusing the existing basename identity keeps one source of truth and avoids a config knob nobody has asked for (YAGNI). If a workflow ever declares `state-branch:`, the qualifier is its trailing path segment, matching what peers already see in `git push origin spacedock-state/<x>`.
- **When to qualify:** the FO qualifies (`{basename}#{ref}`, e.g. `dev#k6`) only when the operating voice spans ≥2 workflows in one context; with a single workflow in play the bare localized form is correct ("task k6, j7"). This is FO discipline carried in the same captain-facing surfaces, proven by the same live-drive rendering.

### Guardrail (load-bearing)

The shared *contract* stays generic "entity": `first-officer-shared-core.md`, `ensign-shared-core.md`, the `dispatch build` ensign-package prose (`entity_path`, the entity-read line), and the *abstraction-level* template field names. Localization touches only the human-facing noun in the operating voice. The guardrail AC below and the localization AC are non-colliding by construction — they name disjoint text regions.

### No spike needed

No spike needed: this composes only already-proven rendering. (1) The FO already reads the README `entity-label` at boot — `first-officer-shared-core.md` step 4 has it read "entity labels" from the README — so the value is in hand; nothing new is parsed. (2) The present-gate skill already renders gates live, proven by every gate this workflow has presented; adding a placeholder it resolves from a value it holds is a substitution into a proven render path, not a new mechanism. (3) The qualifier source reuses the proven `StateBranch` basename resolution verbatim. The riskiest unknown — "does a model-filled template actually surface the declared label to the captain?" — is settled by the live drive in the test plan, which is the proof, not a de-risking spike.

## Out of scope

- **Layer 2** — binary substitution of the label in `dispatch build` / `status` output. The ensign dispatch package and status views stay generic this task; a follow-up may golden-test a `dispatch build` over a non-default-`entity-label` fixture.
- **Graduation bake** — making the `spacedock:roadmap` skill *emit* the placeholder automatically. There is no `roadmap` skill in `skills/` yet; baking it in is a follow-up once that skill exists. This task localizes the commander-dispatch *template text*, not the generator that will one day stamp it.
- Renaming `entity` → `task` (or any label) in the shared first-officer / ensign contracts. The abstraction stays generic.
- A new README short-name field for the qualifier. The dir basename is the source.
- Per-workflow forks of the shared skills.

## Acceptance criteria

- **AC-1 — a gate rendered for a non-default-`entity-label` workflow shows the declared label.** Driving the present-gate render (FO invoking `Skill(skill="spacedock:present-gate")`) over a fixture workflow whose README declares `entity-label: ticket`, the captain-facing gate text speaks the declared noun where it previously said generic "entity" (e.g. the `Decision:` bounce-back prose reads "ticket", not "entity"). **Verified by:** a live present-gate drive over the `entity-label: ticket` fixture, observing "ticket" in the rendered gate output. The expected value ("ticket") comes from the fixture README, not from the template under test — a grep over `SKILL.md` for the placeholder does NOT satisfy this.
- **AC-2 — the commander-dispatch template renders the declared plural label.** A commander-dispatch package authored for the `entity-label: ticket` fixture (the FO filling `dispatch-sprint-execution.md` from the localized template) names the workflow noun as its declared plural ("drive its tickets"), not "members"/"entities". **Verified by:** the rendered package text over the fixture, expected value from the fixture README. Not a grep over the template.
- **AC-3 — the shared contract still reads generic "entity".** `first-officer-shared-core.md` and `ensign-shared-core.md` and the `dispatch build` ensign-package output keep generic "entity" in the abstraction (`entity_path`, the entity-read line, contract prose). **Verified by:** an existing `dispatch build` golden over the dev fixture still passing unchanged (the package says "entity file", not "task file") AND the guardrail being non-colliding with AC-1/AC-2 by region (gate/commander prose vs. contract files). This AC fails if localization leaked into the shared contract.
- **AC-4 — cross-workflow references render qualified; single-workflow render bare.** In a rendered operating-voice surface spanning two workflows (a `dev`+`user-testing` fixture pair), entity references carry the `{basename}#{ref}` qualifier (`dev#k6`, `user-testing#foo`); the same surface over a single workflow renders the bare localized form ("task k6"). The qualifier source is the dir basename (= the `StateBranch` suffix). **Verified by:** a live render over the two-workflow fixture observing the qualified form, vs. the single-workflow fixture observing the bare form — expected qualifiers (`dev`, `user-testing`) derived from the fixture dir names, not from the template. Not a grep.

## Test plan

The behavioral half — a model-filled captain-facing surface actually surfaces the declared label — is proven by **rendering via a live drive**, never by a substring match over the template (proof-policy; matches this workflow's established discipline for skill-rendered surfaces, e.g. `survey-skill-correctness-pass` AC-4). The template-prose edits (placeholder + resolve instruction) are authoring work, not standalone behavioral ACs.

- **Fixtures (throwaway, checked into the test/scratch as needed):**
  - A single-workflow fixture whose README declares `entity-label: ticket` / `entity-label-plural: tickets`, with one entity carrying a `## Stage Report` so a gate can be presented.
  - A two-workflow fixture pair: `dev/` and `user-testing/`, each a minimal workflow with one entity, for the qualifier render.
- **AC-1 / AC-2 — live present-gate + commander-dispatch render (the only proof for the rendered label).** Drive the FO present-gate render over the `entity-label: ticket` fixture and observe "ticket" in the gate's `Decision:` prose; author a commander-dispatch package over the same fixture and observe "tickets" in the drive-split prose. Cost: minutes per drive; needs the FO/present-gate skills and a built `spacedock` for `status`. Expected nouns come from the fixture README.
- **AC-3 — guardrail, code-gated where possible.** Run the existing `dispatch build` golden (`internal/dispatch` golden harness) over the dev fixture; assert it still emits generic "entity" prose unchanged (`go test ./internal/dispatch -run TestBuild`). This is a real code gate over the shared-contract half: the golden's expected text is the frozen package, independent of the skill edits, so a leak of localization into the ensign package fails it. The contract-file half (`*-shared-core.md` unchanged abstraction) is asserted by region-disjointness from AC-1/AC-2, not by a self-referential grep.
- **AC-4 — live two-workflow qualifier render.** Drive the operating-voice surface over the `dev`+`user-testing` fixture pair, observe `dev#k6` / `user-testing#foo`; repeat over a single workflow, observe the bare "ticket k6". Cost: minutes. Qualifier values come from the fixture dir basenames.
- **Estimated total cost/complexity: LOW–MEDIUM.** Two template-prose edits (present-gate SKILL.md, the commander-dispatch template) + a one-line resolve instruction each; two small fixtures; one reused `dispatch build` golden assertion; two live drives (gate + commander) plus one two-workflow drive. No binary/Go changes (Layer 2 is out of scope).

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
