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

	claude := installArgvSequence("spacedock-dev/marketplace@edge", "next")
	wantClaude := []installStep{
		{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}, tolerateExit: true, reportExit: true},
		{argv: []string{"plugin", "marketplace", "update", "spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "uninstall", "alpha@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "alpha@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "alpha@spacedock"}, tolerateExit: true, reportExit: true},
		{argv: []string{"plugin", "uninstall", "beta@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "install", "beta@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "update", "beta@spacedock"}, tolerateExit: true, reportExit: true},
	}
	if len(claude) < len(wantClaude) {
		t.Fatalf("claude sequence = %v, want the companion steps at its end", claude)
	}
	if got := claude[len(claude)-len(wantClaude):]; !reflect.DeepEqual(got, wantClaude) {
		t.Errorf("claude companion steps = %v, want %v", got, wantClaude)
	}

	codex := codexInstallArgvSequence("spacedock-dev/marketplace@edge", "next")
	wantCodex := []installStep{
		{argv: []string{"plugin", "marketplace", "add", "spacedock-dev/marketplace"}, tolerateExit: true},
		{argv: []string{"plugin", "marketplace", "upgrade", "spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "remove", "alpha@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "add", "alpha@spacedock"}, tolerateExit: true},
		{argv: []string{"plugin", "remove", "beta@spacedock-edge"}, tolerateExit: true},
		{argv: []string{"plugin", "add", "beta@spacedock"}, tolerateExit: true},
	}
	if got := codex[len(codex)-len(wantCodex):]; !reflect.DeepEqual(got, wantCodex) {
		t.Errorf("codex companion steps = %v, want %v", got, wantCodex)
	}
}

func TestExtensionPluginsInstallFromTheSpacedockMarketplaceWhateverTheChannel(t *testing.T) {
	saved := extensions
	extensions = fakeExtensions()
	defer func() { extensions = saved }()

	for _, devBranch := range []string{"next", "release-candidate"} {
		for _, id := range companionPluginIDs(devBranch) {
			if !strings.HasSuffix(id, "@spacedock") {
				t.Errorf("devBranch %q: companion id %q is not from the spacedock marketplace", devBranch, id)
			}
		}
		for _, step := range installArgvSequence(channelMarketplaceSource(devBranch), devBranch) {
			if step.argv[1] == "install" && strings.Contains(strings.Join(step.argv, " "), "alpha@spacedock-edge") {
				t.Errorf("devBranch %q: step %v installs a companion from the channel marketplace", devBranch, step.argv)
			}
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

func TestInstallSequencesNameNoCompanionOnStable(t *testing.T) {
	saved := extensions
	extensions = fakeExtensions()
	defer func() { extensions = saved }()

	sequences := map[string][]installStep{
		"claude": installArgvSequence("spacedock-dev/marketplace", "main"),
		"codex":  codexInstallArgvSequence("spacedock-dev/marketplace", "main"),
	}
	for host, steps := range sequences {
		for _, step := range steps {
			if strings.Contains(strings.Join(step.argv, " "), "alpha") || strings.Contains(strings.Join(step.argv, " "), "beta") {
				t.Errorf("%s sequence on stable installs an extension plugin: %v", host, step.argv)
			}
		}
	}
}

func TestInstallGoesOnWhenTheEdgeCopyIsNotThereToRemove(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("stub script uses /bin/sh; not portable to Windows")
	}
	saved := extensions
	extensions = fakeExtensions()
	defer func() { extensions = saved }()
	dir := writeClaudeStub(t, "plugin uninstall alpha@spacedock-edge")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	out, err := execHost{}.Install("claude", "spacedock-dev/marketplace", "next")
	if err != nil {
		t.Fatalf("Install stopped on a missing edge copy: %v\nout=%q", err, out)
	}
	for _, want := range []string{"plugin install alpha@spacedock", "plugin update alpha@spacedock"} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q after the failed removal\nout=%q", want, out)
		}
	}
	if strings.Contains(out, "failed") {
		t.Errorf("a missing edge copy is reported as a failure\nout=%q", out)
	}
}
