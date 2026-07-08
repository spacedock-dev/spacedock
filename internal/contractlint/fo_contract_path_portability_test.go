// ABOUTME: Structural guard for workflow-portable FO contract paths.
// ABOUTME: Distinguishes operational instructions from explicit examples and install hints.
package contractlint

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type foContractPathLeak struct {
	Line int
	Text string
	Kind string
}

func TestFOContractPathPortabilityDiscriminators(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		wantLeak bool
	}{
		{
			name:     "operational docs dev README rule fails",
			body:     "| allowed-process | `docs/dev/README.md`; `{workflow_dir}/README.md` | The FO may edit this process doc. |",
			wantLeak: true,
		},
		{
			name:     "operational docs dev mods rule fails",
			body:     "| blocked-product | `docs/dev/_mods/**` | Mods go through workers. |",
			wantLeak: true,
		},
		{
			name:     "operational universal state checkout rule fails",
			body:     "| allowed-state | `.spacedock-state/**`; `{workflow_dir}/_archive/**` | State writes. |",
			wantLeak: true,
		},
		{
			name:     "explicit state checkout example passes",
			body:     "Example: a workflow may declare `state: .spacedock-state`; use the resolved state checkout, not this literal.",
			wantLeak: false,
		},
		{
			name:     "survey discovery signal passes",
			body:     "Survey discovery signal: probe `.spacedock-state`, `docs/**/.spacedock-state`, and `_mods` as possible workflow hints.",
			wantLeak: false,
		},
		{
			name:     "operational local main drift repo rebuild fails",
			body:     "- **local-main-drift** -> `git -C {repo} fetch origin {drift.trunk} && git -C {repo} merge --ff-only origin/{drift.trunk} && cd {repo} && go build -o spacedock ./cmd/spacedock`.",
			wantLeak: true,
		},
		{
			name:     "source build install hint passes",
			body:     "Install hint: when the launcher binary is absent, build it from a Spacedock source checkout with `go build -o spacedock ./cmd/spacedock`.",
			wantLeak: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			leaks := scanFOContractPathLeaks(tc.name, tc.body)
			if tc.wantLeak && len(leaks) == 0 {
				t.Fatal("expected a path portability leak, got none")
			}
			if !tc.wantLeak && len(leaks) > 0 {
				t.Fatalf("unexpected path portability leaks: %v", leaks)
			}
		})
	}
}

func TestShippedFOContractPathsAreWorkflowPortable(t *testing.T) {
	files := foContractPathPortabilityFiles(t)
	if len(files) == 0 {
		t.Fatal("scanned zero FO contract files; guard would pass vacuously")
	}
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		for _, leak := range scanFOContractPathLeaks(path, string(data)) {
			rel, _ := filepath.Rel(repoRoot(t), path)
			t.Errorf("%s:%d: %s: %s", rel, leak.Line, leak.Kind, strings.TrimSpace(leak.Text))
		}
	}
}

func foContractPathPortabilityFiles(t *testing.T) []string {
	t.Helper()
	root := repoRoot(t)
	files := []string{
		filepath.Join(root, "skills", "first-officer", "SKILL.md"),
		filepath.Join(root, "skills", "fo-dispatch-recovery", "SKILL.md"),
		filepath.Join(root, "skills", "fo-status-viewer", "SKILL.md"),
		filepath.Join(root, "skills", "fo-write-core", "SKILL.md"),
		filepath.Join(root, "skills", "feedback-rejection-flow", "SKILL.md"),
		filepath.Join(root, "skills", "present-gate", "SKILL.md"),
		filepath.Join(root, "skills", "using-legacy-claude-team", "SKILL.md"),
	}
	refs, err := filepath.Glob(filepath.Join(root, "skills", "first-officer", "references", "*.md"))
	if err != nil {
		t.Fatalf("glob FO references: %v", err)
	}
	files = append(files, refs...)
	sort.Strings(files)
	return files
}

func scanFOContractPathLeaks(name, body string) []foContractPathLeak {
	lines := strings.Split(body, "\n")
	var leaks []foContractPathLeak
	for i, line := range lines {
		context := localFOContractContext(lines, i)
		for _, check := range []struct {
			token string
			kind  string
		}{
			{token: "docs/dev/README.md", kind: "hardcoded Spacedock workflow README"},
			{token: "docs/dev/_mods", kind: "hardcoded Spacedock workflow mod path"},
			{token: ".spacedock-state/**", kind: "universal state checkout spelling"},
		} {
			if strings.Contains(line, check.token) && !allowedFOContractExampleContext(context) {
				leaks = append(leaks, foContractPathLeak{Line: i + 1, Text: line, Kind: check.kind})
			}
		}
		if strings.Contains(line, "cd {repo} && go build -o spacedock ./cmd/spacedock") &&
			!allowedFOLauncherBuildContext(context) {
			leaks = append(leaks, foContractPathLeak{Line: i + 1, Text: line, Kind: "managed-repo launcher rebuild"})
		}
	}
	return leaks
}

func localFOContractContext(lines []string, i int) string {
	start := i - 3
	if start < 0 {
		start = 0
	}
	end := i + 4
	if end > len(lines) {
		end = len(lines)
	}
	return strings.ToLower(strings.Join(lines[start:end], "\n"))
}

func allowedFOContractExampleContext(context string) bool {
	for _, marker := range []string{
		"example",
		"placeholder",
		"install hint",
		"source-build",
		"source build",
		"discovery signal",
		"probe",
		"survey",
	} {
		if strings.Contains(context, marker) {
			return true
		}
	}
	return false
}

func allowedFOLauncherBuildContext(context string) bool {
	return strings.Contains(context, "install hint") ||
		strings.Contains(context, "source-build") ||
		strings.Contains(context, "source build")
}
