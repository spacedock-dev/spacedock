package journeymetrics

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type CodexCharacterization struct {
	Model           string              `json:"model,omitempty"`
	EventKinds      []string            `json:"event_kinds"`
	FieldsByEvent   map[string][]string `json:"fields_by_event"`
	ToolCalls       int                 `json:"tool_calls,omitempty"`
	ToolCallsByName map[string]int      `json:"tool_calls_by_name,omitempty"`
}

func CharacterizeCodexExecJSONL(data []byte) (CodexCharacterization, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	fields := map[string]map[string]bool{}
	toolSeen := map[string]bool{}
	toolCallsByName := map[string]int{}
	model := ""
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var row map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return CodexCharacterization{}, fmt.Errorf("line %d: %w", lineNo, err)
		}
		typ := rawString(row["type"])
		if typ == "" {
			continue
		}
		if model == "" {
			model = rawString(row["model"])
		}
		if typ == "tool_call.started" {
			callID := rawString(row["call_id"])
			if callID == "" {
				callID = fmt.Sprintf("line-%d", lineNo)
			}
			if !toolSeen[callID] {
				toolSeen[callID] = true
				toolCallsByName[rawString(row["name"])]++
			}
		}
		if fields[typ] == nil {
			fields[typ] = map[string]bool{}
		}
		for key := range row {
			fields[typ][key] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return CodexCharacterization{}, err
	}
	kinds := make([]string, 0, len(fields))
	fieldsByEvent := make(map[string][]string, len(fields))
	for kind, byField := range fields {
		kinds = append(kinds, kind)
		for field := range byField {
			fieldsByEvent[kind] = append(fieldsByEvent[kind], field)
		}
		sort.Strings(fieldsByEvent[kind])
	}
	sort.Strings(kinds)
	return CodexCharacterization{
		Model:           model,
		EventKinds:      kinds,
		FieldsByEvent:   fieldsByEvent,
		ToolCalls:       len(toolSeen),
		ToolCallsByName: toolCallsByName,
	}, nil
}

func CodexCharacterizedRecord(spec JourneySpec, characterization CodexCharacterization, result BehaviorResult) Record {
	spec.Model = firstNonEmpty(spec.Model, characterization.Model)
	record := BuildRecord(spec, result, Observation{MetricsState: StateCharacterized})
	record.CodexCharacter = &characterization
	return record
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
