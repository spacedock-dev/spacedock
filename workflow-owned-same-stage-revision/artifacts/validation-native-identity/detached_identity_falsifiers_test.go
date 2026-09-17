package ensigncycle

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetachedCapturedIdentityFalsifiers(t *testing.T) {
	root := "/Users/clkao/git/spacedock-research/spacedock-v1/docs/dev/.spacedock-state/workflow-owned-same-stage-revision/artifacts/validation/tip35148797301"
	claude := readFile(t, filepath.Join(root, "claude/claude-stream.jsonl"))
	codex := readFile(t, filepath.Join(root, "codex/codex-native-lifecycle.jsonl"))
	self := strings.ReplaceAll(claude, "a6f8a511fd35edfb4", "a891653619c5fb829")
	routes, _ := claudeRejectionRoutes(self)
	if assertSameStageWorkers(routes, true) == nil {
		t.Fatal("captured self-review accepted")
	}
	var lines []string
	for _, line := range strings.Split(claude, "\n") {
		var x struct{ Subtype, TaskID string }
		var m map[string]any
		_ = x
		json.Unmarshal([]byte(line), &m)
		if m["subtype"] == "task_notification" && m["task_id"] == "a6f8a511fd35edfb4" {
			continue
		}
		lines = append(lines, line)
	}
	routes, _ = claudeRejectionRoutes(strings.Join(lines, "\n"))
	if assertSameStageWorkers(routes, true) == nil {
		t.Fatal("captured reviewer without completion accepted")
	}
	// An otherwise complete follow-up on a different, unspawned native owner cannot continue the correction worker.
	nativeLines := strings.Split(codex, "\n")
	for i := 151; i < len(nativeLines); i++ {
		nativeLines[i] = strings.ReplaceAll(nativeLines[i], "recorded_gate_task_validation", "unowned_validation")
	}
	if assertSameStageWorkers(codexRejectionRoutes(strings.Join(nativeLines, "\n")), false) == nil {
		t.Fatal("captured unowned reuse accepted")
	}
}
func TestDetachedRecorderEvidenceMatrix(t *testing.T) {
	success := "exit=0\tgate record rejection-task --round=validation/1 --briefing b --log l\n"
	for name, log := range map[string]string{"absent": "", "failed-masked": "exit=1\tgate record rejection-task --round validation/1\nexit=0\techo exit=0\n", "wrong-command": "exit=0\techo gate record rejection-task --round validation/1\n", "wrong-entity": "exit=0\tgate record another-task --round validation/1\n", "duplicate": success + success} {
		t.Run(name, func(t *testing.T) {
			if assertSingleRejectionRoundPublication(recorderRejectionRoundPublications(log)) == nil {
				t.Fatal("invalid actual-command evidence accepted")
			}
		})
	}
	if err := assertSingleRejectionRoundPublication(recorderRejectionRoundPublications("exit=1\tgate record rejection-task --round validation/1\n" + success)); err != nil {
		t.Fatal(err)
	}
}
