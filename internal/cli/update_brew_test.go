package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

func TestBrewRepairSelectsOneStableExtensionEntry(t *testing.T) {
	entry := `{"name":"compass","version":"3.0.1","source":{"source":"git-subdir","url":"https://github.com/spacedock-dev/compass.git","path":"compass","ref":"stable"}}`
	for _, tc := range []struct {
		name, entries string
		ok            bool
	}{
		{"stable", entry, true},
		{"duplicate", entry + "," + entry, false},
		{"prerelease", strings.Replace(entry, "3.0.1", "3.0.1-beta.1", 1), false},
		{"ambiguous path", strings.Replace(entry, `"path":"compass"`, `"path":"../compass"`, 1), false},
		{"unverified source", strings.Replace(entry, "https://", "file://", 1), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveRepairEntry([]byte(`{"name":"spacedock","plugins":[`+tc.entries+`]}`), "compass")
			if (err == nil) != tc.ok || tc.ok && got.Version != "3.0.1" {
				t.Fatalf("resolve = %+v, %v", got, err)
			}
		})
	}
}

func TestBrewRepairRejectsAnExpandedOrUnknownPlan(t *testing.T) {
	good := "==> Would upgrade 1 requested outdated package\nspacedock-dev/tap/compass 2.13.0 -> 3.0.1\n"
	for _, tc := range []struct {
		output string
		ok     bool
	}{
		{good, true},
		{good + "==> Would upgrade 1 dependent\nother 1.0.0 -> 1.0.1\n", false},
		{good + "==> Installing 1 dependency\nother\n", false},
		{"", false},
		{strings.ReplaceAll(good, "3.0.1", "3.1.0"), false},
	} {
		if got := singleBrewUpgrade(tc.output, "spacedock-dev/tap/compass", "3.0.1"); got != tc.ok {
			t.Errorf("accepted %q = %v, want %v", tc.output, got, tc.ok)
		}
	}
}

func TestBrewRepairUpdate(t *testing.T) {
	for _, tc := range []struct {
		name, cask, claude, codex, failure, want, changes string
		code                                              int
	}{
		{"already mixed", "spacedock", "3.0.1", "3.0.1", "", "compass 3.0.1 verified", "brew", 0},
		{"edge cask", "spacedock@next", "2.13.0", "2.13.0", "", "compass 3.0.1 verified", "brew\nclaude\ncodex", 0},
		{"archive unavailable", "spacedock", "2.13.0", "2.13.0", "missing-archive", "Kept compass (program 2.13.0, Claude Code plugin 2.13.0, Codex plugin 2.13.0)", "", 1},
		{"corrupt archive", "spacedock", "2.13.0", "2.13.0", "bad-archive", "failed its SHA-256 check", "", 1},
		{"stale host snapshot", "spacedock", "2.13.0", "2.13.0", "stale-codex", "Codex plugin 2.13.0). codex's spacedock snapshot cannot select compass 3.0.1", "", 1},
		{"expanded Brew plan", "spacedock", "2.13.0", "2.13.0", "extra-upgrade", "did not preview only", "", 1},
		{"host failure and retry", "spacedock", "2.13.0", "2.13.0", "fail-codex", "program 3.0.1, Claude Code plugin 3.0.1, Codex plugin 2.13.0). Codex:", "brew\nclaude", 1},
		{"Brew failure after active link", "spacedock", "2.13.0", "2.13.0", "fail-brew", "program 3.0.1, Claude Code plugin 2.13.0, Codex plugin 2.13.0). brew upgrade", "brew", 1},
		{"shadowed executable", "spacedock", "2.13.0", "2.13.0", "unlinked", "does not resolve to the linked", "", 1},
		{"tap already locked", "spacedock", "2.13.0", "2.13.0", "tap-busy", "could not lock the tap", "", 1},
		{"pinned formula", "spacedock", "2.13.0", "2.13.0", "pinned", "pinned, disabled or deprecated", "", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := stageBrewRepair(t, tc.cask, tc.claude, tc.codex, tc.failure)
			out := runBrewRepairFixture(t, root, tc.code)
			if !strings.Contains(out, tc.want) || !strings.Contains(out, "brew upgrade --cask "+tc.cask) {
				t.Fatalf("wanted %q and independent cask advice in output:\n%s", tc.want, out)
			}
			assertBrewRepairWrites(t, root, tc.changes)
			if tc.failure == "fail-codex" || tc.failure == "fail-brew" {
				if err := os.Remove(filepath.Join(root, tc.failure)); err != nil {
					t.Fatal(err)
				}
				mustWrite(t, filepath.Join(root, "changes"), "")
				if out := runBrewRepairFixture(t, root, 0); !strings.Contains(out, "compass 3.0.1 verified") {
					t.Fatal(out)
				}
				want := "codex"
				if tc.failure == "fail-brew" {
					want = "claude\ncodex"
				}
				assertBrewRepairWrites(t, root, want)
			}
		})
	}
}

func TestBrewRepairMissingProgram(t *testing.T) {
	for _, tc := range []struct {
		name, extras, installFrom, want string
		exit                            int
	}{
		{"extras opt-out", "none", "", "compass program skipped by opt-out", 0},
		{"custom install opt-out", "", "preview", "compass program skipped by opt-out", 0},
		{"explicit extras", "compass", "preview", "brew install spacedock-dev/tap/compass", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := stageBrewRepair(t, "spacedock", "2.13.0", "2.13.0", "")
			t.Setenv("SPACEDOCK_EXTRAS", tc.extras)
			t.Setenv("SPACEDOCK_INSTALL_FROM", tc.installFrom)
			if err := os.Remove(filepath.Join(root, "bin", "spacedock-compass")); err != nil {
				t.Fatal(err)
			}
			if out := runBrewRepairFixture(t, root, tc.exit); !strings.Contains(out, tc.want) {
				t.Fatalf("expected %q, got:\n%s", tc.want, out)
			}
			assertBrewRepairWrites(t, root, "")
		})
	}
}

func TestBrewRepairRefreshesCoreWhenListingVersionStaysBehind(t *testing.T) {
	root := stageBrewRepair(t, "spacedock", "3.0.1", "3.0.1", "")
	if err := os.Remove(filepath.Join(root, "bin", "spacedock-compass")); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(root, "market", ".claude-plugin", "marketplace.json"),
		`{"name":"spacedock","plugins":[{"name":"spacedock","version":"1.0.0","source":{"source":"url","url":"https://github.com/example-fixture/core.git","ref":"stable"}}]}`)
	for _, host := range []string{"claude", "codex"} {
		mustWrite(t, filepath.Join(root, host+"-core"), "1.0.0")
	}
	mustWrite(t, filepath.Join(root, "core-release"), "1.1.0")

	runBrewRepairFixture(t, root, 0)
	assertBrewRepairWrites(t, root, "claude core\ncodex core")
	for _, host := range []string{"claude", "codex"} {
		if got, err := os.ReadFile(filepath.Join(root, host+"-core")); err != nil || string(got) != "1.1.0" {
			t.Errorf("%s core = %q, %v", host, got, err)
		}
	}
}

func TestBrewRepairUpdatesSubspaceHostsIndependentlyOfPi(t *testing.T) {
	for _, ref := range []string{"stable", "v1.1.0"} {
		t.Run(ref, func(t *testing.T) {
			root := stageBrewRepair(t, "spacedock", "3.0.1", "3.0.1", "")
			if err := os.Remove(filepath.Join(root, "bin", "spacedock-compass")); err != nil {
				t.Fatal(err)
			}
			writeProgram(t, filepath.Join(root, "bin", "spacedock-review"), "1.1.0")
			tap := filepath.Join(root, "tap")
			for _, args := range [][]string{{"branch", "stable"}, {"tag", "v1.1.0"}} {
				cmd := exec.Command("git", args...)
				cmd.Dir = tap
				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("git %v: %s (%v)", args, output, err)
				}
			}
			t.Setenv("GIT_ALLOW_PROTOCOL", "file")
			t.Setenv("GIT_CONFIG_COUNT", "1")
			t.Setenv("GIT_CONFIG_KEY_0", "url.file://"+tap+".insteadOf")
			t.Setenv("GIT_CONFIG_VALUE_0", "https://github.com/example-fixture/subspace.git")
			transport := http.DefaultTransport
			t.Cleanup(func() { http.DefaultTransport = transport })
			http.DefaultTransport = suiteRoundTrip(func(req *http.Request) (*http.Response, error) {
				if req.URL.Host != "raw.githubusercontent.com" || !strings.HasPrefix(req.URL.Path, "/example-fixture/subspace/") {
					return nil, fmt.Errorf("unexpected fixture request %s", req.URL)
				}
				manifest := `{"name":"subspace","version":"1.1.0"}`
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(manifest))}, nil
			})

			for _, host := range []string{"claude", "codex"} {
				mustWrite(t, filepath.Join(root, host+"-subspace"), "1.0.0")
			}
			t.Setenv("SPACEDOCK_EXTENSION_MARKETPLACE_SOURCE", filepath.Join(root, "market"))
			mustWrite(t, filepath.Join(root, "market", ".claude-plugin", "marketplace.json"),
				fmt.Sprintf(`{"name":"spacedock","plugins":[{"name":"subspace","version":"1.1.0","source":{"source":"url","url":"https://github.com/example-fixture/subspace.git","ref":%q}},{"name":"unrelated","source":"./unrelated"}]}`, ref))
			t.Setenv("PI_CODING_AGENT_DIR", root)
			mustWrite(t, filepath.Join(root, "settings.json"),
				`{"packages":["git:github.com/spacedock-dev/spacedock@stable","git:github.com/example-fixture/subspace@stable","git:github.com/example-fixture/unrelated@keep"]}`)
			pi := filepath.Join(root, "bin", "pi")
			mustWrite(t, pi, `#!/bin/sh
set -eu
cd "$FIXTURE"
case "$*" in
'remove git:github.com/spacedock-dev/spacedock'|'install git:github.com/spacedock-dev/spacedock@stable') ;;
'remove git:github.com/example-fixture/subspace') echo 'pi remove subspace' >> changes ;;
'install git:github.com/example-fixture/subspace@v1.1.0')
	sed 's/subspace@stable/subspace@v1.1.0/' settings.json > settings.new
	mv settings.new settings.json
	echo 'pi install subspace' >> changes ;;
*) printf 'unexpected Pi command: %s\n' "$*" >&2; exit 1 ;;
esac
`)
			if err := os.Chmod(pi, 0755); err != nil {
				t.Fatal(err)
			}

			out := runBrewRepairFixture(t, root, 0)
			if strings.Count(out, "1.0.0 -> 1.1.0") != 2 || strings.Contains(out, "still MIXED on Pi") {
				t.Fatalf("unexpected host or Pi result:\n%s", out)
			}
			want := "claude subspace\ncodex subspace"
			if ref != "stable" {
				want += "\npi remove subspace\npi install subspace"
			}
			assertBrewRepairWrites(t, root, want)
			if sources := piSettingsSources(root); !slices.Contains(sources, "git:github.com/example-fixture/subspace@"+ref) ||
				!slices.Contains(sources, "git:github.com/example-fixture/unrelated@keep") {
				t.Errorf("unexpected Pi registrations: %v", sources)
			}
		})
	}
}

func assertBrewRepairWrites(t *testing.T, root, want string) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(root, "changes"))
	if err != nil || strings.TrimSpace(string(got)) != want {
		t.Fatalf("package writes = %q, %v, want %q", got, err, want)
	}
}

func runBrewRepairFixture(t *testing.T, root string, wantCode int) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmd := newUpdateCommand(context.Background(), os.Environ(), root, &stdout, &stderr)
	var code int
	if err := cmd.RunE(cmd, nil); err != nil {
		var ok bool
		if e, isExit := err.(exitCodeError); isExit {
			code, ok = e.code, true
		}
		if !ok {
			t.Fatal(err)
		}
	}
	out := stdout.String() + stderr.String()
	if code != wantCode {
		t.Fatalf("exit = %d, want %d:\n%s", code, wantCode, out)
	}
	return out
}

func stageBrewRepair(t *testing.T, cask, claude, codex, failure string) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(root, "bin")
	mustMkdir(t, bin)
	launcher := filepath.Join(root, "Caskroom", cask, "fixture", "spacedock")
	writeProgram(t, launcher, "spacedock 1.0.0")
	savedExecutable, savedBranch, savedMarket, savedExtension := executablePath, devBranch, marketplaceSource, extensionMarketplaceSourceOverride
	executablePath = func() (string, error) { return launcher, nil }
	t.Cleanup(func() {
		executablePath, devBranch, marketplaceSource, extensionMarketplaceSourceOverride = savedExecutable, savedBranch, savedMarket, savedExtension
	})
	market := filepath.Join(root, "market")
	t.Setenv("PATH", bin+string(os.PathListSeparator)+"/usr/bin:/bin")
	t.Setenv("SPACEDOCK_DEV_BRANCH", "main")
	t.Setenv("SPACEDOCK_MARKETPLACE_SOURCE", market)
	t.Setenv("SPACEDOCK_EXTENSION_MARKETPLACE_SOURCE", filepath.Join(root, "unused"))
	if cask == "spacedock@next" {
		t.Setenv("SPACEDOCK_DEV_BRANCH", "next")
		t.Setenv("SPACEDOCK_MARKETPLACE_SOURCE", filepath.Join(root, "unused"))
		t.Setenv("SPACEDOCK_EXTENSION_MARKETPLACE_SOURCE", market)
	}
	t.Setenv("SPACEDOCK_EXTRAS", "none")
	t.Setenv("SPACEDOCK_INSTALL_FROM", "")
	t.Setenv("SPACEDOCK_INSTALL_SCRIPT", filepath.Join(root, "must-not-run.sh"))
	mustWrite(t, filepath.Join(root, "must-not-run.sh"), "echo Homebrew invoked install.sh >&2; exit 3\n")
	t.Setenv("FIXTURE", root)
	t.Setenv("HOME", root)
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("PI_CODING_AGENT_DIR", filepath.Join(root, "pi-profile"))
	mustWrite(t, filepath.Join(root, "channel"), cask)
	mustWrite(t, filepath.Join(root, "claude"), claude)
	mustWrite(t, filepath.Join(root, "codex"), codex)
	mustWrite(t, filepath.Join(root, "changes"), "")
	if failure != "" {
		mustWrite(t, filepath.Join(root, failure), "")
	}
	for _, v := range []string{"2.13.0", "3.0.1"} {
		writeProgram(t, filepath.Join(root, "Cellar", "compass", v, "bin", "spacedock-compass"), v)
	}
	active := filepath.Join(root, "Cellar", "compass", "2.13.0", "bin", "spacedock-compass")
	if failure == "unlinked" {
		active = filepath.Join(root, "shadow", "spacedock-compass")
		writeProgram(t, active, "2.13.0")
	}
	if err := os.Symlink(active, filepath.Join(bin, "spacedock-compass")); err != nil {
		t.Fatal(err)
	}
	writeFileWithDirs(t, filepath.Join(market, ".claude-plugin", "marketplace.json"),
		`{"name":"spacedock","plugins":[{"name":"compass","version":"3.0.1","source":"./compass"}]}`)
	for _, host := range []string{"claude", "codex"} {
		writeFileWithDirs(t, filepath.Join(market, "compass", manifestSubpath(host)), `{"name":"compass","version":"3.0.1"}`)
	}
	if failure == "stale-codex" {
		writeFileWithDirs(t, filepath.Join(market, ".agents", "plugins", "marketplace.json"),
			`{"name":"spacedock","plugins":[{"name":"compass","version":"2.13.0","source":"./compass"}]}`)
	}
	archive := "a synthetic Compass archive"
	digest := fmt.Sprintf("%x", sha256.Sum256([]byte(archive)))
	if failure == "bad-archive" {
		archive += " (corrupted)"
	}
	mustWrite(t, filepath.Join(root, "archive"), archive)
	recipe := strings.NewReplacer("{version}", "3.0.1", "{sha256}", digest).Replace(compassRecipeTemplate)
	tap := filepath.Join(root, "tap")
	writeFileWithDirs(t, filepath.Join(tap, "Formula", "compass.rb"), recipe)
	testgit.InitRepo(t, tap, "-q", "-b", "main")
	for _, args := range [][]string{
		{"add", "Formula/compass.rb"}, {"commit", "-qm", "formula"},
		{"remote", "add", "origin", spacedockTapRemote}, {"update-ref", "refs/remotes/origin/main", "HEAD"},
		{"symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main"}, {"branch", "--set-upstream-to=origin/main"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = tap
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("fixture git %v: %s (%v)", args, out, err)
		}
	}
	if failure == "tap-busy" {
		mustWrite(t, filepath.Join(tap, ".git", "index.lock"), "sibling lock")
		t.Cleanup(func() {
			if data, err := os.ReadFile(filepath.Join(tap, ".git", "index.lock")); err != nil || string(data) != "sibling lock" {
				t.Errorf("sibling Git lock changed: %q, %v", data, err)
			}
		})
	}
	info := map[string]any{
		"name": "compass", "tap": spacedockTap, "full_name": spacedockTap + "/compass",
		"linked_keg": "KEG", "installed": []any{map[string]string{"version": "KEG"}},
		"pinned": failure == "pinned", "outdated": true, "versions": map[string]string{"stable": "3.0.1"},
		"ruby_source_path": "Formula/compass.rb", "ruby_source_checksum": map[string]string{"sha256": fmt.Sprintf("%x", sha256.Sum256([]byte(recipe)))},
		"urls": map[string]any{"stable": map[string]string{"url": compassAsset("3.0.1", runtime.GOOS, runtime.GOARCH), "checksum": digest}},
	}
	for _, key := range []string{"dependencies", "build_dependencies", "test_dependencies", "recommended_dependencies", "optional_dependencies", "uses_from_macos", "requirements", "conflicts_with", "link_overwrite"} {
		info[key] = []string{}
	}
	data, err := json.Marshal(map[string]any{"formulae": []any{info}, "casks": []string{}})
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(root, "info"), string(data))
	mustWrite(t, filepath.Join(root, "keg"), "2.13.0")
	for name, script := range map[string]string{"brew": brewRepairFixtureBrew, "claude": brewRepairFixtureHost, "codex": brewRepairFixtureHost} {
		filename := filepath.Join(bin, name)
		mustWrite(t, filename, script)
		if err := os.Chmod(filename, 0755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const brewRepairFixtureBrew = `#!/bin/sh
set -eu
test "$HOMEBREW_NO_AUTO_UPDATE/$HOMEBREW_NO_INSTALL_CLEANUP/$HOMEBREW_NO_INSTALLED_DEPENDENTS_CHECK" = 1/1/1
cd "$FIXTURE"
case "$*" in
'--repo spacedock-dev/tap') printf '%s/tap\n' "$PWD" ;;
'--cellar spacedock-dev/tap/compass') printf '%s/Cellar/compass\n' "$PWD" ;;
'info --json=v2 --formula spacedock-dev/tap/compass') sed "s/KEG/$(cat keg)/g" info ;;
'fetch --formula spacedock-dev/tap/compass') test ! -f missing-archive || { echo 'archive unavailable' >&2; exit 1; } ;;
'--cache --formula spacedock-dev/tap/compass') printf '%s/archive\n' "$PWD" ;;
'upgrade --dry-run --formula spacedock-dev/tap/compass')
	printf '==> Would upgrade 1 requested outdated package\nspacedock-dev/tap/compass 2.13.0 -> 3.0.1\n'
	if test -f extra-upgrade; then printf '==> Would upgrade 1 dependent\nother 1.0.0 -> 1.0.1\n'; fi ;;
'upgrade --formula spacedock-dev/tap/compass')
	test -f tap/.git/index.lock
	echo brew >> changes
	printf 3.0.1 > keg
	ln -sf "$PWD/Cellar/compass/3.0.1/bin/spacedock-compass" bin/spacedock-compass
	test ! -f fail-brew || { echo 'failure after linking' >&2; exit 1; } ;;
*) echo "unexpected Homebrew operation: $*" >&2; exit 1 ;;
esac
`

const brewRepairFixtureHost = `#!/bin/sh
set -eu
host=${0##*/}
cd "$FIXTURE"
v=$(cat "$host")
market=$PWD/market
core=spacedock
if test "$(cat channel)" = spacedock@next; then core=spacedock-edge; fi
case "$*" in
'plugin list --json')
	core_version=1.0.0
	core_enabled=false
	if test -f "$host-core"; then core_version=$(cat "$host-core"); core_enabled=true; fi
	if test "$host" = claude; then
		printf '[{"id":"spacedock@%s","version":"%s","scope":"user","enabled":%s,"installPath":"/core"},{"id":"compass@spacedock","version":"%s","scope":"user","enabled":true,"installPath":"/compass"}' "$core" "$core_version" "$core_enabled" "$v"
		if test -f "$host-subspace"; then printf ',{"id":"subspace@spacedock","version":"%s","scope":"user","enabled":true,"installPath":"/subspace"}' "$(cat "$host-subspace")"; fi
		printf ',{"id":"unrelated@spacedock","version":"1.0.0","scope":"user","enabled":true,"installPath":"/unrelated"}]\n'
	else
		printf '{"installed":[{"pluginId":"spacedock@%s","version":"%s","enabled":%s,"installed":true},{"pluginId":"compass@spacedock","version":"%s","installed":true,"enabled":true}' "$core" "$core_version" "$core_enabled" "$v"
		if test -f "$host-subspace"; then printf ',{"pluginId":"subspace@spacedock","version":"%s","installed":true,"enabled":true}' "$(cat "$host-subspace")"; fi
		printf ',{"pluginId":"unrelated@spacedock","version":"1.0.0","installed":true,"enabled":true}]}\n'
	fi ;;
'plugin marketplace list --json')
	case "$host" in
	claude) printf '[{"name":"spacedock","source":"directory","installLocation":"%s"}]\n' "$market" ;;
	codex)
		printf '{"marketplaces":[{"name":"spacedock","root":"%s","marketplaceSource":{"sourceType":"local","source":"%s"}}]}\n' "$market" "$market" ;;
	esac ;;
'plugin update spacedock@spacedock --scope user'|'plugin add spacedock@spacedock')
	echo "$host core" >> changes
	cp core-release "$host-core" ;;
'plugin update subspace@spacedock --scope user'|'plugin add subspace@spacedock')
	echo "$host subspace" >> changes
	printf 1.1.0 > "$host-subspace" ;;
'plugin update compass@spacedock --scope user'|'plugin add compass@spacedock')
	test ! -f "fail-$host" || { echo 'fixture host failed' >&2; exit 1; }
	echo "$host" >> changes
	printf 3.0.1 > "$host" ;;
*) echo "unexpected host operation: $*" >&2; exit 1 ;;
esac
`
