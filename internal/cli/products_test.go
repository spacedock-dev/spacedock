package cli

import (
	"path/filepath"
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
