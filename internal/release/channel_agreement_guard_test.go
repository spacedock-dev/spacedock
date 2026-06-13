package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// The post-flip agreement invariant: the released stable channel's plugin source
// must settle on `main` across the two BINARY-side surfaces —
//   (1) release.yml's "Stamp plugin manifests" step git switch/push target,
//   (2) .goreleaser.yaml's stable-build cli.devBranch ldflag.
// Under Model B the marketplace manifest moved OUT of the plugin branch into a
// separate marketplace repo, so the former third surface — an in-branch
// .claude-plugin/marketplace.json source.ref that had to be re-settled per release
// — no longer exists. That removal is the AC-2 invariant guarded by
// TestPluginBranchCarriesNoMarketplaceManifest below; the channel pin now lives in
// the marketplace repo's manifest, not the plugin branch. The two binary surfaces
// are parsed out of two real artifacts authored by different changes, so a drift in
// either fails the check — an independent source of truth, not a re-read of the
// value the implementer wrote.
const stableChannelBranch = "main"

// releaseStampTarget extracts the branch the release.yml "Stamp plugin manifests
// to the release version" step switches to and pushes — the surface (1) value.
// It reads the step's run block, finds the `git switch <branch>` and `git push
// origin <branch>` commands, and returns the branch only when BOTH name the same
// branch (a switch/push split would itself be a drift). Returns "" when the step,
// or either command, is absent.
func releaseStampTarget(workflow string) string {
	for _, step := range parseWorkflowSteps(workflow) {
		if step.name != "Stamp plugin manifests to the release version" {
			continue
		}
		var switchTo, pushTo string
		for _, command := range executableShellCommands(step.run) {
			fields := strings.Fields(command)
			if len(fields) == 3 && fields[0] == "git" && fields[1] == "switch" {
				switchTo = fields[2]
			}
			if len(fields) == 4 && fields[0] == "git" && fields[1] == "push" && fields[2] == "origin" {
				pushTo = fields[3]
			}
		}
		if switchTo != "" && switchTo == pushTo {
			return switchTo
		}
		return ""
	}
	return ""
}

// goreleaserStableDevBranch extracts the cli.devBranch value the goreleaser
// STABLE build stamps — the surface (2) value. The stable build is the one whose
// id is spacedock-stable (it ldflags devBranch=main); the edge build
// (spacedock-edge) ldflags devBranch=next. The parse reads the build's ldflags
// for `-X …cli.devBranch=<value>` rather than grepping the file, so an edge-only
// or single-build config yields "" and reds the agreement check loudly.
func goreleaserStableDevBranch(config string) string {
	var doc struct {
		Builds []struct {
			ID      string   `yaml:"id"`
			Ldflags []string `yaml:"ldflags"`
		} `yaml:"builds"`
	}
	if err := yaml.Unmarshal([]byte(config), &doc); err != nil {
		return ""
	}
	for _, b := range doc.Builds {
		if b.ID != "spacedock-stable" {
			continue
		}
		return devBranchLdflag(b.Ldflags)
	}
	return ""
}

// goreleaserEdgeDevBranch extracts the cli.devBranch value the goreleaser EDGE
// build (spacedock-edge) stamps. The edge channel is the sole consumer of `next`
// post-flip, so the edge build must keep stamping `next` even after the stable
// build moves to `main` — a guard that the channels do not collapse to one value.
func goreleaserEdgeDevBranch(config string) string {
	var doc struct {
		Builds []struct {
			ID      string   `yaml:"id"`
			Ldflags []string `yaml:"ldflags"`
		} `yaml:"builds"`
	}
	if err := yaml.Unmarshal([]byte(config), &doc); err != nil {
		return ""
	}
	for _, b := range doc.Builds {
		if b.ID != "spacedock-edge" {
			continue
		}
		return devBranchLdflag(b.Ldflags)
	}
	return ""
}

// devBranchLdflag returns the value of the `-X …cli.devBranch=<value>` entry in
// ldflags, or "" when none is present.
func devBranchLdflag(ldflags []string) string {
	const marker = "cli.devBranch="
	for _, f := range ldflags {
		if i := strings.Index(f, marker); i >= 0 {
			return strings.TrimSpace(f[i+len(marker):])
		}
	}
	return ""
}

func readReleaseWorkflow(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestStableChannelBinaryPairAgreesOnMain locks AC-b's k6d-landable half: the two
// BINARY-side surfaces — release.yml's stamp target (1) and .goreleaser.yaml's
// stable-build devBranch (2) — must both be `main`. Each is parsed out of its own
// real artifact, authored by a different part of this change; a drift in either
// (a stamp step left on `next`, a stable build stamping `next`) reds this test.
// This is the half k6d makes green now, pre-flip, on `next`.
func TestStableChannelBinaryPairAgreesOnMain(t *testing.T) {
	stampTarget := releaseStampTarget(readReleaseWorkflow(t))
	if stampTarget == "" {
		t.Fatal("release.yml has no 'Stamp plugin manifests' step with matching git switch/push branch targets")
	}
	stableDevBranch := goreleaserStableDevBranch(readGoreleaserConfig(t))
	if stableDevBranch == "" {
		t.Fatal(".goreleaser.yaml has no spacedock-stable build with a cli.devBranch ldflag")
	}

	if stampTarget != stableChannelBranch {
		t.Errorf("release.yml stamp target = %q, want %q (stable channel binds the plugin version to main)", stampTarget, stableChannelBranch)
	}
	if stableDevBranch != stableChannelBranch {
		t.Errorf(".goreleaser.yaml stable-build cli.devBranch = %q, want %q", stableDevBranch, stableChannelBranch)
	}
	if stampTarget != stableDevBranch {
		t.Errorf("binary-side pair disagrees: release.yml stamp target %q != .goreleaser.yaml stable devBranch %q", stampTarget, stableDevBranch)
	}
}

// TestEdgeChannelStampsNext locks the channel-separation half: the edge build
// must keep stamping `next` even as the stable build moves to `main`, so the two
// channels resolve distinct plugin sources rather than collapsing to one branch.
func TestEdgeChannelStampsNext(t *testing.T) {
	edgeDevBranch := goreleaserEdgeDevBranch(readGoreleaserConfig(t))
	if edgeDevBranch == "" {
		t.Fatal(".goreleaser.yaml has no spacedock-edge build with a cli.devBranch ldflag")
	}
	if edgeDevBranch != "next" {
		t.Errorf(".goreleaser.yaml edge-build cli.devBranch = %q, want next (the edge channel stays on next)", edgeDevBranch)
	}
	stableDevBranch := goreleaserStableDevBranch(readGoreleaserConfig(t))
	if edgeDevBranch == stableDevBranch {
		t.Errorf("stable and edge builds both stamp devBranch=%q; the two channels collapsed to one source", edgeDevBranch)
	}
}

// TestPluginBranchCarriesNoMarketplaceManifest is the AC-2 git-state guard: the
// Model B decouple moves the marketplace manifest OUT of the plugin branch into a
// separate marketplace repo, so the plugin branch carries NO
// .claude-plugin/marketplace.json. With no in-branch manifest there is no
// per-release source.ref surface to re-settle (the divergence that kept
// `next → main` from being a clean fast-forward on that field is gone). The check
// is git state — the file's presence on disk, an independent fact, not a re-read
// of a value the implementer wrote. (The plugin's own .claude-plugin/plugin.json
// stays; only the marketplace.json moved.) This guards the main-bound branch; the
// next-branch manifest removal + full main/next alignment is the separate trunk
// reconcile.
func TestPluginBranchCarriesNoMarketplaceManifest(t *testing.T) {
	manifest := filepath.Join("..", "..", ".claude-plugin", "marketplace.json")
	if _, err := os.Stat(manifest); err == nil {
		t.Fatalf("%s is present on the plugin branch; Model B moves the marketplace manifest to the separate marketplace repo, so the plugin branch must carry no marketplace.json (and thus no per-release source.ref to re-settle)", manifest)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", manifest, err)
	}

	// The plugin manifest itself stays on the plugin branch — only the marketplace
	// manifest moved out.
	plugin := filepath.Join("..", "..", ".claude-plugin", "plugin.json")
	if _, err := os.Stat(plugin); err != nil {
		t.Fatalf("plugin manifest %s missing: %v (only marketplace.json should move out of the plugin branch)", plugin, err)
	}
}
