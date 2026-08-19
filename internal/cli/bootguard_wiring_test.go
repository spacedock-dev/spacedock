// ABOUTME: AC-1/AC-3 wiring proof — gate record, gate consume, and merge guard
// ABOUTME: each refuse on a stale/missing boot receipt and succeed once one is fresh; this fails if the preflight call site is ever removed.
package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spacedock-dev/spacedock/internal/status"
)

// bootGuardWiringSessionID is the pinned session id these tests resolve their
// receipt under — distinct from the ambient session running this very suite
// (see hermeticEnv), so the guard's real code path is exercised deliberately
// rather than by ambient accident.
const bootGuardWiringSessionID = "wiring-proof-session"

// writeFreshBootReceipt drops a timestamp-only receipt at gitRoot (no
// transcript field — bootGuardVerdict then cannot prove staleness and passes,
// exactly AC-2's "after status --boot re-runs ... succeeds unchanged"). The
// exact receipt format is internal/status/bootguard.go's contract; this test
// writes it directly rather than shelling to `status --boot` so the wiring
// proof does not also depend on transcript-glob resolution.
func writeFreshBootReceipt(t *testing.T, gitRoot string) {
	t.Helper()
	dir := filepath.Join(gitRoot, ".spacedock", "boot")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	line := time.Now().UTC().Format(time.RFC3339) + "\n"
	if err := os.WriteFile(filepath.Join(dir, bootGuardWiringSessionID), []byte(line), 0o644); err != nil {
		t.Fatal(err)
	}
}

// runGuardedVerb runs one guarded command under the pinned session env,
// asserting the stale-receipt refusal (exit BootStaleExitCode, no mutation of
// entityPath), then writing a fresh receipt and asserting the retry succeeds
// with wantSuccessSubstr in stdout. If the CLI's BootGuardRefuse call site were
// ever removed, the first assertion goes red — the verb would exit 0 with no
// receipt at all, which is exactly the incident this guard exists to prevent.
func runGuardedVerb(t *testing.T, gitRoot, entityPath string, args []string, wantSuccessSubstr string) {
	t.Helper()
	env := []string{"CLAUDE_CODE_SESSION_ID=" + bootGuardWiringSessionID}

	var before []byte
	if entityPath != "" {
		var err error
		before, err = os.ReadFile(entityPath)
		if err != nil {
			t.Fatal(err)
		}
	}

	var out, errOut bytes.Buffer
	code := run(context.Background(), args, env, gitRoot, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != status.BootStaleExitCode {
		t.Fatalf("%v with no boot receipt: exit=%d, want %d (stdout=%q stderr=%q)",
			args, code, status.BootStaleExitCode, out.String(), errOut.String())
	}
	if !strings.Contains(errOut.String(), "no boot receipt") || !strings.Contains(errOut.String(), "status --boot") {
		t.Fatalf("%v refusal stderr = %q, want it to name the condition and the status --boot remedy", args, errOut.String())
	}
	if entityPath != "" {
		after, err := os.ReadFile(entityPath)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(before, after) {
			t.Fatalf("%v refusal mutated %s", args, entityPath)
		}
	}

	writeFreshBootReceipt(t, gitRoot)
	out.Reset()
	errOut.Reset()
	code = run(context.Background(), args, env, gitRoot, nil, &out, &errOut, &status.NativeRunner{}, nil)
	if code != 0 {
		t.Fatalf("%v after a fresh boot receipt: exit=%d, want 0 (stdout=%q stderr=%q)", args, code, out.String(), errOut.String())
	}
	if !strings.Contains(out.String(), wantSuccessSubstr) {
		t.Fatalf("%v succeeded but stdout=%q missing %q", args, out.String(), wantSuccessSubstr)
	}
}

// TestBootGuardWiringGateRecord drives `gate record` on a freshly prepared,
// never-closed gate.
func TestBootGuardWiringGateRecord(t *testing.T) {
	root, entity := semanticDecisionFixture(t)
	runGuardedVerb(t, root, entity,
		[]string{"gate", "record", "task", "--workflow-dir", root, "--decision", "approve", "--actor", "person:captain"},
		"recorded gate=")
}

// TestBootGuardWiringGateConsume drives `gate consume` on a gate this test
// closes unguarded first (env=nil resolves no session id, so the setup close
// itself is a guard no-op) — the guard proof below is scoped to consume alone.
// unboundGateRoomFixture (not semanticDecisionFixture) is the base: it is the
// fixture gate_test.go's own record-then-consume chain
// (TestGateRequestLocatorCarriesArbitraryBriefingNameThroughRecordValidateAndEligibility)
// already proves reaches a real eligible consume.
func TestBootGuardWiringGateConsume(t *testing.T) {
	root, entity, _ := unboundGateRoomFixture(t)
	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"gate", "record", "task", "--workflow-dir", root, "--decision", "approve", "--actor", "person:captain"},
		nil, root, nil, &out, &errOut, &status.NativeRunner{}, nil); code != 0 {
		t.Fatalf("setup close exit=%d stderr=%q", code, errOut.String())
	}
	runGuardedVerb(t, root, entity,
		[]string{"gate", "consume", "task", "--workflow-dir", root},
		"route=approved-awaiting-merge")
}

// TestBootGuardWiringMergeGuard drives `merge guard` on a merge: local fixture
// whose first successful call arms the mod-block.
func TestBootGuardWiringMergeGuard(t *testing.T) {
	root := stageMergeFixture(t, "merge-local-workflow")
	entity := filepath.Join(root, "020-no-sentinel.md")
	runGuardedVerb(t, root, entity,
		[]string{"merge", "guard", "020-no-sentinel", "--verdict", "passed", "--workflow-dir", root},
		"armed")
}
