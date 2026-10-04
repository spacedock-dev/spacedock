//go:build live

// ABOUTME: The lane's single independent oracle: it queries the npm registry at
// ABOUTME: check time (never a stored copy) and fails when a pin has drifted.
package pilive

import (
	"encoding/json"
	"os/exec"
	"testing"
)

// TestPiLivePinsMatchRegistry is the lane's independent oracle. It queries the
// registry at check time — it deliberately stores no copy of the current
// numbers, because a stored copy is the duplication this package removes. A pin
// that no longer matches the published `latest` version or the published
// integrity fails here; update the single source in pilive.go when the family
// moves.
func TestPiLivePinsMatchRegistry(t *testing.T) {
	if _, err := exec.LookPath("npm"); err != nil {
		t.Skip("npm is not available; the registry oracle cannot run")
	}
	for _, pkg := range Packages {
		out, err := exec.Command("npm", "view", pkg.Spec, "dist-tags.latest", "version", "dist.integrity", "--json").Output()
		if err != nil {
			t.Fatalf("npm view %s: %v", pkg.Spec, err)
		}
		var meta struct {
			Latest    string `json:"dist-tags.latest"`
			Version   string `json:"version"`
			Integrity string `json:"dist.integrity"`
		}
		if err := json.Unmarshal(out, &meta); err != nil {
			t.Fatalf("parse npm view %s: %v", pkg.Spec, err)
		}
		if meta.Latest != pkg.Version {
			t.Errorf("%s: pinned %s but registry latest is %s; update internal/pilive/pilive.go", pkg.Spec, pkg.Version, meta.Latest)
		}
		if meta.Version != pkg.Version || meta.Integrity != pkg.Integrity {
			t.Errorf("%s: pinned %s %s but registry reports %s %s", pkg.Spec, pkg.Version, pkg.Integrity, meta.Version, meta.Integrity)
		}
	}
}
