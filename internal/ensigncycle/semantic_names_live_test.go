//go:build live

package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

//spacedock:live-proof id=codex-semantic-worker-handles lane=codex-live
func TestLiveSemanticNamesCodex(t *testing.T) {
	runner := newCodexLiveRunner(t)
	root, state, ep, seed := writeSemanticNamesFixture(t)
	before := strings.TrimSpace(git(t, state, "rev-parse", "HEAD"))
	build := func(stage string, advance bool) string {
		args := []string{"dispatch", "build", "--host", "codex", "--workflow-dir", root, "--entity-path", ep, "--stage", stage, "--checklist-file", "-"}
		if advance {
			args = append(args, "--advance")
		}
		cmd := exec.Command(runner.binary, args...)
		cmd.Stdin = strings.NewReader("Append a Stage Report containing SEMANTIC-" + strings.ToUpper(stage) + ". Commit only this entity body; preserve frontmatter exactly. No push.\n")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("build: %v %s", err, out)
		}
		return string(out)
	}
	fresh, advance := build("ideation", false), build("validation", true)
	var envelope struct{ Name, Prompt string }
	if err := json.Unmarshal([]byte(fresh), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Name != "ci-duration-hints-ideation" {
		t.Fatal(envelope.Name)
	}
	prompt := fmt.Sprintf("This is one bounded native worker/reuse proof, not a workflow orchestration exercise. Spawn exactly one native Codex worker with task_name ci_duration_hints_ideation and this generated prompt verbatim: %s\nWait for its completion. Then send this advance envelope's prompt verbatim to the SAME returned worker handle using followup_task: %s\nWait for its second completion, then finish. Do not edit or commit fixture files yourself. Workers may only append their assigned marker Stage Report and commit %s. No frontmatter changes, extra workers, pushes, or workflow transitions. Root: %s", envelope.Prompt, advance, ep, root)
	result, err := runner.run(t, sharedRuntimeScenario{name: "semantic-names-codex"}, root, prompt)
	if err != nil {
		t.Fatalf("%v; artifacts %s", err, result.artifactDir)
	}
	stream, err := codexNativeLifecycleStream(runner.codexHome, result.jsonl, result.artifactDir)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(result.artifactDir, "native-lifecycle.jsonl"), stream)
	writeFile(t, filepath.Join(result.artifactDir, "entity.md"), readFile(t, ep))
	writeFile(t, filepath.Join(result.artifactDir, "state-log.txt"), git(t, state, "log", "--format=%H %s", "--name-only", before+"..HEAD"))
	if err := semanticLifecycle(stream); err != nil {
		t.Fatal(err)
	}
	wrongHandle := strings.Split(stream, "\n")
	for i, line := range wrongHandle {
		if strings.Contains(line, `"name":"followup_task"`) {
			wrongHandle[i] = strings.ReplaceAll(line, "/root/ci_duration_hints_ideation", "/root/other_worker")
		}
	}
	for _, mutant := range []string{strings.ReplaceAll(stream, "ci_duration_hints_ideation", "spacedock_ensign_ci_duration_hints_ideation"), strings.Join(wrongHandle, "\n")} {
		if err := semanticLifecycle(mutant); err == nil {
			t.Fatal("wrong native lifecycle accepted despite report markers")
		}
	}
	body := readFile(t, ep)
	if !strings.HasPrefix(body, seed) {
		t.Fatal("worker changed seed/frontmatter")
	}
	for _, marker := range []string{"SEMANTIC-IDEATION", "SEMANTIC-VALIDATION"} {
		if !strings.Contains(body, marker) {
			t.Fatal("missing marker", marker, body)
		}
	}
	if strings.TrimSpace(git(t, state, "rev-list", "--count", before+"..HEAD")) != "2" {
		t.Fatal("expected two durable report commits")
	}
	if strings.TrimSpace(git(t, state, "status", "--porcelain")) != "" {
		t.Fatal("dirty state")
	}
	for _, path := range strings.Fields(git(t, state, "diff", "--name-only", before, "HEAD")) {
		if path != "ci-duration-hints.md" {
			t.Fatal("unassigned mutation", path)
		}
	}
	t.Logf("native semantic spawn and same-handle advance; two clean report commits; artifacts: %s", result.artifactDir)
}

func semanticLifecycle(stream string) error {
	spawnID, handle, follow, spawns := "", "", "", 0
	for _, line := range strings.Split(stream, "\n") {
		var event struct {
			Type    string
			Payload struct{ Type, Name, Arguments, CallID, Output string }
		}
		// Native rollout fields use snake_case.
		var raw struct {
			Type    string
			Payload json.RawMessage
		}
		if json.Unmarshal([]byte(line), &raw) != nil || raw.Type != "response_item" {
			continue
		}
		json.Unmarshal([]byte(line), &event)
		var payload map[string]json.RawMessage
		json.Unmarshal(raw.Payload, &payload)
		var callID string
		json.Unmarshal(payload["call_id"], &callID)
		if event.Payload.Type == "function_call" {
			var args map[string]string
			json.Unmarshal([]byte(event.Payload.Arguments), &args)
			if event.Payload.Name == "spawn_agent" {
				spawns++
				spawnID = callID
				if args["task_name"] != "ci_duration_hints_ideation" {
					return fmt.Errorf("wrong semantic spawn: %q", args["task_name"])
				}
			}
			if event.Payload.Name == "followup_task" {
				follow = args["target"]
			}
		}
		if event.Payload.Type == "function_call_output" && callID == spawnID && spawnID != "" {
			var output map[string]any
			json.Unmarshal([]byte(event.Payload.Output), &output)
			if value, ok := output["task_name"].(string); ok {
				handle = value
			}
		}
	}
	if spawns != 1 || handle == "" || follow != handle {
		return fmt.Errorf("spawn/reuse mismatch: spawns=%d handle=%q follow=%q", spawns, handle, follow)
	}
	return nil
}

//spacedock:live-fixture id=semantic-names/split-root
func writeSemanticNamesFixture(t *testing.T) (root, state, ep, seed string) {
	root = t.TempDir()
	state = filepath.Join(root, "state")
	ep = filepath.Join(state, "ci-duration-hints.md")
	readme := "---\nentity-type: task\nid-style: slug\nstate: state\nstages:\n  states:\n    - name: ideation\n      initial: true\n    - name: validation\n    - name: done\n      terminal: true\n---\n\n### ideation\n\nAppend the assigned marker report and commit only the entity body.\n\n### validation\n\nAppend the assigned marker report and commit only the entity body.\n\n### done\n\nFinished.\n"
	writeFile(t, filepath.Join(root, "README.md"), readme)
	gitInit(t, root)
	seed = "---\ntitle: CI duration hints\nstatus: ideation\n---\n\nFixture.\n"
	writeFile(t, ep, seed)
	gitInit(t, state)
	return
}
