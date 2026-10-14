package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestExtensionTableAgreesWithInstallScript(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`(?s)\nEXTRAS='([^']*)'\n`).FindStringSubmatch(string(src))
	if m == nil {
		t.Fatal("install.sh declares no EXTRAS table")
	}
	rows := map[string][]string{}
	for _, line := range strings.Split(strings.TrimSpace(m[1]), "\n") {
		fields := strings.Split(line, "|")
		if len(fields) != 6 {
			t.Fatalf("EXTRAS row %q has %d fields, want 6", line, len(fields))
		}
		rows[fields[0]] = fields
	}
	if len(rows) != len(extensions) {
		t.Errorf("install.sh lists %d extensions and internal/cli lists %d", len(rows), len(extensions))
	}
	for _, ext := range extensions {
		fields, ok := rows[ext.name]
		if !ok {
			t.Errorf("%s is in internal/cli but not in install.sh", ext.name)
			continue
		}
		if (fields[5] == "tree") != ext.tree {
			t.Errorf("%s: install.sh kind %q, internal/cli tree=%v", ext.name, fields[5], ext.tree)
		}
		linked := false
		for _, link := range strings.Fields(fields[4]) {
			linked = linked || link[strings.Index(link, ":")+1:] == ext.program
		}
		if !linked {
			t.Errorf("%s: install.sh links %q, internal/cli looks for %q", ext.name, fields[4], ext.program)
		}
	}
}

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

func TestUpdateRemindsToRestartARunningCompassOnlyWhenItsProgramChanged(t *testing.T) {
	report := func(version string) versionReport {
		return versionReport{channel: "edge", rows: []pieceRow{{"compass program", version, ""}}}
	}
	const restart = "spacedock compass --port 4553 --stop && spacedock compass --port 4553 --daemon"

	var changed, unchanged bytes.Buffer
	printUpdateChanges(&changed, report("1.0.0"), report("1.0.1"), "", false)
	printUpdateChanges(&unchanged, report("1.0.1"), report("1.0.1"), "", false)
	if !strings.Contains(changed.String(), restart) {
		t.Errorf("update output lacks the restart line after Compass changed:\n%s", changed.String())
	}
	if strings.Contains(unchanged.String(), restart) {
		t.Errorf("update output tells the user to restart Compass when it did not change:\n%s", unchanged.String())
	}
}
