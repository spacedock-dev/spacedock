//go:build live

package ensigncycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

func FuzzPiSharedLiveDriverEmitsJourneyMetric(f *testing.F) {
	f.Add("pi-metrics-proof")
	f.Fuzz(func(t *testing.T, scenarioName string) {
		dir := t.TempDir()
		t.Setenv("SPACEDOCK_JOURNEY_METRICS_DIR", dir)
		driver := piSharedLiveDriver{modelName: "openai/gpt-5.6-luna:max"}
		driver.emitMetrics(t, sharedRuntimeScenario{name: scenarioName}, liveResult{duration: 2 * time.Second})

		paths, err := filepath.Glob(filepath.Join(dir, "shared-scenarios", "*.json"))
		if err != nil {
			t.Fatal(err)
		}
		if len(paths) != 1 {
			t.Fatalf("Pi journey metric files = %d, want 1", len(paths))
		}
		data, err := os.ReadFile(paths[0])
		if err != nil {
			t.Fatal(err)
		}
		var record journeymetrics.Record
		if len(data) == 0 || json.Unmarshal(data, &record) != nil {
			t.Fatalf("Pi journey metric is empty or invalid: %q", data)
		}
		if record.ScenarioID != scenarioName || record.Runtime != "pi" || record.Model != driver.modelName || record.DurationMS != 2000 {
			t.Fatalf("Pi journey metric = %#v", record)
		}
	})
}

func emitClaudeScenarioMetrics(t *testing.T, scenario sharedRuntimeScenario, result liveResult, model string) {
	t.Helper()
	dir := os.Getenv("SPACEDOCK_JOURNEY_METRICS_DIR")
	if dir == "" {
		return
	}
	parsed, err := journeymetrics.ParseClaudeJSONL([]byte(result.stream))
	if err != nil {
		t.Fatalf("parse Claude journey metrics for %s: %v", scenario.name, err)
	}
	observation := parsed.Observation
	observation.Duration = result.duration
	// Fold the dispatched-ensign sub-agent transcripts' --read adoption onto the FO
	// front-door counts: `status --read` adoption is principally an ensign behavior,
	// and the ensign runs as a separate sub-agent session whose transcript lands on
	// disk, never in result.stream.
	observation, err = journeymetrics.FoldEnsignReadAdoption(observation, ensignTranscripts(result))
	if err != nil {
		t.Fatalf("fold ensign --read adoption for %s: %v", scenario.name, err)
	}
	record := journeymetrics.BuildRecord(journeymetrics.JourneySpec{
		ScenarioID: scenario.name,
		Source:     "live-harness",
		Mode:       journeymetrics.ModeLLMLive,
		Runtime:    "claude",
		Executor:   "llm",
		Host:       "claude",
		Model:      model,
	}, journeymetrics.BehaviorResult{Passed: true}, observation)
	if err := journeymetrics.EmitRecord(filepath.Join(dir, "shared-scenarios"), record); err != nil {
		t.Fatalf("emit Claude journey metrics for %s: %v", scenario.name, err)
	}
}

// ensignTranscripts reads the dispatched-ensign sub-agent transcripts for this
// journey from disk. The ensign runs as a separate sub-agent session whose
// transcript lands under the FO session's subagents dir — never in result.stream —
// so the glob mirrors the proven scanSubagentMeta shape:
// {configDir}/projects/{encode(cwd)}/{FO-session-id}/subagents/agent-*.jsonl. The
// FO session id comes from the stream's system/init event. Returns nil when the run
// did not record a config dir / cwd / session id (e.g. the pty transport), so the
// fold no-ops to FO-front-door counts.
func ensignTranscripts(result liveResult) [][]byte {
	if result.configDir == "" || result.cwd == "" {
		return nil
	}
	sessionID := initEventSessionID(strings.Split(result.stream, "\n"))
	if sessionID == "" {
		return nil
	}
	pattern := filepath.Join(result.configDir, "projects",
		encodeProjectDir(result.cwd), sessionID, "subagents", "agent-*.jsonl")
	matches, _ := filepath.Glob(pattern)
	var transcripts [][]byte
	for _, p := range matches {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		transcripts = append(transcripts, data)
	}
	return transcripts
}

// emitShallowBootWindowMetrics emits the shallow-boot-window observation (AC-1)
// alongside the whole-run "shallow-boot" record emitClaudeScenarioMetrics already
// publishes, into the same SPACEDOCK_JOURNEY_METRICS_DIR/shared-scenarios dir, so
// the two sibling records land together without either overwriting the other.
func emitShallowBootWindowMetrics(t *testing.T, stream string, model string) {
	t.Helper()
	dir := os.Getenv("SPACEDOCK_JOURNEY_METRICS_DIR")
	if dir == "" {
		return
	}
	turns, err := journeymetrics.ParseClaudeTurns([]byte(stream))
	if err != nil {
		t.Fatalf("parse Claude turns for shallow-boot-window: %v", err)
	}
	record, err := BuildShallowBootWindowRecord(turns, model, journeymetrics.ParseClaudeCodeVersion([]byte(stream)), journeymetrics.ParseClaudeInitModel([]byte(stream)))
	if err != nil {
		t.Fatalf("build shallow-boot-window record: %v", err)
	}
	if err := journeymetrics.EmitRecord(filepath.Join(dir, "shared-scenarios"), record); err != nil {
		t.Fatalf("emit shallow-boot-window record: %v", err)
	}
}

func emitCodexScenarioMetrics(t *testing.T, scenario sharedRuntimeScenario, result codexScenarioResult) {
	t.Helper()
	dir := os.Getenv("SPACEDOCK_JOURNEY_METRICS_DIR")
	if dir == "" {
		return
	}
	characterization, err := journeymetrics.CharacterizeCodexExecJSONL([]byte(result.jsonl))
	if err != nil {
		t.Fatalf("characterize Codex journey metrics for %s: %v", scenario.name, err)
	}
	record := journeymetrics.CodexCharacterizedRecord(journeymetrics.JourneySpec{
		ScenarioID: scenario.name,
		Source:     "live-harness",
		Mode:       journeymetrics.ModeLLMLive,
		Runtime:    "codex",
		Executor:   "llm",
		Host:       "codex",
		Model:      characterization.Model,
	}, characterization, journeymetrics.BehaviorResult{Passed: true})
	record.DurationMS = result.duration.Milliseconds()
	record.ToolCalls = characterization.ToolCalls
	record.ToolCallsByName = characterization.ToolCallsByName
	if err := journeymetrics.EmitRecord(filepath.Join(dir, "shared-scenarios"), record); err != nil {
		t.Fatalf("emit Codex journey metrics for %s: %v", scenario.name, err)
	}
}

func emitPiScenarioMetrics(t *testing.T, scenario sharedRuntimeScenario, result liveResult, model string) {
	t.Helper()
	dir := os.Getenv("SPACEDOCK_JOURNEY_METRICS_DIR")
	if dir == "" {
		return
	}
	record := journeymetrics.BuildRecord(journeymetrics.JourneySpec{
		ScenarioID: scenario.name,
		Source:     "live-harness",
		Mode:       journeymetrics.ModeLLMLive,
		Runtime:    "pi",
		Executor:   "llm",
		Host:       "pi",
		Model:      model,
	}, journeymetrics.BehaviorResult{Passed: true}, journeymetrics.Observation{
		Duration: result.duration,
	})
	if err := journeymetrics.EmitRecord(filepath.Join(dir, "shared-scenarios"), record); err != nil {
		t.Fatalf("emit Pi journey metrics for %s: %v", scenario.name, err)
	}
}
