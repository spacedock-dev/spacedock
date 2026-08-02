// ABOUTME: `dispatch build --stamp` — mechanism 3 of the gate-approval-ceremony
// ABOUTME: collapse: folds the ordinary post-gate dispatch steps into one build call.
package dispatch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/statesync"
	"github.com/spacedock-dev/spacedock/internal/status"
)

// runStamp implements `dispatch build --stamp`: fo-dispatch-core's dispatch
// steps 5-7 (frontmatter stamps, state commit+sync, worktree creation) folded
// into the build invocation, run before artifact assembly. It never advances
// `status:` — that stays owned by `gate consume` (gated) or `status --set`
// (non-gated) — so it refuses outright when the entity's status does not
// already equal --stage.
//
// Every failure here is prefixed "dispatch build --stamp:" on stderr so
// fo-dispatch-core's block clause can tell a stamp/sync failure (remedy the
// named problem and rerun the same build, never break-glass) apart from an
// assembly failure (break-glass eligible): a stamp-phase failure burned no
// authority and emitted no envelope.
func runStamp(opts buildOptions, fields map[string]json.RawMessage, stderr io.Writer) int {
	entityPath := jsonString(fields["entity_path"])
	if abs, err := filepath.Abs(entityPath); err == nil {
		entityPath = abs
	}
	workflowDir := opts.WorkflowDir
	if workflowDir == "" {
		workflowDir = jsonString(fields["workflow_dir"])
	}
	if abs, err := filepath.Abs(workflowDir); err == nil {
		workflowDir = abs
	}
	stage := jsonString(fields["stage"])

	if !isFile(entityPath) {
		return stampError(stderr, 1, "entity file not readable at '%s'", entityPath)
	}
	entityFields := status.ParseFrontmatter(entityPath)
	if entityFields["status"] != stage {
		return stampError(stderr, 1,
			"entity status '%s' does not match --stage '%s'; --stamp never advances status "+
				"(that stays owned by gate consume or status --set)", entityFields["status"], stage)
	}

	readmePath := filepath.Join(workflowDir, "README.md")
	readmeData, err := os.ReadFile(readmePath)
	if err != nil {
		return stampError(stderr, 1, "workflow README not found at '%s'", readmePath)
	}
	readmeFields := status.ParseFrontmatterData(readmeData)
	stages, _ := status.ParseStagesWithDefaultsData(readmeData)
	stageIdx := -1
	for i, s := range stages {
		if s.Name == stage {
			stageIdx = i
			break
		}
	}
	if stageIdx < 0 {
		return stampError(stderr, 1, "stage '%s' not found in %s", stage, readmePath)
	}
	stageMeta := stages[stageIdx]

	subagentType := "spacedock:ensign"
	if agent, ok := stageMeta.Agent(); ok {
		subagentType = agent
	}
	workerKey := strings.ReplaceAll(subagentType, ":", "-")
	slug := status.EntitySlug(entityPath)
	worktreeRel := filepath.Join(".worktrees", workerKey+"-"+slug)

	needStarted := entityFields["started"] == ""
	needWorktreeStamp := stageMeta.Worktree && entityFields["worktree"] == ""

	if needStarted || needWorktreeStamp {
		if code := stampFrontmatter(workflowDir, slug, needStarted, needWorktreeStamp, worktreeRel, stderr); code != 0 {
			return code
		}
		if mode, relPath, err := status.ClassifyState(readmeFields["state"]); err == nil && mode == status.StateSplitRoot {
			checkout := filepath.Join(workflowDir, relPath)
			branch, err := status.StateBranch(workflowDir)
			if err != nil {
				return stampError(stderr, 1, "%v", err)
			}
			if code := stampSync(stderr, checkout, branch, entityPath, slug, stage); code != 0 {
				return code
			}
		}
	}

	if stageMeta.Worktree {
		gitRoot := status.FindGitRoot(workflowDir)
		worktreePath := filepath.Join(gitRoot, worktreeRel)
		if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
			branch := workerKey + "/" + slug
			cmd := exec.Command("git", "-C", gitRoot, "worktree", "add", "-b", branch, worktreePath)
			if out, err := cmd.CombinedOutput(); err != nil {
				return stampError(stderr, 1, "git worktree add %s failed: %s", worktreePath, strings.TrimSpace(string(out)))
			}
		}
	}
	return 0
}

// stampFrontmatter stamps `started` (bare, auto-fill) and, for a worktree-
// declaring stage, `worktree=` through the native `status --set` machinery —
// inheriting its mutation guards (mod-block, stage membership, away-status
// refusals) rather than adding a parallel frontmatter writer.
func stampFrontmatter(workflowDir, slug string, needStarted, needWorktree bool, worktreeRel string, stderr io.Writer) int {
	args := []string{"--workflow-dir", workflowDir, "--set", slug}
	if needStarted {
		args = append(args, "started")
	}
	if needWorktree {
		args = append(args, "worktree="+worktreeRel)
	}
	runner := &status.NativeRunner{}
	var setOut, setErr strings.Builder
	code, _ := runner.Run(context.Background(), status.Request{
		Args: args, Dir: workflowDir, Env: os.Environ(), Stdout: &setOut, Stderr: &setErr,
	})
	if code != 0 {
		return stampError(stderr, 1, "stamping %s failed: %s", slug, strings.TrimSpace(setErr.String()))
	}
	return 0
}

// stampSync is mechanism 3's use of the mechanism-1 seam: path-scoped commit +
// publish for the just-stamped entity, message `dispatch: <slug> entering
// <stage>`. A rebase-conflict HALT (exit 3) or a sync failure (exit 1) is
// reported with the same "dispatch build --stamp:" prefix as every other stamp
// failure — the stamp itself already landed durably (locally); this is a
// publication problem, not a lost mutation.
func stampSync(stderr io.Writer, checkout, branch, entityPath, slug, stage string) int {
	msg := fmt.Sprintf("dispatch: %s entering %s", slug, stage)
	_, outcome, err := commitAndPublishEntity(checkout, branch, entityPath, msg)
	if err != nil {
		return stampError(stderr, 1, "state sync failed: %v", err)
	}
	switch outcome.Result {
	case statesync.ResultHalted:
		writeStampHaltStderr(stderr, branch, outcome)
		return 3
	case statesync.ResultFailed:
		return stampError(stderr, 1, "state publication failed; the stamp is already durable locally:\n%s", outcome.Detail)
	default:
		return 0
	}
}

// commitAndPublishEntity stages+commits entityPath's path-scoped commit unit
// (itself, plus a flat entity's companion directory when present) with msg —
// skipping the commit when there is nothing staged — then publishes via the
// shared push/rebase/HALT sequence regardless, so a peer's pending state still
// gets integrated on a clean stamp. This mirrors internal/cli's mechanism-1
// seam (syncActiveEntity) for the one package boundary --stamp cannot cross:
// internal/cli already imports internal/dispatch, so the reverse import would
// cycle.
func commitAndPublishEntity(checkout, branch, entityPath, msg string) (committed bool, outcome statesync.Outcome, err error) {
	pathspecs := literalPathspecs(checkout, entityCommitPaths(checkout, entityPath))
	addArgs := append([]string{"add", "-A", "--"}, pathspecs...)
	if ok, out := runGitRetry(checkout, addArgs...); !ok {
		return false, statesync.Outcome{}, fmt.Errorf("git add failed: %s", strings.TrimSpace(out))
	}
	ok, staged := runGitRetry(checkout, "diff", "--cached", "--name-only", "--no-renames", "-z")
	if !ok {
		return false, statesync.Outcome{}, fmt.Errorf("git diff --cached failed: %s", strings.TrimSpace(staged))
	}
	if staged == "" {
		return false, statesync.Publish(checkout, branch), nil
	}
	commitArgs := append([]string{"commit", "-m", msg, "--"}, pathspecs...)
	if ok, out := runGitRetry(checkout, commitArgs...); !ok {
		return false, statesync.Outcome{}, fmt.Errorf("git commit failed: %s", strings.TrimSpace(out))
	}
	return true, statesync.Publish(checkout, branch), nil
}

// entityCommitPaths returns the path-scoped git commit unit for the entity at
// entityPath (already known to exist): a folder-form entity's whole directory,
// or a flat entity's file plus its companion directory (same slug) when present.
func entityCommitPaths(checkout, entityPath string) []string {
	if filepath.Base(entityPath) == "index.md" {
		return []string{filepath.Dir(entityPath)}
	}
	paths := []string{entityPath}
	slug := strings.TrimSuffix(filepath.Base(entityPath), ".md")
	companion := filepath.Join(checkout, slug)
	if info, err := os.Stat(companion); err == nil && info.IsDir() {
		paths = append(paths, companion)
	}
	return paths
}

// literalPathspecs renders each absolute path relative to checkout as a
// literal (metacharacter-safe) git pathspec.
func literalPathspecs(checkout string, paths []string) []string {
	pathspecs := make([]string, 0, len(paths))
	for _, p := range paths {
		rel, err := filepath.Rel(checkout, p)
		if err != nil {
			rel = p
		}
		pathspecs = append(pathspecs, ":(literal)"+rel)
	}
	return pathspecs
}

// runGitRetry runs a git command in dir, retrying on index.lock contention up
// to a bounded number of times (~2s total) before failing — the shared state
// index is a single non-branched git index a sibling writer's add/commit can
// briefly hold the lock on.
func runGitRetry(dir string, args ...string) (ok bool, out string) {
	const maxRetries = 4
	const wait = 500 * time.Millisecond
	for attempt := 0; ; attempt++ {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		outBytes, err := cmd.CombinedOutput()
		out = string(outBytes)
		ok = err == nil
		if ok || attempt >= maxRetries || !strings.Contains(out, "index.lock") {
			return ok, out
		}
		time.Sleep(wait)
	}
}

// writeStampHaltStderr renders the shared HALT diagnostic (same prose as the
// gate verbs' and `state commit`'s HALT rendering) under the "dispatch build
// --stamp:" prefix.
func writeStampHaltStderr(stderr io.Writer, branch string, outcome statesync.Outcome) {
	reportedPaths := strings.Join(outcome.ConflictingPaths, ", ")
	if reportedPaths == "" {
		reportedPaths = "none reported by Git"
	}
	fmt.Fprintf(stderr, "dispatch build --stamp: HALT — same-entity rebase conflict on %s.\n", branch)
	fmt.Fprintf(stderr, "Conflicting path(s): %s\n", reportedPaths)
	if outcome.PeerCommit != "" {
		fmt.Fprintf(stderr, "Peer commit: %s (origin/%s)\n", outcome.PeerCommit, branch)
	}
	fmt.Fprintln(stderr, "The rebase was aborted (checkout left clean) and nothing was force-pushed; a peer's edit is preserved on origin.")
	fmt.Fprintln(stderr, "Next: HALT dispatch — do not dispatch against this state tree. Surface the conflicting path(s) and peer commit to the operator and stop.")
	fmt.Fprintln(stderr, "Never `git push --force`/`--force-with-lease`; never re-run with `-X ours`/`-X theirs`; never discard either side.")
}

// stampError writes a "dispatch build --stamp:"-prefixed stderr diagnostic and
// emits NO envelope — the discriminator fo-dispatch-core's block clause reads
// to distinguish a stamp/sync failure (remedy + rerun) from an assembly
// failure (break-glass eligible).
func stampError(stderr io.Writer, code int, format string, a ...any) int {
	fmt.Fprintf(stderr, "dispatch build --stamp: "+format+"\n", a...)
	return code
}
