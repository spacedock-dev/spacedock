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
	if len(output.Fetch) != 1 || output.Fetch[0] != wantCommand || !strings.Contains(body, "    "+wantCommand) {
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
	if result.exit != 1 || !strings.Contains(result.stderr, "refusing to write a dispatch artifact") {
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
		"team_name": "fixture-team", "host": host, "advance": advance,
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

func TestBuildHostModeMatrixKeepsPointerTransportAndExactFetchShape(t *testing.T) {
	for _, host := range []string{"claude", "codex", "pi"} {
		for _, advance := range []bool{false, true} {
			t.Run(host+"/"+map[bool]string{false: "fresh", true: "advance"}[advance], func(t *testing.T) {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, "thing.md"), entityFM("Thing", "implementation", ""))
				gitInit(t, root)
				got, body := selfContainedBuild(t, root, host, advance, "STAGE-SENTINEL", "CONTEXT-SENTINEL", []string{"- CHECKLIST-SENTINEL"}, nil)
				if len(got.Fetch) != 1 {
					t.Fatalf("fetch command count=%d, want 1: %v", len(got.Fetch), got.Fetch)
				}
				want := testWorkflowLauncher + " dispatch show-stage-def --workflow-dir " + shlexQuote(root) + " --stage implementation"
				if got.Fetch[0] != want || !strings.Contains(body, "    "+want) {
					t.Fatalf("structured fetch mismatch: got=%v\nbody=%s", got.Fetch, body)
				}
				for _, forbidden := range []string{"STAGE-SENTINEL", "CONTEXT-SENTINEL"} {
					if strings.Contains(got.Prompt, forbidden) || strings.Contains(body, forbidden) {
						t.Errorf("stage payload %q crossed into pointer transport or bootstrap body", forbidden)
					}
				}
				if strings.Contains(got.Prompt, "CHECKLIST-SENTINEL") || !strings.Contains(body, "CHECKLIST-SENTINEL") {
					t.Error("checklist did not remain file-only")
				}
			})
		}
	}
}

func TestBuildMaxLegalDispatchFilenamesRemainExact(t *testing.T) {
	root := t.TempDir()
	entityPath := filepath.Join(root, strings.Repeat("x", 33)+".md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", ""))
	gitInit(t, root)
	seed, _ := selfContainedBuild(t, root, "claude", false, "work", "authority", []string{"- exact"}, map[string]any{"entity_path": entityPath})
	if seed.Name == nil {
		t.Fatal("fresh output omitted worker name")
	}
	for _, tc := range []struct {
		advance bool
		base    int
		suffix  string
	}{{false, 251, ""}, {true, 243, "-advance"}} {
		prefix := strings.Repeat("a", tc.base-len(*seed.Name)-2)
		var paths []string
		for _, final := range []string{"b", "c"} {
			team := prefix + final
			got, _ := selfContainedBuild(t, root, "claude", tc.advance, "work", "authority", []string{"- exact"},
				map[string]any{"entity_path": entityPath, "team_name": team})
			want := team + "-" + *seed.Name + tc.suffix
			if strings.TrimSuffix(filepath.Base(got.DispatchFile), ".md") != want || len(want) != dispatchFileNameMaxLen {
				t.Fatalf("dispatch stem changed or shortened: got=%q want=%q", filepath.Base(got.DispatchFile), want)
			}
			if len(got.Prompt) <= 300 {
				t.Fatalf("max-legal pointer prompt length=%d, want >300", len(got.Prompt))
			}
			paths = append(paths, got.DispatchFile)
		}
		if paths[0] == paths[1] {
			t.Fatalf("distinct max-legal names collided: %s", paths[0])
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err != nil {
				t.Fatal(err)
			}
		}
	}
}
