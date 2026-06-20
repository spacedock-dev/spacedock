// ABOUTME: `spacedock merge guard <slug>` drives the terminal merge ceremony as one
// ABOUTME: ordered envelope: arm the mod-block, detect completion by state delta, clear, terminalize, archive.
package status

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
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
// workflow when no --workflow-dir is passed. The exit domain is {0,1}: 0 for a
// completed signal (armed / blocked / finalized), 1 for a usage error or a
// propagated guard refusal.
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
		if discovered, ok := DiscoverWorkflowDir(dir); ok {
			pipelineDir = discovered
		} else if resolved, rc := discoverWorkflowDownward(dir, stderr); rc != 0 {
			return rc
		} else {
			pipelineDir = resolved
		}
	}

	roots, err := resolveRoots(pipelineDir, dir)
	if err != nil {
		return errExit(stderr, err.Error())
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
	mergeHooks := scanMods(roots.definitionDir)["merge"]
	hookRegistered := len(mergeHooks) > 0

	// State-delta classifier (the one genuinely-new logic): a pure read of
	// pr/mod-block/verdict plus the policy decides armed / blocked / finalize. No
	// new mutation primitive — every transition below emits a proven --set/--archive.
	switch {
	case verdict == "rejected":
		// A rejected entity never merged, so the pr-requirement is vacuous: finalize
		// straight through, clearing an in-flight mod-block standalone first (AC-6).
		return finalize(roots, slug, modBlock, verdict, quiet, asJSON, stdout, stderr)

	case prIndicatesMerged(pr):
		// FINALIZE from a detected-MERGED state. The `pr` field carries a merge
		// sentinel (pr-merge:{n} / local-merge:{sha}) — the LOCAL signal the FO's
		// pr-merge hook records when `gh` detects the PR MERGED. The verb never talks
		// to GitHub; it keys off this sentinel. This finalizes EVEN from a non-armed
		// (empty mod-block) state — the stranded case a re-validation bounce leaves
		// behind (AC-2). The merge-hook guard is satisfied because the sentinel is a
		// non-empty pr that honestly records the landed merge.
		return finalize(roots, slug, modBlock, verdict, quiet, asJSON, stdout, stderr)

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
		return finalize(roots, slug, modBlock, verdict, quiet, asJSON, stdout, stderr)
	}
}

// prIndicatesMerged reports whether the pr field carries a merge sentinel — the
// LOCAL signal that a merge already LANDED, as opposed to a bare/open PR reference
// (#42, owner/repo#42, a URL) that is still in review. The pr-merge hook writes a
// `pr-merge:{number}` sentinel on MERGED detection; the no-PR local fallback writes
// `local-merge:{sha}`. Either honestly records a shipped merge, so the verb may
// finalize. A bare reference does not — finalizing on it would archive a task
// before its PR landed.
func prIndicatesMerged(pr string) bool {
	pr = strings.TrimSpace(pr)
	return strings.HasPrefix(pr, mergedPRSentinelPrefix) ||
		strings.HasPrefix(pr, localMergeSentinelPrefix)
}

// mergedPRSentinelPrefix is the pr-field prefix the pr-merge hook writes on MERGED
// detection (`pr=pr-merge:{number}`) so the local state honestly records that the
// PR landed. It mirrors localMergeSentinelPrefix for the no-PR fallback.
const mergedPRSentinelPrefix = "pr-merge:"

// arm performs Phase A: set mod-block=merge:{hook} in its own --set and signal the
// FO to invoke the hook. The underlying runSet is the proven mutation path.
func arm(roots roots, slug, hook string, quiet, asJSON bool, stdout, stderr io.Writer) int {
	modValue := "merge:" + hook
	if rc := emitSet(roots, slug, []fieldUpdate{{field: "mod-block", value: modValue, hasValue: true}}, stderr); rc != 0 {
		return rc
	}
	return signalArmed(slug, hook, quiet, asJSON, stdout)
}

// finalize performs Phase C: clear an in-flight mod-block in a STANDALONE --set,
// terminalize status+verdict+completed in ONE --set, then archive. Each step is a
// proven guarded path; the verb refuses to proceed to a later step if an earlier
// one's guard refuses, propagating the guard's exit 1 + stderr verbatim and never
// passing --force.
func finalize(roots roots, slug, modBlock, verdict string, quiet, asJSON bool, stdout, stderr io.Writer) int {
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
	if rc := commitArchiveMove(roots.entityDir, slug, stderr); rc != 0 {
		return rc
	}
	return signalFinalized(slug, terminal, verdict, quiet, asJSON, stdout)
}

// commitArchiveMove commits the archive rename PATH-SCOPED so the verb — not the FO
// — owns the commit (AC-3). It stages exactly the two paths the rename touches (the
// vacated source and the new _archive dest) and commits only those, so a sibling
// entity left dirty in the same tree is never swept in by a `git add -A`. When the
// entity root is not under a git work tree the move is a plain on-disk archive and
// the commit is skipped (no error). A real git failure exits 1 with stderr.
func commitArchiveMove(entityDir, slug string, stderr io.Writer) int {
	gitRoot := FindGitRoot(entityDir)
	if !hasGitEntry(gitRoot) {
		return 0
	}
	source := relToGitRoot(gitRoot, entityDir+"/"+slug+".md")
	dest := relToGitRoot(gitRoot, entityDir+"/_archive/"+slug+".md")
	// Stage the vacated source (deletion) and the new dest. git records this as a
	// rename in the commit. `git add -- <path>` is path-scoped; --all is never used.
	if _, err := runGitCmd(gitRoot, "add", "--", source, dest); err != nil {
		return errExit(stderr, fmt.Sprintf("merge guard: failed to stage archive move for %s: %v", slug, err))
	}
	if _, err := runGitCmd(gitRoot, "commit", "-q", "-m", "archive "+slug+" (merge guard)", "--", source, dest); err != nil {
		return errExit(stderr, fmt.Sprintf("merge guard: failed to commit archive move for %s: %v", slug, err))
	}
	return 0
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
// invoke the named merge hook.
func signalArmed(slug, hook string, quiet, asJSON bool, stdout io.Writer) int {
	switch {
	case asJSON:
		emitJSON(stdout, newJSONObj().
			set("command", "merge-guard").set("slug", slug).
			set("signal", "armed").set("action", "invoke-hook").set("hook", hook))
	case quiet:
		fmt.Fprintf(stdout, "merge-guard slug=%s signal=armed action=invoke-hook hook=%s\n", slug, hook)
	default:
		fmt.Fprintf(stdout, "armed: mod-block set to merge:%s — invoke the %s merge hook, then re-run `merge guard %s`.\n", hook, hook, slug)
	}
	return 0
}

// signalBlocked reports the blocked outcome: the hook opened a PR, so the verb
// leaves the mod-block intact and waits.
func signalBlocked(slug, pr string, quiet, asJSON bool, stdout io.Writer) int {
	switch {
	case asJSON:
		emitJSON(stdout, newJSONObj().
			set("command", "merge-guard").set("slug", slug).
			set("signal", "blocked").set("action", "await-pr").set("pr", pr))
	case quiet:
		fmt.Fprintf(stdout, "merge-guard slug=%s signal=blocked action=await-pr pr=%s\n", slug, pr)
	default:
		fmt.Fprintf(stdout, "blocked: PR %s is pending — mod-block left intact. Re-run `merge guard %s` after it lands.\n", pr, slug)
	}
	return 0
}

// signalFinalized reports the finalize outcome: the entity is terminal, verdict
// recorded, archived.
func signalFinalized(slug, terminal, verdict string, quiet, asJSON bool, stdout io.Writer) int {
	switch {
	case asJSON:
		emitJSON(stdout, newJSONObj().
			set("command", "merge-guard").set("slug", slug).
			set("signal", "finalized").set("status", terminal).set("verdict", verdict))
	case quiet:
		fmt.Fprintf(stdout, "merge-guard slug=%s signal=finalized status=%s verdict=%s\n", slug, terminal, verdict)
	default:
		fmt.Fprintf(stdout, "finalized: %s -> %s (verdict %s), archived.\n", slug, terminal, verdict)
	}
	return 0
}
