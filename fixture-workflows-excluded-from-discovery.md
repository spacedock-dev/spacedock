---
title: Fixture workflow READMEs must not be workflow-discovery candidates
status: validation
source: "Live FO session, 2026-07-21, after the refit-content-propagation fixtures landed on main."
id: ab3ma8m7gsm8tra2ksmcdydq
started: 2026-07-21T16:05:13Z
worktree: .worktrees/spacedock-ensign-fixture-workflows-excluded-from-discovery
---

A commissioned-shape README used as a test fixture is counted as a real workflow by downward discovery, so a one-workflow repo now looks like a two-workflow repo to every command that auto-discovers.

## Problem

`internal/status/handlers.go:discoverWorkflows` walks down from the git toplevel and treats any directory whose `README.md` frontmatter starts with `commissioned-by: spacedock@` as a workflow. `fixtures/refit-content-propagation/site-workflow/README.md` is exactly such a README on purpose — it is "what the commission skill emits verbatim" at `spacedock@0.25.0`, the observation the refit-content-propagation fixture exists for (see `fixtures/refit-content-propagation/README.md`). The prune set `discoverIgnoreDirs` (`.git`, `.worktrees`, `node_modules`, `vendor`, `dist`, `build`, `__pycache__`, `tests`, `.spacedock-state`) does not include `fixtures`, and `.gitignore` does not list it (the fixture is committed), so the walk descends into `fixtures/` and reports the fixture as a second workflow.

Reproduced from the repo root (this is the defect, exit 0 with two lines):

```
$ ./spacedock status --discover
/Users/clkao/git/spacedock-research/spacedock-v1/docs/dev
/Users/clkao/git/spacedock-research/spacedock-v1/fixtures/refit-content-propagation/site-workflow
```

Because the repo now looks like a two-workflow repo, every command that resolves the workflow by downward discovery (no `--workflow-dir`) refuses with the multi-workflow ambiguity error and demands `--workflow-dir`, for a repo that operationally has ONE workflow (`docs/dev`). The binary handles the ambiguity *correctly* (clear named diagnostic, non-zero exit — verified live); the ambiguity itself is the bug because a test fixture is not a real workflow. This is the root defect that fans out to every auto-discovering call site below; it is not an error-handling bug.

Affected discovery consumers (all route through `discoverWorkflows`):
- `status --new` (the `new` command) — `internal/status/native_runner.go:87` via `ResolveWorkflowDir` → `discoverWorkflowDownward` → `discoverWorkflows`.
- `state commit` — `internal/cli/state_sync.go:454` via `ResolveWorkflowDir`.
- `status` without `--workflow-dir` (table/read) — `internal/status/native_runner.go:87`.
- `status --discover` — `internal/status/native_runner.go:runDiscover` → `discoverWorkflows`.
- `status --boot --identify` — `internal/status/native_runner.go:resolveIdentifyBootDir` → `discoverWorkflows`.
- The front-door launch banner — `internal/cli/frontdoor.go:199` via `status.DiscoverWorkflows`.

Session note recorded alongside the seed (not part of this fix): the FO's own `cmd 2>&1 | tail -1` piping masked a non-zero exit inside an `&&` chain — a shell-discipline lesson for the FO, not a binary defect.

## Proposed approach

**Captain re-scope (2026-07-22).** The earlier in-place prune — add `fixtures` and `testdata` to `discoverIgnoreDirs` (landed as commit 1f8b666a, since removed) — is SUPERSEDED. The captain chose to REMOVE the top-level `fixtures/` directory rather than legitimize its name in the prune set, and rejected a top-level `testdata/`. The prior cheapest-first option analysis (marker field, entity-scoped resolution) is preserved in the ideation stage report below.

Two moves against the single shared resolver:

- **Relocate every fixture to package-adjacent `testdata/` (Go-idiomatic).** The three fixtures are live/manual-drive scaffolding with NO Go-test consumer (traced: the sole code reference is a reference-only classifier example, not a file read). Their drivers are the `refit` and `first-officer` skills, neither a Go package. `skills/integration` is the one Go test package among the skills (the live-drive harness) and its dir is already excluded from the contractlint instruction-surface walk (`shippedInstructionMarkdown` SkipDir's `integration`), so it is the clean package-adjacent home:
  - `fixtures/refit-content-propagation/` → `skills/integration/testdata/refit-content-propagation/`
  - `fixtures/entity-label-drive/` → `skills/integration/testdata/entity-label-drive/`
  Top-level `fixtures/` ends up empty and is removed.

- **Prune the `testdata` basename ONLY from `discoverIgnoreDirs`** (drop the `fixtures` prune — `fixtures/` no longer exists). `testdata` is the Go-tool-ignored test-scaffolding convention, so a commissioned-shape README carried as a package-adjacent fixture is no longer counted as a real workflow — the same test-scaffolding tradeoff the set already makes for `tests`/`vendor`/`dist`/`build`. Fixing `discoverWorkflows` fixes **every** auto-discovering consumer at once.
  - What it costs: a real workflow deliberately placed *under* a `testdata/` directory would no longer be auto-discovered. Accepted — real workflows live at `docs/dev`; identical tradeoff already made for `tests`. The walk always inspects its start-root regardless of basename (only descent into a `testdata` *child* is pruned), so `--workflow-dir`/`--root` pointed directly at such a path still resolves; only the auto-walk skips it.

The moved fixtures' own READMEs point at their drivers (`skills/refit/SKILL.md`, `skills/first-officer`) and sibling dirs by relative name — all still valid, the skills are unchanged. The one code reference to a fixture path — the FO-write classifier example at `internal/contractlint/fo_write_core_mutation_gate_test.go` — is updated to the new path (still classifies `blocked-product` via `skills/**`).

No `--workflow-dir` / `--root` override behavior changes. The clear-error-on-ambiguity path stays exactly as is (correct, out of scope).

**Documentation:** no doc diff required. `discoverIgnoreDirs` is an internal implementation detail not enumerated in any user-facing doc; this change *restores* the documented single-workflow auto-discovery, so no user-visible command surface or output wording changes.

## Out of scope

- The ambiguity error path itself (clear message + non-zero exit) — already correct, verified live.
- Any `--workflow-dir` / `--root` / `PIPELINE_DIR` override behavior.
- The FO's shell-piping lesson from the session note.
- Removing, moving, or re-marking the `refit-content-propagation` fixture.

## Acceptance criteria

Each AC names a property of the finished task, not a stage action, and how it is verified.

**AC-1 (VALUE) — The top-level `fixtures/` directory is gone AND downward discovery returns EXACTLY the one real workflow.**
Two coupled properties: (i) no top-level `fixtures/` directory remains (its three fixtures relocated under `skills/integration/testdata/`), and (ii) the discovered-workflow set for the repo toplevel is exactly `{<repo>/docs/dev}` (count 1), down from the pre-change 2 (`docs/dev` + `fixtures/refit-content-propagation/site-workflow`). The count is the independent baseline that can move the wrong way: too-aggressive → 0, ineffective → 2.
Verified by: `test ! -d fixtures` at the repo root; and `./spacedock status --discover` from the repo root prints exactly one line, `<repo>/docs/dev`, exit 0 (pre-change printed two lines). Removing the `testdata` entry from `discoverIgnoreDirs` re-surfaces a `testdata`-nested commissioned README (AC-3 reds).

**AC-2 (VALUE) — The auto-discovering commands succeed without `--workflow-dir` in this one-workflow repo.**
`new` and `state commit` (and the read `status`) resolve `docs/dev` by discovery alone, instead of exiting non-zero with the multi-workflow "pass --workflow-dir" refusal.
Verified by: from the repo root with no `--workflow-dir`, `./spacedock status` (read) exits 0 and lists `docs/dev` entities; `new`/`state commit` share the resolver `discoverWorkflows`, exercised by the AC-3 behavior test. Failure mode that would make it red: discovery again returns ≥2, re-triggering the refusal.

**AC-3 — A behavior test proves the `testdata` prune excludes a nested commissioned README while still finding the real workflow.**
A Go test builds a temp tree with a real workflow (`docs/dev/README.md`, commissioned) plus a commissioned README nested under `skills/integration/testdata/…` (the relocated-fixture shape), calls `discoverWorkflows(repo)`, and asserts the result is exactly `[<repo>/docs/dev]`. It is the falsifiable guard: removing the `testdata` entry from `discoverIgnoreDirs` re-introduces the fixture row and reds the test.
Verified by: `TestDiscoverWorkflowsPrunesTestdata` in `internal/status/discover_worktree_noise_test.go` (sibling of `TestDiscoverWorkflowsSkipsNestedCheckout`); `go test ./internal/status/...` green with the prune, red without it (watched RED→GREEN).

**AC-4 — The full Go suite stays green with the fixtures relocated and the prune applied.**
Relocating the fixtures newly exposes their `*.md` frontmatter to the `internal/status` migration-check walk (it walks `skills/`, does not prune `testdata`) and to the `skills/integration` package's own guards; the `testdata` prune touches no start-root `--root`/`--workflow-dir` fixture (those are inspected directly, not as descended children). Nothing breaks.
Verified by: `go test ./...` and `go test ./... -race` pass — including `internal/status` (migration check over the moved fixtures), `skills/integration`, and `internal/contractlint` (the updated FO-write classifier example + the `shippedInstructionMarkdown` guard, which excludes the `integration` dir).

## Test plan

- **Primary (behavior, ~36 LOC):** `TestDiscoverWorkflowsPrunesTestdata`, a Go unit test alongside `internal/status/discover_worktree_noise_test.go`, same shape as `TestDiscoverWorkflowsSkipsNestedCheckout` — write a temp tree with the real `docs/dev` workflow plus a commissioned README nested under `skills/integration/testdata/…`, assert `discoverWorkflows(repo)` returns exactly `[docs/dev]`. Written first and watched fail (returns the fixture row) before the `testdata` prune lands. It covers AC-1(ii)/AC-2/AC-3 because `discoverWorkflows` is the single resolver behind `--discover`, `new`, `state commit`, boot `--identify`, and the front-door banner — driving each command separately re-proves the same resolver.
- **Value check — `fixtures/` gone (AC-1(i), near-zero cost):** `test ! -d fixtures` at the repo root; the git move is a set of pure renames.
- **Real-repo check (AC-1(ii), CLI, near-zero cost):** `./spacedock status --discover` from the repo root returns the single `docs/dev` line, exit 0.
- **Regression (AC-4):** `go test ./...`, then `go test ./... -race`. Includes the migration-check walk over the relocated fixtures (they must parse cleanly) and the `skills/integration` + `contractlint` guards.

The moved fixtures have no automated Go-test consumer (traced), so "run the consumer green" reduces to: the relocated `*.md` still parse under the migration check and the full suite stays green. The simplest alternative to a behavior test — asserting the source contains the string `"testdata"` — is insufficient: it is a static prose check that passes even if the map key is misspelled or the walk never consults it. Exercising `discoverWorkflows` over a real tree is the proof.

## Expected surface

Code (3 files):
- `internal/status/handlers.go` — `"testdata": true` on the `discoverIgnoreDirs` map + expanded comment: ~4 changed lines.
- `internal/status/discover_worktree_noise_test.go` — one appended `TestDiscoverWorkflowsPrunesTestdata`: ~36 LOC.
- `internal/contractlint/fo_write_core_mutation_gate_test.go` — 1 line: the classifier-example path updated to the new location.

Moves (7 files, pure renames, 0 content change):
- `fixtures/refit-content-propagation/**` (2 files) → `skills/integration/testdata/refit-content-propagation/**`
- `fixtures/entity-label-drive/**` (5 files) → `skills/integration/testdata/entity-label-drive/**`

Total ~41 code LOC across 3 files + 7 renames + the removed empty `fixtures/` dir. Tolerance ±15 LOC; the file count is fixed by the fixture set (7 moves) plus the 3 code files.

## Spike (riskiest unverified mechanism)

Riskiest unknown for the re-scope (named in the dispatch): the three fixtures may have no consuming Go package, so "relocate to a consuming package's testdata/" might have no valid home. **Traced and resolved:**

1. **Consumer trace** (grep over all `*.go`/`*.sh`/`*.yml`/`*.md`, `.github`, `go:embed`, Makefile): none of the three fixtures has a Go-test consumer that reads it. `refit-content-propagation/site-workflow` is driven manually against `skills/refit/SKILL.md` Phase 3b; `entity-label-drive/*` by a live first-officer drive. The sole code reference — `fo_write_core_mutation_gate_test.go` — is reference-only (a classifier path-string example; it reads `fo-write-core.md`, never the fixture). Escalated the no-Go-consumer ambiguity to the FO; captain chose `skills/integration/testdata/` (Option B: the one real Go package among the skills, already excluded from the contractlint instruction-surface walk — `shippedInstructionMarkdown` SkipDir's `integration`).
2. **Second riskiest — migration-check exposure:** relocating under `skills/` newly exposes fixture frontmatter to the `internal/status` migration-check walk (walks `skills/`, does not prune `testdata`). Zero-cost pre-check confirmed every fixture frontmatter parses cleanly (top-level scalars agree with a direct yaml.v3 decode; the bare RFC3339 `started` is safe under the node-based comparison). `go test ./internal/status/` GREEN post-move.
3. **Discovery mechanism:** `TestDiscoverWorkflowsPrunesTestdata` watched RED (fixture returned as a second workflow) before the `testdata` prune, GREEN after. Real binary `status --discover` from the repo root → exactly `docs/dev`, exit 0; `fixtures/` removed.

Relied-on proven mechanisms: `discoverWorkflows`'s existing `ignore[ent.Name()]` basename prune (already the mechanism for `tests`/`vendor`/`dist`) and the always-walk-the-start-root behavior.

## Validation note

The **detached adversarial audit trigger applies**: this change touches two of the four high-stakes surfaces — the **front-door launcher** (`frontdoor.go` consumes `DiscoverWorkflows`) and the **`status` mutation/guard/discovery path** (`ResolveWorkflowDir` backs `--set` mutation resolution, the merge guard, and `new`/`state commit`). Validation should run the read-only audit on a throwaway checkout: construct a claim-breaking edit (remove the `testdata` entry from `discoverIgnoreDirs`, or misspell it) and confirm `TestDiscoverWorkflowsPrunesTestdata` reds. Additionally confirm the value AC's move dimension: `fixtures/` is gone and its contents are present under `skills/integration/testdata/` as pure renames (no content drift), and the relocated fixtures parse under the `internal/status` migration-check walk. The AC-provenance sub-trigger does NOT fire — the value AC's expected value (count 1, the known real `docs/dev`, plus the structural `fixtures/`-absent fact) is an independent constant, not derived from the package's own production functions or constants.

## Stage Report: ideation

- DONE: Flesh out the task body: replace the single seed paragraph with structured Problem, Proposed approach, Out of scope, Acceptance criteria, and Test plan sections (the README template shape).
  Body now carries Problem / Proposed approach / Out of scope / Acceptance criteria / Test plan / Expected surface / Spike / Validation note, matching `docs/dev/README.md` task-template + Verified-by house style.
- DONE: Evaluate the fix options cheapest-first and choose one with rationale.
  Chose (a) add `fixtures`/`testdata` to `discoverIgnoreDirs`; rejected (b) marker-field (corrupts the verbatim-commission fixture, adds standing author burden) and (c) entity-scoped resolution (only fixes `state commit`, leaves `new`/discover/banner broken). Rationale recorded in Proposed approach.
- DONE: Locate and name the discovery code path and the auto-discovering commands.
  Root: `internal/status/handlers.go:discoverWorkflows` (+ `discoverIgnoreDirs` at :471). Consumers: `native_runner.go:87` (`new`, read `status`), `state_sync.go:454` (`state commit`), `runDiscover`/`resolveIdentifyBootDir`, `frontdoor.go:199` (banner) — all via `ResolveWorkflowDir`/`DiscoverWorkflows`.
- DONE: Write Acceptance criteria including at least one value-measuring AC.
  AC-1 (VALUE): `--discover` returns exactly `{docs/dev}`, count 2→1 (baseline can move to 0 or 2). AC-2 (VALUE): `new`/`state commit` succeed without `--workflow-dir`.
- DONE: Write the Test plan (behavior test over discovery; go test ./..., go test ./... -race).
  Primary: `discoverWorkflows(repo)` behavior test alongside `TestDiscoverWorkflowsSkipsNestedCheckout`; real-repo `--discover` check; regression `go test ./...` and `-race`.
- DONE: Declare the expected surface (files + LOC) and tolerance.
  ~38 LOC / 2 files (handlers.go ~3, one test ~35); tolerance ±15 LOC and ±1 file.
- DONE: Spike the riskiest unverified mechanism (filter excludes fixture while still finding the real workflow).
  Reproduced 2 workflows on the real binary; throwaway prune edit → `status --discover` returned exactly `docs/dev` (2→1); `go test ./internal/status/...` PASS with the edit; reverted, tree clean. Recorded in Spike section.
- DONE: Note for validation whether the detached adversarial audit trigger applies.
  Applies — touches the front-door launcher and the `status` mutation/guard path (2 of 4 high-stakes surfaces). AC-provenance sub-trigger does not fire. Ambiguity error-handling kept out of scope. Recorded in Validation note.

### Summary

Root cause is `discoverWorkflows` counting the verbatim-commission test fixture `fixtures/refit-content-propagation/site-workflow` as a real workflow because `discoverIgnoreDirs` omits `fixtures`/`testdata`. Chose the 1-line prune-set fix over a fixture marker or entity-scoped resolution because it repairs the single shared resolver and thus every auto-discovering command at once. Spiked the fix live (discovery count 2→1) and confirmed the full status test suite stays green before proposing it; flagged the detached adversarial audit for validation.

## Gate: ideation — APPROVED (FO)

- **Verdict:** approved for implementation. One-line prune-set fix, spiked live (discovery 2→1, `go test ./internal/status/...` green), value-measuring ACs (count baseline can move to 0 or 2), ~38 LOC / 2 files. No captain decision required.
- **Validation carries the detached adversarial audit** — touches the front-door launcher + the `status` discovery/mutation path (2 of 4 high-stakes surfaces). AC-provenance sub-trigger does not fire.
- **Implementation base:** worktree off `origin/main` (`ca136f83`), not local `main` (which carries unrelated durable-decisions work).

## Stage Report: implementation

- DONE: discoverIgnoreDirs includes fixtures and testdata; `./spacedock status --discover` from repo root returns EXACTLY docs/dev (AC-1, count 2->1).
  Added `"fixtures": true, "testdata": true` to `discoverIgnoreDirs` (handlers.go, code commit 1f8b666a). Built the binary in the worktree; `status --discover` from the worktree root printed exactly one line ending `/docs/dev`, exit 0 — the `fixtures/refit-content-propagation/site-workflow` row that the pre-prune walk returned is gone.
- DONE: AC-3 behavior test added in internal/status: a temp tree with commissioned READMEs nested under fixtures/ and testdata/ makes discoverWorkflows return exactly [docs/dev]; written to fail pre-prune and pass post.
  `TestDiscoverWorkflowsPrunesFixturesAndTestdata` (discover_worktree_noise_test.go) builds `docs/dev` + commissioned READMEs under `fixtures/…` and `testdata/…`, asserts `discoverWorkflows(repo) == [docs/dev]`. Ran RED first (returned all three subtrees), GREEN after the prune landed. Reverting the two basenames re-surfaces the fixture rows and reds it.
- DONE: go test ./... and go test ./... -race green; change stays within the ~38 LOC / 2 file expected surface (re-gate if it drifts past ~200 LOC or into a github.ref rewrite).
  `go test ./...` — all packages ok; `go test ./... -race` — all packages ok. Diff is 43 insertions / 2 deletions across handlers.go + discover_worktree_noise_test.go (within ~38 LOC ±15, 2 files). No `github.ref` rewrite, no scope drift.

### Summary

Added `fixtures` and `testdata` to `discoverIgnoreDirs` in `internal/status/handlers.go` — the same basename prune the walk already applies to `tests`/`vendor`/`dist`/`build` — so a commissioned-shape README that exists only as a test fixture is no longer counted as a real workflow, repairing every auto-discovering command (`--discover`, `new`, `state commit`, boot `--identify`, front-door banner) at once via the single shared resolver. TDD: wrote the behavior test first, watched it red returning all three subtrees, then landed the two-basename prune to green; only descent into such a child is pruned, so `--workflow-dir`/`--root` pointed directly at one still resolves. Full suite and `-race` both pass; surface is 43 LOC across the two expected files. Confirmed the `DISCOVER_IGNORE_DIRS` comment reference has no live counterpart to keep in sync.

## Stage Report: validation

- DONE: Verify every AC with reproduced evidence: AC-1, AC-2, AC-3, AC-4.
  AC-1: fixed binary `status --discover` from worktree root → EXACTLY `…/docs/dev`, exit 0, 1 line; pre-prune binary (built from `ca136f83` in throwaway) → 2 lines (`docs/dev` + `fixtures/…/site-workflow`). Count 2→1 on an independent moving baseline. AC-2: fixed read `status` (no `--workflow-dir`) exit 0 (resolved `docs/dev`); pre-prune exit 1 with the "multiple commissioned Spacedock workflows found … pass --workflow-dir" refusal — `new`/`state commit` share the same resolver proven by AC-3. AC-3: `TestDiscoverWorkflowsPrunesFixturesAndTestdata` GREEN with prune. AC-4: `go test ./...` and `go test ./... -race` both all-ok, 0 FAIL / 0 DATA RACE.
- DONE: Detached adversarial audit on a THROWAWAY checkout (not this worktree): construct a claim-breaking edit and confirm the AC-3 behavior test reds. 2 of 4 high-stakes surfaces (front-door launcher + status discovery/mutation).
  Throwaway detached worktree at HEAD `1f8b666a` under scratchpad. Three claim-breaking edits, each reds AC-3 for the right reason: (a) full prune revert (`git checkout ca136f83 -- handlers.go`) → returns all 3 subtrees; (b) misspell `fixtures`→`fixturez` (testdata kept) → `fixtures/…/site-workflow` resurfaces; (c) misspell `testdata`→`testdataX` (fixtures kept) → `testdata/…` resurfaces. Test is sensitive to each basename independently — no single-basename hole. CLEAN audit: no test-strength hole found.
- DONE: Semantic adversarial pass over the changed discovery behavior; classify each finding; recommend PASSED or REJECTED with evidence.
  Blast radius: `discoverIgnoreDirs` has one production consumer (`discoverWorkflows`); no parallel `discover_workflows`/`DISCOVER_IGNORE_DIRS` implementation exists (grep) — the "Matches" comments are dead parity refs, no divergence. Start-root always inspected: `status --discover --root …/fixtures/refit-content-propagation` still resolves `site-workflow` (exit 0), so `--root`/`--workflow-dir` overrides reach a fixtures-nested path; only the toplevel auto-walk prunes it. Perf: +2 O(1) map keys, pruning more subtrees is strictly cheaper — no risk. No material or outcome/evidence defect found.

### Deferred risks

- A real workflow deliberately placed under a `fixtures/` or `testdata/` directory is no longer auto-discovered (such a repo could return 0, or omit it from a multi-workflow set). Trigger is outside the current promise — real workflows live at `docs/dev`; this is the identical, documented tradeoff already accepted for `tests`/`vendor`/`dist`/`build`, and the supported path still satisfies AC-1 (returns exactly `docs/dev`) while `--root`/`--workflow-dir` still reaches such a path. Promotes to material only if the project ever legitimately hosts a real, auto-discoverable workflow under a `fixtures`/`testdata`-named dir.

### Summary

Recommend **PASSED**. Both value ACs (AC-1 count 2→1, AC-2 exit-0-vs-refusal) have valid, reproduced evidence on an independent moving baseline; the AC-3 behavior test is falsifiable and proven strong (reds under full-revert and each single-basename misspell); `go test ./...` and `-race` are green. The detached adversarial audit on a throwaway checkout is CLEAN — every claim-breaking edit reds the guard. One deferred risk recorded (real workflow nested under `fixtures`/`testdata`), outside the promised workflow. No material findings block the gate.

## Stage Report: implementation (cycle 2)

Captain re-scope (2026-07-22): relocate fixtures to package-adjacent testdata (Option B, `skills/integration/testdata/`) + prune `testdata`-only. Supersedes the prune-both implementation (commit 1f8b666a, removed) and its validation reported above.

- DONE: Revise the entity's Proposed approach / Acceptance criteria / Test plan / Expected surface to the re-scoped target.
  All four rewritten (plus Spike and Validation note, for consistency). Value AC-1 now couples "no top-level `fixtures/`" with "discovery returns exactly `docs/dev`"; the prior cheapest-first option analysis is preserved in the ideation stage report.
- DONE: Trace each fixture's consumer first (riskiest unknown).
  Grep over all `*.go`/`*.sh`/`*.yml`/`*.md`, `.github`, `go:embed`, Makefile: no fixture has a Go-test consumer; drivers are the `refit` + `first-officer` skills (manual/live). The sole code ref (`fo_write_core_mutation_gate_test.go`) is reference-only. Escalated the no-Go-consumer ambiguity; captain confirmed Option B.
- DONE: Spike the move mechanism before executing all.
  Adapted (no consumer test exists to "run green"): zero-cost pre-check that every fixture frontmatter parses cleanly under the migration-check yaml.v3 consistency walk, then ran the three newly-exposed packages (`internal/status`, `skills/integration`, `internal/contractlint`) immediately post-move before the full suite. All green.
- DONE: Execute the move + `testdata`-only prune + TDD the behavior test + update references.
  `git mv fixtures/{refit-content-propagation,entity-label-drive}` → `skills/integration/testdata/` (7 pure renames), removed the empty `fixtures/`. Added `"testdata": true` to `discoverIgnoreDirs` (`handlers.go`); the `fixtures` prune was never on this base (origin/main), so nothing to drop. `TestDiscoverWorkflowsPrunesTestdata` watched RED→GREEN. Updated the FO-write classifier example to `skills/integration/testdata/entity-label-drive/README.md` (still `blocked-product` via `skills/**`). Left the `docs/roadmap` debrief's historical "untracked `fixtures/refit-content-propagation/`" line unchanged — an accurate past-incident record, not a live pointer, and `docs/roadmap` is pruned from the migration check + outside the instruction-surface walk.
- DONE: go test ./... + -race green; fixture-consuming tests green; diff review confirms `fixtures/` gone; value AC.
  `go test ./...` → 0 FAIL. `go test ./... -race` → 0 FAIL, 0 DATA RACE (15 packages ok). `status --discover` from the worktree root → exactly `…/docs/dev`, exit 0, 1 line. `test ! -d fixtures` → true. Surface: 43 insertions / 3 deletions across 3 code files + 7 renames (within tolerance).

### Summary

Per the captain's 2026-07-22 re-scope, eliminated the top-level `fixtures/` directory by relocating its three fixtures to `skills/integration/testdata/` — the one Go test package among the skills, already excluded from the contractlint instruction-surface walk — and pruned only the `testdata` basename from `discoverWorkflows`. The riskiest unknown, that the fixtures have no Go-test consumer, was traced and escalated; the captain chose the `skills/integration/testdata/` home. Confirmed the relocated frontmatter parses under the migration-check walk (a new exposure since it walks `skills/`). TDD'd the `testdata` prune (RED→GREEN); full suite and `-race` green; discovery returns exactly `docs/dev` and the top-level `fixtures/` is gone.
