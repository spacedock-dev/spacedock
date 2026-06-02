// ABOUTME: AC-1 migration check — every fixture frontmatter parses through
// ABOUTME: yaml.v3 and the reader's value-map matches a direct decode.
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestMigrationCheckFixturesParseConsistently is the AC-1 divergence guard:
// it walks every `*.md` under internal/status/testdata AND the repo-root
// workflow-definition files (agent prompts, skill scaffolding, hook mods,
// dev workflow READMEs), parses each frontmatter block via the public
// reader (ParseFrontmatter, yaml.v3-backed), then independently decodes
// the in-fence slice through yaml.v3 directly, and asserts the two
// value-maps agree key-by-key. A disagreement would be silent drift
// between the reader and the YAML standard — exactly the failure mode the
// migration check is meant to surface.
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
	walk := func(root string, path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
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
	if checked == 0 {
		t.Fatal("walked no frontmatters — repo layout changed?")
	}
	t.Logf("migration-check verified %d frontmatters across fixtures + live", checked)
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
