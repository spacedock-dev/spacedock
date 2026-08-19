// ABOUTME: `merge guard` routing/help/completion (AC-7) and the end-to-end ceremony
// ABOUTME: through cobra — armed/finalized/refused drive the real status handler.
package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// stageMergeFixture copies a status-package merge fixture into a fresh
// git-initialized temp dir and returns its root, so `merge guard` drives the real
// --set/--archive ceremony end-to-end through cobra (mutations resolve git_root).
func stageMergeFixture(t *testing.T, fixture string) string {
	t.Helper()
	src, err := filepath.Abs(filepath.Join("..", "status", "testdata", fixture))
	if err != nil {
		t.Fatal(err)
	}
	dst := t.TempDir()
	if err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	}); err != nil {
		t.Fatal(err)
	}
	// InitRepo persists a PERSISTENT local identity on the temp repo. The verb's own
	// commits (commitArchiveMove) run plain `git commit` without `-c`, so they must
	// resolve an identity from the repo's config — independent of global/system
	// config and git's auto-detection. A CI lane with no global identity and
	// auto-detection disabled would otherwise exit-128 the verb's commit.
	testgit.InitRepo(t, dst, "-q")
	for _, args := range [][]string{
		{"add", "-A"},
		{"commit", "-q", "-m", "seed"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dst
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return dst
}

// runMergeCLI's env is hermeticEnv(), not os.Environ() (gate_ceremony_count_test.go):
// this suite's fixtures carry no boot receipt, and an ambient
// CLAUDE_CODE_SESSION_ID (this suite dogfoods spacedock and can run inside a
// live Claude Code session) would make the boot guard spuriously refuse.
func runMergeCLI(t *testing.T, dir string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errBuf bytes.Buffer
	code = run(context.Background(), args, hermeticEnv(), dir, nil, &out, &errBuf, &status.NativeRunner{}, nil)
	return out.String(), errBuf.String(), code
}

// TestMergeGuardHelpRenders (AC-7): `merge guard --help` exits 0 with usage text.
func TestMergeGuardHelpRenders(t *testing.T) {
	for _, args := range [][]string{{"merge", "--help"}, {"merge", "guard", "--help"}} {
		out, errOut, code := runMergeCLI(t, t.TempDir(), args...)
		if code != 0 {
			t.Fatalf("%v exit=%d stderr=%q", args, code, errOut)
		}
		if !strings.Contains(out, "merge guard") || !strings.Contains(out, "--verdict") {
			t.Fatalf("%v help missing usage/--verdict:\n%s", args, out)
		}
	}
}

// TestMergeGuardHelpDocumentsWorkflowDirResolution is AC-6: the merge help text
// states the cwd-relative --workflow-dir resolution rule, so an operator reads
// the fix before hitting the issue-#485 foreign-cwd confusion.
func TestMergeGuardHelpDocumentsWorkflowDirResolution(t *testing.T) {
	out, errOut, code := runMergeCLI(t, t.TempDir(), "merge", "guard", "--help")
	if code != 0 {
		t.Fatalf("merge guard --help exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(out, "resolves against the current directory") || !strings.Contains(out, "absolute path") {
		t.Fatalf("merge guard --help missing cwd-relative --workflow-dir resolution rule:\n%s", out)
	}
}

// TestMergeGuardInGroupedHelp (AC-7): `merge guard <slug>` appears in the
// top-level grouped help workflow group.
func TestMergeGuardInGroupedHelp(t *testing.T) {
	out, _, code := runMergeCLI(t, t.TempDir(), "--help")
	if code != 0 {
		t.Fatalf("--help exit=%d", code)
	}
	if !strings.Contains(out, "merge") || !strings.Contains(out, "guard <slug>") {
		t.Fatalf("grouped help missing the merge guard row:\n%s", out)
	}
}

// TestMergeInCompletionScripts (AC-7): both completion scripts list `merge` as a
// verb and complete `guard` under it.
func TestMergeInCompletionScripts(t *testing.T) {
	for _, shell := range []string{"bash", "zsh"} {
		var stdout, stderr bytes.Buffer
		if code := Run([]string{"completion", shell}, &stdout, &stderr); code != 0 {
			t.Fatalf("completion %s exit=%d stderr=%q", shell, code, stderr.String())
		}
		script := stdout.String()
		if !strings.Contains(script, " merge ") && !strings.Contains(script, " merge\n") {
			t.Fatalf("completion %s script does not list merge as a verb:\n%s", shell, script)
		}
		if !strings.Contains(script, "merge)") || !strings.Contains(script, "guard") {
			t.Fatalf("completion %s script missing the merge) guard case:\n%s", shell, script)
		}
	}
}

// TestMergeUnknownSubcommand (AC-7 companion): `merge bogus` is a usage error
// (exit 2), matching the state-verb unknown-subcommand contract.
func TestMergeUnknownSubcommand(t *testing.T) {
	_, errOut, code := runMergeCLI(t, t.TempDir(), "merge", "bogus")
	if code != 2 {
		t.Fatalf("merge bogus must exit 2, got %d", code)
	}
	if !strings.Contains(errOut, "unknown subcommand") {
		t.Fatalf("stderr should name the unknown subcommand, got %q", errOut)
	}
}

// TestMergeGuardEndToEndArmFinalize drives the whole ceremony through cobra on a
// merge: local fixture: the first `merge guard` arms, the second finalizes and
// archives — proving the verb routes to the handler and the real --set/--archive
// paths run end to end.
func TestMergeGuardEndToEndArmFinalize(t *testing.T) {
	root := stageMergeFixture(t, "merge-local-workflow")

	out1, err1, code1 := runMergeCLI(t, root, "merge", "guard", "020-no-sentinel", "--verdict", "passed", "--workflow-dir", root)
	if code1 != 0 {
		t.Fatalf("arm exit=%d stderr=%q", code1, err1)
	}
	if !strings.Contains(out1, "armed") {
		t.Fatalf("first run should arm, got %q", out1)
	}

	out2, err2, code2 := runMergeCLI(t, root, "merge", "guard", "020-no-sentinel", "--verdict", "passed", "--workflow-dir", root)
	if code2 != 0 {
		t.Fatalf("finalize exit=%d stderr=%q", code2, err2)
	}
	if !strings.Contains(out2, "finalized") {
		t.Fatalf("second run should finalize, got %q", out2)
	}
	if !fileExists(filepath.Join(root, "_archive", "020-no-sentinel.md")) {
		t.Fatal("finalize should archive the entity")
	}
}

// TestMergeGuardEndToEndRefusalPropagated drives the AC-5 refusal through cobra:
// under merge: pr the verb auto-arms an empty-mod-block entity, but if the hook is
// never invoked (no merge sentinel) the re-run finalize hits the merge-hook guard,
// whose exit 1 + stderr must reach the process exit code unchanged — the verb never
// bypasses the guard. First run arms (exit 0); the second run refuses (exit 1).
func TestMergeGuardEndToEndRefusalPropagated(t *testing.T) {
	root := stageMergeFixture(t, "merge-pr-workflow")
	out1, err1, code1 := runMergeCLI(t, root, "merge", "guard", "020-no-sentinel", "--verdict", "passed", "--workflow-dir", root)
	if code1 != 0 {
		t.Fatalf("auto-arm should exit 0, got %d (stderr=%q)", code1, err1)
	}
	if !strings.Contains(out1, "armed") {
		t.Fatalf("first run should auto-arm, got %q", out1)
	}
	out2, errOut, code := runMergeCLI(t, root, "merge", "guard", "020-no-sentinel", "--verdict", "passed", "--workflow-dir", root)
	if code != 1 {
		t.Fatalf("armed-but-no-merge refusal should propagate exit 1, got %d (stdout=%q)", code, out2)
	}
	if !strings.Contains(errOut, "cannot advance to terminal") {
		t.Fatalf("guard refusal should reach stderr verbatim, got %q", errOut)
	}
}
