package ensigncycle

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestPiNonSelfDescribingDispatchBuildBodyCarriesProtocol is the offline
// (non-live) guard for the non-self-describing lane (AC-2 body presence +
// AC-3 tautology closure): it builds the dispatch artifact with a checklist
// equal to a real entity's acceptance criteria — no ensign skill path, no
// stage-report heading, no DONE/Summary structure — and asserts the body
// carries the embedded stage-report protocol tokens while the checklist
// stdin does not smuggle any format hint. This runs without the `live` tag
// so the AC-2 body presence and the AC-3 tautology-closure (checklist has no
// format hint) are checked on every test run, not only live dispatches.
func TestPiNonSelfDescribingDispatchBuildBodyCarriesProtocol(t *testing.T) {
	binary := buildRecordedGateBinary(t)
	workflowRoot, stateRoot, entityPath := writePiNonSelfDescribingSmokeWorkflow(t)
	_ = stateRoot
	checklist := []string{
		"- append the smoke marker line `PI-NONSD-SMOKE-MARKER` to the entity file",
		"- commit only the entity path in the state checkout with message 'ensign: pi live smoke' (path-scoped git add/commit for pi-nonsd-smoke/index.md)",
	}
	stdin, err := json.Marshal(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   workflowRoot,
		"stage":          "implementation",
		"checklist":      checklist,
		"bare_mode":      true,
		"host":           "pi",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, "dispatch", "build", "--workflow-dir", workflowRoot)
	cmd.Dir = workflowRoot
	cmd.Stdin = strings.NewReader(string(stdin))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("dispatch build --host pi failed: %v\nstderr:\n%s", err, stderr.String())
	}
	var envelope piSmokeEnvelope
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatalf("dispatch build stdout is not the build envelope: %v\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	body, err := os.ReadFile(envelope.DispatchFile)
	if err != nil {
		t.Fatalf("read dispatch artifact: %v", err)
	}
	bodyStr := string(body)
	// AC-2: the body carries the embedded protocol tokens.
	for _, want := range []string{"### Stage Report format", "## Stage Report:", "- DONE:", "- SKIPPED:", "- FAILED:", "### Summary"} {
		if !strings.Contains(bodyStr, want) {
			t.Fatalf("non-self-describing dispatch body missing embedded protocol token %q:\n%s", want, bodyStr)
		}
	}
	// AC-3 tautology closure: the checklist (stdin) must not name the ensign
	// skill path, the stage-report heading, or the DONE/Summary structure —
	// the body embed is the worker's only format source.
	for _, banned := range []string{"ensign/SKILL.md", "## Stage Report:", "- DONE:", "- SKIPPED:", "- FAILED:", "### Summary", "Stage Report format"} {
		if strings.Contains(string(stdin), banned) {
			t.Fatalf("non-self-describing checklist smuggles format hint %q into the dispatch stdin:\n%s", banned, stdin)
		}
	}
}

// writePiNonSelfDescribingSmokeWorkflow creates a split-root smoke workflow
// whose implementation stage-def names only the real work (append a marker
// line) — no "stage report" mention — so the worker's stage-report format
// source is the embedded dispatch body block, not the stage-def.
//
//spacedock:live-fixture id=pi/non-self-describing-smoke
func writePiNonSelfDescribingSmokeWorkflow(t *testing.T) (workflowRoot, stateRoot, entityPath string) {
	t.Helper()
	workflowRoot = t.TempDir()
	stateRoot = filepath.Join(workflowRoot, ".spacedock-state")
	writeFile(t, filepath.Join(workflowRoot, "README.md"), piNonSelfDescribingSmokeReadme())
	entityPath = filepath.Join(stateRoot, "pi-nonsd-smoke", "index.md")
	writeFile(t, entityPath, piNonSelfDescribingSmokeEntity())
	gitInit(t, workflowRoot)
	gitInit(t, stateRoot)
	return workflowRoot, stateRoot, entityPath
}

func piNonSelfDescribingSmokeReadme() string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"state: .spacedock-state\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: implementation\n" +
		"      initial: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Pi Non-Self-Describing Smoke\n\n" +
		"### implementation\n\n" +
		"Append the live Pi smoke marker line `PI-NONSD-SMOKE-MARKER` as a standalone line to the entity file, then commit only the entity path in the state checkout.\n\n" +
		"- **Outputs:** The marker line present in the entity file and a path-scoped state commit.\n\n" +
		"### done\n\nTerminal state.\n"
}

func piNonSelfDescribingSmokeEntity() string {
	return "---\n" +
		"id: pi-nonsd-smoke\n" +
		"title: Pi Non-Self-Describing Smoke\n" +
		"status: implementation\n" +
		"completed:\n" +
		"verdict:\n" +
		"worktree:\n" +
		"---\n" +
		"# Pi Non-Self-Describing Smoke\n\n" +
		"This entity is mutated only by the Pi subagent non-self-describing live smoke.\n"
}
