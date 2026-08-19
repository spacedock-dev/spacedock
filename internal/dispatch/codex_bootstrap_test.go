// ABOUTME: Codex fresh-dispatch boundary fixture — a supplied child probe loads
// ABOUTME: the real ensign contract before consuming the generated artifact.
package dispatch

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type codexBootstrapFixture struct {
	root         string
	entityPath   string
	dispatchPath string
	promptPath   string
	prompt       string
	artifact     string
}

const codexBootstrapProbeScript = "internal/dispatch/testdata/codex-bootstrap-probe.sh"

func TestCodexFreshDispatchBoundaryFixture(t *testing.T) {
	repo := dispatchTestRepoRoot(t)
	pluginRoot := copyCodexPluginRoot(t, repo)
	spacedock := buildCodexBootstrapBinary(t, repo)

	t.Run("bootstrap loads contract before artifact and report parser accepts", func(t *testing.T) {
		fixture := newCodexBootstrapFixture(t)
		output, code := runCodexBootstrapProbe(t, repo, pluginRoot, fixture, "done")
		if code != 0 {
			t.Fatalf("fresh child probe exit=%d, want 0:\n%s", code, output)
		}
		for _, marker := range []string{
			"bootstrap edge accepted",
			"contract-read:skills/ensign/SKILL.md",
			"contract-read:skills/ensign/references/ensign-shared-core.md",
			"contract-read:skills/ensign/references/codex-ensign-runtime.md",
			"artifact-read",
			"report-written:done",
		} {
			if !strings.Contains(output, marker) {
				t.Fatalf("fresh child output missing %q:\n%s", marker, output)
			}
		}
		artifactRead := strings.Index(output, "artifact-read")
		lastContractRead := strings.LastIndex(output, "contract-read:")
		if artifactRead <= lastContractRead {
			t.Fatalf("child read artifact before all contract files: contract=%d artifact=%d\n%s", lastContractRead, artifactRead, output)
		}

		body := readCodexBootstrapEntity(t, fixture.entityPath)
		if !strings.Contains(body, "\n## Stage Report: implementation\n") {
			t.Fatalf("child did not emit anchored implementation report:\n%s", body)
		}
		if !strings.Contains(body, "\n- DONE:") || !strings.Contains(body, "\n### Summary\n") {
			t.Fatalf("child report is missing DONE accounting or Summary:\n%s", body)
		}
		if strings.Contains(body, "- [x]") {
			t.Fatalf("positive child report contains forbidden checkbox shape:\n%s", body)
		}

		statusOut, statusCode := runCodexBootstrapStatus(t, spacedock, fixture)
		if statusCode != 0 {
			t.Fatalf("existing status report parser rejected compliant child report, exit=%d:\n%s", statusCode, statusOut)
		}
		if !strings.Contains(statusOut, `"status":"DONE"`) {
			t.Fatalf("status parser did not return the DONE item:\n%s", statusOut)
		}
	})

	t.Run("removing bootstrap edge fails before contract or artifact reads", func(t *testing.T) {
		fixture := newCodexBootstrapFixture(t)
		writeFile(t, fixture.promptPath, strings.TrimPrefix(fixture.prompt, "$spacedock:ensign; then "))
		oldFalseSentence := "This file contains the shared ensign discipline entry points plus the stage-specific assignment; this file pointer is the contract surface."
		writeFile(t, fixture.dispatchPath, oldFalseSentence+"\n\n"+fixture.artifact)

		output, code := runCodexBootstrapProbe(t, repo, pluginRoot, fixture, "done")
		if code == 0 {
			t.Fatalf("bootstrap-removal mutation unexpectedly passed:\n%s", output)
		}
		if !strings.Contains(output, "bootstrap edge missing") {
			t.Fatalf("bootstrap-removal mutation failed for the wrong reason:\n%s", output)
		}
		for _, marker := range []string{"contract-read:", "artifact-read", "report-written:"} {
			if strings.Contains(output, marker) {
				t.Fatalf("bootstrap-removal mutation read or wrote after the missing edge (%q):\n%s", marker, output)
			}
		}
	})

	t.Run("generic checkbox report fails the existing parser independently", func(t *testing.T) {
		fixture := newCodexBootstrapFixture(t)
		output, code := runCodexBootstrapProbe(t, repo, pluginRoot, fixture, "checkbox")
		if code != 0 {
			t.Fatalf("checkbox child probe exit=%d, want 0 so the report mutation reaches the parser:\n%s", code, output)
		}
		if !strings.Contains(output, "artifact-read") || !strings.Contains(output, "report-written:checkbox") {
			t.Fatalf("checkbox child did not reach report emission:\n%s", output)
		}

		body := readCodexBootstrapEntity(t, fixture.entityPath)
		if strings.Contains(body, "## Stage Report: implementation") || !strings.Contains(body, "\n## Stage Report\n") {
			t.Fatalf("checkbox mutation unexpectedly emitted an anchored stage report:\n%s", body)
		}
		if !strings.Contains(body, "- [x]") || !strings.Contains(body, "### Summary") {
			t.Fatalf("checkbox mutation did not preserve its independent invalid shape:\n%s", body)
		}

		statusOut, statusCode := runCodexBootstrapStatus(t, spacedock, fixture)
		if statusCode == 0 {
			t.Fatalf("existing status report parser accepted generic checkbox report:\n%s", statusOut)
		}
		if !strings.Contains(statusOut, "no ## Stage Report") {
			t.Fatalf("generic checkbox rejection did not come from the anchored-heading parser:\n%s", statusOut)
		}
	})
}

func newCodexBootstrapFixture(t *testing.T) codexBootstrapFixture {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeNonWorktreeStages())
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", ""))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      []string{"- run the Codex bootstrap probe"},
		"bare_mode":      false,
		"host":           "codex",
	}, nil)
	native := runNativePreservingHostEnv(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("fixture dispatch build exit=%d stderr=%q", native.exit, native.stderr)
	}
	out := decodeBuildOutput(t, native.stdout)
	assertCodexFreshPrompt(t, out.Prompt, out.DispatchFilePath)
	body := readDispatchBody(t, out.DispatchFilePath)
	promptPath := filepath.Join(root, "fresh-prompt.txt")
	writeFile(t, promptPath, out.Prompt)

	return codexBootstrapFixture{
		root:         root,
		entityPath:   entityPath,
		dispatchPath: out.DispatchFilePath,
		promptPath:   promptPath,
		prompt:       out.Prompt,
		artifact:     body,
	}
}

func runCodexBootstrapProbe(t *testing.T, repo, pluginRoot string, fixture codexBootstrapFixture, reportShape string) (string, int) {
	t.Helper()
	script := filepath.Join(repo, codexBootstrapProbeScript)
	cmd := exec.Command("sh", script, fixture.promptPath, fixture.dispatchPath, pluginRoot, fixture.entityPath, "implementation", reportShape)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return string(output), 0
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("run Codex bootstrap probe: %v\n%s", err, output)
	}
	return string(output), exitErr.ExitCode()
}

func runCodexBootstrapStatus(t *testing.T, spacedock string, fixture codexBootstrapFixture) (string, int) {
	t.Helper()
	cmd := exec.Command(spacedock, "status", "--workflow-dir", fixture.root, "--read", fixture.entityPath, "--stage", "implementation", "--checklist", "--json")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return string(output), 0
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("run status report parser: %v\n%s", err, output)
	}
	return string(output), exitErr.ExitCode()
}

func copyCodexPluginRoot(t *testing.T, repo string) string {
	t.Helper()
	pluginRoot := t.TempDir()
	cmd := exec.Command("cp", "-R", filepath.Join(repo, "skills"), pluginRoot)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("copy current checkout skills into supplied plugin root: %v\n%s", err, output)
	}
	return pluginRoot
}

func buildCodexBootstrapBinary(t *testing.T, repo string) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "spacedock")
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/spacedock")
	cmd.Dir = repo
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build spacedock status parser: %v\n%s", err, output)
	}
	return binary
}

func dispatchTestRepoRoot(t *testing.T) string {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller did not return dispatch test path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(source), "..", ".."))
}

func readCodexBootstrapEntity(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Codex bootstrap fixture entity %s: %v", path, err)
	}
	return string(data)
}
