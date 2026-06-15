---
id: 6rtpj5avcp733tb15dfjcbbb
title: Summarize CI artifacts for FO/ensign reads — replace whole-log (143KB) reads with a triage summary
status: ideation
source: FO + 0.20.4 scope survey (2026-06-14, this session) — CI logs (143KB cited in e6a's source as a recurring read sink) are read whole into FO/ensign context for validation and triage, when the agent needs only pass/fail + the failing lines. Tokens scale with the whole log. 0.20.4 read-cost theme; lower-frequency than the entity/README reads e6a covers.
started: 2026-06-15T05:19:11Z
completed:
verdict:
score: 0.30
worktree:
issue:
sprint: 0204-structured-reads
sprint-readiness: ready
---

CI logs are read whole into context for triage when the agent needs only pass/fail, the exit code, and the failing tests/lines. Give the FO/ensign a triage summary they can read instead of the 143KB log, while keeping the full artifact reachable when the summary is not enough.

## Problem

Validation and triage read large CI logs whole (143KB cited). The signal an agent needs is small: did it pass, what failed, where. The rest is noise that fills context. There is no summary surface, so callers load the whole artifact or grep blind.

## Proposed approach

A **read-time `spacedock` helper** that parses a `go test -v` transcript (file path or stdin) into a 10-20 line triage summary: overall verdict + process exit code, per-package pass/fail counts with the failing package names, any panic with its user-code frame, and the top-N failing tests each with `file:line` and the one-line assertion message — followed by a pointer back to the full log. The full artifact is never consumed or replaced; the helper *reads* it, so the whole log stays one `Read`/`grep` away when the summary is not enough.

### Why a read-time helper, not a CI-emitted artifact (decided on spike evidence)

The two candidate directions were (A) a CI step that parses the transcript and uploads a small `summary.md` beside the full log, and (B) a read-time `spacedock` helper the agent runs over a transcript it already has. **B is primary; A is an optional thin add-on that composes.** Rationale, each point grounded in the spike below:

- **The savings is context tokens, and B captures them regardless of how the file arrived.** The 143KB log lives on disk (the agent `gh run download`s it, or pipes `gh run view --log`); disk bytes cost nothing. The cost is loading those bytes into the agent's *context*. B collapses the on-disk 170KB / 2842-line transcript to an 11-line read — the 143KB never enters context. The spike measured exactly this: an 11-line summary carrying every triage signal.
- **B works on every historical log with zero CI change.** A only helps runs produced *after* the CI change ships, and forces the agent to learn a second artifact name. The summary is a pure deterministic function of the transcript text (no live state, no network), so it can be computed at read time over any existing artifact.
- **B keeps the full log reachable by construction.** Because the helper reads the artifact rather than replacing it, "no information loss" (AC-2) is free: the summary's `source:` line names the artifact and its last line gives the exact `grep` (`--- FAIL:` / `panic:`) to pull the full failing region. The triage drive (spike) confirmed the FO can locate every failure from the summary but needs the full log for *root cause* — so the pointer-back is load-bearing, not decorative.
- **stdin composition avoids even the disk write.** `gh run view --log | spacedock <helper> -` streams the transcript through the parser; the agent reads 11 lines and never materializes the 143KB at all. Proven in the spike (`-` reads stdin).

A (CI side) remains a cheap future add-on: the same parser, run in the `if: always()` step, can write the summary to `$GITHUB_STEP_SUMMARY` so the run's web page shows triage at a glance. It reuses B's parser; it is not required for the token savings and is out of scope here (noted below).

### Surface shape (the helper)

- **A new `go test -v` summarizer reached through the existing native CLI surface.** The exact verb/flag is implementation's call, but the shape is fixed: it takes one positional argument — a transcript path, or `-` for stdin — and emits the triage summary to stdout, exit 0 on a clean parse (independent of the run's own verdict; the verdict is *content*, not the helper's exit). Sibling prior art for "a pure-function text reader on the native surface" is the e6a `status --read` helper (`internal/status/section_read.go`); this is the log-shaped analogue (test-transcript grammar instead of markdown-heading grammar). They are distinct grammars and must not be conflated — see Out of scope.
- **Input grammar = `go test -v` lines** (the format every CI transcript here uses; confirmed: no `-json`, no `gotestsum` anywhere — `.github/workflows/runtime-live-e2e.yml:179,198` pipe `go test ... -v 2>&1 | tee`). The parser keys on: `--- FAIL: <name> (Ns)` (top-level and indented subtest), the preceding `    file_test.go:NN: <msg>` location line, `panic: <msg>` plus the first `*_test.go:NN` frame in the goroutine stack (user code), and the per-package `FAIL\t<pkg>\t<dur>` / `ok\t<pkg>\t<dur>` result lines.
- **Output (10-20 lines), in this order:** `verdict=PASS|FAIL exit=N`; a **source line** (`source: <path> (<N> lines)`); package pass/fail counts + failing-package names; a `PANIC:` line with its user frame when present; then up to 10 leaf failures as `TestName  file:line  <truncated msg>`; then the full-log pointer with the exact grep. Subtests roll up to their leaf (`Parent/Sub`, inheriting the parent's captured location); a panicking test links to its panic frame.
- **The source line is the ONLY path-dependent line, and it is normalized for stdin.** When the input is a file path, the source line reads `source: <path> (<N> lines)`. When the input is stdin (`-`), there is no path, so the source line reads `source: <stdin> (<N> lines)` — the literal token `<stdin>`, not the run's path. Everything below the source line (verdict, packages, panic, the failure rows, the pointer-with-grep) is byte-identical between the path and stdin forms over the same bytes. This makes AC-5's equality assertion satisfiable: it asserts byte-identity of the summary *excluding the source line* (the line count is identical; only the path token differs). The pointer-back line names the source generically ("the full transcript") plus the grep, so it carries no path and stays identical too.

## Riskiest-mechanism spike (DONE — exercised against a real failing log, not asserted)

The riskiest path: does a deterministic parser over `go test -v` text reliably extract verdict + exit + the failing tests' `file:line` such that an FO triages CORRECTLY from the summary alone — including the hazards real logs carry (subtests, a panic that aborts the package, buried failures, truncation)? Spiked end-to-end:

**Fixtures (real, not synthetic prose).** Ran the actual repo suite `go test -v ./internal/...` → a real **172KB / 2868-line** all-pass transcript. Then built a throwaway failing module (`/tmp/ci-log-spike/failmod`) with five planted failures and spliced its `go test -v` output into the middle of the pass transcript → a realistic **170KB / 2842-line** failing log with the failures BURIED at lines 1407-1422 (the real triage problem: ~5 signals among 2842 lines). Independent oracle = the planted failures read from `cat -n` of the source, never from the summarizer.

**Parser + run.** A throwaway Go summarizer (`/tmp/ci-log-spike/summarize.go`) produced an **11-line** summary of the failing log and a **5-line** summary of the passing log — both inside the 10-20 line budget. Against the oracle it got every signal right:

| Failure (oracle from `cat -n`) | Ground-truth loc | Summary loc | Match |
|---|---|---|---|
| `TestGammaFails` | a_test.go:10 | a_test.go:10 | yes |
| `TestDeltaFails` | a_test.go:14 | a_test.go:14 | yes |
| `TestZetaSubtests/case_bad` (subtest) | b_test.go:8 | b_test.go:8 | yes |
| `TestEtaPanics` (panic) | b_test.go:13 | b_test.go:13 | yes |
| package verdict / exit | FAIL | `verdict=FAIL exit=1` | yes |

The subtest case and the panic case were the two hazards the spike *caught*: `go test -v` prints the subtest's `file:line` BEFORE the parent's `--- FAIL` marker and the leaf `--- FAIL: Parent/Sub` AFTER it, so a naive "loc attaches to next FAIL" rule mis-attributes; the fix is leaf-inherits-parent-loc. A panic emits no normal per-test loc — the user frame `b_test.go:13` is buried among Go-runtime frames in the goroutine stack; the fix is "first `*_test.go:NN` frame after `panic:`". Both fixes are in the spike and seed the implementation's first tests.

**Oracle is real, not a tautology (mutation-checked).** Breaking the `file:line` regex → all locations drop to "(no assertion loc — see full log)". Removing the panic from the verdict condition → the failing log falsely reports `verdict=PASS`. The tests would catch a broken parser.

**Triage drive (the AC-1 behavioral proof).** Handed ONLY the 11-line summary (not the log) to a fresh agent acting as FO and asked the seven triage questions a validator faces. It answered all seven correctly against the oracle: verdict+exit, failing package, all four test names, every `file:line`, the panic site, and correctly prioritized `b_test.go` (panic + highest fix-density). Crucially it also drew the **information-loss boundary**: the summary locates and classifies every failure but cannot give *root cause* (why `compute()` returns 7, which map is nil) — that needs the full log. This is the empirical justification for AC-2: the summary is a triage surface, not a replacement, so the full artifact must stay reachable.

**Secondary risks cleared.** stdin form (`-`) works (composes with `gh run view --log | spacedock …`, no disk write). Read-time cost is trivial: **0.125s** over the 170KB log. A TRUNCATED log (cut mid-panic before any `FAIL\tpkg` line — real `gh` log truncation, debrief entity 3g) still reports `verdict=FAIL exit=1`, because the `--- FAIL:` markers alone trip the verdict — it does NOT falsely PASS. Pinned as AC-3.

## Out of scope

- **e6a's entity/README section reads** (`status --read`, `internal/status/section_read.go`) — a markdown-HEADING grammar over entity bodies/README. THIS task is a `go test -v` TRANSCRIPT grammar (test markers, locations, panics, package results). Different input, different parser; deliberately kept distinct so neither helper grows a second grammar.
- **`go test -json` / `gotestsum` adoption** — would give structured events, but the repo emits neither anywhere (confirmed) and changing what CI runs is out of scope. The helper parses the `-v` text CI already produces; it does not change the CI command.
- **The CI-side `$GITHUB_STEP_SUMMARY` emission (direction A)** — a future thin add-on reusing this parser in the `if: always()` step. Composes with B, not required for the token savings, not built here.
- **Non-`go test` logs** (install-e2e shell output, release logs, host CLI stream artifacts) — the grammar is `go test -v` specifically. Other log shapes are separate work if they ever become a read sink.
- **Replacing or deleting the full log artifact** — the full transcript upload stays exactly as is; this helper only adds a cheap *read path* over it.

## Acceptance criteria

Each AC names a property of the finished outcome, not a stage action, and how it is verified. The external oracle throughout is the committed failing-log fixture's KNOWN planted failures (their names and `file:line` read independently of the parser), never a prose/regex self-match.

**AC-1 — The helper emits a ≤20-line triage summary that carries every triage signal of a real large failing `go test -v` log: verdict + process exit code, the failing package(s), and every failing test with its correct `file:line` and one-line symptom — including a subtest leaf and a panicking test.**
*Tested by:* a Go unit test running the parser over a committed failing-transcript fixture (the spike's large failing log, frozen under `internal/.../testdata/`) and asserting, table-driven against the fixture's known planted failures, that each failing test appears with the exact `file:line` the fixture's source defines (`a_test.go:10`, `a_test.go:14`, the `b_test.go:8` subtest, the `b_test.go:13` panic frame) and that the summary is ≤20 lines. Mutation-guarded: the spike already showed that breaking the `file:line` regex drops all locations and removing the panic from the verdict flips FAIL→PASS, so the oracle is non-tautological.

**AC-2 — An FO/ensign triages the failing run CORRECTLY from the summary alone (without the whole log), AND the summary points back to the full log so root-cause work stays possible (no information loss).**
*Verified by:* a behavioral triage drive — a fresh agent given ONLY the summary answers the validator's triage questions (verdict/exit, failing package, every failing test + `file:line`, the panic site, fix priority) and they match the fixture's planted failures; the same drive confirms the agent can name what it would need the full log for (root cause) and that the summary carries the reach-back — the `source:` line names the artifact and the last line gives the exact `grep`/`Read` to retrieve the full failing region. The spike already ran this drive (7/7 correct + the information-loss boundary drawn); validation re-runs it against the committed fixture. Proof is the behavioral answer set matching the oracle, never the summary text containing a phrase.

**AC-3 — A truncated or panic-aborted log does not falsely report success.**
*Tested by:* a Go unit test over a truncated copy of the fixture (cut mid-panic, before any `FAIL\t<pkg>` result line) asserting the summary still reports `verdict=FAIL exit=1` — the `--- FAIL:` markers alone trip the verdict. (Real hazard: `gh` log truncation, debrief entity 3g.) Mutation-guarded: a parser that keyed verdict only off the trailing `FAIL\t<pkg>` line would report PASS here and fail this test.

**AC-4 — A clean (all-pass) log reports success in a handful of lines and surfaces no spurious failures.**
*Tested by:* a Go unit test over the committed all-pass fixture (the real 172KB pass transcript, frozen) asserting `verdict=PASS exit=0`, zero failing tests, and ≤ a small line bound. Guards against a parser that mistakes a `=== RUN` / package-name / stack-frame line for a failure.

**AC-5 — The helper reads from both a file path and stdin and produces the same triage (only the source line differs), and its own exit status reflects parse success, not the run's verdict.**
*Tested by:* a Go behavior test driving the binary over the same fixture bytes two ways — `<helper> <fixture-path>` and `cat <fixture> | <helper> -` — and asserting the two summaries are **byte-identical EXCEPT the single `source:` line** (the path form reads `source: <path> (<N> lines)`, the stdin form reads `source: <stdin> (<N> lines)`; the line count and every other line match). The test strips the `source:` line from both and asserts byte-equality of the remainder. Separately it asserts exit status: the helper exits 0 on a well-formed failing transcript (the FAIL is *content*; the helper's exit ≠ the run's exit) and non-zero only on an unreadable/garbage input. (stdin composition with `gh run view --log | spacedock <helper> -` is the no-disk-write path; the spike proved `-` reads stdin.) *Why the carve-out:* stdin has no path, so a naive whole-output byte-identity is unsatisfiable; normalizing the one path-bearing line — the `source:` line, with the pointer-back line kept path-free — makes the equality real and the carve-out exactly one well-defined line.

**AC-6 — An FO actually triages a real failed CI run by fetching the transcript and running the summarizer at the triage site, reading ~15 lines instead of the whole log — proven by that behavioral trace, never by the instruction text.**

This is the e6a live-helper-use bar (e6a AC6: the proof is the FO *using* `status --read` at a wired site, not the contract saying to). The savings only lands if a real triage actually routes through the helper, so the proof is the trace, not the wording. Two parts, because **grep confirms no CI-triage read site exists today** (searched `skills/` + `agents/` for transcript/triage/`gh run`/live-CI read instructions — the only hits are unrelated: pi context-seeding, the dispatch context-budget probe, teardown prose; none tell an FO to read a CI transcript). So this AC owns *creating* the site, not adopting an existing one:

- **AC-6a — The triage-read site exists and is the FO's CI-triage instruction.** The finished entity adds a CI-triage read step to the FO read discipline at the natural home — `skills/first-officer/references/first-officer-shared-core.md` `## Completion and Gates` (where the FO validates a stage that just ran live CI; this is the same read-cost family e6a wired, and the `Working Principles` "prefer a code gate over prose" rule at `:204` is exactly why the wording alone is not the proof). The step instructs: when validating a failed live run, do NOT read the `*-transcript.txt` whole — fetch it (`gh run download` / `gh run view --log`) and run the summarizer, opening the full log only for root cause. The exact wording is in `## Documentation changes`.
  *Tested by:* the existing skill-text/doc-contract guards passing over the edited file (the under-test scaffolding bar — AC-6a only pins the site exists and is guarded; it is NOT the savings proof).

- **AC-6b — The behavioral trace: an FO triages a real failed run through the helper.** A live (or faithfully-replayed) FO, at the `## Completion and Gates` triage step, validates a run whose live CI FAILED: it fetches the real transcript and invokes `${SPACEDOCK_BIN:-spacedock} <ci-summarize> <transcript-or-->`, reads the ~15-line summary, and proceeds to triage (names the failing tests + `file:line`, decides reject/fix) — with the ≥2000-line transcript NEVER loaded whole into its context. The proof is that tool-call trace (the fetch + the summarizer invocation + the short read standing in for the whole-log read), exactly as e6a AC6 demands the `--read` call trace.
  *Verified by:* a validation-stage live FO triage drive against a real failing transcript (a frozen real failing artifact is acceptable as the fetched input — the bar is that the FO *runs the helper and triages from its output*, not that CI is re-run live). The instruction edit (AC-6a) proves nothing on its own; AC-6b is the behavioral evidence the savings lands. The spike already ran the agent-triages-from-summary half (7/7 correct from the 11-line summary, oracle-matched); AC-6b extends that to the FO doing it AT THE WIRED SITE via the helper invocation, closing the loop e6a's AC6 closes for `--read`.

## Test plan

- **Mechanism (already paid):** the spike above — throwaway parser + real 170KB failing fixture + a live FO triage drive + mutation checks + timing + truncation. Seeds AC-1/AC-3's first tests. Cost: done.
- **Fixtures:** freeze the spike's two transcripts under `testdata/` — the 170KB failing log (with its five planted failures) and the 172KB all-pass log. The failing module's source (`a_test.go`, `b_test.go`) is committed beside them so the oracle (`file:line`) is checkable in-repo. Cost: copy the spike artifacts; the planted-failure table is the oracle.
- **Unit (Go):** AC-1 (failing-log signal extraction, table-driven over planted failures), AC-3 (truncated log → still FAIL), AC-4 (clean log → PASS, no spurious failures). Pure-function parser, no I/O beyond reading the fixture; minutes. Mutation-checked as the spike showed (regex break, panic-verdict break).
- **Behavior (Go, drives the binary):** AC-5 — path-vs-stdin equality *excluding the single `source:` line* (strip it from both, assert byte-equality of the remainder) + exit-status semantics, reusing the native-runner test harness pattern.
- **Live/behavioral (validation):** AC-2 — re-run the FO triage drive against the committed fixture (answers match the oracle; information-loss boundary drawn; full-log pointer present). The only thing not provable in a unit test is that an FO genuinely triages from the summary; the spike already ran it, validation confirms against the frozen fixture.
- **Adoption — site + behavioral trace (AC-6, the e6a bar):** AC-6a applies the doc diff that *creates* the CI-triage read step at `first-officer-shared-core.md` `## Completion and Gates` (no such site exists today — confirmed by grep) and runs the existing skill-text/doc-contract guards over the edited file. AC-6b is the live FO triage drive at that wired site: a real failed run, the FO fetches the transcript and invokes the summarizer, triages from the ~15-line read, the whole transcript never loaded. Cost: one live exercise; the instruction edit alone is not the proof — the trace is.

## Documentation changes

This task changes agent-facing read discipline (how an FO/ensign triages a failed CI run), so ideation owns a concrete doc diff; implementation applies it.

### (a) FO triage-read discipline — a NEW step (no such site exists today)

Grep over `skills/` + `agents/` confirms there is no current instruction telling an FO/ensign to read a CI transcript or triage a failed live run (the only transcript/`gh run` hits are unrelated — pi context-seeding, the dispatch context-budget probe, teardown prose). So this is an ADDED step, not a reword. Home: `skills/first-officer/references/first-officer-shared-core.md` `## Completion and Gates` (the FO validates a stage that just ran live CI here; the `Working Principles` "prefer a code gate over prose" rule at `:204` is why the wording is scaffolding and AC-6b's trace is the proof). Proposed added step (exact verb/flag finalized by implementation against the shipped surface):

> **Triaging a failed live CI run.** Do not read the `*-transcript.txt` whole (these run 140KB+ / 2000+ lines). Fetch it (`gh run download`, or pipe `gh run view --log`) and run the transcript summarizer — `${SPACEDOCK_BIN:-spacedock} <ci-summarize> <transcript-path-or-->` — to get verdict + exit code + each failing test's `file:line` and symptom in ~15 lines, then reject/triage from that. Open the full transcript only for root cause; the summary's source line names the artifact and its last line gives the `grep` (`--- FAIL:` / `panic:`) to jump to the failing region.

(The pointer-back line in the summary is path-free — it says "the full transcript" plus the grep, not a re-echoed path — so the stdin and path forms stay byte-identical below the `source:` line, per AC-5.)

### (b) User-facing CLI doc

Add a one-line row for the new helper to whatever enumerates `spacedock` subcommands/flags (the README command area and, if present, the mkdocs command-reference page) — addition only, no rewording of existing entries. Concrete wording deferred to implementation against the shipped verb/flag name.

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
