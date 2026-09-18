package dispatch

import (
	"github.com/spacedock-dev/spacedock/internal/claudeteam"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSemanticNames(t *testing.T) {
	if name, err := semanticName("a", "done"); err != nil || name != "a-done" {
		t.Fatalf("single-character slug: %q %v", name, err)
	}
	name, err := semanticName("ci-duration-hints", "ideation")
	if err != nil || name != "ci-duration-hints-ideation" {
		t.Fatalf("%q %v", name, err)
	}
	a, _ := semanticName(longSlug, "implementation")
	b, _ := semanticName(longSlugShare, "implementation")
	if a == b || len(a) > 56 || !strings.HasPrefix(a, "dispatch-") || !strings.HasSuffix(a, "-implementation") {
		t.Fatal(a, b)
	}
	for _, pair := range [][2]string{{"Bad", "done"}, {longSlug, strings.Repeat("x", 40)}} {
		if _, err := semanticName(pair[0], pair[1]); err == nil {
			t.Fatal(pair)
		}
	}
	entities := map[string]entityRecord{"alpha": {slug: "alpha"}, "alpha-review": {slug: "alpha-review"}, "spacedock-ensign-alpha": {slug: "spacedock-ensign-alpha"}}
	for _, name := range []string{"alpha-review-done", "spacedock-ensign-alpha-done", "alpha-review-2"} {
		if resolveWorker(name, []string{"done", "review-done", "review", "review-2"}, entities).ok {
			t.Fatal("ambiguous", name)
		}
	}
	for _, style := range []string{"sd-b32", "sequential", "slug"} {
		rec := entityRecord{slug: longSlug, id: idAlpha}
		if style == "sequential" {
			rec.id = "001"
		}
		entities := map[string]entityRecord{longSlug: rec}
		legacy := map[string]string{"sd-b32": "spacedock-ensign-367s2zrbkm-implementation", "sequential": "spacedock-ensign-001-implementation", "slug": "spacedock-ensign-dispatch-reconcile-deconf-implementation"}
		for _, name := range []string{a, legacy[style]} {
			d := resolveWorker(name+"-cycle2", []string{"implementation"}, entities)
			if !d.ok || d.slug != longSlug || d.cycle != "cycle2" {
				t.Fatalf("%s: %+v", name, d)
			}
		}
	}
	if validWorkerSuffix("-cycle10") == false || validWorkerSuffix("-cycle100") {
		t.Fatal("suffix budget")
	}
}

func TestSemanticCohorts(t *testing.T) {
	records := map[string]entityRecord{"thing": {slug: "thing", status: "done"}}
	members := []claudeteam.ReconcileMember{{Name: "spacedock-ensign-thing-done"}, {Name: "thing-done-cycle2"}}
	drift := classB(members, []string{"done"}, records, nil)
	if len(drift) != 1 || drift[0].Name != members[0].Name {
		t.Fatal(drift)
	}
	members = append(members, claudeteam.ReconcileMember{Name: "spacedock-ensign-thing-done-2"})
	if drift := classB(members, []string{"done"}, records, nil); len(drift) != 0 {
		t.Fatal("equal-cycle tie", drift)
	}
	if drift := classA(members, []string{"done"}, records, nil); len(drift) != 3 {
		t.Fatal("lost original handles", drift)
	}
}

func TestSemanticCandidateCollision(t *testing.T) {
	// A crafted equal digest candidate must not pick the first entity.
	candidates := map[decomposeResult]bool{{slug: "long-alpha", stage: "done", ok: true}: true, {slug: "long-bravo", stage: "done", ok: true}: true}
	if uniqueWorker(candidates).ok {
		t.Fatal("digest collision acquired an owner")
	}
}

func TestDispatchNameReadOnly(t *testing.T) {
	for _, slug := range []string{"ci-duration-hints", longSlug, longSlugShare} {
		t.Run(slug, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, filepath.Join(root, "README.md"), readmeIDStyle("slug", false))
			ep := filepath.Join(root, slug+".md")
			writeFile(t, ep, entityFMID(idAlpha, "Naming", "backlog"))
			gitInit(t, root)
			args := []string{"name", "--workflow-dir", root, "--entity-path", ep, "--stage", "backlog"}
			got := runNative("", args...)
			if got.exit != 0 {
				t.Fatal(got.stderr)
			}
			name := strings.TrimSpace(got.stdout)
			if slug == "ci-duration-hints" && name != "ci-duration-hints-backlog" {
				t.Fatal(name)
			}
			if len(name) > 56 || !strings.HasSuffix(name, "-backlog") {
				t.Fatal(name)
			}
			if gitOutput(t, root, "status", "--porcelain") != "" {
				t.Fatal("name query mutated fixture")
			}
			built := runNative("- work", "build", "--workflow-dir", root, "--entity-path", ep, "--stage", "backlog", "--checklist-file", "-")
			if built.exit != 0 || nameFromStdout(t, built.stdout) != name {
				t.Fatal("name/build disagree", got, built)
			}
		})
	}
}

func TestDispatchNameRefusesUnsafeIdentity(t *testing.T) {
	for _, kind := range []string{"missing-entity", "missing-workflow", "unknown-stage", "invalid", "ambiguous", "budget", "worktree", "missing-flag", "unknown-flag"} {
		t.Run(kind, func(t *testing.T) {
			root, ep := buildHostFixture(t)
			stage := "backlog"
			switch kind {
			case "missing-entity":
				ep = filepath.Join(root, "absent.md")
			case "missing-workflow":
				os.Remove(filepath.Join(root, "README.md"))
			case "unknown-stage":
				stage = "absent"
			case "invalid":
				ep = filepath.Join(root, "Bad.md")
				writeFile(t, ep, entityFM("Bad", "backlog", ""))
			case "ambiguous":
				ep = filepath.Join(root, "spacedock-ensign-thing.md")
				writeFile(t, ep, entityFM("Alias", "backlog", ""))
			case "budget":
				stage = strings.Repeat("x", 40)
				writeFile(t, filepath.Join(root, "README.md"), strings.ReplaceAll(readmeIDStyle("slug", false), "backlog", stage))
				ep = filepath.Join(root, longSlug+".md")
				writeFile(t, ep, entityFM("Long", stage, ""))
			case "worktree":
				ep = filepath.Join(root, ".worktrees", "copy", "thing.md")
				writeFile(t, ep, entityFM("Copy", "backlog", ""))
			}
			args := []string{"name", "--workflow-dir", root, "--entity-path", ep, "--stage", stage}
			if kind == "missing-flag" {
				args = args[:len(args)-2]
			}
			if kind == "unknown-flag" {
				args = append(args, "--stamp")
			}
			before := gitOutput(t, root, "status", "--porcelain")
			head := gitOutput(t, root, "rev-parse", "HEAD")
			got := runNative("", args...)
			if got.exit == 0 || got.stdout != "" || got.stderr == "" {
				t.Fatal("unsafe query accepted", got)
			}
			if gitOutput(t, root, "status", "--porcelain") != before || gitOutput(t, root, "rev-parse", "HEAD") != head {
				t.Fatal("refusal mutated fixture")
			}
		})
	}
}
