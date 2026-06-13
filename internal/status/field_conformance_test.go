// ABOUTME: Schema-driven field-conformance tests — warn-tier field violations
// ABOUTME: surface on stderr, exit 0, never gate reads; schema bytes are the source.
package status

import (
	"strings"
	"testing"

	spacedock "github.com/spacedock-dev/spacedock"
	"gopkg.in/yaml.v3"
)

// entField builds a sequential-style entity with an extra frontmatter line so a
// single bad field can be planted while id/status stay clean.
func entField(id, status, extraKey, extraVal string) string {
	fm := "---\nid: " + id + "\ntitle: T\nstatus: " + status + "\nscore: \"0.5\"\nsource: x\n"
	if extraKey != "" {
		fm += extraKey + ": " + extraVal + "\n"
	}
	return fm + "---\n# T\n"
}

// TestFieldConformanceWarnsSurface locks AC-1: each per-field schema violation
// the entity schema declares produces a field-named diagnostic on stderr under
// `--validate`. The bad value lives in the fixture; the rule is the schema's.
func TestFieldConformanceWarnsSurface(t *testing.T) {
	env := pinnedEnv(t)
	cases := []struct {
		name      string
		files     map[string]string
		wantField string // the field name the diagnostic must mention
	}{
		{
			name:      "mod-block-no-colon",
			files:     map[string]string{"a.md": entField(`"001"`, "backlog", "mod-block", "noColonHere")},
			wantField: "mod-block",
		},
		{
			name:      "verdict-out-of-enum",
			files:     map[string]string{"a.md": entField(`"001"`, "backlog", "verdict", "MAYBE")},
			wantField: "verdict",
		},
		{
			name:      "score-not-numeric",
			files:     map[string]string{"a.md": entField(`"001"`, "backlog", "score", "notanumber")},
			wantField: "score",
		},
		{
			name:      "started-malformed-iso",
			files:     map[string]string{"a.md": entField(`"001"`, "backlog", "started", "not-a-date")},
			wantField: "started",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			root := validationFixture(t, seqREADME, tc.files)
			_, nErr, nCode := runNative(t, root, env, "--workflow-dir", root, "--validate")
			if nCode != 0 {
				t.Fatalf("warn-only fixture must exit 0, got %d (stderr=%q)", nCode, nErr)
			}
			// The diagnostic must carry the Warning: prefix bound to the field —
			// a field-conformance finding is advisory, not a gating Error:. This
			// pins the warn-vs-error signal the FO reads; flipping the line to
			// Error: must red this assertion.
			wantLine := "Warning: field '" + tc.wantField
			if !strings.Contains(nErr, wantLine) {
				t.Fatalf("stderr missing warn line %q, got %q", wantLine, nErr)
			}
			if strings.Contains(nErr, "Error: field '"+tc.wantField) {
				t.Fatalf("field-conformance finding misreported as Error: (must be Warning:), got %q", nErr)
			}
		})
	}
}

// TestFieldConformanceWarnsDoNotGateExit locks AC-2: a warn-only workflow exits
// 0 with non-empty warning stderr; a structural error (dup id) still exits 1.
func TestFieldConformanceWarnsDoNotGateExit(t *testing.T) {
	env := pinnedEnv(t)

	t.Run("warn-only-exits-0", func(t *testing.T) {
		files := map[string]string{
			"a.md": entField(`"001"`, "backlog", "verdict", "MAYBE"),
			"b.md": entField(`"002"`, "done", "mod-block", "noColon"),
		}
		root := validationFixture(t, seqREADME, files)
		nOut, nErr, nCode := runNative(t, root, env, "--workflow-dir", root, "--validate")
		if nCode != 0 {
			t.Fatalf("warn-only --validate must exit 0, got %d", nCode)
		}
		if strings.TrimSpace(nErr) == "" {
			t.Fatalf("warn-only --validate must print warnings to stderr, got empty")
		}
		if strings.TrimSpace(nOut) != "VALID" {
			t.Fatalf("warn-only --validate stdout must stay VALID, got %q", nOut)
		}
	})

	t.Run("structural-error-exits-1", func(t *testing.T) {
		files := map[string]string{
			"a.md": entField(`"001"`, "backlog", "verdict", "MAYBE"),
			"b.md": entField(`"001"`, "done", "", ""), // dup id == structural error
		}
		root := validationFixture(t, seqREADME, files)
		_, _, nCode := runNative(t, root, env, "--workflow-dir", root, "--validate")
		if nCode != 1 {
			t.Fatalf("structural error (dup id) must exit 1, got %d", nCode)
		}
	})
}

// TestFieldConformanceWarnsDoNotGateReads locks AC-3: a warn-tier field
// violation does not gate the default-table read path — the FO is never locked
// out by a field warning. Mirrors TestNativeValidationGatesReads inverted.
func TestFieldConformanceWarnsDoNotGateReads(t *testing.T) {
	env := pinnedEnv(t)
	files := map[string]string{"a.md": entField(`"001"`, "backlog", "verdict", "MAYBE")}
	root := validationFixture(t, seqREADME, files)

	nOut, _, nCode := runNative(t, root, env, "--workflow-dir", root)
	if nCode != 0 {
		t.Fatalf("default table over a warn-violation fixture must exit 0, got %d", nCode)
	}
	if strings.TrimSpace(nOut) == "" {
		t.Fatalf("default table must still print the table on stdout, got empty")
	}
}

// TestFieldConformanceSchemaDriven locks AC-4: the validator's rule and the
// schema file agree. The test loads the embedded schema bytes independently,
// parses out the mod-block pattern, and asserts the validator's loaded pattern
// equals it — two independently-readable values that could diverge.
func TestFieldConformanceSchemaDriven(t *testing.T) {
	// Independently parse the embedded schema for the mod-block pattern.
	var doc struct {
		Frontmatter struct {
			Fields map[string]struct {
				Pattern      string   `yaml:"pattern"`
				Conventional []string `yaml:"conventional"`
			} `yaml:"fields"`
		} `yaml:"frontmatter"`
	}
	if err := yaml.Unmarshal(spacedock.EntityMDSchema, &doc); err != nil {
		t.Fatalf("parse embedded schema: %v", err)
	}
	wantPattern := doc.Frontmatter.Fields["mod-block"].Pattern
	if wantPattern == "" {
		t.Fatalf("embedded schema has no mod-block pattern")
	}

	schema := loadEntitySchema()
	gotPattern := schema.fields["mod-block"].Pattern
	if gotPattern != wantPattern {
		t.Fatalf("validator mod-block pattern %q != schema file pattern %q", gotPattern, wantPattern)
	}

	// Same for the verdict conventional enum.
	wantEnum := doc.Frontmatter.Fields["verdict"].Conventional
	gotEnum := schema.fields["verdict"].Conventional
	if strings.Join(gotEnum, ",") != strings.Join(wantEnum, ",") {
		t.Fatalf("validator verdict enum %v != schema file enum %v", gotEnum, wantEnum)
	}
}
