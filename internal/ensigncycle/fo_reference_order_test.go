package ensigncycle

import (
	"bufio"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type foReferenceEvent string

const (
	foSharedRead    foReferenceEvent = "shared-read"
	foRuntimeRead   foReferenceEvent = "runtime-read"
	foWriteRead     foReferenceEvent = "write-read"
	foMergeRead     foReferenceEvent = "merge-read"
	foModBlockSeen  foReferenceEvent = "merge-mod-block-seen"
	foEngage        foReferenceEvent = "engage"
	foMutation      foReferenceEvent = "mutation"
	foTerminal      foReferenceEvent = "terminal-mutation"
	foMergeGuard    foReferenceEvent = "merge-guard"
	foMergeAction   foReferenceEvent = "merge-action"
	foWrongPath     foReferenceEvent = "wrong-core-path"
	foFailedRead    foReferenceEvent = "failed-core-read"
	foFailedEngage  foReferenceEvent = "failed-engage"
	foBroadSearch   foReferenceEvent = "broad-core-search"
	foWrapperSkill  foReferenceEvent = "wrapper-skill"
	foRepeatedMerge foReferenceEvent = "repeated-merge-work"
)

type foIndexedEvent struct {
	index int
	kind  foReferenceEvent
}

type claudePendingFOCall struct {
	events []foReferenceEvent
}

var (
	firstOfficerBaseAnnouncementRE = regexp.MustCompile(`Base directory for this skill:[[:space:]]*([^\\"\r\n]+/first-officer)`)
	firstOfficerEntryPathRE        = regexp.MustCompile(`([^"'[:space:];|&]+/skills/first-officer)/SKILL\.md`)
)

func findFirstOfficerSkillBase(stream string) string {
	if match := firstOfficerBaseAnnouncementRE.FindStringSubmatch(stream); len(match) == 2 {
		return filepath.Clean(match[1])
	}
	scanner := bufio.NewScanner(strings.NewReader(stream))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var row struct {
			Type string `json:"type"`
			Item struct {
				Type             string `json:"type"`
				Command          string `json:"command"`
				AggregatedOutput string `json:"aggregated_output"`
				Output           string `json:"output"`
				Status           string `json:"status"`
				ExitCode         *int   `json:"exit_code"`
			} `json:"item"`
		}
		if json.Unmarshal(scanner.Bytes(), &row) != nil || row.Type != "item.completed" || row.Item.Type != "command_execution" {
			continue
		}
		if row.Item.Status == "failed" || (row.Item.ExitCode != nil && *row.Item.ExitCode != 0) {
			continue
		}
		if !strings.Contains(strings.ToLower(row.Item.AggregatedOutput+row.Item.Output), "name: first-officer") {
			continue
		}
		if match := firstOfficerEntryPathRE.FindStringSubmatch(row.Item.Command); len(match) == 2 {
			return filepath.Clean(match[1])
		}
	}
	return ""
}

func normalizeClaudeFOReferenceEvents(stream string) []foReferenceEvent {
	var events []foReferenceEvent
	skillBase := findFirstOfficerSkillBase(stream)
	seenTools := map[string]bool{}
	pending := map[string]claudePendingFOCall{}
	scanner := bufio.NewScanner(strings.NewReader(stream))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var row struct {
			Type    string `json:"type"`
			Message struct {
				ID      string            `json:"id"`
				Content []json.RawMessage `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(scanner.Bytes(), &row) != nil {
			continue
		}
		switch row.Type {
		case "assistant":
			for _, raw := range row.Message.Content {
				var block struct {
					Type  string         `json:"type"`
					ID    string         `json:"id"`
					Name  string         `json:"name"`
					Input map[string]any `json:"input"`
				}
				if json.Unmarshal(raw, &block) != nil || block.Type != "tool_use" {
					continue
				}
				key := row.Message.ID + ":" + block.ID
				if block.ID != "" && seenTools[key] {
					continue
				}
				seenTools[key] = true
				switch block.Name {
				case "Read":
					pending[block.ID] = claudePendingFOCall{events: classifyFORead(fmt.Sprint(block.Input["file_path"]), skillBase)}
				case "Bash":
					pending[block.ID] = claudePendingFOCall{events: classifyFOCommand(fmt.Sprint(block.Input["command"]), skillBase)}
				case "Skill":
					skill := fmt.Sprint(block.Input["skill"])
					if skill == "spacedock:fo-write-core" || skill == "spacedock:fo-merge-core" {
						events = append(events, foWrapperSkill)
					}
				case "Edit", "Write", "NotebookEdit":
					events = append(events, foMutation)
				}
			}
		case "user":
			for _, raw := range row.Message.Content {
				var block struct {
					Type      string `json:"type"`
					ToolUseID string `json:"tool_use_id"`
					IsError   bool   `json:"is_error"`
					Content   any    `json:"content"`
				}
				if json.Unmarshal(raw, &block) != nil || block.Type != "tool_result" {
					continue
				}
				if call, ok := pending[block.ToolUseID]; ok {
					events = append(events, resolvedFOCallEvents(call.events, foCallSucceeded(call.events, !block.IsError, fmt.Sprint(block.Content)))...)
					delete(pending, block.ToolUseID)
				}
				if !block.IsError && containsMergeModBlock(fmt.Sprint(block.Content)) {
					events = append(events, foModBlockSeen)
				}
			}
		}
	}
	return events
}

func normalizeCodexFOReferenceEvents(stream string) []foReferenceEvent {
	var events []foReferenceEvent
	skillBase := findFirstOfficerSkillBase(stream)
	scanner := bufio.NewScanner(strings.NewReader(stream))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var row struct {
			Type string `json:"type"`
			Item struct {
				Type             string `json:"type"`
				Command          string `json:"command"`
				AggregatedOutput string `json:"aggregated_output"`
				Output           string `json:"output"`
				Status           string `json:"status"`
				ExitCode         *int   `json:"exit_code"`
			} `json:"item"`
		}
		if json.Unmarshal(scanner.Bytes(), &row) != nil || row.Type != "item.completed" {
			continue
		}
		switch row.Item.Type {
		case "command_execution":
			succeeded := row.Item.Status != "failed" && (row.Item.ExitCode == nil || *row.Item.ExitCode == 0)
			callEvents := classifyFOCommand(row.Item.Command, skillBase)
			events = append(events, resolvedFOCallEvents(callEvents, foCallSucceeded(callEvents, succeeded, row.Item.AggregatedOutput+row.Item.Output))...)
			if succeeded && containsMergeModBlock(row.Item.AggregatedOutput+row.Item.Output) {
				events = append(events, foModBlockSeen)
			}
		case "file_change":
			events = append(events, foMutation)
		}
	}
	return events
}

func foCallSucceeded(events []foReferenceEvent, transportSucceeded bool, output string) bool {
	hasRead := false
	for _, event := range events {
		switch event {
		case foSharedRead, foRuntimeRead, foWriteRead, foMergeRead:
			hasRead = true
		}
	}
	if !hasRead {
		return transportSucceeded
	}
	anchors := map[foReferenceEvent]string{
		foSharedRead:  "# first officer shared core",
		foRuntimeRead: "first officer runtime",
		foWriteRead:   "# first officer write core",
		foMergeRead:   "# first officer merge core",
	}
	var required []string
	for _, event := range events {
		if anchor := anchors[event]; anchor != "" {
			required = append(required, anchor)
		}
	}
	return ReferenceReadSucceeded(transportSucceeded, output, required...)
}

func resolvedFOCallEvents(events []foReferenceEvent, succeeded bool) []foReferenceEvent {
	if succeeded {
		return events
	}
	out := make([]foReferenceEvent, 0, len(events))
	for _, event := range events {
		switch event {
		case foSharedRead, foRuntimeRead, foWriteRead, foMergeRead:
			out = append(out, foFailedRead)
		case foEngage:
			out = append(out, foFailedEngage)
		default:
			out = append(out, event)
		}
	}
	return out
}

func classifyFORead(target, skillBase string) []foReferenceEvent {
	cleanTarget := filepath.Clean(target)
	expected := func(suffix string) string { return filepath.Join(skillBase, filepath.FromSlash(suffix)) }
	switch {
	case skillBase != "" && cleanTarget == expected("references/first-officer-shared-core.md"):
		return []foReferenceEvent{foSharedRead}
	case skillBase != "" && (cleanTarget == expected("references/claude-first-officer-runtime.md") ||
		cleanTarget == expected("references/codex-first-officer-runtime.md") ||
		cleanTarget == expected("references/pi-first-officer-runtime.md")):
		return []foReferenceEvent{foRuntimeRead}
	case skillBase != "" && cleanTarget == expected("references/fo-write-core.md"):
		return []foReferenceEvent{foWriteRead}
	case skillBase != "" && cleanTarget == expected("references/fo-merge-core.md"):
		return []foReferenceEvent{foMergeRead}
	case strings.Contains(target, "first-officer-shared-core.md"), strings.Contains(target, "first-officer-runtime.md"),
		strings.Contains(target, "fo-write-core.md"), strings.Contains(target, "fo-merge-core.md"):
		return []foReferenceEvent{foWrongPath}
	default:
		return nil
	}
}

func classifyFOCommand(command, skillBase string) []foReferenceEvent {
	lower := strings.ToLower(command)
	indexed := classifyShellFOReadTargets(command, skillBase)
	if (strings.Contains(lower, "fo-write-core") || strings.Contains(lower, "fo-merge-core")) &&
		(strings.Contains(lower, "find ") || strings.Contains(lower, "grep -r") || strings.Contains(lower, "rg --files") || strings.Contains(lower, "ls -r")) {
		indexed = append(indexed, foIndexedEvent{kind: foBroadSearch})
	}
	if strings.Contains(lower, "spacedock:fo-write-core") || strings.Contains(lower, "spacedock:fo-merge-core") {
		indexed = append(indexed, foIndexedEvent{kind: foWrapperSkill})
	}
	for _, repeated := range []string{"gh pr create", "git merge", "git push"} {
		if at := strings.Index(lower, repeated); at >= 0 {
			indexed = append(indexed, foIndexedEvent{index: at, kind: foRepeatedMerge})
		}
	}
	if at := strings.Index(lower, " state ready"); at >= 0 {
		indexed = append(indexed, foIndexedEvent{index: at, kind: foEngage})
	}
	mutationAt := shellMutationIndex(lower)
	for _, needle := range []string{" status ", "spacedock status", " state commit ", "spacedock state commit", " dispatch build ", "spacedock dispatch build", " new ", "spacedock new", " --archive "} {
		if at := strings.Index(lower, needle); at >= 0 && (strings.Contains(lower[at:], "--set") || strings.Contains(needle, "commit") || strings.Contains(needle, "build") || strings.Contains(needle, "new") || strings.Contains(needle, "archive")) {
			if mutationAt < 0 || at < mutationAt {
				mutationAt = at
			}
		}
	}
	if at := strings.Index(lower, "merge guard"); at >= 0 {
		indexed = append(indexed, foIndexedEvent{index: at, kind: foMergeGuard})
		indexed = append(indexed, foIndexedEvent{index: at, kind: foMergeAction})
		if mutationAt < 0 || at < mutationAt {
			mutationAt = at
		}
	}
	if mutationAt >= 0 {
		indexed = append(indexed, foIndexedEvent{index: mutationAt, kind: foMutation})
	}
	if at := terminalMutationIndex(lower); at >= 0 {
		indexed = append(indexed, foIndexedEvent{index: at, kind: foTerminal})
		indexed = append(indexed, foIndexedEvent{index: at, kind: foMergeAction})
	}
	for _, needle := range []string{"mod-block=", "pr=pr-merge:", "pr=local-merge:"} {
		if at := strings.Index(lower, needle); at >= 0 {
			indexed = append(indexed, foIndexedEvent{index: at, kind: foMergeAction})
		}
	}
	sort.SliceStable(indexed, func(i, j int) bool { return indexed[i].index < indexed[j].index })
	events := make([]foReferenceEvent, 0, len(indexed))
	for _, event := range indexed {
		events = append(events, event.kind)
	}
	return events
}

var shellReadCommandRE = regexp.MustCompile(`(?:^|[[:space:]])(?:cat|sed|head|tail|less|more|bat|awk)(?:[[:space:]]|$)`)

func shellReadTargetIndices(command, target string) []int {
	var indices []int
	for from := 0; ; {
		at := strings.Index(command[from:], target)
		if at < 0 {
			break
		}
		at += from
		start := strings.LastIndexAny(command[:at], ";\n|&") + 1
		if shellReadCommandRE.MatchString(command[start:at]) {
			indices = append(indices, at)
		}
		from = at + len(target)
	}
	return indices
}

func classifyShellFOReadTargets(command, skillBase string) []foIndexedEvent {
	targets := []struct {
		name, suffix string
		kind         foReferenceEvent
	}{
		{"first-officer-shared-core.md", "references/first-officer-shared-core.md", foSharedRead},
		{"claude-first-officer-runtime.md", "references/claude-first-officer-runtime.md", foRuntimeRead},
		{"codex-first-officer-runtime.md", "references/codex-first-officer-runtime.md", foRuntimeRead},
		{"pi-first-officer-runtime.md", "references/pi-first-officer-runtime.md", foRuntimeRead},
		{"fo-write-core.md", "references/fo-write-core.md", foWriteRead},
		{"fo-merge-core.md", "references/fo-merge-core.md", foMergeRead},
	}
	var indexed []foIndexedEvent
	for _, target := range targets {
		for _, at := range shellReadTargetIndices(command, target.name) {
			path := shellTokenAt(command, at)
			expected := ""
			if skillBase != "" {
				expected = filepath.Join(skillBase, filepath.FromSlash(target.suffix))
			}
			kind := foWrongPath
			if expected != "" && filepath.Clean(path) == expected {
				kind = target.kind
			}
			indexed = append(indexed, foIndexedEvent{index: at, kind: kind})
		}
	}
	return indexed
}

func shellTokenAt(command string, at int) string {
	isBoundary := func(r byte) bool {
		return r == ' ' || r == '\t' || r == '\r' || r == '\n' || strings.ContainsRune("\"'`;|&() <>", rune(r))
	}
	start := at
	for start > 0 && !isBoundary(command[start-1]) {
		start--
	}
	end := at
	for end < len(command) && !isBoundary(command[end]) {
		end++
	}
	return command[start:end]
}

var (
	sedInPlaceRE   = regexp.MustCompile(`(?:^|[;&|[:space:]])sed[[:space:]]+(?:-[a-z]*i[a-z]*|-i\.[^[:space:]]+)`)
	fileMutationRE = regexp.MustCompile(`(?:^|[;&|[:space:]])(?:mv|cp|rm|touch|mkdir|rmdir|install|apply_patch)[[:space:]]`)
	gitMutationRE  = regexp.MustCompile(`\bgit(?:[[:space:]]+(?:-[cC][[:space:]]+[^[:space:]]+|--[^[:space:]]+))*[[:space:]]+(?:add|commit|push|merge|rebase|reset|checkout|switch|tag|rm|mv|restore|clean|update-index)\b`)
	redirectionRE  = regexp.MustCompile(`(?:^|[^0-9<])>{1,2}[[:space:]]*[^>&]`)
)

func shellMutationIndex(command string) int {
	first := -1
	for _, re := range []*regexp.Regexp{sedInPlaceRE, fileMutationRE, gitMutationRE, redirectionRE} {
		if loc := re.FindStringIndex(command); loc != nil && (first < 0 || loc[0] < first) {
			first = loc[0]
		}
	}
	return first
}

func terminalMutationIndex(command string) int {
	for _, needle := range []string{"status=done", "--archive", "merge guard"} {
		if at := strings.Index(command, needle); at >= 0 {
			return at
		}
	}
	return -1
}

func containsMergeModBlock(output string) bool {
	compact := strings.ReplaceAll(strings.ToLower(output), " ", "")
	return strings.Contains(compact, "mod-block:merge:") || strings.Contains(compact, `"mod-block":"merge:`)
}

func assertFOReferenceJourney(events []foReferenceEvent, journey string) error {
	for _, hazard := range []foReferenceEvent{foWrongPath, foFailedRead, foFailedEngage, foBroadSearch, foWrapperSkill, foRepeatedMerge} {
		if eventIndex(events, hazard) >= 0 {
			return fmt.Errorf("%s contains forbidden %s event: %v", journey, hazard, events)
		}
	}
	if journey == "gate" {
		if eventIndex(events, foSharedRead) < 0 || eventIndex(events, foRuntimeRead) < 0 {
			return fmt.Errorf("cold gate lacks shared/runtime reads: %v", events)
		}
		if eventIndex(events, foWriteRead) >= 0 || eventIndex(events, foMergeRead) >= 0 {
			return fmt.Errorf("cold gate eagerly read a deferred core: %v", events)
		}
		return nil
	}
	writeAt := eventIndex(events, foWriteRead)
	if writeAt < 0 || eventCount(events, foWriteRead) != 1 {
		return fmt.Errorf("%s write-core read count=%d, want exactly 1: %v", journey, eventCount(events, foWriteRead), events)
	}
	mutationAt := eventIndex(events, foMutation)
	if mutationAt < 0 || writeAt >= mutationAt {
		return fmt.Errorf("%s write-core must precede first mutation: %v", journey, events)
	}
	if journey == "filing" || journey == "rejection" {
		if eventIndex(events, foMergeRead) >= 0 {
			return fmt.Errorf("%s read merge-core before a terminal boundary: %v", journey, events)
		}
		return nil
	}
	mergeAt := eventIndex(events, foMergeRead)
	if mergeAt < 0 || eventCount(events, foMergeRead) != 1 || writeAt >= mergeAt {
		return fmt.Errorf("%s requires one write read before one merge read: %v", journey, events)
	}
	guardAt := eventIndex(events, foMergeGuard)
	terminalAt := eventIndex(events, foTerminal)
	if guardAt < 0 && terminalAt < 0 {
		return fmt.Errorf("%s lacks a terminal/merge boundary: %v", journey, events)
	}
	boundary := guardAt
	if boundary < 0 || (terminalAt >= 0 && terminalAt < boundary) {
		boundary = terminalAt
	}
	if mergeAt >= boundary {
		return fmt.Errorf("%s merge-core must precede terminal/merge action: %v", journey, events)
	}
	if journey == "recovery" {
		modAt := eventIndex(events, foModBlockSeen)
		engageAt := eventIndex(events, foEngage)
		actionAt := eventIndex(events, foMergeAction)
		if modAt < 0 || engageAt < 0 || actionAt < 0 || !(modAt < engageAt && engageAt < writeAt && writeAt < mergeAt && mergeAt < actionAt) {
			return fmt.Errorf("recovery order must be boot mod-block → engage → write → merge → first merge action: %v", events)
		}
	}
	return nil
}

func eventIndex(events []foReferenceEvent, want foReferenceEvent) int {
	for i, event := range events {
		if event == want {
			return i
		}
	}
	return -1
}

func eventCount(events []foReferenceEvent, want foReferenceEvent) int {
	count := 0
	for _, event := range events {
		if event == want {
			count++
		}
	}
	return count
}
