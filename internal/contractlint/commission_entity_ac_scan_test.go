// ABOUTME: Feeds every entity body the commission skill generates through the real
// ABOUTME: `status --read --ac-scan` parser, so a commissioned entity always carries an AC section.
package contractlint

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// commissionEntityShape is one fenced entity body the commission skill writes into
// a generated workflow: a seed entity or a README entity template.
type commissionEntityShape struct {
	source string
	body   string
}

// fencedBlockAfter returns the content of the first ``` fenced block that follows
// the heading line in text. A ```` outer fence (the SKILL's README template) is
// stepped into, so a nested ``` block is still found.
func fencedBlockAfter(t *testing.T, text, source, heading string) string {
	t.Helper()
	lines := strings.Split(text, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == heading {
			start = i + 1
			break
		}
	}
	if start < 0 {
		t.Fatalf("%s: heading %q not found", source, heading)
	}
	for i := start; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "```") || strings.HasPrefix(lines[i], "````") {
			continue
		}
		for j := i + 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == "```" {
				return strings.Join(lines[i+1:j], "\n")
			}
		}
		t.Fatalf("%s: unterminated fence after %q", source, heading)
	}
	t.Fatalf("%s: no fenced block after %q", source, heading)
	return ""
}

// templateEntityHeadingRe matches a selected template's entity-template heading
// (`## Task Template`, `## Experiment Template`, `## Artifact Template`).
var templateEntityHeadingRe = regexp.MustCompile(`(?m)^## [A-Za-z]+ Template$`)

func commissionEntityShapes(t *testing.T) []commissionEntityShape {
	t.Helper()
	const skill = "skills/commission/SKILL.md"
	skillText := readRepoFile(t, skill)
	shapes := []commissionEntityShape{
		{skill + " seed entity", fencedBlockAfter(t, skillText, skill, "### 2b. Generate Seed Entities")},
		{skill + " README entity template", fencedBlockAfter(t, skillText, skill, "## {Entity_label} Template")},
	}
	templates, err := filepath.Glob(filepath.Join(repoRoot(t), "skills", "commission", "references", "templates", "*.md"))
	if err != nil || len(templates) == 0 {
		t.Fatalf("no commission templates found: %v", err)
	}
	for _, path := range templates {
		rel, _ := filepath.Rel(repoRoot(t), path)
		text := readRepoFile(t, rel)
		headings := templateEntityHeadingRe.FindAllString(text, -1)
		if len(headings) != 1 {
			t.Fatalf("%s: want exactly one `## <Entity> Template` heading, found %v", rel, headings)
		}
		shapes = append(shapes, commissionEntityShape{rel + " " + headings[0], fencedBlockAfter(t, text, rel, headings[0])})
	}
	return shapes
}

// TestCommissionedEntitiesPassACScan builds an entity from each body shape the
// commission skill generates, appends a stage report citing AC-1, and runs the
// real `status --read --ac-scan` over it. A shape without `## Acceptance criteria`
// fails ac-scan loudly, which halts every non-initial gate the entity reaches.
func TestCommissionedEntitiesPassACScan(t *testing.T) {
	for _, shape := range commissionEntityShapes(t) {
		t.Run(shape.source, func(t *testing.T) {
			dir := t.TempDir()
			readme := "---\ncommissioned-by: spacedock@test\nentity-type: task\nstages:\n  states:\n    - name: backlog\n      initial: true\n    - name: implementation\n    - name: done\n      terminal: true\n---\n\n# Test workflow\n"
			if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
				t.Fatal(err)
			}
			entity := filepath.Join(dir, "seed-entity.md")
			body := strings.TrimRight(shape.body, "\n") + "\n\n## Stage Report: implementation\n\n- DONE: verified AC-1\n"
			if err := os.WriteFile(entity, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			code, err := (&status.NativeRunner{}).Run(context.Background(), status.Request{
				Args:   []string{"--workflow-dir", dir, "--read", entity, "--stage", "implementation", "--ac-scan", "--json"},
				Dir:    dir,
				Env:    []string{"HOME=" + t.TempDir(), "PATH=" + os.Getenv("PATH")},
				Stdout: &stdout,
				Stderr: &stderr,
			})
			if err != nil || code != 0 {
				t.Fatalf("ac-scan exit=%d err=%v stderr=%q\nentity body:\n%s", code, err, stderr.String(), body)
			}
			var doc struct {
				ACs []struct {
					ID string `json:"id"`
				} `json:"acs"`
			}
			if err := json.Unmarshal(stdout.Bytes(), &doc); err != nil {
				t.Fatalf("ac-scan output is not JSON: %v\n%s", err, stdout.String())
			}
			if len(doc.ACs) == 0 || doc.ACs[0].ID != "AC-1" {
				t.Fatalf("ac-scan found no **AC-1** stub in the acceptance criteria; got %s", stdout.String())
			}
		})
	}
}
