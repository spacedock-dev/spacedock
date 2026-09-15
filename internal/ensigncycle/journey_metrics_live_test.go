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
	if len(result.phases) > 0 {
		// Each launch has independent IDs and terminal usage totals.
		observation = journeymetrics.Observation{MetricsState: observation.MetricsState, ClaudeCodeVersion: observation.ClaudeCodeVersion, ResolvedModel: observation.ResolvedModel, ToolCallsByName: map[string]int{}, ModelUsage: map[string]journeymetrics.ModelUsage{}}
		for _, phase := range result.phases {
			parsed, err := journeymetrics.ParseClaudeJSONL([]byte(phase.stream))
			if err != nil {
				t.Fatal(err)
			}
			o := parsed.Observation
			observation.Turns += o.Turns
			observation.ToolCalls += o.ToolCalls
			observation.StatusReadCalls += o.StatusReadCalls
			observation.ScopedReadCalls += o.ScopedReadCalls
			observation.Tokens = sumPhaseTokens(observation.Tokens, o.Tokens)
			observation.TotalCostUSD += o.TotalCostUSD
			for name, count := range o.ToolCallsByName {
				observation.ToolCallsByName[name] += count
			}
			for name, usage := range o.ModelUsage {
				prior := observation.ModelUsage[name]
				prior.Tokens = sumPhaseTokens(prior.Tokens, usage.Tokens)
				prior.CostUSD += usage.CostUSD
				observation.ModelUsage[name] = prior
			}
		}
	}
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
	}, scenarioBehaviorResult(scenario), observation)
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
	if len(result.phases) > 0 {
		var transcripts [][]byte
		for _, phase := range result.phases {
			transcripts = append(transcripts, ensignTranscripts(phase)...)
		}
		return transcripts
	}
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

func emitCodexScenarioMetrics(t *testing.T, scenario sharedRuntimeScenario, result codexScenarioResult, phases ...liveResult) {
	t.Helper()
	dir := os.Getenv("SPACEDOCK_JOURNEY_METRICS_DIR")
	if dir == "" {
		return
	}
	characterization, err := journeymetrics.CharacterizeCodexExecJSONL([]byte(result.jsonl))
	if err != nil {
		t.Fatalf("characterize Codex journey metrics for %s: %v", scenario.name, err)
	}
	if len(phases) > 0 {
		characterization.ToolCalls = 0
		characterization.StatusReadCalls = 0
		characterization.ToolCallsByName = map[string]int{}
		for _, phase := range phases {
			part, err := journeymetrics.CharacterizeCodexExecJSONL([]byte(phase.stream))
			if err != nil {
				t.Fatal(err)
			}
			characterization.ToolCalls += part.ToolCalls
			characterization.StatusReadCalls += part.StatusReadCalls
			for name, count := range part.ToolCallsByName {
				characterization.ToolCallsByName[name] += count
			}
		}
	}
	record := journeymetrics.CodexCharacterizedRecord(journeymetrics.JourneySpec{
		ScenarioID: scenario.name,
		Source:     "live-harness",
		Mode:       journeymetrics.ModeLLMLive,
		Runtime:    "codex",
		Executor:   "llm",
		Host:       "codex",
		Model:      characterization.Model,
	}, characterization, scenarioBehaviorResult(scenario))
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
	}, scenarioBehaviorResult(scenario), journeymetrics.Observation{
		Duration: result.duration,
	})
	if err := journeymetrics.EmitRecord(filepath.Join(dir, "shared-scenarios"), record); err != nil {
		t.Fatalf("emit Pi journey metrics for %s: %v", scenario.name, err)
	}
}

func scenarioBehaviorResult(scenario sharedRuntimeScenario) journeymetrics.BehaviorResult {
	result := journeymetrics.BehaviorResult{Passed: true}
	if scenario.gap.kind == "xfail" {
		result.Outcome = &journeymetrics.Outcome{Status: scenario.grade.status, Owner: scenario.gap.owner, FailureCodes: scenario.grade.codes}
	}
	return result
}

func TestSmallestMechanismPhaseMetrics(t *testing.T) {
	for _, host := range []string{"claude", "codex", "pi"} {
		t.Run(host, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("SPACEDOCK_JOURNEY_METRICS_DIR", dir)
			stream := `{"type":"assistant","message":{"id":"same-id","content":[{"type":"tool_use","id":"same-tool","name":"Read"}]}}` + "\n" + `{"type":"result","usage":{"input_tokens":10,"output_tokens":2},"total_cost_usd":0.5}`
			if host == "codex" {
				stream = `{"type":"tool_call.started","call_id":"same-tool","name":"exec_command"}`
			}
			phase := liveResult{stream: stream, duration: time.Second}
			result := liveResult{stream: stream + "\n" + stream, duration: 2 * time.Second, phases: []liveResult{phase, phase}}
			scenario := sharedRuntimeScenario{name: "phase-metrics"}
			switch host {
			case "claude":
				emitClaudeScenarioMetrics(t, scenario, result, "test")
			case "codex":
				(codexAsLiveDriver{}).emitMetrics(t, scenario, result)
			case "pi":
				emitPiScenarioMetrics(t, scenario, result, "test")
			}
			paths, _ := filepath.Glob(filepath.Join(dir, "shared-scenarios", "*.json"))
			if len(paths) != 1 {
				t.Fatalf("metric paths: %v", paths)
			}
			data, err := os.ReadFile(paths[0])
			if err != nil {
				t.Fatal(err)
			}
			var record journeymetrics.Record
			if err := json.Unmarshal(data, &record); err != nil {
				t.Fatal(err)
			}
			if record.DurationMS != 2000 {
				t.Fatalf("duration: %d", record.DurationMS)
			}
			if host != "pi" && record.ToolCalls != 2 {
				t.Errorf("tool calls = %d, want 2", record.ToolCalls)
			}
			if host == "claude" && (record.Tokens.Total != 24 || record.TotalCostUSD != 1 || record.Turns != 2) {
				t.Errorf("combined usage: %+v", record)
			}
		})
	}
}

func sumPhaseTokens(a, b journeymetrics.TokenTotals) journeymetrics.TokenTotals {
	return journeymetrics.TokenTotals{Input: a.Input + b.Input, Output: a.Output + b.Output, CacheCreation: a.CacheCreation + b.CacheCreation, CacheRead: a.CacheRead + b.CacheRead, Total: a.Total + b.Total}
}
