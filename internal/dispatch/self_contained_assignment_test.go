// ABOUTME: Dispatch fetches pin one resolved executable across every supported host and mode.
// ABOUTME: Pointer prompts remain payload-free while structured fetch commands stay executable.
package dispatch

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type selfContainedOutput struct {
	DispatchFile string   `json:"dispatch_file_path"`
	Prompt       string   `json:"prompt"`
	Fetch        []string `json:"fetch_commands"`
	Name         *string  `json:"name"`
}

func TestGeneratedAssignmentExecutesStageLoadThroughPinnedA(t *testing.T) {
	root := t.TempDir()
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, filepath.Join(root, "README.md"), "---\nstages:\n  states:\n    - name: implementation\n      initial: true\n---\n# Fixture\n\n### implementation\nwork\n")
	writeFile(t, entityPath, entityFM("Thing", "implementation", ""))
	gitInit(t, root)

	aDir := filepath.Join(t.TempDir(), "launcher A")
	if err := os.MkdirAll(aDir, 0o755); err != nil {
		t.Fatal(err)
	}
	aLog := filepath.Join(t.TempDir(), "a.log")
	bLog := filepath.Join(t.TempDir(), "b.log")
	aPath := writeExecutable(t, aDir, "spacedock", "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$A_LOG\"\nprintf 'A-STAGE\\n'\n")
	bDir := t.TempDir()
	bPath := writeExecutable(t, bDir, "spacedock", "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$B_LOG\"\nexit 91\n")
	cPath := writeExecutable(t, t.TempDir(), "product-C", "#!/bin/sh\nprintf 'C:%s\\n' \"$*\" >> \"$A_LOG\"\n")

	fields := map[string]any{
		"schema_version": 2, "entity_path": entityPath, "workflow_dir": root,
		"stage": "implementation", "checklist": []string{"- pin A"}, "host": "codex",
	}
	built := runNativeWithLauncher(mergeStdin(fields, nil), aPath, "build", "--workflow-dir", root, "--host", "codex")
	if built.exit != 0 {
		t.Fatalf("build exit=%d stderr=%s", built.exit, built.stderr)
	}
	var output selfContainedOutput
	if err := json.Unmarshal([]byte(built.stdout), &output); err != nil {
		t.Fatal(err)
	}
	body := readDispatchBody(t, output.DispatchFile)
	wantCommand := shlexQuote(aPath) + " dispatch show-stage-def --workflow-dir " + shlexQuote(root) + " --stage implementation"
	sections := dispatchArtifactSections(t, body)
	if len(output.Fetch) != 1 || output.Fetch[0] != wantCommand || sections["Fetch commands"] != wantCommand {
		t.Fatalf("assignment did not own the exact A command: fetch=%v\n%s", output.Fetch, body)
	}

	env := append(environWithoutSpacedockBin(),
		"PATH="+bDir+":/usr/bin:/bin", "SPACEDOCK_BIN="+bPath, "A_LOG="+aLog, "B_LOG="+bLog)
	cmd := exec.Command("sh", "-c", output.Fetch[0])
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil || string(out) != "A-STAGE\n" {
		t.Fatalf("pinned stage load failed: err=%v out=%q", err, out)
	}
	product := exec.Command(cPath, "version-under-test")
	product.Env = env
	if out, err := product.CombinedOutput(); err != nil {
		t.Fatalf("explicit product C failed: %v\n%s", err, out)
	}
	aCalls, err := os.ReadFile(aLog)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(aCalls); strings.Count(got, "dispatch show-stage-def") != 1 || strings.Count(got, "C:version-under-test") != 1 {
		t.Fatalf("unexpected A/C calls:\n%s", got)
	}
	if bCalls, err := os.ReadFile(bLog); err == nil && len(bCalls) != 0 {
		t.Fatalf("ambient B was invoked:\n%s", bCalls)
	}
}

func TestBuildWithoutResolvedLauncherFailsBeforeWritingArtifact(t *testing.T) {
	root := t.TempDir()
	entityPath := filepath.Join(root, "unique-no-launcher.md")
	writeFile(t, filepath.Join(root, "README.md"), "---\nstages:\n  states:\n    - name: implementation\n      initial: true\n---\n# Fixture\n\n### implementation\nwork\n")
	writeFile(t, entityPath, entityFM("Unique", "implementation", ""))
	gitInit(t, root)
	fields := map[string]any{
		"schema_version": 2, "entity_path": entityPath, "workflow_dir": root,
		"stage": "implementation", "checklist": []string{"- fail closed"}, "host": "codex",
	}
	path := filepath.Join(dispatchFileDir, "spacedock-ensign-unique-no-launcher-implementation.md")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("test artifact already exists: %s", path)
	}
	result := runNativeWithLauncher(mergeStdin(fields, nil), "", "build", "--workflow-dir", root, "--host", "codex")
	if result.exit != 1 || result.stdout != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", result.exit, result.stdout, result.stderr)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("unresolved launcher wrote artifact: %s", path)
	}
}

func environWithoutSpacedockBin() []string {
	filtered := make([]string, 0, len(os.Environ()))
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "SPACEDOCK_BIN=") {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func writeExecutable(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func selfContainedBuild(t *testing.T, root, host string, advance bool, stage, context string, checklist []string, extra map[string]any) (selfContainedOutput, string) {
	t.Helper()
	readme := "---\nentity-type: task\nid-style: slug\nstages:\n  defaults:\n    worktree: false\n    context-sections: [Authority]\n  states:\n    - name: implementation\n      initial: true\n---\n# Fixture\n\n### implementation\n\n" + stage + "\n\n## Authority\n\n" + context + "\n"
	writeFile(t, filepath.Join(root, "README.md"), readme)
	fields := map[string]any{
		"schema_version": 2, "entity_path": filepath.Join(root, "thing.md"),
		"workflow_dir": root, "stage": "implementation", "checklist": checklist,
		"host": host, "advance": advance,
	}
	for key, value := range extra {
		fields[key] = value
	}
	run := runNative(mergeStdin(fields, nil), "build", "--workflow-dir", root, "--host", host)
	if run.exit != 0 {
		t.Fatalf("build exit=%d stderr=%s", run.exit, run.stderr)
	}
	var output selfContainedOutput
	if err := json.Unmarshal([]byte(run.stdout), &output); err != nil {
		t.Fatal(err)
	}
	return output, readDispatchBody(t, output.DispatchFile)
}

// dispatchArtifactSections returns the bounded level-three sections emitted by
// dispatch build. Tests compare structured inputs with their owned sections;
// narrative wording outside those sections is intentionally not an oracle.
func dispatchArtifactSections(t *testing.T, body string) map[string]string {
	t.Helper()
	sections := make(map[string]string)
	var heading string
	var lines []string
	flush := func() {
		if heading != "" {
			sections[heading] = strings.TrimSpace(strings.Join(lines, "\n"))
		}
	}
	for _, line := range strings.Split(body, "\n") {
		if name, ok := strings.CutPrefix(line, "### "); ok {
			flush()
			heading = name
			lines = nil
			continue
		}
		if heading != "" {
			lines = append(lines, strings.TrimPrefix(line, "    "))
		}
	}
	flush()
	return sections
}

func TestBuildHostModeMatrixKeepsPointerTransportAndExactFetchShape(t *testing.T) {
	for _, host := range []string{"claude", "codex", "pi"} {
		for _, advance := range []bool{false, true} {
			t.Run(host+"/"+map[bool]string{false: "fresh", true: "advance"}[advance], func(t *testing.T) {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, "thing.md"), entityFM("Thing", "implementation", ""))
				gitInit(t, root)
				checklist := []string{"- CHECKLIST-SENTINEL"}
				got, body := selfContainedBuild(t, root, host, advance, "STAGE-SENTINEL", "CONTEXT-SENTINEL", checklist, nil)
				if len(got.Fetch) != 1 {
					t.Fatalf("fetch command count=%d, want 1: %v", len(got.Fetch), got.Fetch)
				}
				want := testWorkflowLauncher + " dispatch show-stage-def --workflow-dir " + shlexQuote(root) + " --stage implementation"
				sections := dispatchArtifactSections(t, body)
				if got.Fetch[0] != want || sections["Fetch commands"] != want {
					t.Fatalf("structured fetch mismatch: got=%v\nbody=%s", got.Fetch, body)
				}
				if sections["Completion checklist"] != strings.Join(checklist, "\n") {
					t.Fatalf("checklist section=%q, want exact structured input %q", sections["Completion checklist"], strings.Join(checklist, "\n"))
				}
				if _, exists := sections["Stage definition"]; exists {
					t.Fatal("fetched stage definition was unexpectedly inlined")
				}
				if _, exists := sections["Authority"]; exists {
					t.Fatal("fetched context section was unexpectedly inlined")
				}
			})
		}
	}
}

// The legacy-era sibling of this file, TestBuildMaxLegalDispatchFilenamesRemainExact,
// hit the dispatchFileNameMaxLen boundary via an unbounded team_name prefix. That
// input retired with legacy team mode: the merged floor's session-token prefix is
// capped well under the boundary (sessionTokenMaxLen), so the boundary is no
// longer reachable through any supported input. Session-prefix + -advance suffix
// ordering and no-collision coverage live in TestBuildAdvanceFilenameSuffix and
// TestBuildMergedModeDispatchFileDisambiguator.
