// ABOUTME: AC-2 — the copyable entity stub template is present in the shared
// ABOUTME: constant and surfaced by new's no-frontmatter stdin error, with red-controls.
package status

import (
	"io"
	"strings"
	"testing"
)

// stubTemplateComplete reports whether out carries all three token groups of the
// entity stub template: the frontmatter skeleton (title + status), the id-omitted
// directive, and the post-new commit next-step. AND'd, not any single group — a
// frontmatter skeleton alone does not tell the FO to omit the id or how to commit;
// an id-omitted or commit note alone hands the FO no fillable shape. Mirrors the
// sibling reference_load_anti_hunt_test.go two-token-AND shape.
func stubTemplateComplete(out string) bool {
	hasSkeleton := strings.Contains(out, "title:") && strings.Contains(out, "status:")
	hasIDOmitted := strings.Contains(out, "id OMITTED")
	hasCommitStep := strings.Contains(out, "state commit")
	return hasSkeleton && hasIDOmitted && hasCommitStep
}

// TestEntityStubTemplateHasRequiredTokens is the AC-2 presence check on the shared
// constant: it carries the frontmatter skeleton, the id-omitted directive, and the
// commit next-step.
func TestEntityStubTemplateHasRequiredTokens(t *testing.T) {
	if !stubTemplateComplete(EntityStubTemplate) {
		t.Errorf("EntityStubTemplate is missing a required token group (frontmatter skeleton / id-omitted / commit next-step):\n%s", EntityStubTemplate)
	}
}

// TestStubTemplateGuardFailsOnDeletion is the AC-2 red-control: the check can go
// RED, not just pass. The full constant passes; an empty string reds; and each
// single-group fragment (half the invariant) reds, so the check cannot pass
// vacuously off one group alone.
func TestStubTemplateGuardFailsOnDeletion(t *testing.T) {
	if !stubTemplateComplete(EntityStubTemplate) {
		t.Fatal("control: the full template was NOT flagged complete — the check would pass vacuously")
	}
	if stubTemplateComplete("") {
		t.Error("control: an empty output was wrongly reported as carrying the template")
	}
	skeletonOnly := "title: <t>\nstatus: <s>"
	if stubTemplateComplete(skeletonOnly) {
		t.Error("control: a frontmatter-skeleton-only fragment (no id-omitted, no commit step) was wrongly reported complete")
	}
	idOmittedOnly := "with id OMITTED (new mints it)"
	if stubTemplateComplete(idOmittedOnly) {
		t.Error("control: an id-omitted-only fragment (no skeleton, no commit step) was wrongly reported complete")
	}
	commitOnly := "run spacedock state commit SLUG afterward"
	if stubTemplateComplete(commitOnly) {
		t.Error("control: a commit-step-only fragment (no skeleton, no id-omitted) was wrongly reported complete")
	}
}

// TestNewNoFrontmatterErrorShowsStubTemplate drives runNew (via the native runner)
// with empty and malformed stdin and asserts the error surfaces the full stub
// template, so a filing FO that pipes nothing or a bad stub is handed the shape
// instead of hunting for an example. The terse "no frontmatter" guard is kept.
func TestNewNoFrontmatterErrorShowsStubTemplate(t *testing.T) {
	env := pinnedEnv(t)
	for _, tc := range []struct {
		name  string
		stdin io.Reader
	}{
		{"empty-stdin", reader("")},
		{"malformed-stdin", reader("no frontmatter here\n")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := stageFixture(t, "seq-workflow")
			_, errOut, code := runNativeStdin(t, root, env, tc.stdin, "--workflow-dir", root, "--new", "fumbled")
			if code != 1 {
				t.Fatalf("exit=%d, want 1", code)
			}
			if !strings.Contains(errOut, "no frontmatter") {
				t.Fatalf("stderr dropped the terse guard message: %q", errOut)
			}
			if !stubTemplateComplete(errOut) {
				t.Fatalf("stderr did not surface the stub template:\n%s", errOut)
			}
		})
	}
}

// templateBlockFromPrinted returns the Body-template block a filing FO copies out
// of printed new output: the lines from the "Body template" header to the next
// blank-line section break (or end). It reads the LITERAL printed output, not the
// source constant, so the round-trip below proves what a human actually copies is
// pipe-safe.
func templateBlockFromPrinted(printed string) string {
	const header = "Body template (pipe this on stdin; fill the angle-bracket fields, id omitted):\n"
	i := strings.Index(printed, header)
	if i < 0 {
		return ""
	}
	block := printed[i+len(header):]
	if end := strings.Index(block, "\n\n"); end >= 0 {
		block = block[:end+1]
	}
	return block
}

// TestPrintedStubTemplateBlockIsPipeSafe is the round-trip guard the audit finding
// demands: it extracts the Body-template block AS PRINTED by new's no-frontmatter
// error and pipes it VERBATIM back through runNew, asserting the atomic create
// succeeds. A display indent on the fence line makes contentHasOpeningFence reject
// it, so this reds on exactly the defect the substring/self-referential checks
// missed. Placeholders are piped unfilled: the point is the SHAPE parses, as a
// human's first copy-paste would.
func TestPrintedStubTemplateBlockIsPipeSafe(t *testing.T) {
	env := pinnedEnv(t)

	// Capture the LITERAL printed template from new's own error output.
	root := stageFixture(t, "seq-workflow")
	_, printed, code := runNativeStdin(t, root, env, reader("oops\n"), "--workflow-dir", root, "--new", "trigger")
	if code != 1 {
		t.Fatalf("setup: malformed-stdin new exit=%d, want 1", code)
	}
	block := templateBlockFromPrinted(printed)
	if block == "" {
		t.Fatalf("could not locate the Body-template block in printed output:\n%s", printed)
	}
	if strings.HasPrefix(block, " ") || strings.HasPrefix(block, "\t") {
		t.Fatalf("printed block's first line carries leading whitespace — a literal copy-paste is not pipe-safe:\n%q", block)
	}

	// Pipe that printed block VERBATIM into a fresh new — it must parse and create.
	root2 := stageFixture(t, "seq-workflow")
	out, errOut, code := runNativeStdin(t, root2, env, reader(block), "--workflow-dir", root2, "--new", "piped-back")
	if code != 0 {
		t.Fatalf("piping the printed template block verbatim failed: exit=%d\nstderr=%q\nblock=%q", code, errOut, block)
	}
	if !strings.Contains(out, "created:") {
		t.Fatalf("expected a created entity, got out=%q", out)
	}
}
