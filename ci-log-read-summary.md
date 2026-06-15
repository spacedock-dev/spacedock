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
- **B keeps the full log reachable by construction.** Because the helper reads the artifact rather than replacing it, "no information loss" (AC-2) is free: the summary's last line names the source and the exact `grep` (`--- FAIL:` / `panic:`) to pull the full failing region. The triage drive (spike) confirmed the FO can locate every failure from the summary but needs the full log for *root cause* — so the pointer-back is load-bearing, not decorative.
- **stdin composition avoids even the disk write.** `gh run view --log | spacedock <helper> -` streams the transcript through the parser; the agent reads 11 lines and never materializes the 143KB at all. Proven in the spike (`-` reads stdin).

A (CI side) remains a cheap future add-on: the same parser, run in the `if: always()` step, can write the summary to `$GITHUB_STEP_SUMMARY` so the run's web page shows triage at a glance. It reuses B's parser; it is not required for the token savings and is out of scope here (noted below).

### Surface shape (the helper)

- **A new `go test -v` summarizer reached through the existing native CLI surface.** The exact verb/flag is implementation's call, but the shape is fixed: it takes one positional argument — a transcript path, or `-` for stdin — and emits the triage summary to stdout, exit 0 on a clean parse (independent of the run's own verdict; the verdict is *content*, not the helper's exit). Sibling prior art for "a pure-function text reader on the native surface" is the e6a `status --read` helper (`internal/status/section_read.go`); this is the log-shaped analogue (test-transcript grammar instead of markdown-heading grammar). They are distinct grammars and must not be conflated — see Out of scope.
- **Input grammar = `go test -v` lines** (the format every CI transcript here uses; confirmed: no `-json`, no `gotestsum` anywhere — `.github/workflows/runtime-live-e2e.yml:179,198` pipe `go test ... -v 2>&1 | tee`). The parser keys on: `--- FAIL: <name> (Ns)` (top-level and indented subtest), the preceding `    file_test.go:NN: <msg>` location line, `panic: <msg>` plus the first `*_test.go:NN` frame in the goroutine stack (user code), and the per-package `FAIL\t<pkg>\t<dur>` / `ok\t<pkg>\t<dur>` result lines.
- **Output (10-20 lines), in this order:** `verdict=PASS|FAIL exit=N`; source path + total line count; package pass/fail counts + failing-package names; a `PANIC:` line with its user frame when present; then up to 10 leaf failures as `TestName  file:line  <truncated msg>`; then the full-log pointer with the exact grep. Subtests roll up to their leaf (`Parent/Sub`, inheriting the parent's captured location); a panicking test links to its panic frame.

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
*Verified by:* a behavioral triage drive — a fresh agent given ONLY the summary answers the validator's triage questions (verdict/exit, failing package, every failing test + `file:line`, the panic site, fix priority) and they match the fixture's planted failures; the same drive confirms the agent can name what it would need the full log for (root cause) and that the summary's final line names the source artifact + the exact `grep`/`Read` to retrieve the full failing region. The spike already ran this drive (7/7 correct + the information-loss boundary drawn); validation re-runs it against the committed fixture. Proof is the behavioral answer set matching the oracle, never the summary text containing a phrase.

**AC-3 — A truncated or panic-aborted log does not falsely report success.**
*Tested by:* a Go unit test over a truncated copy of the fixture (cut mid-panic, before any `FAIL\t<pkg>` result line) asserting the summary still reports `verdict=FAIL exit=1` — the `--- FAIL:` markers alone trip the verdict. (Real hazard: `gh` log truncation, debrief entity 3g.) Mutation-guarded: a parser that keyed verdict only off the trailing `FAIL\t<pkg>` line would report PASS here and fail this test.

**AC-4 — A clean (all-pass) log reports success in a handful of lines and surfaces no spurious failures.**
*Tested by:* a Go unit test over the committed all-pass fixture (the real 172KB pass transcript, frozen) asserting `verdict=PASS exit=0`, zero failing tests, and ≤ a small line bound. Guards against a parser that mistakes a `=== RUN` / package-name / stack-frame line for a failure.

**AC-5 — The helper reads from both a file path and stdin, and its own exit status reflects parse success, not the run's verdict.**
*Tested by:* a Go behavior test driving the binary: `<helper> <fixture-path>` and `cat <fixture> | <helper> -` produce byte-identical summaries; the helper exits 0 on a well-formed failing transcript (the FAIL is *content*, exit≠the run's exit) and non-zero only on an unreadable/garbage input. (stdin composition with `gh run view --log | spacedock <helper> -` is the no-disk-write path; the spike proved `-` reads stdin.)

**AC-6 — The FO/ensign read discipline actually points at the helper at the CI-triage site.** The recurring token savings only lands if the validation/triage instruction adopts it. The doc/skill site that today says "read the CI transcript" carries the helper invocation, and the change is under the existing skill-text guards.
*Verified by:* the doc diff in `## Documentation changes` applied at the concrete site, plus the existing skill-text/doc-contract guards passing over the edited files. (Per the proof policy, the instruction edit alone is not proof the savings lands; AC-2's behavioral drive is. AC-6 only pins that the adoption wording exists and is under-test.)

## Test plan

- **Mechanism (already paid):** the spike above — throwaway parser + real 170KB failing fixture + a live FO triage drive + mutation checks + timing + truncation. Seeds AC-1/AC-3's first tests. Cost: done.
- **Fixtures:** freeze the spike's two transcripts under `testdata/` — the 170KB failing log (with its five planted failures) and the 172KB all-pass log. The failing module's source (`a_test.go`, `b_test.go`) is committed beside them so the oracle (`file:line`) is checkable in-repo. Cost: copy the spike artifacts; the planted-failure table is the oracle.
- **Unit (Go):** AC-1 (failing-log signal extraction, table-driven over planted failures), AC-3 (truncated log → still FAIL), AC-4 (clean log → PASS, no spurious failures). Pure-function parser, no I/O beyond reading the fixture; minutes. Mutation-checked as the spike showed (regex break, panic-verdict break).
- **Behavior (Go, drives the binary):** AC-5 — path vs. stdin byte-equality + exit-status semantics, reusing the native-runner test harness pattern.
- **Live/behavioral (validation):** AC-2 — re-run the FO triage drive against the committed fixture (answers match the oracle; information-loss boundary drawn; full-log pointer present). The only thing not provable in a unit test is that an FO genuinely triages from the summary; the spike already ran it, validation confirms against the frozen fixture.
- **Adoption:** AC-6 — apply the doc diff at the triage site; the existing skill-text/doc-contract guards keep it under test.

## Documentation changes

This task changes agent-facing read discipline (how an FO/ensign triages a failed CI run), so ideation owns a concrete doc diff; implementation applies it.

### (a) FO/ensign triage-read discipline

The validation/triage instruction that today implies reading the whole CI transcript should route through the helper. The concrete site is the FO completion-and-gates / validation read discipline in `skills/first-officer/references/first-officer-shared-core.md` (the same family of read-cost sites e6a wired) and the dev README `## Runtime Live CI` / Testing-Resources area. Proposed addition (exact verb/flag finalized by implementation against the shipped surface):

> **Triaging a failed live CI run.** Do not read the 143KB `*-transcript.txt` whole. Fetch it (`gh run download` or `gh run view --log`) and run the transcript summarizer — `${SPACEDOCK_BIN:-spacedock} <ci-summarize> <transcript-or-->` — to get verdict + exit code + each failing test's `file:line` and symptom in ~15 lines. Open the full log only for root cause (the summary's last line names it and the `grep` to find the failing region).

### (b) User-facing CLI doc

Add a one-line row for the new helper to whatever enumerates `spacedock` subcommands/flags (the README command area and, if present, the mkdocs command-reference page) — addition only, no rewording of existing entries. Concrete wording deferred to implementation against the shipped verb/flag name.

## Stage Report: ideation

- DONE: Spike FIRST on a real large CI log (~143KB): produce a 10-20 line triage summary (pass/fail, exit code, top-N failures with file:line), then confirm an FO/ensign can triage correctly from the summary alone, and record the spike result in the body.
  Real fixtures: `go test -v ./internal/...` → 172KB/2868-line all-pass; spliced a throwaway failing module's output into it → 170KB/2842-line failing log, failures buried at L1407-1422. Throwaway parser `/tmp/ci-log-spike/summarize.go` produced an 11-line summary; oracle (`cat -n` of source) confirmed all four planted failures' exact `file:line` (a_test.go:10, :14, subtest b_test.go:8, panic b_test.go:13) + verdict/exit. A fresh FO agent given ONLY the 11 lines answered 7/7 triage questions correctly and drew the information-loss boundary (triage yes, root-cause needs full log). Mutation-checked non-tautological. Recorded in body "Riskiest-mechanism spike".
- DONE: Design the surface on that spike evidence: decide where the summary is produced (CI-emitted artifact versus a read-time `spacedock` helper) and how the full log stays reachable with no information loss, with behavior-first ACs proven over a real failing-log fixture.
  Body "Proposed approach" decides a READ-TIME `spacedock` helper (primary) over a CI-emitted artifact: the savings is context tokens not disk, B captures them over any historical log with zero CI change, and reading-not-replacing keeps the full log reachable by construction (AC-2 free). stdin form composes with `gh run view --log | spacedock … -` (no disk write). CI-side `$GITHUB_STEP_SUMMARY` emission kept as an out-of-scope thin add-on. AC-1..AC-6 are behavior-first over the frozen failing/passing fixtures (oracle = planted failures), incl. AC-3 truncated-log-does-not-falsely-PASS and AC-6 adoption-under-guards. Doc diff for the FO triage-read site included.

### Summary

Designed ci-log-read-summary as a read-time `spacedock` helper that parses `go test -v` transcript text (path or stdin) into a ≤20-line triage summary — verdict+exit, failing packages, panic frame, and every failing test with `file:line` and symptom, plus a pointer back to the full log. Chose read-time over CI-emitted because the cost is context tokens (not disk), the helper works on every historical log with no CI change, and reading-not-replacing makes "full log reachable / no information loss" free. Spiked the riskiest path FIRST on a real 170KB failing transcript: an 11-line summary captured all four planted failures' exact locations (oracle = source `cat -n`, mutation-verified non-tautological), and a fresh FO agent triaged 7/7 correctly from the summary alone while correctly identifying that root-cause still needs the full log — the empirical justification for keeping it reachable. Surfaced and pinned two real hazards the spike caught (subtest loc-attribution; panic frame in the goroutine stack) and one truncation edge (a cut-off log must not falsely report PASS). Scope held distinct from e6a's markdown-heading `status --read` grammar.
