// ABOUTME: name-cap tests — dispatch build caps the worker name at 64 chars for
// ABOUTME: long slugs via an sd-b32 id-prefix, keeping short names byte-identical.
package dispatch

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// readmeIDStyle builds a worktree README declaring the given id-style. The stage
// set mirrors readmeWorktree so the same long-slug entity exercises both id-style
// fixtures; the only difference is the id-style frontmatter line.
func readmeIDStyle(idStyle string, splitRoot bool) string {
	state := ""
	if splitRoot {
		state = "state: state-checkout\n"
	}
	return "---\n" +
		"entity-type: task\n" +
		"id-style: " + idStyle + "\n" +
		state +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"      worktree: true\n" +
		"    - name: validation\n" +
		"      worktree: true\n" +
		"      feedback-to: implementation\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Fixture Workflow\n" +
		"\n" +
		"### backlog\n\nseed.\n\n- **Outputs:** x.\n\n" +
		"### implementation\n\nwork.\n\n- **Outputs:** y.\n\n" +
		"### validation\n\nverify.\n\n- **Outputs:** z.\n\n" +
		"### done\n\nterm.\n"
}

// entityFMID builds an entity file carrying an explicit id frontmatter value (a
// real 24-char sd-b32 id for the id-first fixtures), the given title/status, and
// no worktree (the name-cap fixtures dispatch the non-worktree backlog stage so
// no on-disk worktree dir is needed).
func entityFMID(id, title, status string) string {
	return "---\n" +
		"id: " + id + "\n" +
		"title: " + title + "\n" +
		"status: " + status + "\n" +
		"worktree: \n" +
		"---\n" +
		"# " + title + "\n\nBody.\n"
}

// nameFromStdout pulls the emitted name out of a build stdout JSON.
func nameFromStdout(t *testing.T, stdout string) string {
	t.Helper()
	var out struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, stdout)
	}
	if out.Name == "" {
		t.Fatalf("no name in stdout:\n%s", stdout)
	}
	return out.Name
}

// buildNameCapStdin writes an id-style: sd-b32 workflow with an entity carrying
// the given sd-b32 id and a (possibly long) slug, then returns the workflow dir
// and a well-formed build request for the backlog stage. The slug is the flat
// entity filename stem, so a long slug name overflows the 64-char budget.
func buildNameCapStdin(t *testing.T, root, idStyle, slug, id string) (workflowDir, stdin string) {
	t.Helper()
	wd := root
	writeFile(t, filepath.Join(wd, "README.md"), readmeIDStyle(idStyle, false))
	ep := filepath.Join(wd, slug+".md")
	writeFile(t, ep, entityFMID(id, "Thing", "backlog"))
	gitInit(t, root)
	return wd, mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    ep,
		"workflow_dir":   wd,
		"stage":          "backlog",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
	}, nil)
}

// A real 24-char sd-b32 id (alphabet 0123456789abcdefghjkmnpqrstvwxyz) used by
// the id-first fixtures. Distinct from idBravo so the no-collision AC holds.
const (
	idAlpha = "367s2zrbkm4fcwfff5ac62zz"
	idBravo = "0qv3kaqs062fp3p01tz6re99"
	// A 57-char slug: spacedock-ensign-{slug}-backlog overflows 64 well past the
	// ceiling (17 + 57 + 9 = 83 chars), reproducing the github#366 hazard.
	longSlug      = "dispatch-reconcile-deconflate-repo-hygiene-and-name-overflow"
	longSlugShare = "dispatch-reconcile-deconflate-repo-hygiene-and-name-sibling"
)

// TestBuildNameCapSDB32LongSlug — AC1: a sd-b32 entity whose uncapped name would
// exceed 64 chars emits a name ≤64 that matches namePattern. Fails today (the
// uncapped form is 80+ chars).
func TestBuildNameCapSDB32LongSlug(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd, stdin := buildNameCapStdin(t, root, "sd-b32", longSlug, idAlpha)

	native := runNative(stdin, "build", "--workflow-dir", wd)
	if native.exit != 0 {
		t.Fatalf("exit native=%d, want 0\nstderr:\n%s", native.exit, native.stderr)
	}
	name := nameFromStdout(t, native.stdout)
	if len(name) > 64 {
		t.Errorf("name %q is %d chars, want ≤64", name, len(name))
	}
	if !namePattern.MatchString(name) {
		t.Errorf("name %q does not match namePattern %s", name, namePattern.String())
	}
	// The capped name carries the id-prefix in place of the slug, keeping the
	// worker-key prefix and -{stage} suffix verbatim.
	if !strings.HasPrefix(name, "spacedock-ensign-") {
		t.Errorf("name %q lost worker-key prefix", name)
	}
	if !strings.HasSuffix(name, "-backlog") {
		t.Errorf("name %q lost -{stage} suffix", name)
	}
	if !strings.Contains(name, idAlpha[:sdB32NameIDPrefixLen]) {
		t.Errorf("name %q does not embed the id-prefix %q", name, idAlpha[:sdB32NameIDPrefixLen])
	}
}

// TestBuildNameCapDistinctIDs — AC2: two sd-b32 entities whose slugs share a long
// common prefix produce two different names (distinct because their ids differ).
func TestBuildNameCapDistinctIDs(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	rootA := t.TempDir()
	wdA, stdinA := buildNameCapStdin(t, rootA, "sd-b32", longSlug, idAlpha)
	nativeA := runNative(stdinA, "build", "--workflow-dir", wdA)
	if nativeA.exit != 0 {
		t.Fatalf("alpha exit=%d\nstderr:\n%s", nativeA.exit, nativeA.stderr)
	}
	nameA := nameFromStdout(t, nativeA.stdout)

	rootB := t.TempDir()
	wdB, stdinB := buildNameCapStdin(t, rootB, "sd-b32", longSlugShare, idBravo)
	nativeB := runNative(stdinB, "build", "--workflow-dir", wdB)
	if nativeB.exit != 0 {
		t.Fatalf("bravo exit=%d\nstderr:\n%s", nativeB.exit, nativeB.stderr)
	}
	nameB := nameFromStdout(t, nativeB.stdout)

	if nameA == nameB {
		t.Errorf("distinct entities produced identical name %q (cohort collision)", nameA)
	}
}

// TestBuildNameCapShortUnchanged — AC5: a short slug emits the uncapped
// {workerKey}-{slug}-{stage} byte-for-byte, with no id substitution. Runs for
// every id-style so the cap fires only on overflow regardless of id-style.
func TestBuildNameCapShortUnchanged(t *testing.T) {
	for _, idStyle := range []string{"sd-b32", "slug", "sequential"} {
		t.Run(idStyle, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			root := t.TempDir()
			id := idAlpha
			if idStyle == "sequential" {
				id = "001"
			} else if idStyle == "slug" {
				id = ""
			}
			wd, stdin := buildNameCapStdin(t, root, idStyle, "thing", id)

			native := runNative(stdin, "build", "--workflow-dir", wd)
			if native.exit != 0 {
				t.Fatalf("exit native=%d, want 0\nstderr:\n%s", native.exit, native.stderr)
			}
			name := nameFromStdout(t, native.stdout)
			want := "spacedock-ensign-thing-backlog"
			if name != want {
				t.Errorf("short name = %q, want %q (no cap should fire)", name, want)
			}
		})
	}
}

// TestBuildNameCapCycleHeadroom — AC6: a capped base name plus a realistic FO
// cycle suffix stays ≤64.
func TestBuildNameCapCycleHeadroom(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd, stdin := buildNameCapStdin(t, root, "sd-b32", longSlug, idAlpha)

	native := runNative(stdin, "build", "--workflow-dir", wd)
	if native.exit != 0 {
		t.Fatalf("exit native=%d, want 0\nstderr:\n%s", native.exit, native.stderr)
	}
	name := nameFromStdout(t, native.stdout)
	if len(name+"-cycle3") > 64 {
		t.Errorf("name+cycle3 %q is %d chars, want ≤64 (no cycle headroom)", name+"-cycle3", len(name+"-cycle3"))
	}
}

// TestBuildNameCapSlugFallback — id-less slug-truncation path: an id-style: slug
// long-slug entity (no stored id) still emits a ≤64 namePattern-valid name by
// truncating the slug head. This is the only path that truncates a slug.
func TestBuildNameCapSlugFallback(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	wd, stdin := buildNameCapStdin(t, root, "slug", longSlug, "")

	native := runNative(stdin, "build", "--workflow-dir", wd)
	if native.exit != 0 {
		t.Fatalf("exit native=%d, want 0\nstderr:\n%s", native.exit, native.stderr)
	}
	name := nameFromStdout(t, native.stdout)
	if len(name) > 64 {
		t.Errorf("slug-fallback name %q is %d chars, want ≤64", name, len(name))
	}
	if !namePattern.MatchString(name) {
		t.Errorf("slug-fallback name %q does not match namePattern (trailing hyphen / -- not trimmed?)", name)
	}
	if len(name+"-cycle3") > 64 {
		t.Errorf("slug-fallback name+cycle3 %q is %d chars, want ≤64", name+"-cycle3", len(name+"-cycle3"))
	}
	if !strings.HasPrefix(name, "spacedock-ensign-") || !strings.HasSuffix(name, "-backlog") {
		t.Errorf("slug-fallback name %q lost prefix/suffix", name)
	}
}
