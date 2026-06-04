---
id: e3z4yjk7mfna7mktyxw7xbw7
title: "Model bare mode as a coverage variant + establish a live baseline (contract simplification may have broken it)"
status: ideation
source: "captain (2026-06-04) — bare mode (no-team FO dispatch / Degraded-Mode fallback) has NO live coverage and is not a declared variant axis (the spec's `mode` = codified/llm-live evidence path, not team-vs-bare). This session's operating-contract simplification (zd #291 extracted `## Team Creation` — which INCLUDES the bare-mode entry 'if ToolSearch returns no match, enter bare mode' — into the lazily-loaded using-claude-team skill) may have ORPHANED the bare-mode path: bare mode is the case where team setup is unavailable, so the FO may never load the skill that tells it to go bare. Design how bare mode should be modeled AND do an exploratory baseline run to see where we actually are."
score: "0.34"
worktree:
started: 2026-06-04T07:33:00Z
completed:
verdict:
issue:
---

Bare mode is the FO's no-team fallback: sequential blocking `Agent()` (no `team_name`), no `SendMessage` reuse, Degraded-Mode semantics. It activates on team-infra failure and is conceptually the teams-off path. Today it is offline-tested only (`dispatch build` emits the right bare shape) and has ZERO live coverage; worse, it is not a coverage dimension at all. AND the session's contract decomposition may have broken its bootstrapping.

## Problem

1. **No live coverage.** Every live lane runs `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` (team mode). No live test drives an FO through a workflow in bare mode — we have no evidence an FO completes a cycle without a team.
2. **Not modeled.** The scenario spec's variant axes are `runtime {claude,codex} × mode {codified, llm-live}`, where `mode` = evidence path. Team-vs-bare is not a declared dimension, so the gap is invisible to the coverage matrix.
3. **Possible regression from the contract simplification (the riskiest unknown).** zd (#291) moved `## Team Creation` — including the bare-mode entry rule and Degraded Mode — into the lazy `using-claude-team` skill, loaded via `Skill(skill=...)` AT the team-creation step. Bare mode is the case where team setup is unavailable/fails. If the FO only learns "enter bare mode" from a skill it loads while engaging the team path, a genuinely teams-off environment may never surface the fallback. The decomposition was faithfulness-audited for CONTENT, but the bare-mode BOOTSTRAPPING ordering (does the FO reach the bare-mode instruction without a working team?) is exactly the kind of seam a content diff would not catch.

## Proposed approach (ideation)

1. **Design the variant model.** Add team-vs-bare as a first-class coverage dimension — likely a `dispatch-mode {team, bare}` axis distinct from the spec's evidence-path `mode {codified, llm-live}`. Define how it composes with the existing `runtime × mode` matrix and what the cost ledger / coverage meta-tests should require. Decide whether bare is a per-scenario variant or a dedicated lane.
2. **EXPLORATORY BASELINE RUN (the spike — RUN at ideation).** Actually drive an FO in bare mode live and record where we are. Concretely: launch the FO with `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` UNSET (or otherwise force the no-team path) on a small real workflow, and observe: (a) does the FO correctly DETECT and ENTER bare mode given the post-zd contract structure (the bootstrapping question — does it surface the fallback without a working team?); (b) does it complete a dispatch→gate→(merge) cycle sequentially; (c) where exactly does it break or degrade if it does. Record the baseline verdict (works / partial / broken + the precise failure point) in the entity body. This is the riskiest unknown — exercise it FIRST; the rest of the design depends on what the baseline shows.
3. **Branch on the baseline.** If bare mode WORKS: design the standing live bare lane (AC). If bare mode is BROKEN by the simplification: this entity's deliverable becomes the FIX (re-home the bare-mode entry so it is reachable without a team — e.g. keep the bare-mode detection in the always-on skeleton, with only the team-success path deferred to the skill) PLUS the live lane that would have caught it. Either way, end with a live bare-mode check so the fallback stops being invisible.

## Acceptance criteria (seed — firm at ideation after the baseline)

- **AC-1 (seed):** A documented bare-mode coverage model — team-vs-bare as a declared variant dimension in `docs/specs/scenario-testing-principles.md`, composing with runtime × mode, with the coverage/cost shape stated.
- **AC-2 (seed):** A recorded live baseline — an actual no-team FO run with its verdict (works / partial / broken + failure point), proving the exploratory spike ran (a `Verified by: live <ref>` citation under p4's gate).
- **AC-3 (seed):** A standing live bare-mode check (a bare variant of a shared scenario or a dedicated lane) that grades the bare-mode FO on durable outcomes — OR, if the baseline shows a break, the fix that restores the bare-mode bootstrapping plus the check that pins it.

## Out of scope

- Reverting zd's decomposition (if a break is found, re-home the bare-mode ENTRY to the always-on skeleton; do not undo the team-lifecycle extraction wholesale).

## Test plan (seed)

- The baseline is a live exploratory run (p4's `livescenario` primitive is the natural authoring surface). Local live setup: build, export SPACEDOCK_BIN/SPACEDOCK_REPO_ROOT, rotate ~/.claude/benchmark-token (ping FO on 401), force the no-team path. Offline: the variant-model meta-test + any fix's unit coverage.

## Notes

Does NOT block 0.19.5 (bare mode has been live-untested for many releases — not a release regression). Connects to the live-verification line (p4 shipped the citation gate + primitive; n1a is hardening the team-mode live cycle) and the scenario roadmap (`docs/specs/scenario-testing-principles.md`). Sibling context: zd `extract-team-orchestration-skill` (#291) is the decomposition under suspicion; n1a's 1b bare-mode dispatch-path fallback is offline-tested only — consistent with this hole.

---

# Ideation deliverable (run at ideation)

## Baseline run (the spike — RUN FIRST, the riskiest unknown)

**Verdict: WORKS.** The post-#291 contract structure still reaches the bare-mode entry rule in a genuinely teams-off environment. The captain's orphaning hypothesis is **DISPROVEN** by a live run.

### Why it is not orphaned (the bootstrapping mechanism)

zd's #291 moved `## Team Creation` — including the bare-mode entry rule "if ToolSearch returns no match, enter bare mode" — out of the FO Claude runtime adapter and into the lazily-loaded `using-claude-team` skill. The fear was: bare mode is the case where teams are unavailable, so the FO might never load the skill that tells it to go bare. That fear does not materialize, because the FO runtime adapter (`skills/first-officer/references/claude-first-officer-runtime.md:7-9`) invokes the skill **unconditionally** at the team-creation step — *to decide* whether teams are available:

> At startup (after reading the README, before dispatch), invoke the generic Claude-team-harness discipline: `Skill(skill="spacedock:using-claude-team")`

It is not gated on team success. The skill is loaded precisely so its `ToolSearch(select:TeamCreate)` probe can run and route to bare mode on a no-match. The decision logic and the thing that triggers it live in the same lazily-loaded skill, so loading-to-decide is the path, not loading-only-if-teams-on.

### How the baseline was driven (reproducible)

Local live setup per the dispatch: built the binary (`go build -o spacedock ./cmd/spacedock`, `spacedock 0.19.0 (contract 1)`), launched the real `spacedock claude` front door mirroring the live-runner shape (`--plugin-dir <repo> --skip-contract-check -- -p <prompt> --permission-mode bypassPermissions --output-format stream-json --verbose --model sonnet`), clean isolated `HOME`, OAuth `~/.claude/benchmark-token`, binary first on `PATH`. Teams forced OFF via `env -u CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS -u CLAUDECODE`. Fixtures = the shared-scenario `gate-guardrail` and `rejection-flow` fixtures (`internal/ensigncycle/shared_fixtures_test.go`). All three runs exited 0 with a `result/subtype:success, is_error:false`.

### Three runs, what each proved

1. **gate-guardrail, single-entity `-p` naming the entity.** The FO never reached the team-creation step at all — no dispatch was required, so no `TeamCreate`, no `using-claude-team` load, no probe. The gate guardrail held (presented the review, did not self-approve/mutate/archive). This run does NOT exercise the bootstrapping seam — it only confirms the no-dispatch path is unaffected.

2. **rejection-flow, single-entity `-p` naming the entity.** The FO's env probe saw `CLAUDECODE=1` (the front door re-sets it for the inner session even though the launcher unset it), concluded **single-entity mode**, and per contract (`claude-first-officer-runtime.md:13` "In single-entity mode, skip team creation") went straight to a **bare** `Agent()` dispatch (`team_name` ABSENT, `name` ABSENT). The dispatched ensign applied the fix marker and committed; the FO completed the dispatch→feedback-route→completion cycle and reported correctly. Durable outcome confirmed: the fix marker landed in the entity. **Caveat:** this proved bare *dispatch* works, but it reached bare via the single-entity bypass, NOT via the team-creation probe — so it still does not exercise the #291 seam.

3. **rejection-flow, prompt forcing the non-single-entity team-creation path.** This is the decisive run. The FO ran its full startup, invoked `Skill(skill="spacedock:first-officer")` then `Skill(skill="spacedock:using-claude-team")`, ran `ToolSearch(query="select:TeamCreate")` → **"No matching deferred tools found"**, then narrated verbatim: *"Mode resolved: Bare mode. ToolSearch(query=\"select:TeamCreate\") returned no match … entering bare mode for the entire session,"* and emitted the contract's verbatim `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` restart hint. It then dispatched a **bare** `Agent()` (`team_name` ABSENT, `name` ABSENT, `model` omitted), the ensign applied the fix marker, and the FO completed the cycle and reported the resolved dispatch mode. Durable outcome confirmed: marker present in the entity.

### Precise findings

- **No degradation point.** The probe-and-fallback chain (`first-officer` → `using-claude-team` → `ToolSearch` no-match → bare entry → bare `Agent()` dispatch → blocking completion → report) executed cleanly in a teams-off env.
- **The contract is correct as written.** The unconditional `Skill(skill="spacedock:using-claude-team")` invocation at the team-creation step is what keeps the bare-mode entry reachable. No re-homing of the entry rule is needed.
- **Single-entity mode is a second, independent bare path** that legitimately bypasses team creation (the `CLAUDECODE=1` + named-entity `-p` trigger). A live bare-mode check must drive the **team-creation-probe** path (run 3's shape), not the single-entity path, or it will pass without exercising the #291 seam — this is the single most important authoring constraint for AC-3.

Therefore the deliverable branches to the **WORKS** arm: model the variant + add a standing live check that pins the probe-path bare-mode entry. No fix is needed; the standing check exists to keep it from silently regressing in a future decomposition.

## Variant model design (AC-1)

### The new axis: `dispatch-mode {team, bare}`

Add `dispatch-mode` as a first-class variant axis in `docs/specs/scenario-testing-principles.md`, **distinct from** the existing `mode {codified, llm-live}` (which is the *evidence path*, not the orchestration shape). The spec's `## Variant Axes` becomes:

```text
scenario × {executor, mode, runtime, dispatch-mode, model?}
```

- **dispatch-mode** — the FO's orchestration shape for the observation: `team` (TeamCreate succeeded; concurrent dispatch, SendMessage reuse) or `bare` (no team; sequential blocking `Agent()`, fresh re-dispatch per stage). It is orthogonal to `mode`: a `bare` observation can be graded by either the `codified` or the `llm-live` evidence path. It is orthogonal to `runtime`: `bare` is meaningful on `claude` (the ToolSearch-no-match path) and degenerate on `codex` (Codex has no team concept — every Codex dispatch is already sequential, so `codex` collapses to a single `dispatch-mode` and the matrix should mark `codex × team` as N/A rather than a missing cell).

### How it composes with the existing matrix

The baseline coverage matrix per shared scenario becomes:

```text
runtime {claude, codex} × mode {codified, llm-live} × dispatch-mode {team, bare}
```

with two declared N/A reductions so the matrix is not falsely sparse:
- `codex × * × team` is **N/A** (Codex has no team primitive; its only dispatch-mode is bare-equivalent-sequential). Codex carries one dispatch-mode column.
- `claude × codified × bare` collapses onto the same recorded-stream consumer the `team` codified observation uses **unless** the codified check specifically models the bare dispatch shape (the `bare_mode: true` / `team_name: null` `dispatch build` emission, already offline-tested in `internal/ensigncycle/cycle_test.go` and `feedback_test.go`). The bare dispatch *emission* already has codified coverage; what is missing is the `claude × llm-live × bare` cell — the live producer proof that an actual FO resolves and enters bare mode. That one missing cell is exactly what the baseline filled and what AC-3's standing check pins.

### Per-scenario-variant vs dedicated lane — DECISION: dedicated bare lane (one scenario), not per-scenario variants

Reasoning: bare mode is a *bootstrapping/orchestration* property of the FO, not a per-journey behavior. Running every shared scenario (gate-guardrail, rejection-flow, merge-hook-guardrail) a second time under `dispatch-mode: bare` would roughly double live cost for near-zero marginal signal — the guardrails themselves (gate-hold, merge-hook refusal) are dispatch-mode-independent; they exercise `status`/contract logic, not the team-vs-bare seam. The single load-bearing bare-mode claim is "the FO resolves and enters bare mode via the ToolSearch-no-match probe, then completes a dispatch→completion cycle." **rejection-flow is the right journey** to carry that claim because it forces a real dispatch (run 1 showed gate-guardrail does not). So:

- Add **one** new shared scenario `bare-mode-dispatch` (or extend the table with a `dispatch-mode`-tagged variant of `rejection-flow`) whose `claude` runner launches with teams OFF and the non-single-entity prompt shape (run 3), and grades the durable outcome PLUS the bare-mode entry signal.
- Keep the other shared scenarios at `dispatch-mode: team` (their CI default) — do not fan them out to bare.

### Cost-ledger / coverage-meta shape

- **Cost ledger** keys by `(scenario, mode, runtime, dispatch-mode, model?)`. The host-neutral scenario identity rule (no `claude-`-prefixed keys) extends to dispatch-mode: `dispatch-mode` is a variant column, never folded into the scenario name or into `mode`.
- **Coverage meta-test** (the `sharedRuntimeScenarios()` ↔ doc lock, `internal/ensigncycle/shared_scenarios_docs_test.go` + `shared_coverage_meta_test.go`) gains a `dispatch-mode` dimension: it asserts the bare lane has a `claude` runner and that `codex × team` is declared N/A rather than silently absent. The seed-scenario doc block grows to include the bare lane so the doc↔code lock reds on drift in either direction.

## Acceptance criteria (firm)

- **AC-1:** `docs/specs/scenario-testing-principles.md` declares `dispatch-mode {team, bare}` as a variant axis distinct from `mode`, states how it composes with `runtime × mode` (including the `codex × team` N/A and the `claude × llm-live × bare` missing-cell that the new lane fills), and records the per-scenario-variant-vs-dedicated-lane decision (dedicated bare lane on rejection-flow's journey). *Verified by:* a meta-test in `internal/ensigncycle` (extending `shared_scenarios_docs_test.go`'s doc↔code lock) that parses the spec's variant-axes block and fails if `dispatch-mode` is absent or if the declared bare lane is not present in the code-side scenario table — i.e. the doc claim is bound to a real value, not prose about itself.

- **AC-2:** The live baseline is recorded in this entity body with verdict **WORKS** and the precise probe-path mechanism, citing the run-3 stream evidence (the `ToolSearch(select:TeamCreate)` → "No matching deferred tools found" → "Mode resolved: Bare mode" → bare `Agent()` → marker-applied chain). *Verified by:* this section plus the new standing live lane's first run under p4's citation gate (a `Verified by: live <ref>` citation — the LLM-executor run that only the live path can decide), so the baseline claim is backed by a re-runnable live observation, not a one-time transcript.

- **AC-3:** A standing live bare-mode check exists as a shared scenario whose `claude` runner (a) launches with `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` unset AND a non-single-entity prompt shape so it exercises the **team-creation probe** path (NOT the single-entity bypass), and (b) grades on durable outcomes: the entity's fix marker is applied (the dispatch→completion cycle ran) AND the FO resolved bare mode via the probe. Because the FO's bare-mode entry narration is transcript phrasing (non-deterministic, not gradeable per the spec), the gradeable bare-mode signal is the **bare dispatch shape** observable in the stream — a `team_name`-absent `Agent()` tool_use and the absence of any `TeamCreate` tool_use — captured the way `streamwatch` already classifies bare-mode Agent tool_results (`internal/ensigncycle/streamwatch_test.go:428` `looksLikeBareDone`). *Verified by:* the new `//go:build live` claude runner for the bare lane (red if the marker is absent or a `TeamCreate` tool_use appears), plus the parity/coverage meta-test that fails if the bare lane lacks a claude runner.

## Riskiest-unknown determination

The riskiest unknown was the #291 bootstrapping seam (is the bare-mode entry reachable in a teams-off env after the extraction?). It was exercised FIRST, live, in run 3 — verdict WORKS, mechanism recorded above. The remaining implementation rests on already-proven mechanisms: the `livescenario`/shared-scenario runner harness (p4, already shipping the three seed scenarios live), the offline bare-dispatch emission (`cycle_test.go`/`feedback_test.go`, already green), and `streamwatch`'s bare-mode tool_result classification (already shipping). No further spike needed for those — they are demonstrated, not asserted.

## Test plan (firm)

- **Live (the producer proof, `claude × llm-live × bare`):** the new standing bare lane — one `//go:build live` claude runner driving the rejection-flow fixture with teams OFF + the non-single-entity prompt, asserting marker-applied + no-`TeamCreate`-tool_use + `team_name`-absent `Agent()`. Cost class ≈ the existing rejection-flow live scenario (~4 min, one sub-dispatch). Local setup as recorded in the baseline section.
- **Offline (cheap, always-on):** (a) the AC-1 doc↔code lock meta-test for the `dispatch-mode` axis; (b) the AC-3 parity meta-test (bare lane has a claude runner; `codex × team` declared N/A); (c) reuse the existing offline bare-dispatch-emission tests as the codified executor of the bare lane — no new live spend to prove the *emission* shape, only to prove the live *producer* resolves bare mode.

## Stage Report: ideation

- DONE: Run the EXPLORATORY BASELINE FIRST (the riskiest unknown): drive a real FO in BARE mode live and record where we are.
  3 live `spacedock claude` runs (teams OFF), all exit 0 / result success. Run 3 (non-single-entity prompt) drove the team-creation probe path: `Skill(using-claude-team)` → `ToolSearch(select:TeamCreate)` → "No matching deferred tools found" → "Mode resolved: Bare mode" + verbatim AGENT_TEAMS hint → bare `Agent()` (team_name/name absent) → ensign applied fix marker → cycle completed. Verdict WORKS; mechanism + per-run findings in "# Ideation deliverable / ## Baseline run".
- DONE: Design the bare-mode coverage MODEL: add team-vs-bare as a first-class variant dimension distinct from evidence-path `mode`.
  `dispatch-mode {team, bare}` axis designed in "## Variant model design (AC-1)": composes as runtime × mode × dispatch-mode with `codex × team` = N/A and `claude × llm-live × bare` as the missing cell the baseline fills; decision = dedicated bare lane on rejection-flow's journey (not per-scenario fan-out); cost-ledger + coverage-meta shape stated.
- DONE: Branch the deliverable on the baseline + firm AC-1/AC-2/AC-3.
  Baseline = WORKS → branched to the standing-live-check arm (no fix needed; the unconditional `Skill(using-claude-team)` invocation keeps the entry reachable). AC-1/AC-2/AC-3 firmed with entity-level "Verified by" naming meta-tests + a `//go:build live` claude bare-lane runner; out-of-scope (no zd revert) preserved.

### Summary

Disproved the captain's orphaning hypothesis with a live baseline: after #291 moved the bare-mode entry rule into the lazily-loaded `using-claude-team` skill, the FO still reaches it because the runtime adapter invokes that skill unconditionally at the team-creation step to DECIDE team availability — `ToolSearch(select:TeamCreate)` no-match routes to bare mode and a clean bare dispatch→completion cycle (verdict WORKS, no degradation point). Designed `dispatch-mode {team, bare}` as a variant axis distinct from evidence-path `mode`, chose a dedicated bare lane over per-scenario fan-out, and firmed three entity-level ACs whose proofs are a doc↔code lock meta-test, a parity meta-test, and a live claude bare-lane runner. Key authoring constraint recorded: the standing check MUST use the non-single-entity prompt shape (run 3), since single-entity `-p` reaches bare via a different bypass that does not exercise the #291 seam.
