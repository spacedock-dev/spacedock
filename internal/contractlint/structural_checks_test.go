// ABOUTME: Keeps model-read skill frontmatter, discovery, and reference topology structurally valid.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var userSkills = []string{
	"commission", "debrief", "refit", "survey", "ensign",
	"first-officer", "present-gate", "feedback-rejection-flow",
}

func skillsRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "skills")
}

func frontmatter(doc string) (string, bool) {
	if !strings.HasPrefix(doc, "---\n") {
		return "", false
	}
	rest := doc[len("---\n"):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return "", false
	}
	return rest[:end], true
}

func TestUserSkillsParseWithFrontmatter(t *testing.T) {
	root := skillsRoot(t)
	for _, skill := range userSkills {
		path := filepath.Join(root, skill, "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("user skill %q has no SKILL.md at %s: %v", skill, path, err)
			continue
		}
		fm, ok := frontmatter(string(data))
		if !ok {
			t.Errorf("%s/SKILL.md has no parseable YAML frontmatter block", skill)
			continue
		}
		if !frontmatterHasKey(fm, "name") {
			t.Errorf("%s/SKILL.md frontmatter declares no `name` key", skill)
		}
		if !frontmatterHasKey(fm, "description") {
			t.Errorf("%s/SKILL.md frontmatter declares no `description` key", skill)
		}
	}
}

func frontmatterHasKey(fm, key string) bool {
	prefix := key + ":"
	for _, line := range strings.Split(fm, "\n") {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

func frontmatterField(fm, key string) string {
	prefix := key + ":"
	for _, line := range strings.Split(fm, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
		}
	}
	return ""
}

func discoverUserInvocableSkills(t *testing.T) map[string]string {
	t.Helper()
	root := skillsRoot(t)
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read skills root %s: %v", root, err)
	}
	out := map[string]string{}
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "integration" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, e.Name(), "SKILL.md"))
		if err != nil {
			continue
		}
		fm, ok := frontmatter(string(data))
		if !ok || frontmatterField(fm, "user-invocable") != "true" {
			continue
		}
		name := frontmatterField(fm, "name")
		if name == "" {
			t.Errorf("user-invocable skill dir %q has no name field", e.Name())
			continue
		}
		out[name] = e.Name()
	}
	return out
}

func TestSurveyIsDiscoverableUserCommand(t *testing.T) {
	discovered := discoverUserInvocableSkills(t)
	dir, ok := discovered["survey"]
	if !ok {
		t.Fatalf("survey is not discoverable; discovered user commands: %v", sortedUniqueKeys(discovered))
	}
	if dir != "survey" {
		t.Errorf("survey command resolves from dir %q, want skills/survey", dir)
	}
}

func TestFOGateLifecycleIsDeferredAndAdapterless(t *testing.T) {
	dir := filepath.Join(skillsRoot(t), "fo-gate-lifecycle")
	data, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	fm, ok := frontmatter(string(data))
	if !ok || frontmatterField(fm, "name") != "fo-gate-lifecycle" ||
		frontmatterField(fm, "user-invocable") != "false" {
		t.Fatalf("fo-gate-lifecycle must be named and non-user-invocable: %q", fm)
	}
	if _, found := discoverUserInvocableSkills(t)["fo-gate-lifecycle"]; found {
		t.Fatal("fo-gate-lifecycle leaked into user command discovery")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "SKILL.md" {
		t.Fatalf("fo-gate-lifecycle must remain adapter-less, entries=%v", entries)
	}
}

func sortedUniqueKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return sortedUnique(keys)
}

var referenceRe = regexp.MustCompile(`@?(references/[A-Za-z0-9_./-]+\.md)`)

func TestUserSkillReferenceClosureResolves(t *testing.T) {
	root := skillsRoot(t)
	for _, skill := range userSkills {
		skillDir := filepath.Join(root, skill)
		data, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
		if err != nil {
			t.Errorf("%s: %v", skill, err)
			continue
		}
		for _, m := range referenceRe.FindAllStringSubmatch(string(data), -1) {
			rel := m[1]
			if strings.Contains(rel, "{") {
				parent := filepath.Join(skillDir, filepath.Dir(rel))
				glob, _ := filepath.Glob(filepath.Join(parent, "*.md"))
				if len(glob) == 0 {
					t.Errorf("%s: templated reference %q has no concrete .md under %s", skill, rel, parent)
				}
				continue
			}
			if _, err := os.Stat(filepath.Join(skillDir, rel)); err != nil {
				t.Errorf("%s: dangling reference %q: %v", skill, rel, err)
			}
		}
	}
}

var piRuntimeAdapters = []struct{ skill, ref string }{
	{skill: "first-officer", ref: "references/pi-first-officer-runtime.md"},
	{skill: "ensign", ref: "references/pi-ensign-runtime.md"},
}

func TestPiRuntimeAdaptersResolveOnDisk(t *testing.T) {
	root := skillsRoot(t)
	for _, tc := range piRuntimeAdapters {
		path := filepath.Join(root, tc.skill, tc.ref)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("%s: Pi runtime adapter %s is not loadable on disk: %v", tc.skill, tc.ref, err)
		}
	}
}
