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
	currentStatus := strings.TrimSpace(currentFields["status"])
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

	// Post-dispatch guard: once the normal same-stage dispatch owns a worktree,
	// that entered working stage cannot be changed away until the shared
	// scheduler predicate sees a committed, structurally complete current-stage
	// report. This is intentionally before every --force bypass guard below.
	// Any away status token in a repeated/chained update refuses the whole
	// command byte-clean; the same-stage dispatch mutation itself and unrelated
	// non-status updates remain allowed.
	if strings.TrimSpace(currentFields["worktree"]) != "" {
		for _, stage := range stages {
			if stage.Name != currentStatus || !enteredStageAwaitingCompletion(&entity{path: entityPath}, stage) {
				continue
			}
			for _, u := range set.updates {
				if u.field == "status" && u.hasValue && u.value != currentStatus {
					return errExit(stderr, fmt.Sprintf(
						"entity %s cannot change status away from entered stage %q until a durable, complete "+
							"## Stage Report: %s is committed.",
						slug, currentStatus, currentStatus))
				}
			}
			break
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
				"entity %s has %s. Clear mod-block in a separate --set call. "+
					"(--force bypasses this guard; a refusal usually means a ceremony step was skipped — re-run merge guard %s instead.)",
				slug, reason, slug))
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

	// Finalize-status gate: the finalize action (setting `completed`) requires the
	// resulting status to be a declared terminal stage — either already terminal or
	// set terminal in THIS SAME --set call. This closes the residual hole behind the
	// verdict gate above: `--set completed verdict=X worktree=` alone (no
	// status={terminal}) satisfies the verdict gate yet still finalizes an entity
	// that never advanced past a non-terminal stage — Spike C's exact reproduction
	// of the incomplete-finalize deviation the merge guard's atomic terminalize
	// (status+verdict+completed in one call) exists to prevent. No-op when the
	// workflow declares no stages block (membership cannot be checked). --force
	// bypasses, same idiom as the guards above.
	postUpdateStatus := strings.TrimSpace(currentFields["status"])
	for _, u := range set.updates {
		if u.field == "status" && u.hasValue {
			postUpdateStatus = u.value
		}
	}
	if !force && len(terminalNames) > 0 {
		// Sole-terminal-consumer refusal: while the entity carries a binding
		// pending terminal-target application, merge guard's locked gates write
		// is the ONLY writer allowed to land terminal status — a hand --set to
		// terminal would recreate the done-but-undelivered shape (authority
		// unspent, delivery unproven). --force honors the uniform escape hatch.
		// Refusal-only: no state or record is added.
		terminalSet := false
		for _, u := range set.updates {
			if u.field == "status" && u.hasValue && terminalNames[u.value] {
				terminalSet = true
			}
		}
		if terminalSet {
			// Fail-closed: only a genuinely gate-less entity (ErrNoGateRecord,
			// flattened to pending=false, err=nil) takes the hand-set path.
			// Unreadable/stale authority, a live pending approval, or any
			// other gates-record shape refuses — a classification we cannot
			// make must never default to permission.
			_, pending, classErr := pendingTerminalApproval(entityPath, roots.definitionDir, strings.TrimSpace(currentFields["status"]))
			if classErr != nil {
				return errExit(stderr, fmt.Sprintf(
					"entity %s: terminal status --set refused — its gate authority cannot be classified (%v); "+
						"merge guard %s is the sole terminal consumer while that authority is unreadable, stale, or in force. "+
						"(--force bypasses this guard.)",
					slug, classErr, slug))
			}
			if pending {
				return errExit(stderr, fmt.Sprintf(
					"entity %s carries a pending terminal-target approval — merge guard %s is the sole terminal consumer: "+
						"run `spacedock merge guard %s --verdict passed|rejected` (or `merge guard %s --rework` to send back). "+
						"(--force bypasses this guard.)",
					slug, slug, slug, slug))
			}
		}
	}
	if !force && finalizing && len(stages) > 0 && !terminalNames[postUpdateStatus] {
		return errExit(stderr, fmt.Sprintf(
			"entity %s cannot be finalized ('completed') while status '%s' is not the terminal stage. "+
				"Set status=%s in the same --set (or run 'spacedock merge guard %s'), or use --force.",
			slug, postUpdateStatus, terminalStageName(roots.definitionDir), slug))
	}

	// The rejected carve-out is case-insensitive (isRejectedVerdict): the verdict
	// under test comes from frontmatter or from a --set value, so it may arrive as
	// the `REJECTED` merge guard writes or the `rejected` an older binary wrote.
	if !force && policy != mergeLocal && isTerminalUpdate() && modBlock == "" && postUpdatePR == "" && !isRejectedVerdict(postUpdateVerdict) {
		mergeHooks := scanMods(roots.definitionDir)["merge"]
		if len(mergeHooks) > 0 {
			return errExit(stderr, fmt.Sprintf(
				"entity %s cannot advance to terminal — workflow has merge hook(s) [%s] that have not run "+
					"(pr field is empty and mod-block is empty). Set mod-block=merge:%s and invoke the hook. "+
					"(--force bypasses this guard; a refusal usually means a ceremony step was skipped — re-run merge guard %s instead.)",
				slug, strings.Join(mergeHooks, ", "), mergeHooks[0], slug))
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

	// Conventional-vocabulary admission: a field declaring a schema `conventional`
	// list is closed on write (conventionalViolation), so `--set verdict=superseded`
	// — the exact write that produced the four archived records carrying a token the
	// enum never admitted — refuses instead of storing it. The write path already
	// consulted this list to fold case (canonicalConventional, inside
	// updateFrontmatter); this is that same lookup used for admission, not a second
	// source of truth. `verdict` is the only conventional field today, so the blast
	// radius is exactly it, but nothing here is verdict-specific.
	//
	// The guard lives here rather than in updateFrontmatter because `force` is in
	// scope here alongside the other --set guards, and because the shared write
	// engine stays a pure normaliser for finalize and archive, which pass
	// already-canonical values. It is placed last, immediately before the write, so
	// every already-refused command keeps reporting its existing cause first and the
	// refusal is byte-clean by construction. --force bypasses, the uniform escape
	// the guards above honor.
	if !force {
		for _, u := range set.updates {
			problem := conventionalViolation(u.field, u.value)
			if problem == "" {
				continue
			}
			return errExit(stderr, fmt.Sprintf(
				"entity %s: %s. Pick one of those, clear it with `%s=`, or use --force.",
				slug, problem, u.field))
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
// the tail of main() after the mutation branches. transcriptProbe is the boot
// guard's receipt-write seam (nil on a non-Claude host).
func runRead(probe claudeteam.TeamStateProbe, transcriptProbe claudeteam.TranscriptProbe, roots roots, args []string, e env, whereFilters []whereFilter,
	includeArchive, showNext, showBoot, showNextID, showValidate, identify bool,
	explicitFields []string, allFieldsFlag, asJSON, quiet bool,
	hasArchiveSlug, hasSet, hasResolve bool,
	page, limit int,
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
		// Explicit --validate command opts INTO the external-proof sub-check and
		// the warn-tier per-field schema-conformance sub-check. The read-path
		// pre-check (failOnValidationErrors) passes false, opting out of both.
		errs, warns := validateWorkflow(roots.definitionDir, roots.entityDir, idStyle, true, stderr)
		// Warn-tier per-field schema conformance prints to stderr but does NOT
		// flip the exit code — exit 1 stays reserved for structural errors.
		for _, w := range warns {
			fmt.Fprintln(stderr, w)
		}
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
	materializeGateReadinessWhenReferenced(entities, roots.definitionDir, explicitFields, allFieldsFlag, whereFilters)
	// Machine scheduler reads consume the same readiness reducer as boot
	// identify. Materialize it before filtering so --next can expose ready gates
	// even when no gate-readiness field was explicitly requested.
	if showNext {
		materializeGateReadiness(entities, stages)
	}
	materializeSuppressedBy(entities, stages, explicitFields, whereFilters)
	if err := validateWhereFields(entities, whereFilters); err != nil {
		return errExit(stderr, err.Error())
	}
	entities = applyFilters(entities, whereFilters)

	switch {
	case showBoot:
		if asJSON {
			data, err := gatherBoot(probe, transcriptProbe, entities, stages, roots.definitionDir, roots.entityDir, gitRoot, idStyle, e, stderr, identify)
			if err != nil {
				return 1
			}
			emitJSON(stdout, bootJSON(data))
			return 0
		}
		if err := printBoot(probe, transcriptProbe, stdout, entities, stages, roots.definitionDir, roots.entityDir, gitRoot, idStyle, e, stderr, identify); err != nil {
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
		win := paginate(len(entities), page, limit)
		if asJSON {
			fields := resolveJSONFields(entities, explicitFields, allFieldsFlag, defaultStatusFields)
			emitJSON(stdout, statusJSON(entities, stages, fields, win))
			return 0
		}
		extras := resolveExtraFields(entities, explicitFields, allFieldsFlag, defaultStatusFields, defaultStatusFields)
		printStatusTable(stdout, entities, stages, extras, quiet, win)
	}
	return 0
}

// discoverIgnoreDirs is the baseline prune set for --discover. Matches
// DISCOVER_IGNORE_DIRS, extended with .spacedock-state so a state checkout is
// never walked into and re-surfaced as a second workflow (the native split-root
// state dir holds no own README; a stray symlinked one would otherwise match),
// and with testdata so a commissioned-shape README carried as a package-adjacent
// test fixture (Go ignores testdata/ when building, this repo homes fixtures
// there) is not counted as a real workflow — the same test-scaffolding tradeoff
// the set already makes for tests/vendor/dist/build.
var discoverIgnoreDirs = map[string]bool{
	".git": true, ".worktrees": true, "node_modules": true, "vendor": true,
	"dist": true, "build": true, "__pycache__": true, "tests": true,
	".spacedock-state": true, "testdata": true,
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
			child := filepath.Join(dir, ent.Name())
			// os.Stat follows symlinks, matching os.walk(followlinks=True).
			st, err := os.Stat(child)
			if err != nil || !st.IsDir() {
				continue
			}
			if ignore[ent.Name()] {
				continue
			}
			if prunedStateDirs[realpathOf(child)] {
				continue
			}
			// Skip a nested git checkout (a linked or agent worktree): it is a full
			// copy of the repo, so any commissioned workflow inside it is a duplicate
			// of one the outer repo already carries, not a distinct workflow. This is
			// the host-neutral generalization of the .worktrees prune — it catches
			// agent worktrees under any directory (e.g. `.claude/worktrees/<agent>`)
			// without naming a host-specific path. The start root's own `.git` is never
			// inspected here: this check runs only on descended children.
			if hasGitEntry(child) {
				continue
			}
			walk(child)
		}
	}
	walk(absRoot)

	sort.Strings(results)
	return results
}

// hasGitEntry reports whether dir holds a `.git` entry — a directory (a normal
// repository checkout) or a regular file (a linked/agent worktree's gitlink). A
// dir with one is a self-contained checkout, so the discovery walk does not
// descend into it (its workflows are copies of the outer repo's).
//
// Guarded by TestDiscoverWorkflowsSkipsNestedCheckout (discover_worktree_noise_test.go):
// that test is the sole coverage proving this prune fires, so changing this
// gitlink/checkout detection should be checked against it.
func hasGitEntry(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil && (st.IsDir() || st.Mode().IsRegular())
}

// DiscoverWorkflows returns the commissioned workflow dirs under root (realpath'd,
// sorted), with the same linked/agent-worktree + VCS noise pruned as the
// `--discover` boot walk — so a caller outside this package (e.g. the front-door
// launch banner) resolves the same single real workflow the first officer sees.
func DiscoverWorkflows(root string) []string {
	return discoverWorkflows(root)
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
