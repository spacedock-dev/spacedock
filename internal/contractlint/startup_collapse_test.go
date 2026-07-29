// ABOUTME: Keeps the boot-resident First Officer shared core below its pre-change byte cap.
package contractlint

import (
	"os"
	"path/filepath"
	"testing"
)

// preChangeSharedCoreBytes is the byte size of first-officer-shared-core.md
// measured immediately before the "boot identifies, engage converges" edits, on
// the post-vcm-merge branch this implementation opened from (`wc -c` == 26755).
const preChangeSharedCoreBytes = 26755

func TestSharedCoreRemainsBelowPreChangeByteCap(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "first-officer", "references", "first-officer-shared-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	if got := len(data); got >= preChangeSharedCoreBytes {
		t.Errorf("shared core is %d bytes, want strictly < %d (pre-change)", got, preChangeSharedCoreBytes)
	}
}
