package cli

import (
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestInstallSequencesFollowEachExtensionWithItsOwnPluginStep(t *testing.T) {
	saved := extensions
	extensions = fakeExtensions()
	defer func() { extensions = saved }()

	claude := installArgvSequence("spacedock-dev/marketplace", "next")
	if len(claude) < 4 {
		t.Fatalf("claude sequence = %v, want the companion steps at its end", claude)
	}
	wantClaude := []installStep{
		{argv: []string{"plugin", "install", "alpha@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "alpha@spacedock-edge"}, tolerateExit: true, reportExit: true},
		{argv: []string{"plugin", "install", "beta@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "beta@spacedock-edge"}, tolerateExit: true, reportExit: true},
	}
	if got := claude[len(claude)-4:]; !reflect.DeepEqual(got, wantClaude) {
		t.Errorf("claude companion steps = %v, want %v", got, wantClaude)
	}

	codex := codexInstallArgvSequence("spacedock-dev/marketplace", "main")
	wantCodex := []installStep{
		{argv: []string{"plugin", "add", "alpha@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "add", "beta@spacedock"}, tolerateExit: true},
	}
	if got := codex[len(codex)-2:]; !reflect.DeepEqual(got, wantCodex) {
		t.Errorf("codex companion steps = %v, want %v", got, wantCodex)
	}
}

func TestInstallSequencesNameNoCompanionWhenNoExtensionIsListed(t *testing.T) {
	saved := extensions
	extensions = nil
	defer func() { extensions = saved }()

	for _, step := range installArgvSequence("spacedock-dev/marketplace", "next") {
		if step.reportExit {
			t.Errorf("step %v reports a companion, but no extension is listed", step.argv)
		}
	}
}

func TestInstallNamesAnExtensionPluginThatCannotBeUpdated(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	saved := extensions
	extensions = fakeExtensions()
	defer func() { extensions = saved }()
	dir := writeClaudeStub(t, "plugin update beta@spacedock-edge")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("claude", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install returned error on a failed companion: %v\nout=%q", err, out)
	}
	if want := "claude plugin update beta@spacedock-edge failed"; !strings.Contains(out, want) {
		t.Errorf("output does not name the failed companion, want %q\nout=%q", want, out)
	}
	if strings.Contains(out, "alpha@spacedock-edge failed") {
		t.Errorf("output blames alpha, which updated fine\nout=%q", out)
	}
}
