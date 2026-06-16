---
title: Spike — live Haiku FO drives a hand-simulated simplified loop, before any verb is built
status: ideation
source: 'captain scope-lock (2026-06-15, this session) — 0205 is spike-first. Run a live Haiku FO drive on a hand-simulated simplified loop BEFORE building any binary verb, to prove the loop is Haiku-operable and map the irreducible-judgment boundary. Gates the rest of the 0205 carve (the must-build verb list comes from what Haiku breaks every time). Predicted Haiku failure modes to surface or refute — inventing rebase-conflict auto-recovery; interpreting a drift-class name semantically; bare-dispatching; auto-approving a gate; idle-vs-completion confusion.'
started: 2026-06-16T02:21:21Z
completed:
verdict:
score: 0.6
worktree:
issue:
sprint: 0205-layered-fo
sprint-readiness: ready
id: w4ryf4mg4vn1emwp906vd8yp
---

The 0205 gating exercise. A live Haiku FO drives one real entity through a hand-simulated simplified loop — the FO follows the dispatch->gate->merge prose-functions by hand (guillemets), with a standing level-3 teammate (stronger model) handling every judgment call. Run BEFORE any binary verb is built. The result sets the must-build verb list and the delegation routing table's real boundary, so the rest of 0205 builds against evidence, not a guess.

## Problem

The 0205 thesis — that pushing mechanics into binary verbs makes the FO loop simple enough for Haiku — rests on an unproven assumption: that Haiku can operate even the simplified loop at all, and that the judgment boundary sits where the per-area analysis predicts. Building the verbs before testing that assumption risks building the wrong set.

## Proposed approach

This IS the spike. Ideation designs the drive protocol, the fixture, the level-3 wiring, and the
measurement; the drive itself is the implementation/validation work that follows.

### Forensic framing (decided in ideation)

The breadth-first map (Phase -1, running-research-spikes) settled four facts that reshape the naive
"just run a Haiku FO" plan. All four are empirically grounded in the existing code, not assumed:

1. **`spacedock claude` cannot run a simplified loop.** The front door hard-codes
   `--agent spacedock:first-officer` (`internal/cli/frontdoor.go:319`) — it always loads the *full
   real* contract, with no override flag. Driving Haiku through `spacedock claude` would test "is the
   CURRENT full contract Haiku-operable" (the thesis already assumes NO). It does NOT test the spike's
   question, "is the SIMPLIFIED loop Haiku-operable." Therefore the drive launches `claude` directly
   (bare, no `--agent`) with the hand-authored simplified loop as the `-p` prompt.

2. **The team tooling is native, so the level-3 teammate needs no plugin.** `Agent`, `TeamCreate`,
   `SendMessage` are native Claude Code tools (`skills/using-claude-team/SKILL.md:13`; Pi bans them as
   native in `internal/cli/pi_frontdoor_test.go:90`). A bare `claude --model haiku` session can spawn a
   stronger-model teammate via `Agent(model="opus", …)` — model-per-subagent is a native Agent
   parameter. `TeamCreate`/`SendMessage` are *deferred* (need a `ToolSearch(query="select:…")` hop
   first); that hop is itself a candidate Haiku failure mode (see the level-3 wiring decision).

3. **The prose-functions wrap real binary calls, so the built binary must be on PATH.** The simplified
   loop's bodies call today's shipped surface (`spacedock status --boot --json`,
   `spacedock dispatch build`, `spacedock status --set/--archive`) — the not-yet-built 0205 verbs are
   the *output* of this spike, not its input. The drive reuses the proven `internal/ensigncycle`
   launch substrate: `isolatedClaudeEnv` (clean HOME + OAuth/API-key auth, t.Skip when neither),
   `withBinaryOnPath`, the `streamWatcher` quiet-budget liveness, and the durable-state + tool-stream
   grade.

4. **The spike's own riskiest sub-mechanism is loop-substitution, not FO behavior.** Because the
   simplified-loop prompt is *hand-authored*, the first thing that could invalidate the rest of the
   work is "can a Haiku model reliably FOLLOW a substituted hand-loop at all, instead of improvising or
   reverting to its own idea of what an FO does." Per validating-new-mechanisms (recursively): the
   drive's first exercise is the smallest single-turn slice (boot → one dispatch → consume worker
   report → path-scoped commit → stop), proving the substitution holds, BEFORE the full
   dispatch→gate→merge drive. If the single-turn slice fails, the full drive is meaningless.

### The simplified loop (the hand-followed prose-functions)

The drive prompt declares the loop as prose-functions in guillemets — the FO follows the body by hand
(prototype form from `docs/dev/_proposals/state-repo-pointer-prototype.md`). The loop is the minimum
dispatch→gate→merge spine:

- `«boot»()` → `spacedock status --boot --json`; greet from it.
- `«next»()` → `spacedock status --next --workflow-dir {wd}`; yields the one dispatchable entity+stage.
- `«dispatch»(slug, stage)` → set frontmatter + path-scoped commit, then `spacedock dispatch build`,
  then bare-mode `Agent(prompt=…)` for the worker; wait for the worker's completion.
- `«gate»(slug, stage)` → **route to level-3**: send the stage report + the entity's `## Acceptance
  criteria` to the level-3 teammate; the verdict (APPROVED/REJECTED + chosen direction) comes back from
  level-3 and is recorded. Haiku adjudicates none alone.
- `«merge»(slug)` → the set→invoke→clear→terminalize→archive ceremony (the binary backstop —
  `status --set`/`--archive` already refuse terminal-with-mod-block — is the safety net).

The loop body names, for each step, the future verb it becomes (the `→ becomes spacedock …` migration
target), so the recorded must-build list maps one-to-one onto loop steps.

### Level-3 wiring decision (recommendation, open to the gate)

Two shapes were considered. **Recommendation: per-judgment bare-mode `Agent(model="opus", …)` blocking
call, NOT a standing team member.** Rationale: a single-entity `-p` drive is bare-mode (the contract's
own rule: "single-entity mode → skip team creation → bare-mode dispatch"), and a bare `Agent` call
without `team_name` blocks until the subagent returns — which is exactly the synchronous "ask the brain,
get the verdict" shape a judgment call needs, and it sidesteps the deferred-`TeamCreate`/`SendMessage`
`ToolSearch` hop entirely. The standing-teammate shape (declare a level-3 mod, spawn via
`spawn-standing-all`) is the *production* shape but adds the team-lifecycle surface the spike does not
need to validate. The spike measures the BOUNDARY (what routes), not the residency mechanism. If the
gate prefers the standing shape, the fixture declares a `_mods/level-3.md` with `model: opus` and the
drive spawns it — but that imports the deferred-tool hop as an extra Haiku failure surface, which is
itself worth noting as a finding either way.

### Measurement (the spike's actual product)

Each drive captures the full tool-call stream + the durable workflow end-state (the
`internal/ensigncycle` grade pattern). From the stream + end-state, classify each loop step and each
predicted failure mode:

- **must-build verb list** — a loop step is a "must-build verb" when Haiku breaks it *every time*
  across the drives (e.g. always combines the mod-block clear with terminalize). A step Haiku holds
  reliably stays a prose-function. The list is recorded as `{step → breaks|holds → becomes-verb}`.
- **irreducible-judgment boundary** — what Haiku could NOT do even with the loop spelled out, that the
  level-3 teammate had to decide. This is the delegation routing table's real boundary, recorded as the
  category list (gate verdict, chosen direction, scope, …) actually observed to route.

### The 5 predicted failure modes — reachability triage (decided in ideation)

The predicted modes are HYPOTHESES TO TEST, not foregone. The single-entity `-p` drive can directly
exercise only some; the design states which, and how the rest are surfaced or refuted:

1. **Inventing rebase-conflict auto-recovery** — NOT directly reachable in a single-writer `-p` drive
   (no concurrent writer to cause a conflict). Surfaced via an *injected* conflict: the fixture's
   `«state.commit»` body is fed a pre-staged divergent commit so the path-scoped commit hits a rebase
   conflict, and the drive is graded on whether Haiku HALTS (per the rule) or invents a recovery. If the
   inject is out of scope for the first drive, this mode is recorded as "deferred to a 2-writer harness"
   with rationale — refuting by absence is not allowed.
2. **Interpreting a drift-class name semantically** — NOT reachable: `next-action`/`team_action` is a
   not-yet-built verb (it is the spike's *output*). The simplified loop uses today's `status --next`,
   which returns no drift-class. Recorded as "not exercisable pre-verb; the fully-qualified
   `team_action` schema is designed to PREVENT it, and that schema is what this spike informs."
3. **Bare-dispatching (skipping the checklist / `dispatch build`)** — directly reachable and graded
   from the stream: did the worker `Agent` prompt come from `spacedock dispatch build` output, or did
   Haiku hand-assemble it?
4. **Auto-approving a gate / inferring a chosen direction it should ask for** — the single highest-value
   observation, directly reachable: did Haiku route the verdict to level-3 (stream shows the
   `Agent(model=opus)` judgment call before any terminalize), or did it self-approve?
5. **Idle-vs-completion confusion (re-dispatching under social pressure)** — reachable in bare mode: a
   completed `Agent` returns control; does Haiku treat the return as the worker's completion-signal, or
   re-dispatch the same stage?

Each mode is recorded as SURFACED (Haiku did it, with the stream evidence), REFUTED (Haiku held, with
evidence), or NOT-EXERCISABLE (with the concrete reason and the harness that would test it). A mode with
no evidence either way is a FAILED measurement, not a pass.

### Doc-diff impact

None. The spike ships no user-visible CLI/output/banner change — it is a throwaway fixture drive plus a
recorded findings table in this entity body. The recorded must-build list and routing boundary feed the
*rest of 0205's carve* (the verb members and the delegation table), which is where any doc change lands.
Recorded here so the ideation-gate determination is on the record: "no doc diff — spike output is a
findings table consumed by sibling 0205 members, not a shipped behavior change."

## Out of scope

- **Building any binary verb.** The verbs follow the spike; this entity only RECORDS the must-build
  list as its input. No `spacedock dispatch next-action` / `state` / `merge guard` / `gate` code lands
  here.
- **The codex / pi substrates.** Claude only. The native-team-tool finding above is Claude-specific
  (Pi bans those tools); a Haiku drive on codex/pi is a separate follow-up with different failure modes.
- **The 2y-dependent core members.** The simplified loop wraps TODAY's shipped binary surface; it does
  not depend on the `shared-merge-dispatch-contract` extraction.
- **Validating the standing-teammate residency mechanism.** The spike measures the judgment BOUNDARY
  (what routes to level-3), not how the level-3 member is housed. The recommended per-judgment bare
  `Agent(model=opus)` shape is chosen precisely to keep team-lifecycle out of scope.
- **A real entity in this repo.** The drive runs only against a throwaway temp split-root FIXTURE
  workflow (the `internal/ensigncycle` isolated-repeatable-throwaway pattern). No `docs/dev` entity is
  ever driven.

## Acceptance criteria

Entity-level end-state properties of the finished spike, each with how it is checked.

**AC-1 — The fixture entity reaches its terminal stage through a Haiku-driven hand-loop, with the gate
verdict on the record made by the level-3 teammate, not by Haiku.**
The durable end-state of the throwaway fixture workflow shows the driven entity archived/terminal with a
recorded verdict, the appended stage report(s), and a gate-verdict provenance that traces to the
level-3 `Agent(model=opus,…)` judgment call — not to a Haiku self-decision.
Verified by: a durable-state oracle over the fixture after the drive (entity terminal + archived +
verdict set), CROSS-CHECKED against the captured tool-call stream showing the verdict originated from
the level-3 sub-call. Grading reads on-disk state + stream, never the FO's transcript prose. (Mirrors
the `internal/ensigncycle` durable-state-plus-stream grade.)

**AC-2 — The substitution mechanism holds: a Haiku model FOLLOWS the hand-authored simplified loop
rather than improvising its own FO behavior, proven by the smallest single-turn slice before the full
drive.**
The smallest slice (boot → one `«dispatch»` → consume worker report → path-scoped commit → stop) leaves
a durable end-state consistent with the loop's prescribed steps, and the stream shows Haiku invoked the
loop's prose-functions in order rather than substituting an improvised flow.
Verified by: the single-turn-slice drive's durable state (one path-scoped commit naming only the entity,
a protocol-shaped stage report) + a stream check that the prescribed binary calls appear in order. This
slice runs and passes BEFORE the full drive is attempted (validating-new-mechanisms ordering).

**AC-3 — Each of the 5 predicted failure modes carries a recorded verdict of SURFACED / REFUTED /
NOT-EXERCISABLE with concrete evidence; none is left as an unbacked pass.**
For every predicted mode the entity body records the verdict, the evidence (stream excerpt or end-state
fact for SURFACED/REFUTED; the concrete unreachability reason + the harness that would test it for
NOT-EXERCISABLE). A mode with no evidence either way is a FAILED measurement.
Verified by: the per-mode table in the drive's recorded findings, each row tied to a stream/end-state
artifact path from the run.

**AC-4 — The must-build verb list and the irreducible-judgment routing boundary are recorded as the
typed input to the rest of 0205's carve.**
The entity body holds (a) the `{loop-step → breaks-every-time|holds-reliably → becomes-verb}` table and
(b) the observed `{judgment-category → routed-to-level-3}` boundary list, each derived from the drive
evidence, not from the per-area prediction.
Verified by: the two recorded tables, each row citing the drive observation that produced it; downstream
0205 members (`gate-extract-verbs`, `fo-tier-delegation`) can consume them as a list.

## Test plan

The spike IS the test; there is no separate unit suite to write. The test is a live drive with a
durable-state-plus-stream oracle, reusing the proven `internal/ensigncycle` substrate.

**Shape.** A `//go:build live` test (or a self-contained driver script) that, per the
`internal/ensigncycle` pattern: stands up a temp split-root fixture workflow with one entity, resolves
auth/HOME via `isolatedClaudeEnv` (t.Skip when no credential), puts the built `spacedock` binary on the
child PATH via `withBinaryOnPath`, launches a **bare** `claude --model haiku -p {hand-loop-prompt}`
(NOT `spacedock claude` — no `--agent`), drains the stream via `streamWatcher` (60s quiet budget), and
grades durable end-state + the captured stream. The level-3 teammate is a per-judgment bare
`Agent(model=opus|sonnet,…)` call the hand-loop makes (recommended) or a declared `_mods/level-3.md`
spawned in-team (gate's choice).

**Ordering (validating-new-mechanisms).** (1) Single-turn-slice drive (AC-2) — smallest exercise of the
loop-substitution mechanism, the riskiest path; must pass first. (2) Full dispatch→gate→merge→terminal
drive (AC-1). (3) The injected-conflict variant for failure-mode 1 if in scope (AC-3).

**Cost / complexity.** A live drive lane: a few multi-minute Haiku sessions (the slice is sub-minute;
the full drive a few minutes). Cheap relative to building the verbs. Requires a live credential
(benchmark-token or `ANTHROPIC_API_KEY`); skips cleanly without one. No new on-disk format and no new
binary surface — it consumes today's shipped commands, so the only new mechanism is the hand-loop
prompt itself (covered by AC-2).

**What counts as proof.** Durable fixture end-state (terminal/archived/verdict, path-scoped commits,
stage reports) + the captured tool-call stream. Never the FO's narration. A mode "passes" only with a
stream excerpt or an on-disk fact behind it.

**Spike-first determination.** The riskiest unverified mechanism is the loop-substitution itself (can
Haiku follow a hand-loop), exercised first by the single-turn slice per AC-2 — this is the spike-within
the ideation discipline asks for. The launch/auth/liveness/grade substrate is already proven by the
existing `internal/ensigncycle` live tests, so it is NOT re-spiked: "no spike needed for the launch
substrate — it relies on the proven `isolatedClaudeEnv` + `withBinaryOnPath` + `streamWatcher` +
durable-state-grade path that the shared-scenario live tests already exercise."

## Stage Report: ideation

- DONE: Use TEST FIXTURES, not this repo: design the spike to drive a Haiku FO against a temp split-root FIXTURE workflow (the internal/ensigncycle live-test pattern — isolated, repeatable, throwaway), never a real entity in this repo.
  Test plan + Out-of-scope pin the throwaway temp split-root fixture and reuse the proven internal/ensigncycle substrate (isolatedClaudeEnv / withBinaryOnPath / streamWatcher / durable-state grade); "no real docs/dev entity is ever driven."
- DONE: Design the hand-simulated simplified loop + the standing level-3 teammate wiring + the measurement (which mechanics Haiku breaks vs holds, where the irreducible-judgment boundary sits).
  Proposed approach declares the «boot»/«next»/«dispatch»/«gate»/«merge» prose-functions, a level-3 wiring DECISION (per-judgment bare Agent(model=opus) recommended over a standing member, with rationale and the gate's alternative), and the must-build/routing-boundary measurement.
- DONE: Behavior-first ACs over durable fixture state (entity terminal; level-3 gate verdict on record); AC for the 5 predicted failure modes surfaced or refuted with evidence; recorded must-build verb list + routing boundary as input to the rest of the 0205 carve.
  AC-1 (terminal + level-3-provenance verdict), AC-2 (loop-substitution holds, single-turn slice first), AC-3 (5 modes each SURFACED/REFUTED/NOT-EXERCISABLE with evidence), AC-4 (must-build + routing-boundary tables) — all verified by durable-state + tool-stream, never transcript prose.

### Summary

Designed the 0205 gating spike. The breadth-first map produced three load-bearing forensic findings that reshaped the naive plan: `spacedock claude` hard-codes `--agent spacedock:first-officer` (frontdoor.go:319) so a simplified loop must launch bare `claude --model haiku` with the hand-loop as the `-p` prompt; team tooling (Agent/TeamCreate/SendMessage) is native, so a Haiku FO can spawn a stronger-model level-3 teammate with no plugin; and the spike's own riskiest sub-mechanism is loop-SUBSTITUTION (can Haiku follow a hand-loop at all), which AC-2 exercises with a single-turn slice BEFORE the full drive. Two design decisions surfaced for the gate: per-judgment bare `Agent(model=opus)` over a standing level-3 member (the seed/frontmatter say "standing"; I recommend bare for the throwaway drive and present standing as the alternative), and a reachability triage showing only 3 of the 5 predicted failure modes are directly exercisable in a single-writer `-p` drive (modes 1 and 2 are NOT-EXERCISABLE without an injected conflict / a not-yet-built verb, recorded honestly rather than refuted by absence).
