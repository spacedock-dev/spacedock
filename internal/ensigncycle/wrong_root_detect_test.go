// ABOUTME: Offline test for the wrong-root boot detector — proves it reds on a
// ABOUTME: simulated FO wander off the fixture root and passes a fixture-rooted boot.
package ensigncycle

import (
	"strings"
	"testing"
)

// streamLine builds one assistant-with-Bash-tool_use stream-json line carrying the
// given shell command, the shape detectWrongRootBoot scans. Kept tiny so the cases
// read as "this command in the boot stream" without a fixture file.
func streamLine(command string) string {
	return `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash","input":{"command":` +
		mustJSONString(command) + `}}]}}`
}

func mustJSONString(s string) string {
	// A minimal JSON string encoder for the test fixtures (the commands here carry
	// no control chars), so the line is valid stream-json without importing the
	// encoder into the test body.
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

// TestDetectWrongRootBoot covers the pure detector both ways: it reds on a stream
// where the FO `cd`s off the fixture root (the PR #365 opus wander) or boots a
// workflow-dir outside it, and it passes on a fixture-rooted boot whose only
// real-repo paths are the legitimate --plugin-dir contract reads.
func TestDetectWrongRootBoot(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveEnsignCycle1166216625/002"
	const realRepo = "/home/runner/work/spacedock/spacedock"

	t.Run("cd_away_from_fixture_root_reds", func(t *testing.T) {
		stream := strings.Join([]string{
			streamLine(`echo "CLAUDECODE=${CLAUDECODE:-unset}"`),
			streamLine(`cd ` + realRepo + ` && spacedock --version`),
			streamLine(`spacedock status --discover`),
		}, "\n")

		err := detectWrongRootBoot(stream, fixtureRoot)
		if err == nil {
			t.Fatal("detector passed a stream where the FO cd'd into the real repo — want a wrong-root error")
		}
		if !strings.Contains(err.Error(), fixtureRoot) || !strings.Contains(err.Error(), realRepo) {
			t.Errorf("error must name both expected (%q) and actual (%q): %v", fixtureRoot, realRepo, err)
		}
	})

	t.Run("boot_workflow_dir_outside_fixture_reds", func(t *testing.T) {
		stream := strings.Join([]string{
			streamLine(`spacedock --version`),
			streamLine(`spacedock status --boot --workflow-dir ` + realRepo + `/docs/dev`),
		}, "\n")

		err := detectWrongRootBoot(stream, fixtureRoot)
		if err == nil {
			t.Fatal("detector passed a boot whose --workflow-dir is the real repo — want a wrong-root error")
		}
		if !strings.Contains(err.Error(), fixtureRoot) {
			t.Errorf("error must name the expected fixture root %q: %v", fixtureRoot, err)
		}
	})

	t.Run("workflow_readme_outside_fixture_reds", func(t *testing.T) {
		// The FO discovered cwd from the real repo and read its docs/dev workflow
		// README — the wander even when no explicit --workflow-dir is on the boot cmd.
		stream := strings.Join([]string{
			streamLine(`spacedock --version`),
			streamLine(`spacedock status --discover`),
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"` + realRepo + `/docs/dev/README.md"}}]}}`,
		}, "\n")

		err := detectWrongRootBoot(stream, fixtureRoot)
		if err == nil {
			t.Fatal("detector passed a boot that read the real repo's workflow README — want a wrong-root error")
		}
		if !strings.Contains(err.Error(), fixtureRoot) {
			t.Errorf("error must name the expected fixture root %q: %v", fixtureRoot, err)
		}
	})

	t.Run("contract_skill_read_outside_fixture_passes", func(t *testing.T) {
		// A contract-skill Read from the real-repo --plugin-dir is legitimate, NOT a
		// workflow README, so it must not false-red.
		stream := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"` + realRepo + `/skills/first-officer/references/claude-first-officer-runtime.md"}}]}}`
		if err := detectWrongRootBoot(stream, fixtureRoot); err != nil {
			t.Errorf("detector red a legitimate contract-skill Read from the plugin-dir: %v", err)
		}
	})

	t.Run("fixture_rooted_boot_passes", func(t *testing.T) {
		// A correct boot: the contract skill Reads come from the real-repo
		// --plugin-dir (legitimate), but the workflow boot stays under the fixture.
		stream := strings.Join([]string{
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"` + realRepo + `/skills/first-officer/references/first-officer-shared-core.md"}}]}}`,
			streamLine(`spacedock --version`),
			streamLine(`git rev-parse --show-toplevel`),
			streamLine(`spacedock status --boot --workflow-dir ` + fixtureRoot),
			`{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"` + fixtureRoot + `/README.md"}}]}}`,
		}, "\n")

		if err := detectWrongRootBoot(stream, fixtureRoot); err != nil {
			t.Errorf("detector red a fixture-rooted boot (plugin-dir contract reads are legitimate): %v", err)
		}
	})

	t.Run("cd_to_fixture_subdir_passes", func(t *testing.T) {
		// `cd` INTO the fixture (or a subdir) is not a wander.
		stream := strings.Join([]string{
			streamLine(`cd ` + fixtureRoot + ` && spacedock status --discover`),
		}, "\n")

		if err := detectWrongRootBoot(stream, fixtureRoot); err != nil {
			t.Errorf("detector red a cd into the fixture root itself: %v", err)
		}
	})

	t.Run("compound_cd_into_fixture_passes", func(t *testing.T) {
		// Latest-opus chains its boot in one bash call: `cd <fixtureRoot>; ls; …`.
		// strings.Fields glues the trailing `;` onto the path token (`<fixtureRoot>;`),
		// which must not be read as a wander off the (correct) fixture root.
		stream := streamLine(`cd ` + fixtureRoot + `; ls -la; echo "===README==="; sed -n '1,60p' README.md`)
		if err := detectWrongRootBoot(stream, fixtureRoot); err != nil {
			t.Errorf("detector red a `;`-glued compound cd into the fixture root: %v", err)
		}
	})

	t.Run("compound_cd_into_fixture_amp_passes", func(t *testing.T) {
		// The `&&`-glued-with-no-space form (`cd <fixtureRoot>&& ls`) glues `&&` onto
		// the path token the same way; it must also not flag a wander.
		stream := streamLine(`cd ` + fixtureRoot + `&& ls`)
		if err := detectWrongRootBoot(stream, fixtureRoot); err != nil {
			t.Errorf("detector red a `&&`-glued compound cd into the fixture root: %v", err)
		}
	})

	t.Run("compound_cd_away_from_fixture_reds", func(t *testing.T) {
		// The no-false-negative direction: a genuine off-fixture compound boot must
		// STILL red, naming both roots. Guards the trim against over-stripping and
		// silently disabling detection.
		stream := streamLine(`cd ` + realRepo + `; ls -la; cat README.md`)
		err := detectWrongRootBoot(stream, fixtureRoot)
		if err == nil {
			t.Fatal("detector passed a compound boot cd'ing into the real repo — want a wrong-root error")
		}
		if !strings.Contains(err.Error(), fixtureRoot) || !strings.Contains(err.Error(), realRepo) {
			t.Errorf("error must name both expected (%q) and actual (%q): %v", fixtureRoot, realRepo, err)
		}
	})

	t.Run("empty_stream_does_not_false_red", func(t *testing.T) {
		// No Bash commands at all (e.g. a launch failure stream) is not a wrong-root
		// boot — it is a different failure the caller's own checks surface.
		if err := detectWrongRootBoot("", fixtureRoot); err != nil {
			t.Errorf("detector red an empty stream: %v", err)
		}
	})
}
