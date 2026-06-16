---
title: Model-aware tier self-identification + level-3 delegation routing
status: ideation
source: 'FO shaping 0205 (2026-06-15, this session) — for a weak FO to be safe by construction it must self-identify its tier and route every judgment call to a stronger standing teammate, never adjudicate alone. No such mechanism exists. The gate-half (route all gate verdicts to level-3 when the FO is Haiku and the workflow has gate:true stages) is the no-2y startable slice; the full seam visibility composes with 2y. 0205 layered-FO.'
started: 2026-06-16T02:21:21Z
completed:
verdict:
score: 0.5
worktree:
issue:
sprint: 0205-layered-fo
sprint-readiness: ready
id: 72r2x0nnvx9az9x1adf08svq
---

The mechanism that makes a weak FO safe by construction: it knows it is weak and escalates structurally, never by luck.

## Problem

A Haiku FO has no way to know it should not make a gate verdict, a scope call, or a conflict-recovery decision. The contract today never lets the FO read its own model: boot reads contract version, project root, workflow frontmatter, and `status --boot --json` — none carries the FO's resolved model (`skills/first-officer/references/first-officer-shared-core.md:9-29`; the launcher passes `model=` INTO the FO's `Agent()` dispatch but there is no post-hoc query for it). The FO can read a *worker's* stamped model via `spacedock dispatch context-budget --name {agent}` (`claude-fo-dispatch.md:163`), but never its own.

Without self-identification, the safety of a weak FO depends on the model happening to escalate. The gate flow is where this bites hardest: `present-gate`'s `Recommend {approve | reject}` line (assembly rule 5) is the FO's verdict, written by whatever model is driving. A Haiku FO fills that line itself today. The fix is structural: the FO self-identifies its tier at boot and a named routing table sends each judgment category to a standing level-3 teammate on a stronger model, so the verdict on a `gate: true` stage is *authored* by level-3, not by Haiku.

The layered-FO vision (`docs/roadmap/fo-layered-architecture.md`) already names the three layers — L1 automation (binary), L2 driving (run the contract by the book), L3 judgment (override automation when warranted) — and the principle "a weak model self-identifies as level-2-only and routes all judgment to a higher-model teammate." This task makes that principle real for the one judgment category that is startable without 2y: gate verdicts.

## Proposed approach

Four pieces, expressed in the established prose-function notation (`docs/dev/_proposals/state-repo-pointer-prototype.md`): `«fn»(arg)` guillemets = a body the FO follows by hand; backticks = a shipped binary it calls; each body names the verb it becomes.

**1. Tier self-identification — one source of truth at boot.** `«fo.tier»()` resolves the FO's tier from a single launcher-set source, not a self-probe (a model asked "what model are you" is unreliable) and not a README declaration (the tier is a property of the running session, not the workflow). The launcher (`spacedock claude --model haiku`) already knows the model it dispatched the FO with; it exports `SPACEDOCK_FO_MODEL` into the FO's environment. `«fo.tier»()` reads it and maps to a tier: `haiku → level-2-only`; `sonnet | opus | <unset> → level-3-capable` (unset defaults to capable — a hand-launched FO on a capable session is the common case and must not be hobbled). The map lives in one place so "which models are level-2-only" is data, not scattered prose.

**2. A named routing table — judgment categories → level-3 route.** A single table maps each judgment category to its delegation route:

| Judgment category | Route when level-2-only |
| --- | --- |
| gate verdicts (approve/reject) | level-3 authors the `Recommend` line; FO forwards it |
| design / scope calls | level-3 |
| feedback-cycle-3 escalation | level-3 (before escalating to captain) |
| model-mismatch reuse decision | level-3 |
| rebase-conflict / state recovery | level-3 (the halt itself is mechanical; whether to deviate is judgment) |
| teardown health (shutdown-vs-keep) | level-3 |

The table is the full seam. The **gate-verdicts row is the only one this slice wires live** (the no-2y startable slice); the other five rows ship as the documented table with their routes named but their *live wiring* deferred — they compose once 2y makes the mechanical/judgment seam visible in the restructured dispatch/merge cores. The table is authored now so the seam is legible and the remaining rows are a wiring exercise, not a redesign.

**3. A standing level-3 teammate, spawned on demand.** A `level-3-judge` mod under `{workflow_dir}/_mods/`, declared exactly like `comm-officer` (`standing: true`, a `## Hook: startup` with `model:` and `subagent_type:`, an `## Agent Prompt`). It declares `model: opus` (the stronger model). It is spawned lazily by the existing `spacedock dispatch spawn-standing-all` at first team dispatch (`claude-fo-dispatch.md:21`) — no new spawn machinery. The FO routes to it via `SendMessage(to="level-3-judge", ...)`. Unlike `comm-officer` (best-effort, non-blocking, proceed-if-absent), the level-3 route on a gate verdict is **load-bearing and blocking for a level-2-only FO**: if the teammate is absent or silent, the FO does NOT fabricate the verdict — it surfaces the gate to the captain with "level-3 unavailable, no FO verdict" rather than self-approving. This is the one place the standing-teammate contract differs from prose-polish, and the difference is the whole safety property.

**4. A boot gate.** At the end of `«fo.tier»()`, if `tier == level-2-only` AND the workflow has any `gate: true` stage (already read into boot frontmatter at Startup step 4), the FO emits a boot-line stating that gate verdicts will route to `level-3-judge`, and arms the gate-verdict routing for the session. A level-3-capable FO arms nothing — the table is inert and the gate flow is unchanged. This is the conditional that keeps the mechanism free for the common (capable-FO) case.

**Durable record of an L3-authored verdict.** For the AC to be gradable on durable state (never transcript phrasing), an L3 gate verdict leaves an on-disk mark. The FO appends a `### Gate Verdicts` section to the entity body (mirroring how `### Feedback Cycles` is FO-owned and worktree-side-when-`worktree:`-set), one line per gated stage: `- {stage}: {approve|reject} — decided-by: level-3-judge`. This is grep-able durable state that proves authorship, and it composes with the haiku-loop-spike's terminal-state oracle (the spike greps this section plus the terminal `verdict:` frontmatter). The FO never writes `decided-by: {self}` on a level-2-only session — absence of an L3 line on a gated stage is itself a detectable violation.

### Prose-function declaration (the form the contract restructure will carry)

## «fo.tier»(): self-identify the FO tier at boot and arm level-3 gate routing
- **effect:** read `SPACEDOCK_FO_MODEL` (launcher-set); map model→tier via the tier map (`haiku → level-2-only`; else `level-3-capable`); if `level-2-only` ∧ workflow has a `gate: true` stage, arm gate-verdict routing to `level-3-judge` and emit the boot line
- **done-when:** the FO tier is resolved and (for a level-2-only FO meeting gated stages) gate-verdict routing is armed
- **block:** none — unresolvable model defaults to `level-3-capable` (fail open to the safe-by-luck status quo, never fail into a hobbled capable FO)
- → **prose**, becomes `` `spacedock fo tier` `` (a verb that reads the launcher env and the workflow frontmatter and returns `{tier, route_gate_verdicts: bool, gate_route: "level-3-judge"}`)

## «gate.route-verdict»(slug, stage): obtain the gate verdict from the armed route
- **guard:** gate-verdict routing is armed (a level-2-only FO at a `gate: true` stage)
- **effect:** `SendMessage(to="level-3-judge", ...)` with the stage report + AC section + checklist accounting; receive the `{approve|reject, reason}`; write the `### Gate Verdicts` line `decided-by: level-3-judge`; fill `present-gate`'s `Recommend` line from the returned verdict
- **done-when:** the `Recommend` line is filled from an L3-authored verdict and the `### Gate Verdicts` line is committed
- **block:** teammate absent/silent past timeout → surface the gate to the captain with "level-3 unavailable" and NO FO verdict; never self-author the `Recommend` line on a level-2-only session
- → **prose** (stays prose this sprint; the verdict is judgment, never a binary)

When the FO is `level-3-capable`, `«gate.route-verdict»` is never invoked — the gate flow runs exactly as today (the FO authors its own `Recommend` line), so this mechanism is a no-op for capable sessions.

## Out of scope

- **Live-wiring the five non-gate routing-table rows** (design/scope, feedback-cycle-3, model-mismatch reuse, conflict recovery, teardown health). They ship as the documented table with named routes; their live wiring composes with 2y, which makes the mechanical/judgment seam visible in the restructured dispatch and merge cores. This slice wires only the gate-verdicts row.
- **A separate live Haiku drive.** The riskiest mechanism (does tier self-ID + routing hold in a live Haiku drive?) is exercised in `haiku-loop-spike` — this task coordinates its proof there, not in a second drive (see Spike determination).
- **The `spacedock fo tier` binary verb.** This slice ships the prose-function `«fo.tier»()` (hand-followed) with its named migration target; building the verb follows the spike, like the other 0205 verbs.
- **Codex / Pi substrates.** No team registry on Codex/Pi and different teammate models; their layered-FO operability is a 0205-scoped follow-up, Claude only here.
- **Making the `level-3-judge` model configurable per workflow.** It declares `model: opus` in its mod; per-workflow model selection is a later refinement, not this slice.

## Spike determination

**No separate spike needed for this task: the riskiest mechanism is exercised in `haiku-loop-spike` (`w4ryf4mg4vn1emwp906vd8yp`).** The one claim that would invalidate this design — that a live Haiku FO can self-identify its tier and route a gate verdict to a level-3 teammate such that the verdict is L3-authored, not Haiku-authored — is precisely that spike's AC-1 (one real entity driven dispatch→gate→merge→terminal with the gate verdict made by level-3, graded on durable end-state). This task's AC-1 is satisfied by that same drive observing this task's durable artifact (the `### Gate Verdicts` `decided-by: level-3-judge` line). Designing a second drive would duplicate the spike. The proven-mechanism this design rests on that is NOT the spike's concern — that a launcher can export an env var the FO reads at boot — is a trivial, already-proven mechanism (`SPACEDOCK_BIN` is exported the same way today, `first-officer-shared-core.md:7`), exercised by a fixture test in the test plan below, not a live drive.

## Acceptance criteria

**AC-1 — A level-2-only (Haiku) FO at a `gate: true` stage obtains the gate verdict from `level-3-judge` and the durable entity state shows the verdict was authored by level-3, never by the Haiku FO.**
Verified by: the haiku-loop-spike live drive (coordinated, not duplicated) leaving, in durable workflow state, a `### Gate Verdicts` section with a `decided-by: level-3-judge` line for the gated stage, plus the terminal `verdict:` frontmatter — asserted by grepping the archived entity for `decided-by: level-3-judge` on the gated stage and the absence of any self-authored verdict line, never by reading the drive transcript.

**AC-2 — `«fo.tier»()` resolves the correct tier and arms gate-verdict routing only for a level-2-only FO meeting `gate: true` stages, and is inert for every other combination.**
Verified by: a fixture-driven test of the tier-resolution + arming decision over the four cases — (`SPACEDOCK_FO_MODEL=haiku` + workflow with a gated stage → armed, boot line emitted), (`haiku` + no gated stage → inert), (`opus`/`sonnet` + gated stage → inert), (`unset` + gated stage → inert, fail-open) — asserting the returned `{tier, route_gate_verdicts}` against each case's known input. This proves the boot gate's conditional, the one piece that keeps capable sessions unaffected.

**AC-3 — A `level-3-judge` standing-teammate mod exists in the workflow's `_mods/`, declares `model: opus`, and is spawned by the existing `spawn-standing-all` path without new spawn machinery.**
Verified by: `spacedock dispatch spawn-standing-all --workflow-dir {wd} --team {team}` emitting an Agent spec for `level-3-judge` with `model: opus` in its returned JSON array (asserted against the array contents), confirming the mod is discovered and spawned by the same mechanism as `comm-officer`, not a new code path.

**AC-4 — A level-3-capable FO leaves the gate flow unchanged: it authors its own `Recommend` line and writes no `decided-by: level-3-judge` record.**
Verified by: a fixture/behavior test asserting that with `SPACEDOCK_FO_MODEL` unset or `opus`/`sonnet`, `«fo.tier»()` returns `route_gate_verdicts: false`, so the gate flow's verdict path is the unmodified present-gate flow — proving the mechanism is a no-op for capable sessions and adds no cost to the common case.

## Test plan

- **AC-2 / AC-4 (the boot gate's conditional) — Go fixture test, cheap, runs in CI.** Table-driven over the four `(model, has_gated_stage)` cases asserting `{tier, route_gate_verdicts}`. This is the deterministic core and the smallest exercise of the launcher-env→tier mechanism; it seeds the `spacedock fo tier` verb's first test when that verb is built. Cost: minutes.
- **AC-3 (standing-teammate discovery) — fixture test over `spawn-standing-all`.** A test workflow dir with a `level-3-judge` mod; assert the emitted Agent-spec array contains the teammate with `model: opus`. Reuses the existing comm-officer spawn test shape. Cost: minutes.
- **AC-1 (the live routing) — coordinated with haiku-loop-spike, NOT a separate drive.** The spike's live Haiku drive is the live test; this task's contribution is the durable-state oracle (grep the archived entity for the `decided-by: level-3-judge` line and the terminal `verdict:`). Graded on durable state, never transcript phrasing. Cost: the spike's cost, not additive.
- **Sequencing:** AC-2/AC-4/AC-3 fixture tests gate the mechanism before the live drive (validating-new-mechanisms discipline: the cheap conditional/discovery checks pay the small bill first; the live drive in the spike pays the large one). The launcher-export-of-`SPACEDOCK_FO_MODEL` is the one new contract surface and is covered by the AC-2 fixture, not assumed.

## Documentation changes (concrete doc diff for the gate review)

This changes user-visible FO behavior (a new boot line on a level-2-only FO; the gate-verdict route) and adds a contract surface. The diffs below are recorded for the ideation gate; implementation applies them. Verbs that stay prose-functions this sprint flip notation later, not now.

**`skills/first-officer/references/first-officer-shared-core.md` — add a Startup step 1.5 (after the contract-version gate, before discovery, since tier is launcher-data with no workflow dependency):**

Before — Startup step 1 ends, step 2 begins:
```
   In every class, do NOT proceed to discovery or `--boot`.
2. Discover the project root with `git rev-parse --show-toplevel`.
```
After — insert step 1.5:
```
   In every class, do NOT proceed to discovery or `--boot`.
1.5. **Tier self-identification.** Invoke `«fo.tier»()`: read `SPACEDOCK_FO_MODEL` (launcher-set; unset → level-3-capable), map model→tier (`haiku → level-2-only`; else `level-3-capable`). The workflow's `gate: true` stages are read at step 4; once read, if `tier == level-2-only` AND any stage is `gate: true`, arm gate-verdict routing to the `level-3-judge` standing teammate and emit one boot line: "Level-2-only FO: gate verdicts route to level-3-judge." A level-3-capable FO arms nothing and the gate flow is unchanged. (The arming check folds into the greet once step 4 has read the stage flags.)
2. Discover the project root with `git rev-parse --show-toplevel`.
```

**`skills/first-officer/references/first-officer-shared-core.md` — `## Completion and Gates`, the gated-stage block (lines 115-121), add the verdict-route bullet:**

Before:
```
If the stage is gated:
- never self-approve
- present the stage report by invoking `Skill(skill="spacedock:present-gate")` and following its template + assembly rules
```
After:
```
If the stage is gated:
- never self-approve
- if gate-verdict routing is armed (a level-2-only FO, per Startup 1.5), obtain the verdict via `«gate.route-verdict»(slug, stage)`: route the stage report + AC section to `level-3-judge`, fill the `present-gate` `Recommend` line from the L3-authored verdict, and record a `### Gate Verdicts` line (`{stage}: {verdict} — decided-by: level-3-judge`); if `level-3-judge` is absent/silent, surface the gate to the captain with "level-3 unavailable" and author NO verdict
- present the stage report by invoking `Skill(skill="spacedock:present-gate")` and following its template + assembly rules
```

**`skills/first-officer/references/first-officer-shared-core.md` — `## FO Write Scope`, the `### Feedback Cycles` bullet, add a sibling bullet for `### Gate Verdicts`:**

After (insert after the `### Feedback Cycles` bullet):
```
- **`### Gate Verdicts` section** — in entity bodies, recording the verdict author for each gated stage (`{stage}: {verdict} — decided-by: {level-3-judge | self}`). Same write-scope rule as `### Feedback Cycles`: worktree-side when `worktree:` is set, main-side otherwise. A level-2-only FO writes `decided-by: level-3-judge`; a capable FO needs no line (absence on a level-2-only session is a detectable violation).
```

**New file `docs/dev/_mods/level-3-judge.md`** — declared like `comm-officer.md` (`standing: true`; `## Hook: startup` with `subagent_type: general-purpose`, `name: level-3-judge`, `model: opus`; `## Hook: shutdown` mirroring comm-officer; a `## Routing Usage` block stating the load-bearing-blocking gate-verdict contract; an `## Agent Prompt` instructing the teammate to author gate verdicts — apply the present-gate assembly rules to the forwarded stage report + AC section and return `{approve|reject, one-line reason}`, and to make design/scope/escalation judgment calls when routed). Full mod text authored at implementation; the ideation gate reviews this shape.

## Stage Report: ideation

- DONE: Design «fo.tier»() tier self-identification (resolved runtime model name or a launcher flag — pick one source of truth at boot)
  Chose launcher-set `SPACEDOCK_FO_MODEL` env var as the single source (rejected self-probe as unreliable and README-declaration as wrong-locus); `«fo.tier»()` declared with the model→tier map and fail-open default. Boot insertion point is a new Startup step 1.5 (`first-officer-shared-core.md`).
- DONE: a named routing table mapping judgment categories (gate verdicts, design/scope, feedback-cycle-3 escalation, model-mismatch reuse, conflict recovery, teardown health) to a level-3 route
  Six-row table authored in Proposed approach; gate-verdicts row is the live-wired no-2y slice, the other five named-but-deferred-to-2y.
- DONE: a standing level-3 teammate (stronger model) spawned on demand
  `level-3-judge` mod designed on the `comm-officer` template (`standing: true`, `model: opus`), spawned by the existing `spawn-standing-all` path (no new machinery); AC-3 asserts discovery via that path.
- DONE: a boot gate that routes ALL gate verdicts to level-3 when the FO is Haiku and the workflow has gate:true stages
  Boot-gate conditional designed (Startup 1.5): arm only when `level-2-only` ∧ a `gate: true` stage exists; inert otherwise (AC-2/AC-4). The gate-half is named as the no-2y startable slice; full seam visibility noted as composing with 2y.
- DONE: Coordinate the live-drive proof with haiku-loop-spike + behavior-first ACs graded on durable workflow state
  Spike determination records "no separate spike needed: exercised in haiku-loop-spike"; AC-1 grades on the durable `### Gate Verdicts` `decided-by: level-3-judge` line + terminal `verdict:`, never transcript phrasing. Concrete doc diffs recorded for the ideation gate.

### Summary

Designed the model-aware tier mechanism as four pieces in the established prose-function notation: `«fo.tier»()` self-identifies the FO tier from a single launcher-set `SPACEDOCK_FO_MODEL` source at a new Startup step 1.5; a six-row routing table names every judgment category's level-3 route; a `level-3-judge` standing teammate (`model: opus`) reuses the existing `comm-officer`/`spawn-standing-all` machinery; and a boot gate arms gate-verdict routing only for a level-2-only FO meeting `gate: true` stages, inert otherwise. The no-2y startable slice is the gate-verdicts row wired live; the other five rows ship as the documented table and compose with 2y. The riskiest mechanism — tier self-ID + routing holding in a live Haiku drive — is exercised in haiku-loop-spike, not a duplicate drive; AC-1 grades on a durable `### Gate Verdicts` `decided-by: level-3-judge` artifact. Key decision: the level-3 gate route is load-bearing-blocking (FO surfaces "level-3 unavailable" rather than self-authoring a verdict), which is the entire safety property and the one place it differs from comm-officer's best-effort polish.
