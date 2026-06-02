// ABOUTME: show-stage-def parity vs the oracle across well-formed, decorated,
// ABOUTME: malformed, and missing headings; plus the deferred-subcommand guard.
package dispatch

import (
	"os"
	"path/filepath"
	"testing"
)

// readmeHeadings is a README whose ### subsections exercise the extractor: a
// plain heading, a decorated heading (backticks + parenthetical), and a heading
// whose stage name is NOT the first token (malformed). The frontmatter stages
// block is irrelevant to show-stage-def (it reads the body sections directly).
const readmeHeadings = `---
entity-type: task
id-style: slug
---
# Headings Fixture

### ideation

Ideation body line one.

- **Inputs:** the seed.
- **Outputs:** the design.

### ` + "`implementation`" + ` *(captain-interactive)*

Implementation body.

- **Outputs:** the deliverable.

### review of the validation stage

This heading mentions validation as a non-first token.

### done

Terminal.
`

// TestShowStageDefParity diffs native show-stage-def against the oracle across
// the well-formed, decorated-heading, malformed-heading, and missing-stage
// cases. The malformed-heading diagnostic is a hand-built f-string (not a
// Python str(e)), so it is byte-reproducible and asserted byte-for-byte.
func TestShowStageDefParity(t *testing.T) {
	cases := []struct {
		name  string
		stage string
	}{
		{"well-formed", "ideation"},
		{"decorated-heading", "implementation"},
		{"malformed-heading", "validation"},
		{"missing-stage", "nonesuch"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			root := t.TempDir()
			writeFile(t, filepath.Join(root, "README.md"), readmeHeadings)

			native := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", tc.stage)
			assertGolden(t, "showstagedef-"+tc.name, goldenEnvelope{res: normRun(native, root, home)})
		})
	}
}

// TestShowStageDefMissingReadme locks the workflow-dir / README-not-found
// diagnostics against the oracle.
func TestShowStageDefMissingReadme(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	// A workflow dir that does not exist.
	missing := filepath.Join(root, "nope")
	native := runNative("", "show-stage-def", "--workflow-dir", missing, "--stage", "ideation")
	assertGolden(t, "showstagedef-missing-dir", goldenEnvelope{res: normRun(native, root, home)})

	// A dir that exists but has no README.
	noReadme := filepath.Join(root, "empty")
	if err := os.MkdirAll(noReadme, 0o755); err != nil {
		t.Fatal(err)
	}
	native2 := runNative("", "show-stage-def", "--workflow-dir", noReadme, "--stage", "ideation")
	assertGolden(t, "showstagedef-no-readme", goldenEnvelope{res: normRun(native2, root, home)})
}
