---
title: "Make the filing live-test fixture README representative of a real workflow (add a Task Template)"
status: validation
source: "Captain decision 2026-07-10 on the science-officer dp-premise finding. Verbatim: then fix the fixture's readme to be realistic. new --help is not a priority. Finding it acts on (science-officer-verified): filingReadme() at internal/ensigncycle/shared_fixtures_test.go:313-332 carries no Task Template and no copyable stub (grep for template over it is empty), so dp's AC-1 filing-hunt methodology is structurally blind to a README-based filing approach — it measures a stripped fixture. The real docs/dev/README.md carries a full column-0 pipe-safe ## Task Template. And the FO filing contract in claude-first-officer-runtime.md's ## Filing New Entities never points the FO at the workflow README before filing (grep across skills/first-officer for a README/template-consultation instruction returns zero)."
started: 2026-07-10T12:52:14Z
completed:
verdict:
score: 0.4
worktree: .worktrees/spacedock-ensign-filing-fixture-readme-realistic
issue:
id: c3wxhq3qj94mhakam80g4zxw
---

## Problem

The `filing` live-scenario fixture README is too small to represent a commissioned workflow. `filingReadme()` in `internal/ensigncycle/shared_fixtures_test.go` contains workflow frontmatter, a title, fixture prose, and two stage sections. It contains no `## Task Template` and no copyable entity body. By contrast, `docs/dev/README.md` gives filers a full Task Template. A live filing measurement against the stripped fixture cannot reveal whether a first officer uses the workflow-local template that exists in a real workflow.

The Claude filing contract has the matching blind spot. `skills/first-officer/references/claude-first-officer-runtime.md` teaches `spacedock new`, stdin, and the `--next-id` warning, but it never points the first officer to the workflow README's `## Task Template` at filing time.

## Proposed approach

Use the smallest complete change: enrich the fixture and add one filing-time contract pointer.

### 1. Add a realistic fixture Task Template

Append this section to the body returned by `filingReadme()`, after the stage descriptions. Keep the existing workflow frontmatter, fixture prose, and stage sections unchanged.

````markdown
## Task Template

```yaml
---
title: Task name here
status: backlog
---

Brief description of this task and what it aims to achieve.

## Problem

{What is broken or missing, and why it matters. Ideation fills this in.}

## Proposed approach

{How the task intends to solve the problem. Ideation fills this in.}

## Out of scope

{What this task deliberately does not cover, so the boundary is explicit.}

## Acceptance criteria

{Each criterion names an end-state property and how it is verified.}

## Test plan

{What verifies the implementation, its cost, and whether fixture, CLI, or live tests are needed.}
```
````

The outer fence above only displays the proposed README bytes in this entity. The fixture output contains one YAML code fence. Inside that fence, both `---` lines and every frontmatter field begin at column zero. The stdin stub includes `title` and the workflow's initial `status: backlog`; it omits `id` so `spacedock new` can mint and stamp the sequential ID. The body mirrors the canonical development template's Problem, Proposed approach, Out of scope, Acceptance criteria, and Test plan sections.

### 2. Add the minimal filing-time pointer

In `skills/first-officer/references/claude-first-officer-runtime.md`, add one sentence at the start of `## Filing New Entities`:

> Before filing, read the workflow README's `## Task Template`; use its frontmatter and section scaffolding as the starting shape for the entity body you pipe to `spacedock new`.

This sentence loads the README body only when filing consumes it. It does not add startup work, change `spacedock new`, or extend help output.

### Alternatives considered

1. **Recommended: fixture template plus filing-time pointer.** This makes the fixture representative and gives the new section a contract consumer. It adds one fixture section and one runtime sentence.
2. **Fixture template only.** This fixes structural fidelity but leaves the first officer with no instruction to consult the template. The live fixture would contain realistic information that the filing contract still ignores.
3. **Share one template constant across the fixture, `spacedock new --help`, and commission templates.** This reduces textual duplication but couples workflow-specific scaffolding to a deliberately small generic CLI stub. It also extends the deprioritized help-template work and exceeds this task's scope.

## Out of scope

- Extending or reprioritizing the `spacedock new --help` stub from PR #491.
- Sharing or generating templates across the CLI, commissioned workflows, and fixtures.
- Changing filing prompts, live assertions, entity ID semantics, or the `spacedock new` command.
- Claiming a filing-hunt-rate improvement. PR #491 already measured `filing` at 0/3 hunts, so this task has no live behavioral headroom.

## Acceptance criteria

**AC-1 — The filing fixture exposes the workflow-local entity shape that a real development workflow exposes.** The on-disk README produced by `writeFilingWorkflow` contains one `## Task Template` whose fenced body starts with a column-zero `---`, declares `title` and `status: backlog`, omits `id`, and includes `## Problem`, `## Proposed approach`, `## Out of scope`, `## Acceptance criteria`, and `## Test plan`. The independent baseline is `docs/dev/README.md`'s Task Template section. The value measure can move the wrong way: current `main` exposes zero fixture Task Templates; the finished fixture exposes exactly one complete template. *Verified by:* `TestFilingReadmeTaskTemplateRoundTripsThroughNew` reads the actual `README.md` written by the fixture and asserts the count and shape above.

**AC-2 — The literal fixture template is pipe-safe through the production creation path.** Piping the fenced Task Template body read from the fixture's on-disk README into `status.NativeRunner` with `--new from-template` exits 0, prints `created:`, creates `from-template.md`, and stamps a non-empty sequential `id`. Moving the opening fence off column zero makes the same production path return the `no frontmatter` error. *Verified by:* the round-trip and negative-control cases in `TestFilingReadmeTaskTemplateRoundTripsThroughNew`. The current fixture fails before invocation because it has no Task Template, so the test is red on current `main`.

**AC-3 — The Claude filing contract points to the workflow Task Template at the moment of filing.** `## Filing New Entities` begins with the exact sentence specified above. The pointer remains filing-local; Startup stays unchanged. *Verified by:* focused diff review at the ideation/validation gate plus `git diff --check`. Do not add a semantic phrase-presence test: `internal/contractlint/doc_test.go` expressly forbids prose-grep tests that substitute bytes for instruction meaning.

## Test plan

- Add `internal/ensigncycle/filing_readme_template_test.go`. `TestFilingReadmeTaskTemplateRoundTripsThroughNew` calls `writeFilingWorkflow`, reads the resulting `README.md`, extracts the sole fenced block under `## Task Template`, and drives `status.NativeRunner.Run` with that literal block as stdin. It asserts the complete shape, successful creation output, the written entity, and its minted ID. A subtest uses a fresh fixture root, indents the opening fence, and expects the production `no frontmatter` failure. Cost: offline, deterministic, sub-second.
- Run `go test ./internal/ensigncycle -run 'Test(FilingReadmeTaskTemplateRoundTripsThroughNew|SharedScenarioFixturesAreDiscoverable)' -count=1` to prove the new value and preserve PR #490's fixture-discovery guard.
- Run `go test ./internal/ensigncycle ./internal/status ./internal/contractlint -count=1` for affected-package regression coverage. `internal/status` covers the reused native `new` seam; `internal/contractlint` confirms the instruction quarantine remains green without a prose-grep addition.
- Run the repository gates before completion: `gofmt -w ./cmd ./internal`, `go test ./...`, and `go test ./... -race`.
- Skip live-model lanes. The task claims fixture fidelity and literal parser round-trip behavior, both proven offline. PR #491's live result leaves no filing-hunt-rate delta to measure.

## Exact files and sequencing

Implementation changes exactly these files:

- `internal/ensigncycle/shared_fixtures_test.go` — append the Task Template bytes to `filingReadme()` after the stage sections.
- `internal/ensigncycle/filing_readme_template_test.go` — add the on-disk extraction, structure assertions, native `new` round trip, and indented-fence negative control.
- `skills/first-officer/references/claude-first-officer-runtime.md` — add the one filing-time pointer sentence.

Sequence the implementation after PR #490. PR #490 is merged, but the local `main` inspected during ideation predates it; update or rebase onto current `origin/main` before editing `shared_fixtures_test.go`. PR #490 changes `filingPrompt` to accept `workflowRoot` and adds the shared fixture discovery test, but it does not change `filingReadme()`'s body. Preserve both changes.

PR #491 remains open and edits `claude-first-officer-runtime.md`'s filing paragraph. Sequence after PR #491 merges, or rebase the implementation branch onto PR #491 before editing that paragraph. Keep PR #491's generic CLI stub unchanged.

## Spike determination

No new parser spike is needed. PR #491 commit `205ca99a` already proves the exact risky mechanism: `TestPrintedStubTemplateBlockIsPipeSafe` extracts literal printed template bytes and pipes them through the native `runNew` path. Ideation re-ran that focused test on PR #491; it passed. `status.NativeRunner` exposes the same path to `internal/ensigncycle` through `status.Request`. This task changes only the source of the bytes from help output to the fixture's on-disk README. The new test is the implementation's first red test: current `main` has no `## Task Template`, so extraction fails before creation.

## Documentation impact

The exact skill-text change appears above. No docs-site change is needed because the CLI command, help output, and public command reference remain unchanged. The workflow README fixture is test data; the runtime contract is the user-facing instruction source being changed.

## Stage Report: ideation

- DONE: Specify a realistic, column-zero, pipe-safe workflow Task Template and a round-trip test against the filing fixture's actual README output
  The proposed template declares `title` and `status: backlog`, omits `id`, carries the five development-workflow body sections, and keeps both frontmatter fences at column zero. The focused test writes the fixture, rereads its actual `README.md`, extracts the literal fenced body, and drives `status.NativeRunner --new`; an indented-fence subtest proves the guard can red.
- DONE: Specify the minimal filing-time pointer without extending the deprioritized `spacedock new --help` work
  The exact one-sentence addition sits at the start of `## Filing New Entities` and directs the first officer to the workflow README's Task Template only when filing. `spacedock new`, its help, and its generic stub remain out of scope.
- DONE: Name exact files, sequencing after PR #490, focused tests, and measurable acceptance criteria
  The design names three files, requires a rebase/update past merged PR #490, preserves its `workflowRoot` and discovery changes, and sequences the overlapping contract edit after or onto open PR #491. AC-1 measures complete fixture templates from 0 to 1; AC-2 exercises literal bytes through production creation; AC-3 uses focused contract diff review because repository policy forbids semantic prose-grep guards.
- DONE: Compare alternatives and settle the smallest complete design
  Chosen: fixture plus filing-time pointer. Rejected fixture-only because the contract would still ignore the new section; rejected a shared template constant because it couples workflow-specific scaffolding to generic CLI help and expands the captain-deprioritized work.
- DONE: Resolve the riskiest mechanism before handoff
  Re-ran PR #491's `TestPrintedStubTemplateBlockIsPipeSafe` at commit `205ca99a`; it passed. The new test reuses the same native path and changes only the source to the fixture's on-disk README.

### Summary

Designed a three-file, offline-proven change that makes the filing fixture resemble a real development workflow and teaches the Claude first officer to use that workflow-local shape at filing time. The literal on-disk README template must round-trip through the native `new` path, mint an ID, and fail when its opening fence is indented. The work starts after merged PR #490 and after or atop open PR #491; it does not extend `spacedock new --help` or spend a live lane.

## Stage Report: implementation

- DONE: Write the failing on-disk fixture Task Template round-trip test first, including the indented-fence red control and preservation of PR #490 discovery behavior.
  `TestFilingReadmeTaskTemplateRoundTripsThroughNew` first failed because the written README had zero Task Templates, then passed with its literal native-new round trip and indented-fence rejection; the focused run also kept `TestSharedScenarioFixturesAreDiscoverable` green.
- DONE: Implement exactly the realistic Task Template and the one filing-local workflow README pointer; do not change spacedock new help or unrelated files.
  Commit `11b03fd4` changes exactly the fixture producer, its new test, and the Claude filing runtime pointer; the branch is based atop PR #491 and its help work is untouched.
- DONE: Run gofmt, focused affected-package tests, go test ./..., and go test ./... -race; commit the deliverable and append a complete implementation Stage Report.
  Repository Go sources were formatted (the one unrelated pre-existing formatting delta was restored), the 14-case focused run and affected-package run passed, and both full repository gates exited 0 before commit `11b03fd4`.

### Summary

The filing fixture now writes one complete, column-zero, pipe-safe Task Template matching the development workflow shape, and the Claude first officer is told to consult it when filing. The new offline test proves the exact on-disk fenced body creates an entity with a minted ID and that indenting its frontmatter fence produces the expected `no frontmatter` failure. The implementation is isolated to the three specified files and is ready for independent validation.

## Stage Report: implementation (cycle 2)

- DONE: Transplant only the filing-fixture implementation onto current `origin/main`, excluding PR #491.
  Fetched `origin/main` at `2137832b` and replayed only the implementation change; corrected commit `edde1514` is the branch's sole commit above that base. The prior report's `11b03fd4` SHA and claim that PR #491 ancestry was appropriate are superseded by this report.
- DONE: Prove the branch diff contains exactly the three approved files and no PR #491 changes.
  `git diff --name-only origin/main...HEAD` lists only `internal/ensigncycle/filing_readme_template_test.go`, `internal/ensigncycle/shared_fixtures_test.go`, and `skills/first-officer/references/claude-first-officer-runtime.md`; the stat is 3 files, 155 insertions, 1 deletion.
- DONE: Re-run the focused tests, `go test ./...`, and `go test ./... -race` on the corrected ancestry.
  The focused Task Template/discovery run passed in `internal/ensigncycle`, and both complete repository gates exited 0 after the transplant; `git diff --check origin/main...HEAD` also passed.

### Summary

The deliverable is unchanged, but its ancestry is repaired: commit `edde1514` now sits directly on current `origin/main` with no PR #491 help, skill-promotion, or prose-guard history. The branch and this report now agree on an exact three-file implementation, with all required tests freshly green.

## Stage Report: validation

- DONE: Independently audit the branch against current origin/main. Confirm it is exactly one c3 commit above origin/main, changes only the three approved files, and contains no PR #491 ancestry or behavior.
  After a fresh fetch, `origin/main=2137832b`; `HEAD=edde1514` has that exact parent and `origin/main...HEAD` is `0 1`. The diff is exactly the approved three paths; `git merge-base --is-ancestor 205ca99a HEAD` exited 1, and the PR #491 help-producing paths have no branch diff.
- DONE: Reproduce the filing fixture end value: the on-disk Task Template round-trips exactly, including its indented fenced block; the filing-local README pointer is present; PR #490 discovery behavior remains intact; no spacedock-new help text changed.
  The focused run passed the literal-template create/ID-mint subtest and its indented-frontmatter `no frontmatter`/no-file negative control; the 14-case run kept `TestSharedScenarioFixturesAreDiscoverable` green. Focused diff review places the README pointer only under `## Filing New Entities` and leaves Startup plus `internal/status/new.go`, its stub test, `cmd`, and `internal/cli` unchanged.
- DONE: Run focused affected-package tests, gofmt verification, git diff --check, go test ./..., and go test ./... -race. Record exact commands/results and the reviewed commit in the validation Stage Report.
  Reviewed `edde151478d6c7d946d1404f4de3dbc172d5226c`; `go test ./internal/ensigncycle -run 'Test(FilingReadmeTaskTemplateRoundTripsThroughNew|SharedScenarioFixturesAreDiscoverable)' -count=1` passed 14 cases.
  `go test ./internal/ensigncycle ./internal/status ./internal/contractlint -count=1`, `git diff --check origin/main...HEAD`, `go test ./...`, and `go test ./... -race` all exited 0.
  `gofmt -l internal/ensigncycle/shared_fixtures_test.go internal/ensigncycle/filing_readme_template_test.go` printed nothing. The broader `gofmt -l ./cmd ./internal` reports only unchanged `origin/main` file `internal/release/journeydelta.go`, verified by a zero `git diff --exit-code origin/main --` result.
- DONE: AC-1 — The filing fixture exposes the workflow-local entity shape that a real development workflow exposes.
  The on-disk fixture test observed exactly one Task Template with column-zero frontmatter, `title`, `status: backlog`, no `id`, and all five required body sections; the development README independently supplies the same section scaffold.
- DONE: AC-2 — The literal fixture template is pipe-safe through the production creation path.
  The verbose targeted test passed literal `status.NativeRunner --new` creation with `created:` output and a minted ID, while the indented opening fence exited 1 with `no frontmatter` and created no entity.
- DONE: AC-3 — The Claude filing contract points to the workflow Task Template at the moment of filing.
  Focused diff review shows the specified sentence as the sole addition under `## Filing New Entities`; Startup is unchanged, and no banned prose-grep test was added.
- SKIPPED: Live-model lanes.
  The approved test plan explicitly excludes them because this task claims fixture fidelity and parser round-trip state, not a filing-hunt-rate improvement; PR #491 already left no behavioral headroom for that measure.
- DONE: PASSED recommendation.
  All three acceptance criteria are satisfied by commit `edde1514`; no material validation finding remains in the scoped branch.

### Summary

Validated the corrected one-commit branch independently against freshly fetched `origin/main`. The fixture's actual README template round-trips through native creation, the malformed-indentation control fails as intended, PR #490 discovery remains green, and PR #491 help behavior is absent. Recommend PASSED, with the repository-wide pre-existing `internal/release/journeydelta.go` formatting output recorded but not attributed to this branch.
