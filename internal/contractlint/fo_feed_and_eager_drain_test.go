// ABOUTME: Contract anchors for the two Bridge-seam steering/visibility fixes:
// ABOUTME: (Issue 2) the FO drains captain intent EAGERLY at the top of each loop
// ABOUTME: iteration so a mid-drive `pause` is honored; (Issue 1) the bridge-inbox
// ABOUTME: mod writes a `_bridge/fo-feed.jsonl` narration the fleet-history reads.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEagerCaptainIntentDrain locks Issue 2: the event loop drains the captain
// inbox at the TOP of each iteration, not only at the idle boundary — a Driving FO
// never reaches idle, so an idle-only drain leaves a queued `pause` unread.
func TestEagerCaptainIntentDrain(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "first-officer", "references", "fo-dispatch-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fo-dispatch-core: %v", err)
	}
	c := string(data)
	for _, r := range []string{
		"Drain captain intent + refresh liveness (eager", // the new pre-dispatch step
		"heartbeat refresh",      // RC: a busy FO must not go stale/not-attached
		"never reaches idle",     // why idle-only is insufficient
		"halts further dispatch", // a pause takes effect this iteration
	} {
		if !strings.Contains(c, r) {
			t.Errorf("fo-dispatch-core.md no longer drains captain intent eagerly (Issue 2): missing %q.\n"+
				"A Driving FO would never reach the idle boundary, so a queued pause/redirect sits unread.", r)
		}
	}
}

// TestFOFeedNarration locks Issue 1: the bridge-inbox mod writes a session
// narration to _bridge/fo-feed.jsonl on dispatch/advance/complete, so the
// fleet-history rail shows activity for a local-only workflow that commits no
// dispatch:/advance: git narration.
func TestFOFeedNarration(t *testing.T) {
	path := filepath.Join(repoRoot(t), "docs", "dev", "_mods", "bridge-inbox.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read bridge-inbox mod: %v", err)
	}
	c := string(data)
	for _, r := range []string{
		"_bridge/fo-feed.jsonl", // the feed file Bridge reads
		"dispatch",              // the verbs narrated
		"local-only workflow",   // the case it exists for
	} {
		if !strings.Contains(c, r) {
			t.Errorf("bridge-inbox mod no longer writes the FO feed (Issue 1): missing %q.", r)
		}
	}
}
