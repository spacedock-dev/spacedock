// ABOUTME: The terminal delivery envelope — merge guard's sole consumer of a
// ABOUTME: pending terminal-target approval. One candidate replacement carries
// the authority spend (pending->consumed, or pending->superseded on --rework),
// the status move, and the delivery/terminal fields together.
package status

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

// envelopeWriteFn is the ONE atomic candidate replacement at the delivery
// envelope's write site (atomicWrite in production). Declared as a var so the
// package-level AC-1 assertion can observe that a single candidate carries all
// four terminal field changes at once, and can fail the replacement
// mid-envelope to prove the original bytes stay intact.
var envelopeWriteFn = atomicWrite

// setEnvelope carries merge guard's delivery-envelope extensions over runSet:
// the gates-subtree application mutation (the exactly-once authority spend or
// supersede) rides the SAME candidate replacement as the status/verdict/
// completed/delivery-state scalar fields, so authority, status, and terminal
// fields can never desync. A nil envelope is the plain --set path.
type setEnvelope struct {
	updates      []fieldUpdate
	gatesMutator func(*yaml.Node) (*yaml.Node, error)
	// sole-consumer envelope: merge guard legitimately fuses
	// the delivery-state retirement (mod-block/pr clear) with the terminal
	// transition in its ONE replacement — the hand --set guard that forces
	// those steps apart exists to police the CLI writer merge guard replaces.
	ceremony bool
}

// terminalApproval binds the pending terminal-target application merge guard is
// about to spend (or supersede): it is consume's own eligibility predicate
// (binding approval, digest-current briefing, pending state, status at the
// gated stage) plus the same terminal-target predicate consume routes on.
type terminalApproval struct {
	eligibility gates.Eligibility
}

// bindingTerminalApproval reports the current record's binding, still-pending
// terminal-target application — or nil when the entity carries none (an entity
// with no gates record, a consumed/superseded application, or a non-terminal
// target follows the pre-envelope finalize path unchanged).
func bindingTerminalApproval(entityPath, workflowDir, currentStatus string) *terminalApproval {
	elig, err := gates.EligibilityFileAt(entityPath, workflowDir)
	if err != nil || !elig.Eligible || elig.Action != "advance" || elig.ApplicationState != "pending" {
		return nil
	}
	if !gates.AdvanceTargetTerminal(workflowDir, currentStatus, elig.TargetStage) {
		return nil
	}
	return &terminalApproval{eligibility: elig}
}

// gatesApplicationMutator builds the gates-subtree mutator for the envelope on
// top of the proven exactly-once guard: the mutated document must be the
// byte-exact from/to application-state swap that ValidateApplicationMutation
// proves, with every other gates field untouched.
func gatesApplicationMutator(attemptID, from, to string) func(*yaml.Node) (*yaml.Node, error) {
	return func(gatesNode *yaml.Node) (*yaml.Node, error) {
		encoded, err := yaml.Marshal(gatesNode)
		if err != nil {
			return nil, fmt.Errorf("encode gates for the delivery envelope: %w", err)
		}
		var doc gates.Document
		decoder := yaml.NewDecoder(strings.NewReader(string(encoded)))
		decoder.KnownFields(true)
		if err := decoder.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode gates for the delivery envelope: %w", err)
		}
		found := false
		for ri := range doc.Records {
			for ai := range doc.Records[ri].Attempts {
				app := doc.Records[ri].Attempts[ai].Application
				if doc.Records[ri].Attempts[ai].ID == attemptID && app != nil && app.State == from {
					app.State = to
					found = true
				}
			}
		}
		if !found {
			return nil, fmt.Errorf("application attempt %s is not %s — refusing the envelope write", attemptID, from)
		}
		if err := gates.ValidateApplicationMutation(gatesNode, &doc, attemptID, from, to); err != nil {
			return nil, err
		}
		gatesBlock, err := yaml.Marshal(struct {
			Gates *gates.Document `yaml:"gates"`
		}{Gates: &doc})
		if err != nil {
			return nil, err
		}
		var replaced yaml.Node
		if err := yaml.Unmarshal(gatesBlock, &replaced); err != nil {
			return nil, fmt.Errorf("reparse envelope gates: %w", err)
		}
		if len(replaced.Content) == 0 || len(replaced.Content[0].Content) < 2 {
			return nil, fmt.Errorf("reparse envelope gates: empty mapping")
		}
		return replaced.Content[0].Content[1], nil
	}
}

// deliveryUpdates composes the delivery-state retirement in front of the
// terminal/route scalar fields: only the delivery fields actually recorded
// (mod-block / pr) are cleared, in the same replacement.
func deliveryUpdates(clearModBlock, clearPR bool, core ...fieldUpdate) []fieldUpdate {
	var clears []fieldUpdate
	if clearModBlock {
		clears = append(clears, fieldUpdate{field: "mod-block", value: "", hasValue: true})
	}
	if clearPR {
		clears = append(clears, fieldUpdate{field: "pr", value: "", hasValue: true})
	}
	return append(clears, core...)
}

// spendEnvelope is the success envelope: ONE candidate replacement carrying
// application.state pending->consumed, the terminal status, the verdict, the
// completed stamp, and the delivery-state retirement (mod-block / pr cleared
// when recorded).
func spendEnvelope(attemptID, terminal, verdict string, clearModBlock, clearPR bool) *setEnvelope {
	return &setEnvelope{
		updates: deliveryUpdates(clearModBlock, clearPR,
			fieldUpdate{field: "status", value: terminal, hasValue: true},
			fieldUpdate{field: "verdict", value: verdict, hasValue: true},
			fieldUpdate{field: "completed", hasValue: false},
		),
		gatesMutator: gatesApplicationMutator(attemptID, "pending", "consumed"),
		ceremony:     true,
	}
}

// reworkEnvelope is the --rework envelope: ONE candidate replacement carrying
// application.state pending->superseded, status := the record stage's declared
// feedback-to, and the delivery-state clear (mod-block / pr when recorded).
// Verdict/completed are untouched — never written pre-delivery.
func reworkEnvelope(attemptID, feedbackTo string, clearModBlock, clearPR bool) *setEnvelope {
	return &setEnvelope{
		updates: deliveryUpdates(clearModBlock, clearPR,
			fieldUpdate{field: "status", value: feedbackTo, hasValue: true},
		),
		gatesMutator: gatesApplicationMutator(attemptID, "pending", "superseded"),
	}
}

// declaredReworkTarget validates the --rework route: the record stage's
// declared feedback-to must exist, be a stage the workflow defines, and be
// non-terminal. Unlike the gate-rejection routing (which silently falls back
// to the same stage on a missing feedback-to and never validates the target),
// --rework refuses on every malformed declaration — the send-back must land
// somewhere real.
func declaredReworkTarget(definitionDir, recordStage string) (string, error) {
	stages := parseStagesBlock(filepath.Join(definitionDir, "README.md"))
	var declared string
	found := false
	for _, s := range stages {
		if s.Name == recordStage {
			declared = s.optional["feedback-to"]
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("record stage %q is not defined in workflow %s", recordStage, definitionDir)
	}
	if strings.TrimSpace(declared) == "" {
		return "", fmt.Errorf("record stage %q declares no feedback-to — merge guard --rework refuses to route without a declared send-back target", recordStage)
	}
	for _, s := range stages {
		if s.Name == declared {
			if s.terminal {
				return "", fmt.Errorf("record stage %q declares feedback-to %q, which is a terminal stage — merge guard --rework refuses to route a rework send-back into terminal status", recordStage, declared)
			}
			return declared, nil
		}
	}
	return "", fmt.Errorf("record stage %q declares feedback-to %q, which is not a stage defined in workflow %s", recordStage, declared, definitionDir)
}

// hasConsumedTerminalApplication reports whether the entity carries a consumed
// terminal-target application: the delivery envelope retires the merge sentinel
// (pr) in the same replacement that spends the binding approval, so post-envelope
// the sentinel is gone and the CONSUMED application itself is the delivery proof
// (merge guard is the sole terminal consumer — a consumed terminal application is
// unreachable without its delivery proof). Archive- and guard-side readers use
// this to recognize the post-envelope shape.
func hasConsumedTerminalApplication(sourcePath string) bool {
	doc, _, err := gates.Read(sourcePath)
	if err != nil {
		return false
	}
	for i := range doc.Records {
		if doc.Records[i].ID != doc.Current.Gate {
			continue
		}
		if attempts := doc.Records[i].Attempts; len(attempts) > 0 {
			app := attempts[len(attempts)-1].Application
			return app != nil && app.Action == "advance" && app.State == "consumed"
		}
	}
	return false
}

// pendingTerminalApplicationFor is the shared gate used by the terminal --set
// refusal (runSet): the entity carries a binding pending terminal-target
// approval, so merge guard is the only writer allowed to terminalize it.
func pendingTerminalApplicationFor(entityPath, workflowDir, currentStatus string) bool {
	return bindingTerminalApproval(entityPath, workflowDir, currentStatus) != nil
}
