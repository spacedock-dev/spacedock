// ABOUTME: Cycle-2 parity fixtures — abs-worktree join, str.splitlines separator
// ABOUTME: set, numeric schema_version, and non-object stdin, each vs the oracle.
package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuildAbsoluteWorktreeParity locks fix 1: an ABSOLUTE worktree: frontmatter
// value must resolve via os.path.join semantics (absolute component wins), not
// filepath.Join (which doubles the path under gitRoot). The FO stamps absolute
// worktree: values on live entities, so this is the happy path. Before the fix
// native rejected with a doubled-path "does not exist"; the oracle accepts.
func TestBuildAbsoluteWorktreeParity(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))

	// An ABSOLUTE worktree path (the shape the FO stamps), existing on disk.
	absWorktree := filepath.Join(root, ".worktrees", "spacedock-ensign-thing")
	if err := os.MkdirAll(absWorktree, 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", absWorktree))
	gitInit(t, root)

	stdin := strings.Join([]string{"- a"}, "\n")
	stdinArgs := []string{"build", "--workflow-dir", root, "--entity-path", entityPath, "--stage", "implementation", "--checklist-file", "-"}

	native := runNative(stdin, stdinArgs...)
	nativeBody := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	env := goldenEnvelope{res: normRun(native, root, home), body: normPaths(nativeBody, root, home)}
	assertGolden(t, "build-abs-worktree", env)
	// Explicit: native must exit 0 (full dispatch), not the doubled-path error.
	if native.exit != 0 {
		t.Errorf("abs-worktree native exit=%d, want 0; stderr=%q", native.exit, native.stderr)
	}
	if strings.Contains(native.stderr, "does not exist") {
		t.Errorf("abs-worktree native still doubles the path: %q", native.stderr)
	}
}

// TestShowStageDefSeparatorParity locks fix 2: show-stage-def's line splitting
// matches Python str.splitlines() across the full separator set, not just
// \r\n/\r/\n. Each case embeds one separator inside a ### subsection body; the
// oracle breaks the line there and the native must produce byte-identical
// stdout.
func TestShowStageDefSeparatorParity(t *testing.T) {
	// name -> the separator rune to embed between two body words.
	separators := []struct {
		name string
		sep  string
	}{
		{"VT", "\v"},
		{"FF", "\f"},
		{"FS", "\x1c"},
		{"GS", "\x1d"},
		{"RS", "\x1e"},
		{"NEL", "\u0085"},
		{"LS", "\u2028"},
		{"PS", "\u2029"},
	}
	for _, sc := range separators {
		t.Run(sc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			root := t.TempDir()
			// A README whose ideation body has two words split by the separator.
			readme := "---\nentity-type: task\nid-style: slug\n---\n# Sep\n\n### ideation\n\nalpha" + sc.sep + "beta gamma.\n\n### done\n\nterm.\n"
			writeFile(t, filepath.Join(root, "README.md"), readme)

			native := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", "ideation")
			assertGolden(t, "showstagedef-separator-"+sc.name, goldenEnvelope{res: normRun(native, root, home)})
		})
	}
}
