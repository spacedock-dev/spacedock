// ABOUTME: Tests for the two gates-owned locked terminal delivery operations:
// ABOUTME: one candidate replacement carries authority+terminal fields (with
// ABOUTME: failure-injection byte-clean proof), and the rework route's refusals.
package gates

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// approveTerminal records a captain approval in the applicationTerminalWorkflow
// fixture and returns the entity path with a binding pending terminal-target
// application in force.
func approveTerminal(t *testing.T) (root, entity string) {
	t.Helper()
	root, entity = applicationTerminalWorkflow(t)
	if err := RecordSemantic(entity, RecordInput{Decision: "approve", Actor: "person:captain", WorkflowDir: root}); err != nil {
		t.Fatal(err)
	}
	return root, entity
}

// TestFinalizeTerminalApprovalOneLockedReplacement is the moved envelope-atomics
// assertion, now observed AT THE GATES WRITER: the terminal delivery write is ONE
// candidate locked replacement carrying application.state pending→consumed, the
// terminal status, the verdict, and the completed stamp at once; an injected
// failure at the write site leaves the original entity bytes intact.
func TestFinalizeTerminalApprovalOneLockedReplacement(t *testing.T) {
	root, entity := approveTerminal(t)
	before, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}

	var candidates [][]byte
	inject := fmt.Errorf("injected write failure")
	failOnce := true
	originalWrite := entityDocumentWriteFn
	entityDocumentWriteFn = func(path string, data []byte) error {
		candidates = append(candidates, append([]byte(nil), data...))
		if failOnce {
			failOnce = false
			return inject
		}
		return originalWrite(path, data)
	}
	defer func() { entityDocumentWriteFn = originalWrite }()

	if _, err := FinalizeTerminalApproval(entity, root, "passed", "2026-08-01T00:00:00Z"); !errors.Is(err, inject) {
		t.Fatalf("finalize with failing write = %v, want the injected failure", err)
	}
	after, err := os.ReadFile(entity)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("failed replacement changed the entity (read err=%v)", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("failed finalize produced %d candidate replacements, want 1", len(candidates))
	}
	// The ONE candidate already carries all four field changes.
	candidate := string(candidates[0])
	for _, want := range []string{"status: done", "verdict: passed", "completed: 2026-08-01T00:00:00Z", "state: consumed"} {
		if !strings.Contains(candidate, want) {
			t.Fatalf("candidate replacement missing %q:\n%s", want, candidate)
		}
	}

	target, err := FinalizeTerminalApproval(entity, root, "passed", "2026-08-01T00:00:00Z")
	if err != nil || target != "done" {
		t.Fatalf("finalize = %q, %v", target, err)
	}
	if len(candidates) != 2 {
		t.Fatalf("finalize produced %d candidate replacements total, want exactly 2 (the failed attempt plus ONE successful write)", len(candidates))
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	app := doc.Records[0].Attempts[len(doc.Records[0].Attempts)-1].Application
	if app.State != "consumed" {
		t.Fatalf("application state = %q, want consumed", app.State)
	}
	finalBytes, _ := os.ReadFile(entity)
	for _, want := range []string{"status: done", "verdict: passed", "completed: 2026-08-01T00:00:00Z"} {
		if !strings.Contains(string(finalBytes), want) {
			t.Fatalf("final entity missing %q:\n%s", want, finalBytes)
		}
	}
	if doc.Current.Gate != "gate:task:ideation" {
		t.Fatalf("finalize moved gates.current: %q", doc.Current.Gate)
	}
	if !strings.Contains(string(finalBytes), "Body.\n") {
		t.Fatalf("finalize disturbed the entity body:\n%s", finalBytes)
	}
}

// TestTerminalDeliveryRefusalsByteClean: gate-less entities report
// ErrNoGateRecord (the ONLY legacy-path condition), and spent/superseded
// authority refuses byte-clean for both operations.
func TestTerminalDeliveryRefusalsByteClean(t *testing.T) {
	t.Run("gate-less-entity", func(t *testing.T) {
		root := t.TempDir()
		readme := "---\nstages:\n  states:\n    - name: validation\n      gate: true\n    - name: done\n      terminal: true\n---\n# Workflow\n"
		if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
			t.Fatal(err)
		}
		entity := filepath.Join(root, "task.md")
		content := []byte("---\nstatus: validation\ntitle: No gates\n---\n# Task\n")
		if err := os.WriteFile(entity, content, 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := FinalizeTerminalApproval(entity, root, "passed", "2026-08-01T00:00:00Z"); !errors.Is(err, ErrNoGateRecord) {
			t.Fatalf("finalize on gate-less entity = %v, want ErrNoGateRecord", err)
		}
		if _, err := SupersedeTerminalApproval(entity, root); !errors.Is(err, ErrNoGateRecord) {
			t.Fatalf("supersede on gate-less entity = %v, want ErrNoGateRecord", err)
		}
		after, _ := os.ReadFile(entity)
		if !bytes.Equal(content, after) {
			t.Fatal("gate-less refusals changed the entity")
		}
	})

	t.Run("superseded-authority", func(t *testing.T) {
		root, entity := approveTerminal(t)
		if _, err := SupersedeTerminalApproval(entity, root); err != nil {
			t.Fatal(err)
		}
		before, _ := os.ReadFile(entity)
		if _, err := FinalizeTerminalApproval(entity, root, "passed", "2026-08-01T00:00:00Z"); err == nil ||
			!strings.Contains(err.Error(), "no binding pending terminal-target approval") {
			t.Fatalf("finalize on superseded authority = %v, want byte-clean refusal", err)
		}
		if after, _ := os.ReadFile(entity); !bytes.Equal(before, after) {
			t.Fatal("refused finalize changed the entity")
		}
	})

	t.Run("digest-stale-authority", func(t *testing.T) {
		root, entity := approveTerminal(t)
		before, _ := os.ReadFile(entity)
		// Disturb the retained briefing room: retained-authority validation
		// fails, so both operations refuse byte-clean rather than spend.
		room := filepath.Join(root, "review", "ideation", "briefing-1", "briefing.json")
		if err := os.WriteFile(room, []byte("tampered"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := FinalizeTerminalApproval(entity, root, "passed", "2026-08-01T00:00:00Z"); err == nil {
			t.Fatal("finalize with tampered briefing must refuse")
		}
		if _, err := SupersedeTerminalApproval(entity, root); err == nil {
			t.Fatal("supersede with tampered briefing must refuse")
		}
		if after, _ := os.ReadFile(entity); !bytes.Equal(before, after) {
			t.Fatal("digest-stale refusals changed the entity")
		}
	})
}

// TestSupersedeTerminalApprovalRoutesDeclaredFeedback pins the rework route:
// ONE locked replacement writes pending→superseded + status := declared
// feedback-to, with verdict/completed never written pre-delivery; missing,
// undefined, or terminal feedback-to declarations refuse closed byte-clean.
func TestSupersedeTerminalApprovalRoutesDeclaredFeedback(t *testing.T) {
	root, entity := approveTerminal(t)
	var writes int
	originalWrite := entityDocumentWriteFn
	entityDocumentWriteFn = func(path string, data []byte) error {
		writes++
		return originalWrite(path, data)
	}
	defer func() { entityDocumentWriteFn = originalWrite }()

	feedbackTo, err := SupersedeTerminalApproval(entity, root)
	if err != nil || feedbackTo != "ideation" {
		t.Fatalf("supersede = %q, %v", feedbackTo, err)
	}
	if writes != 1 {
		t.Fatalf("supersede produced %d candidate replacements, want 1", writes)
	}
	finalBytes, _ := os.ReadFile(entity)
	for _, want := range []string{"status: ideation", "state: superseded"} {
		if !strings.Contains(string(finalBytes), want) {
			t.Fatalf("reworked entity missing %q:\n%s", want, finalBytes)
		}
	}
	for _, banned := range []string{"verdict:", "completed:"} {
		if strings.Contains(string(finalBytes), banned) {
			t.Fatalf("--rework wrote terminal field %q pre-delivery:\n%s", banned, finalBytes)
		}
	}
	doc, _, err := Read(entity)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Current.Gate != "gate:task:ideation" {
		t.Fatalf("--rework moved gates.current: %q", doc.Current.Gate)
	}

	// Refusal matrix: each malformed declaration refuses closed, byte-clean.
	for _, tc := range []struct {
		name       string
		stageFlags string
		want       string
	}{
		{"no-feedback-to", "      initial: true\n      gate: true\n", "declares no feedback-to"},
		{"undefined-feedback-to", "      initial: true\n      gate: true\n      feedback-to: neverland\n", "not a stage defined"},
		{"terminal-feedback-to", "      initial: true\n      gate: true\n      feedback-to: done\n", "terminal stage"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, entity := approveTerminal(t)
			readme := "---\nstages:\n  states:\n    - name: ideation\n" + tc.stageFlags + "    - name: done\n      terminal: true\n---\n# Workflow\n"
			if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(readme), 0o644); err != nil {
				t.Fatal(err)
			}
			before, _ := os.ReadFile(entity)
			if _, err := SupersedeTerminalApproval(entity, root); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("supersede with %s = %v, want refusal naming %q", tc.name, err, tc.want)
			}
			if _, err := DeclaredReworkTarget(root, "ideation"); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("DeclaredReworkTarget with %s = %v, want refusal naming %q", tc.name, err, tc.want)
			}
			if after, _ := os.ReadFile(entity); !bytes.Equal(before, after) {
				t.Fatalf("refused supersede with %s changed the entity", tc.name)
			}
		})
	}
}
