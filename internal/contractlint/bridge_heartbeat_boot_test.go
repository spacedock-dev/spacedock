// ABOUTME: Contract anchor — the bridge-inbox liveness heartbeat is a BEFORE-GREET
// ABOUTME: boot step, so a live FO shows attached in Bridge from boot even in a
// ABOUTME: greet-and-stop launch (the regression that read "no FO attached").
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBridgeHeartbeatRunsBeforeGreet locks the boot-time liveness contract that the
// `bridge-inbox` mod's `## Hook: startup` (the per-`$SLUG` `_bridge/fo.$SLUG.json`
// heartbeat + initial drain) is run BEFORE the greet — not deferred to first dispatch
// like the comm-officer spawn. Without this, a greet-and-stop boot (interactive
// step 8, which never enters the event loop) writes no heartbeat, so Bridge reads
// "no FO attached" though a live FO exists. The shipped behavior that produced this
// test: the FO greeted, parked at a gate, and never wrote a heartbeat all session.
func TestBridgeHeartbeatRunsBeforeGreet(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "first-officer", "references", "first-officer-shared-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	content := string(data)

	// The before-greet boot step must exist, name the canonical mod and the per-slug
	// heartbeat file, and state explicitly that a greet-and-stop boot still runs it.
	// Named by mod NAME (not a concrete `_mods/bridge-inbox.md` path): only pr-merge
	// has a canonical `<root>/mods/` copy, so naming a concrete bridge-inbox path
	// would dead-end the boot-resident closure check. The MODS map gives the FO the
	// mod name; it resolves `{workflow_dir}/_mods/<name>.md` at runtime.
	required := []string{
		"Bridge liveness heartbeat (before-greet)",
		"startup` hook for `bridge-inbox",
		"_bridge/fo.$SLUG.json",
		"greet-and-stop boot",
	}
	for _, r := range required {
		if !strings.Contains(content, r) {
			t.Errorf("first-officer-shared-core.md no longer anchors the before-greet bridge-inbox heartbeat: missing %q.\n"+
				"Without a before-greet heartbeat, a greet-and-stop FO never attaches in Bridge.", r)
		}
	}

	// Guard the specific deferral carve-out: the MODS note must NOT claim ALL startup
	// hooks defer, or the FO will again defer the heartbeat to first dispatch.
	if strings.Contains(content, "Startup hooks run deferred:") &&
		!strings.Contains(content, "Startup hooks run deferred EXCEPT") {
		t.Errorf("the MODS startup-hook note reverted to an unconditional 'Startup hooks run deferred:' — " +
			"the bridge-inbox heartbeat must be carved out as before-greet, or greet-and-stop boots show no FO attached.")
	}
}
