// ABOUTME: Offline test for the wrong-root boot detector — proves it reds on a
// ABOUTME: simulated FO wander off the fixture root and passes a fixture-rooted boot.
package ensigncycle

import (
	"os"
	"path/filepath"
	"strconv"
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

func bashToolLine(id, command string) string {
	return `{"type":"assistant","message":{"content":[{"type":"tool_use","id":` + mustJSONString(id) +
		`,"name":"Bash","input":{"command":` + mustJSONString(command) + `}}]}}`
}

func toolResultLine(id string, isError bool, content string) string {
	errValue := "false"
	if isError {
		errValue = "true"
	}
	return `{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":` + mustJSONString(id) +
		`,"content":` + mustJSONString(content) + `,"is_error":` + errValue + `}]}}`
}

func mustJSONString(s string) string {
	return strconv.Quote(s)
}

// TestDetectWrongRootBoot covers the pure detector both ways: it reds on a stream
// where the FO `cd`s off the fixture root (the PR #365 opus wander) or boots a
// workflow-dir outside it, and it passes on a fixture-rooted boot whose only
// real-repo paths are the legitimate --plugin-dir contract reads.
func TestDetectWrongRootBoot(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveEnsignCycle1166216625/002"
	const realRepo = "/home/runner/work/spacedock/spacedock"

	t.Run("cd_away_from_fixture_root_reds", func(t *testing.T) {
		// The cd is corroborated in the SAME command by `status --discover`, so
		// this is a genuine wander, not a probe-only cd (see
		// cd_probe_only_away_from_fixture_root_passes below for the harmless shape).
		stream := strings.Join([]string{
			streamLine(`echo "CLAUDECODE=${CLAUDECODE:-unset}"`),
			streamLine(`cd ` + realRepo + ` && spacedock status --discover`),
		}, "\n")

		err := detectWrongRootBoot(stream, fixtureRoot)
		if err == nil {
			t.Fatal("detector passed a stream where the FO cd'd into the real repo and ran a workflow command — want a wrong-root error")
		}
		if !strings.Contains(err.Error(), fixtureRoot) || !strings.Contains(err.Error(), realRepo) {
			t.Errorf("error must name both expected (%q) and actual (%q): %v", fixtureRoot, realRepo, err)
		}
		if strings.Contains(err.Error(), "CI env leak") {
			t.Errorf("error text must not assert a CI env leak as the cause (disproved by the archived PR #446 evidence): %v", err)
		}
	})

	t.Run("cd_probe_only_away_from_fixture_root_passes", func(t *testing.T) {
		// The real PR #446 shape (reworded here with realRepo for readability):
		// a bare cd escaping the fixture, corroborated by nothing but a version
		// probe. This is sonnet's speculative repo-root sniff, not a wander — see
		// TestDetectWrongRootBootRealPR446Streams for the verbatim captured commands.
		stream := streamLine(`cd ` + realRepo + ` 2>/dev/null || cd .; spacedock --version; pwd; git rev-parse --show-toplevel 2>&1`)

		if err := detectWrongRootBoot(stream, fixtureRoot); err != nil {
			t.Errorf("detector red a probe-only cd with no workflow-operative corroboration in the same command: %v", err)
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
		if strings.Contains(err.Error(), "CI env leak") {
			t.Errorf("error text must not assert a CI env leak as the cause: %v", err)
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

	t.Run("plugin_skill_readme_outside_fixture_passes", func(t *testing.T) {
		// Claude's Skill loader may read a plugin skill root README. It is outside
		// the fixture, but it is not the workflow README and must not be classified
		// as a wrong-root boot.
		stream := `{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Read","input":{"file_path":"/tmp/spacedock-live-plugin-1648179417/skills/first-officer/README.md"}}]}}`
		if err := detectWrongRootBoot(stream, fixtureRoot); err != nil {
			t.Errorf("detector red a plugin skill README read as a workflow wander: %v", err)
		}
	})

	t.Run("failed_parent_new_then_corrected_workflow_dir_passes", func(t *testing.T) {
		parent := "/tmp/TestLiveClaudeSharedScenariosfiling2421552038"
		fixture := parent + "/001"
		stream := strings.Join([]string{
			bashToolLine("toolu_bad", `cd `+parent+` && printf '%s\n' '---' 'title: Wire The Thing' 'status: backlog' '---' 'Wire the thing end to end so it is connected and functional.' | ${SPACEDOCK_BIN:-spacedock} new wire-the-thing`),
			toolResultLine("toolu_bad", true, "Exit code 1\nError: no commissioned Spacedock workflow found in "+parent),
			bashToolLine("toolu_good", `cd `+parent+` && printf '%s\n' '---' 'title: Wire The Thing' 'status: backlog' '---' 'Wire the thing end to end so it is connected and functional.' | ${SPACEDOCK_BIN:-spacedock} new wire-the-thing --workflow-dir `+fixture),
			toolResultLine("toolu_good", false, "created: "+fixture+"/wire-the-thing.md id=001"),
		}, "\n")

		if err := detectWrongRootBoot(stream, fixture); err != nil {
			t.Errorf("detector red the PR #483 opus filing stream shape even though the off-fixture command failed and the corrected --workflow-dir command landed in the fixture: %v", err)
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
		if strings.Contains(err.Error(), "CI env leak") {
			t.Errorf("error text must not assert a CI env leak as the cause: %v", err)
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

// TestDetectWrongRootBootRealPR446Commands runs the detector against the two
// captured sonnet boot commands verbatim (PR #446, run 28466995641, attempts 1
// and 2) — the exact strings that fatal'd `TestLiveClaudeSharedScenarios` at
// claude_live_runner_test.go:122 both times. Both must now pass: neither
// command carries a workflow-operative token alongside its cd, so each is
// sonnet's harmless speculative repo-root sniff, not a wander. This is the
// command-level half of AC-1; TestDetectWrongRootBootRealPR446Streams below is
// the full-stream half.
func TestDetectWrongRootBootRealPR446Commands(t *testing.T) {
	const fixtureRoot = "/tmp/TestLiveClaudeSharedScenariosfeedback-3-cycle-escalation1752814120/001"

	cases := map[string]string{
		"attempt_1": `cd /home/user/spacedock-workflow 2>/dev/null || cd .; ${SPACEDOCK_BIN:-spacedock} --version; echo "---"; pwd; git rev-parse --show-toplevel 2>&1`,
		"attempt_2": `cd /tmp && (echo "--version:"; ${SPACEDOCK_BIN:-spacedock} --version) 2>&1; echo "---"; git rev-parse --show-toplevel 2>&1`,
	}
	for name, command := range cases {
		t.Run(name, func(t *testing.T) {
			if err := detectWrongRootBoot(streamLine(command), fixtureRoot); err != nil {
				t.Errorf("detector red the real PR #446 %s boot command (probe-only, no same-command corroboration): %v", name, err)
			}
		})
	}
}

// TestDetectWrongRootBootRealPR446Streams replays the two full captured sonnet
// streams from PR #446 (run 28466995641, artifact
// runtime-live-e2e-claude-live-sonnet, `claude-shared-scenarios/feedback-3-cycle-escalation/claude-stream.jsonl`)
// through detectWrongRootBoot end to end. This is AC-1's RED->GREEN gate: on
// the pre-fix detector this reproduced both exact CI failures (confirmed by
// the ideation-stage spike); with the corroboration gate, both must return
// nil, matching the streams' actual `result` event (`subtype=success` — both
// runs finished the scenario correctly, only the detector fataled them).
func TestDetectWrongRootBootRealPR446Streams(t *testing.T) {
	cases := []struct {
		name        string
		fixture     string
		fixtureRoot string
	}{
		{
			name:        "attempt_1",
			fixture:     "claude_live_wrong_root_pr446_attempt1.stream.jsonl",
			fixtureRoot: "/tmp/TestLiveClaudeSharedScenariosfeedback-3-cycle-escalation1752814120/001",
		},
		{
			name:        "attempt_2",
			fixture:     "claude_live_wrong_root_pr446_attempt2.stream.jsonl",
			fixtureRoot: "/tmp/TestLiveClaudeSharedScenariosfeedback-3-cycle-escalation1376149796/001",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			stream, err := os.ReadFile(filepath.Join("testdata", c.fixture))
			if err != nil {
				t.Fatalf("open real-stream fixture %q: %v", c.fixture, err)
			}
			if err := detectWrongRootBoot(string(stream), c.fixtureRoot); err != nil {
				t.Errorf("detector red the real PR #446 %s stream — want nil (the run finished the scenario correctly): %v", c.name, err)
			}
		})
	}
}
