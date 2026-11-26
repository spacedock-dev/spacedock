package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/testgit"
)

const compassFixtureRemote = "https://github.com/example-fixture/compass.git"

func compassFixtureListing(version string) string {
	return `{"name":"spacedock","plugins":[{"name":"unrelated","source":"./unrelated"},{"name":"compass","version":"` + version +
		`","source":{"source":"git-subdir","url":"` + compassFixtureRemote + `","path":"compass","ref":"stable"}}]}`
}

func brewFixtureGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %s (%v)", args, out, err)
	}
	return strings.TrimSpace(string(out))
}

// clonePublisher copies repo as a published remote. GitHub holds no remote-tracking
// refs, and ls-remote's HEAD pattern would match the clone's origin/HEAD.
func clonePublisher(t *testing.T, repo, dir string) {
	t.Helper()
	brewFixtureGit(t, filepath.Dir(dir), "clone", "-q", repo, dir)
	brewFixtureGit(t, dir, "remote", "remove", "origin")
	testgit.InitRepo(t, dir)
}

func commitFixtureFile(t *testing.T, repo, file, content, message string) string {
	t.Helper()
	writeFileWithDirs(t, filepath.Join(repo, file), content)
	brewFixtureGit(t, repo, "add", file)
	brewFixtureGit(t, repo, "commit", "-qm", message)
	return brewFixtureGit(t, repo, "rev-parse", "HEAD")
}

// rewriteGitHubRemotes sends Git's reads of each GitHub URL prefix to a local repository.
func rewriteGitHubRemotes(t *testing.T, remotes map[string]string) {
	t.Helper()
	t.Setenv("GIT_ALLOW_PROTOCOL", "file")
	t.Setenv("GIT_CONFIG_COUNT", strconv.Itoa(len(remotes)))
	i := 0
	for remote, dir := range remotes {
		t.Setenv("GIT_CONFIG_KEY_"+strconv.Itoa(i), "url.file://"+dir+".insteadOf")
		t.Setenv("GIT_CONFIG_VALUE_"+strconv.Itoa(i), remote)
		i++
	}
}

// serveGitHubRaw answers raw.githubusercontent.com/<owner>/<repo>/<ref>/<path> from the
// local repository registered for owner/repo, and records every requested path.
func serveGitHubRaw(t *testing.T, repos map[string]string) func() []string {
	t.Helper()
	transport := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = transport })
	var mu sync.Mutex
	var seen []string
	respond := func(code int, body string) *http.Response {
		return &http.Response{StatusCode: code, Body: io.NopCloser(strings.NewReader(body))}
	}
	http.DefaultTransport = suiteRoundTrip(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "raw.githubusercontent.com" {
			return nil, fmt.Errorf("unexpected fixture request %s", req.URL)
		}
		mu.Lock()
		seen = append(seen, req.URL.Path)
		mu.Unlock()
		parts := strings.SplitN(strings.TrimPrefix(req.URL.Path, "/"), "/", 4)
		if len(parts) != 4 {
			return respond(http.StatusBadRequest, ""), nil
		}
		dir, ok := repos[parts[0]+"/"+parts[1]]
		if !ok {
			return respond(http.StatusNotFound, ""), nil
		}
		out, err := exec.Command("git", "-C", dir, "show", parts[2]+":"+parts[3]).Output()
		if err != nil {
			return respond(http.StatusNotFound, ""), nil
		}
		return respond(http.StatusOK, string(out)), nil
	})
	return func() []string {
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(seen)
	}
}

// stageCompassSource publishes compass 3.0.1 manifests on a stable branch and returns
// the repository and its commit.
func stageCompassSource(t *testing.T, root, version string) (string, string) {
	t.Helper()
	src := filepath.Join(root, "compass-source")
	mustMkdir(t, src)
	testgit.InitRepo(t, src, "-q", "-b", "stable")
	for _, host := range []string{"claude", "codex"} {
		writeFileWithDirs(t, filepath.Join(src, "compass", manifestSubpath(host)), `{"name":"compass","version":"`+version+`"}`)
	}
	brewFixtureGit(t, src, "add", "compass")
	brewFixtureGit(t, src, "commit", "-qm", "compass "+version)
	return src, brewFixtureGit(t, src, "rev-parse", "HEAD")
}

func TestBrewRepairRefreshesLaggingHostSnapshotBeforeCompass(t *testing.T) {
	for _, tc := range []struct {
		name, want, changes string
		code                int
	}{
		{"published release", "compass 3.0.1 verified", "brew\nclaude\ncodex", 0},
		{"dirty Claude snapshot", "could not refresh claude's spacedock snapshot: marketplace snapshot has local changes", "", 1},
		{"listing not yet in Git", "claude's spacedock snapshot cannot select compass 3.0.1. Refresh it with: claude plugin marketplace update spacedock", "", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := stageBrewRepair(t, "spacedock", "2.13.0", "2.13.0", "")
			src, commit := stageCompassSource(t, root, "3.0.1")
			publisher := filepath.Join(root, "publisher")
			mustMkdir(t, publisher)
			testgit.InitRepo(t, publisher, "-q", "-b", "main")
			listing := filepath.Join(".claude-plugin", "marketplace.json")
			commitFixtureFile(t, publisher, listing, compassFixtureListing("2.13.0"), "Initial listing")
			rewriteGitHubRemotes(t, map[string]string{
				"https://github.com/example-fixture/market.git": publisher,
				compassFixtureRemote:                            src,
			})
			t.Setenv("SPACEDOCK_MARKETPLACE_SOURCE", "example-fixture/market")
			claude, codex := filepath.Join(root, "claude-profile"), filepath.Join(root, "codex-profile")
			t.Setenv("CLAUDE_CONFIG_DIR", claude)
			t.Setenv("CODEX_HOME", codex)
			snapshots := map[string]string{
				"claude": filepath.Join(claude, "plugins", "marketplaces", "spacedock"),
				"codex":  filepath.Join(codex, ".tmp", "marketplaces", "spacedock"),
			}
			for host, dir := range snapshots {
				brewFixtureGit(t, root, "clone", "-q", "--depth=1", "https://github.com/example-fixture/market.git", dir)
				mustWrite(t, filepath.Join(root, host+"-registered"), dir)
			}
			// GitHub's raw listing may run ahead of the Git ref the hosts fetch.
			served := publisher
			if tc.name == "listing not yet in Git" {
				served = filepath.Join(root, "served")
				clonePublisher(t, publisher, served)
				commitFixtureFile(t, served, listing, compassFixtureListing("3.0.1"), "Publish compass")
			} else {
				commitFixtureFile(t, publisher, listing, compassFixtureListing("3.0.1"), "Publish compass")
			}
			if tc.name == "dirty Claude snapshot" {
				mustWrite(t, filepath.Join(snapshots["claude"], "local-notes"), "preserve")
			}
			requests := serveGitHubRaw(t, map[string]string{"example-fixture/market": served, "example-fixture/compass": src})
			before := map[string]string{}
			for host, dir := range snapshots {
				before[host] = brewFixtureGit(t, dir, "rev-parse", "HEAD")
			}

			out := runBrewRepairFixture(t, root, tc.code)
			if !strings.Contains(out, tc.want) {
				t.Fatalf("wanted %q in output:\n%s", tc.want, out)
			}
			assertBrewRepairWrites(t, root, tc.changes)
			published := brewFixtureGit(t, publisher, "rev-parse", "HEAD")
			for host, dir := range snapshots {
				head := brewFixtureGit(t, dir, "rev-parse", "HEAD")
				if want := before[host]; tc.code == 0 {
					want = published
					if head != want {
						t.Errorf("%s snapshot = %s, want the published listing %s", host, head, want)
					}
				} else if head != want {
					t.Errorf("held repair moved the %s snapshot from %s to %s", host, want, head)
				}
			}
			if tc.code != 0 {
				return
			}
			if !slices.Contains(requests(), "/example-fixture/compass/"+commit+"/compass/.claude-plugin/plugin.json") {
				t.Errorf("manifest was not read at the resolved stable commit: %q", requests())
			}
			runBrewRepairFixture(t, root, 0)
			assertBrewRepairWrites(t, root, tc.changes)
		})
	}
}

func TestBrewRepairAdvancesLaggingTapBeforeUpgrade(t *testing.T) {
	root := stageBrewRepair(t, "spacedock", "2.13.0", "2.13.0", "")
	tap := filepath.Join(root, "tap")
	formula := filepath.Join("Formula", "compass.rb")
	commitFixtureFile(t, tap, formula, "# stable 2.13.0\n", "compass 2.13.0")
	brewFixtureGit(t, tap, "update-ref", "refs/remotes/origin/main", "HEAD")
	publisher := filepath.Join(root, "tap-publisher")
	clonePublisher(t, tap, publisher)
	published := commitFixtureFile(t, publisher, formula, "# stable 3.0.1\n", "compass 3.0.1")
	rewriteGitHubRemotes(t, map[string]string{spacedockTapRemote: publisher})

	if out := runBrewRepairFixture(t, root, 0); !strings.Contains(out, "compass 3.0.1 verified") {
		t.Fatal(out)
	}
	assertBrewRepairWrites(t, root, "brew\nclaude\ncodex")
	if head := brewFixtureGit(t, tap, "rev-parse", "HEAD"); head != published {
		t.Errorf("tap HEAD = %s, want the published %s", head, published)
	}
	if remote := brewFixtureGit(t, tap, "config", "--get", "remote.origin.url"); remote != spacedockTapRemote {
		t.Errorf("tap origin changed to %q", remote)
	}
}

func TestBrewTapAdvanceRefusesUnsafeTaps(t *testing.T) {
	for _, tc := range []struct{ name, want string }{
		{"advances", ""},
		{"diverged", "has diverged"},
		{"different default branch", "could not verify the official tap default branch"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := stageBrewRepair(t, "spacedock", "2.13.0", "2.13.0", "")
			tap := filepath.Join(root, "tap")
			publisher := filepath.Join(root, "tap-publisher")
			clonePublisher(t, tap, publisher)
			published := commitFixtureFile(t, publisher, "Formula/compass.rb", "# stable 3.0.1\n", "publish")
			if tc.name == "different default branch" {
				brewFixtureGit(t, publisher, "checkout", "-qb", "release")
			}
			head := brewFixtureGit(t, tap, "rev-parse", "HEAD")
			if tc.name == "diverged" {
				head = commitFixtureFile(t, tap, "Formula/local.rb", "# local\n", "local work")
			}
			rewriteGitHubRemotes(t, map[string]string{spacedockTapRemote: publisher})

			repair := brewRepair{context.Background(), os.Environ(), root}
			got, err := repair.advanceBrewTap(tap, "main", head)
			if tc.want == "" {
				if err != nil || got != published || brewFixtureGit(t, tap, "rev-parse", "HEAD") != published {
					t.Fatalf("advance = %s, %v; want %s", got, err, published)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("advance error = %v, want %q", err, tc.want)
			}
			if now := brewFixtureGit(t, tap, "rev-parse", "HEAD"); now != head {
				t.Errorf("refused advance moved the tap from %s to %s", head, now)
			}
		})
	}
}

func TestBrewRepairResolvesRemoteRefs(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	src, branch := stageCompassSource(t, root, "3.0.1")
	tagged := commitFixtureFile(t, src, "compass/notes", "release\n", "release notes")
	brewFixtureGit(t, src, "tag", "-a", "-m", "release", "v3.0.1")
	brewFixtureGit(t, src, "tag", "light", branch)
	brewFixtureGit(t, src, "branch", "both", branch)
	brewFixtureGit(t, src, "tag", "both", tagged)
	brewFixtureGit(t, src, "reset", "-q", "--hard", branch)
	rewriteGitHubRemotes(t, map[string]string{compassFixtureRemote: src})
	pinned := strings.Repeat("c", 40)

	repair := brewRepair{context.Background(), os.Environ(), root}
	for _, tc := range []struct{ ref, want, err string }{
		{"stable", branch, ""},
		{"v3.0.1", tagged, ""},
		{"light", branch, ""},
		{pinned, pinned, ""},
		{"both", "", "is ambiguous"},
		{"missing", "", "could not resolve"},
	} {
		t.Run(tc.ref, func(t *testing.T) {
			got, err := repair.resolveRepairRef("example-fixture/compass", tc.ref)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("resolve %s = %s, %v; want error %q", tc.ref, got, err, tc.err)
				}
				return
			}
			if err != nil || got != tc.want {
				t.Fatalf("resolve %s = %s, %v; want %s", tc.ref, got, err, tc.want)
			}
		})
	}
}

func TestBrewRepairVerifiesRemoteManifestAtResolvedCommit(t *testing.T) {
	for _, tc := range []struct {
		name, host, version, err string
		codexManifest            bool
	}{
		{"Claude manifest", "claude", "3.0.1", "", true},
		{"Codex manifest", "codex", "3.0.1", "", true},
		{"Codex falls back to the Claude manifest", "codex", "3.0.1", "", false},
		{"version mismatch", "claude", "3.0.0", "does not identify compass 3.0.1", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
			t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
			src, commit := stageCompassSource(t, root, tc.version)
			if !tc.codexManifest {
				brewFixtureGit(t, src, "rm", "-q", filepath.Join("compass", manifestSubpath("codex")))
				brewFixtureGit(t, src, "commit", "-qm", "Claude manifest only")
				commit = brewFixtureGit(t, src, "rev-parse", "HEAD")
			}
			rewriteGitHubRemotes(t, map[string]string{compassFixtureRemote: src})
			requests := serveGitHubRaw(t, map[string]string{"example-fixture/compass": src})
			entry := repairEntry{Version: "3.0.1", Repo: "example-fixture/compass", Path: "compass", Ref: "stable"}

			repair := brewRepair{context.Background(), os.Environ(), root}
			fingerprint, err := repair.verifyRepairManifest(entry, "", tc.host, "compass")
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("verify = %q, %v; want error %q", fingerprint, err, tc.err)
				}
				return
			}
			if err != nil || !strings.HasPrefix(fingerprint, commit+":") {
				t.Fatalf("verify = %q, %v; want a fingerprint pinned to %s", fingerprint, err, commit)
			}
			for _, path := range requests() {
				if !strings.HasPrefix(path, "/example-fixture/compass/"+commit+"/compass/") {
					t.Errorf("manifest read %s is not pinned to %s", path, commit)
				}
			}
		})
	}
}

func TestBrewRepairReadsGitHubDocuments(t *testing.T) {
	transport := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = transport })
	var paths []string
	http.DefaultTransport = suiteRoundTrip(func(req *http.Request) (*http.Response, error) {
		paths = append(paths, req.URL.Host+req.URL.Path)
		code, body := http.StatusOK, `{"name":"spacedock"}`
		switch {
		case strings.Contains(req.URL.Path, "/broken/"):
			code = http.StatusBadGateway
		case strings.Contains(req.URL.Path, "/missing/"):
			code = http.StatusNotFound
		case strings.Contains(req.URL.Path, "/huge/"):
			body = strings.Repeat("x", 1<<20+1)
		}
		return &http.Response{StatusCode: code, Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	repair := brewRepair{ctx: context.Background()}

	data, err := repair.readReleaseListing("example-fixture/market")
	if err != nil || string(data) != `{"name":"spacedock"}` {
		t.Fatalf("listing = %q, %v", data, err)
	}
	if want := "raw.githubusercontent.com/example-fixture/market/HEAD/.claude-plugin/marketplace.json"; !slices.Equal(paths, []string{want}) {
		t.Errorf("listing requests = %q, want %q", paths, want)
	}
	if _, err := repair.readReleaseListing("not a repository"); err == nil || !strings.Contains(err.Error(), "directory or a GitHub owner/repo") {
		t.Errorf("invalid source error = %v", err)
	}
	for _, tc := range []struct{ repo, want string }{
		{"example-fixture/broken", "GitHub returned HTTP 502"},
		{"example-fixture/huge", "exceeds the supported size"},
	} {
		if _, err := repair.rawGitHub(tc.repo, "HEAD", "file.json"); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s error = %v, want %q", tc.repo, err, tc.want)
		}
	}
	if _, err := repair.rawGitHub("example-fixture/missing", "HEAD", "file.json"); !os.IsNotExist(err) {
		t.Errorf("missing document error = %v, want os.ErrNotExist", err)
	}
}

func TestSnapshotRefreshAdviceNamesHostCommand(t *testing.T) {
	git := repairMarketplace{Source: "github"}
	codexGit := repairMarketplace{}
	codexGit.MarketplaceSource.SourceType = "git"
	for _, tc := range []struct {
		host   string
		market repairMarketplace
		want   string
	}{
		{"claude", git, "Refresh it with: claude plugin marketplace update spacedock"},
		{"codex", codexGit, "Refresh it with: codex plugin marketplace upgrade spacedock"},
		{"claude", repairMarketplace{Source: "directory"}, "Make the matching listing available in this host"},
	} {
		if got := snapshotRefreshAdvice(tc.host, tc.market); got != tc.want {
			t.Errorf("%s advice = %q, want %q", tc.host, got, tc.want)
		}
	}
}
