//go:build live

package ensigncycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLivePiNonSelfDescribingDispatch (AC-1, AC-3) is the tautology-closing
// live lane: it dispatches a Pi worker through `dispatch build --host pi` with
// a checklist equal to a real entity's acceptance criteria — no "First read
// ensign/SKILL.md", no stage-report heading, no DONE/Summary structure — and
// asserts the worker still writes a complete `## Stage Report: implementation`
// (heading + `- DONE:` + `### Summary`) with a clean state-checkout commit.
// The worker's only format source is the embedded `### Stage Report format`
// block the build artifact now carries for host=pi. Reverting the AC-2 body
// embed makes this lane RED (the worker has no format source), while the
// self-describing TestLivePiFrontDoorSmoke stays green — proving this lane
// tests the real mode, not a fixture hint.
//
//spacedock:live-proof id=pi-non-self-describing-dispatch lane=pi-live
func TestLivePiNonSelfDescribingDispatch(t *testing.T) {
	repo := repoRoot(t)
	piSubagentsRoot := piSubagentsPackageRoot(t)
	binary := piSpacedockBinary(t, repo)
	workflowRoot, stateRoot, entityPath, artifactDir, env, model := newPiNonSelfDescribingSmokeFixture(t, "pi-nonsd-smoke", repo, piSubagentsRoot, binary)

	envelope := runPiNonSelfDescribingDispatchBuild(t, binary, workflowRoot, entityPath)
	prompt := piNonSelfDescribingSmokePrompt(repo, workflowRoot, stateRoot, entityPath, envelope)
	runPiLiveCommand(t, artifactDir, workflowRoot, env, binary,
		"pi",
		prompt,
		"--plugin-dir", repo,
		"--",
		"--print",
		"--model", model,
		"--session-dir", filepath.Join(artifactDir, "sessions"),
	)
	assertPiNonSelfDescribingSmokeResult(t, stateRoot, entityPath, artifactDir)
}

func newPiNonSelfDescribingSmokeFixture(t *testing.T, name, repo, piSubagentsRoot, binary string) (workflowRoot, stateRoot, entityPath, artifactDir string, env []string, model string) {
	t.Helper()
	piHome := t.TempDir()
	sessionDir := t.TempDir()
	cleanHome := t.TempDir()
	decision := seedPiLiveAuth(t, piHome, os.Getenv("HOME"), os.Getenv("CODEX_AUTH_JSON"), os.Getenv("OPENAI_API_KEY"), os.Getenv("SPACEDOCK_PI_LIVE_REQUIRED"))
	writeFile(t, filepath.Join(piHome, "settings.json"), fmt.Sprintf("{\"packages\":[%q]}\n", "file:"+repo))
	writePiSubagentsProjectArtifactDir(t, piHome)
	workflowRoot, stateRoot, entityPath = writePiNonSelfDescribingSmokeWorkflow(t)
	artifactDir = filepath.Join(piLiveArtifactDir(t, name), "run")
	if err := os.MkdirAll(filepath.Join(artifactDir, "sessions"), 0o755); err != nil {
		t.Fatal(err)
	}
	env = piLiveEnvForAuth(piHome, sessionDir, cleanHome, filepath.Dir(binary), piSubagentsRoot, os.Getenv("OPENAI_API_KEY"), decision.mode)
	model = piLiveChildModel(decision)
	return workflowRoot, stateRoot, entityPath, artifactDir, env, model
}

// runPiNonSelfDescribingDispatchBuild assembles the initial-dispatch artifact
// for the non-self-describing smoke entity with a checklist equal to a real
// entity's acceptance criteria: no ensign skill path, no stage-report heading,
// no DONE/Summary structure. The worker's only stage-report format source is
// the embedded `### Stage Report format` block the build artifact carries for
// host=pi (AC-2).
func runPiNonSelfDescribingDispatchBuild(t *testing.T, binary, workflowRoot, entityPath string) piSmokeEnvelope {
	t.Helper()
	// A real-entity acceptance-criteria checklist: the work to do and the
	// commit discipline, with zero format or skill-path hints.
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
	cmd := exec.Command(binary, "dispatch", "build", "--workflow-dir", workflowRoot)
	cmd.Dir = workflowRoot
	cmd.Stdin = strings.NewReader(string(stdin))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("dispatch build --host pi failed: %v\nstderr:\n%s", err, stderr.String())
	}
	out := stdout.Bytes()
	var envelope piSmokeEnvelope
	if err := json.Unmarshal(out, &envelope); err != nil {
		t.Fatalf("dispatch build stdout is not the build envelope: %v\n%s\nstderr:\n%s", err, out, stderr.String())
	}
	if envelope.Agent != "worker" || envelope.Skill != "ensign" {
		t.Fatalf("pi build envelope = agent %q skill %q, want worker/ensign:\n%s", envelope.Agent, envelope.Skill, out)
	}
	if envelope.Prompt == "" || envelope.DispatchFile == "" {
		t.Fatalf("pi build envelope missing prompt/dispatch_file_path:\n%s", out)
	}
	// Adversarial guard: confirm the dispatch body carries the embedded
	// stage-report format block (AC-2) so the worker has a format source.
	body, err := os.ReadFile(envelope.DispatchFile)
	if err != nil {
		t.Fatalf("read dispatch artifact: %v", err)
	}
	bodyStr := string(body)
	for _, want := range []string{"### Stage Report format", "## Stage Report:", "- DONE:", "### Summary"} {
		if !strings.Contains(bodyStr, want) {
			t.Fatalf("non-self-describing dispatch body missing embedded protocol token %q:\n%s", want, bodyStr)
		}
	}
	return envelope
}

// piNonSelfDescribingSmokePrompt is the FO prompt for the non-self-describing
// lane. It forwards the dispatch artifact's spawn fields verbatim and, after
// the worker returns, verifies the entity carries a complete
// `## Stage Report: implementation` and the state git log has the worker
// commit. Unlike piLiveSmokePrompt it does NOT tell the FO to verify the
// smoke marker — the marker is the worker's real work, not the proof; the
// stage report (sourced from the embedded body block) is the proof.
func piNonSelfDescribingSmokePrompt(repo, workflowRoot, stateRoot, entityPath string, envelope piSmokeEnvelope) string {
	return fmt.Sprintf(`You are the Spacedock first officer for a live Pi smoke test.

An initial-dispatch artifact was assembled for the entity with `+"`spacedock dispatch build --host pi`"+`; forward it through pi-subagents exactly as emitted — this smoke exists to prove the build artifact's embedded stage-report format drives the worker's report even when the checklist does not name the format.

  agent: %[5]s
  skill: %[6]s
  task: %[7]s

Use the pi-subagents subagent(...) tool exactly once with those fields verbatim (context must be "fresh", working directory %[2]s). Do not use or mention Claude Agent, SendMessage, TeamCreate, or TeamDelete tools. Do not paraphrase, re-order, or extend the task string.

After subagent(...) returns, you as first officer must verify the entity file %[4]s contains a '## Stage Report: implementation' section with at least one '- DONE:' item and a '### Summary' subsection, and verify the state checkout %[3]s git log contains 'ensign: pi live smoke' over pi-nonsd-smoke/index.md. Exit successfully only after those durable checks pass; your final message names the agent and skill values you passed to subagent(...) and the child's run id.

Reference paths: ensign contract at %[1]s/skills/ensign/SKILL.md; Pi ensign adapter at %[1]s/skills/ensign/references/pi-ensign-runtime.md (the worker's dispatch artifact already points at them).`,
		repo, workflowRoot, stateRoot, entityPath, envelope.Agent, envelope.Skill, envelope.Prompt)
}

func assertPiNonSelfDescribingSmokeResult(t *testing.T, stateRoot, entityPath, artifactDir string) {
	t.Helper()
	entity := readFile(t, entityPath)
	// The complete stage report structure (heading + DONE + Summary) plus the
	// durable git commit prove the spawned worker followed the embedded
	// stage-report format block without the checklist naming it.
	for _, want := range []string{"## Stage Report: implementation", "- DONE:", "### Summary"} {
		if !strings.Contains(entity, want) {
			t.Fatalf("entity missing %q after non-self-describing pi subagent smoke; artifacts in %s\n%s", want, artifactDir, entity)
		}
	}
	log := git(t, stateRoot, "log", "--oneline", "--", "pi-nonsd-smoke", "index.md")
	if !strings.Contains(log, "ensign: pi live smoke") {
		t.Fatalf("state checkout git log missing worker commit; artifacts in %s\n%s", artifactDir, log)
	}
	if strings.TrimSpace(git(t, stateRoot, "status", "--short", "--", "pi-nonsd-smoke", "index.md")) != "" {
		t.Fatalf("state checkout entity has uncommitted changes after worker commit; artifacts in %s\n%s", artifactDir, git(t, stateRoot, "status", "--short"))
	}
}
