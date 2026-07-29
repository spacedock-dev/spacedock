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
