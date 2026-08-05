//go:build live

package ensigncycle

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestLiveCodexWaitMatrixFromShippedAdapter(t *testing.T) {
	runner := newCodexLiveRunner(t)
	head := strings.TrimSpace(git(t, repoRoot(t), "rev-parse", "HEAD"))
	if got := strings.TrimSpace(readFile(t, filepath.Join(runner.artifactRoot, "_setup", "source-head.txt"))); got != head {
		t.Fatalf("Codex adapter source head = %q, want exact candidate %q", got, head)
	}
	for _, tc := range []struct {
		state    string
		wantWait bool
	}{{"active unresolved", true}, {"completed", false}, {"errored", false}, {"absent", false}} {
		t.Run(tc.state, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, filepath.Join(root, "README.md"), codexWaitMatrixWorkflow)
			gitInit(t, root)
			scenario := sharedRuntimeScenario{name: "wait-matrix-" + strings.ReplaceAll(tc.state, " ", "-")}
			result, err := runner.run(t, scenario, root, codexWaitMatrixPrompt(root, tc.state))
			if err != nil {
				t.Fatalf("%v\nArtifacts: %s", err, result.artifactDir)
			}
			calls, exact := observedCodexWaitCalls(t, result.jsonl, runner.codexHome)
			if tc.wantWait && (calls != 1 || exact != 1) {
				t.Fatalf("%s: wait calls/exact-300000 = %d/%d, want 1/1; artifacts: %s", tc.state, calls, exact, result.artifactDir)
			}
			if !tc.wantWait && calls != 0 {
				t.Fatalf("%s emitted %d wait_agent call(s), want none; artifacts: %s", tc.state, calls, result.artifactDir)
			}
		})
	}
}

const codexWaitMatrixWorkflow = `---
entity-type: task
id-style: slug
stages: {states: [{name: backlog, initial: true}, {name: done, terminal: true}]}
---`

func codexWaitMatrixPrompt(root, state string) string {
	setup := map[bool]string{true: "First spawn exactly one probe worker whose whole assignment is to run `sleep 20` and then return; its live task is the supplied unresolved worker.", false: "Do not spawn a worker."}[state == "active unresolved"]
	return fmt.Sprintf(`Use $spacedock:first-officer and the current-checkout Codex runtime adapter. This is a bounded conformance exercise in %s. The loop has completed its first empty status query, one idle hook, roster reconcile, and second empty status query, with no dispatchable, gate, mod/PR, or other state work. %s Treat this harness observation as authoritative: worker set is %s. Apply the shipped adapter's wait decision once. If it requires an async wait, make the one adapter-prescribed call and exit after it returns. Otherwise do not call wait_agent and exit. Do not mutate or replace the supplied observation with another roster query.`, root, setup, state)
}

func observedCodexWaitCalls(t *testing.T, execJSONL, codexHome string) (calls, exact int) {
	t.Helper()
	marker := `"thread_id":"`
	start := strings.Index(execJSONL, marker)
	if start < 0 {
		t.Fatal("Codex exec JSONL has no thread id")
	}
	rest := execJSONL[start+len(marker):]
	threadID := rest[:strings.Index(rest, `"`)]
	paths, _ := filepath.Glob(filepath.Join(codexHome, "sessions", "*", "*", "*", "rollout-*"+threadID+".jsonl"))
	if len(paths) != 1 {
		t.Fatalf("session rollout for thread %q = %v, want one", threadID, paths)
	}
	for _, line := range strings.Split(readFile(t, paths[0]), "\n") {
		if strings.Contains(line, `"type":"response_item"`) && strings.Contains(line, `"type":"function_call"`) && strings.Contains(line, `"name":"wait_agent"`) {
			calls++
			if strings.Contains(line, `"arguments":"{\"timeout_ms\":300000}"`) {
				exact++
			}
		}
	}
	return calls, exact
}
