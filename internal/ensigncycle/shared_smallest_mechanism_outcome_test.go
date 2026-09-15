package ensigncycle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSmallestMechanismOutcomes(t *testing.T) {
	for _, mutation := range []string{"", "missing alpha", "wrong alpha", "wrong beta", "wrong strategy", "uncommitted strategy", "codex worker", "claude worker", "pi worker", "structured edit without result", "final drift"} {
		t.Run(mutation, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, filepath.Join(root, "seed"), "known initial state\n")
			gitInit(t, root)
			expected := map[string]string{
				ssmEditFileA:   "# Ladder Note Alpha\n\nStatus: RESOLVED\n",
				ssmEditFileB:   "# Ladder Note Beta\n\nStatus: RESOLVED\n",
				ssmStrategyDoc: "# Roadmap Strategy\n",
			}
			for f, content := range expected {
				writeFile(t, filepath.Join(root, f), content)
			}
			if mutation != "uncommitted strategy" {
				git(t, root, "add", ssmStrategyDoc)
				git(t, root, "commit", "-m", "Add strategy")
			}
			native := ""
			switch mutation {
			case "missing alpha":
				if err := os.Remove(filepath.Join(root, ssmEditFileA)); err != nil {
					t.Fatal(err)
				}
			case "wrong alpha", "structured edit without result":
				writeFile(t, filepath.Join(root, ssmEditFileA), ladderNote("Ladder Note Alpha"))
				native = codexFileChange(ssmEditFileA)
			case "wrong beta":
				writeFile(t, filepath.Join(root, ssmEditFileB), "Status: RESOLVED\n")
			case "wrong strategy":
				writeFile(t, filepath.Join(root, ssmStrategyDoc), "wrong\n")
			case "codex worker":
				native = `{"payload":{"type":"function_call","name":"spawn_agent","arguments":"{\"message\":\"help\"}"}}`
			case "claude worker":
				native = claudeToolUse("Agent", `{"prompt":"help"}`)
			case "pi worker":
				native = `{"message":{"content":[{"type":"toolCall","name":"subagent","arguments":{"task":"help"}}]}}`
			}
			err := assertSmallestMechanismDirect(root, native)
			if mutation == "" || mutation == "final drift" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil {
				t.Fatal("incorrect or delegated direct result passed")
			}
			if mutation == "final drift" {
				writeFile(t, filepath.Join(root, ssmEditFileB), "worker changed this\n")
				if assertSmallestMechanismFiles(root) == nil {
					t.Fatal("phase-2 final drift passed")
				}
			}
		})
	}
}

func TestSmallestMechanismDirectDoesNotDependOnTools(t *testing.T) {
	// Read-only mentions of agent APIs and edit tools are data, never worker calls.
	stream := codexCommand("cat file # spawn_agent apply_patch python3") + "\n" + codexAgentMessage("I will check spawn_agent documentation")
	if smallestMechanismHasWorker(stream) {
		t.Fatal("command/prose counted as native delegation")
	}
	if !smallestMechanismHasWorker(strings.Replace(claudeToolUse("Task", `{"prompt":"help"}`), "help", "any task", 1)) {
		t.Fatal("legacy Task worker missed")
	}
}

func TestSmallestMechanismNativeOutputIsNotSpawn(t *testing.T) {
	for host, stream := range map[string]string{
		"claude": `{"type":"user","message":{"content":[{"type":"tool_result","name":"Agent","content":"spawn_agent"},{"type":"text","text":"Task Agent"}]}}`,
		"pi":     `{"message":{"role":"toolResult","toolName":"subagent","content":[{"type":"text","text":"subagent"}]}}`,
		"codex":  `{"payload":{"type":"function_call_output","name":"spawn_agent","output":"spawn_agent"}}`,
	} {
		t.Run(host, func(t *testing.T) {
			if smallestMechanismHasWorker(stream) {
				t.Fatal("output counted as worker spawn")
			}
		})
	}
}
