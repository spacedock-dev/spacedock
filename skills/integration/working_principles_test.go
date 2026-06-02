// ABOUTME: AC-1 text-presence audit — the team's proven working habits ship in
// ABOUTME: the three instruction files in plain language with zero insider jargon.
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// shippedInstructionFiles is the trio of instruction surfaces this task encodes
// the working principles into: the workflow guide a captain reads, the FO
// operating contract, and the worker (ensign) contract. The map value is a
// human label used in failure messages.
func shippedInstructionFiles(t *testing.T) map[string]string {
	t.Helper()
	root := skillsRoot(t)
	repo := repoRoot(t)
	paths := map[string]string{
		"workflow guide (docs/dev/README.md)":        filepath.Join(repo, "docs", "dev", "README.md"),
		"FO contract (first-officer-shared-core.md)": filepath.Join(root, "first-officer", "references", "first-officer-shared-core.md"),
		"ensign contract (ensign-shared-core.md)":    filepath.Join(root, "ensign", "references", "ensign-shared-core.md"),
	}
	out := make(map[string]string, len(paths))
	for label, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read shipped instruction file %s (%s): %v", label, p, err)
		}
		out[label] = string(b)
	}
	return out
}

// TestShippedInstructionsCarryNoInsiderJargon locks the plain-language guarantee
// of AC-1: the three shipped instruction files contain zero insider-jargon
// tokens. "oracle" is the named token — the design proposal that seeded this work
// uses it pervasively for "external check," and that jargon must not leak into the
// instructions a clean-room contributor reads. The check is case-insensitive so
// "Oracle"/"ORACLE" cannot sneak through.
func TestShippedInstructionsCarryNoInsiderJargon(t *testing.T) {
	bannedJargon := []string{"oracle"}
	for label, content := range shippedInstructionFiles(t) {
		lower := strings.ToLower(content)
		for _, token := range bannedJargon {
			if strings.Contains(lower, token) {
				t.Errorf("%s contains insider-jargon token %q — shipped instructions must be plain language", label, token)
			}
		}
	}
}

// TestFOContractCarriesWorkingPrinciplesSection is a structural lint: it asserts
// the `## Working Principles` section heading exists in the FO contract. The
// heading is a structural anchor — a real on-disk section that other instruction
// text and refits reference by name — so deleting it is a non-paraphrasable
// mutation the lint catches. It deliberately does NOT grep the section's prose:
// wording is doc review, not a Go assertion, and a paraphrase that keeps the
// heading is not something this lint should pass or fail on.
func TestFOContractCarriesWorkingPrinciplesSection(t *testing.T) {
	fo := shippedInstructionFiles(t)["FO contract (first-officer-shared-core.md)"]
	if !strings.Contains(fo, "## Working Principles") {
		t.Errorf("FO contract missing the `## Working Principles` section heading")
	}
}
