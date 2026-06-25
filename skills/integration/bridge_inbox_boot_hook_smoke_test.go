// ABOUTME: Bridge-inbox boot smoke — the binary precondition the FO's before-greet
// ABOUTME: liveness heartbeat (shared core Startup step 7b) relies on: `status --boot`
// ABOUTME: MODS-REPORT surfaces a registered `startup` bridge-inbox hook so the FO
// ABOUTME: knows to write `_bridge/fo.$SLUG.json` before the greet, even greet-and-stop.
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// bridgeInboxMod is a minimal mod carrying a `## Hook: startup` heading — enough for
// the boot MODS-REPORT to register it under the `startup` lifecycle point. The real
// mod's heartbeat/drain body is exercised by the FO at runtime; this fixture locks
// only that the boot surface ADVERTISES the startup hook to the FO.
const bridgeInboxMod = `---
name: bridge-inbox
description: liveness heartbeat + captain-intent drain
---

## Hook: startup

Write the heartbeat, then drain.

## Hook: idle

Refresh the heartbeat, then drain.
`

// TestBootReportsBridgeInboxStartupHook locks the precondition for shared core
// Startup step 7b: a workflow carrying `_mods/bridge-inbox.md` with a `## Hook:
// startup` must show that hook under MODS in `status --boot`, so the FO runs the
// before-greet liveness heartbeat. If the boot stopped surfacing the startup hook,
// the FO could not know to write `_bridge/fo.$SLUG.json` before the greet, and a
// greet-and-stop launch would read "no FO attached" in Bridge.
func TestBootReportsBridgeInboxStartupHook(t *testing.T) {
	root := t.TempDir()
	defDir := filepath.Join(root, "wf-hb")
	if err := os.MkdirAll(filepath.Join(defDir, "_mods"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(defDir, "README.md"), []byte(fleetMemberReadme), 0o644); err != nil {
		t.Fatal(err)
	}
	entity := "---\nid: \"\"\ntitle: HB one\nstatus: backlog\nscore: \"0.50\"\nsource: smoke\n---\n# HB one\n\nSeed entity.\n"
	if err := os.WriteFile(filepath.Join(defDir, "wf-hb-1.md"), []byte(entity), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(defDir, "_mods", "bridge-inbox.md"), []byte(bridgeInboxMod), 0o644); err != nil {
		t.Fatal(err)
	}
	gitInitFixture(t, root)

	out, code := runStatus(t, defDir, "--boot")
	if code != 0 {
		t.Fatalf("status --boot exit %d:\n%s", code, out)
	}

	// MODS-REPORT must carry a `startup:` line that names bridge-inbox.
	var startupLine string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "startup:") {
			startupLine = line
			break
		}
	}
	if startupLine == "" {
		t.Fatalf("status --boot MODS-REPORT has no `startup:` hook line — the FO cannot learn the bridge-inbox startup hook is registered:\n%s", out)
	}
	if !strings.Contains(startupLine, "bridge-inbox") {
		t.Fatalf("status --boot `startup:` line %q does not name bridge-inbox — the before-greet liveness heartbeat (step 7b) precondition is broken:\n%s", startupLine, out)
	}
}
