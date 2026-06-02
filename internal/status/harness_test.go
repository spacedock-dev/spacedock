// ABOUTME: Golden-harness shared helpers — pinned env, golden file IO, and the
// ABOUTME: normalization (timestamps, root prefix, realpath, state-backend) for compares.
package status

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// update regenerates checked-in golden files from the current native output
// instead of comparing against them. Run: go test ./internal/status -update
var update = flag.Bool("update", false, "regenerate golden files from native output")

// pinnedEnv returns the locale/id/timestamp-pinned environment the native runner
// must run under so env-dependent output is reproducible. The values mirror the
// test plan: PYTHONUTF8/LANG for locale, USER/actor/seed and
// SPACEDOCK_TEST_SD_B32_TIMESTAMP for sd-b32, HOME for the team probe, PATH for
// locating git/gh.
func pinnedEnv(t *testing.T) []string {
	t.Helper()
	return []string{
		"PYTHONUTF8=1",
		"LANG=C.UTF-8",
		"LC_ALL=C.UTF-8",
		"USER=pinned-actor",
		"SPACEDOCK_TEST_SD_B32_TIMESTAMP=2026-01-01T00:00:00.000000Z",
		"HOME=" + t.TempDir(), // empty team dir -> deterministic TEAM_STATE
		"PATH=" + os.Getenv("PATH"),
	}
}

// tsRe matches ISO-8601 UTC timestamps in BOTH the second-precision mutation
// shape and the microsecond sd-b32 shape. The trailing (\.\d+)? is required or
// the microsecond timestamp slips through un-normalized.
var tsRe = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z`)

// stateBackendLineRe matches the native-only STATE_BACKEND boot banner. The text
// --boot adds it to surface the state backend for a human reading the boot, but
// the Python oracle has no such line — the same kind of intentional native/oracle
// divergence as the dispatch fetch-line and state-commit guidance. The boot-text
// parity normalizers strip it from BOTH sides (the oracle never emits it) so the
// shared sections still byte-match; the JSON form is covered by json_boot_test.go.
var stateBackendLineRe = regexp.MustCompile(`STATE_BACKEND: [^\n]*\n`)

// stripStateBackend removes the native-only STATE_BACKEND boot banner so a body
// with it and the oracle body without it normalize to the same bytes.
func stripStateBackend(s string) string {
	return stateBackendLineRe.ReplaceAllString(s, "")
}

// bootNextIDLineRe matches ONLY the boot `NEXT_ID:` line's sd-b32 value — the
// minted candidate, which hashes the realpath'd workflow dir and so varies per
// checkout path. Anchored to the line prefix so it does NOT touch a stored
// fixture id elsewhere in the output (the DISPATCHABLE table's id column, a
// --resolve id= field): those are fixed fixture content and must freeze
// literally, so a wrong-but-valid id there fails its golden. The minting
// derivation itself is pinned path-independently by
// TestSDB32CandidateDerivationVector.
var bootNextIDLineRe = regexp.MustCompile(`(?m)^NEXT_ID: [0-9a-hjkmnp-tv-z]{24}$`)

// maskBootNextID replaces the minted value on the boot NEXT_ID line with a
// placeholder, leaving every other sd-b32 token (fixed fixture ids) intact.
func maskBootNextID(s string) string {
	return bootNextIDLineRe.ReplaceAllString(s, "NEXT_ID: <ID>")
}

// normalize applies the test-plan normalization to output before comparison:
// the timestamp placeholder, and root-prefix placeholders for the workflow root
// (both as-spelled and realpath'd, since --resolve workflow= is realpath'd on
// macOS /var->/private/var while path=/archived: are not).
func normalize(s, root string) string {
	s = tsRe.ReplaceAllString(s, "<TS>")
	if root != "" {
		real := realpath(root)
		// Replace the realpath'd spelling first (longer/more specific on macOS),
		// then the as-spelled root, so both map to the same placeholder.
		if real != root {
			s = strings.ReplaceAll(s, real, "<ROOT>")
		}
		s = strings.ReplaceAll(s, root, "<ROOT>")
	}
	return s
}

// realpath resolves symlinks the way the oracle's os.path.realpath does, so the
// expected workflow= prefix accounts for the macOS /var->/private/var rewrite.
func realpath(p string) string {
	resolved, err := filepath.EvalSymlinks(p)
	if err != nil {
		return p
	}
	return resolved
}

// writeGolden / readGolden manage checked-in golden files under testdata/golden.
func goldenPath(name string) string {
	return filepath.Join("testdata", "golden", name)
}

func writeGolden(t *testing.T, name, content string) {
	t.Helper()
	if err := os.WriteFile(goldenPath(name), []byte(content), 0o644); err != nil {
		t.Fatalf("write golden %s: %v", name, err)
	}
}

func readGolden(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(goldenPath(name))
	if err != nil {
		t.Fatalf("read golden %s (regenerate with -update): %v", name, err)
	}
	return string(b)
}
