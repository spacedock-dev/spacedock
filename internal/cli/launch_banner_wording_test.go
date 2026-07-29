// ABOUTME: AC-1..AC-7 oracles for the rewritten pre-launch banner — consistent
// ABOUTME: first-officer framing, no self-serve status, multi-workflow path list.
package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// renderBanner exercises the real launchBanner for host from dir and returns the
// rendered stderr bytes — the proof surface for AC-1/2/3/4/5 (the bytes the
// operator actually sees), never a source-grep of the constant. The sandbox
// inputs are pinned (not selected, binary not found) so the Sandbox: line is a
// stable `unavailable` for these workflow-line oracles; the sandbox three-way
// states are proven independently by TestLaunchBannerSandboxLine.
func renderBanner(host, dir string) string {
	var buf bytes.Buffer
	launchBanner(host, dir, false, bannerEnv(nil), lookMissing, &buf)
	return buf.String()
}

// twoWorkflowRepo builds a git repo holding exactly two top-level commissioned
// workflows, docs/dev and docs/user-testing. DiscoverWorkflows sorts the resolved
// absolute paths, so docs/dev precedes docs/user-testing; the relative-path order
// in the banner is therefore deterministic and is the independent expected value.
func twoWorkflowRepo(t *testing.T) string {
	t.Helper()
	repo := gitRepoFixture(t)
	commissionWorkflowAt(t, filepath.Join(repo, "docs", "dev"))
	commissionWorkflowAt(t, filepath.Join(repo, "docs", "user-testing"))
	return repo
}

// TestLaunchBannerConsistentFirstOfficerMetaphor (AC-1): line 1 frames the host as
// the first officer the launcher starts (`launching {host} as your first officer`),
// never the flipped `first officer launching` that read as spacedock-is-FO. The
// expected/forbidden phrases are written independently of the source constant, so a
// reverted-wording edit fails this assertion.
func TestLaunchBannerConsistentFirstOfficerMetaphor(t *testing.T) {
	for _, host := range []string{"claude", "codex"} {
		t.Run(host, func(t *testing.T) {
			out := renderBanner(host, t.TempDir())
			want := "launching " + host + " as your first officer"
			if !strings.Contains(out, want) {
				t.Fatalf("banner missing consistent line-1 framing %q: %q", want, out)
			}
			if strings.Contains(out, "first officer launching") {
				t.Fatalf("banner still carries the flipped line-1 framing %q: %q", "first officer launching", out)
			}
		})
	}
}

// TestLaunchBannerNoSelfServeStatus (AC-2): the banner tells the operator to ask
// the first officer (which runs status for them) rather than printing a self-serve
// `spacedock status` instruction, on all three workflow shapes (single / multi /
// none). The pitch is "the host is your first officer," so a `spacedock status`
// instruction would contradict it.
func TestLaunchBannerNoSelfServeStatus(t *testing.T) {
	single := gitRepoFixture(t)
	commissionWorkflowAt(t, filepath.Join(single, "docs", "dev"))
	cases := map[string]string{
		"single": single,
		"multi":  twoWorkflowRepo(t),
		"none":   gitRepoFixture(t),
	}
	for name, dir := range cases {
		t.Run(name, func(t *testing.T) {
			out := renderBanner("claude", dir)
			if strings.Contains(out, "spacedock status") {
				t.Fatalf("%s banner carries a self-serve `spacedock status` instruction: %q", name, out)
			}
			if !strings.Contains(out, "ask it for the queue") {
				t.Fatalf("%s banner does not point the operator to ask the first officer: %q", name, out)
			}
		})
	}
}

// TestLaunchBannerMultiWorkflowListsPaths (AC-3): on the multi-workflow path, line 2
// is exactly `Workflows: docs/dev docs/user-testing` — the detected paths, space
// separated, in discovery order — with no count, no "will help you pick", no
// "to pick)", and no `spacedock status`. The fixture's two relative paths are the
// independent expected value.
func TestLaunchBannerMultiWorkflowListsPaths(t *testing.T) {
	out := renderBanner("claude", twoWorkflowRepo(t))
	wantLine := "Workflows: " + filepath.Join("docs", "dev") + " " + filepath.Join("docs", "user-testing")
	if !lineEquals(out, wantLine) {
		t.Fatalf("multi-workflow line 2 is not the space-joined path list %q: %q", wantLine, out)
	}
	for _, forbidden := range []string{"found", "will help you pick", "to pick)", "spacedock status"} {
		if strings.Contains(out, forbidden) {
			t.Fatalf("multi-workflow banner contains forbidden token %q: %q", forbidden, out)
		}
	}
}

// TestLaunchBannerLabelAgreement (AC-4): the singular label `Workflow:` is used for
// the single and none shapes; the plural `Workflows:` only for the multi shape, so
// the label agrees with its content.
func TestLaunchBannerLabelAgreement(t *testing.T) {
	single := gitRepoFixture(t)
	commissionWorkflowAt(t, filepath.Join(single, "docs", "dev"))

	if got := bannerLine2(renderBanner("claude", single)); got != "Workflow: "+filepath.Join("docs", "dev") {
		t.Fatalf("single line 2 = %q, want singular Workflow: docs/dev", got)
	}
	if got := bannerLine2(renderBanner("claude", gitRepoFixture(t))); got != "Workflow: none detected" {
		t.Fatalf("none line 2 = %q, want %q", got, "Workflow: none detected")
	}
	multi := bannerLine2(renderBanner("claude", twoWorkflowRepo(t)))
	if !strings.HasPrefix(multi, "Workflows: ") {
		t.Fatalf("multi line 2 = %q, want plural Workflows: prefix", multi)
	}
}

// TestLaunchBannerSingleWorkflowGolden (AC-5): the single-workflow happy path is the
// exact three lines, in order, with no extra lines — a byte-exact golden, the
// strongest assertion, so any drift (an added line, a reworded line, a separator
// change) fails. The three expected lines are written independently of the source.
func TestLaunchBannerSingleWorkflowGolden(t *testing.T) {
	repo := gitRepoFixture(t)
	commissionWorkflowAt(t, filepath.Join(repo, "docs", "dev"))

	want := "spacedock " + displayVersion() + " · launching claude as your first officer\n" +
		"Workflow: " + filepath.Join("docs", "dev") + "\n" +
		"Sandbox: not wrapping this launch (no .safehouse profile)\n" +
		"claude is your first officer — ask it for the queue and next steps.\n"
	if got := renderBanner("claude", repo); got != want {
		t.Fatalf("single-workflow banner golden mismatch:\n got=%q\nwant=%q", got, want)
	}
}

// TestUnstampedVersionIsNotARelease (AC-6): the package default Version (the value
// before any ldflags stamp) is the `dev` sentinel and is NOT a release-shaped
// semver, so an unstamped `go build`/`go install` binary does not impersonate a
// release. A stamped release still overwrites Version via the ldflags path
// exercised by the release pipeline.
func TestUnstampedVersionIsNotARelease(t *testing.T) {
	if Version != "dev" {
		t.Fatalf("unstamped package Version = %q, want the dev sentinel %q", Version, "dev")
	}
	releaseShaped := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if releaseShaped.MatchString(Version) {
		t.Fatalf("unstamped Version %q matches the release semver pattern; it must read as a dev build", Version)
	}
}

// lineEquals reports whether out has a line equal (whole-line) to want.
func lineEquals(out, want string) bool {
	for _, line := range strings.Split(out, "\n") {
		if line == want {
			return true
		}
	}
	return false
}

// bannerLine2 returns the banner's second line (the Workflow/Workflows line).
func bannerLine2(out string) string {
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		return ""
	}
	return lines[1]
}

// TestPiBannerEmittedBeforeLaunch (AC-7): `spacedock pi` renders the same 3-line
// banner before launch — closing the pi no-orientation gap — and still reaches the
// launch seam. Driven through runPi with a ready stub runtime (no live host, no
// network) via the existing pi fixtures.
func TestPiBannerEmittedBeforeLaunch(t *testing.T) {
	repo := t.TempDir()
	writePiSkillFixtures(t, repo)
	pkg := t.TempDir()
	writePiSubagentsFixtures(t, pkg)
	ops := &fakePiRuntimeOps{
		lookPath:      piHealthyPathFixtures(),
		statOK:        statOKForPiResources(repo, pkg),
		packageStatus: healthyPiPackageStatus(),
	}
	var stdout, stderr bytes.Buffer

	code := runPi(context.Background(), []string{"--plugin-dir", repo}, t.TempDir(), piTestEnv(pkg, t.TempDir()), ops, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	out := stderr.String()
	for _, want := range []string{
		"launching pi as your first officer",
		"pi is your first officer — ask it for the queue",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("pi banner missing %q: %q", want, out)
		}
	}
	if len(ops.launched) == 0 {
		t.Fatalf("launch seam not reached after pi banner")
	}
}
