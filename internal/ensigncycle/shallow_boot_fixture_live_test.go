package ensigncycle

import (
	"os"
	"path/filepath"
	"testing"
)

// writeShallowBootWorkflow seeds the shallow-boot fixture under root and returns
// the gate entity path. The fixture is intentionally mutation-free: merged-PR
// discovery and the registered startup hook are engage work, not greet work.
//
//spacedock:live-fixture id=boot/held-gate
func writeShallowBootWorkflow(t *testing.T, root string) shallowBootFixture {
	t.Helper()
	writeFile(t, filepath.Join(root, "README.md"), shallowBootReadme())
	gatePath := filepath.Join(root, "gate-check.md")
	writeFile(t, gatePath, shallowBootGateEntity())
	gitInit(t, root)

	return shallowBootFixture{gateEntityPath: gatePath}
}

func TestShallowBootFixtureContainsOnlyHeldGate(t *testing.T) {
	root := t.TempDir()
	fixture := writeShallowBootWorkflow(t, root)
	if fixture.gateEntityPath != filepath.Join(root, "gate-check.md") {
		t.Fatalf("gate entity path = %q, want the held gate fixture", fixture.gateEntityPath)
	}
	for _, unexpected := range []string{
		filepath.Join(root, "merged-pr.md"),
		filepath.Join(root, "_mods", "pr-merge.md"),
	} {
		if _, err := os.Stat(unexpected); !os.IsNotExist(err) {
			t.Fatalf("greet-only shallow boot fixture contains engage workload %s: %v", unexpected, err)
		}
	}
}
