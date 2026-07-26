package release

import (
	"os"
	"path/filepath"
	"slices"
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

const releaseVersionStamp = "-X github.com/spacedock-dev/spacedock/internal/cli.Version={{ .Version }}"
const releaseMarkerStamp = "-X github.com/spacedock-dev/spacedock/internal/cli.releaseBuild=true"

func validReleaseIdentityStamps(config string) bool {
	var doc struct {
		Builds []struct {
			ID      string   `yaml:"id"`
			Ldflags []string `yaml:"ldflags"`
		} `yaml:"builds"`
	}
	if err := yaml.Unmarshal([]byte(config), &doc); err != nil {
		return false
	}

	for _, id := range []string{"spacedock-stable", "spacedock-edge"} {
		valid := false
		for _, build := range doc.Builds {
			if build.ID == id &&
				slices.Contains(build.Ldflags, releaseVersionStamp) &&
				slices.Contains(build.Ldflags, releaseMarkerStamp) {
				valid = true
			}
		}
		if !valid {
			return false
		}
	}
	return true
}

func TestGoreleaserBuildsStampReleaseIdentity(t *testing.T) {
	if !validReleaseIdentityStamps(readGoreleaserConfig(t)) {
		t.Fatal(".goreleaser.yaml must stamp Version and releaseBuild=true on stable and edge")
	}
}

func TestReleaseIdentityGuardRejectsChannelMarkerRemoval(t *testing.T) {
	config := readGoreleaserConfig(t)
	needle := "      - " + releaseMarkerStamp + "\n"
	for i, id := range []string{"spacedock-stable", "spacedock-edge"} {
		t.Run(id, func(t *testing.T) {
			at := strings.Index(config, needle)
			if i == 1 {
				at = strings.LastIndex(config, needle)
			}
			if at < 0 || strings.Count(config, needle) != 2 {
				t.Fatal("fixture must carry one marker per channel")
			}
			adversarial := config[:at] + config[at+len(needle):]
			if validReleaseIdentityStamps(adversarial) {
				t.Fatalf("removing %s marker did not trip guard", id)
			}
		})
	}
}

// goosOf returns the distinct goos tokens .goreleaser.yaml builds — the
// independent oracle the release.yml header must not contradict.
func goosOf(targets map[buildTarget]bool) []string {
	seen := map[string]bool{}
	for tgt := range targets {
		seen[tgt.os] = true
	}
	var oses []string
	for o := range seen {
		oses = append(oses, o)
	}
	sort.Strings(oses)
	return oses
}

// leadingCommentBlock returns the file's leading `#` comment lines (up to the
// first non-comment, non-blank line) joined into one string — the workflow
// header the doc-accuracy check inspects.
func leadingCommentBlock(content string) string {
	var header []string
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			header = append(header, line)
			continue
		}
		if trimmed == "" {
			continue
		}
		break
	}
	return strings.Join(header, "\n")
}

// goosMissingFromHeader returns the build OSes that the header text fails to
// name — the doc-vs-config drift the guard catches.
func goosMissingFromHeader(header string, oses []string) []string {
	var missing []string
	for _, o := range oses {
		if !strings.Contains(header, o) {
			missing = append(missing, o)
		}
	}
	return missing
}

// TestReleaseHeaderNamesEveryBuildOS locks AC-2: release.yml's file-header comment
// names every goos .goreleaser.yaml actually builds, so the header cannot claim a
// darwin-only build while the config cross-builds linux too. The oracle is the
// parsed build set (not header prose), so a header that drops `linux` reds even
// though the file still mentions `darwin`.
func TestReleaseHeaderNamesEveryBuildOS(t *testing.T) {
	oses := goosOf(parseGoreleaserBuildTargets(readGoreleaserConfig(t)))
	if len(oses) == 0 {
		t.Fatal("parsed no goos from .goreleaser.yaml; the header check has no oracle to bind")
	}
	header := leadingCommentBlock(readReleaseWorkflow(t))
	if header == "" {
		t.Fatal("release.yml has no leading comment header to check")
	}
	if missing := goosMissingFromHeader(header, oses); len(missing) > 0 {
		t.Errorf("release.yml header omits build OS %s while .goreleaser.yaml builds %s; header:\n%s",
			strings.Join(missing, ", "), strings.Join(oses, ", "), header)
	}
}

// TestReleaseHeaderGuardRejectsDarwinOnly proves the guard is load-bearing: a
// header with `linux` removed (the pre-task darwin-only shape) must fail the
// name-every-OS check against a config that still builds linux.
func TestReleaseHeaderGuardRejectsDarwinOnly(t *testing.T) {
	oses := goosOf(parseGoreleaserBuildTargets(readGoreleaserConfig(t)))
	header := leadingCommentBlock(readReleaseWorkflow(t))

	darwinOnly := strings.ReplaceAll(header, "linux", "")
	if darwinOnly == header {
		t.Fatal("release.yml header has no `linux` token to strip; the load-bearing check cannot bind")
	}
	if missing := goosMissingFromHeader(darwinOnly, oses); len(missing) == 0 {
		t.Fatalf("stripping `linux` from the header did not trip the guard; it is not load-bearing (oses=%s)",
			strings.Join(oses, ", "))
	}
}
