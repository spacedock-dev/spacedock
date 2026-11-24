// ABOUTME: Content-addressed assignment publication and directory-selection behavior.
package dispatch

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func artifactArgs(t *testing.T, host, dir string) []string {
	t.Helper()
	root, ep := buildHostFixture(t)
	return []string{"build", "--host", host, "--workflow-dir", root, "--entity-path", ep, "--stage", "backlog", "--checklist-file", "-", "--dispatch-dir", dir}
}

func assertEmptyArtifacts(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 0 {
		t.Fatalf("unexpected artifacts in %s: %v %v", dir, entries, err)
	}
}

func TestArtifactOverlappingBriefsAndRetry(t *testing.T) {
	t.Setenv("CLAUDE_CODE_SESSION_ID", "same-parent")
	t.Setenv("PASEO_AGENT_ID", "same-parent")
	for _, host := range []string{"claude", "paseo"} {
		t.Run(host, func(t *testing.T) {
			dir := t.TempDir()
			args := artifactArgs(t, host, dir)
			start := make(chan struct{})
			results := make(chan struct {
				i   int
				res runResult
			}, 2)
			for i := 0; i < 2; i++ {
				go func(i int) {
					<-start
					results <- struct {
						i   int
						res runResult
					}{i, runNativeWithLauncher(fmt.Sprintf("- sentinel-%d\n- complete", i), testWorkflowLauncher, args...)}
				}(i)
			}
			close(start)
			paths := make([]string, 2)
			for i := 0; i < 2; i++ {
				r := <-results
				if r.res.exit != 0 {
					t.Fatalf("build: %+v", r.res)
				}
				paths[r.i] = dispatchFilePathFromStdout(t, r.res.stdout)
			}
			if paths[0] == paths[1] {
				t.Fatal("different briefs share one path")
			}
			for i, path := range paths {
				body := readDispatchBody(t, path)
				if !strings.Contains(body, fmt.Sprintf("- sentinel-%d\n- complete\n\n### Summary", i)) {
					t.Fatal("incomplete or replaced brief", body)
				}
			}
			before, err := os.Stat(paths[0])
			if err != nil {
				t.Fatal(err)
			}
			retry := runNativeWithLauncher("- sentinel-0\n- complete", testWorkflowLauncher, args...)
			after, err := os.Stat(paths[0])
			if err != nil {
				t.Fatal(err)
			}
			if retry.exit != 0 || dispatchFilePathFromStdout(t, retry.stdout) != paths[0] || !os.SameFile(before, after) || before.ModTime() != after.ModTime() {
				t.Fatalf("retry replaced file: %+v", retry)
			}
			entries, _ := os.ReadDir(dir)
			if len(entries) != 2 || after.Mode().Perm() != 0644 {
				t.Fatal("temporary files leaked or wrong mode", entries, after.Mode())
			}
		})
	}
}

func TestArtifactPublicationConflict(t *testing.T) {
	for _, kind := range []string{"bytes", "symlink", "directory"} {
		t.Run(kind, func(t *testing.T) {
			dir := t.TempDir()
			args := artifactArgs(t, "paseo", dir)
			initial := runNativeWithLauncher("- original", testWorkflowLauncher, args...)
			if initial.exit != 0 {
				t.Fatal(initial)
			}
			path := dispatchFilePathFromStdout(t, initial.stdout)
			original := readDispatchBody(t, path)
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "bytes":
				writeFile(t, path, "conflicting assignment")
			case "directory":
				if err := os.Mkdir(path, 0755); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				target := filepath.Join(t.TempDir(), "target")
				writeFile(t, target, original)
				if err := os.Symlink(target, path); err != nil {
					t.Fatal(err)
				}
			}
			before, _ := os.Lstat(path)
			res := runNativeWithLauncher("- original", testWorkflowLauncher, args...)
			after, _ := os.Lstat(path)
			if res.exit != 1 || res.stdout != "" || !strings.HasPrefix(res.stderr, "dispatch_file_write_failed: "+path+": ") || !os.SameFile(before, after) {
				t.Fatalf("conflict accepted or replaced: %+v", res)
			}
			if kind == "bytes" {
				b, _ := os.ReadFile(path)
				if string(b) != "conflicting assignment" {
					t.Fatal("conflict overwritten")
				}
			}
			entries, _ := os.ReadDir(dir)
			if len(entries) != 1 {
				t.Fatal("temporary files leaked", entries)
			}
		})
	}
}

func TestArtifactDirectorySelection(t *testing.T) {
	for _, mode := range []string{"flag", "equals", "env", "default", "relative", "spaces", "precedence"} {
		t.Run(mode, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "shared assignments", "nested")
			args := artifactArgs(t, "paseo", dir)
			switch mode {
			case "equals":
				args = append(args[:len(args)-2], "--dispatch-dir="+dir)
			case "env":
				t.Setenv("SPACEDOCK_DISPATCH_DIR", dir)
				args = args[:len(args)-2]
			case "default":
				dir = dispatchFileDir
				args = args[:len(args)-2]
			case "relative":
				originalCwd, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				cwdDir := t.TempDir()
				if err := os.Chdir(cwdDir); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.Chdir(originalCwd); err != nil {
						t.Fatal(err)
					}
				})
				cwd, err := os.Getwd()
				if err != nil {
					t.Fatal(err)
				}
				dir = filepath.Join(cwd, "shared assignments", "nested")
				args[len(args)-1] = filepath.Join("shared assignments", "tmp", "..", "nested")
			case "precedence":
				t.Setenv("SPACEDOCK_DISPATCH_DIR", "\n")
			}
			res := runNativeWithLauncher("- shared sentinel", testWorkflowLauncher, args...)
			if res.exit != 0 {
				t.Fatal(res)
			}
			path := dispatchFilePathFromStdout(t, res.stdout)
			t.Cleanup(func() { _ = os.Remove(path) })
			if !filepath.IsAbs(path) || filepath.Dir(path) != dir || !strings.Contains(res.stdout, path) {
				t.Fatal("wrong locator", res)
			}
			if !strings.Contains(readDispatchBody(t, path), "- shared sentinel") {
				t.Fatal("missing body")
			}
		})
	}
}

func TestArtifactInvalidDirectoryBeforeStamp(t *testing.T) {
	_, wf, state, ep := stampFixture(t, "implementation", "implementation", false)
	before, err := os.ReadFile(ep)
	if err != nil {
		t.Fatal(err)
	}
	head := gitOutput(t, state, "rev-parse", "HEAD")
	occupied := filepath.Join(t.TempDir(), "file")
	writeFile(t, occupied, "occupied")
	for _, source := range []string{"flag", "equals", "env"} {
		for _, value := range []string{"", " \t", "nul\x00", "line\n", "line\r", occupied} {
			if source == "env" && strings.ContainsRune(value, 0) {
				continue
			} // OS environments cannot contain NUL.
			t.Run(source+fmt.Sprintf("-%q", value), func(t *testing.T) {
				args := append(buildStampArgs(wf, ep, "implementation", "-"), "--host", "paseo")
				switch source {
				case "flag":
					args = append(args, "--dispatch-dir", value)
				case "equals":
					args = append(args, "--dispatch-dir="+value)
				case "env":
					t.Setenv("SPACEDOCK_DISPATCH_DIR", value)
				}
				res := runNativeWithLauncher("- stamp", testWorkflowLauncher, args...)
				want := 2
				if value == occupied {
					want = 1
				}
				if res.exit != want || res.stdout != "" {
					t.Fatal(res)
				}
				if want == 2 && res.stderr != fmt.Sprintf("error: invalid dispatch directory %q\n", value) {
					t.Fatal(res.stderr)
				}
				after, _ := os.ReadFile(ep)
				if !bytes.Equal(before, after) || head != gitOutput(t, state, "rev-parse", "HEAD") {
					t.Fatal("invalid directory stamped state")
				}
			})
		}
	}
	args := append(buildStampArgs(wf, ep, "implementation", "-"), "--dispatch-dir")
	res := runNativeWithLauncher("- stamp", testWorkflowLauncher, args...)
	if res.exit != 2 || res.stdout != "" || !strings.Contains(res.stderr, "requires value for --dispatch-dir") {
		t.Fatal(res)
	}
}

func TestArtifactHostsModesAndFilenameBoundary(t *testing.T) {
	for _, host := range []string{"claude", "codex", "pi", "paseo"} {
		for _, mode := range []string{"fresh", "advance", "bare"} {
			if mode == "bare" && (host == "codex" || host == "paseo") {
				continue
			}
			t.Run(host+"-"+mode, func(t *testing.T) {
				args := artifactArgs(t, host, t.TempDir())
				if mode != "fresh" {
					flag := "--advance"
					if mode == "bare" {
						flag = "--bare-mode"
					}
					args = append(args, flag)
				}
				res := runNativeWithLauncher("- all hosts", testWorkflowLauncher, args...)
				if res.exit != 0 {
					t.Fatal(res)
				}
				readDispatchBody(t, dispatchFilePathFromStdout(t, res.stdout))
			})
		}
	}
	for _, oldLen := range []int{242, 243} {
		name, err := dispatchArtifactName(strings.Repeat("a", oldLen), []byte("body"))
		if oldLen == 242 {
			if err != nil || len(name) != 254 {
				t.Fatalf("242-byte stem: %q %v", name, err)
			}
		} else if err == nil || name != "" || !strings.Contains(err.Error(), "exceeds 251 characters") {
			t.Fatalf("243-byte stem: %q %v", name, err)
		}
	}
}
