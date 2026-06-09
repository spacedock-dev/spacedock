package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstallJourneyDocumentsLinuxPath locks AC-4's positive half: the install
// guide carries a runnable `curl … install.sh | sh` invocation so a Linux user
// has a concrete install path (the cask is darwin-only). The check requires the
// raw.githubusercontent install.sh URL piped to a shell — the exact line a user
// copy-pastes — not merely a mention of "install.sh".
func TestInstallJourneyDocumentsLinuxPath(t *testing.T) {
	doc := readInstallJourney(t)
	commands := executableShellCommands(doc)

	found := false
	for _, c := range commands {
		if strings.Contains(c, "install.sh") &&
			strings.Contains(c, "curl ") &&
			strings.Contains(c, "raw.githubusercontent.com/spacedock-dev/spacedock") &&
			(strings.Contains(c, "| sh") || strings.HasSuffix(c, "sh")) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("docs/site/get-started/install.md has no runnable `curl … install.sh | sh` Linux install line")
	}
}

// TestInstallJourneyDoesNotOverclaimLinuxSandbox locks AC-4's honesty half:
// spacedock ships no sandbox — internal/safehouse only detects a profile, checks
// for a `safehouse` binary on PATH, and wraps argv. So the doc must NOT assert an
// unqualified "sandboxed on Linux" claim. We scan for sandbox-claim phrasings
// near "Linux" and reject any that promise sandboxing without the
// requires-a-safehouse-binary qualifier.
func TestInstallJourneyDoesNotOverclaimLinuxSandbox(t *testing.T) {
	// Strip markdown code backticks so `safehouse` binary reads as the plain
	// phrase the qualifier check looks for.
	doc := strings.ToLower(strings.ReplaceAll(readInstallJourney(t), "`", ""))

	// Phrasings that would over-claim a Linux sandbox if stated unqualified.
	for _, claim := range []string{
		"sandboxed on linux",
		"runs sandboxed on linux",
		"sandboxing on linux works",
	} {
		if strings.Contains(doc, claim) {
			t.Errorf("install-journey over-claims a Linux sandbox with %q; spacedock ships no sandbox — sandboxing requires a Linux-capable safehouse binary", claim)
		}
	}

	// If the doc discusses safehouse on Linux at all, it must name the binary
	// dependency so the claim is qualified, not an unconditional promise.
	if strings.Contains(doc, "safehouse") && strings.Contains(doc, "linux") {
		if !strings.Contains(doc, "safehouse binary") {
			t.Errorf("install-journey mentions safehouse + Linux but never names the required `safehouse binary` — the qualifier that keeps the claim honest")
		}
	}
}

func readInstallJourney(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "docs", "site", "get-started", "install.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
