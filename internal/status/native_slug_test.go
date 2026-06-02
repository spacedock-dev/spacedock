// ABOUTME: id-style: slug parity — default table id==slug, --short-id prints
// ABOUTME: the slug, and --next-id is not applicable (exit 1), all vs the oracle.
package status

import (
	"testing"
)

func slugFixture(t *testing.T) string {
	t.Helper()
	readme := "---\nid-style: slug\nstages:\n  states:\n    - name: backlog\n      initial: true\n    - name: done\n      terminal: true\n---\n# Slug WF\n"
	alpha := "---\ntitle: A\nstatus: backlog\nscore: \"0.5\"\nsource: x\n---\n# A\n"
	beta := "---\ntitle: B\nstatus: done\nscore: \"0.7\"\nsource: y\n---\n# B\n"
	dst := t.TempDir()
	writeAll(t, dst, map[string]string{
		"README.md": readme,
		"alpha.md":  alpha,
		"beta.md":   beta,
	})
	gitInit(t, dst)
	return dst
}

func TestNativeSlugStyleParity(t *testing.T) {
	env := pinnedEnv(t)
	cases := []struct {
		name  string
		extra []string
	}{
		{"default", nil},
		{"short-id", []string{"--short-id", "alpha"}},
		{"resolve", []string{"--resolve", "alpha"}},
		{"next", []string{"--next"}},
		{"validate", []string{"--validate"}},
		{"next-id-not-applicable", []string{"--next-id"}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			nativeRoot := slugFixture(t)
			nArgs := append([]string{"--workflow-dir", nativeRoot}, tc.extra...)
			nOut, nErr, nCode := runNative(t, nativeRoot, env, nArgs...)
			assertEnvelopeGolden(t, "native-slug-"+tc.name, goldenEnvelope{
				stdout: normalize(nOut, nativeRoot), stderr: normalize(nErr, nativeRoot), exit: nCode,
			})
		})
	}
}

// TestNativeSlugDuplicateValidation locks slug-style duplicate-id validation: two
// entities with the same slug is impossible on disk, so the duplicate-effective-
// id path is exercised via an archived sibling sharing the active slug.
func TestNativeSlugDuplicateValidation(t *testing.T) {
	env := pinnedEnv(t)
	readme := "---\nid-style: slug\nstages:\n  states:\n    - name: backlog\n      initial: true\n    - name: done\n      terminal: true\n---\n# Slug WF\n"
	ent := "---\ntitle: T\nstatus: backlog\nscore: \"0.5\"\nsource: x\n---\n# T\n"
	dst := t.TempDir()
	writeAll(t, dst, map[string]string{
		"README.md":       readme,
		"dup.md":          ent,
		"_archive/dup.md": ent, // archived sibling shares the slug
	})
	gitInit(t, dst)
	nativeRoot := dst

	nOut, nErr, nCode := runNative(t, nativeRoot, env, "--workflow-dir", nativeRoot, "--validate")

	if nCode != 1 {
		t.Fatalf("exit: native=%d, want 1 (duplicate effective id)", nCode)
	}
	assertEnvelopeGolden(t, "native-slug-dup-validation", goldenEnvelope{
		stdout: normalize(nOut, nativeRoot), stderr: normalize(nErr, nativeRoot), exit: nCode,
	})
	if nOut != "" {
		t.Fatalf("stdout must be empty on validation failure: native=%q", nOut)
	}
}
