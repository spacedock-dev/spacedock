// ABOUTME: verdict case boundary — the CLI takes lowercase `--verdict passed|rejected`,
// ABOUTME: frontmatter carries the schema's PASSED/REJECTED, and reads accept either case.
package status

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestMergeGuardWritesSchemaCasedVerdict pins the write-boundary normalisation.
// `merge guard --verdict` documents (and every caller uses) the lowercase
// `passed|rejected`, but entity.mdschema.yml declares the frontmatter value's
// conventional set as [PASSED REJECTED] — so writing the flag verbatim made
// `status --validate` warn on the entity the verb itself had just finalized, a
// permanent warning because the entity is terminal and archived. The verb
// upper-cases at the write boundary: the flag surface stays lowercase, the
// stored state matches the schema.
//
// The `rejected` leg is not a mirror of `passed`: it also covers the guards that
// read the verdict straight back out of frontmatter. finalize archives through
// runArchive without --force, and the archive-side merge-hook invariant exempts
// a rejected entity — an exemption that must survive the case change.
func TestMergeGuardWritesSchemaCasedVerdict(t *testing.T) {
	for _, tc := range []struct{ name, slug, flag, want string }{
		{"passed", "080-pr-merged", "passed", "PASSED"},
		{"rejected", "040-rejected", "rejected", "REJECTED"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, out, errOut, code := driveMergeGuard(t, "merge-pr-workflow", tc.slug, "--verdict", tc.flag)
			if code != 0 {
				t.Fatalf("finalize should exit 0, got %d (stdout=%q stderr=%q)", code, out, errOut)
			}
			archived := filepath.Join(root, "_archive", tc.slug+".md")
			if got := frontmatterField(t, archived, "verdict"); got != tc.want {
				t.Fatalf("archived verdict=%q, want %q (the schema's conventional value)", got, tc.want)
			}
			vOut, vErr, _ := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--validate")
			if strings.Contains(vOut+vErr, "field 'verdict'") {
				t.Fatalf("--validate must not warn on the verdict the verb just wrote:\nstdout=%s\nstderr=%s", vOut, vErr)
			}
		})
	}
}

// TestValidateAcceptsEitherVerdictCase pins the read side: an entity carrying a
// lowercase verdict — everything a pre-fix binary wrote — validates clean. The
// conventional-enum comparison is case-insensitive, so already-archived entities
// need no migration and no hand-editing.
func TestValidateAcceptsEitherVerdictCase(t *testing.T) {
	for _, verdict := range []string{"passed", "PASSED", "rejected", "REJECTED"} {
		t.Run(verdict, func(t *testing.T) {
			root := stageFixtureWith(t, "merge-pr-workflow", map[string]string{
				"_archive/200-legacy-verdict.md": "---\nid: \"200\"\ntitle: Legacy verdict casing\n" +
					"status: done\nverdict: " + verdict + "\nscore: \"0.5\"\nsource: roadmap\n---\n" +
					"# Legacy verdict casing\n",
			})
			out, errOut, _ := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--validate")
			if strings.Contains(out+errOut, "field 'verdict'") {
				t.Fatalf("verdict %q must validate clean in either case:\nstdout=%s\nstderr=%s", verdict, out, errOut)
			}
		})
	}
}

// TestSetStoresSchemaCasedVerdict pins the OTHER terminal writer. `merge guard`
// is not the only way a verdict reaches frontmatter: the user-facing
// `status --set <slug> verdict=passed` writes one too, and normalising only the
// merge-guard verb would leave stored case depending on which verb finalized the
// entity. updateFrontmatter canonicalises on every write, so both agree.
func TestSetStoresSchemaCasedVerdict(t *testing.T) {
	for _, tc := range []struct{ flag, want string }{
		{"passed", "PASSED"},
		{"PASSED", "PASSED"},
		{"rejected", "REJECTED"},
		{"Rejected", "REJECTED"},
	} {
		t.Run(tc.flag, func(t *testing.T) {
			root := stageFixture(t, "seq-workflow")
			_, errOut, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root,
				"--set", "002-vendor-script", "status=done", "completed", "verdict="+tc.flag, "worktree=")
			if code != 0 {
				t.Fatalf("--set finalize should exit 0, got %d (stderr=%q)", code, errOut)
			}
			if got := frontmatterField(t, filepath.Join(root, "002-vendor-script.md"), "verdict"); got != tc.want {
				t.Fatalf("--set verdict=%s stored %q, want %q", tc.flag, got, tc.want)
			}
		})
	}
}

// TestNonConventionalValuePassesThroughUnchanged pins the deliberate limit on
// canonicalisation. A conventional list is advisory (`invalid_severity: warn`),
// so a field may legitimately carry a value outside it — canonicalisation
// snaps values ONTO the schema's spellings, it does not case-fold everything.
// A blanket strings.ToUpper would silently rewrite `needs-work` to `NEEDS-WORK`,
// editing state the caller wrote on purpose.
func TestNonConventionalValuePassesThroughUnchanged(t *testing.T) {
	root := stageFixture(t, "seq-workflow")
	_, errOut, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root,
		"--set", "002-vendor-script", "verdict=needs-work", "--force")
	if code != 0 {
		t.Fatalf("--set with a non-conventional verdict should exit 0, got %d (stderr=%q)", code, errOut)
	}
	if got := frontmatterField(t, filepath.Join(root, "002-vendor-script.md"), "verdict"); got != "needs-work" {
		t.Fatalf("non-conventional verdict stored %q, want it unchanged as %q", got, "needs-work")
	}
}
