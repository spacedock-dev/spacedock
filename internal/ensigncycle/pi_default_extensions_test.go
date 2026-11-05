package ensigncycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// piExtensionRoots holds the installed package roots the isolated Pi home links
// into its own npm store. Empty means "not installed/registered".
type piExtensionRoots struct {
	subagents string
	intercom  string
}

// piDefaultExtensionRoots resolves the pi-subagents and pi-intercom package roots
// from the real agent directory's settings.json — Pi's own discovery chain
// (settings `packages` -> agentDir/npm/node_modules/<name> -> package.json name).
// It reads the real installed location: npm: entries resolve through the
// agentDir npm path, and local (absolute or agentDir-relative) entries resolve to
// their path. Every candidate is matched by its package.json `name`, never by a
// hard-coded directory or a sibling-directory guess. A missing package yields "".
func piDefaultExtensionRoots(t *testing.T, agentDir string) piExtensionRoots {
	t.Helper()
	var roots piExtensionRoots
	for _, source := range piSettingsPackageSources(t, agentDir) {
		root := piResolveSettingsPackageRoot(source, agentDir, os.Getenv("HOME"))
		if root == "" {
			continue
		}
		switch piPackageManifestName(root) {
		case "pi-subagents":
			if roots.subagents == "" {
				roots.subagents = root
			}
		case "pi-intercom":
			if roots.intercom == "" {
				roots.intercom = root
			}
		}
	}
	return roots
}

// piSettingsPackageSources returns the settings.json `packages` entries as
// strings. Both the bare-string and {source: ...} entry shapes are accepted, the
// same as the launcher's reader.
func piSettingsPackageSources(t *testing.T, agentDir string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(agentDir, "settings.json"))
	if err != nil {
		return nil
	}
	var settings struct {
		Packages []json.RawMessage `json:"packages"`
	}
	if json.Unmarshal(data, &settings) != nil {
		return nil
	}
	var sources []string
	for _, raw := range settings.Packages {
		var s string
		if json.Unmarshal(raw, &s) == nil && s != "" {
			sources = append(sources, s)
			continue
		}
		var obj struct {
			Source string `json:"source"`
		}
		if json.Unmarshal(raw, &obj) == nil && obj.Source != "" {
			sources = append(sources, obj.Source)
		}
	}
	return sources
}

// piPackageManifestName reads the package.json `name` at a package root; "" when
// the manifest is absent or unreadable.
func piPackageManifestName(root string) string {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return ""
	}
	var pkg struct {
		Name string `json:"name"`
	}
	if json.Unmarshal(data, &pkg) != nil {
		return ""
	}
	return pkg.Name
}

// piResolveSettingsPackageRoot mirrors the launcher's resolution for the sources
// this harness registers: npm:<name>[@ver] resolves through
// <agentDir>/npm/node_modules/<name>; absolute, ~, and agentDir-relative paths
// resolve to that path. A file: prefix is stripped to its path — Pi's loader does
// not load file: packages, so callers must register the bare path instead.
func piResolveSettingsPackageRoot(source, agentDir, home string) string {
	s := strings.TrimSpace(source)
	if s == "" {
		return ""
	}
	if name, ok := strings.CutPrefix(s, "npm:"); ok {
		return filepath.Join(agentDir, "npm", "node_modules", piNpmPackageName(name))
	}
	norm := strings.TrimPrefix(s, "file:")
	if norm == "" {
		return ""
	}
	if norm == "~" {
		return home
	}
	if strings.HasPrefix(norm, "~/") {
		return filepath.Join(home, norm[2:])
	}
	if filepath.IsAbs(norm) {
		return norm
	}
	if norm == "." || norm == ".." || strings.HasPrefix(norm, "./") || strings.HasPrefix(norm, "../") {
		return filepath.Join(agentDir, norm)
	}
	return ""
}

// piNpmPackageName strips an optional @version suffix from an npm: source; ""
// when nothing usable remains. "@scope/pkg@1.0" -> "@scope/pkg".
func piNpmPackageName(spec string) string {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return ""
	}
	if i := strings.LastIndex(spec, "@"); i > 0 {
		spec = spec[:i]
	}
	if spec == "" || strings.HasPrefix(spec, "@") && !strings.Contains(spec, "/") {
		return ""
	}
	return spec
}

// piIsolatedExtensionRoots selects the package roots to link into the isolated
// home: each nonempty explicit operator override wins its own package
// independently, and the other package still comes from real-agent discovery. An
// override is never inferred from the sibling package.
func piIsolatedExtensionRoots(t *testing.T, realAgentDir string) piExtensionRoots {
	t.Helper()
	roots := piDefaultExtensionRoots(t, realAgentDir)
	if v := os.Getenv("PI_SUBAGENTS_PACKAGE_ROOT"); v != "" {
		roots.subagents = v
	}
	if v := os.Getenv("PI_INTERCOM_PACKAGE_ROOT"); v != "" {
		roots.intercom = v
	}
	return roots
}

// seedPiDefaultExtensions links both substrate packages under
// <piHome>/npm/node_modules so Pi's npm-path discovery finds them, and writes
// <piHome>/settings.json registering BOTH npm:pi-subagents and npm:pi-intercom.
// The Spacedock checkout stays exactly one absolute path entry (repo), never a
// file: entry. Root symlinks alone do not register the packages with Pi.
func seedPiDefaultExtensions(t *testing.T, piHome, repo string, roots piExtensionRoots) {
	t.Helper()
	if roots.subagents == "" || roots.intercom == "" {
		t.Fatalf("isolated home needs both package roots; got subagents=%q intercom=%q", roots.subagents, roots.intercom)
	}
	store := filepath.Join(piHome, "npm", "node_modules")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, target := range map[string]string{"pi-subagents": roots.subagents, "pi-intercom": roots.intercom} {
		link := filepath.Join(store, name)
		if err := os.Remove(link); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		if err := os.Symlink(target, link); err != nil {
			t.Fatalf("link %s -> %s: %v", link, target, err)
		}
	}
	settings, err := json.Marshal(map[string]any{"packages": []string{"npm:pi-subagents", "npm:pi-intercom", repo}})
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(piHome, "settings.json"), string(settings)+"\n")
}

// seedPiIsolatedHome is the single isolated-home setup contract for Pi live runs.
// It allocates piHome = <cleanHome>/.pi/agent (so the launcher's HOME-based
// package-root preflight and Pi's PI_CODING_AGENT_DIR package loader agree),
// discovers the real installed package roots from realAgentDir, links and
// registers them, and returns the piHome to pass as PI_CODING_AGENT_DIR.
func seedPiIsolatedHome(t *testing.T, cleanHome, realAgentDir, repo string) string {
	t.Helper()
	piHome := filepath.Join(cleanHome, ".pi", "agent")
	seedPiDefaultExtensions(t, piHome, repo, piIsolatedExtensionRoots(t, realAgentDir))
	return piHome
}

// realPiAgentDir is the parent's real Pi agent directory: an explicit
// PI_CODING_AGENT_DIR when set, otherwise the conventional <realHome>/.pi/agent.
// Discovery reads the installed roots from this directory, never from the
// isolated clean HOME.
func realPiAgentDir(realHome string) string {
	if d := os.Getenv("PI_CODING_AGENT_DIR"); d != "" {
		return d
	}
	return filepath.Join(realHome, ".pi", "agent")
}

func writePiPackageRoot(t *testing.T, dir, name string) {
	t.Helper()
	writeFile(t, filepath.Join(dir, "package.json"), `{"name":`+jsonString(name)+`}`+"\n")
}

func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// TestPiDefaultExtensionRootsReadsRealInstalledLocation proves discovery reads
// the real installed location rather than a hard-coded path: with settings in a
// custom agentDir, both roots resolve through that agentDir's npm path.
func TestPiDefaultExtensionRootsReadsRealInstalledLocation(t *testing.T) {
	agentDir := filepath.Join(t.TempDir(), "custom-agent")
	sub := filepath.Join(agentDir, "npm", "node_modules", "pi-subagents")
	ic := filepath.Join(agentDir, "npm", "node_modules", "pi-intercom")
	writePiPackageRoot(t, sub, "pi-subagents")
	writePiPackageRoot(t, ic, "pi-intercom")
	writeFile(t, filepath.Join(agentDir, "settings.json"),
		`{"packages":["npm:pi-subagents@0.0.0-synthetic","npm:pi-intercom"]}`+"\n")

	roots := piDefaultExtensionRoots(t, agentDir)
	if roots.subagents != sub {
		t.Fatalf("subagents root = %q, want the custom agentDir npm path %q", roots.subagents, sub)
	}
	if roots.intercom != ic {
		t.Fatalf("intercom root = %q, want the custom agentDir npm path %q", roots.intercom, ic)
	}
}

// TestPiIsolatedHomeRegistersBothSubstratesAndAbsoluteSpacedock proves the
// isolated home's own settings register both npm substrates and keep the
// Spacedock checkout as one absolute path entry — never a file: entry.
func TestPiIsolatedHomeRegistersBothSubstratesAndAbsoluteSpacedock(t *testing.T) {
	realAgentDir := t.TempDir()
	sub := filepath.Join(realAgentDir, "npm", "node_modules", "pi-subagents")
	ic := filepath.Join(realAgentDir, "npm", "node_modules", "pi-intercom")
	writePiPackageRoot(t, sub, "pi-subagents")
	writePiPackageRoot(t, ic, "pi-intercom")
	writeFile(t, filepath.Join(realAgentDir, "settings.json"),
		`{"packages":["npm:pi-subagents","npm:pi-intercom"]}`+"\n")
	cleanHome := t.TempDir()
	repo := t.TempDir()

	piHome := seedPiIsolatedHome(t, cleanHome, realAgentDir, repo)
	if want := filepath.Join(cleanHome, ".pi", "agent"); piHome != want {
		t.Fatalf("piHome = %q, want %q", piHome, want)
	}
	sources := piSettingsPackageSources(t, piHome)
	want := []string{"npm:pi-subagents", "npm:pi-intercom", repo}
	if len(sources) != len(want) {
		t.Fatalf("isolated settings packages = %q, want %q", sources, want)
	}
	for i, w := range want {
		if sources[i] != w {
			t.Fatalf("isolated settings packages[%d] = %q, want %q", i, sources[i], w)
		}
	}
	for _, s := range sources {
		if strings.HasPrefix(s, "file:") {
			t.Fatalf("isolated settings must not carry a file: entry, got %q", s)
		}
	}
	if !filepath.IsAbs(repo) {
		t.Fatalf("test repo path must be absolute, got %q", repo)
	}
}

// TestPiIsolatedHomeNegativeControlDropsSubstrateRegistrations is the retained
// negative control: both root symlinks present, both npm registrations removed,
// both substrate packages disappear from Pi's discovery. This is the measured
// fact the fold rests on — root symlinks alone do not register resources.
func TestPiIsolatedHomeNegativeControlDropsSubstrateRegistrations(t *testing.T) {
	realAgentDir := t.TempDir()
	sub := filepath.Join(realAgentDir, "npm", "node_modules", "pi-subagents")
	ic := filepath.Join(realAgentDir, "npm", "node_modules", "pi-intercom")
	writePiPackageRoot(t, sub, "pi-subagents")
	writePiPackageRoot(t, ic, "pi-intercom")
	writeFile(t, filepath.Join(realAgentDir, "settings.json"),
		`{"packages":["npm:pi-subagents","npm:pi-intercom"]}`+"\n")
	cleanHome := t.TempDir()
	repo := t.TempDir()

	piHome := seedPiIsolatedHome(t, cleanHome, realAgentDir, repo)
	// Remove both npm registrations, keep both root symlinks.
	writeFile(t, filepath.Join(piHome, "settings.json"), `{"packages":[`+jsonString(repo)+`]}`+"\n")
	for _, name := range []string{"pi-subagents", "pi-intercom"} {
		if _, err := os.Stat(filepath.Join(piHome, "npm", "node_modules", name, "package.json")); err != nil {
			t.Fatalf("negative control expects the %s root symlink present: %v", name, err)
		}
	}
	roots := piDefaultExtensionRoots(t, piHome)
	if roots.subagents != "" || roots.intercom != "" {
		t.Fatalf("without npm registrations both substrates must disappear, got subagents=%q intercom=%q", roots.subagents, roots.intercom)
	}
}
