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
			Tokens:          tokens,
			TotalCostUSD:    terminalCost,
			ModelUsage:      terminalModelUsage,
		},
		AssistantUsage: assistantUsage.withTotal(),
	}, nil
}

type claudeAssistant struct {
	ID      string
	Model   string
	Usage   TokenTotals
	Content []claudeContent
}

type claudeContent struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
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
