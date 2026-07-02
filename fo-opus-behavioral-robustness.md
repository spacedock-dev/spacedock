---
id: 5xst4dae1m2acx76czpe857v
title: FO behavioral-robustness — live opus-4-8 FO intermittently deviates from contract guidance
status: validation
source: '0221 cut (2026-06-19): the live claude-live OPUS lane flaked on essentially every PR with DIFFERENT non-deterministic FO-behavior deviations, while SONNET stayed compliant. Diagnosed diff-unrelated across 3e #398 (sonnet broad-search) and rgq #399 (two reruns: no spacedock-new, then merged-team-mode exit-before-terminalize) and accepted non-blocking under the ci-env conn to land 0.22.1. This is the recorded root-cause follow-up.'
started: 2026-07-02T09:55:07Z
completed:
verdict:
score: 0.6
worktree: .worktrees/spacedock-ensign-fo-opus-behavioral-robustness
issue:
group: cleanup
---

The live opus-4-8 first officer intermittently fails to follow contract guidance precisely — surfaced by the live-e2e lanes doing their job — making the opus claude-live lane unreliable on every PR. This is the strong-model echo of the weak-model concern the layered-FO sprint targets: even opus deviates, so the deferred tier/delegation mechanism (72) and report-and-stop hardening matter.

## Problem

Across the 0.22.1 cut the opus claude-live lane failed on three DISTINCT non-deterministic FO behaviors, none related to the PR under test, while sonnet/codex/pi passed:

1. **No `spacedock new`** — the FO hand-filed a seed instead of the atomic-create path (`claude_live_runner_test.go:122`).
2. **Broad-search after zero-discover** — ran `find <root> -type f` instead of report-and-stop (`zero_discover_live_test.go`, the `tq0` sibling).
3. **Exit before terminalize** — the FO subprocess exited (code 0) before driving the entity to terminal in merged-team-mode (`merged_team_mode_live_test.go:176`).

The lane won't converge by re-running (a different deviation each run), so it blocks every PR's merge.

## Proposed direction (ideation to refine)

Harden the FO contract's behavioral robustness so a strong-model FO reliably complies — and/or right-size any over-strict live assertions where a capable FO legitimately varies. Distinguish, per failing scenario, "contract guidance too weak" vs "test too strict." Relates to `tq0` (broad-search) and to the deferred `72` (tier/delegation forces verdict/judgment compliance structurally).

## Out of scope

- The deferred `72` tier mechanism itself (separate rebuild on the cleaned foundation).

## Ideation — evidence-grounded determinations (2026-07-02)

Evidence read from the failing runs' own artifacts, never from failure labels: GH Actions run 27841539602 (PR #399, attempts 1 and 2 — opus artifacts `claude-stream.jsonl` / `merged-stream.jsonl`, job logs 82401555749 / 82405494944), run 27835552853 (PR #398, attempt 1 — sonnet job log 82382674493), and green main run 28432388663 (opus artifact) for the post-remediation check.

### Record correction

The Problem section above misattributes deviation 2: the broad-search-after-zero-discover failure in the 0.22.1 cut was the **sonnet** leg (PR #398 run 27835552853 attempt 1), not opus. Opus passed `TestLiveZeroDiscoverReportsAndStops` on every attempt in the cut. The class does flip between sonnet/opus legs on other runs (`tq0`'s byte-verified record), so it remains a both-models class — but this cut recorded no opus instance of it. The two opus deviations in the cut were the filing red (attempt 1) and the merged-team-mode red (attempt 2).

### Scenario 1 — "no `spacedock new`" (opus, run 27841539602 attempt 1): live assertion too strict, already fixed

The artifact shows the FO filing exactly the contract-blessed way: `cd {root}`, `B=${SPACEDOCK_BIN:-spacedock}`, `$B new wire-the-thing --workflow-dir . <<'EOF' …` — atomic create, launcher var-capture, no `--next-id`, no `Write`. The then-current `newInvocation` regex required a literal launcher token and the `new` verb on one line, so the var-capture split made correct behavior red. This is the named pattern (unscoped token matching over free model output false-positiving on correct behavior). PR #433 (2026-06-22) added `capturedLauncherFilesViaNew` covering this shape. Spike A (below) replays the archived stream through the current assertion: PASS. Residual: the checked-in fixture is synthetic (no `cd` prefix, no heredoc); the real stream is not pinned.

### Scenario 2 — broad-search after zero-discover (sonnet in this cut): genuine deviation class, owned by siblings

PR #398 attempt 1: after a zero `status --discover`, the sonnet FO ran `find {root} -type f | head -30` — a genuinely banned shape (find over the project root hunting a workflow), not the `4t8` flat-`ls` over-strictness. Determination: contract-side genuine deviation, model-stochastic; the prose lever is spent (#374's finding). Hardening is owned by `tq0`; detector/contract alignment (the flat-`ls` false red) is owned by `4t8`. This task contributes the record correction only — no duplicate deliverable.

### Scenario 3 — "exit before terminalize" (opus, run 27841539602 attempt 2): genuine deviation, mislabeled; structurally remediated by #400; one binary hole remains

The label is wrong: the FO did not walk off the job. The merged-stream shows a complete ceremony — implementation dispatched and reported, then `--set make-it-work completed verdict=PASSED worktree=` (tool result: completed/verdict/worktree stamped), `--archive` (tool result: archived), teardown, reconcile, clean summary, exit 0. The miss is that the terminalize `--set` **omitted `status=done`**; the entity was archived with `status: implementation`, so the test's terminal condition (`^status: done$`) never became true. The test graded the durable end-state correctly — this is a genuine FO deviation (incomplete hand-rolled finalize), not an over-strict assertion.

Structural remediation landed 75 minutes after the failing run: PR #400's `spacedock merge guard <slug>` finalize writes status+verdict+completed in ONE `--set` (merge.go terminalize list includes `status={terminal}`), and `fo-merge-core.md` routes the FO through `«merge.guard»` with an explicit warning that a stale binary silently falls back to the hand ceremony. Live opus exercise of the remediation: green main run 28432388663's merged-team-mode transcript shows the FO reading fo-merge-core and declaring the merge-guard drive, and the entity reached `status: done` on disk.

Residual hole, spiked open today (Spike C): the current binary still accepts the exact failing sequence. `--set {slug} completed verdict=PASSED worktree=` on a non-terminal entity succeeds, and `--archive` then exits 0, leaving `status: implementation` + `verdict: PASSED` in `_archive/`. `runSet` already gates the finalize action (`completed` without a verdict is refused) and `runArchive` mirrors that verdict gate "so finalization-by-archive cannot route around it" — but neither gate covers finalize-without-terminal-status. A drifting or stale-binary FO can still reproduce the deviation end-state with no refusal.

### Post-remediation opus compliance record (the independent baseline)

Across the 21 completed `runtime-live-e2e` runs from 2026-06-22 (post-#433) through 2026-07-02: the opus leg failed once (run 28497354019 — shallow-boot pre-greet deferred-skill invocation, a different class), and the three recorded classes recurred **zero** times on opus. Cut baseline: 3 deviations across 3 opus lane executions (2 genuine + 1 assertion false-red).

## Proposed approach

Scope this task to the two residuals its own evidence owns; everything else is determination-on-the-record.

1. **Finalize-status gate.** Extend the existing finalize-gate pattern in `internal/status/handlers.go` (runSet) and its archive mirror in `mutate.go` (runArchive): the finalize action (`completed` being set) is refused when the resulting `status` is not a declared terminal stage (neither already terminal nor set terminal in the same call); `--archive` refuses an entity carrying a non-`rejected` verdict while `status` is non-terminal. `--force` bypasses both, `verdict: rejected` keeps reject-then-archive on the happy path, bare dispatch-into-terminal (`status=done started`, no `completed`) stays ungated, and `merge guard`'s atomic finalize never trips either gate. Error text follows the #465 pattern — remediation in the binary's own output, e.g.: `Error: entity {slug} cannot be finalized ('completed') while status '{status}' is not the terminal stage. Set status={terminal} in the same --set (or run 'spacedock merge guard {slug}'), or use --force.`
2. **Real-stream regression fixture.** Check in the archived PR #399 attempt-1 opus filing stream (102,607 bytes) as a must-pass fixture for `assertClaudeFilingViaNew` — pinning the real deviation shape (cd-prefix + var-capture + heredoc), per the `mq`/`4t8` practice of preserving captured streams as fixtures.
3. **Measured opus confirmation** (AC-1) after the gate lands.

Non-goals: broad-search hardening (`tq0`/`4t8`), tier delegation (`72`), the shallow-boot deferred-skill class (new observation, worth its own seed — noted for the FO).

## Acceptance criteria

- **AC-1 (MEASURES the end-value):** Across 5 fresh live opus repetitions of the three affected scenarios (filing, merged-team-mode, zero-discover), the deviation count for the two opus classes (filing-path false-red, incomplete-finalize) is 0, against the 0.22.1-cut baseline of 3 deviations across 3 opus lane executions — and any zero-discover red in those reps is triaged per `4t8`'s detector-vs-contract semantics, not auto-counted. The count can move the wrong way (any recurrence is >0 and fails the AC).
  Test: 5 × `SPACEDOCK_LIVE_MODEL=claude-opus-4-8 go test -tags live ./internal/ensigncycle -run 'TestLiveClaudeSharedScenarios/filing|TestLiveMergedTeamModeDispatch|TestLiveZeroDiscoverReportsAndStops'` with run artifacts attached to the task. Cost: measured opus filing run is $0.38/55s (journey metrics); merged-team-mode ≈ $1–2/3–6 min; zero-discover ≈ $0.40/35s → ≈ $2.5–3 and ~8 min per rep, ≈ $15 and ≤45 min wall for all 5.
- **AC-2:** The archived PR #399 attempt-1 opus filing stream is a checked-in fixture that `assertClaudeFilingViaNew` accepts.
  Test: an offline test in `internal/ensigncycle` reads the fixture from testdata and asserts nil error (it fails on the pre-#433 regex — verified by Spike A's replay mechanics); `go test ./internal/ensigncycle -run Filing`.
- **AC-3:** The incomplete-finalize shape is refused by the binary: `--set {slug} completed verdict=PASSED worktree=` on a non-terminal entity exits non-zero with the remediation text, `--archive` refuses a non-`rejected`-verdict + non-terminal entity, and the sanctioned paths stay green — `merge guard` finalize, reject-then-archive, verdict-less non-terminal archive, bare dispatch-into-terminal.
  Test: unit table in `internal/status` whose first red case is Spike C's exact command sequence; existing merge-guard and golden-fixture suites stay green (`go test ./internal/status ./internal/cli`).

## Test plan

Offline first: the AC-3 guard table (minutes, no model spend) and the AC-2 fixture replay (seconds). One live opus merged-team-mode + filing run after the gate lands to confirm the FO's sanctioned path never trips the new refusal (~$2, ~6 min), then the AC-1 5-rep measurement (≈ $15). No new fixture/CLI machinery is needed — all three tests and both assertion suites already exist.

## Spike record (riskiest mechanisms exercised first)

- **Spike A (S1 determination):** replayed the archived attempt-1 filing stream through the current `assertClaudeFilingViaNew` via a throwaway in-package test — PASS ("current assertion ACCEPTS the archived PR #399 attempt-1 opus filing stream"). The #433 fix covers the real deviation, not just its synthetic fixture.
- **Spike B (S3 remediation is live):** green main run 28432388663's opus merged-team-mode transcript shows fo-merge-core loaded and the FO declaring the `merge guard` drive; the entity reached `status: done` (the lane passed).
- **Spike C (the proposed gate's target exists):** replayed the attempt-2 command sequence against a binary built from current main — `--set … completed verdict=PASSED worktree=` then `--archive` both exit 0 and leave `status: implementation` + `verdict: PASSED` archived. This script is the implementation's first red test.

## Documentation impact

None on the docs site: no user-facing doc enumerates the mutation guards (architecture-notes' generic "guards that refuse an unsafe mutation" sentence already covers the new gate); `fo-merge-core.md` already routes the FO through `«merge.guard»` and needs no change. The new error text is specified verbatim in the Proposed approach.

## Stage Report: ideation

- DONE: Each of the three recorded opus deviations (hand-filed seed instead of `spacedock new`, broad-search after zero-discover, exit-before-terminalize in merged-team-mode) gets an evidenced per-scenario determination — "contract guidance too weak" vs "live assertion too strict" — grounded in the actual failing-run artifacts and test code, not inherited labels.
  S1: assertion too strict, fixed by #433 (archived stream shows correct var-capture filing; replay passes now). S2: genuine class, but the cut's instance was SONNET (record corrected), owned by tq0/4t8. S3: genuine deviation mislabeled "exit-before-terminalize" — terminalize `--set` omitted `status=done`; remediated by #400, residual binary hole spiked open. All from run 27841539602 attempts 1/2 and run 27835552853 attempt 1 artifacts.
- DONE: At least one AC MEASURES opus-FO behavioral compliance against an independent baseline that can move the wrong way (e.g. deviation count across N live opus runs of the affected scenarios), with a test plan and cost estimate — not just "contract prose updated".
  AC-1: deviation count across 5 live opus reps of the three scenarios vs the cut baseline (3 deviations / 3 executions); ≈$2.5–3 and ~8 min per rep, ≈$15 total, grounded in journey-metrics ($0.38 measured opus filing run).
- DONE: The riskiest unverified mechanism is exercised first and recorded in the body (or "no spike needed" with the proven mechanisms named) — a candidate contract/assertion change should have at least one live opus exercise behind it before the gate.
  Spike A (archived stream replays green through current assertion), Spike B (green opus run 28432388663 drives merge guard live), Spike C (current binary accepts the incomplete-finalize sequence — the proposed gate's first red test). Live opus exercises behind both candidate changes are the archived attempt-1 stream (AC-2) and run 28432388663 (AC-3's sanctioned path).

### Summary

Diagnosed all three 0.22.1-cut deviations from the runs' own artifacts: one assertion false-red already fixed (#433), one sonnet-not-opus misattribution corrected and deferred to tq0/4t8, and one genuine opus finalize omission structurally remediated by #400 with a remaining binary hole (archive/finalize accepted without terminal status) proven open by replay. Scoped the task to a finalize-status gate + a real-stream regression fixture + a measured 5-rep live opus compliance check; post-remediation baseline recorded (0 recurrences across 21 lane runs since 2026-06-22).

## Stage Report: implementation

- DONE: AC-3 red-first — Spike C's exact command sequence (`--set {slug} completed verdict=PASSED worktree=` on a non-terminal entity, then `--archive`) is the first FAILING test in the internal/status guard table, turned green by the finalize/archive gates, sanctioned paths stay green
  `TestFinalizeStatusGateRefusesSpikeCSequence` confirmed RED pre-fix (exit 0, entity mutated), GREEN after adding the gate; `internal/status/handlers.go` (runSet) and `internal/status/mutate.go` (runArchive) each gained a finalize-status gate; 8 new tests in `internal/status/finalize_status_guard_test.go` cover merge-guard atomic finalize (existing merge_guard_test.go suite unaffected), reject-then-archive, verdict-less non-terminal archive, bare dispatch-into-terminal, and `--force` bypass. Commit 66750ea9.
- DONE: AC-2 with falsifiability — the real archived PR #399 attempt-1 opus filing stream (102,607 bytes) is checked-in testdata that `assertClaudeFilingViaNew` accepts, and the test demonstrably fails against the pre-#433 regex logic
  Downloaded run 27841539602's `runtime-live-e2e-claude-live-claude-opus-4-8` artifact (attempt-1, id 7756226361) via `gh api`; `claude-shared-scenarios/filing/claude-stream.jsonl` is exactly 102,607 bytes and its filing command is the real cd-prefix + var-capture + heredoc shape (`cd {root}\nB=${SPACEDOCK_BIN:-spacedock}\n$B new wire-the-thing --workflow-dir . <<'EOF'…`). Checked in as `internal/ensigncycle/testdata/claude_live_filing_pr399_attempt1.stream.jsonl`; new permanent test `TestAssertClaudeFilingViaNewAcceptsRealPR399Attempt1Stream` asserts nil error. Falsifiability verified by temporarily swapping in the pre-#433 `shared_filing_test.go` (`git show 6ea56f95~1`) — the fixture test went RED (`the FO did not file the seed via a spacedock … new … command`) — then restored (clean diff confirmed via `git status`).
- DONE: The refusal error text carries its remediation in binary output per the #465 pattern (the verbatim text specified in the task body's Proposed approach), and full `go test ./...` is green
  `--set` refusal: `entity {slug} cannot be finalized ('completed') while status '{status}' is not the terminal stage. Set status={terminal} in the same --set (or run 'spacedock merge guard {slug}'), or use --force.` — matches the Proposed approach text verbatim. `go test ./...` green across all 14 packages (`go vet ./...` and `gofmt -l` clean on all changed files).

### Summary

Added a finalize-status gate to `runSet` and its `--archive` mirror in `runArchive`, closing the residual hole Spike C proved open: a finalize (`completed`+verdict) that never advances `status` to a declared terminal stage in the same `--set` call is now refused, and `--archive` independently refuses a non-`rejected`-verdict entity sitting at a non-terminal status (the stale-binary/hand-edit defense-in-depth case). Both gates key off `--force` and were ordered to preserve every existing golden/merge-guard test unchanged (the archive gate sits after the merge-hook invariant so `060-passed-nosentinel`'s existing refusal reason is undisturbed). Also checked in the real archived PR #399 attempt-1 opus filing stream as a must-pass regression fixture for `assertClaudeFilingViaNew`, with falsifiability empirically verified against the pre-#433 logic. `go test ./...` is green; AC-1 (the 5-rep live opus measurement) is intentionally left for validation/a follow-on live run per the stage checklist, which scoped implementation to AC-2/AC-3 plus the verbatim error text.
