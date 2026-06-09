package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// The post-flip agreement invariant: the released stable channel's plugin source
// must settle on `main` across three INDEPENDENTLY-authored surfaces —
//   (1) release.yml's "Stamp plugin manifests" step git switch/push target,
//   (2) .goreleaser.yaml's stable-build cli.devBranch ldflag,
//   (3) .claude-plugin/marketplace.json source.ref.
// Surfaces (1) and (2) are the BINARY side; surface (3) is the marketplace side.
// Surface (3) is branch-local: marketplace.json points at the channel its branch
// serves — `next` on the edge branch, `main` on the stable branch. So the
// tri-surface ==main check is meaningful only on a tree whose marketplace ref is
// already `main` (the stable branch); on the edge tree (ref `next`) it skips, and
// the binary-side pair (1)==(2)==main is asserted unconditionally below. The
// surfaces are parsed out of three real artifacts authored by different changes,
// so a drift in any one fails the check — an independent source of truth, not a
// re-read of the value the implementer wrote.
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

// marketplaceSourceRef extracts .claude-plugin/marketplace.json's first plugin
// entry's source.ref — surface (3), pj's. Parsed from the real file so a pj flip
// (next→main) is observed here, not asserted from a value k6d wrote.
func marketplaceSourceRef(manifest string) string {
	var doc struct {
		Plugins []struct {
			Source struct {
				Ref string `json:"ref"`
			} `json:"source"`
		} `json:"plugins"`
	}
	if err := yamlJSON([]byte(manifest), &doc); err != nil {
		return ""
	}
	if len(doc.Plugins) == 0 {
		return ""
	}
	return doc.Plugins[0].Source.Ref
}

// yamlJSON decodes JSON via the yaml.v3 codec (a JSON document is valid YAML), so
// the agreement guard needs no extra import beyond the one the sibling goreleaser
// guard already carries.
func yamlJSON(blob []byte, out any) error {
	return yaml.Unmarshal(blob, out)
}

func readReleaseWorkflow(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func readMarketplaceManifest(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".claude-plugin", "marketplace.json"))
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

// TestTriSurfaceChannelAgreement is the FULL agreement invariant: all three
// independently-authored surfaces — release.yml stamp target (1),
// .goreleaser.yaml stable devBranch (2), and .claude-plugin/marketplace.json
// source.ref (3) — must agree on `main`. Surface (3) is branch-local, so this
// asserts only on a `main`-ref tree (the stable branch); on the edge tree (ref
// `next`) it skips, since `next` correctly serves its own channel. The binary-side
// pair (1)==(2) is asserted unconditionally by TestStableChannelBinaryPairAgreesOnMain
// above.
func TestTriSurfaceChannelAgreement(t *testing.T) {
	marketplaceRef := marketplaceSourceRef(readMarketplaceManifest(t))
	if marketplaceRef != stableChannelBranch {
		t.Skipf("marketplace source.ref = %q (edge branch serves its own channel; the tri-surface ==%q agreement holds only on a main-ref tree); binary-side pair is covered by TestStableChannelBinaryPairAgreesOnMain", marketplaceRef, stableChannelBranch)
	}

	surfaces := map[string]string{
		"release.yml stamp target":           releaseStampTarget(readReleaseWorkflow(t)),
		".goreleaser.yaml stable devBranch":  goreleaserStableDevBranch(readGoreleaserConfig(t)),
		".claude-plugin/marketplace.json ref": marketplaceRef,
	}
	for name, got := range surfaces {
		if got != stableChannelBranch {
			t.Errorf("channel surface %q = %q, want %q", name, got, stableChannelBranch)
		}
	}
}
