// ABOUTME: D5 sync test — the FO shared-core's stamped minor must equal the
// ABOUTME: vendored plugin manifests' own minor; the D4 tombstones stay frozen.
package contractlint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/contract"
	"github.com/spacedock-dev/spacedock/internal/release"
)

// This is the legitimate quarantined instruction-file read this package exists
// for: it binds TWO INDEPENDENT values that can diverge — the prose literal a
// release stamps into the FO shared-core, and the manifest minor the binary
// actually parses at the gate — so a hand edit to either side, or a forgotten
// stamp call site, goes red immediately. It is not prose-grep: neither value is
// asserted to equal text the check itself wrote.

// vendoredManifest is the subset of plugin.json fields the sync test reads.
type vendoredManifest struct {
	Version string `json:"version"`
}

// readVendoredManifest reads and parses the named repo plugin manifest
// (".claude-plugin" or ".codex-plugin").
func readVendoredManifest(t *testing.T, pluginDir string) vendoredManifest {
	t.Helper()
	path := filepath.Join(repoRoot(t), pluginDir, "plugin.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var m vendoredManifest
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return m
}

// foSharedCoreProse reads the FO shared-core reference file — the boot-resident
// contract carrying D5's single release-stamped "required binary minor" literal.
func foSharedCoreProse(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join(foReferenceDir(t), "first-officer-shared-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

// TestProseMinorMatchesVendoredManifestMinor is the D5 sync test: the FO
// shared-core's stamped minor (release.ProseMinor) must equal the major.minor
// each vendored plugin manifest's own `version` field parses to. These are two
// independent artifacts a release stamps together (D5's "all three stamp call
// sites carry the prose file") — a stamp that forgets one, or a hand edit to
// either, diverges them and this test catches it on every branch.
func TestProseMinorMatchesVendoredManifestMinor(t *testing.T) {
	proseMinor, err := release.ProseMinor(foSharedCoreProse(t))
	if err != nil {
		t.Fatalf("read FO shared-core stamped minor: %v", err)
	}

	for _, pluginDir := range []string{".claude-plugin", ".codex-plugin"} {
		m := readVendoredManifest(t, pluginDir)
		major, minor, ok := contract.ParseMajorMinor(m.Version)
		if !ok {
			t.Fatalf("%s version %q does not parse as major.minor", pluginDir, m.Version)
		}
		manifestMinor := fmt.Sprintf("%d.%d", major, minor)
		if manifestMinor != proseMinor {
			t.Errorf("%s: manifest minor %s does not match the FO shared-core's stamped minor %s — "+
				"a release stamped one side and not the other, or a hand edit drifted them apart",
				pluginDir, manifestMinor, proseMinor)
		}
	}
}
