//go:build live

// ABOUTME: Real Claude First Officer proof of the selected room-only gate handoff.
package ensigncycle

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spacedock-dev/spacedock/internal/gates"
)

const selectedGateSkillSinkRefusal = "SELECTED-GATE-SKILL-SINK-REFUSED"

func writeSelectedGateSkillCallSink(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".claude-plugin", "plugin.json"), `{
  "name": "subspace",
  "version": "0.0.0-test",
  "description": "Test-only native Skill-call sink.",
  "skills": "./skills/"
}
`)
	writeFile(t, filepath.Join(root, "skills", "r", "SKILL.md"), `---
name: r
description: Test-only native Skill-call sink; return one deterministic refusal without reading arguments or running tools.
---

# Selected gate Skill-call sink

The invocation itself is the observed handoff. Do not read the room, validate
capability, materialize sources, emit evidence, simulate a Result or inventory,
run a tool, retry, or choose a fallback. Return exactly
`+selectedGateSkillSinkRefusal+` to the invoking root session.
`)
	return root
}

func selectedGateOverridePrompt(fixture recordedGateFixture) string {
	return "Use $spacedock:first-officer for this whole run.\n\n" +
		"Workflow directory: " + fixture.root + "\n" +
		"Engage only `recorded-gate-task`. This session declares one presentation override: select the installed `subspace:r` override for this gate.\n" +
		"Use the already-committed `.spacedock-state/recorded-gate-task/selected/gate-review.md` Artifact and both already-committed References, `.spacedock-state/recorded-gate-task/selected/entity-snapshot.md` and `recorder-contract.md`; do not replace or regenerate them.\n" +
		"Bring the open validation gate through provider-neutral preparation and its durable bind to the selected presentation channel. Stop after the selected override returns; leave the gate open."
}

func TestLiveClaudeSelectedGateOverride(t *testing.T) {
	runner := newClaudeLiveRunner(t)
	fixture := writePreparedRecordedGateFixture(t)
	before := readFile(t, fixture.entity)
	commandLog := filepath.Join(fixture.root, "evidence", "selected-override-command.log")
	shimDir := writeRecordedGateLoggingShim(t, buildRecordedGateBinary(t), commandLog)
	runner = runner.withStubPATH(shimDir).(claudeLiveRunner)
	shellEnvDir := t.TempDir()
	bashEnv := filepath.Join(shellEnvDir, "selected-override-env.sh")
	writeFile(t, bashEnv, "export SPACEDOCK_BIN="+filepath.Join(shimDir, "spacedock")+"\n")
	writeFile(t, filepath.Join(shellEnvDir, ".zshenv"), readFile(t, bashEnv))
	runner.env = withRecordedGateEnv(runner.env, "BASH_ENV", bashEnv)
	runner.env = withRecordedGateEnv(runner.env, "ZDOTDIR", shellEnvDir)
	runner = runner.withExtraPluginDir(writeSelectedGateSkillCallSink(t))

	scenario := sharedRuntimeScenario{
		name:   "selected-gate-override",
		intent: "FO prepares and commits one gate room, then emits one selected native room-only Skill call without probe or fallback.",
	}
	result := runner.run(t, scenario, fixture.root, selectedGateOverridePrompt(fixture))
	room := filepath.Join(filepath.Dir(fixture.entity), "review", "validation", "briefing-1")
	logBytes, readErr := os.ReadFile(commandLog)
	if readErr != nil {
		t.Fatalf("selected gate command log unavailable: %v\nFinal message:\n%s\nArtifacts: %s", readErr, result.finalMessage, result.artifactDir)
	}
	if err := assertSelectedGateOverride(string(logBytes), result.stream, room); err != nil {
		t.Fatalf("selected gate override graded FAIL: %v\nFinal message:\n%s\nArtifacts: %s", err, result.finalMessage, result.artifactDir)
	}
	after := readFile(t, fixture.entity)
	if before == after || !validatingStatus.MatchString(after) {
		t.Fatalf("selected gate override did not durably bind the open validation gate\nArtifacts: %s", result.artifactDir)
	}
	doc, _, readGateErr := gates.Read(fixture.entity)
	if readGateErr != nil {
		t.Fatalf("selected gate durable binding is invalid: %v\nArtifacts: %s", readGateErr, result.artifactDir)
	}
	current := gates.CurrentSummary(doc)
	if current.State != "open" || current.Briefing == "" || current.Resolution != "" || current.Application != "" {
		t.Fatalf("selected gate did not remain open and unresolved: %#v\nArtifacts: %s", current, result.artifactDir)
	}
	entries, roomErr := os.ReadDir(room)
	if roomErr != nil {
		t.Fatalf("selected gate room unavailable: %v\nArtifacts: %s", roomErr, result.artifactDir)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	if want := []string{"gate-briefing.json", "request.json"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("selected gate room files = %v, want %v\nArtifacts: %s", names, want, result.artifactDir)
	}
	if !strings.Contains(result.finalMessage, selectedGateSkillSinkRefusal) {
		t.Fatalf("Skill-call sink refusal did not return to the root session: %q\nArtifacts: %s", result.finalMessage, result.artifactDir)
	}
}
