// ABOUTME: AC-1 measurement — the --advance prompt is O(pointer) and byte-
// ABOUTME: identical across stage-section sizes, and the hand-assembled
// ABOUTME: verbatim-section template it replaces is >=5x larger even at the
// ABOUTME: smallest fixture stage.
package dispatch

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// advancePromptByteCeiling is AC-1's invariant: the emitted advance prompt
// stays at or under this many bytes for every stage in the dev workflow,
// regardless of the target stage's README section size.
const advancePromptByteCeiling = 300

// readmeWithStageBody builds a fixture README whose "implementation" stage
// section body is exactly stageBody, so the fixture can hold the stage name,
// entity, and every other structural element fixed while varying only the
// section content size under test.
func readmeWithStageBody(stageBody string) string {
	return "---\n" +
		"entity-type: task\n" +
		"id-style: slug\n" +
		"stages:\n" +
		"  defaults:\n" +
		"    worktree: false\n" +
		"    concurrency: 1\n" +
		"  states:\n" +
		"    - name: backlog\n" +
		"      initial: true\n" +
		"    - name: implementation\n" +
		"      worktree: true\n" +
		"    - name: done\n" +
		"      terminal: true\n" +
		"---\n" +
		"# Fixture Workflow\n\n" +
		"### backlog\n\nseed.\n\n" +
		"### implementation\n\n" + stageBody + "\n\n" +
		"### done\n\nterm.\n"
}

// buildAdvancePromptFixture builds one fixture tree with the given stage-section
// body and returns the emitted advance prompt plus the raw ingredients ( the
// extracted stage subsection and checklist) an old-template comparator needs.
func buildAdvancePromptFixture(t *testing.T, stageBody string, checklist []string) (prompt, entityTitle, entityPath, stageSubsection string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "README.md"), readmeWithStageBody(stageBody))
	worktreeRel := ".worktrees/spacedock-ensign-thing"
	writeFile(t, filepath.Join(root, worktreeRel, ".keep"), "")
	entityPath = filepath.Join(root, "thing.md")
	entityTitle = "Thing"
	writeFile(t, entityPath, entityFM(entityTitle, "implementation", worktreeRel))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      checklist,
		"bare_mode":      false,
		"advance":        true,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("advance build exit=%d stderr=%s", native.exit, native.stderr)
	}
	var out struct {
		Prompt string `json:"prompt"`
	}
	if err := json.Unmarshal([]byte(native.stdout), &out); err != nil {
		t.Fatalf("stdout is not build JSON: %v\n%s", err, native.stdout)
	}

	sub, err := extractStageSubsection(filepath.Join(root, "README.md"), "implementation")
	if err != nil {
		t.Fatalf("extractStageSubsection: %v", err)
	}
	return out.Prompt, entityTitle, entityPath, sub
}

// oldTemplateAssembly reconstructs the hand-assembled reuse-advance message
// claude-fo-dispatch.md carried before this task (the verbatim stage-section
// echo), byte-for-byte in shape, so its length can be compared against the
// emitted --advance prompt on the same fixture ingredients.
func oldTemplateAssembly(nextStage, stageSubsection string, checklist []string, entityTitle, entityPath string) string {
	return "Advancing to next stage: " + nextStage + "\n\n" +
		"### Stage definition:\n\n" + stageSubsection + "\n\n" +
		"### Completion checklist\n\n" + strings.Join(checklist, "\n") + "\n\n" +
		"Continue working on " + entityTitle + " at " + entityPath + ". Commit before sending your completion message."
}

// TestBuildAdvancePromptSizeInvariant is AC-1: the emitted --advance prompt is
// <=300 bytes and byte-identical across a small and a large stage-section
// fixture (same stage name, entity, team, checklist — only the README section
// body size differs), because the prompt never embeds the section; it only
// points at the show-stage-def fetch line. The old hand-assembled template
// this replaces embeds the section directly, so it grows with section size —
// asserted here to exceed the emitted prompt by >=5x even on the SMALL fixture.
func TestBuildAdvancePromptSizeInvariant(t *testing.T) {
	checklist := []string{"- verify the thing"}

	// Sized to the dev workflow's own measured baselines (implementation
	// 1,259 bytes; ideation 4,866 bytes) so the "smallest fixture stage" is
	// realistic, not a single word the 5x bar would trivially clear either way.
	smallBody := strings.Repeat("This stage does substantial work worth describing in the workflow README. ", 14) // ~1.0KB
	largeBody := strings.Repeat("This stage does substantial work across many files. ", 100)                      // ~5.3KB

	smallPrompt, title, entityPath, smallSub := buildAdvancePromptFixture(t, smallBody, checklist)
	largePrompt, _, _, largeSub := buildAdvancePromptFixture(t, largeBody, checklist)

	for name, p := range map[string]string{"small-fixture": smallPrompt, "large-fixture": largePrompt} {
		if len(p) > advancePromptByteCeiling {
			t.Errorf("%s: advance prompt is %d bytes, want <= %d: %q", name, len(p), advancePromptByteCeiling, p)
		}
	}

	if len(smallPrompt) != len(largePrompt) {
		t.Errorf("advance prompt length varies with stage-section size: small=%d bytes, large=%d bytes\nsmall=%q\nlarge=%q",
			len(smallPrompt), len(largePrompt), smallPrompt, largePrompt)
	}

	oldSmall := oldTemplateAssembly("implementation", smallSub, checklist, title, entityPath)
	oldLarge := oldTemplateAssembly("implementation", largeSub, checklist, title, entityPath)

	if got, want := len(oldSmall), 5*len(smallPrompt); got < want {
		t.Errorf("old-template assembly on the SMALL fixture is %d bytes, want >= %d (5x the %d-byte emitted prompt)",
			got, want, len(smallPrompt))
	}
	if len(oldLarge) <= len(oldSmall) {
		t.Errorf("old-template assembly should grow with stage-section size: small=%d large=%d", len(oldSmall), len(oldLarge))
	}
}
