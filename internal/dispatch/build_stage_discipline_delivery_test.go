// ABOUTME: Behavior-loss control — proves dev stage discipline rides the dev-shape
// ABOUTME: show-stage-def fetch, not the universal ensign core, so neutralizing the core loses none.
package dispatch

import (
	"path/filepath"
	"strings"
	"testing"
)

// devDisciplineSentinel is a freshly-minted token written ONLY into the fixture
// README's ### ideation / ### implementation subsections. Because the test invents
// it, it cannot appear in the shipped universal ensign core — so its absence from the
// dispatch body (which delivers the core by reference, via the Skill first-action)
// proves the stage prose is not inlined into the assignment, and its presence in
// show-stage-def output proves the fetch delivers it.
const devDisciplineSentinel = "DEV-DISCIPLINE-SENTINEL"

// readmeDevDiscipline is a dev-shape workflow README that plants the sentinel in the
// ideation and implementation stage subsections (and deliberately not in backlog /
// validation / done), so a stage with the discipline and a stage without it can both
// be exercised through show-stage-def.
const readmeDevDiscipline = `---
entity-type: task
id-style: slug
stages:
  defaults:
    worktree: false
    concurrency: 1
  states:
    - name: backlog
      initial: true
    - name: ideation
    - name: implementation
      worktree: true
    - name: validation
      worktree: true
      feedback-to: implementation
    - name: done
      terminal: true
---
# Dev Discipline Fixture

### backlog

seed.

- **Outputs:** x.

### ideation

Flesh out the approach. ` + devDisciplineSentinel + `

- **Outputs:** the design.

### implementation

Produce the deliverable. ` + devDisciplineSentinel + `

- **Outputs:** the deliverable.

### validation

Independently verify.

- **Outputs:** the verdict.

### done

term.
`

// TestBuildStageDisciplineRidesFetchNotInlineAssignment is the behavior-loss control
// for neutralizing the universal ensign core. It proves a dev-workflow ensign still
// receives dev stage discipline through the dev-shape scaffolding — the README
// ### {stage} subsection delivered via the assignment's show-stage-def fetch line —
// and NOT through the universal core. So removing the core's illustrative dev
// stage-name parentheticals loses no discipline: the core never carried the discipline
// to begin with; the fetch does.
func TestBuildStageDisciplineRidesFetchNotInlineAssignment(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeDevDiscipline)
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "ideation", ""))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "ideation",
		"checklist":      []string{"- a"},
		"team_name":      "fixture-team",
		"bare_mode":      false,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit %d, stderr:\n%s", native.exit, native.stderr)
	}
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	// The assignment loads the universal ensign core by REFERENCE (the Skill
	// first-action directive), not inline.
	if !strings.Contains(body, `Skill(skill="spacedock:ensign")`) {
		t.Errorf("dispatch body missing the universal-core load directive (Skill first-action):\n%s", body)
	}

	// Stage discipline rides the show-stage-def fetch line, fetched on demand.
	wantFetch := LauncherCommand() + " dispatch show-stage-def --workflow-dir " + shlexQuote(root) + " --stage ideation"
	if !strings.Contains(body, wantFetch) {
		t.Errorf("dispatch body missing the show-stage-def fetch line\nwant contains: %s\n---body---\n%s", wantFetch, body)
	}

	// The stage prose is NOT inlined into the assignment: the sentinel that lives in
	// the README's ### ideation subsection must not appear in the dispatch body. The
	// universal ensign core that the Skill directive loads never carries this
	// test-minted sentinel either — it exists only in the fixture README.
	if strings.Contains(body, devDisciplineSentinel) {
		t.Errorf("dispatch body INLINED the stage prose — the %s sentinel leaked into the assignment instead of riding the fetch:\n%s", devDisciplineSentinel, body)
	}

	// Running the fetch line returns the ### ideation subsection carrying the sentinel:
	// the dev discipline reaches the ensign through the show-stage-def fetch.
	fetched := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", "ideation")
	if fetched.exit != 0 {
		t.Fatalf("show-stage-def ideation exit %d, stderr:\n%s", fetched.exit, fetched.stderr)
	}
	if !strings.Contains(fetched.stdout, devDisciplineSentinel) {
		t.Errorf("show-stage-def ideation did not return the stage-discipline sentinel:\n%s", fetched.stdout)
	}

	// Perturbation control: a stage WITHOUT the sentinel (validation) returns no
	// sentinel — show-stage-def returns ONLY the requested stage's subsection, so a
	// false green from the sentinel leaking across stages is excluded.
	plain := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", "validation")
	if plain.exit != 0 {
		t.Fatalf("show-stage-def validation exit %d, stderr:\n%s", plain.exit, plain.stderr)
	}
	if strings.Contains(plain.stdout, devDisciplineSentinel) {
		t.Errorf("perturbation control: show-stage-def validation (no sentinel planted) returned the sentinel:\n%s", plain.stdout)
	}
}
