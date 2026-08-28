// ABOUTME: AC-2 — the dispatch build artifact body carries the stage-report
// ABOUTME: protocol template for host=pi, and omits it for claude and codex.
package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPiFirstActionInvokesEnsignSkill(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
	worktreeRel := ".worktrees/spacedock-ensign-first-action"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}
	entityPath := filepath.Join(root, "thing.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))
	gitInit(t, root)

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   root,
		"stage":          "implementation",
		"checklist":      []string{"- a"},
		"bare_mode":      false,
		"host":           "pi",
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", root)
	if native.exit != 0 {
		t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
	}
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	// The false claim that the dispatch file itself carries the ensign
	// discipline entry points must be gone.
	for _, banned := range []string{
		"This file contains the shared ensign discipline entry points",
	} {
		if strings.Contains(body, banned) {
			t.Fatalf("pi First-action still carries false claim %q:\n%s", banned, body)
		}
	}
	// The worker must be told to load the ensign skill before reading the
	// dispatch file.
	hasSkillLoad := strings.Contains(body, "/skill:ensign") ||
		strings.Contains(body, "skills/ensign/SKILL.md")
	if !hasSkillLoad {
		t.Fatalf("pi First-action missing ensign skill-load instruction (/skill:ensign or skills/ensign/SKILL.md):\n%s", body)
	}
	// The skill-load must come before the instruction to read the dispatch
	// file, mirroring Claude and Codex.
	skillIdx := strings.Index(body, "/skill:ensign")
	if skillIdx < 0 {
		skillIdx = strings.Index(body, "skills/ensign/SKILL.md")
	}
	readIdx := strings.Index(body, "read this dispatch file")
	if readIdx < 0 {
		t.Fatalf("pi First-action missing 'read this dispatch file' instruction:\n%s", body)
	}
	if skillIdx >= readIdx {
		t.Fatalf("pi First-action: ensign skill-load must precede 'read this dispatch file' (skillIdx=%d readIdx=%d):\n%s", skillIdx, readIdx, body)
	}
}
