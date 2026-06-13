// ABOUTME: AC-1 — `status --next-id` emits a use-`new` hint on stderr without
// ABOUTME: altering the stdout id, and the --next-id --json path stays hint-free.
package status

import (
	"strings"
	"testing"
)

// TestNextIDPlainTextEmitsNewHintOnStderr pins that the plain-text --next-id
// path writes ONLY the computed id to stdout (unchanged from the parity bytes)
// while pointing the operator at the atomic `spacedock new` path on stderr. The
// expectation source is independent of native_runner.go: the id is what
// computeNextID returns (observed via stdout), not the hint string read back
// out of the runner, so the assertion cannot be satisfied by echoing prose.
func TestNextIDPlainTextEmitsNewHintOnStderr(t *testing.T) {
	root := indSeqWorkflow(t)
	env := indEnv(t)

	out, errOut, code := indRunNative(t, root, env, "", "--workflow-dir", root, "--next-id")
	if code != 0 {
		t.Fatalf("--next-id exit=%d stderr=%q", code, errOut)
	}

	// stdout is exactly the id and nothing else (single trailing newline). The
	// id is computeNextID's output; the seq fixture's next sequential id is 006.
	if out != "006\n" {
		t.Fatalf("--next-id stdout must be exactly the id, got %q", out)
	}

	// stderr carries a hint that points at `spacedock new` as the atomic path.
	if !strings.Contains(errOut, "spacedock new") {
		t.Fatalf("--next-id stderr must hint at `spacedock new`, got %q", errOut)
	}

	// The hint must not leak the id into stderr (the id belongs on stdout only).
	if strings.Contains(errOut, strings.TrimSpace(out)) {
		t.Fatalf("--next-id hint must not echo the id %q onto stderr: %q", strings.TrimSpace(out), errOut)
	}
}

// TestNextIDJSONStaysHintFree pins that the machine-readable --next-id --json
// path is unchanged: stdout is the single-key {"command":"next-id","id":...}
// object and stderr carries NO hint, so a JSON consumer sees no new noise.
func TestNextIDJSONStaysHintFree(t *testing.T) {
	root := indSeqWorkflow(t)
	env := indEnv(t)

	out, errOut, code := indRunNative(t, root, env, "", "--workflow-dir", root, "--next-id", "--json")
	if code != 0 {
		t.Fatalf("--next-id --json exit=%d stderr=%q", code, errOut)
	}
	if got := strings.TrimSpace(out); got != `{"command":"next-id","id":"006"}` {
		t.Fatalf("--next-id --json stdout changed, got %q", got)
	}
	if strings.Contains(errOut, "spacedock new") {
		t.Fatalf("--next-id --json must stay hint-free on stderr, got %q", errOut)
	}
	if errOut != "" {
		t.Fatalf("--next-id --json stderr must be empty, got %q", errOut)
	}
}
