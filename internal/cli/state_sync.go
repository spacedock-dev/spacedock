// ABOUTME: `spacedock state commit <slug>` / `state ready` — path-scoped state
// ABOUTME: commit+push and boot-pull, each HALTing (exit 3) on a same-entity rebase conflict.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spacedock-dev/spacedock/internal/dispatch"
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

// runStateCommit implements `spacedock state commit <slug>`. It resolves <slug> to
// its entity commit unit under the split-root state checkout, then runs the path-scoped
// commit → push → on-reject pull --rebase → re-push sequence. A same-entity rebase
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
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "no-op",
			Reason: "Inline workflow — entities live beside the README; nothing to commit to a state checkout.",
		}, 0)
	}

	entityPath, ok := resolveEntityCommitPath(checkout, slug)
	if !ok {
		fmt.Fprintf(stderr, "spacedock state commit: no entity %q under %s (looked for %s.md and %s/index.md)\n", slug, checkout, slug, slug)
		return 1
	}
	if msg == "" {
		msg = fmt.Sprintf("state: update %s", slug)
	}

	// Path-scoped commit. A clean tree (nothing staged for this entity) is a no-op
	// success, not an error — the FO may call commit defensively after a --set that
	// changed nothing.
	if ok, out := commitEntityPathScoped(checkout, entityPath, msg); !ok {
		fmt.Fprintf(stderr, "spacedock state commit: git commit failed:\n%s\n", out)
		return 1
	} else if out == "" {
		// out=="" is the no-op sentinel from commitEntityPathScoped (nothing staged).
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "no-op", StateBranch: branch,
			Reason: fmt.Sprintf("Nothing to commit for %s — state checkout already up to date.", slug),
		}, 0)
	}

	// No origin → local-only success (mirrors the boot `remote: none` carve-out).
	if !stateHasOrigin(checkout) {
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "local-only", StateBranch: branch,
			Reason: fmt.Sprintf("Committed %s locally; no origin remote — state is local-only until an origin is configured.", slug),
		}, 0)
	}

	// Push; on a non-fast-forward rejection, pull --rebase and re-push.
	if ok, _ := runGit(checkout, "push", "origin", branch); ok {
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state commit", Slug: slug, Result: "pushed", StateBranch: branch,
			Reason: fmt.Sprintf("Committed and pushed %s to %s.", slug, branch),
		}, 0)
	}

	// Push rejected (non-ff or other). Integrate peers' state via pull --rebase,
	// reading git's OWN exit status (never piped through anything that swallows it).
	rebaseOK, rebaseOut := runGit(checkout, "pull", "--rebase", "origin", branch)
	if !rebaseOK {
		if rebaseInProgress(checkout) {
			return haltOnConflict(stdout, stderr, jsonOut, "state commit", slug, branch, checkout, entityPath, rebaseOut)
		}
		// Non-conflict pull failure (network, auth) — leave the local commit, report.
		fmt.Fprintf(stderr, "spacedock state commit: push rejected and pull --rebase failed (not a conflict):\n%s\n", rebaseOut)
		return 1
	}

	// Clean rebase replayed our commit atop the peer's — re-push.
	if ok, out := runGit(checkout, "push", "origin", branch); !ok {
		fmt.Fprintf(stderr, "spacedock state commit: re-push after pull --rebase failed:\n%s\n", out)
		return 1
	}
	return emitSync(stdout, jsonOut, syncResult{
		Command: "state commit", Slug: slug, Result: "pushed", StateBranch: branch,
		Reason: fmt.Sprintf("Committed %s, integrated peers' state, and pushed to %s.", slug, branch),
	}, 0)
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

	// No origin → ready, local-only (no network integration to do).
	if !stateHasOrigin(checkout) {
		return emitSync(stdout, jsonOut, syncResult{
			Command: "state ready", Result: "ready", StateBranch: branch,
			Reason: "State checkout ready (no origin remote — state is local-only).",
		}, 0)
	}

	// Integrate peers' state with one pull --rebase, reading git's own exit status.
	rebaseOK, rebaseOut := runGit(checkout, "pull", "--rebase", "origin", branch)
	if !rebaseOK {
		if rebaseInProgress(checkout) {
			return haltOnConflict(stdout, stderr, jsonOut, "state ready", "", branch, checkout, "", rebaseOut)
		}
		fmt.Fprintf(stderr, "spacedock state ready: pull --rebase failed (not a conflict):\n%s\n", rebaseOut)
		return 1
	}
	return emitSync(stdout, jsonOut, syncResult{
		Command: "state ready", Result: "ready", StateBranch: branch,
		Reason: fmt.Sprintf("State checkout ready — integrated peers' state from %s.", branch),
	}, 0)
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

// haltOnConflict is the shared HALT both verbs run on a same-entity rebase
// conflict: abort the rebase to restore a clean tree, NEVER force-push, NEVER
// auto-resolve, and exit 3 with the conflicting paths named on stderr. The exit
// code is the enforcement: a caller cannot proceed on an unmerged tree.
func haltOnConflict(stdout, stderr io.Writer, jsonOut bool, command, slug, branch, checkout, entityPath, rebaseOut string) int {
	conflicting := conflictingPaths(checkout)
	// The peer commit that survived: the pull's fetch phase already updated
	// origin/{branch} before the rebase conflicted, and --abort does not touch
	// that ref, so this resolves network-free (spiked in ideation).
	peerCommit := peerCommitSHA(checkout, branch)
	// Restore a clean tree so the next operation starts fresh. abort failure is
	// surfaced but the exit is still the halt.
	runGit(checkout, "rebase", "--abort")
	if len(conflicting) == 0 && entityPath != "" {
		// Fall back to the named entity when git's conflict list could not be read.
		conflicting = []string{relToCheckout(checkout, entityPath)}
	}
	fmt.Fprintf(stderr, "spacedock %s: HALT — same-entity rebase conflict on %s.\n", command, branch)
	fmt.Fprintf(stderr, "Conflicting path(s): %s\n", strings.Join(conflicting, ", "))
	if peerCommit != "" {
		fmt.Fprintf(stderr, "Peer commit: %s (origin/%s)\n", peerCommit, branch)
	}
	fmt.Fprintf(stderr, "The rebase was aborted (checkout left clean) and nothing was force-pushed; a peer's edit is preserved on origin.\n")
	fmt.Fprintf(stderr, "Next: HALT dispatch — do not dispatch against this state tree. Surface the conflicting path(s) and peer commit to the operator and stop.\n")
	fmt.Fprintf(stderr, "Never `git push --force`/`--force-with-lease`; never re-run with `-X ours`/`-X theirs`; never discard either side.\n")
	return emitSync(stdout, jsonOut, syncResult{
		Command: command, Slug: slug, Result: "halted", StateBranch: branch,
		ConflictingPaths: conflicting,
		PeerCommit:       peerCommit,
		Reason:           fmt.Sprintf("HALT: same-entity rebase conflict on %s — rebase aborted, nothing force-pushed, manual intervention required.", strings.Join(conflicting, ", ")),
	}, 3)
}

// peerCommitSHA resolves the peer's pushed commit on origin/{branch} — the edit
// preserved when this side's rebase conflicted. Runs BEFORE `rebase --abort`
// clears the conflict, though the ref itself does not depend on abort timing (the
// pull's fetch phase already updated it). Returns "" on any git failure rather
// than surfacing an error for what is purely diagnostic context.
func peerCommitSHA(checkout, branch string) string {
	if branch == "" {
		return ""
	}
	ok, out := runGit(checkout, "rev-parse", "--short", "origin/"+branch)
	if !ok {
		return ""
	}
	return strings.TrimSpace(out)
}

// commitEntityPathScoped stages and commits exactly entityPath with `add -A`
// restricted to its literal pathspec, retrying staging on index.lock contention.
// It returns (true, "") for a
// clean no-op (nothing staged for the entity → success), (true, output) for a
// landed commit, and (false, output) on a real git failure. The `git -C` argv
// form (via runGit) has no rendered command string for a weak model to paraphrase.
func commitEntityPathScoped(checkout, entityPath, msg string) (ok bool, output string) {
	rel := relToCheckout(checkout, entityPath)
	pathspec := literalGitPathspec(rel)
	if ok, out := runGitRetryLock(checkout, "add", "-A", "--", pathspec); !ok {
		return false, out
	}
	// Nothing staged for this entity → clean no-op success.
	if clean, _ := runGit(checkout, "diff", "--cached", "--quiet", "--", pathspec); clean {
		return true, ""
	}
	ok, out := runGitRetryLock(checkout, "commit", "-m", msg, "--", pathspec)
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

// rebaseInProgress reports whether a rebase is in progress in checkout, asking git
// for the worktree-correct path (`git rev-parse --git-path`) rather than assuming
// `<checkout>/.git/rebase-merge` — a linked-worktree state checkout keeps its git
// dir elsewhere, so a hardcoded path would miss the conflict and let the halt slip.
func rebaseInProgress(checkout string) bool {
	for _, name := range []string{"rebase-merge", "rebase-apply"} {
		ok, out := runGit(checkout, "rev-parse", "--git-path", name)
		if !ok {
			continue
		}
		p := strings.TrimSpace(out)
		if !filepath.IsAbs(p) {
			p = filepath.Join(checkout, p)
		}
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

// conflictingPaths returns the entity paths git reports as unmerged during the
// in-progress rebase (`git diff --name-only --diff-filter=U`), relative to the
// checkout. Read BEFORE the rebase --abort clears the conflict state.
func conflictingPaths(checkout string) []string {
	ok, out := runGit(checkout, "diff", "--name-only", "--diff-filter=U")
	if !ok {
		return nil
	}
	var paths []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			paths = append(paths, line)
		}
	}
	return paths
}

// resolveEntityCommitPath finds the commit unit for slug under checkout. Folder
// form wins when present, matching canonical entity discovery, and commits the
// whole folder; flat form commits only `{slug}.md`. A tracked candidate remains
// resolvable after deletion so the exact same commit unit can record its removal.
// EntitySlug remains the authority for the slug represented by either path.
func resolveEntityCommitPath(checkout, slug string) (string, bool) {
	nested := filepath.Join(checkout, slug, "index.md")
	if fileExists(nested) && status.EntitySlug(nested) == slug {
		return filepath.Dir(nested), true
	}
	flat := filepath.Join(checkout, slug+".md")
	if fileExists(flat) && status.EntitySlug(flat) == slug {
		return flat, true
	}
	if gitTracksPath(checkout, nested) && status.EntitySlug(nested) == slug {
		return filepath.Dir(nested), true
	}
	if gitTracksPath(checkout, flat) && status.EntitySlug(flat) == slug {
		return flat, true
	}
	return "", false
}

func gitTracksPath(checkout, path string) bool {
	ok, _ := runGit(checkout, "ls-files", "--error-unmatch", "--", literalGitPathspec(relToCheckout(checkout, path)))
	return ok
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

// stateHasOrigin reports whether the state checkout has a named `origin` remote,
// asking the exact named-remote question via `git remote get-url origin` (network-
// free, true iff exit 0). Mirrors status.stateHasOrigin so these verbs agree with
// boot's STATE_BACKEND origin line by construction. A non-repo dir or any git
// failure reports false → the verb degrades to local-only.
func stateHasOrigin(checkout string) bool {
	ok, _ := runGit(checkout, "remote", "get-url", "origin")
	return ok
}
