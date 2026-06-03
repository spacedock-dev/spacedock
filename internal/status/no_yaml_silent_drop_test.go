// ABOUTME: TestNoSilentYAMLValueDrop pins the writer/reader contract: a
// ABOUTME: non-empty raw frontmatter value must not decode to "" via yaml.v3.
package status

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoSilentYAMLValueDrop walks every index.md fixture under testdata/,
// extracts each top-level frontmatter line's raw post-`:` substring, parses the
// file through ParseFrontmatter, and asserts that no key whose raw value is
// non-empty (and not an explicit empty quoted scalar `""` / `''`) decodes to
// the empty string. The canonical failure is `pr: #N` — yaml.v3 treats `#N` as
// a comment and drops the value silently. The writer-quoting policy in
// mutate.go (needsExplicitQuoting + setScalarValue) prevents this on every
// write; this test catches a regression of that policy on existing fixtures.
func TestNoSilentYAMLValueDrop(t *testing.T) {
	root := "testdata"
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(path) != "index.md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !contentHasOpeningFence(data) {
			return nil
		}
		raw := rawTopLevelValues(data)
		parsed := parseFrontmatterContent(data)
		for key, rawVal := range raw {
			if rawVal == "" {
				continue
			}
			if rawVal == `""` || rawVal == "''" {
				// Explicit empty quoted scalar — yaml.v3 legitimately returns "".
				continue
			}
			if parsed[key] == "" {
				t.Errorf("%s: key %q has raw value %q but decodes to \"\" (yaml.v3 silently dropped it)", path, key, rawVal)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk testdata: %v", err)
	}
}

// rawTopLevelValues scans the frontmatter slice of data and returns the raw
// post-`:` substring (trimmed of leading/trailing whitespace) for each
// top-level (non-indented) `key: value` line. Indented lines (nested mappings)
// are skipped, matching parseFrontmatterContent's top-level-only surface.
func rawTopLevelValues(data []byte) map[string]string {
	out := map[string]string{}
	slice := frontmatterSlice(data)
	if len(slice) == 0 {
		return out
	}
	for _, line := range strings.Split(string(slice), "\n") {
		if line == "" {
			continue
		}
		// Skip indented lines (nested mapping body).
		if line[0] == ' ' || line[0] == '\t' {
			continue
		}
		idx := strings.Index(line, ":")
		if idx <= 0 {
			continue
		}
		key := line[:idx]
		val := strings.TrimSpace(line[idx+1:])
		out[key] = val
	}
	return out
}
