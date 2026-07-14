// ABOUTME: Structural gates for the consolidated bridge seam (mod + hooks + prose).
// ABOUTME: Absence of dropped verbs, reference closure, mod frontmatter, session-id markers.
package contractlint

import (
	"path/filepath"
	"strings"
	"testing"
)

// bridgeSeamProseFiles are the shipped instruction files that carry seam prose.
// The full drain protocol lives in the bridge-seam mod; these speak it too.
var bridgeSeamProseFiles = []string{
	"mods/bridge-seam.md",
	"skills/first-officer/references/first-officer-shared-core.md",
	"skills/first-officer/references/fo-dispatch-core.md",
	"skills/first-officer/references/fo-fleet.md",
	"skills/first-officer/references/fo-bridge.md",
	"skills/first-officer/references/claude-first-officer-runtime.md",
	"skills/first-officer/references/codex-first-officer-runtime.md",
	"skills/first-officer/references/pi-first-officer-runtime.md",
	"skills/present-gate/SKILL.md",
}

// droppedBridgeVerbs are the #445 CLI verbs the consolidation ABSORBED into the
// file protocol. The FO writes the seam files directly, so an instruction that
// tells it to shell one of these is a regression — a producer coupled to a verb
// that no longer exists. (`egress emit`, `inbox check`, and `ingress wake`
// survive as hook/daemon entrypoints and are deliberately NOT listed.)
var droppedBridgeVerbs = []string{
	"bridge inbox drain",
	"bridge inbox ack",
	"bridge inbox commit",
	"bridge alert",
	"bridge initiate",
}

// TestBridgeSeamNoDroppedVerbs is a structural-absence gate: no shipped seam
// instruction may invoke a dropped verb. Keeps the "consumes files, not verbs"
// invariant from drifting back to a CLI coupling.
func TestBridgeSeamNoDroppedVerbs(t *testing.T) {
	for _, rel := range bridgeSeamProseFiles {
		body := readRepoFile(t, filepath.FromSlash(rel))
		for _, verb := range droppedBridgeVerbs {
			if strings.Contains(body, verb) {
				t.Errorf("%s references dropped verb %q — the seam is a direct file write, not a CLI verb", rel, verb)
			}
		}
	}
}

// TestBridgeSeamReferenceClosure is a reference-closure gate: the mod ships at
// both the canonical and dogfood paths, and every deferred reference the seam
// prose points at resolves to a file that exists.
func TestBridgeSeamReferenceClosure(t *testing.T) {
	for _, rel := range []string{
		"mods/bridge-seam.md",
		"docs/dev/_mods/bridge-seam.md",
		"skills/first-officer/references/fo-bridge.md",
		"skills/first-officer/references/fo-fleet.md",
		"docs/dev/bridge-seam.md",
	} {
		if body := readRepoFile(t, filepath.FromSlash(rel)); strings.TrimSpace(body) == "" {
			t.Errorf("required seam file %s is missing or empty", rel)
		}
	}
}

// TestBridgeSeamModFrontmatter is a frontmatter-validity gate: the mod declares a
// name and is NON-standing (a standing:true mod routes through a different, wrong
// execution path — a persistent teammate spawn rather than FO-loop hooks).
func TestBridgeSeamModFrontmatter(t *testing.T) {
	body := readRepoFile(t, filepath.FromSlash("mods/bridge-seam.md"))
	if !strings.HasPrefix(strings.TrimSpace(body), "---") {
		t.Fatal("mods/bridge-seam.md has no frontmatter block")
	}
	fm := body[strings.Index(body, "---")+3:]
	fm = fm[:strings.Index(fm, "---")]
	if !strings.Contains(fm, "name:") {
		t.Error("mods/bridge-seam.md frontmatter missing name:")
	}
	if strings.Contains(fm, "standing:") {
		t.Error("mods/bridge-seam.md must NOT declare standing: — it is an FO-loop hook mod, not a standing teammate")
	}
	for _, h := range []string{"## Hook: startup", "## Hook: idle", "## Agent Prompt"} {
		if !strings.Contains(body, h) {
			t.Errorf("mods/bridge-seam.md missing required section %q", h)
		}
	}
}

// perHostSessionIDMarkers are the harness session-id sources the heartbeat must
// stamp so the Stop-hook inbox check can resolve which slugs belong to a stopping
// session (the load-bearing wake-resolution binding). Portability-marker presence:
// a heartbeat that omits the harness id silently kills Claude intent delivery.
var perHostSessionIDMarkers = []string{
	"CLAUDE_CODE_SESSION_ID", // claude
	"CODEX_THREAD_ID",        // codex
}

// TestBridgeSeamSessionIDBinding is a portability-marker gate: the mod names the
// per-host session-id source, so the wake-resolution binding cannot silently drop
// out of the drain protocol.
func TestBridgeSeamSessionIDBinding(t *testing.T) {
	body := readRepoFile(t, filepath.FromSlash("mods/bridge-seam.md"))
	for _, marker := range perHostSessionIDMarkers {
		if !strings.Contains(body, marker) {
			t.Errorf("mods/bridge-seam.md missing per-host session-id marker %q — the heartbeat must carry the harness session id or the Stop-hook wake resolves nothing", marker)
		}
	}
}
