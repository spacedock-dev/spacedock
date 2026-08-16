// ABOUTME: `spacedock state commit <slug>` / `state ready` — path-scoped state
// ABOUTME: commit+push and boot-pull, each HALTing (exit 3) on a same-entity rebase conflict.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/dispatch"
	"github.com/spacedock-dev/spacedock/internal/statesync"
	"github.com/spacedock-dev/spacedock/internal/status"
)

// syncResult is the small JSON object the state-sync verbs emit under --json. The
// result field is the verb's outcome verb (committed/pushed/local-only/no-op/
// halted/ready); conflictingPaths is populated only on a halt. reason carries a
// human one-liner for the prose and the JSON alike.
type syncResult struct {
	Command          string   `json:"command"`
	Slug             string   `json:"slug,omitempty"`
	Result           string   `json:"result"`
	StateBranch      string   `json:"state_branch,omitempty"`
	ConflictingPaths []string `json:"conflicting_paths,omitempty"`
	PeerCommit       string   `json:"peer_commit,omitempty"`
	Reason           string   `json:"reason"`
}

// emitSync writes the result as JSON (jsonOut) or as a one-line prose summary, then
// returns code. Centralizing the dual rendering keeps every verb's exit path
// identical: the JSON envelope and the prose say the same thing.
func emitSync(stdout io.Writer, jsonOut bool, res syncResult, code int) int {
	if jsonOut {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		enc.Encode(res)
		return code
	}
	fmt.Fprintln(stdout, res.Reason)
	return code
}

// runStateCommit implements `spacedock state commit <slug>`. An inline workflow
// commits in the workflow's own repository and pushes nothing (commitInlineEntity);
// the rest of this function is the split-root path. Active scope commits
// path-scoped; clean archived scope publishes existing history. Both use the
// shared push → on-reject pull --rebase → re-push sequence. A same-entity rebase
// conflict HALTS: the rebase is aborted (clean tree restored), no force-push is
// ever issued, and the verb exits 3 with the conflicting path named — so a caller
// cannot proceed on an unmerged tree. The halt is the exit code, not the model's
// discipline.
func runStateCommit(ctx context.Context, args []string, env []string, dir string, stdout, stderr io.Writer) int {
	slug, workflowDir, msg, jsonOut, code := parseStateCommitArgs(args, dir, stderr)
	if code != 0 {
		return code
	}

	checkout, branch, mode, code := resolveStateCheckout("state commit", workflowDir, stderr)
	if code != 0 {
		return code
	}
	if mode == status.StateInline {
		return commitInlineEntity(stdout, stderr, jsonOut, workflowDir, slug, msg)
	}

	if outcome := statesync.Preflight(checkout, branch); outcome.Result == statesync.ResultHalted {
		return emitSyncHalt(stdout, stderr, jsonOut, "state commit", slug, branch, outcome)
	} else if outcome.Result == statesync.ResultFailed {
		fmt.Fprintf(stderr, "spacedock state commit: state-sync preflight failed: %s\n", outcome.Detail)
		return 1
	}

	target, found, resolveErr := resolveEntityCommitTarget(checkout, slug)
	if resolveErr != nil {
		fmt.Fprintf(stderr, "spacedock state commit: %v\n", resolveErr)
		return 1
	}
	if !found {
		fmt.Fprintf(stderr, "spacedock state commit: no entity %q under %s (looked in active and _archive scope)\n", slug, checkout)
		return 1
	}
	if msg == "" {
		msg = fmt.Sprintf("state: update %s", slug)
	}

	committed := false
	var outcome statesync.Outcome
	if target.scope == entityScopeArchived {
		clean, detail, checkErr := archivedTargetClean(checkout, target.pathspecs)
		if checkErr != nil {
			fmt.Fprintf(stderr, "spacedock state commit: could not inspect archived entity %q before publication: %v\n", slug, checkErr)
			return 1
		}
		if !clean {
			fmt.Fprintf(stderr, "spacedock state commit: archived entity %q is dirty; archived scope is publish-only and will not stage or commit it:\n%s\n", slug, detail)
			return 1
		}
		outcome = statesync.Publish(checkout, branch)
	} else {
		// Active scope retains the path-scoped stage+commit behavior through the
		// mechanism-1 seam shared with the gate record/consume verbs' implicit sync.
		var syncErr error
		committed, outcome, syncErr = syncActiveEntity(checkout, branch, slug, msg)
		if syncErr != nil {
			fmt.Fprintf(stderr, "spacedock state commit: git commit failed:\n%v\n", syncErr)
			return 1
		}
	}

	switch outcome.Result {
	case statesync.ResultHalted:
		return emitSyncHalt(stdout, stderr, jsonOut, "state commit", slug, branch, outcome)
	case statesync.ResultFailed:
		switch {
		case target.scope == entityScopeArchived:
			fmt.Fprintf(stderr, "spacedock state commit: state publication failed; the existing archive commit remains recoverable:\n%s\n", outcome.Detail)
		case committed:
			fmt.Fprintf(stderr, "spacedock state commit: state publication failed; the new local commit remains recoverable:\n%s\n", outcome.Detail)
		default:
			fmt.Fprintf(stderr, "spacedock state commit: state synchronization failed; no new commit was created in this invocation, and existing local history remains recoverable:\n%s\n", outcome.Detail)
		}
		return 1
	case statesync.ResultLocalOnly:
		if target.scope == entityScopeActive && !committed {
			return emitSync(stdout, jsonOut, syncResult{
				Command: "state commit", Slug: slug, Result: "no-op", StateBranch: branch,
				Reason: fmt.Sprintf("Nothing to commit for %s — state checkout already up to date.", slug),
			}, 0)
		}
		reason := fmt.Sprintf("Committed %s locally; no origin remote — state is local-only until an origin is configured.", slug)
		if target.scope == entityScopeArchived {
			reason = fmt.Sprintf("Archived state for %s is committed locally; no origin remote — state is local-only.", slug)
		}
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "local-only", StateBranch: branch,
			Reason: reason,
		}, 0)
	case statesync.ResultNoOp:
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "no-op", StateBranch: branch,
			Reason: fmt.Sprintf("Nothing to commit for %s — state checkout already up to date.", slug),
		}, 0)
	case statesync.ResultPushed:
		if target.scope == entityScopeActive && !committed && !outcome.PublishedLocal {
			return emitSync(stdout, jsonOut, syncResult{
				Command: "state commit", Slug: slug, Result: "no-op", StateBranch: branch,
				Reason: fmt.Sprintf("Nothing to commit for %s — integrated peers' state; checkout is up to date.", slug),
			}, 0)
		}
		reason := fmt.Sprintf("Committed and pushed %s to %s.", slug, branch)
		if target.scope == entityScopeArchived {
			reason = fmt.Sprintf("Published existing archived state for %s to %s.", slug, branch)
		} else if !committed {
			reason = fmt.Sprintf("Published previously committed state for %s to %s.", slug, branch)
		} else if outcome.IntegratedPeers {
			reason = fmt.Sprintf("Committed %s, integrated peers' state, and pushed to %s.", slug, branch)
		}
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "pushed", StateBranch: branch,
			Reason: reason,
		}, 0)
	default:
		fmt.Fprintln(stderr, "spacedock state commit: unexpected state publication result")
		return 1
	}
}

// commitInlineEntity is `state commit`'s inline half. An inline workflow keeps
// its entities beside the README, so the workflow dir IS the entity root and the
// same path-scoped commit-unit seam split-root uses applies to the workflow's own
// repository. It never publishes: an inline workflow repo is the caller's code
// repo, and this verb has no authority to push it. That durability matters
// because `gate prepare` and `status --next`'s `needs-preparation` row both read
// committed bytes only — an uncommitted entity is invisible to the first and
// refused by the second, so a no-op here strands the caller with no verb that can
// clear the condition. A workflow dir outside a Git work tree stays an exit-0
// no-op, but names that as the reason instead of claiming a commit happened.
func commitInlineEntity(stdout, stderr io.Writer, jsonOut bool, workflowDir, slug, msg string) int {
	if !insideGitWorkTree(workflowDir) {
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "no-op",
			Reason: fmt.Sprintf("Inline workflow at %s is not inside a Git work tree — %s has nothing to commit to.", workflowDir, slug),
		}, 0)
	}
	target, found, resolveErr := resolveEntityCommitTarget(workflowDir, slug)
	if resolveErr != nil {
		fmt.Fprintf(stderr, "spacedock state commit: %v\n", resolveErr)
		return 1
	}
	if !found {
		fmt.Fprintf(stderr, "spacedock state commit: no entity %q under %s (looked in active and _archive scope)\n", slug, workflowDir)
		return 1
	}
	// Archived scope is publish-only in split-root, and inline has nothing to
	// publish to — so the dirt refusal still applies (it protects the archive
	// move from a half-staged commit), and clean archived state is a no-op.
	if target.scope == entityScopeArchived {
		clean, detail, checkErr := archivedTargetClean(workflowDir, target.pathspecs)
		if checkErr != nil {
			fmt.Fprintf(stderr, "spacedock state commit: could not inspect archived entity %q: %v\n", slug, checkErr)
			return 1
		}
		if !clean {
			fmt.Fprintf(stderr, "spacedock state commit: archived entity %q is dirty; archived scope is publish-only and will not stage or commit it:\n%s\n", slug, detail)
			return 1
		}
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "no-op",
			Reason: fmt.Sprintf("Archived state for %s is already committed; an inline workflow publishes nothing.", slug),
		}, 0)
	}
	if msg == "" {
		msg = fmt.Sprintf("state: update %s", slug)
	}
	ok, out := commitEntityPathsScoped(workflowDir, target.entityPaths, msg)
	if !ok {
		fmt.Fprintf(stderr, "spacedock state commit: git commit failed:\n%s\n", out)
		return 1
	}
	if out == "" {
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "no-op",
			Reason: fmt.Sprintf("Nothing to commit for %s — inline workflow already up to date.", slug),
		}, 0)
	}
	return emitSync(stdout, jsonOut, syncResult{
		Command: "state commit", Slug: slug, Result: "committed",
		Reason: fmt.Sprintf("Committed %s in the inline workflow repository; nothing pushed.", slug),
	}, 0)
}

// insideGitWorkTree reports whether dir sits inside a Git work tree. A bare repo
// or a plain directory answers no, which keeps the inline commit an honest no-op
// rather than a git failure the caller cannot act on.
func insideGitWorkTree(dir string) bool {
	ok, out := runGitOutput(dir, "rev-parse", "--is-inside-work-tree")
	return ok && strings.TrimSpace(out) == "true"
}

// runStateReady implements `spacedock state ready`. It is the boot gate: an inline
// workflow is a clean no-op; a split-root checkout that is absent is resumed (the
// same fetch+worktree-add resume `state init` does); a present origin-backed
// checkout integrates peers' state with one pull --rebase. A same-entity boot
// conflict HALTS identically to `state commit` (abort, exit 3, name the paths).
func runStateReady(ctx context.Context, args []string, env []string, dir string, stdout, stderr io.Writer) int {
	workflowDir, jsonOut, code := parseWorkflowJSONArgs("state ready", args, dir, stderr)
	if code != 0 {
		return code
	}

	checkout, branch, mode, code := resolveStateCheckout("state ready", workflowDir, stderr)
	if code != 0 {
		return code
	}
	if mode == status.StateInline {
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state ready", Result: "ready",
			Reason: "Inline workflow — entities live beside the README; nothing to sync.",
		}, 0)
	}

	// Absent checkout → resume it (reuse `state init`'s fetch + worktree-add path).
	if !dirExists(checkout) {
		if code := runStateInit(ctx, []string{"--workflow-dir", workflowDir}, env, dir, stdout, stderr); code != 0 {
			return code
		}
		// The re-boot-after-resume sequencing the «state.ensure-ready» prose used to
		// own: a just-linked checkout means the boot read the FO already did (if any)
		// predates the entity dir existing, so it must re-read before greeting.
		// Prose-only (not --json) — it is FO guidance, not part of the result envelope.
		if !jsonOut {
			fmt.Fprintln(stdout, "checkout resumed — re-run `spacedock status --boot` before the greet.")
		}
	}

	outcome := statesync.Pull(checkout, branch)
	switch outcome.Result {
	case statesync.ResultLocalOnly:
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state ready", Result: "ready", StateBranch: branch,
			Reason: "State checkout ready (no origin remote — state is local-only).",
		}, 0)
	case statesync.ResultHalted:
		return emitSyncHalt(stdout, stderr, jsonOut, "state ready", "", branch, outcome)
	case statesync.ResultFailed:
		fmt.Fprintf(stderr, "spacedock state ready: state synchronization failed:\n%s\n", outcome.Detail)
		return 1
	case statesync.ResultReady:
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state ready", Result: "ready", StateBranch: branch,
			Reason: fmt.Sprintf("State checkout ready — integrated peers' state from %s.", branch),
		}, 0)
	default:
		fmt.Fprintln(stderr, "spacedock state ready: unexpected state synchronization result")
		return 1
	}
}

// runStateSweep implements `spacedock state sweep`. It is the state-repo's
// read-only view of entities whose code PR merged but whose state is not yet
// terminalized — the merged-PR iteration the FO sweeps each loop. It makes no
// commit, push, or mutation; it delegates to dispatch.Sweep, which reuses
// reconcile's un-advanced-PR detection rather than re-implementing merged-state.
func runStateSweep(ctx context.Context, args []string, env []string, dir string, stdout, stderr io.Writer) int {
	workflowDir, jsonOut, code := parseWorkflowJSONArgs("state sweep", args, dir, stderr)
	if code != 0 {
		return code
	}
	return dispatch.Sweep(workflowDir, dispatch.GhRunnerExec, jsonOut, stdout, stderr)
}

// emitSyncHalt renders the shared HALT after statesync captured the conflict
// evidence and aborted its same-branch rebase. It never discovers paths or
// mutates Git itself; the exit code enforces that a caller cannot proceed.
func emitSyncHalt(stdout, stderr io.Writer, jsonOut bool, command, slug, branch string, outcome statesync.Outcome) int {
	reportedPaths := writeSyncHaltStderr(stderr, command, branch, outcome)
	return emitSync(stdout, jsonOut, syncResult{
		Command: command, Slug: slug, Result: "halted", StateBranch: branch,
		ConflictingPaths: outcome.ConflictingPaths,
		PeerCommit:       outcome.PeerCommit,
		Reason:           fmt.Sprintf("HALT: same-entity rebase conflict on %s — rebase aborted, nothing force-pushed, manual intervention required.", reportedPaths),
	}, 3)
}

// writeSyncHaltStderr renders the shared HALT stderr diagnostic — command names
// the calling verb in the first line. Shared by `state commit`/`state ready`'s
// JSON/prose envelope (emitSyncHalt) and the gate verbs' single-line
// `sync=halted` rendering (mechanism 1), so the HALT prose stays byte-identical
// across every caller. Returns the rendered conflicting-paths string for reuse
// in a caller's own summary.
func writeSyncHaltStderr(stderr io.Writer, command, branch string, outcome statesync.Outcome) string {
	reportedPaths := strings.Join(outcome.ConflictingPaths, ", ")
	if reportedPaths == "" {
		reportedPaths = "none reported by Git"
	}
	fmt.Fprintf(stderr, "spacedock %s: HALT — same-entity rebase conflict on %s.\n", command, branch)
	fmt.Fprintf(stderr, "Conflicting path(s): %s\n", reportedPaths)
	if outcome.PeerCommit != "" {
		fmt.Fprintf(stderr, "Peer commit: %s (origin/%s)\n", outcome.PeerCommit, branch)
	}
	fmt.Fprintf(stderr, "The rebase was aborted (checkout left clean) and nothing was force-pushed; a peer's edit is preserved on origin.\n")
	fmt.Fprintf(stderr, "Next: HALT dispatch — do not dispatch against this state tree. Surface the conflicting path(s) and peer commit to the operator and stop.\n")
	fmt.Fprintf(stderr, "Never `git push --force`/`--force-with-lease`; never re-run with `-X ours`/`-X theirs`; never discard either side.\n")
	return reportedPaths
}

// syncActiveEntity is the mechanism-1 implicit-sync seam: resolve slug's
// active-scope commit-unit paths, stage+commit them (skip if clean), then
// publish via the shared push/rebase/HALT sequence. Shared by `state commit`'s
// active-scope path (below) and the gate record/consume verbs' implicit sync —
// every caller already knows checkout/branch name a ready StateSplitRoot
// checkout. Archived-scope entities are out of scope for every caller (a gate
// only ever closes/consumes a live, active entity), so this refuses rather than
// silently falling back to an unscoped `git add -A`.
func syncActiveEntity(checkout, branch, slug, msg string) (committed bool, outcome statesync.Outcome, err error) {
	target, found, err := resolveEntityCommitTarget(checkout, slug)
	if err != nil {
		return false, statesync.Outcome{}, err
	}
	if !found {
		return false, statesync.Outcome{}, fmt.Errorf("no entity %q under %s", slug, checkout)
	}
	if target.scope != entityScopeActive {
		return false, statesync.Outcome{}, fmt.Errorf("entity %q resolved to archived scope; this seam only commits active entities", slug)
	}
	ok, out := commitEntityPathsScoped(checkout, target.entityPaths, msg)
	if !ok {
		return false, statesync.Outcome{}, fmt.Errorf("git commit failed:\n%s", out)
	}
	committed = out != ""
	outcome = statesync.Publish(checkout, branch)
	return committed, outcome, nil
}

// commitEntityPathsScoped stages and commits exactly the entity unit with `add
// -A` restricted to literal pathspecs, retrying staging on index.lock contention.
// It returns (true, "") for a
// clean no-op (nothing staged for the entity → success), (true, output) for a
// landed commit, and (false, output) on a real git failure. The `git -C` argv
// form (via runGit) has no rendered command string for a weak model to paraphrase.
func commitEntityPathsScoped(checkout string, entityPaths []string, msg string) (ok bool, output string) {
	pathspecs := make([]string, 0, len(entityPaths))
	for _, entityPath := range entityPaths {
		pathspecs = append(pathspecs, literalGitPathspec(relToCheckout(checkout, entityPath)))
	}
	addArgs := append([]string{"add", "-A", "--"}, pathspecs...)
	if ok, out := runGitRetryLock(checkout, addArgs...); !ok {
		return false, out
	}
	// Commit only paths with staged changes. A flat entity's companion path is
	// part of its commit unit, but an absent, never-tracked companion is not a
	// valid `git commit -- <path>` operand.
	//
	// `--relative` is load-bearing, not tidiness: without it `diff --cached`
	// reports names relative to the REPO ROOT while entityPaths are relative to
	// checkout. Those agree only when checkout IS the repo root — true for a
	// split-root state worktree, false for an inline workflow nested under the
	// repo (docs/dev/...). There the names never match, nothing is selected to
	// commit, and the verb reports "already up to date" having just staged the
	// change — the lying no-op this seam exists to prevent. The commit pathspecs
	// below are resolved against `git -C checkout`, so they must be
	// checkout-relative too.
	ok, stagedNames := runGitOutput(checkout, "diff", "--cached", "--name-only", "--no-renames", "--relative", "-z")
	if !ok {
		return false, stagedNames
	}
	var changedPathspecs []string
	for _, staged := range strings.Split(stagedNames, "\x00") {
		if staged == "" {
			continue
		}
		for _, entityPath := range entityPaths {
			rel := filepath.ToSlash(relToCheckout(checkout, entityPath))
			if staged == rel || strings.HasPrefix(staged, rel+"/") {
				changedPathspecs = append(changedPathspecs, literalGitPathspec(staged))
				break
			}
		}
	}
	if len(changedPathspecs) == 0 {
		return true, ""
	}
	commitArgs := append([]string{"commit", "-m", msg, "--"}, changedPathspecs...)
	ok, out := runGitRetryLock(checkout, commitArgs...)
	if !ok {
		return false, out
	}
	return true, out
}

// runGitRetryLock runs a git command, retrying on index.lock contention up to a
// bounded number of times (~2s total) before failing. The shared state index is a
// single non-branched git index; a sibling writer's add/commit can hold the lock.
func runGitRetryLock(checkout string, gitArgs ...string) (bool, string) {
	const maxRetries = 4
	const wait = 500 * time.Millisecond
	for attempt := 0; ; attempt++ {
		ok, out := runGit(checkout, gitArgs...)
		if ok || attempt >= maxRetries || !strings.Contains(out, "index.lock") {
			return ok, out
		}
		time.Sleep(wait)
	}
}

type entityCommitScope int

const (
	entityScopeActive entityCommitScope = iota
	entityScopeArchived
)

type entityCommitTarget struct {
	scope       entityCommitScope
	entityPaths []string
	pathspecs   []string
}

// resolveEntityCommitTarget validates identity shape before staging or
// publication. Active scope is a commit unit; archive scope is publish-only and
// carries both sides of the canonical archive move for dirt checks.
func resolveEntityCommitTarget(checkout, slug string) (entityCommitTarget, bool, error) {
	activeFlat := filepath.Join(checkout, slug+".md")
	activeIndex := filepath.Join(checkout, slug, "index.md")
	archiveFlat := filepath.Join(checkout, "_archive", slug+".md")
	archiveIndex := filepath.Join(checkout, "_archive", slug, "index.md")

	activeFlatExists := fileExists(activeFlat) && status.EntitySlug(activeFlat) == slug
	activeFolderExists := fileExists(activeIndex) && status.EntitySlug(activeIndex) == slug
	archiveFlatExists := fileExists(archiveFlat) && status.EntitySlug(archiveFlat) == slug
	archiveFolderExists := fileExists(archiveIndex) && status.EntitySlug(archiveIndex) == slug

	if archiveFlatExists && archiveFolderExists {
		return entityCommitTarget{}, false, fmt.Errorf("archive-shape collision for %q: both _archive/%s.md and _archive/%s/index.md exist; remove the invalid duplicate, then rerun state commit", slug, slug, slug)
	}
	if (activeFlatExists || activeFolderExists) && (archiveFlatExists || archiveFolderExists) {
		return entityCommitTarget{}, false, fmt.Errorf("active/archive identity collision for %q: the slug exists in both scopes; resolve the duplicate identity, then rerun state commit", slug)
	}
	if archiveFolderExists {
		return entityCommitTarget{
			scope: entityScopeArchived,
			pathspecs: []string{
				slug,
				filepath.Join("_archive", slug),
			},
		}, true, nil
	}
	if archiveFlatExists {
		return entityCommitTarget{
			scope: entityScopeArchived,
			pathspecs: []string{
				slug + ".md",
				filepath.Join("_archive", slug+".md"),
				slug,
				filepath.Join("_archive", slug),
			},
		}, true, nil
	}
	if activeFolderExists {
		return entityCommitTarget{scope: entityScopeActive, entityPaths: []string{filepath.Dir(activeIndex)}}, true, nil
	}
	if activeFlatExists {
		return entityCommitTarget{scope: entityScopeActive, entityPaths: flatEntityCommitPaths(checkout, activeFlat, slug)}, true, nil
	}
	// Preserve deletion commits for active entities that have not moved into
	// archive scope: tracked-but-missing canonical paths remain commit units.
	if gitTracksPath(checkout, activeIndex) && status.EntitySlug(activeIndex) == slug {
		return entityCommitTarget{scope: entityScopeActive, entityPaths: []string{filepath.Dir(activeIndex)}}, true, nil
	}
	if gitTracksPath(checkout, activeFlat) && status.EntitySlug(activeFlat) == slug {
		return entityCommitTarget{scope: entityScopeActive, entityPaths: flatEntityCommitPaths(checkout, activeFlat, slug)}, true, nil
	}
	return entityCommitTarget{}, false, nil
}

func flatEntityCommitPaths(checkout, flatPath, slug string) []string {
	paths := []string{flatPath}
	companion := filepath.Join(checkout, slug)
	if dirExists(companion) || gitTracksPathTree(checkout, companion) {
		paths = append(paths, companion)
	}
	return paths
}

func archivedTargetClean(checkout string, paths []string) (bool, string, error) {
	args := []string{"status", "--porcelain=v1", "--untracked-files=all", "--"}
	for _, path := range paths {
		args = append(args, literalGitPathspec(path))
	}
	ok, out := runGit(checkout, args...)
	if !ok {
		return false, "", fmt.Errorf("%s", strings.TrimSpace(out))
	}
	detail := strings.TrimSpace(out)
	return detail == "", detail, nil
}

func gitTracksPath(checkout, path string) bool {
	ok, _ := runGit(checkout, "ls-files", "--error-unmatch", "--", literalGitPathspec(relToCheckout(checkout, path)))
	return ok
}

func gitTracksPathTree(checkout, path string) bool {
	ok, out := runGitOutput(checkout, "ls-files", "-z", "--", literalGitPathspec(relToCheckout(checkout, path)))
	return ok && out != ""
}

func literalGitPathspec(path string) string {
	return ":(literal)" + path
}

// validStateEntitySlug accepts the filesystem-level slug shape canonical entity
// discovery can expose: one visible top-level name, never a path alias. It does
// not impose a separate character grammar on otherwise valid entity filenames.
func validStateEntitySlug(slug string) bool {
	return slug != "" && slug != "." && slug != ".." &&
		!filepath.IsAbs(slug) && !strings.HasPrefix(slug, ".") &&
		!strings.ContainsAny(slug, `/\\`) && filepath.Clean(slug) == slug
}

// relToCheckout returns entityPath relative to checkout for the path-scoped git
// pathspec. An already-relative path is returned as-is.
func relToCheckout(checkout, entityPath string) string {
	if !filepath.IsAbs(entityPath) {
		return entityPath
	}
	rel, err := filepath.Rel(checkout, entityPath)
	if err != nil {
		return filepath.Base(entityPath)
	}
	return rel
}

// resolveStateCheckout reads the workflow README to resolve the split-root state
// checkout path and its state branch, classifying the mode. An inline workflow
// returns mode=StateInline with empty checkout/branch (the caller no-ops). A
// missing README or a malformed state: field is a setup failure (exit 1).
func resolveStateCheckout(command, workflowDir string, stderr io.Writer) (checkout, branch string, mode status.StateMode, code int) {
	readme := filepath.Join(workflowDir, "README.md")
	if !fileExists(readme) {
		fmt.Fprintf(stderr, "spacedock %s: no README.md at %s\n", command, workflowDir)
		return "", "", status.StateInline, 1
	}
	fm := status.ParseFrontmatter(readme)
	m, relPath, err := status.ClassifyState(fm["state"])
	if err != nil {
		fmt.Fprintf(stderr, "spacedock %s: %v\n", command, err)
		return "", "", status.StateInline, 1
	}
	if m == status.StateInline {
		return "", "", status.StateInline, 0
	}
	b, err := status.StateBranch(workflowDir)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock %s: %v\n", command, err)
		return "", "", status.StateInline, 1
	}
	return filepath.Join(workflowDir, relPath), b, status.StateSplitRoot, 0
}

// parseStateCommitArgs reads the positional <slug> plus `--workflow-dir DIR`,
// `-m MSG`, and `--json`. A missing slug is a usage error (exit 2).
func parseStateCommitArgs(args []string, dir string, stderr io.Writer) (slug, workflowDir, msg string, jsonOut bool, code int) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--workflow-dir":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "spacedock state commit: --workflow-dir requires a path")
				return "", "", "", false, 2
			}
			workflowDir = args[i+1]
			i++
		case "-m", "--message":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "spacedock state commit: -m requires a message")
				return "", "", "", false, 2
			}
			msg = args[i+1]
			i++
		case "--json":
			jsonOut = true
		default:
			if strings.HasPrefix(args[i], "-") {
				fmt.Fprintf(stderr, "spacedock state commit: unknown flag %q\n", args[i])
				return "", "", "", false, 2
			}
			if slug != "" {
				fmt.Fprintf(stderr, "spacedock state commit: unexpected extra argument %q (one slug only)\n", args[i])
				return "", "", "", false, 2
			}
			slug = args[i]
		}
	}
	if slug == "" {
		fmt.Fprintln(stderr, "spacedock state commit: missing required <slug> argument")
		return "", "", "", false, 2
	}
	if !validStateEntitySlug(slug) {
		fmt.Fprintf(stderr, "spacedock state commit: invalid entity slug %q (expected one canonical top-level slug, not a path)\n", slug)
		return "", "", "", false, 2
	}
	workflowDir, code = resolveWorkflowDir(workflowDir, dir, stderr)
	return slug, workflowDir, msg, jsonOut, code
}

// parseWorkflowJSONArgs reads the `--workflow-dir DIR` + `--json` arg shape shared
// by `state ready` and `state sweep` (neither takes a positional). command names
// the verb in diagnostics.
func parseWorkflowJSONArgs(command string, args []string, dir string, stderr io.Writer) (workflowDir string, jsonOut bool, code int) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--workflow-dir":
			if i+1 >= len(args) {
				fmt.Fprintf(stderr, "spacedock %s: --workflow-dir requires a path\n", command)
				return "", false, 2
			}
			workflowDir = args[i+1]
			i++
		case "--json":
			jsonOut = true
		default:
			fmt.Fprintf(stderr, "spacedock %s: unknown argument %q\n", command, args[i])
			return "", false, 2
		}
	}
	workflowDir, code = resolveWorkflowDir(workflowDir, dir, stderr)
	return workflowDir, jsonOut, code
}

// resolveWorkflowDir resolves the workflow dir: an explicit relative --workflow-dir
// is joined against dir; an empty one is discovered from dir via the shared
// walk-up + downward-fallback resolver, the same one status/new/merge guard use —
// so all five state verbs agree on workflow discovery.
func resolveWorkflowDir(workflowDir, dir string, stderr io.Writer) (string, int) {
	if workflowDir == "" {
		return status.ResolveWorkflowDir(dir, stderr)
	}
	if !filepath.IsAbs(workflowDir) {
		return filepath.Join(dir, workflowDir), 0
	}
	return workflowDir, 0
}
