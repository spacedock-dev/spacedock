// ABOUTME: AC-6 — a two-position structural lint asserting no .github/workflows
// ABOUTME: YAML resolves the integration trunk to `next` post-flip (PR-base + gh-run).
package contractlint

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"gopkg.in/yaml.v3"
)

// CI workflow YAML is parsed by GitHub Actions, NOT ingested by the model as
// instruction prose, so a structural assertion over it is a legitimate non-grep
// oracle (the same non-model-read exemption structural_checks_test.go relies on).
// This lint checks the TWO specific positions where a workflow resolves the
// integration trunk — NOT a blanket `next` token sweep — so it does not
// false-positive on the legitimately-`next` uses (next-publish.yml's edge-channel
// publish, the awk `next` keyword, an @next action pin, `--ref next`, or English
// in comments).

// ghRunListNextRe matches a producer-run query pinned to the `next` branch:
// `gh run list … --branch next`. The `\b` after next keeps `next-foo` from
// matching. Position (b) of the two-position check.
var ghRunListNextRe = regexp.MustCompile(`gh run list[^\n]*--branch\s+next\b`)

// workflowTrunkOffenses returns, for a workflow file, the human-readable reasons
// it resolves the integration trunk to `next` — one per offending position. An
// empty slice means the file is clean. The two positions:
//
//	(a) on.pull_request.branches lists `next` (the PR-base filter), found by a
//	    yaml.Node mapping walk (yaml.v3 preserves `on` as a string key — it does
//	    NOT apply the YAML-1.1 `on`→`true` bool remap — but the walk also tolerates
//	    a bool-remapped key defensively).
//	(b) a `gh run list … --branch next` query in the file body (regex).
func workflowTrunkOffenses(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var offenses []string

	// (a) structural: on.pull_request.branches must not list `next`.
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	if len(doc.Content) > 0 {
		root := doc.Content[0]
		on := mappingValue(root, isOnKey)
		if on != nil {
			pr := mappingValue(on, func(k string) bool { return k == "pull_request" })
			if pr != nil {
				branches := mappingValue(pr, func(k string) bool { return k == "branches" })
				if branches != nil && sequenceHasScalar(branches, "next") {
					offenses = append(offenses, "on.pull_request.branches lists `next`")
				}
			}
		}
	}

	// (b) command: no `gh run list … --branch next` producer-run query.
	if ghRunListNextRe.MatchString(string(data)) {
		offenses = append(offenses, "`gh run list … --branch next` query targets the pre-flip trunk")
	}
	return offenses
}

// mappingValue returns the value node for the first key in a mapping node whose
// scalar value satisfies match, or nil. Non-mapping nodes yield nil.
func mappingValue(node *yaml.Node, match func(key string) bool) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		k := node.Content[i]
		if k.Kind == yaml.ScalarNode && match(k.Value) {
			return node.Content[i+1]
		}
	}
	return nil
}

// isOnKey matches the workflow trigger key. yaml.v3 keeps `on` as the string
// "on", but a YAML-1.1 parser would remap the bare `on:` to the boolean true; we
// tolerate both the "on" and the "true" rendering so the walk finds the trigger
// block regardless of parser version.
func isOnKey(k string) bool { return k == "on" || k == "true" }

// sequenceHasScalar reports whether a sequence node contains a scalar equal to
// want (e.g. `branches: [next, main]` contains "next").
func sequenceHasScalar(node *yaml.Node, want string) bool {
	if node == nil || node.Kind != yaml.SequenceNode {
		return false
	}
	for _, item := range node.Content {
		if item.Kind == yaml.ScalarNode && item.Value == want {
			return true
		}
	}
	return false
}

// workflowYAMLFiles returns every .github/workflows/*.yml file.
func workflowYAMLFiles(t *testing.T) []string {
	t.Helper()
	dir := filepath.Join(repoRoot(t), ".github", "workflows")
	matches, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		t.Fatalf("glob workflows: %v", err)
	}
	return matches
}

// TestNoWorkflowResolvesTrunkToNext (AC-6) is a STANDING structural guard: no
// .github/workflows YAML resolves the integration trunk to `next` post-flip, in
// either of the two trunk-resolving positions. It reds on a re-introduced PR-base
// `next` filter or `gh run list --branch next` query — not a one-time flip.
func TestNoWorkflowResolvesTrunkToNext(t *testing.T) {
	files := workflowYAMLFiles(t)
	if len(files) == 0 {
		t.Fatal("walked zero workflow YAML files — scope bug; the guard would pass vacuously")
	}
	for _, path := range files {
		rel := filepath.Base(path)
		for _, offense := range workflowTrunkOffenses(t, path) {
			t.Errorf("%s resolves the integration trunk to next: %s", rel, offense)
		}
	}
}

// TestWorkflowTrunkLintDiscriminates is the DISCRIMINATOR control proving the
// two-position lint is precise — able to flag the real offending positions yet
// silent on the legitimately-`next` uses. Without this, a lint that flagged
// nothing (or flagged everything) would pass the absence check vacuously.
func TestWorkflowTrunkLintDiscriminates(t *testing.T) {
	dir := t.TempDir()

	// Positive: a workflow with BOTH offending positions is flagged twice.
	offender := filepath.Join(dir, "offender.yml")
	writeWorkflow(t, offender, `name: offender
on:
  pull_request:
    branches: [next, main]
jobs:
  ledger:
    runs-on: ubuntu-latest
    steps:
      - run: |
          run_id="$(gh run list --workflow "X" --branch next --status success --limit 1 || true)"
`)
	if got := workflowTrunkOffenses(t, offender); len(got) != 2 {
		t.Errorf("offender flagged %d positions, want 2 (PR-base + gh-run): %v", len(got), got)
	}

	// Negative: the legitimate `next` uses must NOT be flagged — an edge-channel
	// publish (push to next), an @next action pin, a `--ref next` checkout, the awk
	// `next` keyword, and English mentioning ref `next` in comments/strings.
	legit := filepath.Join(dir, "legit.yml")
	writeWorkflow(t, legit, `name: next-publish
on:
  push:
    branches: [main]
jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v5
        with:
          ref: next
      - uses: some/action@next
      - run: git push origin next
      - run: git cat-file tag "$X" | awk 'p{print; next} /^$/{p=1}' > notes.txt
      - run: echo "must not install from remote next." >&2
`)
	if got := workflowTrunkOffenses(t, legit); len(got) != 0 {
		t.Errorf("legitimate-next workflow false-positived: %v", got)
	}
}

// writeWorkflow writes a workflow YAML fixture under a temp dir.
func writeWorkflow(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
