// ABOUTME: Proves the release-rendered Homebrew casks declare agentsview as a
// ABOUTME: cask dependency, declare no safehouse formula dependency, and keep
// ABOUTME: safehouse in the caveats with its brew install command.
package release

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// the two casks .goreleaser.yaml's homebrew_casks block emits per release, and
// the rendered .rb filename goreleaser writes each into.
var renderedCasks = []struct {
	name string
	file string
}{
	{name: "spacedock", file: "spacedock.rb"},
	{name: "spacedock@next", file: "spacedock@next.rb"},
}

// TestRenderedCaskDependsOnAgentsview locks AC-1 + AC-2: a release-rendered cask
// must carry a `depends_on cask:` stanza naming agentsview (AC-1); must carry NO
// `depends_on formula:` for the tap-qualified safehouse formula (the cross-tap
// dep the spike refuted); and its caveats must name safehouse with its
// `brew install eugene1g/safehouse/agent-safehouse` command and canonical site,
// but NOT agentsview, which is a dependency (AC-2). Both rendered casks (stable +
// edge) are checked. The expected values come from the AC, not from the file
// under test — the rendered .rb is goreleaser OUTPUT, so this reds if the
// dependency is dropped or mis-spelled, if the refuted formula dep is added, or
// if the caveat regresses.
//
// Primary proof renders the real Ruby via goreleaser snapshot. When goreleaser
// cannot run in the environment (not on PATH, sandbox denies its caches), the
// test falls back to parsing .goreleaser.yaml's homebrew_casks declarations,
// which assert the same properties at the config layer.
func TestRenderedCaskDependsOnAgentsview(t *testing.T) {
	casks, rendered := renderCasks(t)
	if !rendered {
		t.Log("goreleaser render unavailable; falling back to .goreleaser.yaml config parse")
		assertCaskConfig(t)
		return
	}
	for _, c := range renderedCasks {
		body := casks[c.name]
		if body == "" {
			t.Fatalf("rendered cask %q (%s) not found in goreleaser output", c.name, c.file)
		}
		assertRenderedCask(t, c.name, body)
	}
}

// assertRenderedCask checks one rendered cask's Ruby body against AC-1 and AC-2.
func assertRenderedCask(t *testing.T, name, body string) {
	t.Helper()

	if !strings.Contains(body, "depends_on cask:") {
		t.Errorf("rendered cask %q has no `depends_on cask:` stanza:\n%s", name, body)
	}
	if !strings.Contains(body, `"agentsview"`) {
		t.Errorf("rendered cask %q does not name agentsview as a dependency:\n%s", name, body)
	}

	// safehouse must NOT be a formula dependency: a cask depends_on formula on
	// the tap-qualified eugene1g/safehouse/agent-safehouse renders, but Homebrew
	// 6.0 tap-trust refuses to auto-load that untrusted third-party tap as a
	// transitive dep — the spike refuted it, so it must not ship. goreleaser
	// emits a formula dep as a `formula: [...]` key inside the depends_on stanza;
	// the tap-qualified name legitimately appears in the caveats install command,
	// so the depends_on-stanza form is the precise signal.
	if depBlock := dependsOnBlock(body); strings.Contains(depBlock, "formula:") {
		t.Errorf("rendered cask %q carries a formula dependency in its depends_on stanza; the cross-tap safehouse dep was refuted and must not ship:\n%s", name, depBlock)
	}

	caveats := caveatsBlock(body)
	if caveats == "" {
		t.Fatalf("rendered cask %q has no caveats block to inspect:\n%s", name, body)
	}
	assertSafehouseCaveat(t, name, caveats)
	if strings.Contains(caveats, "agentsview") {
		t.Errorf("rendered cask %q still lists agentsview in caveats; it moved to depends_on:\n%s", name, caveats)
	}
}

// assertSafehouseCaveat pins AC-2's corrected caveat content: safehouse stays a
// caveat, but with its real third-party-tap install command and canonical site,
// and the stale "not installed by brew" phrasing gone (safehouse IS brew-
// installable, just not as a transitive depends_on).
func assertSafehouseCaveat(t *testing.T, name, caveats string) {
	t.Helper()
	if !strings.Contains(caveats, "safehouse") {
		t.Errorf("cask %q caveats dropped safehouse:\n%s", name, caveats)
	}
	if !strings.Contains(caveats, "https://agent-safehouse.dev") {
		t.Errorf("cask %q caveats does not name the canonical safehouse site:\n%s", name, caveats)
	}
	if !strings.Contains(caveats, "brew install eugene1g/safehouse/agent-safehouse") {
		t.Errorf("cask %q caveats does not name safehouse's brew install command:\n%s", name, caveats)
	}
	if strings.Contains(caveats, "not installed by brew") {
		t.Errorf("cask %q caveats still says safehouse is \"not installed by brew\"; it is brew-installable from its tap:\n%s", name, caveats)
	}
}

// dependsOnBlock extracts the `depends_on …` stanza of a rendered cask up to the
// blank line that separates it from the next stanza, so a `formula:` check sees
// only the dependency declaration and not the tap-qualified name that also
// appears in the caveats install command.
func dependsOnBlock(body string) string {
	start := strings.Index(body, "depends_on")
	if start < 0 {
		return ""
	}
	rest := body[start:]
	if end := strings.Index(rest, "\n\n"); end >= 0 {
		return rest[:end]
	}
	return rest
}

// caveatsBlock extracts the heredoc body of a `caveats <<~EOS … EOS` stanza so
// the agentsview/safehouse assertions inspect only the companions text and not,
// say, a homepage URL elsewhere in the cask.
func caveatsBlock(body string) string {
	start := strings.Index(body, "caveats <<~EOS")
	if start < 0 {
		return ""
	}
	rest := body[start:]
	nl := strings.Index(rest, "\n")
	if nl < 0 {
		return ""
	}
	rest = rest[nl+1:]
	end := strings.Index(rest, "EOS")
	if end < 0 {
		return ""
	}
	return rest[:end]
}

// renderCasks runs `goreleaser release --snapshot --skip=publish --clean` from
// the repo root and returns each rendered cask's Ruby body keyed by cask name.
// rendered=false (the caller falls back to the config parse) when goreleaser is
// not on PATH or the render fails for an environmental reason — goreleaser is a
// release dependency, but CI sandboxes can deny the gem/module caches it needs.
func renderCasks(t *testing.T) (map[string]string, bool) {
	t.Helper()
	if _, err := exec.LookPath("goreleaser"); err != nil {
		return nil, false
	}
	repoRoot := filepath.Join("..", "..")
	cmd := exec.Command("goreleaser", "release", "--snapshot", "--skip=publish", "--clean")
	cmd.Dir = repoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("goreleaser snapshot render failed (falling back to config parse): %v\n%s", err, out)
		return nil, false
	}

	casks := map[string]string{}
	for _, c := range renderedCasks {
		data, err := os.ReadFile(filepath.Join(repoRoot, "dist", "homebrew", "Casks", c.file))
		if err != nil {
			t.Fatalf("goreleaser rendered no %s: %v", c.file, err)
		}
		casks[c.name] = string(data)
	}
	return casks, true
}

// goreleaserCaskConfig is the slice of homebrew_casks fields the config-parse
// fallback inspects: the declared cask-on-cask dependencies and the caveats text.
type goreleaserCaskConfig struct {
	HomebrewCasks []struct {
		Name         string `yaml:"name"`
		SkipUpload   string `yaml:"skip_upload"`
		Caveats      string `yaml:"caveats"`
		Dependencies []struct {
			Cask    string `yaml:"cask"`
			Formula string `yaml:"formula"`
		} `yaml:"dependencies"`
	} `yaml:"homebrew_casks"`
}

// TestCaskReleaseChannelRouting makes sure that a final tag cannot replace the
// edge cask and that a prerelease tag cannot replace the stable cask.
func TestCaskReleaseChannelRouting(t *testing.T) {
	var cfg goreleaserCaskConfig
	if err := yaml.Unmarshal([]byte(readGoreleaserConfig(t)), &cfg); err != nil {
		t.Fatalf("parse .goreleaser.yaml: %v", err)
	}
	want := map[string]string{
		"spacedock":      "auto",
		"spacedock@next": `{{ eq .Prerelease "" }}`,
	}
	for _, cask := range cfg.HomebrewCasks {
		if expected, ok := want[cask.Name]; ok {
			if cask.SkipUpload != expected {
				t.Errorf("cask %q skip_upload = %q, want %q", cask.Name, cask.SkipUpload, expected)
			}
			delete(want, cask.Name)
		}
	}
	for name := range want {
		t.Errorf("missing cask %q", name)
	}
}

// assertCaskConfig is the fallback proof: it parses .goreleaser.yaml's
// homebrew_casks declarations and asserts the same AC-1/AC-2 properties the
// render check proves on the emitted Ruby — each cask declares a `cask:
// agentsview` dependency, declares NO formula dependency (the refuted cross-tap
// safehouse dep), and a caveats block that keeps the corrected safehouse install
// line but drops agentsview.
func assertCaskConfig(t *testing.T) {
	t.Helper()
	var cfg goreleaserCaskConfig
	if err := yaml.Unmarshal([]byte(readGoreleaserConfig(t)), &cfg); err != nil {
		t.Fatalf("parse .goreleaser.yaml: %v", err)
	}
	if len(cfg.HomebrewCasks) != 2 {
		t.Fatalf("expected 2 homebrew_casks (stable + edge), got %d", len(cfg.HomebrewCasks))
	}
	for _, cask := range cfg.HomebrewCasks {
		hasAgentsview := false
		for _, dep := range cask.Dependencies {
			if dep.Cask == "agentsview" {
				hasAgentsview = true
			}
			if dep.Formula != "" {
				t.Errorf("cask %q declares a formula dependency %q; the cross-tap safehouse dep was refuted and must not ship", cask.Name, dep.Formula)
			}
		}
		if !hasAgentsview {
			t.Errorf("cask %q does not declare `cask: agentsview` as a dependency", cask.Name)
		}
		assertSafehouseCaveat(t, cask.Name, cask.Caveats)
		if strings.Contains(cask.Caveats, "agentsview") {
			t.Errorf("cask %q still lists agentsview in caveats; it moved to a dependency:\n%s", cask.Name, cask.Caveats)
		}
	}
}
