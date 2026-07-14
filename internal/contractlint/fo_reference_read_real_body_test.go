package contractlint

import (
	"os"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/ensigncycle"
)

func TestFOReferenceReadAcceptsRealSharedCoreBody(t *testing.T) {
	body, err := os.ReadFile("../../skills/first-officer/references/first-officer-shared-core.md")
	if err != nil {
		t.Fatal(err)
	}
	// Exercise the whole canonical body: it includes ordinary failure-language prose,
	// which must not override the host's structured success signal.
	const anchor = "# First Officer Shared Core"
	if !ensigncycle.ReferenceReadSucceeded(true, string(body), anchor) {
		t.Fatal("structured successful read of the real shared core was rejected")
	}
	if ensigncycle.ReferenceReadSucceeded(false, string(body), anchor) {
		t.Fatal("structured failed read of the real shared core was accepted")
	}
	if ensigncycle.ReferenceReadSucceeded(true, string(body), "# Not The Shared Core") {
		t.Fatal("successful transport without the canonical document anchor was accepted")
	}
}
