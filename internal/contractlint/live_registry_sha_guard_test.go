package contractlint

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var reconciliationSHA = regexp.MustCompile("Registry reconciliation SHA: `([0-9a-f]{40})`")

// TestRuntimeLiveRegistryReconciliationSHA is independent of the semantic
// inventory check. It proves that no watched source changed after the recorded
// reconciliation commit and that a stale base names every changed path.
func TestRuntimeLiveRegistryReconciliationSHA(t *testing.T) {
	repo := repoRoot(t)
	doc, err := os.ReadFile(filepath.Join(repo, "docs", "runtime-live-ci.md"))
	if err != nil {
		t.Fatal(err)
	}
	match := reconciliationSHA.FindSubmatch(doc)
	if len(match) != 2 {
		t.Fatal("docs/runtime-live-ci.md has no 40-character registry reconciliation SHA")
	}
	base := string(match[1])
	if err := assertRegistryReconciliationCurrent(repo, base, "HEAD"); err != nil {
		t.Fatal(err)
	}

	stale := base + "^"
	changed, err := registryWatchedChanges(repo, stale, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) == 0 {
		t.Fatalf("stale reconciliation base %s did not expose a watched-path change", stale)
	}
	err = assertRegistryReconciliationCurrent(repo, stale, "HEAD")
	if err == nil {
		t.Fatalf("stale reconciliation base %s passed", stale)
	}
	for _, path := range changed {
		if !strings.Contains(err.Error(), path) {
			t.Errorf("stale reconciliation diagnostic omitted %s: %v", path, err)
		}
	}
}

func assertRegistryReconciliationCurrent(repo, base, head string) error {
	changed, err := registryWatchedChanges(repo, base, head)
	if err != nil {
		return err
	}
	if len(changed) != 0 {
		return fmt.Errorf("registry reconciliation is stale after %s; watched paths changed: %s", base, strings.Join(changed, ", "))
	}
	return nil
}

func registryWatchedChanges(repo, base, head string) ([]string, error) {
	cmd := exec.Command("git", "-C", repo, "diff", "--name-only", base+".."+head, "--", "internal/ensigncycle", "internal/livescenario", ".github/workflows/runtime-live-e2e.yml")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git diff registry reconciliation: %w: %s", err, strings.TrimSpace(string(out)))
	}
	var paths []string
	for _, path := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if path != "" {
			paths = append(paths, path)
		}
	}
	return paths, nil
}
