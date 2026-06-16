// ABOUTME: SOURCE column is opt-in — the default human table and bare --json omit
// ABOUTME: it; --fields source / --all-fields restore it, asserted by header/key tokens.
package status

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

// sourceWorkflowRoot returns the abs path to the seq-workflow fixture used by
// every SOURCE-conditional assertion.
func sourceWorkflowRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// TestSourceOmittedFromDefaultTables (AC-1) asserts the fixed human table carries
// no SOURCE column for the bare default, --archived, and --where reads. SOURCE is
// no longer a fixed column; it is opt-in only.
func TestSourceOmittedFromDefaultTables(t *testing.T) {
	root := sourceWorkflowRoot(t)
	env := pinnedEnv(t)

	cases := []struct {
		name  string
		extra []string
	}{
		{"default", nil},
		{"archived", []string{"--archived"}},
		{"where", []string{"--where", "status=ideation"}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"--workflow-dir", root}, tc.extra...)
			out, errOut, code := runNative(t, root, env, args...)
			if code != 0 {
				t.Fatalf("native exit=%d stderr=%q", code, errOut)
			}
			header := headerOf(out)
			if c := countToken(header, "SOURCE"); c != 0 {
				t.Fatalf("%s header carries %d SOURCE columns, want 0\nheader: %q", tc.name, c, header)
			}
		})
	}
}

// TestFieldsSourceRestoresColumn (AC-2) asserts `--fields source` renders exactly
// one SOURCE column (now appended as an extra, since source left the fixed set),
// carrying each entity's value.
func TestFieldsSourceRestoresColumn(t *testing.T) {
	root := sourceWorkflowRoot(t)
	env := pinnedEnv(t)
	out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--fields", "source")
	if code != 0 {
		t.Fatalf("native exit=%d stderr=%q", code, errOut)
	}
	header := headerOf(out)
	if c := countToken(header, "SOURCE"); c != 1 {
		t.Fatalf("--fields source header has %d SOURCE columns, want exactly 1\nheader: %q", c, header)
	}
	// The restored column must carry a real value, not an empty extra.
	if !containsToken(out, "roadmap") {
		t.Fatalf("--fields source did not render any source value:\n%s", out)
	}
}

// TestAllFieldsSurfacesSourceAsExtra (AC-2) asserts --all-fields surfaces SOURCE
// as a single extra column (no longer a fixed slot).
func TestAllFieldsSurfacesSourceAsExtra(t *testing.T) {
	root := sourceWorkflowRoot(t)
	env := pinnedEnv(t)
	out, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--all-fields")
	if code != 0 {
		t.Fatalf("native exit=%d stderr=%q", code, errOut)
	}
	header := headerOf(out)
	if c := countToken(header, "SOURCE"); c != 1 {
		t.Fatalf("--all-fields header has %d SOURCE columns, want exactly 1\nheader: %q", c, header)
	}
}

// TestBareJSONOmitsSource (AC-3, direction A) asserts a bare --json entity object
// carries no source key, while --json --fields id,source carries it.
func TestBareJSONOmitsSource(t *testing.T) {
	root := sourceWorkflowRoot(t)
	env := pinnedEnv(t)

	bare, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--json")
	if code != 0 {
		t.Fatalf("bare --json exit=%d stderr=%q", code, errOut)
	}
	for _, e := range parseStatusEntities(t, bare) {
		if _, ok := e["source"]; ok {
			t.Fatalf("bare --json entity carries a source key, want none: %v", e)
		}
	}

	named, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--json", "--fields", "id,source")
	if code != 0 {
		t.Fatalf("--json --fields id,source exit=%d stderr=%q", code, errOut)
	}
	ents := parseStatusEntities(t, named)
	if len(ents) == 0 {
		t.Fatal("--json --fields id,source produced no entities")
	}
	for _, e := range ents {
		if _, ok := e["source"]; !ok {
			t.Fatalf("--json --fields id,source entity missing source key: %v", e)
		}
	}
}

// containsToken reports whether any whitespace-separated token of s equals tok.
func containsToken(s, tok string) bool {
	return countToken(s, tok) > 0
}

// parseStatusEntities parses a --json status envelope into its entity objects.
func parseStatusEntities(t *testing.T, out string) []map[string]string {
	t.Helper()
	var env statusEnvelope
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("parse --json: %v\n%s", err, out)
	}
	return env.Entities
}
