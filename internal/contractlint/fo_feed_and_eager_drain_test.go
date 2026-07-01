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

// TestBridgeConversationReplyContract locks the full Bridge conversation loop:
// Bridge writes stable intent ids and a frozen target_set, and the FO writes a
// post-action reply/ack line to _bridge/fo-replies.jsonl.
func TestBridgeConversationReplyContract(t *testing.T) {
	modPath := filepath.Join(repoRoot(t), "docs", "dev", "_mods", "bridge-inbox.md")
	modData, err := os.ReadFile(modPath)
	if err != nil {
		t.Fatalf("read bridge-inbox mod: %v", err)
	}
	mod := string(modData)
	for _, r := range []string{
		`"id":"<opaque unique string>"`,
		`"target_set":["<workflow-slug>", "..."]`,
		`If ` + "`target_set`" + ` is present, act only when ` + "`\"$SLUG\"`" + ` is in ` + "`target_set`",
		`Ignore ` + "`target`" + ` entirely for routing in that case, including ` + "`target == \"all\"`",
		`If ` + "`target_set`" + ` is absent, preserve old target behavior`,
		`physical ` + "`LINE`" + ` number`,
		`not merely after shell-reading the line`,
		`_bridge/fo-replies.jsonl`,
		`"kind":"reply"|"conn-ack"|"decision-ack"`,
		`"in_reply_to_line":123`,
		`"status":"answered"|"accepted"|"released"|"applied"|"rejected"|"blocked"`,
		`target` + "` is the actual acknowledging workflow slug (`$SLUG`), never " + "`\"all\"`",
		`Cursor remains the delivery/read source of truth`,
		`Bridge folds them by intent id (or legacy line fallback), acknowledging target, and reply kind`,
		`one complete newline-terminated JSON object in one append operation`,
		`Do not rewrite, truncate, sort, or compact ` + "`fo-replies.jsonl`",
		`Append a rejected ack only when the record is addressed to ` + "`\"$SLUG\"`" + ` and has enough valid metadata to produce a valid reply shape`,
		`Unknown-kind or unrouteable records cannot be represented by the reply schema`,
	} {
		if !strings.Contains(mod, r) {
			t.Errorf("bridge-inbox mod no longer pins the reply-loop contract: missing %q.", r)
		}
	}

	contractPath := filepath.Join(repoRoot(t), "docs", "dev", "bridge-egress-contract.md")
	contractData, err := os.ReadFile(contractPath)
	if err != nil {
		t.Fatalf("read bridge egress contract: %v", err)
	}
	contract := string(contractData)
	for _, r := range []string{
		`Bridge writes every new intent with an opaque ` + "`id`",
		`frozen ` + "`target_set`" + ` array of workflow slugs`,
		`When ` + "`target_set`" + ` is present it is authoritative`,
		`## ` + "`_bridge/fo-replies.jsonl`" + ` — captain-intent acknowledgements`,
		`"schema":1`,
		`"target":"<acknowledging workflow slug>"`,
		`"intent_kind":"<tell|conn|decision>"`,
		`"status":"<answered|accepted|released|applied|rejected|blocked>"`,
		`target` + "` is the actual acknowledging workflow slug, never " + "`all`",
		`applied` + "` when a decision field value is present and gate resolution finished or was already satisfied",
		`blocked` + "` when a valid intent could not finish",
		`rejected` + "` when an intent is invalid or unresolvable",
		`not an exactly-once delivery ledger`,
		`Append-only and best-effort: write one complete newline-terminated JSON object in one append operation`,
	} {
		if !strings.Contains(contract, r) {
			t.Errorf("bridge egress contract no longer pins the reply-loop contract: missing %q.", r)
		}
	}
}

// TestFleetModeMentionsFrozenTargetSet guards the high-level FO routing guide:
// fleet-mode drains must obey Bridge's frozen target_set before falling back to
// legacy target behavior.
func TestFleetModeMentionsFrozenTargetSet(t *testing.T) {
	path := filepath.Join(repoRoot(t), "skills", "first-officer", "references", "first-officer-shared-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read first-officer-shared-core: %v", err)
	}
	c := string(data)
	for _, r := range []string{
		`frozen ` + "`target_set`",
		`act only when the member's ` + "`$SLUG`" + ` is in ` + "`target_set`",
		`_bridge/fo-replies.jsonl`,
		`with ` + "`target`" + ` set to the actual member slug, never ` + "`all`",
		`Legacy records without ` + "`target_set`" + ` keep old ` + "`target`" + ` behavior`,
		`Claude, Codex, and Pi have packaged event producers`,
		`deterministic session→entity marker parity remains Claude-proven only`,
	} {
		if !strings.Contains(c, r) {
			t.Errorf("fleet-mode routing guide no longer pins target_set/ack behavior: missing %q.", r)
		}
	}
}
