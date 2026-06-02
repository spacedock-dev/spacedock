// ABOUTME: AC-6 native validation parity — each defect class (dup/bad/missing
// ABOUTME: id, flat/folder conflict, bad stage name) matches the oracle's lines.
package status

import (
	"os"
	"path/filepath"
	"testing"
)

// validationFixture builds a workflow in a temp dir from a README and a set of
// rel-path entity files, git-inits it, and returns the root.
func validationFixture(t *testing.T, readme string, files map[string]string) string {
	t.Helper()
	dst := t.TempDir()
	if err := os.WriteFile(filepath.Join(dst, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	for rel, content := range files {
		p := filepath.Join(dst, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitInit(t, dst)
	return dst
}

const seqREADME = `---
entity-type: task
id-style: sequential
stages:
  states:
    - name: backlog
      initial: true
    - name: done
      terminal: true
---

# Seq Workflow
`

func ent(id, status string) string {
	return "---\nid: " + id + "\ntitle: T\nstatus: " + status + "\nscore: \"0.5\"\nsource: x\n---\n# T\n"
}

func TestNativeValidationParity(t *testing.T) {
	cases := []struct {
		name   string
		readme string
		files  map[string]string
	}{
		{
			name:   "valid",
			readme: seqREADME,
			files:  map[string]string{"a.md": ent(`"001"`, "backlog"), "b.md": ent(`"002"`, "done")},
		},
		{
			name:   "missing-id",
			readme: seqREADME,
			files:  map[string]string{"a.md": "---\ntitle: T\nstatus: backlog\n---\n# T\n"},
		},
		{
			name:   "non-numeric-id",
			readme: seqREADME,
			files:  map[string]string{"a.md": ent("abc", "backlog")},
		},
		{
			name:   "duplicate-id",
			readme: seqREADME,
			files:  map[string]string{"a.md": ent(`"001"`, "backlog"), "b.md": ent(`"001"`, "done")},
		},
		{
			name:   "flat-folder-conflict",
			readme: seqREADME,
			files: map[string]string{
				"dup.md":       ent(`"001"`, "backlog"),
				"dup/index.md": ent(`"002"`, "done"),
			},
		},
		{
			name: "bad-stage-name",
			readme: `---
id-style: sequential
stages:
  states:
    - name: Bad_Stage
      initial: true
    - name: done
      terminal: true
---
# Bad stage
`,
			files: map[string]string{"a.md": ent(`"001"`, "done")},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			env := pinnedEnv(t)
			nativeRoot := validationFixture(t, tc.readme, tc.files)

			nArgs := []string{"--workflow-dir", nativeRoot, "--validate"}
			nOut, nErr, nCode := runNative(t, nativeRoot, env, nArgs...)

			assertEnvelopeGolden(t, "native-validate-"+tc.name, goldenEnvelope{
				stdout: normalize(nOut, nativeRoot), stderr: normalize(nErr, nativeRoot), exit: nCode,
			})
		})
	}
}

// TestNativeValidationGatesReads locks that a defect rejects an enumerate op
// (default table) globally — native exits 1 like the oracle — proving the read-
// path id-strictness rule.
func TestNativeValidationGatesReads(t *testing.T) {
	env := pinnedEnv(t)
	files := map[string]string{"a.md": "---\ntitle: T\nstatus: backlog\n---\n# T\n"}
	nativeRoot := validationFixture(t, seqREADME, files)

	nOut, nErr, nCode := runNative(t, nativeRoot, env, "--workflow-dir", nativeRoot)

	if nCode != 1 {
		t.Fatalf("default table over an id-less workflow must exit 1: native=%d", nCode)
	}
	if nOut != "" {
		t.Fatalf("stdout must be empty on validation failure: native=%q", nOut)
	}
	assertEnvelopeGolden(t, "native-validate-gates-reads", goldenEnvelope{
		stdout: normalize(nOut, nativeRoot), stderr: normalize(nErr, nativeRoot), exit: nCode,
	})
}
