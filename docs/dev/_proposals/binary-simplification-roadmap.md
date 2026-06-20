# Binary Simplification Roadmap — Reducing FO + Ensign Prose-Contract Load

**Status:** Proposal, FO-drafted 2026-06-02 from sprint-end strategic question.
**Author:** First Officer (session 9, post-0.19.4 sprint).
**Context:** Captain observation mid-sprint — "we shipped a lot of deterministic logic in the helper since we started the prose-based contract. what might we simplify further? because right now it is heavy for both FO and ensigns."

## Refresh — 2026-06-04 (post-0.19.5 wave)

This session shipped two token-relevant contract decompositions and established a **second lever** the original draft missed.

### Two levers, not one
1. **Binary-command migration** (the original roadmap below): move HOW-discipline from prose into atomic binary ops. Still unshipped — remains the Phase 1–4 plan.
2. **Contract decomposition into lazy-loaded skills** (NEW, proven this session): lift a *generic, self-contained* block out of the always-loaded boot read into a skill pulled on demand via `Skill(skill=...)`. It doesn't cut total tokens when the path is hit; it defers the load off the boot read and makes the block reusable. zd (#291) proved the mechanism — a cross-skill `Skill()` handoff lands the body in-context mid-run; the `@`-include path was disproven and the skill must ship in the plugin (registers at boot).

### Shipped this session
- **zd #291** — extracted the 4 generic team-lifecycle blocks (Team Creation, Awaiting Completion, Terminal Teardown, Degraded Mode) out of `claude-first-officer-runtime.md` into a lazy `using-claude-team` skill. FO runtime adapter **305→224 lines (~8.8k→6.7k tok)**; the 4.4k-tok team lifecycle now loads only on first team creation, not at boot.
- **ep #290** — re-homed dev-only discipline (TDD, code-only deliverable, "CODE only" worktree) out of the **universal** `ensign-shared-core.md` into `development.md` (opt-in). The ensign shared core is the highest-frequency contract (loaded M× per session, fresh per dispatch), so dehydrating it is the highest per-session leverage of the two.

### Measured contract token state (2026-06-04, chars/4 estimate)
| Contract | ~tokens | load frequency |
|---|---:|---|
| FO boot read = SKILL 372 + **first-officer-shared-core 9,730** + claude-first-officer-runtime 6,737 | **~16,840** | 1× / session |
| `using-claude-team` skill | 4,363 | lazy — first team creation |
| Ensign contract (SKILL 106 + shared-core 2,003 + claude-runtime 639 + codex-runtime 597) | ~3,345 | M× / session (per dispatch) |
| codex-first-officer-runtime | 1,185 | Codex sessions only |

**Top remaining target: `first-officer-shared-core.md` (9,730 tok) is now the single largest boot-read file** — it overtook the runtime adapter after zd. Highest-value next move on BOTH levers:
- *Dehydration:* the original Phase 0.A prose pass, still unshipped on the FO side (ep did the ensign side).
- *Decomposition:* candidate self-contained blocks to lazy-extract — **Gate Presentation** (template + captain-facing assembly rules), **Merge-and-Cleanup** ceremony, **Feedback Rejection Flow**. Each is needed only at specific event-loop points, not every boot.

### Re-prioritization (supersedes the phasing table below for the next sprint)
- **Top move:** FO-side dehydration **+** a decomposition pass on `first-officer-shared-core.md` — 9.7k tok loaded every session, the largest single lever left.
- **Strongest binary-command ROI: #3 `spacedock pr complete`.** This session I ran its ceremony **~7× MANUALLY** (at/n3/2a/am/6b/ep/zd merges) — the heaviest, most error-prone sequence, exactly #3's target. #1 (`state sync`) and #2 (`dispatch advance`) each fired 10+× this session too. The hot-path evidence is now overwhelming.
- The two levers are complementary: decompose generic prose → skills (lever 2); migrate mechanical HOW → binary (lever 1); what remains in the always-loaded core is spacedock-specific WHAT/WHEN policy.

### Third lever — guarantee → code-gate (xa determination, 2026-06-04)

The a9 (`feedback-rejection-flow-skill-extraction`) detached audit surfaced a **third** lever beside the two above: **promote an FO-behavioral guarantee that lives only as contract prose into a binary-enforced gate over durable on-disk state.** A prose-only guarantee has a ceiling of "the wording is present"; where the guarantee is mechanizable, a guard (like the `mod-block` / terminal-transition guard) closes the body-vs-label gap entirely. This is distinct from lever 2 (which *moves* prose to a lazy skill but keeps it prose) — lever 3 *replaces* prose with enforcement.

**xa (`feedback-guarantee-binary-gate`) evaluated the two feedback-rejection guarantees the audit flagged.** Determination (captain decided 2026-06-04 to roadmap the decision + fork a separate build task; xa closed as a roadmap decision, not a PASSED dev task):

- **Candidate 1 — 3-cycle escalation: MECHANIZABLE (the count), via a dedicated cycle-record command — NOT a `--set status` guard.** Two spikes settled the hook: (a) `--set status=implementation` carries no feedback-bounce signal — the disambiguating `is_feedback_reflow` lives on the dispatch-build input path, not as durable state — so a status-transition guard would false-fire on legitimate forward re-entry; (b) section-scoped counting of `### Feedback Cycles` entries is deterministic and tamper-evident. The correct hook is the cycle-record WRITE itself (unambiguously a bounce): a `spacedock status --record-feedback-cycle {slug}` subcommand that owns the `### Feedback Cycles` append, computes the post-append cycle number, stamps a durable escalation marker (frontmatter, e.g. `feedback-escalate: cycle-3`) on the threshold, and refuses a further auto-bounce (exit 1, `--force` override) — the same prose→binary promotion `mod-block` models. The FO's RESPONSE to a refusal (actually escalating) stays a live concern → gq's `feedback-3-cycle-escalation` scenario. **Build forked to `feedback-cycle-record-command` (`bwr6j6edkmfx5sbz73cr2952`).**
- **Candidate 2 — budget-probe fail-safe: NOT MECHANIZABLE — stays prose + live scenario.** The reuse-vs-fresh decision produces identical durable on-disk state either way (the difference is which live runtime handle the FO messages, leaving no on-disk trace), so there is no transition for a guard to sit in front of. Stays reuse-condition-0 prose; live coverage defers to gq (`feedback-nonhappy-live-coverage`).

The mechanizability bar this established (reusable for future lever-3 candidates): a guarantee mechanizes only if (a) there is a durable on-disk transition to sit in front of, AND (b) the guard's decision is computable from durable on-disk state — not from a model intent that leaves no on-disk trace.

## Measured + rejected — checklist externalization to a todo/task tool (3c0b, 2026-06-04)

Considered hinting/forcing ensigns to externalize the dispatch `### Completion checklist` into a todo/task tool (for mid-stage captain visibility + a smaller prompt). Measured empirically before committing — a 3-arm live opus A/B, 3 trials/arm (full data in archived `ensign-todo-task-tool-hint`):

- **No completeness or correctness gain on any arm** (all ~11/11 report-accounting; all pass `go test`).
- **Ensign-side hint** raises externalization 0/3→2/3 but at **+94% cost / +108% turns** for the externalizing trials.
- **Spawner-side pre-creation** (FO `TaskCreate`s, owner=ensign) does NOT invert that — it relocates AND grows the cost (system-wide **42** task-ops vs the hint's 15.7; worker slowest at 137s; **+ a ~7min/trial team-teardown tax**), though it makes externalization reliable (3/3) and shrinks the prompt 62%.
- The ONLY benefit, on every arm, is Claude-team mid-run visibility (Shift+Down) — a human-UX nicety a transcript metric can't value, and never free.

**Decision:** not shipped (parked here). If reliable mid-run visibility is later judged worth it, the spawner arm is the path — as an explicit FO-side opt-in. Reusable artifact: the per-turn jsonl token/wallclock parse + offline negative-case grading is a measurement method for any "does this nudge earn its cost" question. Twice in this spike a plausible assumption (the hint helps; spawner-side inverts the cost) was measured and refuted before it shipped.

## Pattern

The prose contract started thin (capabilities lived in scripts), then grew dense as Spacedock's capabilities became real. Every helper we shipped moved discipline from prose-instruction → code-enforcement:

- `status` mutator → entity-frontmatter discipline owns itself
- `dispatch build` → ensign-prompt assembly owns itself
- `context-budget` → reuse-decision-0 owns itself
- `reconcile` (qs, just landed) → 5-class drift detection owns itself; FO acts on the report

The natural endpoint: **FO + ensigns are thin orchestrators over atomic binary operations.** Prose carries WHAT (state transitions, when to act, captain involvement); the binary owns HOW (ceremony sequencing, retry-on-rejection, idempotency).

## Prioritization framing (captain insight, 2026-06-02)

**Per-load token economics:**
- FO contract: loaded ~1× per session (long-living, no fresh re-load between entities)
- **Ensign contract: loaded ~M× per session** (fresh-spawned each dispatch, no context carryover)
- **Shared contract** (FO-shared-core + ensign-shared-core): loaded **2× per dispatch cycle** (FO when it dispatches, ensign when it boots)

This session: ~3 FO contract loads vs. ~15+ ensign loads vs. ~20+ shared-contract loads. The shared and ensign contracts are higher-leverage targets than the FO event loop because every dispatch pays the cost again.

This roadmap captures the next 11 candidate moves, **prioritized by per-load impact × load frequency**, NOT by FO cognitive load alone.

## Phase 0 — Shared + ensign contract (HIGHEST PER-DISPATCH IMPACT)

These pay off every single ensign dispatch.

### A. Shared-contract prose pass (the qs cycle-3 pattern, generalized)

**Replaces:** audit-trail exposition + cross-file restatement + over-qualification in `first-officer-shared-core.md`, `ensign-shared-core.md`, `claude-first-officer-runtime.md`, `claude-ensign-runtime.md`, `codex-first-officer-runtime.md`, `codex-ensign-runtime.md`. The AI engineer pattern that cut qs reconcile prose by 55% almost certainly applies sweepingly across the existing contract.

**Estimated prose impact:** SHARED-CORE is ~600 lines + ENSIGN-SHARED-CORE is ~200 lines + the 4 runtime adapters are ~250 lines each. If we hit even a 25–30% reduction (more conservative than the 55% qs cycle-3 number, because the existing prose is older and probably tighter), that's **~400 lines removed from EVERY ENSIGN DISPATCH'S read**.

**Logic readiness:** N/A — this is pure prose. The AI engineer agent + the host-neutrality oracle + the portability oracle are the tools. The qs cycle-3 is the proof-of-concept.

**New code needed:** none.

**Work shape:** dispatch a senior AI engineer agent over each contract file, apply the verbatim rewrite pattern. Adversarially verify with the existing oracles (host-neutrality, portability) + a full test run.

**Risk:** prose has accreted load-bearing nuance that a casual cut would lose. Mitigation: AI engineer is read-only-propose; FO routes to an implementer; oracles catch token-level regressions.

### B. `spacedock dispatch report {stage} --item N --result DONE/SKIPPED/FAILED --evidence "..."` (ensign-side)

**Replaces:** Manual stage-report templating in ensigns. The DONE/SKIPPED/FAILED structure + the indented evidence line are pure mechanical templating. Currently ~10 lines of ensign-shared-core prose tells every ensign how to format every line.

**Estimated prose impact:** ~10 lines of ensign contract prose at every dispatch.

**Logic readiness:** MEDIUM. Need a template + an appender that knows the entity body location (worktree vs state checkout per the split-root rule). The split-root logic already exists in `dispatch build`'s state-checkout path injection.

**New code needed:** ~80 LOC.

**Tests:** fixture-based; assert the appended report has the right shape + lands at the right path.

### C. `spacedock dispatch signal-complete {entity_path}` (ensign-side)

**Replaces:** the SendMessage "Done: ..." format every ensign must remember.

**Estimated prose impact:** ~3 lines of ensign contract prose.

**Logic readiness:** HIGH. Format is fixed; just wraps SendMessage to `team-lead` with the templated body.

**New code needed:** ~40 LOC.

### D. `spacedock state commit-and-push {entity_path} -m "..."` (universal — FO + ensign)

**Replaces:** the path-scoped commit + push + pull-rebase-on-rejection + halt-on-conflict pattern. Touched every time the state checkout is written. The discipline is currently encoded in shared-core's State Management section (~30 lines) AND in every dispatch's prose for ensigns.

**Estimated prose impact:** ~30 lines of shared-core prose + ~10 lines of per-dispatch ensign prose. Touched 5+ times per dispatch cycle.

**Logic readiness:** HIGH. The discipline is mechanical: stage path-scoped, commit with given message, push, retry on rejection with pull-rebase, halt on conflict.

**New code needed:** ~80 LOC.

**Tests:** simulated push rejection + simulated rebase conflict.

**Idempotency:** safe to re-run on a no-op state (returns 0 with nothing to do).

## Phase 1 — FO event loop hot path (per-FO-session impact)

Touched many times per session but loaded once per session, so each cycle's marginal gain is lower than Phase 0.

### 1. `spacedock state sync`

**Replaces:** state checkout pull-rebase + retry-on-rejection + halt-on-conflict. Atomic operation; one command instead of three with manual retry. (Could fold into Phase 0 #D if scoped narrowly.)

**Estimated prose impact:** 4 lines per occurrence; touched 8–15× per session. ~30–60 lines of FO cognitive load per session.

**Logic readiness:** HIGH. Discipline is already prose-encoded; git operations are mechanical. Rebase-conflict halt is a documented requirement.

**New code needed:** ~50 LOC.

### 2. `spacedock dispatch advance {slug} --to {stage}`

**Replaces:** `status --set` field update + state commit + state push + worktree-create-if-needed + supersede-shutdown invocation. Every stage transition.

**Estimated prose impact:** 6 lines per occurrence; touched 4–8× per session.

**Logic readiness:** HIGH. Stage-transition rules are encoded in the FO contract + qs reconcile (Class B). Worktree create is `git worktree add` with conventions we control.

**New code needed:** ~150 LOC.

### 3. `spacedock pr complete {slug}`

**Replaces:** Merge-and-Cleanup post-merge half: mod-block clear (separate commit per the standalone rule) + terminalize + archive + worktree-remove + branch-delete + main-reset + binary-rebuild. THE HEAVIEST CEREMONY in the contract.

**Estimated prose impact:** 25 lines per entity merge.

**Logic readiness:** HIGH.

**New code needed:** ~200 LOC.

## Phase 2 — Gate work

### 4. `spacedock gate verify {slug}`

### 1. `spacedock state sync`

**Replaces:** state checkout pull-rebase + retry-on-rejection + halt-on-conflict. Atomic operation; one command instead of three with manual retry.

**Estimated prose impact:** 4 lines per occurrence; touched 8–15× per session. **~30–60 lines of FO cognitive load per session.**

**Logic readiness:** HIGH. Discipline is already prose-encoded; git operations are mechanical. Rebase-conflict halt is a documented requirement (FO contract State Management → Rebase-conflict halt B6).

**New code needed:** ~50 LOC. Wraps `git -C {state} pull --rebase` with retry + halt-on-conflict + clean exit-code surface.

**Tests:** state-checkout fixture + simulated push rejection (existing pattern for git tests in `internal/dispatch/reconcile_e_test.go`).

**Idempotency:** trivial; re-running converges.

### 2. `spacedock dispatch advance {slug} --to {stage}`

**Replaces:** `status --set` field update + state commit + state push + worktree-create-if-needed + supersede-shutdown invocation. Every stage transition.

**Estimated prose impact:** 6 lines per occurrence; touched 4–8× per session. **~30–50 lines per session.**

**Logic readiness:** HIGH. Stage-transition rules are encoded in the FO contract (Completion and Gates) + the qs reconcile helper (Class B). Worktree create is `git worktree add` with conventions we control.

**New code needed:** ~150 LOC. Combines existing `status`, `dispatch build`, and git operations.

**Tests:** behavioral fixtures over each stage-transition class (gated, fresh, reuse, feedback-rejection).

**Idempotency:** advance should be no-op if already at target stage.

### 3. `spacedock pr complete {slug}`

**Replaces:** Merge-and-Cleanup post-merge half: mod-block clear (separate commit per the standalone rule) + terminalize (`--set completed verdict=PASSED worktree=`) + archive + worktree-remove + branch-delete + main-reset + binary-rebuild. THE HEAVIEST CEREMONY in the contract.

**Estimated prose impact:** 25 lines per entity merge. Touched 3–5× per sprint.

**Logic readiness:** HIGH. Every step is an existing binary operation; the only constraint is the SEPARATE-COMMIT discipline for mod-block clear, which the `status --set` guard already enforces.

**New code needed:** ~200 LOC. Orchestrates existing commands in a transactional sequence; surfaces failures with the same shape `pr-merge` mod uses today.

**Tests:** end-to-end fixture from merged-PR state to archived entity. Should pass through every guard.

**Idempotency:** re-running on an already-archived entity should be a clean no-op.

## High leverage (gate work)

### 4. `spacedock gate verify {slug}`

**Replaces:** FO's AC cross-check at every gate. Pulls every `**AC-N**` from entity body, scans stage reports for evidence citations, flags missing-evidence ACs.

**Estimated prose impact:** ~20 lines of FO contract prose at every gate. Touched 8–15× per session.

**Logic readiness:** MEDIUM-HIGH. AC + evidence parsing is similar to what `status --validate` already does for entity-format validation. The parser already exists; the evidence-scanner is the new piece.

**New code needed:** ~100 LOC. AC extractor (already partial in `--validate`) + evidence-citation scanner + report renderer.

**Tests:** fixtures with passing/failing AC-evidence patterns. The existing `--validate` test set is the starting point.

**Stretch:** integrate the dev-workflow self-reference guard (entity `2a`, parked-approved this session) into the same code path.

### 5. `spacedock pr build {slug}`

**Replaces:** Manual extraction of title + body per the pr-merge template; audit-link short-SHA + short-ID + owner/repo computation; word-count check.

**Estimated prose impact:** ~30 lines of pr-merge mod prose per PR. Touched 3–5× per sprint.

**Logic readiness:** MEDIUM. The template is well-specified (pr-merge.md). The extractor needs to handle documented patterns (motivation lead from first paragraph; what-changed from implementation DONE items; evidence from validation stage report).

**New code needed:** ~150 LOC. Section parser + template renderer + GitHub API wrapper for owner/repo/short-SHA.

**Tests:** fixtures with several entity body shapes; golden PR bodies for diff stability.

**Idempotency:** repeated invocation produces byte-identical output.

## Medium leverage (orchestration)

### 6. `spacedock dispatch reconcile --act`

**Replaces:** Per-class action table in FO contract. Reconcile currently DETECTS; this would make it also ACT (gated by `--act` to preserve the dry-run-by-default safety).

**Estimated prose impact:** ~15 lines of per-class action prose. Touched at boot/idle/post-merge.

**Logic readiness:** HIGH (extends qs). Each class action is a well-defined binary operation:
- A (lingering): SendMessage shutdown_request (Claude only — Codex has no analog, document the skip)
- B (superseded): same as A
- C (un-advanced PR): trigger `pr complete` on the slug (depends on #3)
- D (stale branch): `git pull --rebase`
- E (stale local main): `git fetch + reset + rebuild`

**New code needed:** ~80 LOC. Extends reconcile with `--act` flag.

**Tests:** extend `reconcile_test.go` with `--act` assertions; the existing 5-class fixture is the starting point.

**Constraint:** `--act` must be conservative — surface every action it took on stdout so the FO/captain can audit.

## Lower-leverage candidates (diminishing returns; queue last)

### 7. `spacedock ci approve {pr}`

Wraps `gh api repos/{owner}/{repo}/actions/runs/{run}/pending_deployments` with env-ID lookup. ~50 LOC. **Prose impact:** 3 lines per occurrence; touched 2–3× per PR. **Readiness:** HIGH.

### 8. `spacedock audit dispatch {entity}`

Wraps "create detached worktree + spawn audit workflow" for high-stakes entities. ~120 LOC (incl. lens template). **Prose impact:** 10 lines per audit; touched 2–4× per sprint. **Readiness:** MEDIUM (workflow script structure is established but per-entity prompts vary).

### 9. `spacedock dispatch report {stage} --item N --result DONE --evidence "..."` (ensign-side)

Templates stage reports. ~80 LOC. **Prose impact:** ~10 lines of ensign contract prose. **Readiness:** MEDIUM (needs template + appender that knows entity body location — worktree vs state checkout).

### 10. `spacedock dispatch signal-complete {entity}` (ensign-side)

Formats the SendMessage "Done: ..." completion signal. ~40 LOC. **Prose impact:** 3 lines. **Readiness:** HIGH.

## Recommendation — phasing

| Phase | Candidates | LOC est | Prose lines saved | Cadence |
|-------|------------|---------|-------------------|---------|
| **1** (next sprint) | #1 + #2 + #3 | ~400 | ~150 per session | Hot-path event-loop work. Highest cognitive-load win. |
| **2** | #4 + #5 | ~250 | ~60 per gate/PR | Gate work + PR assembly. |
| **3** | #6 | ~80 | ~15 per session | Natural follow-up to qs; touches the same code. |
| **4** | #7 + #8 + #9 + #10 | ~290 | ~25 per occurrence | Ensign-side + niche. Queue when binary-surface-vs-cognitive-load trade-off reverses. |

**Total Phases 1–3:** ~730 LOC of new binary code; ~225 lines of prose contract savings per session.

The savings compound: each command moved removes a class of mistakes the FO/ensign could make.

## Trade-offs

- **Binary surface grows.** Every command needs tests, version compatibility check, contract gate. We already have the discipline for this (the parser swap entity `zj` just landed enforces yaml.v3-canonical formatting; the reconcile entity `qs` just shipped Class A/B/C/D/E detection backed by 7 reconcile tests + fixture coverage).
- **Idempotency thinking required.** Re-running `pr complete` twice should be safe. Re-running `state sync` should converge. These become new contracts in their own right.
- **The prose contract becomes more declarative.** It says WHAT exists and WHEN to act; the binary enforces HOW. This is a one-way migration — once a command exists, the prose contract should NOT duplicate its discipline. The qs cycle-3 prose simplification (PR #273) is the first explicit example of this rule.
- **Lower onboarding cost.** A new FO/ensign instance has less to memorize; the binary does the discipline work. This is consistent with the recalibrated sprint goal "fresh install works from `spacedock-dev/spacedock@next`" — both are about a more self-sufficient binary.

## What this is NOT

- **NOT a rewrite of the prose contract.** The contract stays for policy + when-to + escalation. Only the HOW gets moved.
- **NOT a single sprint.** Multi-sprint roadmap; each phase is its own dispatched entity.
- **NOT a replacement for the human gate.** Captain + FO still make decisions; the binary owns the mechanics around those decisions.
- **NOT a universal-binary-behavior expansion.** Several of these (e.g., `pr complete`, `ci approve`) are dev-workflow-specific. Some (e.g., `state sync`) are universal. The README `merge:` policy pattern is the precedent for opt-in workflow-bound behaviors.

## Next-step concrete proposal

If the captain greenlights Phase 1, file three implementation entities in the backlog (one per command in #1, #2, #3) and dispatch them sequentially since they share `internal/dispatch` + `internal/status` surfaces. The qs reconcile helper is the template — both for the helper shape and for the test-coverage discipline that landed.

If the captain wants a different phasing or scope (e.g., only #3 first because Merge-and-Cleanup is the most error-prone ceremony), that's a single-line direction.
