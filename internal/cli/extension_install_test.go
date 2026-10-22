package cli

import (
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestInstallSequencesFollowEachExtensionWithItsOwnPluginStep(t *testing.T) {
	for _, ch := range []struct {
		name, source, branch string
		addsMarketplace      bool
	}{
		{"edge", "spacedock-dev/marketplace@edge", "next", true},
		{"stable", "spacedock-dev/marketplace", "main", false},
	} {
		t.Run(ch.name, func(t *testing.T) {
			testInstallSequencesFollowEachExtension(t, ch.source, ch.branch, ch.addsMarketplace)
		})
	}
}

func testInstallSequencesFollowEachExtension(t *testing.T, source, branch string, addsMarketplace bool) {
	saved := extensions
	extensions = fakeExtensions()
	defer func() { extensions = saved }()

	claude := installArgvSequence(source, branch)
	wantClaude := []installStep{
		{argv: []string{"plugin", "uninstall", "alpha@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "alpha@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "alpha@spacedock"}, tolerateExit: true, reportExit: true},
		{argv: []string{"plugin", "uninstall", "beta@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "beta@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "beta@spacedock"}, tolerateExit: true, reportExit: true},
	}
	if addsMarketplace {
		wantClaude = append([]installStep{
			{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}, tolerateExit: true, reportExit: true},
			{argv: []string{"plugin", "marketplace", "update", "spacedock"}, tolerateExit: true},
		}, wantClaude...)
	}
	if len(claude) < len(wantClaude) {
		t.Fatalf("claude sequence = %v, want the companion steps at its end", claude)
	}
	if got := claude[len(claude)-len(wantClaude):]; !reflect.DeepEqual(got, wantClaude) {
		t.Errorf("claude companion steps = %v, want %v", got, wantClaude)
	}

	codex := codexInstallArgvSequence(source, branch)
	wantCodex := []installStep{
		{argv: []string{"plugin", "remove", "alpha@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "add", "alpha@spacedock"}, tolerateExit: true, reportExit: true},
		{argv: []string{"plugin", "remove", "beta@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "add", "beta@spacedock"}, tolerateExit: true, reportExit: true},
	}
	if addsMarketplace {
		wantCodex = append([]installStep{
			{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}, tolerateExit: true, reportExit: true},
			{argv: []string{"plugin", "marketplace", "upgrade", "spacedock"}, tolerateExit: true},
		}, wantCodex...)
	}
	if len(codex) < len(wantCodex) {
		t.Fatalf("codex sequence = %v, want the companion steps at its end", codex)
	}
	if got := codex[len(codex)-len(wantCodex):]; !reflect.DeepEqual(got, wantCodex) {
		t.Errorf("codex companion steps = %v, want %v", got, wantCodex)
	}

	wantAdds := 1
	if addsMarketplace {
		wantAdds = 2
	}
	for host, steps := range map[string][]installStep{"claude": claude, "codex": codex} {
		adds := 0
		for _, step := range steps {
			if len(step.argv) > 2 && step.argv[1] == "marketplace" && step.argv[2] == "add" {
				adds++
			}
		}
		if adds != wantAdds {
			t.Errorf("%s sequence adds %d marketplaces, want %d: on stable the channel marketplace already is the extension marketplace", host, adds, wantAdds)
		}
	}
}

func TestExtensionMarketplaceSourceOverrideReplacesOnlyTheExtensionMarketplace(t *testing.T) {
	saved := extensions
	extensions = fakeExtensions()
	savedOverride := extensionMarketplaceSourceOverride
	defer func() { extensions, extensionMarketplaceSourceOverride = saved, savedOverride }()

	applyMarketplaceSourceOverride([]string{"SPACEDOCK_EXTENSION_MARKETPLACE_SOURCE=/tmp/staged"})
	var adds []string
	for _, step := range installArgvSequence("spacedock-dev/marketplace@edge", "next") {
		if len(step.argv) == 4 && step.argv[1] == "marketplace" && step.argv[2] == "add" {
			adds = append(adds, step.argv[3])
		}
	}
	if want := []string{"spacedock-dev/marketplace@edge", "/tmp/staged"}; !reflect.DeepEqual(adds, want) {
		t.Errorf("marketplace adds = %v, want %v", adds, want)
	}
}

func TestInstallNamesAnExtensionPluginThatCannotBeUpdated(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	saved := extensions
	extensions = fakeExtensions()
	defer func() { extensions = saved }()
	dir := writeClaudeStub(t, "plugin update beta@spacedock")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("claude", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install returned error on a failed companion: %v\nout=%q", err, out)
	}
	if want := "claude plugin update beta@spacedock failed"; !strings.Contains(out, want) {
		t.Errorf("output does not name the failed companion, want %q\nout=%q", want, out)
	}
	if strings.Contains(out, "alpha@spacedock failed") {
		t.Errorf("output blames alpha, which updated fine\nout=%q", out)
	}
}

func TestCodexInstallNamesAnExtensionPluginThatCannotBeAdded(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	saved := extensions
	extensions = fakeExtensions()
	defer func() { extensions = saved }()
	dir := writeHostStub(t, "codex", "plugin add beta@spacedock")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("codex", "spacedock-dev/marketplace@edge", "next")
	if err != nil {
		t.Fatalf("Install returned error on a failed companion: %v\nout=%q", err, out)
	}
	if want := "codex plugin add beta@spacedock failed"; !strings.Contains(out, want) {
		t.Errorf("output does not name the failed companion, want %q\nout=%q", want, out)
	}
	if strings.Contains(out, "alpha@spacedock failed") {
		t.Errorf("output blames alpha, which was added fine\nout=%q", out)
	}
}
