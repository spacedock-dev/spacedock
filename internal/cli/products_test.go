package cli

import (
	"path/filepath"
	"reflect"
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

func TestEdgeInstallSequencesEndWithSubspaceAndCompassForClaudeAndCodex(t *testing.T) {
	claude := installArgvSequence("spacedock-dev/marketplace@edge", "next")
	wantClaude := []installStep{
		{argv: []string{"plugin", "uninstall", "subspace@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "subspace@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "subspace@spacedock"}, tolerateExit: true, reportExit: true},
		{argv: []string{"plugin", "uninstall", "compass@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "compass@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "compass@spacedock"}, tolerateExit: true, reportExit: true},
	}
	codex := codexInstallArgvSequence("spacedock-dev/marketplace@edge", "next")
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
