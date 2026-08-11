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
	Outcome *Outcome
}

type Observation struct {
	MetricsState    MetricsState
	Duration        time.Duration
	Turns           int
	ToolCalls       int
	ToolCallsByName map[string]int
	StatusReadCalls int
	ScopedReadCalls int
	Tokens          TokenTotals
	// BaselineTokens is the first turn's full TokenTotals, a reference point for
	// how the observation started. Per-turn usage is NOT cumulative across
	// turns, so only Tokens.Context() - BaselineTokens.Context() is a
	// defensible derived quantity (Context tracks the conversation's
	// accumulating cached prefix); Input/Output are not independently
	// subtractable between two unrelated single requests. Left zero-valued for
	// scenarios that don't measure a pre-observation baseline.
	BaselineTokens TokenTotals
	// PreGreetPeakCacheCreation is the MAXIMUM cache_creation across the turns
	// before the measured turn — e.g. the shallow-boot-window scenario's former
	// teamRecacheSpikeFloor signal. Left zero when the scenario has no
	// pre-observation window (a single-turn observation) or doesn't track it.
	PreGreetPeakCacheCreation int
	TotalCostUSD              float64
	ModelUsage                map[string]ModelUsage
	// ClaudeCodeVersion is the Claude Code CLI client version (the stream's
	// system/init event's claude_code_version field) that produced this
	// observation, distinct from Model — left empty for non-Claude runtimes or
	// streams that predate the field.
	ClaudeCodeVersion string
	// ResolvedModel is the model identifier the runtime actually resolved at
	// boot (the stream's system/init event's model field), which can be more
	// precise than a CI-matrix alias recorded as JourneySpec.Model — e.g.
	// "sonnet" resolves to "claude-sonnet-4-6". Left empty when unavailable.
	ResolvedModel string
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
	Status       string   `json:"status"`
	Failure      string   `json:"failure,omitempty"`
	Owner        string   `json:"owner,omitempty"`
	FailureCodes []string `json:"failure_codes,omitempty"`
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
	SchemaVersion   int            `json:"schema_version"`
	ScenarioID      string         `json:"scenario_id"`
	JourneyID       string         `json:"journey_id,omitempty"`
	Source          string         `json:"source"`
	Mode            string         `json:"mode,omitempty"`
	Runtime         string         `json:"runtime,omitempty"`
	Executor        string         `json:"executor,omitempty"`
	Host            string         `json:"host"`
	Model           string         `json:"model"`
	MetricsState    MetricsState   `json:"metrics_state"`
	Outcome         Outcome        `json:"outcome"`
	DurationMS      int64          `json:"duration_ms,omitempty"`
	Turns           int            `json:"turns,omitempty"`
	ToolCalls       int            `json:"tool_calls,omitempty"`
	ToolCallsByName map[string]int `json:"tool_calls_by_name,omitempty"`
	StatusReadCalls int            `json:"status_read_calls,omitempty"`
	ScopedReadCalls int            `json:"scoped_read_calls,omitempty"`
	Tokens          TokenTotals    `json:"tokens,omitempty"`
	// BaselineTokens is the first turn's full TokenTotals — a reference point
	// for how the observation started. Per-turn usage is NOT cumulative across
	// turns, so only Tokens.Context() - BaselineTokens.Context() is a
	// defensible derived quantity; Input/Output are not independently
	// subtractable between two unrelated single requests. Empty for scenarios
	// that don't record a pre-observation baseline.
	BaselineTokens TokenTotals `json:"baseline_tokens,omitempty"`
	// PreGreetPeakCacheCreation is the MAXIMUM cache_creation across the turns
	// before the measured turn (e.g. shallow-boot-window's former
	// teamRecacheSpikeFloor signal). Zero when there is no pre-observation
	// window or the scenario doesn't track it.
	PreGreetPeakCacheCreation int                    `json:"pre_greet_peak_cache_creation,omitempty"`
	TotalCostUSD              float64                `json:"total_cost_usd,omitempty"`
	ModelUsage                map[string]ModelUsage  `json:"model_usage,omitempty"`
	Budget                    Budget                 `json:"budget,omitempty"`
	BudgetResult              *BudgetResult          `json:"budget_result,omitempty"`
	CodexCharacter            *CodexCharacterization `json:"codex_characterization,omitempty"`
	// RunID, RunURL, and CapturedAt are run-provenance fields stamped at emission
	// time (see EmitRecord's stampProvenance) so a scenario/model that accumulates
	// multiple observations in a published ledger can be traced back to the run
	// that produced each one, and ordered chronologically. A record emitted outside
	// CI (no GITHUB_RUN_ID) simply omits RunID/RunURL.
	RunID      string `json:"run_id,omitempty"`
	RunURL     string `json:"run_url,omitempty"`
	CapturedAt string `json:"captured_at,omitempty"`
	// ClaudeCodeVersion is the Claude Code CLI client version that produced this
	// record (the stream's system/init event's claude_code_version field), so a
	// future trend reader can attribute a boot-cost shift to a client update
	// rather than the FO's own contract. Empty for non-Claude runtimes or
	// streams captured before the field existed.
	ClaudeCodeVersion string `json:"claude_code_version,omitempty"`
	// ResolvedModel is the model identifier the runtime actually resolved at
	// boot, distinct from Model when Model is a CI-matrix alias (e.g. "sonnet"
	// resolves to "claude-sonnet-4-6") — so a trend reader can tell a silent
	// alias-resolution change from a real regression. Empty when unavailable.
	ResolvedModel string `json:"resolved_model,omitempty"`
}

type recordJSON struct {
	SchemaVersion             int                    `json:"schema_version"`
	ScenarioID                string                 `json:"scenario_id"`
	Source                    string                 `json:"source"`
	Mode                      string                 `json:"mode,omitempty"`
	Runtime                   string                 `json:"runtime,omitempty"`
	Executor                  string                 `json:"executor,omitempty"`
	Host                      string                 `json:"host"`
	Model                     string                 `json:"model"`
	MetricsState              MetricsState           `json:"metrics_state"`
	Outcome                   Outcome                `json:"outcome"`
	DurationMS                int64                  `json:"duration_ms,omitempty"`
	Turns                     int                    `json:"turns,omitempty"`
	ToolCalls                 int                    `json:"tool_calls,omitempty"`
	ToolCallsByName           map[string]int         `json:"tool_calls_by_name,omitempty"`
	StatusReadCalls           int                    `json:"status_read_calls,omitempty"`
	ScopedReadCalls           int                    `json:"scoped_read_calls,omitempty"`
	Tokens                    *TokenTotals           `json:"tokens,omitempty"`
	BaselineTokens            *TokenTotals           `json:"baseline_tokens,omitempty"`
	PreGreetPeakCacheCreation int                    `json:"pre_greet_peak_cache_creation,omitempty"`
	TotalCostUSD              float64                `json:"total_cost_usd,omitempty"`
	ModelUsage                map[string]ModelUsage  `json:"model_usage,omitempty"`
	Budget                    *Budget                `json:"budget,omitempty"`
	BudgetResult              *BudgetResult          `json:"budget_result,omitempty"`
	CodexCharacter            *CodexCharacterization `json:"codex_characterization,omitempty"`
	RunID                     string                 `json:"run_id,omitempty"`
	RunURL                    string                 `json:"run_url,omitempty"`
	CapturedAt                string                 `json:"captured_at,omitempty"`
	ClaudeCodeVersion         string                 `json:"claude_code_version,omitempty"`
	ResolvedModel             string                 `json:"resolved_model,omitempty"`
}

func (r Record) MarshalJSON() ([]byte, error) {
	var tokens *TokenTotals
	if !r.Tokens.isZero() {
		t := r.Tokens.withTotal()
		tokens = &t
	}
	var baselineTokens *TokenTotals
	if !r.BaselineTokens.isZero() {
		t := r.BaselineTokens.withTotal()
		baselineTokens = &t
	}
	var budget *Budget
	if r.Budget.MaxTotalTokens != nil || r.Budget.MaxToolCalls != nil {
		b := r.Budget
		budget = &b
	}
	return marshalRecordJSON(recordJSON{
		SchemaVersion:             RecordSchemaVersion,
		ScenarioID:                firstNonEmpty(r.ScenarioID, r.JourneyID),
		Source:                    r.Source,
		Mode:                      r.Mode,
		Runtime:                   firstNonEmpty(r.Runtime, r.Host),
		Executor:                  r.Executor,
		Host:                      r.Host,
		Model:                     r.Model,
		MetricsState:              r.MetricsState,
		Outcome:                   r.Outcome,
		DurationMS:                r.DurationMS,
		Turns:                     r.Turns,
		ToolCalls:                 r.ToolCalls,
		ToolCallsByName:           r.ToolCallsByName,
		StatusReadCalls:           r.StatusReadCalls,
		ScopedReadCalls:           r.ScopedReadCalls,
		Tokens:                    tokens,
		BaselineTokens:            baselineTokens,
		PreGreetPeakCacheCreation: r.PreGreetPeakCacheCreation,
		TotalCostUSD:              r.TotalCostUSD,
		ModelUsage:                r.ModelUsage,
		Budget:                    budget,
		BudgetResult:              r.BudgetResult,
		CodexCharacter:            r.CodexCharacter,
		RunID:                     r.RunID,
		RunURL:                    r.RunURL,
		CapturedAt:                r.CapturedAt,
		ClaudeCodeVersion:         r.ClaudeCodeVersion,
		ResolvedModel:             r.ResolvedModel,
	})
}
