// ABOUTME: AC-3 proof — an old-line patch cut leaves `next`'s tip byte-identical
// ABOUTME: under the real decision guard, and clobbers it once the gate is removed.
package release

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

// patchNoRegressFixture builds the AC-3 scenario in a temp git repo: `next` has
// advanced to the 0.27 dev line with exclusive content (manifest 0.27.0-pre1,
// FO-prose required minor 0.27, a bumped calendar key), while an OLDER line cuts
// a vX.Y.1 patch (0.25.1) whose tree would clobber `next` if reconciled. It
// returns the repo dir, the patch tag, and `next`'s pristine tip so the subtests
// can prove the guard's decision leaves that tip untouched — and that removing
// the gate rewinds it.
type patchNoRegressFixture struct {
	dir           string
	patchTag      string
	releaseCommit string
	nextTip       string
	nextVersion   string
}

func newPatchNoRegressFixture(t *testing.T) patchNoRegressFixture {
	t.Helper()
	dir := t.TempDir()
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return string(out)
	}
	write := func(rel, content string) {
		t.Helper()
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	pluginJSON := func(version string) string {
		return "{\n  \"name\": \"spacedock\",\n  \"version\": \"" + version + "\"\n}\n"
	}
	codexJSON := func(version string) string {
		return "{\n  \"name\": \"spacedock-codex\",\n  \"version\": \"" + version + "\"\n}\n"
	}
	prose := func(minor string) string {
		return "# FO shared core\n\nThese skills require binary minor " + minor + " to boot.\n"
	}
	marketplaceJSON := func(cal string) string {
		return "{\n  \"name\": \"spacedock-edge\",\n  \"plugins\": [\n    {\n      \"name\": \"spacedock\",\n      \"version\": \"" + cal + "\"\n    }\n  ]\n}\n"
	}
	const foProse = "skills/first-officer/references/first-officer-shared-core.md"

	testgit.InitRepo(t, dir, "-q")

	// Shared ancestor both the old release line and `next` descend from.
	write(".claude-plugin/plugin.json", pluginJSON("0.24.0-pre1"))
	write(".codex-plugin/plugin.json", codexJSON("0.24.0-pre1"))
	write(foProse, prose("0.24"))
	write(".claude-plugin/marketplace.json", marketplaceJSON("0.0.2026050101"))
	write("app.txt", "base\n")
	git("add", "-A")
	git("commit", "-q", "-m", "seed")
	git("branch", "-M", "main")
	git("branch", "next")

	// The OLD line advances to a vX.Y.1 patch (0.25.1) with its own content on the
	// same files `next` also touches — a genuine divergence a -X theirs reconcile
	// would resolve in the OLD release's favor.
	patchVersion := "0.25.1"
	write("app.txt", "old-line-patch\n")
	write(".claude-plugin/plugin.json", pluginJSON(patchVersion))
	write(".codex-plugin/plugin.json", codexJSON(patchVersion))
	write(foProse, prose("0.25"))
	git("commit", "-q", "-am", "release: stamp "+patchVersion)
	releaseCommit := strings.TrimSpace(git("rev-parse", "HEAD"))
	patchTag := "v" + patchVersion
	git("tag", "-a", patchTag, releaseCommit, "-m", "patch "+patchVersion)

	// `next` advances PAST the release to the 0.27 dev line with exclusive content.
	git("switch", "-q", "next")
	nextVersion := "0.27.0-pre1"
	write("app.txt", "next-0.27-exclusive\n")
	write(".claude-plugin/plugin.json", pluginJSON(nextVersion))
	write(".codex-plugin/plugin.json", codexJSON(nextVersion))
	write(foProse, prose("0.27"))
	write(".claude-plugin/marketplace.json", marketplaceJSON("0.0.2026060101"))
	git("commit", "-q", "-am", "next: advance to 0.27 dev line")
	nextTip := strings.TrimSpace(git("rev-parse", "HEAD"))

	return patchNoRegressFixture{
		dir:           dir,
		patchTag:      patchTag,
		releaseCommit: releaseCommit,
		nextTip:       nextTip,
		nextVersion:   nextVersion,
	}
}

// runGuardedEdgeAdvance mirrors release.yml's decision-gated stable path against
// the fixture repo: it computes the REAL EdgeAdvanceDecision for the tag, and
// runs the reconcile+stamp+bump ONLY when it says advance. gateRemoved=true drops
// the gate (the adversarial twin), running the reconcile unconditionally to show
// the damage the gate prevents. It returns the decision so the caller can assert
// the guard's verdict.
func runGuardedEdgeAdvance(t *testing.T, f patchNoRegressFixture, branch string, gateRemoved bool) (advance bool) {
	t.Helper()
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = f.dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return string(out)
	}
	advance, _, err := EdgeAdvanceDecision(f.patchTag, f.nextVersion)
	if err != nil {
		t.Fatalf("EdgeAdvanceDecision(%q, %q): %v", f.patchTag, f.nextVersion, err)
	}
	if !advance && !gateRemoved {
		// The gate skips the whole job: `next` is left byte-identical.
		return advance
	}

	// The reconcile the stable path runs (guarded), or the ungated adversarial run.
	git("switch", "-q", "-c", branch, "next")
	git("merge", "-X", "theirs", "--no-edit", f.releaseCommit, "-m", "next: reconcile edge line to "+f.patchTag)
	dev, err := DevPreVersion(strings.TrimPrefix(f.patchTag, "v"))
	if err != nil {
		t.Fatalf("DevPreVersion: %v", err)
	}
	const foProse = "skills/first-officer/references/first-officer-shared-core.md"
	for _, rel := range []string{".claude-plugin/plugin.json", ".codex-plugin/plugin.json", foProse} {
		path := filepath.Join(f.dir, rel)
		data, rerr := os.ReadFile(path)
		if rerr != nil {
			t.Fatal(rerr)
		}
		var out []byte
		if strings.HasSuffix(rel, ".md") {
			out, err = StampProseVersion(data, dev)
		} else {
			out, err = StampVersion(data, dev)
		}
		if err != nil {
			t.Fatalf("stamp %s: %v", rel, err)
		}
		if werr := os.WriteFile(path, out, 0o644); werr != nil {
			t.Fatal(werr)
		}
	}
	git("commit", "-q", "-m", "next: bump dev pre-version to "+dev,
		"--", ".claude-plugin/plugin.json", ".codex-plugin/plugin.json", foProse)
	mktPath := filepath.Join(f.dir, ".claude-plugin/marketplace.json")
	mktData, rerr := os.ReadFile(mktPath)
	if rerr != nil {
		t.Fatal(rerr)
	}
	bumped, berr := BumpCalendarVersion(mktData, time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC))
	if berr != nil {
		t.Fatalf("BumpCalendarVersion: %v", berr)
	}
	if werr := os.WriteFile(mktPath, bumped, 0o644); werr != nil {
		t.Fatal(werr)
	}
	git("commit", "-q", "-m", "next: bump marketplace calendar version",
		"--", ".claude-plugin/marketplace.json")
	return advance
}

func gitRevParse(t *testing.T, dir, rev string) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", rev)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse %s: %v\n%s", rev, err, out)
	}
	return strings.TrimSpace(string(out))
}

// gitShow returns the bytes of path as committed on branch (never the working
// tree), so a branch-scoped read is unaffected by whatever branch is checked out.
func gitShow(t *testing.T, dir, branch, path string) []byte {
	t.Helper()
	cmd := exec.Command("git", "show", branch+":"+path)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git show %s:%s: %v\n%s", branch, path, err, out)
	}
	return out
}

func readEdgeVersionOnBranch(t *testing.T, dir, branch, rel string) string {
	t.Helper()
	v, err := ManifestVersion(gitShow(t, dir, branch, rel))
	if err != nil {
		t.Fatalf("parse %s on %s: %v", rel, branch, err)
	}
	return v
}

func readEdgeCalendarOnBranch(t *testing.T, dir, branch string) string {
	t.Helper()
	var doc struct {
		Plugins []struct {
			Version string `json:"version"`
		} `json:"plugins"`
	}
	if err := json.Unmarshal(gitShow(t, dir, branch, ".claude-plugin/marketplace.json"), &doc); err != nil {
		t.Fatalf("parse marketplace.json on %s: %v", branch, err)
	}
	if len(doc.Plugins) == 0 {
		t.Fatalf("marketplace.json on %s has no plugin entry", branch)
	}
	return doc.Plugins[0].Version
}

func gitTagList(t *testing.T, dir string) []string {
	t.Helper()
	cmd := exec.Command("git", "tag", "--list")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git tag --list: %v\n%s", err, out)
	}
	var tags []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			tags = append(tags, line)
		}
	}
	return tags
}

func proseMinorOnBranch(t *testing.T, dir, branch string) string {
	t.Helper()
	minor, err := ProseMinor(gitShow(t, dir, branch, "skills/first-officer/references/first-officer-shared-core.md"))
	if err != nil {
		t.Fatalf("ProseMinor on %s: %v", branch, err)
	}
	return minor
}

// TestEdgeAdvancePatchDoesNotRegressNext (AC-3) proves the decision guard is
// load-bearing on all four regress counts. Under the real guard an old-line
// vX.Y.1 patch is decided `skip`, so `next`'s tip is left byte-identical — same
// commit; no manifest/gate-line rewind; the marketplace calendar key unchanged —
// and no colliding pre0 tag is cut. The adversarial half removes the gate and
// runs the same reconcile: `-X theirs` clobbers `next`'s newer content, the stamp
// rewinds the manifest (0.27.0-pre1 → 0.26.0-pre1) and the prose minor
// (0.27 → 0.26), and the calendar key advances — the exact damage the gate
// prevents.
func TestEdgeAdvancePatchDoesNotRegressNext(t *testing.T) {
	t.Run("guarded patch skips, next byte-identical", func(t *testing.T) {
		f := newPatchNoRegressFixture(t)
		advance := runGuardedEdgeAdvance(t, f, "edge-advance", false)
		if advance {
			t.Fatalf("EdgeAdvanceDecision advanced for old-line patch %q vs next %q; want skip", f.patchTag, f.nextVersion)
		}
		// next's tip is untouched: same commit, same manifest, prose, calendar.
		if got := gitRevParse(t, f.dir, "next"); got != f.nextTip {
			t.Fatalf("next tip moved on a skip: %s != pristine %s", got, f.nextTip)
		}
		if got := readEdgeVersionOnBranch(t, f.dir, "next", ".claude-plugin/plugin.json"); got != "0.27.0-pre1" {
			t.Fatalf("next manifest rewound on a skip: %q, want 0.27.0-pre1", got)
		}
		if got := proseMinorOnBranch(t, f.dir, "next"); got != "0.27" {
			t.Fatalf("next prose minor rewound on a skip: %q, want 0.27", got)
		}
		if got := readEdgeCalendarOnBranch(t, f.dir, "next"); got != "0.0.2026060101" {
			t.Fatalf("next calendar key churned on a skip: %q, want 0.0.2026060101", got)
		}
		// No colliding vX.(Y+1).0-pre0 tag was cut.
		tags := gitTagList(t, f.dir)
		for _, tag := range tags {
			if strings.HasSuffix(tag, "-pre0") {
				t.Fatalf("a pre0 tag %q was cut on a skip; want none", tag)
			}
		}
	})

	t.Run("gate removed clobbers next", func(t *testing.T) {
		f := newPatchNoRegressFixture(t)
		runGuardedEdgeAdvance(t, f, "edge-advance-ungated", true)
		branch := "edge-advance-ungated"
		if got := gitRevParse(t, f.dir, branch); got == f.nextTip {
			t.Fatal("ungated reconcile left next's tip unchanged; the adversarial twin proved nothing")
		}
		if got := readEdgeVersionOnBranch(t, f.dir, branch, ".claude-plugin/plugin.json"); got != "0.26.0-pre1" {
			t.Fatalf("ungated reconcile manifest = %q, want the rewound 0.26.0-pre1 (proof the gate matters)", got)
		}
		if got := proseMinorOnBranch(t, f.dir, branch); got != "0.26" {
			t.Fatalf("ungated reconcile prose minor = %q, want the rewound 0.26", got)
		}
		if got := readEdgeCalendarOnBranch(t, f.dir, branch); got == "0.0.2026060101" {
			t.Fatal("ungated reconcile left the calendar key unchanged; it should have churned")
		}
	})
}
