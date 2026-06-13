---
id: ndpfbqqvezggnrydnvrxjmh2
title: Prefer `spacedock new` over manual --next-id for filing entities
status: validation
source: captain (2026-06-13) — noticed FO still calls --next-id when filing seeds
started: 2026-06-13T05:52:37Z
completed:
verdict:
score:
worktree: .worktrees/spacedock-ensign-prefer-new-over-next-id
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

## Stage Report: implementation

- DONE: `status --next-id` emits a use-`new` hint on stderr WITHOUT altering the stdout id (stdout stays exactly the id from computeNextID), and the `--next-id --json` path stays hint-free
  native_runner.go: stderr hint before the stdout id write in the plain-text branch only; TestNextIDPlainTextEmitsNewHintOnStderr + TestNextIDJSONStaysHintFree (id from computeNextID via stdout, not read back from prose). Proven end-to-end via the built binary: stdout `001`, stderr the hint; `--json` stdout `{"command":"next-id","id":"001"}`, stderr empty. Commit e2fd0dd8.
- DONE: teach `spacedock new <slug> [--folder] [--id-seed S --id-actor A] < stub` as the blessed atomic-create path in BOTH first-officer-shared-core.md (Status Viewer + ID Styles + FO Write Scope) and claude-first-officer-runtime.md, keeping --next-id documented only for candidate-preview, with the split-root path-scoped commit+push note; PROVE behavior with a NEW `filing` shared scenario in internal/ensigncycle
  New `filing` scenario: empty-workflow fixture + filingPrompt + host-specific stream-grading assertions (assertClaudeFilingViaNew / assertCodexFilingViaNew: filed via `new <slug>`, did NOT emit the `--next-id`+Write pair), wired into both runner maps + Pi coverage + the seed-scenarios doc-lock. TestSharedScenarioRunnerCoverage green. Offline positive/negative assertion tests model the manual-flow drift. Commit e2fd0dd8.
- DONE: keep TestSharedScenarioRunnerCoverage green with a Codex-side runner
  Added runCodexFilingScenario + map entry alongside the Claude runner; live coverage meta-tests pass under -tags live (no model spent).

### Summary

AC-1: the `--next-id` plain-text path now writes a `spacedock new` hint to stderr (the id stays the only thing on stdout) and the `--json` path is unchanged/hint-free; the hint is a native-only divergence stripped in the parity normalizer like STATE_BACKEND, with its presence pinned by a dedicated test. AC-2: the FO contract teaches `spacedock new` as the atomic-create path across both files (Status Viewer, ID Styles, FO Write Scope, Claude runtime), proven by a new `filing` shared scenario that grades the FO's tool-call stream on both hosts. Two corrections surfaced by exercising the real binary: `new` requires a frontmatter-bearing stub on stdin (not plain prose) and writes flat `<slug>.md` with the id stamped into frontmatter — both fixed in the prose and the fixture's expected path. Full offline module suite green (1253), status suite 453, live builds + coverage guards green. Live FO drives (TestLiveClaudeSharedScenarios/TestLiveCodexSharedScenarios `filing`) are the validation stage's job — they spend a model and run in CI/validation, not here.

## Stage Report: validation

VERDICT: PASSED

- DONE: AC-1 reproduce — `go test ./internal/status/ -run NextID` (TestNextIDPlainTextEmitsNewHintOnStderr + TestNextIDJSONStaysHintFree), built-binary stdout/stderr split, status suite + `go test ./...`
  4/4 NextID tests PASS. Built binary (`go build ./cmd/spacedock`): plain `--next-id` stdout = exactly `001\n` (hex `3030310a`, no hint leak), stderr = the use-`new` hint, exit 0; `--json` stdout = `{"command":"next-id","id":"001"}`, stderr = 0 bytes, exit 0. `go test ./internal/status/` 453 PASS; `go test ./...` 1253 PASS across 16 packages.
- DONE: AC-2 LIVE drive (load-bearing) — `TestLiveClaudeSharedScenarios/filing` live, Codex leg, `TestSharedScenarioRunnerCoverage` parity
  Claude live PASS (57.8s, claude-sonnet-4-6): FO filed via `${SPACEDOCK_BIN:-spacedock} new wire-the-thing --workflow-dir . <<'EOF'` with a stdin stub; ZERO `--next-id` Bash commands, ZERO `Write` calls; `new` returned `created: …/wire-the-thing.md id=001`. Codex live PASS (80.3s): FO filed via `${SPACEDOCK_BIN:-spacedock} new wire-the-thing <<'EOF'`; ZERO `--next-id` command_executions (the 14 raw `--next-id` string hits were all the FO reading the contract prose, not running it); entity landed with `id: 001`. Coverage parity `TestSharedScenarioRunnerCoverage` + `TestSeedScenariosDocLock` PASS (no model spend).
- DONE: Detached adversarial audit on a separate throwaway checkout (worktree never mutated)
  Throwaway `git worktree add --detach e2fd0dd8`; three refutations all confirmed the guards bite, then restored+removed. (1) Revert the native_runner stderr hint → TestNextIDPlainTextEmitsNewHintOnStderr REDS ("stderr must hint at `spacedock new`, got """), JSON test stays green. (2) Weaken the filing stream-assertion to drop the manual-pair check → TestAssertClaudeFilingViaNew REDS on the `--next-id`+`Write` negative case; same for the Codex `--next-id` guard. (3) Remove `stripNextIDHint` from the parity normalizer → TestIndReadFlagsSeq/next-id REDS (native carries the hint, golden is empty), proving the strip is load-bearing and the golden was correctly NOT regenerated to absorb the hint.

### Summary

PASSED. AC-1 is proven outside the prose: the built binary's `--next-id` puts only the id on stdout (hex-verified) with the use-`new` hint on stderr, and `--json` stays hint-free with empty stderr — offline suite 1253 green. AC-2's load-bearing live drives confirmed on BOTH hosts that the real FO files via `spacedock new` and never reaches for the `--next-id`+write manual pair (both with-my-own-eyes stream inspection, not just the green assertion). The detached adversarial audit refuted nothing material — all three intended guards (the stderr-hint test, the filing stream-assertion's manual-pair check, the parity strip) genuinely red under their adversarial edits, so none green-lights a future regression. No instruction-file prose-grep in the new tests: the `strings.Contains(errOut, "spacedock new")` checks observe the runtime binary's stderr, not a contract file.
