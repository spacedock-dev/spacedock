//go:build live

package ensigncycle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

func TestSharedScenarioRunnerCoverageFinal(t *testing.T) {
	want := []string{"full-ensign-cycle", "gate-guardrail", "default-headless-gate-stop", "withdrawn-gate-recovery", "recorded-gate-lifecycle", "rejection-flow", "feedback-3-cycle-escalation", "merge-hook-guardrail", "filing", "shallow-boot", "zero-discovery", "auto-continue-after-implementation", "self-evidence-merge-triage", "smallest-sufficient-mechanism", "keep-moving-posture", "ac-value-reanchor"}

	var scenarios []string
	for _, scenario := range sharedRuntimeScenarios() {
		scenarios = append(scenarios, scenario.name)
	}
	if !reflect.DeepEqual(scenarios, want) {
		t.Fatalf("shared scenarios = %v, want %v", scenarios, want)
	}

	runners := sharedScenarioRunners()
	if len(runners) != len(want) {
		t.Fatalf("shared runner count = %d, want %d", len(runners), len(want))
	}
	for _, name := range want {
		if runners[name] == nil {
			t.Errorf("shared journey %q has no runner", name)
		}
	}
	for name := range runners {
		if !containsString(want, name) {
			t.Errorf("runner %q has no registry journey", name)
		}
	}
}

func TestSharedLiveRuntimeSelection(t *testing.T) {
	for _, runtime := range []string{"claude", "codex", "pi"} {
		t.Run(runtime, func(t *testing.T) {
			adapter, err := selectSharedLiveRuntime(runtime)
			if err != nil {
				t.Fatal(err)
			}
			if got := adapter.runtimeName(); got != runtime {
				t.Fatalf("adapter runtime = %q, want %q", got, runtime)
			}
		})
	}
	if _, err := selectSharedLiveRuntime(""); err == nil {
		t.Fatal("empty SPACEDOCK_LIVE_RUNTIME accepted")
	}
	if _, err := selectSharedLiveRuntime("other"); err == nil {
		t.Fatal("unknown SPACEDOCK_LIVE_RUNTIME accepted")
	}
}

func TestSharedLiveEvidenceTargets(t *testing.T) {
	t.Setenv("SPACEDOCK_LIVE_MODEL", "sonnet")
	claude, err := selectSharedLiveRuntime("claude")
	if err != nil {
		t.Fatal(err)
	}
	if got := claude.liveEvidenceTarget(); got != liveEvidenceTargetClaudeSonnet {
		t.Fatalf("Sonnet evidence target = %q, want %q", got, liveEvidenceTargetClaudeSonnet)
	}
	t.Setenv("SPACEDOCK_LIVE_MODEL", "claude-opus-4-8")
	if got := claude.liveEvidenceTarget(); got != liveEvidenceTargetClaudeOpus {
		t.Fatalf("Opus evidence target = %q, want %q", got, liveEvidenceTargetClaudeOpus)
	}
	for runtime, want := range map[string]liveEvidenceTarget{
		"codex": liveEvidenceTargetCodex,
		"pi":    liveEvidenceTargetPi,
	} {
		adapter, err := selectSharedLiveRuntime(runtime)
		if err != nil {
			t.Fatal(err)
		}
		if got := adapter.liveEvidenceTarget(); got != want {
			t.Errorf("%s evidence target = %q, want %q", runtime, got, want)
		}
	}
}

func TestPromotedCommonJourneyEntrypoints(t *testing.T) {
	for _, name := range []string{
		"full-ensign-cycle",
		"default-headless-gate-stop",
		"withdrawn-gate-recovery",
		"zero-discovery",
		"auto-continue-after-implementation",
		"ac-value-reanchor",
	} {
		if sharedScenarioRunners()[name] == nil {
			t.Errorf("promoted journey %q is not selected by the common runner", name)
		}
	}
}

func TestSharedScenarioSequenceStopsAfterFirstFailure(t *testing.T) {
	scenarios := []sharedRuntimeScenario{{name: "first"}, {name: "second"}}
	var ran []string
	runSharedScenarioSequence(scenarios, func(scenario sharedRuntimeScenario) bool {
		ran = append(ran, scenario.name)
		return false
	})
	if want := []string{"first"}; !reflect.DeepEqual(ran, want) {
		t.Fatalf("ran scenarios = %v, want %v", ran, want)
	}
}

func TestSharedCodexAndPiDriversPreserveSpacedockShimAfterFrontDoorPin(t *testing.T) {
	shimDir := t.TempDir()
	want := filepath.Join(shimDir, "spacedock")

	t.Run("codex", func(t *testing.T) {
		realHostDir := t.TempDir()
		realHost := filepath.Join(realHostDir, "codex")
		writeFile(t, realHost, "#!/bin/sh\nprintf %s \"$SPACEDOCK_BIN\"\n")
		if err := os.Chmod(realHost, 0o755); err != nil {
			t.Fatal(err)
		}
		driver := codexAsLiveDriver{runner: codexLiveRunner{
			codexBin: realHost,
			env:      []string{"PATH=" + realHostDir + ":/usr/bin:/bin"},
		}}
		configured := driver.withStubPATH(t, shimDir).(codexAsLiveDriver)
		cmd := exec.Command("/bin/sh", "-c", "codex")
		cmd.Env = withRecordedGateEnv(configured.runner.env, "SPACEDOCK_BIN", "/real/spacedock")
		got, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("launch Codex after front-door pin: %v\n%s", err, got)
		}
		if string(got) != want {
			t.Fatalf("Codex child SPACEDOCK_BIN = %q, want scenario shim %q", got, want)
		}
		if err := os.Remove(filepath.Join(shimDir, "codex")); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command("/bin/sh", "-c", "codex")
		cmd.Env = withRecordedGateEnv(configured.runner.env, "SPACEDOCK_BIN", "/real/spacedock")
		if output, err := cmd.CombinedOutput(); err != nil || string(output) != "/real/spacedock" {
			t.Fatalf("Codex removal control did not expose the front-door pin: output=%q err=%v", output, err)
		}
	})

	t.Run("pi", func(t *testing.T) {
		driver := piSharedLiveDriver{env: []string{"PATH=/usr/bin:/bin"}}
		configured := driver.withStubPATH(t, shimDir).(piSharedLiveDriver)
		cmd := exec.Command("/bin/bash", "-c", `printf %s "$SPACEDOCK_BIN"`)
		cmd.Env = withRecordedGateEnv(configured.env, "SPACEDOCK_BIN", "/real/spacedock")
		got, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("launch Pi shell after front-door pin: %v\n%s", err, got)
		}
		if string(got) != want {
			t.Fatalf("Pi child SPACEDOCK_BIN = %q, want scenario shim %q", got, want)
		}
		cmd = exec.Command("/bin/bash", "-c", `printf %s "$SPACEDOCK_BIN"`)
		withoutPropagation := withoutEnvKey(withoutEnvKey(configured.env, "BASH_ENV"), "ZDOTDIR")
		cmd.Env = withRecordedGateEnv(withoutPropagation, "SPACEDOCK_BIN", "/real/spacedock")
		if output, err := cmd.CombinedOutput(); err != nil || string(output) != "/real/spacedock" {
			t.Fatalf("Pi removal control did not expose the front-door pin: output=%q err=%v", output, err)
		}
	})
}

func TestSharedLiveTODOEvidenceSet(t *testing.T) {
	targets := []liveEvidenceTarget{liveEvidenceTargetClaudeSonnet, liveEvidenceTargetClaudeOpus, liveEvidenceTargetCodex, liveEvidenceTargetPi}
	want := map[liveEvidenceKey]string{{target: liveEvidenceTargetCodex, journey: "full-ensign-cycle"}: codexEnsignContractDefectID, {target: liveEvidenceTargetClaudeSonnet, journey: "default-headless-gate-stop"}: defaultHeadlessGateStopDefectID, {target: liveEvidenceTargetPi, journey: "default-headless-gate-stop"}: defaultHeadlessGateStopDefectID, {target: liveEvidenceTargetClaudeSonnet, journey: "smallest-sufficient-mechanism"}: liveDurableJourneyDefectID, {target: liveEvidenceTargetClaudeSonnet, journey: "keep-moving-posture"}: liveDurableJourneyDefectID, {target: liveEvidenceTargetCodex, journey: "smallest-sufficient-mechanism"}: liveDurableJourneyDefectID, {target: liveEvidenceTargetCodex, journey: "keep-moving-posture"}: liveDurableJourneyDefectID, {target: liveEvidenceTargetPi, journey: "rejection-flow"}: liveRejectionFlowDefectID, {target: liveEvidenceTargetCodex, journey: "withdrawn-gate-recovery"}: liveWithdrawnGateDefectID, {target: liveEvidenceTargetCodex, journey: "rejection-flow"}: liveRejectionFlowDefectID}
	for _, target := range targets {
		t.Run(string(target), func(t *testing.T) {
			found := 0
			for _, scenario := range sharedRuntimeScenarios() {
				key := liveEvidenceKey{target: target, journey: scenario.name}
				reason := liveDurableJourneyTODO(target, scenario.name)
				owner, missing := want[key]
				if missing {
					if reason == "" || !strings.HasPrefix(reason, "TODO("+owner+"):") {
						t.Errorf("TODO %s/%s reason = %q, want exact owner TODO(%s)", target, scenario.name, reason, owner)
					}
					found++
					continue
				}
				if reason != "" {
					t.Errorf("runnable %s/%s unexpectedly has missing-evidence TODO %q", target, scenario.name, reason)
				}
			}
			expected := 0
			for key := range want {
				if key.target == target {
					expected++
				}
			}
			if found != expected {
				t.Fatalf("%s TODO bindings found = %d, want %d", target, found, expected)
			}
		})
	}
}

func TestAutoContinueReadmesAreDiscoverable(t *testing.T) {
	readmes := map[string]func() string{"auto-continue/single-root": autoContinueReadme, "auto-continue/split-root": piAutoContinueReadme}
	for id, readme := range readmes {
		t.Run(id, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, filepath.Join(root, "README.md"), readme())
			if discovered, found := status.DiscoverWorkflowDir(root); !found || discovered != root {
				t.Fatalf("DiscoverWorkflowDir(%s) = %q, %t; want exact fixture root", id, discovered, found)
			}

			withoutMarker := strings.Replace(readme(), "commissioned-by: spacedock@1\n", "", 1)
			if withoutMarker == readme() {
				t.Fatal("commissioned-by marker removal control did not apply")
			}
			control := t.TempDir()
			writeFile(t, filepath.Join(control, "README.md"), withoutMarker)
			if discovered, found := status.DiscoverWorkflowDir(control); found {
				t.Fatalf("markerless %s unexpectedly discovered at %q", id, discovered)
			}
		})
	}
}

func TestAutoContinueFixtureGitBaselines(t *testing.T) {
	t.Run("single-root", func(t *testing.T) {
		root := t.TempDir()
		entityPath := writeAutoContinueWorkflow(t, root)
		assertAutoContinueGitBaseline(t, root, root, entityPath, false)
	})
	t.Run("split-root", func(t *testing.T) {
		root, stateRoot, entityPath := writePiAutoContinueWorkflow(t)
		assertAutoContinueGitBaseline(t, root, stateRoot, entityPath, true)
	})
	t.Run("no-git-removal-controls", func(t *testing.T) {
		single := t.TempDir()
		entityPath, err := writeAutoContinueWorkflowNoGit(single)
		if err != nil {
			t.Fatal(err)
		}
		if autoContinueGitBaselineError(single, single, entityPath, false) == nil {
			t.Fatal("single-root fixture without Git passed the baseline check")
		}

		split := t.TempDir()
		stateRoot, entityPath, err := writePiAutoContinueWorkflowNoGit(split)
		if err != nil {
			t.Fatal(err)
		}
		if autoContinueGitBaselineError(split, stateRoot, entityPath, true) == nil {
			t.Fatal("split-root fixture without Git passed the baseline check")
		}
	})
}

type autoContinueLaunch struct {
	fixtureID     string
	artifactLabel string
	discoverable  bool
	gitClean      bool
}

type recordingAutoContinueDriver struct {
	launches []autoContinueLaunch
	result   *liveResult
	homeDir  string
}

func (d *recordingAutoContinueDriver) run(t *testing.T, scenario sharedRuntimeScenario, root, _ string) liveResult {
	t.Helper()
	if d.result != nil {
		d.result.artifactDir = t.TempDir()
		return *d.result
	}
	readme := readFile(t, filepath.Join(root, "README.md"))
	fixtureID := "auto-continue/single-root"
	stateRoot := root
	entityPath := filepath.Join(root, "auto-continue-task.md")
	if strings.Contains(readme, "state: .spacedock-state") {
		fixtureID = "auto-continue/split-root"
		stateRoot = filepath.Join(root, ".spacedock-state")
		entityPath = filepath.Join(stateRoot, "auto-continue-task", "index.md")
	}
	_, discoverable := status.DiscoverWorkflowDir(root)
	d.launches = append(d.launches, autoContinueLaunch{fixtureID: fixtureID, artifactLabel: scenario.name, discoverable: discoverable, gitClean: autoContinueGitBaselineError(root, stateRoot, entityPath, stateRoot != root) == nil})
	after := strings.Replace(readFile(t, entityPath), "status: implementation", "status: validation", 1) +
		"\n## Stage Report: validation\n\n- DONE: Validate the fixture\n  Validation passed.\n"
	writeFile(t, entityPath, after)
	return liveResult{artifactDir: scenario.name}
}

func (d *recordingAutoContinueDriver) model() string { return "fake" }
func (d *recordingAutoContinueDriver) home() string  { return d.homeDir }
func (d *recordingAutoContinueDriver) withStubPATH(*testing.T, string) liveDriver {
	return d
}

func TestAutoContinueCommonRunnerLaunchesBothVariantsSerially(t *testing.T) {
	driver := &recordingAutoContinueDriver{}
	runAutoContinueJourney(t, driver, sharedRuntimeScenario{name: "auto-continue-after-implementation"})
	want := []string{"auto-continue/single-root", "auto-continue/split-root"}
	got := make([]string, 0, len(driver.launches))
	for _, launch := range driver.launches {
		got = append(got, launch.fixtureID)
		if !launch.discoverable {
			t.Errorf("%s launch was not discoverable", launch.fixtureID)
		}
		if !launch.gitClean {
			t.Errorf("%s launch lacked its committed clean Git/state-root baseline", launch.fixtureID)
		}
		if !strings.Contains(launch.artifactLabel, launch.fixtureID) {
			t.Errorf("%s artifact label %q does not contain its fixture ID", launch.fixtureID, launch.artifactLabel)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("auto-continue fixture launches = %v, want exactly %v in serial order", got, want)
	}
}

func TestShallowBootRuntimeRouting(t *testing.T) {
	greet := shallowBootHeldGateLine + "\n" + shallowBootEngageHintLine
	piGreet := fmt.Sprintf(`{"type":"message","message":{"role":"assistant","content":[{"type":"text","text":%q}]}}`, greet)
	for name, result := range map[string]liveResult{"pi native session": {runtime: "pi", stream: "non-Claude Pi diagnostics", sessionJSONL: piFilingCall("call_boot", "spacedock status") + "\n" + piGreet, finalMessage: greet}, "claude stream": {runtime: "claude", stream: readMeasureFixture(t, "shallow-boot-greet.stream.jsonl"), sessionJSONL: "not Pi evidence", finalMessage: greet}} {
		t.Run(name, func(t *testing.T) {
			driver := &recordingAutoContinueDriver{result: &result, homeDir: t.TempDir()}
			runClaudeShallowBootScenario(t, driver, sharedRuntimeScenario{name: "shallow-boot"})
		})
	}
}

func TestResolveAutoContinueEndState(t *testing.T) {
	validated := func(body string) string {
		return body + "\n## Stage Report: validation\n\n- DONE: Validate the fixture\n  Validation passed.\n"
	}

	t.Run("single-root active worktree", func(t *testing.T) {
		root := t.TempDir()
		worktree := ".worktrees/spacedock-ensign-auto-continue-task"
		pipeline := strings.Replace(autoContinueEntity(), "status: implementation", "status: validation", 1)
		pipeline = strings.Replace(pipeline, "worktree:\n", "worktree: "+worktree+"\n", 1)
		active := validated(pipeline)
		writeFile(t, filepath.Join(root, worktree, "auto-continue-task.md"), active)

		if got := resolveAutoContinueEndState(root, false, pipeline); got != active {
			t.Fatalf("resolved single-root body stayed on the stale pipeline copy\n--- got ---\n%s\n--- want active worktree ---\n%s", got, active)
		}
	})

	t.Run("split-root remains in state checkout", func(t *testing.T) {
		workflowRoot := t.TempDir()
		stateRoot := filepath.Join(workflowRoot, ".spacedock-state")
		worktree := ".worktrees/spacedock-ensign-auto-continue-task"
		stateBody := validated(strings.Replace(autoContinueEntity(), "status: implementation", "status: validation", 1))
		stateBody = strings.Replace(stateBody, "worktree:\n", "worktree: "+worktree+"\n", 1)
		writeFile(t, filepath.Join(stateRoot, worktree, "auto-continue-task.md"), "WRONG CODE WORKTREE COPY")

		if got := resolveAutoContinueEndState(stateRoot, true, stateBody); got != stateBody {
			t.Fatalf("split-root resolver left the state checkout: got %q", got)
		}
	})

	for _, archive := range []string{
		filepath.Join("_archive", "auto-continue-task.md"),
		filepath.Join("_archive", "auto-continue-task", "index.md"),
	} {
		t.Run("archive "+archive, func(t *testing.T) {
			stateRoot := t.TempDir()
			want := validated(strings.Replace(autoContinueEntity(), "status: implementation", "status: done", 1))
			writeFile(t, filepath.Join(stateRoot, archive), want)
			if got := resolveAutoContinueEndState(stateRoot, false, ""); got != want {
				t.Fatalf("archive fallback = %q, want %q", got, want)
			}
		})
	}
}

func assertAutoContinueGitBaseline(t *testing.T, root, stateRoot, entityPath string, split bool) {
	t.Helper()
	if err := autoContinueGitBaselineError(root, stateRoot, entityPath, split); err != nil {
		t.Fatal(err)
	}
}

func autoContinueGitBaselineError(root, stateRoot, entityPath string, split bool) error {
	wantStateRoot := root
	if split {
		wantStateRoot = filepath.Join(root, ".spacedock-state")
	}
	if stateRoot != wantStateRoot {
		return fmt.Errorf("state root = %q, want %q", stateRoot, wantStateRoot)
	}
	if _, err := os.Stat(entityPath); err != nil {
		return fmt.Errorf("entity path: %w", err)
	}
	for _, repo := range []string{root, stateRoot} {
		cmd := exec.Command("git", "-C", repo, "rev-parse", "--verify", "HEAD")
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("Git baseline missing at %s: %v: %s", repo, err, out)
		}
		cmd = exec.Command("git", "-C", repo, "rev-list", "--count", "HEAD")
		if out, err := cmd.CombinedOutput(); err != nil || strings.TrimSpace(string(out)) != "1" {
			return fmt.Errorf("Git baseline at %s has wrong commit count: err=%v count=%q", repo, err, out)
		}
		cmd = exec.Command("git", "-C", repo, "status", "--short")
		if out, err := cmd.CombinedOutput(); err != nil || strings.TrimSpace(string(out)) != "" {
			return fmt.Errorf("Git baseline at %s is not clean: err=%v status=%q", repo, err, out)
		}
		if repo == stateRoot {
			break
		}
	}
	return nil
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
