---
id: 3g6gbbn1bvk41a57tjbe50rv
title: Upload the live-e2e transcript as a CI artifact (diagnose failures past gh log truncation)
status: validation
source: "FO (2026-06-02): 38's streamWatcher tees the FO transcript to t.Log, but on a failed run gh truncates the large job log AND the CI uploads only the spacedock binary, not the transcript — so the sonnet FO-stall on 38's PR was not fully diagnosable. Prerequisite for root-causing the headless FO-drive flake."
started: 2026-06-02T07:57:40Z
completed:
verdict:
score: "0.30"
worktree: .worktrees/spacedock-ensign-live-e2e-transcript-artifact
issue:
mod-block: merge:pr-merge
---

`38` (live-e2e-per-stage-timeouts) made the live test STREAM the FO's stream-json to `t.Log` so a hang names the stalled step. But that diagnosability is defeated downstream: on 38's own PR CI-E2E (#261) the sonnet job failed (FO stalled headless), and the streamed transcript was NOT recoverable — `gh run view --log` truncates the multi-MB streamed log, AND the CI "Upload live artifacts" step (`.github/workflows/runtime-live-e2e.yml`) uploads only the `spacedock` binary, not the test transcript. So the captured-transcript half of 38's diagnosability never reaches a human on a real failure.

This is the PREREQUISITE for root-causing the headless FO-drive flake (the sibling follow-up `headless-fo-drive-flake`): you cannot diagnose WHY the FO stalls without the failing run's transcript.

## Where the transcript lives (grounding)

`38` is already merged into `next` (commit `3a59916a`, `internal/ensigncycle/streamwatch.go` present on `next`, absent on `main`). The live workflow triggers on `pull_request` to `next` and on `workflow_dispatch`, so the test a CI-E2E job actually compiles is 38's streaming-watcher version, not the pre-38 `CombinedOutput` version still on `main`.

Under 38 the watcher's `tee` sink is `func(line string) { t.Log(line) }`: every drained stream-json line flows to `t.Log` line-by-line as it streams (streamwatch.go `drainEntries` → `w.tee(line)`). On a stalled step the watcher trips a `stepTimeout`/`stepFailure` whose `Error()` carries the step label + a transcript tail, and the test calls `t.Fatalf(... err)` — so the labelled failure tail also lands in the test framework's output. `go test … -v` surfaces all of that `t.Log`/`t.Fatalf` output on the step's **stdout**.

So the transcript's destination is already the CI step's stdout. The only thing missing is: that stdout is neither persisted to a file nor uploaded — `gh run view --log` truncates the multi-MB stream, and the upload step grabs only `./spacedock`. This task tees the step's stdout to a file and adds that file to the upload. It composes with 38; it ships no new transcript-production mechanism.

## Design

A CI-only change to `.github/workflows/runtime-live-e2e.yml`. No test-side hook — see "Decision: tee, not a test-side path" below.

**Path.** `./live-e2e-transcript.txt` in the repo root (the job's working directory), matching the existing `./spacedock` relative-path convention in the upload step. No per-model suffix is needed in the filename: each matrix leg (`sonnet` / `claude-opus-4-8`) is a distinct job with its own workspace and its own artifact name (`runtime-live-e2e-claude-live-${{ matrix.model }}`), so the two legs never share a file or an upload.

**Edit 1 — tee the streamed test output, preserving the exit code.** Replace the `Run live ensign cycle` step's `run:` body with:

```yaml
      - name: Run live ensign cycle
        env:
          SPACEDOCK_LIVE_MODEL: ${{ matrix.model }}
        run: |
          set -o pipefail
          go test -tags live -run TestLiveEnsignCycle ./internal/ensigncycle/ -v 2>&1 | tee live-e2e-transcript.txt
```

`set -o pipefail` makes the pipeline's exit status the rightmost non-zero status, so a failing `go test` (a watcher `stepTimeout`/`stepFailure` → `t.Fatalf` → non-zero exit) propagates through `tee` and the step goes RED. `2>&1` folds stderr (Go's SIGQUIT goroutine-dump on the built-in default test timeout, plus any test-framework stderr) into the same captured stream. `tee` writes line-by-line as the pipe streams, so a hung/killed run leaves a partial transcript on disk up to the kill point — capture-as-it-streams, not capture-on-success.

GitHub Actions' default `run:` shell is already `bash --noprofile --norc -eo pipefail {0}`, so `pipefail` is on by default — but the explicit `set -o pipefail` makes the AC-2 guarantee self-evident in the step and immune to a later `shell:` override silently turning it off. It is one cheap line that removes a hidden dependency on GHA's default-shell choice.

**Edit 2 — add the transcript to the existing upload.** Add the path to the `Upload live artifacts` step's `path:` list (the step already has `if: always()`, so it runs on failure):

```yaml
      - name: Upload live artifacts
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: runtime-live-e2e-claude-live-${{ matrix.model }}
          path: |
            ./spacedock
            ./live-e2e-transcript.txt
          if-no-files-found: warn
```

`if-no-files-found: warn` is retained: if a run dies so early the transcript file was never created, the upload warns rather than failing — but on any run that reached `go test`, `tee` has created the file (even if empty-then-partial), so the transcript is present on every real failure.

## Decision: tee, not a test-side path

The assignment leaves room for "a tiny test-side hook if a known transcript path is cleaner than a tee." It is not, here. 38 already routes the full transcript to `t.Log`, which `-v` puts on stdout; a CI-level `tee` captures exactly that with zero Go changes. A test-side file would duplicate the sink, introduce a second source of truth for "where the transcript is," and couple the CI artifact path into Go source — strictly more surface for no diagnostic gain. Tee wins on YAGNI and on keeping the change inside the named scope file. Recorded so the determination is on the record.

## Acceptance criteria

**AC-1 — A failed live-e2e run's artifact contains the non-empty transcript with the failure tail.** On a CI-E2E job that FAILS (e.g. a watcher `stepTimeout`/`stepFailure`), the uploaded `runtime-live-e2e-claude-live-<model>` artifact contains `live-e2e-transcript.txt`, non-empty, carrying the streamed stream-json lines AND the labelled `stepTimeout`/`stepFailure` step name + transcript tail — not just `./spacedock`.
Verified by: inspect the uploaded artifact of a failing CI-E2E run (a real flake or a forced-fail dispatch) and confirm `live-e2e-transcript.txt` is present, non-empty, and contains the labelled failure tail string (e.g. `step "..." made no progress` / `FO subprocess exited (code=...) before step`).

**AC-2 — Exit code is not swallowed by the tee; pass stays green, fail stays red.** The `tee` does NOT mask the test exit code: a passing `go test` run leaves the `Run live ensign cycle` step green; a failing run leaves it RED (the `claude-live` job fails). The file is written as the pipe streams (survives a killed run), not on clean exit only.
Verified by:
- exit-code preservation is exercised offline (machine-independent, no live credential, no GitHub) by `bash --noprofile --norc -eo pipefail -c '<failing-cmd> 2>&1 | tee f; echo "should-not-print"'` returning non-zero AND `f` containing the failing command's output. Already run during ideation (see Spike below) — it returns exit 1, the post-pipe `echo` is suppressed by `-e`, and `f` carries both lines. This is the riskiest mechanism and it is proven.
- end-to-end: a green CI-E2E run shows the step green with the transcript uploaded; a red CI-E2E run (forced or flake) shows the step red with the transcript uploaded.

**AC-3 — Scope contained to the workflow file.** The shipped diff touches only `.github/workflows/runtime-live-e2e.yml` (the two steps above). No change to `internal/status`, no change to the FO runtime, no test-side transcript hook.
Verified by: the PR diff's changed-files list is exactly `[.github/workflows/runtime-live-e2e.yml]`.

## Test plan

- **Mechanism spike (done, offline, throwaway):** `set -o pipefail` + `tee` preserves a non-zero pipeline exit while leaving the transcript file on disk. Ran during ideation; result recorded under Spike below. Cost: seconds. This is the riskiest unknown — that the tee doesn't swallow the exit code — and it is proven before the rest of the plan.
- **YAML well-formedness (offline, cheap):** the edited workflow must parse. A `yamllint`/`actionlint` or `python -c 'import yaml,sys; yaml.safe_load(open(sys.argv[1]))'` parse of the file, or a `gh workflow view` after merge, confirms the file is still valid YAML — this is a parse of a real artifact, not a substring search.
- **End-to-end (live, gated, expensive — captain-driven):** a real CI-E2E dispatch. Green path: the run passes, the step is green, the artifact carries the transcript. Red path: a forced-fail (e.g. temporarily shrink a watcher budget on a throwaway branch, or catch a real flake) leaves the step red and the artifact still carries the partial/failed transcript with the labelled tail. The live run is approval-gated and spends `ANTHROPIC_API_KEY`, so it is the captain's call when to burn it; AC-1/AC-2's end-to-end legs are verified on that run. AC-2's exit-code half and AC-3 are fully verifiable offline without a live run.
- No new Go unit test is warranted: the change adds no Go code (Decision above), and the streaming/tee/`t.Log` behavior it relies on is 38's, already covered by `streamwatch_unit_test.go` on `next`.

## Spike: pipefail + tee exit-code preservation (riskiest unknown — PROVEN)

Ran offline during ideation (mimics the GHA default shell):

```
$ bash --noprofile --norc -eo pipefail -c \
    'go() { echo "line1"; echo "FAIL line"; return 1; }; \
     go test 2>&1 | tee /tmp/t.txt; echo "exit after pipe: $?"'
line1
FAIL line
# outer exit: 1   (the post-pipe echo did NOT print — `-e` aborted on the failed pipeline → step goes RED)
$ cat /tmp/t.txt
line1
FAIL line          (transcript survived the failed pipeline)
```

Result: the failing pipeline returns exit 1 (step would go RED), and the transcript file is on disk with the streamed lines. The trade-off the AC rests on is proven. No other unverified mechanism remains — the transcript-to-`t.Log` and `stepTimeout`/`stepFailure`-to-`t.Fatalf` paths are 38's, already merged on `next` and unit-tested there; `go test -v`→stdout and `actions/upload-artifact` multi-path upload are stock behavior.

## Notes
- `.github/workflows/runtime-live-e2e.yml`. Merge this FIRST (captain) — prerequisite for the `headless-fo-drive-flake` investigation. Targets `next` (where 38 lives and where the live workflow triggers); `main` still has the pre-38 `CombinedOutput` test, so the diagnostic payoff is realized on `next`-bound PRs.
- The transcript can run multi-MB (the same volume `gh` truncates); `actions/upload-artifact` handles that fine and the artifact retention is GitHub's default. No size cap or truncation is added on purpose — the whole point is to escape `gh`'s truncation.

## Stage Report: ideation

- DONE: Design the transcript capture (CI step tees streamed `-v` output to a file surviving a FAILED/killed run; Upload step includes it; named path + exact yaml edits).
  `live-e2e-transcript.txt` at repo root; Edit 1 adds `set -o pipefail` + `2>&1 | tee live-e2e-transcript.txt` to `Run live ensign cycle`; Edit 2 adds `./live-e2e-transcript.txt` to the existing `if: always()` upload `path:` list. Both yaml blocks are in the body Design section.
- DONE: AC — failed run uploads non-empty transcript + labelled stepTimeout/stepFailure tail; passing stays green, failing stays RED (exit not swallowed).
  Hardened as AC-1 (failed-run artifact carries non-empty transcript + labelled failure tail) and AC-2 (pipefail preserves exit; pass green / fail red; tee streams so a killed run leaves a partial file). Each AC names an out-of-body check.
- DONE: Scope confined to `.github/workflows/runtime-live-e2e.yml` (+ tiny test hook only if a known path is cleaner) — NOT status lane, NOT FO runtime; PREREQUISITE merged FIRST.
  Added AC-3 (diff = only the workflow file). Recorded the explicit decision REJECTING a test-side hook (tee is strictly simpler given 38 already routes the transcript to `t.Log`). Notes restate the captain-first / `next`-target ordering.

### Summary

Grounded the design against the real tree: 38 (the streaming watcher whose `tee` sends every line to `t.Log`) is already merged on `next` (commit 3a59916a) and the live workflow triggers on `next`, so the full transcript already reaches the `go test -v` step stdout — only persistence + upload were missing. The change is two edits to `runtime-live-e2e.yml`: `set -o pipefail` + `2>&1 | tee live-e2e-transcript.txt` on the run step, and adding that file to the existing `if: always()` artifact upload. Ran the riskiest unknown offline first (pipefail+tee preserves a non-zero exit AND leaves the transcript on disk — proven, exit 1, file survived); no other unverified mechanism remains, so no Go spike beyond that and no new unit test (the change ships no Go code). Chose a CI-level tee over a test-side path on YAGNI grounds and recorded the determination.

## Stage Report: implementation

- DONE: Exactly two edits to .github/workflows/runtime-live-e2e.yml: (a) on the `Run live ensign cycle` step, add `set -o pipefail` and pipe the run through `2>&1 | tee live-e2e-transcript.txt`; (b) add `./live-e2e-transcript.txt` to the existing `if: always()` upload step's `path:` list. Use the exact yaml from the ideation Design section.
  Commit 2ec22cc4 on branch spacedock-ensign/live-e2e-transcript-artifact; verbatim from Design (run step + upload path). Ruby/Psych YAML.load_file parses the edited file clean and re-renders both edits exactly.
- DONE: Preserve the exit contract (the offline-proven mechanism): pipefail keeps the run step's non-zero exit so a FAILING live cycle stays RED, while `tee` streams so a killed run still leaves a partial transcript on disk. Re-confirm OFFLINE that `set -o pipefail; <false> 2>&1 | tee f` exits non-zero AND leaves f populated; a passing run stays green.
  Re-proven offline (no live creds/GitHub): `bash --noprofile --norc -eo pipefail -c 'gotest(){ echo line1; echo "FAIL line" >&2; return 1; }; gotest 2>&1 | tee sd_pf.txt'` → EXIT=1, sd_pf.txt = 16 bytes / 2 lines (stdout + 2>&1-folded stderr survived). Passing case (`return 0`) → OUTER_EXIT=0 and the post-pipe `echo` printed (step stays green).
- DONE: AC-3 scope: the committed diff touches ONLY runtime-live-e2e.yml — no Go code, no other workflow, no status-lane file. (No Go change, so gofmt/vet/go test are unaffected; note that in the report.)
  `git status --porcelain` and `git diff` show exactly one changed file (`.github/workflows/runtime-live-e2e.yml`, +3 -1). No Go source touched, so gofmt/go vet/go test are unaffected and were not re-run.

### Summary

Implemented the two CI-only edits to `runtime-live-e2e.yml` verbatim from the Design section: `set -o pipefail` + `2>&1 | tee live-e2e-transcript.txt` on the `Run live ensign cycle` step, and `./live-e2e-transcript.txt` added to the existing `if: always()` upload `path:` list. Re-proved the riskiest mechanism offline — pipefail keeps the failing pipeline's non-zero exit (EXIT=1, step would go RED) while `tee` leaves the streamed transcript on disk (16 bytes / 2 lines), and the passing path stays green. The diff is exactly one file (+3 -1); no Go code changed, so the Go toolchain gates are unaffected. End-to-end green/red CI-E2E legs (AC-1 and AC-2's live half) remain the captain's approval-gated live run.

## Stage Report: validation

- DONE: Verify the exit contract by RUNNING it offline (not reading yaml): reproduce the run-step shape — `set -o pipefail; <command that exits non-zero> 2>&1 | tee live-e2e-transcript.txt` exits NON-ZERO (so a failing live cycle stays RED) AND leaves live-e2e-transcript.txt populated; a zero-exit command stays exit 0. Confirm tee streams (a killed mid-run leaves a partial file).
  Ran the run-step shape under `bash --noprofile --norc -eo pipefail`: failing cmd (`return 7`) → PIPELINE_EXIT=7 (rightmost non-zero preserved → step RED) with file populated (29 bytes, 2 lines, both stdout + 2>&1-folded stderr); zero-exit cmd → PIPELINE_EXIT=0 (stays green); a process killed at ~2s of a 5s sleep left a partial file with only the pre-kill lines (`late line` absent) — tee streams as the pipe runs.
- DONE: Confirm the upload wiring: `./live-e2e-transcript.txt` is added to the existing `if: always()` upload step's `path:` list (so a failed/killed run still uploads it), and the workflow YAML is well-formed.
  Ruby/Psych `YAML.load_file` parses clean; assertions over the PARSED step values pass — `Upload live artifacts` keeps `if: always()`, `with.path` parses to `["./spacedock", "./live-e2e-transcript.txt"]`, `if-no-files-found: warn` retained; run step's parsed `run` carries `set -o pipefail`, `2>&1`, and `| tee live-e2e-transcript.txt`.
- PASSED: AC-3 scope: the diff touches ONLY .github/workflows/runtime-live-e2e.yml (no Go code, no other workflow, no status-lane file).
  Branch delta vs merge-base origin/next is exactly one file (`.github/workflows/runtime-live-e2e.yml`, +3 -1); `git diff --name-only` shows no `.go` files, no `internal/status`, no other `.github/workflows/*`. No Go code changed, so go test/vet/gofmt are unaffected and were not run.

### Summary

PASSED. The riskiest property — pipefail exit-preservation through the tee — was adversarially confirmed offline: a non-zero command propagates its exit (a failing live cycle stays RED, it does NOT go green), while a zero-exit command stays exit 0, and tee streams line-by-line so a killed run leaves a partial transcript. Upload wiring verified over parsed YAML values (not substring): `./live-e2e-transcript.txt` is in the `if: always()` upload `path:` alongside `./spacedock`, and the file is well-formed. Scope is contained to the single workflow file; no Go code changed, so the Go toolchain gates are unaffected. The end-to-end live legs of AC-1 and AC-2 remain the captain's approval-gated CI-E2E dispatch, as the spec already designates.
