// ABOUTME: AC-3 scaffold-classifier proof — runs the survey scaffold detector artifact
// ABOUTME: against committed fixture repos and asserts each resolves to its label.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestSurveyScaffoldClassifier is the AC-3 detection-half proof. It runs the survey
// scaffold detector artifact in each committed fixture repo and asserts the
// emitted label matches the scaffold that fixture carries.
//
// The expected label for each case comes from the committed
// fixture file tree (testdata/survey/scaffolds/<name>) — an independent source. The
// classifier reads those files; if its detection logic regresses (a dropped superpowers
// skill-name check, a swapped gsd/superpowers branch, a missing generic fallback), the
// run over the fixture emits the wrong label and this test REDs. The proof is the
// EXECUTION of the detector artifact against known trees, never a substring over SKILL.md.
func TestSurveyScaffoldClassifier(t *testing.T) {
	scaffoldsRoot, err := filepath.Abs(filepath.Join("testdata", "survey", "scaffolds"))
	if err != nil {
		t.Fatal(err)
	}
	detector := filepath.Join(repoRoot(t), "skills", "survey", "bin", "detect-scaffold")

	cases := []struct {
		fixture  string
		wantLine string // the expected first non-marker output line, or a prefix for "similar:"
		exact    bool
	}{
		{fixture: "superpowers", wantLine: "superpowers", exact: true},
		{fixture: "gsd", wantLine: "gsd", exact: true},
		{fixture: "similar", wantLine: "similar:", exact: false}, // generic fallback names the dirs
		{fixture: "none", wantLine: "none", exact: true},
	}

	for _, tc := range cases {
		t.Run(tc.fixture, func(t *testing.T) {
			dir := filepath.Join(scaffoldsRoot, tc.fixture)
			if _, err := os.Stat(dir); err != nil {
				t.Fatalf("missing scaffold fixture %s: %v", dir, err)
			}
			cmd := exec.Command(detector)
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("run scaffold detector %s in %s: %v\n%s", detector, tc.fixture, err, out)
			}
			label := scaffoldLabel(string(out))
			t.Logf("%s -> %q", tc.fixture, label)
			if tc.exact {
				if label != tc.wantLine {
					t.Errorf("fixture %s classified as %q, want %q", tc.fixture, label, tc.wantLine)
				}
			} else {
				if !strings.HasPrefix(label, tc.wantLine) {
					t.Errorf("fixture %s classified as %q, want a line beginning %q", tc.fixture, label, tc.wantLine)
				}
			}
		})
	}
}

// scaffoldLabel pulls the detection block's emitted label: the first output line that
// is not the `## SCAFFOLD` marker or blank.
func scaffoldLabel(out string) string {
	for _, line := range strings.Split(out, "\n") {
		l := strings.TrimSpace(line)
		if l == "" || l == "## SCAFFOLD" {
			continue
		}
		return l
	}
	return ""
}
