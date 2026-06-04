package journeymetrics

import "time"

const (
	RecordSchemaVersion = 2
	LedgerSchemaVersion = 2
)

type MetricsState string

const (
	StateMeasured      MetricsState = "measured"
	StateCharacterized MetricsState = "characterized"
)

const (
	ModeLLMLive  = "llm-live"
	ModeCodified = "codified"
)

type JourneySpec struct {
	// ID is a legacy alias for ScenarioID.
	ID         string
	ScenarioID string
	Source     string
	Mode       string
	Runtime    string
	Executor   string
	Host       string
	Model      string
	Budget     Budget
}

type BehaviorResult struct {
	Passed  bool
	Failure string
}

type Observation struct {
	MetricsState    MetricsState
	Duration        time.Duration
	Turns           int
	ToolCalls       int
	ToolCallsByName map[string]int
	Tokens          TokenTotals
	TotalCostUSD    float64
	ModelUsage      map[string]ModelUsage
}

type TokenTotals struct {
	Input         int `json:"input,omitempty"`
	Output        int `json:"output,omitempty"`
	CacheCreation int `json:"cache_creation,omitempty"`
	CacheRead     int `json:"cache_read,omitempty"`
	Total         int `json:"total,omitempty"`
}

func (t TokenTotals) withTotal() TokenTotals {
	if t.Total == 0 {
		t.Total = t.Input + t.Output + t.CacheCreation + t.CacheRead
	}
	return t
}

func (t TokenTotals) isZero() bool {
	return t.Input == 0 && t.Output == 0 && t.CacheCreation == 0 && t.CacheRead == 0 && t.Total == 0
}

func (t TokenTotals) add(other TokenTotals) TokenTotals {
	return TokenTotals{
		Input:         t.Input + other.Input,
		Output:        t.Output + other.Output,
		CacheCreation: t.CacheCreation + other.CacheCreation,
		CacheRead:     t.CacheRead + other.CacheRead,
	}.withTotal()
}

type ModelUsage struct {
	Tokens  TokenTotals `json:"tokens"`
	CostUSD float64     `json:"cost_usd,omitempty"`
}

type Outcome struct {
	Status  string `json:"status"`
	Failure string `json:"failure,omitempty"`
}

type Budget struct {
	MaxTotalTokens *int `json:"max_total_tokens,omitempty"`
	MaxToolCalls   *int `json:"max_tool_calls,omitempty"`
}

type BudgetResult struct {
	Blocking   bool     `json:"blocking"`
	Violations []string `json:"violations,omitempty"`
}

type Record struct {
	SchemaVersion   int                    `json:"schema_version"`
	ScenarioID      string                 `json:"scenario_id"`
	JourneyID       string                 `json:"journey_id,omitempty"`
	Source          string                 `json:"source"`
	Mode            string                 `json:"mode,omitempty"`
	Runtime         string                 `json:"runtime,omitempty"`
	Executor        string                 `json:"executor,omitempty"`
	Host            string                 `json:"host"`
	Model           string                 `json:"model"`
	MetricsState    MetricsState           `json:"metrics_state"`
	Outcome         Outcome                `json:"outcome"`
	DurationMS      int64                  `json:"duration_ms,omitempty"`
	Turns           int                    `json:"turns,omitempty"`
	ToolCalls       int                    `json:"tool_calls,omitempty"`
	ToolCallsByName map[string]int         `json:"tool_calls_by_name,omitempty"`
	Tokens          TokenTotals            `json:"tokens,omitempty"`
	TotalCostUSD    float64                `json:"total_cost_usd,omitempty"`
	ModelUsage      map[string]ModelUsage  `json:"model_usage,omitempty"`
	Budget          Budget                 `json:"budget,omitempty"`
	BudgetResult    *BudgetResult          `json:"budget_result,omitempty"`
	CodexCharacter  *CodexCharacterization `json:"codex_characterization,omitempty"`
}

type recordJSON struct {
	SchemaVersion   int                    `json:"schema_version"`
	ScenarioID      string                 `json:"scenario_id"`
	Source          string                 `json:"source"`
	Mode            string                 `json:"mode,omitempty"`
	Runtime         string                 `json:"runtime,omitempty"`
	Executor        string                 `json:"executor,omitempty"`
	Host            string                 `json:"host"`
	Model           string                 `json:"model"`
	MetricsState    MetricsState           `json:"metrics_state"`
	Outcome         Outcome                `json:"outcome"`
	DurationMS      int64                  `json:"duration_ms,omitempty"`
	Turns           int                    `json:"turns,omitempty"`
	ToolCalls       int                    `json:"tool_calls,omitempty"`
	ToolCallsByName map[string]int         `json:"tool_calls_by_name,omitempty"`
	Tokens          *TokenTotals           `json:"tokens,omitempty"`
	TotalCostUSD    float64                `json:"total_cost_usd,omitempty"`
	ModelUsage      map[string]ModelUsage  `json:"model_usage,omitempty"`
	Budget          *Budget                `json:"budget,omitempty"`
	BudgetResult    *BudgetResult          `json:"budget_result,omitempty"`
	CodexCharacter  *CodexCharacterization `json:"codex_characterization,omitempty"`
}

func (r Record) MarshalJSON() ([]byte, error) {
	var tokens *TokenTotals
	if !r.Tokens.isZero() {
		t := r.Tokens.withTotal()
		tokens = &t
	}
	var budget *Budget
	if r.Budget.MaxTotalTokens != nil || r.Budget.MaxToolCalls != nil {
		b := r.Budget
		budget = &b
	}
	return marshalRecordJSON(recordJSON{
		SchemaVersion:   RecordSchemaVersion,
		ScenarioID:      firstNonEmpty(r.ScenarioID, r.JourneyID),
		Source:          r.Source,
		Mode:            r.Mode,
		Runtime:         firstNonEmpty(r.Runtime, r.Host),
		Executor:        r.Executor,
		Host:            r.Host,
		Model:           r.Model,
		MetricsState:    r.MetricsState,
		Outcome:         r.Outcome,
		DurationMS:      r.DurationMS,
		Turns:           r.Turns,
		ToolCalls:       r.ToolCalls,
		ToolCallsByName: r.ToolCallsByName,
		Tokens:          tokens,
		TotalCostUSD:    r.TotalCostUSD,
		ModelUsage:      r.ModelUsage,
		Budget:          budget,
		BudgetResult:    r.BudgetResult,
		CodexCharacter:  r.CodexCharacter,
	})
}
