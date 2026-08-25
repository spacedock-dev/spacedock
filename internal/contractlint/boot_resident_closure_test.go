// ABOUTME: Resolves deferred First Officer load-point references against the filesystem.
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var bootResidentBodies = []string{
	filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md"),
	filepath.Join("skills", "first-officer", "references", "claude-first-officer-runtime.md"),
	filepath.Join("skills", "first-officer", "references", "codex-first-officer-runtime.md"),
	filepath.Join("skills", "first-officer", "references", "pi-first-officer-runtime.md"),
}

var deferredModuleBodies = []string{
	filepath.Join("skills", "first-officer", "references", "claude-fo-dispatch.md"),
	filepath.Join("skills", "first-officer", "references", "fo-dispatch-core.md"),
}

// movedSections pins sections relocated out of a boot-resident body into a
// deferred module. A move is only legitimate when the section is still
// reachable at its new trigger, so each row asserts three things together:
// the heading is gone from the source, it exists exactly once in the target,
// and the shared core's deferred load-points block names that target.
var movedSections = []struct {
	heading string
	from    string
	to      string
}{
	{
		heading: "## Completion and Gates",
		from:    filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md"),
		to:      filepath.Join("skills", "first-officer", "references", "fo-dispatch-core.md"),
	},
}

var bodyReferenceRe = regexp.MustCompile(`references/[A-Za-z0-9_./-]+\.md`)
var bodySkillRe = regexp.MustCompile(`spacedock:([a-z0-9-]+)`)
var bodyModRe = regexp.MustCompile(`_mods/([a-z0-9][a-z0-9_.-]*\.md)`)

var lazyLoadSkills = map[string]bool{
	"present-gate":            true,
	"feedback-rejection-flow": true,
	"fo-gate-lifecycle":       true,
	"fo-status-viewer":        true,
	"fo-dispatch-recovery":    true,
}

type deferredLoadPoint struct {
	path  string
	named string
}

func extractDeferredLoadPoints(body string) []deferredLoadPoint {
	foSkillDir := filepath.Join("skills", "first-officer")
	var out []deferredLoadPoint
	seen := map[string]bool{}
	add := func(p deferredLoadPoint) {
		if seen[p.path] {
			return
		}
		seen[p.path] = true
		out = append(out, p)
	}
	for _, m := range bodyReferenceRe.FindAllString(body, -1) {
		if !strings.Contains(m, "{") {
			add(deferredLoadPoint{path: filepath.Join(foSkillDir, m), named: m})
		}
	}
	for _, m := range bodySkillRe.FindAllStringSubmatch(body, -1) {
		name := m[1]
		if lazyLoadSkills[name] {
			add(deferredLoadPoint{path: filepath.Join("skills", name, "SKILL.md"), named: m[0]})
		}
	}
	for _, m := range bodyModRe.FindAllStringSubmatch(body, -1) {
		add(deferredLoadPoint{path: filepath.Join("mods", m[1]), named: m[0]})
	}
	return out
}

func TestBootResidentDeferredLoadPointsResolve(t *testing.T) {
	root := repoRoot(t)
	total := 0
	walked := append(append([]string{}, bootResidentBodies...), deferredModuleBodies...)
	for _, rel := range walked {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("read boot-resident body %s: %v", rel, err)
		}
		for _, p := range extractDeferredLoadPoints(string(data)) {
			total++
			if _, err := os.Stat(filepath.Join(root, p.path)); err != nil {
				t.Errorf("%s names deferred load-point %q which does not resolve at %s: %v", rel, p.named, p.path, err)
			}
		}
	}
	if total == 0 {
		t.Fatal("extracted zero deferred load-points")
	}
}

func TestMovedSectionsLiveBehindANamedTrigger(t *testing.T) {
	root := repoRoot(t)
	headingRe := func(heading string) *regexp.Regexp {
		return regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(heading) + `\s*$`)
	}
	sharedRel := filepath.Join("skills", "first-officer", "references", "first-officer-shared-core.md")
	shared, err := os.ReadFile(filepath.Join(root, sharedRel))
	if err != nil {
		t.Fatalf("read shared core: %v", err)
	}
	block := deferredLoadPointsBlock(t, string(shared))

	for _, mv := range movedSections {
		src, err := os.ReadFile(filepath.Join(root, mv.from))
		if err != nil {
			t.Fatalf("read move source %s: %v", mv.from, err)
		}
		if got := len(headingRe(mv.heading).FindAllString(string(src), -1)); got != 0 {
			t.Errorf("%q still appears %d times in %s; a moved section must not stay boot-resident", mv.heading, got, mv.from)
		}

		dst, err := os.ReadFile(filepath.Join(root, mv.to))
		if err != nil {
			t.Fatalf("read move target %s: %v", mv.to, err)
		}
		if got := len(headingRe(mv.heading).FindAllString(string(dst), -1)); got != 1 {
			t.Errorf("%q appears %d times in %s, want exactly once", mv.heading, got, mv.to)
		}

		// The target must be reachable: the shared core has to name it as a
		// deferred load point, or nothing loads the section at its trigger.
		named := filepath.ToSlash(strings.TrimPrefix(filepath.ToSlash(mv.to), "skills/first-officer/"))
		if !strings.Contains(block, named) {
			t.Errorf("deferred load-points block does not name %q, so %q is unreachable at its trigger", named, mv.heading)
		}
	}
}

func TestBootResidentDeferredLoadPointGuardFailsOnDanglingTarget(t *testing.T) {
	root := repoRoot(t)
	fixture := "read references/claude-fo-this-file-does-not-exist.md\n" +
		"invoke spacedock:fo-dispatch-recovery.\n"
	points := extractDeferredLoadPoints(fixture)
	if len(points) == 0 {
		t.Fatal("control fixture extracted no load-points")
	}
	var sawDangling, sawReal bool
	for _, p := range points {
		_, err := os.Stat(filepath.Join(root, p.path))
		if strings.Contains(p.named, "does-not-exist") {
			if err == nil {
				t.Fatalf("dangling reference %q unexpectedly resolved", p.named)
			}
			sawDangling = true
		} else if err == nil {
			sawReal = true
		}
	}
	if !sawDangling || !sawReal {
		t.Fatalf("closure control missing contrast: dangling=%v real=%v", sawDangling, sawReal)
	}
}
