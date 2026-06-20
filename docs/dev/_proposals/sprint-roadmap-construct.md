# Proposal: sprint / roadmap construct — ship as skill + commission template, defer builtin

Status: proposal / decision record (2026-06-08). Provenance: an adversarial design workflow (5 advocates → 5 adversarial critics → deciding synthesis) over spacedock-gym's roadmap/sprint convention, plus a cross-project scan of `~/git/spacedock-research/*`. Nothing is shipped by this doc; it records the conclusion so it is not lost.

## The construct (observed in spacedock-gym)

A program/strategy layer ABOVE Spacedock's per-entity workflow:

- **Roadmap** — durable strategy + a value-ordered sprint sequence (`docs/roadmap/README.md`). Owns outcome/scope/sequencing/DoD; does NOT track task state.
- **Sprint** = a vertical value-increment that GROUPS entities by frontmatter query (`sprint:` + a `sprint-readiness:` filter), never a hard-coded list. Folder `NNN-<slug>/`: `index.md` (goal/DoD/deliverable), `staff-review.md` (readiness gap analysis), `dispatch-sprint-execution.md` (cold-boot Commander package).
- **Two-tier orchestration** — a **Shaping FO** (strategy, sprint definition, gating ideation, cross-entity coherence check, staff review, packaging) and a **Commander** (takes ONE packaged sprint, boots `spacedock:first-officer`, drives it to its deliverable; escalates only on a 3rd feedback cycle / budget blowout / irrecoverable block / scope fork).
- **Sprint-wide proof** — DoD + integration test + report at the sprint level, above per-entity validation.

It is ~90% prose / convention / frontmatter.

## Decision: contrib skill + commission template — NOT a builtin core construct

Scores (0-10): hybrid (skill + template) **6.5 (winner)** · skill-only 5.5 · defer 5 · template-only 4 · builtin **2.5**.

When pursued, ship:

1. A `spacedock:roadmap` (or `commander`) **contrib skill** carrying the Shaping-FO/Commander discipline, composing with `spacedock:first-officer` (the Commander boots the FO). Built by a dispatched worker under test (it touches guarded product scaffolding).
2. A **commission template** scaffolding the `NNN-<slug>/` folder shape + the `sprint` / `group` / `sprint-readiness` frontmatter, and documenting that membership is a native query.

Do NOT:

- Bump `CONTRACT_VERSION`, add a `sprint` frontmatter recognizer, or add a `--sprint-validate` DoD gate. The cross-entity DoD predicate is domain-specific and binds to FO-authored frontmatter, so a generic gate is tautological (fails the proof-policy — a guard must bind to an independent, divergence-capable source).
- Reinvent membership: the roll-up the builtin case wanted ALREADY EXISTS — `spacedock status --where sprint=<slug> --where 'sprint-readiness != defer' --fields group` (verified live during the design workflow).

## Why defer builtin (the evidence)

- **N=1.** The full construct exists in **spacedock-gym only**. Cross-project scan (2026-06-08): `spacedock-landing` has no roadmap/sprints; `spacedock-v1`'s `docs/roadmap/` is just `bootstrap-roadmap.md` — a strategy doc, not the NNN-slug sprints + `sprint:` frontmatter + Commander construct.
- Baking an N=1, design-phase convention into core ossifies a moving target and imposes a sprint model on every workflow that does not need one.
- It cuts against the codebase's "prove a pattern in prose/skill before mechanizing" culture (cf. the 4q net-removal of over-built standing machinery).

## Graduation triggers (promote a piece to the binary only when these hold)

- **DoD/cross-entity guard:** a delivered sprint exposes a divergence-capable, project-agnostic predicate (a CI exit code / artifact hash — NOT "members terminal-with-verdict," which only binds to sibling frontmatter), recurring across ≥2 sprints AND a second project.
- **Frontmatter recognizer:** the `sprint`/`group`/`sprint-readiness` triplet stabilizes unchanged across 2-3 delivered sprints + a second project, AND a concrete need arises that `--where` cannot serve.
- **Shaping/Commander roles into the shipped contract (CONTRACT_VERSION bump):** only with a new binary-observable surface AND the role boundary proven by a Commander actually driving a sprint to a deliverable.

## Next step

When prioritized (post-flip), file a build task for the `spacedock:roadmap` skill + the commission sprint template. This proposal is the decision record only.
