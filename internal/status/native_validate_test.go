// ABOUTME: AC-6 native validation parity — each defect class (dup/bad/missing
// ABOUTME: id, flat/folder conflict, bad stage name) matches the oracle's lines.
package status

import (
	"os"
	"path/filepath"
	"strings"
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

// TestValidateFlagsSelfRefACs locks AC-2: under `require-external-proof: true`,
// `spacedock status --validate` emits one standard `entityEvidence`-shaped line
// per flagged AC and exits 1; a workflow with no self-referential ACs returns
// VALID exit 0. Drives the testdata/external-proof-workflow/ fixture.
func TestValidateFlagsSelfRefACs(t *testing.T) {
	env := pinnedEnv(t)

	t.Run("flags-self-ref", func(t *testing.T) {
		root := stageExternalProofFixture(t, "require-external-proof: true\n")

		_, nErr, nCode := runNative(t, root, env, "--workflow-dir", root, "--validate")
		if nCode != 1 {
			t.Fatalf("--validate must exit 1 when self-ref ACs exist, got %d (stderr=%q)", nCode, nErr)
		}
		wantSubstr := "Error: self-referential AC proof (AC-1):"
		if !strings.Contains(nErr, wantSubstr) {
			t.Fatalf("stderr missing evidence line %q, got %q", wantSubstr, nErr)
		}
		for _, want := range []string{"workflow=", "scope=active", "slug=010-self-ref-only", "id=010", "path="} {
			if !strings.Contains(nErr, want) {
				t.Fatalf("stderr missing %q in evidence line, got %q", want, nErr)
			}
		}
	})

	t.Run("real-proof-only-valid", func(t *testing.T) {
		root := stageExternalProofFixture(t, "require-external-proof: true\n")
		// Delete the self-ref + force entities so only the real-proof one remains.
		for _, name := range []string{"010-self-ref-only.md", "030-force-bypass.md"} {
			if err := os.Remove(filepath.Join(root, name)); err != nil {
				t.Fatal(err)
			}
		}

		nOut, _, nCode := runNative(t, root, env, "--workflow-dir", root, "--validate")
		if nCode != 0 {
			t.Fatalf("--validate should return VALID with no self-ref ACs, got exit=%d", nCode)
		}
		if strings.TrimSpace(nOut) != "VALID" {
			t.Fatalf("stdout should be VALID, got %q", nOut)
		}
	})
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
