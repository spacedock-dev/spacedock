// ABOUTME: AC-1/AC-2/AC-3 section-scoped presence oracles — the FO contract surfaces
// ABOUTME: forbid a completion-only stop and mandate same-turn advance/dispatch after a non-gated stage.
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// piFirstOfficerRuntime reads the vendored Pi first-officer runtime adapter text.
func piFirstOfficerRuntime(t *testing.T) string {
	t.Helper()
	p := filepath.Join(skillsRoot(t), "first-officer", "references", "pi-first-officer-runtime.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read Pi first-officer runtime: %v", err)
	}
	return string(b)
}

// devWorkflowReadme reads the dev workflow guide the captain reads.
func devWorkflowReadme(t *testing.T) string {
	t.Helper()
	p := filepath.Join(repoRoot(t), "docs", "dev", "README.md")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read dev workflow README: %v", err)
	}
	return string(b)
}

// TestSharedCoreForbidsCompletionOnlyStop locks AC-1: the FO shared contract's
// `## Completion and Gates` section forbids a completion-only stop after a
// non-gated, non-terminal stage and requires advancing/dispatching the next
// stage before ending the turn, naming the halt exceptions (gate / terminal /
// blocker / captain decision).
//
// Oracle: a presence test scoped to the `## Completion and Gates` section via
// sectionAfter, so the clause cannot be satisfied by unrelated prose elsewhere.
// The anchors are the obligation phrasing — "is not a stopping point", the MUST
// advance-before-ending-the-turn duty, and each of the four named halt spans —
// not loose single words a meaning-inverting paraphrase could keep. Deleting or
// gutting the clause reds the test. Proof at the claim's own level: the claim IS
// the contract text. The behavioral half is AC-5's live regression.
func TestSharedCoreForbidsCompletionOnlyStop(t *testing.T) {
	section := sectionAfter(foSharedCore(t), "## Completion and Gates")
	if section == "" {
		t.Fatal("FO shared core missing the `## Completion and Gates` section")
	}

	// The not-a-stopping-point obligation: a completed non-gated, non-terminal
	// stage is not a stopping point, and the FO MUST advance + dispatch before
	// ending its turn (its own next action, not the captain's).
	for _, anchor := range []string{
		"non-gated, non-terminal stage is not a stopping point",
		"MUST advance the entity to the next stage and dispatch it",
		"BEFORE ending its turn",
	} {
		if !strings.Contains(section, anchor) {
			t.Errorf("`## Completion and Gates` is missing the not-a-stopping-point obligation anchor %q", anchor)
		}
	}

	// A bare "advance before ending the turn" obligation with no enumerated halt
	// spans would over-claim (it would forbid stopping even at a gate). The clause
	// must name the four legitimate halts so the duty is scoped, not absolute.
	for _, halt := range []string{
		"`gate: true`",
		"terminal",
		"blocker",
		"captain decision",
	} {
		if !strings.Contains(section, halt) {
			t.Errorf("`## Completion and Gates` not-a-stopping-point clause does not name the halt exception %q", halt)
		}
	}

	// The clause must frame a completion-only stop as the forbidden case, not a
	// permitted one — a paraphrase that drops the prohibition would keep some of
	// the anchors above but lose this.
	if !strings.Contains(section, "stopping after a completion-only report is a contract violation") {
		t.Error("`## Completion and Gates` does not name a completion-only stop as a contract violation")
	}
}

// TestDevReadmeImplementationRoutesToValidation locks AC-2: the dev workflow
// README states implementation report filing is not a stopping point and routes
// immediately to fresh `validation` dispatch.
//
// Oracle: a presence test scoped to the `### `+"`implementation`"+` stage
// subsection via subsectionAfter (the `###`-aware scope helper), so a stray
// mention in the `### `+"`validation`"+` subsection cannot satisfy it. Anchors
// name `validation` and the `fresh` validator and the not-a-stopping-point
// routing; gutting the clause reds the test.
func TestDevReadmeImplementationRoutesToValidation(t *testing.T) {
	section := subsectionAfter(devWorkflowReadme(t), "### `implementation`")
	if section == "" {
		t.Fatal("dev README missing the `### `+\"`implementation`\"+` stage subsection")
	}
	// Guard: the helper must not have swept into the `### `+"`validation`"+`
	// subsection (which legitimately discusses validation). The implementation
	// subsection ends before the validation heading.
	if strings.Contains(section, "A task moves to validation after implementation is complete") {
		t.Fatal("implementation subsection scope leaked into the validation subsection — scoping is broken")
	}

	for _, anchor := range []string{
		"Implementation completion is not a stopping point",
		"routes immediately to independent `validation` dispatch",
		"`fresh: true`",
	} {
		if !strings.Contains(section, anchor) {
			t.Errorf("`### `+\"`implementation`\"+` subsection is missing the routes-to-validation anchor %q", anchor)
		}
	}
	// The routing must name its interrupt exceptions so it is scoped, not absolute,
	// and must state the FO does not park-and-wait.
	if !strings.Contains(section, "The FO does not park a completed implementation and wait") {
		t.Error("`### `+\"`implementation`\"+` subsection does not forbid parking a completed implementation")
	}
	for _, halt := range []string{"gate", "blocker", "terminal", "captain decision"} {
		if !strings.Contains(section, halt) {
			t.Errorf("`### `+\"`implementation`\"+` routes-to-validation clause does not name the interrupt exception %q", halt)
		}
	}
}

// TestPiRuntimePreservesSameTurnLifecycle locks AC-3: the Pi runtime guidance
// preserves the same lifecycle after a `pi-subagents` completion — once the
// stage report is verified for a non-gated, non-terminal stage, the parent must
// continue the shared `## Completion and Gates` lifecycle in the SAME turn
// (advance + dispatch, fresh subagent when `fresh: true`) and not return a
// completion-only result unless gated / terminal / blocked / awaiting a captain
// decision.
//
// Oracle: a presence test scoped to the `## Awaiting Completion` section via
// sectionAfter, so a clause elsewhere cannot satisfy it. Anchors name the
// same-turn obligation and the halt exceptions; gutting the clause reds the test.
func TestPiRuntimePreservesSameTurnLifecycle(t *testing.T) {
	section := sectionAfter(piFirstOfficerRuntime(t), "## Awaiting Completion")
	if section == "" {
		t.Fatal("Pi first-officer runtime missing the `## Awaiting Completion` section")
	}
	for _, anchor := range []string{
		"Verifying the stage report is not the end of the parent's turn",
		"MUST continue the shared `## Completion and Gates` lifecycle in the same turn",
		"fresh subagent when the next stage is `fresh: true`",
	} {
		if !strings.Contains(section, anchor) {
			t.Errorf("`## Awaiting Completion` is missing the same-turn-continuation anchor %q", anchor)
		}
	}
	// The clause must forbid returning a completion-only result for the captain to
	// resume from, scoped by the halt spans — a paraphrase dropping the prohibition
	// keeps some anchors but loses this.
	if !strings.Contains(section, "does not return a completion-only result for the captain to resume from") {
		t.Error("`## Awaiting Completion` does not forbid returning a completion-only result for the captain to resume from")
	}
	for _, halt := range []string{"gated", "terminal", "blocked", "captain decision"} {
		if !strings.Contains(section, halt) {
			t.Errorf("`## Awaiting Completion` same-turn clause does not name the halt span %q", halt)
		}
	}
}
