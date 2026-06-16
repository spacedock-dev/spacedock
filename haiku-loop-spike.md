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
   the *output* of this spike, not its input. The drive reuses the proven `internal/ensigncycle` HELPERS:
   `isolatedClaudeEnv` (clean HOME + OAuth/API-key auth, t.Skip when neither), the `streamWatcher`
   quiet-budget liveness, and the durable-state + tool-stream grade. But the substrate proves only the
   `spacedock claude` FRONT-DOOR launch — bare `claude --model haiku -p` (no front door, contract not
   loaded) is a NEW launch path it does NOT prove, smoke-tested first (AC-1). `withBinaryOnPath` still
   applies: the prose-functions call `spacedock`, so the built binary stays on the child PATH even though
   the launched executable is bare `claude`.

4. **The spike's own riskiest sub-mechanisms are the bare launch and the loop-substitution, not FO
   behavior.** Two things could invalidate the rest of the work before any FO judgment matters: (a) the
   bare-`claude` launch path booting/authing/streaming at all (NEW, unproven by the substrate — AC-1),
   and (b) "can a Haiku model reliably FOLLOW a *hand-authored* substituted loop at all, instead of
   improvising or reverting to its own idea of what an FO does." Per validating-new-mechanisms
   (recursively): the drive exercises (a) with a sub-minute smoke FIRST, then (b) with the smallest
   single-turn slice (boot → one dispatch → consume worker report → path-scoped commit → stop), proving
   the substitution holds, BEFORE the full dispatch→gate→merge drive. If the smoke or the slice fails,
   the full drive is meaningless.

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

**What the spike measures is residency-AGNOSTIC: the routing BOUNDARY — which judgment categories Haiku
cannot decide and must hand to a stronger model — not how the level-3 model is housed.** The seed and
frontmatter say "standing level-3 teammate"; that names a *residency*, and the spike does not need to
validate residency to find the boundary. So the residency choice below is a spike-local convenience, not
a claim about 0205's production shape.

Two shapes were considered. **Recommendation: per-judgment bare-mode `Agent(model="opus", …)` blocking
call, NOT a standing team member.** Rationale: a single-entity `-p` drive is bare-mode (the contract's
own rule: "single-entity mode → skip team creation → bare-mode dispatch"), and a bare `Agent` call
without `team_name` blocks until the subagent returns — which is exactly the synchronous "ask the brain,
get the verdict" shape a judgment call needs, and it sidesteps the deferred-`TeamCreate`/`SendMessage`
`ToolSearch` hop entirely. The standing-teammate shape (declare a level-3 mod, spawn via
`spawn-standing-all`) is the *production* shape but adds the team-lifecycle surface the spike does not
need to validate. If the gate prefers the standing shape, the fixture declares a `_mods/level-3.md` with
`model: opus` and the drive spawns it — but that imports the deferred-tool hop as an extra Haiku failure
surface, which is itself worth noting as a finding either way.

**Cross-member proof boundary (gate-readability).** AC-3's verdict-PROVENANCE oracle is
mechanism-SPECIFIC to whatever residency THIS drive uses (it traces the verdict to a bare
`Agent(model=opus)` sub-call when the recommendation is taken). A sibling 0205 member —
`fo-tier-delegation` — builds a STANDING `level-3-judge` mod with `spawn-standing-all` + `SendMessage`
routing. That standing path does NOT inherit this spike's AC-3 artifact: a different residency is a
different verdict-provenance shape that the sibling must prove on its own. The gate should read this
spike's AC-3 as proof that *the boundary routes under the spike's residency*, not as a cross-member proof
that the standing production wiring works.

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

**AC-1 — The bare-`claude` launch path works in the isolated env, proven by a sub-minute smoke before
any loop drive is attempted.**
A bare `claude --model haiku -p {trivial-prompt}` — launched WITHOUT the `spacedock` front door, so no
`--agent`, no plugin, no skills, the FO contract NOT loaded — boots, authenticates against the isolated
HOME credential, and emits a parseable `stream-json` transcript with a `result`/success event. This is a
NEW launch path: every existing `internal/ensigncycle` live test launches `exec.Command(binary, "claude",
…)` through the `spacedock claude` front door (which hard-codes `--agent spacedock:first-officer`), never
a bare `claude` executable, and `withBinaryOnPath` puts `spacedock` — not `claude` — on PATH. So the bare
launch is unproven by that substrate and must be smoke-tested first.
Verified by: the smoke's captured `stream-json` parsing to a success `result` event (via the existing
`extractClaudeFinalMessage` shape), with auth resolved by `isolatedClaudeEnv` (t.Skip when no credential).
Runs and passes BEFORE AC-2.

**AC-2 — The substitution mechanism holds: a Haiku model FOLLOWS the hand-authored simplified loop
rather than improvising its own FO behavior, proven by the smallest single-turn slice before the full
drive.**
The smallest slice (boot → one `«dispatch»` → consume worker report → path-scoped commit → stop) leaves
a durable end-state consistent with the loop's prescribed steps, and the stream shows Haiku invoked the
loop's prose-functions in order rather than substituting an improvised flow.
Verified by: the single-turn-slice drive's durable state (one path-scoped commit naming only the entity,
a protocol-shaped stage report) + a stream check that the prescribed binary calls appear in order. This
slice runs and passes AFTER AC-1's bare-launch smoke and BEFORE the full drive
(validating-new-mechanisms ordering).

**AC-3 — The fixture entity reaches its terminal stage through a Haiku-driven hand-loop, with the gate
verdict on the record made by the level-3 model, not by Haiku.**
The durable end-state of the throwaway fixture workflow shows the driven entity archived/terminal with a
recorded verdict, the appended stage report(s), and a gate-verdict provenance that traces to the
level-3 `Agent(model=opus,…)` judgment call — not to a Haiku self-decision.
Verified by: a durable-state oracle over the fixture after the drive (entity terminal + archived +
verdict set), CROSS-CHECKED against the captured tool-call stream showing the verdict originated from
the level-3 sub-call. Grading reads on-disk state + stream, never the FO's transcript prose. (Mirrors
the `internal/ensigncycle` durable-state-plus-stream grade.) The provenance oracle is mechanism-SPECIFIC
to this drive's residency (see the Level-3 wiring decision) — it does not stand in for the standing-mod
shape a sibling member builds.

**AC-4 — Each of the 5 predicted failure modes carries a recorded verdict of SURFACED / REFUTED /
NOT-EXERCISABLE with concrete evidence; none is left as an unbacked pass.**
For every predicted mode the entity body records the verdict, the evidence (stream excerpt or end-state
fact for SURFACED/REFUTED; the concrete unreachability reason + the harness that would test it for
NOT-EXERCISABLE). A mode with no evidence either way is a FAILED measurement. Failure-mode 1 (rebase-conflict
HALT) is NOT-EXERCISABLE without the injected-conflict variant; if that variant is out of scope for this
drive, the entity records on the record that mode 1's HALT behavior is validated only DOWNSTREAM — by
`state-verbs`' verb-enforced halt and the final `spacedock claude`-shape haiku-drive-validation — so the
gate sees the residual rather than a silent gap.
Verified by: the per-mode table in the drive's recorded findings, each row tied to a stream/end-state
artifact path from the run.

**AC-5 — The must-build verb list and the irreducible-judgment routing boundary are recorded as
PROVISIONAL input to the rest of 0205's carve, under a pinned decision rule.**
The entity body holds (a) the `{loop-step → breaks-every-time|holds-reliably → becomes-verb}` table and
(b) the observed `{judgment-category → routed-to-level-3}` boundary list, each derived from the drive
evidence under this rule: each drive is repeated N≥3 times; a loop step that breaks on ANY of the N≥3
drives is classified must-build (a stochastic model that breaks sometimes breaks unsafely), and a step is
"holds reliably" ONLY if clean across all N. The recorded list is PROVISIONAL — it sequences the carve of
the held sibling members (`gate-extract-verbs`, `fo-tier-delegation`), then is reconfirmed by the final
`haiku-drive-validation` on the real `spacedock claude` shape (which loads the actual restructured
contract, not this hand-loop). It is not the final typed input; it is the input that lets the siblings
start against evidence.
Verified by: the two recorded tables, each row citing the N≥3 drive observations that produced it and
labelled provisional-pending-final-drive; downstream members consume them as a sequencing input.

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

**Ordering (validating-new-mechanisms).** (1) Bare-launch smoke (AC-1) — sub-minute; proves the NEW
bare-`claude --model haiku -p` launch path boots/auths/streams at all, before spending any drive on it.
(2) Single-turn-slice drive (AC-2) — smallest exercise of the loop-substitution mechanism. (3) Full
dispatch→gate→merge→terminal drive (AC-3), repeated N≥3 for the AC-5 decision rule. (4) The
injected-conflict variant for failure-mode 1 if in scope (AC-4).

**Cost / complexity.** A live drive lane: the smoke is sub-minute; the slice sub-minute; the full drive a
few minutes, run N≥3 (so a few-times-3 budget) for the breaks-every-time decision rule. Cheap relative to
building the verbs. Requires a live credential (benchmark-token or `ANTHROPIC_API_KEY`); skips cleanly
without one. No new on-disk format and no new binary surface — it consumes today's shipped commands, so
the new mechanisms are the bare-launch path (covered by AC-1) and the hand-loop prompt (covered by AC-2).

**What counts as proof.** Durable fixture end-state (terminal/archived/verdict, path-scoped commits,
stage reports) + the captured tool-call stream. Never the FO's narration. A mode "passes" only with a
stream excerpt or an on-disk fact behind it.

**Spike-first determination.** Two unverified mechanisms are spiked, smallest-first: the bare-`claude`
launch path (AC-1 smoke) and the loop-substitution (AC-2 slice). The narrowed substrate-reuse claim:
"no spike needed for the auth/HOME-isolation, liveness, and durable-state-grade HELPERS — they reuse the
proven `isolatedClaudeEnv` + `streamWatcher` + durable-state-grade path the shared-scenario live tests
exercise. But that substrate proves only the `spacedock claude` FRONT-DOOR launch (every internal live
test runs `exec.Command(binary, "claude", …)` through the front door, which hard-codes
`--agent spacedock:first-officer`; `withBinaryOnPath` puts `spacedock`, not `claude`, on PATH). The bare
`claude` executable launch with the contract NOT loaded is unproven by it and is spiked by AC-1."

## Stage Report: ideation

- DONE: Use TEST FIXTURES, not this repo: design the spike to drive a Haiku FO against a temp split-root FIXTURE workflow (the internal/ensigncycle live-test pattern — isolated, repeatable, throwaway), never a real entity in this repo.
  Test plan + Out-of-scope pin the throwaway temp split-root fixture and reuse the proven internal/ensigncycle substrate (isolatedClaudeEnv / withBinaryOnPath / streamWatcher / durable-state grade); "no real docs/dev entity is ever driven."
- DONE: Design the hand-simulated simplified loop + the standing level-3 teammate wiring + the measurement (which mechanics Haiku breaks vs holds, where the irreducible-judgment boundary sits).
  Proposed approach declares the «boot»/«next»/«dispatch»/«gate»/«merge» prose-functions, a level-3 wiring DECISION (per-judgment bare Agent(model=opus) recommended over a standing member, with rationale and the gate's alternative), and the must-build/routing-boundary measurement.
- DONE: Behavior-first ACs over durable fixture state (entity terminal; level-3 gate verdict on record); AC for the 5 predicted failure modes surfaced or refuted with evidence; recorded must-build verb list + routing boundary as input to the rest of the 0205 carve.
  AC-1 (terminal + level-3-provenance verdict), AC-2 (loop-substitution holds, single-turn slice first), AC-3 (5 modes each SURFACED/REFUTED/NOT-EXERCISABLE with evidence), AC-4 (must-build + routing-boundary tables) — all verified by durable-state + tool-stream, never transcript prose.

### Summary

Designed the 0205 gating spike. The breadth-first map produced three load-bearing forensic findings that reshaped the naive plan: `spacedock claude` hard-codes `--agent spacedock:first-officer` (frontdoor.go:319) so a simplified loop must launch bare `claude --model haiku` with the hand-loop as the `-p` prompt; team tooling (Agent/TeamCreate/SendMessage) is native, so a Haiku FO can spawn a stronger-model level-3 teammate with no plugin; and the spike's own riskiest sub-mechanism is loop-SUBSTITUTION (can Haiku follow a hand-loop at all), which AC-2 exercises with a single-turn slice BEFORE the full drive. Two design decisions surfaced for the gate: per-judgment bare `Agent(model=opus)` over a standing level-3 member (the seed/frontmatter say "standing"; I recommend bare for the throwaway drive and present standing as the alternative), and a reachability triage showing only 3 of the 5 predicted failure modes are directly exercisable in a single-writer `-p` drive (modes 1 and 2 are NOT-EXERCISABLE without an injected conflict / a not-yet-built verb, recorded honestly rather than refuted by absence).

## Stage Report: ideation (cycle 2)

- DONE: Fold 1 — bare `claude --model haiku -p` is a NEW launch path not proven by internal/ensigncycle; add a sub-minute smoke before the single-turn slice and narrow the "launch substrate already proven" claim to spacedock-claude only.
  Verified the claim myself: `grep` shows every live test uses `exec.Command(binary, "claude", …)` via the front door (which hard-codes `--agent`), zero bare `exec.Command("claude", …)`; `withBinaryOnPath` puts `spacedock` on PATH, not `claude`. New AC-1 (bare-launch smoke, runs first); forensic finding #3/#4 caveated; Test-plan ordering + spike-first determination narrowed to front-door-only.
- DONE: Fold 2 — pin AC's must-build decision rule (N≥3, any-break→must-build, holds only if clean across all N) and reword the must-build list + routing boundary as PROVISIONAL, reconfirmed by the final haiku-drive-validation on the real spacedock claude shape.
  AC-5 (was AC-4) now pins N≥3 + the conservative any-break rule and labels both tables provisional-pending-final-drive; Test-plan cost reflects the few-times-3 budget.
- DONE: Fold 3 — reconcile residency wording: spike measures the routing BOUNDARY (residency-agnostic); the verdict-provenance oracle is mechanism-SPECIFIC; the standing `level-3-judge` sibling does NOT inherit this AC artifact; failure-mode-1 HALT residual (validated only downstream) is on the record at the gate.
  Level-3 wiring section now leads with the boundary-vs-residency distinction + a "Cross-member proof boundary" note (AC-3 provenance oracle does not transfer to fo-tier-delegation's standing mod); AC-4 records mode-1's downstream-only HALT validation.

### Summary

Folded the three staff-review corrections into the body, ACs, and test plan together. The load-bearing one was fold 1: I had overclaimed the bare-`claude` launch as "already proven" — I re-verified by grep that every internal live test goes through the `spacedock claude` front door (which forces `--agent`), so bare `claude` with the contract unloaded is genuinely new and now gets a sub-minute AC-1 smoke ahead of everything. ACs renumbered to AC-1..AC-5 (smoke → slice → full-drive → failure-modes → provisional-list); cross-references (the residency provenance note, the ordering chain, the spike-first determination) updated to match. Folds 2 and 3 sharpen the must-build rule (N≥3, conservative, provisional pending the final real-shape drive) and the residency framing (boundary is residency-agnostic; the provenance oracle and failure-mode-1 HALT proof do not transfer cross-member), which aligns with the captain's 2026-06-16 re-sequencing that holds the two sibling verbs pending this spike's findings.
