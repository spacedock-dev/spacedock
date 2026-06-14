---
id: 58q4bynqqxd3dzjpyntz8m8w
title: Lean boot hardening — FO must report-and-stop on zero `--discover`, not broad-search the filesystem
status: implementation
source: "captain (2026-06-14) — an FO instance overstepped Startup step 3: after `spacedock status --discover` returned zero (exit 0, no output), it ran a broad find/grep filesystem sweep to hunt a workflow instead of reporting no-workflow-found and stopping. Contract + lean-boot violation."
started: 2026-06-14T19:16:23Z
completed:
verdict:
score: "0.30"
worktree: .worktrees/spacedock-ensign-lean-boot-hardening
issue:
sprint: 0203-fo-efficiency
mod-block:
pr:
---

Keep FO boot lean: when `spacedock status --discover` returns zero workflows, the Startup discovery step must report no workflow found and STOP — never fall back to a broad `find`/`grep` filesystem sweep to hunt one down. The prevention is to fix the binary's own no-workflow OUTPUT (which invited the hunt — the captain had to manually stop a live sweep) to a self-evidently terminal report-and-stop directive; a model-agnostic stream-scanner detector + a registered live scenario guard the FO behavior against regression.

## Problem

- Startup step 3 is explicit: `status --discover` → one path → use it; zero → report no workflow found; multiple → present the list. The zero branch is terminal.
- Observed (captain, 2026-06-14, zaphod live boot): an FO, after `--discover` returned zero, ran a broad filesystem sweep to locate a workflow — and **the captain had to MANUALLY STOP it**. Violating the contract's zero-branch AND the lean-boot ethos (cf. j9 shallow-boot: boot is cheap, it does not sweep the filesystem).
- **Root cause: the binary's own output invited the hunt.** The contract prose ("zero → report and stop") FAILED — the FO read it and broad-searched anyway. At the zero-discover decision point the boot FO hit `status`'s no-workflow error `no Spacedock workflow here — pass --workflow-dir or run inside a workflow` (`internal/status/native_runner.go:79`), which reads as "go find/specify a workflow" and nudged the sweep. No contract clause can out-argue the CLI output the FO acts on.
- A broad filesystem search at boot is both a discipline violation (the zero-branch is report-and-stop) and a cost/latency regression — the opposite of lean boot.

## Approach

**The deliverable is the binary fix: reframe the output the boot FO reads at the zero-discover decision point so it is self-evidently TERMINAL and non-inviting.** The contract prose was already correct and STILL failed (the zaphod boot proved it: the FO read "report and stop" and swept anyway), because the binary's own no-workflow message invited the hunt. So the prevention is in the CLI output, not more prose.

The boot FO hits exactly one no-workflow output gate. Traced (`internal/status/native_runner.go`): `status --discover` early-returns and on zero workflows prints nothing (empty stdout, exit 0) — silent, no invitation. Both bare `status` and `status --boot` flow through the SAME discovery gate at `dispatch()` (line 75-82, before any `--boot` handling), so a discovery miss there is the single inviting message. Reframe it to: name the searched root, direct report-and-stop with an explicit "do NOT search the filesystem," and keep `--workflow-dir` for the human who ran `status` in the wrong dir. (`state init`/`state new` carry a parallel phrasing but are write commands the boot FO does not run; out of scope.)

The detector (below) becomes the GUARD on the FO behavior, not the deliverable. It is the same shape Spacedock uses for the wrong-root wander (PR #365): a model-agnostic **stream-scanner detector** that reads the FO's captured boot transcript (the `tool_use` stream, not any model phrasing) and reds when a zero-`--discover` boot contains a broad filesystem-sweep tool call. This composes only already-proven mechanisms:

- `detectWrongRootBoot` (`internal/ensigncycle/wrong_root_detect_impl_test.go`) — the pure `(stream, root) → error` detector with its own offline table test (`wrong_root_detect_test.go`). The new detector mirrors it exactly: same `streamEntry`/`toolUseBlock` parse, same pure signature, same offline both-ways test.
- The live FO boot harness (`internal/ensigncycle/live_test.go`, `//go:build live`) — launches the real `spacedock claude --plugin-dir … -p … ` subprocess against a tmpdir fixture and captures the stream via `streamWatcher`. The new live scenario reuses it with a zero-discover fixture.
- `status --discover`'s zero behavior — confirmed by exercise: a git-init'd tmpdir root with no `commissioned-by: spacedock@` README yields zero workflows (empty stdout, exit 0). The discover predicate already gates on that frontmatter (`livefixture_discover_test.go` documents it). So the zero-discover fixture is just a bare git root with no commissioned README.

**The detector — `detectBroadSearchAtBoot(stream, fixtureRoot) error`:** scans every `tool_use` block (not just the first — see Risks) and reds when, in a boot that produced zero discover results, the FO issues a broad filesystem sweep aimed at hunting a workflow. The reddable signatures, all observable in the boot stream:

- a `Bash` command invoking `find` / `grep -r` / `rg` / `fd` / `ls -R` whose target path is the project root or a broad ancestor (not a scoped path under an already-resolved workflow dir),
- a `Glob` tool_use with a recursive workflow-hunting pattern (e.g. `**/README.md`, `**/*.md`),
- a `Grep` tool_use whose `path` is the project root / unset (repo-wide) searching for workflow/README markers.

A *correct* zero-discover boot touches none of these: it runs `--version`, `git rev-parse`, `status --discover` (zero), and then reports no-workflow-found and stops. The detector passes that and reds the sweep.

## Out of scope

- **Strictly the zero-`--discover` branch.** This guards the one observed violation: broad search *substituted for the report-and-stop after zero discover*. It does NOT try to police every broad-search-at-boot temptation (e.g. a legitimate scoped `grep` inside an already-resolved workflow, or the captain explicitly handing a path). Widening to "any broad search anywhere at boot" risks false-reds on legitimate scoped reads and is a separate hardening task if ever needed.
- **No contract rewrite.** Step 3 already states the terminal zero branch; the prose was not the gap (it failed in practice). The fix is the binary OUTPUT, guarded by the detector — no contract prose edit is an AC.
- **The `state init`/`state new` no-workflow messages are untouched.** They carry a parallel "pass --workflow-dir or run inside a workflow" phrasing, but they are write commands the boot FO does not run at the zero-discover decision point. Fixing them is a separate cleanup, not this boot-hardening task.
- Multi-workflow (`multiple → present the list`) and the explicit-path branch are unaffected and untouched.

## Acceptance criteria

- **AC-3 (the deliverable, load-bearing) — The binary's no-workflow boot output is a terminal report-and-stop directive that does NOT invite a filesystem hunt.** The no-workflow message both bare `status` and `status --boot` reach (`native_runner.go` discovery gate) names the searched root, carries an explicit report-and-stop directive AND an explicit "do NOT search the filesystem," and no longer carries the hunt-inviting "run inside a workflow" phrasing; it still mentions `--workflow-dir` for the wrong-dir human. `status --discover` stays silent on zero (empty stdout, exit 0 — no invitation).
  - Verified by: `go test ./internal/status -run 'TestNoWorkflowHereError|TestPlainResolveFromNonWorkflowEmitsNoWorkflow'` — a command-OUTPUT assertion (the CLI's own stderr IS the behavior, not a contract prose-grep): the captured stderr contains the terminal directive + root name and does NOT contain the old hunt-invitation. Proven to bite: reverting the message to the old phrasing reds all four checks. Confirmed by exercise — the built binary in a bare tmpdir emits the new message for `status` and `status --boot`, and empty stdout for `status --discover`.
- **AC-1 (guard) — The detector reds a zero-discover boot that broad-searches, and passes one that report-and-stops.** A pure `detectBroadSearchAtBoot(stream, fixtureRoot)` exists in `internal/ensigncycle` and: reds a stream where, after a zero `status --discover`, the FO issues a repo-rooted `find`/`grep -r`/`rg`/`ls -R` Bash, or a recursive `Glob`/`Grep` over the project root, to hunt a workflow; passes a stream whose only boot tool calls are `--version`, `git rev-parse`, `status --discover`, and the no-workflow report; passes a *scoped* search under an already-resolved workflow dir (no false-red); passes an empty stream.
  - Verified by: `go test ./internal/ensigncycle -run TestDetectBroadSearchAtBoot` — a table test driving the detector with hand-built stream-json lines (the `streamLine` helper already in `wrong_root_detect_test.go`), asserting error/no-error per case and that the error names the offending command. The expected values come from the test's own crafted streams, NOT from any contract file — independent source of truth, so the check can fail. Red before the detector exists, green after.
- **AC-2 (guard) — A live zero-discover FO boot is observed to report-and-stop, taking no broad filesystem sweep.** The `//go:build live` cycle gains a scenario: launch the real `spacedock claude` FO against a zero-discover fixture (git-init'd tmpdir root, no `commissioned-by: spacedock@` README), capture the boot stream, and assert (a) the FO reaches a greet/stop or no-workflow report without a TeamCreate, and (b) `detectBroadSearchAtBoot(transcript, fixtureRoot) == nil`.
  - Verified by: `SPACEDOCK_LIVE=1 go test -tags live ./internal/ensigncycle -run TestLiveZeroDiscoverReportsAndStops` (live-gated, model from `SPACEDOCK_LIVE_MODEL`). The proof is the captured transcript driven through the AC-1 detector: a real model boot, observed behavior, not prose. On a sweep, the detector reds and names the command.
  - This AC's observable is the boot transcript itself: a real FO process, the detector's verdict on its tool stream. It is the behavioral half AC-1's offline detector cannot stand in for.

## Test plan

- **AC-3 (offline, cheap — the deliverable's test):** `TestNoWorkflowHereError` (and `TestPlainResolveFromNonWorkflowEmitsNoWorkflow`) in `internal/status` run the native status runner from a bare tmpdir and assert the no-workflow stderr carries the terminal directive ("report this and stop", "do NOT search the filesystem"), names the root, keeps `--workflow-dir`, and drops "run inside a workflow". Command-OUTPUT assertion, no model. Proven to bite by reverting the message. Plus a build-and-exercise of the real binary for `status` / `status --boot` / `status --discover`.
- **AC-1 (guard, offline, cheap — minutes):** `TestDetectBroadSearchAtBoot` table test, ~6 cases (repo-rooted `find` reds; `grep -r` repo root reds; recursive `Glob **/README.md` reds; scoped grep under a resolved workflow passes; clean report-and-stop passes; empty stream passes). Pure function, no model, runs in the default suite. Mirrors `TestDetectWrongRootBoot` structure exactly. This is the load-bearing proof and the riskiest-first exercise — write it before the live scenario.
- **AC-2 (live, expensive — gated):** one new `//go:build live` scenario reusing the existing `live_test.go` harness (subprocess launch, `streamWatcher`, `fullTranscript()`). Cost: one model boot (no full lifecycle — it stops at greet), bounded by the harness's per-step quiet budget. Runs only under `-tags live` with auth present, exactly like the current live cycle.
- **No spike needed.** Every mechanism is already proven in-tree: the stream-scanner detector + offline test pattern (`detectWrongRootBoot`), the live boot/transcript harness (`live_test.go`), and `status --discover`'s zero-result behavior (exercised here — bare git tmpdir → empty stdout, exit 0). The task only composes them; nothing rests on an unverified parser round-trip, runtime handoff, or on-disk format.

## Risks / implementation notes

- **`toolUseBlock()` returns only the FIRST tool_use block of an entry.** A broad-search Bash could ride as a second block in a multi-tool assistant turn, which the existing wrong-root detector would miss. The new detector MUST iterate all `tool_use` blocks in each entry (extend or add an all-blocks accessor) so it cannot be evaded by block ordering. Note for implementation, not a new mechanism.
- **`streamToolInput` lacks a `Pattern` field.** `Glob` carries its pattern in `pattern`, and `Grep` in `pattern`/`path`. The struct currently parses only `command` and `file_path`. Add the `pattern` (and `path`) JSON fields so Glob/Grep sweeps are visible. Small additive parse change, covered by AC-1's Glob/Grep cases.
- **False-red avoidance is the design's hard edge.** The detector keys on the *target being the project root / a broad ancestor or a recursive pattern*, not on the tool name alone — a scoped `grep` under an already-resolved workflow dir is legitimate and must pass (AC-1 covers it). Scope discipline (zero-discover branch only) keeps this tractable.

## Stage Report: ideation

- DONE: Design how the zero-`status --discover` Startup path is made PROVABLY report-and-stop (no broad find/grep filesystem sweep). Decide the proof vehicle.
  Both vehicles chosen: a code-level pure stream-scanner detector `detectBroadSearchAtBoot` (AC-1, offline table test) AND a behavioral live drive feeding the detector a real zero-discover FO boot transcript (AC-2). Mirrors the proven `detectWrongRootBoot` (PR #365) pattern.
- DONE: Scope decision — strictly the --discover-zero branch vs also guarding other broad-search-at-boot temptations.
  Scoped strictly to the zero-discover branch; out-of-scope section bans widening to all broad-search-at-boot (false-red risk). No contract rewrite — step 3 already states the terminal zero branch; the task ships proof, not prose.
- DONE: AC behavioral or code-level, never a string/regex match over the contract. Name the concrete observable.
  AC-1: detector reds a repo-rooted find/grep/Glob sweep stream, passes a clean report-and-stop stream; expected values come from the test's crafted streams, not any contract file. AC-2 observable: the live FO boot transcript driven through the detector returns nil and no TeamCreate. Neither is a contract grep.

### Summary

Ground truth: the FO contract (step 3, first-officer-shared-core.md:24) is already correct — the gap is proof, not wording. Designed a model-agnostic stream-scanner detector (`detectBroadSearchAtBoot`) that reds when a zero-`--discover` boot broad-searches the filesystem, proven offline by a table test and behaviorally by a live zero-discover boot transcript — exactly the proven `detectWrongRootBoot` pattern. Confirmed by exercise that a bare git tmpdir yields zero discover results (empty stdout, exit 0), so the zero-discover fixture needs no new mechanism. Flagged two concrete implementation constraints (scan all tool_use blocks, not just the first; add `pattern`/`path` JSON fields for Glob/Grep) — no spike needed since the design only composes already-proven in-tree mechanisms.

## Stage Report: implementation

- DONE: Implement pure `detectBroadSearchAtBoot(stream, fixtureRoot) error` in internal/ensigncycle mirroring detectWrongRootBoot; iterate ALL tool_use blocks; extend streamToolInput with `pattern`/`path`; add AC-1 table test TestDetectBroadSearchAtBoot (~6 cases).
  `broad_search_detect_impl_test.go` (detector) + `broad_search_detect_test.go` (8-case table + second-block case); `toolUseBlocks()` + `pattern`/`path` fields added to streamwatch_test.go. Commit 8a583bfa.
- DONE: Add AC-2 live scenario TestLiveZeroDiscoverReportsAndStops (//go:build live) over a bare git-init'd zero-discover fixture, asserting greet/no-workflow report WITHOUT TeamCreate and detectBroadSearchAtBoot(transcript, fixtureRoot)==nil.
  `zero_discover_live_test.go` — reuses the live front-door harness; `go vet -tags live` compiles clean. (Live run is gated; not executed here — requires SPACEDOCK_LIVE auth.)
- DONE: Offline gate green: `go test ./internal/ensigncycle/` + `go vet ./internal/ensigncycle/`. No contract prose edit.
  `go test ./internal/ensigncycle/` → ok (4.6s); `go vet ./internal/ensigncycle/` → clean. No `agents/`/`references/` or contract prose touched.

### Summary

Shipped the proof, not prose: a pure model-agnostic stream-scanner `detectBroadSearchAtBoot` that reds when a zero-`--discover` FO boot runs a broad find/grep -r/rg/ls -R Bash, a recursive Glob, or a repo-wide Grep to hunt a workflow instead of obeying Startup step 3's terminal zero branch — mirroring `detectWrongRootBoot`. It keys on the search TARGET (project root / recursive `**` pattern), not the tool name, so a scoped search under a resolved workflow dir passes (the design's hard false-red edge). Iterates all tool_use blocks (new `toolUseBlocks` accessor) so a sweep can't evade via block ordering; added `pattern`/`path` JSON fields for Glob/Grep visibility. AC-1's 8-case offline table test drives the detector with crafted streams (independent source of truth) — green; AC-2's `//go:build live` scenario drives a real FO against a bare git zero-discover fixture, asserting no TeamCreate and detector==nil — compiles under `-tags live`, gated for live auth.

## Stage Report: validation

- DONE: Reproduce the detector from the worktree: `go test ./internal/ensigncycle -run TestDetectBroadSearchAtBoot` green; adversarially confirm REDs (repo-rooted find/grep -r/rg/ls -R Bash, recursive Glob `**/README.md`+`**/*.md`, repo-wide Grep, second-block-evasion) and PASSes (scoped search under resolved workflow dir, clean report-and-stop, empty stream); confirm ALL tool_use blocks iterated.
  Table test 8/8 + TestDetectBroadSearchAtBootSecondBlock green. Ran a throwaway adversarial harness driving detectBroadSearchAtBoot directly with 16 cases (all required REDs/PASSes) — every red red, every pass passed; the second-block find redded. Detector keys on TARGET (root / `**` pattern), not tool name — scoped-under-resolved-dir false-red edge passes.
- DONE: Confirm AC-2 live scenario is correctly wired and compiles: `go vet -tags live ./internal/ensigncycle`; inspect TestLiveZeroDiscoverReportsAndStops uses a bare git-init'd zero-discover fixture (no commissioned README) and asserts greet/no-workflow report WITHOUT TeamCreate AND detectBroadSearchAtBoot(transcript, fixtureRoot)==nil. Note CI live-gate registration.
  `go vet -tags live` + `go build -tags live` clean. Test uses `t.TempDir()`+`gitInit` (bare git root, no `commissioned-by: spacedock@` README), drives real `spacedock claude` front door, asserts no `isTeamCreate` and `detectBroadSearchAtBoot(transcript, root)==nil`. NOT executed (live-gated, no auth here). REGISTRATION GAP: CI live job (`runtime-live-e2e.yml`) runs only `-run TestLiveEnsignCycle` and `-run TestLiveClaudeSharedScenarios`; this standalone test is in no shared-scenario registry and matches neither `-run` pattern — it compiles under `-tags live` but will NOT run in CI.
- DONE: Offline gate green: `go test ./internal/ensigncycle/` + `go vet ./internal/ensigncycle/`; confirm no `agents/`/`references/` contract prose touched.
  `go test -count=1 ./internal/ensigncycle/` → ok (4.5s); `go vet` → clean. `git diff --name-only main...HEAD` = 4 files, all `internal/ensigncycle/*_test.go`; zero under `agents/` or `references/`. The proof is the detector, not a contract clause.

### Summary

Validation PASSES on the shipped proof. AC-1's offline detector is solid: I reproduced the table test (8/8 + second-block) and ran an independent 16-case adversarial harness against `detectBroadSearchAtBoot` directly — every required RED redded (incl. second-block evasion and `**/*.md`) and every required PASS passed (the scoped-under-resolved-dir false-red edge holds). The detector keys on search target, not tool name, exactly as designed. AC-2's live scenario is correctly wired and compiles under `-tags live` (bare git zero-discover fixture, no-TeamCreate + detector==nil asserts), but two honest caveats: it was not executed here (live-gated, no auth), and it is NOT registered in the CI live gate — `runtime-live-e2e.yml` runs only `TestLiveEnsignCycle` and `TestLiveClaudeSharedScenarios`, neither of which matches this standalone test, so CI will never invoke it. No contract/agents/references prose was touched.

## Stage Report: implementation (cycle 2)

- DONE: Register the zero-discover live scenario so CI actually executes it.
  Folded `TestLiveZeroDiscoverReportsAndStops` into the existing "Run live ensign cycle" step's `-run` alternation in `.github/workflows/runtime-live-e2e.yml` (`-run 'TestLiveEnsignCycle|TestLiveZeroDiscoverReportsAndStops'`). Both are cheap boot/front-door drives, so they share the step — no new CI job. Commit 7c02565a.
- SKIPPED: PREFERRED path — integrate as a shared scenario in the `TestLiveClaudeSharedScenarios` registry.
  The shared-scenario harness asserts a workflow LIFECYCLE (before/after entity reads, archival checks) and enforces a Codex+Claude+Pi parity guard (`TestSharedScenarioRunnerCoverage`: every shared ID needs a Codex runner, a Claude runner, AND a Pi coverage entry). A zero-discover boot has NO entity and no durable end-state — only a boot transcript — so it doesn't fit those assertions, and forcing it in would add three host runners + a Pi entry, MORE surface than the fallback. The team-lead authorized the `-run` fallback for exactly this shape mismatch.
- DONE: VERIFY the registration actually wires it in.
  `go test -tags live -list 'TestLiveEnsignCycle|TestLiveZeroDiscoverReportsAndStops' ./internal/ensigncycle/` lists BOTH test names — the exact CI `-run` pattern selects the new test. `grep TestLiveZeroDiscoverReportsAndStops .github/workflows/runtime-live-e2e.yml` → the gate's invocation now names it. Not just authored — proven invoked.
- DONE: `go test ./internal/ensigncycle/` + `go vet -tags live` green.
  `go test ./internal/ensigncycle/` → ok; `go vet -tags live ./internal/ensigncycle/` → clean. AC-1 detector + table test untouched (load-bearing, preserved).

### Summary

Closed the validation-flagged registration gap: `TestLiveZeroDiscoverReportsAndStops` is now invoked by the live-e2e gate, folded into the existing live-cycle canary step's `-run` alternation (no new job). Verified by `go test -list` with the CI `-run` pattern selecting the test and by grepping the yaml — proven invoked, not just authored. Took the team-lead's authorized fallback over the preferred shared-scenario path because that registry asserts a workflow lifecycle and enforces a 3-host parity guard, which a no-entity zero-discover boot doesn't fit (it would add more surface than the fallback). Offline gate + live vet green; the load-bearing AC-1 detector and table test are untouched.

## Stage Report: validation (cycle 2)

- DONE: Reproduce the detector from the worktree: `go test ./internal/ensigncycle -run TestDetectBroadSearchAtBoot` green; adversarially confirm it reds a repo-rooted find/grep -r/recursive Glob/repo-wide Grep after zero --discover, passes a scoped search under a resolved workflow dir (false-red edge), a clean report-and-stop, and empty stream; iterates ALL tool_use blocks (second-block evasion reds).
  Table test 8/8 + `TestDetectBroadSearchAtBootSecondBlock` green. Ran a throwaway 17-case harness driving `detectBroadSearchAtBoot` directly: every required RED redded (find/grep -r/rg/fd/ls -R repo-root, cwd-default `find .`, recursive Glob `**/README.md`+`**/*.md`, repo-wide Grep unset+root path); every required PASS passed (scoped-under-resolved-dir false-red edge for Bash AND Grep tool, clean version/gitrev/discover, empty stream, non-recursive `*.md` glob). Detector keys on TARGET (root / `**` pattern), not tool name. Harness removed.
- DONE: Verify the live scenario is ACTUALLY wired into CI (the cycle-1 gap): confirm `runtime-live-e2e.yml` invocation includes TestLiveZeroDiscoverReportsAndStops in its -run alternation, AND `go test -tags live -list 'TestLiveZeroDiscoverReportsAndStops'` selects it.
  CI gap CLOSED. `runtime-live-e2e.yml:179` — the actual `run:` of the "Run live ensign cycle" step — is `-run 'TestLiveEnsignCycle|TestLiveZeroDiscoverReportsAndStops'` (verified the line is the gate invocation, not a comment). `go test -tags live -list 'TestLiveZeroDiscoverReportsAndStops'` returns the test name — the gate's pattern selects it. Invoked, not just authored.
- DONE: Offline gate green: `go test ./internal/ensigncycle/` + `go vet -tags live ./internal/ensigncycle/`; no agents/references/contract prose touched.
  `go test -count=1 ./internal/ensigncycle/` → ok (4.6s); `go vet -tags live ./internal/ensigncycle/` → clean. `git diff --name-only main...HEAD` = 5 files (`runtime-live-e2e.yml` + 4 `internal/ensigncycle/*_test.go`); grep for `agents/`/`references/`/`skills/*.md` → NONE. The proof is the detector, not a contract clause.

### Summary

Validation PASSES. The cycle-1 registration gap I flagged is genuinely closed: `TestLiveZeroDiscoverReportsAndStops` is now in the live-e2e gate's `-run` alternation (yaml:179, the real `run:` block) and `go test -list` confirms the CI pattern selects it. AC-1's detector reproduced green (table 8/8 + second-block) and survived an independent 17-case adversarial harness covering every checklist edge — including both false-red passes (scoped Bash and Grep under a resolved workflow dir) and the second-block evasion red. Offline gate + live vet green; diff is CI yaml + test files only, no contract/agents/references prose touched.

## Stage Report: implementation (cycle 3)

- DONE: Trace every no-workflow output path a booting FO hits (the `--discover` zero/empty result, the `native_runner.go:79` no-workflow error, and `--boot`).
  `status --discover` early-returns at `dispatch()` line 49-51 → on zero, empty stdout + exit 0 (silent, no invitation). Bare `status` AND `status --boot` both flow through the SINGLE discovery gate at lines 75-82 (before any `--boot`-specific handling) → the one inviting message at line 79. Confirmed by exercise: built binary in a bare tmpdir emits the message for `status` and `status --boot`, empty stdout for `status --discover`.
- DONE: Fix the message(s) to be self-evidently TERMINAL and NON-INVITING.
  `native_runner.go:79` reframed to `no commissioned Spacedock workflow found in {dir} — report this and stop; do NOT search the filesystem for one. If a workflow exists elsewhere, point at it with --workflow-dir <dir>.` Names the root, directs report-and-stop, explicit do-not-search, keeps `--workflow-dir` for the wrong-dir human. Commit e8ddf051.
- DONE: AC-3 (load-bearing) — a Go test asserting the no-workflow command OUTPUT carries the terminal directive and not a hunt-invitation.
  `TestNoWorkflowHereError` rewritten to a command-OUTPUT assertion (stderr contains "report this and stop" + "do NOT search the filesystem" + root name + "--workflow-dir", and NOT "run inside a workflow"). Proven to bite: reverting the message to the old phrasing reds all four checks (verified, then restored). `TestPlainResolveFromNonWorkflowEmitsNoWorkflow` + `root_resolve_discovery` negative guard updated to the new stable marker.
- DONE: AC-1 (detector) + AC-2 (registered live scenario) remain green as the regression guard.
  `go test ./internal/ensigncycle/` green (detector + table test untouched); `go vet -tags live` clean; live scenario still registered in `runtime-live-e2e.yml`.
- DONE: `go test ./...` green.
  Full offline suite: 15 packages ok, 0 failures.

### Summary

Cycle-3 closes the captain's scope expansion: the deliverable is now the BINARY FIX, with the detector demoted to its guard. Traced the boot FO's no-workflow output paths — `status --discover` is already silent on zero (empty stdout, exit 0), and both bare `status` and `status --boot` flow through one discovery-gate message (`native_runner.go:79`) that read as "go hunt a workflow." Reframed it to a self-evidently terminal "report this and stop; do NOT search the filesystem" naming the searched root, keeping `--workflow-dir` for the wrong-dir human. The load-bearing AC-3 is a command-OUTPUT assertion (the CLI's own stderr IS the behavior) proven to bite by reverting the message; AC-1's detector and AC-2's registered live scenario remain the FO-behavior regression guard. Full suite green. Scoped out the parallel `state init`/`state new` phrasing — write commands the boot FO does not run.
