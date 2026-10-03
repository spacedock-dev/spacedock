package ensigncycle

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Pi recorded-gate semantic grade (test support only). It observes a Pi root
// session plus any successor child sessions and reports unmet obligations by
// stable code. Expected values (canonical id/digest, session model, successor
// marker) are read from the trace, a compact canonical digest prefix is accepted,
// and absent evidence fails closed. It is not a production protocol, CLI verb, or
// stored format.

type piRecordedGateFinding struct{ Code, Detail string }

func piRecordedGateFindingsString(findings []piRecordedGateFinding) string {
	parts := make([]string, 0, len(findings))
	for _, f := range findings {
		parts = append(parts, f.Code+": "+f.Detail)
	}
	return strings.Join(parts, "; ")
}

type piRecordedGateTrace struct {
	root     []map[string]any
	children [][]map[string]any
}

// piRecordedGateObservation is the graded evidence; indices are root-session
// record positions, -1 meaning "not observed".
type piRecordedGateObservation struct {
	requestedModel  string
	boundBriefingID string
	canonicalID     string
	canonicalDigest string
	canonicalIndex  int
	reviewIndex     int
	reviewDigest    string
	childReview     bool
	toolReview      bool
	decisionIndex   int
	consumedIndex   int
	spawnCount      int
	spawnModel      string
	completionIndex int
	reportReadIndex int
}

// piGateSnapshotRE captures the Briefing id and rendered digest from a snapshot
// line, in the canonical room file or the presented review prose.
var piGateSnapshotRE = regexp.MustCompile("Reviewed snapshot:?\\s*`?(briefing:[^`\\s]+)`?\\s*(?:@|at)\\s*`?sha256:([0-9a-fA-F]{6,64})")

// piGateCompletionRE matches an observed subagent completion.
var piGateCompletionRE = regexp.MustCompile(`Outcome:\s*\d+\s+complete|State:\s*complete`)

// piGateDigestRenders accepts the canonical digest or a compact prefix of it.
func piGateDigestRenders(rendered, canonical string) bool {
	rendered, canonical = strings.ToLower(strings.TrimSpace(rendered)), strings.ToLower(strings.TrimSpace(canonical))
	return rendered != "" && len(rendered) >= 6 && len(rendered) <= len(canonical) && strings.HasPrefix(canonical, rendered)
}

func piGateMessage(line map[string]any) (string, map[string]any) {
	msg, _ := line["message"].(map[string]any)
	if msg == nil {
		return "", nil
	}
	role, _ := msg["role"].(string)
	return role, msg
}

func piGateContent(msg map[string]any) []any {
	content, _ := msg["content"].([]any)
	return content
}

func piGateText(msg map[string]any) string {
	var sb strings.Builder
	for _, block := range piGateContent(msg) {
		if m, _ := block.(map[string]any); m != nil {
			if typ, _ := m["type"].(string); typ == "text" {
				text, _ := m["text"].(string)
				sb.WriteString(text)
			}
		}
	}
	return sb.String()
}

// piGateCall is one toolCall with its id, name, and argument map.
type piGateCall struct {
	id, name string
	args     map[string]any
}

func piGateCalls(line map[string]any) []piGateCall {
	role, msg := piGateMessage(line)
	if role != "assistant" {
		return nil
	}
	var calls []piGateCall
	for _, block := range piGateContent(msg) {
		m, _ := block.(map[string]any)
		if m == nil {
			continue
		}
		if typ, _ := m["type"].(string); typ != "toolCall" {
			continue
		}
		id, _ := m["id"].(string)
		name, _ := m["name"].(string)
		args, _ := m["arguments"].(map[string]any)
		calls = append(calls, piGateCall{id, name, args})
	}
	return calls
}

func piGateArg(call piGateCall, key string) string {
	s, _ := call.args[key].(string)
	return s
}

// piGateOutputValue reads a "key=value" token from a recorder stdout line.
func piGateOutputValue(stdout, key string) string {
	for _, field := range strings.Fields(stdout) {
		if value, ok := strings.CutPrefix(field, key+"="); ok {
			return value
		}
	}
	return ""
}

// piGateResult is one tool result with its position and error state.
type piGateResult struct {
	index  int
	callID string
	text   string
	failed bool
}

func piGateResults(root []map[string]any) ([]piGateResult, map[string]piGateResult) {
	var ordered []piGateResult
	byID := map[string]piGateResult{}
	for i, line := range root {
		role, msg := piGateMessage(line)
		if role != "toolResult" {
			continue
		}
		id, _ := msg["toolCallId"].(string)
		failed, _ := msg["isError"].(bool)
		result := piGateResult{index: i, callID: id, text: piGateText(msg), failed: failed}
		ordered = append(ordered, result)
		byID[id] = result
	}
	return ordered, byID
}

// piRecordedGateObserve folds a parsed journey into the graded observation.
func piRecordedGateObserve(trace piRecordedGateTrace) piRecordedGateObservation {
	o := piRecordedGateObservation{canonicalIndex: -1, reviewIndex: -1, decisionIndex: -1, consumedIndex: -1, completionIndex: -1, reportReadIndex: -1}
	for _, line := range trace.root {
		if typ, _ := line["type"].(string); typ == "model_change" {
			provider, _ := line["provider"].(string)
			modelID, _ := line["modelId"].(string)
			if o.requestedModel = modelID; provider != "" && modelID != "" {
				o.requestedModel = provider + "/" + modelID
			}
			break
		}
	}
	results, byID := piGateResults(trace.root)
	for _, result := range results {
		if !result.failed {
			if m := piGateSnapshotRE.FindStringSubmatch(result.text); m != nil {
				o.canonicalIndex, o.canonicalID, o.canonicalDigest = result.index, m[1], m[2]
				break
			}
		}
	}
	for i, line := range trace.root {
		role, msg := piGateMessage(line)
		if role == "assistant" && o.reviewIndex < 0 {
			if m := piGateSnapshotRE.FindStringSubmatch(piGateText(msg)); m != nil && (o.canonicalID == "" || m[1] == o.canonicalID) {
				o.reviewIndex, o.reviewDigest = i, m[2]
			}
		}
		for _, call := range piGateCalls(line) {
			result, ok := byID[call.id]
			command := piGateArg(call, "command")
			switch {
			case (call.name == "bash" || call.name == "shell") && strings.Contains(command, "gate record") && strings.Contains(command, "--briefing"):
				if ok && !result.failed && strings.Contains(result.text, "state=open") && o.boundBriefingID == "" {
					o.boundBriefingID = piGateOutputValue(result.text, "briefing")
				}
			case (call.name == "bash" || call.name == "shell") && strings.Contains(command, "gate record") && strings.Contains(command, "--decision"):
				if ok && !result.failed && strings.Contains(result.text, "state=closed") && o.decisionIndex < 0 {
					o.decisionIndex = i
				}
			case (call.name == "bash" || call.name == "shell") && strings.Contains(command, "gate consume"):
				if ok && !result.failed && strings.Contains(result.text, "consumed=true") && o.consumedIndex < 0 {
					o.consumedIndex = i
				}
			case call.name == "subagent" && piGateArg(call, "task") != "":
				o.spawnCount++
				if o.spawnModel == "" {
					o.spawnModel = piGateArg(call, "model")
				}
			}
		}
	}
	for _, child := range trace.children {
		for _, line := range child {
			if role, msg := piGateMessage(line); role == "assistant" {
				if m := piGateSnapshotRE.FindStringSubmatch(piGateText(msg)); m != nil && (o.canonicalID == "" || m[1] == o.canonicalID) {
					o.childReview = true
				}
			}
		}
	}
	if o.reviewIndex < 0 && !o.childReview {
		for _, result := range results {
			if !result.failed {
				if m := piGateSnapshotRE.FindStringSubmatch(result.text); m != nil && (o.canonicalID == "" || m[1] == o.canonicalID) {
					o.toolReview = true
					break
				}
			}
		}
	}
	if o.spawnCount >= 1 {
		for _, result := range results {
			if o.completionIndex < 0 {
				if !result.failed && piGateCompletionRE.MatchString(result.text) {
					o.completionIndex = result.index
				}
				continue
			}
			if result.index > o.completionIndex && !result.failed && strings.Contains(result.text, recordedGateDispatchMarker) {
				o.reportReadIndex = result.index
				break
			}
		}
	}
	return o
}

// piRecordedGateJudge turns an observation into unmet obligations. disabled names
// obligations a control's named falsifying edit removes; the real grade passes nil.
func piRecordedGateJudge(o piRecordedGateObservation, disabled map[string]bool) []piRecordedGateFinding {
	var findings []piRecordedGateFinding
	add := func(code, detail string) {
		if !disabled[code] {
			findings = append(findings, piRecordedGateFinding{code, detail})
		}
	}
	if o.boundBriefingID == "" {
		add("pi-gate-binding-missing", "no successful `gate record --briefing` bind (state=open) with a briefing id")
	}
	if o.canonicalIndex < 0 {
		add("pi-gate-canonical-authority-missing", "no canonical bound Briefing id/digest read from a successful tool result")
	}
	if o.reviewIndex < 0 {
		switch {
		case o.childReview:
			add("pi-gate-review-not-root-child", "captain-facing gate review appears only in a child session")
		case o.toolReview:
			add("pi-gate-review-not-root-tool", "captain-facing gate review appears only in tool output")
		default:
			add("pi-gate-review-missing", "no root-assistant gate review")
		}
	} else {
		if o.canonicalIndex >= 0 && o.canonicalIndex > o.reviewIndex {
			add("pi-gate-review-before-authority", "canonical bound authority was read after the review was presented")
		}
		if o.decisionIndex >= 0 && o.reviewIndex >= o.decisionIndex {
			add("pi-gate-review-late", "root review was rendered at or after the decision tool call")
		}
		if !piGateDigestRenders(o.reviewDigest, o.canonicalDigest) {
			add("pi-gate-review-digest-mismatch", fmt.Sprintf("reviewed digest %q does not render canonical digest %q", o.reviewDigest, o.canonicalDigest))
		}
	}
	if o.decisionIndex < 0 {
		add("pi-gate-approval-not-recorded", "no successful `gate record --decision` observed")
	}
	if o.consumedIndex < 0 {
		add("pi-gate-approval-not-consumed", "no successful `gate consume ... consumed=true` observed")
	}
	if o.spawnCount != 1 {
		add("pi-gate-successor-count", fmt.Sprintf("successor dispatches = %d, want exactly one", o.spawnCount))
	}
	if o.spawnCount == 1 && o.spawnModel != o.requestedModel {
		add("pi-gate-successor-model-mismatch", fmt.Sprintf("successor requested model %q, want %q", o.spawnModel, o.requestedModel))
	}
	if o.spawnCount >= 1 && o.completionIndex < 0 {
		add("pi-gate-successor-completion-missing", "successor completion was never observed")
	}
	if o.completionIndex >= 0 && o.reportReadIndex < 0 {
		add("pi-gate-report-read-missing", "first officer never read the successor's durable report after completion")
	}
	return findings
}

// piRecordedGateFindings observes and judges a journey.
func piRecordedGateFindings(trace piRecordedGateTrace, disabled map[string]bool) []piRecordedGateFinding {
	return piRecordedGateJudge(piRecordedGateObserve(trace), disabled)
}

func piGateHasFinding(findings []piRecordedGateFinding, code string) bool {
	for _, f := range findings {
		if f.Code == code {
			return true
		}
	}
	return false
}

func piGateSession(t *testing.T, content string) []map[string]any {
	t.Helper()
	var lines []map[string]any
	for _, raw := range strings.Split(content, "\n") {
		if strings.TrimSpace(raw) == "" {
			continue
		}
		var line map[string]any
		if err := json.Unmarshal([]byte(raw), &line); err != nil {
			t.Fatalf("Pi recorded-gate session line is not JSON: %v", err)
		}
		lines = append(lines, line)
	}
	return lines
}

// piRecordedGateRetained loads a retained Pi session from testdata.
func piRecordedGateRetained(t *testing.T, name string) []map[string]any {
	t.Helper()
	return piGateSession(t, readFile(t, filepath.Join("testdata", "pi-recorded-gate", name)))
}

// piRecordedGateTraceFromArtifacts loads the live grade's trace: the single root
// session under sessions/, plus successor children under
// sessions/<root>/<handle>/run-*/session.jsonl.
func piRecordedGateTraceFromArtifacts(t *testing.T, artifactDir string) piRecordedGateTrace {
	t.Helper()
	roots, err := filepath.Glob(filepath.Join(artifactDir, "sessions", "*.jsonl"))
	if err != nil || len(roots) != 1 {
		t.Fatalf("Pi recorded-gate root sessions=%d (%v), want exactly one: %v", len(roots), err, roots)
	}
	trace := piRecordedGateTrace{root: piGateSession(t, readFile(t, roots[0]))}
	children, err := filepath.Glob(filepath.Join(artifactDir, "sessions", "*", "*", "run-*", "session.jsonl"))
	if err != nil {
		t.Fatalf("globbing Pi recorded-gate child sessions: %v", err)
	}
	for _, path := range children {
		trace.children = append(trace.children, piGateSession(t, readFile(t, path)))
	}
	return trace
}

// piGateSnapshotDigestReplace rewrites the digest captured by piGateSnapshotRE in
// the first matching text block of the first record of the given role.
func piGateSnapshotDigestReplace(t *testing.T, lines []map[string]any, role, digest string) {
	t.Helper()
	for i, line := range lines {
		recordRole, msg := piGateMessage(line)
		if recordRole != role {
			continue
		}
		content := piGateContent(msg)
		for j, block := range content {
			m, _ := block.(map[string]any)
			if typ, _ := m["type"].(string); typ != "text" {
				continue
			}
			text, _ := m["text"].(string)
			if loc := piGateSnapshotRE.FindStringSubmatchIndex(text); loc != nil {
				m["text"] = text[:loc[4]] + digest + text[loc[5]:]
				content[j] = m
				msg["content"] = content
				lines[i] = line
				return
			}
		}
	}
	t.Fatalf("no %s snapshot line to edit", role)
}

// TestPiRecordedGateRetainedRootsStayRed replays both retained exact-tip roots and
// asserts each is red only for its named obligation. Falsifiable: a grader that
// stops checking review provenance, or that lets the first root's completed
// lifecycle mask its invented digest, fails here.
func TestPiRecordedGateRetainedRootsStayRed(t *testing.T) {
	first := piRecordedGateFindings(piRecordedGateTrace{root: piRecordedGateRetained(t, "first-root.jsonl"), children: [][]map[string]any{piRecordedGateRetained(t, "first-child.jsonl")}}, nil)
	if !piGateHasFinding(first, "pi-gate-review-digest-mismatch") || !piGateHasFinding(first, "pi-gate-review-before-authority") {
		t.Fatalf("first retained root lost its invented-digest/read-order rejection: %s", piRecordedGateFindingsString(first))
	}
	for _, code := range []string{"pi-gate-approval-not-recorded", "pi-gate-approval-not-consumed", "pi-gate-successor-count", "pi-gate-successor-model-mismatch", "pi-gate-successor-completion-missing", "pi-gate-report-read-missing"} {
		if piGateHasFinding(first, code) {
			t.Fatalf("first retained root reported unrelated %s despite a completed lifecycle: %s", code, piRecordedGateFindingsString(first))
		}
	}
	retry := piRecordedGateFindings(piRecordedGateTrace{root: piRecordedGateRetained(t, "retry-root.jsonl")}, nil)
	if !piGateHasFinding(retry, "pi-gate-approval-not-recorded") {
		t.Fatalf("retry retained root lost its early-stop rejection: %s", piRecordedGateFindingsString(retry))
	}
	for _, code := range []string{"pi-gate-review-digest-mismatch", "pi-gate-review-before-authority", "pi-gate-review-not-root-child", "pi-gate-review-not-root-tool"} {
		if piGateHasFinding(retry, code) {
			t.Fatalf("retry retained root reported unrelated %s despite a correct bound review: %s", code, piRecordedGateFindingsString(retry))
		}
	}
}

// TestPiRecordedGateCompactDigestPrefixAccepted asserts a correct compact
// canonical prefix is not falsely rejected. Falsifiable: a grader that requires
// the full 64-hex digest fails here.
func TestPiRecordedGateCompactDigestPrefixAccepted(t *testing.T) {
	lines := piRecordedGateRetained(t, "first-root.jsonl")
	piGateSnapshotDigestReplace(t, lines, "assistant", "0a54f1baec01")
	if findings := piRecordedGateFindings(piRecordedGateTrace{root: lines}, nil); piGateHasFinding(findings, "pi-gate-review-digest-mismatch") {
		t.Fatalf("compact canonical prefix was rejected as a mismatch: %s", piRecordedGateFindingsString(findings))
	}
}

// TestPiRecordedGateFollowsCanonicalValues asserts the grade follows the trace
// rather than shipping fixture constants: with a different canonical digest
// rendered consistently, no mismatch is reported. Falsifiable: a grader with a
// hardcoded retained digest fails here.
func TestPiRecordedGateFollowsCanonicalValues(t *testing.T) {
	const other = "111122223333444455556666777788889999aaaabbbbccccddddeeeeffff0000"
	lines := piRecordedGateRetained(t, "first-root.jsonl")
	piGateSnapshotDigestReplace(t, lines, "toolResult", other)
	piGateSnapshotDigestReplace(t, lines, "assistant", other)
	if findings := piRecordedGateFindings(piRecordedGateTrace{root: lines}, nil); piGateHasFinding(findings, "pi-gate-review-digest-mismatch") {
		t.Fatalf("consistent alternate canonical values were rejected: %s", piRecordedGateFindingsString(findings))
	}
}
