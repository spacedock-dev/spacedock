package cli

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestDoctorNamesAnUnreadableReviewChoiceAndShowsTheChoice(t *testing.T) {
	s := stageInstallWith(t, extensions)

	report := s.report()
	if len(report.problems) != 0 {
		t.Fatalf("a consistent install reported problems: %q", report.problems)
	}
	if got := report.version("review choice"); got != "not chosen" {
		t.Errorf("review choice = %q, want %q", got, "not chosen")
	}

	_, path := reviewSettingPath(s.env, s.project)
	mustMkdir(t, filepath.Dir(path))
	mustWrite(t, path, "{not json")
	problems := s.report().problems
	if len(problems) != 1 || !strings.HasPrefix(problems[0], "UNREADABLE: the review choice for this directory.") {
		t.Fatalf("problems = %q, want exactly the unreadable review choice", problems)
	}
}

func TestInstallSequencesEndWithSubspaceAndCompassForClaudeAndCodexOnEachChannel(t *testing.T) {
	for _, ch := range []struct{ name, source, branch string }{
		{"edge", "spacedock-dev/marketplace@edge", "next"},
		{"stable", "spacedock-dev/marketplace", "main"},
	} {
		t.Run(ch.name, func(t *testing.T) { testInstallSequencesEndWithSubspaceAndCompass(t, ch.source, ch.branch) })
	}
}

func testInstallSequencesEndWithSubspaceAndCompass(t *testing.T, source, branch string) {
	claude := installArgvSequence(source, branch)
	wantClaude := []installStep{
		{argv: []string{"plugin", "uninstall", "subspace@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "subspace@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "subspace@spacedock"}, tolerateExit: true, reportExit: true},
		{argv: []string{"plugin", "uninstall", "compass@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "compass@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "compass@spacedock"}, tolerateExit: true, reportExit: true},
	}
	codex := codexInstallArgvSequence(source, branch)
	wantCodex := []installStep{
		{argv: []string{"plugin", "remove", "subspace@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "add", "subspace@spacedock"}, tolerateExit: true, reportExit: true},
		{argv: []string{"plugin", "remove", "compass@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "add", "compass@spacedock"}, tolerateExit: true, reportExit: true},
	}
	for host, c := range map[string]struct{ got, want []installStep }{
		"claude": {claude, wantClaude},
		"codex":  {codex, wantCodex},
	} {
		if len(c.got) < len(c.want) {
			t.Fatalf("%s sequence = %v, want the product steps at its end", host, c.got)
		}
		if got := c.got[len(c.got)-len(c.want):]; !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s product steps = %v, want %v", host, got, c.want)
		}
	}
}

func TestDoctorOnHomebrewSaysCompassIsNotYetAvailableWhenItsProgramIsMissing(t *testing.T) {
	s := stageInstallWith(t, extensions)
	os.Remove(filepath.Join(s.first, "spacedock-compass"))
	report := collectVersions(versionInputs{env: s.env, dir: s.project, host: s.host, lookPath: s.lookPath,
		selfPath: "/opt/homebrew/Caskroom/spacedock/0.28.0/spacedock", selfVersion: "0.28.0-pre6"})
	want := "MISSING: the compass program is not on PATH. Compass is not yet available from Homebrew."
	if !slices.Contains(report.problems, want) || report.needsUpdate {
		t.Errorf("problems = %q, needsUpdate = %v, want %q and no `spacedock update` remedy", report.problems, report.needsUpdate, want)
	}
}
