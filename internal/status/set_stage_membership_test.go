// ABOUTME: AC-3 --set status= membership parity — a non-member stage value is
// ABOUTME: rejected (exit 1, unchanged frontmatter), a member is accepted, launcher vs oracle.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestSetStatusNonMemberRejected locks the #189 membership check: a `--set
// status=zzz` where zzz is not a declared stage in the workflow's
// stages.states[].name exits non-zero with an actionable error listing the known
// stages, and leaves the entity frontmatter UNCHANGED. Asserted launcher vs
// oracle so the parity suite covers the new guard.
func TestSetStatusNonMemberRejected(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=zzz")
	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 1 {
		t.Fatalf("native exit=%d, want 1 (non-member status must reject)", nCode)
	}
	// The error names the unknown value and lists the known stages.
	for _, want := range []string{"zzz", "backlog", "ideation", "implementation", "done"} {
		if !strings.Contains(nErr, want) {
			t.Fatalf("native stderr = %q, want it to mention %q", nErr, want)
		}
	}
	assertEnvelopeGolden(t, "set-membership-nonmember-rejected", goldenEnvelope{
		stdout: normalize(nOut, root), stderr: normalize(nErr, root), exit: nCode,
	})
	if nOut != "" {
		t.Fatalf("stdout must be empty on rejection: native=%q", nOut)
	}
	// Frontmatter unchanged: status still ideation, not zzz.
	fm := readFrontmatter(t, filepath.Join(root, "002-vendor-script.md"))
	if !strings.Contains(fm, "status: ideation") {
		t.Fatalf("entity status was mutated despite rejection:\n%s", fm)
	}
	if strings.Contains(fm, "status: zzz") {
		t.Fatalf("entity advanced to zzz despite rejection:\n%s", fm)
	}
}

// TestSetStatusMemberAccepted locks the complement: a `--set status=implementation`
// (a declared stage) exits zero and mutates, launcher vs oracle, so the membership
// guard does not reject legitimate transitions.
func TestSetStatusMemberAccepted(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixture(t, "seq-workflow")

	args := append([]string{"--workflow-dir", root}, "--set", "002-vendor-script", "status=implementation")
	nOut, nErr, nCode := runNative(t, root, env, args...)

	if nCode != 0 {
		t.Fatalf("exit: native=%d (%q)", nCode, nErr)
	}
	assertEnvelopeGolden(t, "set-membership-member-accepted", goldenEnvelope{
		stdout: normalize(nOut, root), stderr: normalize(nErr, root), exit: nCode,
	})
	fm := normalize(readFrontmatter(t, filepath.Join(root, "002-vendor-script.md")), root)
	assertTextGolden(t, "set-membership-member-frontmatter", fm)
	if !strings.Contains(fm, "status: implementation") {
		t.Fatalf("member status was not applied:\n%s", fm)
	}
}
