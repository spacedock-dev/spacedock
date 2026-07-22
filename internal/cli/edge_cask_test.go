// ABOUTME: Table test for isEdgeCaskPath — the resolved-path check that flips the
// ABOUTME: too-old-binary remedy to the spacedock@next formula for an edge install.
package cli

import "testing"

// TestIsEdgeCaskPath checks that only a resolved path under a
// `Caskroom/spacedock@next/` segment is classified as the edge cask; the stable
// cask, a source checkout, and an empty path are not.
func TestIsEdgeCaskPath(t *testing.T) {
	cases := []struct {
		name string
		path string
		want bool
	}{
		{"edge cask", "/opt/homebrew/Caskroom/spacedock@next/0.26.0-pre0/spacedock", true},
		{"stable cask", "/usr/local/Caskroom/spacedock/0.25.0/spacedock", false},
		{"source checkout", "/Users/x/git/spacedock/spacedock", false},
		{"empty path", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isEdgeCaskPath(c.path); got != c.want {
				t.Fatalf("isEdgeCaskPath(%q) = %v, want %v", c.path, got, c.want)
			}
		})
	}
}
