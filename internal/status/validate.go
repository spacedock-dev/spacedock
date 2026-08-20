// ABOUTME: Validation matching validate_workflow — flat/folder conflicts,
// ABOUTME: stage-name regex, and per-id-style id rules with evidence lines.
package status

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

// entityEvidence formats the Error: ... workflow= scope= slug= id= [display=]
// path= evidence line. Matches entity_evidence.
func entityEvidence(e *entity, workflowDir, problem, displayID string) string {
	return entityEvidenceLine("Error", e, workflowDir, problem, displayID)
}

// entityEvidenceLine formats an evidence line at the given severity prefix
// ("Error" for gating structural defects, "Warning" for warn-tier field
// conformance). The field shape (workflow=/scope=/slug=/id=/[display=]/path=) is
// identical so the FO locates the entity the same way for either severity.
func entityEvidenceLine(severity string, e *entity, workflowDir, problem, displayID string) string {
	display := displayID
	if display == "" {
		display = e.displayID
	}
	parts := []string{
		fmt.Sprintf("%s: %s:", severity, problem),
		"workflow=" + workflowDir,
		"scope=" + scopeOf(e),
		"slug=" + e.slug,
		"id=" + e.storedID,
	}
	if display != "" {
		parts = append(parts, "display="+display)
	}
	parts = append(parts, "path="+e.path)
	return strings.Join(parts, " ")
}

func scopeOf(e *entity) string {
	if e.scope == "" {
		return "active"
	}
	return e.scope
}

// findEntityFormConflicts returns conflict errors for slugs present as both flat
// and folder entities in directory. directory is the absolute path scanned;
// dirSpelling is the as-passed spelling used in the printed flat_path/folder_path
// and workflow= fields. Matches find_entity_form_conflicts.
func findEntityFormConflicts(directory, dirSpelling, scope string) []string {
	info, err := os.Stat(directory)
	if err != nil || !info.IsDir() {
		return nil
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil
	}
	flat := map[string]bool{}
	folder := map[string]bool{}
	for _, ent := range entries {
		name := ent.Name()
		full := filepath.Join(directory, name)
		st, err := os.Stat(full)
		if err != nil {
			continue
		}
		if st.Mode().IsRegular() {
			if strings.HasSuffix(name, ".md") && name != "README.md" && !strings.HasPrefix(name, ".") {
				flat[strings.TrimSuffix(name, ".md")] = true
			}
		} else if st.IsDir() {
			if !reservedSubdirs[name] && !strings.HasPrefix(name, ".") {
				if isRegularFile(filepath.Join(full, "index.md")) {
					folder[name] = true
				}
			}
		}
	}

	var conflicts []string
	for slug := range flat {
		if folder[slug] {
			conflicts = append(conflicts, slug)
		}
	}
	sort.Strings(conflicts)

	workflowField := dirSpelling
	if scope == "archived" {
		workflowField = pyDirname(dirSpelling)
	}
	var errs []string
	for _, slug := range conflicts {
		folderPath := PyJoin(dirSpelling, slug, "index.md")
		flatPath := PyJoin(dirSpelling, slug+".md")
		errs = append(errs, fmt.Sprintf(
			"Error: flat/folder conflict: workflow=%s scope=%s slug=%s flat_path=%s folder_path=%s",
			workflowField, scope, slug, flatPath, folderPath))
	}
	return errs
}

// kebabSuggestion mirrors _kebab_suggestion.
func kebabSuggestion(name string) string {
	lowered := strings.ToLower(name)
	replaced := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(lowered, "-")
	return strings.Trim(replaced, "-")
}

// validateWorkflowStageNames returns errors for README stage names violating the
// dispatch-name regex. Matches validate_workflow_stage_names.
func validateWorkflowStageNames(definitionDir string) []string {
	readme := filepath.Join(definitionDir, "README.md")
	if !fileExists(readme) {
		return nil
	}
	states := parseStagesBlock(readme)
	if len(states) == 0 {
		return nil
	}
	var errs []string
	for _, s := range states {
		name := s.Name
		if name == "" || stageNameRe.MatchString(name) {
			continue
		}
		suggestion := kebabSuggestion(name)
		if suggestion == "" {
			suggestion = "a-kebab-name"
		}
		errs = append(errs, fmt.Sprintf(
			"workflow '%s': stage name '%s' must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$; rename to '%s' or similar",
			definitionDir, name, suggestion))
	}
	return errs
}

// validateWorkflow returns validation error lines for active + archived
// entities. Matches validate_workflow. definitionDir/entityDir are the absolute
// README and entity roots; the workflow= evidence field uses the absolute
// entityDir (the FO always passes an absolute --workflow-dir for enumerate ops,
// so absolute == the oracle's literal here). includeExternalProof scopes the
// require-external-proof sub-check: only the explicit `--validate` command
// passes true, so the read-path validation pre-check (cwd gate on
// status/--next/--boot/--next-id) never fires the AC classifier. A read path
// failing on a flagged AC would lock the FO out of the very listing they need
// to see the broken entity.
func validateWorkflow(definitionDir, entityDir, idStyle string, includeExternalProof bool, stderr io.Writer) (errs []string, warns []string) {
	errs = append(errs, findEntityFormConflicts(entityDir, entityDir, "active")...)
	errs = append(errs, findEntityFormConflicts(filepath.Join(entityDir, "_archive"), PyJoin(entityDir, "_archive"), "archived")...)
	errs = append(errs, validateWorkflowStageNames(definitionDir)...)

	entities := activeAndArchivedEntities(entityDir, stderr)
	var sdDisplay map[string]string
	if idStyle == "sd-b32" {
		sdDisplay = computeSDB32DisplayIDs(entities)
	}

	switch idStyle {
	case "sequential":
		errs = append(errs, validateSequential(entities, entityDir)...)
	case "slug":
		errs = append(errs, validateSlug(entities, entityDir)...)
	case "sd-b32":
		errs = append(errs, validateSDB32(entities, entityDir, sdDisplay)...)
	}

	if !includeExternalProof {
		return errs, nil
	}
	// Apply display IDs before both gate and schema diagnostics. This keeps the
	// evidence line stable for sd-b32 workflows regardless of which warning tier
	// produced it.
	applyEffectiveIDs(entities, idStyle, entities)

	// Explicit validation also validates the binary-owned gate trees. Ordinary
	// status reads keep their historical projection behavior, but the explicit
	// diagnostic surface must report strict gate errors and bounded application
	// extension warnings with the same entity evidence shape as other findings.
	gateErrs, gateWarns := gateValidationDiagnostics(entities, entityDir)
	errs = append(errs, gateErrs...)
	warns = append(warns, gateWarns...)

	// Warn-tier per-field schema conformance shares the same opt-in as the
	// external-proof sub-check: only the explicit --validate command computes it,
	// reusing the entities already enumerated above (with effective ids applied so
	// the evidence line carries the display id). Warns are returned separately so
	// the caller keeps them out of the exit-code decision.
	warns = append(warns, fieldConformanceWarnings(entities, entityDir)...)

	// require-external-proof sub-check: when the workflow opts in, every
	// active entity is classified and each flagged AC is emitted as a standard
	// entityEvidence line. A typo in the README key is surfaced as the same
	// kind of loud rejection findEntityFormConflicts uses for its own defects.
	policy, perr := resolveExternalProofPolicy(definitionDir)
	if perr != nil {
		errs = append(errs, "Error: "+perr.Error())
		return errs, warns
	}
	if policy == externalProofOn {
		for _, e := range entities {
			if e.scope != "active" {
				continue
			}
			for _, f := range classifyEntityFile(e.path) {
				display := ""
				if idStyle == "sd-b32" {
					display = sdDisplay[e.storedID]
				}
				errs = append(errs, entityEvidence(e, entityDir,
					"self-referential AC proof ("+acLabel(f.Header)+")", display))
			}
		}
	}

	return errs, warns
}

func gateValidationDiagnostics(entities []*entity, workflowDir string) (errs, warns []string) {
	for _, e := range entities {
		_, _, compatibility, err := gates.ReadDiagnostics(e.path)
		if err != nil {
			if errors.Is(err, gates.ErrNoGateRecord) {
				continue
			}
			errs = append(errs, entityEvidenceLine("Error", e, workflowDir, "invalid gates: "+err.Error(), e.displayID))
			continue
		}
		// Warn tier only: archived scope is publish-only, so an application-field
		// extension there can never be cleared by a tool-mediated write. The
		// structural error above stays scope-inclusive.
		if e.scope != "active" {
			continue
		}
		for _, warning := range compatibility {
			problem := fmt.Sprintf("unknown gate application field '%s' at %s", warning.Field, warning.Path)
			warns = append(warns, entityEvidenceLine("Warning", e, workflowDir, problem, e.displayID))
		}
		// A retained room that no longer resolves is the #739 end state: the gate
		// commands fail on it mid-ceremony while every read surface still reports
		// the entity as healthy. Reporting it here is what makes a hand
		// conversion verifiable — the operator can confirm the rewrite landed
		// instead of finding out at the next gate.
		for _, ref := range unresolvedRoomRefs(e.path) {
			warns = append(warns, entityEvidenceLine("Warning", e, workflowDir,
				"retained gate room does not resolve: "+ref, e.displayID))
		}
		// A flat entity that already holds prepared rooms is grandfathered by
		// gate prepare: its refs are ./<slug>/review/... and correct while it
		// stays flat. Moving it to folder form without rewriting them in the
		// same commit makes every retained room unreadable, and nothing else
		// reports that — so the warning carries the whole remedy, not just the
		// finding. Warn tier, and only on explicit --validate: an error here
		// would exit 1 on the plain status read path.
		if filepath.Base(e.path) != "index.md" {
			if _, err := os.Stat(filepath.Join(filepath.Dir(e.path), e.slug, "review")); err == nil {
				warns = append(warns, entityEvidenceLine("Warning", e, workflowDir, fmt.Sprintf(
					"flat entity holds gate rooms in %s/review/; to convert it, `git mv %s.md %s/index.md` AND rewrite every `room-ref: ./%s/` to `room-ref: ./` in the same commit, or every retained room becomes unreadable",
					e.slug, e.slug, e.slug, e.slug), e.displayID))
			}
		}
	}
	return errs, warns
}

func validateSequential(entities []*entity, workflowDir string) []string {
	var errs []string
	byID := map[string][]*entity{}
	var order []string
	for _, e := range entities {
		id := e.storedID
		if id == "" {
			errs = append(errs, entityEvidence(e, workflowDir, "missing required id", ""))
			continue
		}
		if !isDigits(id) {
			errs = append(errs, entityEvidence(e, workflowDir, "non-numeric sequential id", ""))
			continue
		}
		if _, ok := byID[id]; !ok {
			order = append(order, id)
		}
		byID[id] = append(byID[id], e)
	}
	for _, id := range order {
		group := byID[id]
		if len(group) < 2 || !anyActive(group) {
			continue
		}
		for _, e := range group {
			errs = append(errs, entityEvidence(e, workflowDir, "duplicate id", ""))
		}
	}
	return errs
}

func validateSlug(entities []*entity, workflowDir string) []string {
	var errs []string
	seen := map[string]*entity{}
	for _, e := range entities {
		effective := e.slug
		e.displayID = effective
		if first, ok := seen[effective]; ok {
			errs = append(errs, entityEvidence(first, workflowDir, "duplicate effective id", effective))
			errs = append(errs, entityEvidence(e, workflowDir, "duplicate effective id", effective))
		} else {
			seen[effective] = e
		}
	}
	return errs
}

func validateSDB32(entities []*entity, workflowDir string, sdDisplay map[string]string) []string {
	var errs []string
	byID := map[string][]*entity{}
	var order []string
	for _, e := range entities {
		id := e.storedID
		display := sdDisplay[id]
		if id == "" {
			errs = append(errs, entityEvidence(e, workflowDir, "missing required id", display))
			continue
		}
		if !(isValidSDB32ID(id) || isDigits(id)) {
			errs = append(errs, entityEvidence(e, workflowDir, "invalid sd-b32 stored id", display))
			continue
		}
		if _, ok := byID[id]; !ok {
			order = append(order, id)
		}
		byID[id] = append(byID[id], e)
	}
	for _, id := range order {
		group := byID[id]
		if len(group) < 2 || !anyActive(group) {
			continue
		}
		display := sdDisplay[id]
		for _, e := range group {
			errs = append(errs, entityEvidence(e, workflowDir, "duplicate sd-b32 stored id", display))
		}
	}
	return errs
}

func anyActive(group []*entity) bool {
	for _, e := range group {
		if e.scope == "active" {
			return true
		}
	}
	return false
}

// unresolvedRoomRefs returns the retained room-refs of one entity that no longer
// resolve on disk, joined the way the gate commands join them.
func unresolvedRoomRefs(entityPath string) []string {
	doc, _, err := gates.Read(entityPath)
	if err != nil || doc == nil {
		return nil
	}
	var missing []string
	for ri := range doc.Records {
		for ai := range doc.Records[ri].Attempts {
			ref := doc.Records[ri].Attempts[ai].Briefing.RoomRef
			if ref == "" {
				continue
			}
			if _, err := os.Stat(filepath.Join(filepath.Dir(entityPath), filepath.FromSlash(ref))); err != nil {
				missing = append(missing, ref)
			}
		}
	}
	return missing
}
