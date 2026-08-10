package journeymetrics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func marshalRecordJSON(r recordJSON) ([]byte, error) {
	return json.Marshal(r)
}

func Track(dir string, spec JourneySpec, behavior func() BehaviorResult, observe func() Observation) (BehaviorResult, error) {
	start := time.Now()
	result := behavior()
	observation := observe()
	if observation.Duration == 0 {
		observation.Duration = time.Since(start)
	}
	record := BuildRecord(spec, result, observation)
	if dir != "" {
		if err := EmitRecord(dir, record); err != nil {
			return result, err
		}
	}
	return result, nil
}

func BuildRecord(spec JourneySpec, result BehaviorResult, observation Observation) Record {
	state := observation.MetricsState
	if state == "" {
		state = StateMeasured
	}
	outcome := Outcome{Status: "failed", Failure: result.Failure}
	if result.Passed {
		outcome = Outcome{Status: "passed"}
	}
	if result.Outcome != nil {
		outcome = *result.Outcome
	}
	record := Record{
		SchemaVersion:             RecordSchemaVersion,
		ScenarioID:                firstNonEmpty(spec.ScenarioID, spec.ID),
		Source:                    spec.Source,
		Mode:                      spec.Mode,
		Runtime:                   firstNonEmpty(spec.Runtime, spec.Host),
		Executor:                  spec.Executor,
		Host:                      spec.Host,
		Model:                     spec.Model,
		MetricsState:              state,
		Outcome:                   outcome,
		DurationMS:                observation.Duration.Milliseconds(),
		Turns:                     observation.Turns,
		ToolCalls:                 observation.ToolCalls,
		ToolCallsByName:           observation.ToolCallsByName,
		StatusReadCalls:           observation.StatusReadCalls,
		ScopedReadCalls:           observation.ScopedReadCalls,
		Tokens:                    observation.Tokens.withTotal(),
		BaselineTokens:            observation.BaselineTokens.withTotal(),
		PreGreetPeakCacheCreation: observation.PreGreetPeakCacheCreation,
		TotalCostUSD:              observation.TotalCostUSD,
		ModelUsage:                normalizeModelUsage(observation.ModelUsage),
		Budget:                    spec.Budget,
		ClaudeCodeVersion:         observation.ClaudeCodeVersion,
		ResolvedModel:             observation.ResolvedModel,
	}
	if hasBudget(spec.Budget) {
		result := EvaluateBudget(record, spec.Budget)
		record.BudgetResult = &result
	}
	return record
}

func EmitRecord(dir string, record Record) error {
	record = normalizeRecord(record)
	record = stampProvenance(record)
	if strings.TrimSpace(record.ScenarioID) == "" {
		return fmt.Errorf("scenario id is required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, recordFilename(record)), data, 0o644)
}

// stampProvenance fills in the run-provenance fields a caller left unset.
// CapturedAt defaults to the current time, so every freshly emitted record can be
// ordered chronologically once it lands in a published ledger. RunID/RunURL
// default from the GitHub Actions run environment ($GITHUB_RUN_ID /
// $GITHUB_SERVER_URL / $GITHUB_REPOSITORY) when present; outside CI they stay
// empty and the omitempty tags drop them from the JSON. It deliberately runs ONLY
// at emission time, never inside normalizeRecord — normalizeRecord also runs when
// re-reading an already-emitted record (ReadRecordsDir, AggregateLedger), and
// re-stamping CapturedAt there would silently overwrite a historical record's real
// capture time with "now" on every rebuild.
func stampProvenance(record Record) Record {
	if record.CapturedAt == "" {
		record.CapturedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if record.RunID == "" {
		record.RunID = os.Getenv("GITHUB_RUN_ID")
	}
	if record.RunURL == "" && record.RunID != "" {
		serverURL, repo := os.Getenv("GITHUB_SERVER_URL"), os.Getenv("GITHUB_REPOSITORY")
		if serverURL != "" && repo != "" {
			record.RunURL = strings.TrimSuffix(serverURL, "/") + "/" + repo + "/actions/runs/" + record.RunID
		}
	}
	return record
}

func EvaluateBudget(record Record, budget Budget) BudgetResult {
	var violations []string
	if budget.MaxTotalTokens != nil {
		if record.MetricsState != StateMeasured {
			violations = append(violations, "max_total_tokens requires measured metrics")
		} else if record.Tokens.Total > *budget.MaxTotalTokens {
			violations = append(violations,
				fmt.Sprintf("max_total_tokens exceeded: got %d, max %d", record.Tokens.Total, *budget.MaxTotalTokens))
		}
	}
	if budget.MaxToolCalls != nil && record.ToolCalls > *budget.MaxToolCalls {
		violations = append(violations,
			fmt.Sprintf("max_tool_calls exceeded: got %d, max %d", record.ToolCalls, *budget.MaxToolCalls))
	}
	return BudgetResult{Blocking: len(violations) > 0, Violations: violations}
}

func hasBudget(b Budget) bool {
	return b.MaxTotalTokens != nil || b.MaxToolCalls != nil
}

func normalizeRecord(record Record) Record {
	record.SchemaVersion = RecordSchemaVersion
	record.ScenarioID = firstNonEmpty(record.ScenarioID, record.JourneyID)
	record.JourneyID = ""
	record.Runtime = firstNonEmpty(record.Runtime, record.Host)
	record.Tokens = record.Tokens.withTotal()
	record.BaselineTokens = record.BaselineTokens.withTotal()
	record.ModelUsage = normalizeModelUsage(record.ModelUsage)
	return record
}

func recordFilename(record Record) string {
	parts := []string{record.ScenarioID}
	seen := map[string]bool{record.ScenarioID: true}
	for _, part := range []string{
		record.Runtime,
		record.Executor,
		record.Host,
		record.Mode,
		record.Model,
		string(record.MetricsState),
	} {
		if part == "" || seen[part] {
			continue
		}
		seen[part] = true
		parts = append(parts, part)
	}
	return safeFilename(strings.Join(parts, "--")) + ".json"
}

func safeFilename(id string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '-' || r == '_' || r == '.':
			return r
		default:
			return '-'
		}
	}, id)
}

func normalizeModelUsage(in map[string]ModelUsage) map[string]ModelUsage {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]ModelUsage, len(in))
	for model, usage := range in {
		usage.Tokens = usage.Tokens.withTotal()
		out[model] = usage
	}
	return out
}
