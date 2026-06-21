// ABOUTME: AC-3 — `spacedock install --host codex` with NO plugin installed drives
// ABOUTME: the install seam (marketplace add + plugin add), not prose-only.
package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestInitCodexNotInstalledRunsInstall locks AC-3: `spacedock install --host codex`
// MUST drive the install seam with {codex, marketplaceSource, devBranch} — running
// the marketplace add + plugin add and re-pinning the source — instead of printing
// manual prose and running nothing. The recorded {host, source, branch} call is the
// source of truth that the install fired. No "Run these in your shell" prose
// appears, and the post-install doctor (resolving the now-present plugin) reports OK.
func TestInitCodexNotInstalledRunsInstall(t *testing.T) {
	saved := devBranch
	devBranch = "main"
	defer func() { devBranch = saved }()

	// The plugin is present once the seam runs; the doctor resolve after install
	// sees the compatible manifest.
	fake := &fakeHost{manifest: compatibleManifest(t)}
	var stdout, stderr bytes.Buffer

	code := runInit(context.Background(), []string{"--host", "codex"}, fake, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 (stderr=%q)", code, stderr.String())
	}
	wantInstall := []string{"codex", marketplaceSource, "main"}
	if !equalArgv(fake.installCmds, wantInstall) {
		t.Fatalf("install seam = %v, want %v — codex install must run the seam, not print prose", fake.installCmds, wantInstall)
	}
	out := stdout.String()
	if strings.Contains(out, "Run these in your shell") {
		t.Fatalf("codex install must not fall back to manual prose:\n%s", out)
	}
	// After install, init runs doctor — a compatible report on stdout.
	if !strings.Contains(out, "OK") {
		t.Fatalf("codex install should run doctor after install; stdout = %q", out)
	}
}
