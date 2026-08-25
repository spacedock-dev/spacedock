// ABOUTME: AC-2 — the dispatch build artifact body carries the stage-report
// ABOUTME: protocol template for host=pi, and omits it for claude and codex.
package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuildPiArtifactCarriesStageReportProtocol (AC-2) builds an artifact with
// a non-self-describing checklist (one that does NOT mention the ensign skill
// path, the stage-report heading, or the DONE/Summary structure) for host=pi
// and asserts the generated body carries the protocol tokens: the
// `## Stage Report:` template heading, the `- DONE:`/`- SKIPPED:`/`- FAILED:`
// markers, and `### Summary`. The same test asserts host=claude and host=codex
// artifacts do NOT carry the embedded `### Stage Report format` block — the
// embed is Pi-only because Claude's Skill() and Codex's $spacedock:ensign
// bootstrap already supply the format.
func TestBuildPiArtifactCarriesStageReportProtocol(t *testing.T) {
	// A non-self-describing checklist: the entity's real acceptance criteria,
	// with no skill-path, heading, or format hints.
	checklist := []string{
		"- commit the deliverable on the worktree branch",
		"- run go test ./... green",
	}

	for _, host := range []string{"pi", "claude", "codex"} {
		t.Run(host, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, filepath.Join(root, "README.md"), readmeWorktree(false))
			worktreeRel := ".worktrees/spacedock-ensign-stage-report"
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
				"checklist":      checklist,
				"bare_mode":      false,
				"host":           host,
			}, nil)

			native := runNative(stdin, "build", "--workflow-dir", root)
			if native.exit != 0 {
				t.Fatalf("build exit=%d stderr=%q", native.exit, native.stderr)
			}
			body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

			if host == "pi" {
				for _, want := range []string{
					"## Stage Report:",
					"- DONE:",
					"- SKIPPED:",
					"- FAILED:",
					"### Summary",
					"### Stage Report format",
				} {
					if !strings.Contains(body, want) {
						t.Fatalf("pi dispatch body missing protocol token %q:\n%s", want, body)
					}
				}
			} else {
				// Claude and Codex must NOT carry the embedded block; the skill
				// supplies the format for those hosts.
				for _, banned := range []string{
					"### Stage Report format",
					"## Stage Report: {stage}",
				} {
					if strings.Contains(body, banned) {
						t.Fatalf("%s dispatch body must not carry the embedded stage-report block (token %q):\n%s", host, banned, body)
					}
				}
			}
		})
	}
}

// TestBuildPiFirstActionNarrowedToStageReportFormat (AC-1) asserts the Pi
// First-action claim no longer overclaims the full ensign discipline (polling,
// worktree ownership, completion protocol) and instead attributes the
// stage-report format to the body and the rest to the ensign skill.
func TestBuildPiFirstActionNarrowedToStageReportFormat(t *testing.T) {
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

	// The overclaim is gone.
	for _, overclaim := range []string{
		"This file contains the shared ensign discipline entry points",
	} {
		if strings.Contains(body, overclaim) {
			t.Fatalf("pi First-action still overclaims: %q present in body:\n%s", overclaim, body)
		}
	}
	// The narrowed claim attributes the format template to the body and the
	// rest of the discipline to the ensign skill.
	for _, want := range []string{
		"This file carries the stage-report format template",
		"The ensign skill supplies the remaining shared discipline",
		"not auto-loaded",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("pi First-action missing narrowed claim %q:\n%s", want, body)
		}
	}
}
