// ABOUTME: Live-run guard tests — the lead-with-`live` runtime-observable AC
// ABOUTME: classifier and the three-way ci-run:/session: citation resolution.
package status

import (
	"os"
	"path/filepath"
	"testing"
)

// TestClassifyLiveACs locks the explicit-convention classifier: a runtime-
// observable AC declares itself by leading its proof clause with `live`, and the
// classifier returns one flag per such AC carrying the trailing citation ref
// (`ci-run:<id>` / `session:<path>` / the empty/placeholder text). An offline
// AC (any non-`live` proof clause) is never returned; this is the explicit
// `live`-lead convention, not prose keyword inference.
func TestClassifyLiveACs(t *testing.T) {
	cases := []struct {
		name      string
		body      string
		wantFlags int
		wantCite  string // citation of the first flag, when wantFlags == 1
	}{
		{
			name:      "offline-go-test-not-live",
			body:      "**AC-1 — Producer exits.**\nVerified by: a Go unit test TestRoundTrip asserts the exit code.\n",
			wantFlags: 0,
		},
		{
			name:      "live-uncited-placeholder",
			body:      "**AC-1 — The model exits on the contract.**\nVerified by: live <no artifact yet>\n",
			wantFlags: 1,
			wantCite:  "<no artifact yet>",
		},
		{
			name:      "live-cited-ci-run",
			body:      "**AC-1 — The model exits on the contract.**\nVerified by: live ci-run:12345\n",
			wantFlags: 1,
			wantCite:  "ci-run:12345",
		},
		{
			name:      "live-cited-session",
			body:      "**AC-1 — The model exits on the contract.**\nVerified by: live session:/tmp/run.jsonl\n",
			wantFlags: 1,
			wantCite:  "session:/tmp/run.jsonl",
		},
		{
			// yy's shape: a runtime-observable AC declared live but cited only by
			// an offline-proof narrative / a pending placeholder.
			name:      "yy-shape-offline-narrative",
			body:      "**AC-1 — The fix actually closes the flake at runtime.**\nVerified by: live <none — passed offline + 3 audit cycles, merged pending-live-run>\n",
			wantFlags: 1,
			wantCite:  "<none — passed offline + 3 audit cycles, merged pending-live-run>",
		},
		{
			name:      "live-case-insensitive-and-no-leak-to-offline",
			body:      "**AC-1 — Runtime.**\nVerified by: Live ci-run:7\n\n**AC-2 — Offline.**\nVerified by: a `go vet ./...` pass.\n",
			wantFlags: 1,
			wantCite:  "ci-run:7",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			flags := classifyLiveACs(tc.body)
			if len(flags) != tc.wantFlags {
				t.Fatalf("want %d live flags, got %d: %+v", tc.wantFlags, len(flags), flags)
			}
			if tc.wantFlags == 1 && flags[0].Citation != tc.wantCite {
				t.Fatalf("want citation %q, got %q", tc.wantCite, flags[0].Citation)
			}
		})
	}
}

// TestResolveLiveCitationThreeWay locks AC-4's spike split for `ci-run:<id>`:
// a real id resolves (citedAndReal), a definitive 404 refuses (definitivelyAbsent),
// and a connectivity/auth error is INDETERMINATE (a tooling error, NOT a
// refusal) so a network blip cannot masquerade as a missing live run. The
// `gh`-shelling resolver is swapped for a stub so the unit test runs offline.
func TestResolveLiveCitationThreeWay(t *testing.T) {
	cases := []struct {
		name string
		stub ciRunResolver
		want liveResolutionKind
	}{
		{
			name: "real-id-cited-and-real",
			stub: func(id string) (bool, error) { return true, nil },
			want: citedAndReal,
		},
		{
			name: "404-definitively-absent-refuse",
			stub: func(id string) (bool, error) { return false, nil },
			want: definitivelyAbsent,
		},
		{
			name: "connectivity-error-indeterminate",
			stub: func(id string) (bool, error) { return false, errConnectivity },
			want: indeterminate,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := resolveLiveCitation("ci-run:12345", tc.stub)
			if res.kind != tc.want {
				t.Fatalf("want kind %v, got %v (err=%v)", tc.want, res.kind, res.err)
			}
		})
	}
}

// TestResolveLiveCitationSession locks the offline session: citation: a present
// .jsonl on disk resolves (citedAndReal), an absent path refuses
// (definitivelyAbsent). No `gh`, no network.
func TestResolveLiveCitationSession(t *testing.T) {
	dir := t.TempDir()
	present := filepath.Join(dir, "run.jsonl")
	if err := os.WriteFile(present, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	absent := filepath.Join(dir, "missing.jsonl")

	if res := resolveLiveCitation("session:"+present, nil); res.kind != citedAndReal {
		t.Fatalf("present session jsonl should resolve, got %v", res.kind)
	}
	if res := resolveLiveCitation("session:"+absent, nil); res.kind != definitivelyAbsent {
		t.Fatalf("absent session jsonl should refuse, got %v", res.kind)
	}
}

// TestResolveLiveCitationPlaceholder locks that a placeholder / empty citation
// is definitivelyAbsent (refused) — turning pending-live-run from a hopeful
// label into an enforced state.
func TestResolveLiveCitationPlaceholder(t *testing.T) {
	for _, cite := range []string{"<no artifact yet>", "<pending>", "", "<none — passed offline>"} {
		if res := resolveLiveCitation(cite, nil); res.kind != definitivelyAbsent {
			t.Fatalf("placeholder citation %q should be definitivelyAbsent, got %v", cite, res.kind)
		}
	}
}
