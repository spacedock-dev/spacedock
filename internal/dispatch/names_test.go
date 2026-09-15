package dispatch

import (
	"github.com/spacedock-dev/spacedock/internal/claudeteam"
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
