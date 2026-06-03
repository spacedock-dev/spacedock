//go:build live

package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

func emitClaudeScenarioMetrics(t *testing.T, scenario sharedRuntimeScenario, result claudeScenarioResult, model string) {
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
	record := journeymetrics.BuildRecord(journeymetrics.JourneySpec{
		ID:     "claude-" + scenario.name,
		Source: "live-harness",
		Host:   "claude",
		Model:  model,
	}, journeymetrics.BehaviorResult{Passed: true}, observation)
	if err := journeymetrics.EmitRecord(filepath.Join(dir, "shared-scenarios"), record); err != nil {
		t.Fatalf("emit Claude journey metrics for %s: %v", scenario.name, err)
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
		ID:     "codex-" + scenario.name,
		Source: "live-harness",
		Host:   "codex",
		Model:  characterization.Model,
	}, characterization, journeymetrics.BehaviorResult{Passed: true})
	record.DurationMS = result.duration.Milliseconds()
	record.ToolCalls = characterization.ToolCalls
	record.ToolCallsByName = characterization.ToolCallsByName
	if err := journeymetrics.EmitRecord(filepath.Join(dir, "shared-scenarios"), record); err != nil {
		t.Fatalf("emit Codex journey metrics for %s: %v", scenario.name, err)
	}
}
