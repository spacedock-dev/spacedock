//go:build live

package ensigncycle

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/journeymetrics"
)

type piMetricsParseResult struct {
	provider    string
	model       string
	observation journeymetrics.Observation
}

func parsePiSessionMetrics(data []byte) (piMetricsParseResult, error) {
	result := piMetricsParseResult{}
	seenMessages := map[string]bool{}
	seenTools := map[string]bool{}
	result.observation.MetricsState = journeymetrics.StateMeasured
	result.observation.ToolCallsByName = map[string]int{}
	result.observation.ModelUsage = map[string]journeymetrics.ModelUsage{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var row struct {
			Type    string `json:"type"`
			ID      string `json:"id"`
			Message struct {
				Role, Provider, Model string
				Usage                 struct {
					Input, Output, CacheRead, CacheWrite, TotalTokens *int
					Cost                                              struct{ Total *float64 }
				}
				Content []struct{ Type, ID, Name string }
			}
		}
		if json.Unmarshal(scanner.Bytes(), &row) != nil || row.Type != "message" || row.Message.Role != "assistant" || row.ID == "" || seenMessages[row.ID] {
			continue
		}
		seenMessages[row.ID] = true
		u := row.Message.Usage
		if u.Input == nil || u.Output == nil || u.CacheRead == nil || u.CacheWrite == nil || u.TotalTokens == nil || u.Cost.Total == nil {
			return piMetricsParseResult{}, fmt.Errorf("Pi assistant usage is missing token or total-cost fields")
		}
		if row.Message.Provider == "" || row.Message.Model == "" {
			return piMetricsParseResult{}, fmt.Errorf("Pi assistant usage is missing provider/model")
		}
		if result.provider != "" && (result.provider != row.Message.Provider || result.model != row.Message.Model) {
			return piMetricsParseResult{}, fmt.Errorf("Pi session mixes provider/models: %s/%s and %s/%s", result.provider, result.model, row.Message.Provider, row.Message.Model)
		}
		result.provider, result.model = row.Message.Provider, row.Message.Model
		tokens := journeymetrics.TokenTotals{Input: *u.Input, Output: *u.Output, CacheCreation: *u.CacheWrite, CacheRead: *u.CacheRead, Total: *u.TotalTokens}
		result.observation.Turns++
		addPiTokens(&result.observation.Tokens, tokens)
		result.observation.TotalCostUSD += *u.Cost.Total
		key := result.provider + "/" + result.model
		modelUsage := result.observation.ModelUsage[key]
		addPiTokens(&modelUsage.Tokens, tokens)
		modelUsage.CostUSD += *u.Cost.Total
		result.observation.ModelUsage[key] = modelUsage
		for _, block := range row.Message.Content {
			if block.Type == "toolCall" && block.ID != "" && !seenTools[block.ID] {
				seenTools[block.ID] = true
				result.observation.ToolCalls++
				result.observation.ToolCallsByName[block.Name]++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return piMetricsParseResult{}, err
	}
	if result.observation.Turns == 0 {
		return piMetricsParseResult{}, fmt.Errorf("Pi session contains no assistant usage")
	}
	return result, nil
}

func addPiTokens(dst *journeymetrics.TokenTotals, src journeymetrics.TokenTotals) {
	dst.Input += src.Input
	dst.Output += src.Output
	dst.CacheCreation += src.CacheCreation
	dst.CacheRead += src.CacheRead
	dst.Total += src.Total
}

func TestEmitPiScenarioMetricsUsesNativeSessionForAttributionAndUsage(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("SPACEDOCK_JOURNEY_METRICS_DIR", dir)
	t.Setenv("GITHUB_RUN_ID", "31016570689")
	t.Setenv("GITHUB_SERVER_URL", "https://github.com")
	t.Setenv("GITHUB_REPOSITORY", "spacedock-dev/spacedock")
	session := `{"type":"message","id":"turn-1","message":{"role":"assistant","provider":"openai","model":"gpt-5.4","usage":{"input":3,"output":2,"cacheRead":10,"cacheWrite":1,"reasoning":4,"totalTokens":20,"cost":{"total":0.25}},"content":[{"type":"toolCall","id":"call-1","name":"bash","arguments":{"command":"spacedock status"}}]}}`
	emitPiScenarioMetrics(t, sharedRuntimeScenario{name: "zero-discovery"}, []liveResult{{
		runtime: "pi", sessionJSONL: session, artifactDir: filepath.Join(dir, "pi-shared-scenarios", "zero-discovery"), duration: 1500 * time.Millisecond,
	}})
	matches, err := filepath.Glob(filepath.Join(dir, "shared-scenarios", "*.json"))
	if err != nil || len(matches) != 1 {
		t.Fatalf("Pi metric files = %v (error %v), want exactly one", matches, err)
	}
	if strings.Contains(filepath.Base(matches[0]), "claude") || !strings.Contains(filepath.Base(matches[0]), "--pi--") {
		t.Fatalf("Pi metric filename has wrong runtime attribution: %s", filepath.Base(matches[0]))
	}
	data, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	var record journeymetrics.Record
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	if record.Runtime != "pi" || record.Host != "pi" || record.Model != "openai/gpt-5.4" || record.DurationMS != 1500 {
		t.Fatalf("Pi attribution/duration = runtime=%q host=%q model=%q duration=%d", record.Runtime, record.Host, record.Model, record.DurationMS)
	}
	if record.RunID != "31016570689" || record.RunURL != "https://github.com/spacedock-dev/spacedock/actions/runs/31016570689" {
		t.Fatalf("Pi artifact provenance = run %q URL %q", record.RunID, record.RunURL)
	}
	if record.Tokens != (journeymetrics.TokenTotals{Input: 3, Output: 2, CacheCreation: 1, CacheRead: 10, Total: 20}) || record.TotalCostUSD != 0.25 {
		t.Fatalf("Pi usage = tokens=%+v cost=%v", record.Tokens, record.TotalCostUSD)
	}
}

func TestParsePiSessionMetricsRequiresCompleteAttributedLargeRow(t *testing.T) {
	valid := `{"type":"message","id":"turn-large","message":{"role":"assistant","provider":"openai","model":"gpt-5.4","usage":{"input":1,"output":1,"cacheRead":0,"cacheWrite":0,"totalTokens":2,"cost":{"total":0}},"content":[{"type":"text","text":"` + strings.Repeat("x", 106000) + `"}]}}`
	if _, err := parsePiSessionMetrics([]byte(valid)); err != nil {
		t.Fatalf("representative >64KiB row with explicit zero fields rejected: %v", err)
	}
	for name, data := range map[string]string{
		"missing cost":        strings.Replace(valid, `,"cost":{"total":0}`, "", 1),
		"missing totalTokens": strings.Replace(valid, `,"totalTokens":2`, "", 1),
		"missing cacheWrite":  strings.Replace(valid, `,"cacheWrite":0`, "", 1),
		"missing provider":    strings.Replace(valid, `"provider":"openai",`, "", 1),
		"mixed model":         valid + "\n" + strings.Replace(strings.Replace(valid, "turn-large", "turn-two", 1), "gpt-5.4", "gpt-5.5", 1),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := parsePiSessionMetrics([]byte(data)); err == nil {
				t.Fatal("partial or mixed Pi usage was accepted")
			}
		})
	}
}

func emitClaudeScenarioMetrics(t *testing.T, scenario sharedRuntimeScenario, result liveResult, model string) {
	t.Helper()
	if result.runtime == "pi" {
		return
	}
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

func emitPiScenarioMetrics(t *testing.T, scenario sharedRuntimeScenario, results []liveResult) {
	t.Helper()
	dir := os.Getenv("SPACEDOCK_JOURNEY_METRICS_DIR")
	if dir == "" {
		return
	}
	if len(results) == 0 {
		t.Fatalf("emit Pi journey metrics for %s: no archived root session", scenario.name)
	}
	observation := journeymetrics.Observation{
		MetricsState:    journeymetrics.StateMeasured,
		ToolCallsByName: map[string]int{},
		ModelUsage:      map[string]journeymetrics.ModelUsage{},
	}
	provider, model := "", ""
	for _, result := range results {
		parsed, err := parsePiSessionMetrics([]byte(result.sessionJSONL))
		if err != nil {
			t.Fatalf("parse Pi journey metrics for %s (%s): %v", scenario.name, result.artifactDir, err)
		}
		if provider != "" && (provider != parsed.provider || model != parsed.model) {
			t.Fatalf("Pi journey metrics for %s mixed provider/models: %s/%s and %s/%s", scenario.name, provider, model, parsed.provider, parsed.model)
		}
		provider, model = parsed.provider, parsed.model
		observation.Duration += result.duration
		observation.Turns += parsed.observation.Turns
		observation.ToolCalls += parsed.observation.ToolCalls
		addPiTokens(&observation.Tokens, parsed.observation.Tokens)
		observation.TotalCostUSD += parsed.observation.TotalCostUSD
		for name, count := range parsed.observation.ToolCallsByName {
			observation.ToolCallsByName[name] += count
		}
		for name, usage := range parsed.observation.ModelUsage {
			current := observation.ModelUsage[name]
			addPiTokens(&current.Tokens, usage.Tokens)
			current.CostUSD += usage.CostUSD
			observation.ModelUsage[name] = current
		}
	}
	resolvedModel := provider + "/" + model
	record := journeymetrics.BuildRecord(journeymetrics.JourneySpec{
		ScenarioID: scenario.name,
		Source:     "live-harness",
		Mode:       journeymetrics.ModeLLMLive,
		Runtime:    "pi",
		Executor:   "llm",
		Host:       "pi",
		Model:      resolvedModel,
	}, journeymetrics.BehaviorResult{Passed: true}, observation)
	if err := journeymetrics.EmitRecord(filepath.Join(dir, "shared-scenarios"), record); err != nil {
		t.Fatalf("emit Pi journey metrics for %s: %v", scenario.name, err)
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
