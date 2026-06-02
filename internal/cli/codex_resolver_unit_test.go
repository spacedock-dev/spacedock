// ABOUTME: Unit tests for the codex resolver helpers — version-dir selection
// ABOUTME: (semver order), install-listing parse, CODEX_HOME, cache degradation.
package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLatestVersionDirSemverOrder pins the semver-aware selection: with a stale
// cache holding both 0.9.0 and 0.10.0, the resolver must pick 0.10.0 (the newer
// install) rather than the lexically-greater 0.9.0.
func TestLatestVersionDirSemverOrder(t *testing.T) {
	root := t.TempDir()
	for _, v := range []string{"0.9.0", "0.10.0", "0.12.1"} {
		if err := os.Mkdir(filepath.Join(root, v), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	got, err := latestVersionDir(root)
	if err != nil {
		t.Fatalf("latestVersionDir errored: %v", err)
	}
	if filepath.Base(got) != "0.12.1" {
		t.Fatalf("latestVersionDir = %q, want the 0.12.1 dir", got)
	}
}

// TestLatestVersionDirSemverNotLexical isolates the lexical-vs-semver hazard:
// 0.9.0 sorts lexically AFTER 0.10.0, so a string compare would wrongly pick the
// older version. The semver compare must pick 0.10.0.
func TestLatestVersionDirSemverNotLexical(t *testing.T) {
	root := t.TempDir()
	for _, v := range []string{"0.9.0", "0.10.0"} {
		if err := os.Mkdir(filepath.Join(root, v), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	got, err := latestVersionDir(root)
	if err != nil {
		t.Fatalf("latestVersionDir errored: %v", err)
	}
	if filepath.Base(got) != "0.10.0" {
		t.Fatalf("latestVersionDir = %q, want the 0.10.0 dir (lexical compare would wrongly pick 0.9.0)", got)
	}
}

// TestLatestVersionDirSingleVersion is the live-cache shape today: a single
// version dir is returned as-is.
func TestLatestVersionDirSingleVersion(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "0.12.1"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := latestVersionDir(root)
	if err != nil {
		t.Fatalf("latestVersionDir errored: %v", err)
	}
	if filepath.Base(got) != "0.12.1" {
		t.Fatalf("latestVersionDir = %q, want the 0.12.1 dir", got)
	}
}

// TestLatestVersionDirAbsentRoot: a missing cache root is not an error — it is
// the no-install degradation state, returning "" with no error.
func TestLatestVersionDirAbsentRoot(t *testing.T) {
	got, err := latestVersionDir(filepath.Join(t.TempDir(), "no-such-cache"))
	if err != nil {
		t.Fatalf("latestVersionDir errored on absent root: %v", err)
	}
	if got != "" {
		t.Fatalf("latestVersionDir = %q, want empty for an absent cache root", got)
	}
}

// TestLatestVersionDirNoSubdirs: a present root with no version subdirectories
// (only files) returns "" with no error.
func TestLatestVersionDirNoSubdirs(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a-file"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := latestVersionDir(root)
	if err != nil {
		t.Fatalf("latestVersionDir errored: %v", err)
	}
	if got != "" {
		t.Fatalf("latestVersionDir = %q, want empty when root has no subdirectories", got)
	}
}

// TestCodexEntryInstalled exercises the `codex plugin list` text parse on the
// tolerated legacy paren form. Current codex renders the comma/table form
// `<id>  installed, enabled  <ver>  <path>` (covered by
// TestCodexEntryInstalledRealFormats); the paren form `<id> (installed[, enabled])`
// is the older rendering the predicate still accepts. An installed entry's
// status field reads `installed` after stripping surrounding `()`; a
// not-installed entry, or a listing without the id, does not match.
func TestCodexEntryInstalled(t *testing.T) {
	cases := []struct {
		name    string
		listing string
		want    bool
	}{
		{"installed", "  spacedock@spacedock (installed)", true},
		{"installed-enabled", "  spacedock@spacedock (installed, enabled)", true},
		{"not-installed", "  spacedock@spacedock (not installed)", false},
		{"other-plugin", "  other@market (installed)", false},
		{"empty", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := codexEntryInstalled(tc.listing, "spacedock@spacedock"); got != tc.want {
				t.Fatalf("codexEntryInstalled(%q) = %v, want %v", tc.listing, got, tc.want)
			}
		})
	}
}

// TestCodexEntryInstalledRealFormats pins the parse to the formats codex
// actually emits. The current comma/table form (`<id>  installed, enabled  <ver>
// <path>`) renders the status as `installed,` with NO parens, inside a
// column-aligned table with a header row and a marketplace PATH cell that
// contains the bare marketplace word `spacedock`. The predicate must field-match
// the id exactly and read the next field as the status, so it accepts the comma
// form and the legacy paren form, rejects `not installed`, a foreign id, and the
// PATH line, and never matches the bare substring `installed` inside
// `not installed`.
func TestCodexEntryInstalledRealFormats(t *testing.T) {
	// Live-captured block from `codex plugin list` (codex-cli 0.136.0): a
	// column-aligned table with a header row, the spacedock data row (comma
	// status), a not-installed row, and the marketplace PATH cell containing
	// the bare word `spacedock`.
	commaTable := "" +
		"PLUGIN               STATUS              VERSION  PATH\n" +
		"browser@openai-bundled  installed, enabled  26.519.22136  /Users/clkao/.codex/.tmp/bundled-marketplaces/openai-bundled/plugins/browser\n" +
		"chrome@openai-bundled  not installed  /Users/clkao/.codex/.tmp/bundled-marketplaces/openai-bundled/plugins/chrome\n" +
		"\n" +
		"PLUGIN               STATUS              VERSION  PATH\n" +
		"spacedock@spacedock  installed, enabled  0.12.1   /Users/clkao/.codex/.tmp/marketplaces/spacedock/plugins/spacedock\n"

	cases := []struct {
		name    string
		listing string
		want    bool
	}{
		// The real comma/table form codex emits today — this is the case the
		// paren-literal predicate misses, shipping the no-plugin-found drift.
		{"real-comma-table", commaTable, true},
		// The marketplace PATH cell contains the bare word `spacedock` but the
		// id `spacedock@spacedock` is not a field on this line — no false match.
		{"marketplace-path-line", "  not a row but a path /Users/clkao/.codex/.tmp/marketplaces/spacedock/plugins/spacedock", false},
		// `not installed` must reject: the status field after the id is `not`,
		// and the bare substring `installed` inside `not installed` must NOT match.
		{"not-installed-comma", "spacedock@spacedock  not installed  0.12.1  /some/path", false},
		// Legacy paren forms stay accepted (codex-revert tolerance).
		{"legacy-paren-installed", "  spacedock@spacedock (installed)", true},
		{"legacy-paren-enabled", "  spacedock@spacedock (installed, enabled)", true},
		// A foreign id that is not spacedock@spacedock must not match even when installed.
		{"foreign-id-installed", "other@market  installed, enabled  1.0.0  /some/path", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := codexEntryInstalled(tc.listing, "spacedock@spacedock"); got != tc.want {
				t.Fatalf("codexEntryInstalled(%q) = %v, want %v", tc.listing, got, tc.want)
			}
		})
	}
}

// TestCodexCacheManifestResolvesCachedPluginJSON drives the cache-resolution
// tail of resolveCodexManifest — the half the predicate test does not cover.
// With a temp CODEX_HOME holding the real install layout
// (plugins/cache/spacedock/spacedock/<ver>/.codex-plugin/plugin.json), the
// resolver must return that exact plugin.json path. This is what proves the
// predicate fix alone is insufficient: the front door only launches if this
// cache tail lands on a real manifest.
func TestCodexCacheManifestResolvesCachedPluginJSON(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CODEX_HOME", home)
	manifest := filepath.Join(home, "plugins", "cache", "spacedock", "spacedock", "0.12.1", ".codex-plugin", "plugin.json")
	if err := os.MkdirAll(filepath.Dir(manifest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifest, []byte(`{"name":"spacedock"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := codexCacheManifest()
	if err != nil {
		t.Fatalf("codexCacheManifest errored: %v", err)
	}
	if got != manifest {
		t.Fatalf("codexCacheManifest = %q, want %q", got, manifest)
	}
}

// TestCodexCacheManifestAbsentCache: with no cached install under CODEX_HOME,
// the cache tail degrades to "" (no error) — the no-cached-manifest state.
func TestCodexCacheManifestAbsentCache(t *testing.T) {
	t.Setenv("CODEX_HOME", t.TempDir())
	got, err := codexCacheManifest()
	if err != nil {
		t.Fatalf("codexCacheManifest errored on absent cache: %v", err)
	}
	if got != "" {
		t.Fatalf("codexCacheManifest = %q, want empty for an absent cache", got)
	}
}

// TestCodexHomeFromEnv: CODEX_HOME takes precedence over ~/.codex.
func TestCodexHomeFromEnv(t *testing.T) {
	t.Setenv("CODEX_HOME", "/tmp/custom-codex")
	if got := codexHome(); got != "/tmp/custom-codex" {
		t.Fatalf("codexHome() = %q, want /tmp/custom-codex", got)
	}
}

// TestCodexHomeDefault: with CODEX_HOME unset, codexHome resolves <home>/.codex.
func TestCodexHomeDefault(t *testing.T) {
	t.Setenv("CODEX_HOME", "")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("no user home dir: %v", err)
	}
	want := filepath.Join(home, ".codex")
	if got := codexHome(); got != want {
		t.Fatalf("codexHome() = %q, want %q", got, want)
	}
}
