// ABOUTME: --validate warns on a flat entity holding gate rooms, carrying the
// ABOUTME: conversion remedy, while the plain status read path stays unaffected.
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// hybridFlatRoomsFixture plants two flat entities holding prepared rooms and one
// clean folder-form entity, so the finding is exercised against both the shape
// it must flag and the shape it must not.
func hybridFlatRoomsFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	readme := "---\nid-style: slug\nstages:\n  states:\n    - name: ideation\n      initial: true\n      gate: true\n    - name: implementation\n---\n# Workflow\n"
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	body := func(slug, ref string) string {
		return "---\nid: " + slug + "\nstatus: ideation\ntitle: Task\ngates:\n  version: 1\n  records:\n    - id: gate:" + slug +
			":ideation\n      stage: ideation\n      attempts:\n        - id: attempt:" + slug +
			":ideation\n          briefing: {id: briefing:" + slug + ":ideation:attempt-1:revision-1, digest: sha256:" +
			strings.Repeat("1", 64) + ", room-ref: " + ref + "}\n---\n# Task\n"
	}
	for _, slug := range []string{"alpha", "beta"} {
		if err := os.WriteFile(filepath.Join(root, slug+".md"),
			[]byte(body(slug, "./"+slug+"/review/ideation/briefing-1")), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(root, slug, "review", "ideation", "briefing-1"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "gamma", "review", "ideation", "briefing-1"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "gamma", "index.md"),
		[]byte(body("gamma", "./review/ideation/briefing-1")), 0o644); err != nil {
		t.Fatal(err)
	}
	gitInit(t, root)
	return root
}

// The remedy is the payload: a reader who only sees "this is a hybrid" still
// hand-converts and destroys every retained room, so the line must name the move
// AND the ref rewrite. Asserting both tokens fails if either half is dropped.
func TestValidateWarnsFlatEntityHoldingGateRoomsWithConversionRemedy(t *testing.T) {
	root := hybridFlatRoomsFixture(t)
	out, stderr, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--validate")
	if code != 0 || strings.TrimSpace(out) != "VALID" {
		t.Fatalf("--validate exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if got := strings.Count(stderr, "flat entity holds gate rooms"); got != 2 {
		t.Fatalf("want 2 hybrid warnings, got %d: %q", got, stderr)
	}
	for _, want := range []string{
		"slug=alpha", "slug=beta",
		"git mv alpha.md alpha/index.md", "rewrite every `room-ref: ./alpha/` to `room-ref: ./`",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("hybrid warning missing %q: %q", want, stderr)
		}
	}
	if strings.Contains(stderr, "slug=gamma") {
		t.Fatalf("folder-form entity flagged: %q", stderr)
	}
}

// The finding must not live in findEntityFormConflicts: that runs on the read
// path too, where an error exits 1 and locks the first officer out of the very
// listing that shows the broken entity.
func TestHybridFindingLeavesPlainStatusReadPathUnaffected(t *testing.T) {
	root := hybridFlatRoomsFixture(t)
	out, stderr, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root)
	if code != 0 {
		t.Fatalf("plain status exit=%d stdout=%q stderr=%q", code, out, stderr)
	}
	if strings.Contains(stderr, "flat entity holds gate rooms") || strings.Contains(stderr, "flat/folder conflict") {
		t.Fatalf("read path emitted a hybrid finding: %q", stderr)
	}
	for _, slug := range []string{"alpha", "beta", "gamma"} {
		if !strings.Contains(out, slug) {
			t.Fatalf("read path dropped %s: %q", slug, out)
		}
	}
}

// The residual grandfathering leaves is an operator who converts by hand and
// cannot tell whether the ref rewrite landed. This is the check that answers
// them: the doubled-slug end state of #739 becomes a validate-time finding
// instead of a mid-ceremony gate failure.
func TestValidateReportsRetainedRoomThatNoLongerResolves(t *testing.T) {
	root := hybridFlatRoomsFixture(t)
	// The hand conversion, ref rewrite forgotten: alpha.md -> alpha/index.md
	// leaves ./alpha/review/... resolving at alpha/alpha/review/...
	body, err := os.ReadFile(filepath.Join(root, "alpha.md"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "alpha", "index.md"), body, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, "alpha.md")); err != nil {
		t.Fatal(err)
	}
	_, stderr, code := runNative(t, root, pinnedEnv(t), "--workflow-dir", root, "--validate")
	if code != 0 {
		t.Fatalf("--validate exit=%d stderr=%q", code, stderr)
	}
	if !strings.Contains(stderr, "retained gate room does not resolve: ./alpha/review/ideation/briefing-1") {
		t.Fatalf("botched conversion not reported: %q", stderr)
	}
	// beta is untouched and its rooms resolve; gamma is folder-form and clean.
	// Exactly one unresolved-ref finding, and it is alpha's.
	if got := strings.Count(stderr, "does not resolve"); got != 1 {
		t.Fatalf("want exactly 1 unresolved-ref finding, got %d: %q", got, stderr)
	}
	if strings.Contains(stderr, "slug=gamma") {
		t.Fatalf("false positive on the clean folder entity: %q", stderr)
	}
}
