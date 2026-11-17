package updatecheck

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// lookup answers the first request with a 302 to location (when set) and every
// other request with status and body (404 unless it went to location), and
// returns Latest's answer and the number of requests made.
func lookup(channel, location string, status int, body string) (string, int) {
	calls := 0
	client := &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
		calls++
		res := &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
		if location != "" && (calls == 1 || location == loop) {
			res.StatusCode = http.StatusFound
			res.Header.Set("Location", strings.TrimPrefix(location, "loop:"))
		} else if location != "" && r.URL.String() != location {
			res.StatusCode = http.StatusNotFound
		}
		return res, nil
	})}
	return Latest(context.Background(), client, channel, "test"), calls
}

const (
	stableAPI = "https://api.github.com/repos/spacedock-dev/spacedock/releases/latest"
	edgeAPI   = "https://api.github.com/repos/spacedock-dev/spacedock/releases?per_page=1"
	loop      = "loop:" + stableAPI
)

func TestUserAgentInstallerParity(t *testing.T) {
	data, err := os.ReadFile("../../install.sh")
	if err != nil {
		t.Fatal(err)
	}
	source := string(data)
	function := func(name string) string {
		start := strings.Index(source, name+"() {")
		end := strings.Index(source[max(start, 0):], "\n}\n")
		if start < 0 || end < 0 {
			t.Fatalf("missing installer function %s", name)
		}
		return source[start : start+end+3]
	}
	for _, tc := range []struct{ os, arch, ci, internal string }{
		{"Darwin", "aarch64", " \tTrUe\n", "1"}, {"Linux", "x86_64", "yes", "true"},
		{"Linux", "x64", " ON ", "0"}, {"Linux", "amd64", "false", ""},
		{"Darwin", "arm64", "enabled; evil", "1 "}, {"Linux", "arm64", " true ", ""},
	} {
		t.Run(tc.os+tc.arch+tc.ci, func(t *testing.T) {
			bin := t.TempDir()
			tr, err := exec.LookPath("tr")
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(tr, filepath.Join(bin, "tr")); err != nil {
				t.Fatal(err)
			}
			uname := "#!/bin/sh\ncase $1 in -s) printf '%s' '" + tc.os + "';; -m) printf '%s' '" + tc.arch + "';; esac\n"
			if err := os.WriteFile(filepath.Join(bin, "uname"), []byte(uname), 0o700); err != nil {
				t.Fatal(err)
			}
			// Only these three extracted functions run: no installer, curl or network.
			script := function("detect_os") + function("detect_arch") + function("lookup_user_agent") + "\nos=$(detect_os)\narch=$(detect_arch)\nlookup_user_agent\n"
			cmd := exec.Command("/bin/sh", "-c", script)
			cmd.Env = []string{"PATH=" + bin, "INSTALL_DIR=" + bin, "CI=" + tc.ci, "SPACEDOCK_INTERNAL=" + tc.internal, "LC_ALL=C"}
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("installer functions: %v %s", err, out)
			}
			goos, goarch := strings.ToLower(tc.os), "amd64"
			if tc.arch == "aarch64" || tc.arch == "arm64" {
				goarch = "arm64"
			}
			env := map[string]string{"CI": tc.ci, "SPACEDOCK_INTERNAL": tc.internal}
			want := strings.Replace(strings.TrimSpace(string(out)), "spacedock-install/", "spacedock/", 1)
			if got := UserAgent("dev", goos, goarch, func(k string) string { return env[k] }); got != want {
				t.Fatalf("Go UA=%q shell UA=%q", got, want)
			}
		})
	}
	for version, want := range map[string]string{
		"v1.2.3": "1.2.3", "1.2.3-pre10+build": "1.2.3-pre10+build", "1.2.3\nInjected: value": "unknown",
		"1.2.3-" + strings.Repeat("a", 26): "1.2.3-" + strings.Repeat("a", 26), "1.2.3-" + strings.Repeat("a", 27): "unknown",
	} {
		if got := UserAgent(version, "linux", "arm64", func(string) string { return "" }); got != "spacedock/"+want+" (linux; arm64)" {
			t.Errorf("version=%q UA=%q", version, got)
		}
	}
}

func TestLatestAnswers(t *testing.T) {
	for _, tc := range []struct {
		name, channel string
		status        int
		body, want    string
	}{
		{"stable object", "stable", 200, `{"tag_name":"v1.5.0","prerelease":false}`, "1.5.0"},
		{"edge array", "edge", 200, `[{"tag_name":"v2.0.0-pre10","prerelease":true}]`, "2.0.0-pre10"},
		{"503 with a valid body", "stable", 503, `{"tag_name":"v9.0.0"}`, ""},
		{"junk", "stable", 200, `<html>v99.0.0</html>`, ""},
		{"junk tag", "edge", 200, `{"tag_name":"v2.0.0\nrun evil"}`, ""},
		{"valid JSON over the size cap", "stable", 200, `{"tag_name":"v9.0.0"}` + strings.Repeat(" ", maxAnswer), ""},
		{"draft", "edge", 200, `{"tag_name":"v2.0.0","draft":true}`, ""},
		{"prerelease flag on stable", "stable", 200, `{"tag_name":"v9.0.0","prerelease":true}`, ""},
		{"prerelease tag on stable", "stable", 200, `{"tag_name":"v9.0.0-pre1"}`, ""},
		{"build tag", "stable", 200, `{"tag_name":"v9.0.0+build"}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got, calls := lookup(tc.channel, "", tc.status, tc.body); got != tc.want || calls != 1 {
				t.Fatalf("version=%q calls=%d, want %q", got, calls, tc.want)
			}
		})
	}
}

func TestLatestRedirects(t *testing.T) {
	for _, tc := range []struct {
		channel, location, want string
		calls                   int
	}{
		{"stable", stableAPI, "9.0.0", 2},
		{"edge", edgeAPI, "9.0.0", 2},
		{"stable", edgeAPI, "", 1},
		{"edge", stableAPI, "", 1},
		{"stable", loop, "", 5},
		{"stable", "https://github.com/spacedock-dev/spacedock/releases/tag/v9.0.0", "", 1},
		{"stable", "https://evil.example/repos/spacedock-dev/spacedock/releases/latest", "", 1},
		{"stable", "https://api.github.com.evil.example/repos/spacedock-dev/spacedock/releases/latest", "", 1},
		{"stable", "https://api.github.com:443/repos/spacedock-dev/spacedock/releases/latest", "", 1},
		{"stable", "http://api.github.com/repos/spacedock-dev/spacedock/releases/latest", "", 1},
		{"stable", "https://user@api.github.com/repos/spacedock-dev/spacedock/releases/latest", "", 1},
		{"stable", stableAPI + "#fragment", "", 1},
		{"stable", stableAPI + "?count=1", "", 1},
		{"stable", stableAPI + "?", "", 1},
		{"stable", "https://api.github.com/repos/spacedock-dev/compass/releases/latest", "", 1},
		{"edge", edgeAPI + "&count=1", "", 1},
		{"edge", "https://api.github.com/repos/spacedock-dev/spacedock/releases?per_page=2", "", 1},
		{"edge", "https://api.github.com/repos/spacedock-dev/spacedock/releases?per_page=%31", "", 1},
		{"stable", "https://api.github.com/repos/spacedock-dev/compass/../spacedock/releases/latest", "", 1},
		{"stable", "https://api.github.com/repos/spacedock-dev/%73pacedock/releases/latest", "", 1},
		{"stable", "https://api.github.com/repos/spacedock-dev/spacedock/releases%2flatest", "", 1},
	} {
		t.Run(tc.channel+" "+tc.location, func(t *testing.T) {
			if got, calls := lookup(tc.channel, tc.location, 200, `{"tag_name":"v9.0.0"}`); got != tc.want || calls != tc.calls {
				t.Fatalf("version=%q calls=%d, want %q after %d", got, calls, tc.want, tc.calls)
			}
		})
	}
}
