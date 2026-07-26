// ABOUTME: `spacedock merge guard <slug>` drives the terminal merge ceremony as one
// ABOUTME: ordered envelope: arm the mod-block, detect completion by state delta, clear, terminalize, archive.
package status

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/statesync"
)

// MergeGuard is the entry point for `spacedock merge guard <slug>`. It owns the
// ordered mod-block set->invoke->clear->terminalize sequence so a weak first
// officer cannot combine, skip, or reorder the steps. It does NOT invoke the
// merge hook, make the merge verdict (passed in via --verdict), or run the host
// agent-teardown; it emits the already-proven --set/--archive paths in the
// mandatory order and propagates — never bypasses — their guard refusals.
//
// args is the post-`guard` argv: the positional <slug>, --workflow-dir DIR,
// --verdict {passed|rejected}, and --json/--quiet. dir is the working directory
// used to resolve a relative --workflow-dir and to discover the enclosing
// workflow when no --workflow-dir is passed. The exit domain is {0,1,3}: 0 for a
// completed signal (armed / blocked / finalized), 1 for a usage error or a
// propagated guard/publication refusal, and 3 for an aborted state-rebase HALT.
func MergeGuard(args []string, dir string, stdout, stderr io.Writer) int {
	workflowDir, rest, err := parseWorkflowDir(args)
	if err != nil {
		return errExit(stderr, err.Error())
	}
	asJSON := contains(rest, "--json")
	quiet := contains(rest, "--quiet")

	verdict, err := parseSingleArg(rest, "--verdict", "passed|rejected")
	if err != nil {
		return errExit(stderr, err.Error())
	}
	// Verdict gate: a verdict-less finalize is structurally unreachable — the verb
	// requires the FO/captain's decision up front so it can never write a terminal
	// status without a verdict in the same mutation (the team-mode-verdict-omission
	// fix). Refuse BEFORE resolving the workflow so the omission is reported even on
	// a malformed workflow.
	if verdict == "" {
		return errExit(stderr, "merge guard requires --verdict {passed|rejected}")
	}
	if verdict != "passed" && verdict != "rejected" {
		return errExit(stderr, fmt.Sprintf("merge guard --verdict must be 'passed' or 'rejected', not '%s'", verdict))
	}

	slug := ""
	for i := 0; i < len(rest); i++ {
		a := rest[i]
		switch {
		case a == "--json" || a == "--quiet":
			continue
		case a == "--verdict":
			i++ // skip the flag's value, already parsed by parseSingleArg
			continue
		case strings.HasPrefix(a, "--"):
			return errExit(stderr, "unknown argument: "+a)
		case slug != "":
			return errExit(stderr, "merge guard takes a single <slug> argument")
		default:
			slug = a
		}
	}
	if slug == "" {
		return errExit(stderr, "merge guard requires a <slug> argument")
	}

	pipelineDir := workflowDir
	if pipelineDir == "" {
		resolved, rc := ResolveWorkflowDir(dir, stderr)
		if rc != 0 {
			return rc
		}
		pipelineDir = resolved
	}

	roots, err := resolveRoots(pipelineDir, dir)
	if err != nil {
		return errExit(stderr, err.Error())
	}

	if rc := mergeRootsGuard(workflowDir, roots, dir, stderr); rc != 0 {
		return rc
	}
	if rc := preflightMergeState(roots, slug, asJSON, stdout, stderr); rc != 0 {
		return rc
	}

	resolved, rc := resolveMutationEntity(roots, slug, stderr)
	if rc != 0 {
		return rc
	}
	slug = resolved.slug
	entityPath := resolveEntityPath(roots.entityDir, slug, stderr)
	if entityPath == "" {
		return errExit(stderr, "entity not found: "+slug)
	}

	// Validate the merge: policy up front so an invalid value (a workflow-config
	// typo) fails fast before any mutation. The classifier below no longer branches
	// on the policy — auto-arm and the merge-sentinel finalize apply under both
	// merge: local and merge: pr — but the parse still guards the config error.
	if _, perr := resolveMergePolicy(roots.definitionDir); perr != nil {
		return errExit(stderr, perr.Error())
	}

	fields := ParseFrontmatter(entityPath)
	modBlock := strings.TrimSpace(fields["mod-block"])
	pr := strings.TrimSpace(fields["pr"])
	worktree := strings.TrimSpace(fields["worktree"])
	mergeHooks := scanMods(roots.definitionDir)["merge"]
	hookRegistered := len(mergeHooks) > 0

	// State-delta classifier (the one genuinely-new logic): a pure read of
	// pr/mod-block/verdict plus the policy decides armed / blocked / finalize. No
	// new mutation primitive — every transition below emits a proven --set/--archive.
	switch {
	case verdict == "rejected":
		// A rejected entity never merged, so the pr-requirement is vacuous: finalize
		// straight through, clearing an in-flight mod-block standalone first (AC-6).
		return finalize(roots, slug, modBlock, pr, verdict, worktree, hookRegistered, quiet, asJSON, stdout, stderr)

	case prIndicatesMerged(pr):
		// FINALIZE from a detected-MERGED state. The `pr` field carries a merge
		// sentinel (pr-merge:{n} / local-merge:{sha}) — the LOCAL signal the FO's
		// pr-merge hook records when `gh` detects the PR MERGED. The verb never talks
		// to GitHub; it keys off this sentinel. This finalizes EVEN from a non-armed
		// (empty mod-block) state — the stranded case a re-validation bounce leaves
		// behind (AC-2). The merge-hook guard is satisfied because the sentinel is a
		// non-empty pr that honestly records the landed merge.
		return finalize(roots, slug, modBlock, pr, verdict, worktree, hookRegistered, quiet, asJSON, stdout, stderr)

	case modBlockNamesMissingMergeMod(modBlock, mergeHooks):
		// A mod-block naming a merge mod that no longer exists under _mods/, with no
		// merge sentinel recorded and a non-rejected verdict — the mid-flight state a
		// deleted mod file leaves behind. The default case below would otherwise clear
		// the block (standalone --set) and terminalize unguarded once no OTHER merge
		// hook happens to be registered, silently finalizing without the hook ever
		// having run. Refuse BEFORE any mutation — this case is checked ahead of both
		// the blocked and default cases so no --set runs.
		return refuseMissingMergeMod(roots.definitionDir, slug, modBlock, stderr)

	case pr != "":
		// Phase B blocked: a bare/open PR reference (e.g. #42) — the hook opened a PR
		// that has NOT merged. NEVER finalize on pr-presence alone; archiving here
		// would strand the task before its PR landed (the premature-finalize bug).
		// Leave mod-block + pr intact and wait for the merge sentinel.
		return signalBlocked(slug, pr, quiet, asJSON, stdout)

	case modBlock == "" && hookRegistered:
		// Phase A AUTO-ARM. Entering terminal with an empty mod-block and a merge hook
		// registered is the START of the merge ceremony, under BOTH merge: local AND
		// merge: pr (AC-1): the verb owns arming the mod-block (mod-block=merge:{hook})
		// and signals the FO to invoke the hook. The verb never invokes the hook or
		// local-merges; the FO runs it (merge: local does the --no-ff merge; merge: pr
		// opens the captain-gated PR). Ceremony integrity holds downstream: a merge: pr
		// entity can only FINALIZE once a merge sentinel records the landed PR, so an
		// arm-then-immediate-finalize is impossible without the hook running.
		return arm(roots, slug, mergeHooks[0], quiet, asJSON, stdout, stderr)

	default:
		// Phase C finalize. Under merge: local a cleared/clearable mod-block stands as
		// completion. Under merge: pr with no pr, no mod-block, and no hook registered
		// the terminalize --set is unguarded and succeeds. A merge: pr with a hook
		// registered never reaches here with an empty mod-block — auto-arm above claims
		// that state — so the merge-hook guard cannot strand a finalize.
		return finalize(roots, slug, modBlock, pr, verdict, worktree, hookRegistered, quiet, asJSON, stdout, stderr)
	}
}

// mergeRootsGuard fails MergeGuard closed, before entity resolution, when the
// resolved roots cannot possibly hold entities: the resolved definition dir does
// not exist, or a declared split-root state checkout is missing. Either shape
// today ends at the same misleading "entity not found: <slug>" — this names the
// real problem instead. workflowDirSpelling is the --workflow-dir value as
// literally passed (before defaulting via ResolveWorkflowDir), used both to name
// the as-passed spelling in the refusal and to gate the did-you-mean hint, which
// only fires for a relative spelling. A resolvable roots pair (or an inline
// workflow, which can never suffer the split-root shape) returns 0.
func mergeRootsGuard(workflowDirSpelling string, roots roots, dir string, stderr io.Writer) int {
	hint := didYouMeanHint(workflowDirSpelling, dir)

	if info, statErr := os.Stat(roots.definitionDir); statErr != nil || !info.IsDir() {
		msg := fmt.Sprintf(
			"merge guard: --workflow-dir %s resolves to %s (a relative --workflow-dir resolves against the current directory), which does not exist",
			workflowDirSpelling, roots.definitionDir)
		return errExit(stderr, appendHint(msg, hint))
	}

	mode, relPath, _ := ClassifyState(ParseFrontmatter(filepath.Join(roots.definitionDir, "README.md"))["state"])
	if mode != StateSplitRoot {
		return 0
	}
	if info, statErr := os.Stat(roots.entityDir); statErr != nil || !info.IsDir() {
		msg := fmt.Sprintf(
			"merge guard: workflow at %s declares state: %s but the state checkout is missing at %s",
			roots.definitionDir, relPath, roots.entityDir)
		return errExit(stderr, appendHint(msg, hint))
	}
	return 0
}

// preflightMergeState runs before entity resolution or mutation so a process
// restart with Git already mid-rebase cannot terminalize against an unmerged
// index. Inline workflows have no separate state branch to preflight.
func preflightMergeState(roots roots, slug string, asJSON bool, stdout, stderr io.Writer) int {
	mode, _, err := ClassifyState(ParseFrontmatter(filepath.Join(roots.definitionDir, "README.md"))["state"])
	if err != nil {
		return errExit(stderr, err.Error())
	}
	if mode == StateInline {
		return 0
	}
	branch, err := StateBranch(roots.definitionDir)
	if err != nil {
		return errExit(stderr, err.Error())
	}
	outcome := statesync.Preflight(roots.entityDir, branch)
	if outcome.Result == statesync.ResultHalted {
		return signalMergeStateHalt(slug, branch, outcome, asJSON, stdout, stderr)
	}
	if outcome.Result == statesync.ResultFailed {
		return errExit(stderr, "merge guard: state-sync preflight failed: "+outcome.Detail)
	}
	return 0
}

// didYouMeanHint computes the merge-guard corrective hint named in the roots
// guard's refusals: when workflowDirSpelling is a relative --workflow-dir, derive
// the enclosing main-checkout root as the parent of `git rev-parse
// --git-common-dir` run from dir, re-join the as-passed spelling there, and
// return the hint sentence only if that candidate validates as a resolvable
// workflow (via the same check --validate uses). A relative `--git-common-dir`
// output (the main-checkout-cwd case) resolves against dir, since git prints it
// relative to the invocation dir itself. Returns "" when the spelling is
// absent/absolute, git fails (non-repo cwd, bare), or the candidate does not
// validate — the refusal then stands without naming an unproven recovery path.
func didYouMeanHint(workflowDirSpelling, dir string) string {
	if workflowDirSpelling == "" || filepath.IsAbs(workflowDirSpelling) {
		return ""
	}
	out, err := runGitCmd(dir, "rev-parse", "--git-common-dir")
	if err != nil {
		return ""
	}
	commonDir := strings.TrimSpace(out)
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(dir, commonDir)
	}
	mainRoot := filepath.Dir(commonDir)

	candidate, rerr := resolveRoots(workflowDirSpelling, mainRoot)
	if rerr != nil {
		return ""
	}
	if rc := validateRootsOrExit(candidate, "", io.Discard); rc != 0 {
		return ""
	}
	return fmt.Sprintf(
		"did you mean --workflow-dir %s? (the current directory is a linked git worktree; a relative --workflow-dir resolves against it)",
		candidate.definitionDir)
}

// appendHint appends a did-you-mean hint on its own line when non-empty, else
// returns msg unchanged.
func appendHint(msg, hint string) string {
	if hint == "" {
		return msg
	}
	return msg + "\n" + hint
}

// prIndicatesMerged reports whether the pr field carries a WELL-FORMED merge
// sentinel — the LOCAL signal that a merge already LANDED, as opposed to a bare/open
// PR reference (#42, owner/repo#42, a URL) that is still in review. The pr-merge hook
// writes a `pr-merge:{number}` sentinel on MERGED detection; the no-PR local fallback
// writes `local-merge:{sha}`. A finalize+archive is irreversible, so the prefix match
// is not enough — the suffix must validate: a pr-merge sentinel finalizes only when
// its suffix parses as a positive PR number, and a local-merge sentinel only when its
// suffix is a non-empty SHA-like token. A bare reference, a garbage suffix, or an
// empty suffix returns false — the safe, fail-CLOSED direction, since finalizing on
// one would archive a task whose PR never landed.
func prIndicatesMerged(pr string) bool {
	pr = strings.TrimSpace(pr)
	if suffix := strings.TrimPrefix(pr, mergedPRSentinelPrefix); suffix != pr {
		n, err := strconv.Atoi(suffix)
		return err == nil && n > 0
	}
	if suffix := strings.TrimPrefix(pr, localMergeSentinelPrefix); suffix != pr {
		return isSHALike(suffix)
	}
	return false
}

// isSHALike reports whether s is a non-empty token of hex digits — the shape of the
// short merge-commit SHA the local-merge fallback records. It rejects an empty suffix
// and any non-hex character so a malformed local-merge sentinel does not finalize.
func isSHALike(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		isHex := (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
		if !isHex {
			return false
		}
	}
	return true
}

// mergedPRSentinelPrefix is the pr-field prefix the pr-merge hook writes on MERGED
// detection (`pr=pr-merge:{number}`) so the local state honestly records that the
// PR landed. It mirrors localMergeSentinelPrefix for the no-PR fallback.
const mergedPRSentinelPrefix = "pr-merge:"

// modBlockNamesMissingMergeMod reports whether modBlock is a merge-ceremony block
// (`merge:{name}`) naming a mod that is NOT among the currently registered merge
// hooks — a mod file deleted while its block was still in flight. An empty
// mod-block, or one naming a mod still present in mergeHooks, returns false.
func modBlockNamesMissingMergeMod(modBlock string, mergeHooks []string) bool {
	name, ok := strings.CutPrefix(modBlock, "merge:")
	if !ok || name == "" {
		return false
	}
	for _, hook := range mergeHooks {
		if hook == name {
			return false
		}
	}
	return true
}

// refuseMissingMergeMod signals D5: the mod-block names a merge mod missing from
// definitionDir/_mods/, with no merge sentinel recorded — the ceremony cannot
// resolve without the mod's guidance, and clearing the block would finalize
// without the hook ever having run. Exits 1, mutating nothing.
func refuseMissingMergeMod(definitionDir, slug, modBlock string, stderr io.Writer) int {
	name := strings.TrimPrefix(modBlock, "merge:")
	return errExit(stderr, fmt.Sprintf(
		"blocking mod %s is missing from %s/_mods/ — the entity %s is stuck. Restore the mod file, or have the operator clear the block with --force.",
		name, definitionDir, slug))
}

// arm performs Phase A: set mod-block=merge:{hook} in its own --set and signal the
// FO to invoke the hook. The underlying runSet is the proven mutation path.
func arm(roots roots, slug, hook string, quiet, asJSON bool, stdout, stderr io.Writer) int {
	modValue := "merge:" + hook
	if rc := emitSet(roots, slug, []fieldUpdate{{field: "mod-block", value: modValue, hasValue: true}}, stderr); rc != 0 {
		return rc
	}
	return signalArmed(roots.definitionDir, slug, hook, quiet, asJSON, stdout)
}

// finalize performs Phase C: clear an in-flight mod-block in a STANDALONE --set,
// terminalize status+verdict+completed in ONE --set, then archive. Each step is a
// proven guarded path; the verb refuses to proceed to a later step if an earlier
// one's guard refuses, propagating the guard's exit 1 + stderr verbatim and never
// passing --force.
//
// The archive move + commit is ATOMIC: the on-disk rename (runArchive) and its commit
// (commitArchiveMove) must both land or neither does. If the commit fails (a failing
// pre-commit hook, a lock), the entity would otherwise be stranded — mod-block
// cleared, status terminal, file moved into _archive where resolveEntityPath cannot
// find it on a re-run. So finalize snapshots the entity's pre-finalize bytes and live
// location before mutating, and on commit failure reverses the move and restores the
// original content, returning the entity to its exact pre-finalize state.
func finalize(roots roots, slug, modBlock, pr, verdict, worktree string, hookRegistered bool, quiet, asJSON bool, stdout, stderr io.Writer) int {
	// Snapshot the pre-finalize state up front — before any mutation — so a failed
	// archive commit can be rolled back to exactly the state the FO would re-run
	// against. The live path and form resolve here while the file still sits at its
	// live location; the bytes capture the untouched content.
	snapshot, snapErr := captureArchiveState(roots.entityDir, slug)
	if snapErr != nil {
		return errExit(stderr, fmt.Sprintf("merge guard: failed to snapshot %s before finalize: %v", slug, snapErr))
	}
	if modBlock != "" {
		if rc := emitSet(roots, slug, []fieldUpdate{{field: "mod-block", value: "", hasValue: true}}, stderr); rc != 0 {
			return rc
		}
	}
	terminal := terminalStageName(roots.definitionDir)
	if terminal == "" {
		return errExit(stderr, fmt.Sprintf("workflow %s declares no terminal stage — cannot finalize", roots.definitionDir))
	}
	terminalize := []fieldUpdate{
		{field: "status", value: terminal, hasValue: true},
		{field: "verdict", value: verdict, hasValue: true},
		{field: "completed", hasValue: false},
	}
	if rc := emitSet(roots, slug, terminalize, stderr); rc != 0 {
		return rc
	}
	if rc := runArchive(roots.definitionDir, roots.entityDir, roots.entityDirSpelling, slug, false, true, false, io.Discard, stderr); rc != 0 {
		return rc
	}
	if rc := commitArchiveMove(roots.entityDir, slug, snapshot, stderr); rc != 0 {
		if rbErr := rollbackArchive(roots.entityDir, slug, snapshot); rbErr != nil {
			// The commit failed AND the rollback could not fully restore the
			// pre-finalize state — the entity may be half-mutated with a partial archive
			// or a leaked index entry. This needs human attention, not a silent rc=1, so
			// say so unmistakably and name the entity and its expected live path.
			fmt.Fprintf(stderr,
				"merge guard: CRITICAL — archive commit for %s failed and rollback did NOT fully restore it: %v\n"+
					"  the entity may be in an inconsistent state (expected live at %s). "+
					"Inspect `git status`/`git diff --cached` and the entity location before re-running.\n",
				slug, rbErr, snapshot.livePath)
		}
		return rc
	}
	durability, syncOutcome, rc := publishMergeArchive(roots, slug, stdout, stderr)
	if durability == "" {
		durability = "unpublished"
	}
	signalFinalized(roots.definitionDir, slug, terminal, verdict, worktree, hookRegistered, prIndicatesMerged(pr), durability, syncOutcome, quiet, asJSON, stdout)
	return rc
}

func publishMergeArchive(roots roots, slug string, stdout, stderr io.Writer) (string, statesync.Outcome, int) {
	mode, _, err := ClassifyState(ParseFrontmatter(filepath.Join(roots.definitionDir, "README.md"))["state"])
	if err != nil {
		return "", statesync.Outcome{}, errExit(stderr, err.Error())
	}
	if mode == StateInline {
		return "inline", statesync.Outcome{}, 0
	}
	branch, err := StateBranch(roots.definitionDir)
	if err != nil {
		return "", statesync.Outcome{}, errExit(stderr, err.Error())
	}
	outcome := statesync.Publish(roots.entityDir, branch)
	switch outcome.Result {
	case statesync.ResultPushed, statesync.ResultNoOp:
		return "pushed", outcome, 0
	case statesync.ResultLocalOnly:
		return "local-only", outcome, 0
	case statesync.ResultHalted:
		return "halted", outcome, signalMergeStateHalt(slug, branch, outcome, false, stdout, stderr)
	case statesync.ResultFailed:
		fmt.Fprintf(stderr,
			"merge guard: archive commit for %s is durable locally, but split-root publication failed:\n%s\n"+
				"Settle unrelated state-checkout dirt if Git refused rebase, then resume with `spacedock state commit %s --workflow-dir %s`; archived resume publishes the existing commit without creating another archive commit.\n",
			slug, outcome.Detail, slug, roots.definitionDir)
		return "unpublished", outcome, 1
	default:
		fmt.Fprintf(stderr, "merge guard: archive commit for %s was not published (unexpected state-sync result %s)\n", slug, outcome.Result)
		return "unpublished", outcome, 1
	}
}

func signalMergeStateHalt(slug, branch string, outcome statesync.Outcome, asJSON bool, stdout, stderr io.Writer) int {
	paths := strings.Join(outcome.ConflictingPaths, ", ")
	if paths == "" {
		paths = "none reported by Git"
	}
	fmt.Fprintf(stderr, "merge guard: HALT — same-entity rebase conflict on %s.\n", branch)
	fmt.Fprintf(stderr, "Conflicting path(s): %s\n", paths)
	if outcome.PeerCommit != "" {
		fmt.Fprintf(stderr, "Peer commit: %s (origin/%s)\n", outcome.PeerCommit, branch)
	}
	fmt.Fprintln(stderr, "The rebase was aborted and nothing was force-pushed; the local archive commit and peer edit remain recoverable.")
	fmt.Fprintln(stderr, "Next: HALT dispatch — surface the conflicting paths and peer commit to the operator; never force-push or auto-resolve.")
	if asJSON {
		emitJSON(stdout, newJSONObj().
			set("command", "merge-guard").set("slug", slug).set("signal", "halted").
			set("result", "halted").set("state_branch", branch).
			setValue("conflicting_paths", jsonStrArr(outcome.ConflictingPaths)).set("peer_commit", outcome.PeerCommit).
			set("reason", "HALT: rebase aborted; manual conflict resolution required."))
	}
	return 3
}

// archiveSnapshot captures an entity's pre-archive state so finalize can reverse a
// failed archive commit: the live source path (where the entity sat before the
// rename), whether it was folder-form, the flat companion's filesystem and Git
// membership, the source file's exact bytes, and its mode.
type archiveSnapshot struct {
	livePath         string
	isFolder         bool
	companionPresent bool
	companionTracked bool
	content          []byte
	mode             os.FileMode
}

// captureArchiveState reads an entity's pre-archive state. It resolves the live path
// the way runArchive does (folder form wins over flat) and snapshots the source
// file's bytes and mode for a byte- and mode-faithful restore.
func captureArchiveState(entityDir, slug string) (archiveSnapshot, error) {
	folderIndex := filepath.Join(entityDir, slug, "index.md")
	flatPath := filepath.Join(entityDir, slug+".md")
	livePath := flatPath
	isFolder := false
	if isRegularFile(folderIndex) {
		livePath = folderIndex
		isFolder = true
	}
	info, err := os.Stat(livePath)
	if err != nil {
		return archiveSnapshot{}, err
	}
	content, err := os.ReadFile(livePath)
	if err != nil {
		return archiveSnapshot{}, err
	}
	companion := filepath.Join(entityDir, slug)
	gitRoot := FindGitRoot(entityDir)
	return archiveSnapshot{
		livePath:         livePath,
		isFolder:         isFolder,
		companionPresent: !isFolder && directoryExists(companion),
		companionTracked: !isFolder && trackedArchivePath(gitRoot, companion),
		content:          content,
		mode:             info.Mode().Perm(),
	}, nil
}

// rollbackArchive reverses runArchive AND its index staging after a failed commit. It
// must return the entity to its EXACT pre-finalize state on every axis the commit
// touched: the working tree (move the archived entity back, restore the pre-archive
// bytes+mode) AND the git index (unstage the rename `commitArchiveMove` staged before
// the commit failed). Without the unstage the index keeps a phantom staged rename to
// _archive that a later plain commit would sweep into HEAD — committing the entity
// ONLY at _archive and orphaning the live file — and that breaks the recovery re-run's
// `git add` (exit 128). The flat form moves _archive/{slug}.md → {slug}.md; the folder
// form moves the whole _archive/{slug}/ → {slug}/.
//
// It attempts every step even if an earlier one fails, joining the errors, so a
// partial failure surfaces all of what went wrong rather than masking later steps.
func rollbackArchive(entityDir, slug string, snap archiveSnapshot) error {
	var errs []error

	if snap.isFolder {
		archivedFolder := filepath.Join(entityDir, "_archive", slug)
		liveFolder := filepath.Join(entityDir, slug)
		if err := os.Rename(archivedFolder, liveFolder); err != nil {
			errs = append(errs, fmt.Errorf("reverse folder rename: %w", err))
		}
	} else {
		archivedFile := filepath.Join(entityDir, "_archive", slug+".md")
		if err := os.Rename(archivedFile, snap.livePath); err != nil {
			errs = append(errs, fmt.Errorf("reverse file rename: %w", err))
		}
		if snap.companionPresent {
			archivedCompanion := filepath.Join(entityDir, "_archive", slug)
			liveCompanion := filepath.Join(entityDir, slug)
			if err := os.Rename(archivedCompanion, liveCompanion); err != nil {
				errs = append(errs, fmt.Errorf("reverse companion rename: %w", err))
			}
		}
	}
	// Restore the pre-archive content+mode only if the file is back at its live path
	// (a failed reverse-rename leaves nothing to write to).
	if isRegularFile(snap.livePath) {
		if err := os.WriteFile(snap.livePath, snap.content, snap.mode); err != nil {
			errs = append(errs, fmt.Errorf("restore content: %w", err))
		}
	}
	// Unstage the leaked rename — the same pathspecs commitArchiveMove staged, under
	// the same git-worktree guard. `git reset -- <paths>` only touches the index, not
	// the working tree we just restored.
	if gitRoot := FindGitRoot(entityDir); hasGitEntry(gitRoot) {
		pathspecs := archiveMovePathspecs(
			gitRoot,
			entityDir,
			slug,
			snap.isFolder,
			snap.companionTracked,
			snap.companionPresent,
		)
		args := append([]string{"reset", "-q", "--"}, pathspecs...)
		if _, err := runGitCmd(gitRoot, args...); err != nil {
			errs = append(errs, fmt.Errorf("unstage archive rename: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// commitArchiveMove commits the archive rename PATH-SCOPED so the verb — not the FO
// — owns the commit (AC-3). It stages exactly the paths the rename touches (the
// vacated source and the new _archive dest) and commits only those, so a sibling
// entity left dirty in the same tree is never swept in by a `git add -A`. The
// pathspecs mirror runArchive's form resolution: a flat-form entity moves
// {slug}.md → _archive/{slug}.md (two file paths); a folder-form entity moves the
// whole {slug}/ → _archive/{slug}/ (two directory paths, capturing index.md and any
// siblings). The form is read from where the entity landed in _archive, since the
// source has already been moved by the time this runs. The pre-archive snapshot
// retains a tracked-but-deleted flat companion as a source-only pathspec. When the
// entity root is not under a git work tree the move is a plain on-disk archive and
// the commit is skipped (no error). A real git failure exits 1 with stderr.
func commitArchiveMove(entityDir, slug string, snapshot archiveSnapshot, stderr io.Writer) int {
	gitRoot := FindGitRoot(entityDir)
	if !hasGitEntry(gitRoot) {
		return 0
	}
	isFolder := archivedAsFolder(entityDir, slug)
	pathspecs := archiveMovePathspecs(
		gitRoot,
		entityDir,
		slug,
		isFolder,
		!isFolder && snapshot.companionTracked,
		!isFolder && snapshot.companionPresent,
	)
	// Stage the vacated source (deletion) and the new dest. git records this as a
	// rename in the commit. `archiveMovePathspecs` returns literal pathspecs so
	// valid entity names that look like Git magic stay exact; --all is never used.
	addArgs := append([]string{"add", "--"}, pathspecs...)
	if _, err := runGitCmd(gitRoot, addArgs...); err != nil {
		return errExit(stderr, fmt.Sprintf("merge guard: failed to stage archive move for %s: %v", slug, err))
	}
	commitArgs := append([]string{"commit", "-q", "-m", "archive " + slug + " (merge guard)", "--"}, pathspecs...)
	if _, err := runGitCmd(gitRoot, commitArgs...); err != nil {
		return errExit(stderr, fmt.Sprintf("merge guard: failed to commit archive move for %s: %v", slug, err))
	}
	return 0
}

// archiveMovePathspecs returns the source and dest pathspecs (relative to gitRoot)
// for the archive rename, mirroring runArchive's flat/folder resolution. A folder-
// form entity renames the whole {slug}/ directory; a flat-form entity renames the
// {slug}.md file. Flat companion source and destination membership are independent
// so a tracked source deletion can be committed without naming an absent archive
// destination. The caller supplies the form explicitly so the pathspecs are stable
// regardless of where the entity currently sits on disk.
func archiveMovePathspecs(gitRoot, entityDir, slug string, isFolder, companionSource, companionDest bool) []string {
	if isFolder {
		return []string{
			literalArchivePathspec(relToGitRoot(gitRoot, filepath.Join(entityDir, slug))),
			literalArchivePathspec(relToGitRoot(gitRoot, filepath.Join(entityDir, "_archive", slug))),
		}
	}
	pathspecs := []string{
		literalArchivePathspec(relToGitRoot(gitRoot, filepath.Join(entityDir, slug+".md"))),
		literalArchivePathspec(relToGitRoot(gitRoot, filepath.Join(entityDir, "_archive", slug+".md"))),
	}
	if companionSource {
		pathspecs = append(pathspecs, literalArchivePathspec(relToGitRoot(gitRoot, filepath.Join(entityDir, slug))))
	}
	if companionDest {
		pathspecs = append(pathspecs, literalArchivePathspec(relToGitRoot(gitRoot, filepath.Join(entityDir, "_archive", slug))))
	}
	return pathspecs
}

func literalArchivePathspec(path string) string {
	return ":(literal)" + filepath.ToSlash(path)
}

func trackedArchivePath(gitRoot, path string) bool {
	if !hasGitEntry(gitRoot) {
		return false
	}
	out, err := runGitCmd(gitRoot, "ls-files", "--", literalArchivePathspec(relToGitRoot(gitRoot, path)))
	return err == nil && strings.TrimSpace(out) != ""
}

// archivedAsFolder reports whether the entity landed in _archive as a folder
// (_archive/{slug}/index.md), used by commitArchiveMove to pick the rename pathspecs
// from where the entity sits post-move.
func archivedAsFolder(entityDir, slug string) bool {
	return isRegularFile(filepath.Join(entityDir, "_archive", slug, "index.md"))
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// relToGitRoot renders path relative to gitRoot for a path-scoped `git add`/`commit`
// pathspec. It falls back to the absolute path when the relativization fails so the
// pathspec still resolves.
func relToGitRoot(gitRoot, path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	rel, err := filepath.Rel(gitRoot, abs)
	if err != nil {
		return abs
	}
	return rel
}

// emitSet runs one proven --set mutation through runSet, discarding its success
// narration (the verb emits its own phase signal) while propagating a guard
// refusal's stderr verbatim. It never passes --force.
func emitSet(roots roots, slug string, updates []fieldUpdate, stderr io.Writer) int {
	set := &setUpdate{slug: slug, updates: updates}
	return runSet(roots, set, nil, nil,
		false, false, false, false, false, false, true, false,
		io.Discard, stderr)
}

// terminalStageName returns the workflow's terminal stage name, or "" when the
// README declares none. The verb terminalizes into this stage by name rather than
// hardcoding "done", so a workflow with a differently-named terminal stage works.
func terminalStageName(definitionDir string) string {
	readme := definitionDir + "/README.md"
	if !fileExists(readme) {
		return ""
	}
	for _, s := range parseStagesBlock(readme) {
		if s.terminal {
			return s.Name
		}
	}
	return ""
}

// signalArmed reports the arm outcome: the mod-block is set and the FO must now
// invoke the named merge hook. The default prose names the hook's file path and
// the re-run command — the FO's next action, carried at fire time instead of
// resident prose (D3).
func signalArmed(definitionDir, slug, hook string, quiet, asJSON bool, stdout io.Writer) int {
	switch {
	case asJSON:
		emitJSON(stdout, newJSONObj().
			set("command", "merge-guard").set("slug", slug).
			set("signal", "armed").set("action", "invoke-hook").set("hook", hook))
	case quiet:
		fmt.Fprintf(stdout, "merge-guard slug=%s signal=armed action=invoke-hook hook=%s\n", slug, hook)
	default:
		fmt.Fprintf(stdout, "armed: mod-block set to merge:%s — invoke the %s merge hook (%s/_mods/%s.md; merge guard never invokes it), then re-run `merge guard %s`.\n",
			hook, hook, definitionDir, hook, slug)
	}
	return 0
}

// signalBlocked reports the blocked outcome: the hook opened a PR, so the verb
// leaves the mod-block intact and waits. The default prose names the never-
// finalize-on-open-PR invariant and the sentinel format that unlocks finalize
// (D3), carried at fire time instead of resident prose.
func signalBlocked(slug, pr string, quiet, asJSON bool, stdout io.Writer) int {
	switch {
	case asJSON:
		emitJSON(stdout, newJSONObj().
			set("command", "merge-guard").set("slug", slug).
			set("signal", "blocked").set("action", "await-pr").set("pr", pr))
	case quiet:
		fmt.Fprintf(stdout, "merge-guard slug=%s signal=blocked action=await-pr pr=%s\n", slug, pr)
	default:
		fmt.Fprintf(stdout, "blocked: PR %s is pending — mod-block left intact, never finalize on an open PR. "+
			"When gh reports it MERGED, record the sentinel (pr=pr-merge:{number}) and re-run `merge guard %s`.\n", pr, slug)
	}
	return 0
}

// signalFinalized reports the finalize outcome: the entity is terminal, verdict
// recorded, archived. The default prose appends the FO's next-step lines (D3),
// built from the pre-terminalize frontmatter: a recorded worktree names its own
// removal/branch-cleanup/teardown sequence; an entity finalizing with no merge
// hook registered and no merge sentinel (the default-local-merge path) also
// names the manual `--no-ff` merge onto trunk that nothing automated.
func signalFinalized(definitionDir, slug, terminal, verdict, worktree string, hookRegistered, hasSentinel bool, durability string, syncOutcome statesync.Outcome, quiet, asJSON bool, stdout io.Writer) int {
	base := fmt.Sprintf("finalized: %s -> %s (verdict %s), archived.", slug, terminal, verdict)
	switch durability {
	case "pushed":
		base += " State durability: pushed to the split-root origin."
	case "local-only":
		base += " State durability: local-only (no origin remote)."
	case "unpublished":
		base += " State durability: unpublished; the local archive commit is retained."
	case "halted":
		base += " State durability: HALT after the local archive commit; remote publication requires manual conflict resolution."
	}
	var next []string
	if worktree != "" {
		next = append(next, fmt.Sprintf(
			"Next: push; remove the worktree (`git worktree remove %s`, no --force — if it refuses on untracked files, audit them with the operator before any --force); "+
				"delete the local branch (`git branch -d`); keep the remote branch while a PR references it; tear down the entity's workers per your runtime adapter.",
			worktree))
	}
	if !hookRegistered && !hasSentinel && verdict != "rejected" {
		branch := "the stage branch"
		if worktree != "" {
			if ref, ok := abbrevRefHead(worktree); ok && ref != "" {
				branch = ref
			}
		}
		next = append(next, fmt.Sprintf(
			"no merge hook registered — merge %s onto %s with --no-ff if not already merged.", branch, resolveMergeTrunk(definitionDir)))
	}
	switch {
	case asJSON:
		doc := newJSONObj().
			set("command", "merge-guard").set("slug", slug).
			set("signal", "finalized").set("status", terminal).set("verdict", verdict).set("result", durability)
		if durability == "halted" {
			doc.setValue("conflicting_paths", jsonStrArr(syncOutcome.ConflictingPaths)).
				set("peer_commit", syncOutcome.PeerCommit).
				set("reason", "HALT: local archive committed; rebase aborted; manual conflict resolution required.")
		}
		emitJSON(stdout, doc)
	case quiet:
		fmt.Fprintf(stdout, "merge-guard slug=%s signal=finalized status=%s verdict=%s result=%s\n", slug, terminal, verdict, durability)
	default:
		fmt.Fprintln(stdout, base)
		for _, line := range next {
			fmt.Fprintln(stdout, line)
		}
	}
	return 0
}

// resolveMergeTrunk mirrors dispatch.resolveTrunk's resolution (the top-level
// README `trunk:` key, falling back to "main"). Duplicated rather than imported —
// internal/dispatch already imports internal/status, so the reverse import would
// cycle.
func resolveMergeTrunk(definitionDir string) string {
	if t := strings.TrimSpace(ParseFrontmatter(filepath.Join(definitionDir, "README.md"))["trunk"]); t != "" {
		return t
	}
	return "main"
}

// abbrevRefHead best-effort resolves worktree's current branch name via
// `git -C {worktree} rev-parse --abbrev-ref HEAD`. ok is false on any git failure
// (detached HEAD, a removed worktree) — the caller falls back to a branch-
// agnostic phrase rather than surfacing an error for a purely cosmetic name.
func abbrevRefHead(worktree string) (ref string, ok bool) {
	out, err := runGitCmd(worktree, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", false
	}
	ref = strings.TrimSpace(out)
	return ref, ref != ""
}
