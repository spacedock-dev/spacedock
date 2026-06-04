// ABOUTME: Fixture-backed smoke for safehouse-shaped environment propagation.
// ABOUTME: Proves SPACEDOCK_BIN can reach an inner child through the wrapper argv shape.
package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafehouseShapePreservesSpacedockBinToInnerChildFixture(t *testing.T) {
	dir := t.TempDir()
	safehouse := filepath.Join(dir, "safehouse")
	if err := os.WriteFile(safehouse, []byte(`#!/bin/sh
while [ "$#" -gt 0 ] && [ "$1" != "--" ]; do shift; done
if [ "$1" = "--" ]; then shift; fi
exec "$@"
`), 0o755); err != nil {
		t.Fatal(err)
	}
	inner := filepath.Join(dir, "inner-env")
	if err := os.WriteFile(inner, []byte(`#!/bin/sh
printf '%s\n' "SPACEDOCK_BIN=${SPACEDOCK_BIN:-}"
`), 0o755); err != nil {
		t.Fatal(err)
	}

	bin := filepath.Join(dir, "spacedock")
	cmd := exec.Command(safehouse, "--trust-workdir-config", "--", inner)
	cmd.Env = append(os.Environ(), spacedockBinEnv+"="+bin)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("safehouse-shaped fixture failed: %v\n%s", err, out)
	}
	if got, want := strings.TrimSpace(string(out)), spacedockBinEnv+"="+bin; got != want {
		t.Fatalf("inner child env = %q, want %q", got, want)
	}
}
