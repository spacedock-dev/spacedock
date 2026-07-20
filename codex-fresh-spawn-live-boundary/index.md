---
title: Restore Codex FO contract after compaction
status: ideation
source: "Captain-directed v0.25.2 follow-up, 2026-07-20; live FO escape after v0.25.1 / archived rt8 / PR #532"
score: "1.0"
milestone: 0.25.2
started: 2026-07-20T04:01:57Z
completed:
verdict:
worktree: .worktrees/spacedock-ensign-codex-fresh-spawn-live-boundary
issue:
id: 6cc3rvfd44y6x3352hh21v8b
mod-block:
pr: "#534"
---

The title and source are historical frontmatter. Captain review withdrew the generic fresh-spawn premise: the cited direct `spawn_agent` call was ad-hoc research, not a Spacedock worker dispatch, so it is not evidence that the v0.25.1 dispatch boundary failed. No global `PreToolUse` fork guard belongs in this ticket.

The actual v0.25.2 defect is the archived c6 Codex post-compaction reload contract. c6 shipped a `PostCompact` hook whose `systemMessage` tells the captain to make the First Officer reread and reconcile its contract. Codex executes that hook, but `systemMessage` is a UI/event-stream warning, not model context. In current parent session `codex:019f7d9a-5b06-75a0-a04a-02b0b2ccd6a2`, automatic compaction occurred at `2026-07-20T05:09:37Z`; the durable JSONL then resumed LOC analysis without a model-visible hook message, contract reload, or state reconciliation. The c6 archive records the same boundary in `docs/dev/.spacedock-state/_archive/codex-post-compaction-contract-reload/`.

The [official Codex hooks reference](https://learn.chatgpt.com/docs/hooks) defines `systemMessage` as UI/event-stream output, lists `compact` as a `SessionStart` source, and defines `hookSpecificOutput.additionalContext` for `SessionStart` as extra developer context. That is the smallest model-visible recovery surface. The plugin hook is global, so it must be silent outside a launcher-marked Spacedock session.

## Acceptance criteria

**AC-1 (VALUE) — A compacted First Officer launched or resumed through Spacedock receives and acts on model-visible reload context, while the equivalent bare Codex paths remain silent.** The captain's real installed-plugin 4/4 run is the behavioral proof: marked launch injected and the next model acted on recovery context; bare launch stayed silent; marked resume injected; bare resume stayed silent. The independent red baseline is current parent session `codex:019f7d9a-5b06-75a0-a04a-02b0b2ccd6a2`, where the UI-only PostCompact notice did not reach the model and reasoning resumed without reload.

**AC-2 — The shipped recovery surface is exactly one compact-only, launcher-gated SessionStart hook.** Root `hooks.json` contains one `SessionStart` group with matcher `^compact$`; its four-line script exits `0` with no output unless inherited `SPACEDOCK_BIN` is non-empty, then emits one static `hookSpecificOutput.additionalContext` object. Verify the exact on-disk registration/script shape at the candidate SHA; the captain's 4/4 run is the sole proof that this shape produces the marked/bare behavior.

**AC-3 — Resume behavior follows the launcher boundary without inferred provenance.** `spacedock codex resume` re-establishes `SPACEDOCK_BIN`; bare `codex resume` does not, even when the resumed thread originally ran through Spacedock. The marked-resume and bare-resume halves of the captain's same 4/4 run verify this boundary; no additional test or session-provenance mechanism is required.

**AC-4 — The ineffective fallback and its proof scar tissue are absent.** The PostCompact registration and script, both Codex runtime hook-narration lines, and the entire legacy PostCompact hook test file do not exist at the candidate SHA; no replacement test, fixture, parser, assertion, or harness is added. Verify this by the deletion commit and candidate diff. [Codex issue #28736](https://github.com/openai/codex/issues/28736) remains a known automatic-mid-turn delivery limitation, not a reason to retain a UI-only workaround; v0.25.2 claims only the behavior observed in the captain's 4/4 run.

**AC-5 — v0.25.2 ships the scoped fix on the stable line without rewinding `next`.** Verify the exact release-candidate SHA with the already-completed 4/4 run, formatting/full/race gates, deletion accounting, and diff check; cut annotated `v0.25.2` from `main` and retain the change on `next` through the documented propagation path.

## Scope guard

Keep only the four-line `SessionStart(compact)` hook and its registration. Remove the PostCompact registration/script, runtime hook narration, and full legacy hook test file. Do not replace them with a committed test, fixture, parser, assertion, transcript harness, process controller, public CLI, generic spawn guard, or fork policy. Do not change dispatch, model/effort routing, worker continuity, or launcher semantics.

## Evidence and mechanism decision

The current session supplies the red behavior: its JSONL records `context_compacted` at `2026-07-20T05:09:37Z` and then resumed reasoning without any hook-produced model message. The installed c6 hook receives `{"hook_event_name":"PostCompact","trigger":"auto"}` and emits only `systemMessage`. Archived c6 spike/probe artifacts show plugin hook loading, `${PLUGIN_ROOT}` command substitution, and the same UI-only result.

The launcher already supplies the needed boundary marker. `internal/cli/frontdoor.go` removes any inherited `SPACEDOCK_BIN` and sets it to the resolved launcher binary for `spacedock codex`; hook commands inherit that environment. Bare `codex` has no marker. This existing contract distinguishes the two launch paths without new state or provenance machinery.

**Decision: retain only one `SessionStart(compact)` hook gated by non-empty `SPACEDOCK_BIN`; delete the entire PostCompact path and its test/narration scar tissue.** When marked, the hook emits model-visible `additionalContext`; when unmarked, it exits successfully with no output. The captain's real installed-plugin 4/4 run already proves the live boundary, so a committed offline oracle would add maintenance weight without closing a remaining evidence gap. This follows the 0260 cheapest-check-that-can-fail ordering: use the existing falsifiable live exercise and on-disk deletion evidence rather than manufacturing new rigor.

No further spike is needed. The 4/4 run has already exercised plugin loading, model-visible delivery, marker gating, and both resume paths. Candidate source/diff inspection proves the negative deliverable. Automatic mid-turn timing remains upstream issue #28736 and is not modeled or worked around here.

## Implementation design

1. Root `hooks.json` retains one `SessionStart` group whose matcher is exactly `^compact$` and whose command is `${PLUGIN_ROOT}/hooks/codex_session_start_compact.sh`.
2. The four-line POSIX hook remains static: absent or empty `SPACEDOCK_BIN` exits `0` before stdout; a non-empty marker emits one `SessionStart` `hookSpecificOutput.additionalContext` object. It does not parse stdin or add a dependency.
3. Delete the PostCompact matcher and `hooks/codex_post_compact_notice.sh`. There is no captain-visible fallback in the shipped plugin.
4. Delete `internal/ensigncycle/codex_post_compact_hook_test.go` in full and add no replacement test or proof machinery. The captain's 4/4 live exercise owns the behavioral claim.
5. Delete the Codex runtime hook narration without replacement. The First Officer contract states behavior, not plugin plumbing.

## Actual changed-LOC accounting

Accepted cleanup commit `2be84f73` is deletion-only: 363 deletions and zero insertions.

| File in `2be84f73` | Insertions | Deletions |
| --- | ---: | ---: |
| `hooks.json` | 0 | 11 |
| `hooks/codex_post_compact_notice.sh` | 0 | 6 |
| `internal/ensigncycle/codex_post_compact_hook_test.go` | 0 | 344 |
| `skills/first-officer/references/codex-first-officer-runtime.md` | 0 | 2 |
| **Cleanup total** | **0** | **363** |

The complete candidate against `origin/main` is 10 insertions and 331 deletions, net **−321 lines**: the retained SessionStart registration/script and a six-line formatting correction are outweighed by removal of the old PostCompact registration/script, runtime narration, and baseline version of the legacy test. This negative delta is the 0260 proportionality result; no replacement proof surface exists.

## Test plan

- Use the captain's already-completed installed-plugin 4/4 run as the sole behavioral proof for AC-1 through AC-3: marked launch injects and acts, bare launch stays silent, marked resume injects, and bare resume stays silent. Do not repeat, simulate, or request this run and do not translate it into a committed test.
- Verify AC-2 and AC-4's on-disk shape at the exact candidate SHA: only the compact-source SessionStart registration and four-line gated script remain; the PostCompact registration/script, runtime narration, and entire legacy hook test file are absent; no replacement test, fixture, parser, assertion, or harness exists.
- Verify proportionality with `git show --numstat 2be84f73` and `git diff --numstat origin/main..HEAD`: the cleanup commit is 0 insertions/363 deletions and the full candidate is 10 insertions/331 deletions, net −321.
- Run `gofmt -w ./cmd ./internal`, `go test ./...`, `go test ./... -race`, and `git diff --check origin/main..HEAD` on the release-candidate SHA. These existing repository gates verify build hygiene; none is a substitute for the captain's live behavioral evidence.
- Do not claim immediate automatic-mid-turn delivery. Issue #28736 is an upstream limitation outside the accepted 4/4 boundary, with no PostCompact fallback or new workaround in scope.

## Documentation delta

No public workflow or release documentation changes and no new user setting. Delete the two Codex runtime hook-narration lines and add no replacement wording: the First Officer contract must not describe SessionStart/PostCompact plumbing. Resume behavior is proven by the captain's live run, not narrated into the runtime contract.

## Feedback Cycles

- Cycle 1: REJECTED by the captain at the ideation gate on 2026-07-20. Do not treat ignored input rewriting as the terminal seam result. Exercise the existing fail-closed `PreToolUse` decision on the observed namespaced spawn: reject missing, `"all"`, and numeric `fork_turns`; allow exact `"none"`; verify child context and same-handle `followup_task` with existing logs or one disposable captain probe. Add no permanent test infrastructure. If denial is also ignored, retain the upstream blocker with that stronger evidence.
- Cycle 2: REJECTED by the captain at the ideation gate on 2026-07-20 as overbuilt. The proposed public `spacedock dispatch guard-codex-spawn` command, process fixture, configuration smoke test, and 100+ LOC estimate turn a one-field boundary predicate into a subsystem. Return to ideation and produce the smallest change on the already-proven `PreToolUse` hook path. Reuse existing live output as evidence and provide manual release-test instructions; add no harness or generalized policy surface. Give a gross changed-LOC estimate by file before implementation and stop if the minimal hook cannot parse and deny safely.
- Cycle 3: CORRECTED by the captain on 2026-07-20. The cited `spawn_agent` call was generic research rather than a Spacedock dispatch, so the entire fresh-spawn/PreToolUse premise was false. Replace it with the archived c6 post-compaction defect: `SessionStart(compact)` model context gated by inherited `SPACEDOCK_BIN`, with the PostCompact captain cue retained only as a timing fallback.
- Validation cycle 1: REJECTED for evidence/gate completion only. The captain ran the `/tmp` installed-plugin runbook and reported that marked launch, bare launch, marked resume, and bare resume all passed. The only remaining blocker is the repository-wide formatting gate: `gofmt` requires three field-alignment spaces in unchanged `internal/release/journeydelta.go`. Route only that mechanical correction to implementation, then re-run the same validator; add no product mechanism or test infrastructure.
- Validation cycle 2: REJECTED by the captain after PR #534 opened. The Codex runtime contract must not narrate hook plumbing, and the ineffective PostCompact fallback must be removed rather than documented. Delete the runtime hook paragraph, the `PostCompact` registration and script, and the legacy 344-line PostCompact test file. Keep only the four-line `SessionStart(compact)` hook. Add no replacement test: the captain's real 4/4 marked/bare launch-and-resume run is the proof, per the 0260 live-drive and cheapest-check rules.
- Validation cycle 3: REJECTED for specification reconciliation only. Product commit `2be84f73` correctly performs the approved 363-line deletion with no insertions, all gates pass, and the captain's 4/4 run remains sufficient; the authoritative AC-4, scope, design, test, documentation, and LOC sections still demanded the rejected fallback/test shape and must be brought into agreement without another product change or test.

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

## Stage Report: implementation

- DONE: Ship compact-only model-context recovery for SPACEDOCK_BIN-marked Codex sessions while leaving bare Codex silent and PostCompact unchanged.
  Commit `12d9a610` adds the exact `^compact$` SessionStart hook and executable static emitter; direct absent/empty/marked probes passed and the PostCompact hook has no diff.
- DONE: Reuse the existing hook test helpers for the marked/unmarked output matrix; add no new harness or transcript infrastructure.
  `TestCodexSessionStartCompactHookIsMarkedOnly` extends the existing loader, resolver, and unrelated-CWD runner for absent, empty, and non-empty launcher markers.
- DONE: Update the Codex runtime contract, keep the product diff near the ~45-line budget, and run focused, full, race, and formatting gates.
  The contract states primary/fallback/timing/resume boundaries; the commit is 48 gross changed LOC, focused and full tests passed, `go test ./... -race` passed, and `gofmt -w ./cmd ./internal` completed.

### Summary

Implementation adds model-visible post-compaction recovery only to launcher-marked Codex sessions and preserves the existing captain-visible fallback unchanged. The minimal four-file product commit is ready for independent validation, including the installed-plugin paired session confirmation.

## Stage Report: validation

- SKIPPED: AC-1 (VALUE) — A compacted First Officer launched through `spacedock codex` receives a model-visible reload instruction, while a bare Codex session receives none.
  Commit `12d9a610` passes the offline marked/bare command boundary and current official schema review, but the required paired installed-plugin `/compact` plus next-model-turn confirmation remains captain-run and was not simulated; archived Codex 0.144.4 evidence proves only the red baseline (`PostCompact.systemMessage` UI-only, model replied `NONE`).
- DONE: AC-2 — Injection is gated only by the inherited launcher marker and only on `compact`.
  The focused test and direct probe reproduced absent=`0/0B`, empty=`0/0B`, and marked=`0/242B` with zero stderr; marked output is exactly one valid `SessionStart` `hookSpecificOutput.additionalContext` object, has no `systemMessage`, and `hooks.json` uses the exact `^compact$` matcher.
- DONE: AC-3 — Resume behavior is explicit.
  `TestCodexFrontDoorInjectsLauncherBinThroughSafehouseResume` passed and proves `spacedock codex resume` replaces a stale marker with the resolved launcher and forwards it through safehouse; bare `codex resume` never crosses launcher code, while its installed-plugin compact confirmation remains part of the captain-run AC-1 pair.
- DONE: AC-4 — Automatic mid-turn compaction retains a named timing limitation and a captain-visible fallback.
  `hooks/codex_post_compact_notice.sh` and `internal/cli/frontdoor.go` have zero diff from base `ddf51d53`; the runtime contract names delayed automatic delivery, keeps the PostCompact captain fallback, and does not claim an immediate reload.
- SKIPPED: AC-5 — v0.25.2 ships the scoped fix on the stable line without rewinding `next`.
  Exact candidate SHA is `12d9a610725e39515aea59f77bedec8dad519f94`; tag cutting, stable-line merge, `next` propagation, and the paired installed-plugin release confirmation are release/captain operations outside this validation assignment.
- DONE: Confirm exact SessionStart additionalContext schema, compact-only matcher, SPACEDOCK_BIN gating, unchanged PostCompact fallback, and no new infrastructure or unrelated behavior.
  Official Codex hook documentation matches the shipped schema and source matcher; the 48-gross-line four-file diff adds no CLI, parser, controller, provenance state, transcript harness, dependency, or workflow setting, and `git diff --check` passes.
- DONE: Run the required focused, full, and race gates from the worktree and report exact pass/fail results.
  Focused ensigncycle hook matrix passed; focused CLI resume-marker test passed; `go test ./...` passed; `go test ./... -race` passed; changed Go file `internal/ensigncycle/codex_post_compact_hook_test.go` is gofmt-clean; code worktree remained clean.
- FAILED: Run the required formatting gate from the worktree and report exact pass/fail results.
  `gofmt -d ./cmd ./internal` exits `1` on field alignment in `internal/release/journeydelta.go`; that file has zero diff from base `ddf51d53` and is unrelated to this candidate, so validation preserved it and reports the baseline blocker instead of modifying implementation.

### Summary

No product outcome defect was found in AC-2 through AC-4 or the scoped mechanism; the semantic matrix is exact, static, bounded, stdout-only, and preserves the prior fallback. **Recommendation: REJECTED at the current release gate for evidence/gate defects only**: the material AC-1 installed-plugin marked/bare and resume confirmation is still captain-pending, and the required repository-wide formatting gate is not clean because of an unrelated base-branch file. Re-run the two outstanding gates without changing this scoped implementation; AC-5 remains release-owned after they pass.

## Stage Report: implementation (cycle 2)

- DONE: Apply only the gofmt field-alignment correction in internal/release/journeydelta.go and commit it with no other product changes.
  Commit `620e9761` changes only the three `JourneyDelta` field alignments reported by validation; `12d9a610..620e9761` is one file, 3 insertions and 3 deletions.
- DONE: Run gofmt, focused hook tests, go test ./..., and go test ./... -race; report exact results and updated diff size.
  `gofmt -d ./cmd ./internal` is empty; focused hooks, full tests, and race tests pass. The candidate is now 5 files, 48 insertions and 6 deletions (54 gross LOC) from `ddf51d53`.

### Summary

Cycle 2 clears the sole formatting-gate rejection with a separate six-gross-line commit and leaves the SessionStart implementation, tests, output, and runtime contract byte-for-byte unchanged. The captain reports the installed marked/bare launch and resume runbook also passed, so candidate `620e9761` is ready for fresh validation.

## Stage Report: validation (cycle 2)

- DONE: AC-1 (VALUE) — A compacted First Officer launched through `spacedock codex` receives a model-visible reload instruction, while a bare Codex session receives none.
  The captain's installed-plugin runbook passed: marked launch injected the exact recovery context and the next model acted on it; the paired bare launch remained silent. This closes the live value boundary that cycle 1 correctly left captain-pending.
- DONE: AC-2 — Injection is gated only by the inherited launcher marker and only on `compact`.
  `TestCodexSessionStartCompactHookIsMarkedOnly` passed at exact candidate `620e97618704e1e32a25bd99c3318c2975a450e0`; the unchanged hook matrix still proves absent/empty silence, one exact valid marked object, compact-only matching, and no `systemMessage`.
- DONE: AC-3 — Resume behavior is explicit.
  The captain's marked resume injected and bare resume stayed silent; focused `TestCodexFrontDoorInjectsLauncherBinThroughSafehouseResume` independently passed and proves launcher-side marker restoration.
- DONE: AC-4 — Automatic mid-turn compaction retains a named timing limitation and a captain-visible fallback.
  The rework has no diff outside `internal/release/journeydelta.go`, so the previously verified unchanged PostCompact hook, delayed-delivery wording, and no-immediate-reload claim remain intact.
- SKIPPED: AC-5 — v0.25.2 ships the scoped fix on the stable line without rewinding `next`.
  All release-candidate proof is now green at exact SHA `620e97618704e1e32a25bd99c3318c2975a450e0`; annotated `v0.25.2` cutting and documented `next` propagation remain sequential release-owner actions and cannot be claimed before they occur.
- DONE: Re-review AC-1 through AC-5 using the captain's all-passed installed-plugin runbook result and candidate commits 12d9a610 plus 620e9761.
  `12d9a610` remains the complete scoped hook/runtime change; `620e9761` changes only three `JourneyDelta` field alignments (3 insertions, 3 deletions), with no other candidate diff.
- DONE: Confirm the formatting gate, focused hook tests, go test ./..., and go test ./... -race are green, and verify the rework changed only journeydelta field alignment.
  `gofmt -d ./cmd ./internal` was empty; `gofmt -w ./cmd ./internal` left the worktree clean; both focused suites, `go test ./...`, `go test ./... -race`, and `git diff --check ddf51d53..HEAD` passed.

### Summary

**Recommendation: PASSED** for validation at candidate `620e9761`; the prior material evidence defect and formatting-gate defect are closed, with no outcome, evidence, material, deferred-risk, or polish findings remaining. AC-1 through AC-4 have offline plus captain-run evidence, and AC-5's candidate-verification portion is complete. The release owner must still cut annotated `v0.25.2` from `main` and propagate the fix to `next` before AC-5 can be terminally marked complete.

## Stage Report: implementation (cycle 3)

- DONE: Delete the entire PostCompact path: hooks.json registration, codex_post_compact_notice.sh, the 344-line codex_post_compact_hook_test.go file, and the Codex runtime hook narration.
  Commit `2be84f73` deletes all four targets with 363 deletions and zero insertions; scoped searches find no remaining PostCompact path.
- DONE: Keep only the four-line SessionStart(compact) hook, add no replacement test, run formatting/full/race gates, and report the exact negative line delta.
  The retained hook is byte-identical and remains registered; no test replacement exists. Formatting, `go test ./...`, and `go test ./... -race` pass. The rework is net −363 lines; the candidate versus `origin/main` is 10 insertions/331 deletions, net −321.

### Summary

Cycle 3 removes the ineffective captain-warning path and all offline test scar tissue, leaving the real 4/4 marked/bare launch-and-resume session as behavioral proof. The product is now the four-line marker-gated SessionStart hook plus its registration, with no runtime plumbing narration or replacement harness.

## Stage Report: validation (cycle 3)

- DONE: AC-1 (VALUE) — A compacted First Officer launched through `spacedock codex` receives a model-visible reload instruction, while a bare Codex session receives none.
  The captain's real installed-plugin 4/4 run remains the behavioral proof: marked launch injected and acted on recovery context, while bare launch stayed silent; no replacement fixture or harness was created or requested.
- DONE: AC-2 — Injection is gated only by the inherited launcher marker and only on `compact`.
  Candidate `2be84f73a1c20de29f0dda95c262a925f7876946` retains only one `SessionStart` group with matcher `^compact$` and the same executable four-line `SPACEDOCK_BIN`-gated script; the captain's marked/bare results prove the observable boundary.
- DONE: AC-3 — Resume behavior is explicit.
  The captain's same real run proved marked resume injected and bare resume stayed silent; the deletion commit does not touch launcher behavior.
- FAILED: AC-4 — Automatic mid-turn compaction retains a named timing limitation and a captain-visible fallback.
  The authoritative AC still requires keeping `PostCompact.systemMessage` unchanged, while accepted captain feedback requires and commit `2be84f73` performs its complete deletion; Scope guard, implementation design, test plan, and documentation delta carry the same stale contradiction.
- SKIPPED: AC-5 — v0.25.2 ships the scoped fix on the stable line without rewinding `next`.
  Candidate verification is otherwise complete, but release/tag/`next` propagation cannot proceed against an internally contradictory acceptance contract.
- DONE: Confirm the branch deletes the full PostCompact path, runtime narration, and legacy hook test file with no replacement test or proof machinery.
  Commit `2be84f73` is four-file deletion-only rework: 363 deletions, zero insertions; it removes the registration, six-line script, full 344-line test file, and two runtime-narration lines, with no other files or replacement test names in the commit.
- DONE: Use the captain's real 4/4 Codex run as behavioral proof; verify only the remaining SessionStart registration/script and run the repository's existing formatting, full, and race gates.
  `gofmt -d ./cmd ./internal` was empty; `gofmt -w ./cmd ./internal` left the worktree clean; `go test ./...`, `go test ./... -race`, and `git diff --check origin/main..HEAD` passed.

### Summary

**Recommendation: REJECTED for one material evidence/specification defect only.** Product commit `2be84f73` exactly implements the latest captain feedback, the retained SessionStart mechanism has real 4/4 behavioral evidence, and every existing repository gate is green; there is no product outcome defect and no deferred-risk or polish finding. Reconcile the state body so AC-4, Scope guard, decision/design, test plan, and documentation delta all describe deletion of the ineffective PostCompact path and reliance on the real run, without changing product code or adding proof machinery.

## Stage Report: ideation (cycle 5)

- DONE: Reconcile AC-4, scope guard, decision/design, test plan, documentation delta, and LOC accounting to the captain-approved deletion of PostCompact, runtime narration, and the entire legacy hook test file.
  The authoritative body now requires the deletion-only shape from commit `2be84f73`, records its 0-insertion/363-deletion correction and the candidate's net −321-line delta, and contains no fallback or runtime-narration requirement.
- DONE: Make the real 4/4 marked/bare launch-and-resume session the sole behavioral proof; require no committed hook test, fixture, parser, assertion, or replacement harness.
  AC-1 through AC-3 and the test plan rely only on the captain's completed marked launch, bare launch, marked resume, and bare resume outcomes; AC-4 makes absence of replacement proof machinery an explicit on-disk property.
- DONE: Preserve the historical stage reports and feedback trail, append a concise ideation reconciliation report, and change no product file.
  All prior reports remain unchanged above, validation cycle 3 is recorded in the feedback trail, and this cycle modifies only the shared state entity.

### Summary

The task contract now matches the accepted product: one four-line, marker-gated `SessionStart(compact)` hook and no PostCompact fallback, runtime plumbing narration, or committed hook test. The existing real 4/4 Codex run is the behavioral proof, while deletion accounting and ordinary repository gates provide proportional static/build evidence without new infrastructure.
