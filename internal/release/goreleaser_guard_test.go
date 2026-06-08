package release

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// buildTarget is a single goos/goarch build pair, the cross-product element a
// `v*` release tarball is published for.
type buildTarget struct {
	os   string
	arch string
}

func (t buildTarget) String() string { return t.os + "/" + t.arch }

// parseGoreleaserBuildTargets expands the goreleaser config's `builds[].goos` ×
// `builds[].goarch` cross-product into the flat set of build targets goreleaser
// produces — the same expansion goreleaser does when it emits one tarball per
// goos/goarch pair. A config that does not parse yields a nil set (the guard's
// superset assertion then fails loudly).
func parseGoreleaserBuildTargets(config string) map[buildTarget]bool {
	var doc struct {
		Builds []struct {
			Goos   []string `yaml:"goos"`
			Goarch []string `yaml:"goarch"`
		} `yaml:"builds"`
	}
	if err := yaml.Unmarshal([]byte(config), &doc); err != nil {
		return nil
	}
	targets := map[buildTarget]bool{}
	for _, b := range doc.Builds {
		for _, o := range b.Goos {
			for _, a := range b.Goarch {
				targets[buildTarget{os: o, arch: a}] = true
			}
		}
	}
	return targets
}

// TestGoreleaserBuildsLinuxAndDarwin locks AC-1: the release config's build
// target set is a superset of {linux,darwin}×{amd64,arm64}, so a `v*` release
// publishes the two linux tarballs alongside the two darwin ones. The expected
// set is written here independently of the YAML it checks — a future edit that
// drops `linux` from `builds.goos` reds this test, exactly as the AC requires.
func TestGoreleaserBuildsLinuxAndDarwin(t *testing.T) {
	got := parseGoreleaserBuildTargets(readGoreleaserConfig(t))

	want := []buildTarget{
		{os: "linux", arch: "amd64"},
		{os: "linux", arch: "arm64"},
		{os: "darwin", arch: "amd64"},
		{os: "darwin", arch: "arm64"},
	}
	var missing []string
	for _, target := range want {
		if !got[target] {
			missing = append(missing, target.String())
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Errorf(".goreleaser.yaml build target set is missing %s; got %s",
			strings.Join(missing, ", "), targetSetString(got))
	}
}

// TestGoreleaserBuildGuardRejectsDroppedLinux proves the guard is load-bearing:
// a config with `linux` removed from `builds.goos` must fail the superset check.
// The check parses the YAML rather than grepping its text, so the guard tracks
// what goreleaser actually builds, not what the file happens to mention.
func TestGoreleaserBuildGuardRejectsDroppedLinux(t *testing.T) {
	config := readGoreleaserConfig(t)
	// Drop the linux goos line the production config carries. The remaining
	// darwin-only config is the pre-task shape this task replaces.
	adversarial := strings.Replace(config, "      - linux\n", "", 1)
	if adversarial == config {
		t.Fatal("fixture .goreleaser.yaml has no `      - linux` goos line to drop")
	}

	got := parseGoreleaserBuildTargets(adversarial)
	if got[buildTarget{os: "linux", arch: "amd64"}] || got[buildTarget{os: "linux", arch: "arm64"}] {
		t.Fatalf("dropping the linux goos line still left a linux target in the set: %s", targetSetString(got))
	}
}

func targetSetString(targets map[buildTarget]bool) string {
	var names []string
	for target := range targets {
		names = append(names, target.String())
	}
	sort.Strings(names)
	return "{" + strings.Join(names, ", ") + "}"
}

func readGoreleaserConfig(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".goreleaser.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
