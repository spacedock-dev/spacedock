package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Deterministic controls for the Pi recorded-gate grade. Each control applies one
// named edit to a retained trace so that one obligation is violated, then asserts
// (a) the full grade names that obligation's code — the historical run's other
// errors cannot mask it because the code is asserted directly — and (b) a grade
// with that single rejection removed no longer reports it. Assertion (b) is the
// falsifier: deleting the check makes its source-named control fail here.

// piGateSnapshotBlockIndex returns the (line, block) indices of the first
// assistant snapshot text block, or (-1, -1).
func piGateSnapshotBlockIndex(t *testing.T, lines []map[string]any) (int, int) {
	t.Helper()
	for i, line := range lines {
		role, msg := piGateMessage(line)
		if role != "assistant" {
			continue
		}
		for j, block := range piGateContent(msg) {
			m, _ := block.(map[string]any)
			if typ, _ := m["type"].(string); typ == "text" {
				if text, _ := m["text"].(string); piGateSnapshotRE.MatchString(text) {
					return i, j
				}
			}
		}
	}
	return -1, -1
}

// piGateRemoveRootReview removes the first root-assistant snapshot text block and
// returns its text.
func piGateRemoveRootReview(t *testing.T, lines []map[string]any) string {
	t.Helper()
	i, j := piGateSnapshotBlockIndex(t, lines)
	if i < 0 {
		t.Fatalf("no root review block to remove")
	}
	_, msg := piGateMessage(lines[i])
	content := piGateContent(msg)
	text, _ := content[j].(map[string]any)["text"].(string)
	msg["content"] = append(append([]any{}, content[:j]...), content[j+1:]...)
	return text
}

// piGateTextRecord builds a Pi session record carrying one text block.
func piGateTextRecord(role, text string) map[string]any {
	return map[string]any{"type": "message", "message": map[string]any{
		"role": role, "content": []any{map[string]any{"type": "text", "text": text}},
	}}
}

// piGateAppendAssistantText appends a plain assistant text record.
func piGateAppendAssistantText(lines []map[string]any, text string) []map[string]any {
	return append(lines, piGateTextRecord("assistant", text))
}

// piGateAppendToolResult appends a non-error toolResult record carrying text.
func piGateAppendToolResult(lines []map[string]any, text string) []map[string]any {
	record := piGateTextRecord("toolResult", text)
	record["message"].(map[string]any)["toolCallId"] = "call-injected-review"
	record["message"].(map[string]any)["toolName"] = "read"
	record["message"].(map[string]any)["isError"] = false
	return append(lines, record)
}

// piGateDropResults removes every toolResult record whose text matches.
func piGateDropResults(lines []map[string]any, match func(string) bool) []map[string]any {
	out := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		if role, msg := piGateMessage(line); role == "toolResult" && match(piGateText(msg)) {
			continue
		}
		out = append(out, line)
	}
	return out
}

// piGateSetSuccessorModel rewrites the requested model on the successor spawn.
func piGateSetSuccessorModel(t *testing.T, lines []map[string]any, model string) {
	t.Helper()
	for _, line := range lines {
		for _, call := range piGateCalls(line) {
			if call.name == "subagent" && piGateArg(call, "task") != "" {
				call.args["model"] = model
				return
			}
		}
	}
	t.Fatalf("no successor spawn to edit")
}

// TestPiRecordedGateControls runs the single-obligation falsifier table.
func TestPiRecordedGateControls(t *testing.T) {
	type control struct {
		obligation string
		absent     []string // sibling codes the edit must not also trigger
		trace      func(t *testing.T) piRecordedGateTrace
	}
	controls := map[string]control{
		// Retained first root: completed lifecycle, invented review digest.
		"accept-wrong-digest": {"pi-gate-review-digest-mismatch", nil, func(t *testing.T) piRecordedGateTrace {
			return piRecordedGateTrace{root: piRecordedGateRetained(t, "first-root.jsonl")}
		}},
		// Retained retry root: correct bound review, stops before recording approval.
		"accept-early-stop": {"pi-gate-approval-not-recorded", nil, func(t *testing.T) piRecordedGateTrace {
			return piRecordedGateTrace{root: piRecordedGateRetained(t, "retry-root.jsonl")}
		}},
		"accept-child-review": {"pi-gate-review-not-root-child", []string{"pi-gate-review-not-root-tool"}, func(t *testing.T) piRecordedGateTrace {
			root := piRecordedGateRetained(t, "first-root.jsonl")
			review := piGateRemoveRootReview(t, root)
			child := piGateAppendAssistantText(piRecordedGateRetained(t, "first-child.jsonl"), review)
			return piRecordedGateTrace{root: root, children: [][]map[string]any{child}}
		}},
		"accept-tool-review": {"pi-gate-review-not-root-tool", []string{"pi-gate-review-not-root-child"}, func(t *testing.T) piRecordedGateTrace {
			root := piRecordedGateRetained(t, "first-root.jsonl")
			review := piGateRemoveRootReview(t, root)
			return piRecordedGateTrace{root: piGateAppendToolResult(root, review)}
		}},
		"accept-late-review": {"pi-gate-review-late", []string{"pi-gate-review-before-authority", "pi-gate-review-digest-mismatch"}, func(t *testing.T) piRecordedGateTrace {
			root := piRecordedGateRetained(t, "first-root.jsonl")
			// Fix the retained invented digest first so this edit isolates placement.
			canonical := piRecordedGateObserve(piRecordedGateTrace{root: root}).canonicalDigest
			piGateSnapshotDigestReplace(t, root, "assistant", canonical)
			review := piGateRemoveRootReview(t, root)
			return piRecordedGateTrace{root: piGateAppendAssistantText(root, review)}
		}},
		"ignore-model": {"pi-gate-successor-model-mismatch", []string{"pi-gate-successor-count"}, func(t *testing.T) piRecordedGateTrace {
			root := piRecordedGateRetained(t, "first-root.jsonl")
			piGateSetSuccessorModel(t, root, "someone-else/model")
			return piRecordedGateTrace{root: root}
		}},
		"ignore-completion": {"pi-gate-successor-completion-missing", []string{"pi-gate-report-read-missing"}, func(t *testing.T) piRecordedGateTrace {
			return piRecordedGateTrace{root: piGateDropResults(piRecordedGateRetained(t, "first-root.jsonl"), piGateCompletionRE.MatchString)}
		}},
		"ignore-report-read": {"pi-gate-report-read-missing", []string{"pi-gate-successor-completion-missing"}, func(t *testing.T) piRecordedGateTrace {
			root := piGateDropResults(piRecordedGateRetained(t, "first-root.jsonl"), func(text string) bool {
				return strings.Contains(text, recordedGateDispatchMarker)
			})
			return piRecordedGateTrace{root: root}
		}},
	}
	for name, c := range controls {
		t.Run(name, func(t *testing.T) {
			full := piRecordedGateFindings(c.trace(t), nil)
			if !piGateHasFinding(full, c.obligation) {
				t.Fatalf("full grade did not report %s: %s", c.obligation, piRecordedGateFindingsString(full))
			}
			for _, code := range c.absent {
				if piGateHasFinding(full, code) {
					t.Fatalf("control for %s also triggered sibling %s: %s", c.obligation, code, piRecordedGateFindingsString(full))
				}
			}
			mutant := piRecordedGateFindings(c.trace(t), map[string]bool{c.obligation: true})
			if piGateHasFinding(mutant, c.obligation) {
				t.Fatalf("falsifying edit that removes %s still reported it: %s", c.obligation, piRecordedGateFindingsString(mutant))
			}
		})
	}
}

// TestPiRecordedGateCleanObservationPasses is the positive control at the judge
// seam: a fully satisfied observation is accepted. It is a unit-level control, not
// a hand-authored passing transcript or a claim of live conduct — no live passing
// Pi trace exists in this checkout.
func TestPiRecordedGateCleanObservationPasses(t *testing.T) {
	observation := piRecordedGateObservation{
		requestedModel: "openrouter/openai/gpt-5.4", boundBriefingID: "briefing:sample:stage:attempt-1:revision-1",
		canonicalID: "briefing:sample:stage:attempt-1:revision-1", canonicalDigest: "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
		canonicalIndex: 4, reviewIndex: 8, reviewDigest: "abcdef0123456789", decisionIndex: 12, consumedIndex: 14,
		spawnCount: 1, spawnModel: "openrouter/openai/gpt-5.4", completionIndex: 18, reportReadIndex: 20,
	}
	if findings := piRecordedGateJudge(observation, nil); len(findings) != 0 {
		t.Fatalf("a fully satisfied observation was rejected: %s", piRecordedGateFindingsString(findings))
	}
}

// TestPiRecordedGateArtifactLoader pins the live callback's session discovery:
// one root directly under sessions/ and one child under
// sessions/<root>/<handle>/run-*/session.jsonl.
func TestPiRecordedGateArtifactLoader(t *testing.T) {
	dir := t.TempDir()
	rootDir := filepath.Join(dir, "sessions")
	childDir := filepath.Join(rootDir, "2026-01-01T00-00-00-000Z_aaaa", "bbbb", "run-0")
	for _, path := range []string{rootDir, childDir} {
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(rootDir, "root.jsonl"), "{\"type\":\"session\"}\n")
	writeFile(t, filepath.Join(childDir, "session.jsonl"), "{\"type\":\"session\"}\n")
	trace := piRecordedGateTraceFromArtifacts(t, dir)
	if len(trace.root) != 1 || len(trace.children) != 1 {
		t.Fatalf("artifact loader root=%d children=%d, want 1/1", len(trace.root), len(trace.children))
	}
}
