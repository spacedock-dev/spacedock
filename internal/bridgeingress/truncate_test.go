package bridgeingress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// writeReplyFixture writes one fo-replies.jsonl line per element of inReplyToLines,
// each a minimal record carrying that in_reply_to_line (the only field truncateReplies
// reads). Returns the _bridge dir.
func writeReplyFixture(t *testing.T, root string, inReplyToLines []int) string {
	t.Helper()
	dir := filepath.Join(root, "_bridge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir _bridge: %v", err)
	}
	var b strings.Builder
	for _, n := range inReplyToLines {
		fmt.Fprintf(&b, `{"in_reply_to_line":%d}`+"\n", n)
	}
	if err := os.WriteFile(filepath.Join(dir, "fo-replies.jsonl"), []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write fo-replies.jsonl: %v", err)
	}
	return dir
}

// writeCursorFile writes a per-slug committed cursor (the same "<n>\n" shape Commit
// writes), so truncateReplies derives a retention floor from it.
func writeCursorFile(t *testing.T, dir, slug string, cursor int) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".inbox-cursor."+slug), []byte(strconv.Itoa(cursor)+"\n"), 0o644); err != nil {
		t.Fatalf("write cursor %s: %v", slug, err)
	}
}

// readReplyInReplyToLines reads fo-replies.jsonl back as the ordered slice of
// in_reply_to_line values it carries.
func readReplyInReplyToLines(t *testing.T, dir string) []int {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "fo-replies.jsonl"))
	if err != nil {
		t.Fatalf("read fo-replies.jsonl: %v", err)
	}
	trimmed := strings.TrimRight(string(data), "\n")
	if trimmed == "" {
		return nil
	}
	var out []int
	for _, line := range strings.Split(trimmed, "\n") {
		var v struct {
			N int `json:"in_reply_to_line"`
		}
		if err := json.Unmarshal([]byte(line), &v); err != nil {
			t.Fatalf("kept line does not parse: %q: %v", line, err)
		}
		out = append(out, v.N)
	}
	return out
}

func countValue(vals []int, target int) int {
	n := 0
	for _, v := range vals {
		if v == target {
			n++
		}
	}
	return n
}

// TestTruncateRepliesNoOpWithoutCursor: with no committed per-slug cursor the floor
// is unknown, so nothing is dropped even well past the size cap (dropping blind would
// risk re-draining an already-answered intent).
func TestTruncateRepliesNoOpWithoutCursor(t *testing.T) {
	root := t.TempDir()
	n := maxReplyLines + 100
	lines := make([]int, n)
	for i := range lines {
		lines[i] = 1
	}
	dir := writeReplyFixture(t, root, lines)

	truncateReplies(dir)

	if got := len(readReplyInReplyToLines(t, dir)); got != n {
		t.Errorf("no cursor present: want file untouched (%d lines), got %d", n, got)
	}
}

// TestTruncateRepliesNoOpUnderCap: at or below the size cap the file is left intact
// even when every line sits below the floor and outside the recency window.
func TestTruncateRepliesNoOpUnderCap(t *testing.T) {
	root := t.TempDir()
	n := maxReplyLines // exactly the cap → the > cap guard is false
	lines := make([]int, n)
	for i := range lines {
		lines[i] = 1
	}
	dir := writeReplyFixture(t, root, lines)
	writeCursorFile(t, dir, "alpha", 999999)

	truncateReplies(dir)

	if got := len(readReplyInReplyToLines(t, dir)); got != n {
		t.Errorf("under cap: want file untouched (%d lines), got %d", n, got)
	}
}

// TestTruncateRepliesDropsBelowFloorOutsideWindow is the core safety rule: over the
// cap, a reply is dropped only when it is BOTH below the retention floor AND outside
// the recency window. A below-floor line inside the window is kept; an at-or-above-floor
// line outside the window is kept (dedup preservation); a below-floor line outside the
// window is the only thing dropped. Original order is preserved.
func TestTruncateRepliesDropsBelowFloorOutsideWindow(t *testing.T) {
	root := t.TempDir()
	n := maxReplyLines + 100 // 2100
	windowStart := n - keepReplyLines
	const floor = 500
	const keepMarker = 9999 // >= floor, placed outside the window: must survive
	const dropMarker = 10   // < floor, outside the window: must be dropped
	const windowMarker = 1  // < floor, inside the window: kept by recency

	lines := make([]int, n)
	for i := range lines {
		switch {
		case i >= windowStart:
			lines[i] = windowMarker
		case i == 0:
			lines[i] = keepMarker
		default:
			lines[i] = dropMarker
		}
	}
	dir := writeReplyFixture(t, root, lines)
	writeCursorFile(t, dir, "alpha", floor)

	truncateReplies(dir)

	kept := readReplyInReplyToLines(t, dir)
	if got, want := len(kept), 1+keepReplyLines; got != want {
		t.Fatalf("kept %d lines, want %d (the one at/above-floor outside line + the %d-line window)", got, want, keepReplyLines)
	}
	if kept[0] != keepMarker {
		t.Errorf("first kept line = %d, want the at/above-floor line %d preserved in original order", kept[0], keepMarker)
	}
	if c := countValue(kept, dropMarker); c != 0 {
		t.Errorf("below-floor outside-window lines survived: %d remain, want 0", c)
	}
	if c := countValue(kept, windowMarker); c != keepReplyLines {
		t.Errorf("recency window not fully retained: %d window lines, want %d", c, keepReplyLines)
	}
}

// TestTruncateRepliesKeepsOnlyWindowWhenAllBelowFloor: over the cap with every line
// below the floor, exactly the recency window survives.
func TestTruncateRepliesKeepsOnlyWindowWhenAllBelowFloor(t *testing.T) {
	root := t.TempDir()
	n := maxReplyLines + 100
	lines := make([]int, n)
	for i := range lines {
		lines[i] = 1
	}
	dir := writeReplyFixture(t, root, lines)
	writeCursorFile(t, dir, "alpha", 999999)

	truncateReplies(dir)

	if got := len(readReplyInReplyToLines(t, dir)); got != keepReplyLines {
		t.Errorf("all below floor: want exactly the %d-line window, got %d", keepReplyLines, got)
	}
}

// TestTruncateRepliesFloorIsLowestAcrossSlugs: the floor is the MINIMUM committed
// cursor across every per-slug file, so a reply at/above the lowest slug's cursor is
// retained even though it is below another slug's cursor.
func TestTruncateRepliesFloorIsLowestAcrossSlugs(t *testing.T) {
	root := t.TempDir()
	n := maxReplyLines + 100
	windowStart := n - keepReplyLines
	const lowFloor = 50
	const atLowFloor = 100 // >= 50 (lowest) but well below 500: retained via the min floor
	const belowAll = 10    // < 50: dropped

	lines := make([]int, n)
	for i := range lines {
		switch {
		case i >= windowStart:
			lines[i] = 1
		case i == 0:
			lines[i] = atLowFloor
		default:
			lines[i] = belowAll
		}
	}
	dir := writeReplyFixture(t, root, lines)
	writeCursorFile(t, dir, "alpha", 500)
	writeCursorFile(t, dir, "beta", lowFloor)

	truncateReplies(dir)

	kept := readReplyInReplyToLines(t, dir)
	if countValue(kept, atLowFloor) != 1 {
		t.Errorf("line at the lowest slug's floor was dropped — floor is not the min across slugs")
	}
	if countValue(kept, belowAll) != 0 {
		t.Errorf("line below every floor survived, want 0")
	}
}
