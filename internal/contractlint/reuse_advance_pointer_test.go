// ABOUTME: AC-3 structural guard — the reuse-advance path routes through
// ABOUTME: `dispatch build --advance`, not a hand-assembled verbatim-section message.
package contractlint

import (
	"path/filepath"
	"strings"
	"testing"
)

var (
	foDispatchCoreRel   = filepath.Join("skills", "first-officer", "references", "fo-dispatch-core.md")
	claudeFODispatchRel = filepath.Join("skills", "first-officer", "references", "claude-fo-dispatch.md")
	ensignSharedCoreRel = filepath.Join("skills", "ensign", "references", "ensign-shared-core.md")
)

// TestFODispatchCoreReuseNamesAdvanceHelper is AC-3's first anchor:
// fo-dispatch-core.md's "If reuse" step routes through `«dispatch.build»`
// `--advance` mode, and the retired claim that the reuse path bypasses the
// helper entirely is gone (structural absence — a stale claim left in place
// would contradict the routing this task ships).
func TestFODispatchCoreReuseNamesAdvanceHelper(t *testing.T) {
	text := readRepoFile(t, foDispatchCoreRel)

	section := markdownSectionFromText(t, text, "## Reuse and Fresh Dispatch")
	if !strings.Contains(section, "--advance") {
		t.Errorf("%s `## Reuse and Fresh Dispatch` does not name --advance — the reuse path must route through `dispatch build --advance`", foDispatchCoreRel)
	}
	if !strings.Contains(section, "«dispatch.build»") {
		t.Errorf("%s `## Reuse and Fresh Dispatch` does not name «dispatch.build» in its \"If reuse\" step", foDispatchCoreRel)
	}

	retired := "does NOT route through `«dispatch.build»` — assemble the advancement message directly"
	if strings.Contains(text, retired) {
		t.Errorf("%s still carries the retired claim %q — the reuse path routes through the helper now", foDispatchCoreRel, retired)
	}
	retiredClosingLine := "the reuse-advance path assembles its message directly"
	if strings.Contains(text, retiredClosingLine) {
		t.Errorf("%s still carries the retired closing-line claim %q", foDispatchCoreRel, retiredClosingLine)
	}
}

// TestClaudeFODispatchReuseAdvanceRoutesThroughHelper is AC-3's second anchor:
// claude-fo-dispatch.md's reuse-advance handle is `SendMessage(to={handle},
// message=output.prompt)` — the helper's emitted pointer, forwarded verbatim —
// with the old hand-assembled verbatim-stage-section template demoted to a
// break-glass fallback (retained, not deleted, but no longer the live path).
func TestClaudeFODispatchReuseAdvanceRoutesThroughHelper(t *testing.T) {
	text := readRepoFile(t, claudeFODispatchRel)

	if !strings.Contains(text, `SendMessage(to="{live worker handle from session roster}", message=output.prompt)`) {
		t.Errorf("%s missing the reuse-advance handle's live call shape: SendMessage(to={handle}, message=output.prompt)", claudeFODispatchRel)
	}

	breakGlassHeading := "**Break-glass reuse-advance"
	breakGlassIdx := strings.Index(text, breakGlassHeading)
	if breakGlassIdx < 0 {
		t.Fatalf("%s missing the %q break-glass heading for the reuse-advance fallback", claudeFODispatchRel, breakGlassHeading)
	}

	verbatimTemplateMarker := "[STAGE_DEFINITION — copy the full ### stage subsection from the README verbatim]"
	verbatimIdx := strings.Index(text, verbatimTemplateMarker)
	if verbatimIdx < 0 {
		t.Errorf("%s no longer carries the verbatim-section template at all — it must be RETAINED as the break-glass fallback, not deleted", claudeFODispatchRel)
	}
	if verbatimIdx >= 0 && verbatimIdx < breakGlassIdx {
		t.Errorf("%s carries the verbatim-section template BEFORE the break-glass heading — it must live under break-glass, not the live reuse-advance path", claudeFODispatchRel)
	}

	// The reuse-advance handle paragraph (before break-glass) must be the one
	// naming the live call shape, not the verbatim template.
	liveSection := text[:breakGlassIdx]
	if !strings.Contains(liveSection, "message=output.prompt") {
		t.Errorf("%s live reuse-advance handle section (before break-glass) does not route through output.prompt", claudeFODispatchRel)
	}
	if strings.Contains(liveSection, verbatimTemplateMarker) {
		t.Errorf("%s live reuse-advance handle section (before break-glass) still inlines the verbatim-section template", claudeFODispatchRel)
	}
}

// TestEnsignSharedCoreCarriesAdvanceBootstrapClause is AC-3's third anchor:
// ensign-shared-core.md's DISPATCH_FILE Bootstrap section gains the mid-session
// advance clause (Read the pointed file on an "Advancing to next stage:"
// message) with the same DISPATCH_FILE_MISSING failure shape as the initial
// bootstrap.
func TestEnsignSharedCoreCarriesAdvanceBootstrapClause(t *testing.T) {
	text := readRepoFile(t, ensignSharedCoreRel)

	section := markdownSectionFromText(t, text, "## DISPATCH_FILE Bootstrap")
	if !strings.Contains(section, "Advancing to next stage:") {
		t.Errorf("%s `## DISPATCH_FILE Bootstrap` does not name the mid-session advance pointer shape (\"Advancing to next stage:\")", ensignSharedCoreRel)
	}
	if !strings.Contains(section, "next-stage assignment") {
		t.Errorf("%s `## DISPATCH_FILE Bootstrap` does not name the next-stage assignment the advance pointer resolves to", ensignSharedCoreRel)
	}
	if strings.Count(section, "DISPATCH_FILE_MISSING") < 2 {
		t.Errorf("%s `## DISPATCH_FILE Bootstrap` must carry the DISPATCH_FILE_MISSING failure shape for BOTH the initial and advance bootstrap paths", ensignSharedCoreRel)
	}
}
