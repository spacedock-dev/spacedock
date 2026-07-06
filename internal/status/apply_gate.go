// ABOUTME: status --apply-gate applies a human gate verdict to workflow state
// ABOUTME: through the existing guarded --set mutation path.
package status

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

func runApplyGate(roots roots, req *applyGateRequest, quiet, asJSON bool, stdout, stderr io.Writer) int {
	readme := filepath.Join(roots.definitionDir, "README.md")
	stages := parseStagesBlock(readme)
	if len(stages) == 0 {
		return errExit(stderr, "README.md has no stages block. --apply-gate requires stage metadata.")
	}

	resolved, rc := resolveMutationEntity(roots, req.entity, stderr)
	if rc != 0 {
		return rc
	}
	entityPath := resolveEntityPath(roots.entityDir, resolved.slug, stderr)
	if entityPath == "" {
		return errExit(stderr, "entity not found: "+resolved.slug)
	}

	fields := ParseFrontmatter(entityPath)
	currentStatus := strings.TrimSpace(fields["status"])
	currentStage, idx, ok := stageByName(stages, currentStatus)
	if !ok {
		return errExit(stderr, fmt.Sprintf(
			"entity %s status %q is not a defined stage in workflow %s",
			resolved.slug, currentStatus, roots.definitionDir))
	}
	if !currentStage.gate {
		return errExit(stderr, fmt.Sprintf(
			"entity %s is not at a gate stage (status=%s)",
			resolved.slug, currentStatus))
	}

	nextStatus, err := applyGateTargetStatus(stages, idx, req.verdict)
	if err != nil {
		return errExit(stderr, fmt.Sprintf("entity %s: %v", resolved.slug, err))
	}

	updates := []fieldUpdate{
		{field: "status", value: nextStatus, hasValue: true},
		{field: "gate-id", value: req.gateID, hasValue: true},
		{field: "gate-verdict", value: req.verdict, hasValue: true},
	}
	set := &setUpdate{slug: resolved.slug, updates: updates}
	rc = runSet(roots, set, nil, nil,
		false, false, false, false, false, false, true, false,
		io.Discard, stderr)
	if rc != 0 {
		return rc
	}

	switch {
	case asJSON:
		emitJSON(stdout, newJSONObj().
			set("command", "apply-gate").
			set("slug", resolved.slug).
			set("gate_id", req.gateID).
			set("verdict", req.verdict).
			set("status_old", currentStatus).
			set("status_new", nextStatus))
	case quiet:
		fmt.Fprintf(stdout, "apply-gate slug=%s gate=%s verdict=%s status=%s->%s\n",
			resolved.slug, req.gateID, req.verdict, currentStatus, nextStatus)
	default:
		fmt.Fprintf(stdout, "apply-gate slug=%s gate=%s verdict=%s status=%s->%s\n",
			resolved.slug, req.gateID, req.verdict, currentStatus, nextStatus)
	}
	return 0
}

func stageByName(stages []Stage, name string) (Stage, int, bool) {
	for i, s := range stages {
		if s.Name == name {
			return s, i, true
		}
	}
	return Stage{}, -1, false
}

func applyGateTargetStatus(stages []Stage, currentIndex int, verdict string) (string, error) {
	current := stages[currentIndex]
	switch verdict {
	case "approve":
		if currentIndex+1 >= len(stages) {
			return "", fmt.Errorf("cannot approve gate %s because it has no next stage", current.Name)
		}
		return stages[currentIndex+1].Name, nil
	case "revise", "reject":
		if target, ok := current.optional["feedback-to"]; ok {
			target = strings.TrimSpace(target)
			if target == "" {
				return current.Name, nil
			}
			if _, _, exists := stageByName(stages, target); !exists {
				return "", fmt.Errorf("feedback target %q is not a defined stage", target)
			}
			return target, nil
		}
		return current.Name, nil
	default:
		return "", fmt.Errorf("unsupported verdict %q", verdict)
	}
}
