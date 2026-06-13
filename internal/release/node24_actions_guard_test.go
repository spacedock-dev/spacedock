// ABOUTME: Guards that every GitHub Actions `uses:` pin sits at or above the
// ABOUTME: node24 major GitHub forces from 2026-06-16, so no workflow regresses.
package release

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// node24MinMajor is the independent oracle: the lowest major release of each
// action that runs on node24. Sourced from GitHub's runner-deprecation changelog
// (revised deadline 2026-06-16) and each action's release notes, NOT from the
// workflow files this test checks — so a pin that drops below its minimum reds.
// Actions absent from this map (deploy-pages, upload-pages-artifact) have no
// node24 release published yet and are left at their current node20 pins; they
// sit off the release-cut path (docs.yml only) and are not enforced here.
var node24MinMajor = map[string]int{
	"actions/checkout":             5,
	"actions/setup-go":             6,
	"actions/setup-node":           6,
	"actions/setup-python":         6,
	"actions/upload-artifact":      5,
	"goreleaser/goreleaser-action": 7,
}

// usesPin matches a workflow `uses: owner/action@vN` step ref and captures the
// `owner/action` slug and the integer major from a `@vN` tag. A `@<sha>` or
// `@<branch>` pin (no `vN`) does not match — those carry no major to compare.
var usesPin = regexp.MustCompile(`uses:\s*([\w.-]+/[\w.-]+)@v(\d+)`)

// actionPin is one parsed `uses:` occurrence with its source location, so a
// failure names the exact workflow file and line that regressed.
type actionPin struct {
	file  string
	line  int
	slug  string
	major int
}

// parseWorkflowActionPins extracts every `uses: owner/action@vN` pin from a
// workflow file's text, with the line number of each, so the guard can report
// precisely where a sub-minimum pin lives.
func parseWorkflowActionPins(file, content string) []actionPin {
	var pins []actionPin
	for i, line := range strings.Split(content, "\n") {
		if m := usesPin.FindStringSubmatch(line); m != nil {
			major, _ := strconv.Atoi(m[2])
			pins = append(pins, actionPin{file: file, line: i + 1, slug: m[1], major: major})
		}
	}
	return pins
}

// readWorkflowFiles returns the path + text of every `.github/workflows/*.yml`
// file, the full set the deprecation guard must cover.
func readWorkflowFiles(t *testing.T) map[string]string {
	t.Helper()
	dir := filepath.Join("..", "..", ".github", "workflows")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yml") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		files[e.Name()] = string(data)
	}
	if len(files) == 0 {
		t.Fatalf("no .github/workflows/*.yml files found under %s", dir)
	}
	return files
}

// belowNode24Minimum returns the pins in the given set that are pinned below
// their recorded node24-minimum major. Pins for actions absent from the oracle
// map (the exempt pages actions) are ignored.
func belowNode24Minimum(pins []actionPin) []actionPin {
	var below []actionPin
	for _, p := range pins {
		min, tracked := node24MinMajor[p.slug]
		if tracked && p.major < min {
			below = append(below, p)
		}
	}
	return below
}

// TestNode24ActionsPinnedAtMinimum locks AC-5: every workflow's `uses:` pin for
// a node24-tracked action sits at or above the major GitHub forces to node24 on
// 2026-06-16. A re-pin below the minimum (e.g. checkout back to @v4) reds here.
func TestNode24ActionsPinnedAtMinimum(t *testing.T) {
	var all []actionPin
	for name, content := range readWorkflowFiles(t) {
		all = append(all, parseWorkflowActionPins(name, content)...)
	}

	// Every tracked action must actually appear somewhere, so the guard fails if
	// a workflow stops using an action the oracle still expects to find pinned.
	seen := map[string]bool{}
	for _, p := range all {
		seen[p.slug] = true
	}
	for slug := range node24MinMajor {
		if !seen[slug] {
			t.Errorf("no `uses: %s@vN` pin found in any workflow; the node24 guard is not covering it", slug)
		}
	}

	for _, p := range belowNode24Minimum(all) {
		t.Errorf("%s:%d pins %s@v%d, below the node24 minimum @v%d (node-20 forced off 2026-06-16)",
			p.file, p.line, p.slug, p.major, node24MinMajor[p.slug])
	}
}

// TestNode24GuardRejectsRevertedPin proves the guard is load-bearing: rewriting
// one node24-pinned `uses:` line back to @v4 must make belowNode24Minimum flag
// it. The check compares parsed majors against the in-test oracle, so it tracks
// the actual pinned version, not whatever the file happens to mention.
func TestNode24GuardRejectsRevertedPin(t *testing.T) {
	const file = "release.yml"
	content := readWorkflowFiles(t)[file]

	reverted := strings.Replace(content, "uses: actions/checkout@v5", "uses: actions/checkout@v4", 1)
	if reverted == content {
		t.Fatalf("fixture %s has no `uses: actions/checkout@v5` line to revert", file)
	}

	below := belowNode24Minimum(parseWorkflowActionPins(file, reverted))
	if len(below) == 0 {
		t.Fatal("reverting checkout to @v4 did not trip belowNode24Minimum; the guard is not load-bearing")
	}
	found := false
	for _, p := range below {
		if p.slug == "actions/checkout" && p.major == 4 {
			found = true
		}
	}
	if !found {
		t.Errorf("expected the reverted actions/checkout@v4 pin to be flagged; got %v", below)
	}
}
