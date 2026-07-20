// ABOUTME: Structural-absence guard — runtime adapters carry no negative host-contrast
// ABOUTME: wording and no mutable absolute step ordinals, with its non-vacuity control.
package contractlint

import (
	"fmt"
	"strings"
	"testing"
)

// stepOrdinalPhrases are the absolute step-number couplings an adapter must not
// carry: a step ordinal goes stale the moment the referenced procedure is
// renumbered, so adapters name capabilities instead.
var stepOrdinalPhrases = []string{
	"Dispatch step",
	"Event Loop step",
	"Merge-and-Cleanup step",
	"step 10",
}

// negativeContrastPhrases are the "this host is not that host" smell phrases the
// runtime-support positive-binding sweep removed, per adapter. An adapter may name
// its own host; it must not define itself by what another host lacks. A blanket
// host-name ban is deliberately not used — the Pi FO adapter legitimately names
// hosts in transport prose — so each entry is a phrase the sweep actually removed.
// Absence guards only; the runtime-meaning claims are bound by the three binding
// tests, and the unowned ones are listed under UNCOVERED RUNTIME TOKENS.
var negativeContrastPhrases = map[string][]string{
	codexEnsignRel: {
		"Claude",
		"another host",
		"Codex declares none",
		"Do not reconstruct",
		"Do not send follow-up",
	},
	piFORuntimeRel: {
		"does not expose Claude Code team-tool signatures",
		"Do not call or ask workers to call Claude team tools",
		"Pi has no such enum",
		"Claude-centric enum",
		"no Claude-centric",
		"Do not create teams",
		"Codex declares none",
	},
	piEnsignRel: {
		"Do not assume Claude team tools",
		"Claude team tools exist in Pi",
	},
}

// stepOrdinalFiles are the adapters held to the no-absolute-step-ordinal rule.
var stepOrdinalFiles = []string{piFORuntimeRel}

// lifecycleHeadings are the per-adapter lifecycle sections the runtime-binding-block
// migration folded into `## Runtime implementation`; a re-introduced one gives a cold
// agent two competing lifecycle stories. Deleted once on the unverified reasoning that
// the capability-set equality subsumed it — wrong, because that set is extracted from
// `## Runtime implementation` alone and never sees a sibling section.
var lifecycleHeadings = map[string][]string{
	codexFORuntimeRel: {
		"## Dispatch\n",
		"## Awaiting Completion\n",
		"## Reuse And Feedback Routing\n",
	},
	piFORuntimeRel: {
		"## Runtime Shape\n",
		"## Dispatch\n",
		"## Awaiting Completion\n",
		"## Follow-up and Reuse\n",
		"## Shutdown\n",
		"### Model Resolution\n",
		"### Canonical Model Space\n",
	},
}

func proseAbsenceViolations(text string, banned []string, kind string) []string {
	var out []string
	for _, phrase := range banned {
		if strings.Contains(text, phrase) {
			out = append(out, fmt.Sprintf("contains %s %q", kind, phrase))
		}
	}
	return out
}

// TestRuntimeAdaptersAvoidNegativeContrastAndStepOrdinals holds the Codex and Pi
// adapters to three structural absences: no negative host-contrast wording, no
// absolute step ordinal coupling an adapter to a numbered procedure, and no
// re-introduced lifecycle section competing with the runtime-binding block.
func TestRuntimeAdaptersAvoidNegativeContrastAndStepOrdinals(t *testing.T) {
	if len(negativeContrastPhrases) == 0 || len(stepOrdinalPhrases) == 0 || len(lifecycleHeadings) == 0 {
		t.Fatal("banned-phrase tables are empty — the guards would pass vacuously")
	}
	for rel, banned := range negativeContrastPhrases {
		for _, msg := range proseAbsenceViolations(readRepoFile(t, rel), banned, "negative host-contrast wording") {
			t.Errorf("%s %s", rel, msg)
		}
	}
	for _, rel := range stepOrdinalFiles {
		for _, msg := range proseAbsenceViolations(readRepoFile(t, rel), stepOrdinalPhrases, "mutable absolute step ordinal") {
			t.Errorf("%s %s", rel, msg)
		}
	}
	for rel, banned := range lifecycleHeadings {
		for _, msg := range proseAbsenceViolations(readRepoFile(t, rel), banned, "retired lifecycle heading") {
			t.Errorf("%s %s", rel, msg)
		}
	}
}

// TestRuntimeAdapterAbsenceGuardsDiscriminate is the non-vacuity control:
// positive-binding adapter prose PASSES, while a re-introduced host-contrast
// sentence, step ordinal, and lifecycle section each RED.
func TestRuntimeAdapterAbsenceGuardsDiscriminate(t *testing.T) {
	pass := []struct {
		why, text, kind string
		banned          []string
	}{
		{"positive Codex ensign binding", "`«context-budget»` is unavailable unless a future probe binds it.\n", "negative host-contrast wording", negativeContrastPhrases[codexEnsignRel]},
		{"capability-named teardown", "Teardown runs through `«worker.shutdown»`.\n", "mutable absolute step ordinal", stepOrdinalPhrases},
	}
	for _, c := range pass {
		if v := proseAbsenceViolations(c.text, c.banned, c.kind); len(v) != 0 {
			t.Fatalf("control: the %s was wrongly flagged: %v", c.why, v)
		}
	}

	red := []struct {
		why, text, kind string
		banned          []string
	}{
		{"re-introduced host-contrast sentence", "The context budget is unavailable here; Codex declares none.\n", "negative host-contrast wording", negativeContrastPhrases[codexEnsignRel]},
		{"re-introduced enum contrast", "Pi has no such enum, so pass the model string through.\n", "negative host-contrast wording", negativeContrastPhrases[piFORuntimeRel]},
		{"re-introduced step ordinal", "Tear the worker down at Merge-and-Cleanup step 10.\n", "mutable absolute step ordinal", stepOrdinalPhrases},
		{"re-introduced lifecycle section", "## Awaiting Completion\n\nPoll the worker until it reports.\n", "retired lifecycle heading", lifecycleHeadings[codexFORuntimeRel]},
	}
	for _, c := range red {
		if v := proseAbsenceViolations(c.text, c.banned, c.kind); len(v) == 0 {
			t.Fatalf("control: the %s was not flagged — the guard stopped biting", c.why)
		}
	}
}
