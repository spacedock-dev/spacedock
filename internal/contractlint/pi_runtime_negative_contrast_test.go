package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPiFirstOfficerRuntimeAvoidsNegativeHostContrast is the Pi FO equivalent of
// TestCodexFirstOfficerRuntimeAvoidsNegativeHostContrast. It guards
// skills/first-officer/references/pi-first-officer-runtime.md against the
// negative host-contrast wording that the runtime-support.md positive-binding
// sweep removed (see the pi-fo-runtime-runtime-support-compliance entity).
//
// A blanket "Claude" substring ban is intentionally NOT used: the Pi FO adapter
// legitimately names "Claude" in a transport instruction (line 16: "without
// rewriting it into Claude syntax"), comparative technical contrast (lines
// 30/36: model-resolution delegation), and a runtime-support.md-conformant
// teardown note (line 66: "Do not emulate Claude team deletion"). This test
// instead bans the specific smell phrases the sweep removed, and asserts the
// positive Pi bindings that replaced them — binding two independent values (the
// runtime-support.md principle and the adapter text), not a prose-grep tautology.
func TestPiFirstOfficerRuntimeAvoidsNegativeHostContrast(t *testing.T) {
	root := repoRoot(t)
	rel := filepath.Join("skills", "first-officer", "references", "pi-first-officer-runtime.md")
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	text := string(data)

	for _, banned := range []string{
		"does not expose Claude Code team-tool signatures",
		"Do not call or ask workers to call Claude team tools",
		"Pi has no such enum",
		"Claude-centric enum",
		"no Claude-centric",
		"Merge-and-Cleanup step 10",
		"Merge-and-Cleanup step",
	} {
		if strings.Contains(text, banned) {
			t.Errorf("%s contains negative host-contrast wording %q", rel, banned)
		}
	}

	for _, want := range []string{
		// FIX 2 — the «worker.shutdown» capability binding replaced the mutable
		// step-number coupling ("Merge-and-Cleanup step 10").
		"`«worker.shutdown»`",
		// FIX 3 — the positive Pi model-space binding replaced the Claude-centric
		// enum contrast.
		"Pi's model-space binding is provider/model strings",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("%s missing positive Pi capability wording %q", rel, want)
		}
	}
}

// TestPiEnsignRuntimeAvoidsNegativeHostContrast is the Pi ensign equivalent of
// TestCodexEnsignRuntimeAvoidsNegativeHostContrast. The Pi ensign adapter is a
// compact binding block keyed by the Pi-specific concerns (clarification,
// completion signal); the shared ensign core owns the rest. The adapter has zero
// legitimate "Claude" mentions, so a targeted phrase ban is safe here: it passes
// today and would catch a re-introduced negative host-contrast sentence. The
// positive assertions pin the binding-block bullets so the host-specific Pi
// completion and clarification bindings survive the shared-core duplication trim.
func TestPiEnsignRuntimeAvoidsNegativeHostContrast(t *testing.T) {
	root := repoRoot(t)
	rel := filepath.Join("skills", "ensign", "references", "pi-ensign-runtime.md")
	data, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	text := string(data)

	for _, banned := range []string{
		"Do not assume Claude team tools",
		"Claude team tools exist in Pi",
	} {
		if strings.Contains(text, banned) {
			t.Errorf("%s contains negative host-contrast wording %q", rel, banned)
		}
	}

	for _, want := range []string{
		// The completion-signal binding bullet — the Pi-specific completion signal
		// the shared core does not carry.
		"- `Completion signal` ->",
		"Completion is reported by the worker's final result in the Pi turn or by the active Pi adapter's task-completion notification.",
		// The clarification binding bullet — the Pi-specific contact_supervisor
		// surface the shared core does not carry.
		"- `Clarification` ->",
		"contact_supervisor` with `reason: \"need_decision\"`",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("%s missing positive Pi ensign binding wording %q", rel, want)
		}
	}
}
