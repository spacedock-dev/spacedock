// ABOUTME: help-routing guards for dispatch subcommands with required runtime flags.
// ABOUTME: Help must render before operational flag validation runs.
package dispatch

import (
	"strings"
	"testing"
)

func TestDispatchBuildHelpBeforeRequiredFlags(t *testing.T) {
	for _, helpFlag := range []string{"--help", "-h"} {
		t.Run(helpFlag, func(t *testing.T) {
			res := runNative("", "build", helpFlag)
			if res.exit != 0 {
				t.Fatalf("dispatch build %s exit=%d, want 0\nstderr=%q", helpFlag, res.exit, res.stderr)
			}
			if res.stderr != "" {
				t.Fatalf("dispatch build %s stderr=%q, want empty", helpFlag, res.stderr)
			}
			assertContainsAll(t, res.stdout,
				"Usage:",
				"spacedock dispatch build --workflow-dir DIR",
				"stdin JSON",
				"--workflow-dir",
				"schema_version",
				"entity_path",
				"workflow_dir",
				"stage",
				"checklist",
			)
			assertNotContains(t, res.stdout, "requires --workflow-dir")
		})
	}
}

func TestDispatchShowStageDefHelpBeforeRequiredFlags(t *testing.T) {
	for _, helpFlag := range []string{"--help", "-h"} {
		t.Run(helpFlag, func(t *testing.T) {
			res := runNative("", "show-stage-def", helpFlag)
			if res.exit != 0 {
				t.Fatalf("dispatch show-stage-def %s exit=%d, want 0\nstderr=%q", helpFlag, res.exit, res.stderr)
			}
			if res.stderr != "" {
				t.Fatalf("dispatch show-stage-def %s stderr=%q, want empty", helpFlag, res.stderr)
			}
			assertContainsAll(t, res.stdout,
				"Usage:",
				"spacedock dispatch show-stage-def --workflow-dir DIR --stage STAGE",
				"--workflow-dir",
				"--stage",
			)
			assertNotContains(t, res.stdout, "requires --workflow-dir")
			assertNotContains(t, res.stdout, "requires --workflow-dir and --stage")
		})
	}
}

func assertContainsAll(t *testing.T, got string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}

func assertNotContains(t *testing.T, got, bad string) {
	t.Helper()
	if strings.Contains(got, bad) {
		t.Fatalf("output contains %q:\n%s", bad, got)
	}
}
