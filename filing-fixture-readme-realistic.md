---
title: "Make the filing live-test fixture README representative of a real workflow (add a Task Template)"
status: implementation
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
