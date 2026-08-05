//go:build live

package ensigncycle

import (
	"os/exec"
	"strings"
	"testing"
)

func skipUnlessMergedHost(t *testing.T) {
	t.Helper()
	out, err := exec.Command("claude", "--version").Output()
	if err != nil {
		t.Logf("merged Claude proof: `claude --version` failed (%v); proceeding so the selected proof fails if the host is unsupported", err)
		return
	}
	merged, parsed := mergedClaudeHost(string(out))
	if parsed && !merged {
		t.Skipf("merged Claude proof requires Claude Code 2.1.%d or newer; installed version is %q", mergedFloorPatch, strings.TrimSpace(string(out)))
	}
}
