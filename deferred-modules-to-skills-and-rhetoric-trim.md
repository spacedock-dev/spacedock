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

## Design

### Part A — the skill conversion

`fo-status-viewer.md` and `fo-write-core.md` become non-user-invocable skills at `skills/fo-status-viewer/SKILL.md` and `skills/fo-write-core/SKILL.md` (mirroring `present-gate` / `feedback-rejection-flow`: `user-invocable: false`, a `description` frontmatter that carries the when-to-load, a `# Title` header, then the body verbatim). The reference files are deleted; their section bodies (`## Status Viewer` / `### Captain-Facing State Display` / `## Issue Filing`; `## FO Write Scope` / `## ID Styles`) move into the SKILL.md unchanged. `fo-dispatch-core.md` and `fo-merge-core.md` stay references — they compose with a per-host runtime-adapter section (dispatch/terminal-teardown), which is the captain's discriminator: **adapter-less → skill; host-composed → reference.**

Names stay `fo-status-viewer` / `fo-write-core` (matches the seed, the re-pointed AC-2 assertion, and the closure oracle's resolution; the `fo-` prefix legibly scopes them as first-officer-internal, distinct from user skills). Naming churn is out of scope.

Concrete before/after — the boot-core pointers turning into skill invocations:

**fo-status-viewer header (`fo-status-viewer.md:1-3`) → `skills/fo-status-viewer/SKILL.md` frontmatter.** The self-describing paragraph is duplicated by the skill's own `description`, so it is cut, not moved.
BEFORE:
```
# First Officer Status Viewer (host-neutral)

The `status` query/mutate command surface — the launcher invocation, ... plus the GitHub-issue-filing approval gate. Lazily loaded (named by the boot-resident core) at the FIRST ad-hoc status question, `--set` mutation, `--next-id`/`--resolve` lookup, or issue filing; a greet-and-stop boot never reads it. Host-neutral: pure `status`-command and filing reference, identical on Claude, Codex, and Pi, with no per-host adapter.
```
AFTER:
```
---
name: fo-status-viewer
description: "First-officer status query/mutate/display surface — the `status` command flag docs, `--set` field docs, canonical captain-facing invocations, the Captain-Facing State Display rendering, and the GitHub-issue-filing approval gate. Invoke at the first ad-hoc status question, `--set` mutation, `--next-id`/`--resolve` lookup, or issue filing."
user-invocable: false
---

# First Officer Status Viewer
```
(`fo-write-core.md:1-7` header → analogous frontmatter; `## FO Write Scope` / `## ID Styles` bodies verbatim.)

**Deferred-modules registry (`first-officer-shared-core.md:35-44`).** The four-entry `module → realization → core file → load-point → greet-guard` block collapses: the two skills leave the registry (skill metadata carries when-to-load), the greet-guard is stated ONCE not per-row, and the arrow-taxonomy (`→ host-neutral (no per-host adapter) →` vs `→ runtime-binding →`) dissolves — the realization is self-evident from the invocation shape (`Skill(...)` vs `references/*.md`), so no label is needed. This is where the seed's "84/k4/z4 arrow-taxonomy punt" finally resolves.
BEFORE: the `## Deferred Modules (registry)` block, its greet-guard/each-row-narration paragraph, four `→`-taxonomy bullets, and the "Doubles as the navigation index" line (~1.3 KB).
AFTER:
```
## Deferred load points

A greet-and-stop boot loads NONE of these — it composes its summary from `«state.boot»` JSON + README frontmatter (Startup step 8) and presents any ready gate via `present-gate`. Each loads only at its trigger:

- `Skill(skill="spacedock:fo-status-viewer")` — first status query (`--set` / `--next-id` / `--resolve` / issue filing).
- `Skill(skill="spacedock:fo-write-core")` — first write to main (`status --set`, `spacedock new`, archive move, `### Feedback Cycles` write).
- `references/fo-dispatch-core.md` — first worker dispatch.
- `references/fo-merge-core.md` — terminal boundary.
```

**State-Management pointer (`first-officer-shared-core.md:102`).**
BEFORE: `... full write-authority scope in the deferred write reference `references/fo-write-core.md`, indexed in the Deferred Modules registry).`
AFTER: `... full write-authority scope in `Skill(skill="spacedock:fo-write-core")`, loaded at first write to main).`

The two `Skill(skill="spacedock:...")` tokens keep the load points reachable through the SAME closure mechanism that resolves `present-gate` today (`bodySkillRe` → `skills/<name>/SKILL.md`).

### Part B — the rhetoric trim (concrete hit-list)

| Rhetoric | Sites | Action |
|---|---|---|
| `(host-neutral)` header label | `fo-status-viewer.md:1`, `fo-write-core.md:1`, `fo-dispatch-core.md:1`, `fo-merge-core.md:1` | CUT from all four headers |
| `host-neutral` self-narration | `fo-status-viewer.md:3` ("Host-neutral: ... no per-host adapter"), `shared-core:37,39,41`, `claude-runtime:7,13` | CUT (skills) / fold into terse index (dispatch, merge) |
| Self-describing "Lazily loaded ... a greet-and-stop boot never reads it" headers | `fo-status-viewer.md:3`, `fo-write-core.md:3-7`, `fo-dispatch-core.md:3`, `fo-merge-core.md:3` | CUT for skills (metadata carries it); trim dispatch/merge to a one-line "what" |
| `(named by the boot-resident core)` parentheticals | `fo-status-viewer.md:3`, `fo-dispatch-core.md:3`, `fo-write-core.md:5`, `fo-merge-core.md:3`, `claude-runtime:7,13` | CUT (each file narrating who named it adds nothing) |
| `## X (deferred module)` heading ceremony + arrow/done-when/guard triplet | the `## Deferred Modules (registry)` block, `fo-write-core.md:5` back-reference | Collapse per Part A |
| `lazily-loaded` / `boot-resident` self-narration | `shared-core:3`, `claude-runtime:3,7,13` | Trim to the kernel (keep the greet-guard once, drop the repeated "neither is read at boot" restatements) |
| internal cross-refs to the moved files | `fo-status-viewer.md:17`, `claude-runtime:35`, `feedback-rejection-flow/SKILL.md:23` | re-point `references/fo-write-core.md` → the fo-write-core skill |

**Sweep result (both cores).** The **ensign core is clean** — `grep` for `host-neutral` / `deferred` / `lazily` / `boot-resident` / `named by` over `ensign/references/*` returns nothing; no ensign-side rhetoric to cut. The FO codex/pi runtime adapters carry a mild "This file defines how the shared first-officer core executes on Pi/Codex. The shared core owns invocation timing; this adapter owns the X realization." opener — this is BORDERLINE (a per-host reader is oriented by it). Recommendation: keep a one-line purpose, drop only the redundant "The shared core owns invocation timing; this adapter owns the X realization" meta-sentence. `claude-first-officer-runtime.md:3` ("It is the boot-resident runtime adapter — ...; the dispatch and merge machinery live in lazily-loaded references named below; neither is read at boot") is the FO-side instance of the same self-narration and is trimmed.

### Keep (load-bearing — do NOT cut)
- The **greet-guard** ("a greet-and-stop boot loads none of these") — the live shallow-boot lane's actual assertion. Stated ONCE in the `## Deferred load points` block, not restated per file.
- The **load-point trigger** — carried by skill `description` (skills) or the terse index line (dispatch/merge references).
- `fo-dispatch-core` / `fo-merge-core` as references — they genuinely compose with a per-host adapter section.

## Spike — mechanism validation (PROVEN, no further spike needed)

The riskiest unverified claim is AC-2's re-point: *does a non-user-invocable skill load at the FO's trigger, and is its `Skill(...)` invocation observable in the transcript so the greet-independence lane can assert on it?* Exercised end-to-end against real, shipped artifacts and transcripts:

1. **A non-user-invocable skill loads at the FO's trigger via `Skill(skill="spacedock:...")`** — already in production. `present-gate` (`user-invocable: false`) is invoked via `Skill(skill="spacedock:present-gate")` in the shipped `first-officer-shared-core.md:32,88` and appears in REAL captured session transcripts (`~/.claude/projects/.../181055aa-….jsonl`, `018c3c8d-….jsonl`). Not asserted — demonstrated.
2. **The `Skill(...)` invocation is observable with its skill NAME.** A real `Skill` tool_use block is `{"type":"tool_use","name":"Skill","input":{"skill":"spacedock:present-gate"}}` (verified in fixture `sonnet_teamdelete_marker_continues.stream.jsonl` and live transcripts). `journeymetrics.ParseClaudeTurns` already records `block.Name` into `ClaudeTurn.ToolNames`; the `skill` argument is extractable with one line — `jsonStringField(input, "skill")` — exactly as `readToolTarget` extracts `file_path` for `Read` today.
3. **Non-user-invocable skill descriptions ARE surfaced to the agent's skill catalog** (so the FO knows to invoke at the trigger, and the `description` genuinely carries the when-to-load). Proven from this very session: `spacedock:present-gate` and `spacedock:feedback-rejection-flow` appear in the available-skills list despite `user-invocable: false`; `user-invocable: false` blocks the USER typing `/name`, not agent-side `Skill()` invocation or catalog visibility.
4. **Sharp constraint the spike surfaced:** the greet legitimately invokes `Skill(skill="spacedock:present-gate")` at the greet (Startup step 8 / `runClaudeShallowBootScenario` presents a ready gate). So AC-2's assertion CANNOT be "no `Skill` call before greet" — it MUST key on the skill ARGUMENT (`fo-status-viewer` / `fo-write-core`), with `present-gate` explicitly allowed. This is why the `skill`-argument extraction (not mere tool-name detection) is load-bearing.

The throwaway extraction seeds the implementation's first test (the `jsonStringField(input,"skill")` extractor + its RED fixture).

## Acceptance criteria

- **AC-1 (VALUE — measured byte delta vs the post-0240 `main` baseline that can move the wrong way).** The cumulative signed `wc -c` delta over the affected contract-file SET is NEGATIVE. The set: `first-officer-shared-core.md`, `claude-first-officer-runtime.md`, `fo-dispatch-core.md`, `fo-merge-core.md` (edited in place), `fo-status-viewer.md` + `fo-write-core.md` (→ 0, deleted), and the two new `skills/{fo-status-viewer,fo-write-core}/SKILL.md` (added). Measured with the 0240 DoD methodology — `git show <baseline-sha>:<path> | wc -c` per baseline file vs working-tree `wc -c`, summed. The stage report records every per-file figure and the signed cumulative delta. This is the value gate: the skill frontmatter is byte-ADDING (~450-500 B each), so the net-negative must come from the rhetoric cut EXCEEDING that overhead — a feasibility estimate puts the cut near −1.9 KB (registry collapse −0.85 KB, four self-describing headers −0.7 KB net of frontmatter, host-neutral/named-by parentheticals −0.35 KB), but the AC measures it, it does not assume it. **Test:** shell one-liner over the file set, asserted `< 0`; figures in the stage report.
- **AC-2 (VALUE/behavioral — greet-independence re-pointed).** On the live shallow-boot lane, no pre-greet turn invokes `Skill(skill=…fo-status-viewer…)` or `Skill(skill=…fo-write-core…)`; `present-gate` remains allowed pre-greet. The old "no Read of `fo-status-viewer.md`" oracle is retired (its target file no longer exists) and replaced by a skill-invocation oracle. **Test:** a new `assertGreetInvokesNoDeferredFOSkill(stream)` walking a new `ClaudeTurn.SkillNames` field, wired into `claude_live_runner_test.go` (replacing the `assertGreetReadsNoDeferredStatusReference` call at :346), plus offline RED/GREEN unit fixtures — a pre-greet `Skill(skill="spacedock:fo-status-viewer")` on a LATER delta must go RED (mirroring the existing multi-delta negative fixture), and the clean shallow-boot fixture stays GREEN. `greet-reads-status-viewer.stream.jsonl` is replaced by a `greet-invokes-fo-status-viewer-skill.stream.jsonl` negative fixture.
- **AC-3 (mechanism — reachability closure; paired with AC-1 + AC-2).** No behavior loss: every moved rule resolves at its trigger, proven by the contractlint closure suite going GREEN with the two modules resolved as SKILLS — a structural os.Stat + skill-resolution check, NOT a prose-grep. This is a mechanism-only AC and counts only paired with the value it serves (AC-1's leaner contract, AC-2's still-independent greet). **Test:** `go test ./internal/contractlint/`, specifically `boot_resident_closure_test.go` updated so `fo-status-viewer`/`fo-write-core` move out of `foReferenceCores` into `lazyLoadSkills` (the shared core names them via `spacedock:…`, `bodySkillRe` resolves them to `skills/<name>/SKILL.md`), a skill-anchor check confirms each SKILL.md still carries its ceremony sections (`## Status Viewer`/`### Captain-Facing State Display`/`## Issue Filing`; `## FO Write Scope`/`## ID Styles`), and `deferredReferenceFiles` is re-targeted to the skill files so the prose-pointer dangler guard still walks them. The existing dangling-target and prose-pointer CONTROLS must stay RED-capable.
- **AC-4 (regression gate).** `go build ./...` and `go test ./internal/contractlint/ ./internal/ensigncycle/` are green from the root — including the re-pointed greet-independence lane (AC-2), the closure suite (AC-3), and the measured shallow-boot window (`assertShallowBootMeasured`, greet-turn context stays below the 60k ceiling; the two short skill descriptions added to the always-on catalog are ≈200 tokens, negligible against the ceiling). **Test:** the two `go test` invocations green; the measured-window oracle unchanged and passing.
- **AC-5 (authority-doc alignment).** `docs/runtime-support.md` states that an adapter-less deferred module realizes as a non-user-invocable skill (invoked via `Skill(skill="spacedock:…")`), while a runtime-binding deferred module (dispatch/merge) stays a core file the boot core names by path — so the `→ runtime-binding` / deferred-module convention text (`:21`) does not drift from the shipped two-skills-two-references shape. **Test:** cross-check the documented convention against the shipped skills/references; a one-paragraph doc edit.

Mechanism→value pairing (for `«gate.ac-cross-check»`): AC-3 (closure ships) serves AC-1 (leaner) + AC-2 (greet-independent); AC-3 passes only if AC-1's delta is NEGATIVE and AC-2 is GREEN. A closure that resolves but a contract that GREW is a REJECT.

## Test plan
- **Offline Go unit (minutes):** the `jsonStringField(input,"skill")` extractor + `ClaudeTurn.SkillNames`; the `assertGreetInvokesNoDeferredFOSkill` oracle with RED (pre-greet `Skill(fo-status-viewer)` on a later delta) and GREEN (clean shallow boot; a pre-greet `Skill(present-gate)` must NOT trip it) fixtures; the `boot_resident_closure_test.go` edits (`lazyLoadSkills` += the two skills, skill-anchor check, `deferredReferenceFiles` re-target) with the existing controls held RED-capable. This is the riskiest-path-first slice.
- **Live workflow smoke (multi-minute, already in the suite):** re-run the `shallow-boot` shared scenario with AC-2's re-pointed assertion wired in; confirm the greet invokes no FO deferred skill and the merged/gate durable-state grade (`assertShallowBoot`) + measured-window (`assertShallowBootMeasured`) stay green. Positive half ("a status path DOES load fo-status-viewer") is covered by the closure resolution (AC-3) plus the spike's transcript evidence; a dedicated live status-query micro-scenario is optional and only if cheap.
- **Byte measurement (minutes):** the AC-1 shell one-liner over the file set, signed cumulative delta reported.
- **No live cost for AC-1/AC-3/AC-5** (offline/measurement); one live shallow-boot re-run for AC-2/AC-4.

## Doc / user-visible surface
No user-visible CLI output, startup banner, or host-integration behavior changes — this is contract scaffolding (skills/references) + a test-harness oracle. The only doc touched is the dev-facing authority doc `docs/runtime-support.md` (AC-5). No user-docs-site diff is required.

## Staff review — WARRANTED
This is a contract restructure (the deferred-modules registry + boot-core pointers) plus skill integration (two new skills, closure-oracle surgery, a re-pointed live lane) — two of the README's named complexity triggers ("skill integration", contract restructuring). The FO should request an independent review before presenting the ideation gate: design soundness, whether the closure-oracle re-registration truly preserves reachability (not just resolves paths), and that the riskiest mechanism (AC-2's skill-invocation observability) was exercised first — which it was, against real transcripts (see Spike).

## Related
- Skill precedent: `skills/present-gate/SKILL.md`, `skills/feedback-rejection-flow/SKILL.md` (non-user-invocable, FO-invoked).
- Shipped reference form to convert: `skills/first-officer/references/fo-status-viewer.md` (84), `fo-write-core.md` (k4).
- z4's four-entry `## Deferred Modules (registry)` block + the arrow taxonomy — both simplify/dissolve once adapter-less modules are skills.
- Sibling backlog: `next-post-release-preversion-bump`, `codex-dev-plugin-launch-ergonomics`.

## Stage Report: ideation

- DONE: Flesh the seed into a gate-ready design — skill-conversion + rhetoric-trim hit-list, concrete before/after for the boot-core pointers turning into skill invocations, sweep of both cores
  `## Design` (Part A conversion + Part B hit-list table with sites/actions) and the boot-core `## Deferred load points` before/after; sweep result recorded (ensign core clean, FO codex/pi openers flagged borderline).
- DONE: Measured ACs, each with its test — (a) net-NEGATIVE byte delta, (b) greet-independence re-pointed to no-Skill-before-greet, (c) no behavior loss proven behaviorally
  `## Acceptance criteria` AC-1 (signed `wc -c` delta over the affected set, 0240 DoD methodology), AC-2 (`assertGreetInvokesNoDeferredFOSkill` live + RED/GREEN fixtures), AC-3 (contractlint closure GREEN, os.Stat + skill-resolution, not prose-grep); mechanism→value pairing stated for `«gate.ac-cross-check»`.
- DONE: Spike the riskiest mechanism FIRST + flag staff review
  `## Spike` — PROVEN against real shipped artifacts and transcripts (present-gate is `user-invocable:false` yet invoked via `Skill()` and observable as `{"name":"Skill","input":{"skill":"spacedock:present-gate"}}` in live jsonl; `ParseClaudeTurns` already captures the name, the skill arg is one `jsonStringField` away). Staff review flagged WARRANTED (`## Staff review`).

### Summary
Designed the two-part change: (A) convert the adapter-less `fo-status-viewer` / `fo-write-core` references into `user-invocable:false` skills whose `description` carries the when-to-load, keeping host-composed `fo-dispatch-core` / `fo-merge-core` as references; (B) cut the self-referential rhetoric (host-neutral labels, self-describing headers, named-by-boot-resident parentheticals, the four-entry deferred-modules registry + arrow-taxonomy), collapsing the boot pointers to a single `## Deferred load points` block with `Skill(skill="spacedock:…")` invocations. The riskiest mechanism (AC-2's re-point to a skill-invocation transcript oracle) was exercised first and proven from real transcripts — no further spike needed; the `jsonStringField(input,"skill")` extractor seeds implementation's first test. Key decision: AC-2 must key on the skill ARGUMENT, not tool-name, because the greet legitimately invokes `present-gate`. Staff review is warranted (contract restructure + skill integration).
