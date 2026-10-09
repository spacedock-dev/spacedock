package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

func TestStateStorageRefusalAndRecovery(t *testing.T) {
	for _, kind := range []string{"missing", "ordinary", "wrong-branch", "detached", "subdirectory"} {
		t.Run(kind, func(t *testing.T) {
			root, wf := writeSplitReadmeRepo(t)
			checkout := filepath.Join(wf, ".spacedock-state")
			if kind == "ordinary" || kind == "subdirectory" {
				os.MkdirAll(checkout, 0755)
			}
			if kind == "wrong-branch" || kind == "detached" {
				if c, o, e := execStateNew(t, root, wf); c != 0 {
					t.Fatalf("birth %d %s %s", c, o, e)
				}
				if kind == "wrong-branch" {
					git(t, checkout, "checkout", "-qb", "wrong")
				} else {
					git(t, checkout, "checkout", "--detach")
				}
			}
			code := "state-checkout-invalid"
			if kind == "missing" {
				code = "state-checkout-missing"
			}
			snapshot := func() string {
				result := git(t, root, "show-ref") + git(t, root, "ls-files", "--stage")
				err := filepath.WalkDir(wf, func(path string, d os.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.Name() == ".git" {
						if d.IsDir() {
							return filepath.SkipDir
						}
						return nil
					}
					if !d.IsDir() {
						b, err := os.ReadFile(path)
						if err != nil {
							return err
						}
						result += fmt.Sprintf("%s:%x\n", path, b)
					}
					return nil
				})
				if err != nil {
					t.Fatal(err)
				}
				if kind == "wrong-branch" || kind == "detached" {
					result += git(t, checkout, "ls-files", "--stage")
				}
				return result
			}
			before := snapshot()
			for _, args := range [][]string{{"status"}, {"status", "--json"}, {"status", "--boot", "--json"}, {"new", "unsafe"}, {"status", "--new", "unsafe"}} {
				var out, errOut strings.Builder
				args = append(args, "--workflow-dir", wf)
				c := run(context.Background(), args, os.Environ(), root, strings.NewReader("---\nstatus: ideation\n---\n# Unsafe\n"), &out, &errOut, &status.NativeRunner{}, nil)
				if c != 1 || !strings.Contains(out.String()+errOut.String(), code) || strings.Contains(out.String(), `"entities"`) {
					t.Fatalf("%v: %d %s %s", args, c, out.String(), errOut.String())
				}
			}
			after := snapshot()
			if before != after {
				t.Fatal("refusal mutated filesystem/index/refs")
			}
			if _, err := os.Stat(filepath.Join(checkout, "unsafe.md")); !os.IsNotExist(err) {
				t.Fatal("unsafe filing created state")
			}
			if kind == "missing" {
				if _, err := os.Stat(checkout); !os.IsNotExist(err) {
					t.Fatal("read created checkout")
				}
			}
			if c, o, e := terminalInvoke(t, root, "status", "--workflow-dir", wf, "--read", filepath.Join(wf, "README.md"), "--json"); c != 0 {
				t.Fatalf("definition read %d %s %s", c, o, e)
			}
			if c, o, e := terminalInvoke(t, root, "status", "--discover", "--root", root); c != 0 {
				t.Fatalf("discovery %d %s %s", c, o, e)
			}
			if kind == "missing" {
				if c, o, e := execStateNew(t, root, wf); c != 0 {
					t.Fatalf("recovery birth %d %s %s", c, o, e)
				}
				if c, o, e := terminalInvoke(t, root, "status", "--workflow-dir", wf, "--json"); c != 0 || !strings.Contains(o, `"entities":[]`) {
					t.Fatalf("empty status %d %s %s", c, o, e)
				}
				var out, errOut strings.Builder
				c := run(context.Background(), []string{"new", "safe", "--workflow-dir", wf}, os.Environ(), root, strings.NewReader("---\nstatus: ideation\n---\n# Safe\n"), &out, &errOut, &status.NativeRunner{}, nil)
				if c != 0 {
					t.Fatalf("safe filing %d %s %s", c, out.String(), errOut.String())
				}
				if c, o, e := runStateCommitCmd(t, root, wf, "safe"); c != 0 {
					t.Fatalf("durable filing %d %s %s", c, o, e)
				}
				clone := filepath.Join(t.TempDir(), "clone")
				git(t, root, "clone", "-q", root, clone)
				cloneWF := filepath.Join(clone, "docs", "dev")
				if c, _, _ := terminalInvoke(t, clone, "status", "--workflow-dir", cloneWF); c != 1 {
					t.Fatal("clone falsely healthy before state init")
				}
				stateInit(t, clone, cloneWF)
				if c, o, e := terminalInvoke(t, clone, "status", "--workflow-dir", cloneWF, "--json"); c != 0 || !strings.Contains(o, "safe") {
					t.Fatalf("init recovery %d %s %s", c, o, e)
				}

			}
		})
	}
}
