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
//
//	(1) release.yml's "Stamp plugin manifests" step git switch/push target,
//	(2) .goreleaser.yaml's stable-build cli.devBranch ldflag.
//
// Under Model B the marketplace manifest will move OUT of the plugin branch into a
// separate marketplace repo, retiring the former third surface — an in-branch
// .claude-plugin/marketplace.json source.ref re-settled per release. That removal is
// deferred until after the v0.20.1 cutover (the released v0.20.0 binary still resolves
// its install from main's marketplace.json), so a transitional bridge manifest stays on
// main meanwhile, guarded by TestMainCarriesMarketplaceBridgeManifest below; the channel
// pin for the new binary lives in the marketplace repo's manifest. The two binary surfaces
// are parsed out of two real artifacts authored by different changes, so a drift in
// either fails the check — an independent source of truth, not a re-read of the
// value the implementer wrote.
const stableChannelBranch = "main"

// releaseStampTarget extracts the branch the release.yml "Stamp plugin manifests
// to the release version" step switches to and pushes — the surface (1) value.
// It reads the step's run block, finds the `git switch <branch>` and `git push
// origin <branch>` commands, and returns the branch only when BOTH name the same
// branch (a switch/push split would itself be a drift). Returns "" when the step,
// or either command, is absent. The step also pushes the stable channel ref
// (`git push origin main:refs/heads/stable`); that refspec push (target contains
// `:`) is skipped here so it does not shadow the bare-branch stamp target.
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
			if len(fields) == 4 && fields[0] == "git" && fields[1] == "push" && fields[2] == "origin" && !strings.Contains(fields[3], ":") {
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

// stampStepAdvancesStableRef reports whether the "Stamp plugin manifests" step
// pushes the stamped commit to the stable channel ref. It looks for a
// `git push origin <src>:refs/heads/stable` command in the step's run block.
func stampStepAdvancesStableRef(workflow string) bool {
	for _, step := range parseWorkflowSteps(workflow) {
		if step.name != "Stamp plugin manifests to the release version" {
			continue
		}
		for _, command := range executableShellCommands(step.run) {
			fields := strings.Fields(command)
			if len(fields) == 4 && fields[0] == "git" && fields[1] == "push" && fields[2] == "origin" && strings.HasSuffix(fields[3], ":refs/heads/stable") {
				return true
			}
		}
	}
	return false
}

// TestStampStepAdvancesStableRef locks the stable-channel publish mechanism: the
// release stamp step MUST push the release commit to the `stable` ref, because the
// spacedock-dev/marketplace stable entry pins source.ref=stable. Without this push
// the stable channel would freeze at the prior release forever (a fresh
// `spacedock@spacedock` install would resolve the old commit), since the marketplace
// manifest is intentionally static and no longer hand-edited per release. The
// command is parsed out of the real release.yml, so dropping the push reds this.
func TestStampStepAdvancesStableRef(t *testing.T) {
	if !stampStepAdvancesStableRef(readReleaseWorkflow(t)) {
		t.Error("release.yml stamp step does not push to refs/heads/stable; the stable marketplace channel (source.ref=stable) would never advance past the prior release")
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

// TestMainCarriesMarketplaceBridgeManifest guards the transitional bridge that keeps
// the released v0.20.0 binary's install path working. That binary resolves its
// install from main's .claude-plugin/marketplace.json (`claude plugin marketplace
// add spacedock-dev/spacedock@main`). Model B's removal of this manifest — it moves
// to the standalone marketplace repo, retiring the per-release source.ref re-settle —
// is deferred until the v0.20.1 cutover ships the binary that points at that repo and
// existing v0.20.0 installs have migrated. Until then the bridge MUST stay on main, or
// v0.20.0 installs break, so this guards its PRESENCE (git state — the file on disk,
// an independent fact). The plugin's own .claude-plugin/plugin.json stays alongside it.
func TestMainCarriesMarketplaceBridgeManifest(t *testing.T) {
	manifest := filepath.Join("..", "..", ".claude-plugin", "marketplace.json")
	if _, err := os.Stat(manifest); err != nil {
		t.Fatalf("bridge marketplace manifest %s missing: %v — the released v0.20.0 binary resolves its install from main's marketplace.json; removing it before the v0.20.1 cutover breaks v0.20.0 installs", manifest, err)
	}

	plugin := filepath.Join("..", "..", ".claude-plugin", "plugin.json")
	if _, err := os.Stat(plugin); err != nil {
		t.Fatalf("plugin manifest %s missing: %v", plugin, err)
	}
}

// TestChannelSurfacesDoNotDivergeAfterDecouple re-expresses the old tri-surface
// agreement invariant for Model B. Before the decouple, the channel a release
// served had THREE surfaces that had to agree (release.yml stamp target,
// goreleaser stable devBranch, and an in-branch marketplace.json source.ref), and
// a drift on the manifest ref was a real per-release re-settle hazard. The decouple
// will retire the in-branch ref surface (deferred post-cutover; a transitional bridge
// manifest stays on main meanwhile, guarded above), so the new binary's channel is
// determined SOLELY by the binary's devBranch stamp selecting a marketplace ENTRY
// NAME (stable=main→spacedock, edge=next→spacedock-edge). The surviving invariant —
// independent values that CAN disagree, so not a tautology — is that the two
// channel devBranch stamps are both present and DISTINCT: if the stable and edge
// builds collapsed to one devBranch, both channels would select the same entry and
// the stable/edge split would vanish. The two values are parsed out of the real
// .goreleaser.yaml builds, so a config that drops the edge build, or sets both to
// the same branch, reds here.
func TestChannelSurfacesDoNotDivergeAfterDecouple(t *testing.T) {
	config := readGoreleaserConfig(t)
	stable := goreleaserStableDevBranch(config)
	edge := goreleaserEdgeDevBranch(config)
	if stable == "" {
		t.Fatal(".goreleaser.yaml has no spacedock-stable build with a cli.devBranch ldflag")
	}
	if edge == "" {
		t.Fatal(".goreleaser.yaml has no spacedock-edge build with a cli.devBranch ldflag")
	}
	if stable == edge {
		t.Fatalf("stable and edge builds both stamp devBranch=%q; post-decouple the channel IS the devBranch-selected entry name, so identical stamps collapse the two channels onto one marketplace entry", stable)
	}
}
