// ABOUTME: --set mutation flow, the read (table/next/boot/validate) flow, and
// ABOUTME: workflow discovery / git helpers backing the native runner.
package status

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/claudeteam"
)

// runSet handles the --set mutation flow with mod-block / merge-hook /
// terminal-transition guards and the field: old -> new narration. Matches the
// --set branch of main().
func runSet(roots roots, set *setUpdate, args []string, whereFilters []whereFilter,
	includeArchive, showNext, showBoot, showNextID, showValidate, hasFieldsFlag, quiet, asJSON bool,
	stdout, stderr io.Writer) int {

	var incompatible []string
	if showNext {
		incompatible = append(incompatible, "--next")
	}
	if includeArchive {
		incompatible = append(incompatible, "--archived")
	}
	if showBoot {
		incompatible = append(incompatible, "--boot")
	}
	if showNextID {
		incompatible = append(incompatible, "--next-id")
	}
	if len(whereFilters) > 0 {
		incompatible = append(incompatible, "--where")
	}
	if hasFieldsFlag {
		incompatible = append(incompatible, "--fields/--all-fields")
	}
	if showValidate {
		incompatible = append(incompatible, "--validate")
	}
	if len(incompatible) > 0 {
		return errExit(stderr, "--set cannot be combined with "+strings.Join(incompatible, ", "))
	}

	force := contains(args, "--force")
	resolved, rc := resolveMutationEntity(roots, set.slug, stderr)
	if rc != 0 {
		return rc
	}
	slug := resolved.slug
	mainEntityPath := resolveEntityPath(roots.entityDir, slug, stderr)
	if mainEntityPath == "" {
		return errExit(stderr, "entity not found: "+slug)
	}
	// Single-root stage: no worktree overlay, entity path is the main path.
	entityPath := mainEntityPath

	currentFields := ParseFrontmatter(entityPath)
	modBlock := strings.TrimSpace(currentFields["mod-block"])
	currentPR := strings.TrimSpace(currentFields["pr"])
	currentVerdict := strings.TrimSpace(currentFields["verdict"])
	clearingModBlock := false
	for _, u := range set.updates {
		if u.field == "mod-block" && u.hasValue && u.value == "" {
			clearingModBlock = true
		}
	}

	readme := filepath.Join(roots.definitionDir, "README.md")
	var stages []Stage
	if fileExists(readme) {
		stages = parseStagesBlock(readme)
	}
	terminalNames := map[string]bool{}
	for _, s := range stages {
		if s.terminal {
			terminalNames[s.Name] = true
		}
	}

	// Membership guard: a --set status=X must name a stage the workflow
	// declares. This is distinct from validateWorkflowStageNames' format regex —
	// it checks the *value* against stages.states[].name membership. When the
	// workflow declares no stages block we cannot validate membership, so the
	// guard is a no-op. A non-member value exits non-zero and leaves the
	// frontmatter unchanged.
	if len(stages) > 0 {
		stageNames := make([]string, len(stages))
		stageNameSet := map[string]bool{}
		for i, s := range stages {
			stageNames[i] = s.Name
			stageNameSet[s.Name] = true
		}
		for _, u := range set.updates {
			if u.field == "status" && u.hasValue && !stageNameSet[u.value] {
				return errExit(stderr, fmt.Sprintf(
					"'%s' is not a defined stage in workflow %s — known stages: [%s]",
					u.value, roots.definitionDir, strings.Join(stageNames, ", ")))
			}
		}
	}

	isTerminalUpdate := func() bool {
		for _, u := range set.updates {
			if u.field == "status" && u.hasValue && terminalNames[u.value] {
				return true
			}
			if u.field == "completed" || u.field == "verdict" {
				return true
			}
			if u.field == "worktree" && u.hasValue && u.value == "" {
				return true
			}
		}
		return false
	}

	postUpdatePR := currentPR
	for _, u := range set.updates {
		if u.field == "pr" {
			postUpdatePR = strings.TrimSpace(u.value)
		}
	}
	postUpdateVerdict := currentVerdict
	for _, u := range set.updates {
		if u.field == "verdict" {
			postUpdateVerdict = strings.TrimSpace(u.value)
		}
	}
	// finalizing reports whether this --set performs the FINALIZE action: setting
	// `completed` from empty to non-empty. `completed` is a timestamp field, so a
	// bare `completed` (no value) auto-fills now() when currently empty — both the
	// bare and the valued form finalize. This is the contract terminalize shape
	// (shared-core step 7: `completed verdict={verdict} worktree=`); it is the
	// signal the verdict gate keys on, NOT `status==terminal` alone, so a bare
	// dispatch-into-terminal (`status=done started`, no `completed`) is not gated.
	finalizing := false
	for _, u := range set.updates {
		if u.field != "completed" {
			continue
		}
		newCompleted := u.value
		if !u.hasValue { // bare timestamp auto-fills now() only when currently empty
			newCompleted = "filled"
		}
		if strings.TrimSpace(currentFields["completed"]) == "" && strings.TrimSpace(newCompleted) != "" {
			finalizing = true
		}
	}

	if (modBlock != "" || clearingModBlock) && !force {
		if isTerminalUpdate() {
			var reason string
			if modBlock != "" && !clearingModBlock {
				reason = fmt.Sprintf("pending mod-block (%s)", modBlock)
			} else if clearingModBlock {
				reason = fmt.Sprintf("combined mod-block clear with terminal transition (current mod-block='%s')", modBlock)
			} else {
				reason = "mod-block transition"
			}
			return errExit(stderr, fmt.Sprintf(
				"entity %s has %s. Clear mod-block in a separate --set call, or use --force.",
				slug, reason))
		}
	}

	// Merge-hook invariant: when the workflow registers a merge hook, refuse a
	// terminal transition that has skipped it (empty pr AND empty mod-block).
	// Under `merge: local` the workflow has declared it merges locally, so the
	// pr-requirement is exempted — a cleared mod-block stands as proof the
	// ceremony ran. This relaxes the guard CHECK only; the mod-block-pending and
	// combined-clear-with-terminal guards above are policy-independent, so the
	// mandatory set->invoke->clear ceremony structure is unaffected.
	policy, perr := resolveMergePolicy(roots.definitionDir)
	if perr != nil {
		return errExit(stderr, perr.Error())
	}

	// Verdict gate: a finalized entity carries a verdict, so refuse the FINALIZE
	// action (setting `completed`) with no verdict (neither already set nor set in
	// the same call). This is the mechanism behind shared-core's terminalize step
	// (`completed verdict={verdict} worktree=`): without it an FO can finalize and
	// leave verdict empty. It keys on the finalize action, NOT on `status==terminal`
	// alone — a bare dispatch-into-terminal (`status=done started`, no `completed`)
	// is a legitimate verdict-less transition (the verdict is the outcome of work
	// that has not happened yet), so it passes; gating on status alone wrongly
	// blocked the dispatch step of a backlog→done workflow and made the live cycle
	// flaky. --force is the same escape hatch the mod-block / merge-hook guards
	// honor. Policy-independent and composes with those guards. Placed AFTER the
	// merge-policy parse so an invalid `merge:` value (a workflow-config error) is
	// still reported first.
	if !force && finalizing && postUpdateVerdict == "" {
		return errExit(stderr, fmt.Sprintf(
			"entity %s cannot be finalized (`completed`) without a verdict. "+
				"Set verdict in the same --set call, or use --force.",
			slug))
	}
	if !force && policy != mergeLocal && isTerminalUpdate() && modBlock == "" && postUpdatePR == "" && postUpdateVerdict != "rejected" {
		mergeHooks := scanMods(roots.definitionDir)["merge"]
		if len(mergeHooks) > 0 {
			return errExit(stderr, fmt.Sprintf(
				"entity %s cannot advance to terminal — workflow has merge hook(s) [%s] that have not run "+
					"(pr field is empty and mod-block is empty). Set mod-block=merge:%s and invoke the hook, or use --force to bypass.",
				slug, strings.Join(mergeHooks, ", "), mergeHooks[0]))
		}
	}

	if force && modBlock != "" {
		fmt.Fprintf(stderr, "Warning: --force overriding mod-block (%s) on entity %s\n", modBlock, slug)
	}

	// require-external-proof guard: when the workflow opts in, a terminal --set
	// on an entity whose ACs are only self-referentially proven exits 1 and
	// leaves the frontmatter untouched. Layered after mod-block + merge-hook so
	// the order of guard messages mirrors the order of declared invariants.
	// --force bypasses with a warning in the same idiom as the mod-block bypass.
	proofPolicy, perr := resolveExternalProofPolicy(roots.definitionDir)
	if perr != nil {
		return errExit(stderr, perr.Error())
	}
	if proofPolicy == externalProofOn && isTerminalUpdate() {
		flags := classifyEntityFile(entityPath)
		if len(flags) > 0 {
			if force {
				fmt.Fprintf(stderr, "Warning: --force overriding require-external-proof on entity %s\n", slug)
			} else {
				return errExit(stderr, fmt.Sprintf(
					"entity %s cannot advance to terminal — AC(s) [%s] have self-referential proof "+
						"(no test, command, file, or on-disk-state cited). Add an external-proof clause to each, or use --force to bypass.",
					slug, flaggedACLabels(flags)))
			}
		}
	}

	// live-run guard: a runtime-observable AC (one whose truth can only be
	// decided by RUNNING the real producer) declares itself with the explicit
	// `Verified by: live <ref>` convention. Under the SAME require-external-proof
	// opt-in, a terminal --set is refused unless every such AC cites a RESOLVABLE
	// live run. Resolution is three-way: a ci-run:/session: that resolves passes;
	// a placeholder / definitive 404 / absent .jsonl refuses; a connectivity-or-
	// auth error reaching GitHub is INDETERMINATE — surfaced as a tooling error,
	// NOT a refusal, so a network blip never masquerades as a missing live run
	// (the yy slip). An offline-checkable AC (any non-`live` clause) is untouched.
	// --force bypasses with a loud, risk-naming warning. Layered after the self-
	// referential check so it rides the same opt-in and the same terminal guard.
	if proofPolicy == externalProofOn && isTerminalUpdate() {
		data, rdErr := os.ReadFile(entityPath)
		if rdErr == nil {
			for _, lf := range classifyLiveACs(stripFrontmatter(data)) {
				res := resolveLiveCitation(lf.Citation, nil)
				switch res.kind {
				case indeterminate:
					return errExit(stderr, fmt.Sprintf(
						"entity %s: live-run citation for %s could not be checked: %v "+
							"(tooling error, not a refusal — retry when connectivity is restored, or use --force).",
						slug, acLabel(lf.Header), res.err))
				case definitivelyAbsent:
					if force {
						fmt.Fprintf(stderr,
							"Warning: --force overriding live-run requirement on entity %s — runtime-observable AC %s "+
								"is being terminalized without a cited live run (this is exactly the yy slip)\n",
							slug, acLabel(lf.Header))
					} else {
						return errExit(stderr, fmt.Sprintf(
							"entity %s cannot advance to terminal — runtime-observable AC %s declares `Verified by: live …` "+
								"but cites no resolvable live run (citation %q). Cite a resolvable ci-run:<id> or session:<path>, "+
								"OR if the run exists check your token's repo scope (a private/unscoped-token repo returns a masked 404), "+
								"or use --force to bypass.",
							slug, acLabel(lf.Header), lf.Citation))
					}
				}
			}
		}
	}

	resolvedFields, err := updateFrontmatter(entityPath, set.updates)
	if err != nil {
		return errExit(stderr, err.Error())
	}

	switch {
	case asJSON:
		changes := make(jsonArr, 0, len(resolvedFields.keys()))
		for _, field := range resolvedFields.keys() {
			val, _ := resolvedFields.get(field)
			changes = append(changes, newJSONObj().
				set("field", field).set("old", currentFields[field]).set("new", val))
		}
		emitJSON(stdout, newJSONObj().set("command", "set").set("slug", slug).setValue("changes", changes))
	case quiet:
		var pairs []string
		for _, field := range resolvedFields.keys() {
			val, _ := resolvedFields.get(field)
			pairs = append(pairs, fmt.Sprintf("%s=%s->%s", field, currentFields[field], val))
		}
		fmt.Fprintf(stdout, "set slug=%s %s\n", slug, strings.Join(pairs, " "))
	default:
		for _, field := range resolvedFields.keys() {
			val, _ := resolvedFields.get(field)
			oldValue := currentFields[field]
			fmt.Fprintf(stdout, "%s: %s -> %s\n", field, oldValue, val)
		}
	}
	return 0
}

// runRead handles the table / --next / --boot / --validate read flows. Matches
// the tail of main() after the mutation branches.
func runRead(probe claudeteam.TeamStateProbe, roots roots, args []string, e env, whereFilters []whereFilter,
	includeArchive, showNext, showBoot, showNextID, showValidate bool,
	explicitFields []string, allFieldsFlag, asJSON, quiet bool,
	hasArchiveSlug, hasSet, hasResolve bool,
	stdout, stderr io.Writer) int {

	readme := filepath.Join(roots.definitionDir, "README.md")
	var stages []Stage
	if fileExists(readme) {
		stages = parseStagesBlock(readme)
	}

	if showNext && stages == nil {
		return errExit(stderr, "README.md has no stages block. --next requires stage metadata.")
	}
	if showBoot && stages == nil {
		return errExit(stderr, "README.md has no stages block. --boot requires stage metadata.")
	}

	gitRoot := FindGitRoot(roots.definitionDir)
	idStyle, err := workflowIDStyle(roots.definitionDir)
	if err != nil {
		return errExit(stderr, err.Error())
	}

	if showValidate {
		var incompatible []string
		if showNext {
			incompatible = append(incompatible, "--next")
		}
		if includeArchive {
			incompatible = append(incompatible, "--archived")
		}
		if showBoot {
			incompatible = append(incompatible, "--boot")
		}
		if showNextID {
			incompatible = append(incompatible, "--next-id")
		}
		if len(whereFilters) > 0 {
			incompatible = append(incompatible, "--where")
		}
		if explicitFields != nil || allFieldsFlag {
			incompatible = append(incompatible, "--fields/--all-fields")
		}
		if hasArchiveSlug {
			incompatible = append(incompatible, "--archive")
		}
		if hasSet {
			incompatible = append(incompatible, "--set")
		}
		if len(incompatible) > 0 {
			return errExit(stderr, "--validate cannot be combined with "+strings.Join(incompatible, ", "))
		}
		// Explicit --validate command opts INTO the external-proof sub-check.
		// The read-path pre-check (failOnValidationErrors) passes false.
		errs := validateWorkflow(roots.definitionDir, roots.entityDir, idStyle, true, stderr)
		if len(errs) > 0 {
			for _, er := range errs {
				fmt.Fprintln(stderr, er)
			}
			if asJSON {
				emitJSON(stdout, singletonJSON("validate", "valid", "false"))
			}
			return 1
		}
		if asJSON {
			emitJSON(stdout, singletonJSON("validate", "valid", "true"))
			return 0
		}
		fmt.Fprintln(stdout, "VALID")
		return 0
	}

	if rc := failOnValidationErrors(roots, idStyle, stderr); rc != 0 {
		return rc
	}

	allEntities := activeAndArchivedEntities(roots.entityDir, stderr)
	entities := scanEntitiesActive(roots.entityDir, stderr)
	if includeArchive {
		entities = append(entities, archiveEntities(roots.entityDir, stderr)...)
	}

	applyEffectiveIDs(allEntities, idStyle, allEntities)
	applyEffectiveIDs(entities, idStyle, allEntities)
	materializeSuppressedBy(entities, stages, explicitFields, whereFilters)
	entities = applyFilters(entities, whereFilters)

	switch {
	case showBoot:
		if asJSON {
			data, err := gatherBoot(probe, entities, stages, roots.definitionDir, roots.entityDir, gitRoot, idStyle, e, stderr)
			if err != nil {
				return 1
			}
			emitJSON(stdout, bootJSON(data))
			return 0
		}
		if err := printBoot(probe, stdout, entities, stages, roots.definitionDir, roots.entityDir, gitRoot, idStyle, e, stderr); err != nil {
			return 1
		}
	case showNext:
		if asJSON {
			emitJSON(stdout, nextJSON(entities, stages, explicitFields, allFieldsFlag))
			return 0
		}
		extras := resolveExtraFields(entities, explicitFields, allFieldsFlag, defaultNextFields, nextFixedFields)
		printNextTable(stdout, entities, stages, extras, quiet)
	default:
		if asJSON {
			fields := resolveJSONFields(entities, explicitFields, allFieldsFlag, defaultStatusFields)
			emitJSON(stdout, statusJSON(entities, stages, fields))
			return 0
		}
		extras := resolveExtraFields(entities, explicitFields, allFieldsFlag, defaultStatusFields, defaultStatusFields)
		printStatusTable(stdout, entities, stages, extras, quiet)
	}
	return 0
}

// discoverIgnoreDirs is the baseline prune set for --discover. Matches
// DISCOVER_IGNORE_DIRS, extended with .spacedock-state so a state checkout is
// never walked into and re-surfaced as a second workflow (the native split-root
// state dir holds no own README; a stray symlinked one would otherwise match).
var discoverIgnoreDirs = map[string]bool{
	".git": true, ".worktrees": true, "node_modules": true, "vendor": true,
	"dist": true, "build": true, "__pycache__": true, "tests": true,
	".spacedock-state": true,
}

// readGitignoreDirBasenames returns basenames of directory-pattern entries in
// {root}/.gitignore. Matches _read_gitignore_dir_basenames.
func readGitignoreDirBasenames(root string) map[string]bool {
	names := map[string]bool{}
	data, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		return names
	}
	for _, line := range strings.Split(string(data), "\n") {
		entry := strings.TrimSpace(line)
		if entry == "" || strings.HasPrefix(entry, "#") || strings.HasPrefix(entry, "!") {
			continue
		}
		if !strings.HasSuffix(entry, "/") {
			continue
		}
		base := filepath.Base(strings.TrimRight(entry, "/"))
		if base != "" {
			names[base] = true
		}
	}
	return names
}

// discoverWorkflows walks root and returns workflow dirs (README with a
// commissioned-by: spacedock@ frontmatter), realpath'd and sorted. Matches
// discover_workflows.
func discoverWorkflows(root string) []string {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		absRoot = root
	}
	ignore := map[string]bool{}
	for k := range discoverIgnoreDirs {
		ignore[k] = true
	}
	for k := range readGitignoreDirBasenames(absRoot) {
		ignore[k] = true
	}

	seenReal := map[string]bool{}
	var results []string
	resultSet := map[string]bool{}
	// prunedStateDirs holds the realpath of each discovered workflow's state
	// checkout (its `state:` target). A state checkout is the mutable-entity
	// source of a workflow already counted via its main README; it must never be
	// reported as a separate workflow even if a stray README symlink lingers
	// inside it. This complements the .spacedock-state basename prune for
	// workflows that name their state dir something else.
	prunedStateDirs := map[string]bool{}

	var walk func(dir string)
	walk = func(dir string) {
		realDir := realpathOf(dir)
		if seenReal[realDir] {
			return
		}
		seenReal[realDir] = true
		if prunedStateDirs[realDir] {
			return
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		readmePath := filepath.Join(dir, "README.md")
		if isRegularFile(readmePath) {
			fields := ParseFrontmatter(readmePath)
			if strings.HasPrefix(fields["commissioned-by"], "spacedock@") {
				resolved := realpathOf(dir)
				if !resultSet[resolved] {
					resultSet[resolved] = true
					results = append(results, resolved)
				}
				if mode, relPath, err := ClassifyState(fields["state"]); err == nil && mode == StateSplitRoot {
					prunedStateDirs[realpathOf(filepath.Join(dir, relPath))] = true
				}
			}
		}
		for _, ent := range entries {
			// os.Stat follows symlinks, matching os.walk(followlinks=True).
			st, err := os.Stat(filepath.Join(dir, ent.Name()))
			if err != nil || !st.IsDir() {
				continue
			}
			if ignore[ent.Name()] {
				continue
			}
			if prunedStateDirs[realpathOf(filepath.Join(dir, ent.Name()))] {
				continue
			}
			walk(filepath.Join(dir, ent.Name()))
		}
	}
	walk(absRoot)

	sort.Strings(results)
	return results
}

// runGitCmd runs git in dir and returns stdout, or an error on failure.
func runGitCmd(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
