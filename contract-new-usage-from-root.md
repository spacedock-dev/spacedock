---
id: r0n903prs88q29vjkmcmtw2x
title: FO contract documents the from-root `spacedock new` invocation (--workflow-dir) so filing an entity needs no trial-and-error
status: implementation
source: "captain (2026-06-14, this session) — filing an entity from the project root hit no-Spacedock-workflow-here / pass --workflow-dir. FO Write Scope (first-officer-shared-core.md:171) and Filing New Entities (claude-first-officer-runtime.md:43) both document `spacedock new <slug> [--folder] [--id-seed S --id-actor A]` WITHOUT --workflow-dir, but the Working Directory rule keeps the FO at the project root where `new` cannot auto-discover the workflow. `new --help` returns the general menu, not per-command usage, so the FO must trial-and-error the flags. Same fo-efficiency / friction-reduction class as the 0.20.4 structured-reads and 0203 work."
started: 2026-06-15T00:32:39Z
completed:
verdict:
score: 0.30
worktree: .worktrees/spacedock-ensign-contract-new-usage-from-root
issue:
sprint: 0203-fo-efficiency
---

The documented `spacedock new` invocation is incomplete for the FO's standing position. An FO at the project root (per the Working Directory rule) must pass `--workflow-dir`, but the contract's `new` examples omit it, so the first filing attempt fails. Close the gap so a fresh FO files correctly on the first try.

## Problem

Two contract sites document filing:
- `first-officer-shared-core.md:171` (FO Write Scope): "seed task creation via `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < stub`".
- `claude-first-officer-runtime.md:43` (Filing New Entities): "Use `spacedock new <slug> [--folder] [--id-seed S --id-actor A]` via Bash".

Neither shows `--workflow-dir {workflow_dir}`. The Working Directory rule ("Stay at the project root") means the FO runs `new` from root, where the binary errors `no Spacedock workflow here — pass --workflow-dir`. `new --help` prints the top-level command menu rather than per-command usage, so the only way to learn the required flag is to fail and read stderr.

## Spike: what the binary actually does (riskiest unknown, exercised first)

Built current source (`go build -o /tmp/spacedock-current ./cmd/spacedock`) and exercised the real paths. The installed brew binary is 0.20.2; current source is ahead, so all findings below are against current source.

1. **`new` already auto-discovers — by walking UP, which cannot reach a nested workflow.** Current source wires `new` as a pure alias for `status --new` (`cli.go:315`), and the runner already runs discovery when no `--workflow-dir` is given (`native_runner.go:75`). But that discovery is `DiscoverWorkflowDir` (`discover_walkup.go:16`), which walks UP from the cwd to an *enclosing* commissioned workflow. The workflow here lives DOWN at `docs/dev/` (its README declares `commissioned-by: spacedock@…` + `state: .spacedock-state`); the project root README has no such frontmatter. So from the project root the upward walk finds nothing and errors `no commissioned Spacedock workflow found in /…/spacedock-v1 …`. Reproduced: `printf … | /tmp/spacedock-current new <slug>` from root errors; the same from `docs/dev/` SUCCEEDS (`created: …/.spacedock-state/<slug>.md id=…`).

   This invalidates option (a) **as originally framed**. `status --discover` does NOT use the upward walk — it uses `discoverWorkflows(root)` (`handlers.go:478`), a DOWNWARD scan from the repo toplevel. Reproduced: `status --discover` from the project root prints `…/docs/dev`. So "make `new` discover the same way `status --discover` does" means adding the downward scan as a fallback, not reusing the upward walk-up. The mechanism exists and is reusable.

2. **`new --help` prints the general menu — root cause found.** `root.SetHelpFunc(printHelp)` (`cli.go:133`) is inherited by every subcommand in cobra, so `new --help` → `cmd.Help()` (`cli.go:323`) → the inherited func → `printHelp(stdout)`, the top-level menu. Reproduced: `new --help` emits the Launch/Setup/Workflow menu, no `--workflow-dir/--folder/--id-seed/--id-actor` flag surface. The captain's literal friction.

## Proposed approach

Two small, independent fixes; ship both.

**Fix A (behavioral, the real one) — `new` falls back to a downward scan when the upward walk misses.** In `native_runner.go` after `DiscoverWorkflowDir(dir)` returns `!ok` (currently the error at line 79), before erroring, run `discoverWorkflows(repoRoot)` where `repoRoot` is `git rev-parse --show-toplevel` from `dir` (matching `runDiscover`'s default root at `native_runner.go:~525`). If it returns exactly ONE workflow, use it as `pipelineDir`. If zero or more-than-one, keep an error — but reword it to point at the resolved candidates / `--workflow-dir`. This is scoped to the no-flag path only (explicit `--workflow-dir`/`PIPELINE_DIR`/`--root` still short-circuit, unchanged). After this, `spacedock new <slug>` from the project root SUCCEEDS with no `--workflow-dir`, and BOTH contract examples stay correct as written (no `--workflow-dir` needed). The ambiguity guard (≥2 workflows) preserves the existing "report and stop, do NOT search" discipline by refusing to silently guess.

**Fix C (help text) — `new --help` prints the command's own usage.** Give the `new` command a `Long`/usage string carrying its flag surface and stop the inherited root help func from swallowing it. Smallest form: in `newNewCommand`, on `wantsHelp(args)` print a dedicated per-command usage block (the `new [--folder] SLUG` synopsis plus `--workflow-dir`, `--folder`, `--id-seed`, `--id-actor` with one-line descriptions) instead of `cmd.Help()`. This is the FO's fallback when discovery is genuinely ambiguous (≥2 workflows) and `--workflow-dir` is required.

**Contract docs:** With Fix A landed, the from-root single-workflow case needs no `--workflow-dir`, so the two contract `new` examples (`first-officer-shared-core.md:171`, `claude-first-officer-runtime.md:43`) stay correct as-is. Add ONE clause to each noting that when the repo holds more than one commissioned workflow, `new` reports the candidates and the FO passes `--workflow-dir {workflow_dir}` (discoverable via `new --help` after Fix C). This documents the real residual requirement without making `--workflow-dir` look mandatory in the common case.

## Out of scope

- Generalizing the help fix to every subcommand. Scope Fix C to `new`. (If the same `SetHelpFunc`-inheritance bug is trivially fixed once for all subcommands, implementation may, but the AC only binds `new`.)
- Changing `status --discover` itself — it already works; Fix A reuses its scan.
- Multi-workflow auto-selection heuristics. Two-or-more is an explicit error that requires `--workflow-dir`; we do not guess.

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 — A from-root `spacedock new <slug>` SUCCEEDS in a single-workflow repo with no `--workflow-dir`.** Run from the repo toplevel (the FO's standing directory) where the only commissioned workflow is a nested subdir, `new <slug>` reading a valid body from stdin creates the entity file at the workflow's entity dir and exits 0.
Verified by: a behavior/CLI test (Go, in `internal/status`) that builds a fixture repo whose toplevel README is NOT a workflow and whose nested subdir IS one, runs the `--new` path from the toplevel with no `--workflow-dir`, and asserts (a) exit 0, (b) the `created: <path> id=<id>` line, (c) the file exists under the nested workflow's entity dir. Reds today (current source errors `no commissioned Spacedock workflow found`). This is the riskiest path and must red before the fix.

**AC-2 — A from-root `spacedock new <slug>` in a MULTI-workflow repo refuses and names the candidates.** With two-or-more commissioned workflows under the toplevel, `new <slug>` from root exits non-zero with a message listing the candidate workflow dirs and instructing `--workflow-dir`; passing `--workflow-dir <one>` then SUCCEEDS.
Verified by: a Go CLI test with a two-workflow fixture asserting the ambiguity error (non-zero exit + both candidate paths in stderr), plus that `--workflow-dir <dir>` resolves to exit 0 + `created:`. Guards against silently guessing.

**AC-3 — `spacedock new --help` prints the command's own usage and full flag surface.** Its stdout shows the `new [--folder] SLUG` synopsis and the flags `--workflow-dir`, `--folder`, `--id-seed`, `--id-actor` — not the top-level Launch/Setup/Workflow menu.
Verified by: a Go CLI test asserting `new --help` stdout contains the per-command synopsis and each flag name, AND does NOT contain a general-menu marker (e.g. the `Launch` section header). Command-output assertion, not a contract prose-grep.

**AC-4 — Both contract `new` examples are accurate for the post-fix behavior.** `first-officer-shared-core.md:171` and `claude-first-officer-runtime.md:43` show the no-`--workflow-dir` invocation for the common single-workflow case, plus a one-clause note that a multi-workflow repo requires `--workflow-dir {workflow_dir}`.
Verified by: implementation applies the documented doc diff (below); accuracy is anchored to AC-1/AC-2 behavior, so the prose cannot drift from the binary without a behavior test reddening. (No standalone prose-grep AC.)

## Documentation changes (doc diff — implementation applies)

**`skills/first-officer/references/first-officer-shared-core.md:171`** — append one clause after the existing `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < stub`:
> … `< stub` (runs from the project root; `new` discovers the single commissioned workflow automatically — if the repo holds more than one, `new` reports the candidates and you pass `--workflow-dir {workflow_dir}`).

**`skills/first-officer/references/claude-first-officer-runtime.md:43`** — same one-clause addition after the existing `spacedock new <slug> [--folder] [--id-seed S --id-actor A]`:
> … via Bash from the project root; `new` auto-discovers the lone workflow, else pass `--workflow-dir {workflow_dir}` (see `spacedock new --help`).

Exact final wording is pinned at implementation against the resolved file lines; the intent above is the contract.

## Test plan

- **Fix A (AC-1, AC-2):** Go behavior tests in `internal/status` (sibling to `native_new_test.go`, `native_discover_test.go`). Fixtures: a tmp-dir repo with a non-workflow toplevel README + one (AC-1) or two (AC-2) nested workflow READMEs (`commissioned-by: spacedock@…`). Drive the `--new` path through the same `dispatch(...)` harness the existing native tests use; assert exit code, `created:`/error stdout/stderr, and resulting on-disk file. The downward-scan reuse (`discoverWorkflows` + `git rev-parse --show-toplevel`) is already covered for `--discover`; these tests cover the new fallback wiring in `new`. Cost: low (unit-level, no live workflow).
- **Fix C (AC-3):** Go CLI test (sibling to `native_usage_test.go` / `help_test.go`) asserting `new --help` stdout content. Cost: trivial.
- **AC-4:** No standalone test — anchored to AC-1/AC-2 so prose accuracy is enforced by behavior. Implementation applies the doc diff.
- No live workflow smoke test needed: all claims are command-level and covered by binary-driven fixtures. Riskiest path (downward-scan fallback in `new`) is exercised first by AC-1, which reds against current source.

**No spike needed beyond what's recorded above:** the downward-scan mechanism (`discoverWorkflows`), the repo-root resolution (`git rev-parse --show-toplevel`), the `new`→`status --new` alias, and the `SetHelpFunc` inheritance root-cause were all exercised against the freshly built current-source binary (see Spike section).

## Stage Report: ideation

- DONE: Decide the fix and weigh options (a)/(b)/(c); recommend ONE; exercise riskiest unknown (can `new` reuse the discover path) first.
  Spike against freshly-built current source: `new` already runs upward discovery (`native_runner.go:75` via `DiscoverWorkflowDir`) but that walks UP and cannot reach the nested `docs/dev` workflow; `status --discover` uses a separate DOWNWARD scan (`discoverWorkflows`, `handlers.go:478`). Recommended Fix A (downward-scan fallback in `new`) + Fix C (`new --help`); option (b)-prose-only rejected as it papers over a real behavioral gap.
- DONE: Concrete before/after — name the binary change for (a) and the help wording for (c); state which contract examples change.
  Fix A: in `native_runner.go` no-flag miss branch, fall back to `discoverWorkflows(git-toplevel)`, use the sole match, error+name candidates on ≥2. Fix C: `newNewCommand` emits a per-command usage block instead of inherited `printHelp`. Contract examples at `first-officer-shared-core.md:171` + `claude-first-officer-runtime.md:43` stay as-is (no `--workflow-dir` in common case) + one multi-workflow clause each — doc diff recorded in body.
- DONE: AC bound to BEHAVIOR (from-root `new` SUCCEEDS, exit 0 + created file, reds if broken); help-text AC asserts `new --help` per-command output; never a contract substring-grep.
  AC-1 from-root single-workflow success (reds today), AC-2 multi-workflow refusal+resolution, AC-3 `new --help` command-output assertion (incl. negative: no general-menu marker), AC-4 doc accuracy anchored to AC-1/AC-2 not a prose-grep. Test plan: Go fixtures in `internal/status` via the existing `dispatch(...)` harness.

### Summary

Reframed the captain's lean: option (a) is viable but NOT via the upward walk-up the dispatch implied — `new` already walks up and can't reach a nested workflow; the real fix reuses `status --discover`'s downward `discoverWorkflows` scan as a fallback (Fix A), gated to a single match with an ambiguity error otherwise. Also root-caused the `new --help` menu bug to cobra's inherited `SetHelpFunc` (Fix C). Both small, independent, behavior-tested; contract examples stay correct as written plus one multi-workflow clause. Riskiest mechanism (downward-scan fallback) exercised first against a freshly built binary; AC-1 reds against current source.

## Stage Report: implementation

- DONE: Fix A (behavioral) — no-flag miss branch falls back to discoverWorkflows(git-toplevel); 1 match resolves, 0 keeps terminal no-workflow msg, ≥2 refuses naming candidates + --workflow-dir; scoped to no-flag path only.
  `discoverWorkflowDownward` in native_runner.go, wired at the `DiscoverWorkflowDir(dir)` !ok branch (commit 3bbf3698). Explicit --workflow-dir/PIPELINE_DIR/--root still short-circuit unchanged.
- DONE: AC-1 from-root single-workflow `new` SUCCEEDS (exit 0 + `created:` + file under nested entity dir; red against current source first).
  TestNewFromRootSingleWorkflowSucceeds in native_new_from_root_test.go; confirmed red (no-workflow error) before fix, green after.
- DONE: AC-2 multi-workflow refusal naming candidates, then --workflow-dir succeeds.
  TestNewFromRootMultiWorkflowRefuses; asserts both realpath'd candidate dirs + --workflow-dir in stderr, non-zero exit, then exit 0 + created on disambiguation.
- DONE: Fix C (help) — newNewCommand emits a per-command usage block instead of the inherited root help func.
  setNewHelp in help.go (SetHelpFunc on the command, mirroring setSetupHelp); wired in newNewCommand (cli.go).
- DONE: AC-3 `new --help` carries the per-command synopsis + each flag AND not the general-menu marker.
  TestNewHelpRendersCommandUsage in new_help_test.go; asserts synopsis + 4 flags present, "Launch" header absent.
- DONE: Doc diff — one-clause multi-workflow note appended to first-officer-shared-core.md + claude-first-officer-runtime.md.
  Examples stay no---workflow-dir for the common single-workflow case; AC-4 anchored to AC-1/AC-2 behavior.
- DONE: `go test ./...` green.
  Full suite green; one transient flake (ensigncycle streamwatch replay, unrelated to this diff) passed 3/3 standalone and on re-run.

### Summary

Implemented Fix A (downward-scan fallback in the no-flag `new`/status path, reusing `discoverWorkflows` + `git rev-parse --show-toplevel`) and Fix C (per-command `new --help`), plus the two contract doc clauses. AC-1 reddened against current source before the fix as required; AC-1/AC-2/AC-3 are Go fixture/CLI tests via the existing dispatch/Run harnesses. The single full-suite failure was a pre-existing timing flake in `internal/ensigncycle`, not touched by this diff.
