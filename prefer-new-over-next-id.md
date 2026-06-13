---
id: ndpfbqqvezggnrydnvrxjmh2
title: Prefer `spacedock new` over manual --next-id for filing entities
status: ideation
source: captain (2026-06-13) — noticed FO still calls --next-id when filing seeds
started: 2026-06-13T05:52:37Z
completed:
verdict:
score:
worktree:
issue:
sprint: 0202-survey-improvements
group: cleanup
sprint-readiness: ready
---

Make `spacedock new` the contract-blessed atomic-create path for filing entities, and have `status --next-id` emit a hint pointing to it.

## Problem

`spacedock new [--folder] SLUG` (alias `status --new`) has existed since #242 (2026-05-31): it mints the id and atomically writes a valid stamped entity in one call. But the FO operating contract never adopted it — both `first-officer-shared-core.md` (Status/Dispatch sections) and `claude-first-officer-runtime.md` still teach the manual flow: call `status --next-id`, hand-assemble the frontmatter, `Write` the file. The result: every FO files entities the slow, drift-prone way (a `--next-id` candidate is not a reservation, so the fetched id can drift from what finally lands), and operators get no signal the atomic path exists.

## Proposed approach

Two coupled changes:
1. **Contract** (`skills/first-officer/references/first-officer-shared-core.md` + `claude-first-officer-runtime.md`): teach `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < body` as the blessed atomic-create path for seed filing. Keep `--next-id` documented only for its candidate-preview use. Note that for split-root state checkouts the FO still does the path-scoped commit + push after `new` (new writes, it does not commit).
2. **Binary** (`internal/status`): `status --next-id` emits a stderr hint that `spacedock new` files atomically (so an operator reading next-id output is pointed at the better path). Hint must not pollute the stdout id value that callers parse.

## Out of scope

Changing what `new` itself does (it already mints + atomic-writes correctly). Auto-committing from `new` (commit/push stays the caller's, per split-root concurrency-safety rules).

## Acceptance criteria

Each AC names a property of the finished entity, not a stage action, and how it is verified.

**AC-1 - `status --next-id` emits a use-`new` hint on stderr without altering the stdout id, and the `--json` path stays hint-free.**
Verified by: a Go test in `internal/status` asserting (a) for the plain-text path, stdout is exactly the id (unchanged from current) while stderr contains the `new` hint; (b) for the `--next-id --json` path, stdout is the same single-key `{"command":"next-id","id":...}` object with no hint, since a stderr hint would not pollute it but the test pins that the machine-readable path is unchanged; exit code unchanged in both. The expectation source is independent of the file under test: the stdout id comes from `computeNextID`, the test does not read the hint string back out of the contract prose.

**AC-2 - The FO contract teaches `spacedock new` as the atomic-create path and a live drive shows the FO filing via `new`, not the manual --next-id+Write flow.**
Verified by: a live first-officer drive in the `internal/ensigncycle` shared-scenario harness — a new `filing` scenario (fixture + prompt asking the FO to file one seed, plus a durable-state assertion) added in the same shape as the existing `gate-guardrail` / `rejection-flow` scenarios. The assertion grades the FO's recorded tool-call stream: the FO emitted a `spacedock … new <slug>` invocation and did NOT emit the `--next-id` + `Write` pair. Grading the stream (the FO's actual behavior) — not a grep of the contract prose, and not just the end-state file (which looks the same either way) — is what keeps this AC honest. The contract prose change alone does not satisfy this AC.

## Riskiest mechanism

No spike needed. The risky assumptions both resolve by reading existing code:

- **stdout/stderr separation in `--next-id`** — `internal/status/native_runner.go:224` already writes only the id via `fmt.Fprintln(stdout, id)`, with `stderr` an in-scope writer alongside it (the JSON branch at :220-223 is separate). Adding a `fmt.Fprintln(stderr, hint)` before the stdout write is a clean, already-separated split — no parser round-trip or format risk.
- **`spacedock new` is the real atomic-create path** — confirmed at `internal/status/new.go` (mint id-style-appropriate id, stamp into stdin frontmatter, temp+rename atomic write) and wired as a top-level verb aliasing `status --new` at `internal/cli/cli.go:308`. The proposed contract text `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < body` matches the actual interface.

The proven mechanisms this task relies on: the existing stdout/stderr split in `runStatus`, the existing `runNew` atomic-create, and the existing `internal/ensigncycle` live-drive harness.

## Files carrying the text change

The FO-contract prose change lands in two files, BOTH under `skills/first-officer/references/` — `first-officer-shared-core.md` (the `Status/Dispatch` and `FO Write Scope` sections) and `claude-first-officer-runtime.md`. These are product scaffolding under `skills/`, off-limits for direct FO edits on main per the scaffolding guardrail (`first-officer-shared-core.md:258`); they ship as the deliverable and are built by a worker under test at implementation. This is exactly why AC-2's proof is a live drive, not a prose-grep: the same worker authoring the contract text cannot self-certify that the FO obeys it.

## Test plan

AC-1 is a cheap Go unit test in `internal/status` (stdout/stderr separation plus the `--json`-stays-clean assertion). AC-2 is the expensive half: a new `filing` shared scenario in `internal/ensigncycle` — fixture workflow, a prompt asking the FO to file one seed, and a stream-grading assertion (`spacedock new` invoked, `--next-id`+`Write` absent), built in the same shape as the existing scenarios (`writeXWorkflow` + `XPrompt()` + `assertX`, registered in `claudeScenarioRunners` / `sharedRuntimeScenarios`). The harness has no existing FO-filing scenario to reuse, so this adds one rather than extending an existing run; the shared coverage meta-test then requires a Codex-side runner for parity. Contract prose edits are authoring work but are not an AC on their own — the behavior is proven only by the live drive.

## Stage Report: ideation

- DONE: Firm AC-1 (status --next-id emits a use-`new` hint on stderr WITHOUT altering the stdout id — a Go test asserting stdout is exactly the id and stderr carries the hint)
  AC-1 sharpened with independent-expectation note (id from computeNextID, not read back from contract) and a `--json`-stays-clean clause; grounded at native_runner.go:224 (stdout-only id write).
- DONE: Firm AC-2 (the FO contract teaches `spacedock new` as the atomic-create path, proven by a LIVE first-officer drive observing the FO file via `new`, never a prose-grep of the skill)
  AC-2 now specifies a new `filing` shared scenario in internal/ensigncycle grading the FO's tool-call stream (new invoked, --next-id+Write absent) — not the end-state file, not a prose grep.
- DONE: Confirm which contract/skill files carry the text change (product scaffolding, built by a worker under test)
  Both under skills/first-officer/references/: first-officer-shared-core.md + claude-first-officer-runtime.md; flagged as scaffolding-guardrail files (shared-core.md:258), worker-built under test.
- DONE: Riskiest mechanism — record "no spike needed" naming the existing stdout/stderr split
  Recorded "no spike needed": stdout/stderr already split at native_runner.go:224; runNew atomic-create confirmed at new.go; `new` verb aliases status --new at cli.go:308.
- DONE: Note the text-vs-behavior split: the contract prose edit is real authoring but NOT an AC on its own; AC-2's behavior is proven only by the live drive
  Captured in both the AC-2 body and the "Files carrying the text change" section.

### Summary

Verified the entity's design against the live codebase rather than its prose: the stdout/stderr split (native_runner.go:224), the real `runNew` atomic-create (new.go), and the `new` verb wiring (cli.go:308) all confirm the approach with no spike needed. Corrected one imprecision — the ensigncycle harness has no FO-filing scenario to "reuse," so AC-2 now calls for a NEW `filing` shared scenario graded on the FO's tool-call stream (not the indistinguishable end-state file, not a prose grep). Flagged both contract files as scaffolding-guardrail product (skills/first-officer/references/), which is exactly why AC-2 must be a live drive and not self-certification by the authoring worker.
