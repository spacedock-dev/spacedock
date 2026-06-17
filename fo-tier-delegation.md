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

**1. Tier self-identification — one source of truth at boot.** `«fo.tier»()` resolves the FO's tier from a single launcher-set source, not a self-probe (a model asked "what model are you" is unreliable) and not a README declaration (the tier is a property of the running session, not the workflow). The launcher (`spacedock claude --model haiku`) parses the `--model` flag, normalizes aliases, and exports `SPACEDOCK_FO_MODEL` into the FO's environment (requires launcher change: parse passthrough, normalize, export, add `--env-pass SPACEDOCK_FO_MODEL` safehouse). `«fo.tier»()` reads it and maps to a tier: `haiku → level-2-only`; `sonnet | opus → level-3-capable`; `<unset or unresolvable> → level-2-only` (fail-safe: an unset var defaults to weak tier, requiring explicit capability assertion). The map lives in one place so "which models are level-2-only" is data, not scattered prose.

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
- **effect:** read `SPACEDOCK_FO_MODEL` (launcher-set; unset or unresolvable → treat as unset); map model→tier via the tier map (`haiku → level-2-only`; `sonnet | opus → level-3-capable`; unset/unresolvable → `level-2-only`); if `level-2-only` ∧ workflow has a `gate: true` stage, arm gate-verdict routing to `level-3-judge` and emit the boot line
- **done-when:** the FO tier is resolved and (for a level-2-only FO meeting gated stages) gate-verdict routing is armed
- **block:** none — unresolvable model defaults to `level-2-only` (fail-safe: require explicit capability assertion rather than hobbling a capable FO that loses the env var by accident)
- → **prose**, becomes `` `spacedock fo tier` `` (a verb that reads the launcher env and the workflow frontmatter and returns `{tier, route_gate_verdicts: bool, gate_route: "level-3-judge"}`)

## «gate.route-verdict»(slug, stage): obtain the gate verdict from the armed route
- **guard:** gate-verdict routing is armed (a level-2-only FO at a `gate: true` stage)
- **effect:** route the stage report + AC section + checklist to `level-3-judge` via `SendMessage(to="level-3-judge", ...)`; receive the `{approve|reject, reason}`; write the `### Gate Verdicts` line `decided-by: level-3-judge`; fill `present-gate`'s `Recommend` line from the returned verdict
- **done-when:** the `Recommend` line is filled from an L3-authored verdict and the `### Gate Verdicts` line is committed
- **block:** teammate absent/silent past timeout → surface the gate to the captain with "level-3 unavailable" and NO FO verdict; never self-author the `Recommend` line on a level-2-only session
- **Implementation note (temporary hand-followed prose):** This prose-function currently asks the Haiku FO to hand-follow multiple steps: pull the right stage-report section, extract the AC section, format the message, send via SendMessage, parse the plain-text response, write the `### Gate Verdicts` line, and commit path-scoped. This creates 8 points of manual failure (no schema safety). Post-2y, `spacedock gate route-verdict` will be a shipped binary verb that owns this complexity as a single atomic operation. Until then, the hand-followed prose is accepted as a temporary migration step per the prose-function-as-substrate model (see Out-of-scope, line 72).
- → **prose** (stays prose this sprint; the verdict is judgment, never a binary; the routing operation becomes a verb post-2y)

When the FO is `level-3-capable`, `«gate.route-verdict»` is never invoked — the gate flow runs exactly as today (the FO authors its own `Recommend` line), so this mechanism is a no-op for capable sessions.

## Out of scope

- **Live-wiring the five non-gate routing-table rows** (design/scope, feedback-cycle-3, model-mismatch reuse, conflict recovery, teardown health). They ship as the documented table with named routes; their live wiring composes with 2y, which makes the mechanical/judgment seam visible in the restructured dispatch and merge cores. This slice wires only the gate-verdicts row.
- **The `spacedock fo tier` binary verb.** This slice ships the prose-function `«fo.tier»()` (hand-followed) with its named migration target; building the verb follows the spike, like the other 0205 verbs.
- **The `spacedock gate route-verdict` binary verb.** `«gate.route-verdict»()` stays a hand-followed prose-function this sprint; the routing operation becomes a verb post-2y (the accepted-fragility decision is recorded in the Implementation note on `«gate.route-verdict»` and in the Rework Summary). The verdict itself is irreducible judgment and never becomes a binary.
- **Codex / Pi substrates.** No team registry on Codex/Pi and different teammate models; their layered-FO operability is a 0205-scoped follow-up, Claude only here.
- **Making the `level-3-judge` model configurable per workflow.** It declares `model: opus` in its mod; per-workflow model selection is a later refinement, not this slice.

## Spike determination

**This task owns its own live drive of the standing `level-3-judge` route — `haiku-loop-spike` does NOT provide it.** The earlier ideation borrowed the spike's AC-1 as this task's AC-1 proof; that was rejected at the gate and is wrong on the spike's own terms. Verified against `haiku-loop-spike.md`: the spike (a) measures the routing BOUNDARY (which judgment categories Haiku cannot decide), which it states explicitly is "residency-AGNOSTIC" (lines 88–92); (b) RECOMMENDS a per-judgment bare `Agent(model=opus)` blocking call, NOT a standing mod (lines 94–103); (c) lists "Validating the standing-teammate residency mechanism" as OUT OF SCOPE (lines 173–175); (d) grades verdict provenance through the captured tool-call STREAM, not an on-disk artifact, and writes no `### Gate Verdicts` line (zero mentions of `### Gate Verdicts`, `decided-by`, or `level-3-judge` anywhere in the spike); and (e) tells the gate, in its own "Cross-member proof boundary" note (lines 105–112), that the standing `level-3-judge` sibling "does NOT inherit this spike's AC-3 artifact." So the spike proves the BOUNDARY (what routes to L3); it does NOT prove this task's standing residency (the standing `level-3-judge` mod + `SendMessage` route authoring the `### Gate Verdicts` durable line). That proof is this task's own, sequenced as the `haiku-drive-validation` member (AC-1).

**The `SPACEDOCK_FO_MODEL` launcher export is NOT trivial and is built by this slice — not borrowed.** The earlier ideation called the export "trivial, already-proven … `SPACEDOCK_BIN` is exported the same way today." Verified false against the live tree: `frontdoor.go` carries `--model` as opaque passthrough (the value-taking-flag table, line 554; appended verbatim with the rest of `fd.passthrough`) and never resolves it to a tier; `launchEnv` (lines 56–63) exports ONLY `SPACEDOCK_BIN`; and `SPACEDOCK_BIN` survives the safehouse boundary only via an explicit `--env-pass SPACEDOCK_BIN` flag (`launcherBinEnvPassFlags`, lines 64–78; wired at line 343). So `«fo.tier»()` reading a launcher-set tier source requires NEW launcher work, owned by this slice: parse `--model` out of the passthrough, normalize aliases (`default`/`haiku`/`sonnet[1m]`/full-id/absent), resolve to and export `SPACEDOCK_FO_MODEL`, AND add `--env-pass SPACEDOCK_FO_MODEL` so the var survives the sandbox. The launcher already knows the model VALUE; it does not expose any resolved tier the FO can read. This is the Test plan's "Launcher surface (prerequisite)" line, covered by a test that the var actually reaches a launched FO. The fail-safe default (an unset/scrubbed var → `level-2-only`) is precisely the protection against this surface being incompletely wired: a missing `--env-pass`, a resume, or a hand-launched bare `claude --model haiku` all leave the var unset inside the sandbox, and the FO must route rather than self-approve in every one of those cases.

## Acceptance criteria

**AC-1 — A level-2-only (Haiku) FO at a `gate: true` stage obtains the gate verdict from the standing `level-3-judge` teammate (via `SendMessage`) and the durable entity state shows the verdict was authored by level-3, not by Haiku.**
Verified by: a live Haiku FO drive (the "haiku-drive-validation" member, sequenced after 2y lands and all other members are built) at a `gate: true` stage routing the verdict to the standing `level-3-judge` mod and writing the `### Gate Verdicts` section with a `decided-by: level-3-judge` line for the gated stage. The durable state is graded by grepping the archived entity for the line and the terminal `verdict:` frontmatter, never by reading the drive transcript. **Note on spike coordination:** The haiku-loop-spike measures the routing BOUNDARY (what judgment categories Haiku cannot decide, residency-agnostic) and may use a different residency mechanism (bare per-judgment `Agent(model=opus)` vs. the standing mod). The spike proves the BOUNDARY; this AC proves the STANDING residency that 0205 uses in production. Do not claim this AC is satisfied by the spike's drive — it has its own live validation.

**AC-2 — `«fo.tier»()` resolves the correct tier and arms gate-verdict routing only for a level-2-only FO meeting `gate: true` stages.**
Verified by: a fixture-driven test of the tier-resolution + arming decision over five cases — (`SPACEDOCK_FO_MODEL=haiku` + workflow with a gated stage → armed, boot line emitted), (`haiku` + no gated stage → inert), (`opus`/`sonnet` + gated stage → inert), (`unset` + gated stage → **armed, fail-safe** — treats as level-2-only, boots "level-2-only FO: gate verdicts route to level-3-judge"), (a garbage/unrecognized model string + gated stage → **armed, fail-safe** — an unresolvable value is treated as weak, not silently capable) — asserting the returned `{tier, route_gate_verdicts}` against each case's known input. The unset and garbage rows are the safety-critical ones: an un-provable session must NOT default to capable (that would let Haiku self-approve a verdict); it defaults to level-2-only and routes. This proves the boot gate's conditional and the fail-safe default. The fixture proves the resolution LOGIC; the live behavioral proof that a Haiku-tier session with the var actually unset ROUTES rather than self-approves rides AC-1's drive (the fixture cannot prove the FO's runtime behavior, only the tier function's return).

**AC-3 — A `level-3-judge` standing-teammate mod exists in the workflow's `_mods/`, declares `model: opus`, and is spawned by the existing `spawn-standing-all` path without new spawn machinery.**
Verified by: `spacedock dispatch spawn-standing-all --workflow-dir {wd} --team {team}` emitting an Agent spec for `level-3-judge` with `model: opus` in its returned JSON array (asserted against the array contents), confirming the mod is discovered and spawned by the same mechanism as `comm-officer`, not a new code path.

**AC-4 — A level-3-capable FO leaves the gate flow unchanged: it authors its own `Recommend` line and writes no `decided-by: level-3-judge` record.**
Verified by: a fixture/behavior test asserting that with `SPACEDOCK_FO_MODEL` unset or `opus`/`sonnet`, `«fo.tier»()` returns `route_gate_verdicts: false`, so the gate flow's verdict path is the unmodified present-gate flow — proving the mechanism is a no-op for capable sessions and adds no cost to the common case.

## Test plan

- **Launcher surface (prerequisite).** The `«fo.tier»()` mechanism depends on `SPACEDOCK_FO_MODEL` being exported into the FO's environment by the launcher. This requires launcher changes (parse `--model` from passthrough, normalize aliases, export the var, add `--env-pass SPACEDOCK_FO_MODEL` to safehouse). Must be built before AC-2/AC-4 fixtures can run. Ownership: ideation-rework; gate ratification: captain-ratify-at-gate (new contract surface).
- **AC-2 / AC-4 (the boot gate's conditional) — Go fixture test, cheap, runs in CI.** Table-driven over the four `(model, has_gated_stage)` cases asserting `{tier, route_gate_verdicts}`. The unset case is the critical safety test (must default to level-2-only, not level-3-capable). This is the deterministic core and the smallest exercise of the launcher-env→tier mechanism; it seeds the `spacedock fo tier` verb's first test when that verb is built. Cost: minutes.
- **AC-3 (standing-teammate discovery) — fixture test over `spawn-standing-all`.** A test workflow dir with a `level-3-judge` mod; assert the emitted Agent-spec array contains the teammate with `model: opus`. Reuses the existing comm-officer spawn test shape. Cost: minutes.
- **AC-1 (the live standing-mod routing) — the "haiku-drive-validation" member, sequenced LAST.** Not proven by the spike; has its own live Haiku FO drive exercising the standing `level-3-judge` mod and the `### Gate Verdicts` durable artifact. Graded on durable state (grep for the `decided-by: level-3-judge` line and terminal `verdict:`), never transcript. Cost: the final validation drive's cost.
- **Sequencing:** Launcher surface (prerequisite) → AC-2/AC-4/AC-3 fixture tests (validate the mechanism) → AC-1 live drive (final validation in haiku-drive-validation, sequenced last after 2y lands).

## Documentation changes (concrete doc diff for the gate review)

This changes user-visible FO behavior (a new boot line on a level-2-only FO; the gate-verdict route) and adds a contract surface. The diffs below are recorded for the ideation gate; implementation applies them. Verbs that stay prose-functions this sprint flip notation later, not now.

**`skills/first-officer/references/first-officer-shared-core.md` — Startup: add tier self-identification invocation**

After the contract-version gate and before discovery, invoke `«fo.tier»()` to read `SPACEDOCK_FO_MODEL` (launcher-set), resolve the FO tier (`haiku → level-2-only`; `sonnet | opus → level-3-capable`; unset/unresolvable → `level-2-only`, fail-safe), and defer the arming decision to after step 4 reads the workflow's gate stages. Insert:
```
   («fo.tier»(): read launcher SPACEDOCK_FO_MODEL, resolve tier per the tier map; if level-2-only ∧ a gate:true stage exists (known after step 4), arm gate-verdict routing to level-3-judge and emit the boot line)
```

**`skills/first-officer/references/first-officer-shared-core.md` — `## Completion and Gates`, the gated-stage block, add the verdict-route bullet:**

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
  Chose launcher-set `SPACEDOCK_FO_MODEL` env var as the single source (rejected self-probe as unreliable and README-declaration as wrong-locus); `«fo.tier»()` declared with the model→tier map and fail-SAFE default (unset/unresolvable → level-2-only). Boot insertion point is a new Startup step 1.5 (`first-officer-shared-core.md`). The launcher export of the var is owned by this slice (parse `--model`, normalize, export, `--env-pass`), not assumed.
- DONE: a named routing table mapping judgment categories (gate verdicts, design/scope, feedback-cycle-3 escalation, model-mismatch reuse, conflict recovery, teardown health) to a level-3 route
  Six-row table authored in Proposed approach; gate-verdicts row is the live-wired no-2y slice, the other five named-but-deferred-to-2y.
- DONE: a standing level-3 teammate (stronger model) spawned on demand
  `level-3-judge` mod designed on the `comm-officer` template (`standing: true`, `model: opus`), spawned by the existing `spawn-standing-all` path (no new machinery); AC-3 asserts discovery via that path.
- DONE: a boot gate that routes ALL gate verdicts to level-3 when the FO is Haiku and the workflow has gate:true stages
  Boot-gate conditional designed (Startup 1.5): arm only when `level-2-only` ∧ a `gate: true` stage exists; inert otherwise (AC-2/AC-4). The gate-half is named as the no-2y startable slice; full seam visibility noted as composing with 2y.
- DONE: Coordinate the live-drive proof with haiku-loop-spike + behavior-first ACs graded on durable workflow state
  Spike determination records that this task OWNS its own standing-residency drive (the spike proves only the residency-agnostic boundary, disclaims the standing mod, and writes no `### Gate Verdicts` artifact); AC-1 is the `haiku-drive-validation` member graded on the durable `### Gate Verdicts` `decided-by: level-3-judge` line + terminal `verdict:`, never transcript phrasing. Concrete doc diffs recorded for the ideation gate.

### Summary

Designed the model-aware tier mechanism as four pieces in the established prose-function notation: `«fo.tier»()` self-identifies the FO tier from a single launcher-set `SPACEDOCK_FO_MODEL` source at a new Startup step 1.5; a six-row routing table names every judgment category's level-3 route; a `level-3-judge` standing teammate (`model: opus`) reuses the existing `comm-officer`/`spawn-standing-all` machinery; and a boot gate arms gate-verdict routing only for a level-2-only FO meeting `gate: true` stages, inert otherwise. The no-2y startable slice is the gate-verdicts row wired live; the other five rows ship as the documented table and compose with 2y. The tier default is fail-SAFE: an unset/unresolvable var resolves to level-2-only (route to L3), never self-approve; the `SPACEDOCK_FO_MODEL` launcher export is owned by this slice, not assumed. The riskiest mechanism — the standing `level-3-judge` route authoring the verdict on-disk in a live Haiku drive — is this task's OWN proof (AC-1, the `haiku-drive-validation` member), NOT borrowed from haiku-loop-spike, which proves only the residency-agnostic routing boundary and disclaims the standing mod. Key decision: the level-3 gate route is load-bearing-blocking (FO surfaces "level-3 unavailable" rather than self-authoring a verdict), which is the entire safety property and the one place it differs from comm-officer's best-effort polish.

### Rework Summary (cycle 1 → cycle 2)

Reworked against the Opus level-3-judge REJECT (three material gaps), corroborated by the 0205 staff review (findings #1/#4/#5). Implementation details only; the core design — tier self-ID at boot + standing `level-3-judge` route authoring gate verdicts — is unchanged.

**1. AC-1 proof gap → this task OWNS its own live drive (does NOT borrow the spike's).** Re-read `haiku-loop-spike` in full: it measures the residency-AGNOSTIC routing boundary, RECOMMENDS a per-judgment bare `Agent(model=opus)` (not a standing mod), lists "Validating the standing-teammate residency mechanism" as out of scope, writes no `### Gate Verdicts` artifact, and tells the gate in its own cross-member note that the standing `level-3-judge` sibling "does NOT inherit this spike's AC-3 artifact." Borrowing its drive was wrong on the spike's own terms. AC-1 now states the spike proves the BOUNDARY and this task proves the STANDING residency via its own drive (the `haiku-drive-validation` member, sequenced last). The Spike-determination section and the Out-of-scope "separate live Haiku drive" bullet were rewritten to match; the Summary and Stage Report bullet 5 were de-staled.

**2. Tier default fails open → flipped to fail-SAFE.** `unset`/unresolvable `SPACEDOCK_FO_MODEL` now resolves to `level-2-only` (route to L3), never `level-3-capable`. The asymmetry is recorded: a needlessly-routing capable FO costs recoverable latency (the captain is the ultimate L3); a weak FO self-approving under a dropped var is an unrecoverable safety breach — the exact failure the sprint exists to prevent. Capable is the explicit opt-out (a resolved `sonnet`/`opus` in the var). Updated: the tier map (piece 1), the `«fo.tier»()` effect/block bullets, the doc-diff Startup invocation, AC-2 (now five cases incl. a garbage-model fail-safe row, replacing the old tautological "unset → inert"), AC-4, the test plan, and the Stage Report.

**3. Prose-function fragility → ACCEPTED as a temporary hand-followed step, with the binary-verb path named.** The decision: `«gate.route-verdict»()` stays prose this sprint. The verdict CONTENT is irreducible judgment and never becomes a binary; the surrounding DATA-MOVING (pull report section, extract AC, format, send, parse, write the durable line, commit) is the fragile part and becomes the `spacedock gate route-verdict` binary verb post-2y. The deterministic extraction half (which report section, the AC section, the checklist) is exactly what the sibling `gate-extract-verbs` (`6re`) verbs produce, so the prose body shrinks to "call the extract verbs → send their structured output → parse the verdict → write the line" once those land — the fragility is bounded, not permanent. Recorded as the Implementation note on `«gate.route-verdict»` and a new Out-of-scope bullet. NOT sequenced as a hard prerequisite (the slice is gradable without it via AC-1's drive), but the path-to-binary is explicit.

**4. Launcher `SPACEDOCK_FO_MODEL` export → owned by this slice, not assumed.** The earlier "trivial, already-proven" claim was verified false: `frontdoor.go` carries `--model` as opaque passthrough (line 554), `launchEnv` exports only `SPACEDOCK_BIN`, and the var survives the safehouse only via an explicit `--env-pass`. The slice now owns: parse `--model` from passthrough, normalize aliases, resolve+export `SPACEDOCK_FO_MODEL`, add `--env-pass SPACEDOCK_FO_MODEL`. Recorded as the Test plan's "Launcher surface (prerequisite)" line with a test that the var reaches a launched FO. This is also why the fail-safe default matters: a missing env-pass / resume / bare hand-launch all leave the var unset, and the FO must route, not self-approve.

**Captain-ratify-at-gate (noted, not reworked here):** the new `SPACEDOCK_FO_MODEL` contract surface; sequencing/co-design of this slice's three shared-core edits with the post-2y `prose-function-restructure` (which re-touches the same Startup / Completion-and-Gates / FO-Write-Scope regions); and narrowing the routing table's honest claim to the gate-verdicts row (the other five rows are a post-2y member's concern).

### Feedback Cycles

#### Cycle 1: Opus level-3-judge ideation gate review (2026-06-16)

**Verdict:** REJECT

**Findings:**

Level-3 gate review (Opus model) on ideation design. Three material gaps block ideation gate:

**1. AC-1 proof gap (architecture mismatch with spike)**
- AC-1 claims the haiku-loop-spike live drive produces `### Gate Verdicts` `decided-by: level-3-judge` durable artifact
- Verified: w4 (haiku-loop-spike) as ideated uses per-judgment bare `Agent(model=opus)`, NOT a standing `level-3-judge` mod
- w4 explicitly lists "Validating the standing-teammate residency mechanism" as OUT OF SCOPE
- w4 produces no `### Gate Verdicts` artifact; grading is via captured tool-call stream, not durable state
- **Action:** Either reconcile w4 to standing-mod shape and the durable artifact, or give 72 its own live drive

**2. Tier default fails open (inverted safety)**
- Current: `unset SPACEDOCK_FO_MODEL` → `level-3-capable` → Haiku self-approves verdicts
- This violates the 0205 thesis: "weak FO self-identifies and escalates structurally"
- **Action:** Flip to fail-safe — `unset` → `level-2-only` / require explicit capability assertion

**3. Prose-function approach is fragile (integration friction)**
- Design asks the Haiku FO to hand-follow `«gate.route-verdict»(slug, stage)` prose-function
- Actual work in the body:
  - Pull stage report section (which one? all 11 interleaved sections in real entities?)
  - Pull AC section
  - Format a routing message
  - SendMessage with no timeout/retry mechanism
  - Parse plain-text response (no schema, fragile to format changes)
  - Write `### Gate Verdicts` line
  - Path-scoped commit
- All 8 steps in a prose-function body = 8 points of failure
- Compare: a `spacedock gate route-verdict` binary verb owns this complexity
- **Action:** Acknowledge the fragility; either accept it for this slice (post-2y verbs will fix it), or sequence the binary verb as a prerequisite

**Spike integration note:** The integration spike (spawning Opus level-3-judge, routing 72's ideation gate) confirmed that the Opus model correctly identifies the material gaps. The message-passing worked; the verdict format was clear. But the manual prose-function approach requires the FO to handle data gathering, formatting, parsing, and durable-record writing without a safety schema.

**Rework targets (for next ideation cycle):**
1. Fix AC-1: reconcile with w4 architecture or own separate proof
2. Fix tier default: unset → level-2-only
3. Document the prose-function fragility and accept/mitigate it (or defer full binary verb to post-2y)
4. Clarify the launcher `SPACEDOCK_FO_MODEL` export work (staff review finding #2)

