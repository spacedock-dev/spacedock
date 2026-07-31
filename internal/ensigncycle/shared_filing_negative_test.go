package ensigncycle

import (
	"os"
	"os/exec"
	"testing"
)

func TestAssertFilingViaNewUsesExecutedArgv(t *testing.T) {
	ledger := newTestInvocationLedger(t, writeSuccessfulLedgerTarget(t))
	runLedgerShell(t, ledger, `echo 'spacedock new wire-the-thing'; echo 'printf body | "$SPACEDOCK_BIN" new wire-the-thing'`)
	if err := assertFilingViaNew(ledger.read(t), filingSlug); err == nil {
		t.Fatal("echoed or quoted pipeline text satisfied atomic filing without an execution")
	}

	runLedgerShell(t, ledger, `"$SPACEDOCK_BIN" new wire-the-thing --workflow-dir .`)
	if err := assertFilingViaNew(ledger.read(t), filingSlug); err != nil {
		t.Fatalf("actual shim-recorded atomic filing failed: %v", err)
	}
}

func TestAssertFilingViaNewRecordsPR680QuotedLauncherExecution(t *testing.T) {
	ledger := newTestInvocationLedger(t, writeSuccessfulLedgerTarget(t))
	runLedgerShell(t, ledger, `printf '%s\n' '---' 'title: Wire The Thing' 'status: backlog' '---' | "${SPACEDOCK_BIN:-spacedock}" new wire-the-thing`)
	if err := assertFilingViaNew(ledger.read(t), filingSlug); err != nil {
		t.Fatalf("PR #680 quoted launcher execution was not graded from argv: %v", err)
	}
}

func TestAssertFilingViaNewRejectsMalformedQuotedLauncher(t *testing.T) {
	ledger := newTestInvocationLedger(t, writeSuccessfulLedgerTarget(t))
	cmd := exec.Command("/bin/sh", "-c", `printf body | "${SPACEDOCK_BIN:-spacedock} new wire-the-thing"`)
	cmd.Env = ledger.instrumentEnv(os.Environ())
	if err := cmd.Run(); err == nil {
		t.Fatal("malformed quoted launcher unexpectedly executed")
	}
	if err := assertFilingViaNew(ledger.read(t), filingSlug); err == nil {
		t.Fatal("malformed quoted launcher satisfied filing without executed argv")
	}
}

func TestAssertFilingViaNewRejectsManualIDPreview(t *testing.T) {
	ledger := newTestInvocationLedger(t, writeSuccessfulLedgerTarget(t))
	runLedgerShell(t, ledger, `"$SPACEDOCK_BIN" new wire-the-thing --workflow-dir .; "$SPACEDOCK_BIN" status --next-id --workflow-dir .`)
	if err := assertFilingViaNew(ledger.read(t), filingSlug); err == nil {
		t.Fatal("actual --next-id execution was hidden by an atomic new execution")
	}
}

func TestAssertFilingViaNewAcceptsStatusAlias(t *testing.T) {
	ledger := newTestInvocationLedger(t, writeSuccessfulLedgerTarget(t))
	runLedgerShell(t, ledger, `"$SPACEDOCK_BIN" status --new wire-the-thing --workflow-dir .`)
	if err := assertFilingViaNew(ledger.read(t), filingSlug); err != nil {
		t.Fatalf("actual shim-recorded --new alias failed: %v", err)
	}
}
