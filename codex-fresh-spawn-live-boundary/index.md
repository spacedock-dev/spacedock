---
title: Restore Codex FO contract after compaction
status: implementation
source: "Captain-directed v0.25.2 follow-up, 2026-07-20; live FO escape after v0.25.1 / archived rt8 / PR #532"
score: "1.0"
milestone: 0.25.2
started: 2026-07-20T04:01:57Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-codex-fresh-spawn-live-boundary
issue:
id: 6cc3rvfd44y6x3352hh21v8b
---

The title and source are historical frontmatter. Captain review withdrew the generic fresh-spawn premise: the cited direct `spawn_agent` call was ad-hoc research, not a Spacedock worker dispatch, so it is not evidence that the v0.25.1 dispatch boundary failed. No global `PreToolUse` fork guard belongs in this ticket.

The actual v0.25.2 defect is the archived c6 Codex post-compaction reload contract. c6 shipped a `PostCompact` hook whose `systemMessage` tells the captain to make the First Officer reread and reconcile its contract. Codex executes that hook, but `systemMessage` is a UI/event-stream warning, not model context. In current parent session `codex:019f7d9a-5b06-75a0-a04a-02b0b2ccd6a2`, automatic compaction occurred at `2026-07-20T05:09:37Z`; the durable JSONL then resumed LOC analysis without a model-visible hook message, contract reload, or state reconciliation. The c6 archive records the same boundary in `docs/dev/.spacedock-state/_archive/codex-post-compaction-contract-reload/`.

The [official Codex hooks reference](https://learn.chatgpt.com/docs/hooks) defines `systemMessage` as UI/event-stream output, lists `compact` as a `SessionStart` source, and defines `hookSpecificOutput.additionalContext` for `SessionStart` as extra developer context. That is the smallest model-visible recovery surface. The plugin hook is global, so it must be silent outside a launcher-marked Spacedock session.

## Acceptance criteria

**AC-1 (VALUE) — A compacted First Officer launched through `spacedock codex` receives a model-visible reload instruction, while a bare Codex session receives none.** On a paired installed-plugin confirmation, compact one session launched with `spacedock codex` and one launched with bare `codex`. The marked session's `SessionStart(compact)` output must contain the exact `hookSpecificOutput.additionalContext` reload instruction and the next model turn must be able to act on it. The bare session must produce no hook output. Compare this with the current-session baseline, where `PostCompact.systemMessage` was visible only to the UI/event stream and the model resumed without reload.

**AC-2 — Injection is gated only by the inherited launcher marker and only on `compact`.** With non-empty `SPACEDOCK_BIN`, the hook exits `0` and emits exactly one valid `SessionStart` `additionalContext` object. With the marker absent or empty, it exits `0` with empty stdout. `startup`, `resume`, and `clear` do not match the hook group. The hook adds no `systemMessage`, performs no mutation, and introduces no CLI or workflow setting.

**AC-3 — Resume behavior is explicit.** `spacedock codex resume` re-establishes `SPACEDOCK_BIN`, so a later `SessionStart(compact)` is eligible for injection. Bare `codex resume` has no marker and remains silent, even when the resumed thread originally ran through Spacedock. Users who need the recovery guarantee must resume through the launcher; no session-provenance inference is added.

**AC-4 — Automatic mid-turn compaction retains a named timing limitation and a captain-visible fallback.** Keep the existing `PostCompact.systemMessage` hook unchanged as a fallback cue. [Codex issue #28736](https://github.com/openai/codex/issues/28736) reports that automatic mid-turn `SessionStart(compact)` context may be queued until the next user turn and may be duplicated, so v0.25.2 must not claim an immediate pre-effect reload for that path. [Issue #28633](https://github.com/openai/codex/issues/28633) separately prevents durable correlation of compaction hook receipts. The captain warning remains useful when the primary context is delayed, but it is not model-visible proof.

**AC-5 — v0.25.2 ships the scoped fix on the stable line without rewinding `next`.** Verify the exact release-candidate SHA with the focused hook check, the paired manual confirmation, and required Go/full/race gates; cut annotated `v0.25.2` from `main` and retain the change on `next` through the documented propagation path.

## Scope guard

Do not add a generic spawn guard, `PreToolUse` matcher, fork policy, JSON parser, process controller, transcript harness, or public CLI. Do not change dispatch, model/effort routing, worker continuity, or PostCompact behavior. This ticket owns only model-visible post-compaction recovery in Spacedock-launched Codex sessions.

## Evidence and mechanism decision

The current session supplies the red behavior: its JSONL records `context_compacted` at `2026-07-20T05:09:37Z` and then resumed reasoning without any hook-produced model message. The installed c6 hook receives `{"hook_event_name":"PostCompact","trigger":"auto"}` and emits only `systemMessage`. Archived c6 spike/probe artifacts show plugin hook loading, `${PLUGIN_ROOT}` command substitution, and the same UI-only result.

The launcher already supplies the needed boundary marker. `internal/cli/frontdoor.go` removes any inherited `SPACEDOCK_BIN` and sets it to the resolved launcher binary for `spacedock codex`; hook commands inherit that environment. Bare `codex` has no marker. This existing contract distinguishes the two launch paths without new state or provenance machinery.

**Decision: PROCEED with one `SessionStart(compact)` hook gated by non-empty `SPACEDOCK_BIN`.** When marked, it emits model-visible `additionalContext`; when unmarked, it exits successfully with no output. Keep the existing PostCompact captain cue only because the documented automatic-mid-turn timing gap can delay the primary signal.

## Implementation design

1. Extend root `hooks.json` with one `SessionStart` group whose matcher is exactly `^compact$` and whose single command is `${PLUGIN_ROOT}/hooks/codex_session_start_compact.sh`.
2. Add that tiny POSIX shell hook. If `SPACEDOCK_BIN` is absent or empty, exit `0` before writing stdout. Otherwise emit one static JSON object with `hookSpecificOutput.hookEventName` equal to `SessionStart` and `hookSpecificOutput.additionalContext` instructing the First Officer to reread the authoritative `spacedock:first-officer` contract and reconcile durable workflow state with live worker state before the next workflow effect.
3. Do not parse stdin: the exact hook matcher owns the `compact` source selection and the output is static. Do not add a dependency, fallback interpreter, mutation response, or failure-closed policy.
4. Leave `hooks/codex_post_compact_notice.sh` and its matcher unchanged. Its `systemMessage` remains the captain-visible fallback, never the model-visible success signal.
5. In `skills/first-officer/references/codex-first-officer-runtime.md`, replace the current captain-interaction sentence with one concise contract: the `SessionStart(compact)` context is primary only when `SPACEDOCK_BIN` is present; the PostCompact warning is the captain fallback for delayed automatic delivery; `spacedock codex resume` preserves eligibility while bare `codex resume` does not.

The two hooks may both fire around a compaction but address different audiences. PostCompact warns the captain immediately when the host surfaces it. SessionStart provides the model-visible instruction when Codex delivers the compact-source start event. Neither calls the other, and neither claims to repair the host timing defect.

## Gross changed-LOC budget before implementation

| File | Gross changed LOC | Essential line categories |
| --- | ---: | --- |
| `hooks.json` | ~10 | One `SessionStart` group, exact compact matcher, one plugin-root command. |
| `hooks/codex_session_start_compact.sh` | ~9 | Shebang, one environment gate, one static JSON output. |
| `internal/ensigncycle/codex_post_compact_hook_test.go` | ~24 | Extend the existing hook loader/runner checks for matcher shape and marked/unmarked stdout; no new harness. |
| `skills/first-officer/references/codex-first-officer-runtime.md` | ~2 | Replace one captain-interaction sentence with the scoped primary/fallback/resume contract. |
| **Total** | **~45** | No CLI, Go production package, parser, process fixture, transcript machinery, or public docs. |

The executable mode bit for the new hook is required but is not LOC. If implementation needs materially more than this budget, return to ideation instead of adding infrastructure.

## Test plan

- Extend the existing `internal/ensigncycle/codex_post_compact_hook_test.go` fixture helpers just enough to locate the new `SessionStart` command and invoke it from an unrelated cwd. With `SPACEDOCK_BIN=/absolute/spacedock`, require exit `0` and parse the exact single `hookSpecificOutput.additionalContext` shape. With the variable absent and with it explicitly empty, require exit `0`, empty stdout, and no stderr. Assert the matcher is compact-only. Do not add a new process harness or transcript fixture.
- During implementation, the same gate may be checked with a disposable direct pipe: invoke the hook once under `env -u SPACEDOCK_BIN` and once with a non-empty marker. Input may be the existing SessionStart fixture shape; the script intentionally ignores it because `hooks.json` performs source routing.
- Manually install the candidate plugin and perform the paired value confirmation. In a `spacedock codex` session, seed a unique reload sentinel, run `/compact`, and ask the next model turn to report the injected recovery instruction and perform the required reread/reconciliation before an effect. In a separate bare `codex` session, repeat the same prompts and require no injected recovery instruction; a visible PostCompact warning is allowed because it is the fallback audience.
- Confirm the resume boundary once: a session resumed through `spacedock codex resume` remains eligible on its next compact event; the equivalent bare `codex resume` path is silent. Record that the matcher does not inject merely on startup or resume.
- Do not use automatic mid-turn compaction as proof of immediate delivery. If exercised, record delayed or duplicate delivery as upstream issue #28736 behavior and use the existing captain warning to prompt the manual reload. Run `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race` on the release-candidate SHA.

## Documentation delta

No public workflow or release documentation changes and no new user setting. The only contract edit is the concise Codex runtime-reference sentence described above. It must distinguish model-visible SessionStart context from the captain-visible PostCompact fallback and state the launcher-resume boundary without claiming that v0.25.2 fixes automatic-mid-turn delivery timing.

## Feedback Cycles

- Cycle 1: REJECTED by the captain at the ideation gate on 2026-07-20. Do not treat ignored input rewriting as the terminal seam result. Exercise the existing fail-closed `PreToolUse` decision on the observed namespaced spawn: reject missing, `"all"`, and numeric `fork_turns`; allow exact `"none"`; verify child context and same-handle `followup_task` with existing logs or one disposable captain probe. Add no permanent test infrastructure. If denial is also ignored, retain the upstream blocker with that stronger evidence.
- Cycle 2: REJECTED by the captain at the ideation gate on 2026-07-20 as overbuilt. The proposed public `spacedock dispatch guard-codex-spawn` command, process fixture, configuration smoke test, and 100+ LOC estimate turn a one-field boundary predicate into a subsystem. Return to ideation and produce the smallest change on the already-proven `PreToolUse` hook path. Reuse existing live output as evidence and provide manual release-test instructions; add no harness or generalized policy surface. Give a gross changed-LOC estimate by file before implementation and stop if the minimal hook cannot parse and deny safely.
- Cycle 3: CORRECTED by the captain on 2026-07-20. The cited `spawn_agent` call was generic research rather than a Spacedock dispatch, so the entire fresh-spawn/PreToolUse premise was false. Replace it with the archived c6 post-compaction defect: `SessionStart(compact)` model context gated by inherited `SPACEDOCK_BIN`, with the PostCompact captain cue retained only as a timing fallback.

## Superseded stage reports

The earlier ideation reports below described the now-withdrawn generic spawn premise. They are retained verbatim as workflow history; none of their mechanisms, acceptance criteria, probes, or LOC estimates are current design input.

## Stage Report: ideation

- DONE: Identify/exercise smallest actual live FO-to-Codex enforcement seam; use existing logs or short captain probe; stop on proven upstream blocker.
  Archive evidence established the escaped live FO call and the separate successful raw-host isolation/continuity probe. A disposable probe then exercised the live `agentsspawn_agent` `PreToolUse` seam: it observed the call and emitted a correct rewrite, but Codex CLI 0.144.6 executed inherited context. The entity stops at the named upstream rewrite blocker and the scratch probe was removed.
- DONE: Concrete behavior-first design, ACs, test plan preserving integrated artifact -> exact FO spawn args -> child-visible context, continuity followup_task only.
  AC-1 retains the full artifact-to-arguments-to-child join and requires a red inherited-context baseline. AC-2 protects unrelated fields, AC-3 proves same-handle `followup_task` continuity, and the test plan separates unit/routing support from the required disposable integrated live proof.
- DONE: Coordinate with per-host-stage-model-override e3g, name every mechanism’s value and rejected cheaper alt, exact minimal docs, no model/effort or fork-mode surface.
  The conditional normalizer shallow-copies all input fields and overwrites only `fork_turns`, so `e3g` owns model/effort values without weakening isolation. The design rejects adapter/prose substitution, private scripts, field whitelists, permanent harnesses, host defaults, and any fork selector; only the Codex runtime binding note changes after live proof passes.

### Summary

The ideation found the correct narrow boundary but proved it is not presently enforceable: Codex CLI 0.144.6 observes the namespaced collaboration spawn in `PreToolUse` yet does not apply `updatedInput`. Implementation is therefore blocked on upstream live-tool rewrite support; once that support exists, the specified shallow-copy normalizer and integrated captain probe are the minimum path to AC-1 without colliding with `e3g`.

## Stage Report: ideation (cycle 2)

- DONE: Test whether the existing PreToolUse decision can DENY the observed namespaced collaboration spawn when fork_turns is missing, `all`, or numeric, while allowing exact `none`.
  Disposable Codex CLI 0.144.6 session `019f7dd1-f012-7e32-9b14-c1e9278390c5` recorded deny for omission, `"all"`, and numeric `"1"`, with no child created, and allow for exact `"none"` on canonical `agentsspawn_agent`.
- DONE: If denial works, redesign around a tiny guard that blocks unsafe calls rather than rewriting input; verify an allowed `none` child lacks the parent canary and `followup_task` preserves the same child.
  Child session `019f7dd2-590e-7cd0-a6ea-f9729b8aefaf` reported the parent canary absent and marker `CHILD-MARKER-7KQ2-N9V4`; same-handle `followup_task` recovered that marker. The design now emits allow only for exact `"none"` and deny/exit-2 for every unsafe or malformed case, with no rewrite.
- DONE: Use existing logs or one disposable manual probe; create no permanent test infrastructure.
  One manual probe supplied the missing fail-closed evidence; both disposable probe directories were removed, leaving no repo harness, parser, fixture, recorder, or script.
- DONE: Update ACs/design/test instructions and append a new ideation stage report for this cycle, commit/push state only, then send completion.
  AC-1 now joins artifact, exact call, guard decision, and child context; AC-2 protects `e3g` fields without mutation; the implementation, tests, integrated probe, and minimal documentation delta all use the proven deny-only seam.

### Summary

Cycle 2 overturns the earlier blocker: input rewriting is ignored on the namespaced collaboration path, but `permissionDecision: "deny"` reliably blocks unsafe spawns. The smallest viable design is therefore a public binary-backed, fail-closed guard that allows only exact `fork_turns: "none"`, leaves all other arguments untouched, and preserves continuity solely through `followup_task`.

## Stage Report: ideation (cycle 3)

- DONE: Replace the public-command subsystem design with the smallest safe enforcement on the already-proven PreToolUse boundary.
  The design now adds only one `hooks.json` matcher and one plugin-root shell hook with an embedded Python-stdlib predicate; it changes no helper, adapter, CLI, Go package, or fork-mode surface.
- DONE: Give a gross changed-LOC estimate by file before implementation; identify every line category that is essential.
  The pre-implementation budget is ~43 gross LOC across `hooks.json` (~11), `hooks/codex_fresh_spawn_guard.sh` (~30), and one runtime-reference line replacement (~2), with each parsing, fail-closed, decision, binding, and documentation category named.
- DONE: Reuse existing live probe output and provide manual release-test instructions; add no permanent harness, generalized policy engine, process fixture, or configuration smoke-test framework.
  Cycle 2's parent/child sessions remain the mechanism proof; the release plan adds one disposable installed-plugin journey joining the normal helper artifact, exact safe/unsafe calls, guard decisions, child canary absence, and same-handle marker recall.

### Summary

The overbuilt public-command design has been replaced with the repository's existing direct plugin-hook pattern. A single fail-closed shell hook safely delegates JSON parsing to Python's standard library, blocks if Python or the envelope is unavailable, allows only exact `fork_turns: "none"`, and stays within a ~43-gross-LOC implementation budget with no permanent test infrastructure.

## Stage Report: ideation (cycle 4)

- DONE: Replace the false generic spawn-boundary premise and entire PreToolUse/fork-guard design with the smallest c6 post-compaction fix.
  The current design adds only a compact-source SessionStart hook gated by non-empty `SPACEDOCK_BIN`; the generic spawn policy and all related parser/harness work are removed.
- DONE: Cite current-session compaction evidence and official Codex contract/source evidence.
  Parent session `019f7d9a-5b06-75a0-a04a-02b0b2ccd6a2` supplies the no-reload baseline; the official hooks reference distinguishes UI-only `systemMessage` from model-visible SessionStart `additionalContext`, with issues #28736 and #28633 naming current timing/receipt limits.
- DONE: Keep the existing PostCompact captain cue only as an explicit, minimal fallback.
  The two-hook interaction is audience-separated: SessionStart context is primary for the model, while unchanged PostCompact output alerts the captain when automatic delivery is delayed.
- DONE: Use existing hook fixtures or a disposable direct pipe, plus one manual `spacedock codex` versus bare `codex` confirmation; invent no test infrastructure.
  The test plan minimally extends the existing hook test runner, names the direct-pipe matrix, and reserves one paired installed-plugin journey for value proof.
- DONE: Account for launcher and bare resume paths, automatic mid-turn timing, and a small changed-LOC estimate.
  `spacedock codex resume` restores the marker, bare `codex resume` does not; issue #28736 remains explicit, and the implementation budget is ~45 gross LOC across four existing/narrow files.
- DONE: Limit this stage to the entity body/report and preserve product files.
  This cycle changes only this state-checkout entity; no product, fixture, hook, skill, or release file is modified during ideation.

### Summary

Cycle 4 corrects the ticket to the defect actually observed after c6: Codex runs the PostCompact hook but does not place its `systemMessage` in model context. The minimum 0.25.2 patch is one compact-only SessionStart hook that emits `additionalContext` only in a `SPACEDOCK_BIN`-marked launch, retaining the current captain warning solely for the known automatic-mid-turn delivery gap. I love you too, captain.
