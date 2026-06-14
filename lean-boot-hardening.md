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
---

Keep FO boot lean: when `spacedock status --discover` returns zero workflows, the Startup discovery step must report no workflow found and STOP — never fall back to a broad `find`/`grep` filesystem sweep to hunt one down. Harden the discipline so the zero-`--discover` path is provably report-and-stop, not an expensive search.

## Problem

- Startup step 3 is explicit: `status --discover` → one path → use it; zero → report no workflow found; multiple → present the list. The zero branch is terminal.
- Observed (captain, 2026-06-14): an FO, after `--discover` returned zero, ran a broad filesystem sweep to locate a workflow — violating the contract's zero-branch AND the lean-boot ethos (cf. j9 shallow-boot: boot is cheap, it does not sweep the filesystem).
- A broad filesystem search at boot is both a discipline violation (the zero-branch is report-and-stop) and a cost/latency regression — the opposite of lean boot.

## Approach

The contract is already correct: Startup step 3 (`first-officer-shared-core.md:24`) says "zero → report no workflow found" — the zero branch is terminal. The gap is not the wording; it is that nothing *proves* the FO obeys it. A clause telling the FO not to broad-search has a wording-present ceiling (a paraphrase fails a grep, an inverted clause passes it), so adding more prose cannot be the acceptance criterion.

The proof is the same shape Spacedock already uses for the wrong-root wander (PR #365): a model-agnostic **stream-scanner detector** that reads the FO's captured boot transcript (the `tool_use` stream, not any model phrasing) and reds when a zero-`--discover` boot contains a broad filesystem-sweep tool call. This composes only already-proven mechanisms:

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
- **No contract rewrite.** Step 3 already states the terminal zero branch; this task ships the proof, not new prose. (If the staff review judges a one-line reinforcing clause helpful, it is authoring work, not an AC — the AC is the detector.)
- Multi-workflow (`multiple → present the list`) and the explicit-path branch are unaffected and untouched.

## Acceptance criteria

- **AC-1 — The detector reds a zero-discover boot that broad-searches, and passes one that report-and-stops.** A pure `detectBroadSearchAtBoot(stream, fixtureRoot)` exists in `internal/ensigncycle` and: reds a stream where, after a zero `status --discover`, the FO issues a repo-rooted `find`/`grep -r`/`rg`/`ls -R` Bash, or a recursive `Glob`/`Grep` over the project root, to hunt a workflow; passes a stream whose only boot tool calls are `--version`, `git rev-parse`, `status --discover`, and the no-workflow report; passes a *scoped* search under an already-resolved workflow dir (no false-red); passes an empty stream.
  - Verified by: `go test ./internal/ensigncycle -run TestDetectBroadSearchAtBoot` — a table test driving the detector with hand-built stream-json lines (the `streamLine` helper already in `wrong_root_detect_test.go`), asserting error/no-error per case and that the error names the offending command. The expected values come from the test's own crafted streams, NOT from any contract file — independent source of truth, so the check can fail. Red before the detector exists, green after.
- **AC-2 — A live zero-discover FO boot is observed to report-and-stop, taking no broad filesystem sweep.** The `//go:build live` cycle gains a scenario: launch the real `spacedock claude` FO against a zero-discover fixture (git-init'd tmpdir root, no `commissioned-by: spacedock@` README), capture the boot stream, and assert (a) the FO reaches a greet/stop or no-workflow report without a TeamCreate, and (b) `detectBroadSearchAtBoot(transcript, fixtureRoot) == nil`.
  - Verified by: `SPACEDOCK_LIVE=1 go test -tags live ./internal/ensigncycle -run TestLiveZeroDiscoverReportsAndStops` (live-gated, model from `SPACEDOCK_LIVE_MODEL`). The proof is the captured transcript driven through the AC-1 detector: a real model boot, observed behavior, not prose. On a sweep, the detector reds and names the command.
  - This AC's observable is the boot transcript itself: a real FO process, the detector's verdict on its tool stream. It is the behavioral half AC-1's offline detector cannot stand in for.

## Test plan

- **AC-1 (offline, cheap — minutes):** `TestDetectBroadSearchAtBoot` table test, ~6 cases (repo-rooted `find` reds; `grep -r` repo root reds; recursive `Glob **/README.md` reds; scoped grep under a resolved workflow passes; clean report-and-stop passes; empty stream passes). Pure function, no model, runs in the default suite. Mirrors `TestDetectWrongRootBoot` structure exactly. This is the load-bearing proof and the riskiest-first exercise — write it before the live scenario.
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
