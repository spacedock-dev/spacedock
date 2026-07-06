// ABOUTME: --checklist/--ac-scan --stage defaulting to the entity's current
// ABOUTME: frontmatter status when --stage is omitted (AC-1 through AC-5).
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestDefaultStageMatchesExplicitCurrent (AC-1, AC-2) asserts a bare --checklist /
// --ac-scan call (no --stage) against an entity whose current status is
// "validation" succeeds and produces byte-identical output to the same call with
// --stage validation given explicitly, across both text and --json modes. The
// interleaved fixture's frontmatter status is "validation" and it carries a
// matching "## Stage Report: validation" section.
func TestDefaultStageMatchesExplicitCurrent(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	fixture := interleavedFixturePath(t)

	cases := []struct {
		name string
		flag string
	}{
		{"checklist", "--checklist"},
		{"ac-scan", "--ac-scan"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name+" text", func(t *testing.T) {
			omitted, _, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, tc.flag)
			if code != 0 {
				t.Fatalf("bare %s exit=%d, want 0", tc.flag, code)
			}
			if strings.Contains(omitted, "requires --stage") {
				t.Fatalf("bare %s emitted the retired 'requires --stage' wording: %s", tc.flag, omitted)
			}
			explicit, _, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--stage", "validation", tc.flag)
			if code != 0 {
				t.Fatalf("explicit %s exit=%d, want 0", tc.flag, code)
			}
			if omitted != explicit {
				t.Fatalf("bare %s output != explicit --stage validation output\n--- omitted ---\n%s\n--- explicit ---\n%s", tc.flag, omitted, explicit)
			}
		})
		t.Run(tc.name+" json", func(t *testing.T) {
			omitted, _, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, tc.flag, "--json")
			if code != 0 {
				t.Fatalf("bare %s --json exit=%d, want 0", tc.flag, code)
			}
			explicit, _, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--stage", "validation", tc.flag, "--json")
			if code != 0 {
				t.Fatalf("explicit %s --json exit=%d, want 0", tc.flag, code)
			}
			if omitted != explicit {
				t.Fatalf("bare %s --json output != explicit --stage validation output\n--- omitted ---\n%s\n--- explicit ---\n%s", tc.flag, omitted, explicit)
			}
			if !strings.Contains(omitted, `"stage":"validation"`) {
				t.Fatalf("bare %s --json envelope stage != validation: %s", tc.flag, omitted)
			}
		})
	}
}

// TestExplicitStageOverridesDefault (AC-3) asserts --stage <name> with a name
// other than the current status still selects that stage's report, unaffected
// by the current-status default. The interleaved fixture's current status is
// "validation"; --stage implementation must still select the LATEST
// implementation cycle, exactly as it does without any default in play.
func TestExplicitStageOverridesDefault(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	fixture := interleavedFixturePath(t)

	out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, "--stage", "implementation", "--checklist", "--json")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	if !strings.Contains(out, "rework after the gate") {
		t.Fatalf("explicit --stage implementation did not select the latest implementation cycle: %s", out)
	}
	if strings.Contains(out, "write the first cut") {
		t.Fatalf("explicit --stage implementation leaked the earlier cycle-1 item: %s", out)
	}
	if !strings.Contains(out, `"stage":"implementation"`) {
		t.Fatalf("json envelope stage != implementation despite current status being validation: %s", out)
	}
}

// TestDefaultStageNoMatchingReport (AC-4) asserts that when --stage is omitted
// and the entity's current status names a stage with no matching Stage Report,
// the command exits 1 with a diagnostic naming the stage as the current/defaulted
// one — never a silent empty emit. The section-reader fixture's status is
// "implementation" but it only carries a "## Stage Report: ideation" section.
func TestDefaultStageNoMatchingReport(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	fixture := fixturePath(t)

	for _, flag := range []string{"--checklist", "--ac-scan"} {
		flag := flag
		t.Run(flag, func(t *testing.T) {
			out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", fixture, flag)
			if code == 0 {
				t.Fatalf("exit=0 (want non-zero) — silent/partial emit\nstdout=%q", out)
			}
			if strings.TrimSpace(out) != "" {
				t.Fatalf("stdout = %q, want empty on a loud failure", out)
			}
			if !strings.Contains(stderr, `"implementation"`) {
				t.Fatalf("stderr = %q, want it to name the defaulted stage \"implementation\"", stderr)
			}
			if !strings.Contains(stderr, "--stage omitted") {
				t.Fatalf("stderr = %q, want it to note --stage was omitted (defaulting transparency)", stderr)
			}
		})
	}
}

// TestDefaultStageNoStatusField (AC-5) asserts that when --stage is omitted and
// the resolved --read target has no status frontmatter field, the command exits
// 1 with a diagnostic naming that the stage could not be defaulted — never the
// retired bare "requires --stage" wording, and never a silent emit. The real
// docs/dev/README.md has no status: frontmatter field.
func TestDefaultStageNoStatusField(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("testdata", "seq-workflow"))
	if err != nil {
		t.Fatal(err)
	}
	env := pinnedEnv(t)
	readme := devReadmePath(t)

	for _, flag := range []string{"--checklist", "--ac-scan"} {
		flag := flag
		t.Run(flag, func(t *testing.T) {
			out, stderr, code := runNative(t, root, env, "--workflow-dir", root, "--read", readme, flag)
			if code == 0 {
				t.Fatalf("exit=0 (want non-zero) — silent/partial emit\nstdout=%q", out)
			}
			if strings.TrimSpace(out) != "" {
				t.Fatalf("stdout = %q, want empty on a loud failure", out)
			}
			if strings.Contains(stderr, "requires --stage") {
				t.Fatalf("stderr = %q, reused the retired 'requires --stage' wording (misleading now that the flag defaults)", stderr)
			}
			if !strings.Contains(stderr, "--stage omitted") || !strings.Contains(strings.ToLower(stderr), "status") {
				t.Fatalf("stderr = %q, want a diagnostic naming that --stage was omitted and there is no status to default from", stderr)
			}
		})
	}
}
