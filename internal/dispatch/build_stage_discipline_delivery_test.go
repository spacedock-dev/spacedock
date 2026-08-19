// ABOUTME: Behavior-loss control — proves dev stage discipline rides the dev-shape
// ABOUTME: show-stage-def fetch, not the universal ensign core.
package dispatch

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// devDisciplineSentinel is a freshly-minted token written ONLY into the fixture
// README's ### ideation / ### implementation subsections. The exact selected source
// output makes any cross-stage leak fail without searching generated prose.
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

// TestBuildStageDisciplineRidesExactFetchCommand proves the selected README stage
// prose stays behind a generated full-path fetch command and out of the universal
// ensign core, dispatch body, and outer pointer prompt.
func TestBuildStageDisciplineRidesExactFetchCommand(t *testing.T) {
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
		"bare_mode":      false,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit %d, stderr:\n%s", native.exit, native.stderr)
	}
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	var envelope struct {
		Prompt string   `json:"prompt"`
		Fetch  []string `json:"fetch_commands"`
	}
	if err := json.Unmarshal([]byte(native.stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	wantFetch := testWorkflowLauncher + " dispatch show-stage-def --workflow-dir " + shlexQuote(root) + " --stage ideation"
	sections := dispatchArtifactSections(t, body)
	if len(envelope.Fetch) != 1 || envelope.Fetch[0] != wantFetch || sections["Fetch commands"] != wantFetch {
		t.Errorf("assignment did not render the exact full-path stage command: envelope=%#v\n%s", envelope.Fetch, body)
	}
	if sections["Completion checklist"] != "- a" {
		t.Errorf("assignment checklist differs from structured input: %q", sections["Completion checklist"])
	}

	// The public inspection command remains available and returns the same selection.
	fetched := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", "ideation")
	if fetched.exit != 0 {
		t.Fatalf("show-stage-def ideation exit %d, stderr:\n%s", fetched.exit, fetched.stderr)
	}
	wantFetched := "### ideation\n\nFlesh out the approach. " + devDisciplineSentinel + "\n\n- **Outputs:** the design.\n"
	if fetched.stdout != wantFetched {
		t.Errorf("show-stage-def ideation differs from the selected source section:\n got=%q\nwant=%q", fetched.stdout, wantFetched)
	}

	// Perturbation control: a stage WITHOUT the sentinel (validation) returns no
	// sentinel — show-stage-def returns ONLY the requested stage's subsection, so a
	// false green from the sentinel leaking across stages is excluded.
	plain := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", "validation")
	if plain.exit != 0 {
		t.Fatalf("show-stage-def validation exit %d, stderr:\n%s", plain.exit, plain.stderr)
	}
	if plain.stdout != "### validation\n\nIndependently verify.\n\n- **Outputs:** the verdict.\n" {
		t.Errorf("show-stage-def validation differs from the selected source section: %q", plain.stdout)
	}
}

func TestDeclaredContextBuildAndHostNeutralFetch(t *testing.T) {
	root := t.TempDir()
	readme := "---\nentity-type: task\nid-style: slug\nstages:\n  defaults:\n" +
		"    context-sections: [Pølicy]\n  states:\n    - name: ideation\n      initial: true\n---\n" +
		"# Workflow\n\n### ideation\r\nstage β\r\n\r\n## Pølicy\r\nα\rβ\u0085γ\r\n \t\r\n## Unrelated\nnope\n"
	writeFile(t, filepath.Join(root, "README.md"), readme)
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "ideation", ""))
	gitInit(t, root)

	for _, host := range []string{"claude", "codex", "pi"} {
		t.Run(host, func(t *testing.T) {
			req := mergeStdin(map[string]any{
				"schema_version": 2, "entity_path": entityPath, "workflow_dir": root,
				"stage": "ideation", "checklist": []string{"- prove context"},
			}, nil)
			built := runNative(req, "build", "--workflow-dir", root, "--host", host)
			if built.exit != 0 {
				t.Fatalf("build exit=%d stderr=%q", built.exit, built.stderr)
			}
			var envelope struct {
				Fetch []string `json:"fetch_commands"`
			}
			if err := json.Unmarshal([]byte(built.stdout), &envelope); err != nil {
				t.Fatal(err)
			}
			wantCommand := testWorkflowLauncher + " dispatch show-stage-def --workflow-dir " + shlexQuote(root) + " --stage ideation"
			if len(envelope.Fetch) != 1 || envelope.Fetch[0] != wantCommand {
				t.Fatalf("fetch commands = %#v, want one exact full-path show-stage-def command", envelope.Fetch)
			}
			want := "### ideation\nstage β\n\n## Pølicy\nα\nβ\nγ\n"
			got := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", "ideation")
			if got.exit != 0 || got.stdout != want {
				t.Fatalf("show-stage-def exit=%d stderr=%q\n got=%q\nwant=%q", got.exit, got.stderr, got.stdout, want)
			}
		})
	}
}

func TestDeclaredContextInheritanceReplacementAndClear(t *testing.T) {
	root := t.TempDir()
	readme := "---\nstages:\n  defaults:\n    context-sections: [Authority]\n  states:\n" +
		"    - name: inherited\n    - name: replaced\n      context-sections: [Safety, Authority]\n" +
		"    - name: cleared\n      context-sections: []\n---\n# Workflow\n\n" +
		"### inherited\ninherited body\n\n### replaced\nreplacement body\n\n### cleared\nclear body\n\n" +
		"## Authority\nauthority\n\n## Safety\nsafety\n\n## Unrelated\nnever\n"
	writeFile(t, filepath.Join(root, "README.md"), readme)
	cases := map[string]string{
		"inherited": "### inherited\ninherited body\n\n## Authority\nauthority\n",
		"replaced":  "### replaced\nreplacement body\n\n## Safety\nsafety\n\n## Authority\nauthority\n",
		"cleared":   "### cleared\nclear body\n",
	}
	for stage, want := range cases {
		got := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", stage)
		if got.exit != 0 || got.stdout != want {
			t.Fatalf("%s exit=%d stderr=%q\n got=%q\nwant=%q", stage, got.exit, got.stderr, got.stdout, want)
		}
	}
}

func TestDeclaredContextInvalidBuildPreflight(t *testing.T) {
	cases := []struct {
		name, declaration, body, want string
	}{
		{"malformed", "context-sections: [", "## Authority\nok\n", "malformed YAML"},
		{"wrong kind", "context-sections: Authority", "## Authority\nok\n", "want sequence"},
		{"missing", "context-sections: [Missing]", "## Authority\nok\n", `selector "Missing" matches 0 headings`},
		{"repeated", "context-sections: [Authority, Authority]", "## Authority\nok\n", `repeated selector "Authority"`},
		{"ambiguous", "context-sections: [Authority]", "## Authority\none\n## Authority\ntwo\n", "matches 2 headings"},
		{"parent contains stage", "context-sections: [Parent]", "## Parent\n### ideation\nwork\n## End\n", "overlap"},
		{"child inside stage", "context-sections: [Child]", "### ideation\nwork\n#### Child\ninside\n## End\n", "overlap"},
		{"selected parent child", "context-sections: [Parent, Child]", "### ideation\nwork\n## Parent\n### Child\ninside\n## End\n", "overlap"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			readme := "---\nentity-type: task\nid-style: slug\nstages:\n  defaults:\n    worktree: false\n" +
				"    " + tc.declaration + "\n  states:\n    - name: ideation\n      initial: true\n---\n# Workflow\n\n" + tc.body
			if !strings.Contains(tc.body, "### ideation") {
				readme += "\n### ideation\nwork\n"
			}
			writeFile(t, filepath.Join(root, "README.md"), readme)
			entityPath := filepath.Join(root, "thing.md")
			writeFile(t, entityPath, entityFM("Thing", "ideation", ""))
			gitInit(t, root)
			req := mergeStdin(map[string]any{
				"schema_version": 2, "entity_path": entityPath, "workflow_dir": root,
				"stage": "ideation", "checklist": []string{"- x"},
			}, nil)
			got := runNative(req, "build", "--workflow-dir", root)
			if got.exit == 0 || !strings.Contains(got.stderr, tc.want) ||
				!strings.Contains(got.stderr, filepath.Join(root, "README.md")) ||
				!strings.Contains(got.stderr, "ideation") {
				t.Fatalf("exit=%d stdout=%q stderr=%q, want failure containing %q, path, and stage",
					got.exit, got.stdout, got.stderr, tc.want)
			}
		})
	}
}

func TestDeclaredContextLiveReadAdoptsValidAndRejectsInvalidCurrent(t *testing.T) {
	root := t.TempDir()
	readme := func(declaration, authority, body string) string {
		return "---\nentity-type: task\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      " +
			declaration + "\n---\n# Workflow\n\n" + body + "\n## Authority\n" + authority + "\n"
	}
	path := filepath.Join(root, "README.md")
	writeFile(t, path, readme("context-sections: [Authority]", "authority-v1", "### ideation\nwork"))
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "ideation", ""))
	gitInit(t, root)
	req := mergeStdin(map[string]any{
		"schema_version": 2, "entity_path": entityPath, "workflow_dir": root,
		"stage": "ideation", "checklist": []string{"- x"},
	}, nil)
	if got := runNative(req, "build", "--workflow-dir", root); got.exit != 0 {
		t.Fatalf("v1 build exit=%d stderr=%q", got.exit, got.stderr)
	}

	writeFile(t, path, readme("context-sections: [Authority]", "authority-v2", "### ideation\nwork"))
	if got := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", "ideation"); got.exit != 0 ||
		!strings.Contains(got.stdout, "authority-v2") || strings.Contains(got.stdout, "authority-v1") {
		t.Fatalf("valid current edit not adopted: exit=%d stdout=%q stderr=%q", got.exit, got.stdout, got.stderr)
	}
	for name, current := range map[string]string{
		"wrong-kind": readme("context-sections: Authority", "authority-v3", "### ideation\nwork"),
		"missing":    readme("context-sections: [Missing]", "authority-v3", "### ideation\nwork"),
		"overlap":    readme("context-sections: [Parent]", "authority-v3", "## Parent\n### ideation\nwork\n## End"),
	} {
		t.Run(name, func(t *testing.T) {
			writeFile(t, path, current)
			got := runNative("", "show-stage-def", "--workflow-dir", root, "--stage", "ideation")
			if got.exit == 0 || got.stdout != "" {
				t.Fatalf("invalid current edit exit=%d stdout=%q stderr=%q", got.exit, got.stdout, got.stderr)
			}
		})
	}
}

func TestMixedRawSpanReconciliationAndFencedStageCompatibility(t *testing.T) {
	fenced := []byte("```md\n### ideation\n\nfenced-body\n```\n### ideation\nreal\n")
	got, span, err := extractStageSubsectionBytes(fenced, "ideation")
	if err != nil || got != "### ideation\n\nfenced-body\n```" ||
		string(fenced[span.start:span.end]) != "### ideation\n\nfenced-body\n```\n" {
		t.Fatalf("fenced legacy result=%q span=%+v slice=%q err=%v", got, span, fenced[span.start:span.end], err)
	}

	mixed := []byte("## Pärent\r\nα\r\n### `build` *(captain)*\rstage β\vcontinuation\r\n\r\n#### Child\r\nchild γ\r\r## Sibling\r\nδ\r\n")
	_, stage, err := extractStageSubsectionBytes(mixed, "build")
	if err != nil {
		t.Fatal(err)
	}
	spans, err := status.FindSectionSpans(mixed, []string{"Pärent", "Child", "Sibling"})
	if err != nil {
		t.Fatal(err)
	}
	parent := sourceSpan{spans[0].Start, spans[0].End}
	child := sourceSpan{spans[1].Start, spans[1].End}
	sibling := sourceSpan{spans[2].Start, spans[2].End}
	if !spansIntersect(stage, parent) || !spansIntersect(stage, child) ||
		spansIntersect(stage, sibling) || spansIntersect(parent, sibling) {
		t.Fatalf("unexpected intersection matrix: stage=%+v spans=%+v", stage, spans)
	}
}
