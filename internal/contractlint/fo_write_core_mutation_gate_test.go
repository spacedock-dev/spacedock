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
	classifier := parseFOWriteClassifierTable(t, table)

	cases := []struct {
		name string
		ctx  foWriteWorkflowContext
		path string
		want string
	}{
		{
			name: "resolved synthetic state entity",
			ctx:  syntheticFOWriteContext("/tmp/acme-flow", "/tmp/acme-flow/.state"),
			path: "/tmp/acme-flow/.state/task/index.md",
			want: "allowed-state",
		},
		{
			name: "resolved synthetic archive root",
			ctx:  syntheticFOWriteContext("/tmp/acme-flow", "/tmp/acme-flow/.state"),
			path: "/tmp/acme-flow/.state/_archive/done/index.md",
			want: "allowed-state",
		},
		{
			name: "absolute workflow README",
			ctx:  syntheticFOWriteContext("/tmp/acme-flow", "/tmp/acme-flow/.state"),
			path: "/tmp/acme-flow/README.md",
			want: "allowed-process",
		},
		{
			name: "relative workflow README",
			ctx:  syntheticFOWriteContext("workflows/acme", "workflows/acme/.state"),
			path: "workflows/acme/README.md",
			want: "allowed-process",
		},
		{
			name: "docs dev README is not special outside the discovered workflow",
			ctx:  syntheticFOWriteContext("/tmp/acme-flow", "/tmp/acme-flow/.state"),
			path: "docs/dev/README.md",
			want: "blocked-product",
		},
		{
			name: "registered mods are product work",
			ctx:  syntheticFOWriteContext("/tmp/acme-flow", "/tmp/acme-flow/.state"),
			path: "/tmp/acme-flow/_mods/pr-merge.md",
			want: "blocked-product",
		},
		{
			name: "unmatched code defaults to product",
			ctx:  syntheticFOWriteContext("/tmp/acme-flow", "/tmp/acme-flow/.state"),
			path: "internal/status/mutate.go",
			want: "blocked-product",
		},
		{
			name: "shipped skill scaffolding defaults to product",
			ctx:  syntheticFOWriteContext("/tmp/acme-flow", "/tmp/acme-flow/.state"),
			path: "skills/fo-write-core/SKILL.md",
			want: "blocked-product",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifyFOWriteTarget(classifier, tc.ctx, tc.path); got != tc.want {
				t.Errorf("classify %q = %q, want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestFOWriteCoreMutationGateRequiresExactOverride(t *testing.T) {
	table := readFOWriteClassifierTable(t)
	classifier := parseFOWriteClassifierTable(t, table)
	ctx := syntheticFOWriteContext("/tmp/acme-flow", "/tmp/acme-flow/.state")

	target := "internal/status/mutate.go"
	if foWriteOverrideAllows(target, "you may fix the code directly") {
		t.Fatalf("broad direct-edit text unexpectedly allowed %q", target)
	}
	if foWriteMayWrite(classifier, ctx, target, "you may fix the code directly") {
		t.Fatalf("broad direct-edit text must not allow blocked product target %q", target)
	}
	if !foWriteMayWrite(classifier, ctx, target, "you may directly edit internal/status/mutate.go for this task") {
		t.Fatalf("exact target grant must allow %q after blocked-product classification", target)
	}
	if foWriteMayWrite(classifier, ctx, "internal/status/parse.go", "you may directly edit internal/status/mutate.go for this task") {
		t.Fatalf("exact grant for mutate.go must not allow a different product path")
	}
}

func readFOWriteClassifierTable(t *testing.T) string {
	t.Helper()
	path := filepath.Join(repoRoot(t), "skills", "fo-write-core", "SKILL.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	body := string(data)
	start := strings.Index(body, foWriteClassifierStart)
	end := strings.Index(body, foWriteClassifierEnd)
	if start < 0 || end < 0 || end <= start {
		t.Fatalf("skills/fo-write-core/SKILL.md must carry a machine-readable classifier block bounded by %s and %s", foWriteClassifierStart, foWriteClassifierEnd)
	}
	return body[start+len(foWriteClassifierStart) : end]
}

func parseFOWriteClassifierTable(t *testing.T, table string) map[string][]string {
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
				patterns = append(patterns, pattern)
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

type foWriteWorkflowContext struct {
	workflowDir   string
	stateCheckout string
	entityPath    string
	archiveRoots  []string
	modRoots      []string
}

func syntheticFOWriteContext(workflowDir, stateCheckout string) foWriteWorkflowContext {
	return foWriteWorkflowContext{
		workflowDir:   cleanFOWritePath(workflowDir),
		stateCheckout: cleanFOWritePath(stateCheckout),
		entityPath:    cleanFOWritePath(filepath.Join(stateCheckout, "task", "index.md")),
		archiveRoots: []string{
			cleanFOWritePath(filepath.Join(stateCheckout, "_archive")),
		},
		modRoots: []string{
			cleanFOWritePath(filepath.Join(workflowDir, "_mods")),
		},
	}
}

func classifyFOWriteTarget(classifier map[string][]string, ctx foWriteWorkflowContext, target string) string {
	target = cleanFOWritePath(target)
	if classifierHasSource(classifier, "blocked-product", "registered mods plus every target not classified as state/process") &&
		anyFOWritePathPrefix(target, ctx.modRoots) {
		return "blocked-product"
	}
	if classifierHasSource(classifier, "allowed-state", "resolved state/entity/archive paths") &&
		(target == ctx.entityPath || hasFOWritePathPrefix(target, ctx.stateCheckout) || anyFOWritePathPrefix(target, ctx.archiveRoots)) {
		return "allowed-state"
	}
	if classifierHasSource(classifier, "allowed-process", "{workflow_dir}/README.md only") &&
		target == cleanFOWritePath(filepath.Join(ctx.workflowDir, "README.md")) {
		return "allowed-process"
	}
	return "blocked-product"
}

func classifierHasSource(classifier map[string][]string, class, want string) bool {
	for _, source := range classifier[class] {
		if source == want {
			return true
		}
	}
	return false
}

func cleanFOWritePath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

func anyFOWritePathPrefix(target string, roots []string) bool {
	for _, root := range roots {
		if hasFOWritePathPrefix(target, root) {
			return true
		}
	}
	return false
}

func hasFOWritePathPrefix(target, root string) bool {
	root = strings.TrimSuffix(cleanFOWritePath(root), "/")
	return target == root || strings.HasPrefix(target, root+"/")
}

func foWriteMayWrite(classifier map[string][]string, ctx foWriteWorkflowContext, target, grant string) bool {
	switch classifyFOWriteTarget(classifier, ctx, target) {
	case "allowed-state", "allowed-process":
		return true
	default:
		return foWriteOverrideAllows(target, grant)
	}
}

func foWriteOverrideAllows(target, grant string) bool {
	return strings.Contains(grant, target)
}
