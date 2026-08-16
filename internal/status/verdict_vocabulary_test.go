// ABOUTME: the verdict-vocabulary ruling — `superseded` stays OUT of the enum, and
// ABOUTME: a schema `conventional` list is advisory on read but closed on --set write.
package status

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// verdictEntity renders a minimal active-scope entity for the enum-scope fixture.
// An empty verdict omits the field entirely — the supported no-verdict shape, and
// what the supported supersede path leaves behind.
func verdictEntity(slug, title, stage, verdict string) string {
	verdictLine := ""
	if verdict != "" {
		verdictLine = "verdict: " + verdict + "\n"
	}
	return fmt.Sprintf("---\nid: %s\ntitle: %s\nstatus: %s\nscore: \"0.50\"\n"+
		"source: verdict-vocabulary fixture\n%s---\n\nFixture entity for the verdict-vocabulary ruling.\n",
		slug, title, stage, verdictLine)
}

// warningLines returns the `Warning:` diagnostics in a validate stderr. The count
// is load-bearing for AC-1 (the 1 -> 0 delta) and AC-3 (exactly one), so the tests
// assert on the lines rather than on a substring anywhere in stderr.
func warningLines(stderr string) []string {
	var out []string
	for _, line := range strings.Split(stderr, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "Warning:") {
			out = append(out, strings.TrimSpace(line))
		}
	}
	return out
}

// TestSetRefusesNonConventionalToken pins AC-2: the ruling is enforced where
// verdicts are WRITTEN. A field declaring a schema `conventional` list is closed
// on write, so `--set verdict=superseded` — the exact write that produced the four
// archived records — refuses non-zero and leaves the entity byte-identical.
//
// Byte-identity is the assertion that matters: an exit code alone would pass a
// guard that mutates and then errors. The table is keyed on field so a future
// field that declares a `conventional` list is one new row, not a new test.
func TestSetRefusesNonConventionalToken(t *testing.T) {
	env := pinnedEnv(t)
	cases := []struct {
		name    string
		field   string
		token   string
		allowed string
	}{
		{"the token that produced the four records", "verdict", "superseded", "[PASSED REJECTED]"},
		{"the writer gap is general, not superseded-specific", "verdict", "banana", "[PASSED REJECTED]"},
		{"refusal is case-insensitive, like the case-fold it extends", "verdict", "SUPERSEDED", "[PASSED REJECTED]"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := stageFixtureWith(t, "enum-scope-workflow", map[string]string{
				"probe.md": verdictEntity("probe", "Write-path probe", "backlog", ""),
			})
			path := filepath.Join(root, "probe.md")
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			out, errOut, code := runNative(t, root, env,
				"--workflow-dir", root, "--set", "probe", tc.field+"="+tc.token)

			if code == 0 {
				t.Fatalf("--set %s=%s must exit non-zero; got 0 (stdout=%q)", tc.field, tc.token, out)
			}
			for _, want := range []string{tc.field, tc.token, tc.allowed} {
				if !strings.Contains(errOut, want) {
					t.Errorf("refusal must name %q; stderr=%q", want, errOut)
				}
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(before) {
				t.Fatalf("refusal must be byte-clean — entity changed.\nbefore:\n%s\nafter:\n%s", before, after)
			}
		})
	}
}

// TestSetForceWritesNonConventionalToken pins AC-2's escape: --force is the
// uniform bypass every other --set guard honors, so the admission check must not
// become the one guard with no escape hatch. Fails if --force stops writing the
// non-conventional token.
func TestSetForceWritesNonConventionalToken(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "enum-scope-workflow", map[string]string{
		"probe.md": verdictEntity("probe", "Write-path probe", "backlog", ""),
	})

	_, errOut, code := runNative(t, root, env,
		"--workflow-dir", root, "--set", "probe", "verdict=superseded", "--force")
	if code != 0 {
		t.Fatalf("--force must bypass the admission check; exit=%d stderr=%q", code, errOut)
	}

	got := ParseFrontmatter(filepath.Join(root, "probe.md"))["verdict"]
	if strings.TrimSpace(got) != "superseded" {
		t.Fatalf("--force must write the token verbatim; verdict=%q", got)
	}
}

// TestSetClearAndCaseFoldStillPass pins the other half of AC-2: the admission
// check must not swallow the writes the ruling depends on. Clearing is the
// supported supersede shape, the case-fold is the behavior the guard's own schema
// lookup was already performing, and a field with no `conventional` list is
// outside the guard's blast radius entirely.
func TestSetClearAndCaseFoldStillPass(t *testing.T) {
	env := pinnedEnv(t)
	cases := []struct {
		name       string
		startsWith string
		assign     string
		wantStored string
	}{
		{"clearing always passes — the supported supersede shape", "PASSED", "verdict=", ""},
		{"lowercase CLI spelling still folds to the schema spelling", "", "verdict=passed", "PASSED"},
		{"the other conventional token folds too", "", "verdict=rejected", "REJECTED"},
		{"an already-canonical token is admitted unchanged", "", "verdict=PASSED", "PASSED"},
		{"a field with no conventional list is untouched by the guard", "", "source=anything at all", "anything at all"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := stageFixtureWith(t, "enum-scope-workflow", map[string]string{
				"probe.md": verdictEntity("probe", "Write-path probe", "backlog", tc.startsWith),
			})

			_, errOut, code := runNative(t, root, env,
				"--workflow-dir", root, "--set", "probe", tc.assign)
			if code != 0 {
				t.Fatalf("--set %s must succeed; exit=%d stderr=%q", tc.assign, code, errOut)
			}

			field := strings.SplitN(tc.assign, "=", 2)[0]
			got := strings.TrimSpace(ParseFrontmatter(filepath.Join(root, "probe.md"))[field])
			if got != tc.wantStored {
				t.Fatalf("--set %s stored %s=%q, want %q", tc.assign, field, got, tc.wantStored)
			}
		})
	}
}

// TestSupersededVerdictBaselineWarnsAndBlocks pins the baseline half of AC-1 —
// the two numbers that must move. An active entity carrying `verdict: superseded`
// emits exactly ONE validate warning AND makes `--archive` exit 1, because the
// archive guard reads every non-empty, non-`rejected` verdict as approval-style:
// the field written to say "this was superseded" is exactly what blocks the
// superseding. The fixture is written directly rather than through `--set`,
// since `--set` now refuses the token.
//
// Fails when the baseline stops being measurable — i.e. when the "wrong way"
// numbers stop being 1 and 1, which is what makes the 1 -> 0 delta in
// TestSupersedePathIsWarningFree an end-value rather than a tautology.
func TestSupersededVerdictBaselineWarnsAndBlocks(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "enum-scope-workflow", map[string]string{
		"stranded.md": verdictEntity("stranded", "Superseded the unsupported way", "backlog", "superseded"),
	})

	_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--validate")
	if code != 0 {
		t.Fatalf("--validate exit=%d, want 0 (field conformance is warn-tier); stderr=%q", code, errOut)
	}
	warns := warningLines(errOut)
	if len(warns) != 1 {
		t.Fatalf("baseline warning count = %d, want 1:\n%s", len(warns), errOut)
	}
	if !strings.Contains(warns[0], "superseded") || !strings.Contains(warns[0], "[PASSED REJECTED]") {
		t.Fatalf("baseline warning must name the token and the allowed set: %q", warns[0])
	}

	_, archErr, archCode := runNative(t, root, env, "--workflow-dir", root, "--archive", "stranded")
	if archCode != 1 {
		t.Fatalf("baseline --archive exit=%d, want 1 (the token strands the entity); stderr=%q", archCode, archErr)
	}
	if !strings.Contains(archErr, "superseded") {
		t.Fatalf("baseline archive refusal must name the blocking verdict: %q", archErr)
	}
}

// TestSupersedePathIsWarningFree pins the value half of AC-1: an intentional
// supersede of an active entity completes with zero validate warnings and zero
// archive refusals, through the supported path the ruling names — leave `verdict`
// empty and `--archive`. Read against the baseline above, this is the pair
// (warning count 1 -> 0, archive exit 1 -> 0).
//
// Fails if the supported path regains a warning, or if `--archive` starts
// refusing an empty-verdict entity.
func TestSupersedePathIsWarningFree(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "enum-scope-workflow", map[string]string{
		"stranded.md": verdictEntity("stranded", "Superseded the supported way", "backlog", "superseded"),
	})

	// The supported supersede: clear the verdict, which the admission check must
	// keep passing, then archive.
	_, clearErr, clearCode := runNative(t, root, env,
		"--workflow-dir", root, "--set", "stranded", "verdict=")
	if clearCode != 0 {
		t.Fatalf("clearing verdict must succeed; exit=%d stderr=%q", clearCode, clearErr)
	}

	_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--validate")
	if code != 0 {
		t.Fatalf("--validate exit=%d, want 0; stderr=%q", code, errOut)
	}
	if warns := warningLines(errOut); len(warns) != 0 {
		t.Fatalf("supported supersede must validate with zero warnings; got %d:\n%s", len(warns), errOut)
	}

	_, archErr, archCode := runNative(t, root, env, "--workflow-dir", root, "--archive", "stranded")
	if archCode != 0 {
		t.Fatalf("supported supersede must archive cleanly; exit=%d stderr=%q", archCode, archErr)
	}
	if _, err := os.Stat(filepath.Join(root, "_archive", "stranded.md")); err != nil {
		t.Fatalf("archived entity must land in _archive/: %v", err)
	}
}

// TestNonConventionalVerdictStillWarns is AC-3, the negative control. AC-1 asks
// for zero warnings, and the cheapest way to "achieve" that is to silence the
// checker: setting `verdict.invalid_severity: error` makes isWarnSeverity return
// false and the field is skipped entirely (there is no error path for field
// conformance), and deleting the `conventional` list has the same effect while
// also discarding the case-fold the write path depends on. Either cheat drops the
// warning count to 0 and fails this test.
func TestNonConventionalVerdictStillWarns(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "enum-scope-workflow", map[string]string{
		"odd.md": verdictEntity("odd", "Non-conventional verdict on an active entity", "backlog", "banana"),
	})

	_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--validate")
	if code != 0 {
		t.Fatalf("--validate exit=%d, want 0 (warn-tier, never exit 1); stderr=%q", code, errOut)
	}
	warns := warningLines(errOut)
	if len(warns) != 1 {
		t.Fatalf("an active non-conventional verdict must produce exactly one warning; got %d:\n%s", len(warns), errOut)
	}
	for _, want := range []string{"verdict", "banana", "[PASSED REJECTED]", "odd"} {
		if !strings.Contains(warns[0], want) {
			t.Errorf("warning must name %q: %q", want, warns[0])
		}
	}
}

// TestArchivedScopeStaysSilent is AC-4's fixture half: archived scope is
// publish-only, so an archived entity carrying `verdict: superseded` emits
// nothing. This is what lets the ruling leave the four archived records
// untouched. Fails if archived scope starts warning again.
//
// AC-4's state half — that those four records still read `verdict: superseded`
// byte-for-byte — is a check against the state checkout, not this module, so it
// is verified and cited in the stage report rather than tested here.
func TestArchivedScopeStaysSilent(t *testing.T) {
	env := pinnedEnv(t)
	root := stageFixtureWith(t, "enum-scope-workflow", map[string]string{
		filepath.Join("_archive", "sup-arch.md"): verdictEntity(
			"sup-arch", "Archived record carrying the retired token", "done", "superseded"),
	})

	_, errOut, code := runNative(t, root, env, "--workflow-dir", root, "--validate")
	if code != 0 {
		t.Fatalf("--validate exit=%d, want 0; stderr=%q", code, errOut)
	}
	if warns := warningLines(errOut); len(warns) != 0 {
		t.Fatalf("archived scope must stay silent; got %d warning(s):\n%s", len(warns), errOut)
	}
}
