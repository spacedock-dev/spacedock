// ABOUTME: The namespaced agent string stays out of the pi dispatch path:
// ABOUTME: pi-subagents resolves by directory basename, so `spacedock:ensign`
// ABOUTME: names nothing there.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPiDispatchPathBansNamespacedAgentName: pi-subagents resolves skills and
// agents by directory basename, so the namespaced `spacedock:ensign` string
// names nothing on pi — it must not appear in the pi adapter text or the
// piruntime transport wrapper (the build pi branch is pinned by dispatch tests;
// the host-neutral subagent_type identity is intentional and out of this ban).
func TestPiDispatchPathBansNamespacedAgentName(t *testing.T) {
	const banned = "spacedock:ensign"
	if text := readRepoFile(t, piFORuntimeRel); strings.Contains(text, banned) {
		t.Errorf("%s must not name a pi-subagents agent %q", piFORuntimeRel, banned)
	}
	root := repoRoot(t)
	entries, err := os.ReadDir(filepath.Join(root, "internal", "piruntime"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		rel := filepath.Join("internal", "piruntime", e.Name())
		if strings.Contains(readRepoFile(t, rel), banned) {
			t.Errorf("%s must not contain %q in the pi dispatch path", rel, banned)
		}
	}
}
