// ABOUTME: AC-1 live-entities-parse check — every fixture + live frontmatter
// ABOUTME: parses through yaml.v3; helper-level corruption tests cover the seam.
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestMigrationCheckFixturesParseConsistently is the AC-1 live-entities-parse
// guard. It walks every `*.md` under internal/status/testdata AND the
// repo-root workflow-definition files (agent prompts, skill scaffolding, hook
// mods, dev workflow READMEs), parses each frontmatter block via the public
// reader (ParseFrontmatter, yaml.v3-backed), then independently decodes the
// in-fence slice through yaml.v3 directly, and asserts the two value-maps
// agree key-by-key on top-level scalars.
//
// SCOPE NOTE (cycle 1 audit, finding L1-b): both the reader and the
// independent decode here go through yaml.v3, so this is a yaml.v3-against-
// yaml.v3 INTRA-VERSION CONSISTENCY check — NOT a parser-migration parity
// oracle against the retired Python reader. AC-1 claims "live entities still
// parse," which this asserts: every live + fixture frontmatter parses cleanly
// through yaml.v3 and the reader's value-map is consistent with a direct
// yaml.v3 decode against the same fence-extracted bytes. The
// fence-extraction chain itself (`contentHasOpeningFence`,
// `frontmatterSlice`, `splitLines`, `normalizeNewlines`) is the upstream
// helper-level pipeline; bugs in those helpers would affect BOTH sides of
// this consistency check, so they are covered by the corruption-sensitive
// unit tests below (TestFrontmatterSliceHelperCorruption,
// TestSplitLinesNormalizeHelperCorruption) rather than by this walk.
//
// Accepted exceptions: files with NO opening fence (body-only markdown,
// READMEs without frontmatter) are skipped, not failed. The spike found
// exactly one colon-space hazard in the live corpus (an archived entity),
// which is documented as divergence #1 (writer-quote going forward); the
// fixture corpus is curated, so we expect no accepted-exception entries
// here.
func TestMigrationCheckFixturesParseConsistently(t *testing.T) {
	testdata, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	// Repo-root workflow-definition trees: agent definitions, skill
	// scaffolding, hook mods, and dev workflow definitions. These are the
	// live frontmatters the reader sees outside of fixture corpora — the
	// migration check covers them so a YAML quirk in a hand-edited live
	// entity is caught here, not in production.
	roots := []string{
		testdata,
		filepath.Join(repoRoot, "agents"),
		filepath.Join(repoRoot, "skills"),
		filepath.Join(repoRoot, "mods"),
		filepath.Join(repoRoot, "docs"),
	}
	var checked int
	var visited []string
	walk := func(root string, path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if action, dir := migrationCheckWalkDir(path, info); dir {
			return action
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		visited = append(visited, path)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			return nil
		}
		if !contentHasOpeningFence(data) {
			return nil
		}
		rel, _ := filepath.Rel(repoRoot, path)
		checked++

		// 1. The public reader returns a value-map.
		got := parseFrontmatterContent(data)

		// 2. yaml.v3 directly decodes the in-fence slice.
		slice := frontmatterSlice(data)
		if len(slice) == 0 {
			// An entity with an opening fence but empty body — both readers
			// must agree on empty.
			if len(got) != 0 {
				t.Errorf("%s: empty-body slice but reader returned %v", rel, got)
			}
			return nil
		}
		var direct map[string]any
		if err := yaml.Unmarshal(slice, &direct); err != nil {
			t.Errorf("%s: yaml.v3 direct unmarshal failed: %v", rel, err)
			return nil
		}

		// 3. Agree key-by-key on TOP-LEVEL SCALARS (the reader surface).
		for k, v := range got {
			raw, ok := direct[k]
			if !ok {
				t.Errorf("%s: reader has key %q but direct decode does not", rel, k)
				continue
			}
			want := scalarString(raw)
			if want != v {
				t.Errorf("%s: key %q reader=%q direct=%q", rel, k, v, want)
			}
		}
		// Every top-level direct key must be visible to the reader (as scalar
		// or as the indented-lines-ignored empty string for nested values).
		for k := range direct {
			if _, ok := got[k]; !ok {
				t.Errorf("%s: direct decode has key %q but reader does not", rel, k)
			}
		}
		return nil
	}
	for _, root := range roots {
		if _, statErr := os.Stat(root); statErr != nil {
			continue
		}
		root := root
		if err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
			return walk(root, path, info, walkErr)
		}); err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
	}
	// Positive prune proof on the live corpus: on a dev machine the split-root
	// state checkout materializes under docs/dev/.spacedock-state and the walk
	// roots include docs/. If the prune regressed, the ~100 machine-local
	// entity files there would show up in visited. Assert no visited path
	// crosses the state tree. (On CI the tree is absent, so this is a no-op
	// there; the hermetic temp-fixture proof lives in
	// TestMigrationCheckPrunesStateTree.)
	stateSep := string(filepath.Separator) + ".spacedock-state" + string(filepath.Separator)
	for _, p := range visited {
		if strings.Contains(p, stateSep) {
			t.Errorf("walk descended into the pruned .spacedock-state tree: %s", p)
		}
	}
	// Non-vacuous guard: the walk MUST still validate real entity files so the
	// consistency check above can actually fire. A prune that dropped checked to
	// zero would silently disable the migration check.
	if checked == 0 {
		t.Fatal("walked no frontmatters — repo layout changed or the prune over-reached?")
	}
	t.Logf("migration-check verified %d frontmatters across fixtures + live", checked)
}

// TestMigrationCheckPrunesStateTree is the hermetic, CI-stable positive proof
// that the migration-check walk prunes the non-entity trees wholesale rather
// than name-matching individual subdirs. It plants a temp tree holding a real
// entity alongside both pruned subtrees — a .spacedock-state checkout and a
// docs/roadmap strategy tree — drives filepath.Walk through the same
// migrationCheckWalkDir step the production migration check uses, and asserts
// neither pruned subtree is visited while the real entity is — which keeps the
// checked>0 guard in TestMigrationCheckFixturesParseConsistently non-vacuous
// (the prune does not also swallow legitimate entities).
func TestMigrationCheckPrunesStateTree(t *testing.T) {
	root := t.TempDir()

	// A real entity outside the state tree — the walk MUST visit this one.
	liveDir := filepath.Join(root, "live")
	if err := os.MkdirAll(liveDir, 0o755); err != nil {
		t.Fatal(err)
	}
	liveEntity := filepath.Join(liveDir, "entity.md")
	if err := os.WriteFile(liveEntity, []byte("---\nid: live-1\nstatus: done\n---\n\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// A machine-local entity inside the pruned state tree — the walk must NOT
	// visit this one. Its frontmatter is a colon-space hazard that the live
	// migration-check consistency walk would flag if it ever read it, which is
	// exactly why the state tree is pruned.
	stateDir := filepath.Join(root, "docs", "dev", ".spacedock-state")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		t.Fatal(err)
	}
	poison := filepath.Join(stateDir, "index.md")
	if err := os.WriteFile(poison, []byte("---\nsource: a: b: c\n---\n\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// A session debrief inside the pruned docs/roadmap strategy tree — the walk
	// must NOT visit this one either. Its bare-YAML session-date scalar decodes as
	// time.Time directly but as a string through the reader, so the consistency
	// walk would flag it if it ever read it. The roadmap tree is the strategy
	// layer (non-entity by design), pruned wholesale for the same reason.
	roadmapDir := filepath.Join(root, "docs", "roadmap", "0198-pre-flip-hardening")
	if err := os.MkdirAll(roadmapDir, 0o755); err != nil {
		t.Fatal(err)
	}
	debrief := filepath.Join(roadmapDir, "debrief.md")
	if err := os.WriteFile(debrief, []byte("---\nsession-date: 2026-06-08\n---\n\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var visited []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if action, dir := migrationCheckWalkDir(path, info); dir {
			return action
		}
		if strings.HasSuffix(path, ".md") {
			visited = append(visited, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}

	visitedLive := false
	for _, p := range visited {
		if p == poison {
			t.Errorf("walk descended into the pruned .spacedock-state tree: visited %s", p)
		}
		if p == debrief {
			t.Errorf("walk descended into the pruned docs/roadmap tree: visited %s", p)
		}
		if p == liveEntity {
			visitedLive = true
		}
	}
	if !visitedLive {
		t.Fatalf("walk did not visit the real entity %s — the prune over-reached", liveEntity)
	}
}

// scalarString renders a yaml.v3-decoded scalar back to the string the
// reader returns: a nil/empty becomes "", a map or slice becomes "" (the
// nested-lines-ignored semantic), anything else is its string form.
func scalarString(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case map[string]any, []any:
		return ""
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		// Numbers, etc. — render via yaml.v3 itself for consistency.
		out, err := yaml.Marshal(v)
		if err != nil {
			return ""
		}
		s := strings.TrimRight(string(out), "\n")
		return s
	}
}
