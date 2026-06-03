// ABOUTME: AC-4 prose-inflator lock-in oracles for the FO + ensign contract files.
// ABOUTME: Catches audit-trail-style exposition and cross-file restatement.
package hostneutrality

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// contractProseFiles is the four-file surface the AC-4 oracles police.
var contractProseFiles = []string{
	filepath.Join("..", "..", "skills", "first-officer", "references", "first-officer-shared-core.md"),
	filepath.Join("..", "..", "skills", "ensign", "references", "ensign-shared-core.md"),
	filepath.Join("..", "..", "skills", "ensign", "references", "claude-ensign-runtime.md"),
	filepath.Join("..", "..", "skills", "ensign", "references", "codex-ensign-runtime.md"),
}

// sharedCorePaths are the universal cores. They must not restate one another
// against the per-runtime adapter cores.
var sharedCorePaths = []string{
	filepath.Join("..", "..", "skills", "first-officer", "references", "first-officer-shared-core.md"),
	filepath.Join("..", "..", "skills", "ensign", "references", "ensign-shared-core.md"),
}

// runtimeAdapterPaths are the per-host ensign adapter cores tested against the
// shared cores for restatement.
var runtimeAdapterPaths = []string{
	filepath.Join("..", "..", "skills", "ensign", "references", "claude-ensign-runtime.md"),
	filepath.Join("..", "..", "skills", "ensign", "references", "codex-ensign-runtime.md"),
}

// auditTrailLiterals are exact phrases banned across all four contract files.
// These are the audit-trail-exposition class — history-as-meta-comment that
// inflates without conveying contract content.
var auditTrailLiterals = []string{
	"audit-trail",
	"the audit said",
	"the auditor flagged",
	"the audit returned",
	"now we do X because",
}

// auditTrailRegexes catch the parameterized variants of the audit-trail class.
// The cycle-N pattern is scoped to audit/sweep context — a bare "cycle 3" in
// the feedback-flow contract clause (a load-bearing workflow threshold) is
// not audit-trail exposition; "cycle-1 audit" or "the cycle 2 sweep" is.
var auditTrailRegexes = []*regexp.Regexp{
	// cycle-N variants in audit context: "cycle-N audit", "the cycle N sweep",
	// "cycleN reconcile", etc. The "audit"/"sweep"/"reconcile" qualifier is what
	// makes it audit-trail prose; bare "cycle 3" stays allowed for contract use.
	regexp.MustCompile(`(?i)cycle[-\s]?\d+\b[^.\n]{0,40}(audit|sweep|reconcile)\b`),
	regexp.MustCompile(`(?i)\b(audit|sweep|reconcile)\b[^.\n]{0,40}cycle[-\s]?\d+\b`),
	// w-prefix WfRunID-shaped task IDs (8-10 chars), the leak vector for
	// specific audit task IDs in prose. Requires a digit somewhere after
	// the leading w to exclude English words like "with", "want".
	regexp.MustCompile(`\bw[a-z0-9]*\d[a-z0-9]{6,9}\b`),
}

// TestNoAuditTrailExposition locks AC-4: the four contract prose files MUST
// NOT carry audit-trail-style exposition — phrases like "audit-trail",
// "cycle-N", "the auditor flagged", WfRunID-shaped task IDs ("wobgtdlsp"),
// or "now we do X because" restatements. These inflate the load every
// ensign dispatch re-pays without conveying contract content.
//
// Scope: full files per captain (Path B authorization). A passing test means
// the swept HEAD is clean; a deliberately-inserted regression of any banned
// phrase fails the test (positive proof of lock-in).
func TestNoAuditTrailExposition(t *testing.T) {
	for _, path := range contractProseFiles {
		t.Run(filepath.Base(path), func(t *testing.T) {
			body, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			text := string(body)
			lowerText := strings.ToLower(text)

			for _, banned := range auditTrailLiterals {
				if strings.Contains(lowerText, strings.ToLower(banned)) {
					t.Errorf("%s contains banned audit-trail literal %q — re-inflation of swept prose", path, banned)
				}
			}
			for _, re := range auditTrailRegexes {
				if loc := re.FindStringIndex(text); loc != nil {
					match := text[loc[0]:loc[1]]
					t.Errorf("%s contains banned audit-trail pattern %q (matched %q) — re-inflation of swept prose",
						path, re.String(), match)
				}
			}
		})
	}
}

// TestNoCrossFileRestatement locks AC-4: a 12-word n-gram from a runtime
// adapter core (claude-ensign-runtime.md, codex-ensign-runtime.md) MUST NOT
// reappear verbatim in either shared core. Cross-file restatement inflates
// the dispatch-load (the FO/ensign reads BOTH shared core and a runtime
// adapter; restated content is double-paid).
//
// Exceptions:
//   - heading-marked spans (e.g. `### Backstop (Claude)`, `## Sequencing rule`)
//   - fenced code blocks (``` ... ```)
//   - markdown tables (lines starting with `|`)
//   - host-qualified contrast spans (per spanHostQualified — naming both Codex
//     and Claude in the same span)
//
// Threshold: 12 contiguous words.
func TestNoCrossFileRestatement(t *testing.T) {
	// Build the n-gram set from the runtime adapter cores, excluding
	// exception spans.
	adapterNGrams := map[string]string{}
	for _, path := range runtimeAdapterPaths {
		spans := parseProseSpansForOverlap(t, path)
		for _, sp := range spans {
			words := tokenizeForOverlap(sp.text)
			for i := 0; i+12 <= len(words); i++ {
				gram := strings.Join(words[i:i+12], " ")
				if _, exists := adapterNGrams[gram]; !exists {
					adapterNGrams[gram] = filepath.Base(path)
				}
			}
		}
	}

	// Walk the shared cores looking for any of those n-grams appearing in a
	// non-exception span.
	for _, path := range sharedCorePaths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			spans := parseProseSpansForOverlap(t, path)
			for _, sp := range spans {
				words := tokenizeForOverlap(sp.text)
				for i := 0; i+12 <= len(words); i++ {
					gram := strings.Join(words[i:i+12], " ")
					if source, hit := adapterNGrams[gram]; hit {
						t.Errorf("%s span starting line %d restates a 12-word n-gram from %s: %q",
							path, sp.startLine, source, gram)
					}
				}
			}
		})
	}
}

// proseSpan is one span with location and exception-class info.
type proseSpan struct {
	text           string
	startLine      int
	excludeOverlap bool
}

// parseProseSpansForOverlap walks a markdown file and returns spans suitable for
// the n-gram restatement oracle. Spans inside fenced code blocks, markdown
// tables, host-qualified contrast spans, or under sections marked with
// known-headed exceptions are flagged with excludeOverlap=true, and the
// function filters them out — only the includable spans are returned.
func parseProseSpansForOverlap(t *testing.T, path string) []proseSpan {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	lines := strings.Split(string(data), "\n")

	var spans []proseSpan
	var cur []string
	curStart := 0
	inFence := false
	headingMarksException := false

	flush := func() {
		if len(cur) == 0 {
			return
		}
		text := strings.Join(cur, "\n")
		// Skip table-only spans (all non-blank lines start with `|`).
		if isMarkdownTable(text) {
			cur = nil
			return
		}
		// Skip host-qualified contrast spans.
		if spanHostQualified(text) {
			cur = nil
			return
		}
		// Skip spans under headings flagged as exception spans.
		if headingMarksException {
			cur = nil
			return
		}
		spans = append(spans, proseSpan{text: text, startLine: curStart})
		cur = nil
	}

	for i, line := range lines {
		lineNo := i + 1
		trimmed := strings.TrimSpace(line)

		// Headings flip section-exception state but are not spans themselves.
		if strings.HasPrefix(trimmed, "#") {
			flush()
			headingMarksException = isExceptionHeading(trimmed)
			continue
		}

		if strings.HasPrefix(trimmed, "```") {
			// Fenced code blocks are exception spans — flush prose first, skip
			// fence body until close, do not include it.
			flush()
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}

		if trimmed == "" {
			flush()
			continue
		}

		if len(cur) == 0 {
			curStart = lineNo
		}
		cur = append(cur, line)
	}
	flush()
	return spans
}

// isMarkdownTable reports whether a span's non-blank lines are all table rows
// (start with `|`) — table rows often share `| Header | Header |` cells across
// files without being prose restatement.
func isMarkdownTable(text string) bool {
	lines := strings.Split(text, "\n")
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if t == "" {
			continue
		}
		if !strings.HasPrefix(t, "|") {
			return false
		}
	}
	return true
}

// isExceptionHeading reports whether a heading marks a section whose contents
// are exempt from the cross-file restatement oracle. These are the named
// host-contrast sections where each runtime adapter is expected to mirror
// shared-core terminology by design.
func isExceptionHeading(heading string) bool {
	exceptions := []string{
		"Backstop (Claude)",
		"Backstop (Codex)",
		"Standing teammates",
		"Sequencing rule",
	}
	for _, ex := range exceptions {
		if strings.Contains(heading, ex) {
			return true
		}
	}
	return false
}

// tokenizeForOverlap splits a span into lowercase word tokens, stripping
// markdown formatting characters (`*`, `_`, backticks) and trailing
// punctuation. Code identifiers stay intact (`status --set`, `pull --rebase`).
func tokenizeForOverlap(text string) []string {
	// Replace markdown emphasis and inline-code backticks with spaces, keep
	// the content. Hyphens and dots inside identifiers stay as one token.
	cleaned := text
	for _, ch := range []string{"*", "_", "`", "(", ")", "[", "]"} {
		cleaned = strings.ReplaceAll(cleaned, ch, " ")
	}
	// Strip trailing punctuation tokens like `.`, `,`, `;`, `:` by inserting
	// spaces around them.
	for _, ch := range []string{",", ";", ":", "—"} {
		cleaned = strings.ReplaceAll(cleaned, ch, " ")
	}
	fields := strings.Fields(cleaned)
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.TrimRight(f, ".!?")
		if f == "" {
			continue
		}
		out = append(out, strings.ToLower(f))
	}
	return out
}
