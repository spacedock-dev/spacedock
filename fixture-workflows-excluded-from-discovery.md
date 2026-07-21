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

Fix the root cause in discovery, cheapest-first evaluation:

- **(a) Chosen — add `fixtures` and `testdata` to `discoverIgnoreDirs`.** One line in `internal/status/handlers.go`, plus a comment noting these are test-scaffolding basenames (`testdata` is the Go-tool-ignored convention; `fixtures` is this repo's), matching the philosophy the set already applies to `tests`, `vendor`, `dist`, `build`. Fixing `discoverWorkflows` fixes **every** consumer above at once — the walk simply does not descend into a `fixtures/` or `testdata/` subtree. Cheapest-sufficient because it targets the single shared root and needs no per-call-site or per-fixture work.
  - What it costs: a real workflow deliberately placed *under* a directory named `fixtures/` or `testdata/` would no longer be discovered. Accepted — no real workflow lives there (real ones live at `docs/dev`), and this is the identical tradeoff already made for `tests`. The start-root of a walk is always inspected regardless of its own basename (only descent into a `fixtures`/`testdata` *child* is pruned), so pointing `--workflow-dir`/`--root` directly at such a path still works; only the auto-walk skips them.

- **(b) Rejected — a marker field on the fixture README that discovery respects.** More expensive *and* self-defeating: the marker would have to go on `site-workflow/README.md`, whose entire reason for existing is to be a *verbatim* commission output; adding a `fixture: true` field corrupts that faithfulness and breaks the refit-content-propagation test the fixture serves. A marker on the parent doc README does not help — discovery reads each directory's own README frontmatter, not an ancestor's. It also adds a new discovery predicate plus a standing "every future fixture author must remember the marker" failure mode (forget it → the bug returns silently).

- **(c) Rejected — entity-scoped commands resolve the workflow from the entity's state checkout.** Narrower *and* more expensive. It only helps `state commit`, which has an existing entity in a state checkout to resolve *from*. `new` creates a new entity — there is nothing to resolve from — so it stays broken, as do `status --discover`, `status --boot --identify`, and the front-door banner. It leaves the root defect (a fixture counts as a workflow) in place, so the bug resurfaces at every non-entity-scoped call site, for more code than option (a).

No `--workflow-dir` / `--root` override behavior changes. The current clear-error-on-ambiguity path stays exactly as is (it is correct and out of scope).

**Documentation:** no doc diff required. `discoverIgnoreDirs` is an internal implementation detail not enumerated in any user-facing doc (the docs-site describes `spacedock status --workflow-dir docs/dev` and boot as single-workflow; grep of `docs/` and `skills/` finds the prune list only in archived state files and roadmap notes, none rendered). This change *restores* the documented single-workflow auto-discovery, so no user-visible command surface or output wording changes.

## Out of scope

- The ambiguity error path itself (clear message + non-zero exit) — already correct, verified live.
- Any `--workflow-dir` / `--root` / `PIPELINE_DIR` override behavior.
- The FO's shell-piping lesson from the session note.
- Removing, moving, or re-marking the `refit-content-propagation` fixture.

## Acceptance criteria

Each AC names a property of the finished task, not a stage action, and how it is verified.

**AC-1 (VALUE) — In this repo, downward discovery returns EXACTLY the one real workflow.**
The discovered-workflow set for the repo toplevel is exactly `{<repo>/docs/dev}` (count 1), down from the current 2 (`docs/dev` + `fixtures/refit-content-propagation/site-workflow`). The count is the independent baseline that can move the wrong way: too-aggressive → 0, ineffective → 2.
Verified by: `./spacedock status --discover` from the repo root prints exactly one line, `<repo>/docs/dev`, exit 0 (currently prints two lines). A test-strength edit that reverts the prune line makes this red.

**AC-2 (VALUE) — The auto-discovering commands succeed without `--workflow-dir` in this one-workflow repo.**
`new` and `state commit` (and the read `status`) resolve `docs/dev` by discovery alone, instead of exiting non-zero with the multi-workflow "pass --workflow-dir" refusal.
Verified by: from the repo root with no `--workflow-dir`, `./spacedock status` (read) exits 0 and lists `docs/dev` entities; a `new`/`state commit` resolution over a fixture-bearing tree resolves to the real workflow — exercised by the discovery behavior test in AC-3, since `discoverWorkflows` is the shared resolver all three call. Failure mode that would make it red: discovery again returns ≥2, re-triggering the refusal.

**AC-3 — A behavior test proves the fixtures/testdata prune excludes a nested commissioned README while still finding the real workflow.**
A Go test builds a temp tree with a real workflow (`docs/dev/README.md`, commissioned) plus a commissioned README nested under `fixtures/…` and under `testdata/…`, calls `discoverWorkflows(repo)`, and asserts the result is exactly `[<repo>/docs/dev]`. It is the falsifiable guard: reverting the prune line re-introduces the fixture rows and reds the test.
Verified by: the new test in `internal/status/discover_worktree_noise_test.go` (sibling of `TestDiscoverWorkflowsSkipsNestedCheckout`); `go test ./internal/status/...` green with the prune, red without it.

**AC-4 — The full Go suite stays green with the prune applied.**
Adding `fixtures`/`testdata` to the prune set breaks no existing test (in particular the `testdata/*-workflow` fixtures many status tests point `--root`/`--workflow-dir` at directly, which the prune does not touch — those are start-roots, not descended children).
Verified by: `go test ./...` and `go test ./... -race` pass. (Pre-verified during ideation: `go test ./internal/status/...` passed with the throwaway prune applied.)

## Test plan

- **Primary (behavior, cheap, ~35 LOC):** a Go unit test alongside `internal/status/discover_worktree_noise_test.go`, same shape as `TestDiscoverWorkflowsSkipsNestedCheckout` — write a temp tree with the real `docs/dev` workflow plus commissioned READMEs nested under `fixtures/` and `testdata/`, assert `discoverWorkflows(repo)` returns exactly `[docs/dev]`. This is written first and watched fail (returns the fixture rows) before the prune line lands. This one test covers AC-1/AC-2/AC-3 because `discoverWorkflows` is the single resolver behind `--discover`, `new`, `state commit`, boot `--identify`, and the front-door banner — driving each command separately would re-prove the same resolver, so the simplest sufficient proof is one test at the resolver plus the real-repo `--discover` check.
- **Real-repo check (AC-1, CLI, near-zero cost):** `./spacedock status --discover` from the repo root returns the single `docs/dev` line, exit 0.
- **Regression (AC-4):** `go test ./...`, then `go test ./... -race`. No live-workflow test needed — the claim is discovery filtering, not runtime integration.

The only new mechanism is the two added prune basenames; the value AC it serves is AC-1 (discovered count 2→1). The simplest alternative to a behavior test — asserting the source contains the strings `"fixtures"`/`"testdata"` — is insufficient: it is a static prose check that passes even if the map key is misspelled or the walk never consults it, so it cannot fail for the right reason. Exercising `discoverWorkflows` over a real tree is the proof.

## Expected surface

- `internal/status/handlers.go` — +1 changed line (the two basenames on the existing `discoverIgnoreDirs` line) plus ~2 comment lines: ~3 LOC.
- `internal/status/discover_worktree_noise_test.go` — one appended `Test…` function: ~35 LOC.

Total ~38 LOC across 2 files, almost all test. Tolerance ±15 LOC and ±1 file (if the test lands as a new `discover_fixtures_prune_test.go` with its own ABOUTME header instead of appending, add ~4 header LOC and one file).

## Spike (riskiest unverified mechanism)

Riskiest unverified claim: the chosen filter actually excludes the fixture while still finding the real workflow, and does not break the existing discovery tests. **Spiked and proven:**

1. Reproduced the defect against the real binary: `./spacedock status --discover` from the repo root returns 2 workflows (`docs/dev` + `fixtures/refit-content-propagation/site-workflow`), exit 0.
2. Applied the throwaway edit (`"fixtures": true, "testdata": true` on the `discoverIgnoreDirs` line), rebuilt `./cmd/spacedock`, re-ran `status --discover` → returned EXACTLY `/…/docs/dev` (count 2→1). Reverted the edit (`git checkout --`), tree clean.
3. With the same throwaway edit applied, `go test ./internal/status/...` → PASS (16.3s) — the `testdata/*-workflow` fixtures are start-roots the walk inspects directly, not descended `testdata` children, so the prune does not touch them. Reverted.

The throwaway tree seeds the AC-3 test. Relied-on proven mechanisms: `discoverWorkflows`'s existing `ignore[ent.Name()]` basename prune (already the mechanism for `tests`/`vendor`/`dist`) and the always-walk-the-start-root behavior.

## Validation note

The **detached adversarial audit trigger applies**: this change touches two of the four high-stakes surfaces — the **front-door launcher** (`frontdoor.go` consumes `DiscoverWorkflows`) and the **`status` mutation/guard/discovery path** (`ResolveWorkflowDir` backs `--set` mutation resolution, the merge guard, and `new`/`state commit`). Validation should run the read-only audit on a throwaway checkout: construct an adversarial edit (e.g., revert the prune line, or misspell a basename) and confirm the AC-3 behavior test reds. The AC-provenance sub-trigger does NOT fire — the value AC's expected value (count 1, the known real `docs/dev`) is an independent constant, not derived from the package's own production functions or constants.

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
