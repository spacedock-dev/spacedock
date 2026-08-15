// ABOUTME: Structural guard for fo-write-core's machine-readable mutation gate.
// ABOUTME: The path fixtures are independent of the skill prose, so drift reds.
package contractlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	foWriteClassifierStart = "<!-- FO-WRITE-CLASSIFIER:START -->"
	foWriteClassifierEnd   = "<!-- FO-WRITE-CLASSIFIER:END -->"
)

func TestFOWriteCoreMutationGateClassifiesTargets(t *testing.T) {
	table := readFOWriteClassifierTable(t)

	// Two workflow dirs: a table that only works for the repo it shipped from
	// passes docs/dev and fails ops/release.
	for _, workflowDir := range []string{"docs/dev", "ops/release"} {
		classifier := parseFOWriteClassifierTable(t, table, workflowDir)

		cases := []struct {
			path string
			want string
		}{
			{".spacedock-state/task/index.md", "allowed-state"},
			{workflowDir + "/_archive/shipped-task.md", "allowed-state"},
			{workflowDir + "/README.md", "allowed-process"},
			{"cmd/spacedock/main.go", "blocked-product"},
			{"internal/status/mutate.go", "blocked-product"},
			{"internal/status/mutate_test.go", "blocked-product"},
			{"skills/fo-write-core/SKILL.md", "blocked-product"},
			{"agents/first-officer.md", "blocked-product"},
			{"references/legacy.md", "blocked-product"},
			{"plugin.json", "blocked-product"},
			{".github/workflows/runtime-live-e2e.yml", "blocked-product"},
			{"docs/site/reference/command-reference.md", "blocked-product"},
			{"docs/specs/state-behavior-extension.md", "blocked-product"},
			{"docs/roadmap/0250-fo-behavioral-discipline/index.md", "blocked-product"},
			{"skills/integration/testdata/entity-label-drive/README.md", "blocked-product"},
			{workflowDir + "/_mods/pr-merge.md", "blocked-product"},
		}
		for _, tc := range cases {
			got := classifyFOWriteTarget(classifier, tc.path)
			if got != tc.want {
				t.Errorf("workflow_dir=%q: classify %q = %q, want %q", workflowDir, tc.path, got, tc.want)
				continue
			}
			// classifyFOWriteTarget defaults to blocked-product, so a bare
			// blocked-product expectation is satisfied by a table that never
			// mentions the path at all. Require an explicit pattern match.
			if tc.want == "blocked-product" && !foWriteClassMatches(classifier, "blocked-product", tc.path) {
				t.Errorf("workflow_dir=%q: %q is blocked only by default-deny, not an explicit blocked-product pattern", workflowDir, tc.path)
			}
		}
	}
}

func TestFOWriteCoreMutationGateRequiresExactOverride(t *testing.T) {
	table := readFOWriteClassifierTable(t)
	classifier := parseFOWriteClassifierTable(t, table, "docs/dev")

	target := "internal/status/mutate.go"
	if foWriteOverrideAllows(target, "you may fix the code directly") {
		t.Fatalf("broad direct-edit text unexpectedly allowed %q", target)
	}
	if foWriteMayWrite(classifier, target, "you may fix the code directly") {
		t.Fatalf("broad direct-edit text must not allow blocked product target %q", target)
	}
	if !foWriteMayWrite(classifier, target, "you may directly edit internal/status/mutate.go for this task") {
		t.Fatalf("exact target grant must allow %q after blocked-product classification", target)
	}
	if foWriteMayWrite(classifier, "internal/status/parse.go", "you may directly edit internal/status/mutate.go for this task") {
		t.Fatalf("exact grant for mutate.go must not allow a different product path")
	}
}

func readFOWriteClassifierTable(t *testing.T) string {
	t.Helper()
	path := filepath.Join(repoRoot(t), "skills", "first-officer", "references", "fo-write-core.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	body := string(data)
	start := strings.Index(body, foWriteClassifierStart)
	end := strings.Index(body, foWriteClassifierEnd)
	if start < 0 || end < 0 || end <= start {
		t.Fatalf("canonical fo-write-core reference must carry a machine-readable classifier block bounded by %s and %s", foWriteClassifierStart, foWriteClassifierEnd)
	}
	return body[start+len(foWriteClassifierStart) : end]
}

func parseFOWriteClassifierTable(t *testing.T, table, workflowDir string) map[string][]string {
	t.Helper()
	out := map[string][]string{}
	for _, line := range strings.Split(table, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") || strings.Contains(line, "---") || strings.HasPrefix(strings.ToLower(line), "| class ") {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}
		class := strings.TrimSpace(parts[1])
		if class == "" {
			continue
		}
		var patterns []string
		for _, raw := range strings.Split(parts[2], ";") {
			pattern := strings.Trim(strings.TrimSpace(raw), "`")
			if pattern != "" {
				patterns = append(patterns, strings.ReplaceAll(pattern, "{workflow_dir}", workflowDir))
			}
		}
		out[class] = patterns
	}
	for _, class := range []string{"allowed-state", "allowed-process", "blocked-product", "override"} {
		if len(out[class]) == 0 {
			t.Fatalf("classifier table missing class %q with at least one pattern/rule; parsed=%v", class, out)
		}
	}
	return out
}

func classifyFOWriteTarget(classifier map[string][]string, target string) string {
	for _, class := range []string{"blocked-product", "allowed-state", "allowed-process"} {
		if foWriteClassMatches(classifier, class, target) {
			return class
		}
	}
	return "blocked-product"
}

func foWriteClassMatches(classifier map[string][]string, class, target string) bool {
	for _, pattern := range classifier[class] {
		if pathPatternMatches(pattern, target) {
			return true
		}
	}
	return false
}

func pathPatternMatches(pattern, target string) bool {
	pattern = strings.TrimSpace(pattern)
	switch {
	case strings.HasSuffix(pattern, "/**"):
		return strings.HasPrefix(target, strings.TrimSuffix(pattern, "**"))
	case strings.HasPrefix(pattern, "**/*"):
		return strings.HasSuffix(target, strings.TrimPrefix(pattern, "**/*"))
	default:
		return target == pattern
	}
}

func foWriteMayWrite(classifier map[string][]string, target, grant string) bool {
	switch classifyFOWriteTarget(classifier, target) {
	case "allowed-state", "allowed-process":
		return true
	default:
		return foWriteOverrideAllows(target, grant)
	}
}

func foWriteOverrideAllows(target, grant string) bool {
	return strings.Contains(grant, target)
}
