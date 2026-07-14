// ABOUTME: `spacedock state init` resumes a cloned split-root workflow; `state new`
// ABOUTME: births the orphan state branch + linked worktree around a present README.
package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// Resume lifecycle seams let real-Git tests place non-locking commands around
// private restoration and public checkout publication. Production leaves them
// nil.
var stateResumeBeforeRestoreHook func(string)
var stateResumeBeforePublishHook func(string)
var stateResumeAfterRegistrationRepairHook func(string)
var stateResumeAfterPublishHook func(string)
var stateResumeAfterReadyHook func(string)
var stateResumeWaitHook func(string)

// runStateInit implements `spacedock state init`. It reads the workflow's
// `state:` and `state-branch:` from the README, and for a split-root workflow
// whose state checkout is ABSENT, fetches the orphan state branch from origin and
// checks it out as a linked worktree at the state path. The path-exists guard
// makes a second run a no-op (a raw 2nd `git worktree add` fatals "already
// exists"). An inline/empty workflow has nothing to init.
func runStateInit(ctx context.Context, args []string, env []string, dir string, stdout, stderr io.Writer) int {
	workflowDir, code := parseStateInitArgs(args, dir, stderr)
	if code != 0 {
		return code
	}

	readme := filepath.Join(workflowDir, "README.md")
	if !fileExists(readme) {
		fmt.Fprintf(stderr, "spacedock state init: no README.md at %s\n", workflowDir)
		return 1
	}
	fm := status.ParseFrontmatter(readme)
	mode, relPath, err := status.ClassifyState(fm["state"])
	if err != nil {
		fmt.Fprintf(stderr, "spacedock state init: %v\n", err)
		return 1
	}
	if mode == status.StateInline {
		fmt.Fprintln(stdout, "Inline workflow — entities live beside the README; nothing to init.")
		return 0
	}

	branch, err := status.StateBranch(workflowDir)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock state init: %v\n", err)
		return 1
	}
	statePath, err := resolveSplitRootCheckout(workflowDir, relPath)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock state init: cannot resolve state checkout: %v\n", err)
		return 1
	}

	observedAbsent := !dirExists(statePath)
	code, lockErr := withStateResumeLock(workflowDir, statePath, func(waitedForState bool) int {
		if dirExists(statePath) {
			if err := validateExistingStateCheckout(workflowDir, statePath, branch); err != nil {
				fmt.Fprintf(stderr, "spacedock state init: refusing invalid state checkout: %v\n", err)
				return 1
			}
		}
		// A caller that observed absence, or waited on this checkout's lock, follows
		// the preceding convergence result. The per-checkout lock distinguishes this
		// from waiting behind an unrelated workflow's repository-wide worktree work.
		if dirExists(statePath) {
			if observedAbsent || waitedForState {
				if result, err := readStateResumeOutcome(workflowDir, statePath); err != nil || result != "ready" {
					fmt.Fprintln(stderr, "spacedock state init: concurrent state resume did not converge")
					return 1
				} else {
					fmt.Fprintf(stdout, "State checkout already initialized at %s (branch %s).\n", statePath, branch)
					return 0
				}
			}
			if err := writeStateResumeOutcome(workflowDir, statePath, "pending"); err != nil {
				fmt.Fprintf(stderr, "spacedock state init: cannot record convergence outcome: %v\n", err)
				return 1
			}
			if !observedAbsent {
				if fetchOK, _ := runGit(statePath, "fetch", "origin", branch); fetchOK {
					if pullOK, pullOut := runGit(statePath, "pull", "--rebase", "origin", branch); !pullOK {
						if rebaseInProgress(statePath) {
							runGit(statePath, "rebase", "--abort")
						}
						writeStateResumeOutcome(workflowDir, statePath, "failed")
						fmt.Fprintf(stderr, "spacedock state init: pull --rebase failed:\n%s\n", pullOut)
						return 1
					}
				}
			}
			if err := writeStateResumeOutcome(workflowDir, statePath, "ready"); err != nil {
				fmt.Fprintf(stderr, "spacedock state init: cannot commit convergence outcome: %v\n", err)
				return 1
			}
			fmt.Fprintf(stdout, "State checkout already initialized at %s (branch %s).\n", statePath, branch)
			return 0
		}
		if err := writeStateResumeOutcome(workflowDir, statePath, "pending"); err != nil {
			fmt.Fprintf(stderr, "spacedock state init: cannot record resume outcome: %v\n", err)
			return 1
		}
		resumeCode, _ := resumeAbsentSplitRootCheckoutLocked(workflowDir, branch, statePath, stdout, stderr)
		if resumeCode != 0 {
			writeStateResumeOutcome(workflowDir, statePath, "failed")
			return resumeCode
		}
		if err := writeStateResumeOutcome(workflowDir, statePath, "ready"); err != nil {
			fmt.Fprintf(stderr, "spacedock state init: cannot commit resume outcome: %v\n", err)
			return 1
		}
		if stateResumeAfterReadyHook != nil {
			stateResumeAfterReadyHook(statePath)
		}
		return resumeCode
	})
	if lockErr != nil {
		fmt.Fprintf(stderr, "spacedock state init: cannot lock state resume: %v\n", lockErr)
		return 1
	}
	return code
}

// resumeAbsentSplitRootCheckoutLocked implements the RESUME half of `state init`,
// shared with `state ready`'s absent-checkout resume. Callers confirm statePath's
// directory is ABSENT and hold the repository resume lock before calling. It
// converges the state branch in a unique adjacent private worktree, then renames
// the fully ready directory into statePath and repairs its Git metadata. The
// public path is therefore never visible with a pre-convergence tree.
//
// originConverged reports whether this resume fetched and integrated origin —
// false for a local-only or unreachable-origin fallback.
func resumeAbsentSplitRootCheckoutLocked(workflowDir, branch, statePath string, stdout, stderr io.Writer) (code int, originConverged bool) {
	if stateResumeBeforeRestoreHook != nil {
		stateResumeBeforeRestoreHook(statePath)
	}
	if dirExists(statePath) {
		fmt.Fprintf(stderr, "spacedock state init: state checkout appeared concurrently at %s; left it untouched\n", statePath)
		return 1, false
	}

	hasOrigin := stateHasOrigin(workflowDir)
	var fetchOK bool
	var fetchOut string
	if hasOrigin {
		fetchOK, fetchOut = runGit(workflowDir, "fetch", "origin", branch)
	}
	localBranchExists, _ := runGit(workflowDir, "rev-parse", "--verify", "--quiet", "refs/heads/"+branch)

	switch {
	case fetchOK:
		privatePath, stale, err := addPrivateStateWorktree(workflowDir, statePath, branch)
		if err != nil {
			fmt.Fprintf(stderr, "spacedock state init: %v\n", err)
			return 1, false
		}
		privateActive := true
		defer func() {
			if privateActive {
				if ok, out := runGit(workflowDir, "worktree", "remove", privatePath); !ok {
					fmt.Fprintf(stderr, "spacedock state init: preserving private resume checkout %s after cleanup failure:\n%s\n", privatePath, out)
				}
			}
		}()
		if ok, out := runGit(privatePath, "pull", "--rebase", "origin", branch); !ok {
			if rebaseInProgress(privatePath) {
				runGit(privatePath, "rebase", "--abort")
			}
			fmt.Fprintf(stderr, "spacedock state init: pull --rebase origin %s failed in private checkout:\n%s\n", branch, out)
			return 1, false
		}
		published, err := publishPrivateStateWorktree(workflowDir, privatePath, statePath, branch, stale)
		if published {
			privateActive = false
		}
		if err != nil {
			fmt.Fprintf(stderr, "spacedock state init: %v\n", err)
			return 1, false
		}
		if stale {
			fmt.Fprintf(stdout, "Repaired stale worktree registration for %s.\n", statePath)
		}
		fmt.Fprintf(stdout, "Initialized state checkout at %s (branch %s).\n", statePath, branch)
		return 0, true

	case localBranchExists:
		privatePath, stale, err := addPrivateStateWorktree(workflowDir, statePath, branch)
		if err != nil {
			fmt.Fprintf(stderr, "spacedock state init: %v\n", err)
			return 1, false
		}
		privateActive := true
		defer func() {
			if privateActive {
				if ok, out := runGit(workflowDir, "worktree", "remove", privatePath); !ok {
					fmt.Fprintf(stderr, "spacedock state init: preserving private resume checkout %s after cleanup failure:\n%s\n", privatePath, out)
				}
			}
		}()
		published, err := publishPrivateStateWorktree(workflowDir, privatePath, statePath, branch, stale)
		if published {
			privateActive = false
		}
		if err != nil {
			fmt.Fprintf(stderr, "spacedock state init: %v\n", err)
			return 1, false
		}
		if stale {
			fmt.Fprintf(stdout, "Repaired stale worktree registration for %s.\n", statePath)
		}
		if !hasOrigin {
			fmt.Fprintf(stdout, "Initialized state checkout at %s (branch %s) — no origin remote, state is local-only.\n", statePath, branch)
		} else {
			fmt.Fprintf(stdout, "Warning: git fetch origin %s failed; resumed from the local branch — peers' state may be missing:\n%s\n", branch, fetchOut)
			fmt.Fprintf(stdout, "Initialized state checkout at %s (branch %s).\n", statePath, branch)
		}
		return 0, false

	default:
		// Neither a fetchable origin branch nor a local branch. Disambiguate
		// "never birthed" (no way the branch exists anywhere) from "origin
		// unreachable" (indeterminate — the branch may exist, we just can't see
		// it) so the hint doesn't send an operator to `state new` when a peer's
		// birth is merely unreachable right now.
		reachable, found := false, false
		if hasOrigin {
			reachable, found = remoteBranchStatus(workflowDir, branch)
		}
		neverBirthed := !found && (!hasOrigin || reachable)
		if neverBirthed {
			fmt.Fprintf(stderr, "spacedock state init: state branch %s not found locally or on origin at %s — the workflow has not been birthed here. Run `spacedock state new --workflow-dir %s`.\n",
				branch, statePath, workflowDir)
		} else {
			fmt.Fprintf(stderr, "spacedock state init: git fetch origin %s failed:\n%s\n"+
				"Manual fallback: git fetch origin %s && git worktree add %s %s\n",
				branch, fetchOut, branch, statePath, branch)
		}
		return 1, false
	}
}

// parseStateInitArgs reads `--workflow-dir DIR`, resolving a relative path
// against dir. With no flag it discovers the enclosing workflow from dir.
func parseStateInitArgs(args []string, dir string, stderr io.Writer) (workflowDir string, code int) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--workflow-dir":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "spacedock state init: --workflow-dir requires a path")
				return "", 2
			}
			workflowDir = args[i+1]
			i++
		default:
			fmt.Fprintf(stderr, "spacedock state init: unknown argument %q\n", args[i])
			return "", 2
		}
	}
	return resolveWorkflowDir(workflowDir, dir, stderr)
}

// runStateNew implements `spacedock state new`. It is the BIRTH half of the
// split-root lifecycle and the mechanical inverse of `state init`: reading the same
// README `state:`/`state-branch:`, it appends the `.gitignore` entry, births the
// orphan state branch, and checks it out as a linked worktree at the state path.
// The README is a precondition (commission's interactive job), not an output. The
// prechecks refuse an already-birthed or already-remote workflow instead of letting
// git fatal, and an inline README is an operator error (a mutating verb cannot
// no-op the way `state init` does). The orphan push is best-effort: a repo with no
// origin births locally and warns the operator to push later.
func runStateNew(ctx context.Context, args []string, env []string, dir string, stdout, stderr io.Writer) int {
	workflowDir, code := parseStateNewArgs(args, dir, stderr)
	if code != 0 {
		return code
	}

	readme := filepath.Join(workflowDir, "README.md")
	if !fileExists(readme) {
		fmt.Fprintf(stderr, "spacedock state new: no README.md at %s\n", workflowDir)
		return 1
	}
	fm := status.ParseFrontmatter(readme)
	mode, relPath, err := status.ClassifyState(fm["state"])
	if err != nil {
		fmt.Fprintf(stderr, "spacedock state new: %v\n", err)
		return 1
	}
	if mode == status.StateInline {
		fmt.Fprintf(stderr, "spacedock state new: %s is an inline workflow — `state new` only births split-root workflows (set `state:` to a checkout path in the README first).\n", workflowDir)
		return 1
	}

	branch, err := status.StateBranch(workflowDir)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock state new: %v\n", err)
		return 1
	}
	statePath, err := resolveSplitRootCheckout(workflowDir, relPath)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock state new: cannot resolve state checkout: %v\n", err)
		return 1
	}

	// Refusal prechecks (AC-3): an occupied path or an existing orphan (local or
	// remote) means the workflow is already birthed; refuse with a tailored message
	// rather than letting git fatal. The remote case redirects to `state init`.
	if dirExists(statePath) {
		fmt.Fprintf(stderr, "spacedock state new: %s already present — the workflow is already birthed here. Run `spacedock status` to inspect it.\n", statePath)
		return 1
	}
	if ok, _ := runGit(workflowDir, "rev-parse", "--verify", "refs/heads/"+branch); ok {
		fmt.Fprintf(stderr, "spacedock state new: orphan branch %s already present in local refs — the workflow is already birthed. Run `spacedock state init` to check it out at %s.\n", branch, statePath)
		return 1
	}
	if orphanOnRemote(workflowDir, branch) {
		fmt.Fprintf(stderr, "spacedock state new: orphan branch %s already present on origin — this workflow was birthed elsewhere. Run `spacedock state init --workflow-dir %s` to resume it.\n", branch, workflowDir)
		return 1
	}

	// Birth: append the gitignore entry (idempotent), then create the orphan branch
	// + linked worktree.
	if err := appendGitignoreEntry(workflowDir, relPath); err != nil {
		fmt.Fprintf(stderr, "spacedock state new: %v\n", err)
		return 1
	}
	pushed, err := birthOrphanState(workflowDir, statePath, branch)
	if err != nil {
		fmt.Fprintf(stderr, "spacedock state new: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Birthed split-root state at %s (branch %s).\n", statePath, branch)
	fmt.Fprintf(stdout, "Next: commit the .gitignore edit on your code branch (git add .gitignore && git commit).\n")
	if !pushed {
		fmt.Fprintf(stdout, "Warning: could not push %s to origin (no remote or push failed). Push the orphan branch later to share state: git -C %s push origin %s.\n", branch, statePath, branch)
	}
	return 0
}

// orphanOnRemote reports whether branch exists on origin, parsing `git ls-remote
// --heads origin {branch}`. A missing origin or any ls-remote failure reports false
// — the local birth still proceeds best-effort.
func orphanOnRemote(workflowDir, branch string) bool {
	ok, out := runGit(workflowDir, "ls-remote", "--heads", "origin", branch)
	return ok && lsRemoteHasBranch(out, branch)
}

// remoteBranchStatus reports whether `git ls-remote --heads origin {branch}`
// could reach origin at all (reachable) and, if so, whether branch was found
// there (found). Unlike orphanOnRemote, an unreachable origin is distinguished
// from a reachable one with no matching branch: `ls-remote` exits non-zero on an
// unreachable/misconfigured remote but exits 0 with empty output when reachable
// with no such branch — so reachable=false means the caller cannot conclude
// anything about the branch's existence, only that it couldn't check.
func remoteBranchStatus(workflowDir, branch string) (reachable, found bool) {
	ok, out := runGit(workflowDir, "ls-remote", "--heads", "origin", branch)
	if !ok {
		return false, false
	}
	return true, lsRemoteHasBranch(out, branch)
}

// lsRemoteHasBranch reports whether `git ls-remote --heads` output names branch.
// ls-remote prints one `<sha>\trefs/heads/<branch>` line per matching head, so a
// `refs/heads/<branch>` reference in any line is a match. Empty output (no such
// head, the common case) is no match. Branch names can contain `/`
// (spacedock-state/dev), so the whole `refs/heads/<branch>` ref-path is matched,
// not a bare suffix.
func lsRemoteHasBranch(lsRemoteOutput, branch string) bool {
	ref := "refs/heads/" + branch
	for _, line := range strings.Split(lsRemoteOutput, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == ref {
			return true
		}
	}
	return false
}

// appendGitignoreEntry idempotently appends `{relPath}/` to the code branch's
// tracked .gitignore (repo root), so collaborators and fresh clones ignore the
// state checkout. The entry is left uncommitted for the operator to commit
// alongside their README — committing on the operator's code branch is an
// outward-facing side effect the operator should own.
func appendGitignoreEntry(workflowDir, relPath string) error {
	root, err := repoRoot(workflowDir)
	if err != nil {
		return err
	}
	placement, err := status.ResolveRepositoryPlacement(workflowDir)
	if err != nil {
		return err
	}
	// The .gitignore entry is the workflow's repo-relative prefix + the state
	// relPath, so a docs/dev workflow ignores docs/dev/.spacedock-state/. Derive the
	// prefix from git's `--show-prefix` rather than filepath.Rel against the
	// toplevel: git resolves symlinks in --show-toplevel (e.g. /var → /private/var
	// on macOS) while workflowDir is the caller's unresolved path, so Rel between the
	// two divergent prefixes yields a garbage `../../…` path.
	prefix, err := repoRelPrefix(workflowDir)
	if err != nil {
		return err
	}
	entry := filepath.ToSlash(filepath.Join(prefix, relPath)) + "/"
	// The checkout is physically created under the main worktree even when this
	// command runs from a linked worktree. The tracked .gitignore edit belongs to
	// the invoking code branch, but until that edit is merged the main worktree
	// would otherwise see the state checkout as untracked content. Mirror the
	// physical path into the shared repository exclude before creating it.
	if placement.Linked {
		exclude := filepath.Join(placement.CommonGitDir, "info", "exclude")
		if err := appendIgnoreEntry(exclude, "/"+entry); err != nil {
			return fmt.Errorf("updating repository exclude: %w", err)
		}
	}

	gitignore := filepath.Join(root, ".gitignore")
	return appendIgnoreEntry(gitignore, entry)
}

func appendIgnoreEntry(path, entry string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, line := range strings.Split(string(existing), "\n") {
		if strings.TrimSpace(line) == entry {
			return nil // already present — idempotent
		}
	}
	body := string(existing)
	if len(body) > 0 && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	body += entry + "\n"
	return os.WriteFile(path, []byte(body), 0o644)
}

// birthOrphanState creates the orphan branch in a temp detached worktree (clearing
// the inherited tree so the orphan starts empty), seeds an empty commit, pushes
// best-effort to origin, then checks the branch out as a linked worktree at
// statePath. Returns whether the orphan branch reached origin. `worktree add
// --detach` tolerates a dirty code tree (it is a separate checkout), so the birth
// does not require a clean working tree.
func birthOrphanState(workflowDir, statePath, branch string) (pushed bool, err error) {
	root, err := repoRoot(workflowDir)
	if err != nil {
		return false, err
	}
	tmpWT, err := os.MkdirTemp("", "spacedock-orphan-birth-")
	if err != nil {
		return false, fmt.Errorf("creating temp worktree dir: %w", err)
	}
	// os.MkdirTemp creates the dir, but `git worktree add` needs a non-existent
	// path; remove the placeholder so git owns the checkout.
	if err := os.Remove(tmpWT); err != nil {
		return false, fmt.Errorf("preparing temp worktree path: %w", err)
	}
	defer runGit(root, "worktree", "remove", "--force", tmpWT)

	if ok, out := runGit(root, "worktree", "add", "--detach", tmpWT); !ok {
		return false, fmt.Errorf("git worktree add --detach failed:\n%s", out)
	}
	if ok, out := runGit(tmpWT, "checkout", "--orphan", branch); !ok {
		return false, fmt.Errorf("git checkout --orphan %s failed:\n%s", branch, out)
	}
	// Clear the inherited index and working tree so the orphan branch starts empty.
	if ok, out := runGit(tmpWT, "rm", "-rf", "--cached", "."); !ok {
		return false, fmt.Errorf("git rm --cached failed:\n%s", out)
	}
	if err := clearWorktreeFiles(tmpWT); err != nil {
		return false, err
	}
	if ok, out := runGit(tmpWT, "commit", "-q", "--allow-empty", "-m", "seed state"); !ok {
		return false, fmt.Errorf("git commit (seed state) failed:\n%s", out)
	}
	pushed, _ = runGit(tmpWT, "push", "origin", branch)

	// Remove the temp worktree before adding the real one — git refuses a 2nd
	// worktree on the same branch.
	if ok, out := runGit(root, "worktree", "remove", "--force", tmpWT); !ok {
		return pushed, fmt.Errorf("git worktree remove failed:\n%s", out)
	}
	if ok, out := runGit(root, "worktree", "add", statePath, branch); !ok {
		return pushed, fmt.Errorf("git worktree add %s %s failed:\n%s", statePath, branch, out)
	}
	return pushed, nil
}

// clearWorktreeFiles removes every entry in dir except .git, leaving the orphan
// branch's working tree empty.
func clearWorktreeFiles(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Name() == ".git" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(dir, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

// repoRoot returns the absolute path of the git repository containing workflowDir.
func repoRoot(workflowDir string) (string, error) {
	ok, out := runGit(workflowDir, "rev-parse", "--show-toplevel")
	if !ok {
		return "", fmt.Errorf("not a git repository at %s:\n%s", workflowDir, out)
	}
	return status.TrimGitLineTerminator(out), nil
}

// repoRelPrefix returns workflowDir's path relative to its repo root, as git sees
// it (`git rev-parse --show-prefix`) — a trailing-slashed path for a subdir, empty
// when workflowDir is the repo root. Unlike filepath.Rel against --show-toplevel,
// this is symlink-resolution-safe (both sides are git's own view).
func repoRelPrefix(workflowDir string) (string, error) {
	ok, out := runGit(workflowDir, "rev-parse", "--show-prefix")
	if !ok {
		return "", fmt.Errorf("not a git repository at %s:\n%s", workflowDir, out)
	}
	return status.TrimGitLineTerminator(out), nil
}

// parseStateNewArgs reads `--workflow-dir DIR`, resolving a relative path against
// dir. With no flag it discovers the enclosing workflow from dir. (Mirrors
// parseStateInitArgs — the two commands take the same arguments.)
func parseStateNewArgs(args []string, dir string, stderr io.Writer) (workflowDir string, code int) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--workflow-dir":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "spacedock state new: --workflow-dir requires a path")
				return "", 2
			}
			workflowDir = args[i+1]
			i++
		default:
			fmt.Fprintf(stderr, "spacedock state new: unknown argument %q\n", args[i])
			return "", 2
		}
	}
	return resolveWorkflowDir(workflowDir, dir, stderr)
}

// runGit runs a git command in dir, returning success and combined output. On
// failure it does NOT print — the caller decides whether the failure is fatal (a
// fresh fetch) or tolerable (a refresh pull on an already-initialized checkout).
func runGit(dir string, gitArgs ...string) (bool, string) {
	cmd := exec.Command("git", append([]string{"-C", dir}, gitArgs...)...)
	out, err := cmd.CombinedOutput()
	return err == nil, string(out)
}

// fileExists reports whether path is an existing regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// dirExists reports whether path is an existing directory.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// addPrivateStateWorktree checks out branch beside statePath without publishing
// the final path. A stale final-path registration requires --force solely to let
// the same branch enter this private worktree; the stale record remains until
// publication.
func addPrivateStateWorktree(workflowDir, statePath, branch string) (privatePath string, stale bool, err error) {
	if dirExists(statePath) {
		return "", false, fmt.Errorf("state checkout appeared concurrently at %s; left it untouched", statePath)
	}
	ok, out := runGit(workflowDir, "worktree", "list", "--porcelain", "-z")
	if !ok {
		return "", false, fmt.Errorf("listing worktree registrations failed:\n%s", out)
	}
	records, err := status.ParseWorktreePorcelainZ([]byte(out))
	if err != nil {
		return "", false, fmt.Errorf("parsing worktree registrations: %w", err)
	}
	stale, err = classifyStaleStateRegistration(records, statePath, branch)
	if err != nil {
		return "", false, err
	}
	privatePath, err = os.MkdirTemp(filepath.Dir(statePath), ".spacedock-state-resume-")
	if err != nil {
		return "", false, fmt.Errorf("creating private resume path: %w", err)
	}
	if err := os.Remove(privatePath); err != nil {
		return "", false, fmt.Errorf("preparing private resume path: %w", err)
	}
	args := []string{"worktree", "add"}
	if stale {
		args = append(args, "--force")
	}
	args = append(args, privatePath, branch)
	if ok, out := runGit(workflowDir, args...); !ok {
		return "", false, fmt.Errorf("git worktree add private checkout for %s failed:\n%s", statePath, out)
	}
	return privatePath, stale, nil
}

func classifyStaleStateRegistration(records []status.WorktreeRecord, statePath, branch string) (bool, error) {
	target := status.RealpathOf(statePath)
	var matches []status.WorktreeRecord
	for _, record := range records {
		if status.RealpathOf(record.Path) == target {
			matches = append(matches, record)
		}
	}
	if len(matches) == 0 {
		return false, nil
	}
	if len(matches) != 1 {
		return false, fmt.Errorf("state checkout %s has %d registrations; refusing ambiguous stale repair", statePath, len(matches))
	}
	record := matches[0]
	wantBranch := "refs/heads/" + branch
	if record.Branch != wantBranch {
		return false, fmt.Errorf("state checkout %s registration is on %q, expected %q; refusing stale repair", statePath, record.Branch, wantBranch)
	}
	if record.Bare || !record.Prunable {
		return false, fmt.Errorf("state checkout %s registration is active or not prunable; refusing stale repair", statePath)
	}
	return true, nil
}

func removeStaleStateRegistration(workflowDir, statePath, branch string) error {
	if _, err := os.Lstat(statePath); err == nil {
		return fmt.Errorf("state checkout appeared concurrently at %s; left it untouched", statePath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking state checkout before registration repair %s: %w", statePath, err)
	}
	ok, out := runGit(workflowDir, "worktree", "list", "--porcelain", "-z")
	if !ok {
		return fmt.Errorf("listing worktree registrations before repair failed:\n%s", out)
	}
	records, err := status.ParseWorktreePorcelainZ([]byte(out))
	if err != nil {
		return fmt.Errorf("parsing worktree registrations before repair: %w", err)
	}
	stale, err := classifyStaleStateRegistration(records, statePath, branch)
	if err != nil {
		return err
	}
	if !stale {
		return fmt.Errorf("state checkout %s no longer has the expected prunable registration; refusing stale repair", statePath)
	}
	if ok, out := runGit(workflowDir, "worktree", "remove", statePath); !ok {
		return fmt.Errorf("removing stale state worktree registration failed:\n%s", out)
	}
	return nil
}

// publishPrivateStateWorktree atomically publishes a fully converged private
// directory without replacing any destination entry. published reports that
// the rename occurred even
// if subsequent metadata repair/validation fails, so callers never clean the
// now-public directory through its former private name.
func publishPrivateStateWorktree(workflowDir, privatePath, statePath, branch string, stale bool) (published bool, err error) {
	if stateResumeBeforePublishHook != nil {
		stateResumeBeforePublishHook(statePath)
	}
	if stale {
		if err := removeStaleStateRegistration(workflowDir, statePath, branch); err != nil {
			return false, err
		}
	}
	if stateResumeAfterRegistrationRepairHook != nil {
		stateResumeAfterRegistrationRepairHook(statePath)
	}
	if err := renameNoReplace(privatePath, statePath); err != nil {
		return false, fmt.Errorf("publishing converged state checkout at %s: %w", statePath, err)
	}
	published = true
	if ok, out := runGit(workflowDir, "worktree", "repair", statePath); !ok {
		return true, fmt.Errorf("repairing published state checkout metadata failed:\n%s", out)
	}
	if err := validateExistingStateCheckout(workflowDir, statePath, branch); err != nil {
		return true, fmt.Errorf("published state checkout validation failed: %w", err)
	}
	if stateResumeAfterPublishHook != nil {
		stateResumeAfterPublishHook(statePath)
	}
	return true, nil
}
