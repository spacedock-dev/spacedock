// ABOUTME: behavioral guard — unknown dispatch subcommands and missing required
// ABOUTME: flags exit non-zero with a usage diagnostic on stderr.
package dispatch

import (
	"strings"
	"testing"
)

// TestUnknownSubcommandGuard asserts an unknown subcommand and a bare dispatch
// invocation both exit non-zero with usage on stderr.
func TestUnknownSubcommandGuard(t *testing.T) {
	t.Run("unknown", func(t *testing.T) {
		res := runNative("", "frobnicate")
		if res.exit == 0 {
			t.Errorf("unknown subcommand exited 0")
		}
		if !strings.Contains(res.stderr, "unknown dispatch subcommand") {
			t.Errorf("unknown subcommand lacks diagnostic:\n%q", res.stderr)
		}
	})
	t.Run("no-subcommand", func(t *testing.T) {
		res := runNative("")
		if res.exit == 0 {
			t.Errorf("bare dispatch exited 0")
		}
		if !strings.Contains(res.stderr, "Usage:") {
			t.Errorf("bare dispatch lacks usage:\n%q", res.stderr)
		}
	})
}

// TestRequiredFlagGuard asserts build and show-stage-def reject a missing
// required flag with a non-zero exit and a diagnostic.
func TestRequiredFlagGuard(t *testing.T) {
	t.Run("build-missing-workflow-dir", func(t *testing.T) {
		res := runNative(`{}`, "build")
		if res.exit != 2 {
			t.Errorf("build without --workflow-dir exit=%d, want 2", res.exit)
		}
		if res.stdout != "" {
			t.Errorf("build without --workflow-dir stdout=%q, want empty", res.stdout)
		}
		if !strings.Contains(res.stderr, "requires --workflow-dir") {
			t.Errorf("build missing-flag diagnostic lacks required text:\n%q", res.stderr)
		}
	})
	t.Run("show-stage-def-missing-stage", func(t *testing.T) {
		res := runNative("", "show-stage-def", "--workflow-dir", "/tmp")
		if res.exit != 2 {
			t.Errorf("show-stage-def without --stage exit=%d, want 2", res.exit)
		}
		if res.stdout != "" {
			t.Errorf("show-stage-def without --stage stdout=%q, want empty", res.stdout)
		}
		if !strings.Contains(res.stderr, "requires --workflow-dir and --stage") {
			t.Errorf("show-stage-def missing-flag diagnostic lacks required text:\n%q", res.stderr)
		}
	})
}
