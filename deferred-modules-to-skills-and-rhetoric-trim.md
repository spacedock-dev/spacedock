---
title: Turn adapter-less deferred modules into non-user-invocable skills + cut self-referential contract rhetoric
status: ideation
source: "Captain critique 2026-06-30/07-01 (0240 Commander session, post-lean-contract). The lean-contract deferrals delivered the real win (boot core -4402 B vs v0.22.0) but wrapped it in self-referential ceremony: every deferred reference announces itself as a `## X (deferred module)` with a arrow/done-when/guard triplet, re-explains itself in a file header, and both are labelled `(host-neutral)` — a label about the contract's own architecture, meaningless to FO behavior. Separately: the adapter-less deferred modules (fo-status-viewer, fo-write-core) are `references/*.md` the FO must recall-a-path-and-Read at a trigger, when they are architecturally identical to `present-gate`/`feedback-rejection-flow` (host-neutral, FO-invoked, not user-typed) — they should be non-user-invocable SKILLS, whose SKILL.md metadata carries the when-to-load the pointer prose duplicates."
group: tooling
id: nt982bbkf04r0ypbsc8er74s
started: 2026-07-01T00:28:25Z
---

## Problem — two intertwined issues, one fix

**A. Adapter-less deferred modules should be skills, not read-references.** `present-gate` and `feedback-rejection-flow` are non-user-invocable skills the FO invokes via `Skill(skill="...")` at a trigger point. `fo-status-viewer.md` (84) and `fo-write-core.md` (k4) are the SAME kind of thing — host-neutral, FO-invoked, not user-typed, adapter-less — yet shipped as `references/*.md` the FO must recall-a-path-and-`Read`. The captain's discriminator: **runtime-overridable (host-composed, pairs with a per-host adapter) -> stays a reference; not-runtime-overridable (host-neutral, adapter-less) -> becomes a skill.** So: fo-status-viewer + fo-write-core -> skills; fo-dispatch-core + fo-merge-core -> stay references (they compose with a runtime-adapter section).

A skill's own metadata (the SKILL.md description/frontmatter) carries the "when to load" — which is exactly what the `## X (deferred module)` pointer prose duplicates today. So the pointer AND the self-describing header both disappear; `Skill(skill="...")` becomes the observable trigger boundary (the greet-independence live lane asserts "no Skill(fo-status-viewer) before greet" — sharper than "no Read of a path"). It also dissolves the `-> reference` vs `-> runtime-binding` arrow-taxonomy question 84/k4/z4 kept punting.

**B. Cut the self-referential rhetoric** the deferrals accreted (concrete instances on `main` after 0240):
- `(host-neutral)` labels — 10 occurrences (every reference header `# First Officer X Core (host-neutral)`; every boot-core pointer `(host-neutral; no per-host adapter)`). A label about the contract's own architecture; meaningless to FO behavior. CUT.
- Self-describing file headers — every reference re-explains itself in a header paragraph that restates the boot-core pointer verbatim ("Lazily loaded (named by the boot-resident core) at the FIRST ...; a greet-and-stop boot never reads it. Host-neutral: ..."). CUT (skills carry it in metadata; residual references keep at most a one-line "what").
- `(named by the boot-resident core)` cross-reference parentheticals — each file narrating who named it. CUT.
- `## X (deferred module)` heading ceremony + the per-module arrow/done-when/guard triplet — replaced by skill metadata (skills) or collapsed to one terse index line (dispatch/merge references); z4's four-entry registry simplifies with it.
- `no per-host adapter` / `lazily-loaded` / repeated `boot-resident` self-narration. Trim to the kernel.
- Ideation sweeps for more of the same across BOTH the FO and ensign cores.

## Keep — the load-bearing kernel (do NOT cut)
- The **greet-guard** ("a greet-and-stop boot reads none of these") — an actual assertion the live shallow-boot lane checks (84 AC-2), not decoration. Keep it in ONE place, not restated per file.
- The **load-point trigger** (when to load) — needed by a cold FO; carried by skill metadata (skills) or a terse index line (references).
- The host-composed cores (fo-dispatch-core, fo-merge-core) as references — they genuinely compose with a per-host runtime adapter.

## Rough acceptance sketch (ideation tightens into measured ACs + a test plan)
- fo-status-viewer + fo-write-core are non-user-invocable skills invoked via `Skill(skill="...")` at their trigger; the `## X (deferred module)` pointers + self-describing headers are gone. Verified behaviorally: a greet-and-stop live boot invokes NEITHER skill (greet-independence lane re-pointed from "no Read of path" to "no Skill(...) before greet"), and a status-query / write path DOES load them.
- Measured net-NEGATIVE contract-byte delta from the rhetoric cut (file-delta vs post-0240 main, the 0240 DoD methodology).
- No behavior loss: every moved rule still resolves at its trigger (skill/reference carries it verbatim).

## Related
- Skill precedent: `skills/present-gate/SKILL.md`, `skills/feedback-rejection-flow/SKILL.md` (non-user-invocable, FO-invoked).
- Shipped reference form to convert: `skills/first-officer/references/fo-status-viewer.md` (84), `fo-write-core.md` (k4).
- z4's four-entry `## Deferred Modules (registry)` block + the arrow taxonomy — both simplify/dissolve once adapter-less modules are skills.
- Sibling backlog: `next-post-release-preversion-bump`, `codex-dev-plugin-launch-ergonomics`.
