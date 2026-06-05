package release

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowsPreserveAndPublishJourneyCosts(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	release := readWorkflow(t, "release.yml")

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(live); err != nil {
		t.Fatal(err)
	}

	if err := assertReleaseWorkflowPublishesJourneyCosts(release); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeLiveWorkflowGuardRejectsCommentOnlyJourneyMetricUpload(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.ReplaceAll(live,
		`live-artifacts/journey-metrics/**`,
		`# live-artifacts/journey-metrics/**`)

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted journey metrics paths that were only inert upload text")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsMissingCodexJourneyMetricUpload(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	first := strings.Index(live, `live-artifacts/journey-metrics/**`)
	if first < 0 {
		t.Fatal("fixture workflow missing first journey metrics upload path")
	}
	second := strings.Index(live[first+1:], `live-artifacts/journey-metrics/**`)
	if second < 0 {
		t.Fatal("fixture workflow missing second journey metrics upload path")
	}
	second += first + 1
	adversarial := live[:second] + `live-artifacts/codex/**` + live[second+len(`live-artifacts/journey-metrics/**`):]

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted a workflow where only one live job uploads journey metrics")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsMissingSharedScenarioRun(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.Replace(live,
		`go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle/ -v 2>&1 | tee claude-shared-scenarios-transcript.txt`,
		`# go test -tags live -count=1 -timeout 40m -run TestLiveClaudeSharedScenarios ./internal/ensigncycle/ -v 2>&1 | tee claude-shared-scenarios-transcript.txt`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing Claude shared scenario run command")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted a workflow without an executable Claude shared scenario run")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsUnscopedPiPackage(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.Replace(live,
		`npm install -g @earendil-works/pi-coding-agent --before="$NPM_BEFORE" --ignore-scripts --no-audit --no-fund --omit=dev`,
		`npm install -g pi-coding-agent --before="$NPM_BEFORE" --ignore-scripts --no-audit --no-fund --omit=dev`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing scoped Pi CLI install command")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted the wrong unscoped Pi CLI package")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsMissingPiBeforeAgeGate(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.ReplaceAll(live, ` --before="$NPM_BEFORE"`, ``)
	if adversarial == live {
		t.Fatal("fixture workflow missing npm --before install flags")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted Pi npm installs without --before")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsObsoletePiMinReleaseAgeProbe(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.Replace(live,
		`NPM_BEFORE="$(node -e 'console.log(new Date(Date.now() - 24*60*60*1000).toISOString())')"
          echo "Using npm --before age gate for pi-live installs: $NPM_BEFORE"`,
		`npm config get min-release-age
          npm config set min-release-age 1440`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing npm --before age-gate timestamp")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted obsolete min-release-age probing")
	}
}

func TestRuntimeLiveWorkflowGuardRejectsUnverifiedPiPackageInstall(t *testing.T) {
	live := readWorkflow(t, "runtime-live-e2e.yml")
	adversarial := strings.Replace(live,
		`node -e "const p=require('$global_npm_root/@earendil-works/pi-coding-agent/package.json'); if (p.name !== '@earendil-works/pi-coding-agent') throw new Error('unexpected Pi package name '+p.name); if (!p.bin || p.bin.pi !== 'dist/cli.js') throw new Error('unexpected Pi bin '+JSON.stringify(p.bin)); console.log('verified '+p.name+'@'+p.version+' bin pi='+p.bin.pi)"`,
		`echo "skipping Pi package verification"`,
		1)
	if adversarial == live {
		t.Fatal("fixture workflow missing Pi CLI package verification command")
	}

	if err := assertRuntimeLiveWorkflowUploadsRawJourneyMetrics(adversarial); err == nil {
		t.Fatal("runtime live workflow guard accepted an unverified Pi CLI package install")
	}
}

func TestReleaseWorkflowGuardRejectsCommentOnlyJourneyCostBuilder(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	adversarial := strings.Replace(release,
		`go run ./cmd/spacedock-release journey-costs "$RELEASE_VERSION" \`,
		`# go run ./cmd/spacedock-release journey-costs "$RELEASE_VERSION" \`,
		1)
	adversarial = strings.Replace(adversarial,
		`test -s "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`,
		`printf '{}' > "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"
          test -s "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`,
		1)

	if err := assertReleaseWorkflowPublishesJourneyCosts(adversarial); err == nil {
		t.Fatal("release workflow guard accepted a commented builder plus fake JSON output")
	}
}

func TestReleaseWorkflowGuardRejectsCommentOnlyJourneyCostPublish(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	adversarial := strings.Replace(release,
		`gh release upload "$GITHUB_REF_NAME" "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json" --clobber`,
		`# gh release upload "$GITHUB_REF_NAME" "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json" --clobber
          echo "$RUNNER_TEMP/journey-costs-v${RELEASE_VERSION}.json"`,
		1)

	if err := assertReleaseWorkflowPublishesJourneyCosts(adversarial); err == nil {
		t.Fatal("release workflow guard accepted a commented publish command")
	}
}

// TestReleaseWorkflowGuardRejectsGoreleaserNeedsJourneyLedger injects the exact
// re-coupling bug this separation closes: a `needs: journey-ledger` edge on the
// goreleaser job re-blocks the cut on the never-fired Runtime-Live-E2E run. The
// hardened guard binds `needs:` to the OWNING job (the one carrying the
// goreleaser action) and must RED on this edge — while the SAFE one-way
// `journey-ledger: needs: goreleaser` edge in the real file does NOT trip it.
// The edge is injected in every YAML shape `needs:` can take (scalar, flow
// sequence, and block list) so no single re-coupling syntax evades the guard.
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
		{"flow sequence", "    needs: [journey-ledger]\n"},
		{"block list", "    needs:\n      - journey-ledger\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			adversarial := strings.Replace(release,
				"  goreleaser:\n    runs-on: macos-latest",
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

// TestReleaseWorkflowSkipsLedgerWhenNoProducerRun proves the JOB-LEVEL skip
// consequence of POLICY 1: when the download step finds no producer run it must
// not merely exit 0 — the downstream Build/Publish-ledger steps must be GATED on
// the download step's producer-found output, so a producer-less / empty-dir cut
// SKIPS those steps (journey-ledger job green) instead of hard-failing the
// `journey-costs` builder over an empty dir (journey-ledger job RED). Parsed
// from the real release.yml, not matched against prose.
func TestReleaseWorkflowSkipsLedgerWhenNoProducerRun(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	if err := assertReleaseLedgerStepsSkipWhenNoProducerRun(release); err != nil {
		t.Fatal(err)
	}
}

// TestReleaseWorkflowGuardRejectsUngatedLedgerBuild is the adversarial twin of
// the skip-consequence proof: strip the producer-found `if:` gate off the Build
// step and the guard must RED, because an ungated builder runs `journey-costs`
// over an empty dir on a producer-less cut and REDs the journey-ledger job —
// the exact failure POLICY 1 closes.
func TestReleaseWorkflowGuardRejectsUngatedLedgerBuild(t *testing.T) {
	release := readWorkflow(t, "release.yml")
	if err := assertReleaseLedgerStepsSkipWhenNoProducerRun(release); err != nil {
		t.Fatalf("real release.yml unexpectedly fails the skip-consequence guard before mutation: %v", err)
	}
	adversarial := strings.Replace(release,
		"      - name: Build journey cost ledger\n        if: steps.download_metrics.outputs.found == 'true'\n",
		"      - name: Build journey cost ledger\n",
		1)
	if adversarial == release {
		t.Fatal("fixture workflow missing the gated Build journey cost ledger step to mutate")
	}

	if err := assertReleaseLedgerStepsSkipWhenNoProducerRun(adversarial); err == nil {
		t.Fatal("release workflow guard accepted an ungated Build step that would RED the journey-ledger job on a producer-less cut")
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
