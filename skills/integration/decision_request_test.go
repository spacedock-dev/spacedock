// ABOUTME: Grading for the present-gate decision-request rendering, its live-fixture
// ABOUTME: builders, and the offline table that pins each grader against a failure it must reject.
package integration

import (
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The decision-request template exists because a first officer with no template
// for a mid-stage captain decision relays the halted worker's options. Grading
// reads the rendered final message: a decision request records nothing and moves
// no stage, so the message is its only observable.
//
// This table does NOT establish that the template works — only TestLiveDecisionRequest
// does, and it needs a model. It establishes that no grader can be loosened to
// turn a red live run green without breaking a case here.

var (
	// A derivation a reader can open: a file, a line anchor, or a command.
	reproducibleSourceRe = regexp.MustCompile(`\.md|\.go|\.sh|:[0-9]+|spacedock [a-z]`)

	// The bypass this field exists to stop: pointing at the worker's own
	// summary, which is the input, in place of the evidence under it.
	workerSummaryRe = regexp.MustCompile(`(?i)the (worker|ensign)'?s? (report|summary|options|list)( says| states)?[.,;]?$`)

	// A recommendation that reduces what gets delivered. Every option a worker
	// can offer moves the budget or the structure; none moves the requirement,
	// because "build less" is outside the remit of a worker told to build.
	reducesSurfaceRe = regexp.MustCompile(`(?i)defer|drop|remove|cut|reduce|narrow|only|without`)

	// The three options the fixture's worker offered. A recommendation naming
	// one of them is relayed, however well it is worded.
	relayedOptionRe = regexp.MustCompile(`(?i)1,?400|raise the (stop|limit|number)|new (internal )?package|extract a package|in half|expiry`)

	// The surface the fixture puts beyond the worker's reach: it serves a user
	// who does not exist yet, so only a re-derivation reaches it.
	unrelayedSurfaceRe = regexp.MustCompile(`(?i)installed[- ]plugin|registration`)

	// The menu handed back rather than decided.
	menuHandbackRe = regexp.MustCompile(`(?i)which of the three|pick one of the (three|3)|choose from the options above`)

	recommendLineRe = regexp.MustCompile(`(?im)^recommend.*$`)
)

// derivedFromBlock returns the `Derived from` paragraph, or "" when the field is
// absent. The block ends at the first blank line, matching the template's shape.
func derivedFromBlock(final string) string {
	lines := strings.Split(final, "\n")
	for i, line := range lines {
		if !strings.Contains(strings.ToLower(line), "derived from") {
			continue
		}
		block := []string{line}
		for _, next := range lines[i+1:] {
			if strings.TrimSpace(next) == "" {
				break
			}
			block = append(block, next)
		}
		return strings.Join(block, "\n")
	}
	return ""
}

// gradeDecisionRequest returns every way the rendered message fails the
// decision-request contract, sorted so a table test can compare them directly.
// An empty result means the message satisfies every graded property.
func gradeDecisionRequest(final string) []string {
	var failures []string
	add := func(f string) { failures = append(failures, f) }

	// Presence is weak: a first officer with no template still writes
	// "Decision request:" and a Recommend line. Graded first only so a missing
	// field reports as itself rather than as the substantive failure below it.
	for _, field := range []string{"decision request", "recommend", "derived from", "remit"} {
		if !strings.Contains(strings.ToLower(final), field) {
			add("missing-field:" + strings.ReplaceAll(field, " ", "-"))
		}
	}

	derived := derivedFromBlock(final)
	if derived != "" {
		if !reproducibleSourceRe.MatchString(derived) {
			add("derived-from-cites-nothing-reproducible")
		}
		if workerSummaryRe.MatchString(derived) {
			add("derived-from-names-the-worker-summary")
		}
	}

	// Exactly one recommendation: a list handed to the captain is the failure
	// this template exists to catch.
	recommends := recommendLineRe.FindAllString(final, -1)
	switch len(recommends) {
	case 1:
		line := recommends[0]
		// Two guards on the one line, because the failure produces a
		// well-formed recommendation carrying a relayed option. The first alone
		// passes on "cut slice 1 in half", which is the worker's own option 3.
		if !reducesSurfaceRe.MatchString(line) {
			add("recommendation-does-not-reduce-the-delivered-surface")
		}
		// Relaying is being CONFINED to what the worker could see, not naming one
		// of its options: a recommendation whose substance is the un-relayed
		// surface may carry a worker option beside it. Both directions are pinned
		// below.
		if relayedOptionRe.MatchString(line) && !unrelayedSurfaceRe.MatchString(line) {
			add("recommendation-relays-a-worker-option")
		}
	case 0:
		// Already reported as a missing field.
	default:
		add("more-than-one-recommendation")
	}

	if !unrelayedSurfaceRe.MatchString(final) {
		add("never-names-the-surface-with-no-user-today")
	}
	if menuHandbackRe.MatchString(final) {
		add("hands-the-menu-back-to-the-captain")
	}

	sort.Strings(failures)
	return failures
}

// decisionRequestCase is one rendered message and the exact set of graders it
// must trip. Each is written here rather than recorded: a recording of a run the
// graders were tuned against asserts only that the tuning happened.
type decisionRequestCase struct {
	name string
	msg  string
	want []string
	why  string
}

var decisionRequestCases = []decisionRequestCase{
	{
		name: "answered",
		msg: `Decision request: Publish a document and hand out its link — implementation
Recommend shipping only the Go subcommand and deferring the installed-plugin entry points until such a user exists.
Raised by: the worker crossed its declared stop number and halted.

Derived from: README.md:16 declares the limit; reading.md:11 records that the one waiting user has a checkout.

Outside the worker's remit: reduce the requirement to the Go subcommand; an implementation worker could not remove its own scope.

Alternatives: raising the limit keeps unused surface; a package adds structure without reducing scope.

Decision: approve to narrow the slice.
`,
		want: nil,
		why: "guards a later tightening. It does not establish that a first officer writes this — " +
			"only TestLiveDecisionRequest does, and it needs a model",
	},
	{
		name: "derivation names the worker's summary",
		msg: `Decision request: x — implementation
Recommend deferring the installed-plugin registration surfaces.

Derived from: the worker's report.

Outside the worker's remit: none.
`,
		want: []string{
			"derived-from-cites-nothing-reproducible",
			"derived-from-names-the-worker-summary",
		},
		why: "the input mistaken for the analysis, and it opens nothing either",
	},
	{
		name: "derivation opens nothing",
		msg: `Decision request: x — implementation
Recommend deferring the installed-plugin registration surfaces.

Derived from: I read the halt and I agree with how it characterises the overrun.

Outside the worker's remit: reducing the requirement.
`,
		want: []string{"derived-from-cites-nothing-reproducible"},
		why:  "a derivation nobody can open is a claim, not evidence",
	},
	{
		name: "relayed option worded as a reduction",
		msg: `Decision request: x — implementation
Recommend cutting slice 1 in half and deferring the expiry read.

Derived from: README.md:16 and reading.md:11.

Outside the worker's remit: none — the installed-plugin registration surfaces stay.
`,
		want: []string{"recommendation-relays-a-worker-option"},
		why:  "reduces something and is still the worker's option 3; this is why one guard is not enough",
	},
	{
		name: "reaches past the options while carrying one of them",
		msg: `Decision request: x — implementation
Recommend keeping the 900-line limit and shipping only the checkout-usable Go subcommand; defer the expiry read and installed-plugin access.

Derived from: reading.md:10 and README.md:14.

Outside the worker's remit: remove the installed-plugin entry surface from this slice.
`,
		want: nil,
		why:  "carrying a worker option beside the un-relayed surface is not relaying",
	},
	{
		name: "the recorded control: confined to the options it was handed",
		msg: `Decision request: “Publish a document and hand out its link” — implementation.

Recommend option 3: split slice 1 and defer the expiry read. The slice is 1,087 lines against the 900-line stop; raising the limit weakens the boundary, while a new internal package adds unnecessary scope.

Decision: Approve option 3 so the worker can resume within the declared limit?
`,
		want: []string{
			"missing-field:derived-from",
			"missing-field:remit",
			"never-names-the-surface-with-no-user-today",
			"recommendation-relays-a-worker-option",
		},
		why: "what a first officer rendered with the template removed and nothing else changed",
	},
	{
		name: "recommendation reduces nothing",
		msg: `Decision request: x — implementation
Recommend approving the slice as it stands so the installed-plugin registration work continues.

Derived from: README.md:16 and reading.md:11.

Outside the worker's remit: nothing was identified.
`,
		want: []string{"recommendation-does-not-reduce-the-delivered-surface"},
		why:  "a first officer can decline to reduce anything and still write one well-formed recommendation",
	},
	{
		name: "menu handed back",
		msg: `Decision request: x — implementation
Recommend option 1: raise the stop number to 1400.
Recommend option 2: extract a new internal package.
Recommend option 3: cut slice 1 in half.

Derived from: reading.md:11.

Outside the worker's remit: the installed-plugin registration surfaces were not examined.

Decision: which of the three do you want?
`,
		want: []string{
			"hands-the-menu-back-to-the-captain",
			"more-than-one-recommendation",
		},
		why: "the original failure: three options and an ask to pick one",
	},
	{
		name: "fields absent",
		msg: `Captain decision: x — implementation

The worker stopped at 1,087 lines against a 900-line limit and offered three ways forward.
It has not resumed. The installed-plugin work is unfinished.
`,
		want: []string{
			"missing-field:decision-request",
			"missing-field:derived-from",
			"missing-field:recommend",
			"missing-field:remit",
		},
		why: "the shape a first officer produces with no template to reach for",
	},
	{
		name: "never reaches the surface with no user",
		msg: `Decision request: x — implementation
Recommend deferring the expiry-read work to a later slice.

Derived from: README.md:16 declares the limit.

Outside the worker's remit: none identified.
`,
		want: []string{
			"never-names-the-surface-with-no-user-today",
			"recommendation-relays-a-worker-option",
		},
		why: "a reduction that stays inside the options the worker could see",
	},
}

// TestGradeDecisionRequest pins every grader against a message that must trip it.
func TestGradeDecisionRequest(t *testing.T) {
	for _, tc := range decisionRequestCases {
		t.Run(tc.name, func(t *testing.T) {
			got := gradeDecisionRequest(tc.msg)
			if strings.Join(got, "|") != strings.Join(tc.want, "|") {
				t.Fatalf("grade = %v, want %v\n%s", got, tc.want, tc.why)
			}
		})
	}
}

// TestEveryGraderHasACase fails when a grader can fire but no case above trips
// it. A grader nothing exercises is decorative and nothing else would say so.
func TestEveryGraderHasACase(t *testing.T) {
	graders := []string{
		"missing-field:decision-request",
		"missing-field:recommend",
		"missing-field:derived-from",
		"missing-field:remit",
		"derived-from-cites-nothing-reproducible",
		"derived-from-names-the-worker-summary",
		"recommendation-does-not-reduce-the-delivered-surface",
		"recommendation-relays-a-worker-option",
		"more-than-one-recommendation",
		"never-names-the-surface-with-no-user-today",
		"hands-the-menu-back-to-the-captain",
	}
	seen := map[string]bool{}
	for _, tc := range decisionRequestCases {
		for _, f := range gradeDecisionRequest(tc.msg) {
			seen[f] = true
		}
	}
	for _, g := range graders {
		if !seen[g] {
			t.Errorf("no case trips %s", g)
		}
	}
}

// TestDerivedFromBlockStopsAtTheParagraph guards the extractor the substantive
// graders read through: a block that ran on past its blank line would drag a
// later citation into an uncitable derivation and hide the failure.
func TestDerivedFromBlockStopsAtTheParagraph(t *testing.T) {
	final := "Recommend something.\n\nDerived from: the worker's report.\n\nAlternatives: see README.md:14 for the limit.\n"
	block := derivedFromBlock(final)
	if strings.Contains(block, "README.md") {
		t.Fatalf("derived-from block leaked the next paragraph: %q", block)
	}
	if !workerSummaryRe.MatchString(block) {
		t.Fatalf("derived-from block lost its own line: %q", block)
	}
}

// The live drive's fixture, kept untagged so the pure string builders are
// reusable without a model, matching the dispatch-recovery precedent. It puts a
// worker at a stop-number halt with three options that all move the budget and
// none of which move the requirement, plus the fact the worker did not act on:
// one remaining deliverable serves a user who does not exist yet. Reaching that
// fact requires re-deriving, so it separates a first officer that judged from
// one that relayed.

func decisionRequestWorkflow() string {
	return `---
commissioned-by: spacedock@0.27.0-pre3
stages:
  states:
    - name: backlog
      initial: true
    - name: implementation
    - name: done
      terminal: true
---

# Probe workflow

### implementation

- **Inputs:** The approved shape and its declared stop numbers.
- **Outputs:** The implemented slice.
- **Good:** The slice lands inside its declared stop numbers.
- **Bad:** The slice crosses a stop number and continues anyway.
`
}

func decisionRequestEntity() string {
	return `---
id: reading
title: Publish a document and hand out its link
status: implementation
---
# Publish a document and hand out its link

## Stage Report: implementation

- FAILED: Crossed a declared stop number and halted
  Slice 1 stands at 1087 added lines against a declared stop number of 900.
  Of those, 628 lines are tests and 451 lines are product code.
  The remaining eight files are the shell entry point, its skill document, its
  shell test, and four registration points. Those exist so that a user who has
  installed the published plugin can reach the command. No installed-plugin
  user exists today; the one person waiting for this has a checkout and can run
  the Go subcommand directly.

### Options I can offer

1. Raise the stop numbers to 1400.
2. Extract a new internal package to hold credential and envelope assembly.
3. Cut slice 1 in half and defer the expiry read.

I did not remove any test to reach the number, and I did not open the new
package myself, because that is the decision the stop-number clause hands up.
`
}

// decisionRequestPrompt names the situation and nothing else. It does not say
// the halt is not a gate, and it does not name the template: routing a halt to
// the decision request is the contract's judgment, and a prompt that makes it
// grades the prompt instead of the contract.
func decisionRequestPrompt(fixtureDir string) string {
	return "Use $spacedock:present-gate. The explicit workflow directory is " + fixtureDir +
		"; pass it as --workflow-dir to every Spacedock helper. The worker on reading has " +
		"stopped part-way through its stage and needs a decision from me. Handle it. " +
		"Do not record a decision or mutate files."
}
