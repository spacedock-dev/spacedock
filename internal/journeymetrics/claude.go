package journeymetrics

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type ClaudeParseResult struct {
	Observation    Observation
	AssistantUsage TokenTotals
}

func ParseClaudeJSONL(data []byte) (ClaudeParseResult, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	assistantIDs := map[string]bool{}
	assistantUsageSeen := map[string]bool{}
	toolIDs := map[string]string{}
	toolCallsByName := map[string]int{}
	statusReadCalls := 0
	scopedReadCalls := 0
	var assistantUsage TokenTotals
	var terminalUsage *TokenTotals
	var terminalCost float64
	var terminalModelUsage map[string]ModelUsage
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// The live claude runner folds stderr into the stream-json pipe, so
		// non-JSON lines (the front-door launch banner, a 401 launch error, log
		// noise) can appear before or among the JSONL events. Skip them; only
		// JSON object lines are stream-json events.
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var row map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return ClaudeParseResult{}, fmt.Errorf("line %d: %w", lineNo, err)
		}
		typ := rawString(row["type"])
		switch typ {
		case "assistant":
			msg, err := parseClaudeAssistant(row["message"])
			if err != nil {
				return ClaudeParseResult{}, fmt.Errorf("line %d assistant: %w", lineNo, err)
			}
			id := msg.ID
			if id == "" {
				id = fmt.Sprintf("line-%d", lineNo)
			}
			assistantIDs[id] = true
			if !assistantUsageSeen[id] {
				assistantUsageSeen[id] = true
				assistantUsage = assistantUsage.add(msg.Usage)
			}
			for _, block := range msg.Content {
				if block.Type != "tool_use" {
					continue
				}
				toolID := block.ID
				if toolID == "" {
					toolID = id + ":" + block.Name
				}
				if _, seen := toolIDs[toolID]; seen {
					continue
				}
				toolIDs[toolID] = block.Name
				toolCallsByName[block.Name]++
				switch block.Name {
				case "Bash":
					if commandInvokesStatusRead(bashCommand(block.Input)) {
						statusReadCalls++
					}
				case "Read":
					if readInputIsScoped(block.Input) {
						scopedReadCalls++
					}
				}
			}
		case "result":
			result, err := parseClaudeResult(row)
			if err != nil {
				return ClaudeParseResult{}, fmt.Errorf("line %d result: %w", lineNo, err)
			}
			if !result.Usage.isZero() {
				usage := result.Usage.withTotal()
				terminalUsage = &usage
			}
			terminalCost = result.TotalCostUSD
			terminalModelUsage = result.ModelUsage
		}
	}
	if err := scanner.Err(); err != nil {
		return ClaudeParseResult{}, err
	}

	tokens := assistantUsage.withTotal()
	if terminalUsage != nil {
		tokens = *terminalUsage
	}
	if terminalCost == 0 {
		for _, usage := range terminalModelUsage {
			terminalCost += usage.CostUSD
		}
	}
	return ClaudeParseResult{
		Observation: Observation{
			MetricsState:    StateMeasured,
			Turns:           len(assistantIDs),
			ToolCalls:       len(toolIDs),
			ToolCallsByName: toolCallsByName,
			StatusReadCalls: statusReadCalls,
			ScopedReadCalls: scopedReadCalls,
			Tokens:          tokens,
			TotalCostUSD:    terminalCost,
			ModelUsage:      terminalModelUsage,
		},
		AssistantUsage: assistantUsage.withTotal(),
	}, nil
}

// ClaudeTurn is one deduped assistant turn from a stream-json transcript: its
// per-message token usage and the names of the tool_use blocks it emitted. Unlike
// ParseClaudeJSONL (which SUMS usage across the whole run and prefers the terminal
// result usage), this preserves each turn's usage so a caller can measure a single
// turn's context window — e.g. the greet turn's boot-window context.
type ClaudeTurn struct {
	ID    string
	Usage TokenTotals
	// ToolNames is the names of the tool_use blocks in this turn, in order. A
	// turn that dispatches a worker carries an "Agent" (or "TeamCreate") name, so a
	// caller can split the transcript at the first dispatch turn.
	ToolNames []string
	// ReadTargets is the file paths / commands of this turn's read-like tool_use
	// blocks (Read's file_path, Grep's path, Bash's command), so a caller can detect
	// whether a turn read a specific file before some boundary turn.
	ReadTargets []string
	// SkillNames is the skill arguments of this turn's Skill tool_use blocks
	// (Skill's input.skill, e.g. "spacedock:present-gate"), so a caller can detect
	// whether a turn invoked a specific skill before some boundary turn.
	SkillNames []string
}

// Context returns this turn's context-window size as the boot analysis defines it:
// input + cache_read + cache_creation (output is generation, not context).
func (t ClaudeTurn) Context() int {
	return t.Usage.Input + t.Usage.CacheRead + t.Usage.CacheCreation
}

// ParseClaudeTurns walks the stream-json transcript per assistant turn and returns
// one ClaudeTurn per distinct assistant message in stream order. Real runner streams
// are MULTI-DELTA: the same message id appears on several `assistant` rows, the first
// delta carries a `thinking`/`text` block, and the tool_use block(s) land on LATER
// deltas; the per-delta `usage` is identical across the deltas of a message. So this
// MERGES every delta's tool_use names into the turn (a first-delta-only parse drops
// the tool_use entirely — which would make a TeamCreate invisible) while keeping the
// first delta's usage. It reuses the rawTokenUsage field extraction; it does NOT sum
// or prefer the terminal result usage, so each turn's context window is recoverable.
// Non-JSON lines (folded stderr) are skipped, matching ParseClaudeJSONL.
func ParseClaudeTurns(data []byte) ([]ClaudeTurn, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	index := map[string]int{} // message id -> position in turns
	// seenTool dedups tool_use blocks by their unique block id within a message, so a
	// repeated delta (the same tool_use carried again) does not double-count, while a
	// genuinely later, additive tool_use (a different block id) is merged in.
	seenTool := map[string]bool{}
	var turns []ClaudeTurn
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var row map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNo, err)
		}
		if rawString(row["type"]) != "assistant" {
			continue
		}
		msg, err := parseClaudeAssistant(row["message"])
		if err != nil {
			return nil, fmt.Errorf("line %d assistant: %w", lineNo, err)
		}
		id := msg.ID
		if id == "" {
			id = fmt.Sprintf("line-%d", lineNo)
		}
		var names []string
		var readTargets []string
		var skillNames []string
		for _, block := range msg.Content {
			if block.Type != "tool_use" {
				continue
			}
			toolKey := id + ":" + block.ID
			if block.ID != "" && seenTool[toolKey] {
				continue
			}
			seenTool[toolKey] = true
			names = append(names, block.Name)
			if target := readToolTarget(block.Name, block.Input); target != "" {
				readTargets = append(readTargets, target)
			}
			if block.Name == "Skill" {
				if skill := jsonStringField(block.Input, "skill"); skill != "" {
					skillNames = append(skillNames, skill)
				}
			}
		}
		if pos, ok := index[id]; ok {
			// A later delta of a message already seen: merge its NEW tool_use names
			// (the per-block dedup above keeps a repeated delta from double-counting).
			// Usage is identical across deltas, so the first-delta usage is kept.
			turns[pos].ToolNames = append(turns[pos].ToolNames, names...)
			turns[pos].ReadTargets = append(turns[pos].ReadTargets, readTargets...)
			turns[pos].SkillNames = append(turns[pos].SkillNames, skillNames...)
			continue
		}
		index[id] = len(turns)
		turns = append(turns, ClaudeTurn{ID: id, Usage: msg.Usage, ToolNames: names, ReadTargets: readTargets, SkillNames: skillNames})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return turns, nil
}

type claudeAssistant struct {
	ID      string
	Model   string
	Usage   TokenTotals
	Content []claudeContent
}

type claudeContent struct {
	Type  string          `json:"type"`
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`
}

func parseClaudeAssistant(raw json.RawMessage) (claudeAssistant, error) {
	var msg struct {
		ID      string          `json:"id"`
		Model   string          `json:"model"`
		Usage   rawTokenUsage   `json:"usage"`
		Content []claudeContent `json:"content"`
	}
	if len(raw) == 0 {
		return claudeAssistant{}, nil
	}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return claudeAssistant{}, err
	}
	return claudeAssistant{
		ID:      msg.ID,
		Model:   msg.Model,
		Usage:   msg.Usage.totals(),
		Content: msg.Content,
	}, nil
}

type claudeResult struct {
	Usage        TokenTotals
	TotalCostUSD float64
	ModelUsage   map[string]ModelUsage
}

func parseClaudeResult(row map[string]json.RawMessage) (claudeResult, error) {
	var usage rawTokenUsage
	if raw := row["usage"]; len(raw) != 0 {
		if err := json.Unmarshal(raw, &usage); err != nil {
			return claudeResult{}, err
		}
	}
	var cost float64
	if raw := row["total_cost_usd"]; len(raw) != 0 {
		_ = json.Unmarshal(raw, &cost)
	}
	if cost == 0 {
		if raw := row["totalCostUSD"]; len(raw) != 0 {
			_ = json.Unmarshal(raw, &cost)
		}
	}
	modelUsage, err := parseModelUsage(firstRaw(row, "modelUsage", "model_usage"))
	if err != nil {
		return claudeResult{}, err
	}
	return claudeResult{Usage: usage.totals(), TotalCostUSD: cost, ModelUsage: modelUsage}, nil
}

func parseModelUsage(raw json.RawMessage) (map[string]ModelUsage, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var byModel map[string]rawModelUsage
	if err := json.Unmarshal(raw, &byModel); err != nil {
		return nil, err
	}
	out := make(map[string]ModelUsage, len(byModel))
	for model, usage := range byModel {
		out[model] = ModelUsage{Tokens: usage.rawTokenUsage.totals(), CostUSD: usage.cost()}
	}
	return out, nil
}

type rawModelUsage struct {
	rawTokenUsage
	CostSnake float64 `json:"cost_usd"`
	CostCamel float64 `json:"costUSD"`
}

func (u rawModelUsage) cost() float64 {
	if u.CostSnake != 0 {
		return u.CostSnake
	}
	return u.CostCamel
}

type rawTokenUsage struct {
	InputSnake         int `json:"input_tokens"`
	OutputSnake        int `json:"output_tokens"`
	CacheCreationSnake int `json:"cache_creation_input_tokens"`
	CacheReadSnake     int `json:"cache_read_input_tokens"`
	TotalSnake         int `json:"total_tokens"`

	InputCamel         int `json:"inputTokens"`
	OutputCamel        int `json:"outputTokens"`
	CacheCreationCamel int `json:"cacheCreationInputTokens"`
	CacheReadCamel     int `json:"cacheReadInputTokens"`
	TotalCamel         int `json:"totalTokens"`
}

func (u rawTokenUsage) totals() TokenTotals {
	t := TokenTotals{
		Input:         preferInt(u.InputSnake, u.InputCamel),
		Output:        preferInt(u.OutputSnake, u.OutputCamel),
		CacheCreation: preferInt(u.CacheCreationSnake, u.CacheCreationCamel),
		CacheRead:     preferInt(u.CacheReadSnake, u.CacheReadCamel),
		Total:         preferInt(u.TotalSnake, u.TotalCamel),
	}
	return t.withTotal()
}

func preferInt(a, b int) int {
	if a != 0 {
		return a
	}
	return b
}

func rawString(raw json.RawMessage) string {
	var s string
	_ = json.Unmarshal(raw, &s)
	return s
}

func firstRaw(row map[string]json.RawMessage, keys ...string) json.RawMessage {
	for _, key := range keys {
		if raw := row[key]; len(raw) != 0 {
			return raw
		}
	}
	return nil
}
