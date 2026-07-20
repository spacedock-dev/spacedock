// ABOUTME: The filing live fixture exposes the same workflow-local Task Template shape as a real workflow.
// ABOUTME: Its literal fenced body round-trips through NativeRunner new, while indented frontmatter stays invalid.
package ensigncycle

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/status"
)

func TestFilingReadmeTaskTemplateRoundTripsThroughNew(t *testing.T) {
	root := t.TempDir()
	writeFilingWorkflow(t, root)

	readmeBytes, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	template := filingTaskTemplate(t, string(readmeBytes))

	if !strings.HasPrefix(template, "---\n") {
		t.Fatalf("Task Template frontmatter must begin at column zero:\n%q", template)
	}
	if strings.Contains(template, "\nid:") || strings.HasPrefix(template, "id:") {
		t.Fatalf("Task Template must omit id so new can mint it:\n%s", template)
	}

	t.Run("literal template creates entity and mints id", func(t *testing.T) {
		stdout, stderr, code := runFilingTemplateNew(t, root, "from-template", template)
		if code != 0 {
			t.Fatalf("new exit=%d stderr=%q", code, stderr)
		}
		if !strings.Contains(stdout, "created:") {
			t.Fatalf("new stdout=%q, want created narration", stdout)
		}

		entityBytes, err := os.ReadFile(filepath.Join(root, "from-template.md"))
		if err != nil {
			t.Fatal(err)
		}
		entity := string(entityBytes)
		if !strings.Contains(entity, "title: Task name here\n") || !strings.Contains(entity, "status: backlog\n") {
			t.Fatalf("created entity did not preserve template frontmatter:\n%s", entity)
		}
		if !hasNonEmptyFrontmatterField(entity, "id") {
			t.Fatalf("created entity has no minted id:\n%s", entity)
		}
	})

	t.Run("indented frontmatter fence is rejected", func(t *testing.T) {
		controlRoot := t.TempDir()
		writeFilingWorkflow(t, controlRoot)
		indented := "  " + template
		_, stderr, code := runFilingTemplateNew(t, controlRoot, "from-template", indented)
		if code != 1 {
			t.Fatalf("new exit=%d, want 1 for indented opening fence", code)
		}
		if !strings.Contains(stderr, "no frontmatter") {
			t.Fatalf("new stderr=%q, want no frontmatter error", stderr)
		}
		if _, err := os.Stat(filepath.Join(controlRoot, "from-template.md")); !os.IsNotExist(err) {
			t.Fatalf("indented control created an entity; stat error=%v", err)
		}
	})
}

func filingTaskTemplate(t *testing.T, readme string) string {
	t.Helper()
	const heading = "## Task Template\n"
	if got := strings.Count(readme, heading); got != 1 {
		t.Fatalf("README has %d Task Template headings, want 1", got)
	}
	section := readme[strings.Index(readme, heading)+len(heading):]
	if got := strings.Count(section, "```"); got != 2 {
		t.Fatalf("Task Template section has %d fence markers, want 2:\n%s", got, section)
	}
	const opening = "```yaml\n"
	start := strings.Index(section, opening)
	if start < 0 {
		t.Fatalf("Task Template has no YAML opening fence:\n%s", section)
	}
	body := section[start+len(opening):]
	end := strings.Index(body, "```\n")
	if end < 0 {
		t.Fatalf("Task Template has no closing fence:\n%s", section)
	}
	return body[:end]
}

func runFilingTemplateNew(t *testing.T, root, slug, body string) (string, string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	runner := &status.NativeRunner{}
	code, err := runner.Run(context.Background(), status.Request{
		Args:   []string{"--workflow-dir", root, "--new", slug},
		Dir:    root,
		Env:    os.Environ(),
		Stdin:  strings.NewReader(body),
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		t.Fatal(err)
	}
	return stdout.String(), stderr.String(), code
}

func hasNonEmptyFrontmatterField(entity, name string) bool {
	for _, line := range strings.Split(entity, "\n") {
		if key, value, ok := strings.Cut(line, ":"); ok && key == name {
			return strings.TrimSpace(value) != ""
		}
	}
	return false
}
