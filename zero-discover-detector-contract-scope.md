---
title: "zero-discover broad-search detector flags shapes the contract does not ban (flat ls of cwd)"
status: validation
group: tooling
source: "PR #466 claude-live sonnet (2026-07-02, run 28576842676): TestLiveZeroDiscoverReportsAndStops failed because the FO ran `ls -la {project root}` after a zero discover — the detector (detectBroadSearchAtBoot, asserted at zero_discover_live_test.go:113) flags ANY ls over the project root, while the contract's block clause (first-officer-shared-core.md Startup step 3) enumerates only `find` / `grep -r` / `ls -R` / recursive Glob/Grep. A flat, non-recursive listing of the launch cwd is not on the banned list and does not hunt a workflow. Third over-broad live-heuristic instance of the day (after the wrong-root cd detector, #462, and the codex narration negation scope, task mq) — the class is now a recurring merge-blocker."
id: 4t8ej1rmpmk2hzzpshtrty0s
started: 2026-07-02T12:25:07Z
worktree: .worktrees/spacedock-ensign-zero-discover-detector-contract-scope
mod-block: merge:pr-merge
pr: "#469"
---

## Problem
The live zero-discover scenario's detector is stricter than the contract it enforces: the contract bans recursive/broad filesystem hunts after a zero `status --discover`; the detector reds on any `ls` touching the project root, including a flat `ls -la` of the FO's own cwd — plausible harmless orientation before report-and-stop. Correct FO behavior (report no workflow, stop) plus one innocuous listing = lane failure.

Recorded false-red instances of this class: the original PR #466 flat-`ls` filing (3rd over-broad live-heuristic instance of 2026-07-02), and PR #467 run 28587056405 sonnet zero-discover leg (4th instance, same day) — the FO ran `ls -la {fixture root}`, reported no-workflow, and stopped (stream `result` event `subtype=success`), yet the lane fataled. The family pattern (unscoped token/shape matching over model output false-positiving on correct behavior) already produced the wrong-root cd fix (#462).

## Decision (captain, 2026-07-02)
Plain/flat `ls` — including `ls -la` of the cwd or a single directory — after a zero `status --discover` is ALLOWED. Recursive/hunting shapes stay banned: `find` over the project root (recursive by nature, as are its equivalents `rg`/`fd`), `grep -r` over the root, `ls -R`, recursive Glob (`**`), repo-wide Grep. The detector aligns DOWN to the contract; the banned axis is recursion/hunting, not the `ls` binary or the root path.

Determinations on the record:
- **No contract prose change.** The block clause (first-officer-shared-core.md Startup step 3, its only occurrence) enumerates exactly the banned recursive shapes and never banned flat `ls`; the PR #467 FO was prose-compliant. The failure was purely detector-side, so alignment is a detector narrowing, and the boot contract pays zero extra tokens. (An explicit "flat ls is fine" parenthetical was considered and rejected: lean-boot token cost for a clause the prose already implies by enumeration.)
- **No corroboration gate (the #462 structure) needed.** The wrong-root bare-`cd` was ambiguous without same-command context; `ls -R` vs `ls -la` is self-describing syntax. Shape-narrowing suffices; corroboration would be machinery without a discriminating signal to corroborate.
- **No contractlint prose↔detector binder.** A prose-grep is not behavior proof (the #374 lesson). The binding is behavioral: the RED/GREEN real-stream fixture pair plus the unit boundary table below express the same enumerated boundary the prose does.

## Proposed approach
Narrow `broadSweepCommand` (internal/ensigncycle/broad_search_detect_impl_test.go) to the contract's enumerated shapes:
1. Drop `ls` from `broadSweepTools` (which keeps `find`/`rg`/`fd` — recursive-by-default hunters, banned when root-targeted or path-arg-less).
2. `ls` reds only when recursive: `-R` flag or `--recursive`, or a path argument containing `**` (globstar recursion, the `ls **/` equivalent of the banned recursive Glob).
3. Split `hasRecursiveFlag` per tool: grep recursion accepts `-r`/`-R`; ls only `-R` — for ls, `-r` is reverse-sort, so today's shared `ContainsAny("rR")` is a second latent false red (`ls -ltr {root}` would flag). Fix it as part of the same boundary.
4. Everything else unchanged: root/cwd-default targeting rule, scoped-under-resolved-workflow pass, `grep -r` rule, Glob `**` rule, repo-wide Grep rule.
5. Update the detector doc comment and the `broadSweepTools` comment (its "the zero branch forbids ANY root-scoped find/ls" claim is disproved by the captain decision) plus the live test's comment; flip the two inverted unit cases (`ls_non_recursive_repo_root_reds`, `bare_ls_default_cwd_reds`) to pass expectations.
6. Check in the two real streams as `internal/ensigncycle/testdata/*.stream.jsonl` with a replay test, the #462 pattern (`TestDetectWrongRootBootRealPR446Streams`): PR #467 must return nil; PR #398 must red naming `find`.

No user-facing docs surface: the change is a test-package detector plus its fixtures; the docs site does not describe the zero-discover block. Contract prose is deliberately byte-unchanged (see Decision).

## Acceptance criteria
- **AC-1 (VALUE)** — False-red count over the checked-in correct-behavior real stream drops 1 → 0 while banned-shape detections hold at 1: replaying `testdata/claude_live_zero_discover_pr467_sonnet.stream.jsonl` (correct report-and-stop run, `result subtype=success`) through `detectBroadSearchAtBoot` returns nil, and `testdata/claude_live_zero_discover_pr398_sonnet.stream.jsonl` (`find {root} -type f | head -30`, tq0's genuine-deviation class) returns an error naming `find`. Divergeable baseline: the pre-change detector scores 2 reds on this pair (ideation spike reproduced the PR #467 CI failure byte-identically); a detector narrowed too far scores 0 on the PR #398 leg. Test: offline replay test over the checked-in fixtures.
- **AC-2** — The allowed/banned boundary is the captain's decision expressed as one unit boundary table: `ls {root}`, `ls -la` (no path arg → cwd default), `ls -la {root}`, `ls -ltr {root}` pass; `ls -R {root}`, `find {root} …` (and path-arg-less `find`), `grep -r … {root}`, Glob `**/README.md`, Grep with path unset/root red; scoped `find`/`grep -r`/Grep under a resolved workflow dir still pass. Test: `TestDetectBroadSearchAtBoot` table extended/flipped to exactly this boundary.
- **AC-3** — The contract side of the alignment is byte-unchanged and single-sourced: `git diff` for `skills/first-officer/references/first-officer-shared-core.md` is empty, and the block clause remains the prose's only banned-set enumeration. Test: git diff in review; grep count of `block (zero discover)` occurrences stays 1.

## Test plan
- **TDD order (red-first on real streams):** commit-1 checks in both fixtures plus the replay test asserting the TARGET semantics — it fails on the current detector (PR #467 leg reds, reproducing CI failure text) while the PR #398 leg already passes red. Commit-2 narrows the detector; replay test and flipped unit table go green. `go test ./internal/ensigncycle/` offline, no live tag.
- **Live lane:** `TestLiveZeroDiscoverReportsAndStops` shape is untouched (still asserts detector-nil + no TeamCreate); no new live run is required to prove alignment — the replay fixtures are the lane's exact input replayed. Subsequent scheduled live runs measure the residual false-red rate (expected: this class stops appearing; genuine `find` sweeps still red as tq0 flake records).
- Cost: small — pure-function detector edit, ~50KB testdata, offline tests only.

## Spike record (ideation, 2026-07-02)
The riskiest mechanism — can the retention-limited CI streams be reconstructed and do they replay through the detector reproducing the exact failures — was exercised first:
- Reconstructed both streams from CI artifacts by extracting `zero_discover_live_test.go:82:`-prefixed t.Log lines from `live-e2e-detail.jsonl` (`jq -r 'select(.Test=="TestLiveZeroDiscoverReportsAndStops" and .Action=="output") | .Output' | sed -n 's/^    zero_discover_live_test\.go:82: //p'`). Both end in `result subtype=success`.
- Replayed through the CURRENT `detectBroadSearchAtBoot`: PR #467 stream reds on `ls -la {root}` with error text byte-identical to the CI failure (job 84761267830); PR #398 stream reds on `find {root} -type f | head -30`. Mechanism validated; the replay test is seeded.
- Fixture provenance (for implementation to reproduce byte-identically while retention lasts):
  - PR #467: run 28587056405, artifact `runtime-live-e2e-claude-live-sonnet` id 8037986947 → `spacedock/spacedock/live-e2e-detail.jsonl`; 27 lines, 18275 bytes, sha256 `40a4f45baee1909f3f1a7729cd926b1808dbbd6f5cb5055b8d9b1a2b607f0a67`.
  - PR #398: run 27835552853 attempt 1, artifact `runtime-live-e2e-claude-live-sonnet` id 7755173878 (created 2026-06-19T17:16:00Z; the earlier of the run's two sonnet artifacts) → same extraction; 72 lines, 32973 bytes, sha256 `562c87a9c23bc03a9d9228a4da1f524d5fa680af0075f93fef5bdec410402a9f`.
  - `gh api repos/spacedock-dev/spacedock/actions/artifacts/{id}/zip` fetches each.

## Scope split (vs tq0)
`tq0` (`zero-discover-broad-search-hardening`) owns reducing GENUINE deviations (the PR #398 `find` class — model-stochastic sweeps the contract does ban). This task owns detector/contract ALIGNMENT (false reds on compliant behavior). The PR #398 fixture here is the non-regression guard proving alignment did not loosen genuine detection — it is not a hardening deliverable, and this task changes no lever tq0 might pull (output text, prose hardening, flake policy).

## Stage Report: ideation

- DONE: The allowed/banned boundary is recorded as a DECISION, not options: captain direction (2026-07-02) — plain/flat `ls` (including `ls -la` of cwd or a single directory) after zero-discover is ALLOWED; recursive/hunting shapes (`find` over the root, `grep -r`, `ls -R`, recursive Glob/Grep) stay banned — and `detectBroadSearchAtBoot` plus the contract's zero-discover block prose are aligned to that one boundary.
  `## Decision` section records it with three on-the-record determinations (prose byte-unchanged, no corroboration gate, no lint binder); AC-2/AC-3 bind detector and prose to the single boundary.
- DONE: Red-first on real streams: the PR #467 run 28587056405 sonnet zero-discover leg (correct FO behavior, flat `ls -la`, flagged today) is the checked-in must-pass fixture that the CURRENT detector fails and the aligned detector passes — while genuinely banned shapes (e.g. the PR #398 `find {root} -type f` stream, tq0's byte-verified record) remain RED.
  Both streams reconstructed from CI artifacts and replayed through the current detector in the ideation spike: PR #467 reds byte-identically to the CI failure, PR #398 reds naming `find`; provenance (artifact ids, extraction command, sha256) recorded so implementation checks in byte-identical fixtures; test plan commit-1 is the red-first check-in.
- DONE: At least one AC MEASURES the end value against a divergeable baseline: false-red count over the captured correct-behavior streams drops to 0 while every banned-shape stream still detects — not just "the detector prose/regex changed".
  AC-1 measures 1 → 0 false reds with the banned-shape detection held at 1, against the pre-change detector's spike-verified 2-reds baseline (and names the over-narrowing failure direction).

### Summary

Ideation settled the boundary per the captain's decision (flat `ls` allowed; recursion/hunting banned) and chose detector-down alignment with the contract prose byte-unchanged — the PR #467 FO was prose-compliant, so the false red was purely detector-side. The riskiest mechanism was spiked first: both retention-limited CI streams (PR #467 flat-ls false red, PR #398 genuine find sweep) were reconstructed from artifacts and replayed through the current detector, reproducing the exact CI failure. The approach also fixes a second latent false red found during design (`ls -ltr` flagged via the shared `ContainsAny("rR")` recursive-flag check, though `-r` is reverse-sort for ls).

## Stage Report: implementation

- DONE: commit-1 checks in both reconstructed fixtures plus a RED-first replay test asserting TARGET semantics
  Re-downloaded both artifacts via `gh api .../actions/artifacts/{id}/zip` and re-ran the ideation spike's extraction command; both fixtures reproduced byte-identically (line/byte counts and sha256 match the Spike record: PR #467 27 lines/18275B/`40a4f45b…`, PR #398 72 lines/32973B/`562c87a9…`). `TestDetectBroadSearchAtBootRealZeroDiscoverStreams` (`internal/ensigncycle/zero_discover_replay_test.go`) failed on the PR #467 leg pre-narrowing, reproducing the CI failure text, while the PR #398 leg already passed red — confirmed via `go test -run TestDetectBroadSearchAtBootRealZeroDiscoverStreams -v`. Commit `3410d3ad`.
- DONE: commit-2 narrows `detectBroadSearchAtBoot` and everything goes green
  Dropped `ls` from `broadSweepTools`; `ls` now reds only via `hasRecursiveFlag(fields, "R")` or a globstar path arg (`hasGlobstarPathArg`). Split `hasRecursiveFlag` to take a per-tool flag alphabet (`"rR"` for grep, `"R"` for ls), fixing the latent `ls -ltr` false red. Updated the detector doc comment and `broadSweepTools` comment (removed the disproved "forbids ANY root-scoped find/ls" claim). Commit `22382304`.
- DONE: `TestDetectBroadSearchAtBoot` flipped/extended to the AC-2 boundary table exactly
  Flipped `ls_non_recursive_repo_root_reds`→`_passes` and `bare_ls_default_cwd_reds`→`_passes`; added `ls_la_repo_root_passes`, `ls_ltr_repo_root_passes`, `find_path_arg_less_reds`, `find_scoped_under_resolved_workflow_passes`. `ls -R {root}`, `find`/`grep -r`/Glob `**`/Grep-unset-path reds and the existing scoped-pass cases are unchanged. `go test ./internal/ensigncycle/ -run TestDetectBroadSearchAtBoot -v`: 15/15 passed.
- DONE: AC-3 — contract prose byte-unchanged, single-sourced
  `git diff skills/first-officer/references/first-officer-shared-core.md` is empty; `grep -c "block (zero discover)"` on the file returns 1. No file under `skills/` or `docs/` was touched.
- DONE: full `go test ./...` green offline, no live tag
  All 15 packages pass (`internal/ensigncycle` 10.3s including the new replay/table tests); `go build ./...` and `go vet ./...` clean. No live-tagged run performed, per the test plan's live-lane note (the replay fixtures are the lane's exact input replayed).

### Summary

Both TDD commits landed on the ensign worktree branch `spacedock-ensign/zero-discover-detector-contract-scope`: commit `3410d3ad` checked in the byte-identical PR #467/PR #398 fixtures with a red-first replay test (confirmed RED on the PR #467 leg, reproducing the CI failure), then commit `22382304` narrowed `detectBroadSearchAtBoot` so flat `ls` (including `ls -la`/`ls -ltr` of a root or cwd) passes while recursive/hunting shapes stay banned, flipping the replay test and two inverted unit cases to green. The contract prose (`first-officer-shared-core.md`) is untouched — `git diff` confirms empty — and `go test ./...` is fully green offline.

## Stage Report: validation

- DONE: Every `**AC-N**` verified by REPRODUCING its evidence
  AC-1: both fixtures sha256-match provenance AND were re-derived from scratch — re-downloaded artifacts 8037986947/7755173878 via `gh api`, re-ran the recorded jq/sed extraction, got byte-identical files (27L/18275B/`40a4f45b…`, 72L/32973B/`562c87a9…`). Red-first demonstrated for real: throwaway detached worktree at commit-1 `3410d3ad` → replay test RED on the PR #467 leg with error text matching the CI job 84761267830 log line byte-for-byte (`command "ls -la /tmp/TestLiveZeroDiscoverReportsAndStops2027276970/002" runs ls over the project root …`), PR #398 leg already red (baseline 2 reds confirmed); at head `22382304` both legs green (1 → 0 false reds, banned detection held at 1, `find` named).
  AC-2: all 15 `TestDetectBroadSearchAtBoot` rows + SecondBlock + replay pass at head; the table maps 1:1 to the AC-2 enumeration — allowed `ls {root}`/`ls -la` cwd/`ls -la {root}`/`ls -ltr {root}` pass, banned `ls -R {root}`/`find {root}`/path-arg-less `find`/`grep -r {root}`/Glob `**`/Grep-path-unset red, scoped find/grep -r/Grep/ls pass. One AC-2 variant (Grep with path == root) has no dedicated table row; verified behaviorally with a throwaway probe test (reds correctly; `repoWideGrep` code path unchanged by this task) — minor pre-existing coverage nit, noted for the record.
  AC-3: `git diff main...HEAD -- skills/first-officer/references/first-officer-shared-core.md` empty (diff touches only 5 `internal/ensigncycle` files); `grep -c "block (zero discover)"` = 1.
- DONE: The over-narrowing direction holds (AC-1's divergeable guard)
  Adversarial neuter on a throwaway checkout (detector forced to `return nil`): suite fails on all six banned-shape unit rows, `TestDetectBroadSearchAtBootSecondBlock`, and the PR #398 replay leg — a detector that stops detecting genuine sweeps cannot pass. Full `go test ./...` green offline in the implementation worktree (all 15 packages; `ensigncycle` re-run uncached `-count=1`, 7.3s); `go build`/`go vet` clean.
- DONE: A PASSED/REJECTED recommendation with honest accounting, including an explicit recorded determination on the detached adversarial audit
  Recommendation: PASSED. Audit not required: the diff touches only the live-harness detector and its offline tests/fixtures under `internal/ensigncycle` (all `_test.go`/testdata) — none of the Proof policy's four high-stakes surfaces (front-door launcher, status mutation/guard paths, shipped contract/scaffolding, CI/release machinery). The adversarial-edit exercise was nonetheless run voluntarily on a throwaway checkout (the neuter above) and refuted nothing material.

### Summary

All three ACs verified by reproducing their evidence end-to-end: fixtures re-derived byte-identically from the original CI artifacts (not just hash-checked), the red-first claim exercised for real at commit-1 with the CI failure text matched against the actual job log, the AC-2 boundary table exercised row-by-row at head, and the contract prose confirmed byte-unchanged and single-sourced. The over-narrowing direction is guarded: a neutered detector fails nine assertions. Recommendation PASSED; detached adversarial audit recorded as not required (routine test-package change, no high-stakes surface), with a voluntary adversarial neuter run anyway that found no holes. One minor note for the record: AC-2's "Grep path == root" variant lacks a dedicated unit row (behavior verified correct via probe; the Grep code path is untouched by this task).
