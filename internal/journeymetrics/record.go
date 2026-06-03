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
	record := Record{
		SchemaVersion:   RecordSchemaVersion,
		JourneyID:       spec.ID,
		Source:          spec.Source,
		Host:            spec.Host,
		Model:           spec.Model,
		MetricsState:    state,
		Outcome:         outcome,
		DurationMS:      observation.Duration.Milliseconds(),
		Turns:           observation.Turns,
		ToolCalls:       observation.ToolCalls,
		ToolCallsByName: observation.ToolCallsByName,
		Tokens:          observation.Tokens.withTotal(),
		TotalCostUSD:    observation.TotalCostUSD,
		ModelUsage:      normalizeModelUsage(observation.ModelUsage),
		Budget:          spec.Budget,
	}
	if hasBudget(spec.Budget) {
		result := EvaluateBudget(record, spec.Budget)
		record.BudgetResult = &result
	}
	return record
}

func EmitRecord(dir string, record Record) error {
	if strings.TrimSpace(record.JourneyID) == "" {
		return fmt.Errorf("journey id is required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, safeFilename(record.JourneyID)+".json"), data, 0o644)
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
