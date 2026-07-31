// ABOUTME: AC-1 envelope atomics (package level): the terminal delivery envelope
// ABOUTME: must be ONE candidate replacement carrying application+status+verdict+
// completed, and a mid-envelope failure leaves the original bytes intact.
package status

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

// envelopeFixture builds an inline workflow with a registered merge hook and an
// entity whose validation gate carries a real recorded binding approval:
// implementation (initial) -> validation (gate, feedback-to) -> done (terminal).
// It returns the workflow root and the entity path. The repos are git-initialized
// because the ceremony's mutations resolve git_root.
func envelopeFixture(t *testing.T, extraReadme string) (string, string) {
	t.Helper()
	root := t.TempDir()
	readme := "---\nid-style: slug\n" + extraReadme + "stages:\n  states:\n" +
		"    - name: implementation\n      initial: true\n" +
		"    - name: validation\n      gate: true\n      feedback-to: implementation\n" +
		"    - name: done\n      terminal: true\n---\n# Workflow\n"
	writeFile(t, filepath.Join(root, "README.md"), readme)
	writeFile(t, filepath.Join(root, "_mods", "pr-merge.md"), "---\nname: pr-merge\ndescription: stub merge hook.\n---\n\n# PR Merge\n\n## Hook: merge\n\n(stub — registration only)\n")
	room := filepath.Join(root, "review", "validation", "briefing-1")
	briefing := filepath.Join(room, "briefing.json")
	writeFile(t, briefing,
		`{"type":"Briefing","version":"1","id":"briefing:task:validation:attempt-1:revision-1","question":"ship?","artifacts":[{"id":"artifact:1","uri":"artifact.md","rev":"sha256:`+strings.Repeat("a", 64)+`"}]}`)
	briefingBytes, err := os.ReadFile(briefing)
	if err != nil {
		t.Fatal(err)
	}
	digest, err := gates.CanonicalDigest(briefingBytes)
	if err != nil {
		t.Fatal(err)
	}
	entity := filepath.Join(root, "task.md")
	writeFile(t, entity, "---\nid: task\nstatus: validation\ntitle: Task\ngates:\n"+
		"  version: 1\n  current: {gate: 'gate:task:validation'}\n  records:\n"+
		"    - id: gate:task:validation\n      stage: validation\n      attempts:\n"+
		"        - id: gate-attempt:task-validation-1\n"+
		"          briefing: {id: 'briefing:task:validation:attempt-1:revision-1', digest: '"+digest+"', digest-domain: canonical-bytes, room-ref: ./review/validation/briefing-1/briefing.json}\n"+
		"---\n# Task\n")
	if err := gates.RecordSemantic(entity, gates.RecordInput{
		Decision: "approve", Actor: "person:captain", WorkflowDir: root,
	}); err != nil {
		t.Fatal(err)
	}
	gitInit(t, root)
	return root, entity
}

func driveEnvelopeGuard(t *testing.T, root string, args ...string) (string, string, int) {
	t.Helper()
	var out, errBuf bytes.Buffer
	full := append([]string{"--workflow-dir", root}, args...)
	code := MergeGuard(full, root, &out, &errBuf)
	return out.String(), errBuf.String(), code
}

// TestDeliveryEnvelopeIsOneAtomicCandidate (AC-1 envelope atomics): with delivery
// proven (pr-merge sentinel), the guard's finalize lands application.state
// pending->consumed + status=done + verdict + completed + delivery-state
// retirement in ONE candidate replacement: the write site is invoked exactly
// once, its candidate carries all four field changes at once, and the entity's
// pre-envelope bytes stay intact until the candidate applies. A split-write
// envelope (application then status/verdict/completed) turns this red.
func TestDeliveryEnvelopeIsOneAtomicCandidate(t *testing.T) {
	root, entity := envelopeFixture(t, "")
	if rc := emitSetForTest(t, root, "task", []fieldUpdate{{field: "pr", value: "pr-merge:77", hasValue: true}}); rc != 0 {
		t.Fatalf("record merge sentinel rc=%d", rc)
	}
	original, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}

	type call struct {
		path string
		data []byte
	}
	var calls []call
	prev := envelopeWriteFn
	envelopeWriteFn = func(path string, data []byte) error {
		// The original bytes must still be on disk at the moment the candidate
		// replacement applies (no pre-write mutation).
		stillOriginal, rerr := os.ReadFile(path)
		if rerr != nil || !bytes.Equal(stillOriginal, original) {
			t.Errorf("original bytes were gone when the candidate replacement applied (rerr=%v)", rerr)
		}
		calls = append(calls, call{path, data})
		return prev(path, data)
	}
	t.Cleanup(func() { envelopeWriteFn = prev })

	out, errOut, code := driveEnvelopeGuard(t, root, "task", "--verdict", "passed")
	if code != 0 {
		t.Fatalf("envelope finalize exit=%d stdout=%q stderr=%q", code, out, errOut)
	}
	if len(calls) != 1 {
		t.Fatalf("envelope write site invoked %d times, want exactly 1 candidate replacement", len(calls))
	}
	candidate := string(calls[0].data)
	for _, want := range []string{"state: consumed", "status: done", "verdict: passed", "\npr:\n"} {
		if !strings.Contains(candidate, want) {
			t.Errorf("candidate replacement missing %q:\n%s", want, candidate)
		}
	}
	if !strings.Contains(candidate, "completed: 2") {
		t.Errorf("candidate replacement carries no completed stamp:\n%s", candidate)
	}
	if strings.Contains(candidate, "pr-merge:77") {
		t.Errorf("candidate replacement did not retire recorded delivery state:\n%s", candidate)
	}
	final, err := os.ReadFile(entity)
	if err == nil && !bytes.Equal(final, calls[0].data) {
		t.Fatalf("final entity bytes differ from the single candidate replacement")
	}
}

// TestDeliveryEnvelopeFailureLeavesOriginalBytesIntact (AC-1 envelope atomics):
// failing the one replacement mid-envelope leaves the entity's original bytes
// untouched — no partial lead writes (no application flip, no terminal status).
func TestDeliveryEnvelopeFailureLeavesOriginalBytesIntact(t *testing.T) {
	root, entity := envelopeFixture(t, "")
	if rc := emitSetForTest(t, root, "task", []fieldUpdate{{field: "pr", value: "pr-merge:77", hasValue: true}}); rc != 0 {
		t.Fatalf("record merge sentinel rc=%d", rc)
	}
	original, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	prev := envelopeWriteFn
	envelopeWriteFn = func(path string, data []byte) error {
		return errors.New("injected mid-envelope replacement failure")
	}
	t.Cleanup(func() { envelopeWriteFn = prev })

	out, errOut, code := driveEnvelopeGuard(t, root, "task", "--verdict", "passed")
	if code == 0 {
		t.Fatalf("envelope with an injected replacement failure must not exit 0 (stdout=%q)", out)
	}
	if out != "" && errOut == "" {
		t.Errorf("fail path lost stderr: %q", out)
	}
	after, err := os.ReadFile(entity)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, original) {
		t.Fatalf("mid-envelope failure changed the entity bytes:\nbefore:\n%s\nafter:\n%s", original, after)
	}
}

// emitSetForTest stages a merge-sentinel pr field through the real --set path.
func emitSetForTest(t *testing.T, root, slug string, updates []fieldUpdate) int {
	t.Helper()
	r, err := resolveRoots(root, "")
	if err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	rc := emitSet(r, slug, updates, &stderr)
	if rc != 0 {
		t.Logf("emitSetForTest stderr: %s", stderr.String())
	}
	return rc
}
