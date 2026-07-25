// ABOUTME: Mutation controls for the observed selected-gate-override oracle.
package ensigncycle

import (
	"fmt"
	"strings"
	"testing"
)

func selectedOverrideAssistant(block string) string {
	return `{"type":"assistant","parent_tool_use_id":null,"message":{"content":[` + block + `]}}`
}

func selectedOverrideTool(name, input string) string {
	return selectedOverrideAssistant(`{"type":"tool_use","name":"` + name + `","input":{` + input + `}}`)
}

func TestSelectedGateOverrideObservedBehavior(t *testing.T) {
	const room = "/tmp/workflow/.spacedock-state/recorded-gate-task/review/validation/briefing-1"
	shimLog := func(exit int, command string) string {
		return fmt.Sprintf("begin\t%s\nexit=%d\t%s\n", command, exit, command)
	}
	commandLog := shimLog(0, "gate --help") +
		shimLog(0, "gate prepare recorded-gate-task --question Advance?") +
		shimLog(0, "state commit recorded-gate-task --workflow-dir /tmp/workflow")
	prepare := selectedOverrideTool("Bash", `"command":"spacedock gate prepare recorded-gate-task --question Advance?"`)
	commit := selectedOverrideTool("Bash", `"command":"spacedock state commit recorded-gate-task --workflow-dir /tmp/workflow"`)
	override := selectedOverrideTool("Skill", `"skill":"subspace:r","args":"gate `+room+`"`)
	refusal := selectedOverrideAssistant(`{"type":"text","text":"SELECTED-GATE-PROVIDER-REFUSED"}`)
	valid := strings.Join([]string{prepare, commit, override, refusal}, "\n")
	if err := assertSelectedGateOverride(commandLog, valid, room); err != nil {
		t.Fatalf("baseline: %v", err)
	}

	cases := map[string]struct {
		log    string
		stream string
	}{
		"missing-help":   {strings.Replace(commandLog, shimLog(0, "gate --help"), "", 1), valid},
		"duplicate-help": {commandLog + shimLog(0, "gate --help"), valid},
		"failed-help-only": {
			strings.Replace(commandLog, shimLog(0, "gate --help"), shimLog(1, "gate --help"), 1),
			valid,
		},
		"failed-help-before-success": {
			shimLog(1, "gate --help") + commandLog,
			valid,
		},
		"failed-prepare":     {strings.Replace(commandLog, "exit=0\tgate prepare", "exit=1\tgate prepare", 1), valid},
		"pre-bind-override":  {commandLog, strings.Join([]string{prepare, override, commit, refusal}, "\n")},
		"duplicate-override": {commandLog, strings.Join([]string{prepare, commit, override, override, refusal}, "\n")},
		"reconstructed-authority": {
			commandLog,
			strings.Replace(valid, `"args":"gate `+room+`"`, `"args":"gate `+room+` --entity recorded-gate-task"`, 1),
		},
		"provider-probe": {
			commandLog,
			strings.Join([]string{prepare, commit, selectedOverrideTool("Bash", `"command":"subspace --version"`), override, refusal}, "\n"),
		},
		"agent-probe": {
			commandLog,
			strings.Join([]string{prepare, commit, selectedOverrideTool("Agent", `"description":"probe provider"`), override, refusal}, "\n"),
		},
		"agent-probe-before-prepare": {
			commandLog,
			strings.Join([]string{selectedOverrideTool("Agent", `"description":"probe provider"`), prepare, commit, override, refusal}, "\n"),
		},
		"chat-fallback": {
			commandLog,
			valid + "\n" + selectedOverrideAssistant(fmt.Sprintf(`{"type":"text","text":%q}`, recordedGateReview())),
		},
		"chat-review-before-override": {
			commandLog,
			strings.Join([]string{prepare, commit, selectedOverrideAssistant(fmt.Sprintf(`{"type":"text","text":%q}`, recordedGateReview())), override, refusal}, "\n"),
		},
		"record-after-handoff": {
			commandLog + "exit=0\tgate record recorded-gate-task --result /tmp/result.json\n",
			valid + "\n" + selectedOverrideTool("Bash", `"command":"spacedock gate record recorded-gate-task --result /tmp/result.json"`),
		},
		"consume-after-handoff": {
			commandLog + "exit=0\tgate consume recorded-gate-task\n",
			valid + "\n" + selectedOverrideTool("Bash", `"command":"spacedock gate consume recorded-gate-task"`),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if err := assertSelectedGateOverride(tc.log, tc.stream, room); err == nil {
				t.Fatal("mutant graded PASS")
			}
		})
	}
}
