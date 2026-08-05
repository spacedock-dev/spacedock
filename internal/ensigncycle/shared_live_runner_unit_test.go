//go:build live

package ensigncycle

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSharedScenarioRunnerCoverageFinal(t *testing.T) {
	want := []string{
		"full-ensign-cycle",
		"gate-guardrail",
		"default-headless-gate-stop",
		"withdrawn-gate-recovery",
		"recorded-gate-lifecycle",
		"rejection-flow",
		"feedback-3-cycle-escalation",
		"merge-hook-guardrail",
		"filing",
		"shallow-boot",
		"zero-discovery",
		"auto-continue-after-implementation",
		"self-evidence-merge-triage",
		"smallest-sufficient-mechanism",
		"keep-moving-posture",
		"ac-value-reanchor",
	}

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
	targets := []liveEvidenceTarget{
		liveEvidenceTargetClaudeSonnet,
		liveEvidenceTargetClaudeOpus,
		liveEvidenceTargetCodex,
		liveEvidenceTargetPi,
	}
	want := map[liveEvidenceKey]string{
		{target: liveEvidenceTargetClaudeSonnet, journey: "default-headless-gate-stop"}:    defaultHeadlessGateStopDefectID,
		{target: liveEvidenceTargetClaudeSonnet, journey: "smallest-sufficient-mechanism"}: liveDurableJourneyDefectID,
		{target: liveEvidenceTargetClaudeSonnet, journey: "keep-moving-posture"}:           liveDurableJourneyDefectID,
		{target: liveEvidenceTargetCodex, journey: "smallest-sufficient-mechanism"}:        liveDurableJourneyDefectID,
		{target: liveEvidenceTargetCodex, journey: "keep-moving-posture"}:                  liveDurableJourneyDefectID,
	}
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

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
