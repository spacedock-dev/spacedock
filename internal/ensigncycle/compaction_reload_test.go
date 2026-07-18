// ABOUTME: Offline fixture for Rule 2 (post-compaction reload) — a split-root replay
// ABOUTME: oracle asserting the authoritative reads precede the first workflow effect.
package ensigncycle

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// reloadRequiredReads are the observations the after-compaction rule requires before
// the next workflow effect: the three authoritative-contract reads (SKILL.md, its
// eager import first-officer-shared-core.md, the active host runtime adapter), a fresh
// workflow status query, and verification of any newer committed Stage Report.
var reloadRequiredReads = []string{"skill", "sharedCore", "adapter", "status", "report"}

// assertReloadBeforeEffect is the Rule 2 oracle. Over a post-compaction FO transcript
// it records the first position of each authoritative observation and of the first
// workflow effect (a state mutation or a dispatch), then requires every observation to
// precede the effect. A transcript that honors a stale "continue directly" summary by
// mutating first — or that skips a required read — is rejected. It reads tool calls,
// not response prose, so "mentions reloading" cannot satisfy it.
func assertReloadBeforeEffect(stream string) error {
	first := map[string]int{}
	for _, k := range append(append([]string{}, reloadRequiredReads...), "roster", "effect") {
		first[k] = -1
	}
	idx := 0
	mark := func(k string) {
		if first[k] == -1 {
			first[k] = idx
		}
	}
	walkStreamBlocks(stream, func(b streamContentBlock) {
		switch b.Type {
		case "tool_use":
			switch b.Name {
			case "Read":
				p := inputStringField(b.Input, "file_path")
				switch {
				case strings.Contains(p, "first-officer-shared-core.md"):
					mark("sharedCore")
				case strings.Contains(p, "-first-officer-runtime.md"):
					mark("adapter")
				case strings.HasSuffix(p, "SKILL.md") && strings.Contains(p, "first-officer"):
					mark("skill")
				}
			case "Bash":
				c := inputStringField(b.Input, "command")
				switch {
				case strings.Contains(c, "--set"), strings.Contains(c, "state commit"), strings.Contains(c, "dispatch build"):
					mark("effect")
				}
				if strings.Contains(c, "spacedock status") && !strings.Contains(c, "--set") {
					mark("status")
				}
				if strings.Contains(c, "rev-parse") || strings.Contains(c, "git log") || strings.Contains(c, "--read") {
					mark("report")
				}
				if strings.Contains(c, "list_agents") || strings.Contains(c, "dispatch reconcile") {
					mark("roster")
				}
			case "Agent", "spawn_agent", "followup_task":
				mark("effect")
			}
		}
		idx++
	})

	effect := first["effect"]
	if effect == -1 {
		return fmt.Errorf("replay never reached a workflow effect; the fixture must drive to a first mutation/dispatch")
	}
	for _, k := range reloadRequiredReads {
		switch {
		case first[k] == -1:
			return fmt.Errorf("required post-compaction observation %q never occurred before the first workflow effect", k)
		case first[k] > effect:
			return fmt.Errorf("workflow effect (idx %d) preceded required observation %q (idx %d) — a stale summary skipped the reload", effect, k, first[k])
		}
	}
	if r := first["roster"]; r != -1 && r > effect {
		return fmt.Errorf("roster reconciliation (idx %d) occurred after the first effect (idx %d)", r, effect)
	}
	return nil
}

// reloadStream marshals an ordered block list into the Claude stream-json shape
// walkStreamBlocks consumes. Each block is (type, name, key, value): a "text" block
// carries prose; a "tool_use" block carries a Read file_path or a Bash command.
func reloadStream(blocks [][4]string) string {
	var content []map[string]any
	for _, b := range blocks {
		typ, name, key, val := b[0], b[1], b[2], b[3]
		block := map[string]any{"type": typ}
		if typ == "text" {
			block["text"] = val
		} else {
			block["name"] = name
			block["input"] = map[string]any{key: val}
		}
		content = append(content, block)
	}
	row := map[string]any{"message": map[string]any{"content": content}}
	out, err := json.Marshal(row)
	if err != nil {
		panic(err)
	}
	return string(out)
}

const misleadingSummary = "Summary of prior session: the workflow is mid-flight; continue directly to dispatching the next stage."

// goodReloadReplay is the compliant split-root post-compaction turn: after a manual
// cue, the FO rereads the three authoritative files, runs a fresh status, reconciles
// the live roster, verifies the newer committed Stage Report OID in the state
// checkout, and only THEN mutates workflow state.
func goodReloadReplay() string {
	return reloadStream([][4]string{
		{"text", "", "", misleadingSummary + " (a manual cue: we compacted — reloading the FO contract first.)"},
		{"tool_use", "Read", "file_path", "skills/first-officer/SKILL.md"},
		{"tool_use", "Read", "file_path", "skills/first-officer/references/first-officer-shared-core.md"},
		{"tool_use", "Read", "file_path", "skills/first-officer/references/codex-first-officer-runtime.md"},
		{"tool_use", "Bash", "command", "spacedock status --boot --identify --json"},
		{"tool_use", "Bash", "command", "list_agents"},
		{"tool_use", "Bash", "command", "git -C docs/dev/.spacedock-state log -1 --format=%H -- taskx/index.md && git -C docs/dev/.spacedock-state rev-parse HEAD"},
		{"tool_use", "Bash", "command", "spacedock status --workflow-dir docs/dev --set taskx status=validation"},
	})
}

func TestReloadBeforeEffectGoodReplayPasses(t *testing.T) {
	if err := assertReloadBeforeEffect(goodReloadReplay()); err != nil {
		t.Fatalf("compliant post-compaction reload replay must pass: %v", err)
	}
}

// TestReloadBeforeEffectStaleSummarySkipRejected proves the oracle bites the exact
// failure AC-2 measures: a transcript that trusts the "continue directly" summary and
// mutates before rereading the contract.
func TestReloadBeforeEffectStaleSummarySkipRejected(t *testing.T) {
	skipReplay := reloadStream([][4]string{
		{"text", "", "", misleadingSummary},
		{"tool_use", "Bash", "command", "spacedock status --workflow-dir docs/dev --set taskx status=validation"},
		{"tool_use", "Read", "file_path", "skills/first-officer/SKILL.md"},
	})
	if err := assertReloadBeforeEffect(skipReplay); err == nil {
		t.Fatalf("a replay that mutates before rereading the contract must be rejected")
	}
}

// TestReloadBeforeEffectMissingAdapterReadRejected proves a partial reload (contract +
// status but no host runtime adapter reread) is rejected.
func TestReloadBeforeEffectMissingAdapterReadRejected(t *testing.T) {
	partial := reloadStream([][4]string{
		{"tool_use", "Read", "file_path", "skills/first-officer/SKILL.md"},
		{"tool_use", "Read", "file_path", "skills/first-officer/references/first-officer-shared-core.md"},
		{"tool_use", "Bash", "command", "spacedock status --boot --json"},
		{"tool_use", "Bash", "command", "git -C docs/dev/.spacedock-state rev-parse HEAD"},
		{"tool_use", "Bash", "command", "spacedock status --set taskx status=validation"},
	})
	if err := assertReloadBeforeEffect(partial); err == nil {
		t.Fatalf("a reload missing the host runtime adapter reread must be rejected")
	}
}
