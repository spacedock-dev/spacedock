//go:build live

package ensigncycle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func piScenarioRunners() map[string]func(*testing.T, sharedRuntimeScenario) {
	return map[string]func(*testing.T, sharedRuntimeScenario){
		"multi-workflow-boot": runPiMultiWorkflowBootScenario,
	}
}

func TestLivePiMultiWorkflowBootScenario(t *testing.T) {
	var selected *sharedRuntimeScenario
	for _, scenario := range sharedRuntimeScenarios() {
		if scenario.name == "multi-workflow-boot" {
			copy := scenario
			selected = &copy
			break
		}
	}
	if selected == nil {
		t.Fatal("shared multi-workflow-boot scenario is not defined")
	}
	run := piScenarioRunners()[selected.name]
	if run == nil {
		t.Fatal("shared multi-workflow-boot scenario has no Pi live runner")
	}
	run(t, *selected)
}

func runPiMultiWorkflowBootScenario(t *testing.T, scenario sharedRuntimeScenario) {
	t.Helper()
	piBin := piBinaryOrSkip(t)
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := piSpacedockBinary(t, repo)

	piHome := t.TempDir()
	sessionDir := t.TempDir()
	cleanHome := t.TempDir()
	seedPiLocalAuth(t, piHome, os.Getenv("HOME"))
	projectRoot := t.TempDir()
	fixture := writeMultiWorkflowBootFixture(t, projectRoot)
	ledger := newTestInvocationLedger(t, binary)
	artifactDir := filepath.Join(piLiveArtifactDir(t, "pi-shared-scenarios"), scenario.name)
	if err := os.MkdirAll(filepath.Join(artifactDir, "sessions"), 0o755); err != nil {
		t.Fatal(err)
	}
	env := piLiveEnv(piHome, sessionDir, cleanHome, filepath.Dir(binary), piSubagentsRoot)
	env = ledger.instrumentEnv(env)

	runPiLiveCommand(t, artifactDir, projectRoot, env, piBin,
		"--mode", "json",
		"--print",
		"--session-dir", filepath.Join(artifactDir, "sessions"),
		"--skill", filepath.Join(repo, "skills", "first-officer"),
		multiWorkflowBootPrompt(projectRoot),
	)

	finalMessage, err := parsePiFinalMessage(readFile(t, filepath.Join(artifactDir, "pi-stdout.txt")))
	if err != nil {
		t.Fatalf("parse Pi JSON trace: %v; artifacts in %s", err, artifactDir)
	}
	invocations := ledger.read(t)
	writeInvocationLedgerArtifact(t, artifactDir, invocations)
	obs := gatherMultiWorkflowBootObservation(t, fixture, invocations, finalMessage)
	if err := assertMultiWorkflowBoot(obs); err != nil {
		t.Fatalf("Pi multi-workflow boot scenario graded FAIL: %v; artifacts in %s\nFinal message:\n%s", err, artifactDir, finalMessage)
	}
}

func parsePiFinalMessage(jsonl string) (string, error) {
	finalMessage := ""
	for _, line := range strings.Split(jsonl, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event struct {
			Type    string `json:"type"`
			Message struct {
				Role    string `json:"role"`
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return "", err
		}
		if event.Type == "message_end" && event.Message.Role == "assistant" {
			var text strings.Builder
			for _, block := range event.Message.Content {
				if block.Type == "text" {
					text.WriteString(block.Text)
				}
			}
			if text.Len() > 0 {
				finalMessage = text.String()
			}
		}
	}
	if finalMessage == "" {
		return "", fmt.Errorf("Pi trace contains no terminal assistant message")
	}
	return finalMessage, nil
}
