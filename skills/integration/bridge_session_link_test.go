// ABOUTME: Bridge session-link smoke — runs the EXACT shell the Claude ensign
// ABOUTME: runtime documents (extracted from the doc, no drift) and asserts it
// ABOUTME: writes the _bridge/sessions/<session_id>.json map Bridge joins against
// ABOUTME: events.jsonl to show the ensign's entity as RUNNING in real time.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// firstFencedBlock returns the first ``` ... ``` block that follows the given
// heading line in a markdown file. Running the doc's own snippet (rather than a
// copy) keeps this test honest: if the documented shell is removed or its contract
// changes, the extraction fails or the assertions break.
func firstFencedBlock(t *testing.T, path, heading string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	lines := strings.Split(string(data), "\n")
	i := 0
	for ; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == heading {
			break
		}
	}
	if i == len(lines) {
		t.Fatalf("heading %q not found in %s", heading, path)
	}
	// Find the opening fence after the heading.
	for ; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
			break
		}
	}
	if i == len(lines) {
		t.Fatalf("no fenced block after %q in %s", heading, path)
	}
	var body []string
	for j := i + 1; j < len(lines); j++ {
		if strings.HasPrefix(strings.TrimSpace(lines[j]), "```") {
			return strings.Join(body, "\n")
		}
		body = append(body, lines[j])
	}
	t.Fatalf("unterminated fenced block after %q in %s", heading, path)
	return ""
}

// gitInitBare returns a fresh temp dir initialized as a git repo (no commit needed)
// so the documented snippet's `git rev-parse --show-toplevel` resolves to it.
func gitInitBare(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	return dir
}

// TestEnsignBridgeSessionLink runs the documented session-link shell and proves it
// writes the session→entity map keyed by CLAUDE_CODE_SESSION_ID at the repo root —
// the producer half of Bridge's live "running" badge.
func TestEnsignBridgeSessionLink(t *testing.T) {
	doc := filepath.Join("..", "ensign", "references", "claude-ensign-runtime.md")
	snippet := firstFencedBlock(t, doc, "## Bridge Session Link")
	// Substitute the placeholders an ensign fills from its assignment.
	snippet = strings.ReplaceAll(snippet, "ENTITY_SLUG", "drc-3339")
	snippet = strings.ReplaceAll(snippet, "STAGE_NAME", "review")

	root := gitInitBare(t) // so `git rev-parse --show-toplevel` resolves to root
	const sid = "ses-abc-123"

	cmd := exec.Command("bash", "-c", snippet)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "CLAUDE_CODE_SESSION_ID="+sid)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("documented session-link shell failed: %v\n%s", err, out)
	}

	marker := filepath.Join(root, "_bridge", "sessions", sid+".json")
	data, err := os.ReadFile(marker)
	if err != nil {
		t.Fatalf("session→entity marker not written at %s: %v", marker, err)
	}
	got := string(data)
	for _, want := range []string{`"session_id":"ses-abc-123"`, `"entity":"drc-3339"`, `"stage":"review"`} {
		if !strings.Contains(got, want) {
			t.Errorf("marker missing %s\ngot: %s", want, got)
		}
	}

	// Safety contract: an unset session id must be a clean no-op (no stray file).
	root2 := gitInitBare(t)
	cmd2 := exec.Command("bash", "-c", snippet)
	cmd2.Dir = root2
	cmd2.Env = append(os.Environ(), "CLAUDE_CODE_SESSION_ID=")
	if out, err := cmd2.CombinedOutput(); err != nil {
		t.Fatalf("snippet must no-op (exit 0) on unset id, got: %v\n%s", err, out)
	}
	if entries, _ := os.ReadDir(filepath.Join(root2, "_bridge", "sessions")); len(entries) != 0 {
		t.Errorf("unset session id should write nothing, found %d files", len(entries))
	}
}
