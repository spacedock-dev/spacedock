package release

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// Run real curl against failed routes. Removing a timeout or accepting a bad
// response must prevent the direct GitHub fallback from resolving v9.9.9.
func TestInstallLookupRouteFailureUsesDirectGitHub(t *testing.T) {
	realCurl, err := exec.LookPath("curl")
	if err != nil {
		t.Fatal(err)
	}
	for _, channel := range []string{"stable", "edge"} {
		t.Run(channel, func(t *testing.T) {
			for _, route := range []string{"404", "503", "down", "slow", "loop", "html", "foreign", "malformed", "traversal", "nested-tag-path", "whitespace", "metacharacters", "multiline", "prerelease"} {
				if channel == "edge" && (route == "nested-tag-path" || route == "prerelease") {
					continue
				}
				t.Run(route, func(t *testing.T) {
					r := newExtrasRun(t, "1.0.1")
					var mu sync.Mutex
					var requests []string
					server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						mu.Lock()
						requests = append(requests, req.Host+req.URL.RequestURI())
						mu.Unlock()
						if req.Host == "spacedock.md" {
							switch route {
							case "404", "503":
								status, _ := strconv.Atoi(route)
								w.WriteHeader(status)
							case "down":
								conn, _, _ := w.(http.Hijacker).Hijack()
								conn.Close()
							case "slow":
								select {
								case <-time.After(8 * time.Second):
								case <-req.Context().Done():
								}
							case "loop":
								http.Redirect(w, req, "https://spacedock.md"+req.URL.RequestURI(), http.StatusFound)
							case "html":
								fmt.Fprintln(w, `<html>"tag_name": "v8.8.8"</html>`)
							case "foreign":
								http.Redirect(w, req, "https://other.example/releases/tag/v8.8.8", http.StatusFound)
							default:
								tag := map[string]string{"malformed": "garbage", "traversal": "v1.2.3/../../evil", "nested-tag-path": "v1.2.3/releases/tag/v8.8.8", "whitespace": "v1.2.3 bad", "metacharacters": "v1.2.3-$(evil)", "multiline": "v1.2.3\nv8.8.8", "prerelease": "v8.8.8-pre1"}[route]
								if channel == "edge" {
									fmt.Fprintf(w, "[{\"tag_name\": %q}]\n", tag)
								} else {
									http.Redirect(w, req, "https://github.com/spacedock-dev/spacedock/releases/tag/"+tag, http.StatusFound)
								}
							}
							return
						}
						if req.Host == "other.example" {
							fmt.Fprintln(w, `[{"tag_name": "v8.8.8"}]`)
						} else if req.URL.Path == "/spacedock-dev/spacedock/releases/latest" {
							http.Redirect(w, req, "https://github.com/spacedock-dev/spacedock/releases/tag/v9.9.9", http.StatusFound)
						} else {
							fmt.Fprintln(w, `[{"tag_name": "v9.9.9"}]`)
						}
					}))
					defer server.Close()
					defer server.CloseClientConnections()
					address := strings.TrimPrefix(server.URL, "https://")
					body := "exec '" + realCurl + "' --noproxy '*' -k"
					for _, host := range []string{"spacedock.md", "github.com", "api.github.com", "other.example"} {
						body += " --connect-to '" + host + ":443:" + address + "'"
					}
					r.addProgram("curl", body+" \"$@\"")
					r.addProgram("git", "exit 1")
					start := time.Now()
					stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_CHANNEL="+channel)
					if code != 0 || !strings.Contains(stderr, "resolved "+channel+" channel to v9.9.9\n") {
						t.Fatalf("exit=%d; direct GitHub fallback failed:\n%s", code, stderr)
					}
					if route == "slow" && time.Since(start) > 7*time.Second {
						t.Error("route delayed the fallback for more than seven seconds")
					}
					mu.Lock()
					defer mu.Unlock()
					fallback := "github.com/spacedock-dev/spacedock/releases/latest"
					if channel == "edge" {
						fallback = "api.github.com/repos/spacedock-dev/spacedock/releases?per_page=1"
					}
					if len(requests) < 2 || !strings.HasPrefix(requests[0], "spacedock.md/latest-release/spacedock") || requests[len(requests)-1] != fallback {
						t.Errorf("requests=%q, want route then %s", requests, fallback)
					}
				})
			}
		})
	}
}

// Keep curl's request handling real; only route DNS to this test's TLS listener.
func localLookupCurl(r *extrasRun, handler http.HandlerFunc) {
	r.t.Helper()
	realCurl, err := exec.LookPath("curl")
	if err != nil {
		r.t.Fatal(err)
	}
	server := httptest.NewTLSServer(handler)
	r.t.Cleanup(server.Close)
	address := strings.TrimPrefix(server.URL, "https://")
	body := "exec '" + realCurl + "' -q --noproxy '*' -k"
	for _, host := range []string{"spacedock.md", "github.com", "api.github.com"} {
		body += " --connect-to '" + host + ":443:" + address + "'"
	}
	r.addProgram("curl", body+" \"$@\"")
	r.addProgram("git", "exit 1")
}

func realSpacedock(t *testing.T, version string) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "spacedock")
	cmd := exec.Command("go", "build", "-buildvcs=false", "-ldflags", "-X github.com/spacedock-dev/spacedock/internal/cli.Version="+version, "-o", binary, "./cmd/spacedock")
	cmd.Dir = filepath.Join("..", "..")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build spacedock: %v\n%s", err, out)
	}
	return binary
}

// The real installer must resolve each shipped product once, before replacement.
func TestInstallLookupEachProductOnce(t *testing.T) {
	binary, err := os.ReadFile(realSpacedock(t, "9.9.9"))
	if err != nil {
		t.Fatal(err)
	}
	for _, channel := range []string{"stable", "edge"} {
		t.Run(channel, func(t *testing.T) {
			r := newExtrasRun(t, "1.0.1")
			r.script = installScript()
			osName, arch := goosArch(t)
			for _, suffix := range []string{"", "_edge"} {
				writeTarGzFiles(t, filepath.Join(r.dist, "spacedock_9.9.9_"+osName+"_"+arch+suffix+".tar.gz"), map[string]string{"spacedock": string(binary)})
			}
			r.writeChecksums(r.dist)
			writeTarGzFiles(t, filepath.Join(r.extras, "subspace_0.16.0_"+osName+"_"+arch+".tar.gz"), map[string]string{"subspace-tui": "#!/bin/sh\necho 0.16.0\n", "subspace-relay-client": "#!/bin/sh\necho 0.16.0\n"})
			writeTarGzFiles(t, filepath.Join(r.extras, "compass_1.0.1_"+osName+"_"+arch+".tar.gz"), map[string]string{"compass/compass": "#!/bin/sh\necho compass\n"})
			r.writeChecksums(r.extras)
			var mu sync.Mutex
			var lookups []string
			localLookupCurl(r, func(w http.ResponseWriter, req *http.Request) {
				if req.Host == "spacedock.md" {
					mu.Lock()
					lookups = append(lookups, req.URL.RequestURI())
					mu.Unlock()
					if req.Method != "GET" || !strings.Contains(req.UserAgent(), "/unknown (") {
						t.Errorf("lookup = %s %s, want GET with the version from before replacement", req.Method, req.UserAgent())
					}
					product := strings.TrimPrefix(req.URL.Path, "/latest-release/")
					if req.URL.RawQuery == "prerelease=1" {
						http.Redirect(w, req, "https://api.github.com/repos/spacedock-dev/"+product+"/releases?per_page=1", http.StatusFound)
					} else {
						http.Redirect(w, req, "https://github.com/spacedock-dev/"+product+"/releases/latest", http.StatusFound)
					}
					return
				}
				product := strings.Split(strings.TrimPrefix(strings.TrimPrefix(req.URL.Path, "/"), "repos/"), "/")[1]
				tag := map[string]string{"spacedock": "v9.9.9", "subspace": "v0.16.0", "compass": "v1.0.1"}[product]
				switch {
				case strings.HasSuffix(req.URL.Path, "/latest"):
					http.Redirect(w, req, "https://github.com/spacedock-dev/"+product+"/releases/tag/"+tag, http.StatusFound)
				case req.Host == "api.github.com":
					fmt.Fprintf(w, `[{"tag_name":%q}]`, tag)
				case strings.Contains(req.URL.Path, "/download/"):
					dir := r.extras
					if product == "spacedock" {
						dir = r.dist
					}
					http.ServeFile(w, req, filepath.Join(dir, filepath.Base(req.URL.Path)))
				default:
					// A headers-only GET must not wait for a release page body.
					w.Header().Set("Content-Length", "1000000")
					w.WriteHeader(http.StatusOK)
					w.(http.Flusher).Flush()
					<-req.Context().Done()
				}
			})
			stderr, code := r.runWith("SPACEDOCK_PLUGINS=none", "SPACEDOCK_CHANNEL="+channel)
			if code != 0 {
				t.Fatalf("exit=%d:\n%s", code, stderr)
			}
			suffix := ""
			if channel == "edge" {
				suffix = "?prerelease=1"
			}
			want := "/latest-release/spacedock" + suffix + "\n/latest-release/subspace" + suffix + "\n/latest-release/compass" + suffix
			mu.Lock()
			got := strings.Join(lookups, "\n")
			mu.Unlock()
			if got != want {
				t.Errorf("release lookups=%q, want %q", got, want)
			}
			for _, program := range []string{"spacedock", "subspace-tui", "subspace-relay-client", "compass", "spacedock-review", "spacedock-compass"} {
				if _, err := os.Stat(filepath.Join(r.bin, program)); err != nil {
					t.Errorf("%s was not installed: %v", program, err)
				}
			}
			cmd := exec.Command(filepath.Join(r.bin, "spacedock"), "--version")
			cmd.Env = []string{"CI=1", "HOME=" + r.home, "PATH=/usr/bin:/bin"}
			if out, err := cmd.CombinedOutput(); err != nil || !strings.HasPrefix(string(out), "spacedock 9.9.9\n") {
				t.Errorf("installed binary: %v: %s", err, out)
			}
		})
	}
}

// Changing version discovery, normalization, or marker rules breaks a literal UA.
func TestInstallLookupUserAgent(t *testing.T) {
	binary := realSpacedock(t, "0.27.3")
	for _, tc := range []struct {
		name, os, arch, version, ci, internal, want string
		inPath, readsInput, realBinary              bool
		fails, slow                                 bool
	}{
		{name: "real installed binary", realBinary: true, want: "spacedock-install/0.27.3 (" + runtime.GOOS + "; " + runtime.GOARCH + ")"},
		{name: "installed version", os: "Darwin", arch: "aarch64", version: "spacedock v0.27.3\nOS: darwin/arm64", want: "spacedock-install/0.27.3 (darwin; arm64)"},
		{name: "version cannot read installer input", os: "Linux", arch: "amd64", version: "spacedock 0.27.3", readsInput: true, want: "spacedock-install/0.27.3 (linux; amd64)"},
		{name: "version from PATH", os: "Linux", arch: "x86_64", version: "spacedock 0.28.0-pre10", inPath: true, want: "spacedock-install/0.28.0-pre10 (linux; amd64)"},
		{name: "case folded platform", os: "dArWiN", arch: "X64", want: "spacedock-install/unknown (darwin; amd64)"},
		{name: "CI one", os: "Linux", arch: "amd64", ci: "1", want: "spacedock-install/unknown (linux; amd64; ci)"},
		{name: "CI true trimmed", os: "Linux", arch: "amd64", ci: " \tTrUe\n", want: "spacedock-install/unknown (linux; amd64; ci)"},
		{name: "CI yes", os: "Linux", arch: "amd64", ci: "YES", want: "spacedock-install/unknown (linux; amd64; ci)"},
		{name: "CI on and internal", os: "Linux", arch: "amd64", ci: "on", internal: "1", want: "spacedock-install/unknown (linux; amd64; ci; internal)"},
		{name: "internal only", os: "Linux", arch: "amd64", ci: "false", internal: "1", want: "spacedock-install/unknown (linux; amd64; internal)"},
		{name: "markers require explicit values", os: "Linux", arch: "amd64", ci: "enabled", internal: "true", want: "spacedock-install/unknown (linux; amd64)"},
		{name: "unrecognized version", os: "Linux", arch: "arm64", version: "not a version", want: "spacedock-install/unknown (linux; arm64)"},
		{name: "development version is unknown", os: "Linux", arch: "arm64", version: "spacedock dev", ci: "true", want: "spacedock-install/unknown (linux; arm64; ci)"},
		{name: "failed version command", os: "Linux", arch: "amd64", version: "spacedock 1.2.3", fails: true, want: "spacedock-install/unknown (linux; amd64)"},
		{name: "slow version command", os: "Linux", arch: "amd64", slow: true, want: "spacedock-install/unknown (linux; amd64)"},
		{name: "version cannot inject a curl option", os: "Linux", arch: "amd64", version: "spacedock 1.2.3' -A evil; $(exit 1)", want: "spacedock-install/unknown (linux; amd64)"},
		{name: "maximum version length", os: "Linux", arch: "amd64", version: "spacedock 1.2.3-" + strings.Repeat("a", 26), want: "spacedock-install/1.2.3-" + strings.Repeat("a", 26) + " (linux; amd64)"},
		{name: "version exceeds contract limit", os: "Linux", arch: "amd64", version: "spacedock 1.2.3-" + strings.Repeat("a", 27), ci: "1", want: "spacedock-install/unknown (linux; amd64; ci)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newExtrasRun(t, "1.0.1")
			r.script = installScript()
			var mu sync.Mutex
			var agents []string
			localLookupCurl(r, func(w http.ResponseWriter, req *http.Request) {
				if req.Host == "spacedock.md" {
					mu.Lock()
					agents = append(agents, req.UserAgent())
					mu.Unlock()
					http.Redirect(w, req, "https://github.com/spacedock-dev/spacedock/releases/tag/v9.9.9", http.StatusFound)
				}
			})
			if tc.realBinary {
				data, err := os.ReadFile(binary)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(r.bin, "spacedock"), data, 0o755); err != nil {
					t.Fatal(err)
				}
			} else {
				r.addProgram("uname", "case \"$1\" in -s) echo '"+tc.os+"' ;; -m) echo '"+tc.arch+"' ;; esac")
			}
			if tc.version != "" || tc.slow {
				r.stdin = "installer input\n"
				body := "printf '%s\\n' '" + strings.ReplaceAll(tc.version, "'", "'\"'\"'") + "'"
				if tc.fails {
					body += "\nexit 1"
				}
				if tc.slow {
					body = "exec sleep 8"
				}
				if tc.readsInput {
					body = "if read -r input; then exit 1; fi\n" + body
				}
				if tc.inPath {
					r.addProgram("spacedock", body)
				} else if err := os.WriteFile(filepath.Join(r.bin, "spacedock"), []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			start := time.Now()
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "CI="+tc.ci, "SPACEDOCK_INTERNAL="+tc.internal)
			if code != 0 {
				t.Fatalf("exit=%d:\n%s", code, stderr)
			}
			if tc.slow && time.Since(start) > 4*time.Second {
				t.Error("version probe delayed the lookup for more than four seconds")
			}
			mu.Lock()
			defer mu.Unlock()
			if len(agents) != 1 || agents[0] != tc.want {
				t.Errorf("wire User-Agents=%q, want [%q]", agents, tc.want)
			}
		})
	}
}

// A pinned inspection has no release lookup and must not execute an old binary.
func TestInstallLookupSkipsVersionProbeForPinnedSource(t *testing.T) {
	r := newExtrasRun(t, "1.0.1")
	r.addProgram("spacedock", "echo probed > '"+r.log+"'")
	stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_INSTALL_FROM="+r.dist)
	if code != 0 {
		t.Fatalf("exit=%d: %s", code, stderr)
	}
	if _, err := os.Stat(r.log); !os.IsNotExist(err) {
		t.Errorf("pinned inspection executed spacedock: %v", err)
	}
}

func lookupRun(t *testing.T) (*extrasRun, string) {
	t.Helper()
	r := newExtrasRun(t, "1.0.1")
	osName, arch := goosArch(t)
	return r, r.stubCurl(osName, arch)
}

func TestInstallResolvesTheNewestTagByVersionWhenGitHubDoesNotAnswer(t *testing.T) {
	cases := []struct {
		name, channel, tags, want, api, redirect string
	}{
		{"a release outranks its own prereleases", "edge", "v0.16.0 v0.16.0-beta.5 v0.16.0-beta.1", "v0.16.0", "fail", "norelease"},
		{"prerelease numbers compare as numbers", "edge", "v0.28.0-pre9 v0.28.0-pre10 v0.27.3", "v0.28.0-pre10", "fail", "norelease"},
		{"minor versions compare as numbers and non-version tags are ignored", "edge", "v1.9.0 v1.10.0 v1.10.0-rc.1 archive/v0", "v1.10.0", "fail", "norelease"},
		{"null API tag falls back to git", "edge", "v0.16.0", "v0.16.0", "null", "norelease"},
		{"unsafe API tag falls back to git", "edge", "v0.16.0", "v0.16.0", "traversal", "norelease"},
		{"stable skips prereleases", "stable", "v0.28.0-pre5 v0.27.3 v0.9.9 v0.27.10", "v0.27.10", "fail", "norelease"},
		{"direct redirect rejects a foreign host", "stable", "v0.27.3", "v0.27.3", "fail", "foreign"},
		{"direct redirect validates the complete tag path", "stable", "v0.27.3", "v0.27.3", "fail", "nested"},
		{"direct redirect rejects an unsafe tag", "stable", "v0.27.3", "v0.27.3", "fail", "traversal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := lookupRun(t)
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_CHANNEL="+tc.channel,
				"STUB_API="+tc.api, "STUB_REDIRECT="+tc.redirect, "STUB_TAGS="+tc.tags)
			if code != 0 || !strings.Contains(stderr, "resolved "+tc.channel+" channel to "+tc.want+"\n") {
				t.Errorf("exit=%d, want %s:\n%s", code, tc.want, stderr)
			}
			if !strings.Contains(stderr, "reading the tags with git ls-remote") {
				t.Errorf("the fallback is not announced:\n%s", stderr)
			}
		})
	}
}

func TestInstallFailsNamingTheRateLimitWhenEveryRouteFails(t *testing.T) {
	for _, tc := range []struct{ name, channel, tags string }{
		{"stable without a tag", "stable", ""},
		{"edge without a tag", "edge", ""},
		{"stable with only a prerelease tag", "stable", "v0.28.0-pre5"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := lookupRun(t)
			stderr, code := r.runWith("SPACEDOCK_PRINT_TARGET=1", "SPACEDOCK_CHANNEL="+tc.channel,
				"STUB_API=fail", "STUB_REDIRECT=fail", "STUB_TAGS="+tc.tags)
			if code != 1 || !strings.Contains(stderr, "rate limit") || !strings.Contains(stderr, "spacedock-dev/spacedock") {
				t.Errorf("exit=%d, want 1 and a message naming the rate limit and the repo:\n%s", code, stderr)
			}
			if strings.Contains(stderr, "resolved ") {
				t.Errorf("a run that resolved nothing claims a version:\n%s", stderr)
			}
		})
	}
}

func TestInstallFinishesTheOtherStepsThenFailsWhenAnExtensionLookupFails(t *testing.T) {
	r, _ := lookupRun(t)
	r.addProgram("claude", "exit 0")

	stderr, code := r.runWith("SPACEDOCK_INSTALL_FROM="+r.dist, "SPACEDOCK_EXTRAS=1", "STUB_FAIL_REPO=spacedock-dev/alpha")
	if code != 1 {
		t.Errorf("exit = %d, want 1 when an extension release cannot be looked up", code)
	}
	for _, want := range []string{"rate limit", "spacedock-dev/alpha", "installed beta v1.0.1", "installed spacedock"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("stderr lacks %q:\n%s", want, stderr)
		}
	}
	if !strings.Contains(r.calls(), "install --host claude") {
		t.Errorf("the plugin step did not run after the failed lookup: %q", r.calls())
	}
	if _, err := os.Stat(filepath.Join(r.bin, "alpha-tui")); err == nil {
		t.Errorf("alpha was installed although its release could not be looked up")
	}
}

func TestInstallExitsZeroWhenAnExtensionDownloadFailsAfterItsTagResolved(t *testing.T) {
	r, _ := lookupRun(t)
	r.addProgram("claude", "exit 0")
	osName, arch := goosArch(t)
	os.Remove(filepath.Join(r.extras, "alpha_0.16.0_"+osName+"_"+arch+".tar.gz"))

	stderr, code := r.runWith("SPACEDOCK_INSTALL_FROM="+r.dist, "SPACEDOCK_EXTRAS=1")
	if code != 0 || !strings.Contains(stderr, "alpha: download failed") {
		t.Errorf("exit=%d, want 0 with the download failure reported:\n%s", code, stderr)
	}
}
