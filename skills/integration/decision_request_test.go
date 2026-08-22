// ABOUTME: Grading for the present-gate decision-request rendering, plus the offline
// ABOUTME: table test over recorded first-officer final messages in testdata.
package integration

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The decision-request template exists because a first officer with no template
// for a mid-stage captain decision relays the halted worker's options. Grading
// therefore reads the rendered final message: a decision request records
// nothing and moves no stage, so the message is its only observable. These
// graders are the whole checkable surface, which is why they are Go with
// offline fixtures rather than regexes reachable only through a live run.

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

	// Presence is form, and the recorded pair shows how little of it: without the
	// template the first officer still wrote "Decision request:" and a Recommend
	// line, so those two checks separated nothing. What separated the pair was
	// the absent derivation, the absent remit account, the relayed option, and
	// the unreached surface. Graded first only so a missing field reports as
	// itself rather than as the substantive failure downstream of it.
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
	// this template exists to catch, and so is silence.
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
		if relayedOptionRe.MatchString(line) {
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

func readFixture(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join("testdata", "decision-request", name)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	return string(body)
}

// TestGradeDecisionRequest pins each grader against a recorded or constructed
// message. `with-template.txt` and `without-template.txt` are real first-officer
// output captured from the same fixture, prompt, binary, and base revision,
// differing only in whether present-gate carried the template — so the pair is
// the evidence that the template, and not the prompt, produces the behavior.
// The remaining fixtures are the bypasses a relaying first officer produces;
// each exists because a grader that cannot reject it is decorative.
func TestGradeDecisionRequest(t *testing.T) {
	cases := []struct {
		fixture string
		want    []string
		why     string
	}{
		{
			fixture: "with-template.txt",
			want:    nil,
			why:     "live output with the template present satisfies every graded property",
		},
		{
			fixture: "without-template.txt",
			want: []string{
				"missing-field:derived-from",
				"missing-field:remit",
				"never-names-the-surface-with-no-user-today",
				"recommendation-relays-a-worker-option",
			},
			why: "same input without the template: the recommendation is the worker's own option 3",
		},
		{
			fixture: "derived-from-worker-summary.txt",
			want: []string{
				"derived-from-cites-nothing-reproducible",
				"derived-from-names-the-worker-summary",
			},
			why: "citing the worker's report is the input mistaken for the analysis, and it opens nothing either",
		},
		{
			fixture: "derived-from-uncitable.txt",
			want:    []string{"derived-from-cites-nothing-reproducible"},
			why:     "a derivation nobody can open is a claim, not evidence",
		},
		{
			fixture: "relayed-option-worded-as-reduction.txt",
			want:    []string{"recommendation-relays-a-worker-option"},
			why:     "'cut slice 1 in half' reduces something and is still relayed; this is why one guard is not enough",
		},
		{
			fixture: "menu-handback.txt",
			want: []string{
				"hands-the-menu-back-to-the-captain",
				"more-than-one-recommendation",
			},
			why: "three recommendations and an explicit ask to pick one",
		},
	}

	for _, tc := range cases {
		t.Run(tc.fixture, func(t *testing.T) {
			got := gradeDecisionRequest(readFixture(t, tc.fixture))
			if strings.Join(got, "|") != strings.Join(tc.want, "|") {
				t.Fatalf("grade(%s) = %v, want %v\n%s", tc.fixture, got, tc.want, tc.why)
			}
		})
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
