// ABOUTME: AC-1/AC-2 from-root `new` discovery — the no-flag path falls back to a
// ABOUTME: downward scan, succeeding on a single workflow and refusing on ambiguity.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// plainRootREADME is a non-workflow toplevel README — no commissioned-by — so the
// upward walk from the repo root finds nothing and the downward fallback is what
// must resolve the nested workflow.
const plainRootREADME = "---\ntitle: Plain Repo\n---\n# Plain Repo\n"

// TestNewFromRootSingleWorkflowSucceeds (AC-1) drives the --new path from the
// repo toplevel with no --workflow-dir, where the only commissioned workflow is a
// nested subdir. The downward-scan fallback must resolve it: exit 0, the
// `created:` line, and the file landing under the nested workflow's entity dir.
func TestNewFromRootSingleWorkflowSucceeds(t *testing.T) {
	env := pinnedEnv(t)
	root := t.TempDir()
	writeAll(t, root, map[string]string{
		"README.md":          plainRootREADME,
		"docs/dev/README.md": commissionedSeqREADME,
	})
	gitInit(t, root)

	out, errOut, code := runNativeStdin(t, root, env, reader(newBody), "--new", "from-root-task")
	if code != 0 {
		t.Fatalf("--new from root exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(out, "created:") {
		t.Fatalf("--new from root narration %q should report created:", out)
	}
	written := filepath.Join(root, "docs/dev", "from-root-task.md")
	if !isRegularFile(written) {
		t.Fatalf("entity should exist under the nested workflow dir: %s", written)
	}
}

// TestNewFromRootMultiWorkflowRefuses (AC-2) puts two commissioned workflows
// under the toplevel; from-root `new` must refuse (non-zero) and name both
// candidate dirs, then succeed when --workflow-dir disambiguates.
func TestNewFromRootMultiWorkflowRefuses(t *testing.T) {
	env := pinnedEnv(t)
	root := t.TempDir()
	writeAll(t, root, map[string]string{
		"README.md":          plainRootREADME,
		"docs/dev/README.md": commissionedSeqREADME,
		"docs/ops/README.md": commissionedSeqREADME,
	})
	gitInit(t, root)

	_, errOut, code := runNativeStdin(t, root, env, reader(newBody), "--new", "ambiguous-task")
	if code == 0 {
		t.Fatalf("--new from root with two workflows must refuse, got exit 0")
	}
	dev := filepath.Join(root, "docs/dev")
	ops := filepath.Join(root, "docs/ops")
	if !strings.Contains(errOut, realpathOf(dev)) || !strings.Contains(errOut, realpathOf(ops)) {
		t.Fatalf("ambiguity error must name both candidates, got: %q", errOut)
	}
	if !strings.Contains(errOut, "--workflow-dir") {
		t.Fatalf("ambiguity error must instruct --workflow-dir, got: %q", errOut)
	}

	out, errOut, code := runNativeStdin(t, root, env, reader(newBody), "--workflow-dir", dev, "--new", "ambiguous-task")
	if code != 0 {
		t.Fatalf("--new --workflow-dir exit=%d stderr=%q", code, errOut)
	}
	if !strings.Contains(out, "created:") {
		t.Fatalf("disambiguated --new narration %q should report created:", out)
	}
	if !isRegularFile(filepath.Join(dev, "ambiguous-task.md")) {
		t.Fatalf("entity should exist under the chosen workflow dir")
	}
}
