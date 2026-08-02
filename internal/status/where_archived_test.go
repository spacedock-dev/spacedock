// ABOUTME: AC-2 regression pin — --archived composes with --where as
// ABOUTME: active-plus-archived, not an archived-only scope swap.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestWhereArchivedComposesAsActivePlusArchived pins AC-2: --archived is
// already inclusive (activeAndArchivedEntities = scanEntitiesActive +
// archiveEntities), so a --where clause matching one active and one archived
// entity returns only the active member without --archived, and both members
// with it. enum-scope-workflow carries exactly this shape: top-placed (active)
// and arch-placed (archived) share status=backlog, differing only by
// placement — this must fail if --archived ever becomes archived-only or stops
// composing with --where.
func TestWhereArchivedComposesAsActivePlusArchived(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "enum-scope-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)

	t.Run("active-only", func(t *testing.T) {
		out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--where", "status=backlog")
		if code != 0 {
			t.Fatalf("exit=%d, want 0 (err=%q)", code, errOut)
		}
		if !strings.Contains(out, "top-placed") {
			t.Fatalf("active-only --where status=backlog must include top-placed:\n%s", out)
		}
		if strings.Contains(out, "arch-placed") {
			t.Fatalf("active-only --where status=backlog must NOT include arch-placed:\n%s", out)
		}
	})

	t.Run("active-plus-archived", func(t *testing.T) {
		out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--where", "status=backlog", "--archived")
		if code != 0 {
			t.Fatalf("exit=%d, want 0 (err=%q)", code, errOut)
		}
		if !strings.Contains(out, "top-placed") {
			t.Fatalf("--where status=backlog --archived must include top-placed:\n%s", out)
		}
		if !strings.Contains(out, "arch-placed") {
			t.Fatalf("--where status=backlog --archived must include arch-placed:\n%s", out)
		}
	})
}
