package ensigncycle

import "testing"

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
