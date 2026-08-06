// ABOUTME: Self-contained dispatch keeps assignment payload in the file across every host.
// ABOUTME: Fresh and advance pointer prompts remain invariant as selected payload grows.
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

func TestGeneratedAssignmentRoutesWorkflowReadToPinnedA(t *testing.T) {
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
	aPath := writeExecutable(t, aDir, "spacedock", "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$A_LOG\"\nprintf '{}\\n'\n")
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
	launcher := "    " + shlexQuote(aPath)
	if !strings.Contains(body, launcher) || len(output.Fetch) != 0 {
		t.Fatalf("assignment did not retain pinned A or retained fetches: fetch=%v\n%s", output.Fetch, body)
	}

	env := append(environWithoutSpacedockBin(),
		"PATH="+bDir+":/usr/bin:/bin", "SPACEDOCK_BIN="+bPath, "A_LOG="+aLog, "B_LOG="+bLog)
	cmd := exec.Command("sh", "-c", shlexQuote(aPath)+" status --read "+shlexQuote(entityPath)+" --json")
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("pinned status read failed: %v\n%s", err, out)
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
	if got := string(aCalls); strings.Count(got, "status --read") != 1 || strings.Count(got, "C:version-under-test") != 1 {
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

func TestBuildAssignmentPayloadStaysFileOnlyAcrossHosts(t *testing.T) {
	const stage, context, checklist = "STAGE-ONLY", "CONTEXT-ONLY", "- CHECKLIST-ONLY"
	for _, host := range []string{"claude", "codex", "pi"} {
		for _, advance := range []bool{false, true} {
			t.Run(host+"/"+map[bool]string{false: "fresh", true: "advance"}[advance], func(t *testing.T) {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, "thing.md"), entityFM("Thing", "implementation", ""))
				gitInit(t, root)
				base, baseBody := selfContainedBuild(t, root, host, advance, "base", "base", []string{"- base"}, nil)
				variants := []struct {
					name, stage, context string
					checklist            []string
					extra                map[string]any
					wants                []string
				}{
					{"stage", "base\n" + stage, "base", []string{"- base"}, nil, []string{stage}},
					{"context", "base", "base\n" + context, []string{"- base"}, nil, []string{context}},
					{"checklist", "base", "base", []string{"- base", checklist}, nil, []string{checklist}},
					{"scope-feedback", "base", "base", []string{"- base"}, map[string]any{
						"scope_notes": "SCOPE-ONLY", "feedback_context": "FEEDBACK-ONLY", "is_feedback_reflow": true,
					}, []string{"SCOPE-ONLY", "FEEDBACK-ONLY"}},
				}
				for _, variant := range variants {
					got, body := selfContainedBuild(t, root, host, advance, variant.stage, variant.context, variant.checklist, variant.extra)
					if got.Prompt != base.Prompt {
						t.Errorf("%s changed pointer prompt\nbase=%q\ngot=%q", variant.name, base.Prompt, got.Prompt)
					}
					if len(got.Fetch) != 0 || strings.Contains(body, "### Fetch commands") {
						t.Errorf("%s restored worker-time resolution: fetch=%v", variant.name, got.Fetch)
					}
					for _, sentinel := range variant.wants {
						if !strings.Contains(body, sentinel) || strings.Contains(got.Prompt, sentinel) {
							t.Errorf("%s sentinel %q crossed file/prompt boundary", variant.name, sentinel)
						}
					}
					if variant.name != "scope-feedback" && len(body) <= len(baseBody) {
						t.Errorf("%s payload did not grow dispatch file", variant.name)
					}
				}
				writeMods(t, root, map[string]string{"helper.md": standingMod("helper", "sonnet", "STANDING-ONLY", "")})
				got, body := selfContainedBuild(t, root, host, advance, "base", "base", []string{"- base"}, nil)
				if got.Prompt != base.Prompt || !strings.Contains(body, "STANDING-ONLY") || strings.Contains(got.Prompt, "STANDING-ONLY") {
					t.Error("standing payload crossed file/prompt boundary")
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
