// ABOUTME: Pure fail-closed gate-application eligibility and atomic consumption.
// ABOUTME: Consumption spends authorization with status in one entity replacement.
package gates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// EvaluateEligibility is a pure read over a validated gate record, the entity's
// current status, and the caller's digest comparison. It never queries another
// entity or attempts to perform workflow effects.
func EvaluateEligibility(doc *Document, status string, reviewedInputCurrent bool) Eligibility {
	result := Eligibility{Condition: "ineligible"}
	if doc == nil {
		return result
	}
	record, err := recordForStage(doc, status)
	if err != nil || len(record.Attempts) == 0 {
		return result
	}
	attempt := &record.Attempts[len(record.Attempts)-1]
	result.Gate, result.Attempt = record.ID, attempt.ID
	app := attempt.Application
	if app == nil {
		if attempt.Resolution != nil {
			switch attempt.Resolution.Decision {
			case "hold":
				result.Condition = "not-applicable"
			case "revise":
				result.Condition = "feedback-pending"
			}
		}
		return result
	}
	// The enclosing approval Resolution fixes the operation. Keep the existing
	// CLI/readiness vocabulary by deriving its action rather than storing it.
	result.Action, result.TargetStage = "advance", app.TargetStage
	result.ApplicationState = app.State
	switch app.State {
	case "consumed":
		result.Condition = "consumed"
		return result
	case "superseded":
		result.Condition = "superseded"
		return result
	case "pending":
	default:
		return result
	}
	if !reviewedInputCurrent {
		result.Condition = "stale"
		return result
	}
	if status != record.Stage || attemptState(attempt) != "closed" ||
		attempt.Resolution.Briefing != attempt.Briefing.ID ||
		attempt.Resolution.Decision != "approve" ||
		strings.TrimSpace(app.TargetStage) == "" || app.TargetStage == record.Stage {
		return result
	}
	result.Condition = "approved-pending"
	result.Eligible = true
	return result
}

func EligibilityFile(path string) (Eligibility, error) {
	return eligibilityFileAt(path, nearestWorkflowDir(filepath.Dir(path)))
}

func EligibilityFileAt(path, workflowDir string) (Eligibility, error) {
	return eligibilityFileAt(path, workflowDir)
}

func eligibilityFileAt(path, workflowDir string) (Eligibility, error) {
	doc, _, err := Read(path)
	if err != nil {
		return Eligibility{}, err
	}
	if err := validateRetainedAuthority(path, workflowDir, doc); err != nil {
		return Eligibility{}, err
	}
	status, err := entityStatus(path)
	if err != nil {
		return Eligibility{}, err
	}
	inputState := reviewedInputUnknown
	if record, err := recordForStage(doc, status); err == nil && len(record.Attempts) > 0 {
		inputState = inspectReviewedInput(path, record.Attempts[len(record.Attempts)-1].Briefing)
	}
	result := EvaluateEligibility(doc, status, inputState == reviewedInputCurrent)
	if inputState == reviewedInputUnknown && result.Condition == "stale" {
		result.Condition = "ineligible"
	}
	if result.Eligible && !applicationTargetMatches(workflowDir, status, result.TargetStage) {
		result.Eligible = false
		result.Condition = "ineligible"
	}
	return result, nil
}

// Consume spends an eligible approval once. A stale pending approval is marked
// superseded without changing status. Other ineligible states are read-only.
func Consume(path string) (ConsumeResult, error) {
	return ConsumeAt(path, nearestWorkflowDir(filepath.Dir(path)))
}

func ConsumeAt(path, workflowDir string) (ConsumeResult, error) {
	unlock, err := lockEntity(path)
	if err != nil {
		return ConsumeResult{}, err
	}
	defer unlock()
	doc, oldNode, err := Read(path)
	if err != nil {
		return ConsumeResult{}, err
	}
	if err := validateRetainedAuthority(path, workflowDir, doc); err != nil {
		return ConsumeResult{}, err
	}
	status, err := entityStatus(path)
	if err != nil {
		return ConsumeResult{}, err
	}
	record, recordErr := recordForStage(doc, status)
	inputState := reviewedInputUnknown
	if recordErr == nil && len(record.Attempts) > 0 {
		inputState = inspectReviewedInput(path, record.Attempts[len(record.Attempts)-1].Briefing)
	}
	eligibility := EvaluateEligibility(doc, status, inputState == reviewedInputCurrent)
	if inputState == reviewedInputUnknown && eligibility.Condition == "stale" {
		eligibility.Condition = "ineligible"
	}
	result := ConsumeResult{Eligibility: eligibility}
	if record == nil || len(record.Attempts) == 0 {
		return result, nil
	}
	attempt := &record.Attempts[len(record.Attempts)-1]
	if eligibility.Condition == "stale" && attempt.Application != nil && attempt.Application.State == "pending" {
		attempt.Application.State = "superseded"
		if err := validateApplicationMutation(oldNode, doc, attempt.ID, "pending", "superseded"); err != nil {
			return ConsumeResult{}, err
		}
		if err := writeDocument(path, oldNode, doc); err != nil {
			return ConsumeResult{}, err
		}
		result.ApplicationState = "superseded"
		result.Wrote = true
		return result, nil
	}
	if !eligibility.Eligible {
		return result, nil
	}
	if !applicationTargetMatches(workflowDir, status, eligibility.TargetStage) {
		result.Eligible = false
		result.Condition = "ineligible"
		return result, nil
	}
	// A terminal-target approval is NOT spent here: consume leaves the
	// application pending and the status untouched; readiness then projects the
	// approved-awaiting-merge route (reported via CurrentStageReadiness —
	// ConsumeResult carries no duplicate route field). The terminal merge
	// ceremony (merge guard) is the sole terminal consumer; the ceremony
	// trigger hangs off this consume-produced pending approval, not off the
	// stage. Consume carries no merge-hook knowledge: no hook discovery, no
	// arming, no scanner lookup — merge guard discovers and arms the delivery
	// mechanism when it acts. Re-consuming the same still-pending terminal
	// application leaves it pending again: routing is an at-least-once effect
	// without authority movement.
	if advanceTargetTerminal(workflowDir, status, eligibility.TargetStage) {
		return result, nil
	}
	attempt.Application.State = "consumed"
	if err := validateApplicationMutation(oldNode, doc, attempt.ID, "pending", "consumed"); err != nil {
		return ConsumeResult{}, err
	}
	if err := writeDocumentAndStatus(path, oldNode, status, doc, eligibility.TargetStage); err != nil {
		return ConsumeResult{}, err
	}
	result.Consumed = true
	result.ApplicationState = "consumed"
	result.Wrote = true
	return result, nil
}

func nearestWorkflowDir(start string) string {
	for dir := filepath.Clean(start); ; dir = filepath.Dir(dir) {
		if _, err := applicationStages(filepath.Join(dir, "README.md")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
	}
}

func applicationTargetMatches(workflowDir, current, target string) bool {
	if workflowDir == "" {
		return false
	}
	stages, err := applicationStages(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return false
	}
	i := applicationStageIndex(stages, current)
	return i >= 0 && i+1 < len(stages) && stages[i+1].Name == target
}

// advanceTargetTerminal reports whether current's declared advance target (its
// immediate README successor) is the terminal stage. Callers establish the
// successor match first (applicationTargetMatches); this only reads the
// successor's terminal flag. An unparseable workflow reports false, matching
// applicationTargetMatches' fail-closed shape.
func advanceTargetTerminal(workflowDir, current, target string) bool {
	if workflowDir == "" {
		return false
	}
	stages, err := applicationStages(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return false
	}
	i := applicationStageIndex(stages, current)
	return i >= 0 && i+1 < len(stages) && stages[i+1].Name == target && stages[i+1].Terminal
}

type reviewedInputCheck int

const (
	reviewedInputUnknown reviewedInputCheck = iota
	reviewedInputCurrent
	reviewedInputStale
)

func inspectReviewedInput(entityPath string, binding Briefing) reviewedInputCheck {
	var digest string
	var err error
	var data []byte
	if binding.RequestDigest == "" {
		// boundBriefingPath resolves a prepared room's locator and a legacy ref
		// to the Briefing file itself. It does not compare the digest, so a
		// drifted Briefing still reports stale here rather than unknown.
		var path string
		if path, err = boundBriefingPath(entityPath, binding); err == nil {
			data, err = os.ReadFile(path)
		}
	} else {
		data, _, err = boundBriefingBytes(entityPath, binding)
	}
	if err != nil {
		return reviewedInputUnknown
	}
	digest, err = CanonicalDigest(data)
	if err == nil && digest == binding.Digest {
		return reviewedInputCurrent
	}
	return reviewedInputStale
}

// validateApplicationMutation proves the selected application's state is the
// only gates-record change. This is the narrow exception to closed-attempt
// freezing used by consumption, staleness marking, and the terminal merge
// ceremony's locked delivery writes in delivery.go (pending->consumed with
// delivery proof, pending->superseded on the --rework route).
func validateApplicationMutation(oldNode *yaml.Node, next *Document, attemptID, from, to string) error {
	var old Document
	if err := oldNode.Decode(&old); err != nil {
		return err
	}
	found := false
	for ri := range old.Records {
		for ai := range old.Records[ri].Attempts {
			a := &old.Records[ri].Attempts[ai]
			if a.ID != attemptID {
				continue
			}
			if a.Application == nil || a.Application.State != from {
				return fmt.Errorf("attempt %s application is not %s", attemptID, from)
			}
			a.Application.State = to
			found = true
		}
	}
	if !found {
		return fmt.Errorf("application attempt %s does not exist", attemptID)
	}
	if err := Validate(next); err != nil {
		return err
	}
	left, _ := yaml.Marshal(&old)
	right, _ := yaml.Marshal(next)
	if !bytes.Equal(left, right) {
		return fmt.Errorf("application mutation changed fields outside attempt %s state", attemptID)
	}
	return nil
}

func writeDocumentAndStatus(path string, expected *yaml.Node, expectedStatus string, doc *Document, status string) error {
	return writeEntityDocument(path, expected, &expectedStatus, doc, &status)
}
