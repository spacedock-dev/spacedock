package release

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedger injects the exact
// re-coupling bug this separation closes: a `needs: journey-ledger` edge on the
// goreleaser job re-blocks the cut on the never-fired Runtime-Live-E2E run. The
// hardened guard binds `needs:` to the OWNING job (the one carrying the
// goreleaser action) and must RED on this edge — while the SAFE one-way
// `journey-ledger: needs: goreleaser` edge in the real file does NOT trip it.
// The edge is injected in every YAML shape `needs:` can take — scalar, flow
// sequence, block list, the same three with a trailing inline `# comment`, and a
// block list split by a blank line — so no syntactic variation evades the guard.
func TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedger(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	if err := assertReleaseWorkflowPublishesJourneyCosts(release); err != nil {
		t.Fatalf("real release.yml unexpectedly fails the guard before mutation: %v", err)
	}

	for _, tc := range []struct {
		name      string
		needsForm string
	}{
		{"scalar", "    needs: journey-ledger\n"},
		{"flow sequence", "    needs: [e2e-gate, journey-ledger]\n"},
		{"block list", "    needs:\n      - e2e-gate\n      - journey-ledger\n"},
		{"scalar with inline comment", "    needs: journey-ledger  # required for the upload\n"},
		{"flow sequence with inline comment", "    needs: [e2e-gate, journey-ledger]  # required for the upload\n"},
		{"block list with inline comment", "    needs:\n      - e2e-gate\n      - journey-ledger  # required for the upload\n"},
		{"block list split by a blank line", "    needs:\n      - e2e-gate\n\n      - journey-ledger\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			adversarial := strings.Replace(release,
				"  goreleaser:\n    needs: e2e-gate\n    runs-on: macos-latest",
				"  goreleaser:\n"+tc.needsForm+"    runs-on: macos-latest",
				1)
			if adversarial == release {
				t.Fatal("fixture workflow missing the goreleaser job header to mutate")
			}

			if err := assertReleaseWorkflowPublishesJourneyCosts(adversarial); err == nil {
				t.Fatalf("release workflow guard accepted a goreleaser job that needs the journey-ledger job via %s form (re-coupling the cut to the never-fired run)", tc.name)
			}
		})
	}
}

// TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedgerViaJobIdentityShapes
// attacks the OTHER half of the edge — the job-identity end — with YAML shapes a
// line-walk over the raw text cannot see but the real GHA YAML resolver does:
//   - alias: an `&anchor` on journey-ledger's needs and a `*anchor` on
//     goreleaser's needs — GHA resolves the alias to `[goreleaser, journey-ledger]`
//     so goreleaser really does need journey-ledger.
//   - quoted job key: `"goreleaser":` is the same job key as `goreleaser:` to a
//     YAML parser, so a `needs: journey-ledger` under it is a real re-coupling.
//
// Both must RED. Because the whole job graph (names + alias-resolved needs +
// step attribution) comes from one yaml.v3 pass, no job-identity shape evades.
func TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedgerViaJobIdentityShapes(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	if err := assertReleaseWorkflowPublishesJourneyCosts(release); err != nil {
		t.Fatalf("real release.yml unexpectedly fails the guard before mutation: %v", err)
	}

	for _, tc := range []struct {
		name   string
		mutate func(string) string
	}{
		{
			"anchor/alias needs", func(s string) string {
				// Anchor a list containing journey-ledger on the reverse edge, then
				// alias it onto goreleaser — GHA resolves *grx so goreleaser needs
				// journey-ledger (alongside the legitimate e2e-gate it still needs).
				s = strings.Replace(s,
					"  journey-ledger:\n    needs: goreleaser\n",
					"  journey-ledger:\n    needs: &grx [e2e-gate, journey-ledger]\n", 1)
				return strings.Replace(s,
					"  goreleaser:\n    needs: e2e-gate\n    runs-on: macos-latest",
					"  goreleaser:\n    needs: *grx\n    runs-on: macos-latest", 1)
			},
		},
		{
			"quoted job key", func(s string) string {
				return strings.Replace(s,
					"  goreleaser:\n    needs: e2e-gate\n    runs-on: macos-latest",
					"  \"goreleaser\":\n    needs: [e2e-gate, journey-ledger]\n    runs-on: macos-latest", 1)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			adversarial := tc.mutate(release)
			if adversarial == release {
				t.Fatal("fixture workflow shape to mutate not found")
			}

			if err := assertReleaseWorkflowPublishesJourneyCosts(adversarial); err == nil {
				t.Fatalf("release workflow guard accepted a goreleaser→journey-ledger re-coupling via %s (the job-identity end evaded the guard)", tc.name)
			}
		})
	}
}

// goreleaserCarrierJob is a second job that ALSO runs the goreleaser action,
// used to build multi-carrier fixtures. The `%s` is its `needs:` block (may be
// empty). It must be authored before the real goreleaser job's text so the
// flat document-order builder<=goreleaser guard stays satisfied.
const goreleaserCarrierJob = `  goreleaser-extra:
%s    runs-on: macos-latest
    steps:
      - name: Run goreleaser (extra carrier)
        uses: goreleaser/goreleaser-action@v6
        with:
          version: "~> v2"
          args: release --clean
`

// insertGoreleaserCarrier authors a second goreleaser-action job (with the given
// needs block) just before the real goreleaser job in the workflow text.
func insertGoreleaserCarrier(t *testing.T, workflow, needsBlock string) string {
	t.Helper()
	const anchor = "  goreleaser:\n    needs: e2e-gate\n    runs-on: macos-latest"
	if !strings.Contains(workflow, anchor) {
		t.Fatal("fixture workflow missing the goreleaser job header to anchor a second carrier before")
	}
	carrier := fmt.Sprintf(goreleaserCarrierJob, needsBlock)
	return strings.Replace(workflow, anchor, carrier+anchor, 1)
}

// TestReleaseWorkflowGuardRejectsMultiCarrierGoreleaserNeedsJourneyLedger proves
// the guard is DETERMINISTIC across multiple goreleaser-action carriers. With two
// jobs running the goreleaser action — one declaring `needs: journey-ledger` —
// a guard that inspected only the LAST map-iterated carrier would greenlight on
// ~the fraction of runs that iterate the safe carrier last (Go map order is
// random). The collect-all guard rejects on ANY goreleaser carrier → ledger
// edge, so it must RED on EVERY run. Asserted over many iterations so a
// last-wins regression cannot hide behind a lucky map order.
func TestReleaseWorkflowGuardRejectsMultiCarrierGoreleaserNeedsJourneyLedger(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	adversarial := insertGoreleaserCarrier(t, release, "    needs: journey-ledger\n")
	if adversarial == release {
		t.Fatal("multi-carrier mutation did not apply")
	}

	const runs = 200
	for i := 0; i < runs; i++ {
		if err := assertReleaseWorkflowPublishesJourneyCosts(adversarial); err == nil {
			t.Fatalf("run %d/%d: guard accepted a multi-carrier workflow where a goreleaser-action job needs journey-ledger (last-wins map-order flakiness)", i+1, runs)
		}
	}
}

// TestReleaseWorkflowGuardToleratesMultiCarrierSafeShape is the safe twin: two
// goreleaser-action carriers, NEITHER needing the journey-ledger job (the extra
// carrier needs nothing). The collect-all guard must stay GREEN on every run —
// multiple carriers are not themselves a re-coupling.
func TestReleaseWorkflowGuardToleratesMultiCarrierSafeShape(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	safe := insertGoreleaserCarrier(t, release, "")
	if safe == release {
		t.Fatal("multi-carrier safe mutation did not apply")
	}

	const runs = 200
	for i := 0; i < runs; i++ {
		if err := assertReleaseWorkflowPublishesJourneyCosts(safe); err != nil {
			t.Fatalf("run %d/%d: guard wrongly rejected a safe multi-carrier workflow (extra goreleaser carrier needs nothing): %v", i+1, runs, err)
		}
	}
}

// TestReleaseWorkflowGuardToleratesSafeReverseEdgeInEveryShape is the direction
// twin of the rejection test: the SAFE one-way edge `journey-ledger: needs:
// goreleaser` (the upload waits for the Release goreleaser creates) must stay
// GREEN no matter which YAML shape it is written in. The guard binds `needs:` to
// the OWNING job, so rewriting the real file's reverse edge into scalar, flow,
// block-list, their inline-comment variants, or a blank-line-split block list
// must NOT trip the goreleaser→journey-ledger check.
func TestReleaseWorkflowGuardToleratesSafeReverseEdgeInEveryShape(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	const realEdge = "  journey-ledger:\n    needs: goreleaser\n"
	if !strings.Contains(release, realEdge) {
		t.Fatalf("fixture workflow missing the journey-ledger needs: goreleaser edge to rewrite")
	}

	for _, tc := range []struct {
		name      string
		needsForm string
	}{
		// The flow/block/comment/blank-line shapes differ from the real file's
		// scalar edge, so the rewrite is exercised; the scalar baseline itself is
		// the real file the parent guard already passes.
		{"flow sequence", "    needs: [goreleaser]\n"},
		{"block list", "    needs:\n      - goreleaser\n"},
		{"scalar with inline comment", "    needs: goreleaser  # upload waits for the Release\n"},
		{"flow sequence with inline comment", "    needs: [goreleaser]  # upload waits for the Release\n"},
		{"block list with inline comment", "    needs:\n      - goreleaser  # upload waits for the Release\n"},
		{"block list split by a blank line", "    needs:\n      - goreleaser\n\n      - some-other-gate\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rewritten := strings.Replace(release, realEdge, "  journey-ledger:\n"+tc.needsForm, 1)
			if rewritten == release {
				t.Fatal("rewrite of the journey-ledger needs edge did not apply")
			}

			if err := assertReleaseWorkflowPublishesJourneyCosts(rewritten); err != nil {
				t.Fatalf("release workflow guard wrongly rejected the SAFE reverse edge via %s form: %v", tc.name, err)
			}
		})
	}
}

// TestReleaseWorkflowGuardToleratesSafeReverseEdgeViaJobIdentityShapes is the
// direction twin of the job-identity rejection test: an anchor/alias on the SAFE
// reverse edge, or a quoted journey-ledger job key, must NOT trip the guard —
// the edge still points journey-ledger → goreleaser, the required upload order,
// not the cut-blocking reverse.
func TestReleaseWorkflowGuardToleratesSafeReverseEdgeViaJobIdentityShapes(t *testing.T) {
	release := readWorkflow(t, "release.yml")

	for _, tc := range []struct {
		name   string
		mutate func(string) string
	}{
		{
			"anchor/alias needs", func(s string) string {
				// Anchor goreleaser on a one-off carrier job, alias it onto the
				// journey-ledger reverse edge — still journey-ledger → goreleaser.
				return strings.Replace(s,
					"  journey-ledger:\n    needs: goreleaser\n",
					"  journey-ledger:\n    needs: &grx [goreleaser]\n", 1)
			},
		},
		{
			"quoted job key", func(s string) string {
				return strings.Replace(s,
					"  journey-ledger:\n    needs: goreleaser\n",
					"  \"journey-ledger\":\n    needs: goreleaser\n", 1)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rewritten := tc.mutate(release)
			if rewritten == release {
				t.Fatal("rewrite of the safe reverse edge did not apply")
			}

			if err := assertReleaseWorkflowPublishesJourneyCosts(rewritten); err != nil {
				t.Fatalf("release workflow guard wrongly rejected the SAFE reverse edge via %s: %v", tc.name, err)
			}
		})
	}
}

// TestReleaseDownloadSkipBranchExitsZeroOnEmptyRunList EXERCISES the download
// step's missing-run branch: it runs the REAL download script extracted from
// release.yml against a stubbed `gh` that returns an empty run list, and asserts
// the script exits 0 (non-fatal skip, not the old exit 1) and emits
// `found=false` to $GITHUB_OUTPUT so the downstream gate skips. This is the
// behavior-level proof of the skip branch, not a substring check.
func TestReleaseDownloadSkipBranchExitsZeroOnEmptyRunList(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	script := extractStepRun(t, release, "Download latest journey metrics artifacts")

	dir := t.TempDir()
	// Stub `gh` on PATH so the download step's `gh run list ... --jq '.[0].databaseId'`
	// yields the empty-result string the real CLI returns when no run matches.
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	ghStub := "#!/bin/sh\nexit 0\n"
	if err := os.WriteFile(filepath.Join(binDir, "gh"), []byte(ghStub), 0o755); err != nil {
		t.Fatalf("write gh stub: %v", err)
	}
	outPath := filepath.Join(dir, "github_output")
	if err := os.WriteFile(outPath, nil, 0o644); err != nil {
		t.Fatalf("seed github_output: %v", err)
	}

	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"RUNNER_TEMP="+dir,
		"GITHUB_OUTPUT="+outPath,
		"GH_TOKEN=stub",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("download skip branch exited non-zero (want exit 0 on empty run list): %v\n%s", err, out)
	}

	gotOutput, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read github_output: %v", err)
	}
	if !strings.Contains(string(gotOutput), "found=false") {
		t.Errorf("download skip branch did not emit found=false to $GITHUB_OUTPUT; got:\n%s", gotOutput)
	}
}

// TestReleaseDownloadSkipBranchToleratesGhError EXERCISES the download step's
// degradation when the `gh run list` query itself fails — a gh error or a
// missing gh on PATH. Under `set -euo pipefail` an un-guarded command
// substitution would abort the step (RED the journey-ledger job); the `|| true`
// degrades it to an empty run_id so the no-run skip branch fires (exit 0,
// found=false) and the cut is never blocked by a best-effort ledger query.
func TestReleaseDownloadSkipBranchToleratesGhError(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	script := extractStepRun(t, release, "Download latest journey metrics artifacts")

	for _, tc := range []struct {
		name   string
		ghStub string // empty ghStub means: do not put gh on PATH at all
	}{
		{"gh exits non-zero", "#!/bin/sh\necho 'gh: API rate limit exceeded' >&2\nexit 1\n"},
		{"gh missing from PATH", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			binDir := filepath.Join(dir, "bin")
			if err := os.MkdirAll(binDir, 0o755); err != nil {
				t.Fatalf("mkdir bin: %v", err)
			}
			// PATH carries system utilities (mkdir/find/awk) plus binDir, but NOT
			// the directories a real gh lives in — so the gh seen by the script is
			// exactly the stub we write here, or nothing when we write none.
			path := binDir + string(os.PathListSeparator) + "/usr/bin" + string(os.PathListSeparator) + "/bin"
			if tc.ghStub != "" {
				if err := os.WriteFile(filepath.Join(binDir, "gh"), []byte(tc.ghStub), 0o755); err != nil {
					t.Fatalf("write gh stub: %v", err)
				}
			}
			outPath := filepath.Join(dir, "github_output")
			if err := os.WriteFile(outPath, nil, 0o644); err != nil {
				t.Fatalf("seed github_output: %v", err)
			}

			cmd := exec.Command("bash", "-c", script)
			cmd.Dir = dir
			cmd.Env = append(os.Environ(),
				"PATH="+path,
				"RUNNER_TEMP="+dir,
				"GITHUB_OUTPUT="+outPath,
				"GH_TOKEN=stub",
			)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("download step aborted on a gh query failure (want exit 0, found=false): %v\n%s", err, out)
			}

			gotOutput, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("read github_output: %v", err)
			}
			if !strings.Contains(string(gotOutput), "found=false") {
				t.Errorf("download step did not emit found=false to $GITHUB_OUTPUT on a gh query failure; got:\n%s", gotOutput)
			}
		})
	}
}

// ghDownloadStepStub is the stubbed `gh` binary for the download step's
// positive multi-run path: `gh release list` returns a since-tag anchor (not
// the epoch fallback), `gh run list` returns two run ids, and `gh run
// download` populates each run's own artifacts directory with a SAME-filename,
// DIFFERENT-content journey-metrics record — the exact collision shape
// journeymetrics.recordFilename produces for two runs of the same
// scenario/model (no run-distinguishing component), so the download step's own
// per-run-subdirectory copy is what has to keep them apart, not an accidental
// filename difference.
const ghDownloadStepStub = `#!/bin/sh
case "$1" in
  release)
    echo "2026-06-01T00:00:00Z"
    ;;
  run)
    case "$2" in
      list)
        echo "1001"
        echo "1002"
        ;;
      download)
        run_id="$3"
        shift 3
        dir=""
        while [ $# -gt 0 ]; do
          case "$1" in
            --dir) dir="$2"; shift 2 ;;
            *) shift ;;
          esac
        done
        mkdir -p "$dir/journey-metrics"
        printf '{"scenario_id":"shallow-boot-window","run_id":"%s"}' "$run_id" > "$dir/journey-metrics/shallow-boot-window--claude--llm--llm-live--claude-sonnet-4-6--measured.json"
        ;;
    esac
    ;;
esac
exit 0
`

// runDownloadStep EXECUTES script (the extracted, possibly mutated, "Download
// latest journey metrics artifacts" run block) against ghDownloadStepStub's two
// stubbed runs and returns the RUNNER_TEMP directory the script populated, so
// callers can inspect the resulting $RUNNER_TEMP/journey-metrics/ layout.
func runDownloadStep(t *testing.T, script string) string {
	t.Helper()
	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	if err := os.WriteFile(filepath.Join(binDir, "gh"), []byte(ghDownloadStepStub), 0o755); err != nil {
		t.Fatalf("write gh stub: %v", err)
	}
	outPath := filepath.Join(dir, "github_output")
	if err := os.WriteFile(outPath, nil, 0o644); err != nil {
		t.Fatalf("seed github_output: %v", err)
	}

	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"RUNNER_TEMP="+dir,
		"GITHUB_OUTPUT="+outPath,
		"GH_TOKEN=stub",
		"GITHUB_REF_NAME=v0.99.0",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("download step exited non-zero with two stubbed runs: %v\n%s", err, out)
	}

	gotOutput, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read github_output: %v", err)
	}
	if !strings.Contains(string(gotOutput), "found=true") {
		t.Fatalf("download step did not emit found=true with two stubbed runs; got:\n%s", gotOutput)
	}
	return dir
}

// TestReleaseDownloadStepProducesPerRunSubdirectories EXERCISES the real
// download script's positive multi-run path (the AC-2 headline mechanism this
// task exists to ship): two stubbed runs, each carrying a journey-metrics
// record with the SAME on-disk filename but different content. It asserts the
// script's OWN per-run-subdirectory copy — not merely its Go-side consumer
// (journeymetrics.AggregateLedger) — keeps both runs' records intact and
// separate, so a positive multi-run cut actually aggregates N observations
// instead of collapsing to one via silent overwrite.
func TestReleaseDownloadStepProducesPerRunSubdirectories(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	script := extractStepRun(t, release, "Download latest journey metrics artifacts")

	runnerTemp := runDownloadStep(t, script)

	for _, tc := range []struct{ runID, wantContains string }{
		{"1001", `"run_id":"1001"`},
		{"1002", `"run_id":"1002"`},
	} {
		path := filepath.Join(runnerTemp, "journey-metrics", tc.runID, "shallow-boot-window--claude--llm--llm-live--claude-sonnet-4-6--measured.json")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("run %s's record missing from its own subdirectory: %v", tc.runID, err)
		}
		if !strings.Contains(string(data), tc.wantContains) {
			t.Fatalf("run %s's record does not carry its own run id (overwritten by another run's content) — %s:\n%s", tc.runID, path, data)
		}
	}
}

// TestReleaseDownloadStepGuardRejectsFlatCopyRegression is the adversarial
// twin: reverting the per-run cp target back to the flat "$RUNNER_TEMP/journey-
// metrics/" directory (the exact bug AC-2 fixes, and the mutation the
// validation-stage adversarial audit applied by hand) must collapse the two
// stubbed runs' same-named records into ONE file. This proves
// TestReleaseDownloadStepProducesPerRunSubdirectories actually discriminates a
// fixed script from a broken one, rather than passing vacuously.
func TestReleaseDownloadStepGuardRejectsFlatCopyRegression(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	const fixedCopy = `find "$RUNNER_TEMP/runtime-live-artifacts/$run_id" -path '*/journey-metrics/*.json' -type f -exec cp {} "$run_dir/" \;`
	const flatCopy = `find "$RUNNER_TEMP/runtime-live-artifacts/$run_id" -path '*/journey-metrics/*.json' -type f -exec cp {} "$RUNNER_TEMP/journey-metrics/" \;`
	adversarial := strings.Replace(release, fixedCopy, flatCopy, 1)
	if adversarial == release {
		t.Fatal("fixture workflow missing the per-run cp target to mutate")
	}
	script := extractStepRun(t, adversarial, "Download latest journey metrics artifacts")

	runnerTemp := runDownloadStep(t, script)

	matches, err := filepath.Glob(filepath.Join(runnerTemp, "journey-metrics", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("flat-copy regression should collapse both runs' same-named records into exactly one flat file, got %d: %v", len(matches), matches)
	}
	for _, runID := range []string{"1001", "1002"} {
		subdirMatches, err := filepath.Glob(filepath.Join(runnerTemp, "journey-metrics", runID, "*.json"))
		if err != nil {
			t.Fatal(err)
		}
		if len(subdirMatches) != 0 {
			t.Fatalf("run %s's subdirectory unexpectedly has records under the flat-copy regression: %v", runID, subdirMatches)
		}
	}
}

// TestJourneyDeltaLocateMetricsStepFindsNestedJourneyMetricsFiles is the AC-3
// production-bug proof: runtime-live-e2e.yml's journey-delta-comment job used to
// hardcode --metrics-dir to an exact subpath under the download-artifact
// destination, but a real run's downloaded artifact zip nests the journey-metrics
// JSON several directories deeper (verified against run 28432388663) — the exact
// hardcoded path is EMPTY, so journeymetrics.ReadRecordsDir errored and REDed
// every PR's delta-comment job under set -euo pipefail. This exercises the REAL
// "Locate this run's journey metrics" step (extracted from the live workflow,
// not a reimplementation) against that realistic nested layout and proves it
// finds the file regardless of the exact nesting depth.
func TestJourneyDeltaLocateMetricsStepFindsNestedJourneyMetricsFiles(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	script := extractStepRun(t, live, "Locate this run's journey metrics")

	dir := t.TempDir()
	// The verified real nesting: several directories deeper than the
	// "live-artifacts/journey-metrics/" subpath the removed hardcoded
	// --metrics-dir assumed.
	nested := filepath.Join(dir, "current-run-artifacts", "spacedock", "spacedock", "live-artifacts", "journey-metrics", "claude", "claude-opus-4-8")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	metricsFile := filepath.Join(nested, "shallow-boot--claude--llm--llm-live--claude-opus-4-8--measured.json")
	if err := os.WriteFile(metricsFile, []byte(`{"scenario_id":"shallow-boot"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	// Prove the bug was real: the OLD hardcoded path finds nothing in this
	// realistic layout.
	oldHardcodedPath := filepath.Join(dir, "current-run-artifacts", "live-artifacts", "journey-metrics")
	if _, err := os.Stat(oldHardcodedPath); err == nil {
		t.Fatalf("test fixture is unrealistic: the OLD hardcoded path %s must NOT exist", oldHardcodedPath)
	}

	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "RUNNER_TEMP="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("locate-metrics step exited non-zero: %v\n%s", err, out)
	}

	matches, err := filepath.Glob(filepath.Join(dir, "current-run-metrics", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected the nested journey-metrics JSON to be located and copied into current-run-metrics, got %d matches: %v", len(matches), matches)
	}
}

// TestJourneyDeltaLocateAndPostStepsShareMetricsDir closes the cross-step
// wiring gap validation's cycle-3 adversarial audit found: the "Locate this
// run's journey metrics" step's own mechanism is tested in isolation above,
// but nothing previously verified that the "Post the journey-cost delta PR
// comment" step's --metrics-dir argument stays wired to the Locate step's
// actual output directory. A one-line regression reverting ONLY the Post
// step's --metrics-dir back to the old hardcoded broken path left the full
// suite green, because each step's script was only ever checked in isolation.
// This derives BOTH directories from the REAL extracted step scripts (not a
// hardcoded assumption on either side) and asserts they're the same value.
func TestJourneyDeltaLocateAndPostStepsShareMetricsDir(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	locateScript := extractStepRun(t, live, "Locate this run's journey metrics")
	postScript := extractStepRun(t, live, "Post the journey-cost delta PR comment")

	locateDir := firstQuotedArg(t, locateScript, `mkdir -p "`)
	postDir := firstQuotedArg(t, postScript, `--metrics-dir "`)
	if locateDir != postDir {
		t.Fatalf("Locate step writes to %q but Post step's --metrics-dir reads from %q — the two steps are no longer wired together", locateDir, postDir)
	}
}

// firstQuotedArg extracts the double-quoted value immediately following the
// given prefix (e.g. `mkdir -p "`) in script, failing the test if the prefix
// or its closing quote is not found.
func firstQuotedArg(t *testing.T, script, prefix string) string {
	t.Helper()
	start := strings.Index(script, prefix)
	if start < 0 {
		t.Fatalf("script does not contain %q:\n%s", prefix, script)
	}
	start += len(prefix)
	end := strings.Index(script[start:], `"`)
	if end < 0 {
		t.Fatalf("unterminated quoted argument after %q:\n%s", prefix, script)
	}
	return script[start : start+end]
}

// extractStepRun pulls a named step's `run: |` block out of a workflow document,
// dedented to a runnable shell script, so tests exercise the EXACT script CI
// runs rather than a hand-copied duplicate.
func extractStepRun(t *testing.T, workflow, stepName string) string {
	t.Helper()
	for _, step := range parseWorkflowSteps(workflow) {
		if step.name == stepName {
			if step.run == "" {
				t.Fatalf("step %q has no run block", stepName)
			}
			return dedent(step.run)
		}
	}
	t.Fatalf("workflow has no step named %q", stepName)
	return ""
}

// dedent strips the common leading-space indent off every non-blank line of a
// `run:` block so the extracted YAML script runs under a shell.
func dedent(block string) string {
	lines := strings.Split(block, "\n")
	min := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		n := len(line) - len(strings.TrimLeft(line, " "))
		if min < 0 || n < min {
			min = n
		}
	}
	if min <= 0 {
		return block
	}
	for i, line := range lines {
		if len(line) >= min {
			lines[i] = line[min:]
		}
	}
	return strings.Join(lines, "\n")
}

func readWorkflow(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
