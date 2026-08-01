// ABOUTME: AC-4 pin — the pi FO adapter's «worker.spawn» bullet binds the
// ABOUTME: artifact's agent/skill fields to the subagent spawn call, and the
// ABOUTME: namespaced agent string stays out of the pi dispatch path.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The two load-bearing sentences implementation added to the pi adapter's
// «worker.spawn» bullet. Contract-drift candidates: rewording either sentence
// silently unbinds the build artifact from the pi-subagents spawn call.
const (
	piSpawnAgentBindingSentence = "Pass the artifact's `agent` field as the subagent `agent` parameter verbatim."
	piSpawnSkillBindingSentence = "When the artifact carries a `skill` field, pass it as the subagent `skill` parameter verbatim"
)

func TestPiAdapterWorkerSpawnBindsArtifactAgentAndSkill(t *testing.T) {
	text := readRepoFile(t, piFORuntimeRel)
	for _, want := range []string{piSpawnAgentBindingSentence, piSpawnSkillBindingSentence} {
		if !strings.Contains(text, want) {
			t.Errorf("%s «worker.spawn» binding missing sentence %q", piFORuntimeRel, want)
		}
	}
}

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
