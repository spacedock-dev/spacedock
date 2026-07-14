// ABOUTME: AC-1 value check — the FO Startup recipe collapses to <=4 numbered prose
// ABOUTME: steps AND the boot-resident shared core is strictly smaller than its pre-change size.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// preChangeSharedCoreBytes is the byte size of first-officer-shared-core.md
// measured immediately before the "boot identifies, engage converges" edits, on
// the post-vcm-merge branch this implementation opened from (`wc -c` == 26755).
// AC-1's value half asserts the post-change file is STRICTLY smaller: a shorter
// recipe that grew the file fails.
// Re-baselined by the consolidated bridge-seam PR: shared-core gains the
// before-greet `«bridge.boot-liveness»` step + function (the one liveness signal
// that cannot wait for the event loop), mirroring the before-greet heartbeat step
// #435 added here. The recipe stayed lean (no new top-level Startup step; 2b is a
// sub-step); this is the new floor the shrink assertion guards from.
const preChangeSharedCoreBytes = 27200

// startupStepRe matches a top-level numbered Startup step: a line beginning with
// `N.` at column zero. Sub-bullets (indented `-`) and the discovery sub-cases are
// not top-level steps, so they do not count against the <=4 budget.
var startupStepRe = regexp.MustCompile(`(?m)^[0-9]+\. `)

// startupBlock returns the `## Startup` section body (heading to the next `## `),
// so the step count is scoped to the recipe, not the whole file.
func startupBlock(t *testing.T, body string) string {
	t.Helper()
	loc := regexp.MustCompile(`(?m)^## Startup$`).FindStringIndex(body)
	if loc == nil {
		t.Fatal("shared core has no `## Startup` section")
	}
	rest := body[loc[1]:]
	if end := regexp.MustCompile(`(?m)^## `).FindStringIndex(rest); end != nil {
		rest = rest[:end[0]]
	}
	return rest
}

// TestStartupRecipeCollapsedAndLeaner (AC-1) is the two-sided value check: the
// Startup recipe carries <=4 top-level numbered steps (down from 8) AND the whole
// boot-resident shared core is strictly fewer bytes than before the change. Both
// halves must hold — the sprint's principal prose-remover must not merely reshape
// the recipe while leaving the resident core the same size or larger.
func TestStartupRecipeCollapsedAndLeaner(t *testing.T) {
	root := repoRoot(t)
	data, err := os.ReadFile(filepath.Join(root, sharedCorePath()))
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}

	steps := startupStepRe.FindAllString(startupBlock(t, string(data)), -1)
	if len(steps) == 0 {
		t.Fatal("Startup section has no top-level numbered steps — extractor bug; the count check would pass vacuously")
	}
	if len(steps) > 4 {
		t.Errorf("Startup recipe has %d top-level numbered steps, want <=4", len(steps))
	}

	if got := len(data); got >= preChangeSharedCoreBytes {
		t.Errorf("shared core is %d bytes, want strictly < %d (pre-change) — the collapse must shrink the resident core, not grow it", got, preChangeSharedCoreBytes)
	}
}
