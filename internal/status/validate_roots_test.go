package status

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func splitRootValidationFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	def := filepath.Join(root, "docs", "dev")
	state := filepath.Join(def, ".spacedock-state")
	writeFile(t, filepath.Join(def, "README.md"), `---
commissioned-by: spacedock@test
entity-type: task
id-style: sequential
state: .spacedock-state
stages:
  states:
    - name: backlog
      initial: true
    - name: done
      terminal: true
---

# Dev Workflow
`)
	if err := os.MkdirAll(state, 0o755); err != nil {
		t.Fatal(err)
	}
	gitInit(t, root)
	return def, state
}

func hookEnv(t *testing.T, base []string, state string) []string {
	t.Helper()
	gitDir := filepath.Join(state, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	index := filepath.Join(t.TempDir(), "index")
	return append(append([]string{}, base...),
		"GIT_DIR="+gitDir,
		"GIT_INDEX_FILE="+index,
	)
}

func commitAll(t *testing.T, dir, msg string) {
	t.Helper()
	for _, args := range [][]string{
		{"-c", "user.email=t@t", "-c", "user.name=t", "add", "-A"},
		{"-c", "user.email=t@t", "-c", "user.name=t", "commit", "-q", "-m", msg},
	} {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
}

func TestValidateFailsClosedForMisresolvedSplitRoot(t *testing.T) {
	def, state := splitRootValidationFixture(t)
	baseEnv := pinnedEnv(t)
	envs := map[string][]string{
		"clean": baseEnv,
		"hook":  hookEnv(t, baseEnv, state),
	}
	cases := map[string]string{
		"state-cwd":      state,
		"definition-cwd": def,
	}

	for envName, env := range envs {
		for cwdName, cwd := range cases {
			t.Run(envName+"/"+cwdName, func(t *testing.T) {
				out, errOut, code := runNative(t, cwd, env, "--workflow-dir", "docs/dev", "--validate")
				if code == 0 {
					t.Fatalf("wrong root returned success: stdout=%q stderr=%q", out, errOut)
				}
				if strings.Contains(out, "VALID") {
					t.Fatalf("wrong root must not print VALID, stdout=%q", out)
				}
				if !strings.Contains(errOut, "--validate requires --workflow-dir to resolve to a commissioned Spacedock workflow") {
					t.Fatalf("stderr missing root-resolution error, got %q", errOut)
				}
			})
		}
	}
}

func TestValidateReportsUntrackedBlankIDFlatAndFolderEntities(t *testing.T) {
	def, state := splitRootValidationFixture(t)
	writeFile(t, filepath.Join(state, "blank-flat.md"), "---\nid:\ntitle: Flat\nstatus: backlog\n---\n# Flat\n")
	writeFile(t, filepath.Join(state, "blank-folder", "index.md"), "---\nid:\ntitle: Folder\nstatus: backlog\n---\n# Folder\n")
	baseEnv := pinnedEnv(t)
	envs := map[string][]string{
		"clean": baseEnv,
		"hook":  hookEnv(t, baseEnv, state),
	}

	for envName, env := range envs {
		t.Run(envName, func(t *testing.T) {
			for i := 0; i < 5; i++ {
				out, errOut, code := runNative(t, state, env, "--workflow-dir", def, "--validate")
				if code == 0 {
					t.Fatalf("run %d returned success: stdout=%q stderr=%q", i, out, errOut)
				}
				if strings.Contains(out, "VALID") {
					t.Fatalf("run %d must not print VALID, stdout=%q", i, out)
				}
				for _, want := range []string{
					"Error: missing required id:",
					"slug=blank-flat",
					"slug=blank-folder",
				} {
					if !strings.Contains(errOut, want) {
						t.Fatalf("run %d stderr missing %q, got %q", i, want, errOut)
					}
				}
			}
		})
	}
}

func TestValidateKeepsTrackedStagedCommittedEntitiesWhileFailingBadUntracked(t *testing.T) {
	def, state := splitRootValidationFixture(t)
	writeFile(t, filepath.Join(state, "committed.md"), ent(`"001"`, "backlog"))
	commitAll(t, filepath.Dir(filepath.Dir(def)), "committed state entity")
	writeFile(t, filepath.Join(state, "staged.md"), ent(`"002"`, "backlog"))
	cmd := exec.Command("git", "-C", filepath.Dir(filepath.Dir(def)), "add", filepath.Join(state, "staged.md"))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add staged entity: %v\n%s", err, out)
	}
	writeFile(t, filepath.Join(state, "untracked-good.md"), ent(`"003"`, "done"))
	writeFile(t, filepath.Join(state, "bad-untracked.md"), "---\nid:\ntitle: Bad\nstatus: backlog\n---\n# Bad\n")

	out, errOut, code := runNative(t, state, pinnedEnv(t), "--workflow-dir", def, "--validate")
	if code == 0 {
		t.Fatalf("bad untracked sibling should fail validation: stdout=%q stderr=%q", out, errOut)
	}
	if strings.Contains(out, "VALID") {
		t.Fatalf("validation failure must not print VALID, stdout=%q", out)
	}
	for _, want := range []string{"slug=bad-untracked", "Error: missing required id:"} {
		if !strings.Contains(errOut, want) {
			t.Fatalf("stderr missing %q, got %q", want, errOut)
		}
	}
	for _, notWant := range []string{"slug=committed", "slug=staged", "slug=untracked-good"} {
		if strings.Contains(errOut, notWant) {
			t.Fatalf("valid tracked/staged/committed entity reported as invalid: %q", errOut)
		}
	}
}
