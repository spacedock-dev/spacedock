package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type suiteRoundTrip func(*http.Request) (*http.Response, error)

func (f suiteRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

const suiteHint = "Spacedock 9.0.0 is available. Run `spacedock update`.\n"

// suiteTestDeps points the check at transport, a fixed clock and a terminal,
// with a fresh state directory and an installed 1.0.0.
func suiteTestDeps(t *testing.T, transport http.RoundTripper) *suiteDeps {
	t.Helper()
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Setenv("SPACEDOCK_NO_UPDATE_CHECK", "")
	withVersion(t, "1.0.0")
	d := &suiteDeps{
		client:   &http.Client{Transport: transport},
		now:      func() time.Time { return time.Date(2026, 11, 17, 23, 59, 0, 0, time.UTC) },
		terminal: func(io.Writer) bool { return true },
	}
	old := suiteCheckDeps
	suiteCheckDeps = func() *suiteDeps { return d }
	t.Cleanup(func() { suiteCheckDeps = old })
	return d
}

// suiteAnswer answers every lookup with body and records each request.
func suiteAnswer(body string, seen *[]*http.Request) http.RoundTripper {
	return suiteRoundTrip(func(r *http.Request) (*http.Response, error) {
		*seen = append(*seen, r)
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
}

// suiteLaunch runs a claude, codex or pi front door with a ready plugin fixture.
func suiteLaunch(t *testing.T, host string, args ...string) string {
	t.Helper()
	var out, errout bytes.Buffer
	code := 0
	if host == "pi" {
		repo, pkg := t.TempDir(), t.TempDir()
		writePiSkillFixtures(t, repo)
		writePiSubagentsFixtures(t, pkg)
		env := append(piTestEnv(pkg, t.TempDir()), "XDG_STATE_HOME="+os.Getenv("XDG_STATE_HOME"),
			"SPACEDOCK_NO_UPDATE_CHECK="+os.Getenv("SPACEDOCK_NO_UPDATE_CHECK"))
		code = runPi(context.Background(), append([]string{"--plugin-dir", repo}, args...), t.TempDir(), env, piSafehouseReadyOps(repo, pkg), &out, &errout)
	} else {
		door := runClaude
		if host == "codex" {
			door = runCodex
		}
		code = door(context.Background(), append([]string{"--skip-compat-check"}, args...), t.TempDir(), &fakeHost{}, func(string) (string, error) { return "", os.ErrNotExist }, &out, &errout)
	}
	if code != 0 {
		t.Fatalf("%s launch=%d stderr=%s", host, code, &errout)
	}
	return out.String()
}

func TestSuiteCheckFrontDoorArgs(t *testing.T) {
	quiet := []string{"--json", "--jsonl", "--output-format=stream-json", "--output-format", "-p", "-phello", "--print",
		"--input-format", "exec", "e", "review", "completion", "completions", "--non-interactive", "--mode",
		"--version", "-v", "-V", "--help", "-h"}
	type row struct {
		host string
		args []string
		env  string
		hint bool
	}
	rows := []row{
		{"claude", nil, "", true}, {"codex", nil, "", true}, {"pi", nil, "", true},
		{"claude", []string{"task", "--", "--model", "x"}, "", true},
		{"claude", []string{"--no-update-check"}, "", false}, {"pi", []string{"--no-update-check"}, "", false},
		{"claude", nil, "1", false}, {"pi", nil, "1", false},
	}
	for _, arg := range quiet {
		rows = append(rows, row{"claude", []string{"--", arg}, "", false})
	}
	for _, tc := range rows {
		t.Run(tc.host+" "+strings.Join(tc.args, " ")+" env="+tc.env, func(t *testing.T) {
			var seen []*http.Request
			suiteTestDeps(t, suiteAnswer(`{"tag_name":"v9.0.0"}`, &seen))
			t.Setenv("SPACEDOCK_NO_UPDATE_CHECK", tc.env)
			want, requests := "", 0
			if tc.hint {
				want, requests = suiteHint, 1
			}
			if got := suiteLaunch(t, tc.host, tc.args...); got != want || len(seen) != requests {
				t.Fatalf("hint=%q requests=%d, want %q after %d", got, len(seen), want, requests)
			}
		})
	}
}

func TestSuiteCheckOncePerUTCDay(t *testing.T) {
	var seen []*http.Request
	var d *suiteDeps
	d = suiteTestDeps(t, suiteRoundTrip(func(r *http.Request) (*http.Response, error) {
		day := d.now().UTC().Format("2006-01-02")
		if _, err := os.Stat(filepath.Join(os.Getenv("XDG_STATE_HOME"), "spacedock", "version-checks", day)); err != nil {
			t.Errorf("day not claimed before the lookup: %v", err)
		}
		return suiteAnswer(`{"tag_name":"v9.0.0"}`, &seen).RoundTrip(r)
	}))
	for _, step := range []struct {
		name string
		now  time.Time
		hint string
	}{
		{"first launch", time.Date(2026, 11, 17, 23, 59, 0, 0, time.UTC), suiteHint},
		{"same day", time.Date(2026, 11, 17, 23, 59, 30, 0, time.UTC), ""},
		{"next local day, same UTC day", time.Date(2026, 11, 18, 0, 30, 0, 0, time.FixedZone("+1", 3600)), ""},
		{"same local day, next UTC day", time.Date(2026, 11, 17, 23, 30, 0, 0, time.FixedZone("-1", -3600)), suiteHint},
	} {
		d.now = func() time.Time { return step.now }
		if got := suiteLaunch(t, "claude"); got != step.hint {
			t.Fatalf("%s: hint=%q, want %q", step.name, got, step.hint)
		}
	}
	if len(seen) != 2 {
		t.Fatalf("requests=%d, want 2", len(seen))
	}
}

func TestSuiteCheckRequest(t *testing.T) {
	for _, tc := range []struct{ branch, version, body, query, hint string }{
		{"main", "1.0.0", `{"tag_name":"v9.0.0"}`, "", suiteHint},
		{"main", "1.0.0", `{"tag_name":"v9.0.0-pre1"}`, "", ""},
		{"main", "v1.0.0", `{"tag_name":"v9.0.0"}`, "", suiteHint},
		{"main", "9.0.0", `{"tag_name":"v9.0.0"}`, "", ""},
		{"next", "9.0.0-pre2", `[{"tag_name":"v9.0.0-pre10"}]`, "prerelease=1", "Spacedock 9.0.0-pre10 is available. Run `spacedock update`.\n"},
		{"next", "9.0.0", `[{"tag_name":"v9.0.0-pre1"}]`, "prerelease=1", ""},
	} {
		t.Run(tc.branch+" "+tc.version+" "+tc.body, func(t *testing.T) {
			old := devBranch
			devBranch = tc.branch
			t.Cleanup(func() { devBranch = old })
			var seen []*http.Request
			suiteTestDeps(t, suiteAnswer(tc.body, &seen))
			withVersion(t, tc.version)
			t.Setenv("CI", " \tTrUe\n")
			t.Setenv("SPACEDOCK_INTERNAL", "1")
			if got := suiteLaunch(t, "claude"); got != tc.hint || len(seen) != 1 {
				t.Fatalf("hint=%q requests=%d, want %q", got, len(seen), tc.hint)
			}
			r := seen[0]
			ua := fmt.Sprintf("spacedock/%s (%s; %s; ci; internal)", strings.TrimPrefix(tc.version, "v"), runtime.GOOS, runtime.GOARCH)
			if r.Method != "GET" || r.URL.Host != "spacedock.md" || r.URL.Path != "/latest-release/spacedock" || r.URL.RawQuery != tc.query ||
				r.Body != nil || len(r.Header) != 1 || r.Header.Get("User-Agent") != ua {
				t.Fatalf("request %s %s headers=%v, want query %q and User-Agent %q", r.Method, r.URL, r.Header, tc.query, ua)
			}
		})
	}
}

// A hung lookup ends at the deadline, prints nothing, and is not retried that day.
func TestSuiteCheckDeadline(t *testing.T) {
	var requests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		// Answer late rather than never, so a missing deadline fails on time
		// instead of hanging the suite.
		select {
		case <-r.Context().Done():
		case <-time.After(3 * suiteCheckDeadline):
		}
	}))
	defer srv.Close()
	suiteTestDeps(t, suiteRoundTrip(func(r *http.Request) (*http.Response, error) {
		local := r.Clone(r.Context())
		local.URL.Scheme, local.URL.Host = "http", strings.TrimPrefix(srv.URL, "http://")
		return srv.Client().Transport.RoundTrip(local)
	}))
	start := time.Now()
	if got := suiteLaunch(t, "claude"); got != "" || time.Since(start) > suiteCheckDeadline+500*time.Millisecond {
		t.Fatalf("hung lookup: hint=%q after %s", got, time.Since(start))
	}
	if got := suiteLaunch(t, "claude"); got != "" || requests.Load() != 1 {
		t.Fatalf("retried on the same day: hint=%q requests=%d", got, requests.Load())
	}
}
