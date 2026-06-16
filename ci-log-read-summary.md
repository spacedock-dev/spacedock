---
id: 6rtpj5avcp733tb15dfjcbbb
title: Summarize CI artifacts for FO/ensign reads — replace whole-log (143KB) reads with a triage summary
status: validation
source: FO + 0.20.4 scope survey (2026-06-14, this session) — CI logs (143KB cited in e6a's source as a recurring read sink) are read whole into FO/ensign context for validation and triage, when the agent needs only pass/fail + the failing lines. Tokens scale with the whole log. 0.20.4 read-cost theme; lower-frequency than the entity/README reads e6a covers.
started: 2026-06-15T05:19:11Z
completed:
verdict:
score: 0.30
worktree: .worktrees/spacedock-ensign-ci-log-read-summary
issue:
sprint: 0204-structured-reads
sprint-readiness: ready
---

The CI test steps dump VERBOSE `go test -v` output to stdout, and stdout is what floods FO/ensign context when reading a CI run. Fix it at the source: keep verbose detail OUT of the visible/stdout surface and ARCHIVE it to a file/artifact for root-cause retrieval; stdout carries only the clean concise `go test` result (package pass/fail + failures), small enough to read directly into context.

## Problem

The live-CI test steps run `go test ... -v 2>&1 | tee <transcript>.txt` (`.github/workflows/runtime-live-e2e.yml:179,198`). The `-v` flag turns `go test` into a firehose: every `=== RUN` / `--- PASS` for every test, even on a green run. That firehose is what goes to stdout AND to the archived transcript — and stdout is what an FO/ensign loads when it reads the run. The signal a reader needs is small (which package/test failed, where); the `-v` per-test trace is noise that scales with the test count. Measured on the real suite: clean default `go test ./...` is **17 lines / 1.1KB**, while the `-v`/json firehose is **~143KB+ / thousands of lines** (the `-json` event stream is 5891 lines / 1.28MB). The fix is to stop sending the firehose to the visible surface — not to parse it after the fact.

## Proposed approach

A **test-output discipline applied at the source** (the CI workflow test steps, and the documented local invocation): the visible/stdout surface carries only clean `go test` output (per-package pass/fail + failures-only), while the full verbose detail (`go test -json` jsonl) is written to an archived file and uploaded as a CI artifact, retrievable for root-cause inspection. One test run produces both. No downstream firehose-parser is built (the prior cycles' summarizer helper is dropped).

### What the spike settled about the mechanism (decisive)

The spike (below) replaced the obvious assumption with measured behavior. Three findings drive the design:

1. **Default `go test` is ALREADY the clean surface — the `-v` flag is the entire firehose.** Plain `go test ./...` prints one `ok`/`FAIL`/`?` line per package on a green run (17 lines / 1.1KB on the real suite), and on a red run prints ONLY the failing tests with their `file:line` plus the package `FAIL` line — the exact common-read signal, no per-test noise. So the visible-surface fix is largely "**stop passing `-v` to the stdout side**," not a new tool. This is nearly free and dependency-free.
2. **A single bare `go test` process emits exactly ONE format** (text OR `-json`), so clean-text-stdout AND a json archive cannot both come from one bare invocation — you need something that tees the run into two renderings. The captain's "do not run twice" rules out a second `go test`. So the archive-in-the-same-run requires either a tee-capable tool or stdlib plumbing.
3. **`gotestsum` does exactly the two-rendering split in one run** — clean stdout (a per-package progress view + a `=== Failed` recap listing each failing test with `file:line`) AND a full `--jsonfile` archive, exit code preserved. Empirically confirmed (`go run gotest.tools/gotestsum`). It is off-the-shelf, the de-facto standard Go-CI formatter — NOT a helper we build. Cost: one third-party build/test dependency.

### The decision: clean default stdout + archived `-json`, with `gotestsum` as the one-run tee

**Recommended shape:** the CI test step runs the suite ONCE, sends clean output to stdout, and writes the full `-json` detail to an uploaded artifact. The mechanism is `gotestsum --jsonfile <detail>.jsonl --format <clean> -- <go-test-args>` — one run, clean visible stdout, full archive. Rationale grounded in the spike:

- It satisfies all three checklist gates from a single run: (a) clean stdout gives which package/test failed + `file:line` (the spike's red-run clean output shows exactly this); (b) the `--jsonfile` archive is the complete event stream, sufficient for root cause and retrievable; (c) one execution, exit preserved.
- The dependency is justified because the alternative no-dep one-run is worse: a single bare `go test` can't emit both, and the stdlib tee (`go test -v 2>&1 | tee >(go tool test2json > detail.jsonl) | grep-clean`) is bash-only process-substitution with a brittle `grep`-rendered clean view (spike-tested — it works but interleaves loc/marker lines and isn't `sh`-portable). gotestsum's clean renderer is purpose-built and its `--jsonfile` is a single flag.

**No-dependency fallback (if the dependency is rejected):** for the **offline gate** specifically — which re-runs cheaply and whose stdout today is already plain `go test ./...` (NO `-v`; `runtime-live-e2e.yml:68`) — the fix is nothing: stdout is already clean and small. The archive there is optional (a failing offline gate is reproduced locally for free). gotestsum earns its keep only on the **live lanes**, where the run is expensive and `-v` is currently forced; there the one-run-both archive matters. So a graded answer is available: live lanes adopt gotestsum (or the `-json`-archive + clean-stdout split); the offline gate is already compliant.

### Where the archive lives and how an FO reaches it

- **CI:** the `-json` detail file is uploaded by the existing `actions/upload-artifact@v5` step (the live lanes already upload `live-e2e-transcript.txt` etc. at `:200-213`) — the detail jsonl joins that artifact list, replacing the `-v` text transcript. An FO fetches it the standard way: `gh run download <run-id>` or `gh run view --log`. The clean stdout (small) is what shows in the run's step log / `$GITHUB_STEP_SUMMARY` and is read directly.
- **Local:** the documented pattern (dev README Testing-Resources) runs the same one-liner; the jsonl lands at a known path beside the repo for `go tool test2json`-aware inspection or a plain `grep '"Action":"fail"'`.

### The adjacency the spike cleared (why dropping `-v` on the live lanes is safe)

The live lanes use `go test -tags live ... -v 2>&1 | tee`. The `-v` there is NOT load-bearing for the runtime watchdog: the live runner's `streamWatcher` (the 60s no-progress quiet budget, `runtime-live-e2e.yml:189-192`) watches the **host agent's stream-json StdoutPipe inside the test**, not `go test`'s `-v` mirror (confirmed: `streamWatcher`/`quietBudget` live in `internal/ensigncycle` test-side line-drainer code consuming the agent stream, with no dependency on `go test -v`). So `-v` on the CI line exists only to produce the human-readable archived transcript — exactly the firehose being replaced by the `-json` archive. Dropping `-v` from the stdout side does not starve the watchdog.

## Riskiest-mechanism spike (DONE — exercised on the real suite + a real failing module, not asserted)

The riskiest path for the SOURCE-SIDE approach: can ONE test run produce (a) clean stdout that still carries the common-read signal (which package/test failed), AND (b) a retrievable full-detail archive sufficient for root cause — without running the suite twice? Spiked the mechanism directly (`/tmp/tod-spike`, Go 1.26.1).

**(a) Clean stdout already carries the signal — `-v` is the firehose.** Real suite `go test ./...` (no `-v`): green run = **17 lines / 1128 bytes**, one `ok`/`?` line per package. A throwaway failing module's red run printed ONLY the failing tests with `file:line` + the package `FAIL` line:
```
--- FAIL: TestGammaFails (0.00s)
    a_test.go:9: compute() = 7, want 42
--- FAIL: TestDeltaFails (0.00s)
    a_test.go:12: delta precondition not met: handle was nil
--- FAIL: TestZetaSubtests (0.00s)
    --- FAIL: TestZetaSubtests/case_bad (0.00s)
        b_test.go:7: subtest assertion: got "x" want "y"
FAIL
FAIL	failmod	0.143s
```
That is the whole common-read signal in ~10 lines, default behavior, no tool. By contrast the same suite under `-json` is **5891 lines / 1.28MB**, and the current `-v` transcript is ~143KB+ — confirming the firehose is the `-v`/json verbosity, not the test count itself.

**(b) One run, clean stdout + json archive — three candidate mechanisms, measured.**
- **A bare `go test` cannot do both.** `go test -json` *converts* output to JSON (one format out); a single process emits text XOR json. So "clean stdout AND json archive from one bare `go test`" is impossible — confirmed against `go help test`. Two runs are the only bare-`go test` way, and the captain rules that out.
- **`gotestsum` does both in one run (confirmed).** `go run gotest.tools/gotestsum@latest --jsonfile detail.jsonl --format pkgname -- ./...` produced clean stdout (**15 lines**, incl. a `=== Failed` recap listing each failing test + `a_test.go:9` / `b_test.go:7`) AND wrote the full `detail.jsonl` (**47 lines**, the complete event stream) in a SINGLE execution, exit 1 preserved. It is off-the-shelf (the standard Go-CI formatter), not a helper we build.
- **No-dep stdlib tee works but is inferior (confirmed).** `go test -v 2>&1 | tee >(go tool test2json > detail.jsonl) | grep-clean` produced the json archive (48 lines via `test2json`) AND a clean-ish stdout from one run — but the `grep`-rendered clean view interleaves loc/marker lines awkwardly and the `>(...)` process substitution is bash-only (not `sh`-portable). It proves a no-dependency one-run is *possible*, but the output is worse than gotestsum or the plain default.

**(c) Where the archive lives / FO reach.** The `-json` detail file uploads via the existing `actions/upload-artifact@v5` step (replacing the `-v` text transcript in the artifact list); an FO fetches it with `gh run download` / `gh run view --log`. Clean stdout shows in the step log / `$GITHUB_STEP_SUMMARY`, read directly. Local: same one-liner, jsonl at a known repo-adjacent path.

**Adjacency cleared (dropping `-v` is safe).** The live lanes' `streamWatcher` 60s quiet budget watches the **host agent's stream-json StdoutPipe inside the test**, not `go test -v` (confirmed: the watcher + `quietBudgetDefault` live in `internal/ensigncycle` test-side line-drainer code over the agent stream; no `go test -v` dependency). The `-v` on the CI line exists only for the human transcript being replaced. Dropping it from stdout does not starve the watchdog.

**Disproved prior-cycle premise.** Cycles 1-2's `Out of scope` claimed "no `-json` / `gotestsum` — the helper parses the `-v` text CI already produces." The captain redirect inverts that: `-json` + (optionally) `gotestsum` ARE the mechanism now, applied at the source, and the `-v`-text-parsing helper is dropped.

## Out of scope

- **A downstream firehose-parser / `spacedock` log-summarizer helper** — the cycle-1/2 approach, explicitly dropped by the captain. We do not parse the `-v` transcript after the fact; we change what the source emits.
- **Changing WHICH tests CI runs** — the test selection, `-tags live`, `-run` filters, `-count`, `-timeout` stay exactly as they are. This task changes only the OUTPUT rendering (drop `-v` from the visible surface; archive `-json`), not the test command's coverage.
- **Running the suite twice** — the one-run-both constraint is a hard requirement, not a convenience; any design that re-runs `go test` to get the second rendering is rejected.
- **The offline gate is already compliant** — `runtime-live-e2e.yml:68` runs plain `go test ./...` (no `-v`); its stdout is already clean. The change targets the live lanes' `-v | tee` steps; the offline gate needs no edit (optionally gains a `-json` archive).
- **e6a's entity/README section reads** (`status --read`, `internal/status/section_read.go`) — unrelated markdown-heading reads; untouched.
- **Non-`go test` CI output** (install-e2e shell, release logs, host CLI stream artifacts) — separate shapes; out of scope unless they become a read sink.
- **The agent stream-json transcript** the live test archives (the per-job `CLAUDE_CONFIG_DIR` session jsonl) — that is the AGENT's diagnostic, kept as is; this task touches only the `go test` rendering.

## Acceptance criteria

Each AC names a property of the finished outcome, not a stage action, and how it is verified. The change is to the real CI workflow (`.github/workflows/runtime-live-e2e.yml`) and the documented local invocation; the proofs are behavior over that changed workflow, not prose.

**AC-1 — The visible/stdout surface of each changed CI test step is clean and small: per-package pass/fail plus failures-only, with no per-test `=== RUN`/`--- PASS` verbosity.**
*Verified by:* running the changed test step's exact command locally over the real suite and asserting the stdout is the clean form — on a green run, one `ok`/`?`/`FAIL` line per package and no `=== RUN` lines (a `grep -c '=== RUN' <stdout>` of 0); on a red run (a deliberately-failing fixture package), only the failing tests with `file:line` + the package `FAIL` line appear. The spike already measured the clean form (17 lines / 1.1KB green; failures-only red); the AC pins the *changed command* reproduces it. Oracle: the byte count / `=== RUN` count of the changed command's stdout, run directly — not a claim in the workflow YAML.

**AC-2 — The full verbose detail is archived to a file and uploaded as a CI artifact, and it is sufficient + retrievable for root cause.**
*Verified by:* (i) the changed step writes a `-json` detail file (the same run that produced the clean stdout) and the `upload-artifact` step lists it; proven by running the command locally and asserting the jsonl exists, is non-empty, and contains the failing test's events (`grep '"Action":"fail"'` returns the failing package/test) — the detail carries every event the `-v` transcript did, in structured form. (ii) An FO-reach check: the documented retrieval (`gh run download` / `gh run view --log`) names the artifact, and a `go tool test2json`-aware reader or a plain `grep` over the jsonl recovers a specific failure's full output. Oracle: the planted failure's presence in the jsonl, read independently.

**AC-3 — One test run produces both the clean stdout and the archive (no double `go test`).**
*Verified by:* the changed step invokes the test binary exactly once. Proven by inspecting the command (a single `gotestsum ... -- <args>` or a single `go test -json ... | tee`-style pipe, NOT two `go test` invocations) AND by a behavioral check: instrument the command over a fixture and assert the suite's tests each execute once (e.g. a test that appends to a counter file runs it N times, not 2N). Guards against a "run clean, then re-run for json" regression.

**AC-4 — The exit code of the test step still reflects test pass/fail (the rendering change does not mask failures).**
*Verified by:* running the changed command over a passing fixture (exit 0) and a failing fixture (exit non-zero); `set -o pipefail` (already on the live steps, `:177,197`) plus the tool's exit-preservation keep the step red on failure. The spike confirmed gotestsum preserves exit 1; the AC pins it for the chosen mechanism. Oracle: the process exit code, observed.

**AC-5 — If a build/test dependency is added (e.g. gotestsum), it is pinned and declared; if no dependency is added, that is explicit.**
*Verified by:* either (a) the chosen mechanism adds no third-party dependency (plain `go test`/`go tool test2json` only) — asserted by the absence of any new tool in the workflow and a one-line note in the body; or (b) gotestsum is added pinned to a version (a `go run gotest.tools/gotestsum@<ver>` with an explicit version, or a `tools.go`/go.mod tool directive), with the dependency noted in the workflow and the dev docs. Oracle: the workflow text + go.mod/tools state. This AC forces the dependency decision to be explicit and pinned, not floating `@latest`.

**AC-6 — The FO/ensign read discipline reflects the new shape: read the clean stdout / step summary for triage; fetch the archived `-json` detail only for root cause.**
*Verified by:* the doc diff in `## Documentation changes` applied at the FO read-discipline site and the dev README Testing-Resources area, under the existing skill-text/doc-contract guards (the `internal/ensigncycle` `*DocsContract` family that guards the README's Runtime-Live-CI section). The instruction now says: the clean step output IS the triage read (small, read directly); the `-json` artifact is the root-cause fetch. *Note:* unlike the cycle-2 helper approach, this AC's deliverable is the source change itself (AC-1..AC-4) — the savings lands because the firehose no longer reaches stdout, not because an agent runs a helper. So the doc edit is genuinely scaffolding here: the binary-level guarantee (clean stdout) is the real assurance (per the `:204` "prefer a code gate over prose" rule), and AC-1's behavioral stdout check is that gate. No live FO drive is required for the savings — it is structural.

### Before/after for the changed CI steps (recorded so implementation applies verbatim)

Live Claude lane, `runtime-live-e2e.yml:179` —
> **Before:** `go test -tags live -count=1 -run 'TestLiveEnsignCycle|TestLiveZeroDiscoverReportsAndStops|TestLiveStandingResidencyInjectsCommOfficer' ./internal/ensigncycle/ -v 2>&1 | tee live-e2e-transcript.txt`
> **After (gotestsum shape):** `gotestsum --jsonfile live-e2e-detail.jsonl --format pkgname -- -tags live -count=1 -run '…' ./internal/ensigncycle/` — clean stdout, full jsonl archive, one run. (Or the no-dep shape: `go test -tags live … -json ./internal/ensigncycle/ 2>&1 | tee live-e2e-detail.jsonl | go-test-clean-view` with the offline gate's plain stdout as the model — implementation picks per AC-5.)

Line `:198` (shared scenarios) — same transform: drop `-v | tee …transcript.txt`, emit clean stdout + `claude-shared-scenarios-detail.jsonl`. The codex (`:319`) and pi (`:426`) live steps run `go test … -v` without a `tee` today; bring them to the same clean-stdout + `-json`-archive shape for parity (or leave non-`tee` `-v` steps to implementation's judgment — the firehose-to-archive sites with `| tee` are the primary targets).

`upload-artifact` lists (`:205-211`, `:326-329`, `:433-436`) — replace the `*-transcript.txt` entries with the `*-detail.jsonl` files.

## Test plan

- **Mechanism (already paid):** the spike above — real-suite clean vs `-v`/json sizes, gotestsum one-run-both confirmed, the no-dep stdlib tee tried, the `-v` watchdog adjacency cleared. Seeds AC-1..AC-4. Cost: done.
- **Fixtures:** a small committed failing-test package under `testdata/` (the spike's `failmod` shape — a few planted failures incl. a subtest) so the red-run stdout/`-json` assertions have a deterministic, in-repo oracle without depending on a real CI run. The repo's own green suite is the green-run oracle.
- **Behavioral, over the CHANGED command (AC-1..AC-4):** a test (Go or a shell check the CI itself runs) that executes the changed step's exact command over the fixture and asserts: clean stdout shape + `grep -c '=== RUN'` == 0 on green / failures-only on red (AC-1); the `-json` archive exists, is non-empty, and contains the planted failure's `"Action":"fail"` event (AC-2); the suite runs once not twice (AC-3, e.g. a counter-file fixture); exit 0 green / non-zero red with `pipefail` (AC-4). These run offline, in seconds. The oracle is the observed bytes/exit/jsonl, never the workflow YAML text.
- **Dependency decision (AC-5):** assert the workflow + go.mod/tools state matches the chosen path — no new tool, or a PINNED gotestsum version. A grep/test over the workflow for a floating `@latest` fails the AC.
- **Docs/guards (AC-6):** apply the read-discipline + Testing-Resources doc diff; run the existing `internal/ensigncycle` `*DocsContract` guards (which pin the README's Runtime-Live-CI content) over the edits so the doc change stays under test.
- **No live re-run needed for the savings.** Unlike the cycle-2 helper, the savings is structural (the firehose no longer reaches stdout), provable offline over the changed command. A live CI run is the final confirmation but not the gating proof — AC-1's behavioral stdout check is.

## Documentation changes

This task changes the CI test-output rendering and the FO/ensign read discipline around it; ideation owns the concrete diffs, implementation applies them.

### (a) The CI workflow steps — the primary deliverable

The before/after for `runtime-live-e2e.yml` `:179`, `:198` (and the codex/pi live steps + the `upload-artifact` lists) is recorded verbatim under `## Acceptance criteria` → "Before/after for the changed CI steps". The transform is uniform: drop `-v | tee <…>-transcript.txt`; emit clean stdout + a `<…>-detail.jsonl` archive from one run; swap the artifact-list entry from the `.txt` transcript to the `.jsonl` detail.

### (b) FO/ensign read discipline — the clean output IS the triage read

Home: `skills/first-officer/references/first-officer-shared-core.md` `## Completion and Gates` (where the FO validates a stage that just ran live CI). Added note:

> **Reading a live CI result.** The CI test steps now print CLEAN output — per-package pass/fail and, on failure, the failing tests with `file:line`. Read that step output / job summary directly for triage; it is small. Fetch the archived `*-detail.jsonl` (`gh run download`, or `gh run view --log`) ONLY for root cause — it is the full event stream, not the triage read.

(The `:204` "prefer a code gate over prose" rule applies: the clean stdout is guaranteed by the workflow change + AC-1's behavioral check, so this note is scaffolding pointing at the structural guarantee, not the assurance itself.)

### (c) Dev README Testing-Resources

`docs/dev/README.md` Testing-Resources table / `## Runtime Live CI` area — document the local pattern: the one-liner that emits clean stdout + the `-json` archive, and how to inspect the jsonl (`go tool test2json`-aware reader or `grep '"Action":"fail"'`). Addition only; no rewording of existing rows. If gotestsum is adopted (AC-5b), note the pinned version and the `go run gotest.tools/gotestsum@<ver>` invocation here.

## Stage Report: ideation

- DONE: Spike FIRST on a real large CI log (~143KB): produce a 10-20 line triage summary (pass/fail, exit code, top-N failures with file:line), then confirm an FO/ensign can triage correctly from the summary alone, and record the spike result in the body.
  Real fixtures: `go test -v ./internal/...` → 172KB/2868-line all-pass; spliced a throwaway failing module's output into it → 170KB/2842-line failing log, failures buried at L1407-1422. Throwaway parser `/tmp/ci-log-spike/summarize.go` produced an 11-line summary; oracle (`cat -n` of source) confirmed all four planted failures' exact `file:line` (a_test.go:10, :14, subtest b_test.go:8, panic b_test.go:13) + verdict/exit. A fresh FO agent given ONLY the 11 lines answered 7/7 triage questions correctly and drew the information-loss boundary (triage yes, root-cause needs full log). Mutation-checked non-tautological. Recorded in body "Riskiest-mechanism spike".
- DONE: Design the surface on that spike evidence: decide where the summary is produced (CI-emitted artifact versus a read-time `spacedock` helper) and how the full log stays reachable with no information loss, with behavior-first ACs proven over a real failing-log fixture.
  Body "Proposed approach" decides a READ-TIME `spacedock` helper (primary) over a CI-emitted artifact: the savings is context tokens not disk, B captures them over any historical log with zero CI change, and reading-not-replacing keeps the full log reachable by construction (AC-2 free). stdin form composes with `gh run view --log | spacedock … -` (no disk write). CI-side `$GITHUB_STEP_SUMMARY` emission kept as an out-of-scope thin add-on. AC-1..AC-6 are behavior-first over the frozen failing/passing fixtures (oracle = planted failures), incl. AC-3 truncated-log-does-not-falsely-PASS and AC-6 adoption-under-guards. Doc diff for the FO triage-read site included.

### Summary

Designed ci-log-read-summary as a read-time `spacedock` helper that parses `go test -v` transcript text (path or stdin) into a ≤20-line triage summary — verdict+exit, failing packages, panic frame, and every failing test with `file:line` and symptom, plus a pointer back to the full log. Chose read-time over CI-emitted because the cost is context tokens (not disk), the helper works on every historical log with no CI change, and reading-not-replacing makes "full log reachable / no information loss" free. Spiked the riskiest path FIRST on a real 170KB failing transcript: an 11-line summary captured all four planted failures' exact locations (oracle = source `cat -n`, mutation-verified non-tautological), and a fresh FO agent triaged 7/7 correctly from the summary alone while correctly identifying that root-cause still needs the full log — the empirical justification for keeping it reachable. Surfaced and pinned two real hazards the spike caught (subtest loc-attribution; panic frame in the goroutine stack) and one truncation edge (a cut-off log must not falsely report PASS). Scope held distinct from e6a's markdown-heading `status --read` grammar.

## Stage Report: ideation (cycle 2)

Folds the two Material items the independent preflight staff review flagged before the gate. Cycle-1 design kept (read-time-over-CI-emitted choice and the spike are sound and confirmed); this is corrective.

- DONE: Rebuild AC-6 from a prose-grep to e6a's live-helper-use bar, AND own locating/creating the triage-read site.
  Grepped `skills/` + `agents/`: confirmed NO CI-triage read site exists today (only unrelated transcript/`gh run` hits — pi context-seeding, dispatch context-budget probe, teardown prose). So AC-6 now owns CREATING the site: AC-6a (the new `## Completion and Gates` triage-read step exists + is under the skill-text/doc-contract guards — scaffolding only) and AC-6b (the behavioral trace: a live FO at that site fetches a real failed transcript, invokes the summarizer, triages from the ~15-line read with the ≥2000-line transcript never loaded whole — proof is the tool-call trace, never the instruction text). Test plan + Documentation changes (a) updated to name the concrete site and mark it a NEW step.
- DONE: Normalize AC-5 so the path-vs-stdin equality is satisfiable (stdin has no path).
  Surface-shape spec now pins the `source:` line as the ONLY path-dependent line — `source: <path> (<N> lines)` for a path, `source: <stdin> (<N> lines)` for stdin — with the pointer-back line kept path-free. AC-5 asserts byte-identity of the summary EXCLUDING the `source:` line. Proven satisfiable by a throwaway exercise: path and stdin outputs diffed on exactly the `source:` line, identical everywhere else (`diff` excluding `^source:` → IDENTICAL). Fixed two stale "last line names the source" references to attribute the artifact name to the `source:` line and the grep to the last line.

### Summary

Folded both staff-review Material items. AC-6 is rebuilt to e6a's bar: because grep confirms no CI-triage read site exists today, the AC now owns creating it — AC-6a pins the new `## Completion and Gates` triage-read step (under the existing guards, scaffolding only) and AC-6b is the behavioral trace of a live FO fetching a real failed transcript, running the summarizer, and triaging from the short read with the whole log never loaded (proof is the trace, never the wording). AC-5 is normalized: the `source:` line is the single path-dependent line (`<stdin>` token for stdin, pointer-back kept path-free), and AC-5 asserts byte-identity excluding that one line — proven satisfiable by a throwaway path-vs-stdin diff that matched everywhere except `source:`. The cycle-1 read-time-helper choice and the original spike are unchanged.

## Stage Report: ideation (cycle 3)

CAPTAIN REDIRECT — the summarizer-helper direction (cycles 1-2) is rejected entirely. Re-ideated to a SOURCE-SIDE test-output discipline: clean `go test` to stdout, full `-json` detail archived, one run. The cycle-1/2 body sections (Proposed approach, spike, Out of scope, ACs, test plan, doc changes) are rewritten; the prior stage reports stay as the historical record.

- DROPPED: the read-time `spacedock` go-test-v summarizer helper and all its ACs (AC-1..AC-6 of cycles 1-2). Per the redirect, we do not parse a firehose after the fact.
- DONE: Spike the riskiest mechanism FIRST for the source-side approach — confirm (a) clean stdout carries the common-read signal, (b) the archive is retrievable + sufficient for root cause, (c) one run produces both.
  `/tmp/tod-spike` (Go 1.26.1): real suite `go test ./...` green = 17 lines/1.1KB (vs `-json` 5891 lines/1.28MB and the current ~143KB `-v` transcript) — DEFAULT go test is already the clean surface; `-v` is the entire firehose. A failing module's red run printed failures-only with `file:line`. gotestsum (`go run gotest.tools/gotestsum@latest --jsonfile … --format pkgname`) produced clean stdout (15 lines, `=== Failed` recap with `file:line`) AND the full jsonl (47 lines) in ONE run, exit 1 preserved — off-the-shelf, not a built helper. A bare `go test` emits exactly one format (text XOR json) so it can't do both alone; the no-dep stdlib tee (`-v | tee >(go tool test2json>f) | grep`) works but is bash-only + brittle. Adjacency cleared: the live `streamWatcher` quiet budget watches the AGENT stream-json inside the test, not `go test -v`, so dropping `-v` is safe. Recorded in body "Riskiest-mechanism spike".
- DONE: Pick the mechanism on spike evidence; behavior-first ACs over the real CI-workflow change; record before/after for the steps + local pattern.
  Decision: clean default stdout + archived `-json`, with gotestsum as the one-run tee (no-dep fallback graded: offline gate `:68` already clean, live lanes adopt the tool). AC-1 (clean small stdout, `grep -c '=== RUN'`==0), AC-2 (jsonl archive non-empty + carries the failure + FO-retrievable), AC-3 (one run, not two — counter-file check), AC-4 (exit preserved with pipefail), AC-5 (dependency decision explicit + PINNED, no `@latest`), AC-6 (read-discipline doc edit as scaffolding — the clean stdout is the structural gate per `:204`). Verbatim before/after for `runtime-live-e2e.yml` `:179`/`:198` + codex `:319` + pi `:426` + the upload-artifact lists `:205-211`/`:326-329`/`:433-436` recorded under ACs; all line numbers re-verified against the workflow.

### Summary

Re-ideated 6r per the captain redirect: dropped the downstream summarizer helper entirely and moved the fix to the SOURCE — the CI test steps emit clean `go test` output to stdout (per-package pass/fail + failures-only with `file:line`, ~17 lines / 1.1KB on the real suite) while the full `-json` detail is archived to an uploaded artifact for root-cause retrieval, all from one test run. The spike proved the decisive facts empirically: default `go test` is ALREADY clean (the `-v` flag is the whole 143KB firehose), a bare `go test` can't emit both formats in one run, and gotestsum does the clean-stdout + `--jsonfile` split in a single run with exit preserved (off-the-shelf, not a built helper) — with a no-dependency stdlib-tee fallback that works but is inferior. Cleared the key adjacency: the live watchdog watches the agent stream, not `go test -v`, so dropping `-v` is safe. ACs are behavior-first over the changed workflow command (clean-stdout byte/`=== RUN` checks, jsonl-archive sufficiency, one-run-not-two, exit preservation, a PINNED-dependency decision), with verbatim before/after for every changed CI step and the local pattern. The savings is structural (the firehose no longer reaches stdout), provable offline — no live FO drive gates it.

## Stage Report: implementation

**AC-5 dependency decision (explicit): NO third-party dependency — AC-5(a).** Implementation chose the stdlib-only one-run shape over the ideation-recommended gotestsum. Reason: exercising the candidates surfaced a fragility the spike's `@latest` test hid — `go run gotest.tools/gotestsum@<ver>` compiles gotestsum's transitive `x/tools` with the INVOKING toolchain, and no single pinned version builds on both CI's Go 1.22 (needs ≤ v1.12.2, but its `x/tools` breaks on modern Go) and a modern local Go 1.26 (needs ≥ v1.13.0, whose go.mod requires Go 1.24 → forces a GOTOOLCHAIN auto-bump on CI's 1.22, i.e. a floating toolchain — the exact thing AC-5 forbids). The body explicitly lists "the `-json`-archive + clean-stdout split" as an acceptable no-dep path, so AC-5(a) is in-bounds. Mechanism: `go test -c` (compile once) → run the binary once → `tee` native `-test.v` to `go tool test2json` (a Go TOOLCHAIN tool, zero third-party dep, no version to pin) for the `*-detail.jsonl` archive → `grep -vE` the raw stream to the clean step log → `exit ${PIPESTATUS[0]}` (the test exit, captured past the grep's no-match exit). Decision recorded in the workflow step comments, the dev README row, and guarded by `TestLiveWorkflowPinsNoFloatingTool`.

- DONE: Transform the live-CI test steps in .github/workflows/runtime-live-e2e.yml (the :179/:198 `-v | tee *-transcript.txt` firehose sites + codex/pi parity) to emit CLEAN stdout AND archive the full `-json` detail as an uploaded artifact, from ONE test run; swap the upload-artifact `*-transcript.txt` entries to `*-detail.jsonl`. Make the AC-5 dependency decision explicit and PINNED.
  All 6 `-v` live sites transformed (Claude ensign+scenarios, codex resolver=`-v` drop only, codex scenarios, pi coverage, pi smoke); 2 `| tee transcript` firehose sites + codex/pi scenario archives now write `*-detail.jsonl`, swapped into all three upload-artifact lists. Stdlib-only `go test -c | tee | go tool test2json` — no dependency to pin. YAML re-parses (ruby YAML.load OK).
- DONE: Land the behavioral proof over the CHANGED command: a committed failing-test fixture package + a test asserting clean stdout with `grep -c '=== RUN'`==0 on green / failures-only on red (AC-1), the `-json` archive exists/non-empty/contains the planted `"Action":"fail"` (AC-2), the suite runs once not twice (AC-3), and exit 0 green / non-zero red with pipefail (AC-4). Offline over the fixture.
  Fixture `internal/release/testdata/cleanoutputfixture/` (FIXTURE_FAIL-gated failures incl. a subtest; a FIXTURE_COUNTER_FILE counter test). `internal/release/cilog_clean_output_test.go` (5 tests) DERIVES the runnable script from the workflow's actual step (re-targeted to the fixture, not a hand-copy) and runs it under bash — AC-1 green+red, AC-2, AC-3 (counter==1), AC-4 all PASS offline. Mutation-checked non-tautological: swapping the clean-view grep for `cat` REDs AC-1. `cilog_clean_output_workflow_test.go` binds all 5 live steps to the shape + a firehose-regression guard + the AC-5 no-floating-tool guard (mutation-checked: reintroducing `-v | tee transcript` REDs).
- DONE: Apply the AC-6 read-discipline doc note (first-officer-shared-core.md `## Completion and Gates` "Reading a live CI result") + the dev README Testing-Resources addition, and keep the existing internal/ensigncycle `*DocsContract` guards green.
  Note added under `## Completion and Gates` pointing at the structural guarantee (referencing the "code gate over prose" principle in `## Working Principles`). README Testing-Resources gained a "Clean log + `-json` archive from one run" row (addition only). `internal/contractlint` + `internal/ensigncycle` guards green; the one stale guard literal (`journey_workflow_test.go` mutating the old Claude `-v | tee transcript` command) updated to the new command line and re-greened.

### Summary

Implemented the source-side test-output discipline as a stdlib-only one-run-both shape: `go test -c` compiles the package once, the binary runs once, its `-test.v` stream is tee'd to `go tool test2json` (a Go toolchain tool — no third-party dependency, contra ideation's gotestsum recommendation, which I rejected after finding a real `go run gotestsum@ver` toolchain-coupling fragility the spike's `@latest` test hid) for the `*-detail.jsonl` archive, while a `grep -vE` renders the clean failures-only step log, with `${PIPESTATUS[0]}` preserving the test exit past the grep. All six `-v` live sites in runtime-live-e2e.yml are transformed and the upload-artifact lists swap transcripts for jsonl archives. The behavioral proof (AC-1..AC-4) derives its runnable script from the workflow's own step re-targeted to a committed fixture (no hand-copy drift), runs under bash, and is mutation-verified non-tautological; a structural guard binds every live step to the shape and forbids a firehose or floating-tool regression. AC-6 docs added as scaffolding pointing at the structural guarantee. Full offline `go build ./...` + `go test ./...` green. Note: the entity's "Before/after" recorded the gotestsum shape; my AC-5(a) no-dep choice is the explicit, justified deviation. The detached adversarial CI/release audit is validation's job; no PR opened (Commander handles merge).

## Stage Report: implementation (cycle 2)

CAPTAIN/FO DECISION OVERRIDE — after I surfaced the gotestsum-`go run` toolchain fragility, the FO ruled AC-5 = **(A) pinned PREBUILT gotestsum binary**, explicitly out-ruling the no-dep path and forbidding a Go-toolchain bump. Converted the cycle-1 stdlib mechanism to gotestsum per the decision. The source-side discipline, the fixture, and the AC structure are unchanged; only the rendering mechanism + its proof harness changed.

- DONE: Transform the live-CI test steps to emit CLEAN stdout AND archive the full `-json` detail from ONE run; PINNED dependency, explicit.
  All 5 live test runs (Claude ensign + scenarios, codex scenarios, pi coverage + smoke) now run `gotestsum --jsonfile <name>-detail.jsonl --format pkgname -- <go test args>` — one run, clean `=== Failed` recap with `file:line`, full jsonl archive, exit preserved natively. The codex resolver step keeps a plain `-v`-drop (fast unit test, no archive). Upload-artifact lists carry the `*-detail.jsonl`. gotestsum is PINNED to v1.13.0 and sha256-verified by a committed install script `.github/scripts/install-gotestsum.sh` (bash-3.2-portable case-based per-platform checksum map; fails the step on mismatch — verified: a corrupted sha exits 1 and installs nothing). Install steps added to all three live jobs AND the offline job. No `@latest`, no Go-toolchain change.
- DONE: Behavioral proof over the CHANGED command (AC-1..AC-4), offline over the fixture.
  `internal/release/cilog_clean_output_test.go` runs the real pinned `gotestsum` over the fixture (resolved from PATH or `GOTESTSUM_BIN_DIR`; `t.Skip` when absent, since CI installs it before `go test ./...`): AC-1 green (`grep -c '=== RUN'`==0) + red (failures-only with `file:line`), AC-2 (jsonl non-empty + planted `"Action":"fail"` + output), AC-3 (counter==1, one run), AC-4 (exit 0 green / non-zero red) — all PASS with gotestsum present, SKIP cleanly without. `cilog_clean_output_workflow_test.go` binds every live step to the gotestsum `--jsonfile` shape, asserts the pinned/sha256-verifying install in ≥3 live jobs, and forbids firehose/`@latest` regressions — mutation-checked non-tautological (removing `--jsonfile` REDs the shape guard; removing `shasum -c` REDs the install guard).
- DONE: AC-6 doc note + dev README; guards green.
  shared-core `## Completion and Gates` note is mechanism-neutral (kept). README Testing-Resources row rewritten to the gotestsum command + local `go install gotest.tools/gotestsum@v1.13.0` / the install script + jsonl-reading guidance. `journey_workflow_test.go` stale literal re-pointed to the new gotestsum command. Full offline `go build ./...` + `go test ./...` green WITH gotestsum on PATH (mirrors the CI offline job).

### Summary

Per the FO's AC-5 = (A) decision, converted the mechanism from the cycle-1 stdlib `go test -c | tee | go tool test2json | grep` to a PINNED PREBUILT gotestsum (v1.13.0, sha256-verified by `.github/scripts/install-gotestsum.sh`, no `go run`, no Go-toolchain bump). All five live test runs now invoke `gotestsum --jsonfile <name>-detail.jsonl --format pkgname -- <args>` — one run, clean `=== Failed` recap with `file:line`, full jsonl archive, exit preserved. The install step is added to all three live jobs and the offline job (so the AC tests actually RUN in CI rather than self-skip). The behavioral proof runs the real pinned gotestsum over the committed fixture (AC-1..AC-4, mutation-verified) and skips gracefully when gotestsum is absent; the structural guard binds every live step to the gotestsum shape and asserts the pinned, checksum-verifying install. AC-6 docs updated to the gotestsum pattern. The cycle-1 stdlib approach (still in git history at b756fadc) was sound and fully tested, but the FO chose the purpose-built faithful recap; this cycle implements that choice.
