// ABOUTME: Behavioral coverage for split-root state-commit guidance — both the
// ABOUTME: worktree and non-worktree branches must emit resolved, brace-free paths.
package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestStateCommitPathAbsoluteFromRelativeWorkflowDir pins AC-0: a worktree-stage
// dispatch built with a RELATIVE --workflow-dir must still emit an ABSOLUTE
// `git -C <stateCheckout> add …` state-commit path. The worktree worker runs with
// its cwd inside .worktrees/…, where a relative `git -C docs/dev/.spacedock-state`
// resolves nowhere — so the emitted path must be cwd-independent. runBuild
// absolutizes workflowDir against the process cwd, so the relative flag value
// resolves to an absolute base for every downstream join.
func TestStateCommitPathAbsoluteFromRelativeWorkflowDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()

	// Split-root workflow under root. Drive --workflow-dir as a RELATIVE path
	// computed from the process cwd, so runBuild's filepath.Abs resolution is the
	// thing under test. (t.Chdir is go1.24+; this module targets go1.22, so the
	// relative spelling is derived rather than realized via a cwd change.)
	workflowDir := filepath.Join(root, "wf")
	stateCheckout := filepath.Join(workflowDir, "state-checkout")
	writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(true))

	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}

	entityPath := filepath.Join(stateCheckout, "thing", "index.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))

	gitInit(t, root)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	workflowRel, err := filepath.Rel(cwd, workflowDir)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.IsAbs(workflowRel) {
		t.Fatalf("derived workflow-dir spelling is not relative: %q", workflowRel)
	}

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityPath,
		"workflow_dir":   workflowRel,
		"stage":          "implementation",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      false,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", workflowRel)
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	// Pull the `git -C <path>` token out of the emitted state-commit command and
	// assert it is absolute — a relative `git -C wf/state-checkout` is the defect.
	const marker = "git -C "
	idx := strings.Index(body, marker)
	if idx < 0 {
		t.Fatalf("body missing %q state-commit command\n--- body ---\n%s", marker, body)
	}
	rest := body[idx+len(marker):]
	gitCPath := rest[:strings.IndexByte(rest, ' ')]
	if !filepath.IsAbs(gitCPath) {
		t.Errorf("emitted `git -C` path is not absolute: %q\n--- body ---\n%s", gitCPath, body)
	}
}

// TestEntityPathAbsoluteFromRelativeInput pins AC-1: a dispatch built with a
// RELATIVE entity_path must emit an ABSOLUTE path in BOTH the dispatch-body
// entity-read line ("Read the entity file at <path>") AND the completion-signal
// line ("Report written to <path>."). The dispatched ensign is a separate agent
// whose cwd is not pinned to the workflow root, so a relative entity_path
// resolves to a different / nonexistent file there — the "two entity files"
// divergence. runBuild absolutizes entity_path against the process cwd, mirroring
// the workflow_dir absolutization, so every emitted reference is cwd-independent.
func TestEntityPathAbsoluteFromRelativeInput(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := t.TempDir()

	// Split-root workflow: the entity lives in the state checkout, which is the
	// live-cycle shape this fix targets.
	workflowDir := filepath.Join(root, "wf")
	stateCheckout := filepath.Join(workflowDir, "state-checkout")
	writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(true))

	worktreeRel := ".worktrees/spacedock-ensign-thing"
	if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
		t.Fatal(err)
	}

	entityPath := filepath.Join(stateCheckout, "thing", "index.md")
	writeFile(t, entityPath, entityFM("Thing", "implementation", worktreeRel))

	gitInit(t, root)

	// Derive a RELATIVE entity_path spelling from the process cwd, so runBuild's
	// filepath.Abs resolution is the thing under test. (t.Chdir is go1.24+; this
	// module targets go1.22, so the relative spelling is derived, not realized via
	// a cwd change.)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	entityRel, err := filepath.Rel(cwd, entityPath)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.IsAbs(entityRel) {
		t.Fatalf("derived entity_path spelling is not relative: %q", entityRel)
	}

	stdin := mergeStdin(map[string]any{
		"schema_version": 2,
		"entity_path":    entityRel,
		"workflow_dir":   workflowDir,
		"stage":          "implementation",
		"checklist":      []string{"- a", "- b"},
		"bare_mode":      false,
	}, nil)

	native := runNative(stdin, "build", "--workflow-dir", workflowDir)
	body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

	// Pull the entity-read line's path and assert it is absolute.
	readPath := extractTokenAfter(t, body, "Read the entity file at ", body)
	if !filepath.IsAbs(readPath) {
		t.Errorf("entity-read line path is not absolute: %q\n--- body ---\n%s", readPath, body)
	}

	// Pull the completion-signal line's path and assert it is absolute. The
	// signal renders "Report written to <path>." with a trailing period, so trim
	// it before the IsAbs check.
	signalPath := strings.TrimSuffix(extractTokenAfter(t, body, "Report written to ", body), ".")
	if !filepath.IsAbs(signalPath) {
		t.Errorf("completion-signal line path is not absolute: %q\n--- body ---\n%s", signalPath, body)
	}
}

// The legacy-era sibling of this file, TestDispatchFilePathCollisionFreeAcrossTeams,
// pinned dispatch-file collision-avoidance keyed on a TeamCreate-issued team name.
// That mechanism retired with legacy team mode; the merged floor's equivalent
// disambiguator ($CLAUDE_CODE_SESSION_ID) is covered by
// TestBuildMergedModeDispatchFileDisambiguator in build_merged_mode_test.go.

// extractTokenAfter returns the whitespace-delimited token immediately following
// marker in body, failing the test if the marker is absent.
func extractTokenAfter(t *testing.T, body, marker, full string) string {
	t.Helper()
	idx := strings.Index(body, marker)
	if idx < 0 {
		t.Fatalf("body missing %q\n--- body ---\n%s", marker, full)
	}
	rest := body[idx+len(marker):]
	end := strings.IndexAny(rest, " \n")
	if end < 0 {
		end = len(rest)
	}
	return rest[:end]
}

// TestSingleRootNoStateCommitGuidance (AC-3) pins the single-root negative: a
// workflow with NO state: field must emit NEITHER the "This workflow is
// split-root" sentence NOR a `git -C` state-commit command — for both a worktree
// stage and a non-worktree stage. The cross-product parity test strips the
// guidance line before comparing, so only this dedicated negative guards against
// a single-root guidance LEAK.
func TestSingleRootNoStateCommitGuidance(t *testing.T) {
	cases := []struct {
		name     string
		worktree bool
		stage    string
	}{
		{name: "worktree-stage", worktree: true, stage: "implementation"},
		{name: "nonworktree-stage", worktree: false, stage: "backlog"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			root := t.TempDir()

			// Single root: no state: field, so the entity lives beside the README in
			// the workflow dir and no state-commit guidance applies.
			workflowDir := root
			writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(false))

			worktreeRel := ""
			if tc.worktree {
				worktreeRel = ".worktrees/spacedock-ensign-thing"
				if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
					t.Fatal(err)
				}
			}

			entityPath := filepath.Join(workflowDir, "thing.md")
			writeFile(t, entityPath, entityFM("Thing", tc.stage, worktreeRel))

			gitInit(t, root)

			stdin := mergeStdin(map[string]any{
				"schema_version": 2,
				"entity_path":    entityPath,
				"workflow_dir":   workflowDir,
				"stage":          tc.stage,
				"checklist":      []string{"- a", "- b"},
				"bare_mode":      false,
			}, nil)

			native := runNative(stdin, "build", "--workflow-dir", workflowDir)
			body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

			if strings.Contains(body, "This workflow is split-root") {
				t.Errorf("%s: single-root body leaks the split-root guidance sentence\n--- body ---\n%s",
					tc.name, body)
			}
			if strings.Contains(body, "git -C ") {
				t.Errorf("%s: single-root body leaks a `git -C` state-commit command\n--- body ---\n%s",
					tc.name, body)
			}
			if strings.Contains(body, "pull --rebase") {
				t.Errorf("%s: single-root body leaks a state-branch push/`pull --rebase` reminder\n--- body ---\n%s",
					tc.name, body)
			}
		})
	}
}

// TestStateCommitGuidanceResolvesPaths pins both defects on the split-root
// state-commit-guidance surface: a worktree stage AND a non-worktree stage must
// each emit the path-scoped state-commit instruction with the resolved absolute
// state-checkout and entity paths substituted in — and never a literal
// {state_checkout} or {entity_path} brace token.
func TestStateCommitGuidanceResolvesPaths(t *testing.T) {
	cases := []struct {
		name     string
		worktree bool
		stage    string
	}{
		{name: "worktree-stage", worktree: true, stage: "implementation"},
		{name: "nonworktree-stage", worktree: false, stage: "backlog"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			root := t.TempDir()

			// Split root: the workflow dir is the README/definition root; the
			// README's `state: state-checkout` field diverges the entity/state
			// checkout to workflowDir/state-checkout. CODE (if any) lives in the
			// worktree under root.
			workflowDir := root
			stateCheckout := filepath.Join(workflowDir, "state-checkout")
			writeFile(t, filepath.Join(workflowDir, "README.md"), readmeWorktree(true))

			worktreeRel := ""
			if tc.worktree {
				worktreeRel = ".worktrees/spacedock-ensign-thing"
				if err := os.MkdirAll(filepath.Join(root, worktreeRel), 0o755); err != nil {
					t.Fatal(err)
				}
			}

			entityPath := filepath.Join(stateCheckout, "thing", "index.md")
			writeFile(t, entityPath, entityFM("Thing", tc.stage, worktreeRel))

			gitInit(t, root)
			// This case asserts the push / pull-rebase reminder IS emitted, which is
			// the origin-backed contract — so the state checkout must have an origin.
			// The no-origin degrade is covered separately in
			// build_state_no_origin_test.go.
			gitAddOrigin(t, root)

			stdin := mergeStdin(map[string]any{
				"schema_version": 2,
				"entity_path":    entityPath,
				"workflow_dir":   workflowDir,
				"stage":          tc.stage,
				"checklist":      []string{"- a", "- b"},
				"bare_mode":      false,
			}, nil)

			native := runNative(stdin, "build", "--workflow-dir", workflowDir)
			body := readDispatchBody(t, dispatchFilePathFromStdout(t, native.stdout))

			// POSITIVE: the resolved path-scoped state-commit command targets the
			// resolved state checkout (workflowDir/state-checkout), not the bare
			// workflow/definition dir, both halves.
			wantAdd := "git -C " + stateCheckout + " add " + entityPath
			wantCommit := "git -C " + stateCheckout + " commit -m \"...\" -- " + entityPath
			if !strings.Contains(body, wantAdd) {
				t.Errorf("%s: body missing resolved add command %q\n--- body ---\n%s",
					tc.name, wantAdd, body)
			}
			if !strings.Contains(body, wantCommit) {
				t.Errorf("%s: body missing resolved commit command %q\n--- body ---\n%s",
					tc.name, wantCommit, body)
			}

			// CONCURRENCY-SAFETY PHRASE: the split-root guidance must keep the
			// "never a bare `git add -A`" warning that guards the shared state index.
			const concurrencyPhrase = "never a bare `git add -A`"
			if !strings.Contains(body, concurrencyPhrase) {
				t.Errorf("%s: body missing concurrency-safety phrase %q\n--- body ---\n%s",
					tc.name, concurrencyPhrase, body)
			}

			// PUSH REMINDER: split-root commits land on a shared orphan state branch
			// peers also write, so the guidance must remind the worker to push that
			// branch and `pull --rebase` on a rejection. The resolved branch name
			// (spacedock-state/<workflow-basename>) is named verbatim.
			stateBranch := "spacedock-state/" + filepath.Base(workflowDir)
			wantPush := "git -C " + stateCheckout + " push origin " + stateBranch
			if !strings.Contains(body, wantPush) {
				t.Errorf("%s: body missing resolved push reminder %q\n--- body ---\n%s",
					tc.name, wantPush, body)
			}
			if !strings.Contains(body, "pull --rebase") {
				t.Errorf("%s: body missing `pull --rebase` on-rejection reminder\n--- body ---\n%s",
					tc.name, body)
			}

			// NEGATIVE (the encoded failure): no literal brace tokens survive.
			if strings.Contains(body, "{state_checkout}") {
				t.Errorf("%s: body still carries literal {state_checkout}\n--- body ---\n%s",
					tc.name, body)
			}
			if strings.Contains(body, "{entity_path}") {
				t.Errorf("%s: body still carries literal {entity_path}\n--- body ---\n%s",
					tc.name, body)
			}
		})
	}
}
