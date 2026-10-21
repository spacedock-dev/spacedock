// ABOUTME: Checks the argument branches and the choice wording of skills/review/SKILL.md:
// ABOUTME: no argument only shows or changes the choice, a non-file argument writes first.
package integration

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func reviewSkillSections(t *testing.T) map[string]string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), "skills", "review", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	sections := map[string]string{"": ""}
	current := ""
	for _, line := range strings.Split(string(data), "\n") {
		if name, ok := strings.CutPrefix(line, "## "); ok {
			current = name
			sections[current] = ""
			continue
		}
		sections[current] += line + "\n"
	}
	return sections
}

func reviewSection(t *testing.T, name string) string {
	t.Helper()
	body, ok := reviewSkillSections(t)[name]
	if !ok {
		t.Fatalf("skills/review/SKILL.md has no %q section", name)
	}
	return body
}

func TestReviewSkillNoArgumentOnlyShowsAndChangesTheChoice(t *testing.T) {
	body := reviewSection(t, "Show or change the choice")
	for _, want := range []string{
		"review-mode get --json", "review-mode set local", "review-mode clear",
		"keep it", "switch", "clear it", "leave", "Publish nothing",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("no-argument branch lacks %q", want)
		}
	}
	for _, forbidden := range []string{"subspace:r", "subspace:remote-review", "--supports"} {
		if strings.Contains(body, forbidden) {
			t.Errorf("no-argument branch reaches a hand-off: it mentions %q", forbidden)
		}
	}
	if pick := reviewSection(t, "Pick the branch from the argument"); !strings.Contains(pick, "No argument: do **Show or change the choice**") {
		t.Errorf("the branch list does not send a missing argument to the no-argument branch:\n%s", pick)
	}
}

func TestReviewSkillChoiceIsWordedViewFirst(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(repoRoot(t), "skills", "review", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, want := range []string{
		"Terminal (TUI): the file stays on this machine\n",
		"Web: uploads to review.spacedock.md; anyone with the link can open it\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("skill does not word a choice as %q", strings.TrimSpace(want))
		}
	}
	impliesViewer := regexp.MustCompile(`(?i)\blocal\b[^\n.]*\b(opens?|viewer)\b|\b(opens?|viewer)\b[^\n.]*\blocal\b`)
	for _, line := range strings.Split(text, "\n") {
		if impliesViewer.MatchString(line) {
			t.Errorf("a sentence ties local to opening a viewer: %q", line)
		}
	}
}

func TestReviewSkillNonFileArgumentWritesTheArtifactFirst(t *testing.T) {
	pick := reviewSection(t, "Pick the branch from the argument")
	for _, want := range []string{
		"An argument that names an existing file (`test -f`): do **Review the file**",
		"Any other argument is a request, not a path: do **Write first, then review**",
	} {
		if !strings.Contains(pick, want) {
			t.Errorf("branch list lacks %q", want)
		}
	}
	body := reviewSection(t, "Write first, then review")
	for _, want := range []string{"Write the artifact", "Never overwrite an existing file", "**Review the file**"} {
		if !strings.Contains(body, want) {
			t.Errorf("non-file branch lacks %q", want)
		}
	}
	if strings.Index(body, "Write the artifact") > strings.Index(body, "**Review the file**") {
		t.Error("non-file branch reviews before it writes")
	}
}

func TestReviewSkillFileBranchStillAsksOnlyWhenNoChoiceIsRecorded(t *testing.T) {
	body := reviewSection(t, "Review the file")
	for _, want := range []string{"If `mode` is empty, ask the user", "review-mode set local", "subspace:r", "subspace:remote-review"} {
		if !strings.Contains(body, want) {
			t.Errorf("file branch lacks %q", want)
		}
	}
}
