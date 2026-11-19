# 0221 — layered FO (0.22.1)

> **Renamed `0205-layered-fo` → `0221-layered-fo` (captain, 2026-06-19)** — realigns the `NNN`=version slug convention (`0198`→0.19.8, `0203`→0.20.3). The slug encoded the original 0.20.5 target, but the version line jumped to 0.22.0 at the m4 cut, so this sprint ships **0.22.1**. Active members are re-stamped `sprint: 0221-layered-fo`; archived members (`m4`, `w4`), dated debriefs, and worker-authored entity bodies keep the historical `0205` slug as a record. Sibling `0206-dispatch-driver` → `0222-dispatch-driver` (0.22.2).

> **Re-scoped 2026-06-18/19 (captain): this sprint now targets 0.22.1.** `m4` (live-team-mode-terminal-harness) + the CI version-unpin were carved OUT and **SHIPPED as `v0.22.0`** (PR #390 merged 2026-06-19; tag published; Homebrew + the `stable` branch advanced) — the team-mode/CI foundation this arc rides on. 0.22.0 delivered: m4's MERGED-semantics live lane green on current/unpinned Claude Code, and CI unpinned from 2.1.177 to installer-`latest` (the #395 keystone pin retired; `9243`/#396 `using-claude-team` merged-model support is the contract it rests on). Legacy `TeamCreate` semantics stay **best-effort in code** (the conditionally-loaded `using-legacy-claude-team` path), retired when STABLE Claude Code catches up to the merged floor. The layered-FO members below (verb core, tier-delegation, restructure, the haiku-drive gate proof) target **0.22.1**, on top of the unpinned/merged floor 0.22.0 established. See `m4-readiness.md`. (Follow-up filed from the 0.22.0 cut: `wrong-root-guard-compound-boot` — the opus e2e-gate-waiver cause.)

> **Scope-locked 2026-06-15, re-sequenced 2026-06-16 (captain) after the preflight staff review.**
> TRUE spike-first: the live Haiku drive (`w4`, haiku-loop-spike) is the only 0221 member driven now.
> The staff review (`staff-review.md`) found the two verbs cannot be safely built before the spike
> reports — the spike maps the residency and the local-extraction-vs-route-raw boundary they depend on —
> so `gate-extract-verbs` and `fo-tier-delegation` are HELD (`sprint-readiness: defer`) pending the
> spike's must-build list, then reworked against its findings AND real on-disk shapes. fo-tier's tier
> default is flipped to FAIL-SAFE (unset → level-2-only) and must own the unbuilt `SPACEDOCK_FO_MODEL`
> launcher surface — both folded into its post-spike rework. The 2y-dependent core (`next-action`,
> `state-verbs`, `merge-finalize`, the restructure) carves once 2y lands (finalizing, last for 0.20.3).
> Drivable now: `--where sprint=0221-layered-fo --where 'sprint-readiness != defer'` → the spike alone.

**Sprint:** (after scope-lock) the entities matching `sprint: 0221-layered-fo` — list current members with `spacedock status --workflow-dir docs/dev --where sprint=0221-layered-fo`. Membership and per-task state are the query, never enumerated or tracked in this doc.
**Theme:** make a weak model (Haiku) operable as the first officer by pushing the mechanical dispatch into binary verbs, so level-2 driving needs almost no judgment and level-3 judgment is delegated.

## Goal (success criterion)

A LIVE Haiku FO drives one entity end-to-end — dispatch → gate → merge → terminal — because the mechanical loop is reduced to "ask the binary what is next, do the one thing, repeat." Every level-3 judgment call (gate verdicts, design and scope, feedback-cycle escalation, conflict recovery) is DELEGATED to a standing teammate on a stronger model; no gate verdict is made by Haiku alone. Proven by the live drive, never a grep over the restructured contract.

## Why

The `fo-layered-architecture.md` vision names three layers with clean seams: L1 automation (the binary), L2 driving (run the contract by the book), L3 high-function judgment (override the automation when warranted). Today the FO contract interleaves all three. A capable model carries the whole thing and the seams stay invisible. The cost shows up two ways: an expensive model is spent on mechanical dispatch that needs no reasoning, and there is no structural guard that a weak model will not make a gate call it cannot make.

The lever is to push mechanics OUT of the model into binary verbs and trivial prose-function invocations. The per-area analysis is consistent: the event loop, state-repo ops, checklist authoring, gate extraction, and merge ceremony are MECHANICALLY COMPLEX and JUDGMENT-FREE — the load is in how many rules, not in whether any rule is right. That load is exactly what a binary absorbs. What remains after the absorption is the genuine judgment, and that gets a named delegation route. A Haiku FO that knows it is weak escalates structurally rather than by luck.

The prose-function form (prototyped in `docs/dev/_proposals/state-repo-pointer-prototype.md`) is the migration mechanism: declare an op as `«fn»(arg)` the FO follows by hand, name the binary verb it becomes, then flip the notation to backticks when the verb ships. The flow that invokes it never changes; only the body collapses from a hand-followed recipe to "call the verb." This sprint declares the prose-functions and ships the first tranche of verbs behind them.

## Definition of Done

0.22.1 ships when, merged to `main` (on top of the 0.22.0 unpinned/merged floor):

- **The `«dispatch.next-action»` mechanic as a prose-function (the binary is 0222's).** The FO's next-iteration computation (reconcile sweep, PR-state check, mod-block scan, dispatchable query) ships this sprint as `czw`'s `«dispatch.next-action»` prose-function carrying a `→ spacedock dispatch next-action` migration target. The w4 spike found `«next»` HOLDs (Haiku-operable as prose), so 0.22.1 does NOT binarize it — the `spacedock dispatch next-action` binary (the fully-qualified `{action, slug?, stage?, team_action?, reason}` return schema, with `team_action` a resolved instruction not a drift-class) is the **`0222-dispatch-driver`** sprint's deliverable. (Descoped from 0.22.1 per the 2026-06-19 sprint-wide staff review: the binary moved to 0222; the Haiku-operable prose-function form ships here in the restructure.)
- **State-repo verbs** — `spacedock state ready`, `spacedock state sweep`, and `spacedock state commit <slug>` ship behind the `«state.ensure-ready»` / `«state.sweep-merged»` / `«state.commit»` prose-functions. The commit verb is path-scoped, pushes with retry-on-reject-rebase, and REFUSES to return on a rebase conflict (exit non-zero, stderr naming the paths, state repo left in rebase-abort state) so the FO cannot proceed even if it wanted to. The rebase-conflict halt is enforced by the verb, not by FO discipline.
- **A finalize / merge-guard verb** — `spacedock merge guard <slug>` enforces the mod-block set→invoke→clear sequence atomically (set before hook, detect completion by state delta, clear in a standalone `--set`, terminalize only after), so the single highest Haiku merge risk — combining, skipping, or reordering those steps — is owned by the binary. The status tool's existing refusal of terminal-with-mod-block-set is the backstop.
- **A model-aware delegation mechanism** — the FO self-identifies its tier at boot (`«fo.tier»()` from the resolved model name or a launcher flag), and a named routing table sends every judgment-call category (gate verdicts, design/scope, feedback-cycle-3 escalation, model-mismatch reuse, conflict recovery, teardown health) to a standing level-3 teammate on a stronger model. When the FO is Haiku and the workflow has `gate: true` stages, ALL gate verdicts route to level-3; Haiku adjudicates none alone.
- **The prose-function restructure** — the boot-resident core and dispatch core express the FO loop as prose-function invocations (`«state.boot»`, `«dispatch.next-action»`, `«gate.assemble-verdict»`, `«merge.guard»`, `«feedback.route»`) with bodies that either name a shipped verb (backticks) or carry the hand-followed recipe (guillemets). The flow reads as L3 intention; the L2 mechanics live in the bodies. Verbs that do not ship this sprint stay as prose-functions with a named migration target.
- **The validating live Haiku drive** — a Haiku FO, with a standing Sonnet/Opus level-3 teammate, drives one real entity from dispatch through gate through merge to terminal. The gate verdict is made by the level-3 teammate and forwarded; the mechanical steps are made by Haiku calling the verbs. This drive is the gate proof for the sprint, run AFTER the riskiest-mechanism spike below proves the simplified loop is Haiku-operable.
- `v0.22.1` cut after the pre-cut antipattern audit is clean.

Membership (the specific verbs and the restructure slices) is the query above once carved, not enumerated here.

## Candidate members (pending scope-lock)

The decomposition below dedups the per-area proposals into buildable members for the scope-lock conversation. Each is a candidate, not a committed member. The captain decides which are in 0221 and which defer; the carve then stamps `sprint` / `group` / `sprint-readiness` and the roster becomes the query.

- **next-action-driver** (binary-verb) — `spacedock dispatch next-action`: the event-loop computation in one deterministic call, returning a fully-qualified `{action, slug?, stage?, team_action?, reason}`. Folds the reconcile sweep, PR-state check, mod-block scan, dispatchable query, and the context-budget reuse probe. Riskiest member: the return schema is the FO-facing contract and must be fully-qualified (no class-name interpretation left to the FO). Depends on 2y.
- **state-verbs** (binary-verb) — `spacedock state ready` / `sweep` / `commit <slug>`: split-root gate, merged-PR iteration, path-scoped concurrency-safe commit, and the rebase-conflict halt enforced by the verb refusing to return. Depends on 2y.
- **merge-finalize** (binary-verb) — `spacedock merge guard <slug>` (and optionally a `merge ceremony` / `worktree safe-remove` companion): atomic mod-block set→invoke→clear and the PR-vs-Ship-Local ceremony sequence. Depends on 2y (the mod-block prose lands in the host-neutral merge core post-2y).
- **gate-extract-verbs** (binary-verb) — `spacedock gate checklist-extract` / `ac-scan` / `reviewer-parse`: structured extraction of the stage report (DONE/SKIPPED/FAILED, line ranges, `chosen_direction_required` flag), AC coverage with natural-place flags for L3 routing, and tiered reviewer findings. No 2y dependency. The verdict stays judgment and routes to L3.
- **fo-tier-delegation** (delegation-mechanism) — `«fo.tier»()` self-identification, the level-3 routing table prose (judgment categories → routes), the standing level-3 teammate class spawned on demand, and the model-aware boot gate that warns when a Haiku FO meets `gate: true` stages. The mechanism that makes "weak FO, safe by construction" real. Partly depends on 2y (the seam between mechanical and judgment becomes visible once 2y lands).
- **prose-function-restructure** (contract-restructure) — express the boot-resident core, dispatch core, gate flow, and merge flow as prose-function invocations with bodies that name shipped verbs or carry the hand-followed recipe. The migration substrate; every verb member flips one body from guillemets to backticks. Depends on 2y for the merge/dispatch cores it restructures.
- **haiku-drive-validation** (validation) — the live Haiku FO drive end-to-end with the level-3 teammate, the gate proof for the sprint. Sequenced LAST; depends on every other member shipping. Depends on 2y transitively.

## Out of scope

Defined as the sprint forms (ideation / the Commander):

- **L1 binary-module-ification beyond the named verbs.** The full prose-function-to-binary migration is a multi-sprint arc; 0221 ships the first tranche (next-action, state, merge-guard, gate-extract) and leaves the rest as declared prose-functions with named targets. Verbs not in the carved roster stay guillemets.
- **Codex / Pi Haiku operability.** The live drive validates the Claude substrate (team registry, reconcile, context-budget). Codex (mailbox) and Pi (subagent) have different failure modes and no team registry; their layered-FO operability is a follow-up, not this sprint. 0221 may add short failure-mode surface notes to their adapters to prevent applying Claude-team patterns where they do not fit, but does not validate a Haiku drive on them.
- **`spacedock gate assemble-verdict` as a shipped binary — RATIFIED kept-prose (captain, 2026-06-19).** The verdict decision tree stays a prose-function this sprint: `«gate.assemble-verdict»`'s body names (a) `6re`'s extract modes for the deterministic half and (b) the L3 / `present-gate` route for the verdict half — neither a new binary. The three extract verbs ship; the verdict-assembly stays prose with L3 escalation. The verdict is judgment and the L3 route is the safety mechanism regardless of whether assembly is binarized.
- **Reuse-policy and dispatchable-policy redesign.** Concurrency caps, hold-dispatch-until-approval, and similar are workflow-defined policies, not FO inventions. 0221 makes the binary enforce policies the README expresses; it does not design new policy surfaces.

## Dependencies

**2y upstream — `shared-merge-dispatch-contract` (`2yfsf01jf15fmts7xt7w71m2`).** Finalizing on 0203 — the last item for the 0.20.3 cut (status `validation`, PR #385 open, mod-blocked on `merge:pr-merge`). It extracts the host-neutral merge and dispatch ceremony into shared cores and defines where team/registry logic ends and the FO's operating loop begins. Because it lands imminently, the 2y-dependent members start clean against merged `main` rather than against a branch. That extraction is what makes the seam between "mechanical steps Haiku can do safely" and "judgment steps it must delegate" visible. The next-action driver, state-verbs, merge-finalize, and the prose-function restructure all build ON the cores 2y produces.

**Sequencing.** 2y is MERGED and the riskiest-mechanism spike (`w4`) PASSED, so the 2y-dependent members build clean against `main` and the empirical irreducible-judgment boundary is recorded. The haiku-drive-validation is LAST — it depends on every other member.

**Boot-region landing order — RATIFIED `6re` → `72` → restructure (captain, 2026-06-19).** `6re` (step-4 rewrite), `72` (step-1.5 insert), and `czw` (the restructure) all edit the same three boot-resident regions (`## Startup`, `## Completion and Gates`, `## FO Write Scope`). They MUST land in this order: `6re` and `72` land their raw-prose edits first; the restructure lands LAST and re-expresses the merged edits into prose-function form (folding `«fo.tier»()` in as a sixth declared prose-function). This order is a Commander drive-time constraint, not just `czw`'s prose — landing the restructure before `6re`/`72` would force a rewrite collision and risk silently dropping `72`'s fail-safe property. The roster shows all three `ready`; the Commander enforces the order from this note and the package.

## Riskiest mechanism (spike first)

Before building any verb, run a LIVE Haiku drive on a hand-simulated simplified loop — the FO calls the prose-functions by hand (guillemets), a level-3 teammate stands by, and a real entity is driven dispatch → gate → merge. The spike answers the one question that would invalidate the rest of the build: **is the simplified loop actually Haiku-operable, and what is irreducibly judgment?**

The per-area analysis predicts the failure modes the spike must surface or refute:
- Haiku inventing an auto-recovery on a rebase conflict instead of halting (state-repo area's single highest risk).
- Haiku interpreting a drift-class name semantically and dispatching the wrong action (event-loop area's risk; the fix is the fully-qualified `team_action`).
- Haiku skipping the checklist or the `dispatch build` helper and dispatching bare (worker-dispatch area's risk).
- Haiku auto-approving a gate instead of routing the verdict to L3, or inferring a chosen direction it should ask for (gate area's high-consequence risk).
- Haiku confusing an idle notification with a completion signal and re-dispatching under social pressure to "do something" (team-lifecycle area's risk).

The spike is cheap relative to the build (a few drives versus several verbs). It tells us which mechanics the binary MUST own (because Haiku breaks them every time) versus which it holds reliably, and it gives the delegation routing table its real boundary instead of a guessed one. Per the validating-new-mechanisms discipline: the smallest end-to-end exercise of the riskiest path goes first, the comprehensive build after.

### Spike outcome — w4 haiku-loop-spike, PASSED 2026-06-17 (PR #393, archived)

The simplified loop IS Haiku-operable, reliably (N=3): a bare `claude --model haiku -p` FO followed «boot»→«next»→«dispatch»→«gate»→«merge» in order, routed the gate verdict to opus every run (verdict originated at L3, zero silently-absorbed judgment), reached terminal+archived+verdict 3/3. Provisional carve input (reconfirmed by the final `haiku-drive-validation` on the real `spacedock claude` shape):

- **Must-build (4 of 6 loop steps):** `«dispatch»` (Haiku forwarded the dispatch-build prompt verbatim but DROPPED `subagent_type` from the emitted Agent envelope 3/3 — deterministic; the verb must EMIT the spawn, not hand the FO an envelope to re-type), `«state.commit»` (`git -C` form paraphrased to `cd && git add` 2/3 + carries the unexercised rebase-HALT trigger), `«gate»` and `«merge»` (held mechanically but carry judgment triggers the external prior presumes must-build). `«boot»`/`«next»` HOLD → stay prose.
- **Routing boundary (narrow):** the gate verdict is the one true judgment point (routed to L3 3/3); chosen-direction / scope / ambiguity NOT-EXERCISABLE on the linear happy path; completion-vs-idle absorbed correctly (no re-dispatch).
- **Failure modes:** 3/4/5 REFUTED 3/3; modes 1 (rebase-HALT) & 2 (drift-class) NOT-EXERCISABLE, validated downstream by `state-verbs` + the final drive.

**Residency caveat (the deferred tmux team — load-bearing for the carve).** The spike used the per-judgment bare `Agent(model=opus)` residency, which works under headless `-p`. But `m4` (live-team-mode-terminal-harness, validation / PR #390) established that headless `claude -p` exposes NO team tools (TeamCreate/SendMessage) post-2.1.178 — only the interactive tmux/pty drive does, and it stays resident. So `72`'s production STANDING `level-3-judge` residency (`spawn-standing-all` + `SendMessage`) needs the tmux backend, NOT `-p`; the spike's bare-Agent provenance does NOT transfer to it. The final `haiku-drive-validation` must run under the real (tmux) backend, and team-tool availability is version/env-sensitive — to be verified in the CI env with latest Claude Code before `72`'s residency is locked.

## Open questions

- **Self-identification mechanism.** Declared model tier, a self-probe, or a launcher flag from the resolved session model? `«fo.tier»()` needs one source of truth at boot.
- **Hand-off protocol shape (spike-informed).** The spike confirms the gate verdict is the one judgment point that must route to L3; design/scope/ambiguity were NOT-EXERCISABLE on the happy path. RESIDENCY is the open fork: per-judgment bare `Agent(model=opus)` (proven under headless `-p` by the spike) vs. a STANDING `level-3-judge` teammate (`72`'s production shape — needs the tmux backend per `m4`, since `-p` has no team tools post-2.1.178). Locked by the latest-claude-code CI-env verification.
- **next-action return schema authority.** Who computes `drift.owned` for stale-branch drift — the binary (from roster membership) or a captain lookup (session semantics: bare mode? single-entity scope?)? And does a newly-unblocked entity appearing after idle hooks fire warrant a second `status --next` in the same iteration, or is one idle cycle enough?
- **Model-selection helper placement.** The stage-model + runtime-config lookup that yields the effective model for the reuse model-match check — binary or adapter? Pending 2y clarification.
- **`assemble-verdict` binarization timing.** Ship the verdict decision tree as a binary this sprint, or keep it a prose-function with L3 escalation and ship only the three extract verbs? The L3 route is the safety mechanism either way.
- **2y landing vs branch-start.** Does 0221 wait for 2y to merge to `main`, or start against the 2y branch with a rebase plan? Affects the sprint start date and the rebase risk.
