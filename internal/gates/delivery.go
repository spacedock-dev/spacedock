// ABOUTME: The two gates-owned terminal delivery operations: merge guard's
// ABOUTME: sole-consumer finalize (pending→consumed + terminal status + verdict
// ABOUTME: + completed) and the --rework supersede-and-route, each ONE locked
// ABOUTME: compare-before-replace candidate over the io.go write machinery.
package gates

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FinalizeTerminalApproval spends a binding pending terminal-target approval
// with delivery proof in hand: in ONE locked compare-before-replace candidate
// it writes application.state pending→consumed (through the exactly-once
// guarded mutation), status := the terminal target stage, verdict, and the
// completed stamp. completedStamp is caller-supplied so the stamp format stays
// owned by the status package that defines the completed field.
//
// Error contract (byte-clean, nothing written): the entity carries no gates
// record (ErrNoGateRecord — the only case a caller may fall back to a legacy
// gate-free path), the authority is unreadable or digest-stale, or the current
// application is not a binding pending advance onto the terminal stage.
func FinalizeTerminalApproval(path, workflowDir, verdict, completedStamp string) (string, error) {
	unlock, doc, oldNode, status, eligibility, err := lockPendingTerminalApproval(path, workflowDir)
	if err != nil {
		return "", err
	}
	defer unlock()
	record, err := recordForStage(doc, status)
	if err != nil {
		return "", err
	}
	attempt := &record.Attempts[len(record.Attempts)-1]
	attempt.Application.State = "consumed"
	if err := validateApplicationMutation(oldNode, doc, attempt.ID, "pending", "consumed"); err != nil {
		return "", err
	}
	target := eligibility.TargetStage
	if err := writeEntityDocument(path, oldNode, &status, doc, &target,
		entityField{key: "verdict", value: verdict},
		entityField{key: "completed", value: completedStamp},
	); err != nil {
		return "", err
	}
	return target, nil
}

// SupersedeTerminalApproval routes a binding pending terminal-target approval
// back for rework: in ONE locked compare-before-replace candidate it writes
// application.state pending→superseded (the same guarded mutation the drift
// path uses) and status := the record stage's validated declared feedback-to.
// the superseded attempt is frozen history —
// re-entry runs a successor attempt with a fresh approval. verdict/completed
// are never written pre-delivery.
//
// Error contract (byte-clean, nothing written): same authority contract as
// FinalizeTerminalApproval, plus a declared feedback-to that is missing,
// undefined, or terminal refuses closed.
func SupersedeTerminalApproval(path, workflowDir string) (string, error) {
	unlock, doc, oldNode, status, _, err := lockPendingTerminalApproval(path, workflowDir)
	if err != nil {
		return "", err
	}
	defer unlock()
	// Eligibility guarantees current status == the record stage, so the
	// record's declared feedback-to is read off the current status stage.
	feedbackTo, err := declaredReworkTarget(workflowDir, status)
	if err != nil {
		return "", err
	}
	record, err := recordForStage(doc, status)
	if err != nil {
		return "", err
	}
	attempt := &record.Attempts[len(record.Attempts)-1]
	attempt.Application.State = "superseded"
	if err := validateApplicationMutation(oldNode, doc, attempt.ID, "pending", "superseded"); err != nil {
		return "", err
	}
	if err := writeEntityDocument(path, oldNode, &status, doc, &feedbackTo); err != nil {
		return "", err
	}
	return feedbackTo, nil
}

// DeclaredReworkTarget is the read-only form of the --rework route validation:
// merge guard consults it BEFORE clearing delivery state so a malformed
// feedback-to refuses byte-clean, and SupersedeTerminalApproval re-validates
// under the entity lock before writing.
func DeclaredReworkTarget(workflowDir, recordStage string) (string, error) {
	return declaredReworkTarget(workflowDir, recordStage)
}

// declaredReworkTarget validates the --rework route: the record stage's
// declared feedback-to must exist, be a stage the workflow defines, and be
// non-terminal. Unlike the gate-rejection routing (which silently falls back
// to the same stage on a missing feedback-to and never validates the target),
// --rework refuses on every malformed declaration — the send-back must land
// somewhere real.
func declaredReworkTarget(workflowDir, recordStage string) (string, error) {
	stages, err := applicationStages(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return "", err
	}
	i := applicationStageIndex(stages, recordStage)
	if i < 0 {
		return "", fmt.Errorf("record stage %q is not defined in workflow %s", recordStage, workflowDir)
	}
	declared := strings.TrimSpace(stages[i].FeedbackTo)
	if declared == "" {
		return "", fmt.Errorf("record stage %q declares no feedback-to — merge guard --rework refuses to route without a declared send-back target", recordStage)
	}
	j := applicationStageIndex(stages, declared)
	if j < 0 {
		return "", fmt.Errorf("record stage %q declares feedback-to %q, which is not a stage defined in workflow %s", recordStage, declared, workflowDir)
	}
	if stages[j].Terminal {
		return "", fmt.Errorf("record stage %q declares feedback-to %q, which is a terminal stage — merge guard --rework refuses to route a rework send-back into terminal status", recordStage, declared)
	}
	return declared, nil
}

// ApprovedAwaitingMergeRoute projects the entity's durable gate state through
// CurrentStageReadiness and reports whether it reads approved-awaiting-merge —
// the ONE readiness vocabulary the CLI prints for a held (unspent)
// terminal-target approval after consume routes it. Any unreadable or
// unrouted shape reports false (fail-closed: no route, non-zero consume).
func ApprovedAwaitingMergeRoute(path, workflowDir string) bool {
	doc, _, err := Read(path)
	if err != nil {
		return false
	}
	status, err := entityStatus(path)
	if err != nil {
		return false
	}
	stages, err := applicationStages(filepath.Join(workflowDir, "README.md"))
	if err != nil {
		return false
	}
	taxonomy := make([]ReadinessStage, 0, len(stages))
	for _, s := range stages {
		taxonomy = append(taxonomy, ReadinessStage{Name: s.Name, Gate: s.Gate, Terminal: s.Terminal})
	}
	return CurrentStageReadiness(doc, status, taxonomy) == RouteApprovedAwaitingMerge
}

// lockPendingTerminalApproval takes the entity lock and establishes the
// binding pending terminal-target approval both delivery operations spend or
// supersede: consume's own eligibility predicate (binding approval,
// digest-current briefing, pending state, status at the gated stage) plus the
// terminal-target predicate. Fail-closed: any shortfall — no gates record
// (ErrNoGateRecord), unreadable/digest-stale retained authority, or an
// application that is consumed/superseded/never-pending or not aimed at the
// terminal stage — returns an error with the lock released and no byte
// written. On success the caller owns the returned unlock.
func lockPendingTerminalApproval(path, workflowDir string) (unlock func(), doc *Document, oldNode *yaml.Node, status string, eligibility Eligibility, err error) {
	unlock, err = lockEntity(path)
	if err != nil {
		return nil, nil, nil, "", Eligibility{}, err
	}
	fail := func(err error) (func(), *Document, *yaml.Node, string, Eligibility, error) {
		unlock()
		return nil, nil, nil, "", Eligibility{}, err
	}
	doc, oldNode, err = Read(path)
	if err != nil {
		return fail(err)
	}
	if err = validateRetainedAuthority(path, workflowDir, doc); err != nil {
		return fail(err)
	}
	status, err = entityStatus(path)
	if err != nil {
		return fail(err)
	}
	record, err := recordForStage(doc, status)
	if err != nil {
		return fail(err)
	}
	inputState := reviewedInputUnknown
	if record != nil && len(record.Attempts) > 0 {
		inputState = inspectReviewedInput(path, record.Attempts[len(record.Attempts)-1].Briefing)
	}
	eligibility = EvaluateEligibility(doc, status, inputState == reviewedInputCurrent)
	if inputState == reviewedInputUnknown && eligibility.Condition == "stale" {
		eligibility.Condition = "ineligible"
	}
	if record == nil || len(record.Attempts) == 0 || !eligibility.Eligible ||
		eligibility.Action != "advance" || eligibility.ApplicationState != "pending" {
		return fail(fmt.Errorf("entity carries no binding pending terminal-target approval (condition %q)", eligibility.Condition))
	}
	if !applicationTargetMatches(workflowDir, status, eligibility.TargetStage) ||
		!advanceTargetTerminal(workflowDir, status, eligibility.TargetStage) {
		return fail(fmt.Errorf("pending application %s targets %q, which is not the workflow's terminal stage", eligibility.Attempt, eligibility.TargetStage))
	}
	return unlock, doc, oldNode, status, eligibility, nil
}
