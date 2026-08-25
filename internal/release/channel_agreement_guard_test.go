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
//	(1) release.yml's main-stamp step git switch/push target,
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

// releaseStampTarget extracts the branch that release.yml's main-stamp step
// switches to and pushes — the surface (1) value. It reads the step's run block
// and finds the `git switch <branch>` and `git push origin <branch>` commands.
// It gives the branch only when BOTH commands name the same branch, because a
// switch/push split is itself a drift. It gives "" when the step, or either
// command, is absent. A refspec push (a target that contains `:`) is skipped, so
// a ref-advance push cannot shadow the bare-branch stamp target.
func releaseStampTarget(workflow string) string {
	for _, step := range parseWorkflowSteps(workflow) {
		if step.name != mainStampStepName {
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
		t.Fatalf("release.yml has no %q step with matching git switch/push branch targets", mainStampStepName)
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

// stableRefPushSource gives the push source that the stable-advance step moves
// the stable channel ref to — the part before `:refs/heads/stable`, with any
// shell quotes stripped. It gives "" when the step has no such push. It looks
// for a `git push origin <src>:refs/heads/stable` command in the step's run
// block.
func stableRefPushSource(workflow string) string {
	for _, step := range parseWorkflowSteps(workflow) {
		if step.name != stableAdvanceStepName {
			continue
		}
		for _, command := range executableShellCommands(step.run) {
			fields := strings.Fields(command)
			if len(fields) != 4 || fields[0] != "git" || fields[1] != "push" || fields[2] != "origin" {
				continue
			}
			refspec := strings.Trim(fields[3], `"'`)
			if src, ok := strings.CutSuffix(refspec, ":refs/heads/stable"); ok {
				return src
			}
		}
	}
	return ""
}

// TestStampStepAdvancesStableRefToTaggedCommit locks the stable-channel publish
// mechanism AND the tag-binding: the stable-advance step MUST push the TAGGED commit
// ($RELEASE_COMMIT) to the `stable` ref, because the spacedock-dev/marketplace
// stable entry pins source.ref=stable. Without this push the stable channel would
// freeze at the prior release forever (a fresh `spacedock@spacedock` install would
// resolve the old commit). Pushing the tagged SHA rather than `main` keeps stable
// and the tag the SAME commit even if main advanced after the tag fired — the
// former `main:refs/heads/stable` form could point stable at a different commit
// than the tag. The command is parsed out of the real release.yml, so dropping the
// push, or regressing to a `main` source, reds this.
func TestStampStepAdvancesStableRefToTaggedCommit(t *testing.T) {
	src := stableRefPushSource(readReleaseWorkflow(t))
	if src == "" {
		t.Fatal("release.yml stable-advance step does not push to refs/heads/stable; the stable marketplace channel (source.ref=stable) would never advance past the prior release")
	}
	if src != "$RELEASE_COMMIT" {
		t.Errorf("stable ref is advanced to %q, not the tagged commit $RELEASE_COMMIT; stable and the tag could diverge if main advances after the tag fires", src)
	}
}

// TestStableRefGuardRejectsMainSource is the adversarial twin: regress the stable
// push back to the divergeable `main:refs/heads/stable` form and the guard must
// RED, because that source resolves to whatever main HEAD is at push time rather
// than the tagged commit.
func TestStableRefGuardRejectsMainSource(t *testing.T) {
	workflow := readReleaseWorkflow(t)
	adversarial := strings.Replace(workflow,
		`git push origin "$RELEASE_COMMIT:refs/heads/stable"`,
		`git push origin main:refs/heads/stable`,
		1)
	if adversarial == workflow {
		t.Fatal("fixture release.yml missing the `$RELEASE_COMMIT:refs/heads/stable` push to regress")
	}
	if src := stableRefPushSource(adversarial); src == "$RELEASE_COMMIT" {
		t.Fatal("guard still saw the tagged-commit source after regressing to a main source")
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
