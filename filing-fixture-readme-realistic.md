---
title: "Make the filing live-test fixture README representative of a real workflow (add a Task Template)"
status: ideation
source: "Captain decision 2026-07-10 on the science-officer dp-premise finding. Verbatim: then fix the fixture's readme to be realistic. new --help is not a priority. Finding it acts on (science-officer-verified): filingReadme() at internal/ensigncycle/shared_fixtures_test.go:313-332 carries no Task Template and no copyable stub (grep for template over it is empty), so dp's AC-1 filing-hunt methodology is structurally blind to a README-based filing approach — it measures a stripped fixture. The real docs/dev/README.md carries a full column-0 pipe-safe ## Task Template. And the FO filing contract in claude-first-officer-runtime.md's ## Filing New Entities never points the FO at the workflow README before filing (grep across skills/first-officer for a README/template-consultation instruction returns zero)."
started: 2026-07-10T12:52:14Z
completed:
verdict:
score: 0.4
worktree:
issue:
id: c3wxhq3qj94mhakam80g4zxw
---

## Problem

The `filing` live-scenario fixture README (`filingReadme()`, `internal/ensigncycle/shared_fixtures_test.go:313-332`) is unrealistically minimal: frontmatter (`commissioned-by`, `id-style: sequential`, two stages), a `# Filing Fixture` heading, fixture-explanation prose, and `### backlog`/`### done` stage blurbs — no `## Task Template`, no copyable entity stub. A real workflow README (the canonical model is `docs/dev/README.md`) carries a full `## Task Template`: a column-0, pipe-safe `---` YAML frontmatter stub (title + initial `status`, `id` omitted) followed by `## Problem` / `## Proposed approach` / `## Acceptance criteria` / `## Test plan` scaffolding. Because the fixture strips this, any live measurement of filing behavior against it (dp's AC-1 among them) is testing a fixture that structurally cannot detect a README-based filing approach — a real workflow's README already carries the answer, the fixture's does not. Compounding it: the FO contract never tells the FO the workflow README is a filing reference (`claude-first-officer-runtime.md`'s `## Filing New Entities` teaches `spacedock new` + stdin stub and warns off `--next-id`, but never names the workflow README's `## Task Template`).

## Proposed approach

**Primary — make the fixture README representative.** Rewrite `filingReadme()` to carry a `## Task Template` modeled on `docs/dev/README.md`'s: a column-0 (pipe-safe) `---` stub with `title` + `status: backlog` (`id` omitted) and the standard scaffolding sections. Keep the existing frontmatter (`commissioned-by: spacedock@1`, `id-style: sequential`, the backlog/done stages) and the fixture-explanation prose so discoverability and the existing filing assertions are unaffected.

Before/after (the fixture's README body):
- BEFORE: `# Filing Fixture` + fixture prose + `### backlog` / `### done` blurbs, no template.
- AFTER: the same, plus a `## Task Template` section — a column-0 `---` / `title:` / `status: backlog` / `---` / one-paragraph-body stub, then `## Problem` / `## Proposed approach` / `## Acceptance criteria` / `## Test plan` scaffolding — mirroring `docs/dev/README.md`'s Task Template shape.

**Secondary (folded in) — close the contract's README-blindness.** Add a filing-time pointer to `claude-first-officer-runtime.md`'s `## Filing New Entities`: before filing, consult the workflow README's `## Task Template` for the entity's shape and section scaffolding. This is belt-and-suspenders with the already-shipped `new --help` stub — the README carries the richer, workflow-specific sections `new --help` deliberately omits — and it gives the now-realistic fixture something that actually consults it.

## Out of scope

- dp's `spacedock new --help` stub (shipped, PASSED, deprioritized by the captain — do NOT extend it).
- Reducing the live filing-hunt rate: dp already took `filing` to 0/3 hunts, so there is no live behavioral headroom; this entity claims NO live hunt-rate delta and runs no live lane.

## Acceptance criteria

**AC-1 (value — the filing fixture README is representative of a real workflow README; offline, deterministic, independent baseline).** After the change, `filingReadme()` carries a `## Task Template` with (a) a column-0 `---` frontmatter stub naming `title` and the initial `status` with `id` omitted, (b) the `## Problem`/`## Proposed approach`/`## Acceptance criteria`/`## Test plan` scaffolding, matching `docs/dev/README.md`'s Task Template shape (the independent baseline), AND (c) the stub is pipe-safe — extracting the fixture README's `---…---` block and piping it through `runNew` creates the entity (`created:`), not a `no frontmatter found` error. Baseline moves the wrong way: current HEAD's `filingReadme()` has no Task Template (test reds on HEAD); an indented (non-column-0) fence reds the pipe-safe round-trip. *Test:* an offline Go test in `internal/ensigncycle` (or `internal/status`) that builds the fixture, asserts the template + scaffolding sections are present, and round-trips the stub through `runNew` — reusing the pipe-safe round-trip pattern dp's cycle-3 audit established (`TestPrintedStubTemplateBlockIsPipeSafe`).

**AC-2 (mechanism — the FO contract names the workflow README's Task Template as a filing reference; offline, cheap; serves AC-1's representativeness by giving the fixture a consumer).** `claude-first-officer-runtime.md`'s `## Filing New Entities` instructs the FO to consult the workflow README's `## Task Template` before filing. *Test:* an offline `internal/contractlint` presence guard with a red-control (removing the pointer reds; the pointer is absent on HEAD, so the guard reds on HEAD).

## Test plan

- Offline only (seconds, no model spend): AC-1 fixture-structure + pipe-safe round-trip test (`internal/ensigncycle`/`internal/status`); AC-2 contract-pointer presence guard (`internal/contractlint`) with red-control. Confirm existing fixture tests stay green — `TestSharedScenarioFixturesAreDiscoverable` (the enriched README must keep `commissioned-by`), the `filing` scenario assertions, and `go test ./internal/ensigncycle/... ./internal/status/... ./internal/contractlint/...`.
- No live-workflow test: no behavioral headroom post-dp; representativeness is proven structurally against the real README.

## Sequencing dependencies (flag for dispatch)

- Edits `internal/ensigncycle/shared_fixtures_test.go` (`filingReadme`), which **sc5 (#490) also edits** (fix A threads `workflowRoot` through the prompt builders; fix B touches fixture READMEs). Sequence AFTER sc5 merges, or rebase onto it, to avoid a conflict in that file. (Worker-dispatched: test code.)
- Edits `skills/first-officer/references/claude-first-officer-runtime.md`, which **dp (#491) also edits**. Sequence AFTER dp merges, or rebase onto it. (Worker-dispatched: shipped `skills/**` product.)
