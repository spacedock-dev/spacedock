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
	want := map[string]string{
		"default-headless-gate-stop":         defaultHeadlessGateStopDefectID,
		"auto-continue-after-implementation": liveDurableJourneyDefectID,
		"smallest-sufficient-mechanism":      liveDurableJourneyDefectID,
		"keep-moving-posture":                liveDurableJourneyDefectID,
	}
	for _, scenario := range sharedRuntimeScenarios() {
		reason := liveDurableJourneyTODO(scenario.name)
		owner, missing := want[scenario.name]
		if missing {
			if reason == "" || !strings.HasPrefix(reason, "TODO("+owner+"):") {
				t.Errorf("TODO journey %q reason = %q, want exact owner TODO(%s)", scenario.name, reason, owner)
			}
			delete(want, scenario.name)
			continue
		}
		if reason != "" {
			t.Errorf("implemented journey %q unexpectedly has missing-evidence TODO %q", scenario.name, reason)
		}
	}
	if len(want) != 0 {
		t.Fatalf("TODO journeys are hidden from the shared scenario table: %v", want)
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
